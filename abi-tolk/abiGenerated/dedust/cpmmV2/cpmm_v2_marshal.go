// Code generated - DO NOT EDIT.

package abiCpmmV2

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *PoolStatus) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	vx, err := tlb.UnmarshalT[tlb.Uint2](c, decoder)
	*v = PoolStatus(vx)
	return err
}

func (v PoolStatus) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return PoolStatus.MarshalTLB(c, encoder)
}
func (v *Q120X120) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Value.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Value: %v", err)
	}
	return nil
}
func (v Q120X120) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Value.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Value: %v", err)
	}
	return nil
}
func (v Q120X120) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PoolReward) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.RemainingTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RemainingTime: %v", err)
	}
	if err = v.RemainingBudget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RemainingBudget: %v", err)
	}
	if err = v.RewardsPerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RewardsPerToken: %v", err)
	}
	if err = v.LastUpdate.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LastUpdate: %v", err)
	}
	return nil
}
func (v PoolReward) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.RemainingTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RemainingTime: %v", err)
	}
	if err = v.RemainingBudget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RemainingBudget: %v", err)
	}
	if err = v.RewardsPerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RewardsPerToken: %v", err)
	}
	if err = v.LastUpdate.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LastUpdate: %v", err)
	}
	return nil
}
func (v PoolReward) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DedustCpmmV2GetPoolData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Status.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Status: %v", err)
	}
	if v.DepositActive, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .DepositActive: %v", err)
	}
	if v.SwapActive, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .SwapActive: %v", err)
	}
	if err = v.AssetX.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetX: %v", err)
	}
	if err = v.AssetY.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetY: %v", err)
	}
	if err = v.WalletsByAssets.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WalletsByAssets: %v", err)
	}
	if err = v.AssetsByWallets.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AssetsByWallets: %v", err)
	}
	if err = v.WalletsByResolutions.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WalletsByResolutions: %v", err)
	}
	if err = v.BaseFeeBPS.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BaseFeeBPS: %v", err)
	}
	if err = v.ReserveX.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReserveX: %v", err)
	}
	if err = v.ReserveY.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReserveY: %v", err)
	}
	if err = v.Liquidity.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Liquidity: %v", err)
	}
	if err = v.ProtocolFeeX.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFeeX: %v", err)
	}
	if err = v.ProtocolFeeY.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFeeY: %v", err)
	}
	if err = v.CreatorFeeX.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CreatorFeeX: %v", err)
	}
	if err = v.CreatorFeeY.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CreatorFeeY: %v", err)
	}
	if err = v.XLPFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XLPFeePerToken: %v", err)
	}
	if err = v.YLPFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YLPFeePerToken: %v", err)
	}
	if err = v.Rewards.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Rewards: %v", err)
	}
	return nil
}
func (v DedustCpmmV2GetPoolData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Status.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Status: %v", err)
	}
	if err = c.WriteBit(v.DepositActive); err != nil {
		return fmt.Errorf("failed to .DepositActive: %v", err)
	}
	if err = c.WriteBit(v.SwapActive); err != nil {
		return fmt.Errorf("failed to .SwapActive: %v", err)
	}
	if err = v.AssetX.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetX: %v", err)
	}
	if err = v.AssetY.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetY: %v", err)
	}
	if err = v.WalletsByAssets.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WalletsByAssets: %v", err)
	}
	if err = v.AssetsByWallets.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AssetsByWallets: %v", err)
	}
	if err = v.WalletsByResolutions.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WalletsByResolutions: %v", err)
	}
	if err = v.BaseFeeBPS.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BaseFeeBPS: %v", err)
	}
	if err = v.ReserveX.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReserveX: %v", err)
	}
	if err = v.ReserveY.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReserveY: %v", err)
	}
	if err = v.Liquidity.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Liquidity: %v", err)
	}
	if err = v.ProtocolFeeX.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFeeX: %v", err)
	}
	if err = v.ProtocolFeeY.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFeeY: %v", err)
	}
	if err = v.CreatorFeeX.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CreatorFeeX: %v", err)
	}
	if err = v.CreatorFeeY.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CreatorFeeY: %v", err)
	}
	if err = v.XLPFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XLPFeePerToken: %v", err)
	}
	if err = v.YLPFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YLPFeePerToken: %v", err)
	}
	if err = v.Rewards.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Rewards: %v", err)
	}
	return nil
}
func (v DedustCpmmV2GetPoolData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DedustCpmmV2GetPoolData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Rewards, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .Rewards: %v", err)
	}
	if err = v.YLPFeePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .YLPFeePerToken: %v", err)
	}
	if err = v.XLPFeePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .XLPFeePerToken: %v", err)
	}
	if err = v.CreatorFeeY.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CreatorFeeY: %v", err)
	}
	if err = v.CreatorFeeX.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CreatorFeeX: %v", err)
	}
	if err = v.ProtocolFeeY.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProtocolFeeY: %v", err)
	}
	if err = v.ProtocolFeeX.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProtocolFeeX: %v", err)
	}
	if err = v.Liquidity.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Liquidity: %v", err)
	}
	if err = v.ReserveY.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ReserveY: %v", err)
	}
	if err = v.ReserveX.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ReserveX: %v", err)
	}
	if err = v.BaseFeeBPS.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BaseFeeBPS: %v", err)
	}
	if v.WalletsByResolutions, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .WalletsByResolutions: %v", err)
	}
	if v.AssetsByWallets, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .AssetsByWallets: %v", err)
	}
	if v.WalletsByAssets, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .WalletsByAssets: %v", err)
	}
	if err = v.AssetY.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AssetY: %v", err)
	}
	if err = v.AssetX.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AssetX: %v", err)
	}
	if v.SwapActive, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .SwapActive: %v", err)
	}
	if v.DepositActive, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .DepositActive: %v", err)
	}
	if err = v.Status.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Status: %v", err)
	}
	return nil
}
