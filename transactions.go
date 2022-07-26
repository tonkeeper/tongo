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
type Message[T any] struct {
	Info CommonMsgInfo
	Init tlb.Maybe[tlb.Either[StateInit, tlb.Ref[StateInit]]]
	Body tlb.Either[T, tlb.Ref[T]]
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
	SplitDepth tlb.Maybe[struct {
		Depth uint64 `tlb:"5bits"`
	}]
	Special tlb.Maybe[TickTock]
	Code    tlb.Maybe[tlb.Ref[boc.Cell]]
	Data    tlb.Maybe[tlb.Ref[boc.Cell]]
	Library tlb.HashmapE[SimpleLib] `tlb:"256bits"`
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
		Len             uint32 `tlb:"9bits"`
		ExternalAddress boc.BitString
	} `tlbSumType:"addr_extern$01"`
	AddrStd struct {
		Anycast     tlb.Maybe[Anycast]
		WorkchainId int32 `tlb:"8bits"`
		Address     Hash
	} `tlbSumType:"addr_std$10"`
	AddrVar struct {
		Anycast     tlb.Maybe[Anycast]
		AddrLen     uint32 `tlb:"9bits"`
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
		a.AddrExtern.Len = uint32(ln)
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
		a.AddrStd.WorkchainId = int32(workchain)
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
		a.AddrVar.AddrLen = uint32(ln)
		a.AddrVar.Address = addr
		a.AddrVar.WorkchainId = int32(workchain)
		a.AddrVar.Anycast = anycast
		a.SumType = "AddrVar"
	}
	return fmt.Errorf("invalid tag")
}

func (a MsgAddress) AccountId() (*AccountID, error) {
	switch a.SumType {
	case "AddrNone":
		return nil, nil
	case "AddrStd":
		return &AccountID{Workchain: a.AddrStd.WorkchainId, Address: a.AddrStd.Address}, nil
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
	Root   tlb.Ref[boc.Cell]
}

// Transaction
// transaction$0111 account_addr:bits256 lt:uint64
// prev_trans_hash:bits256 prev_trans_lt:uint64 now:uint32
// outmsg_cnt:uint15
// orig_status:AccountStatus end_status:AccountStatus
// ^[ in_msg:(Maybe ^(Message Any)) out_msgs:(HashmapE 15 ^(Message Any)) ]
// total_fees:CurrencyCollection state_update:^(HASH_UPDATE Account)
// description:^TransactionDescr = Transaction;
type Transaction struct {
	tlb.SumType
	Transaction struct {
		AccountAddr   Hash
		Lt            uint64
		PrevTransHash Hash
		PrevTransLt   uint64
		Now           uint32
		OutmsgCnt     uint32 `tlb:"15bits"`
		OrigStatus    AccountStatus
		EndStatus     AccountStatus
		Msgs          tlb.Ref[struct {
			InMsg   tlb.Maybe[tlb.Ref[Message[tlb.Any]]]
			OutMsgs tlb.HashmapE[tlb.Ref[Message[tlb.Any]]] `tlb:"15bits"`
		}]
		TotalFees   CurrencyCollection
		StateUpdate tlb.Ref[HashUpdate]
		Description tlb.Ref[TransactionDescr]
	} `tlbSumType:"transaction$0111"`
}

// TransactionDescr
// trans_ord$0000 credit_first:Bool
// storage_ph:(Maybe TrStoragePhase)
// credit_ph:(Maybe TrCreditPhase)
// compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
// aborted:Bool bounce:(Maybe TrBouncePhase)
// destroyed:Bool
// = TransactionDescr;
// trans_storage$0001 storage_ph:TrStoragePhase
// = TransactionDescr;
// trans_tick_tock$001 is_tock:Bool storage_ph:TrStoragePhase
// compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
// aborted:Bool destroyed:Bool = TransactionDescr;
type TransactionDescr struct {
	tlb.SumType
	TransOrd struct {
		CreditFirst bool
		StoragePh   tlb.Maybe[TrStoragePhase]
		CreditPh    tlb.Maybe[TrCreditPhase]
		ComputePh   TrComputePhase
		Action      tlb.Maybe[tlb.Ref[TrActionPhase]]
		Aborted     bool
		Bounce      tlb.Maybe[TrBouncePhase]
		Destroyed   bool
	} `tlbSumType:"trans_ord$0000"`
	TransStorage struct {
		StoragePh TrStoragePhase
	} `tlbSumType:"trans_storage$0001"`
	TransTickTock struct {
		IsTock    bool
		StoragePh TrStoragePhase
		ComputePh TrComputePhase
		Action    tlb.Maybe[tlb.Ref[TrActionPhase]]
		Aborted   bool
		Destroyed bool
	} `tlbSumType:"trans_tick_tock$001"`
}

// TrStoragePhase
// tr_phase_storage$_ storage_fees_collected:Grams
// storage_fees_due:(Maybe Grams)
// status_change:AccStatusChange
// = TrStoragePhase;
type TrStoragePhase struct {
	StorageFeesCollected Grams
	StorageFeesDue       tlb.Maybe[Grams]
	StatusChange         AccStatusChange
}

// AccStatusChange
// acst_unchanged$0 = AccStatusChange;  // x -> x
// acst_frozen$10 = AccStatusChange;    // init -> frozen
// acst_deleted$11 = AccStatusChange;   // frozen -> deleted
type AccStatusChange string

const (
	AccStatusChangeUnchanged AccStatusChange = "acst_unchanged"
	AccStatusChangeFrozen    AccStatusChange = "acst_frozen"
	AccStatusChangeDeleted   AccStatusChange = "acst_deleted"
)

func (a AccStatusChange) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("AccStatusChange marshaling not implemented")
}

func (a *AccStatusChange) UnmarshalTLB(c *boc.Cell, tag string) error {
	f, err := c.ReadBit()
	if err != nil {
		return err
	}
	if f {
		s, err := c.ReadBit()
		if err != nil {
			return err
		}
		if s {
			*a = AccStatusChangeDeleted
			return nil
		}
		*a = AccStatusChangeFrozen
		return nil
	}
	*a = AccStatusChangeUnchanged
	return nil
}

// TrCreditPhase
// tr_phase_credit$_ due_fees_collected:(Maybe Grams)
// credit:CurrencyCollection = TrCreditPhase;
type TrCreditPhase struct {
	DueFeesCollected tlb.Maybe[Grams]
	Credit           CurrencyCollection
}

// TrComputePhase
// tr_phase_compute_skipped$0 reason:ComputeSkipReason
// = TrComputePhase;
// tr_phase_compute_vm$1 success:Bool msg_state_used:Bool
// account_activated:Bool gas_fees:Grams
// ^[ gas_used:(VarUInteger 7)
// gas_limit:(VarUInteger 7) gas_credit:(Maybe (VarUInteger 3))
// mode:int8 exit_code:int32 exit_arg:(Maybe int32)
// vm_steps:uint32
// vm_init_state_hash:bits256 vm_final_state_hash:bits256 ]
// = TrComputePhase;
type TrComputePhase struct {
	tlb.SumType
	TrPhaseComputeSkipped struct {
		Reason ComputeSkipReason
	} `tlbSumType:"tr_phase_compute_skipped$0"`
	TrPhaseComputeVm struct {
		Success          bool
		MsgStateUsed     bool
		AccountActivated bool
		GasFees          Grams
		Vm               tlb.Ref[struct {
			GasUsed   tlb.VarUInteger `tlb:"7bytes"`
			GasLimit  tlb.VarUInteger `tlb:"7bytes"`
			GasCredit tlb.Maybe[struct {
				Val tlb.VarUInteger `tlb:"3bytes"`
			}]
			Mode             int32 `tlb:"8bits"`
			ExitCode         int32
			ExitArg          tlb.Maybe[int32]
			VmSteps          uint32
			VmInitStateHash  Hash
			VmFinalStateHash Hash
		}]
	} `tlbSumType:"tr_phase_compute_vm$1"`
}

// ComputeSkipReason
// cskip_no_state$00 = ComputeSkipReason;
// cskip_bad_state$01 = ComputeSkipReason;
// cskip_no_gas$10 = ComputeSkipReason;
type ComputeSkipReason string

const (
	ComputeSkipReasonNoState  ComputeSkipReason = "cskip_no_state"
	ComputeSkipReasonBadState ComputeSkipReason = "cskip_bad_state"
	ComputeSkipReasonNoGas    ComputeSkipReason = "cskip_no_gas"
)

func (a ComputeSkipReason) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("ComputeSkipReason marshaling not implemented")
}

func (a *ComputeSkipReason) UnmarshalTLB(c *boc.Cell, tag string) error {
	t, err := c.ReadUint(2)
	if err != nil {
		return err
	}
	switch t {
	case 0:
		*a = ComputeSkipReasonNoState
	case 1:
		*a = ComputeSkipReasonBadState
	case 2:
		*a = ComputeSkipReasonNoGas
	}
	return nil
}

// TrActionPhase
// tr_phase_action$_ success:Bool valid:Bool no_funds:Bool
// status_change:AccStatusChange
// total_fwd_fees:(Maybe Grams) total_action_fees:(Maybe Grams)
// result_code:int32 result_arg:(Maybe int32) tot_actions:uint16
// spec_actions:uint16 skipped_actions:uint16 msgs_created:uint16
// action_list_hash:bits256 tot_msg_size:StorageUsedShort
// = TrActionPhase;
type TrActionPhase struct {
	Success         bool
	Valid           bool
	NoFunds         bool
	StatusChange    AccStatusChange
	TotalFwdFees    tlb.Maybe[Grams]
	TotalActionFees tlb.Maybe[Grams]
	ResultCode      int32
	ResultArg       tlb.Maybe[int32]
	TotActions      uint32 `tlb:"16bits"`
	SpecActions     uint32 `tlb:"16bits"`
	SkippedActions  uint32 `tlb:"16bits"`
	MsgsCreated     uint32 `tlb:"16bits"`
	ActionListHash  Hash
	TotMsgSize      StorageUsedShort
}

// StorageUsedShort
// storage_used_short$_ cells:(VarUInteger 7)
// bits:(VarUInteger 7) = StorageUsedShort;
type StorageUsedShort struct {
	Cells tlb.VarUInteger `tlb:"7bytes"`
	Bits  tlb.VarUInteger `tlb:"7bytes"`
}

// TrBouncePhase
// tr_phase_bounce_negfunds$00 = TrBouncePhase;
// tr_phase_bounce_nofunds$01 msg_size:StorageUsedShort
// req_fwd_fees:Grams = TrBouncePhase;
// tr_phase_bounce_ok$1 msg_size:StorageUsedShort
// msg_fees:Grams fwd_fees:Grams = TrBouncePhase;
type TrBouncePhase struct {
	tlb.SumType
	TrPhaseBounceNegfunds struct {
	} `tlbSumType:"tr_phase_bounce_negfunds$00"`
	TrPhaseBounceNofunds struct {
		MsgSize    StorageUsedShort
		ReqFwdFees Grams
	} `tlbSumType:"tr_phase_bounce_nofunds$01"`
	TrPhaseBounceOk struct {
		MsgSize StorageUsedShort
		MsgFees Grams
		FwdFees Grams
	} `tlbSumType:"tr_phase_bounce_ok$1"`
}

func (tx Transaction) IsSuccess() bool {
	success := true
	switch tx.Transaction.Description.Value.SumType {
	case "TransStorage":
		return true // TODO: check logic
	case "TransOrd":
		{
			if tx.Transaction.Description.Value.TransOrd.ComputePh.SumType == "TrPhaseComputeVm" {
				success = tx.Transaction.Description.Value.TransOrd.ComputePh.TrPhaseComputeVm.Success
			}
			if !tx.Transaction.Description.Value.TransOrd.Action.Null {
				success = success && tx.Transaction.Description.Value.TransOrd.Action.Value.Value.Success
			}
		}
	case "TransTickTock":
		{
			if tx.Transaction.Description.Value.TransTickTock.ComputePh.SumType == "TrPhaseComputeVm" {
				success = tx.Transaction.Description.Value.TransTickTock.ComputePh.TrPhaseComputeVm.Success
			}
			if !tx.Transaction.Description.Value.TransTickTock.Action.Null {
				success = success && tx.Transaction.Description.Value.TransTickTock.Action.Value.Value.Success
			}
		}
	}
	return success
}

func CreateExternalMessage(address AccountID, body *boc.Cell, init *StateInit, importFee Grams) (Message[tlb.Any], error) {
	// TODO: add either selection algorithm
	var msg = Message[tlb.Any]{
		Info: CommonMsgInfo{
			SumType: "ExtInMsgInfo",
		},
		Body: tlb.Either[tlb.Any, tlb.Ref[tlb.Any]]{
			IsRight: true,
			Right:   tlb.Ref[tlb.Any]{Value: tlb.Any(*body)},
		},
	}
	if init != nil {
		msg.Init.Null = false
		msg.Init.Value.IsRight = true
		msg.Init.Value.Right.Value = *init
	} else {
		msg.Init.Null = true
	}
	msg.Info.ExtInMsgInfo.Src = MsgAddressFromAccountID(nil)
	msg.Info.ExtInMsgInfo.Dest = MsgAddressFromAccountID(&address)
	msg.Info.ExtInMsgInfo.ImportFee = importFee
	return msg, nil
}
