// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClmmPayloadCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.RangeCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .RangeCode: %v", err)
	}
	if v.LpAccCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpAccCode: %v", err)
	}
	if v.LpPositionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpPositionCode: %v", err)
	}
	return nil
}
func (v ClmmPayloadCode) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.AddRef(&v.RangeCode); err != nil {
		return fmt.Errorf("failed to .RangeCode: %v", err)
	}
	if err = c.AddRef(&v.LpAccCode); err != nil {
		return fmt.Errorf("failed to .LpAccCode: %v", err)
	}
	if err = c.AddRef(&v.LpPositionCode); err != nil {
		return fmt.Errorf("failed to .LpPositionCode: %v", err)
	}
	return nil
}
func (v ClmmPayloadCode) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolDeploy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolDeploy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.InitBin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitBin: %v", err)
	}
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	if err = v.PayloadCode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PayloadCode: %v", err)
	}
	return nil
}
func (v ClmmPoolDeploy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolDeploy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.InitBin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitBin: %v", err)
	}
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	if err = v.PayloadCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PayloadCode: %v", err)
	}
	return nil
}
func (v ClmmPoolDeploy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolCollectProtocolFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolCollectProtocolFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.RangeNum.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeNum: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	return nil
}
func (v ClmmPoolCollectProtocolFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolCollectProtocolFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.RangeNum.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeNum: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	return nil
}
func (v ClmmPoolCollectProtocolFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolUpdateRangeWallets) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolUpdateRangeWallets); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
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
func (v ClmmPoolUpdateRangeWallets) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolUpdateRangeWallets, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
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
func (v ClmmPoolUpdateRangeWallets) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolUpdateRangeFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolUpdateRangeFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
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
func (v ClmmPoolUpdateRangeFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolUpdateRangeFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
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
func (v ClmmPoolUpdateRangeFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolUpgradeRange) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolUpgradeRange); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
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
func (v ClmmPoolUpgradeRange) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolUpgradeRange, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
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
func (v ClmmPoolUpgradeRange) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmPoolUpgrade); err != nil {
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
func (v ClmmPoolUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmPoolUpgrade, 32); err != nil {
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
func (v ClmmPoolUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClmmPoolStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Factory.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Factory: %v", err)
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
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	if v.RangeCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .RangeCode: %v", err)
	}
	if v.LpAccCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpAccCode: %v", err)
	}
	if v.LpPositionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpPositionCode: %v", err)
	}
	return nil
}
func (v ClmmPoolStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Factory.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Factory: %v", err)
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
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	if err = c.AddRef(&v.RangeCode); err != nil {
		return fmt.Errorf("failed to .RangeCode: %v", err)
	}
	if err = c.AddRef(&v.LpAccCode); err != nil {
		return fmt.Errorf("failed to .LpAccCode: %v", err)
	}
	if err = c.AddRef(&v.LpPositionCode); err != nil {
		return fmt.Errorf("failed to .LpPositionCode: %v", err)
	}
	return nil
}
func (v ClmmPoolStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg ClmmPoolDeploy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmPoolCollectProtocolFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmPoolUpdateRangeWallets) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmPoolUpdateRangeFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmPoolUpgradeRange) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmPoolUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
