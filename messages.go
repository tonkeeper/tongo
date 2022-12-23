package tongo

import (
	"fmt"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// Message
// message$_ {X:Type} info:CommonMsgInfo
// init:(Maybe (Either StateInit ^StateInit))
// body:(Either X ^X) = Message X;
type Message struct {
	Info CommonMsgInfo
	Init tlb.Maybe[tlb.EitherRef[StateInit]]
	Body tlb.EitherRef[tlb.Any]
}

// CommonMsgInfo
// int_msg_info$0 ihr_disabled:Bool bounce:Bool bounced:Bool
// src:MsgAddressInt dest:MsgAddressInt
// value:CurrencyCollection ihr_fee:Grams fwd_fee:Grams
// created_lt:uint64 created_at:uint32 = CommonMsgInfo;
// ext_in_msg_info$10 src:MsgAddressExt dest:MsgAddressInt
// import_fee:Grams = CommonMsgInfo;
// ext_out_msg_info$11 src:MsgAddressInt dest:MsgAddressExt
// created_lt:uint64 created_at:uint32 = CommonMsgInfo;
type CommonMsgInfo struct {
	tlb.SumType
	IntMsgInfo struct {
		IhrDisabled bool
		Bounce      bool
		Bounced     bool
		Src         MsgAddress
		Dest        MsgAddress
		Value       CurrencyCollection
		IhrFee      Grams
		FwdFee      Grams
		CreatedLt   uint64
		CreatedAt   uint32
	} `tlbSumType:"int_msg_info$0"`
	ExtInMsgInfo struct {
		Src       MsgAddress
		Dest      MsgAddress
		ImportFee Grams
	} `tlbSumType:"ext_in_msg_info$10"`
	ExtOutMsgInfo struct {
		Src       MsgAddress
		Dest      MsgAddress
		CreatedLt uint64
		CreatedAt uint32
	} `tlbSumType:"ext_out_msg_info$11"`
}

// StateInit
// _ split_depth:(Maybe (## 5)) special:(Maybe TickTock)
// code:(Maybe ^Cell) data:(Maybe ^Cell)
// library:(HashmapE 256 SimpleLib) = StateInit;
type StateInit struct {
	SplitDepth tlb.Maybe[tlb.Uint5]
	Special    tlb.Maybe[TickTock]
	Code       tlb.Maybe[tlb.Ref[boc.Cell]]
	Data       tlb.Maybe[tlb.Ref[boc.Cell]]
	Library    tlb.HashmapE[tlb.Size256, SimpleLib]
}

// Anycast
// anycast_info$_ depth:(#<= 30) { depth >= 1 }
// rewrite_pfx:(bits depth) = Anycast;
type Anycast struct {
	Depth      uint32
	RewritePfx boc.BitString
}

func (a Anycast) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteLimUint(int(a.Depth), 30)
	if err != nil {
		return err
	}
	err = c.WriteBitString(a.RewritePfx)
	if err != nil {
		return err
	}
	return nil
}

func (a *Anycast) UnmarshalTLB(c *boc.Cell, tag string) error {
	depth, err := c.ReadLimUint(30)
	if err != nil {
		return err
	}
	if depth < 1 {
		return fmt.Errorf("invalid anycast depth")
	}
	pfx, err := c.ReadBits(int(depth))
	if err != nil {
		return err
	}
	a.Depth = uint32(depth)
	a.RewritePfx = pfx
	return nil
}

// MsgAddress
// addr_none$00 = MsgAddressExt;
// addr_extern$01 len:(## 9) external_address:(bits len)
// = MsgAddressExt;
// addr_std$10 anycast:(Maybe Anycast)
// workchain_id:int8 address:bits256  = MsgAddressInt;
// addr_var$11 anycast:(Maybe Anycast) addr_len:(## 9)
// workchain_id:int32 address:(bits addr_len) = MsgAddressInt;
// _ _:MsgAddressInt = MsgAddress;
// _ _:MsgAddressExt = MsgAddress;
type MsgAddress struct {
	tlb.SumType
	AddrNone struct {
	} `tlbSumType:"addr_none$00"`
	AddrExtern struct {
		Len             tlb.Uint9
		ExternalAddress boc.BitString
	} `tlbSumType:"addr_extern$01"`
	AddrStd struct {
		Anycast     tlb.Maybe[Anycast]
		WorkchainId int8
		Address     Hash
	} `tlbSumType:"addr_std$10"`
	AddrVar struct {
		Anycast     tlb.Maybe[Anycast]
		AddrLen     tlb.Uint9
		WorkchainId int32
		Address     boc.BitString
	} `tlbSumType:"addr_var$11"`
}

func (a *MsgAddress) UnmarshalTLB(c *boc.Cell, tag string) error {
	t, err := c.ReadUint(2)
	if err != nil {
		return err
	}
	switch t {
	case 0:
		a.SumType = "AddrNone"
		return nil
	case 1:
		ln, err := c.ReadUint(9)
		if err != nil {
			return err
		}
		addr, err := c.ReadBits(int(ln))
		if err != nil {
			return err
		}
		a.AddrExtern.Len = tlb.Uint9(ln)
		a.AddrExtern.ExternalAddress = addr
		a.SumType = "AddrExtern"
		return nil
	case 2:
		var anycast tlb.Maybe[Anycast]
		err := anycast.UnmarshalTLB(c, "")
		if err != nil {
			return err
		}
		workchain, err := c.ReadInt(8)
		if err != nil {
			return err
		}
		address, err := c.ReadBytes(32)
		if err != nil {
			return err
		}
		a.AddrStd.Anycast = anycast
		a.AddrStd.WorkchainId = int8(workchain)
		copy(a.AddrStd.Address[:], address)
		a.SumType = "AddrStd"
		return nil
	case 3:
		var anycast tlb.Maybe[Anycast]
		err := anycast.UnmarshalTLB(c, "")
		if err != nil {
			return err
		}
		ln, err := c.ReadUint(9)
		if err != nil {
			return err
		}
		workchain, err := c.ReadInt(32)
		if err != nil {
			return err
		}
		addr, err := c.ReadBits(int(ln))
		if err != nil {
			return err
		}
		a.AddrVar.AddrLen = tlb.Uint9(ln)
		a.AddrVar.Address = addr
		a.AddrVar.WorkchainId = int32(workchain)
		a.AddrVar.Anycast = anycast
		a.SumType = "AddrVar"
	}
	return fmt.Errorf("invalid tag")
}

func (a MsgAddress) AccountID() (*AccountID, error) {
	switch a.SumType {
	case "AddrNone":
		return nil, nil
	case "AddrStd":
		return &AccountID{Workchain: int32(a.AddrStd.WorkchainId), Address: a.AddrStd.Address}, nil
	}
	return nil, fmt.Errorf("can not convert not std address to AccountId")
}

// TickTock
// tick_tock$_ tick:Bool tock:Bool = TickTock;
type TickTock struct {
	Tick bool
	Tock bool
}

// SimpleLib
// simple_lib$_ public:Bool root:^Cell = SimpleLib;
type SimpleLib struct {
	Public bool
	Root   boc.Cell `tlb:"^"`
}

func CreateExternalMessage(address AccountID, body *boc.Cell, init *StateInit, importFee Grams) (Message, error) {
	// TODO: add either selection algorithm
	var msg = Message{
		Info: CommonMsgInfo{
			SumType: "ExtInMsgInfo",
		},
		Body: tlb.EitherRef[tlb.Any]{
			IsRight: true,
			Value:   tlb.Any(*body),
		},
	}
	if init != nil {
		msg.Init.Null = false
		msg.Init.Value.IsRight = true
		msg.Init.Value.Value = *init
	} else {
		msg.Init.Null = true
	}
	msg.Info.ExtInMsgInfo.Src = MsgAddressFromAccountID(nil)
	msg.Info.ExtInMsgInfo.Dest = MsgAddressFromAccountID(&address)
	msg.Info.ExtInMsgInfo.ImportFee = importFee
	return msg, nil
}

// msg_import_ext$000 msg:^(Message Any) transaction:^Transaction
//
//	= InMsg;
//
// msg_import_ihr$010 msg:^(Message Any) transaction:^Transaction
//
//	ihr_fee:Grams proof_created:^Cell = InMsg;
//
// msg_import_imm$011 in_msg:^MsgEnvelope
//
//	transaction:^Transaction fwd_fee:Grams = InMsg;
//
// msg_import_fin$100 in_msg:^MsgEnvelope
//
//	transaction:^Transaction fwd_fee:Grams = InMsg;
//
// msg_import_tr$101  in_msg:^MsgEnvelope out_msg:^MsgEnvelope
//
//	transit_fee:Grams = InMsg;
//
// msg_discard_fin$110 in_msg:^MsgEnvelope transaction_id:uint64
//
//	fwd_fee:Grams = InMsg;
//
// msg_discard_tr$111 in_msg:^MsgEnvelope transaction_id:uint64
//
//	fwd_fee:Grams proof_delivered:^Cell = InMsg;
type InMsg struct {
	tlb.SumType
	MsgImportExt struct {
		Msg         Message     `tlb:"^"`
		Transaction Transaction `tlb:"^"`
	} `tlbSumType:"msg_import_ext$000"`
	MsgImportIhr struct {
		Msg          Message     `tlb:"^"`
		Transaction  Transaction `tlb:"^"`
		IhrFee       Grams
		ProofCreated boc.Cell `tlb:"^"`
	} `tlbSumType:"msg_import_ihr$010"`
	MsgImportImm struct {
		InMsg       MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
		FwdFee      Grams
	} `tlbSumType:"msg_import_imm$011"`
	MsgImportFin struct {
		InMsg       MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
		FwdFee      Grams
	} `tlbSumType:"msg_import_fin$100"`
	MsgImportTr struct {
		InMsg      MsgEnvelope `tlb:"^"`
		OutMsg     MsgEnvelope `tlb:"^"`
		TransitFee Grams
	} `tlbSumType:"msg_import_tr$101"`
	MsgDiscardFin struct {
		InMsg         MsgEnvelope `tlb:"^"`
		TransactionId uint64
		FwdFee        Grams
	} `tlbSumType:"msg_discard_fin$110"`
	MsgDiscardTr struct {
		InMsg          MsgEnvelope `tlb:"^"`
		TransactionId  uint64
		FwdFee         Grams
		ProofDelivered boc.Cell `tlb:"^"`
	} `tlbSumType:"msg_discard_tr$111"`
}

// import_fees$_ fees_collected:Grams
//
//	value_imported:CurrencyCollection = ImportFees;
type ImportFees struct {
	FeesCollected Grams
	ValueImported CurrencyCollection
}



// msg_export_ext$000 msg:^(Message Any)
//
//	transaction:^Transaction = OutMsg;
//
// msg_export_imm$010 out_msg:^MsgEnvelope
//
//	transaction:^Transaction reimport:^InMsg = OutMsg;
//
// msg_export_new$001 out_msg:^MsgEnvelope
//
//	transaction:^Transaction = OutMsg;
//
// msg_export_tr$011  out_msg:^MsgEnvelope
//
//	imported:^InMsg = OutMsg;
//
// msg_export_deq$1100 out_msg:^MsgEnvelope
//
//	import_block_lt:uint63 = OutMsg;
//
// msg_export_deq_short$1101 msg_env_hash:bits256
//
//	next_workchain:int32 next_addr_pfx:uint64
//	import_block_lt:uint64 = OutMsg;
//
// msg_export_tr_req$111 out_msg:^MsgEnvelope
//
//	imported:^InMsg = OutMsg;
//
// msg_export_deq_imm$100 out_msg:^MsgEnvelope
//
//	reimport:^InMsg = OutMsg;
type OutMsg struct {
	tlb.SumType
	MsgExportExt struct {
		Msg         Message     `tlb:"^"`
		Transaction Transaction `tlb:"^"`
	} `tlbSumType:"msg_export_ext$000"`
	MsgExportImm struct {
		OutMsg      MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
		Reimport    InMsg       `tlb:"^"`
	} `tlbSumType:"msg_export_imm$010"`
	MsgExportNew struct {
		OutMsg      MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
	} `tlbSumType:"msg_export_new$001"`
	MsgExportTr struct {
		OutMsg   MsgEnvelope `tlb:"^"`
		Imported InMsg       `tlb:"^"`
	} `tlbSumType:"msg_export_tr$011"`
	MsgExportDeq struct {
		OutMsg      MsgEnvelope `tlb:"^"`
		ImportBlock tlb.Uint63
	} `tlbSumType:"msg_export_deq$1100"`
	MsgExportDeqShort struct {
		MsgEnvHash     Hash
		NextWorkchain  uint32
		NextAddrPrefix uint64
		ImportBlockLt  uint64
	} `tlbSumType:"msg_export_deq_short$1101"`
	MsgExportTrReq struct {
		OutMsg   MsgEnvelope `tlb:"^"`
		Imported InMsg       `tlb:"^"`
	} `tlbSumType:"msg_export_tr_req$111"`
	MsgExportDeqImm struct {
		OutMsg   MsgEnvelope `tlb:"^"`
		Reimport InMsg       `tlb:"^"`
	} `tlbSumType:"msg_export_deq_imm$100"`
}

// _ out_queue:OutMsgQueue proc_info:ProcessedInfo
// ihr_pending:IhrPendingInfo = OutMsgQueueInfo;
type OutMsgQueueInfo struct {
	OutQueue  tlb.HashmapAugE[tlb.Size352, EnqueuedMsg, uint64]
	ProcInfo  tlb.HashmapE[tlb.Size96, ProcessedUpto]
	IhrPendig tlb.HashmapE[tlb.Size320, IhrPendingSince]
}



// _ enqueued_lt:uint64 out_msg:^MsgEnvelope = EnqueuedMsg;
type EnqueuedMsg struct {
	EnqueuedLt uint64
	OutMsg     MsgEnvelope `tlb:"^"`
}

//		msg_envelope#4 cur_addr:IntermediateAddress
//	 next_addr:IntermediateAddress fwd_fee_remaining:Grams
//	 msg:^(Message Any) = MsgEnvelope;
type MsgEnvelope struct {
	Magic           tlb.Magic `tlb:"msg_envelope#4"`
	CurrentAddress  IntermediateAddress
	NextAddress     IntermediateAddress
	FwdFeeRemaining Grams
	Msg             Message `tlb:"^"`
}

// interm_addr_regular$0 use_dest_bits:(#<= 96) = IntermediateAddress;
// interm_addr_simple$10 workchain_id:int8 addr_pfx:uint64 = IntermediateAddress;
// interm_addr_ext$11 workchain_id:int32 addr_pfx:uint64 = IntermediateAddress;
type IntermediateAddress struct {
	tlb.SumType
	IntermediateAddressRegular struct {
		UseDestBits tlb.Uint7
	} `tlbSumType:"interm_addr_regular$0"`
	IntermediateAddressSimple struct {
		WorkchainId   int8
		AddressPrefix uint64
	} `tlbSumType:"interm_addr_simple$10"`
	IntermediateAddressExt struct {
		WorkchainId   int32
		AddressPrefix uint64
	} `tlbSumType:"interm_addr_ext$11"`
}
