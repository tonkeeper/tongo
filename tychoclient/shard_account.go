package tychoclient

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tychoclient/proto"
)

// GetShardAccount gets account state from a specific shard at the latest block
func (c *Client) GetShardAccount(ctx context.Context, workchain int32, address []byte, withProof bool) (*ShardAccountInfo, error) {
	return c.getShardAccountInternal(ctx, workchain, address, withProof, &proto.GetShardAccountRequest_Latest{
		Latest: &proto.LatestBlock{},
	})
}

// GetShardAccountAtSeqno gets account state at a specific block seqno
func (c *Client) GetShardAccountAtSeqno(ctx context.Context, workchain int32, address []byte, withProof bool, blockWorkchain int32, shard uint64, seqno uint32) (*ShardAccountInfo, error) {
	return c.getShardAccountInternal(ctx, workchain, address, withProof, &proto.GetShardAccountRequest_BySeqno{
		BySeqno: &proto.BlockBySeqno{
			Workchain: blockWorkchain,
			Shard:     shard,
			Seqno:     seqno,
		},
	})
}

// GetShardAccountByBlockId gets account state at a specific block
func (c *Client) GetShardAccountByBlockId(ctx context.Context, workchain int32, address []byte, withProof bool, blockId *proto.BlockId) (*ShardAccountInfo, error) {
	return c.getShardAccountInternal(ctx, workchain, address, withProof, &proto.GetShardAccountRequest_ById{
		ById: &proto.BlockById{
			Id: blockId,
		},
	})
}

// getShardAccountInternal is the internal implementation for getting shard account
func (c *Client) getShardAccountInternal(ctx context.Context, workchain int32, address []byte, withProof bool, atBlock interface{}) (*ShardAccountInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	req := &proto.GetShardAccountRequest{
		Workchain: workchain,
		Address:   address,
		WithProof: withProof,
	}

	// Set the atBlock field based on the interface type
	switch v := atBlock.(type) {
	case *proto.GetShardAccountRequest_Latest:
		req.AtBlock = v
	case *proto.GetShardAccountRequest_BySeqno:
		req.AtBlock = v
	case *proto.GetShardAccountRequest_ById:
		req.AtBlock = v
	default:
		return nil, fmt.Errorf("unsupported atBlock type")
	}

	resp, err := c.client.GetShardAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard account: %w", err)
	}

	switch result := resp.Account.(type) {
	case *proto.GetShardAccountResponse_BlockNotFound:
		mcInfo := result.BlockNotFound.McStateInfo
		return nil, fmt.Errorf("block not found at MC seqno %d (LT: %d, utime: %d)",
			mcInfo.McSeqno, mcInfo.Lt, mcInfo.Utime)
	case *proto.GetShardAccountResponse_Accessed:
		info := &ShardAccountInfo{
			McStateInfo:  result.Accessed.McStateInfo,
			AccountState: result.Accessed.AccountState,
			Proof:        result.Accessed.Proof,
		}

		// Parse account state if present
		if len(result.Accessed.AccountState) > 0 {
			parsed, err := ParseShardAccount(result.Accessed.AccountState)
			if err != nil {
				return nil, fmt.Errorf("failed to parse account state: %w", err)
			}
			info.ParsedAccountState = parsed
		}

		return info, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", result)
	}
}

// ShardAccountInfo contains account state information
type ShardAccountInfo struct {
	McStateInfo        *proto.McStateInfo
	AccountState       []byte    // BOC-encoded ShardAccount (if found)
	Proof              []byte    // BOC-encoded collection of proofs (if withProof was true)
	ParsedAccountState *boc.Cell // Parsed BOC cell
}

// ParseShardAccount parses BOC-encoded ShardAccount data
func ParseShardAccount(bocData []byte) (*boc.Cell, error) {
	if len(bocData) == 0 {
		return nil, fmt.Errorf("empty BOC data")
	}

	cells, err := boc.DeserializeBoc(bocData)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize BOC: %w", err)
	}

	if len(cells) == 0 {
		return nil, fmt.Errorf("no cells in BOC")
	}

	// Return the first root cell which contains the ShardAccount
	return cells[0], nil
}

// String returns a human-readable representation of ShardAccountInfo
func (s *ShardAccountInfo) String() string {
	if s == nil {
		return "<nil>"
	}

	result := "ShardAccountInfo:\n"
	if s.McStateInfo != nil {
		result += fmt.Sprintf("  MC Seqno: %d\n", s.McStateInfo.McSeqno)
		result += fmt.Sprintf("  MC LT: %d\n", s.McStateInfo.Lt)
		result += fmt.Sprintf("  MC Utime: %d\n", s.McStateInfo.Utime)
	}

	if len(s.AccountState) > 0 {
		result += fmt.Sprintf("  Account State: %d bytes\n", len(s.AccountState))
		result += fmt.Sprintf("  Account State (hex): %s...\n", hex.EncodeToString(s.AccountState[:min(32, len(s.AccountState))]))
	} else {
		result += "  Account State: not found\n"
	}

	if len(s.Proof) > 0 {
		result += fmt.Sprintf("  Proof: %d bytes\n", len(s.Proof))
	}

	if s.ParsedAccountState != nil {
		result += fmt.Sprintf("  Parsed Account:\n%s\n", s.ParsedAccountState.ToString())
	}

	return result
}
