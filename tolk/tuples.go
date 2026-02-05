package tolk

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type TupleValues []Value

func (v *TupleValues) Unmarshal(cell *boc.Cell, ty tolkParser.TupleWith, decoder *Decoder) error {
	list := make(TupleValues, len(ty.Items))
	for i, item := range ty.Items {
		inner := Value{}
		err := inner.Unmarshal(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	*v = list
	return nil
}

func (v *TupleValues) Marshal(cell *boc.Cell, ty tolkParser.TupleWith, encoder *Encoder) error {
	for i, item := range []Value(*v) {
		err := item.Marshal(cell, ty.Items[i], encoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *TupleValues) Equal(other any) bool {
	otherTupleValues, ok := other.(TupleValues)
	if !ok {
		return false
	}
	wV := *v
	if len(otherTupleValues) != len(wV) {
		return false
	}
	for i := range wV {
		if !wV[i].Equal(otherTupleValues[i]) {
			return false
		}
	}
	return true
}

type TensorValues []Value

func (v *TensorValues) Unmarshal(cell *boc.Cell, ty tolkParser.Tensor, decoder *Decoder) error {
	list := make(TensorValues, len(ty.Items))
	for i, item := range ty.Items {
		inner := Value{}
		err := inner.Unmarshal(cell, item, decoder)
		if err != nil {
			return err
		}
		list[i] = inner
	}
	*v = list
	return nil
}

func (v *TensorValues) Marshal(cell *boc.Cell, ty tolkParser.Tensor, encoder *Encoder) error {
	for i, item := range []Value(*v) {
		err := item.Marshal(cell, ty.Items[i], encoder)
		if err != nil {
			return err
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
