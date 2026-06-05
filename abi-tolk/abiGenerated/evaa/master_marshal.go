// Code generated - DO NOT EDIT.

package abiEvaa

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *UpgradeConfigResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MasterCodeVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MasterCodeVersion: %v", err)
	}
	if err = v.UserCodeVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserCodeVersion: %v", err)
	}
	if err = v.Timeout.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Timeout: %v", err)
	}
	if err = v.UpdateTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UpdateTime: %v", err)
	}
	if err = v.FreezeTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FreezeTime: %v", err)
	}
	if v.UserCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UserCode: %v", err)
	}
	if v.NewMasterCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewMasterCode: %v", err)
	}
	if v.NewUserCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewUserCode: %v", err)
	}
	return nil
}
func (v UpgradeConfigResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MasterCodeVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MasterCodeVersion: %v", err)
	}
	if err = v.UserCodeVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserCodeVersion: %v", err)
	}
	if err = v.Timeout.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Timeout: %v", err)
	}
	if err = v.UpdateTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UpdateTime: %v", err)
	}
	if err = v.FreezeTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FreezeTime: %v", err)
	}
	if err = c.AddRef(&v.UserCode); err != nil {
		return fmt.Errorf("failed to .UserCode: %v", err)
	}
	if err = v.NewMasterCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewMasterCode: %v", err)
	}
	if err = v.NewUserCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewUserCode: %v", err)
	}
	return nil
}
func (v UpgradeConfigResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpgradeConfigResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.NewUserCode, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .NewUserCode: %v", err)
	}
	if v.NewMasterCode, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .NewMasterCode: %v", err)
	}
	if v.UserCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .UserCode: %v", err)
	}
	if err = v.FreezeTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .FreezeTime: %v", err)
	}
	if err = v.UpdateTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UpdateTime: %v", err)
	}
	if err = v.Timeout.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Timeout: %v", err)
	}
	if err = v.UserCodeVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UserCodeVersion: %v", err)
	}
	if err = v.MasterCodeVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MasterCodeVersion: %v", err)
	}
	return nil
}
func (v *AssetTrackingInfo) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.TrackingSupplyIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TrackingSupplyIndex: %v", err)
	}
	if err = v.TrackingBorrowIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TrackingBorrowIndex: %v", err)
	}
	if err = v.LastAccrual.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LastAccrual: %v", err)
	}
	return nil
}
func (v AssetTrackingInfo) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.TrackingSupplyIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TrackingSupplyIndex: %v", err)
	}
	if err = v.TrackingBorrowIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TrackingBorrowIndex: %v", err)
	}
	if err = v.LastAccrual.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LastAccrual: %v", err)
	}
	return nil
}
func (v AssetTrackingInfo) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AssetTrackingInfo) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.LastAccrual.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .LastAccrual: %v", err)
	}
	if err = v.TrackingBorrowIndex.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TrackingBorrowIndex: %v", err)
	}
	if err = v.TrackingSupplyIndex.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TrackingSupplyIndex: %v", err)
	}
	return nil
}
func (v *SbRate) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SRate.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SRate: %v", err)
	}
	if err = v.BRate.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BRate: %v", err)
	}
	return nil
}
func (v SbRate) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SRate.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SRate: %v", err)
	}
	if err = v.BRate.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BRate: %v", err)
	}
	return nil
}
func (v SbRate) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SbRate) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.BRate.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BRate: %v", err)
	}
	if err = v.SRate.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .SRate: %v", err)
	}
	return nil
}
func (v *AssetTotals) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.TotalSupply.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TotalSupply: %v", err)
	}
	if err = v.TotalBorrow.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TotalBorrow: %v", err)
	}
	return nil
}
func (v AssetTotals) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.TotalSupply.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TotalSupply: %v", err)
	}
	if err = v.TotalBorrow.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TotalBorrow: %v", err)
	}
	return nil
}
func (v AssetTotals) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AssetTotals) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.TotalBorrow.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TotalBorrow: %v", err)
	}
	if err = v.TotalSupply.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TotalSupply: %v", err)
	}
	return nil
}
