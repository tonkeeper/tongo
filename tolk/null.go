package tolk

import (
	"github.com/tonkeeper/tongo/boc"
)

type NullValue struct{}

func (n *NullValue) Equal(other any) bool {
	_, ok := other.(NullValue)
	if !ok {
		return false
	}
	return true
}

func (n *NullValue) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	return nil
}

func (n *NullValue) Marshal(cell *boc.Cell, encoder *Encoder) error {
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

func (v *VoidValue) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	return nil
}

func (v *VoidValue) Marshal(cell *boc.Cell, encoder *Encoder) error {
	return nil
}

func (v VoidValue) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
