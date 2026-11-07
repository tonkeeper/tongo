// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClmmFactoryCreatePool) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryCreatePool); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.XMaster.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XMaster: %v", err)
	}
	if err = v.YMaster.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YMaster: %v", err)
	}
	if err = v.BinStep.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BinStep: %v", err)
	}
	if err = v.InitBin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitBin: %v", err)
	}
	return nil
}
func (v ClmmFactoryCreatePool) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryCreatePool, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.XMaster.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XMaster: %v", err)
	}
	if err = v.YMaster.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YMaster: %v", err)
	}
	if err = v.BinStep.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BinStep: %v", err)
	}
	if err = v.InitBin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitBin: %v", err)
	}
	return nil
}
func (v ClmmFactoryCreatePool) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryUpgradePool) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryUpgradePool); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if v.NewCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	if v.NewData, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgradePool) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryUpgradePool, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.NewCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	if err = v.NewData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgradePool) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryUpdateRangeWallets) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryUpdateRangeWallets); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.RangeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeNum: %v", err)
	}
	if err = v.XWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XWallet: %v", err)
	}
	if err = v.YWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YWallet: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpdateRangeWallets) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryUpdateRangeWallets, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.RangeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeNum: %v", err)
	}
	if err = v.XWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XWallet: %v", err)
	}
	if err = v.YWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YWallet: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpdateRangeWallets) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryUpdateRangeFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryUpdateRangeFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.RangeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeNum: %v", err)
	}
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpdateRangeFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryUpdateRangeFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.RangeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeNum: %v", err)
	}
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpdateRangeFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryUpgradeRange) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryUpgradeRange); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.RangeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeNum: %v", err)
	}
	if v.NewCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	if v.NewData, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgradeRange) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryUpgradeRange, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.RangeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeNum: %v", err)
	}
	if err = v.NewCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	if err = v.NewData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgradeRange) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryCollectProtocolFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryCollectProtocolFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.RangeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeNum: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	return nil
}
func (v ClmmFactoryCollectProtocolFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryCollectProtocolFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.RangeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeNum: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	return nil
}
func (v ClmmFactoryCollectProtocolFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmFactoryUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmFactoryUpgrade); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	if v.NewData, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmFactoryUpgrade, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	if err = v.NewData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	return nil
}
func (v ClmmFactoryUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg ClmmFactoryCreatePool) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryUpgradePool) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryUpdateRangeWallets) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryUpdateRangeFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryUpgradeRange) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryCollectProtocolFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmFactoryUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
