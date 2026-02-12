package tychoclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tychoclient/proto"
)

// GetShardAccount gets account state and proof LT from a specific shard at the latest block
func (c *Client) GetShardAccount(ctx context.Context, workchain int32, address []byte) (*ShardAccountInfo, uint64, error) {
	return c.getShardAccountInternal(ctx, workchain, address, &proto.GetShardAccountRequest_Latest{
		Latest: &proto.LatestBlock{},
	})
}

// GetShardAccountAtSeqno gets account state and proof LT at a specific block seqno
func (c *Client) GetShardAccountAtSeqno(ctx context.Context, workchain int32, address []byte, blockWorkchain int32, shard uint64, seqno uint32) (*ShardAccountInfo, uint64, error) {
	return c.getShardAccountInternal(ctx, workchain, address, &proto.GetShardAccountRequest_BySeqno{
		BySeqno: &proto.BlockBySeqno{
			Workchain: blockWorkchain,
			Shard:     shard,
			Seqno:     seqno,
		},
	})
}

// GetShardAccountByBlockId gets account state and proof LT at a specific block
func (c *Client) GetShardAccountByBlockId(ctx context.Context, workchain int32, address []byte, blockId *proto.BlockId) (*ShardAccountInfo, uint64, error) {
	return c.getShardAccountInternal(ctx, workchain, address, &proto.GetShardAccountRequest_ById{
		ById: &proto.BlockById{
			Id: blockId,
		},
	})
}

// getShardAccountInternal is the internal implementation for getting shard account
func (c *Client) getShardAccountInternal(ctx context.Context, workchain int32, address []byte, atBlock interface{}) (*ShardAccountInfo, uint64, error) {
	ctx, cancel := limitedContext(ctx)
	defer cancel()

	req := &proto.GetShardAccountRequest{
		Workchain: workchain,
		Address:   address,
		WithProof: true,
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
		return nil, 0, fmt.Errorf("unsupported atBlock type")
	}

	resp, err := c.client.GetShardAccount(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get shard account: %w", err)
	}

	switch result := resp.Account.(type) {
	case *proto.GetShardAccountResponse_BlockNotFound:
		mcInfo := result.BlockNotFound.McStateInfo
		return nil, 0, fmt.Errorf("block not found at MC seqno %d (LT: %d, utime: %d)",
			mcInfo.McSeqno, mcInfo.Lt, mcInfo.Utime)
	case *proto.GetShardAccountResponse_Accessed:
		info := ShardAccountInfo{
			McStateInfo:  result.Accessed.McStateInfo,
			AccountState: result.Accessed.AccountState,
			Proof:        result.Accessed.Proof,
		}
		parsed, proofLT, err := ParseShardAccount(result.Accessed.AccountState, result.Accessed.Proof, address)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse account state: %w", err)
		}
		info.ParsedAccountState = &parsed
		return &info, proofLT, nil
	default:
		return nil, 0, fmt.Errorf("unexpected response type: %T", result)
	}
}

// GetShardAccountRaw gets account state without parsing - used for debugging
func (c *Client) GetShardAccountRaw(ctx context.Context, workchain int32, address []byte, withProof bool) (*ShardAccountInfo, error) {
	req := &proto.GetShardAccountRequest{
		Workchain: workchain,
		Address:   address,
		WithProof: withProof,
		AtBlock: &proto.GetShardAccountRequest_Latest{
			Latest: &proto.LatestBlock{},
		},
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

		// Do not parse account state - leave ParsedAccountState as nil

		return info, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", result)
	}
}

// ShardAccountInfo contains account state information
type ShardAccountInfo struct {
	McStateInfo        *proto.McStateInfo
	AccountState       []byte            // BOC-encoded ShardAccount (if found)
	Proof              []byte            // BOC-encoded collection of proofs (if withProof was true)
	ParsedAccountState *tlb.ShardAccount // Parsed TLB ShardAccount
}

// ParseShardAccount parses BOC-encoded ShardAccount data
func ParseShardAccount(state, proof, accountAddress []byte) (tlb.ShardAccount, uint64, error) {
	lt, hash, dataProofLT, err := decodeAccountDataFromProof(proof, accountAddress)
	if err != nil {
		return tlb.ShardAccount{}, 0, fmt.Errorf("failed to decode account data from proof: %w", err)
	}
	if len(state) == 0 {
		return tlb.ShardAccount{Account: tlb.Account{SumType: "AccountNone"}}, dataProofLT, nil
	}
	cell, err := boc.DeserializeSingleRootBoc(state)
	if err != nil {
		return tlb.ShardAccount{}, 0, fmt.Errorf("failed to deserialize cell: %w", err)
	}
	var acc tlb.Account
	err = tlb.Unmarshal(cell, &acc)
	if err != nil {
		return tlb.ShardAccount{}, 0, fmt.Errorf("failed to unmarshal to tlb.Account: %w", err)
	}
	return tlb.ShardAccount{Account: acc, LastTransHash: hash, LastTransLt: lt}, dataProofLT, nil
}

func decodeAccountDataFromProof(proofBytes []byte, accountAddress []byte) (uint64, tlb.Bits256, uint64, error) {
	cells, err := boc.DeserializeBoc(proofBytes)
	if err != nil {
		return 0, tlb.Bits256{}, 0, err
	}
	if len(cells) < 2 {
		return 0, tlb.Bits256{}, 0, fmt.Errorf("must be at least one root cell")
	}
	var proof struct {
		Proof tlb.MerkleProof[TychoShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[1], &proof) // cells order must be strictly defined
	if err != nil {
		return 0, tlb.Bits256{}, 0, err
	}
	dataProofLT := proof.Proof.VirtualRoot.GenLt
	values := proof.Proof.VirtualRoot.Accounts.Values()
	keys := proof.Proof.VirtualRoot.Accounts.Keys()
	for i, k := range keys {
		if bytes.Equal(k[:], accountAddress) {
			return values[i].LastTransLt, values[i].LastTransHash, dataProofLT, nil
		}
	}
	return 0, tlb.Bits256{}, dataProofLT, nil
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
		result += fmt.Sprintf("  Parsed Account: available (LT: %d, Hash: %x)\n",
			s.ParsedAccountState.LastTransLt, s.ParsedAccountState.LastTransHash)
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
