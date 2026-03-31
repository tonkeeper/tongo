package tolk

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type ArrayOf []Value

func (a *ArrayOf) Unmarshal(cell *boc.Cell, ty parser.ArrayOf, decoder *Decoder) error {
	ln, err := cell.ReadUint(8)
	if err != nil {
		return fmt.Errorf("failed to get array size: %w", err)
	}
	if ln == 0 {
		_, err = cell.ReadBit()
		if err != nil {
			return fmt.Errorf("failed to get array has ref flag: %w", err)
		}
		return nil
	}
	list := make(ArrayOf, ln)
	curr := cell
	for i := 0; i < int(ln); i++ {
		hasRef, err := cell.ReadBit()
		if err != nil {
			return fmt.Errorf("failed to get array has ref flag: %w", err)
		}
		if hasRef {
			curr, err = curr.NextRef()
			if err != nil {
				return fmt.Errorf("failed to get array next ref: %w", err)
			}
		}
		for curr.BitsAvailableForRead() > 0 {
			var val Value
			if err := val.Unmarshal(cell, ty.Inner, decoder); err != nil {
				return fmt.Errorf("failed to unmarshal array[%v] element: %w", i, err)
			}
			list = append(list, val)
		}
	}
	*a = list

	return nil
}

func (a *ArrayOf) Marshal(cell *boc.Cell, ty parser.ArrayOf, encoder *Encoder) error {
	arr := []Value(*a)
	err := cell.WriteUint(uint64(len(arr)), 8)
	if err != nil {
		return fmt.Errorf("failed to write array size: %w", err)
	}
	err = cell.WriteBit(true)
	if err != nil {
		return fmt.Errorf("failed to write array has ref flag: %w", err)
	}
	curr := cell
	for i, item := range arr {
		err = item.Marshal(curr, ty.Inner, encoder)
		if err != nil {
			if errors.Is(err, boc.ErrBitStingOverflow) {
				next := boc.NewCell()
				err = curr.AddRef(next)
				if err != nil {
					return fmt.Errorf("failed to add ref for array: %w", err)
				}
				curr = next

				err = item.Marshal(curr, ty.Inner, encoder)
				if err == nil {
					continue
				}
			}
			return fmt.Errorf("failed to marshal %v tuple's value: %w", i, err)
		}
	}
	return nil
}

func (a *ArrayOf) Equal(other any) bool {
	otherArrayValues, ok := other.(ArrayOf)
	if !ok {
		return false
	}
	wA := *a
	if len(otherArrayValues) != len(wA) {
		return false
	}
	for i := range wA {
		if !wA[i].Equal(otherArrayValues[i]) {
			return false
		}
	}
	return true
}

type LispListOf []Value

func (ll *LispListOf) Unmarshal(cell *boc.Cell, ty parser.LispListOf, decoder *Decoder) error {
	list := make(LispListOf, 0)
	var err error
	curr := cell
	for curr.RefsAvailableForRead() > 0 {
		curr, err = curr.NextRef()
		if err != nil {
			return fmt.Errorf("failed to get lisp list next ref: %w", err)
		}
		var val Value
		if err := val.Unmarshal(cell, ty.Inner, decoder); err != nil {
			return fmt.Errorf("failed to unmarshal lisp list element: %w", err)
		}
		list = append(list, val)
	}

	slices.Reverse(list)
	*ll = list

	return nil
}

func (ll *LispListOf) Marshal(cell *boc.Cell, ty parser.LispListOf, encoder *Encoder) error {
	var err error
	arr := []Value(*ll)
	curr := cell
	for i, item := range arr {
		ref := boc.NewCell()
		err = curr.AddRef(ref)
		if err != nil {
			return fmt.Errorf("failed to add ref for lisp list element: %w", err)
		}
		curr = ref
		err = item.Marshal(curr, ty.Inner, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal %v tuple's value: %w", i, err)
		}
	}
	return nil
}

func (ll *LispListOf) Equal(other any) bool {
	otherLispListValues, ok := other.(ArrayOf)
	if !ok {
		return false
	}
	wLL := *ll
	if len(otherLispListValues) != len(wLL) {
		return false
	}
	for i := range wLL {
		if !wLL[i].Equal(otherLispListValues[i]) {
			return false
		}
	}
	return true
}

type TensorValues []Value

func (v *TensorValues) Unmarshal(cell *boc.Cell, ty parser.Tensor, decoder *Decoder) error {
	list := make(TensorValues, len(ty.Items))
	for i, item := range ty.Items {
		inner := Value{}
		err := inner.Unmarshal(cell, item, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal %v tensor's value: %w", i, err)
		}
		list[i] = inner
	}
	*v = list
	return nil
}

func (v *TensorValues) Marshal(cell *boc.Cell, ty parser.Tensor, encoder *Encoder) error {
	for i, item := range []Value(*v) {
		err := item.Marshal(cell, ty.Items[i], encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal %v tensor's value: %w", i, err)
		}
	}
	return nil
}

func (v *TensorValues) Equal(other any) bool {
	otherTensorValues, ok := other.(TensorValues)
	if !ok {
		return false
	}
	wV := *v
	if len(otherTensorValues) != len(wV) {
		return false
	}
	for i := range wV {
		if !wV[i].Equal(otherTensorValues[i]) {
			return false
		}
	}
	return true
}

func (v TensorValues) MarshalJSON() ([]byte, error) {
	var s strings.Builder
	s.WriteRune('[')
	for i, item := range []Value(v) {
		data, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %v tensor's value: %w", i, err)
		}
		s.Write(data)
		if i != len([]Value(v))-1 {
			s.WriteRune(',')
		}
	}
	s.WriteRune(']')
	return []byte(s.String()), nil
}
