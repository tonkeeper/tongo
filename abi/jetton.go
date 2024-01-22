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
	EmptyJettonOp   = ""
	UnknownJettonOp = "Cell"
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
	if cell.BitsAvailableForRead() == 0 && cell.RefsAvailableForRead() == 0 {
		return nil
	}
	tempCell := cell.CopyRemaining()
	op64, err := tempCell.ReadUint(32)
	if errors.Is(err, boc.ErrNotEnoughBits) {
		j.SumType = UnknownJettonOp
		j.Value = cell.CopyRemaining()
		return nil
	}
	op := uint32(op64)
	j.OpCode = &op
	f, ok := funcJettonDecodersMapping[JettonOpCode(op64)]

	if ok && f != nil {
		err = f(j, tempCell)
		if err == nil {
			return nil
		}
	}

	j.SumType = UnknownJettonOp
	j.Value = cell.CopyRemaining()

	return nil
}
