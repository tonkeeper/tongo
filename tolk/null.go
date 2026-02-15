package tolk

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type NullValue struct{}

func (n *NullValue) Equal(other any) bool {
	_, ok := other.(NullValue)
	if !ok {
		return false
	}
	return true
}

func (n *NullValue) Unmarshal(cell *boc.Cell, ty tolkParser.NullLiteral, decoder *Decoder) error {
	return nil
}

func (n *NullValue) Marshal(cell *boc.Cell, ty tolkParser.NullLiteral, encoder *Encoder) error {
	return nil
}

func (n *NullValue) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

type VoidValue struct{}

func (v *VoidValue) Equal(other any) bool {
	_, ok := other.(VoidValue)
	if !ok {
		return false
	}
	return true
}

func (v *VoidValue) Unmarshal(cell *boc.Cell, ty tolkParser.Void, decoder *Decoder) error {
	return nil
}

func (v *VoidValue) Marshal(cell *boc.Cell, ty tolkParser.Void, encoder *Encoder) error {
	return nil
}

func (v *VoidValue) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
