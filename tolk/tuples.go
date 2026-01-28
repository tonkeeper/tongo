package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type Tensor struct {
	Items []Ty `json:"items"`
}

func (Tensor) SetValue(v *Value, val any) error {
	t, ok := val.(TensorValues)
	if !ok {
		return fmt.Errorf("value is not a tensor")
	}
	v.tensor = &t
	return nil
}

func (t Tensor) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	list := make(TensorValues, len(t.Items))
	for i, item := range t.Items {
		inner := Value{}
		err := item.UnmarshalTolk(cell, &inner, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	err := v.SetValue(list)
	if err != nil {
		return err
	}
	return nil
}

func (Tensor) Equal(v Value, o Value) bool {
	return false
}

type TupleWith struct {
	Items []Ty `json:"items"`
}

func (TupleWith) SetValue(v *Value, val any) error {
	t, ok := val.(TupleValues)
	if !ok {
		return fmt.Errorf("value is not a tuple")
	}
	v.tupleWith = &t
	return nil
}

func (t TupleWith) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	list := make(TupleValues, len(t.Items))
	for i, item := range t.Items {
		inner := Value{}
		err := item.UnmarshalTolk(cell, &inner, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	err := v.SetValue(list)
	if err != nil {
		return err
	}

	return nil
}

func (TupleWith) Equal(v Value, o Value) bool {
	return false
}

type TupleAny struct{}

func (TupleAny) SetValue(v *Value, val any) error {
	return fmt.Errorf("tuple any is not supported")
}

func (t TupleAny) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return fmt.Errorf("tuple any is not supported")
}

func (TupleAny) Equal(v Value, o Value) bool {
	return false
}
