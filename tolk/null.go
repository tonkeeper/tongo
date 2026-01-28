package tolk

import "github.com/tonkeeper/tongo/boc"

type NullLiteral struct{}

func (NullLiteral) SetValue(v *Value, val any) error {
	return nil
}

func (NullLiteral) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return nil
}

func (NullLiteral) Equal(v Value, o Value) bool {
	return false
}

type Void struct{}

func (Void) SetValue(v *Value, val any) error {
	return nil
}

func (Void) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return nil
}

func (Void) Equal(v Value, o Value) bool {
	return false
}
