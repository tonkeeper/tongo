// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *GramboLaunch) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboLaunch); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Buyer.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Buyer: %v", err)
	}
	if err = v.VirtualTonReserve.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonTarget: %v", err)
	}
	if err = v.PrebuyAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PrebuyAmount: %v", err)
	}
	return nil
}

func (v GramboLaunch) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboLaunch, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Buyer.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Buyer: %v", err)
	}
	if err = v.VirtualTonReserve.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonTarget: %v", err)
	}
	if err = v.PrebuyAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PrebuyAmount: %v", err)
	}
	return nil
}

func (v GramboLaunch) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboSellParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MinTonOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinTonOut: %v", err)
	}
	if err = v.Referrer.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Referrer: %v", err)
	}
	if v.CustomPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboSellParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MinTonOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinTonOut: %v", err)
	}
	if err = v.Referrer.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Referrer: %v", err)
	}
	if err = v.CustomPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboSellParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboSellNotification) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboSellNotification); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.From.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .From: %v", err)
	}
	if err = v.ResponseDestination.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResponseDestination: %v", err)
	}
	if err = v.SellParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SellParams: %v", err)
	}
	return nil
}

func (v GramboSellNotification) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboSellNotification, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.From.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .From: %v", err)
	}
	if err = v.ResponseDestination.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResponseDestination: %v", err)
	}
	if err = v.SellParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SellParams: %v", err)
	}
	return nil
}

func (v GramboSellNotification) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboActivateWallet) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboActivateWallet); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v GramboActivateWallet) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboActivateWallet, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v GramboActivateWallet) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
