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

type NFTPayload struct {
	SumType string
	OpCode  *uint32
	Value   any
}

type NFTOpName = string

const (
	EmptyNFTOp   NFTOpName = ""
	UnknownNFTOp NFTOpName = "Cell"
)

// NFTOpCode is the first 4 bytes of a message body identifying an operation to be performed.
type NFTOpCode = uint32

func (p *NFTPayload) UnmarshalJSON(data []byte) error {
	var r struct {
		SumType string
		OpCode  *uint32
		Value   json.RawMessage
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	p.SumType = r.SumType
	p.OpCode = r.OpCode
	if p.SumType == EmptyNFTOp {
		return nil
	}
	if p.SumType == UnknownNFTOp {
		c := boc.NewCell()
		err := json.Unmarshal(r.Value, c)
		if err != nil {
			return err
		}
		p.Value = c
		return nil
	}
	t, ok := KnownNFTTypes[p.SumType]
	if !ok {
		return fmt.Errorf("unknown jetton payload type '%v'", p.SumType)
	}
	o := reflect.New(reflect.TypeOf(t))
	err := json.Unmarshal(r.Value, o.Interface())
	if err != nil {
		return err
	}
	p.Value = o.Elem().Interface()
	return nil
}

func (p NFTPayload) MarshalJSON() ([]byte, error) {
	if p.SumType == EmptyNFTOp {
		return []byte("{}"), nil
	}
	buf := bytes.NewBufferString(`{"SumType": "` + p.SumType + `",`)
	if p.OpCode != nil {
		fmt.Fprintf(buf, `"OpCode":%v,`, *p.OpCode)
	}
	buf.WriteString(`"Value":`)
	if p.SumType == UnknownNFTOp {
		c, ok := p.Value.(*boc.Cell)
		if !ok {
			return nil, fmt.Errorf("unknown NFTPayload should be Cell")
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
	if KnownNFTTypes[p.SumType] == nil {
		return nil, fmt.Errorf("unknown NFTPayload type %v", p.SumType)
	}
	b, err := json.Marshal(p.Value)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (p NFTPayload) MarshalTLB(c *boc.Cell, e *tlb.Encoder) error {
	if p.SumType == EmptyNFTOp {
		return nil
	}
	if p.OpCode != nil {
		err := c.WriteUint(uint64(*p.OpCode), 32)
		if err != nil {
			return err
		}
	} else if op, ok := NFTOpCodes[p.SumType]; ok {
		err := c.WriteUint(uint64(op), 32)
		if err != nil {
			return err
		}
	}
	return tlb.Marshal(c, p.Value)
}

func (j *NFTPayload) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error {
	if completedRead(cell) {
		return nil
	}
	tempCell := cell.CopyRemaining()
	op64, err := tempCell.ReadUint(32)
	if errors.Is(err, boc.ErrNotEnoughBits) {
		j.SumType = UnknownNFTOp
		j.Value = cell.CopyRemaining()
		return nil
	}
	op := uint32(op64)
	j.OpCode = &op
	f, ok := funcNFTDecodersMapping[NFTOpCode(op64)]
	if ok && f != nil {
		err = f(j, tempCell)
		if err == nil {
			return nil
		}
	}
	j.SumType = UnknownNFTOp
	j.Value = cell.CopyRemaining()

	return nil
}
