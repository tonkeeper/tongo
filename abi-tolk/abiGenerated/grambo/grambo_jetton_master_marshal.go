// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (v *GramboTradeFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ReferrerFee.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReferrerFee: %v", err)
	}
	if err = v.Fee1.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Fee1: %v", err)
	}
	if err = v.Fee2.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Fee2: %v", err)
	}
	return nil
}

func (v GramboTradeFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ReferrerFee.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReferrerFee: %v", err)
	}
	if err = v.Fee1.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Fee1: %v", err)
	}
	if err = v.Fee2.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Fee2: %v", err)
	}
	return nil
}

func (v GramboTradeFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboBuy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboBuy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonAmount: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if err = v.MinTokensOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinTokensOut: %v", err)
	}
	if err = v.Referrer.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Referrer: %v", err)
	}
	if v.CustomPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboBuy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboBuy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonAmount: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = v.MinTokensOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinTokensOut: %v", err)
	}
	if err = v.Referrer.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Referrer: %v", err)
	}
	if err = v.CustomPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboBuy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboProvideWalletAddress) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboProvideWalletAddress); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if v.IncludeAddress, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IncludeAddress: %v", err)
	}
	return nil
}

func (v GramboProvideWalletAddress) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboProvideWalletAddress, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = c.WriteBit(v.IncludeAddress); err != nil {
		return fmt.Errorf("failed to .IncludeAddress: %v", err)
	}
	return nil
}

func (v GramboProvideWalletAddress) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboUnlockWallet) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboUnlockWallet); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	return nil
}

func (v GramboUnlockWallet) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboUnlockWallet, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	return nil
}

func (v GramboUnlockWallet) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboWithdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboWithdraw); err != nil {
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

func (v GramboWithdraw) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboWithdraw, 32); err != nil {
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

func (v GramboWithdraw) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MasterIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = MasterIncomingMessageKind(prefix)
	switch v.SumType {
	case MasterIncomingMessageKind_GramboLaunch:
		v.GramboLaunch = new(GramboLaunch)
		return v.GramboLaunch.UnmarshalTLB(c, decoder)
	case MasterIncomingMessageKind_GramboBuy:
		v.GramboBuy = new(GramboBuy)
		return v.GramboBuy.UnmarshalTLB(c, decoder)
	case MasterIncomingMessageKind_GramboSellNotification:
		v.GramboSellNotification = new(GramboSellNotification)
		return v.GramboSellNotification.UnmarshalTLB(c, decoder)
	case MasterIncomingMessageKind_GramboProvideWalletAddress:
		v.GramboProvideWalletAddress = new(GramboProvideWalletAddress)
		return v.GramboProvideWalletAddress.UnmarshalTLB(c, decoder)
	case MasterIncomingMessageKind_GramboUnlockWallet:
		v.GramboUnlockWallet = new(GramboUnlockWallet)
		return v.GramboUnlockWallet.UnmarshalTLB(c, decoder)
	case MasterIncomingMessageKind_GramboWithdraw:
		v.GramboWithdraw = new(GramboWithdraw)
		return v.GramboWithdraw.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v MasterIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case MasterIncomingMessageKind_GramboLaunch:
		if v.GramboLaunch == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboLaunch is nil")
		}
		return v.GramboLaunch.MarshalTLB(c, encoder)
	case MasterIncomingMessageKind_GramboBuy:
		if v.GramboBuy == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboBuy is nil")
		}
		return v.GramboBuy.MarshalTLB(c, encoder)
	case MasterIncomingMessageKind_GramboSellNotification:
		if v.GramboSellNotification == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboSellNotification is nil")
		}
		return v.GramboSellNotification.MarshalTLB(c, encoder)
	case MasterIncomingMessageKind_GramboProvideWalletAddress:
		if v.GramboProvideWalletAddress == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboProvideWalletAddress is nil")
		}
		return v.GramboProvideWalletAddress.MarshalTLB(c, encoder)
	case MasterIncomingMessageKind_GramboUnlockWallet:
		if v.GramboUnlockWallet == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboUnlockWallet is nil")
		}
		return v.GramboUnlockWallet.MarshalTLB(c, encoder)
	case MasterIncomingMessageKind_GramboWithdraw:
		if v.GramboWithdraw == nil {
			return fmt.Errorf("MasterIncomingMessage.GramboWithdraw is nil")
		}
		return v.GramboWithdraw.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown MasterIncomingMessage variant: %v", v.SumType)
	}
}

func (v MasterIncomingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *JettonInternalTransfer) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixJettonInternalTransfer); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	if err = v.From.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .From: %v", err)
	}
	if err = v.ResponseAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResponseAddress: %v", err)
	}
	if err = v.ForwardTonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardTonAmount: %v", err)
	}
	if v.ForwardPayload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .ForwardPayload: %v", err)
	}
	return nil
}

func (v JettonInternalTransfer) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixJettonInternalTransfer, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	if err = v.From.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .From: %v", err)
	}
	if err = v.ResponseAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResponseAddress: %v", err)
	}
	if err = v.ForwardTonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardTonAmount: %v", err)
	}
	if err = c.AddRef(&v.ForwardPayload); err != nil {
		return fmt.Errorf("failed to .ForwardPayload: %v", err)
	}
	return nil
}

func (v JettonInternalTransfer) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboTakeWalletAddress) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboTakeWalletAddress); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.OwnerAddress, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	return nil
}

func (v GramboTakeWalletAddress) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboTakeWalletAddress, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	return nil
}

func (v GramboTakeWalletAddress) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *JettonExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixJettonExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}

func (v JettonExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixJettonExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}

func (v JettonExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MasterOutgoingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = MasterOutgoingMessageKind(prefix)
	switch v.SumType {
	case MasterOutgoingMessageKind_JettonInternalTransfer:
		v.JettonInternalTransfer = new(JettonInternalTransfer)
		return v.JettonInternalTransfer.UnmarshalTLB(c, decoder)
	case MasterOutgoingMessageKind_GramboActivateWallet:
		v.GramboActivateWallet = new(GramboActivateWallet)
		return v.GramboActivateWallet.UnmarshalTLB(c, decoder)
	case MasterOutgoingMessageKind_GramboTakeWalletAddress:
		v.GramboTakeWalletAddress = new(GramboTakeWalletAddress)
		return v.GramboTakeWalletAddress.UnmarshalTLB(c, decoder)
	case MasterOutgoingMessageKind_JettonExcess:
		v.JettonExcess = new(JettonExcess)
		return v.JettonExcess.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v MasterOutgoingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case MasterOutgoingMessageKind_JettonInternalTransfer:
		if v.JettonInternalTransfer == nil {
			return fmt.Errorf("MasterOutgoingMessage.JettonInternalTransfer is nil")
		}
		return v.JettonInternalTransfer.MarshalTLB(c, encoder)
	case MasterOutgoingMessageKind_GramboActivateWallet:
		if v.GramboActivateWallet == nil {
			return fmt.Errorf("MasterOutgoingMessage.GramboActivateWallet is nil")
		}
		return v.GramboActivateWallet.MarshalTLB(c, encoder)
	case MasterOutgoingMessageKind_GramboTakeWalletAddress:
		if v.GramboTakeWalletAddress == nil {
			return fmt.Errorf("MasterOutgoingMessage.GramboTakeWalletAddress is nil")
		}
		return v.GramboTakeWalletAddress.MarshalTLB(c, encoder)
	case MasterOutgoingMessageKind_JettonExcess:
		if v.JettonExcess == nil {
			return fmt.Errorf("MasterOutgoingMessage.JettonExcess is nil")
		}
		return v.JettonExcess.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown MasterOutgoingMessage variant: %v", v.SumType)
	}
}

func (v MasterOutgoingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboSwapEvent) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboSwapEvent); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Recipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Recipient: %v", err)
	}
	if v.IsBuy, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsBuy: %v", err)
	}
	if err = v.NewCollectedTon.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewCollectedTon: %v", err)
	}
	if err = v.TonAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonAmount: %v", err)
	}
	if err = v.TokenAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokenAmount: %v", err)
	}
	if err = v.Referrer.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Referrer: %v", err)
	}
	if err = v.Fees.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Fees: %v", err)
	}
	if v.CustomPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboSwapEvent) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboSwapEvent, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Recipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Recipient: %v", err)
	}
	if err = c.WriteBit(v.IsBuy); err != nil {
		return fmt.Errorf("failed to .IsBuy: %v", err)
	}
	if err = v.NewCollectedTon.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewCollectedTon: %v", err)
	}
	if err = v.TonAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonAmount: %v", err)
	}
	if err = v.TokenAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokenAmount: %v", err)
	}
	if err = v.Referrer.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Referrer: %v", err)
	}
	if err = v.Fees.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Fees: %v", err)
	}
	if err = v.CustomPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CustomPayload: %v", err)
	}
	return nil
}

func (v GramboSwapEvent) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboGraduateEvent) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboGraduateEvent); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TonToDex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonToDex: %v", err)
	}
	if err = v.TokensToDex.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensToDex: %v", err)
	}
	return nil
}

func (v GramboGraduateEvent) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboGraduateEvent, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TonToDex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonToDex: %v", err)
	}
	if err = v.TokensToDex.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensToDex: %v", err)
	}
	return nil
}

func (v GramboGraduateEvent) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *MasterExternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = MasterExternalMessageKind(prefix)
	switch v.SumType {
	case MasterExternalMessageKind_GramboSwapEvent:
		v.GramboSwapEvent = new(GramboSwapEvent)
		return v.GramboSwapEvent.UnmarshalTLB(c, decoder)
	case MasterExternalMessageKind_GramboGraduateEvent:
		v.GramboGraduateEvent = new(GramboGraduateEvent)
		return v.GramboGraduateEvent.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v MasterExternalMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case MasterExternalMessageKind_GramboSwapEvent:
		if v.GramboSwapEvent == nil {
			return fmt.Errorf("MasterExternalMessage.GramboSwapEvent is nil")
		}
		return v.GramboSwapEvent.MarshalTLB(c, encoder)
	case MasterExternalMessageKind_GramboGraduateEvent:
		if v.GramboGraduateEvent == nil {
			return fmt.Errorf("MasterExternalMessage.GramboGraduateEvent is nil")
		}
		return v.GramboGraduateEvent.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown MasterExternalMessage variant: %v", v.SumType)
	}
}

func (v MasterExternalMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboCurveState) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.CollectedTon.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CollectedTon: %v", err)
	}
	if err = v.TokensSold.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensSold: %v", err)
	}
	if v.Graduated, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Graduated: %v", err)
	}
	if err = v.FeeAddr1.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FeeAddr1: %v", err)
	}
	if err = v.FeeAddr2.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FeeAddr2: %v", err)
	}
	return nil
}

func (v GramboCurveState) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.CollectedTon.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CollectedTon: %v", err)
	}
	if err = v.TokensSold.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensSold: %v", err)
	}
	if err = c.WriteBit(v.Graduated); err != nil {
		return fmt.Errorf("failed to .Graduated: %v", err)
	}
	if err = v.FeeAddr1.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FeeAddr1: %v", err)
	}
	if err = v.FeeAddr2.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FeeAddr2: %v", err)
	}
	return nil
}

func (v GramboCurveState) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboDexConfig) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.StonfiRouter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StonfiRouter: %v", err)
	}
	if err = v.PtonWallet.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PtonWallet: %v", err)
	}
	return nil
}

func (v GramboDexConfig) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.StonfiRouter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StonfiRouter: %v", err)
	}
	if err = v.PtonWallet.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PtonWallet: %v", err)
	}
	return nil
}

func (v GramboDexConfig) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboMasterStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.VirtualTonReserve.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonTarget: %v", err)
	}
	if err = v.Factory.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Factory: %v", err)
	}
	if err = v.TotalSupply.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TotalSupply: %v", err)
	}
	if err = v.Seed.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seed: %v", err)
	}
	if err = v.CurveState.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CurveState: %v", err)
	}
	if v.WalletCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .WalletCode: %v", err)
	}
	if v.Metadata, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Metadata: %v", err)
	}
	if err = v.DexConfig.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DexConfig: %v", err)
	}
	return nil
}

func (v GramboMasterStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.VirtualTonReserve.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonTarget: %v", err)
	}
	if err = v.Factory.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Factory: %v", err)
	}
	if err = v.TotalSupply.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TotalSupply: %v", err)
	}
	if err = v.Seed.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seed: %v", err)
	}
	if err = v.CurveState.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CurveState: %v", err)
	}
	if err = c.AddRef(&v.WalletCode); err != nil {
		return fmt.Errorf("failed to .WalletCode: %v", err)
	}
	if err = c.AddRef(&v.Metadata); err != nil {
		return fmt.Errorf("failed to .Metadata: %v", err)
	}
	if err = v.DexConfig.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DexConfig: %v", err)
	}
	return nil
}

func (v GramboMasterStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetCurveStateResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.VirtualTonReserve.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonTarget: %v", err)
	}
	if err = v.CollectedTon.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CollectedTon: %v", err)
	}
	if err = v.TokensSold.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensSold: %v", err)
	}
	if v.Graduated, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Graduated: %v", err)
	}
	return nil
}

func (v GetCurveStateResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.VirtualTonReserve.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VirtualTonReserve: %v", err)
	}
	if err = v.TonTarget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonTarget: %v", err)
	}
	if err = v.CollectedTon.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CollectedTon: %v", err)
	}
	if err = v.TokensSold.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensSold: %v", err)
	}
	if err = c.WriteBit(v.Graduated); err != nil {
		return fmt.Errorf("failed to .Graduated: %v", err)
	}
	return nil
}

func (v GetCurveStateResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetCurveStateResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Graduated, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .Graduated: %v", err)
	}
	if err = v.TokensSold.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TokensSold: %v", err)
	}
	if err = v.CollectedTon.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .CollectedTon: %v", err)
	}
	if err = v.TonTarget.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TonTarget: %v", err)
	}
	if err = v.VirtualTonReserve.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .VirtualTonReserve: %v", err)
	}
	return nil
}

func (v *GetJettonDataResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.TotalSupply.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TotalSupply: %v", err)
	}
	if err = v.Mintable.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Mintable: %v", err)
	}
	if err = v.AdminAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdminAddress: %v", err)
	}
	if v.JettonContent, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .JettonContent: %v", err)
	}
	if v.JettonWalletCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .JettonWalletCode: %v", err)
	}
	return nil
}

func (v GetJettonDataResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.TotalSupply.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TotalSupply: %v", err)
	}
	if err = v.Mintable.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Mintable: %v", err)
	}
	if err = v.AdminAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdminAddress: %v", err)
	}
	if err = c.AddRef(&v.JettonContent); err != nil {
		return fmt.Errorf("failed to .JettonContent: %v", err)
	}
	if err = c.AddRef(&v.JettonWalletCode); err != nil {
		return fmt.Errorf("failed to .JettonWalletCode: %v", err)
	}
	return nil
}

func (v GetJettonDataResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetJettonDataResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.JettonWalletCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .JettonWalletCode: %v", err)
	}
	if v.JettonContent, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .JettonContent: %v", err)
	}
	if err = v.AdminAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .AdminAddress: %v", err)
	}
	if err = v.Mintable.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Mintable: %v", err)
	}
	if err = v.TotalSupply.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TotalSupply: %v", err)
	}
	return nil
}

func (msg GramboSwapEvent) ToExternal(address ton.AccountID, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return ton.CreateExternalMessageTWithState(address, msg, init, tlb.VarUInteger16{})
}

func (msg GramboGraduateEvent) ToExternal(address ton.AccountID, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return ton.CreateExternalMessageTWithState(address, msg, init, tlb.VarUInteger16{})
}

func (msg GramboLaunch) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboBuy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboSellNotification) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboProvideWalletAddress) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboUnlockWallet) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboWithdraw) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MasterIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboMasterStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
