// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *FactoryIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = FactoryIncomingMessageKind(prefix)
	switch v.SumType {
	case FactoryIncomingMessageKind_MinterInitTransfer:
		v.MinterInitTransfer = new(MinterInitTransfer)
		return v.MinterInitTransfer.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterRefundRequest:
		v.MinterRefundRequest = new(MinterRefundRequest)
		return v.MinterRefundRequest.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterInternalLock:
		v.MinterInternalLock = new(MinterInternalLock)
		return v.MinterInternalLock.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterInternalUnlock:
		v.MinterInternalUnlock = new(MinterInternalUnlock)
		return v.MinterInternalUnlock.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterInternalWithdrawTokens:
		v.MinterInternalWithdrawTokens = new(MinterInternalWithdrawTokens)
		return v.MinterInternalWithdrawTokens.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterGiveProtocolOwnership:
		v.MinterGiveProtocolOwnership = new(MinterGiveProtocolOwnership)
		return v.MinterGiveProtocolOwnership.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterUpdateProtocolTier:
		v.MinterUpdateProtocolTier = new(MinterUpdateProtocolTier)
		return v.MinterUpdateProtocolTier.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterTakeProtocolOwnership:
		v.MinterTakeProtocolOwnership = new(MinterTakeProtocolOwnership)
		return v.MinterTakeProtocolOwnership.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterCancelNextProtocolOwner:
		v.MinterCancelNextProtocolOwner = new(MinterCancelNextProtocolOwner)
		return v.MinterCancelNextProtocolOwner.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterResetGas:
		v.MinterResetGas = new(MinterResetGas)
		return v.MinterResetGas.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterLockPayload:
		v.MinterLockPayload = new(MinterLockPayload)
		return v.MinterLockPayload.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterUnlockPayload:
		v.MinterUnlockPayload = new(MinterUnlockPayload)
		return v.MinterUnlockPayload.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_MinterDepositVault:
		v.MinterDepositVault = new(MinterDepositVault)
		return v.MinterDepositVault.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v FactoryIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case FactoryIncomingMessageKind_MinterInitTransfer:
		if v.MinterInitTransfer == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterInitTransfer is nil")
		}
		return v.MinterInitTransfer.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterRefundRequest:
		if v.MinterRefundRequest == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterRefundRequest is nil")
		}
		return v.MinterRefundRequest.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterInternalLock:
		if v.MinterInternalLock == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterInternalLock is nil")
		}
		return v.MinterInternalLock.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterInternalUnlock:
		if v.MinterInternalUnlock == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterInternalUnlock is nil")
		}
		return v.MinterInternalUnlock.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterInternalWithdrawTokens:
		if v.MinterInternalWithdrawTokens == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterInternalWithdrawTokens is nil")
		}
		return v.MinterInternalWithdrawTokens.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterGiveProtocolOwnership:
		if v.MinterGiveProtocolOwnership == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterGiveProtocolOwnership is nil")
		}
		return v.MinterGiveProtocolOwnership.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterUpdateProtocolTier:
		if v.MinterUpdateProtocolTier == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterUpdateProtocolTier is nil")
		}
		return v.MinterUpdateProtocolTier.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterTakeProtocolOwnership:
		if v.MinterTakeProtocolOwnership == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterTakeProtocolOwnership is nil")
		}
		return v.MinterTakeProtocolOwnership.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterCancelNextProtocolOwner:
		if v.MinterCancelNextProtocolOwner == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterCancelNextProtocolOwner is nil")
		}
		return v.MinterCancelNextProtocolOwner.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterResetGas:
		if v.MinterResetGas == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterResetGas is nil")
		}
		return v.MinterResetGas.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterLockPayload:
		if v.MinterLockPayload == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterLockPayload is nil")
		}
		return v.MinterLockPayload.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterUnlockPayload:
		if v.MinterUnlockPayload == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterUnlockPayload is nil")
		}
		return v.MinterUnlockPayload.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_MinterDepositVault:
		if v.MinterDepositVault == nil {
			return fmt.Errorf("FactoryIncomingMessage.MinterDepositVault is nil")
		}
		return v.MinterDepositVault.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown FactoryIncomingMessage variant: %v", v.SumType)
	}
}

func (v FactoryIncomingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *EscrowData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.GlobalId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GlobalId: %v", err)
	}
	if err = v.Protocol.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Protocol: %v", err)
	}
	if err = v.ProtocolTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolTier: %v", err)
	}
	if v.ItemCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ItemCode: %v", err)
	}
	if v.VaultCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .VaultCode: %v", err)
	}
	return nil
}

func (v EscrowData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.GlobalId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GlobalId: %v", err)
	}
	if err = v.Protocol.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Protocol: %v", err)
	}
	if err = v.ProtocolTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolTier: %v", err)
	}
	if err = c.AddRef(&v.ItemCode); err != nil {
		return fmt.Errorf("failed to .ItemCode: %v", err)
	}
	if err = c.AddRef(&v.VaultCode); err != nil {
		return fmt.Errorf("failed to .VaultCode: %v", err)
	}
	return nil
}

func (v EscrowData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *EscrowData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.VaultCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .VaultCode: %v", err)
	}
	if v.ItemCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .ItemCode: %v", err)
	}
	if err = v.ProtocolTier.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProtocolTier: %v", err)
	}
	if err = v.Protocol.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Protocol: %v", err)
	}
	if err = v.GlobalId.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .GlobalId: %v", err)
	}
	return nil
}

func (v *GetVersionResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Major.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Major: %v", err)
	}
	if err = v.Minor.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Minor: %v", err)
	}
	if v.Development, err = c.ReadStringRefTail(); err != nil {
		return fmt.Errorf("failed to read .Development: %v", err)
	}
	return nil
}

func (v GetVersionResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Major.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Major: %v", err)
	}
	if err = v.Minor.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Minor: %v", err)
	}
	if err = (func() error { _, err := c.WriteStringRefTail(v.Development); return err })(); err != nil {
		return fmt.Errorf("failed to .Development: %v", err)
	}
	return nil
}

func (v GetVersionResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetVersionResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Development, err = stack.ReadStringTail(); err != nil {
		return fmt.Errorf("failed to read .Development: %v", err)
	}
	if err = v.Minor.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Minor: %v", err)
	}
	if err = v.Major.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Major: %v", err)
	}
	return nil
}

func (v *BidPaymentData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.UsingVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UsingVault: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.BidJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BidJettonWallet: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v BidPaymentData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteBit(v.UsingVault); err != nil {
		return fmt.Errorf("failed to .UsingVault: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.BidJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BidJettonWallet: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v BidPaymentData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *AskRefundData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	if v.UsingVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UsingVault: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}

func (v AskRefundData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	if err = c.WriteBit(v.UsingVault); err != nil {
		return fmt.Errorf("failed to .UsingVault: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}

func (v AskRefundData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *UserPaymentData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.User.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .User: %v", err)
	}
	if err = v.UserReceiveJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserReceiveJettonWallet: %v", err)
	}
	if err = v.UserReceiveAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserReceiveAmount: %v", err)
	}
	if v.UserDepositViaVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserDepositViaVault: %v", err)
	}
	if err = v.UserForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserForwardParams: %v", err)
	}
	if err = v.UserExcesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserExcesses: %v", err)
	}
	if err = v.UserSafeDepositAndForwardValue.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserSafeDepositAndForwardValue: %v", err)
	}
	return nil
}

func (v UserPaymentData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.User.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .User: %v", err)
	}
	if err = v.UserReceiveJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserReceiveJettonWallet: %v", err)
	}
	if err = v.UserReceiveAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserReceiveAmount: %v", err)
	}
	if err = c.WriteBit(v.UserDepositViaVault); err != nil {
		return fmt.Errorf("failed to .UserDepositViaVault: %v", err)
	}
	if err = v.UserForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserForwardParams: %v", err)
	}
	if err = v.UserExcesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserExcesses: %v", err)
	}
	if err = v.UserSafeDepositAndForwardValue.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserSafeDepositAndForwardValue: %v", err)
	}
	return nil
}

func (v UserPaymentData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterRefundRequestExtraFields) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
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
	if err = v.LockRejectDest.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockRejectDest: %v", err)
	}
	return nil
}

func (v MinterRefundRequestExtraFields) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
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
	if err = v.LockRejectDest.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockRejectDest: %v", err)
	}
	return nil
}

func (v MinterRefundRequestExtraFields) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *InternalLockExtraNested) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AskJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.OwnerPubkey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerPubkey: %v", err)
	}
	return nil
}

func (v InternalLockExtraNested) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AskJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonWallet: %v", err)
	}
	if err = v.OwnerPubkey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerPubkey: %v", err)
	}
	return nil
}

func (v InternalLockExtraNested) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *InternalLockExtra) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.RefFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.RefFeeTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SafeDepositAndForwardValue: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if v.UserFillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
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
	if err = v.Nested.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Nested: %v", err)
	}
	return nil
}

func (v InternalLockExtra) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = v.RefFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFee: %v", err)
	}
	if err = v.RefFeeTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFeeTier: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.SafeDepositAndForwardValue.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SafeDepositAndForwardValue: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = c.WriteBit(v.UserFillToVault); err != nil {
		return fmt.Errorf("failed to .UserFillToVault: %v", err)
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
	if err = v.Nested.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Nested: %v", err)
	}
	return nil
}

func (v InternalLockExtra) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *UnlockAdditionalData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if v.FillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .FillToVault: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	if err = v.UnlockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockArgs: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v UnlockAdditionalData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = c.WriteBit(v.FillToVault); err != nil {
		return fmt.Errorf("failed to .FillToVault: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	if err = v.UnlockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockArgs: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v UnlockAdditionalData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterLockPayloadExtraFieldsMore) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OrderOwner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OrderOwner: %v", err)
	}
	if err = v.AskJettonMinter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonMinter: %v", err)
	}
	return nil
}

func (v MinterLockPayloadExtraFieldsMore) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OrderOwner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OrderOwner: %v", err)
	}
	if err = v.AskJettonMinter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonMinter: %v", err)
	}
	return nil
}

func (v MinterLockPayloadExtraFieldsMore) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterLockPayloadExtraFields) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.AskJettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonWallet: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
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

func (v MinterLockPayloadExtraFields) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.AskJettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonWallet: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
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

func (v MinterLockPayloadExtraFields) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *UnlockPayloadExtraFields) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	if v.FillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .FillToVault: %v", err)
	}
	if v.RefundToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .RefundToVault: %v", err)
	}
	return nil
}

func (v UnlockPayloadExtraFields) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	if err = c.WriteBit(v.FillToVault); err != nil {
		return fmt.Errorf("failed to .FillToVault: %v", err)
	}
	if err = c.WriteBit(v.RefundToVault); err != nil {
		return fmt.Errorf("failed to .RefundToVault: %v", err)
	}
	return nil
}

func (v UnlockPayloadExtraFields) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterInitTransfer) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterInitTransfer); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.TransferType.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TransferType: %v", err)
	}
	if err = v.RefFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFee: %v", err)
	}
	if err = v.RefFeeTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefFeeTier: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.BidPaymentData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .BidPaymentData: %v", err)
	}
	if err = v.AskRefundData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskRefundData: %v", err)
	}
	if err = v.UserPaymentData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserPaymentData: %v", err)
	}
	return nil
}

func (v MinterInitTransfer) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterInitTransfer, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = v.TransferType.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TransferType: %v", err)
	}
	if err = v.RefFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFee: %v", err)
	}
	if err = v.RefFeeTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefFeeTier: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.BidPaymentData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .BidPaymentData: %v", err)
	}
	if err = v.AskRefundData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskRefundData: %v", err)
	}
	if err = v.UserPaymentData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserPaymentData: %v", err)
	}
	return nil
}

func (v MinterInitTransfer) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterRefundRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterRefundRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if v.UseVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UseVault: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.ExitCode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExitCode: %v", err)
	}
	if err = v.PrevMessage.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PrevMessage: %v", err)
	}
	if err = v.ExtraFields.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExtraFields: %v", err)
	}
	return nil
}

func (v MinterRefundRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterRefundRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = c.WriteBit(v.UseVault); err != nil {
		return fmt.Errorf("failed to .UseVault: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.ExitCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExitCode: %v", err)
	}
	if err = v.PrevMessage.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PrevMessage: %v", err)
	}
	if err = v.ExtraFields.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExtraFields: %v", err)
	}
	return nil
}

func (v MinterRefundRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterInternalLock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterInternalLock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.LockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockArgs: %v", err)
	}
	if v.UnlockCondition, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UnlockCondition: %v", err)
	}
	if err = v.AdditionalData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdditionalData: %v", err)
	}
	return nil
}

func (v MinterInternalLock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterInternalLock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.LockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockArgs: %v", err)
	}
	if err = c.AddRef(&v.UnlockCondition); err != nil {
		return fmt.Errorf("failed to .UnlockCondition: %v", err)
	}
	if err = v.AdditionalData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdditionalData: %v", err)
	}
	return nil
}

func (v MinterInternalLock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterInternalUnlock) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterInternalUnlock); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.AdditionalData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdditionalData: %v", err)
	}
	return nil
}

func (v MinterInternalUnlock) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterInternalUnlock, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.AdditionalData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdditionalData: %v", err)
	}
	return nil
}

func (v MinterInternalUnlock) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterInternalWithdrawTokens) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterInternalWithdrawTokens); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.JettonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonWallet: %v", err)
	}
	if err = v.CollectedFees.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CollectedFees: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterInternalWithdrawTokens) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterInternalWithdrawTokens, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.JettonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonWallet: %v", err)
	}
	if err = v.CollectedFees.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CollectedFees: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterInternalWithdrawTokens) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterGiveProtocolOwnership) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterGiveProtocolOwnership); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewProtocolAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewProtocolAddress: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	return nil
}

func (v MinterGiveProtocolOwnership) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterGiveProtocolOwnership, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewProtocolAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewProtocolAddress: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	return nil
}

func (v MinterGiveProtocolOwnership) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterUpdateProtocolTier) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterUpdateProtocolTier); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewTier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewTier: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	return nil
}

func (v MinterUpdateProtocolTier) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterUpdateProtocolTier, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewTier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewTier: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	return nil
}

func (v MinterUpdateProtocolTier) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterTakeProtocolOwnership) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterTakeProtocolOwnership); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	return nil
}

func (v MinterTakeProtocolOwnership) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterTakeProtocolOwnership, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	return nil
}

func (v MinterTakeProtocolOwnership) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterCancelNextProtocolOwner) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterCancelNextProtocolOwner); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	return nil
}

func (v MinterCancelNextProtocolOwner) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterCancelNextProtocolOwner, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	return nil
}

func (v MinterCancelNextProtocolOwner) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterResetGas) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterResetGas); err != nil {
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

func (v MinterResetGas) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterResetGas, 32); err != nil {
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

func (v MinterResetGas) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterLockPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterLockPayload); err != nil {
		return err
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
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.TonSafeDeposit.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonSafeDeposit: %v", err)
	}
	if v.UserFillToVault, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .UserFillToVault: %v", err)
	}
	if err = v.LockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .LockArgs: %v", err)
	}
	if v.UnlockCondition, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .UnlockCondition: %v", err)
	}
	if err = v.ExtraFields.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExtraFields: %v", err)
	}
	return nil
}

func (v MinterLockPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterLockPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
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
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.TonSafeDeposit.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonSafeDeposit: %v", err)
	}
	if err = c.WriteBit(v.UserFillToVault); err != nil {
		return fmt.Errorf("failed to .UserFillToVault: %v", err)
	}
	if err = v.LockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .LockArgs: %v", err)
	}
	if err = c.AddRef(&v.UnlockCondition); err != nil {
		return fmt.Errorf("failed to .UnlockCondition: %v", err)
	}
	if err = v.ExtraFields.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExtraFields: %v", err)
	}
	return nil
}

func (v MinterLockPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterUnlockPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterUnlockPayload); err != nil {
		return err
	}
	if err = v.QuoteId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QuoteId: %v", err)
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ExtraFields.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExtraFields: %v", err)
	}
	if err = v.UnlockArgs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockArgs: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterUnlockPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterUnlockPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QuoteId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QuoteId: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ExtraFields.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExtraFields: %v", err)
	}
	if err = v.UnlockArgs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockArgs: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterUnlockPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MinterDepositVault) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMinterDepositVault); err != nil {
		return err
	}
	if err = v.Excesses.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Excesses: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterDepositVault) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMinterDepositVault, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Excesses.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Excesses: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}

func (v MinterDepositVault) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg MinterInitTransfer) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterRefundRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterInternalLock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterInternalUnlock) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterInternalWithdrawTokens) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterGiveProtocolOwnership) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterUpdateProtocolTier) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterTakeProtocolOwnership) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterCancelNextProtocolOwner) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterResetGas) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterLockPayload) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterUnlockPayload) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MinterDepositVault) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg FactoryIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
