package tolk

import (
	"encoding/json"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Any boc.Cell

func (a *Any) Unmarshal(cell *boc.Cell, ty tolkParser.Cell, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	*a = Any(*ref)

	return nil
}

func (a *Any) Marshal(cell *boc.Cell, ty tolkParser.Cell, encoder *Encoder) error {
	c := boc.Cell(*a)
	ref := c.CopyRemaining()
	err := cell.AddRef(ref)
	if err != nil {
		return err
	}

	return nil
}

func (a *Any) Equal(o any) bool {
	other, ok := o.(Any)
	if !ok {
		return false
	}
	cellV := boc.Cell(*a)
	vHash, err := cellV.HashString()
	if err != nil {
		return false
	}
	cellO := boc.Cell(other)
	oHash, err := cellO.HashString()
	if err != nil {
		return false
	}
	return oHash == vHash
}

type RemainingValue boc.Cell

func (r *RemainingValue) Unmarshal(cell *boc.Cell, ty tolkParser.Remaining, decoder *Decoder) error {
	rem := cell.CopyRemaining()
	if rem != nil {
		*r = RemainingValue(*rem)
	}
	return nil
}

func (r *RemainingValue) Marshal(cell *boc.Cell, ty tolkParser.Remaining, encoder *Encoder) error {
	c := boc.Cell(*r)
	err := cell.WriteBitString(c.ReadRemainingBits())
	if err != nil {
		return err
	}
	for _, ref := range c.Refs() {
		err = cell.AddRef(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RemainingValue) Equal(o any) bool {
	other, ok := o.(RemainingValue)
	if !ok {
		return false
	}
	cellV := boc.Cell(*r)
	vHash, err := cellV.HashString()
	if err != nil {
		return false
	}
	cellO := boc.Cell(other)
	oHash, err := cellO.HashString()
	if err != nil {
		return false
	}
	return oHash == vHash
}

type OptValue struct {
	IsExists bool
	Val      Value
}

func (o *OptValue) Unmarshal(cell *boc.Cell, ty tolkParser.Nullable, decoder *Decoder) error {
	isExists, err := cell.ReadBit()
	if err != nil {
		return err
	}
	o.IsExists = isExists
	if isExists {
		err = o.Val.Unmarshal(cell, ty.Inner, decoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OptValue) Marshal(cell *boc.Cell, ty tolkParser.Nullable, encoder *Encoder) error {
	err := cell.WriteBit(o.IsExists)
	if err != nil {
		return err
	}
	if o.IsExists {
		return o.Val.Marshal(cell, ty.Inner, encoder)
	}

	return nil
}

func (o *OptValue) Equal(other any) bool {
	otherOptValue, ok := other.(OptValue)
	if !ok {
		return false
	}
	if o.IsExists != otherOptValue.IsExists {
		return false
	}
	if o.IsExists {
		return o.Val.Equal(otherOptValue.Val)
	}
	return true
}

func (o *OptValue) MarshalJSON() ([]byte, error) {
	if o.IsExists {
		return json.Marshal(o.Val)
	}
	return []byte("null"), nil
}

type RefValue Value

func (r *RefValue) Unmarshal(cell *boc.Cell, ty tolkParser.CellOf, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	innerV := Value{}
	err = innerV.Unmarshal(ref, ty.Inner, decoder)
	if err != nil {
		return err
	}
	*r = RefValue(innerV)

	return nil
}

func (r *RefValue) Marshal(cell *boc.Cell, ty tolkParser.CellOf, encoder *Encoder) error {
	val := Value(*r)
	ref := boc.NewCell()
	err := val.Marshal(ref, ty.Inner, encoder)
	if err != nil {
		return err
	}
	err = cell.AddRef(ref)
	if err != nil {
		return err
	}

	return nil
}

func (r *RefValue) Equal(other any) bool {
	otherRefValue, ok := other.(RefValue)
	if !ok {
		return false
	}
	return r.Equal(otherRefValue)
}
