// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *CpmmPositionContent) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Lp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Lp: %v", err)
	}
	if err = v.XFeesLp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XFeesLp: %v", err)
	}
	if err = v.YFeesLp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YFeesLp: %v", err)
	}
	return nil
}
func (v CpmmPositionContent) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Lp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Lp: %v", err)
	}
	if err = v.XFeesLp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XFeesLp: %v", err)
	}
	if err = v.YFeesLp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YFeesLp: %v", err)
	}
	return nil
}
func (v CpmmPositionContent) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionDeploy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPositionDeploy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Content.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if v.Locked, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Locked: %v", err)
	}
	return nil
}
func (v CpmmPositionDeploy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPositionDeploy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Content.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	if err = c.WriteBit(v.Locked); err != nil {
		return fmt.Errorf("failed to .Locked: %v", err)
	}
	return nil
}
func (v CpmmPositionDeploy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionWithdrawLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPositionWithdrawLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPositionWithdrawLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionWithdrawFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPositionWithdrawFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPositionWithdrawFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionWithdrawFeesUpdate) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPositionWithdrawFeesUpdate); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.XFeesLp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .XFeesLp: %v", err)
	}
	if err = v.YFeesLp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .YFeesLp: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawFeesUpdate) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPositionWithdrawFeesUpdate, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.XFeesLp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .XFeesLp: %v", err)
	}
	if err = v.YFeesLp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .YFeesLp: %v", err)
	}
	return nil
}
func (v CpmmPositionWithdrawFeesUpdate) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionLockLiquidity) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCpmmPositionLockLiquidity); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v CpmmPositionLockLiquidity) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCpmmPositionLockLiquidity, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v CpmmPositionLockLiquidity) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CpmmPositionStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Pool.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Pool: %v", err)
	}
	if err = v.User.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .User: %v", err)
	}
	if err = v.Data.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Data: %v", err)
	}
	if v.Frozen, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Frozen: %v", err)
	}
	if v.Locked, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Locked: %v", err)
	}
	return nil
}
func (v CpmmPositionStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Pool.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Pool: %v", err)
	}
	if err = v.User.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .User: %v", err)
	}
	if err = v.Data.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Data: %v", err)
	}
	if err = c.WriteBit(v.Frozen); err != nil {
		return fmt.Errorf("failed to .Frozen: %v", err)
	}
	if err = c.WriteBit(v.Locked); err != nil {
		return fmt.Errorf("failed to .Locked: %v", err)
	}
	return nil
}
func (v CpmmPositionStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg CpmmPositionDeploy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPositionWithdrawLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPositionWithdrawFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPositionWithdrawFeesUpdate) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CpmmPositionLockLiquidity) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
