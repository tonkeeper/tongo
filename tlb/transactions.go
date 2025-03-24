package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

// Transaction
// transaction$0111 account_addr:bits256 lt:uint64
// prev_trans_hash:bits256 prev_trans_lt:uint64 now:uint32
// outmsg_cnt:uint15
// orig_status:AccountStatus end_status:AccountStatus
// ^[ in_msg:(Maybe ^(Message Any)) out_msgs:(HashmapE 15 ^(Message Any)) ]
// total_fees:CurrencyCollection state_update:^(HASH_UPDATE Account)
// description:^TransactionDescr = Transaction;
type Transaction struct {
	Magic         Magic `tlb:"transaction$0111"`
	AccountAddr   Bits256
	Lt            uint64
	PrevTransHash Bits256
	PrevTransLt   uint64
	Now           uint32
	OutMsgCnt     Uint15
	OrigStatus    AccountStatus
	EndStatus     AccountStatus
	Msgs          struct {
		InMsg   Maybe[Ref[Message]]
		OutMsgs HashmapE[Uint15, Ref[Message]]
	} `tlb:"^"`
	TotalFees   CurrencyCollection
	StateUpdate HashUpdate       `tlb:"^"`
	Description TransactionDescr `tlb:"^"`

	hash Bits256

	lazySourceBoc func() ([]byte, error)
}

// Hash returns a hash of this transaction.
func (tx *Transaction) Hash() Bits256 {
	return tx.hash
}

// SourceBoc returns a BOC of this transaction.
// It works only if the transaction was unmarshalled from a cell.
func (tx *Transaction) SourceBoc() ([]byte, error) {
	if tx.lazySourceBoc != nil {
		return tx.lazySourceBoc()
	}
	return nil, fmt.Errorf("transaction was not unmarshalled from cell")
}

func (tx *Transaction) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var (
		hash []byte
		err  error
	)
	if decoder.hasher != nil {
		tx.lazySourceBoc = func() ([]byte, error) {
			c.ResetCounters()
			return c.ToBocCustomWithHasher(decoder.Hasher(), false, false, false, 0)
		}
		hash, err = decoder.hasher.Hash(c)
	} else {
		tx.lazySourceBoc = func() ([]byte, error) {
			c.ResetCounters()
			return boc.SerializeBoc(c, false, false, false, 0)
		}
		hash, err = c.Hash()
	}
	if err != nil {
		return err
	}
	copy(tx.hash[:], hash[:])
	c.ResetCounters()

	sumType, err := c.ReadUint(4)
	if err != nil {
		return err
	}
	if sumType != 0b0111 {
		return fmt.Errorf("invalid tag")
	}
	if err = decoder.Unmarshal(c, &tx.AccountAddr); err != nil {
		return err
	}
	if err = decoder.Unmarshal(c, &tx.Lt); err != nil {
		return err
	}
	if err = decoder.Unmarshal(c, &tx.PrevTransHash); err != nil {
		return err
	}
	if err = decoder.Unmarshal(c, &tx.PrevTransLt); err != nil {
		return err
	}
	if err = decoder.Unmarshal(c, &tx.Now); err != nil {
		return err
	}
	outMsgCnt, err := c.ReadUint(15)
	if err != nil {
		return err
	}
	tx.OutMsgCnt = Uint15(outMsgCnt)
	if err = tx.OrigStatus.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = tx.EndStatus.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	c1, err := c.NextRef()
	if err != nil {
		return err
	}
	var msgs struct {
		InMsg   Maybe[Ref[Message]]
		OutMsgs HashmapE[Uint15, Ref[Message]]
	}
	if err = decoder.Unmarshal(c1, &msgs); err != nil {
		return err
	}
	tx.Msgs = msgs
	if err = decoder.Unmarshal(c, &tx.TotalFees); err != nil {
		return err
	}
	c2, err := c.NextRef()
	if err != nil {
		return err
	}
	if err = decoder.Unmarshal(c2, &tx.StateUpdate); err != nil {
		return err
	}
	c3, err := c.NextRef()
	if err != nil {
		return err
	}
	if err = decoder.Unmarshal(c3, &tx.Description); err != nil {
		return err
	}
	return nil
}

// trans_ord$0000 credit_first:Bool
//   storage_ph:(Maybe TrStoragePhase)
//   credit_ph:(Maybe TrCreditPhase)
//   compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
//   aborted:Bool bounce:(Maybe TrBouncePhase)
//   destroyed:Bool
//   = TransactionDescr;

// trans_storage$0001 storage_ph:TrStoragePhase
//   = TransactionDescr;

// trans_tick_tock$001 is_tock:Bool storage_ph:TrStoragePhase
//   compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
//   aborted:Bool destroyed:Bool = TransactionDescr;
// //

// trans_split_prepare$0100 split_info:SplitMergeInfo
//   storage_ph:(Maybe TrStoragePhase)
//   compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
//   aborted:Bool destroyed:Bool
//   = TransactionDescr;
// trans_split_install$0101 split_info:SplitMergeInfo
//   prepare_transaction:^Transaction
//   installed:Bool = TransactionDescr;

// trans_merge_prepare$0110 split_info:SplitMergeInfo
//
//	storage_ph:TrStoragePhase aborted:Bool
//	= TransactionDescr;
//
// trans_merge_install$0111 split_info:SplitMergeInfo
//
//	prepare_transaction:^Transaction
//	storage_ph:(Maybe TrStoragePhase)
//	credit_ph:(Maybe TrCreditPhase)
//	compute_ph:TrComputePhase action:(Maybe ^TrActionPhase)
//	aborted:Bool destroyed:Bool
//	= TransactionDescr;
type TransactionDescr struct {
	SumType
	TransOrd struct {
		CreditFirst bool
		StoragePh   Maybe[TrStoragePhase]
		CreditPh    Maybe[TrCreditPhase]
		ComputePh   TrComputePhase
		Action      Maybe[Ref[TrActionPhase]]
		Aborted     bool
		Bounce      Maybe[TrBouncePhase]
		Destroyed   bool
	} `tlbSumType:"trans_ord$0000"`
	TransStorage struct {
		StoragePh TrStoragePhase
	} `tlbSumType:"trans_storage$0001"`
	TransTickTock struct {
		IsTock    bool
		StoragePh TrStoragePhase
		ComputePh TrComputePhase
		Action    Maybe[Ref[TrActionPhase]]
		Aborted   bool
		Destroyed bool
	} `tlbSumType:"trans_tick_tock$001"`
	TransSplitPrepare *struct {
		SplitInfo SplitMergeInfo
		StoragePh Maybe[TrStoragePhase]
		ComputePh TrComputePhase
		Action    Maybe[Ref[TrActionPhase]]
		Aborted   bool
		Destroyed bool
	} `tlbSumType:"trans_split_prepare$0100"`
	TransSplitInstall *struct {
		SplitInfo          SplitMergeInfo
		PrepareTransaction Any `tlb:"^"`
		Installed          bool
	} `tlbSumType:"trans_split_install$0101"`
	TransMergePrepare *struct {
		SplitInfo SplitMergeInfo
		StoragePh TrStoragePhase
		Aborted   bool
	} `tlbSumType:"trans_merge_prepare$0110"`
	TransMergeInstall *struct {
		SplitInfo          SplitMergeInfo
		PrepareTransaction Any `tlb:"^"` //Transaction]
		StoragePh          Maybe[TrStoragePhase]
		CreditPh           Maybe[TrCreditPhase]
		ComputePh          TrComputePhase
		Action             Maybe[Ref[TrActionPhase]]
		Aborted            bool
		Destroyed          bool
	} `tlbSumType:"trans_merge_install$0111"`
}

// split_merge_info$_ cur_shard_pfx_len:(## 6)
//
//	acc_split_depth:(## 6) this_addr:bits256 sibling_addr:bits256
//	= SplitMergeInfo;
type SplitMergeInfo struct {
	CurSHardPfxLen Uint6
	AccSplitDepth  Uint6
	ThisAddr       Bits256
	SiblingAddr    Bits256
}

// TrStoragePhase
// tr_phase_storage$_ storage_fees_collected:Grams
// storage_fees_due:(Maybe Grams)
// status_change:AccStatusChange
// = TrStoragePhase;
type TrStoragePhase struct {
	StorageFeesCollected Grams
	StorageFeesDue       Maybe[Grams]
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

func (a AccStatusChange) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	if a == AccStatusChangeUnchanged {
		return c.WriteBit(false)
	}
	if err := c.WriteBit(true); err != nil {
		return err
	}
	if a == AccStatusChangeDeleted {
		return c.WriteBit(true)
	}
	return c.WriteBit(false)
}

func (a *AccStatusChange) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
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
	DueFeesCollected Maybe[Grams]
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
	SumType
	TrPhaseComputeSkipped struct {
		Reason ComputeSkipReason
	} `tlbSumType:"tr_phase_compute_skipped$0"`
	TrPhaseComputeVm struct {
		Success          bool
		MsgStateUsed     bool
		AccountActivated bool
		GasFees          Grams
		Vm               struct {
			GasUsed          VarUInteger7
			GasLimit         VarUInteger7
			GasCredit        Maybe[VarUInteger3]
			Mode             int8
			ExitCode         int32
			ExitArg          Maybe[int32]
			VmSteps          uint32
			VmInitStateHash  Bits256
			VmFinalStateHash Bits256
		} `tlb:"^"`
	} `tlbSumType:"tr_phase_compute_vm$1"`
}

// ComputeSkipReason
// cskip_no_state$00 = ComputeSkipReason;
// cskip_bad_state$01 = ComputeSkipReason;
// cskip_no_gas$10 = ComputeSkipReason;
// cskip_suspended$110 = ComputeSkipReason;
type ComputeSkipReason string

const (
	ComputeSkipReasonNoState  ComputeSkipReason = "cskip_no_state"
	ComputeSkipReasonBadState ComputeSkipReason = "cskip_bad_state"
	ComputeSkipReasonNoGas    ComputeSkipReason = "cskip_no_gas"
	ComputeSkipSuspended      ComputeSkipReason = "cskip_suspended"
)

func (a ComputeSkipReason) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	switch a {
	case ComputeSkipReasonNoState:
		return c.WriteUint(0, 2)
	case ComputeSkipReasonBadState:
		return c.WriteUint(1, 2)
	case ComputeSkipReasonNoGas:
		return c.WriteUint(2, 2)
	case ComputeSkipSuspended:
		if err := c.WriteUint(3, 2); err != nil {
			return err
		}
		return c.WriteUint(0, 1)
	}
	return nil
}

func (a *ComputeSkipReason) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
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
	case 3:
		nextBit, err := c.ReadUint(1)
		if err != nil {
			return err
		}
		if nextBit == 0 {
			*a = ComputeSkipSuspended
			return nil
		}
		return fmt.Errorf("unknown ComputeSkipReason")
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
	TotalFwdFees    Maybe[Grams]
	TotalActionFees Maybe[Grams]
	ResultCode      int32
	ResultArg       Maybe[int32]
	TotActions      uint16
	SpecActions     uint16
	SkippedActions  uint16
	MsgsCreated     uint16
	ActionListHash  Bits256
	TotMsgSize      StorageUsed
}

// StorageUsed
// storage_used$_ cells:(VarUInteger 7)
// bits:(VarUInteger 7) = StorageUsed;
type StorageUsed struct {
	Cells VarUInteger7
	Bits  VarUInteger7
}

// TrBouncePhase
// tr_phase_bounce_negfunds$00 = TrBouncePhase;
// tr_phase_bounce_nofunds$01 msg_size:StorageUsedShort
// req_fwd_fees:Grams = TrBouncePhase;
// tr_phase_bounce_ok$1 msg_size:StorageUsedShort
// msg_fees:Grams fwd_fees:Grams = TrBouncePhase;
type TrBouncePhase struct {
	SumType
	TrPhaseBounceNegfunds struct {
	} `tlbSumType:"tr_phase_bounce_negfunds$00"`
	TrPhaseBounceNofunds struct {
		MsgSize    StorageUsed
		ReqFwdFees Grams
	} `tlbSumType:"tr_phase_bounce_nofunds$01"`
	TrPhaseBounceOk struct {
		MsgSize StorageUsed
		MsgFees Grams
		FwdFees Grams
	} `tlbSumType:"tr_phase_bounce_ok$1"`
}

func (tx Transaction) IsSuccess() bool {
	switch tx.Description.SumType {
	case "TransOrd":
		o := tx.Description.TransOrd
		if o.Bounce.Exists {
			return false
		}
		cph := o.ComputePh
		if cph.SumType == "TrPhaseComputeSkipped" && cph.TrPhaseComputeSkipped.Reason != ComputeSkipReasonNoState {
			return false
		}
		if cph.SumType == "TrPhaseComputeVm" && (!cph.TrPhaseComputeVm.Success || (cph.TrPhaseComputeVm.Vm.ExitCode != 0 && cph.TrPhaseComputeVm.Vm.ExitCode != 1)) {
			return false
		}
		if o.Action.Exists && !o.Action.Value.Value.Success {
			return false
		}
		return true
	case "TransTickTock":
		t := tx.Description.TransTickTock
		cph := t.ComputePh
		if cph.SumType == "TrPhaseComputeSkipped" && cph.TrPhaseComputeSkipped.Reason != ComputeSkipReasonNoState {
			return false
		}
		if cph.SumType == "TrPhaseComputeVm" && (!cph.TrPhaseComputeVm.Success || (cph.TrPhaseComputeVm.Vm.ExitCode != 0 && cph.TrPhaseComputeVm.Vm.ExitCode != 1)) {
			return false
		}
		if t.Action.Exists && !t.Action.Value.Value.Success {
			return false
		}
		return true
	default:
		return true //todo: add logic for over types
	}
}
