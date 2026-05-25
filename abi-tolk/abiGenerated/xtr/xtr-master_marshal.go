// Code generated - DO NOT EDIT.

package abiXtr

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *PushXTR) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixPushXTR); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}
func (v PushXTR) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixPushXTR, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}
func (v PushXTR) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg PushXTR) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
