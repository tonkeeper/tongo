// Code generated - DO NOT EDIT.

package abiFfVault

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *UnstakeRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.AmountWantedToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountWantedToUnstake: %v", err)
	}
	return nil
}
func (v UnstakeRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.AmountWantedToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountWantedToUnstake: %v", err)
	}
	return nil
}
func (v UnstakeRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnstakeExecuteInternal) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeExecuteInternal); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.AmountToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountToUnstake: %v", err)
	}
	if err = v.XautPriceActual.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XautPriceActual: %v", err)
	}
	return nil
}
func (v UnstakeExecuteInternal) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeExecuteInternal, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.AmountToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountToUnstake: %v", err)
	}
	if err = v.XautPriceActual.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XautPriceActual: %v", err)
	}
	return nil
}
func (v UnstakeExecuteInternal) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnstakeExecuteCancel) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeExecuteCancel); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ExcessGetter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExcessGetter: %v", err)
	}
	if err = v.AmountToAdd.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountToAdd: %v", err)
	}
	return nil
}
func (v UnstakeExecuteCancel) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeExecuteCancel, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ExcessGetter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExcessGetter: %v", err)
	}
	if err = v.AmountToAdd.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountToAdd: %v", err)
	}
	return nil
}
func (v UnstakeExecuteCancel) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *NftData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Init, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Init: %v", err)
	}
	if err = v.Index.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Index: %v", err)
	}
	if err = v.Collection_address.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Collection_address: %v", err)
	}
	if err = v.Owner_address.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner_address: %v", err)
	}
	if v.Content, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	return nil
}
func (v NftData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteBit(v.Init); err != nil {
		return fmt.Errorf("failed to .Init: %v", err)
	}
	if err = v.Index.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Index: %v", err)
	}
	if err = v.Collection_address.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Collection_address: %v", err)
	}
	if err = v.Owner_address.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner_address: %v", err)
	}
	if err = v.Content.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	return nil
}
func (v NftData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *NftData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Content, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (boc.Cell, error) {
		return stack.ReadCell()
	}); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if err = v.Owner_address.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Owner_address: %v", err)
	}
	if err = v.Collection_address.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Collection_address: %v", err)
	}
	if err = v.Index.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Index: %v", err)
	}
	if v.Init, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .Init: %v", err)
	}
	return nil
}
func (v *StakePositionInfo) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.LockedAssetAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockedAssetAmount: %v", err)
	}
	if err = v.AmountWantedToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountWantedToUnstake: %v", err)
	}
	if err = v.PreviousStablePerAssetCollected.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PreviousStablePerAssetCollected: %v", err)
	}
	if err = v.PriceStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PriceStake: %v", err)
	}
	return nil
}
func (v StakePositionInfo) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.LockedAssetAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockedAssetAmount: %v", err)
	}
	if err = v.AmountWantedToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountWantedToUnstake: %v", err)
	}
	if err = v.PreviousStablePerAssetCollected.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PreviousStablePerAssetCollected: %v", err)
	}
	if err = v.PriceStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceStake: %v", err)
	}
	return nil
}
func (v StakePositionInfo) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *StakePositionInfo) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.PriceStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PriceStake: %v", err)
	}
	if err = v.PreviousStablePerAssetCollected.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PreviousStablePerAssetCollected: %v", err)
	}
	if err = v.AmountWantedToUnstake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AmountWantedToUnstake: %v", err)
	}
	if err = v.LockedAssetAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .LockedAssetAmount: %v", err)
	}
	return nil
}

func (msg UnstakeRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnstakeExecuteInternal) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnstakeExecuteCancel) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
