package tolk

import (
	"github.com/tonkeeper/tongo/boc"
)

type NullLiteral struct{}

func (NullLiteral) SetValue(v *Value, val any) error {
	v.sumType = "nullLiteral"
	return nil
}

func (NullLiteral) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return nil
}

func (n *NullValue) UnmarshalTolk(cell *boc.Cell, ty NullLiteral, decoder *Decoder) error {
	return nil
}

func (NullLiteral) MarshalTolk(cell *boc.Cell, v *Value) error {
	return nil
}

func (NullLiteral) Equal(v Value, o Value) bool {
	return false
}

type Void struct{}

func (Void) SetValue(v *Value, val any) error {
	v.sumType = "void"
	return nil
}

func (Void) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return nil
}

func (v *VoidValue) UnmarshalTolk(cell *boc.Cell, ty Void, decoder *Decoder) error {
	return nil
}

func (Void) MarshalTolk(cell *boc.Cell, v *Value) error {
	return nil
}

func (Void) Equal(v Value, o Value) bool {
	return false
}
