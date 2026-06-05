// Code generated - DO NOT EDIT.

package abiEvaa

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *UserScData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.CodeVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CodeVersion: %v", err)
	}
	if err = v.MasterAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MasterAddress: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if v.UserPrincipals, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UserPrincipals: %v", err)
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if v.UserRewards, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .UserRewards: %v", err)
	}
	if v.BackupCell1, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .BackupCell1: %v", err)
	}
	if v.BackupCell2, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .BackupCell2: %v", err)
	}
	return nil
}
func (v UserScData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.CodeVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CodeVersion: %v", err)
	}
	if err = v.MasterAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MasterAddress: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = c.AddRef(&v.UserPrincipals); err != nil {
		return fmt.Errorf("failed to .UserPrincipals: %v", err)
	}
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.UserRewards.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserRewards: %v", err)
	}
	if err = v.BackupCell1.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BackupCell1: %v", err)
	}
	if err = v.BackupCell2.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BackupCell2: %v", err)
	}
	return nil
}
func (v UserScData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UserScData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.BackupCell2, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .BackupCell2: %v", err)
	}
	if v.BackupCell1, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .BackupCell1: %v", err)
	}
	if v.UserRewards, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .UserRewards: %v", err)
	}
	if err = v.State.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if v.UserPrincipals, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .UserPrincipals: %v", err)
	}
	if err = v.OwnerAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.MasterAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MasterAddress: %v", err)
	}
	if err = v.CodeVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CodeVersion: %v", err)
	}
	return nil
}
func (v *AggregatedBalances) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Supply.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Supply: %v", err)
	}
	if err = v.Borrow.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Borrow: %v", err)
	}
	return nil
}
func (v AggregatedBalances) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Supply.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Supply: %v", err)
	}
	if err = v.Borrow.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Borrow: %v", err)
	}
	return nil
}
func (v AggregatedBalances) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AggregatedBalances) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.Borrow.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Borrow: %v", err)
	}
	if err = v.Supply.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Supply: %v", err)
	}
	return nil
}
