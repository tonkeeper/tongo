package tongo

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type BlockID = ton.BlockID

type BlockIDExt = ton.BlockIDExt

func NewTonBlockId(fileHash, rootHash Bits256, seqno uint32, shard uint64, workchain int32) *BlockIDExt {
	return ton.NewTonBlockId(fileHash, rootHash, seqno, shard, workchain)
}

func GetParents(i tlb.BlockInfo) ([]BlockIDExt, error) {
	return ton.GetParents(i)
}

func ToBlockId(s tlb.ShardDesc, workchain int32) BlockIDExt {
	return ton.ToBlockId(s, workchain)
}

func CreateExternalMessage(address AccountID, body *boc.Cell, init *tlb.StateInit, importFee tlb.VarUInteger16) (tlb.Message, error) {
	return ton.CreateExternalMessage(address, body, init, importFee)
}

// ShardIDs returns a list of IDs of shard blocks this block refers to.
func ShardIDs(blk *tlb.Block) []BlockIDExt {
	return ton.ShardIDs(blk)
}

// ParseBlockID tries to construct BlockID from the given string.
func ParseBlockID(s string) (BlockID, error) {
	return ton.ParseBlockID(s)
}

func MustParseBlockID(s string) BlockID {
	return ton.MustParseBlockID(s)
}
