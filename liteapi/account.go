package liteapi

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func checkAccountProof(accountID ton.AccountID, blockID ton.BlockIDExt, shardBlock ton.BlockIDExt, proof []*boc.Cell, shardProof []byte, stateProof *boc.Cell) (*tlb.ShardAccount, *tlb.ShardStateUnsplit, error) {
	var blockHash ton.Bits256
	if accountID.Workchain == -1 && blockID.Workchain == -1 || blockID == shardBlock {
		blockHash = blockID.RootHash
	} else {
		if len(shardProof) == 0 {
			return nil, nil, errors.New("empty shard proof")
		}
		shardHash := shardBlock.RootHash
		if _, err := checkShardInMasterProof(blockID, shardProof, accountID.Workchain, shardHash); err != nil {
			return nil, nil, fmt.Errorf("shard proof is incorrect: %w", err)
		}
		blockHash = shardHash
	}
	cellsMap := make(map[[32]byte]*boc.Cell)
	hash, err := stateProof.Hash256WithLevel(0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get state hash: %w", err)
	}
	cellsMap[hash] = stateProof

	shardState, err := checkBlockShardStateProof(proof, blockHash, cellsMap)
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

	return nil, nil, errors.New("invalid account state")
}

func checkShardInMasterProof(master ton.BlockIDExt, shardProof []byte, workchain int32, shardRootHash ton.Bits256) (*tlb.McStateExtra, error) {
	shardProofCells, err := boc.DeserializeBoc(shardProof)
	if err != nil {
		return nil, err
	}
	shardState, err := checkBlockShardStateProof(shardProofCells, master.RootHash, nil)
	if err != nil {
		return nil, fmt.Errorf("check block proof failed: %w", err)
	}
	if !shardState.ShardStateUnsplit.Custom.Exists {
		return nil, fmt.Errorf("not a masterchain block")
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
					return &stateExtra, nil
				}
			case "New":
				if int32(k) == workchain && ton.Bits256(b.New.RootHash) == shardRootHash {
					return &stateExtra, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("required shard hash not found in proof")
}

func checkBlockShardStateProof(proof []*boc.Cell, blockRootHash ton.Bits256, cellsMap map[[32]byte]*boc.Cell) (*tlb.ShardStateUnsplit, error) {
	if len(proof) != 2 {
		return nil, errors.New("must be two root cells")
	}
	block, err := checkProof[tlb.Block](*proof[0], blockRootHash, nil)
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
	err = decoder.Unmarshal(proof[1], &stateProof)
	if err != nil {
		return nil, err
	}
	if stateProof.Proof.VirtualHash != block.VirtualRoot.StateUpdate.ToHash {
		return nil, errors.New("invalid virtual hash")
	}
	return &stateProof.Proof.VirtualRoot, nil
}

func checkProof[T any](proof boc.Cell, hash ton.Bits256, decoder *tlb.Decoder) (*tlb.MerkleProof[T], error) {
	if decoder == nil {
		decoder = tlb.NewDecoder()
	}
	var res tlb.MerkleProof[T]
	err := decoder.Unmarshal(&proof, &res) // merkle hash and depth checks inside
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
	}
	if ton.Bits256(res.VirtualHash) != hash {
		return nil, fmt.Errorf("invalid hash")
	}
	return &res, nil // return new_hash field of MerkleUpdate of ShardState
}

func checkTxProof(
	shardAccount tlb.HashmapAugE[tlb.Bits256, tlb.AccountBlock, tlb.CurrencyCollection],
	accountAddr tlb.Bits256,
	lt uint64,
	hash ton.Bits256,
) error {
	var accountBlock tlb.AccountBlock
	for i, key := range shardAccount.Keys() {
		if key.Equal(accountAddr) {
			accountBlock = shardAccount.Values()[i]
			break
		}
	}

	var accountTx *tlb.Transaction
	for i, key := range accountBlock.Transactions.Keys() {
		if uint64(key) == lt {
			accountTx = &accountBlock.Transactions.Values()[i].Value
		}
	}

	if accountTx == nil {
		return fmt.Errorf("tx not found")
	}

	txHash := accountTx.Hash()
	if !bytes.Equal(txHash[:], hash[:]) {
		return fmt.Errorf("invalid tx hash")
	}

	return nil
}

type prunedBlock struct {
	Magic       tlb.Magic `tlb:"block#11ef55aa"`
	GlobalId    int32
	Info        tlb.BlockInfo                    `tlb:"^"`
	ValueFlow   tlb.ValueFlow                    `tlb:"^"`
	StateUpdate tlb.MerkleUpdate[tlb.ShardState] `tlb:"^"`
	Extra       prunedBlockExtra                 `tlb:"^"`
}

type prunedBlockExtra struct {
	Magic           tlb.Magic                                                                `tlb:"block_extra#4a33f6fd"`
	InMsgDescrCell  boc.Cell                                                                 `tlb:"^"`
	OutMsgDescrCell boc.Cell                                                                 `tlb:"^"`
	AccountBlocks   tlb.HashmapAugE[tlb.Bits256, prunedAccountBlock, tlb.CurrencyCollection] `tlb:"^"`
	RandSeed        tlb.Bits256
	CreatedBy       tlb.Bits256
	Custom          tlb.Maybe[tlb.Ref[tlb.McBlockExtra]]
}

type prunedAccountBlock struct {
	Magic        tlb.Magic `tlb:"acc_trans#5"`
	AccountAddr  tlb.Bits256
	Transactions tlb.HashmapAug[tlb.Uint64, boc.Cell, tlb.CurrencyCollection]
	StateUpdate  tlb.HashUpdate `tlb:"^"`
}

// this function doing same proof as checkTxProof but for pruned transactions that cannot be resolved
func checkPrunedTxProof(
	shardAccount tlb.HashmapAugE[tlb.Bits256, prunedAccountBlock, tlb.CurrencyCollection],
	accountAddr tlb.Bits256,
	lt uint64,
	hash ton.Bits256,
) error {
	var accountBlock prunedAccountBlock
	for i, key := range shardAccount.Keys() {
		if key.Equal(accountAddr) {
			accountBlock = shardAccount.Values()[i]
			break
		}
	}

	var accountTx *boc.Cell
	for i, key := range accountBlock.Transactions.Keys() {
		if uint64(key) == lt {
			accountTx = &accountBlock.Transactions.Values()[i]
		}
	}

	if accountTx == nil {
		return fmt.Errorf("tx not found")
	}
	if len(accountTx.Refs()) == 0 {
		return fmt.Errorf("shard account must have at least one ref in value")
	}

	txHash, err := accountTx.Refs()[0].Hash256WithLevel(0)
	if err != nil {
		return err
	}
	if !bytes.Equal(txHash[:], hash[:]) {
		return fmt.Errorf("invalid tx hash")
	}

	return nil
}

func getShardHashesHash(proofBlock boc.Cell, merkleProof *tlb.MerkleProof[tlb.Block]) (ton.Bits256, error) {
	if !merkleProof.VirtualRoot.Extra.Custom.Exists {
		return ton.Bits256{}, fmt.Errorf("mc block extra is missing in block")
	}

	mcExtraCell := proofBlock.Refs()[0].Refs()[3].Refs()[3]
	mcExtraCell.ResetCounters()
	err := mcExtraCell.Skip(17) // 16 + 1
	if err != nil {
		return ton.Bits256{}, err
	}
	hasShardHashesMap, err := mcExtraCell.ReadBit()
	if err != nil {
		return ton.Bits256{}, err
	}

	mapCell := boc.NewCell()
	err = mapCell.WriteBit(hasShardHashesMap)
	if err != nil {
		return ton.Bits256{}, err
	}

	if hasShardHashesMap {
		err = mapCell.AddRef(mcExtraCell.Refs()[0])
		if err != nil {
			return ton.Bits256{}, err
		}
	}

	return mapCell.Hash256WithLevel(0)
}
