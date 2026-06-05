// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *VaultIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = VaultIncomingMessageKind(prefix)
	switch v.SumType {
	case VaultIncomingMessageKind_VaultLock:
		v.VaultLock = new(VaultLock)
		return v.VaultLock.UnmarshalTLB(c, decoder)
	case VaultIncomingMessageKind_VaultUnlock:
		v.VaultUnlock = new(VaultUnlock)
		return v.VaultUnlock.UnmarshalTLB(c, decoder)
	case VaultIncomingMessageKind_VaultWithdrawTokens:
		v.VaultWithdrawTokens = new(VaultWithdrawTokens)
		return v.VaultWithdrawTokens.UnmarshalTLB(c, decoder)
	case VaultIncomingMessageKind_VaultDepositTokens:
		v.VaultDepositTokens = new(VaultDepositTokens)
		return v.VaultDepositTokens.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v VaultIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case VaultIncomingMessageKind_VaultLock:
		if v.VaultLock == nil {
			return fmt.Errorf("VaultIncomingMessage.VaultLock is nil")
		}
		return v.VaultLock.MarshalTLB(c, encoder)
	case VaultIncomingMessageKind_VaultUnlock:
		if v.VaultUnlock == nil {
			return fmt.Errorf("VaultIncomingMessage.VaultUnlock is nil")
		}
		return v.VaultUnlock.MarshalTLB(c, encoder)
	case VaultIncomingMessageKind_VaultWithdrawTokens:
		if v.VaultWithdrawTokens == nil {
			return fmt.Errorf("VaultIncomingMessage.VaultWithdrawTokens is nil")
		}
		return v.VaultWithdrawTokens.MarshalTLB(c, encoder)
	case VaultIncomingMessageKind_VaultDepositTokens:
		if v.VaultDepositTokens == nil {
			return fmt.Errorf("VaultIncomingMessage.VaultDepositTokens is nil")
		}
		return v.VaultDepositTokens.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown VaultIncomingMessage variant: %v", v.SumType)
	}
}

func (v VaultIncomingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetVaultDataResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Minter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Minter: %v", err)
	}
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	return nil
}

func (v GetVaultDataResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Minter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Minter: %v", err)
	}
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
	}
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	return nil
}

func (v GetVaultDataResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetVaultDataResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.Owner.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.Balance.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.JettonWallet.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.Minter.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Minter: %v", err)
	}
	return nil
}

func (v *VaultLockAdditionalData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.TonSafeDeposit.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonSafeDeposit: %v", err)
	}
	if v.UserFillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
	}
	if err = v.AskJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.OwnerPubkey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerPubkey: %v", err)
	}
	if err = v.More.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .More: %v", err)
	}
	if err = v.UserUnlockForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserUnlockForwardParams: %v", err)
	}
	if err = v.LockForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockForwardParams: %v", err)
	}
	return nil
}

func (v VaultLockAdditionalData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.TonSafeDeposit.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonSafeDeposit: %v", err)
	}
	if err = c.WriteBit(v.UserFillToVault); err != nil {
		return fmt.Errorf("failed to .UserFillToVault: %v", err)
	}
	if err = v.AskJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonWallet: %v", err)
	}
	if err = v.OwnerPubkey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerPubkey: %v", err)
	}
	if err = v.More.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .More: %v", err)
	}
	if err = v.UserUnlockForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserUnlockForwardParams: %v", err)
	}
	if err = v.LockForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockForwardParams: %v", err)
	}
	return nil
}

func (v VaultLockAdditionalData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultUnlockExtraFields) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	return nil
}

func (v VaultUnlockExtraFields) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	return nil
}

func (v VaultUnlockExtraFields) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultLock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVaultLock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if err = v.RefFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.RefFeeTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.AdditionalData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdditionalData: %v", err)
	}
	if err = v.LockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockArgs: %v", err)
	}
	if v.UnlockCondition, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UnlockCondition: %v", err)
	}
	return nil
}

func (v VaultLock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultLock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = v.RefFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFee: %v", err)
	}
	if err = v.RefFeeTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFeeTier: %v", err)
	}
	if err = v.AdditionalData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdditionalData: %v", err)
	}
	if err = v.LockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockArgs: %v", err)
	}
	if err = c.AddRef(&v.UnlockCondition); err != nil {
		return fmt.Errorf("failed to .UnlockCondition: %v", err)
	}
	return nil
}

func (v VaultLock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultUnlock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVaultUnlock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if v.FillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .FillToVault: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if err = v.UnlockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockArgs: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ExtraFields.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExtraFields: %v", err)
	}
	if err = v.UnlockForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockForwardParams: %v", err)
	}
	return nil
}

func (v VaultUnlock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultUnlock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = c.WriteBit(v.FillToVault); err != nil {
		return fmt.Errorf("failed to .FillToVault: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = v.UnlockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockArgs: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ExtraFields.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExtraFields: %v", err)
	}
	if err = v.UnlockForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockForwardParams: %v", err)
	}
	return nil
}

func (v VaultUnlock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultWithdrawTokens) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVaultWithdrawTokens); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.WithdrawForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WithdrawForwardParams: %v", err)
	}
	return nil
}

func (v VaultWithdrawTokens) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultWithdrawTokens, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.WithdrawForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WithdrawForwardParams: %v", err)
	}
	return nil
}

func (v VaultWithdrawTokens) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultDepositTokens) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVaultDepositTokens); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ForwardTonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardTonAmount: %v", err)
	}
	if v.ForwardPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .ForwardPayload: %v", err)
	}
	return nil
}

func (v VaultDepositTokens) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultDepositTokens, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ForwardTonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardTonAmount: %v", err)
	}
	if err = v.ForwardPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardPayload: %v", err)
	}
	return nil
}

func (v VaultDepositTokens) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *VaultDepositNotification) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixVaultDepositNotification); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.ForwardPayload.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardPayload: %v", err)
	}
	return nil
}

func (v VaultDepositNotification) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixVaultDepositNotification, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.ForwardPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardPayload: %v", err)
	}
	return nil
}

func (v VaultDepositNotification) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg VaultLock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg VaultUnlock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg VaultWithdrawTokens) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg VaultDepositTokens) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg VaultIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
