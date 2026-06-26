// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *SourceContent) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Version.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Version: %v", err)
	}
	if err = v.Url.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Url: %v", err)
	}
	return nil
}

func (v SourceContent) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Version.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Version: %v", err)
	}
	if err = v.Url.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Url: %v", err)
	}
	return nil
}

func (v SourceContent) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
