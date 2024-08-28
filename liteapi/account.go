package liteapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// GetAccountWithProof
// For safe operation, always use GetAccountWithProof with WithBlock(proofedBlock ton.BlockIDExt), as the proof of masterchain cashed blocks is not implemented yet!
func (c *Client) GetAccountWithProof(ctx context.Context, accountID ton.AccountID) (*tlb.ShardAccount, *tlb.ShardStateUnsplit, error) {
	res, err := c.GetAccountStateRaw(ctx, accountID) // TODO: add proof check for masterHead
	if err != nil {
		return nil, nil, err
	}
	blockID := res.Id.ToBlockIdExt()
	if len(res.Proof) == 0 {
		return nil, nil, errors.New("empty proof")
	}

	var blockHash ton.Bits256
	if (accountID.Workchain == -1 && blockID.Workchain == -1) || blockID == res.Shardblk.ToBlockIdExt() {
		blockHash = blockID.RootHash
	} else {
		if len(res.ShardProof) == 0 {
			return nil, nil, errors.New("empty shard proof")
		}
		if res.Shardblk.RootHash == [32]byte{} { // TODO: how to check for empty shard?
			return nil, nil, errors.New("shard block not passed")
		}
		shardHash := ton.Bits256(res.Shardblk.RootHash)
		if err := checkShardInMasterProof(blockID, res.ShardProof, accountID.Workchain, shardHash); err != nil {
			return nil, nil, fmt.Errorf("shard proof is incorrect: %w", err)
		}
		blockHash = shardHash
	}
	cellsMap := make(map[[32]byte]*boc.Cell)
	if len(res.State) > 0 {
		stateCells, err := boc.DeserializeBoc(res.State)
		if err != nil {
			return nil, nil, fmt.Errorf("state deserialization failed: %w", err)
		}
		hash, err := stateCells[0].Hash256()
		if err != nil {
			return nil, nil, fmt.Errorf("get hash err: %w", err)
		}
		cellsMap[hash] = stateCells[0]
	}
	shardState, err := checkBlockShardStateProof(res.Proof, blockHash, cellsMap)
	if err != nil {
		return nil, nil, fmt.Errorf("incorrect block proof: %w", err)
	}
	values := shardState.ShardStateUnsplit.Accounts.Values()
	keys := shardState.ShardStateUnsplit.Accounts.Keys()
	for i, k := range keys {
		if k == accountID.Address {
			return &values[i], shardState, nil
		}
	}
	if len(res.State) == 0 {
		return &tlb.ShardAccount{Account: tlb.Account{SumType: "AccountNone"}}, shardState, nil
	}
	return nil, nil, errors.New("invalid account state")
}

func checkShardInMasterProof(master ton.BlockIDExt, shardProof []byte, workchain int32, shardRootHash ton.Bits256) error {
	shardState, err := checkBlockShardStateProof(shardProof, master.RootHash, nil)
	if err != nil {
		return fmt.Errorf("check block proof failed: %w", err)
	}
	if !shardState.ShardStateUnsplit.Custom.Exists {
		return fmt.Errorf("not a masterchain block")
	}
	stateExtra := shardState.ShardStateUnsplit.Custom.Value.Value
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

func checkBlockShardStateProof(proof []byte, blockRootHash ton.Bits256, cellsMap map[[32]byte]*boc.Cell) (*tlb.ShardStateUnsplit, error) {
	proofCells, err := boc.DeserializeBoc(proof)
	if err != nil {
		return nil, err
	}
	if len(proofCells) != 2 {
		return nil, errors.New("must be two root cells")
	}
	newStateHash, err := checkBlockProof(proofCells[0], blockRootHash)
	if err != nil {
		return nil, fmt.Errorf("incorrect block proof: %w", err)
	}
	var stateProof struct {
		Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
	}
	decoder := tlb.NewDecoder()
	if cellsMap != nil {
		decoder = decoder.WithPrunedResolver(func(hash tlb.Bits256) (*boc.Cell, error) {
			cell, ok := cellsMap[hash]
			if ok {
				return cell, nil
			}
			return nil, errors.New("not found")
		})
	}
	err = decoder.Unmarshal(proofCells[1], &stateProof)
	if err != nil {
		return nil, err
	}
	if stateProof.Proof.VirtualHash != *newStateHash {
		return nil, errors.New("invalid virtual hash")
	}
	return &stateProof.Proof.VirtualRoot, nil
}

func checkBlockProof(proof *boc.Cell, blockRootHash ton.Bits256) (*tlb.Bits256, error) {
	var res tlb.MerkleProof[tlb.Block]
	err := tlb.Unmarshal(proof, &res) // merkle hash and depth checks inside
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block proof: %w", err)
	}
	if ton.Bits256(res.VirtualHash) != blockRootHash {
		return nil, fmt.Errorf("invalid block root hash")
	}
	return &res.VirtualRoot.StateUpdate.ToHash, nil // return new_hash field of MerkleUpdate of ShardState
}
