package ton

import (
	"encoding/binary"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type BlockID struct {
	Workchain int32
	Shard     uint64
	Seqno     uint32
}

type BlockIDExt struct {
	BlockID
	RootHash Bits256
	FileHash Bits256
}

func (id BlockIDExt) MarshalTL() ([]byte, error) {
	payload := make([]byte, 80)
	binary.LittleEndian.PutUint32(payload[:4], uint32(id.Workchain))
	binary.LittleEndian.PutUint64(payload[4:12], id.Shard)
	binary.LittleEndian.PutUint32(payload[12:16], id.Seqno)
	copy(payload[16:48], id.RootHash[:])
	copy(payload[48:80], id.FileHash[:])
	return payload, nil
}

func (id *BlockIDExt) UnmarshalTL(data []byte) error {
	if len(data) != 80 {
		return fmt.Errorf("invalid data length")
	}
	id.Workchain = int32(binary.LittleEndian.Uint32(data[:4]))
	id.Shard = binary.LittleEndian.Uint64(data[4:12])
	id.Seqno = binary.LittleEndian.Uint32(data[12:16])
	copy(id.RootHash[:], data[16:48])
	copy(id.FileHash[:], data[48:80])
	return nil
}

func NewTonBlockId(fileHash, rootHash Bits256, seqno uint32, shard uint64, workchain int32) *BlockIDExt {
	return &BlockIDExt{
		BlockID: BlockID{
			Workchain: workchain,
			Shard:     shard,
			Seqno:     seqno,
		},
		FileHash: fileHash,
		RootHash: rootHash,
	}
}

func (id BlockIDExt) String() string {
	return fmt.Sprintf("(%d,%x,%d,%x,%x)", id.Workchain, id.Shard, id.Seqno, id.RootHash, id.FileHash)
}
func (id BlockID) String() string {
	return fmt.Sprintf("(%d,%x,%d)", id.Workchain, id.Shard, id.Seqno)
}

func getParents(blkPrevInfo tlb.BlkPrevInfo, afterSplit, afterMerge bool, shard uint64, workchain int32) ([]BlockIDExt, error) {
	var parents []BlockIDExt
	if !afterMerge {
		if blkPrevInfo.SumType != "PrevBlkInfo" {
			return nil, fmt.Errorf("two parent blocks may be only after merge")
		}
		blockID := BlockIDExt{
			BlockID: BlockID{
				Workchain: workchain,
				Seqno:     blkPrevInfo.PrevBlkInfo.Prev.SeqNo,
			},
			FileHash: Bits256(blkPrevInfo.PrevBlkInfo.Prev.FileHash),
			RootHash: Bits256(blkPrevInfo.PrevBlkInfo.Prev.RootHash),
		}
		if afterSplit {
			blockID.Shard = shardParent(shard)
			return []BlockIDExt{blockID}, nil
		}
		blockID.Shard = shard
		return []BlockIDExt{blockID}, nil
	}

	if blkPrevInfo.SumType != "PrevBlksInfo" {
		return nil, fmt.Errorf("two parent blocks must be after merge")
	}

	parents = append(parents, BlockIDExt{
		BlockID: BlockID{
			Seqno:     blkPrevInfo.PrevBlksInfo.Prev1.SeqNo,
			Shard:     shardChild(shard, true),
			Workchain: workchain,
		},
		FileHash: Bits256(blkPrevInfo.PrevBlksInfo.Prev1.FileHash),
		RootHash: Bits256(blkPrevInfo.PrevBlksInfo.Prev1.RootHash),
	})

	parents = append(parents, BlockIDExt{
		FileHash: Bits256(blkPrevInfo.PrevBlksInfo.Prev2.FileHash),
		RootHash: Bits256(blkPrevInfo.PrevBlksInfo.Prev2.RootHash),
		BlockID: BlockID{
			Seqno:     blkPrevInfo.PrevBlksInfo.Prev2.SeqNo,
			Shard:     shardChild(shard, false),
			Workchain: workchain,
		},
	})

	return parents, nil
}

// td::uint64 x = td::lower_bit64(shard) >> 1;
// return left ? shard - x : shard + x;
func shardChild(shard uint64, left bool) uint64 {
	x := (shard & (^shard + 1)) >> 1
	if left {
		return shard - x
	}
	return shard + x
}

// td::uint64 x = td::lower_bit64(shard);
// return (shard - x) | (x << 1);
func shardParent(shard uint64) uint64 {
	x := shard & (^shard + 1)
	return (shard - x) | (x << 1)
}

func convertShardIdent(si tlb.ShardIdent) (workchain int32, shard uint64) {
	shard = si.ShardPrefix
	pow2 := uint64(1) << (63 - si.ShardPfxBits)
	shard |= pow2
	return si.WorkchainID, shard
}

func GetParents(i tlb.BlockInfo) ([]BlockIDExt, error) {
	workchain, shard := convertShardIdent(i.Shard)
	return getParents(i.PrevRef, i.AfterSplit, i.AfterMerge, shard, workchain)
}

func ToBlockId(s tlb.ShardDesc, workchain int32) BlockIDExt {
	if s.SumType == "Old" {
		return BlockIDExt{
			BlockID: BlockID{
				Workchain: workchain,
				Shard:     uint64(s.Old.NextValidatorShard),
				Seqno:     s.Old.SeqNo,
			},
			RootHash: Bits256(s.Old.RootHash),
			FileHash: Bits256(s.Old.FileHash),
		}
	} else {
		return BlockIDExt{
			BlockID: BlockID{
				Workchain: workchain,
				Shard:     uint64(s.New.NextValidatorShard),
				Seqno:     s.New.SeqNo,
			},
			RootHash: Bits256(s.New.RootHash),
			FileHash: Bits256(s.New.FileHash),
		}
	}
}

func CreateExternalMessage(address AccountID, body *boc.Cell, init *tlb.StateInit, importFee tlb.VarUInteger16) (tlb.Message, error) {
	// TODO: add either selection algorithm
	var msg = tlb.Message{
		Info: tlb.CommonMsgInfo{
			SumType: "ExtInMsgInfo",
			ExtInMsgInfo: &struct {
				Src       tlb.MsgAddress
				Dest      tlb.MsgAddress
				ImportFee tlb.VarUInteger16
			}{
				Src:       (*AccountID)(nil).ToMsgAddress(),
				Dest:      address.ToMsgAddress(),
				ImportFee: importFee,
			},
		},
		Body: tlb.EitherRef[tlb.Any]{
			IsRight: true,
			Value:   tlb.Any(*body),
		},
	}
	if init != nil {
		msg.Init.Exists = true
		msg.Init.Value.IsRight = true
		msg.Init.Value.Value = *init
	}
	return msg, nil
}

// ShardIDs returns a list of IDs of shard blocks this block refers to.
func ShardIDs(blk *tlb.Block) []BlockIDExt {
	if !blk.Extra.Custom.Exists {
		return nil
	}
	items := blk.Extra.Custom.Value.Value.ShardHashes.Items()
	shardsCount := 0
	for _, item := range items {
		for _, x := range item.Value.Value.BinTree.Values {
			if x.SeqNo() == 0 {
				continue
			}
			shardsCount += 1
		}
	}
	if shardsCount == 0 {
		return nil
	}
	shards := make([]BlockIDExt, 0, shardsCount)
	for _, item := range items {
		for _, shardDesc := range item.Value.Value.BinTree.Values {
			if shardDesc.SeqNo() == 0 {
				continue
			}
			shards = append(shards, ToBlockId(shardDesc, int32(item.Key)))
		}
	}
	return shards
}

// ParseBlockID tries to construct BlockID from the given string.
func ParseBlockID(s string) (BlockID, error) {
	var id BlockID
	_, err := fmt.Sscanf(s, "(%d,%x,%d)", &id.Workchain, &id.Shard, &id.Seqno)
	if err != nil {
		return BlockID{}, err
	}
	return id, nil
}

func MustParseBlockID(s string) BlockID {
	id, err := ParseBlockID(s)
	if err != nil {
		panic(err)
	}
	return id
}
