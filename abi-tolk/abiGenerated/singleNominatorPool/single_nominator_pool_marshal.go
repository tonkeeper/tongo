// Code generated - DO NOT EDIT.

package abiSingleNominatorPool

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *Withdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixWithdraw); err != nil {
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
func (v Withdraw) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixWithdraw, 32); err != nil {
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
func (v Withdraw) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangeValidatorAddress) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeValidatorAddress); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewValidatorAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewValidatorAddress: %v", err)
	}
	return nil
}
func (v ChangeValidatorAddress) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeValidatorAddress, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewValidatorAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewValidatorAddress: %v", err)
	}
	return nil
}
func (v ChangeValidatorAddress) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SendRawMsg) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSendRawMsg); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Mode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Mode: %v", err)
	}
	if v.Msg, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	return nil
}
func (v SendRawMsg) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSendRawMsg, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Mode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Mode: %v", err)
	}
	if err = c.AddRef(&v.Msg); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	return nil
}
func (v SendRawMsg) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Upgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpgrade); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.Code, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Code: %v", err)
	}
	return nil
}
func (v Upgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpgrade, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.Code); err != nil {
		return fmt.Errorf("failed to .Code: %v", err)
	}
	return nil
}
func (v Upgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Roles) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.Validator.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Validator: %v", err)
	}
	return nil
}
func (v Roles) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.Validator.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Validator: %v", err)
	}
	return nil
}
func (v Roles) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Roles) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.Validator.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Validator: %v", err)
	}
	if err = v.Owner.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	return nil
}
func (v *PoolConfig) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MinStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinStake: %v", err)
	}
	if err = v.DepositFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DepositFee: %v", err)
	}
	if err = v.WithdrawFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WithdrawFee: %v", err)
	}
	if err = v.PoolFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PoolFee: %v", err)
	}
	if err = v.ReceiptPrice.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReceiptPrice: %v", err)
	}
	return nil
}
func (v PoolConfig) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MinStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinStake: %v", err)
	}
	if err = v.DepositFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DepositFee: %v", err)
	}
	if err = v.WithdrawFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WithdrawFee: %v", err)
	}
	if err = v.PoolFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PoolFee: %v", err)
	}
	if err = v.ReceiptPrice.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReceiptPrice: %v", err)
	}
	return nil
}
func (v PoolConfig) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PoolConfig) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.ReceiptPrice.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ReceiptPrice: %v", err)
	}
	if err = v.PoolFee.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PoolFee: %v", err)
	}
	if err = v.WithdrawFee.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .WithdrawFee: %v", err)
	}
	if err = v.DepositFee.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .DepositFee: %v", err)
	}
	if err = v.MinStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinStake: %v", err)
	}
	return nil
}
func (v *PoolData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.NominatorsCount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NominatorsCount: %v", err)
	}
	if err = v.StakeAmountSent.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StakeAmountSent: %v", err)
	}
	if err = v.ValidatorAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorAmount: %v", err)
	}
	if err = v.PoolConfig.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PoolConfig: %v", err)
	}
	if v.Nominators, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Nominators: %v", err)
	}
	if v.WithdrawRequests, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .WithdrawRequests: %v", err)
	}
	if err = v.StakeAt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StakeAt: %v", err)
	}
	if err = v.SavedValidatorSetHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SavedValidatorSetHash: %v", err)
	}
	if err = v.ValidatorSetChangesCount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorSetChangesCount: %v", err)
	}
	if err = v.ValidatorSetChangeTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidatorSetChangeTime: %v", err)
	}
	if err = v.StakeHeldFor.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StakeHeldFor: %v", err)
	}
	if v.ConfigProposalVotings, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .ConfigProposalVotings: %v", err)
	}
	return nil
}
func (v PoolData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.NominatorsCount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NominatorsCount: %v", err)
	}
	if err = v.StakeAmountSent.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StakeAmountSent: %v", err)
	}
	if err = v.ValidatorAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorAmount: %v", err)
	}
	if err = v.PoolConfig.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PoolConfig: %v", err)
	}
	if err = v.Nominators.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Nominators: %v", err)
	}
	if err = v.WithdrawRequests.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WithdrawRequests: %v", err)
	}
	if err = v.StakeAt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StakeAt: %v", err)
	}
	if err = v.SavedValidatorSetHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SavedValidatorSetHash: %v", err)
	}
	if err = v.ValidatorSetChangesCount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorSetChangesCount: %v", err)
	}
	if err = v.ValidatorSetChangeTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidatorSetChangeTime: %v", err)
	}
	if err = v.StakeHeldFor.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StakeHeldFor: %v", err)
	}
	if err = v.ConfigProposalVotings.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ConfigProposalVotings: %v", err)
	}
	return nil
}
func (v PoolData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *PoolData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.ConfigProposalVotings, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .ConfigProposalVotings: %v", err)
	}
	if err = v.StakeHeldFor.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .StakeHeldFor: %v", err)
	}
	if err = v.ValidatorSetChangeTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ValidatorSetChangeTime: %v", err)
	}
	if err = v.ValidatorSetChangesCount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ValidatorSetChangesCount: %v", err)
	}
	if err = v.SavedValidatorSetHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .SavedValidatorSetHash: %v", err)
	}
	if err = v.StakeAt.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .StakeAt: %v", err)
	}
	if v.WithdrawRequests, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .WithdrawRequests: %v", err)
	}
	if v.Nominators, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .Nominators: %v", err)
	}
	if err = v.PoolConfig.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PoolConfig: %v", err)
	}
	if err = v.ValidatorAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ValidatorAmount: %v", err)
	}
	if err = v.StakeAmountSent.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .StakeAmountSent: %v", err)
	}
	if err = v.NominatorsCount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .NominatorsCount: %v", err)
	}
	if err = v.State.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	return nil
}

func (msg Withdraw) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ChangeValidatorAddress) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SendRawMsg) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg Upgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
