// Code generated - DO NOT EDIT.

package abiStonkspump

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *AskToPresaleSell) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixAskToPresaleSell); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.JettonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonAmount: %v", err)
	}
	if err = v.TonsRecipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonsRecipient: %v", err)
	}
	if err = v.MinTonOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinTonOut: %v", err)
	}
	return nil
}
func (v AskToPresaleSell) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixAskToPresaleSell, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.JettonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonAmount: %v", err)
	}
	if err = v.TonsRecipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonsRecipient: %v", err)
	}
	if err = v.MinTonOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinTonOut: %v", err)
	}
	return nil
}
func (v AskToPresaleSell) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PresaleSellNotificationForMinter) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixPresaleSellNotificationForMinter); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.JettonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonAmount: %v", err)
	}
	if err = v.SellInitiator.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SellInitiator: %v", err)
	}
	if err = v.TonsRecipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonsRecipient: %v", err)
	}
	if err = v.MinTonOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinTonOut: %v", err)
	}
	return nil
}
func (v PresaleSellNotificationForMinter) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixPresaleSellNotificationForMinter, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.JettonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonAmount: %v", err)
	}
	if err = v.SellInitiator.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SellInitiator: %v", err)
	}
	if err = v.TonsRecipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonsRecipient: %v", err)
	}
	if err = v.MinTonOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinTonOut: %v", err)
	}
	return nil
}
func (v PresaleSellNotificationForMinter) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BuyFromPresale) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixBuyFromPresale); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonAmount: %v", err)
	}
	if err = v.MinJettonsOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinJettonsOut: %v", err)
	}
	return nil
}
func (v BuyFromPresale) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixBuyFromPresale, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonAmount: %v", err)
	}
	if err = v.MinJettonsOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinJettonsOut: %v", err)
	}
	return nil
}
func (v BuyFromPresale) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PresaleTradeEvent) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixPresaleTradeEvent); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.IsBuy, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsBuy: %v", err)
	}
	if err = v.AmountIn.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountIn: %v", err)
	}
	if err = v.AmountOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountOut: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.Leftover.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Leftover: %v", err)
	}
	return nil
}
func (v PresaleTradeEvent) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixPresaleTradeEvent, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.WriteBit(v.IsBuy); err != nil {
		return fmt.Errorf("failed to .IsBuy: %v", err)
	}
	if err = v.AmountIn.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountIn: %v", err)
	}
	if err = v.AmountOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountOut: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.Leftover.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Leftover: %v", err)
	}
	return nil
}
func (v PresaleTradeEvent) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BondingData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ProgressBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProgressBps: %v", err)
	}
	if err = v.MigrationTonTarget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MigrationTonTarget: %v", err)
	}
	if err = v.CumulativeTonReserve.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CumulativeTonReserve: %v", err)
	}
	if err = v.CumulativeJettonReserve.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CumulativeJettonReserve: %v", err)
	}
	if err = v.TonsCollected.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonsCollected: %v", err)
	}
	if err = v.MaxBuyPercent.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxBuyPercent: %v", err)
	}
	if err = v.TokensSold.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensSold: %v", err)
	}
	if err = v.PresaleOpen.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PresaleOpen: %v", err)
	}
	return nil
}
func (v BondingData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ProgressBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProgressBps: %v", err)
	}
	if err = v.MigrationTonTarget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MigrationTonTarget: %v", err)
	}
	if err = v.CumulativeTonReserve.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CumulativeTonReserve: %v", err)
	}
	if err = v.CumulativeJettonReserve.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CumulativeJettonReserve: %v", err)
	}
	if err = v.TonsCollected.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonsCollected: %v", err)
	}
	if err = v.MaxBuyPercent.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxBuyPercent: %v", err)
	}
	if err = v.TokensSold.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensSold: %v", err)
	}
	if err = v.PresaleOpen.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PresaleOpen: %v", err)
	}
	return nil
}
func (v BondingData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BondingData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.PresaleOpen.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PresaleOpen: %v", err)
	}
	if err = v.TokensSold.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TokensSold: %v", err)
	}
	if err = v.MaxBuyPercent.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MaxBuyPercent: %v", err)
	}
	if err = v.TonsCollected.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TonsCollected: %v", err)
	}
	if err = v.CumulativeJettonReserve.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CumulativeJettonReserve: %v", err)
	}
	if err = v.CumulativeTonReserve.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CumulativeTonReserve: %v", err)
	}
	if err = v.MigrationTonTarget.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MigrationTonTarget: %v", err)
	}
	if err = v.ProgressBps.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProgressBps: %v", err)
	}
	return nil
}

func (msg PresaleSellNotificationForMinter) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg BuyFromPresale) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
