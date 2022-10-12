package tongo

import (
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// MerkleProof
// !merkle_proof#03 {X:Type} virtual_hash:bits256 depth:uint16 virtual_root:^X = MERKLE_PROOF X;
type MerkleProof[T any] struct {
	tlb.SumType
	MerkleProof struct {
		VirtualHash Hash
		Depth       boc.BitString `tlb:"16bits"`
		VirtualRoot T             `tlb:"^"`
	} `tlbSumType:"!merkle_proof#03"`
}

// ShardStateUnsplit
// shard_state#9023afe2 global_id:int32
// shard_id:ShardIdent
// seq_no:uint32 vert_seq_no:#
// gen_utime:uint32 gen_lt:uint64
// min_ref_mc_seqno:uint32
// out_msg_queue_info:^OutMsgQueueInfo
// before_split:(## 1)
// accounts:^ShardAccounts
// ^[ overload_history:uint64 underload_history:uint64
// total_balance:CurrencyCollection
// total_validator_fees:CurrencyCollection
// libraries:(HashmapE 256 LibDescr)
// master_ref:(Maybe BlkMasterInfo) ]
// custom:(Maybe ^McStateExtra)
// = ShardStateUnsplit;
type ShardStateUnsplit struct {
	tlb.SumType
	ShardStateUnsplit struct {
		GlobalID        int32
		ShardID         ShardIdent
		SeqNo           uint32
		VertSeqNo       uint32
		GenUtime        uint32
		GenLt           uint64
		MinRefMcSeqno   uint32
		OutMsgQueueInfo OutMsgQueueInfo `tlb:"^"` // TODO: implement decoding OutMsgQueueInfo fields
		BeforeSplit     bool
		Accounts        tlb.Ref[tlb.Any] // TODO: implement decoding Accounts fields
		Other           tlb.Ref[tlb.Any] // TODO: implement decoding Other fields
		Custom          tlb.Maybe[tlb.Ref[McStateExtra]]
	} `tlbSumType:"shard_state#9023afe2"`
	// SplitState struct{} `tlbSumType:"split_state#5f327da5"` // rare case
}

// _ ShardStateUnsplit = ShardState;
// split_state#5f327da5 left:^ShardStateUnsplit right:^ShardStateUnsplit = ShardState;
type ShardState struct {
	tlb.SumType
	ShardStateUnsplit
	SplitState struct {
		Left  tlb.Ref[ShardStateUnsplit]
		Right tlb.Ref[ShardStateUnsplit]
	} `tlbSumType:"split_state#5f327da5"`
}

// ShardIdent
// shard_ident$00 shard_pfx_bits:(#<= 60)
// workchain_id:int32 shard_prefix:uint64 = ShardIdent;
type ShardIdent struct {
	tlb.SumType
	ShardIdent struct {
		ShardPfxBits uint64 `tlb:"6bits"` // TODO: implement lim uint tag
		WorkchainID  int32
		ShardPrefix  uint64
	} `tlbSumType:"shardident$00"`
}

// func (s *ShardIdent) UnmarshalTLB(c *boc.Cell, tag string) error {
// 	t, err := c.ReadUint(2)
// 	if err != nil {
// 		return err
// 	}
// 	if t != 0 {
// 		return fmt.Errorf("invalid tag")
// 	}
// 	prefixBits, err := c.ReadLimUint(60)
// 	if err != nil {
// 		return err
// 	}
// 	workchain, err := c.ReadInt(32)
// 	if err != nil {
// 		return err
// 	}
// 	prefix, err := c.ReadUint(64)
// 	if err != nil {
// 		return err
// 	}
// 	s.SumType = "ShardIdent"
// 	s.ShardIdent.ShardPfxBits = uint64(prefixBits)
// 	s.ShardIdent.WorkchainID = int32(workchain)
// 	s.ShardIdent.ShardPrefix = prefix
// 	return nil
// }

// _ out_queue:OutMsgQueue proc_info:ProcessedInfo
// ihr_pending:IhrPendingInfo = OutMsgQueueInfo;
type OutMsgQueueInfo struct {
	OutQueue  OutMsgQueue
	ProcInfo  ProcessedInfo
	IhrPendig IhrPendingInfo
}

// _ (HashmapAugE 352 EnqueuedMsg uint64) = OutMsgQueue;
type OutMsgQueue struct {
	//Queue tlb.HashmapE[EnqueuedMsg] `tlb:"352bits"`
	Queue tlb.HashmapAugE[EnqueuedMsg, uint64] `tlb:"352bits"`
}

// _ enqueued_lt:uint64 out_msg:^MsgEnvelope = EnqueuedMsg;
type EnqueuedMsg struct {
	EnqueuedLt uint64
	OutMsg     tlb.Ref[MsgEnvelope]
}

// 	msg_envelope#4 cur_addr:IntermediateAddress
//  next_addr:IntermediateAddress fwd_fee_remaining:Grams
//  msg:^(Message Any) = MsgEnvelope;
type MsgEnvelope struct {
	tlb.SumType
	MsgEnvelope struct {
		CurrentAddress  IntermediateAddress
		NextAddress     IntermediateAddress
		FwdFeeRemaining Grams
		Msg             tlb.Ref[Message[tlb.Any]]
	} `tlbSumType:"msg_envelope#4"`
}

// interm_addr_regular$0 use_dest_bits:(#<= 96) = IntermediateAddress;
// interm_addr_simple$10 workchain_id:int8 addr_pfx:uint64 = IntermediateAddress;
// interm_addr_ext$11 workchain_id:int32 addr_pfx:uint64 = IntermediateAddress;
type IntermediateAddress struct {
	tlb.SumType
	IntermediateAddressRegular struct {
		UseDestBits boc.BitString `tlb:"7bits"`
	} `tlbSumType:"interm_addr_regular$0"`
	IntermediateAddressSimple struct {
		WorkchainId   boc.BitString `tlb:"8bits"`
		AddressPrefix uint64
	} `tlbSumType:"interm_addr_simple$10"`
	IntermediateAddressExt struct {
		WorkchainId   int32
		AddressPrefix uint64
	} `tlbSumType:"interm_addr_ext$11"`
}

// _ (HashmapE 96 ProcessedUpto) = ProcessedInfo;
type ProcessedInfo struct {
	ProcessedInfo tlb.HashmapE[ProcessedUpto] `tlb:"96bits"`
}

// processed_upto$_ last_msg_lt:uint64 last_msg_hash:bits256 = ProcessedUpto;
type ProcessedUpto struct {
	LastMsg     uint64
	LastMsgHash Hash
}

// _ (HashmapE 320 IhrPendingSince) = IhrPendingInfo;
type IhrPendingInfo struct {
	IhrPendingInfo tlb.HashmapE[IhrPendingSince] `tlb:"320bits"`
}

// ihr_pending$_ import_lt:uint64 = IhrPendingSince;
type IhrPendingSince struct {
	ImportLt uint64
}

// _ (HashmapAugE 256 ShardAccount DepthBalanceInfo) = ShardAccounts;
type ShardAccounts struct {
	//Accounts tlb.HashmapE[ShardAccount] `tlb:"256bits"`
	Accounts tlb.HashmapAugE[ShardAccount, DepthBalanceInfo] `tlb:"256bits"`
}

// depth_balance$_ split_depth:(#<= 30) balance:CurrencyCollection = DepthBalanceInfo;
type DepthBalanceInfo struct {
	SplitDepth boc.BitString `tlb:"5bits"`
	Balance    CurrencyCollection
}

// ^[ overload_history:uint64 underload_history:uint64
// total_balance:CurrencyCollection
// total_validator_fees:CurrencyCollection
// libraries:(HashmapE 256 LibDescr)
// master_ref:(Maybe BlkMasterInfo) ]
type ShardStateUnsplitOther struct {
	OverloadHistory    uint64
	UnderloadHistory   uint64
	TotalBalance       CurrencyCollection
	TotalValidatorFees CurrencyCollection
	Libraries          tlb.HashmapE[LibDescr] `tlb:"256bits"`
	MasterRef          tlb.Maybe[BlkMasterInfo]
}

// shared_lib_descr$00 lib:^Cell publishers:(Hashmap 256 True)
//   = LibDescr;
type LibDescr struct {
	tlb.SumType
	LibDescr struct {
		Lib        tlb.Ref[boc.Cell]
		Publishers LibDescrPub
	} `tlbSumType:"shared_lib_descr$00"`
}

type LibDescrPub struct {
	Pub tlb.Hashmap[bool] `tlb:"256bits"`
}

// func (s *ShardIdent) UnmarshalTLB(c *boc.Cell, tag string) error {
// 	t, err := c.ReadUint(2)
// 	if err != nil {
// 		return err
// 	}
// 	if t != 0 {
// 		return fmt.Errorf("invalid tag")
// 	}
// 	prefixBits, err := c.ReadLimUint(60)
// 	if err != nil {
// 		return err
// 	}
// 	workchain, err := c.ReadInt(32)
// 	if err != nil {
// 		return err
// 	}
// 	prefix, err := c.ReadUint(64)
// 	if err != nil {
// 		return err
// 	}
// 	s.SumType = "ShardIdent"
// 	s.ShardIdent.ShardPfxBits = uint64(prefixBits)
// 	s.ShardIdent.WorkchainID = int32(workchain)
// 	s.ShardIdent.ShardPrefix = prefix
// 	return nil
// }

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
type McStateExtra struct {
	tlb.SumType
	McStateExtra struct {
		ShardHashes   ShardHashes
		Config        ConfigParams
		Other         McStateExtraOther `tlb:"^"`
		GlobalBalance CurrencyCollection
	} `tlbSumType:"masterchain_state_extra#cc26"`
}

// ShardHashes
// _ (HashmapE 32 ^(BinTree ShardDescr)) = ShardHashes;
type ShardHashes struct {
	Hashes tlb.HashmapE[tlb.Ref[ShardInfoBinTree]] `tlb:"32bits"`
}

type ConfigHashMap struct {
	Hashmap tlb.Hashmap[tlb.Ref[boc.Cell]] `tlb:"32bits"`
}

// ConfigParams
// _ config_addr:bits256 config:^(Hashmap 32 ^Cell)
// = ConfigParams;
type ConfigParams struct {
	ConfigAddr Hash
	Config     tlb.Ref[ConfigHashMap] //`tlb:"32bits"` // TODO: implement decoding config
}

// type ConfigParams struct {
// 	ConfigAddr Hash
// 	Config     tlb.Ref[tlb.Any] `tlb:"32bits"` // TODO: implement decoding config
// }

// ^[ flags:(## 16) { flags <= 1 }
// validator_info:ValidatorInfo
// prev_blocks:OldMcBlocksInfo
// after_key_block:Bool
// last_key_block:(Maybe ExtBlkRef)
// block_create_stats:(flags . 0)?BlockCreateStats ]
type McStateExtraOther struct {
	tlb.SumType
	// Flags boc.BitString `tlb:"16bits"`
	McStateExtra struct {
		ValidatorInfo ValidatorInfo
		PrevBlocks    OldMcBlocksInfo
		AfterKeyBlock bool
		LastKeyBlock  tlb.Maybe[ExtBlkRef]
	} `tlbSumType:"flags#0000"`
	McStateExtraWithBlockStats struct {
		ValidatorInfo    ValidatorInfo
		PrevBlocks       OldMcBlocksInfo
		AfterKeyBlock    bool
		LastKeyBlock     tlb.Maybe[ExtBlkRef]
		BlockCreateStats BlockCreateStats
	} `tlbSumType:"flags#0001"`
}

// func (mc *McStateExtraOther) UnmarshalTLB(c *boc.Cell, tag string) error {
// 	err := mc.Flags.UnmarshalTLB(c, "16bits")
// 	if err != nil {
// 		return err
// 	}
// 	err = tlb.Unmarshal(c, &mc.ValidatorInfo)
// 	if err != nil {
// 		return err
// 	}
// 	err = tlb.Unmarshal(c, &mc.PrevBlocks)
// 	if err != nil {
// 		return err
// 	}
// 	err = tlb.Unmarshal(c, &mc.AfterKeyBlock)
// 	if err != nil {
// 		return err
// 	}
// 	err = tlb.Unmarshal(c, &mc.LastKeyBlock)
// 	if err != nil {
// 		return err
// 	}
// 	if (binary.BigEndian.Uint16(mc.Flags.Buffer()) & 0x1) == 0x1 {
// 		err = tlb.Unmarshal(c, &mc.BlockCreateStats)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// validator_info$_
//   validator_list_hash_short:uint32
//   catchain_seqno:uint32
//   nx_cc_updated:Bool
// = ValidatorInfo;
type ValidatorInfo struct {
	ValidatorListHashShort uint32
	CatchainSeqno          uint32
	NxCcUpdated            bool
}

// _ (HashmapAugE 32 KeyExtBlkRef KeyMaxLt) = OldMcBlocksInfo
type OldMcBlocksInfo struct {
	//Info tlb.HashmapE[KeyExtBlkRef] `tlb:"32bits"`
	Info tlb.HashmapAugE[KeyExtBlkRef, KeyMaxLt] `tlb:"32bits"`
}

// func (mc *OldMcBlocksInfo) UnmarshalTLB(c *boc.Cell, tag string) error {
// 	// hashmapaug
// 	r, err := c.ReadBit()
// 	if err != nil {
// 		return err
// 	}

// 	if r {
// 		// cc := boc.NewCell()
// 		// cc.WriteBit(true)
// 		// cc.WriteBitString(c.ReadRemainingBits())
// 		// c.WriteBitString(cc.ReadRemainingBits())

// 		err := mc.Info.UnmarshalTLB(c, "32bits")
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	var keyMaxLt KeyMaxLt
// 	err = tlb.Unmarshal(c, &keyMaxLt)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// _ key:Bool max_end_lt:uint64 = KeyMaxLt;
type KeyMaxLt struct {
	Key      bool
	MaxEndLt uint64
}

// _ key:Bool blk_ref:ExtBlkRef = KeyExtBlkRef;
type KeyExtBlkRef struct {
	Key    bool
	BlkRef ExtBlkRef
}

// block_create_stats#17 counters:(HashmapE 256 CreatorStats) = BlockCreateStats;
// block_create_stats_ext#34 counters:(HashmapAugE 256 CreatorStats uint32) = BlockCreateStats;
type BlockCreateStats struct {
	tlb.SumType
	BlockCreateStats struct {
		Counters tlb.HashmapE[CreatorStats] `tlb:"256bits"`
	} `tlbSumType:"block_create_stats#17"`
	BlockCreateStatsExt struct {
		Counters tlb.HashmapAugE[CreatorStats, uint32] `tlb:"256bits"`
	} `tlbSumType:"block_create_statsext#34"`
}

// creator_info#4 mc_blocks:Counters shard_blocks:Counters = CreatorStats;
type CreatorStats struct {
	tlb.SumType
	CreatorStats struct {
		McBlocks    Counters
		ShardBlocks Counters
	} `tlbSumType:"creator_info#4"`
}

// counters#_ last_updated:uint32 total:uint64 cnt2048:uint64 cnt65536:uint64 = Counters;

type Counters struct {
	LastUpdated uint32
	Total       uint64
	Cnt2048     uint64
	Cnt65536    uint64
}
