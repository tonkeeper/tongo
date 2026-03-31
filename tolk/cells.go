package tolk

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Any boc.Cell

func (a *Any) Unmarshal(cell *boc.Cell, ty parser.Cell, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return fmt.Errorf("failed to get next ref: %w", err)
	}
	*a = Any(*ref)

	return nil
}

func (a *Any) Marshal(cell *boc.Cell, ty parser.Cell, encoder *Encoder) error {
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

func (a Any) MarshalJSON() ([]byte, error) {
	data, err := boc.Cell(a).MarshalJSON()
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

type RemainingValue struct {
	IsRef bool
	Value boc.Cell
}

func (r *RemainingValue) Unmarshal(cell *boc.Cell, ty parser.Remaining, decoder *Decoder) error {
	rem := cell.CopyRemaining()
	cell.ReadRemainingBits()
	if rem != nil {
		isRef := cell.BitsAvailableForRead() == 0 && cell.RefsAvailableForRead() > 0
		*r = RemainingValue{
			IsRef: isRef,
			Value: *rem,
		}
	}
	return nil
}

func (r *RemainingValue) Marshal(cell *boc.Cell, ty parser.Remaining, encoder *Encoder) error {
	c := r.Value
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
	cellV := r.Value
	vHash, err := cellV.HashString()
	if err != nil {
		return false
	}
	cellO := other.Value
	oHash, err := cellO.HashString()
	if err != nil {
		return false
	}
	return oHash == vHash && r.IsRef == other.IsRef
}

func (r RemainingValue) MarshalJSON() ([]byte, error) {
	cellData, err := json.Marshal(r.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal remainings data: %w", err)
	}
	if r.Value.BitsAvailableForRead() == 0 && r.Value.RefsAvailableForRead() == 0 {
		cellData = []byte("\"\"")
	}
	if len(cellData) < 2 {
		return nil, fmt.Errorf("invalid remaining cell data: %v", cellData)
	}

	var jsonData = struct {
		IsRef bool   `json:"isRef"`
		Value string `json:"value"`
	}{
		IsRef: r.IsRef,
		Value: string(cellData[1 : len(cellData)-1]),
	}
	data, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal remaining: %w", err)
	}
	return data, nil
}

type OptValue struct {
	IsExists bool
	Val      Value
}

func (o *OptValue) Unmarshal(cell *boc.Cell, ty parser.Nullable, decoder *Decoder) error {
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

func (o *OptValue) Marshal(cell *boc.Cell, ty parser.Nullable, encoder *Encoder) error {
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

func (o OptValue) MarshalJSON() ([]byte, error) {
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

func (r *RefValue) Unmarshal(cell *boc.Cell, ty parser.CellOf, decoder *Decoder) error {
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

func (r *RefValue) Marshal(cell *boc.Cell, ty parser.CellOf, encoder *Encoder) error {
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

func (r RefValue) MarshalJSON() ([]byte, error) {
	v := Value(r)
	data, err := v.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ref: %w", err)
	}
	return data, nil
}

type SnakeString string

func (s *SnakeString) Unmarshal(cell *boc.Cell, ty parser.String, decoder *Decoder) error {
	var res boc.BitString
	res.Append(cell.ReadRemainingBits())
	ref := *cell
	for ref.RefsSize() > 0 {
		nxt, err := ref.NextRef()
		if err != nil {
			return fmt.Errorf("failed to get next ref for snake string: %w", err)
		}
		res.Append(nxt.ReadRemainingBits())
	}

	if res.BitsAvailableForRead()%8 != 0 {
		return fmt.Errorf("incorrect number of bits for snake string: %v", res.BitsAvailableForRead())
	}
	buf, err := res.GetTopUppedArray()
	if err != nil {
		return fmt.Errorf("failed to get top upped array: %w", err)
	}
	if !utf8.Valid(buf) {
		return fmt.Errorf("invalid UTF-8 in snake string")
	}
	*s = SnakeString(buf)

	return nil
}

func (s *SnakeString) Marshal(cell *boc.Cell, ty parser.String, encoder *Encoder) error {
	curr := cell
	for _, r := range string(*s) {
		l := utf8.RuneLen(r)
		if cell.BitsAvailableForWrite() < l {
			next := boc.NewCell()
			err := curr.AddRef(next)
			if err != nil {
				return fmt.Errorf("failed to add ref: %w", err)
			}
			curr = next
		}

		if err := curr.WriteUint(uint64(r), l); err != nil {
			return fmt.Errorf("failed to write rune: %w", err)
		}
	}

	return nil
}

func (s *SnakeString) Equal(o any) bool {
	other, ok := o.(SnakeString)
	if !ok {
		return false
	}
	return string(*s) == string(other)
}
