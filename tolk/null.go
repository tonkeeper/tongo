package tolk

import (
	"fmt"

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

func (n *NullValue) Unmarshal(cell *boc.Cell, ty parser.NullLiteral, decoder *Decoder) error {
	return nil
}

func (n *NullValue) Marshal(cell *boc.Cell, ty parser.NullLiteral, encoder *Encoder) error {
	return nil
}

func (n NullValue) MarshalJSON() ([]byte, error) {
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

func (v *VoidValue) Unmarshal(cell *boc.Cell, ty parser.Void, decoder *Decoder) error {
	return nil
}

func (v *VoidValue) Marshal(cell *boc.Cell, ty parser.Void, encoder *Encoder) error {
	return nil
}

func (v VoidValue) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

type Unknown struct{}

func (u *Unknown) Equal(other any) bool {
	_, ok := other.(Unknown)
	if !ok {
		return false
	}
	return true
}

func (u *Unknown) Unmarshal(cell *boc.Cell, ty parser.Unknown, decoder *Decoder) error {
	return fmt.Errorf("cannot unmarshal unknown type")
}

func (u *Unknown) Marshal(cell *boc.Cell, ty parser.Unknown, encoder *Encoder) error {
	return fmt.Errorf("cannot marshal unknown type")
}

func (u Unknown) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal unknown type")
}
