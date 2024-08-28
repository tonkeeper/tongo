package liteapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// blockTrimmed stripped-down version of the tlb.Block with pruned cell instead of ShardState and skip some decoding
type blockTrimmed struct {
	Magic       tlb.Magic `tlb:"block#11ef55aa"`
	GlobalId    int32
	Info        boc.Cell                   `tlb:"^"`
	ValueFlow   boc.Cell                   `tlb:"^"`
	StateUpdate tlb.MerkleUpdate[boc.Cell] `tlb:"^"`
	Extra       boc.Cell                   `tlb:"^"`
}

// shardAccountPruned stripped-down version of the tlb.ShardAccount with pruned cell instead of Account
type shardAccountPruned struct {
	Account       boc.Cell `tlb:"^"`
	LastTransHash tlb.Bits256
	LastTransLt   uint64
}

// shardStateUnsplitTrimmed stripped-down version of the ShardStateUnsplit structure for extracting Account proof
type shardStateUnsplitTrimmed struct {
	Magic           tlb.Magic `tlb:"shard_state#9023afe2"`
	GlobalID        int32
	ShardID         tlb.ShardIdent
	SeqNo           uint32
	VertSeqNo       uint32
	GenUtime        uint32
	GenLt           uint64
	MinRefMcSeqno   uint32
	OutMsgQueueInfo boc.Cell `tlb:"^"`
	BeforeSplit     bool
	Accounts        tlb.HashmapAugE[tlb.Bits256, shardAccountPruned, tlb.DepthBalanceInfo] `tlb:"^"`
	Other           tlb.ShardStateUnsplitOther                                             `tlb:"^"`
	Custom          tlb.Maybe[tlb.Ref[boc.Cell]]
}

// TODO: use proof merger instead of trimmed
type shardStateUnsplit struct {
	ShardStateTrimmed shardStateUnsplitTrimmed
	ShardState        tlb.ShardStateUnsplit
}

// GetAccountWithProof
// For safe operation, always use GetAccountWithProof with WithBlock(proofedBlock ton.BlockIDExt), as the proof of masterchain cashed blocks is not implemented yet!
func (c *Client) GetAccountWithProof(ctx context.Context, accountID ton.AccountID) (*tlb.ShardAccount, *tlb.ShardStateUnsplit, error) { // TODO: return merged tlb.ShardStateUnsplit (proof+account)
	res, err := c.GetAccountStateRaw(ctx, accountID) // TODO: add proof check for masterHead
	if err != nil {
		return nil, nil, err
	}
	blockID := res.Id.ToBlockIdExt()
	if len(res.Proof) == 0 {
		return nil, nil, errors.New("empty proof")
	}
	var shardHash *ton.Bits256
	if accountID.Workchain != -1 { // TODO: set masterchain constant
		if len(res.ShardProof) == 0 {
			return nil, nil, errors.New("empty shard proof")
		} // TODO: change logic for shard blockID (proof shard to proofed master)
		if res.Shardblk.RootHash == [32]byte{} { // TODO: how to check for empty shard?
			return nil, nil, errors.New("shard block not passed")
		}
		h := ton.Bits256(res.Shardblk.RootHash)
		shardHash = &h
	}
	blockHash := blockID.RootHash
	if shardHash != nil { // we need shard proof only for not masterchain
		if err := checkShardInMasterProof(blockID, res.ShardProof, accountID.Workchain, *shardHash); err != nil {
			return nil, nil, fmt.Errorf("shard proof is incorrect: %w", err)
		}
		blockHash = *shardHash
	}
	shardState, err := checkBlockShardStateProof(res.Proof, blockHash)
	if err != nil {
		return nil, nil, fmt.Errorf("incorrect block proof: %w", err)
	}
	values := shardState.ShardStateTrimmed.Accounts.Values()
	keys := shardState.ShardStateTrimmed.Accounts.Keys()
	for i, k := range keys {
		if bytes.Equal(k[:], accountID.Address[:]) {
			acc, err := decodeAccount(res.State, values[i])
			if err != nil {
				return nil, nil, err
			}
			return acc, &shardState.ShardState, nil
		}
	}
	if len(res.State) == 0 {
		return &tlb.ShardAccount{Account: tlb.Account{SumType: "AccountNone"}}, &shardState.ShardState, nil
	}
	return nil, nil, errors.New("invalid account state")
}

func decodeAccount(state []byte, shardAccount shardAccountPruned) (*tlb.ShardAccount, error) {
	stateCells, err := boc.DeserializeBoc(state)
	if err != nil {
		return nil, err
	}
	if len(stateCells) != 1 {
		return nil, boc.ErrNotSingleRoot
	}
	accountHash, err := stateCells[0].Hash256()
	if err != nil {
		return nil, err
	}
	shardAccountHash, err := shardAccount.Account.Hash256WithLevel(0)
	if err != nil {
		return nil, err
	}
	if accountHash != shardAccountHash {
		return nil, errors.New("invalid account hash")
	}
	var acc tlb.Account
	err = tlb.Unmarshal(stateCells[0], &acc)
	if err != nil {
		return nil, err
	}
	// do not check account balance from tlb.DepthBalanceInfo
	res := tlb.ShardAccount{Account: acc, LastTransHash: shardAccount.LastTransHash, LastTransLt: shardAccount.LastTransLt}
	return &res, nil
}

func checkShardInMasterProof(master ton.BlockIDExt, shardProof []byte, workchain int32, shardRootHash ton.Bits256) error {
	shardState, err := checkBlockShardStateProof(shardProof, master.RootHash)
	if err != nil {
		return fmt.Errorf("check block proof failed: %w", err)
	}
	if !shardState.ShardState.ShardStateUnsplit.Custom.Exists {
		return fmt.Errorf("not a masterchain block")
	}
	stateExtra := shardState.ShardState.ShardStateUnsplit.Custom.Value.Value
	keys := stateExtra.ShardHashes.Keys()
	values := stateExtra.ShardHashes.Values()
	for i, k := range keys {
		binTreeValues := values[i].Value.BinTree.Values
		for _, b := range binTreeValues {
			switch b.SumType {
			case "Old":
				if int32(k) == workchain && ton.Bits256(b.Old.RootHash) == shardRootHash {
					return nil
				}
			case "New":
				if int32(k) == workchain && ton.Bits256(b.New.RootHash) == shardRootHash {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("required shard hash not found in proof")
}

func checkBlockShardStateProof(proof []byte, blockRootHash ton.Bits256) (*shardStateUnsplit, error) {
	proofCells, err := boc.DeserializeBoc(proof)
	if err != nil {
		return nil, err
	}
	if len(proofCells) != 2 {
		return nil, errors.New("must be two root cells")
	}
	block, err := checkBlockProof(proofCells[0], blockRootHash)
	if err != nil {
		return nil, fmt.Errorf("incorrect block proof: %w", err)
	}
	var stateTrimmedProof struct {
		Proof tlb.MerkleProof[shardStateUnsplitTrimmed]
	}
	err = tlb.Unmarshal(proofCells[1], &stateTrimmedProof) // cells order must be strictly defined
	if err != nil {
		return nil, err
	}
	proofCells[1].ResetCounters()
	var stateProof struct {
		Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(proofCells[1], &stateProof)
	if err != nil {
		return nil, err
	}
	toRootHash, err := block.StateUpdate.ToRoot.Hash256WithLevel(0)
	if err != nil {
		return nil, err
	}
	if stateTrimmedProof.Proof.VirtualHash != toRootHash {
		return nil, errors.New("invalid virtual hash")
	}
	res := shardStateUnsplit{
		ShardStateTrimmed: stateTrimmedProof.Proof.VirtualRoot,
		ShardState:        stateProof.Proof.VirtualRoot,
	}
	return &res, nil
}

func checkBlockProof(proof *boc.Cell, blockRootHash ton.Bits256) (*blockTrimmed, error) {
	var res tlb.MerkleProof[blockTrimmed]
	err := tlb.Unmarshal(proof, &res) // merkle hash and depth checks inside
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block proof: %w", err)
	}
	if ton.Bits256(res.VirtualHash) != blockRootHash {
		return nil, fmt.Errorf("invalid block root hash")
	}
	block := res.VirtualRoot
	return &block, nil
}
