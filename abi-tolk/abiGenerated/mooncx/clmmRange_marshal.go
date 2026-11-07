// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClmmRangeInitPrice) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.InitBin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitBin: %v", err)
	}
	if err = v.InitSqrtPrice.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitSqrtPrice: %v", err)
	}
	return nil
}

func (v ClmmRangeInitPrice) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.InitBin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitBin: %v", err)
	}
	if err = v.InitSqrtPrice.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitSqrtPrice: %v", err)
	}
	return nil
}

func (v ClmmRangeInitPrice) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeDeploy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeDeploy); err != nil {
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
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	if err = v.InitPrice.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitPrice: %v", err)
	}
	if v.LpAccCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpAccCode: %v", err)
	}
	if v.LpPositionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpPositionCode: %v", err)
	}
	return nil
}

func (v ClmmRangeDeploy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeDeploy, 32); err != nil {
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
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	if err = v.InitPrice.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitPrice: %v", err)
	}
	if err = c.AddRef(&v.LpAccCode); err != nil {
		return fmt.Errorf("failed to .LpAccCode: %v", err)
	}
	if err = c.AddRef(&v.LpPositionCode); err != nil {
		return fmt.Errorf("failed to .LpPositionCode: %v", err)
	}
	return nil
}

func (v ClmmRangeDeploy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmTonSwapParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SwapAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SwapAmount: %v", err)
	}
	if err = v.Excess.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excess: %v", err)
	}
	if err = v.SqrtPriceLimit.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SqrtPriceLimit: %v", err)
	}
	if v.NextFulfill, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NextFulfill: %v", err)
	}
	if v.NextReject, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NextReject: %v", err)
	}
	return nil
}

func (v ClmmTonSwapParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SwapAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SwapAmount: %v", err)
	}
	if err = v.Excess.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excess: %v", err)
	}
	if err = v.SqrtPriceLimit.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SqrtPriceLimit: %v", err)
	}
	if err = v.NextFulfill.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextFulfill: %v", err)
	}
	if err = v.NextReject.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextReject: %v", err)
	}
	return nil
}

func (v ClmmTonSwapParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeSwap) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeSwap); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.SwapParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SwapParams: %v", err)
	}
	return nil
}

func (v ClmmRangeSwap) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeSwap, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.SwapParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SwapParams: %v", err)
	}
	return nil
}

func (v ClmmRangeSwap) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeDepositLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeDepositLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.DepositAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositAmount: %v", err)
	}
	if v.IsSingle, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsSingle: %v", err)
	}
	if err = v.DepositBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeDepositLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.DepositAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositAmount: %v", err)
	}
	if err = c.WriteBit(v.IsSingle); err != nil {
		return fmt.Errorf("failed to .IsSingle: %v", err)
	}
	if err = v.DepositBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeProvideLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeProvideLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	if err = v.XAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XAmount: %v", err)
	}
	if err = v.YAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YAmount: %v", err)
	}
	if err = v.DepositBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeProvideLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeProvideLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	if err = v.XAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XAmount: %v", err)
	}
	if err = v.YAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YAmount: %v", err)
	}
	if err = v.DepositBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeProvideLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeWithdrawFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeWithdrawFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	if err = v.NftIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NftIndex: %v", err)
	}
	if err = v.WithdrawBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WithdrawBins: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeWithdrawFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	if err = v.NftIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NftIndex: %v", err)
	}
	if err = v.WithdrawBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WithdrawBins: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeWithdrawLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeWithdrawLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	if err = v.NftIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NftIndex: %v", err)
	}
	if err = v.WithdrawBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WithdrawBins: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeWithdrawLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	if err = v.NftIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NftIndex: %v", err)
	}
	if err = v.WithdrawBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WithdrawBins: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeCancelDeposit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeCancelDeposit); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.XAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XAmount: %v", err)
	}
	if err = v.YAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YAmount: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	return nil
}

func (v ClmmRangeCancelDeposit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeCancelDeposit, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.XAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XAmount: %v", err)
	}
	if err = v.YAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YAmount: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	return nil
}

func (v ClmmRangeCancelDeposit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeCollectProtocolFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeCollectProtocolFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	return nil
}

func (v ClmmRangeCollectProtocolFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeCollectProtocolFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	return nil
}

func (v ClmmRangeCollectProtocolFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeUpdateFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeUpdateFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	return nil
}

func (v ClmmRangeUpdateFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeUpdateFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	return nil
}

func (v ClmmRangeUpdateFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeUpdateWallets) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeUpdateWallets); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.XWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XWallet: %v", err)
	}
	if err = v.YWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YWallet: %v", err)
	}
	return nil
}

func (v ClmmRangeUpdateWallets) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeUpdateWallets, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.XWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XWallet: %v", err)
	}
	if err = v.YWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YWallet: %v", err)
	}
	return nil
}

func (v ClmmRangeUpdateWallets) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeUpgrade); err != nil {
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

func (v ClmmRangeUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeUpgrade, 32); err != nil {
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

func (v ClmmRangeUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmJettonSwapParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Excess.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excess: %v", err)
	}
	if err = v.SqrtPriceLimit.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SqrtPriceLimit: %v", err)
	}
	if v.NextFulfill, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NextFulfill: %v", err)
	}
	if v.NextReject, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NextReject: %v", err)
	}
	return nil
}

func (v ClmmJettonSwapParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Excess.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excess: %v", err)
	}
	if err = v.SqrtPriceLimit.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SqrtPriceLimit: %v", err)
	}
	if err = v.NextFulfill.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextFulfill: %v", err)
	}
	if err = v.NextReject.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextReject: %v", err)
	}
	return nil
}

func (v ClmmJettonSwapParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeSwapPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeSwapPayload); err != nil {
		return err
	}
	if err = v.SwapParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SwapParams: %v", err)
	}
	return nil
}

func (v ClmmRangeSwapPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeSwapPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.SwapParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SwapParams: %v", err)
	}
	return nil
}

func (v ClmmRangeSwapPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeDepositLiquidityPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeDepositLiquidityPayload); err != nil {
		return err
	}
	if v.IsSingle, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsSingle: %v", err)
	}
	if err = v.DepositBins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeDepositLiquidityPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.WriteBit(v.IsSingle); err != nil {
		return fmt.Errorf("failed to .IsSingle: %v", err)
	}
	if err = v.DepositBins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositBins: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeDepositLiquidityExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeDepositLiquidityExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeDepositLiquidityExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeDepositLiquidityRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeDepositLiquidityRefund); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityRefund) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeDepositLiquidityRefund, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeDepositLiquidityRefund) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeWithdrawLiquidityPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeWithdrawLiquidityPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawLiquidityPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeWithdrawLiquidityPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawLiquidityPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeWithdrawFeesPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeWithdrawFeesPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawFeesPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeWithdrawFeesPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeWithdrawFeesPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeCancelDepositRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeCancelDepositRefund); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeCancelDepositRefund) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeCancelDepositRefund, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeCancelDepositRefund) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeProtocolFeesPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmRangeProtocolFeesPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeProtocolFeesPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmRangeProtocolFeesPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmRangeProtocolFeesPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClmmExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v ClmmExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClmmExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v ClmmExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeState) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.CurrentBin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CurrentBin: %v", err)
	}
	if err = v.CurrentSqrtP.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CurrentSqrtP: %v", err)
	}
	if err = v.XProtocolFees.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XProtocolFees: %v", err)
	}
	if err = v.YProtocolFees.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YProtocolFees: %v", err)
	}
	return nil
}

func (v ClmmRangeState) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.CurrentBin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CurrentBin: %v", err)
	}
	if err = v.CurrentSqrtP.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CurrentSqrtP: %v", err)
	}
	if err = v.XProtocolFees.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XProtocolFees: %v", err)
	}
	if err = v.YProtocolFees.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YProtocolFees: %v", err)
	}
	return nil
}

func (v ClmmRangeState) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ClmmRangeStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
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
	if err = v.BinStep.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BinStep: %v", err)
	}
	if err = v.LpFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpFee: %v", err)
	}
	if err = v.ProtocolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolFee: %v", err)
	}
	if err = v.PositionIndex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PositionIndex: %v", err)
	}
	if err = v.RangeState.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RangeState: %v", err)
	}
	if err = v.Bins.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Bins: %v", err)
	}
	if v.LpAccountCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpAccountCode: %v", err)
	}
	if v.LpPositionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpPositionCode: %v", err)
	}
	return nil
}

func (v ClmmRangeStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
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
	if err = v.BinStep.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BinStep: %v", err)
	}
	if err = v.LpFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpFee: %v", err)
	}
	if err = v.ProtocolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolFee: %v", err)
	}
	if err = v.PositionIndex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PositionIndex: %v", err)
	}
	if err = v.RangeState.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RangeState: %v", err)
	}
	if err = v.Bins.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Bins: %v", err)
	}
	if err = c.AddRef(&v.LpAccountCode); err != nil {
		return fmt.Errorf("failed to .LpAccountCode: %v", err)
	}
	if err = c.AddRef(&v.LpPositionCode); err != nil {
		return fmt.Errorf("failed to .LpPositionCode: %v", err)
	}
	return nil
}

func (v ClmmRangeStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg ClmmRangeDeploy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeSwap) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeDepositLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeProvideLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeWithdrawFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeWithdrawLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeCancelDeposit) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeCollectProtocolFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeUpdateFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeUpdateWallets) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ClmmRangeUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
