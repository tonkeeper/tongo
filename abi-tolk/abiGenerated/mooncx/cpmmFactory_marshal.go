// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *CpmmFactoryCreatePool) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryCreatePool); err != nil {
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
	if err = v.FeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FeeBps: %v", err)
	}
	return nil
}
func (v CpmmFactoryCreatePool) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryCreatePool, 32); err != nil {
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
	if err = v.FeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FeeBps: %v", err)
	}
	return nil
}
func (v CpmmFactoryCreatePool) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryProvidePoolAddress) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryProvidePoolAddress); err != nil {
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
	if err = v.FeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FeeBps: %v", err)
	}
	if err = v.Initiator.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Initiator: %v", err)
	}
	return nil
}
func (v CpmmFactoryProvidePoolAddress) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryProvidePoolAddress, 32); err != nil {
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
	if err = v.FeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FeeBps: %v", err)
	}
	if err = v.Initiator.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Initiator: %v", err)
	}
	return nil
}
func (v CpmmFactoryProvidePoolAddress) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryTakePoolAddress) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryTakePoolAddress); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.XMaster.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XMaster: %v", err)
	}
	if err = v.YMaster.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YMaster: %v", err)
	}
	if err = v.FeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FeeBps: %v", err)
	}
	return nil
}
func (v CpmmFactoryTakePoolAddress) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryTakePoolAddress, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.XMaster.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XMaster: %v", err)
	}
	if err = v.YMaster.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YMaster: %v", err)
	}
	if err = v.FeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FeeBps: %v", err)
	}
	return nil
}
func (v CpmmFactoryTakePoolAddress) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryCollectProtocolFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryCollectProtocolFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	return nil
}
func (v CpmmFactoryCollectProtocolFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryCollectProtocolFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	return nil
}
func (v CpmmFactoryCollectProtocolFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryUpdatePoolWallets) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryUpdatePoolWallets); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.XWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XWallet: %v", err)
	}
	if err = v.YWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YWallet: %v", err)
	}
	return nil
}
func (v CpmmFactoryUpdatePoolWallets) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryUpdatePoolWallets, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.XWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XWallet: %v", err)
	}
	if err = v.YWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YWallet: %v", err)
	}
	return nil
}
func (v CpmmFactoryUpdatePoolWallets) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryUpgradePool) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryUpgradePool); err != nil {
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
func (v CpmmFactoryUpgradePool) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryUpgradePool, 32); err != nil {
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
func (v CpmmFactoryUpgradePool) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmFactoryUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmFactoryUpgrade); err != nil {
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
func (v CpmmFactoryUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmFactoryUpgrade, 32); err != nil {
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
func (v CpmmFactoryUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg CpmmFactoryCreatePool) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmFactoryProvidePoolAddress) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmFactoryCollectProtocolFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmFactoryUpdatePoolWallets) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmFactoryUpgradePool) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmFactoryUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
