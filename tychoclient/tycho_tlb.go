package tychoclient

import (
	"fmt"
	"sort"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Tycho Block - uses automatic TLB parsing with struct tags
// block_tycho#11ef55bb global_id:int32
// info:^BlockInfo value_flow:^ValueFlow
// ^[
//
//	state_update:^(MERKLE_UPDATE ShardState)
//	out_msg_queue_updates:OutMsgQueueUpdates
//
// ]
// extra:^BlockExtra = Block;
type TychoBlock struct {
	Magic     tlb.Magic `tlb:"block_tycho#11ef55bb"`
	GlobalId  int32
	Info      TychoBlockInfo `tlb:"^"`
	ValueFlow tlb.ValueFlow  `tlb:"^"`
	Other     struct {
		StateUpdate        tlb.MerkleUpdate[TychoShardState] `tlb:"^"`
		OutMsgQueueUpdates OutMsgQueueUpdates
	} `tlb:"^"`
	Extra TychoBlockExtra `tlb:"^"`
}

// Tycho BlockInfo - adds gen_utime_ms field
// block_info_tycho#9bc7a988 version:uint32
// not_master:(## 1) after_merge:(## 1) before_split:(## 1) after_split:(## 1)
// want_split:Bool want_merge:Bool key_block:Bool vert_seqno_incr:(## 1)
// flags:(## 8) { flags <= 1 }
// seq_no:# vert_seq_no:#
// shard:ShardIdent gen_utime:uint32
// gen_utime_ms:uint16
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
type TychoBlockInfo struct {
	TychoBlockInfoPart
	GenSoftware *tlb.GlobalVersion
	MasterRef   *tlb.BlkMasterInfo
	PrevRef     tlb.BlkPrevInfo
	PrevVertRef *tlb.BlkPrevInfo
}

type TychoBlockInfoPart struct {
	Version                   uint32
	NotMaster                 bool
	AfterMerge                bool
	BeforeSplit               bool
	AfterSplit                bool
	WantSplit                 bool
	WantMerge                 bool
	KeyBlock                  bool
	VertSeqnoIncr             bool
	Flags                     uint8
	SeqNo                     uint32
	VertSeqNo                 uint32
	Shard                     tlb.ShardIdent
	GenUtime                  uint32
	GenUtimeMs                uint16 // NEW in Tycho: millisecond precision
	StartLt                   uint64
	EndLt                     uint64
	GenValidatorListHashShort uint32
	GenCatchainSeqno          uint32
	MinRefMcSeqno             uint32
	PrevKeyBlockSeqno         uint32
}

// UnmarshalTLB implements custom TLB unmarshaling for Tycho BlockInfo
func (i *TychoBlockInfo) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var data struct {
		Magic     tlb.Magic `tlb:"block_info_tycho#9bc7a988"`
		BlockInfo TychoBlockInfoPart
	}
	err := decoder.Unmarshal(c, &data)
	if err != nil {
		return err
	}

	var res TychoBlockInfo
	res.TychoBlockInfoPart = data.BlockInfo

	if res.Flags&1 == 1 {
		var gs tlb.GlobalVersion
		err = decoder.Unmarshal(c, &gs)
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
		res.MasterRef = &tlb.BlkMasterInfo{}
		err = decoder.Unmarshal(c1, res.MasterRef)
		if err != nil {
			return err
		}
	}

	c1, err := c.NextRef()
	if err != nil {
		return err
	}
	err = res.PrevRef.UnmarshalTLB(c1, data.BlockInfo.AfterMerge, decoder)
	if err != nil {
		return err
	}

	if data.BlockInfo.VertSeqnoIncr {
		c1, err = c.NextRef()
		if err != nil {
			return err
		}
		res.PrevVertRef = &tlb.BlkPrevInfo{}
		err = res.PrevVertRef.UnmarshalTLB(c1, false, decoder)
		if err != nil {
			return err
		}
	}

	*i = res
	return nil
}

// Tycho BlockExtra
// block_extra_tycho#4a33f6fc in_msg_descr:^InMsgDescr
// out_msg_descr:^OutMsgDescr
// account_blocks:^ShardAccountBlocks
// rand_seed:bits256
// created_by:bits256
// custom:(Maybe ^McBlockExtra) = BlockExtra;
type TychoBlockExtra struct {
	Magic           tlb.Magic                                                              `tlb:"block_extra_tycho#4a33f6fc"`
	InMsgDescrCell  boc.Cell                                                               `tlb:"^"`
	OutMsgDescrCell boc.Cell                                                               `tlb:"^"`
	AccountBlocks   tlb.HashmapAugE[tlb.Bits256, tlb.AccountBlock, tlb.CurrencyCollection] `tlb:"^"`
	RandSeed        tlb.Bits256
	CreatedBy       tlb.Bits256
	Custom          tlb.Maybe[tlb.Ref[McBlockExtraTycho]]
}

// InMsgDescr returns the parsed InMsgDescr hashmap
func (extra *TychoBlockExtra) InMsgDescr() (tlb.HashmapAugE[tlb.Bits256, tlb.InMsg, tlb.ImportFees], error) {
	var hashmap tlb.HashmapAugE[tlb.Bits256, tlb.InMsg, tlb.ImportFees]
	if err := tlb.Unmarshal(&extra.InMsgDescrCell, &hashmap); err != nil {
		return tlb.HashmapAugE[tlb.Bits256, tlb.InMsg, tlb.ImportFees]{}, err
	}
	return hashmap, nil
}

// OutMsgDescr returns the parsed OutMsgDescr hashmap
func (extra *TychoBlockExtra) OutMsgDescr() (tlb.HashmapAugE[tlb.Bits256, tlb.OutMsg, tlb.CurrencyCollection], error) {
	var hashmap tlb.HashmapAugE[tlb.Bits256, tlb.OutMsg, tlb.CurrencyCollection]
	if err := tlb.Unmarshal(&extra.OutMsgDescrCell, &hashmap); err != nil {
		return tlb.HashmapAugE[tlb.Bits256, tlb.OutMsg, tlb.CurrencyCollection]{}, err
	}
	return hashmap, nil
}

// ParseTychoBlock parses a Tycho block from BOC data
func ParseTychoBlock(bocData []byte) (*TychoBlock, error) {
	cells, err := boc.DeserializeBoc(bocData)
	if err != nil {
		return nil, err
	}
	if len(cells) == 0 {
		return nil, fmt.Errorf("no cells in BOC")
	}

	var block TychoBlock
	decoder := tlb.NewDecoder()
	decoder.WithDebug()
	err = decoder.Unmarshal(cells[0], &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// OutMsgQueueUpdates - Tycho's queue update tracking
// out_msg_queue_updates#1 diff_hash:bits256 tail_len:uint32 = OutMsgQueueUpdates;
type OutMsgQueueUpdates struct {
	Magic    tlb.Magic `tlb:"#1"`
	DiffHash tlb.Bits256
	TailLen  uint32
}

// TychoShardState - Tycho's shard state with ProcessedUptoInfo
// shard_state_tycho#9023aeee global_id:int32
// shard_id:ShardIdent
// seq_no:uint32 vert_seq_no:#
// gen_utime:uint32 gen_utime_ms:uint16 gen_lt:uint64
// min_ref_mc_seqno:uint32
// processed_upto:^ProcessedUptoInfo
// before_split:(## 1)
// accounts:^ShardAccounts
// ^[ overload_history:uint64 underload_history:uint64
// total_balance:CurrencyCollection
// total_validator_fees:CurrencyCollection
// libraries:(HashmapE 256 LibDescr)
// master_ref:(Maybe BlkMasterInfo) ]
// custom:(Maybe ^McStateExtra)
// = ShardStateUnsplit;
type TychoShardState struct {
	tlb.SumType
	UnsplitState struct {
		Value TychoShardStateUnsplit
	} `tlbSumType:"_"`
	SplitState struct {
		Left  TychoShardStateUnsplit `tlb:"^"`
		Right TychoShardStateUnsplit `tlb:"^"`
	} `tlbSumType:"split_state#5f327da5"`
}

func (s *TychoShardState) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	sumType, err := c.ReadUint(32)
	if err != nil {
		return err
	}
	switch sumType {
	case 0x5f327da5: // split_state
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		if c1.CellType() != boc.PrunedBranchCell {
			err = decoder.Unmarshal(c1, &s.SplitState.Left)
			if err != nil {
				return err
			}
		} else {
			s.SplitState.Left = TychoShardStateUnsplit{}
		}
		c1, err = c.NextRef()
		if err != nil {
			return err
		}
		if c1.CellType() != boc.PrunedBranchCell {
			if err := decoder.Unmarshal(c1, &s.SplitState.Right); err != nil {
				return err
			}
		} else {
			s.SplitState.Right = TychoShardStateUnsplit{}
		}
		s.SumType = "SplitState"
	case 0x9023aeee: // shard_state_tycho
		// Reset and unmarshal the full struct
		c.ResetCounters()
		err = decoder.Unmarshal(c, &s.UnsplitState.Value)
		if err != nil {
			return err
		}
		s.SumType = "UnsplitState"
	default:
		return fmt.Errorf("invalid TychoShardState tag: 0x%x", sumType)
	}
	return nil
}

// TychoShardStateUnsplit - Tycho's unsplit shard state
type TychoShardStateUnsplit struct {
	Magic         tlb.Magic `tlb:"shard_state_tycho#9023aeee"`
	GlobalId      int32
	ShardId       tlb.ShardIdent
	SeqNo         uint32
	VertSeqNo     uint32
	GenUtime      uint32
	GenUtimeMs    uint16 // Tycho-specific millisecond precision
	GenLt         uint64
	MinRefMcSeqno uint32
	ProcessedUpto ProcessedUptoInfo `tlb:"^"`
	BeforeSplit   bool
	Accounts      tlb.HashmapAugE[tlb.Bits256, tlb.ShardAccount, tlb.DepthBalanceInfo] `tlb:"^"`
	Other         TychoShardStateUnsplitOther                                          `tlb:"^"`
	Custom        tlb.Maybe[tlb.Ref[McStateExtraTycho]]
}

// TychoShardStateUnsplitOther - nested cell in TychoShardStateUnsplit
type TychoShardStateUnsplitOther struct {
	OverloadHistory    uint64
	UnderloadHistory   uint64
	TotalBalance       tlb.CurrencyCollection
	TotalValidatorFees tlb.CurrencyCollection
	Libraries          tlb.HashmapE[tlb.Bits256, tlb.LibDescr]
	MasterRef          tlb.Maybe[tlb.BlkMasterInfo]
}

// ProcessedUptoInfo - Tycho's processed messages tracking
// processed_upto_info#00
// partitions:(HashmapE 16 ProcessedUptoPartition)
// msgs_exec_params:(Maybe ^MsgsExecutionParams)
// = ProcessedUptoInfo;
type ProcessedUptoInfo struct {
	Partitions     tlb.HashmapE[tlb.Uint16, ProcessedUptoPartition]
	MsgsExecParams tlb.Maybe[tlb.Ref[MsgsExecutionParams]]
}

// ProcessedUptoPartition - per-partition processing status
// processedUptoPartition#00
// externals:ExternalsProcessedUpto
// internals:InternalsProcessedUpto
// = ProcessedUptoPartition;
type ProcessedUptoPartition struct {
	Externals ExternalsProcessedUpto
	Internals InternalsProcessedUpto
}

// ExternalsProcessedUpto - external messages processing status
// externalsProcessedUpto#00
// processed_to_anchor_id:uint32
// processed_to_msgs_offset:uint64
// ranges:(HashmapE 32 ExternalsRange)
// = ExternalsProcessedUpto;
type ExternalsProcessedUpto struct {
	ProcessedToAnchorId   uint32
	ProcessedToMsgsOffset uint64
	Ranges                tlb.HashmapE[tlb.Uint32, ExternalsRange]
}

// ExternalsRange - range of external messages
// externalsRange#00
// from_anchor_id:uint32
// from_msgs_offset:uint64
// to_anchor_id:uint32
// to_msgs_offset:uint64
// chain_time:uint64
// skip_offset:uint32
// processed_offset:uint32
// = ExternalsRange;
type ExternalsRange struct {
	FromAnchorId    uint32
	FromMsgsOffset  uint64
	ToAnchorId      uint32
	ToMsgsOffset    uint64
	ChainTime       uint64
	SkipOffset      uint32
	ProcessedOffset uint32
}

// InternalsProcessedUpto - internal messages processing status
// internalsProcessedUpto#00
// processed_to:(HashmapE 96 ProcessedUpto)
// ranges:(HashmapE 32 InternalsRange)
// = InternalsProcessedUpto;
type InternalsProcessedUpto struct {
	ProcessedTo tlb.HashmapE[tlb.Bits96, tlb.ProcessedUpto]
	Ranges      tlb.HashmapE[tlb.Uint32, InternalsRange]
}

// InternalsRange - range of internal messages
// internalsRange#00
// skip_offset:uint32
// processed_offset:uint32
// shards:(HashmapE 96 ShardRange)
// = InternalsRange;
type InternalsRange struct {
	SkipOffset      uint32
	ProcessedOffset uint32
	Shards          tlb.HashmapE[tlb.Bits96, ShardRange]
}

// ShardRange - range within a shard
// shardRange#00 from:uint64 to:uint64 = ShardRange;
type ShardRange struct {
	From uint64
	To   uint64
}

// MsgsExecutionParams - message execution parameters
// msgs_execution_params_tycho#00 or #01
type MsgsExecutionParams struct {
	BufferLimit            uint32
	GroupLimit             uint16
	GroupVertSize          uint16
	ExternalsExpireTimeout uint16
	OpenRangesLimit        uint16
	Par0IntMsgsCountLimit  uint32
	Par0ExtMsgsCountLimit  uint32
	GroupSlotsFractions    tlb.HashmapE[tlb.Uint16, uint8]
	RangeMessagesLimit     uint32 // Only in v2 (#01)
}

// TychoShardDescr - Tycho's shard descriptor
// shard_descr_tycho#a seq_no:uint32 reg_mc_seqno:uint32
// start_lt:uint64 end_lt:uint64
// root_hash:bits256 file_hash:bits256
// before_split:Bool before_merge:Bool
// want_split:Bool want_merge:Bool
// nx_cc_updated:Bool top_sc_block_updated:Bool flags:(## 2) { flags = 0 }
// next_catchain_seqno:uint32 ext_processed_to_anchor_id:uint32
// min_ref_mc_seqno:uint32 gen_utime:uint32
// split_merge_at:FutureSplitMerge
// ^[ fees_collected:CurrencyCollection
//
//	funds_created:CurrencyCollection ] = ShardDescr;
type TychoShardDescr struct {
	Magic                  tlb.Magic `tlb:"shard_descr_tycho#a"`
	SeqNo                  uint32
	RegMcSeqno             uint32
	StartLT                uint64
	EndLT                  uint64
	RootHash               tlb.Bits256
	FileHash               tlb.Bits256
	BeforeSplit            bool
	BeforeMerge            bool
	WantSplit              bool
	WantMerge              bool
	NXCCUpdated            bool
	TopScBlockUpdated      bool // NEW in Tycho
	Flags                  tlb.Uint2
	NextCatchainSeqNo      uint32
	ExtProcessedToAnchorId uint32 // NEW in Tycho (replaces NextValidatorShard)
	MinRefMcSeqNo          uint32
	GenUTime               uint32
	// TODO: split_merge_at:FutureSplitMerge - not implemented yet
	FeesCollectedFundsCreated struct {
		FeesCollected tlb.CurrencyCollection
		FundsCreated  tlb.CurrencyCollection
	} `tlb:"^"`
}

// TychoShardInfoBinTree - Binary tree of Tycho shard descriptors
type TychoShardInfoBinTree struct {
	BinTree tlb.BinTree[TychoShardDescr]
}

// McBlockExtraTycho - Tycho's masterchain block extra with TychoShardDescr
// masterchain_block_extra#cca5 key_block:Bool
// shard_hashes:ShardHashes
// shard_fees:ShardFees
// ^[ prev_blk_signatures:(HashmapE 16 CryptoSignaturePair)
//
//	recover_create_msg:(Maybe ^InMsg)
//	mint_msg:(Maybe ^InMsg) ]
//
// config:key_block?ConfigParams
// = McBlockExtra;
type McBlockExtraTycho struct {
	Magic        tlb.Magic `tlb:"masterchain_block_extra#cca5"`
	KeyBlock     bool
	ShardHashes  tlb.HashmapE[tlb.Uint32, tlb.Ref[TychoShardInfoBinTree]] // Use Tycho version
	ShardFees    tlb.ShardFees
	McExtraOther struct {
		PrevBlkSignatures tlb.HashmapE[tlb.Uint16, tlb.CryptoSignaturePair]
		RecoverCreate     tlb.Maybe[tlb.Ref[tlb.InMsg]]
		MintMsg           tlb.Maybe[tlb.Ref[tlb.InMsg]]
	} `tlb:"^"`
	Config tlb.ConfigParams
}

// McStateExtra
// masterchain_state_extra#cc26
// shard_hashes:ShardHashes
// config:ConfigParams
// ^[ flags:(## 16) { flags <= 1 }
// validator_info:ValidatorInfo
// prev_blocks:OldMcBlocksInfo
// after_key_block:Bool
// last_key_block:(Maybe ExtBlkRef)
// block_create_stats:(flags . 0)?BlockCreateStats ]
// global_balance:CurrencyCollection
// = McStateExtra;
type McStateExtraTycho struct {
	Magic         tlb.Magic `tlb:"masterchain_state_extra#cc26"`
	ShardHashes   tlb.HashmapE[tlb.Uint32, tlb.Ref[TychoShardInfoBinTree]]
	Config        tlb.ConfigParams
	Other         McStateExtraOtherTycho `tlb:"^"`
	GlobalBalance tlb.CurrencyCollection
}

type GenesisInfo struct {
	StartRound    uint32
	GenesisMillis uint64
}

type ConsensusInfo struct {
	VsetSwitchRound         uint32
	PrevVsetSwitchRound     uint32
	GenesisInfo             GenesisInfo
	PrevShuffleMcValidators bool
}

type McStateExtraOtherTycho struct {
	Flags            uint16
	ValidatorInfo    tlb.ValidatorInfo
	PrevBlocks       tlb.HashmapAugE[tlb.Uint32, tlb.KeyExtBlkRef, tlb.KeyMaxLt]
	AfterKeyBlock    bool
	LastKeyBlock     tlb.Maybe[tlb.ExtBlkRef]
	BlockCreateStats *tlb.BlockCreateStats
	ConsensusInfo    *ConsensusInfo
}

func (m *McStateExtraOtherTycho) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	flags, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	m.Flags = uint16(flags)
	err = decoder.Unmarshal(c, &m.ValidatorInfo)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.PrevBlocks)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.AfterKeyBlock)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.LastKeyBlock)
	if err != nil {
		return err
	}
	if m.Flags&1<<0 > 0 {
		err = decoder.Unmarshal(c, &m.BlockCreateStats)
		if err != nil {
			return err
		}
	}
	if m.Flags&1<<2 > 0 {
		err = decoder.Unmarshal(c, &m.ConsensusInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *McBlockExtraTycho) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	sumType, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	if sumType != 0xcca5 {
		return fmt.Errorf("invalid tag for McBlockExtraTycho: 0x%x", sumType)
	}

	err = decoder.Unmarshal(c, &m.KeyBlock)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.ShardHashes)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.ShardFees)
	if err != nil {
		return err
	}
	c1, err := c.NextRef()
	if err != nil && err != boc.ErrNotEnoughRefs {
		return err
	}

	if c1 != nil {
		err = decoder.Unmarshal(c1, &m.McExtraOther)
		if err != nil {
			return err
		}
	}
	if m.KeyBlock {
		err = decoder.Unmarshal(c, &m.Config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TychoShardState) AccountBalances() map[tlb.Bits256]tlb.CurrencyCollection {
	switch s.SumType {
	case "UnsplitState":
		accounts := s.UnsplitState.Value.Accounts.Keys()
		balances := make(map[tlb.Bits256]tlb.CurrencyCollection, len(accounts))
		for i, shardAccount := range s.UnsplitState.Value.Accounts.Values() {
			c, ok := shardAccount.Account.CurrencyCollection()
			if !ok {
				continue
			}
			balances[accounts[i]] = c
		}
		return balances
	default:
		leftAccounts := s.SplitState.Left.Accounts.Keys()
		rightAccounts := s.SplitState.Right.Accounts.Keys()
		balances := make(map[tlb.Bits256]tlb.CurrencyCollection, len(leftAccounts)+len(rightAccounts))
		for i, shardAccount := range s.SplitState.Left.Accounts.Values() {
			c, ok := shardAccount.Account.CurrencyCollection()
			if !ok {
				continue
			}
			balances[leftAccounts[i]] = c
		}
		for i, shardAccount := range s.SplitState.Right.Accounts.Values() {
			c, ok := shardAccount.Account.CurrencyCollection()
			if !ok {
				continue
			}
			balances[rightAccounts[i]] = c
		}
		return balances
	}
}

func (b *TychoBlock) TransactionsQuantity() int {
	quantity := 0
	for _, accountBlock := range b.Extra.AccountBlocks.Values() {
		quantity += len(accountBlock.Transactions.Keys())
	}
	return quantity
}

func (b *TychoBlock) AllTransactions() []*tlb.Transaction {
	transactions := make([]*tlb.Transaction, 0, b.TransactionsQuantity())
	for _, accountBlock := range b.Extra.AccountBlocks.Values() {
		for i := range accountBlock.Transactions.Values() {
			transactions = append(transactions, &accountBlock.Transactions.Values()[i].Value)
		}
	}
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Lt < transactions[j].Lt
	})
	return transactions
}

func (extra *TychoBlockExtra) OutMsgDescrLength() (int, error) {
	// TODO: refactor
	// construct a dummy extra with necessary fields for BlockExtra.OutMsgDescrLength
	dummyExtra := tlb.BlockExtra{
		OutMsgDescrCell: extra.OutMsgDescrCell,
	}
	return dummyExtra.OutMsgDescrLength()
}

func (extra *TychoBlockExtra) InMsgDescrLength() (int, error) {
	// TODO: refactor
	// construct a dummy extra with necessary fields for BlockExtra.InMsgDescrLength
	dummyExtra := tlb.BlockExtra{
		InMsgDescrCell: extra.InMsgDescrCell,
	}
	return dummyExtra.InMsgDescrLength()
}

func GetParents(i TychoBlockInfo) ([]ton.BlockIDExt, error) {
	// TODO: refactor
	// construct a dummy block with necessary fields for tongo.GetParents
	var dummyBlock tlb.BlockInfo
	dummyBlock.Shard = i.Shard
	dummyBlock.PrevRef = i.PrevRef
	dummyBlock.AfterSplit = i.AfterSplit
	dummyBlock.AfterMerge = i.AfterMerge
	return ton.GetParents(dummyBlock)
}

func ShardIDs(extra McBlockExtraTycho) []ton.BlockIDExt {
	items := extra.ShardHashes.Items()
	shardsCount := 0
	for _, item := range items {
		for _, x := range item.Value.Value.BinTree.Values {
			if x.SeqNo == 0 {
				continue
			}
			shardsCount += 1
		}
	}
	if shardsCount == 0 {
		return nil
	}
	shards := make([]ton.BlockIDExt, 0, shardsCount)
	for _, item := range items {
		for _, shardDesc := range item.Value.Value.BinTree.Values {
			if shardDesc.SeqNo == 0 {
				continue
			}
			shards = append(shards, toBlockID(shardDesc, int32(item.Key)))
		}
	}
	return shards
}

func toBlockID(s TychoShardDescr, workchain int32) ton.BlockIDExt {
	return ton.BlockIDExt{
		BlockID: ton.BlockID{
			Workchain: workchain,
			Shard:     0x8000000000000000, // no sharding, the field containing shard information was removed from the block
			Seqno:     s.SeqNo,
		},
		RootHash: ton.Bits256(s.RootHash),
		FileHash: ton.Bits256(s.FileHash),
	}
}
