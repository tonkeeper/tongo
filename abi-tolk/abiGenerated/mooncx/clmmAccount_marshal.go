// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClmmAccountDepositLiquidityRecord) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmAccountDepositLiquidityRecord); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if v.XIn, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .XIn: %v", err)
	}
	if err = v.DepositBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmAccountDepositLiquidityRecord) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmAccountDepositLiquidityRecord, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = c.WriteBit(v.XIn); err != nil {
		return fmt.Errorf("failed to .XIn: %v", err)
	}
	if err = v.DepositBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmAccountDepositLiquidityRecord) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg ClmmAccountDepositLiquidityRecord) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
