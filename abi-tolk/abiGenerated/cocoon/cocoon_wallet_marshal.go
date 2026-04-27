// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (v *OwnerWalletSendMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerWalletSendMessage); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Mode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Mode: %v", err)
	}
	if v.Body, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Body: %v", err)
	}
	return nil
}
func (v OwnerWalletSendMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerWalletSendMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Mode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Mode: %v", err)
	}
	if err = c.AddRef(&v.Body); err != nil {
		return fmt.Errorf("failed to .Body: %v", err)
	}
	return nil
}
func (v OwnerWalletSendMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *TextCommand) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixTextCommand); err != nil {
		return err
	}
	if err = v.Action.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Action: %v", err)
	}
	return nil
}
func (v TextCommand) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixTextCommand, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Action.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Action: %v", err)
	}
	return nil
}
func (v TextCommand) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ForwardMsg) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Mode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Mode: %v", err)
	}
	if v.Msg, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	return nil
}
func (v ForwardMsg) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Mode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Mode: %v", err)
	}
	if err = c.AddRef(&v.Msg); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	return nil
}
func (v ForwardMsg) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExternalSignedMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SubwalletId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SubwalletId: %v", err)
	}
	if err = v.ValidUntil.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidUntil: %v", err)
	}
	if err = v.MsgSeqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MsgSeqno: %v", err)
	}
	if err = v.Forward.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Forward: %v", err)
	}
	return nil
}
func (v ExternalSignedMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SubwalletId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SubwalletId: %v", err)
	}
	if err = v.ValidUntil.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidUntil: %v", err)
	}
	if err = v.MsgSeqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MsgSeqno: %v", err)
	}
	if err = v.Forward.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Forward: %v", err)
	}
	return nil
}
func (v ExternalSignedMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *WalletExternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.Message.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Message: %v", err)
	}
	return nil
}
func (v WalletExternalMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.Message.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Message: %v", err)
	}
	return nil
}
func (v WalletExternalMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *WalletStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if err = v.SubwalletId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SubwalletId: %v", err)
	}
	if err = v.PublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PublicKey: %v", err)
	}
	if err = v.Status.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Status: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	return nil
}
func (v WalletStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = v.SubwalletId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SubwalletId: %v", err)
	}
	if err = v.PublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PublicKey: %v", err)
	}
	if err = v.Status.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Status: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	return nil
}
func (v WalletStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg WalletExternalMessage) ToExternal(address ton.AccountID, init *tlb.StateInitT[*WalletStorage]) (tlb.Message, error) {
	return ton.CreateExternalMessageTWithState(address, msg, init, tlb.VarUInteger16{})
}

func (msg OwnerWalletSendMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*WalletStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg TextCommand) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*WalletStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
