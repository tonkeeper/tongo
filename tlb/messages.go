package tlb

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

// Message
// message$_ {X:Type} info:CommonMsgInfo
// init:(Maybe (Either StateInit ^StateInit))
// body:(Either X ^X) = Message X;
type Message struct {
	Info CommonMsgInfo
	Init Maybe[EitherRef[StateInit]]
	Body EitherRef[Any]

	hash Bits256
}

// Hash returns a hash of this Message.
func (m *Message) Hash() Bits256 {
	if m.Info.SumType != "ExtInMsgInfo" {
		return m.hash
	}
	// normalize ExtIn message
	c := boc.NewCell()
	_ = c.WriteUint(2, 2)                           // message$_ -> info:CommonMsgInfo -> ext_in_msg_info$10
	_ = c.WriteUint(0, 2)                           // message$_ -> info:CommonMsgInfo -> src:MsgAddressExt -> addr_none$00
	_ = m.Info.ExtInMsgInfo.Dest.MarshalTLB(c, nil) // message$_ -> info:CommonMsgInfo -> dest:MsgAddressInt
	_ = c.WriteUint(0, 4)                           // message$_ -> info:CommonMsgInfo -> import_fee:Grams -> 0
	_ = c.WriteBit(false)                           // message$_ -> init:(Maybe (Either StateInit ^StateInit)) -> nothing$0
	_ = c.WriteBit(true)                            // message$_ -> body:(Either X ^X) -> right$1
	body := boc.Cell(m.Body.Value)
	_ = c.AddRef(&body)
	hash, _ := c.Hash256()
	return hash
}

func (m *Message) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var (
		hash []byte
		err  error
	)
	if decoder.hasher != nil {
		hash, err = decoder.hasher.Hash(c)
	} else {
		hash, err = c.Hash()
	}
	if err != nil {
		return err
	}
	copy(m.hash[:], hash[:])
	c.ResetCounters()

	var msg struct {
		Info CommonMsgInfo
		Init Maybe[EitherRef[StateInit]]
		Body EitherRef[Any]
	}
	if err := decoder.Unmarshal(c, &msg); err != nil {
		return err
	}
	m.Info = msg.Info
	m.Init = msg.Init
	m.Body = msg.Body
	return nil
}

func (m Message) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	if err := encoder.Marshal(c, m.Info); err != nil {
		return err
	}
	if err := encoder.Marshal(c, m.Init); err != nil {
		return err
	}
	if err := encoder.Marshal(c, m.Body); err != nil {
		return err
	}
	return nil
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
	SumType
	IntMsgInfo *struct {
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
	ExtInMsgInfo *struct {
		Src       MsgAddress
		Dest      MsgAddress
		ImportFee VarUInteger16
	} `tlbSumType:"ext_in_msg_info$10"`
	ExtOutMsgInfo *struct {
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
	SplitDepth Maybe[Uint5]
	Special    Maybe[TickTock]
	Code       Maybe[Ref[boc.Cell]]
	Data       Maybe[Ref[boc.Cell]]
	Library    HashmapE[Bits256, SimpleLib]
}

// Anycast
// anycast_info$_ depth:(#<= 30) { depth >= 1 }
// rewrite_pfx:(bits depth) = Anycast;
type Anycast struct {
	Depth      uint32
	RewritePfx uint32
}

func (a Anycast) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	err := c.WriteLimUint(int(a.Depth), 30)
	if err != nil {
		return err
	}
	err = c.WriteUint(uint64(a.RewritePfx), int(a.Depth))
	if err != nil {
		return err
	}
	return nil
}

func (a *Anycast) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	depth, err := c.ReadLimUint(30)
	if err != nil {
		return err
	}
	if depth < 1 {
		return fmt.Errorf("invalid anycast depth")
	}
	pfx, err := c.ReadUint(int(depth))
	if err != nil {
		return err
	}
	a.Depth = uint32(depth)
	a.RewritePfx = uint32(pfx)
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
	SumType
	AddrNone struct {
	} `tlbSumType:"addr_none$00"`
	AddrExtern *boc.BitString `tlbSumType:"addr_extern$01"`
	AddrStd    struct {
		Anycast     Maybe[Anycast]
		WorkchainId int8
		Address     Bits256
	} `tlbSumType:"addr_std$10"`
	AddrVar *struct {
		Anycast     Maybe[Anycast]
		AddrLen     Uint9
		WorkchainId int32
		Address     boc.BitString
	} `tlbSumType:"addr_var$11"`
}

func (a *MsgAddress) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
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
		a.SumType = "AddrExtern"
		a.AddrExtern = &addr
		return nil
	case 2:
		var anycast Maybe[Anycast]
		err := anycast.UnmarshalTLB(c, decoder)
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
		a.SumType = "AddrStd"
		a.AddrStd.Anycast = anycast
		a.AddrStd.WorkchainId = int8(workchain)
		copy(a.AddrStd.Address[:], address)
		return nil
	case 3:
		var anycast Maybe[Anycast]
		err := anycast.UnmarshalTLB(c, decoder)
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
		a.SumType = "AddrVar"
		a.AddrVar = &struct {
			Anycast     Maybe[Anycast]
			AddrLen     Uint9
			WorkchainId int32
			Address     boc.BitString
		}{Anycast: anycast, AddrLen: Uint9(ln), WorkchainId: int32(workchain), Address: addr}
		return nil
	}
	return fmt.Errorf("invalid tag")
}

func (a MsgAddress) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	switch a.SumType {
	case "AddrNone":
		return c.WriteUint(0, 2)
	case "AddrExtern":
		err := c.WriteUint(1, 2)
		if err != nil {
			return err
		}
		a.AddrExtern.ResetCounter()
		l := a.AddrExtern.BitsAvailableForRead()
		if l > 511 {
			return fmt.Errorf("external address is too long")
		}
		err = c.WriteUint(uint64(l), 9)
		if err != nil {
			return err
		}
		return c.WriteBitString(*a.AddrExtern)
	case "AddrStd":
		if err := c.WriteUint(2, 2); err != nil {
			return err
		}
		if err := a.AddrStd.Anycast.MarshalTLB(c, encoder); err != nil {
			return err
		}
		if err := c.WriteInt(int64(a.AddrStd.WorkchainId), 8); err != nil {
			return err
		}
		return c.WriteBytes(a.AddrStd.Address[:])
	case "AddrVar":
		if err := c.WriteUint(3, 2); err != nil {
			return err
		}
		if err := a.AddrVar.Anycast.MarshalTLB(c, encoder); err != nil {
			return err
		}
		if err := c.WriteUint(uint64(a.AddrVar.AddrLen), 9); err != nil {
			return err
		}
		if err := c.WriteInt(int64(a.AddrVar.WorkchainId), 32); err != nil {
			return err
		}
		return c.WriteBitString(a.AddrVar.Address)
	}
	return fmt.Errorf("invalid tag")
}

func (a MsgAddress) MarshalJSON() ([]byte, error) {
	var x string
	var extra string
	switch a.SumType {
	case "AddrExtern":
		// we assume that AddrExtern.ExternalAddress has exactly AddrExtern.Len bits
		// that's always true, if the current MsgAddress was deserialized from TL-B.
		x = a.AddrExtern.ToFiftHex()
	case "AddrStd":
		if a.AddrStd.Anycast.Exists {
			extra = fmt.Sprintf(":Anycast(%d,%d)", a.AddrStd.Anycast.Value.Depth, a.AddrStd.Anycast.Value.RewritePfx)
		}
		x = fmt.Sprintf("%d:%s", a.AddrStd.WorkchainId, a.AddrStd.Address.Hex())
	case "AddrVar":
		if a.AddrVar.Anycast.Exists {
			extra = fmt.Sprintf(":Anycast(%d,%d)", a.AddrVar.Anycast.Value.Depth, a.AddrVar.Anycast.Value.RewritePfx)
		}
		// we assume that AddrVar.Address has exactly AddrVar.Len bits
		// that's always true, if the current MsgAddress was deserialized from TL-B.
		x = fmt.Sprintf("%d:%s", a.AddrVar.WorkchainId, a.AddrVar.Address.ToFiftHex())
	}
	return []byte(fmt.Sprintf(`"%s%s"`, x, extra)), nil
}

func (a *MsgAddress) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" {
		*a = MsgAddress{SumType: "AddrNone"}
		return nil
	}
	parts := strings.Split(value, ":")
	if len(parts) == 1 {
		externalAddr, err := boc.BitStringFromFiftHex(value)
		if err != nil {
			return err
		}
		*a = MsgAddress{
			SumType:    "AddrExtern",
			AddrExtern: externalAddr,
		}
		return nil
	}
	if len(parts) != 2 && len(parts) != 3 {
		return fmt.Errorf("unknown MsgAddress format")
	}

	var anycast *Anycast
	if len(parts) == 3 {
		if !strings.HasPrefix(parts[2], "Anycast(") || !strings.HasSuffix(parts[2], ")") {
			return fmt.Errorf("unknown MsgAddress format")
		}
		var depth uint32
		var prefix uint32
		depthAndPrefix := parts[2][len("Anycast(") : len(parts[2])-1]
		if _, err := fmt.Sscanf(depthAndPrefix, "%d,%d", &depth, &prefix); err != nil {
			return fmt.Errorf("failed to parse Anycast in MsgAddress: %w", err)
		}
		anycast = &Anycast{
			Depth:      depth,
			RewritePfx: prefix,
		}
	}
	// try AddrStd first
	num, err := strconv.ParseInt(parts[0], 10, 32)
	isWorkchainInt8 := err == nil && num >= int64(math.MinInt8) && num <= int64(math.MaxInt8)
	if len(parts[1]) == 64 && isWorkchainInt8 && !strings.HasSuffix(parts[1], "_") {
		var dst [32]byte
		_, err := hex.Decode(dst[:], []byte(parts[1]))
		if err != nil {
			return err
		}
		workchain, err := strconv.ParseInt(parts[0], 10, 8)
		if err != nil {
			return fmt.Errorf("failed to parse %v workchain: %w", parts[0], err)
		}
		*a = MsgAddress{
			SumType: "AddrStd",
			AddrStd: struct {
				Anycast     Maybe[Anycast]
				WorkchainId int8
				Address     Bits256
			}{
				WorkchainId: int8(workchain),
				Address:     dst,
			},
		}
		if anycast != nil {
			a.AddrStd.Anycast = Maybe[Anycast]{
				Exists: true,
				Value:  *anycast,
			}
		}
		return nil
	}
	bitstr, err := boc.BitStringFromFiftHex(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse fift hex: %w", err)
	}
	workchain, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse %v workchain: %w", parts[0], err)
	}
	*a = MsgAddress{
		SumType: "AddrVar",
		AddrVar: &struct {
			Anycast     Maybe[Anycast]
			AddrLen     Uint9
			WorkchainId int32
			Address     boc.BitString
		}{
			AddrLen:     Uint9(bitstr.BitsAvailableForRead()),
			WorkchainId: int32(workchain),
			Address:     *bitstr,
		},
	}
	if anycast != nil {
		a.AddrVar.Anycast = Maybe[Anycast]{
			Exists: true,
			Value:  *anycast,
		}
	}
	return nil
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

// msg_import_ext$000 msg:^(Message Any) transaction:^Transaction
//	= InMsg;

// msg_import_ihr$010 msg:^(Message Any) transaction:^Transaction
//	ihr_fee:Grams proof_created:^Cell = InMsg;

// msg_import_imm$011 in_msg:^MsgEnvelope
//	transaction:^Transaction fwd_fee:Grams = InMsg;

// msg_import_fin$100 in_msg:^MsgEnvelope
//	transaction:^Transaction fwd_fee:Grams = InMsg;

// msg_import_tr$101  in_msg:^MsgEnvelope out_msg:^MsgEnvelope
//	transit_fee:Grams = InMsg;

// msg_discard_fin$110 in_msg:^MsgEnvelope transaction_id:uint64
//	fwd_fee:Grams = InMsg;

// msg_discard_tr$111 in_msg:^MsgEnvelope transaction_id:uint64
//	fwd_fee:Grams proof_delivered:^Cell = InMsg;

//msg_import_deferred_fin$00100 in_msg:^MsgEnvelope
//    transaction:^Transaction fwd_fee:Grams = InMsg;

// msg_import_deferred_tr$00101 in_msg:^MsgEnvelope out_msg:^MsgEnvelope = InMsg;
type InMsg struct {
	SumType
	MsgImportExt *struct {
		Msg         Message     `tlb:"^"`
		Transaction Transaction `tlb:"^"`
	} `tlbSumType:"msg_import_ext$000"`
	MsgImportIhr *struct {
		Msg          Message     `tlb:"^"`
		Transaction  Transaction `tlb:"^"`
		IhrFee       Grams
		ProofCreated boc.Cell `tlb:"^"`
	} `tlbSumType:"msg_import_ihr$010"`
	MsgImportImm *struct {
		InMsg       MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
		FwdFee      Grams
	} `tlbSumType:"msg_import_imm$011"`
	MsgImportFin *struct {
		InMsg       MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
		FwdFee      Grams
	} `tlbSumType:"msg_import_fin$100"`
	MsgImportTr *struct {
		InMsg      MsgEnvelope `tlb:"^"`
		OutMsg     MsgEnvelope `tlb:"^"`
		TransitFee Grams
	} `tlbSumType:"msg_import_tr$101"`
	MsgDiscardFin *struct {
		InMsg         MsgEnvelope `tlb:"^"`
		TransactionId uint64
		FwdFee        Grams
	} `tlbSumType:"msg_discard_fin$110"`
	MsgDiscardTr *struct {
		InMsg          MsgEnvelope `tlb:"^"`
		TransactionId  uint64
		FwdFee         Grams
		ProofDelivered boc.Cell `tlb:"^"`
	} `tlbSumType:"msg_discard_tr$111"`
	MsgImportDeferredFin *struct {
		InMsg         MsgEnvelope `tlb:"^"`
		TransactionId Transaction `tlb:"^"`
		FwdFee        Grams
	} `tlbSumType:"msg_import_deferred_fin$00100"`
	MsgImportDeferredTr *struct {
		InMsg  MsgEnvelope `tlb:"^"`
		OutMsg MsgEnvelope `tlb:"^"`
	} `tlbSumType:"msg_import_deferred_tr$00101"`
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
//
// msg_export_new_defer$10100 out_msg:^MsgEnvelope
//
//	transaction:^Transaction = OutMsg;
//
// msg_export_deferred_tr$10101  out_msg:^MsgEnvelope
//
//	imported:^InMsg = OutMsg;
type OutMsg struct {
	SumType
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
		ImportBlock Uint63
	} `tlbSumType:"msg_export_deq$1100"`
	MsgExportDeqShort struct {
		MsgEnvHash     Bits256
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
	MsgExportNewDefer *struct {
		OutMsg      MsgEnvelope `tlb:"^"`
		Transaction Transaction `tlb:"^"`
	} `tlbSumType:"msg_export_new_defer$10100"`
	MsgExportDeferredTr *struct {
		OutMsg   MsgEnvelope `tlb:"^"`
		Imported InMsg       `tlb:"^"`
	} `tlbSumType:"msg_export_deferred_tr$10101"`
}

// dispatch_queue:DispatchQueue out_queue_size:(Maybe uint48) = OutMsgQueueExtra;
type OutMsgQueueExtra struct {
	Magic Magic `tlb:"out_msg_queue_extra#0"`
	// key - sender address, aug - min created_lt
	DispatchQueue HashmapAugE[Bits256, AccountDispatchQueue, uint64]
	OutQueueSize  Maybe[Uint48]
}

// _ messages:(HashmapE 64 EnqueuedMsg) count:uint48 = AccountDispatchQueue;
type AccountDispatchQueue struct {
	Messages HashmapE[Uint64, EnqueuedMsg]
	Count    Uint48
}

// _ out_queue:OutMsgQueue proc_info:ProcessedInfo
// ihr_pending:IhrPendingInfo = OutMsgQueueInfo;
type OutMsgQueueInfo struct {
	OutQueue HashmapAugE[Bits352, EnqueuedMsg, uint64]
	ProcInfo HashmapE[Bits96, ProcessedUpto]
	Extra    Maybe[OutMsgQueueExtra]
}

// _ enqueued_lt:uint64 out_msg:^MsgEnvelope = EnqueuedMsg;
type EnqueuedMsg struct {
	EnqueuedLt uint64
	OutMsg     MsgEnvelope `tlb:"^"`
}

//		msg_envelope#4 cur_addr:IntermediateAddress
//	 next_addr:IntermediateAddress fwd_fee_remaining:Grams
//	 msg:^(Message Any) = MsgEnvelope;
//
// msg_envelope_v2#5 cur_addr:IntermediateAddress
//
//	next_addr:IntermediateAddress fwd_fee_remaining:Grams
//	msg:^(Message Any)
//	emitted_lt:(Maybe uint64)
//	metadata:(Maybe MsgMetadata) = MsgEnvelope;
type MsgEnvelope struct {
	SumType SumType
	V1      struct {
		CurrentAddress  IntermediateAddress
		NextAddress     IntermediateAddress
		FwdFeeRemaining Grams
		Msg             Message `tlb:"^"`
	} `tlbSumType:"msg_envelope#4"`
	V2 struct {
		CurrentAddress  IntermediateAddress
		NextAddress     IntermediateAddress
		FwdFeeRemaining Grams
		Msg             Message      `tlb:"^"`
		EmittedLT       *uint64      `tlb:"maybe"`
		Metadata        *MsgMetadata `tlb:"maybe"`
	} `tlbSumType:"msg_envelope_v2#5"`
}

// msg_metadata#0 depth:uint32 initiator_addr:MsgAddressInt initiator_lt:uint64 = MsgMetadata;
type MsgMetadata struct {
	Magic         Magic `tlb:"msg_metadata#0"`
	Depth         uint32
	InitiatorAddr MsgAddress
	InitiatorLT   uint64
}

// interm_addr_regular$0 use_dest_bits:(#<= 96) = IntermediateAddress;
// interm_addr_simple$10 workchain_id:int8 addr_pfx:uint64 = IntermediateAddress;
// interm_addr_ext$11 workchain_id:int32 addr_pfx:uint64 = IntermediateAddress;
type IntermediateAddress struct {
	SumType
	IntermediateAddressRegular struct {
		UseDestBits Uint7
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
