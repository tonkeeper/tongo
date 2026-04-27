// Code generated - DO NOT EDIT.

package abiFfVault

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *AssetDeposit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixAssetDeposit); err != nil {
		return err
	}
	if v.PriceUpdateData, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .PriceUpdateData: %v", err)
	}
	return nil
}
func (v AssetDeposit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixAssetDeposit, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.PriceUpdateData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PriceUpdateData: %v", err)
	}
	return nil
}
func (v AssetDeposit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *StableDeposit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixStableDeposit); err != nil {
		return err
	}
	if v.DistributeBetweenStakers, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .DistributeBetweenStakers: %v", err)
	}
	return nil
}
func (v StableDeposit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixStableDeposit, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.WriteBit(v.DistributeBetweenStakers); err != nil {
		return fmt.Errorf("failed to .DistributeBetweenStakers: %v", err)
	}
	return nil
}
func (v StableDeposit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *StakeOperation) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixStakeOperation); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.DepositAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositAddress: %v", err)
	}
	if err = v.AssetAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetAmount: %v", err)
	}
	return nil
}
func (v StakeOperation) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixStakeOperation, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.DepositAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositAddress: %v", err)
	}
	if err = v.AssetAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetAmount: %v", err)
	}
	return nil
}
func (v StakeOperation) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnstakeExecute) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeExecute); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.PriceUpdateData, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .PriceUpdateData: %v", err)
	}
	if err = v.AmountToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountToUnstake: %v", err)
	}
	if err = v.ItemIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ItemIndex: %v", err)
	}
	return nil
}
func (v UnstakeExecute) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeExecute, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.PriceUpdateData); err != nil {
		return fmt.Errorf("failed to .PriceUpdateData: %v", err)
	}
	if err = v.AmountToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountToUnstake: %v", err)
	}
	if err = v.ItemIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ItemIndex: %v", err)
	}
	return nil
}
func (v UnstakeExecute) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnstakeOperation) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeOperation); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.AmountToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountToUnstake: %v", err)
	}
	if err = v.ItemIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ItemIndex: %v", err)
	}
	return nil
}
func (v UnstakeOperation) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeOperation, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.AmountToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountToUnstake: %v", err)
	}
	if err = v.ItemIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ItemIndex: %v", err)
	}
	return nil
}
func (v UnstakeOperation) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PositionData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.XautPriceStaked.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XautPriceStaked: %v", err)
	}
	if err = v.PreviousStablePerAssetCollected.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PreviousStablePerAssetCollected: %v", err)
	}
	return nil
}
func (v PositionData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.XautPriceStaked.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XautPriceStaked: %v", err)
	}
	if err = v.PreviousStablePerAssetCollected.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PreviousStablePerAssetCollected: %v", err)
	}
	return nil
}
func (v PositionData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UnstakeExecuteInternalCallback) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUnstakeExecuteInternalCallback); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.ItemIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ItemIndex: %v", err)
	}
	if err = v.AmountToUnstake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AmountToUnstake: %v", err)
	}
	if err = v.XautPriceActual.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XautPriceActual: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.PositionData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PositionData: %v", err)
	}
	return nil
}
func (v UnstakeExecuteInternalCallback) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUnstakeExecuteInternalCallback, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ItemIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ItemIndex: %v", err)
	}
	if err = v.AmountToUnstake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AmountToUnstake: %v", err)
	}
	if err = v.XautPriceActual.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XautPriceActual: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.PositionData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PositionData: %v", err)
	}
	return nil
}
func (v UnstakeExecuteInternalCallback) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *WithdrawJetton) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixWithdrawJetton); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.ToAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ToAddress: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}
func (v WithdrawJetton) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixWithdrawJetton, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
	}
	if err = v.ToAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ToAddress: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}
func (v WithdrawJetton) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VaultStaticData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.PYTHOracle.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PYTHOracle: %v", err)
	}
	if err = v.AssetJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetJettonWallet: %v", err)
	}
	if err = v.StableJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StableJettonWallet: %v", err)
	}
	if v.NftItemCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NftItemCode: %v", err)
	}
	if err = v.StableDecimals.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StableDecimals: %v", err)
	}
	if err = v.AssetDecimals.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetDecimals: %v", err)
	}
	return nil
}
func (v VaultStaticData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.PYTHOracle.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PYTHOracle: %v", err)
	}
	if err = v.AssetJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetJettonWallet: %v", err)
	}
	if err = v.StableJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StableJettonWallet: %v", err)
	}
	if err = c.AddRef(&v.NftItemCode); err != nil {
		return fmt.Errorf("failed to .NftItemCode: %v", err)
	}
	if err = v.StableDecimals.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StableDecimals: %v", err)
	}
	if err = v.AssetDecimals.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetDecimals: %v", err)
	}
	return nil
}
func (v VaultStaticData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VaultStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(8, PrefixVaultStorage); err != nil {
		return err
	}
	if err = v.AssetStaked.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetStaked: %v", err)
	}
	if err = v.DistributedStables.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DistributedStables: %v", err)
	}
	if err = v.MaxPriceChange.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxPriceChange: %v", err)
	}
	if err = v.NextItemIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NextItemIndex: %v", err)
	}
	if err = v.Admin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Admin: %v", err)
	}
	if err = v.AssetAmountOnVault.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetAmountOnVault: %v", err)
	}
	if err = v.StableAmountOnVault.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StableAmountOnVault: %v", err)
	}
	if err = v.StaticData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StaticData: %v", err)
	}
	return nil
}
func (v VaultStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultStorage, 8); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.AssetStaked.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetStaked: %v", err)
	}
	if err = v.DistributedStables.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DistributedStables: %v", err)
	}
	if err = v.MaxPriceChange.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxPriceChange: %v", err)
	}
	if err = v.NextItemIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextItemIndex: %v", err)
	}
	if err = v.Admin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Admin: %v", err)
	}
	if err = v.AssetAmountOnVault.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetAmountOnVault: %v", err)
	}
	if err = v.StableAmountOnVault.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StableAmountOnVault: %v", err)
	}
	if err = v.StaticData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StaticData: %v", err)
	}
	return nil
}
func (v VaultStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CollectionData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Next_item_index.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Next_item_index: %v", err)
	}
	if v.Content, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	return nil
}
func (v CollectionData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Next_item_index.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Next_item_index: %v", err)
	}
	if err = c.AddRef(&v.Content); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	return nil
}
func (v CollectionData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CollectionData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.Owner.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if v.Content, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if err = v.Next_item_index.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Next_item_index: %v", err)
	}
	return nil
}
func (v *StakingData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.TokensStaked.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensStaked: %v", err)
	}
	if err = v.DistributedStables.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DistributedStables: %v", err)
	}
	if err = v.MaxPriceChange.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxPriceChange: %v", err)
	}
	return nil
}
func (v StakingData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.TokensStaked.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensStaked: %v", err)
	}
	if err = v.DistributedStables.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DistributedStables: %v", err)
	}
	if err = v.MaxPriceChange.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxPriceChange: %v", err)
	}
	return nil
}
func (v StakingData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *StakingData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.MaxPriceChange.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MaxPriceChange: %v", err)
	}
	if err = v.DistributedStables.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .DistributedStables: %v", err)
	}
	if err = v.TokensStaked.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TokensStaked: %v", err)
	}
	return nil
}
func (v *Balance) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AssetAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetAmount: %v", err)
	}
	if err = v.StableAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StableAmount: %v", err)
	}
	return nil
}
func (v Balance) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AssetAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetAmount: %v", err)
	}
	if err = v.StableAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StableAmount: %v", err)
	}
	return nil
}
func (v Balance) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Balance) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.StableAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .StableAmount: %v", err)
	}
	if err = v.AssetAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AssetAmount: %v", err)
	}
	return nil
}

func (msg StakeOperation) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VaultStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnstakeExecute) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VaultStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnstakeOperation) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VaultStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg UnstakeExecuteInternalCallback) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VaultStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg WithdrawJetton) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VaultStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
