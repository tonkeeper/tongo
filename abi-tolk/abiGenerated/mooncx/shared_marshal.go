// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClmmDepositBin) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.X.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .X: %v", err)
	}
	if err = v.Y.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Y: %v", err)
	}
	return nil
}

func (v ClmmDepositBin) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.X.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .X: %v", err)
	}
	if err = v.Y.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Y: %v", err)
	}
	return nil
}

func (v ClmmDepositBin) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmBin) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.L.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .L: %v", err)
	}
	if err = v.XFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XFee: %v", err)
	}
	if err = v.YFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YFee: %v", err)
	}
	return nil
}

func (v ClmmBin) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.L.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .L: %v", err)
	}
	if err = v.XFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XFee: %v", err)
	}
	if err = v.YFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YFee: %v", err)
	}
	return nil
}

func (v ClmmBin) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *LpSplit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ARecipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ARecipient: %v", err)
	}
	if err = v.ABps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ABps: %v", err)
	}
	if err = v.BRecipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BRecipient: %v", err)
	}
	return nil
}

func (v LpSplit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ARecipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ARecipient: %v", err)
	}
	if err = v.ABps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ABps: %v", err)
	}
	if err = v.BRecipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BRecipient: %v", err)
	}
	return nil
}

func (v LpSplit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
