// Code generated - DO NOT EDIT.

package abiStonkspump

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *StonksPumpVirtualFactoryMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var vx CreateVirtualLiquidityJetton
	if err := vx.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	*v = StonksPumpVirtualFactoryMessage(vx)
	return nil
}
func (v StonksPumpVirtualFactoryMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return CreateVirtualLiquidityJetton(v).MarshalTLB(c, encoder)
}
func (v *CreateVirtualLiquidityJetton) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCreateVirtualLiquidityJetton); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.Content, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if err = v.DevBuyTonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DevBuyTonAmount: %v", err)
	}
	if err = v.MaxBuyPercent.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxBuyPercent: %v", err)
	}
	if err = v.DevWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DevWallet: %v", err)
	}
	return nil
}
func (v CreateVirtualLiquidityJetton) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCreateVirtualLiquidityJetton, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.Content); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	if err = v.DevBuyTonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DevBuyTonAmount: %v", err)
	}
	if err = v.MaxBuyPercent.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxBuyPercent: %v", err)
	}
	if err = v.DevWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DevWallet: %v", err)
	}
	return nil
}
func (v CreateVirtualLiquidityJetton) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg StonksPumpVirtualFactoryMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
