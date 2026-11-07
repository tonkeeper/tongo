// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *CpmmPoolDeploy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolDeploy); err != nil {
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
	if v.LpAccCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpAccCode: %v", err)
	}
	if v.LpPositionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .LpPositionCode: %v", err)
	}
	return nil
}

func (v CpmmPoolDeploy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolDeploy, 32); err != nil {
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
	if err = c.AddRef(&v.LpAccCode); err != nil {
		return fmt.Errorf("failed to .LpAccCode: %v", err)
	}
	if err = c.AddRef(&v.LpPositionCode); err != nil {
		return fmt.Errorf("failed to .LpPositionCode: %v", err)
	}
	return nil
}

func (v CpmmPoolDeploy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolProvideReady) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolProvideReady); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Attempt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Attempt: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideReady) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolProvideReady, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Attempt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Attempt: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideReady) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolTakeReady) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolTakeReady); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Attempt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Attempt: %v", err)
	}
	if v.IsReady, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsReady: %v", err)
	}
	return nil
}

func (v CpmmPoolTakeReady) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolTakeReady, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Attempt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Attempt: %v", err)
	}
	if err = c.WriteBit(v.IsReady); err != nil {
		return fmt.Errorf("failed to .IsReady: %v", err)
	}
	return nil
}

func (v CpmmPoolTakeReady) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmTonSwapParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SwapAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SwapAmount: %v", err)
	}
	if err = v.Excess.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excess: %v", err)
	}
	if err = v.MinOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinOut: %v", err)
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

func (v CpmmTonSwapParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SwapAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SwapAmount: %v", err)
	}
	if err = v.Excess.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excess: %v", err)
	}
	if err = v.MinOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinOut: %v", err)
	}
	if err = v.NextFulfill.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextFulfill: %v", err)
	}
	if err = v.NextReject.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextReject: %v", err)
	}
	return nil
}

func (v CpmmTonSwapParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolSwap) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolSwap); err != nil {
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

func (v CpmmPoolSwap) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolSwap, 32); err != nil {
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

func (v CpmmPoolSwap) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SwapFulfill) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSwapFulfill); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}

func (v SwapFulfill) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSwapFulfill, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}

func (v SwapFulfill) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SwapReject) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSwapReject); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}

func (v SwapReject) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSwapReject, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}

func (v SwapReject) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolDepositLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolDepositLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.MinLpOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinLpOut: %v", err)
	}
	if v.Locked, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Locked: %v", err)
	}
	if err = v.Split.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolDepositLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolDepositLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.MinLpOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinLpOut: %v", err)
	}
	if err = c.WriteBit(v.Locked); err != nil {
		return fmt.Errorf("failed to .Locked: %v", err)
	}
	if err = v.Split.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolDepositLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolProvideLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolProvideLiquidity); err != nil {
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
	if err = v.MinLpOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinLpOut: %v", err)
	}
	if v.Locked, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Locked: %v", err)
	}
	if err = v.Split.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolProvideLiquidity, 32); err != nil {
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
	if err = v.MinLpOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinLpOut: %v", err)
	}
	if err = c.WriteBit(v.Locked); err != nil {
		return fmt.Errorf("failed to .Locked: %v", err)
	}
	if err = v.Split.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolWithdrawLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolWithdrawLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	if err = v.LpPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpPos: %v", err)
	}
	if err = v.XFeesPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XFeesPos: %v", err)
	}
	if err = v.YFeesPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YFeesPos: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolWithdrawLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	if err = v.LpPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpPos: %v", err)
	}
	if err = v.XFeesPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XFeesPos: %v", err)
	}
	if err = v.YFeesPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YFeesPos: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolWithdrawFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolWithdrawFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Provider.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Provider: %v", err)
	}
	if err = v.LpPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LpPos: %v", err)
	}
	if err = v.XFeesPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XFeesPos: %v", err)
	}
	if err = v.YFeesPos.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YFeesPos: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolWithdrawFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Provider.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Provider: %v", err)
	}
	if err = v.LpPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LpPos: %v", err)
	}
	if err = v.XFeesPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XFeesPos: %v", err)
	}
	if err = v.YFeesPos.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YFeesPos: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolCancelDeposit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolCancelDeposit); err != nil {
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
	return nil
}

func (v CpmmPoolCancelDeposit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolCancelDeposit, 32); err != nil {
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
	return nil
}

func (v CpmmPoolCancelDeposit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolBoost) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolBoost); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolBoost) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolBoost, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolBoost) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolCollectProtocolFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolCollectProtocolFees); err != nil {
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

func (v CpmmPoolCollectProtocolFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolCollectProtocolFees, 32); err != nil {
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

func (v CpmmPoolCollectProtocolFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolUpdateWallets) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolUpdateWallets); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewXWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewXWallet: %v", err)
	}
	if err = v.NewYWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewYWallet: %v", err)
	}
	return nil
}

func (v CpmmPoolUpdateWallets) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolUpdateWallets, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewXWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewXWallet: %v", err)
	}
	if err = v.NewYWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewYWallet: %v", err)
	}
	return nil
}

func (v CpmmPoolUpdateWallets) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolUpgrade); err != nil {
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

func (v CpmmPoolUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolUpgrade, 32); err != nil {
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

func (v CpmmPoolUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmJettonSwapParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Excess.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excess: %v", err)
	}
	if err = v.MinOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinOut: %v", err)
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

func (v CpmmJettonSwapParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Excess.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excess: %v", err)
	}
	if err = v.MinOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinOut: %v", err)
	}
	if err = v.NextFulfill.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextFulfill: %v", err)
	}
	if err = v.NextReject.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextReject: %v", err)
	}
	return nil
}

func (v CpmmJettonSwapParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolSwapPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolSwapPayload); err != nil {
		return err
	}
	if err = v.SwapParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SwapParams: %v", err)
	}
	return nil
}

func (v CpmmPoolSwapPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolSwapPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.SwapParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SwapParams: %v", err)
	}
	return nil
}

func (v CpmmPoolSwapPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SwapFulfillPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSwapFulfillPayload); err != nil {
		return err
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}

func (v SwapFulfillPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSwapFulfillPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}

func (v SwapFulfillPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SwapRejectPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSwapRejectPayload); err != nil {
		return err
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}

func (v SwapRejectPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSwapRejectPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}

func (v SwapRejectPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolDepositLiquidityPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolDepositLiquidityPayload); err != nil {
		return err
	}
	if err = v.MinLpOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinLpOut: %v", err)
	}
	if v.Locked, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Locked: %v", err)
	}
	if err = v.Split.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolDepositLiquidityPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolDepositLiquidityPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.MinLpOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinLpOut: %v", err)
	}
	if err = c.WriteBit(v.Locked); err != nil {
		return fmt.Errorf("failed to .Locked: %v", err)
	}
	if err = v.Split.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Split: %v", err)
	}
	return nil
}

func (v CpmmPoolDepositLiquidityPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolBoostPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolBoostPayload); err != nil {
		return err
	}
	return nil
}

func (v CpmmPoolBoostPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolBoostPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	return nil
}

func (v CpmmPoolBoostPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolProvideLiquidityRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolProvideLiquidityRefund); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidityRefund) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolProvideLiquidityRefund, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidityRefund) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolProvideLiquidityExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolProvideLiquidityExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidityExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolProvideLiquidityExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProvideLiquidityExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolWithdrawLiquidityPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolWithdrawLiquidityPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawLiquidityPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolWithdrawLiquidityPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawLiquidityPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolWithdrawFeesPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolWithdrawFeesPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawFeesPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolWithdrawFeesPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolWithdrawFeesPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolCancelDepositRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolCancelDepositRefund); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolCancelDepositRefund) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolCancelDepositRefund, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolCancelDepositRefund) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmPoolProtocolFeesPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPoolProtocolFeesPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProtocolFeesPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPoolProtocolFeesPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmPoolProtocolFeesPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *CpmmExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v CpmmExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v CpmmExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg CpmmPoolDeploy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolProvideReady) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolSwap) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SwapFulfill) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolDepositLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolProvideLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolWithdrawLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolWithdrawFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolCancelDeposit) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolBoost) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolCollectProtocolFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolUpdateWallets) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPoolUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
