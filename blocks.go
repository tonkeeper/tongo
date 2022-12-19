package tongo

import (
	"encoding/binary"
	"fmt"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

type TonNodeBlockId struct {
	Workchain int32
	Shard     int64
	Seqno     int32
}

type TonNodeBlockIdExt struct {
	TonNodeBlockId
	RootHash Hash
	FileHash Hash
}

func (id TonNodeBlockIdExt) MarshalTL() ([]byte, error) {
	payload := make([]byte, 80)
	binary.LittleEndian.PutUint32(payload[:4], uint32(id.Workchain))
	binary.LittleEndian.PutUint64(payload[4:12], uint64(id.Shard))
	binary.LittleEndian.PutUint32(payload[12:16], uint32(id.Seqno))
	copy(payload[16:48], id.RootHash[:])
	copy(payload[48:80], id.FileHash[:])
	return payload, nil
}

func (id *TonNodeBlockIdExt) UnmarshalTL(data []byte) error {
	if len(data) != 80 {
		return fmt.Errorf("invalid data length")
	}
	id.Workchain = int32(binary.LittleEndian.Uint32(data[:4]))
	id.Shard = int64(binary.LittleEndian.Uint64(data[4:12]))
	id.Seqno = int32(binary.LittleEndian.Uint32(data[12:16]))
	copy(id.RootHash[:], data[16:48])
	copy(id.FileHash[:], data[48:80])
	return nil
}

func NewTonBlockId(fileHash, rootHash Hash, seqno int32, shard int64, workchain int32) *TonNodeBlockIdExt {
	return &TonNodeBlockIdExt{
		TonNodeBlockId: TonNodeBlockId{
			Workchain: workchain,
			Shard:     shard,
			Seqno:     seqno,
		},
		FileHash: fileHash,
		RootHash: rootHash,
	}
}

func (id TonNodeBlockIdExt) String() string {
	return fmt.Sprintf("(%d,%x,%d,%x,%x)", id.Workchain, uint64(id.Shard), id.Seqno, id.RootHash, id.FileHash)
}
func (id TonNodeBlockId) String() string {
	return fmt.Sprintf("(%d,%x,%d)", id.Workchain, uint64(id.Shard), id.Seqno)
}

// BlockInfo
// block_info#9bc7a987 version:uint32
// not_master:(## 1)
// after_merge:(## 1) before_split:(## 1)
// after_split:(## 1)
// want_split:Bool want_merge:Bool
// key_block:Bool vert_seqno_incr:(## 1)
// flags:(## 8) { flags <= 1 }
// seq_no:# vert_seq_no:# { vert_seq_no >= vert_seqno_incr }
// { prev_seq_no:# } { ~prev_seq_no + 1 = seq_no }
// shard:ShardIdent gen_utime:uint32
// start_lt:uint64 end_lt:uint64
// gen_validator_list_hash_short:uint32
// gen_catchain_seqno:uint32
// min_ref_mc_seqno:uint32
// prev_key_block_seqno:uint32
// gen_software:flags . 0?GlobalVersion
// master_ref:not_master?^BlkMasterInfo
// prev_ref:^(BlkPrevInfo after_merge)
// prev_vert_ref:vert_seqno_incr?^(BlkPrevInfo 0)
// = BlockInfo;
type BlockInfo struct {
	blockInfoPart
	GenSoftware *GlobalVersion
	MasterRef   *BlkMasterInfo
	PrevRef     BlkPrevInfo
	PrevVertRef *BlkPrevInfo
}

type blockInfoPart struct {
	Version                   uint32
	NotMaster                 bool
	AfterMerge                bool
	BeforeSplit               bool
	AfterSplit                bool
	WantSplit                 bool
	WantMerge                 bool
	KeyBlock                  bool
	VertSeqnoIncr             bool
	Flags                     uint32 `tlb:"8bits"`
	SeqNo                     uint32
	VertSeqNo                 uint32
	Shard                     ShardIdent
	GenUtime                  uint32
	StartLt                   uint64
	EndLt                     uint64
	GenValidatorListHashShort uint32
	GenCatchainSeqno          uint32
	MinRefMcSeqno             uint32
	PrevKeyBlockSeqno         uint32
}

func (i *BlockInfo) GetParents() ([]TonNodeBlockIdExt, error) {
	workchain, shard := convertShardIdent(i.Shard)
	return getParents(i.PrevRef, i.AfterSplit, i.AfterMerge, shard, workchain)
}

func (i *BlockInfo) UnmarshalTLB(c *boc.Cell, tag string) error {
	var data struct {
		Magic     tlb.Magic `tlb:"block_info#9bc7a987"`
		BlockInfo blockInfoPart
	} // for partial decoding
	err := tlb.Unmarshal(c, &data)
	if err != nil {
		return err
	}
	var res BlockInfo
	res.blockInfoPart = data.BlockInfo

	if res.Flags&1 == 1 {
		var gs GlobalVersion
		err = tlb.Unmarshal(c, &gs)
		if err != nil {
			return err
		}
		res.GenSoftware = &gs
	}

	if data.BlockInfo.NotMaster {
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		res.MasterRef = &BlkMasterInfo{}
		err = tlb.Unmarshal(c1, res.MasterRef)
		if err != nil {
			return err
		}
	}

	c1, err := c.NextRef()
	if err != nil {
		return err
	}
	err = res.PrevRef.UnmarshalTLB(c1, data.BlockInfo.AfterMerge)
	if err != nil {
		return err
	}

	if data.BlockInfo.VertSeqnoIncr {
		c1, err = c.NextRef()
		if err != nil {
			return err
		}
		res.PrevVertRef = &BlkPrevInfo{}
		err = res.PrevVertRef.UnmarshalTLB(c1, false)
		if err != nil {
			return err
		}
	}
	*i = res
	return nil
}

// GlobalVersion
// capabilities#c4 version:uint32 capabilities:uint64 = GlobalVersion;
type GlobalVersion struct {
	Magic        tlb.Magic `tlb:"capabilities#c4"`
	Version      uint32
	Capabilities uint64
}

// ExtBlkRef
// ext_blk_ref$_ end_lt:uint64 seq_no:uint32 root_hash:bits256 file_hash:bits256 = ExtBlkRef;
type ExtBlkRef struct {
	EndLt    uint64
	SeqNo    uint32
	RootHash Hash
	FileHash Hash
}

// BlkMasterInfo
// master_info$_ master:ExtBlkRef = BlkMasterInfo;
// ext_blk_ref$_ end_lt:uint64 seq_no:uint32 root_hash:bits256 file_hash:bits256 = ExtBlkRef;
type BlkMasterInfo struct {
	Master ExtBlkRef
}

// BlkPrevInfo
// prev_blk_info$_ prev:ExtBlkRef = BlkPrevInfo 0;
// prev_blks_info$_ prev1:^ExtBlkRef prev2:^ExtBlkRef = BlkPrevInfo 1;
type BlkPrevInfo struct { // only manual decoding
	tlb.SumType
	PrevBlkInfo struct {
		Prev ExtBlkRef
	} `tlbSumType:"prev_blk_info$_"`
	PrevBlksInfo struct {
		Prev1 ExtBlkRef // ^ but decodes manually
		Prev2 ExtBlkRef // ^ but decodes manually
	} `tlbSumType:"prev_blks_info$_"`
}

func (i *BlkPrevInfo) UnmarshalTLB(c *boc.Cell, isBlks bool) error { // custom unmarshaler. Not for automatic decoder.
	var res BlkPrevInfo
	if isBlks {
		var prev1, prev2 ExtBlkRef
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c1, &prev1)
		if err != nil {
			return err
		}
		c2, err := c.NextRef()
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c2, &prev2)
		if err != nil {
			return err
		}
		res.SumType = "PrevBlksInfo"
		res.PrevBlksInfo.Prev1 = prev1
		res.PrevBlksInfo.Prev2 = prev2
		return nil
	}
	var prev ExtBlkRef
	err := tlb.Unmarshal(c, &prev)
	if err != nil {
		return err
	}
	res.SumType = "PrevBlkInfo"
	res.PrevBlkInfo.Prev = prev
	*i = res
	return nil
}

// Block
// block#11ef55aa global_id:int32
// info:^BlockInfo value_flow:^ValueFlow
// state_update:^(MERKLE_UPDATE ShardState)
// extra:^BlockExtra = Block;
type Block struct {
	Magic       tlb.Magic `tlb:"block#11ef55aa"`
	GlobalId    int32
	Info        BlockInfo  `tlb:"^"`
	ValueFlow   tlb.Any    `tlb:"^"` // ValueFlow
	StateUpdate tlb.Any    `tlb:"^"` //MerkleUpdate[ShardState] `tlb:"^"` //
	Extra       BlockExtra `tlb:"^"`
}

// TODO: clarify the description of the structure
type BlockHeader struct {
	Magic    tlb.Magic `tlb:"block#11ef55aa"`
	GlobalId int32
	Info     BlockInfo `tlb:"^"`
}

// block_proof#c3 proof_for:BlockIdExt root:^Cell signatures:(Maybe ^BlockSignatures) = BlockProof;
type BlockProof struct {
	Magic      tlb.Magic `tlb:"block_proof#c3"`
	ProofFor   BlockIdExt
	Root       boc.Cell `tlb:"^"`
	Signatures tlb.Maybe[tlb.Ref[BlockSignatures]]
}

// block_signatures#11 validator_info:ValidatorBaseInfo pure_signatures:BlockSignaturesPure = BlockSignatures;
type BlockSignatures struct {
	Magic          tlb.Magic `tlb:"block_signatures#11"`
	ValidatorInfo  ValidatorBaseInfo
	PureSignatures BlockSignaturesPure
}

// block_signatures_pure#_ sig_count:uint32 sig_weight:uint64
//   signatures:(HashmapE 16 CryptoSignaturePair) = BlockSignaturesPure;

type BlockSignaturesPure struct {
	SigCount   uint32
	SigWeight  uint64
	Signatures tlb.HashmapE[CryptoSignaturePair] `tlb:"16bits"`
}

// block_id_ext$_ shard_id:ShardIdent seq_no:uint32
// root_hash:bits256 file_hash:bits256 = BlockIdExt;
type BlockIdExt struct {
	ShardId  ShardIdent
	SeqNo    uint32
	RootHash Hash
	FileHash Hash
}

// ValueFlow
// value_flow ^[ from_prev_blk:CurrencyCollection
// to_next_blk:CurrencyCollection
// imported:CurrencyCollection
// exported:CurrencyCollection ]
// fees_collected:CurrencyCollection
// ^[
// fees_imported:CurrencyCollection
// recovered:CurrencyCollection
// created:CurrencyCollection
// minted:CurrencyCollection
// ] = ValueFlow;
type ValueFlow struct {
	Magic   tlb.Magic `tlb:"value_flow#b8e48dfb"`
	Values1 struct {
		FromPrevBlk CurrencyCollection
		ToNextBlk   CurrencyCollection
		Imported    CurrencyCollection
		Exported    CurrencyCollection
	} `tlb:"^"`
	FeesCollected CurrencyCollection
	Values2       struct {
		FeesImported CurrencyCollection
		Recovered    CurrencyCollection
		Created      CurrencyCollection
		Minted       CurrencyCollection
	} `tlb:"^"`
}

// BlockExtra
// block_extra in_msg_descr:^InMsgDescr
// out_msg_descr:^OutMsgDescr
// account_blocks:^ShardAccountBlocks
// rand_seed:bits256
// created_by:bits256
// custom:(Maybe ^McBlockExtra) = BlockExtra;
type BlockExtra struct {
	Magic         tlb.Magic          `tlb:"block_extra#4a33f6fd"`
	InMsgDescr    InMsgDescr         `tlb:"^"` // tlb.Any `tlb:"^"`
	OutMsgDescr   OutMsgDescr        `tlb:"^"` // tlb.Any `tlb:"^"`
	AccountBlocks ShardAccountBlocks `tlb:"^"` // tlb.Any     `tlb:"^"` //
	RandSeed      Hash
	CreatedBy     Hash
	Custom        tlb.Maybe[tlb.Ref[McBlockExtra]]
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

func convertShardIdent(si ShardIdent) (workchain int32, shard uint64) {
	shard = si.ShardPrefix
	pow2 := uint64(1) << (63 - si.ShardPfxBits)
	shard |= pow2
	return si.WorkchainID, shard
}

func getParents(blkPrevInfo BlkPrevInfo, afterSplit, afterMerge bool, shard uint64, workchain int32) ([]TonNodeBlockIdExt, error) {
	var parents []TonNodeBlockIdExt
	if !afterMerge {
		if blkPrevInfo.SumType != "PrevBlkInfo" {
			return nil, fmt.Errorf("two parent blocks may be only after merge")
		}
		blockID := TonNodeBlockIdExt{
			TonNodeBlockId: TonNodeBlockId{
				Workchain: workchain,
				Seqno:     int32(blkPrevInfo.PrevBlkInfo.Prev.SeqNo),
			},
			FileHash: blkPrevInfo.PrevBlkInfo.Prev.FileHash,
			RootHash: blkPrevInfo.PrevBlkInfo.Prev.RootHash,
		}
		if afterSplit {
			blockID.Shard = int64(shardParent(shard))
			return []TonNodeBlockIdExt{blockID}, nil
		}
		blockID.Shard = int64(shard)
		return []TonNodeBlockIdExt{blockID}, nil
	}

	if blkPrevInfo.SumType != "PrevBlksInfo" {
		return nil, fmt.Errorf("two parent blocks must be after merge")
	}

	parents = append(parents, TonNodeBlockIdExt{
		TonNodeBlockId: TonNodeBlockId{
			Seqno:     int32(blkPrevInfo.PrevBlksInfo.Prev1.SeqNo),
			Shard:     int64(shardChild(shard, true)),
			Workchain: workchain,
		},
		FileHash: blkPrevInfo.PrevBlksInfo.Prev1.FileHash,
		RootHash: blkPrevInfo.PrevBlksInfo.Prev1.RootHash,
	})

	parents = append(parents, TonNodeBlockIdExt{
		FileHash: blkPrevInfo.PrevBlksInfo.Prev2.FileHash,
		RootHash: blkPrevInfo.PrevBlksInfo.Prev2.RootHash,
		TonNodeBlockId: TonNodeBlockId{
			Seqno:     int32(blkPrevInfo.PrevBlksInfo.Prev2.SeqNo),
			Shard:     int64(shardChild(shard, false)),
			Workchain: workchain,
		},
	})

	return parents, nil
}

// masterchain_block_extra#cca5
//   key_block:(## 1)
//   shard_hashes:ShardHashes
//   shard_fees:ShardFees
//   ^[ prev_blk_signatures:(HashmapE 16 CryptoSignaturePair)
//      recover_create_msg:(Maybe ^InMsg)
//      mint_msg:(Maybe ^InMsg) ]
//   config:key_block?ConfigParams
// = McBlockExtra;
type McBlockExtra struct {
	Magic        tlb.Magic `tlb:"masterchain_block_extra#cca5"`
	KeyBlock     bool
	ShardHashes  ShardHashes
	ShardFees    ShardFees
	McExtraOther struct {
		PrevBlkSignatures tlb.HashmapE[CryptoSignaturePair] `tlb:"16bits"`
		RecoverCreate     tlb.Maybe[tlb.Ref[InMsg]]
		MintMsg           tlb.Maybe[tlb.Ref[InMsg]]
	} `tlb:"^"`
	Config ConfigParams
}

func (m *McBlockExtra) UnmarshalTLB(c *boc.Cell, tag string) error {
	sumType, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	if sumType != 0xcca5 {
		return fmt.Errorf("invalid tag")
	}

	err = tlb.Unmarshal(c, &m.KeyBlock)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.ShardHashes)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.ShardFees)
	if err != nil {
		return err
	}
	c1, err := c.NextRef()
	if err != nil && err != boc.ErrNotEnoughRefs {
		return err
	}

	if c1 != nil {
		err = tlb.Unmarshal(c1, &m.McExtraOther)
		if err != nil {
			return err
		}
	}
	if m.KeyBlock {
		err = tlb.Unmarshal(c, &m.Config)
		if err != nil {
			return err
		}
	}
	return nil
}
