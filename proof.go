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
		VirtualRoot tlb.Ref[T]
	} `tlbSumType:"!merkle_proof#03"`
}

// !merkle_update#02 {X:Type} old_hash:bits256 new_hash:bits256
//   old:^X new:^X = MERKLE_UPDATE X;
type MerkleUpdate[T any] struct {
	tlb.SumType
	MerkleUpdate struct {
		OldHash Hash
		NewHash Hash
		Old     tlb.Ref[T]
		New     tlb.Ref[T]
	} `tlbSumType:"!merkle_update#02"`
}

// block_info#9bc7a987 version:uint32
//   not_master:(## 1)
//   after_merge:(## 1) before_split:(## 1)
//   after_split:(## 1)
//   want_split:Bool want_merge:Bool
//   key_block:Bool vert_seqno_incr:(## 1)
//   flags:(## 8) { flags <= 1 }
//   seq_no:# vert_seq_no:# { vert_seq_no >= vert_seqno_incr }
//   { prev_seq_no:# } { ~prev_seq_no + 1 = seq_no }
//   shard:ShardIdent gen_utime:uint32
//   start_lt:uint64 end_lt:uint64
//   gen_validator_list_hash_short:uint32
//   gen_catchain_seqno:uint32
//   min_ref_mc_seqno:uint32
//   prev_key_block_seqno:uint32
//   gen_software:flags . 0?GlobalVersion
//   master_ref:not_master?^BlkMasterInfo
//   prev_ref:^(BlkPrevInfo after_merge)
//   prev_vert_ref:vert_seqno_incr?^(BlkPrevInfo 0)
//   = BlockInfo;
type BlockInfo struct {
	tlb.SumType
	BlockInfo struct {
		Version                   uint32
		NotMaster                 bool
		AfterMerge                bool
		BeforeSplit               bool
		AfterSplit                bool
		WantSplit                 bool
		WantMerge                 bool
		KeyBlock                  bool
		VertSeqNoIncr             bool
		Flags                     uint32 `tlb:"8bits"`
		SeqNo                     uint32
		VertSeqNo                 uint32
		Shard                     ShardIdent
		GenUtime                  uint32
		StartLt                   uint64
		EndLt                     uint64
		GenValidatorListHashShort uint32
		GenCatchainSeqNo          uint32
		MinRefMcSeqNo             uint32
		PrevKeyBlockSeqNo         uint32
		GenSoftware               GlobalVersion
		MasterRef                 tlb.Ref[BlkMasterInfo]
		PrevRefNotAfterMerge      tlb.Ref[BlkPrevInfo0]
		PrevRefAfterMerge         tlb.Ref[BlkPrevInfo1]
		PrevVertRef               tlb.Ref[BlkPrevInfo0]
	} `tlbSumType:"block_info#9bc7a987"`
}

func (b *BlockInfo) UnmarshalTLB(c *boc.Cell, tag string) error {
	sumType, err := c.ReadUint(32)
	if err != nil {
		return err
	}
	if sumType != 0x87a9c79b {
		return nil
	}
	b.SumType = "BlockInfo"
	err = tlb.Unmarshal(c, &b.BlockInfo.Version)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.NotMaster)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.AfterMerge)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.BeforeSplit)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.AfterSplit)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.WantSplit)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.WantMerge)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.KeyBlock)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.VertSeqNoIncr)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.Flags)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.SeqNo)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.VertSeqNo)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.Shard)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.GenUtime)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.StartLt)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.EndLt)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.GenValidatorListHashShort)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.GenCatchainSeqNo)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.MinRefMcSeqNo)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &b.BlockInfo.PrevKeyBlockSeqNo)
	if err != nil {
		return err
	}
	if b.BlockInfo.Flags == 1 {
		err = tlb.Unmarshal(c, &b.BlockInfo.GenSoftware)
		if err != nil {
			return err
		}
	}
	if b.BlockInfo.NotMaster {
		err = tlb.Unmarshal(c, &b.BlockInfo.MasterRef)
		if err != nil {
			return err
		}
	}
	if b.BlockInfo.AfterMerge {
		err = tlb.Unmarshal(c, &b.BlockInfo.PrevRefAfterMerge)
		if err != nil {
			return err
		}
	} else {
		err = tlb.Unmarshal(c, &b.BlockInfo.PrevRefNotAfterMerge)
		if err != nil {
			return err
		}
	}
	if b.BlockInfo.VertSeqNoIncr {
		err = tlb.Unmarshal(c, &b.BlockInfo.PrevVertRef)
		if err != nil {
			return err
		}
	}

	return nil
}

// prev_blk_info$_ prev:ExtBlkRef = BlkPrevInfo 0;
// prev_blks_info$_ prev1:^ExtBlkRef prev2:^ExtBlkRef = BlkPrevInfo 1;
type BlkPrevInfo0 struct {
	Prev ExtBlkRef
}

type BlkPrevInfo1 struct {
	Prev1 tlb.Ref[ExtBlkRef]
	Prev2 tlb.Ref[ExtBlkRef]
}

// capabilities#c4 version:uint32 capabilities:uint64 = GlobalVersion;
type GlobalVersion struct {
	GlobalVersion struct {
		Version      uint32
		Capabilities uint64
	} `tlbSumType:"capabilities#c4"`
}

// block#11ef55aa global_id:int32
//   info:^BlockInfo value_flow:^ValueFlow
//   state_update:^(MERKLE_UPDATE ShardState)
//   extra:^BlockExtra = Block;
type Block struct {
	tlb.SumType
	Block struct {
		GlobalId    int32
		Info        tlb.Ref[BlockInfo]
		ValueFollow tlb.Ref[ValueFlow]
		StateUpdate tlb.Ref[MerkleUpdate[ShardState]]
		Extra       tlb.Ref[BlockExtra]
	} `tlbSumType:"block#11ef55aa"`
}

// block_extra in_msg_descr:^InMsgDescr
//   out_msg_descr:^OutMsgDescr
//   account_blocks:^ShardAccountBlocks
//   rand_seed:bits256
//   created_by:bits256
//   custom:(Maybe ^McBlockExtra) = BlockExtra;
type BlockExtra struct {
	InMsgDescr    tlb.Ref[InMsgDescr]
	OutMsgDescr   tlb.Ref[OutMsgDescr]
	AccountBlocks tlb.Ref[ShardAccountBlocks]
	RandSeed      Hash
	CreatedBy     Hash
	Custom        tlb.Maybe[McBlockExtra]
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
	tlb.SumType
	McBlockExtra struct {
		KeyBlock     bool
		ShardHashes  ShardHashes
		ShardFees    ShardFees
		McExtraOther tlb.Ref[struct {
			PrevBlkSignatures tlb.HashmapE[CryptoSignaturePair] `tlb:"16bits"`
			RecoverCreate     tlb.Maybe[tlb.Ref[InMsg]]
			MintMsg           tlb.Maybe[tlb.Ref[InMsg]]
		}]
		Config ConfigParams
	} `tlbSumType:"masterchain_block_extra#cca5"`
}

func (m *McBlockExtra) UnmarshalTLB(c *boc.Cell, tag string) error {
	sumType, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	if sumType != 0xa5cc {
		return nil
	}
	m.SumType = "McBlockExtra"
	err = tlb.Unmarshal(c, &m.McBlockExtra.KeyBlock)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.McBlockExtra.ShardHashes)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.McBlockExtra.ShardFees)
	if err != nil && err != boc.ErrNotEnoughRefs {
		return err
	}
	if c != nil {
		err = m.McBlockExtra.McExtraOther.Value.PrevBlkSignatures.UnmarshalTLB(c, "16bits")
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c, &m.McBlockExtra.McExtraOther.Value.RecoverCreate)
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c, &m.McBlockExtra.McExtraOther.Value.MintMsg)
		if err != nil {
			return err
		}
	}
	if m.McBlockExtra.KeyBlock {
		err = tlb.Unmarshal(c, &m.McBlockExtra.Config)
		if err != nil {
			return err
		}
	}
	return nil
}

// sig_pair$_ node_id_short:bits256 sign:CryptoSignature = CryptoSignaturePair;  // 256+x ~ 772 bits
type CryptoSignaturePair struct {
	NodeIdShort Hash
	Sign        CryptoSignature
}

// ed25519_signature#5 R:bits256 s:bits256 = CryptoSignatureSimple;  // 516 bits
// _ CryptoSignatureSimple = CryptoSignature;
type CryptoSignatureSimple struct {
	tlb.SumType
	CryptoSignatureSimple struct {
		R Hash
		S Hash
	} `tlbSumType:"ed25519_signature#5"`
}
type CryptoSignature CryptoSignatureSimple

// _ fees:CurrencyCollection create:CurrencyCollection = ShardFeeCreated;
type ShardFeeCreated struct {
	Fees   CurrencyCollection
	Create CurrencyCollection
}

// _ (HashmapAugE 96 ShardFeeCreated ShardFeeCreated) = ShardFees;
type ShardFees struct {
	Hashmap tlb.HashmapAugE[ShardFeeCreated, ShardFeeCreated] `tlb:"96bits"`
}

// _ (HashmapAugE 256 InMsg ImportFees) = InMsgDescr;
type InMsgDescr struct {
	Hashmap tlb.HashmapAugE[InMsg, ImportFees] `tlb:"256bits"`
}

// msg_import_ext$000 msg:^(Message Any) transaction:^Transaction
//               = InMsg;
// msg_import_ihr$010 msg:^(Message Any) transaction:^Transaction
//     ihr_fee:Grams proof_created:^Cell = InMsg;
// msg_import_imm$011 in_msg:^MsgEnvelope
//     transaction:^Transaction fwd_fee:Grams = InMsg;
// msg_import_fin$100 in_msg:^MsgEnvelope
//     transaction:^Transaction fwd_fee:Grams = InMsg;
// msg_import_tr$101  in_msg:^MsgEnvelope out_msg:^MsgEnvelope
//     transit_fee:Grams = InMsg;
// msg_discard_fin$110 in_msg:^MsgEnvelope transaction_id:uint64
//     fwd_fee:Grams = InMsg;
// msg_discard_tr$111 in_msg:^MsgEnvelope transaction_id:uint64
//     fwd_fee:Grams proof_delivered:^Cell = InMsg;
type InMsg struct {
	tlb.SumType
	MsgImportExt struct {
		Msg         tlb.Ref[Message[tlb.Any]]
		Transaction tlb.Ref[Transaction]
	} `tlbSumType:"msg_import_ext$000"`
	MsgImportIhr struct {
		Msg          tlb.Ref[Message[tlb.Any]]
		Transaction  tlb.Ref[Transaction]
		IhrFee       Grams
		ProofCreated tlb.Ref[boc.Cell]
	} `tlbSumType:"msg_import_ihr$010"`
	MsgImportImm struct {
		InMsg       tlb.Ref[MsgEnvelope]
		Transaction tlb.Ref[Transaction]
		FwdFee      Grams
	} `tlbSumType:"msg_import_imm$011"`
	MsgImportFin struct {
		InMsg       tlb.Ref[MsgEnvelope]
		Transaction tlb.Ref[Transaction]
		FwdFee      Grams
	} `tlbSumType:"msg_import_fin$100"`
	MsgImportTr struct {
		InMsg      tlb.Ref[MsgEnvelope]
		OutMsg     tlb.Ref[MsgEnvelope]
		TransitFee Grams
	} `tlbSumType:"msg_import_tr$101"`
	MsgDiscardFin struct {
		InMsg         tlb.Ref[MsgEnvelope]
		TransactionId uint64
		FwdFee        Grams
	} `tlbSumType:"msg_discard_fin$110"`
	MsgDiscardTr struct {
		InMsg          tlb.Ref[MsgEnvelope]
		TransactionId  uint64
		FwdFee         Grams
		ProofDelivered tlb.Ref[boc.Cell]
	} `tlbSumType:"msg_discard_tr$111"`
}

// import_fees$_ fees_collected:Grams
//   value_imported:CurrencyCollection = ImportFees;
type ImportFees struct {
	FeesCollected Grams
	ValueImported CurrencyCollection
}

// _ (HashmapAugE 256 OutMsg CurrencyCollection) = OutMsgDescr
type OutMsgDescr struct {
	Hashmap tlb.HashmapAugE[OutMsg, ImportFees] `tlb:"256bits"`
}

// msg_export_ext$000 msg:^(Message Any)
//     transaction:^Transaction = OutMsg;
// msg_export_imm$010 out_msg:^MsgEnvelope
//     transaction:^Transaction reimport:^InMsg = OutMsg;
// msg_export_new$001 out_msg:^MsgEnvelope
//     transaction:^Transaction = OutMsg;
// msg_export_tr$011  out_msg:^MsgEnvelope
//     imported:^InMsg = OutMsg;
// msg_export_deq$1100 out_msg:^MsgEnvelope
//     import_block_lt:uint63 = OutMsg;
// msg_export_deq_short$1101 msg_env_hash:bits256
//     next_workchain:int32 next_addr_pfx:uint64
//     import_block_lt:uint64 = OutMsg;
// msg_export_tr_req$111 out_msg:^MsgEnvelope
//     imported:^InMsg = OutMsg;
// msg_export_deq_imm$100 out_msg:^MsgEnvelope
//     reimport:^InMsg = OutMsg;
type OutMsg struct {
	tlb.SumType
	MsgExportExt struct {
		Msg         tlb.Ref[Message[tlb.Any]]
		Transaction tlb.Ref[Transaction]
	} `tlbSumType:"msg_export_ext$000"`
	MsgExportImm struct {
		OutMsg      tlb.Ref[MsgEnvelope]
		Transaction tlb.Ref[Transaction]
		Reimport    tlb.Ref[InMsg]
	} `tlbSumType:"msg_export_imm$010"`
	MsgExportNew struct {
		OutMsg      tlb.Ref[MsgEnvelope]
		Transaction tlb.Ref[Transaction]
	} `tlbSumType:"msg_export_new$001"`
	MsgExportTr struct {
		OutMsg   tlb.Ref[MsgEnvelope]
		Imported tlb.Ref[InMsg]
	} `tlbSumType:"msg_export_tr$011"`
	MsgExportDeq struct {
		OutMsg      tlb.Ref[MsgEnvelope]
		ImportBlock boc.BitString `tlb:"63bits"`
	} `tlbSumType:"msg_export_deq$1100"`
	MsgExportDeqShort struct {
		MsgEnvHash     Hash
		NextWorkchain  uint32
		NextAddrPrefix uint64
		ImportBlockLt  uint64
	} `tlbSumType:"msg_export_deq_short$1101"`
	MsgExportTrReq struct {
		OutMsg   tlb.Ref[MsgEnvelope]
		Imported tlb.Ref[InMsg]
	} `tlbSumType:"msg_export_tr_req$111"`
	MsgExportDeqImm struct {
		OutMsg   tlb.Ref[MsgEnvelope]
		Reimport tlb.Ref[InMsg]
	} `tlbSumType:"msg_export_deq_imm$100"`
}

// _ (HashmapAugE 256 AccountBlock CurrencyCollection) = ShardAccountBlocks
type ShardAccountBlocks struct {
	Hashmap tlb.HashmapAugE[AccountBlock, ImportFees] `tlb:"256bits"`
}

// acc_trans#5 account_addr:bits256
//             transactions:(HashmapAug 64 ^Transaction CurrencyCollection)
//             state_update:^(HASH_UPDATE Account)
//           = AccountBlock;
type AccountBlock struct {
	tlb.SumType
	AccountBlock struct {
		AccountAddr  Hash
		Transactions tlb.HashmapAugE[tlb.Ref[Transaction], CurrencyCollection] `tlb:"64bits"`
		StateUpdate  tlb.Ref[HashUpdate]
	} `tlbSumType:"acc_trans#5"`
}

// value_flow#b8e48dfb ^[ from_prev_blk:CurrencyCollection
//   to_next_blk:CurrencyCollection
//   imported:CurrencyCollection
//   exported:CurrencyCollection ]
//   fees_collected:CurrencyCollection
//   ^[
//   fees_imported:CurrencyCollection
//   recovered:CurrencyCollection
//   created:CurrencyCollection
//   minted:CurrencyCollection
//   ] = ValueFlow;
type ValueFlow struct {
	ValueFlow struct {
		Block1 tlb.Ref[struct {
			FromPrevBlk   CurrencyCollection
			ToNextBlk     CurrencyCollection
			Imported      CurrencyCollection
			Exported      CurrencyCollection
			FeesCollected CurrencyCollection
		}]
		Block2 tlb.Ref[struct {
			FeesImported CurrencyCollection
			Recovered    CurrencyCollection
			Created      CurrencyCollection
			Minted       CurrencyCollection
		}]
	} `tlbSumType:"value_flow#b8e48dfb"`
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
		OutMsgQueueInfo tlb.Ref[OutMsgQueueInfo] // [tlb.Any] //
		BeforeSplit     bool
		Accounts        tlb.Ref[ShardAccounts]          //[tlb.Any] //
		Other           tlb.Ref[ShardStateUnsplitOther] //[tlb.Any]       //
		Custom          tlb.Maybe[tlb.Ref[McStateExtra]]
	} `tlbSumType:"shard_state#9023afe2"`
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

// master_info$_ master:ExtBlkRef = BlkMasterInfo;
type BlkMasterInfo struct {
	Master ExtBlkRef
}

// ext_blk_ref$_ end_lt:uint64
// seq_no:uint32 root_hash:bits256 file_hash:bits256
// = ExtBlkRef;
type ExtBlkRef struct {
	EndLt    uint64
	SeqNo    uint32
	RootHash Hash
	FileHash Hash
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
		Other         tlb.Ref[McStateExtraOther]
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
