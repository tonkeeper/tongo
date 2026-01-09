package abitolk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func pointer[T any](t T) *T {
	return &t
}

func decodeMsg(tag tlb.Tag, name MsgOpName, bodyType any) msgDecoderFunc {
	return func(cell *boc.Cell) (*uint32, *MsgOpName, any, error) {
		cell.ResetCounters()
		readTag, err := cell.ReadUint(tag.Len)
		if err != nil {
			return nil, nil, nil, err
		}
		var uintTag *uint32
		if readTag != tag.Val {
			return nil, nil, nil, fmt.Errorf("invalid tag")
		}
		if tag.Len == 32 {
			uintTag = pointer(uint32(readTag))
		}
		body := reflect.New(reflect.TypeOf(bodyType))
		err = tlb.Unmarshal(cell, body.Interface())
		if err == nil {
			return uintTag, pointer(name), body.Elem().Interface(), nil
		}
		return uintTag, nil, nil, err
	}
}

func decodeMultipleMsgs(funcs []msgDecoderFunc, tag string) msgDecoderFunc {
	return func(cell *boc.Cell) (*uint32, *MsgOpName, any, error) {
		for _, f := range funcs {
			tag, opName, object, err := f(cell)
			if err == nil && completedRead(cell) {
				return tag, opName, object, err
			}
		}
		return nil, nil, nil, fmt.Errorf("no one message can be unmarshled for %v", tag)
	}
}

type InMsgBody struct {
	SumType MsgOpName
	OpCode  *uint32
	Value   any
}

func (body InMsgBody) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if body.SumType == EmptyMsgOp {
		return nil
	}
	if body.OpCode != nil {
		err := c.WriteUint(uint64(*body.OpCode), 32)
		if err != nil {
			return err
		}
	}
	return tlb.Marshal(c, body.Value)
}

func (body *InMsgBody) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	body.SumType = UnknownMsgOp
	tag, name, value, err := InternalMessageDecoder(cell, nil)
	if err != nil {
		return err
	}
	if tag == nil && name == nil {
		body.SumType = EmptyMsgOp
		return nil
	}
	body.OpCode = tag
	if name != nil {
		body.SumType = *name
		body.Value = value
	} else {
		if cell.BitsAvailableForRead() != cell.BitSize() {
			cell = cell.CopyRemaining() //because body can be part of the message cell
		}
		body.Value = cell
	}
	return nil
}

func (body *InMsgBody) UnmarshalJSON(data []byte) error {
	var r struct {
		SumType string
		OpCode  *uint32
		Value   json.RawMessage
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	body.SumType = r.SumType
	body.OpCode = r.OpCode
	if body.SumType == EmptyMsgOp {
		return nil
	}
	if body.SumType == UnknownMsgOp {
		c := boc.NewCell()
		err := json.Unmarshal(r.Value, c)
		if err != nil {
			return err
		}
		body.Value = c
		return nil
	}
	t, ok := KnownMsgInTypes[body.SumType]
	if !ok {
		return fmt.Errorf("unknown message body type '%v'", body.SumType)
	}
	o := reflect.New(reflect.TypeOf(t))
	err := json.Unmarshal(r.Value, o.Interface())
	if err != nil {
		return err
	}
	body.Value = o.Elem().Interface()
	return nil
}

func (body InMsgBody) MarshalJSON() ([]byte, error) {
	if body.SumType == EmptyMsgOp {
		return []byte("{}"), nil
	}
	buf := bytes.NewBufferString(`{"SumType": "` + body.SumType + `",`)
	if body.OpCode != nil {
		fmt.Fprintf(buf, `"OpCode":%v,`, *body.OpCode)
	}
	buf.WriteString(`"Value":`)
	if body.SumType == UnknownMsgOp {
		c, ok := body.Value.(*boc.Cell)
		if !ok {
			return nil, fmt.Errorf("unknown MsgBody should be Cell")
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
	if KnownMsgInTypes[body.SumType] == nil {
		return nil, fmt.Errorf("unknown MsgBody type %v", body.SumType)
	}
	b, err := json.Marshal(body.Value)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type ExtOutMsgBody struct {
	SumType string
	OpCode  *uint32
	Value   any
}

func (b *ExtOutMsgBody) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	b.SumType = UnknownMsgOp
	tag, name, value, err := ExtOutMessageDecoder(cell, nil, tlb.MsgAddress{SumType: "AddrNone"})
	if err != nil {
		return err
	}
	if tag == nil && name == nil {
		b.SumType = EmptyMsgOp
	}
	b.OpCode = tag
	if name != nil {
		b.SumType = *name
		b.Value = value
	} else {
		if cell.BitsAvailableForRead() != cell.BitSize() {
			cell = cell.CopyRemaining() //because body can be part of the message cell
		}
		b.Value = cell
	}
	return nil
}

func (body *ExtOutMsgBody) UnmarshalJSON(data []byte) error {
	var r struct {
		SumType string
		OpCode  *uint32
		Value   json.RawMessage
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	body.SumType = r.SumType
	body.OpCode = r.OpCode
	if body.SumType == EmptyMsgOp {
		return nil
	}
	if body.SumType == UnknownMsgOp {
		c := boc.NewCell()
		err := json.Unmarshal(r.Value, c)
		if err != nil {
			return err
		}
		body.Value = c
		return nil
	}
	t, ok := KnownMsgExtOutTypes[body.SumType]
	if !ok {
		return fmt.Errorf("unknown message body type '%v'", body.SumType)
	}
	o := reflect.New(reflect.TypeOf(t))
	err := json.Unmarshal(r.Value, o.Interface())
	if err != nil {
		return err
	}
	body.Value = o.Elem().Interface()
	return nil
}

func (body ExtOutMsgBody) MarshalJSON() ([]byte, error) {
	if body.SumType == EmptyMsgOp {
		return []byte("{}"), nil
	}
	buf := bytes.NewBufferString(`{"SumType": "` + body.SumType + `",`)
	if body.OpCode != nil {
		fmt.Fprintf(buf, `"OpCode":%v,`, *body.OpCode)
	}
	buf.WriteString(`"Value":`)
	if body.SumType == UnknownMsgOp {
		c, ok := body.Value.(*boc.Cell)
		if !ok {
			return nil, fmt.Errorf("unknown MsgBody should be Cell")
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
	if KnownMsgExtOutTypes[body.SumType] == nil {
		return nil, fmt.Errorf("unknown MsgBody type %v", body.SumType)
	}
	b, err := json.Marshal(body.Value)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type msgDecoderFunc func(cell *boc.Cell) (*uint32, *MsgOpName, any, error)

// InternalMessageDecoder takes in a message body as a cell and tries to decode it based on the contract type or the first 4 bytes.
// It returns an opcode, an operation name and a decoded body.
func InternalMessageDecoder(cell *boc.Cell, interfaces []ContractInterface) (*MsgOpCode, *MsgOpName, any, error) {
	if cell.BitsAvailableForRead() < 32 {
		return nil, nil, nil, nil
	}
	if cell.BitsAvailableForRead() != cell.BitSize() {
		cell = cell.CopyRemaining() //because body can be part of the message cell
	}
	for _, i := range interfaces {
		for _, f := range i.IntMsgs() {
			t, m, v, err := f(cell)
			if err == nil {
				return t, m, v, nil
			}
		}
	}
	cell.ResetCounters()
	tag64, err := cell.PickUint(32)
	if err != nil {
		return nil, nil, nil, err
	}
	tag := uint32(tag64)
	f := opcodedMsgInDecodeFunctions[tag]
	if f != nil {
		t, o, b, err := f(cell)
		if err == nil {
			return t, o, b, nil
		}
	}
	return &tag, nil, nil, nil
}

func ExtInMessageDecoder(cell *boc.Cell, interfaces []ContractInterface) (*MsgOpCode, *MsgOpName, any, error) {
	if cell.BitsAvailableForRead() < 32 {
		return nil, nil, nil, nil
	}
	if cell.BitsAvailableForRead() != cell.BitSize() {
		cell = cell.CopyRemaining() //because body can be part of the message cell
	}
	for _, i := range interfaces {
		for _, f := range i.ExtInMsgs() {
			t, m, v, err := f(cell)
			if err == nil {
				return t, m, v, nil
			}
		}
	}
	cell.ResetCounters()
	tag64, err := cell.PickUint(32)
	if err != nil {
		return nil, nil, nil, err
	}
	tag := uint32(tag64)
	f := opcodedMsgExtInDecodeFunctions[tag]
	if f != nil {
		t, o, b, err := f(cell)
		if err == nil {
			return t, o, b, nil
		}
	}
	return &tag, nil, nil, nil
}

func ExtOutMessageDecoder(cell *boc.Cell, interfaces []ContractInterface, dest tlb.MsgAddress) (*MsgOpCode, *MsgOpName, any, error) { //todo: change to tlb.MsgAddressExt
	if cell.BitsAvailableForRead() < 32 {
		return nil, nil, nil, nil
	}
	if cell.BitsAvailableForRead() != cell.BitSize() {
		cell = cell.CopyRemaining() //because body can be part of the message cell
	}
	for _, i := range interfaces {
		for _, f := range i.ExtOutMsgs() {
			t, m, v, err := f(cell)
			if err == nil {
				return t, m, v, nil
			}
		}
	}
	cell.ResetCounters()
	tag64, err := cell.PickUint(32)
	if err != nil {
		return nil, nil, nil, err
	}
	tag := uint32(tag64)
	f := opcodedMsgExtOutDecodeFunctions[tag]
	if f != nil {
		t, o, b, err := f(cell)
		if err == nil {
			return t, o, b, nil
		}
	}
	return &tag, nil, nil, nil
}

func completedRead(cell *boc.Cell) bool {
	return cell.RefsAvailableForRead() == 0 && cell.BitsAvailableForRead() == 0
}

// MsgOpName is a human-friendly name for a message's operation which is identified by the first 4 bytes of the message's body.
type MsgOpName = string

const (
	UnknownMsgOp MsgOpName = "Unknown"
	EmptyMsgOp   MsgOpName = ""
)

// MsgOpCode is the first 4 bytes of a message body identifying an operation to be performed.
type MsgOpCode = uint32
