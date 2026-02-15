package tolk

import (
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Any boc.Cell

func (a *Any) Unmarshal(cell *boc.Cell, ty tolkParser.Cell, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return fmt.Errorf("failed to get next ref: %w", err)
	}
	*a = Any(*ref)

	return nil
}

func (a *Any) Marshal(cell *boc.Cell, ty tolkParser.Cell, encoder *Encoder) error {
	c := boc.Cell(*a)
	ref := c.CopyRemaining()
	err := cell.AddRef(ref)
	if err != nil {
		return fmt.Errorf("failed to add ref: %w", err)
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

func (a *Any) MarshalJSON() ([]byte, error) {
	data, err := boc.Cell(*a).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal any: %w", err)
	}
	return data, nil
}

func (a *Any) UnmarshalJSON(b []byte) error {
	v := &boc.Cell{}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("failed to unmarshal ref: %w", err)
	}
	*a = Any(*v)
	return nil
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
		return fmt.Errorf("failed to write remaining bits: %w", err)
	}
	for i, ref := range c.Refs() {
		err = cell.AddRef(ref)
		if err != nil {
			return fmt.Errorf("failed to add %v remaining ref: %w", i, err)
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

func (r *RemainingValue) MarshalJSON() ([]byte, error) {
	data, err := boc.Cell(*r).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal remainings: %w", err)
	}
	return data, nil
}

func (r *RemainingValue) UnmarshalJSON(b []byte) error {
	v := &boc.Cell{}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("failed to unmarshal remainigs: %w", err)
	}
	*r = RemainingValue(*v)
	return nil
}

type OptValue struct {
	IsExists bool
	Val      Value
}

func (o *OptValue) Unmarshal(cell *boc.Cell, ty tolkParser.Nullable, decoder *Decoder) error {
	isExists, err := cell.ReadBit()
	if err != nil {
		return fmt.Errorf("failed to read optinal value existance bit: %w", err)
	}
	o.IsExists = isExists
	if isExists {
		err = o.Val.Unmarshal(cell, ty.Inner, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal optinal value: %w", err)
		}
	}
	return nil
}

func (o *OptValue) Marshal(cell *boc.Cell, ty tolkParser.Nullable, encoder *Encoder) error {
	err := cell.WriteBit(o.IsExists)
	if err != nil {
		return fmt.Errorf("failed to write optinal value existance bit: %w", err)
	}
	if o.IsExists {
		err = o.Val.Marshal(cell, ty.Inner, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optinal value: %w", err)
		}
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
	var jsonOptValue = struct {
		IsExists bool   `json:"isExists"`
		Val      *Value `json:"value,omitempty"`
	}{
		IsExists: o.IsExists,
	}
	if o.IsExists {
		jsonOptValue.Val = &o.Val
	}

	data, err := json.Marshal(jsonOptValue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal optinal value: %w", err)
	}
	return data, nil
}

func (o *OptValue) UnmarshalJSON(b []byte) error {
	var jsonOptValue = struct {
		IsExists bool   `json:"isExists"`
		Val      *Value `json:"value,omitempty"`
	}{}
	if err := json.Unmarshal(b, &jsonOptValue); err != nil {
		return fmt.Errorf("failed to unmarshal optinal value: %w", err)
	}
	o.IsExists = jsonOptValue.IsExists
	if o.IsExists {
		o.Val = *jsonOptValue.Val
	}

	return nil
}

type RefValue Value

func (r *RefValue) Unmarshal(cell *boc.Cell, ty tolkParser.CellOf, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return fmt.Errorf("failed to get next ref: %w", err)
	}
	innerV := Value{}
	err = innerV.Unmarshal(ref, ty.Inner, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ref: %w", err)
	}
	*r = RefValue(innerV)

	return nil
}

func (r *RefValue) Marshal(cell *boc.Cell, ty tolkParser.CellOf, encoder *Encoder) error {
	val := Value(*r)
	ref := boc.NewCell()
	err := val.Marshal(ref, ty.Inner, encoder)
	if err != nil {
		return fmt.Errorf("failed to marshal ref: %w", err)
	}
	err = cell.AddRef(ref)
	if err != nil {
		return fmt.Errorf("failed to add ref: %w", err)
	}

	return nil
}

func (r *RefValue) Equal(other any) bool {
	otherRefValue, ok := other.(RefValue)
	if !ok {
		return false
	}
	v := Value(*r)
	return v.Equal(Value(otherRefValue))
}
