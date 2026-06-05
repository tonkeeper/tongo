// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (v *PositionIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = PositionIncomingMessageKind(prefix)
	switch v.SumType {
	case PositionIncomingMessageKind_ItemInternalLock:
		v.ItemInternalLock = new(ItemInternalLock)
		return v.ItemInternalLock.UnmarshalTLB(c, decoder)
	case PositionIncomingMessageKind_ItemInternalUnlock:
		v.ItemInternalUnlock = new(ItemInternalUnlock)
		return v.ItemInternalUnlock.UnmarshalTLB(c, decoder)
	case PositionIncomingMessageKind_ItemWithdraw:
		v.ItemWithdraw = new(ItemWithdraw)
		return v.ItemWithdraw.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v PositionIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case PositionIncomingMessageKind_ItemInternalLock:
		if v.ItemInternalLock == nil {
			return fmt.Errorf("PositionIncomingMessage.ItemInternalLock is nil")
		}
		return v.ItemInternalLock.MarshalTLB(c, encoder)
	case PositionIncomingMessageKind_ItemInternalUnlock:
		if v.ItemInternalUnlock == nil {
			return fmt.Errorf("PositionIncomingMessage.ItemInternalUnlock is nil")
		}
		return v.ItemInternalUnlock.MarshalTLB(c, encoder)
	case PositionIncomingMessageKind_ItemWithdraw:
		if v.ItemWithdraw == nil {
			return fmt.Errorf("PositionIncomingMessage.ItemWithdraw is nil")
		}
		return v.ItemWithdraw.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown PositionIncomingMessage variant: %v", v.SumType)
	}
}

func (v PositionIncomingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *PositionExternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = PositionExternalMessageKind(prefix)
	switch v.SumType {
	case PositionExternalMessageKind_ExternalItemWithdraw:
		v.ExternalItemWithdraw = new(ExternalItemWithdraw)
		return v.ExternalItemWithdraw.UnmarshalTLB(c, decoder)
	case PositionExternalMessageKind_ExternalCronTrigger:
		v.ExternalCronTrigger = new(ExternalCronTrigger)
		return v.ExternalCronTrigger.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v PositionExternalMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case PositionExternalMessageKind_ExternalItemWithdraw:
		if v.ExternalItemWithdraw == nil {
			return fmt.Errorf("PositionExternalMessage.ExternalItemWithdraw is nil")
		}
		return v.ExternalItemWithdraw.MarshalTLB(c, encoder)
	case PositionExternalMessageKind_ExternalCronTrigger:
		if v.ExternalCronTrigger == nil {
			return fmt.Errorf("PositionExternalMessage.ExternalCronTrigger is nil")
		}
		return v.ExternalCronTrigger.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown PositionExternalMessage variant: %v", v.SumType)
	}
}

func (v PositionExternalMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetOrderDataResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Status.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Status: %v", err)
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.Minter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Minter: %v", err)
	}
	if err = v.BidJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BidJettonWallet: %v", err)
	}
	if err = v.BidJettonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BidJettonAmount: %v", err)
	}
	if err = v.BidFromUser.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BidFromUser: %v", err)
	}
	if err = v.RefFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.RefFeeTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.PendingUserReceiveAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PendingUserReceiveAmount: %v", err)
	}
	if err = v.AskJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.AskJettonMinter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonMinter: %v", err)
	}
	if v.ConditionState, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ConditionState: %v", err)
	}
	if v.ConditionCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ConditionCode: %v", err)
	}
	if err = v.UserUnlockForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserUnlockForwardParams: %v", err)
	}
	if v.UserUnlockForwardSuccessPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .UserUnlockForwardSuccessPayload: %v", err)
	}
	if v.UserFillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SafeDepositAndForwardValue: %v", err)
	}
	if err = v.ExcessesReceiver.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExcessesReceiver: %v", err)
	}
	return nil
}

func (v GetOrderDataResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Status.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Status: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = v.Minter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Minter: %v", err)
	}
	if err = v.BidJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BidJettonWallet: %v", err)
	}
	if err = v.BidJettonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BidJettonAmount: %v", err)
	}
	if err = v.BidFromUser.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BidFromUser: %v", err)
	}
	if err = v.RefFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFee: %v", err)
	}
	if err = v.RefFeeTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFeeTier: %v", err)
	}
	if err = v.PendingUserReceiveAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PendingUserReceiveAmount: %v", err)
	}
	if err = v.AskJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonWallet: %v", err)
	}
	if err = v.AskJettonMinter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonMinter: %v", err)
	}
	if err = c.AddRef(&v.ConditionState); err != nil {
		return fmt.Errorf("failed to .ConditionState: %v", err)
	}
	if err = c.AddRef(&v.ConditionCode); err != nil {
		return fmt.Errorf("failed to .ConditionCode: %v", err)
	}
	if err = v.UserUnlockForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserUnlockForwardParams: %v", err)
	}
	if err = v.UserUnlockForwardSuccessPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserUnlockForwardSuccessPayload: %v", err)
	}
	if err = c.WriteBit(v.UserFillToVault); err != nil {
		return fmt.Errorf("failed to .UserFillToVault: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SafeDepositAndForwardValue: %v", err)
	}
	if err = v.ExcessesReceiver.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExcessesReceiver: %v", err)
	}
	return nil
}

func (v GetOrderDataResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetOrderDataResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.ExcessesReceiver.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ExcessesReceiver: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .SafeDepositAndForwardValue: %v", err)
	}
	if v.UserFillToVault, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
	}
	if v.UserUnlockForwardSuccessPayload, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	}); err != nil {
		return fmt.Errorf("failed to read .UserUnlockForwardSuccessPayload: %v", err)
	}
	if v.UserUnlockForwardParams, err = tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value tlb.Coins, err error) {
		err = value.ReadFromStack(stack)
		return
	}); err != nil {
		return fmt.Errorf("failed to read .UserUnlockForwardParams: %v", err)
	}
	if v.ConditionCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .ConditionCode: %v", err)
	}
	if v.ConditionState, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .ConditionState: %v", err)
	}
	if err = v.AskJettonMinter.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AskJettonMinter: %v", err)
	}
	if err = v.AskJettonWallet.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.PendingUserReceiveAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PendingUserReceiveAmount: %v", err)
	}
	if err = v.RefFeeTier.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.RefFee.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.BidFromUser.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BidFromUser: %v", err)
	}
	if err = v.BidJettonAmount.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BidJettonAmount: %v", err)
	}
	if err = v.BidJettonWallet.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BidJettonWallet: %v", err)
	}
	if err = v.Minter.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Minter: %v", err)
	}
	if err = v.QuoteId.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.Status.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Status: %v", err)
	}
	return nil
}

func (v *GetCronInfoResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.NextCallTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NextCallTime: %v", err)
	}
	if err = v.Reward.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Reward: %v", err)
	}
	if err = v.BalanceMinusAmounts.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BalanceMinusAmounts: %v", err)
	}
	if err = v.RepeatEvery.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RepeatEvery: %v", err)
	}
	return nil
}

func (v GetCronInfoResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.NextCallTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NextCallTime: %v", err)
	}
	if err = v.Reward.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Reward: %v", err)
	}
	if err = v.BalanceMinusAmounts.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BalanceMinusAmounts: %v", err)
	}
	if err = v.RepeatEvery.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RepeatEvery: %v", err)
	}
	return nil
}

func (v GetCronInfoResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetCronInfoResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.RepeatEvery.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .RepeatEvery: %v", err)
	}
	if err = v.BalanceMinusAmounts.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .BalanceMinusAmounts: %v", err)
	}
	if err = v.Reward.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Reward: %v", err)
	}
	if err = v.NextCallTime.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .NextCallTime: %v", err)
	}
	return nil
}

func (v *ItemAdditionalFieldMore) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AskJettonMinter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonMinter: %v", err)
	}
	if err = v.OrderOwner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OrderOwner: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	return nil
}

func (v ItemAdditionalFieldMore) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AskJettonMinter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonMinter: %v", err)
	}
	if err = v.OrderOwner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OrderOwner: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	return nil
}

func (v ItemAdditionalFieldMore) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemAdditionalField) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AskJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.RefFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.RefFeeTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SafeDepositAndForwardValue: %v", err)
	}
	if v.UserFillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
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

func (v ItemAdditionalField) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AskJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonWallet: %v", err)
	}
	if err = v.RefFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFee: %v", err)
	}
	if err = v.RefFeeTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFeeTier: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SafeDepositAndForwardValue: %v", err)
	}
	if err = c.WriteBit(v.UserFillToVault); err != nil {
		return fmt.Errorf("failed to .UserFillToVault: %v", err)
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

func (v ItemAdditionalField) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemInternalUnlockExtraFields) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	return nil
}

func (v ItemInternalUnlockExtraFields) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	return nil
}

func (v ItemInternalUnlockExtraFields) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ExternalItemWithdrawPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.EscrowItem.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .EscrowItem: %v", err)
	}
	if err = v.SignatureTTL.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SignatureTTL: %v", err)
	}
	if v.WithdrawArgs, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .WithdrawArgs: %v", err)
	}
	return nil
}

func (v ExternalItemWithdrawPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.EscrowItem.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .EscrowItem: %v", err)
	}
	if err = v.SignatureTTL.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SignatureTTL: %v", err)
	}
	if err = c.AddRef(&v.WithdrawArgs); err != nil {
		return fmt.Errorf("failed to .WithdrawArgs: %v", err)
	}
	return nil
}

func (v ExternalItemWithdrawPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *EscrowWithdrawSignMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixEscrowWithdrawSignMessage); err != nil {
		return err
	}
	if err = v.SchemaHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SchemaHash: %v", err)
	}
	if err = v.Timestamp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Timestamp: %v", err)
	}
	if err = v.UserWalletAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserWalletAddress: %v", err)
	}
	if v.Domain, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Domain: %v", err)
	}
	if err = v.Payload.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}

func (v EscrowWithdrawSignMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixEscrowWithdrawSignMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.SchemaHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SchemaHash: %v", err)
	}
	if err = v.Timestamp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Timestamp: %v", err)
	}
	if err = v.UserWalletAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserWalletAddress: %v", err)
	}
	if err = c.AddRef(&v.Domain); err != nil {
		return fmt.Errorf("failed to .Domain: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}

func (v EscrowWithdrawSignMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemInternalLock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixItemInternalLock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TokenAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokenAddress: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if err = v.AdditionalField.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdditionalField: %v", err)
	}
	if err = v.LockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockArgs: %v", err)
	}
	if v.UnlockCondition, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UnlockCondition: %v", err)
	}
	return nil
}

func (v ItemInternalLock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixItemInternalLock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TokenAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokenAddress: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = v.AdditionalField.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdditionalField: %v", err)
	}
	if err = v.LockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockArgs: %v", err)
	}
	if err = c.AddRef(&v.UnlockCondition); err != nil {
		return fmt.Errorf("failed to .UnlockCondition: %v", err)
	}
	return nil
}

func (v ItemInternalLock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemInternalUnlock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixItemInternalUnlock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.FillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .FillToVault: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if err = v.Resolver.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Resolver: %v", err)
	}
	if err = v.ResolverSentJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResolverSentJettonWallet: %v", err)
	}
	if err = v.ResolverSentAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResolverSentAmount: %v", err)
	}
	if err = v.UnlockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockArgs: %v", err)
	}
	if err = v.ExtraFields.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExtraFields: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v ItemInternalUnlock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixItemInternalUnlock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.WriteBit(v.FillToVault); err != nil {
		return fmt.Errorf("failed to .FillToVault: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = v.Resolver.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Resolver: %v", err)
	}
	if err = v.ResolverSentJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResolverSentJettonWallet: %v", err)
	}
	if err = v.ResolverSentAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResolverSentAmount: %v", err)
	}
	if err = v.UnlockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockArgs: %v", err)
	}
	if err = v.ExtraFields.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExtraFields: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v ItemInternalUnlock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemWithdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixItemWithdraw); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.WithdrawArgs, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .WithdrawArgs: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v ItemWithdraw) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixItemWithdraw, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.WithdrawArgs); err != nil {
		return fmt.Errorf("failed to .WithdrawArgs: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v ItemWithdraw) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ExternalItemWithdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExternalItemWithdraw); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.SignMessage.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SignMessage: %v", err)
	}
	return nil
}

func (v ExternalItemWithdraw) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExternalItemWithdraw, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.SignMessage.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SignMessage: %v", err)
	}
	return nil
}

func (v ExternalItemWithdraw) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ExternalCronTrigger) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExternalCronTrigger); err != nil {
		return err
	}
	if err = v.RewardAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RewardAddress: %v", err)
	}
	if err = v.Salt.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Salt: %v", err)
	}
	return nil
}

func (v ExternalCronTrigger) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExternalCronTrigger, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.RewardAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RewardAddress: %v", err)
	}
	if err = v.Salt.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Salt: %v", err)
	}
	return nil
}

func (v ExternalCronTrigger) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ItemLockSuccess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixItemLockSuccess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.ForwardPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .ForwardPayload: %v", err)
	}
	return nil
}

func (v ItemLockSuccess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixItemLockSuccess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.ForwardPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardPayload: %v", err)
	}
	return nil
}

func (v ItemLockSuccess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg ExternalItemWithdraw) ToExternal(address ton.AccountID, init *tlb.StateInit) (tlb.Message, error) {
	return ton.CreateExternalMessageT(address, msg, init, tlb.VarUInteger16{})
}

func (msg ExternalCronTrigger) ToExternal(address ton.AccountID, init *tlb.StateInit) (tlb.Message, error) {
	return ton.CreateExternalMessageT(address, msg, init, tlb.VarUInteger16{})
}

func (msg ItemInternalLock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ItemInternalUnlock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ItemWithdraw) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg PositionIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
