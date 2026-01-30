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
	v.sumType = "tensor"
	return nil
}

func (t Tensor) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	list := make(TensorValues, len(t.Items))
	for i, item := range t.Items {
		inner := Value{}
		err := inner.UnmarshalTolk(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	err := t.SetValue(v, list)
	if err != nil {
		return err
	}
	return nil
}

func (v *TensorValues) UnmarshalTolk(cell *boc.Cell, ty Tensor, decoder *Decoder) error {
	list := make(TensorValues, len(ty.Items))
	for i, item := range ty.Items {
		inner := Value{}
		err := inner.UnmarshalTolk(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	*v = list
	return nil
}

func (Tensor) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.tensor == nil {
	//	return fmt.Errorf("tensor is nil")
	//}
	//
	//for _, tv := range []Value(*v.tensor) {
	//	err := tv.valType.MarshalTolk(cell, &tv)
	//	if err != nil {
	//		return err
	//	}
	//}

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
	v.sumType = "tupleWith"
	return nil
}

func (t TupleWith) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	list := make(TupleValues, len(t.Items))
	for i, item := range t.Items {
		inner := Value{}
		err := inner.UnmarshalTolk(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	err := t.SetValue(v, list)
	if err != nil {
		return err
	}

	return nil
}

func (v *TupleValues) UnmarshalTolk(cell *boc.Cell, ty TupleWith, decoder *Decoder) error {
	list := make(TupleValues, len(ty.Items))
	for i, item := range ty.Items {
		inner := Value{}
		err := inner.UnmarshalTolk(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	*v = list
	return nil
}

func (TupleWith) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.tupleWith == nil {
	//	return fmt.Errorf("tupleWith is nil")
	//}
	//
	//for _, tv := range []Value(*v.tupleWith) {
	//	err := tv.valType.MarshalTolk(cell, &tv)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (TupleWith) Equal(v Value, o Value) bool {
	return false
}

type TupleAny struct{}

func (TupleAny) SetValue(v *Value, val any) error {
	return fmt.Errorf("tuple any is not supported")
}

func (TupleAny) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return fmt.Errorf("tuple any is not supported")
}

func (TupleAny) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("tuple any is not supported")
}

func (TupleAny) Equal(v Value, o Value) bool {
	return false
}
