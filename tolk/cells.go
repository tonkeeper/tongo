package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type Cell struct{}

func (Cell) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	return nil
}

func (Cell) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	err = v.SetValue(Any(*ref))
	if err != nil {
		return err
	}
	return nil
}

func (Cell) Equal(v Value, o Value) bool {
	return false
}

type Slice struct{}

func (Slice) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	return nil
}

func (Slice) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return fmt.Errorf("slice is not supported")
}

func (Slice) Equal(v Value, o Value) bool {
	return false
}

type Builder struct{}

func (Builder) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	return nil
}

func (Builder) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return fmt.Errorf("builder is not supported")
}

func (Builder) Equal(v Value, o Value) bool {
	return false
}

type Callable struct{}

func (Callable) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	return nil
}

func (Callable) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return fmt.Errorf("callable is not supported")
}

func (Callable) Equal(v Value, o Value) bool {
	return false
}

type Remaining struct{}

func (Remaining) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	return nil
}

func (Remaining) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	rem := cell.CopyRemaining()
	if rem != nil {
		err := v.SetValue(Any(*rem))
		if err != nil {
			return err
		}
	}
	return nil
}

func (Remaining) Equal(v Value, o Value) bool {
	return false
}

type Nullable struct {
	Inner Ty `json:"inner"`
}

func (Nullable) SetValue(v *Value, val any) error {
	o, ok := val.(OptValue)
	if !ok {
		return fmt.Errorf("value is not an optional value")
	}
	v.optionalValue = &o
	return nil
}

func (n Nullable) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	isExists, err := cell.ReadBit()
	if err != nil {
		return err
	}
	innerV := Value{
		valType: n.Inner,
	}
	optV := OptValue{
		IsExists: isExists,
		Val:      innerV,
	}
	if isExists {
		err = n.Inner.UnmarshalTolk(cell, &innerV, decoder)
		if err != nil {
			return err
		}
		optV.Val = innerV
	}
	err = v.SetValue(optV)
	if err != nil {
		return err
	}
	return nil
}

func (Nullable) Equal(v Value, o Value) bool {
	return false
}

type CellOf struct {
	Inner Ty `json:"inner"`
}

func (CellOf) SetValue(v *Value, val any) error {
	r, ok := val.(RefValue)
	if !ok {
		return fmt.Errorf("value is not a ref value")
	}
	v.refValue = &r
	return nil
}

func (c CellOf) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	innerV := Value{}
	err = c.Inner.UnmarshalTolk(ref, &innerV, decoder)
	if err != nil {
		return err
	}
	err = v.SetValue(RefValue(innerV))
	if err != nil {
		return err
	}

	return nil
}

func (CellOf) Equal(v Value, o Value) bool {
	return false
}
