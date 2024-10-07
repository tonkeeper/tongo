package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type JettonPayload struct {
	SumType string
	OpCode  *uint32
	Value   any
}

type JettonOpName = string

const (
	EmptyJettonOp   JettonOpName = ""
	UnknownJettonOp JettonOpName = "Cell"
)

// JettonOpCode is the first 4 bytes of a message body identifying an operation to be performed.
type JettonOpCode = uint32

func (j *JettonPayload) UnmarshalJSON(data []byte) error {
	var r struct {
		SumType string
		OpCode  *uint32
		Value   json.RawMessage
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	j.SumType = r.SumType
	j.OpCode = r.OpCode
	if j.SumType == EmptyJettonOp {
		return nil
	}
	if j.SumType == UnknownJettonOp {
		c := boc.NewCell()
		err := json.Unmarshal(r.Value, c)
		if err != nil {
			return err
		}
		j.Value = c
		return nil
	}
	t, ok := KnownJettonTypes[j.SumType]
	if !ok {
		return fmt.Errorf("unknown jetton payload type '%v'", j.SumType)
	}
	o := reflect.New(reflect.TypeOf(t))
	err := json.Unmarshal(r.Value, o.Interface())
	if err != nil {
		return err
	}
	j.Value = o.Elem().Interface()
	return nil
}

func (j JettonPayload) MarshalJSON() ([]byte, error) {
	if j.SumType == EmptyJettonOp {
		return []byte("{}"), nil
	}
	buf := bytes.NewBufferString(`{"SumType": "` + j.SumType + `",`)
	if j.OpCode != nil {
		fmt.Fprintf(buf, `"OpCode":%v,`, *j.OpCode)
	}
	buf.WriteString(`"Value":`)
	if j.SumType == UnknownJettonOp {
		c, ok := j.Value.(*boc.Cell)
		if !ok {
			return nil, fmt.Errorf("unknown JettonPayload should be Cell")
		}
		b, err := c.ToBoc()
		if err != nil {
			return nil, err
		}
		buf.WriteRune('"')
		hex.NewEncoder(buf).Write(b)
		buf.WriteString(`"}`)
		return buf.Bytes(), nil
	}
	if KnownJettonTypes[j.SumType] == nil {
		return nil, fmt.Errorf("unknown JettonPayload type %v", j.SumType)
	}
	b, err := json.Marshal(j.Value)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (j JettonPayload) MarshalTLB(c *boc.Cell, e *tlb.Encoder) error {
	if j.SumType == EmptyJettonOp {
		return nil
	}
	if j.OpCode != nil {
		err := c.WriteUint(uint64(*j.OpCode), 32)
		if err != nil {
			return err
		}
	} else if op, ok := JettonOpCodes[j.SumType]; ok {
		err := c.WriteUint(uint64(op), 32)
		if err != nil {
			return err
		}
	}
	return tlb.Marshal(c, j.Value)
}

func (j *JettonPayload) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	if completedRead(cell) {
		return nil
	}
	tempCell := cell.CopyRemaining()
	op64, err := tempCell.ReadUint(32)
	if errors.Is(err, boc.ErrNotEnoughBits) {
		j.SumType = UnknownJettonOp
		j.Value = cell.CopyRemaining()
		cell.ReadRemainingBits()
		return nil
	}
	op := uint32(op64)
	j.OpCode = &op
	f, ok := funcJettonDecodersMapping[JettonOpCode(op64)]

	if ok && f != nil {
		err = f(j, tempCell)
		if err == nil {
			cell.ReadRemainingBits()
			return nil
		}
	}

	j.SumType = UnknownJettonOp
	j.Value = cell.CopyRemaining()
	cell.ReadRemainingBits()

	return nil
}

func (j *JettonNotifyMsgBody) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	var prefix struct {
		QueryId uint64
		Amount  tlb.VarUInteger16
		Sender  tlb.MsgAddress
	}
	err := decoder.Unmarshal(cell, &prefix)
	if err != nil {
		return err
	}
	j.QueryId = prefix.QueryId
	j.Amount = prefix.Amount
	j.Sender = prefix.Sender
	j.ForwardPayload = failsafeJettonPayloadEitherRef(cell, decoder)
	return nil
}

func (j *JettonTransferMsgBody) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	var prefix struct {
		QueryId             uint64
		Amount              tlb.VarUInteger16
		Destination         tlb.MsgAddress
		ResponseDestination tlb.MsgAddress
	}
	err := decoder.Unmarshal(cell, &prefix)
	if err != nil {
		return err
	}
	isCustomPayload, err := cell.PickUint(1)
	if err != nil {
		return err
	}
	var (
		customPayload struct {
			CustomPayload *tlb.Any `tlb:"maybe^"`
		}
		forwardTonAmount tlb.VarUInteger16
	)
	if isCustomPayload == 1 && cell.RefsAvailableForRead() > 0 {
		err = decoder.Unmarshal(cell, &customPayload)
		if err != nil {
			return err
		}
	} else {
		_ = cell.Skip(1)
	}
	err = decoder.Unmarshal(cell, &forwardTonAmount)
	if err != nil {
		return err
	}
	j.QueryId = prefix.QueryId
	j.Amount = prefix.Amount
	j.Destination = prefix.Destination
	j.ResponseDestination = prefix.ResponseDestination
	j.CustomPayload = customPayload.CustomPayload
	j.ForwardTonAmount = forwardTonAmount
	j.ForwardPayload = failsafeJettonPayloadEitherRef(cell, decoder)
	return nil
}

func (j *JettonInternalTransferMsgBody) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	var res struct {
		QueryId          uint64
		Amount           tlb.VarUInteger16
		From             tlb.MsgAddress
		ResponseAddress  tlb.MsgAddress
		ForwardTonAmount tlb.VarUInteger16
	}
	err := decoder.Unmarshal(cell, &res)
	if err != nil {
		return err
	}
	j.QueryId = res.QueryId
	j.Amount = res.Amount
	j.From = res.From
	j.ResponseAddress = res.ResponseAddress
	j.ForwardTonAmount = res.ForwardTonAmount
	j.ForwardPayload = failsafeJettonPayloadEitherRef(cell, decoder)
	return nil
}

func failsafeJettonPayloadEitherRef(cell *boc.Cell, decoder *tlb.Decoder) tlb.EitherRef[JettonPayload] {
	isRight, err := cell.PickUint(1)
	switch {
	case err != nil:
		return tlb.EitherRef[JettonPayload]{} // empty either
	case isRight == 1 && cell.RefsAvailableForRead() < 1: // invalid either
		return tlb.EitherRef[JettonPayload]{
			IsRight: true,
		}
	default:
		var res tlb.EitherRef[JettonPayload]
		err = decoder.Unmarshal(cell, &res)
		if err != nil {
			return tlb.EitherRef[JettonPayload]{
				IsRight: isRight == 1,
			}
		}
		return res
	}
}
