// Code generated - DO NOT EDIT.

package abiWalletTg

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (v *Storage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(8, PrefixStorage); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if err = v.SubwalletId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SubwalletId: %v", err)
	}
	if err = v.PublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PublicKey: %v", err)
	}
	return nil
}
func (v Storage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixStorage, 8); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = v.SubwalletId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SubwalletId: %v", err)
	}
	if err = v.PublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PublicKey: %v", err)
	}
	return nil
}
func (v Storage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *Revision) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	return (*tlb.Uint8)(v).UnmarshalTLB(c, decoder)
}

func (v Revision) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return tlb.Uint8(v).MarshalTLB(c, encoder)
}

func (v *Revision) ReadFromStack(stack *tlb.VmStack) error {
	return (*tlb.Uint8)(v).ReadFromStack(stack)
}
func (v *MessageToSend) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SendMode.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendMode: %v", err)
	}
	if v.MessageCell, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .MessageCell: %v", err)
	}
	return nil
}
func (v MessageToSend) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SendMode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendMode: %v", err)
	}
	if err = c.AddRef(&v.MessageCell); err != nil {
		return fmt.Errorf("failed to .MessageCell: %v", err)
	}
	return nil
}
func (v MessageToSend) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RawArrayOfMessagesToSend) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	vx, err := tlb.LoadTolkArray[MessageToSend](c, decoder)
	if err != nil {
		return err
	}
	*v = RawArrayOfMessagesToSend(vx)
	return nil
}
func (v RawArrayOfMessagesToSend) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return tlb.StoreTolkArray(c, encoder, []MessageToSend([]MessageToSend(v)))
}
func (v *KeyRotationProofPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Tag.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Tag: %v", err)
	}
	if err = v.WalletWorkchain.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WalletWorkchain: %v", err)
	}
	if err = v.WalletAddrHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WalletAddrHash: %v", err)
	}
	return nil
}
func (v KeyRotationProofPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Tag.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Tag: %v", err)
	}
	if err = v.WalletWorkchain.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WalletWorkchain: %v", err)
	}
	if err = v.WalletAddrHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WalletAddrHash: %v", err)
	}
	return nil
}
func (v KeyRotationProofPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SignedInternalRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.Request.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Request: %v", err)
	}
	return nil
}
func (v SignedInternalRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.Request.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Request: %v", err)
	}
	return nil
}
func (v SignedInternalRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SignedExternalRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.Request.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Request: %v", err)
	}
	return nil
}
func (v SignedExternalRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.Request.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Request: %v", err)
	}
	return nil
}
func (v SignedExternalRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AllowedInternalMsg) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = AllowedInternalMsgKind(prefix)
	switch v.SumType {
	case AllowedInternalMsgKind_SendOneMessageRequestI:
		v.SendOneMessageRequestI = new(SendOneMessageRequestI)
		return v.SendOneMessageRequestI.UnmarshalTLB(c, decoder)
	case AllowedInternalMsgKind_SendBulkMessagesRequestI:
		v.SendBulkMessagesRequestI = new(SendBulkMessagesRequestI)
		return v.SendBulkMessagesRequestI.UnmarshalTLB(c, decoder)
	case AllowedInternalMsgKind_ChangePublicKeyRequestI:
		v.ChangePublicKeyRequestI = new(ChangePublicKeyRequestI)
		return v.ChangePublicKeyRequestI.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}
func (v AllowedInternalMsg) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case AllowedInternalMsgKind_SendOneMessageRequestI:
		if v.SendOneMessageRequestI == nil {
			return fmt.Errorf("AllowedInternalMsg.SendOneMessageRequestI is nil")
		}
		return v.SendOneMessageRequestI.MarshalTLB(c, encoder)
	case AllowedInternalMsgKind_SendBulkMessagesRequestI:
		if v.SendBulkMessagesRequestI == nil {
			return fmt.Errorf("AllowedInternalMsg.SendBulkMessagesRequestI is nil")
		}
		return v.SendBulkMessagesRequestI.MarshalTLB(c, encoder)
	case AllowedInternalMsgKind_ChangePublicKeyRequestI:
		if v.ChangePublicKeyRequestI == nil {
			return fmt.Errorf("AllowedInternalMsg.ChangePublicKeyRequestI is nil")
		}
		return v.ChangePublicKeyRequestI.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown AllowedInternalMsg variant: %v", v.SumType)
	}
}
func (v AllowedInternalMsg) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *AllowedExternalMsg) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = AllowedExternalMsgKind(prefix)
	switch v.SumType {
	case AllowedExternalMsgKind_SendOneMessageRequestE:
		v.SendOneMessageRequestE = new(SendOneMessageRequestE)
		return v.SendOneMessageRequestE.UnmarshalTLB(c, decoder)
	case AllowedExternalMsgKind_SendBulkMessagesRequestE:
		v.SendBulkMessagesRequestE = new(SendBulkMessagesRequestE)
		return v.SendBulkMessagesRequestE.UnmarshalTLB(c, decoder)
	case AllowedExternalMsgKind_ChangePublicKeyRequestE:
		v.ChangePublicKeyRequestE = new(ChangePublicKeyRequestE)
		return v.ChangePublicKeyRequestE.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}
func (v AllowedExternalMsg) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case AllowedExternalMsgKind_SendOneMessageRequestE:
		if v.SendOneMessageRequestE == nil {
			return fmt.Errorf("AllowedExternalMsg.SendOneMessageRequestE is nil")
		}
		return v.SendOneMessageRequestE.MarshalTLB(c, encoder)
	case AllowedExternalMsgKind_SendBulkMessagesRequestE:
		if v.SendBulkMessagesRequestE == nil {
			return fmt.Errorf("AllowedExternalMsg.SendBulkMessagesRequestE is nil")
		}
		return v.SendBulkMessagesRequestE.MarshalTLB(c, encoder)
	case AllowedExternalMsgKind_ChangePublicKeyRequestE:
		if v.ChangePublicKeyRequestE == nil {
			return fmt.Errorf("AllowedExternalMsg.ChangePublicKeyRequestE is nil")
		}
		return v.ChangePublicKeyRequestE.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown AllowedExternalMsg variant: %v", v.SumType)
	}
}
func (v AllowedExternalMsg) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SeqnoHeader) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SubwalletId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SubwalletId: %v", err)
	}
	if err = v.ValidUntil.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidUntil: %v", err)
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	return nil
}
func (v SeqnoHeader) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.SubwalletId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SubwalletId: %v", err)
	}
	if err = v.ValidUntil.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidUntil: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	return nil
}
func (v SeqnoHeader) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SendOneMessageRequestI) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSendOneMessageRequestI); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.Msg.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	return nil
}
func (v SendOneMessageRequestI) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSendOneMessageRequestI, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.Msg.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	return nil
}
func (v SendOneMessageRequestI) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SendOneMessageRequestE) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSendOneMessageRequestE); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.Msg.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	return nil
}
func (v SendOneMessageRequestE) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSendOneMessageRequestE, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.Msg.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	return nil
}
func (v SendOneMessageRequestE) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SendBulkMessagesRequestI) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSendBulkMessagesRequestI); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.MsgArr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MsgArr: %v", err)
	}
	return nil
}
func (v SendBulkMessagesRequestI) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSendBulkMessagesRequestI, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.MsgArr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MsgArr: %v", err)
	}
	return nil
}
func (v SendBulkMessagesRequestI) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *SendBulkMessagesRequestE) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSendBulkMessagesRequestE); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.MsgArr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MsgArr: %v", err)
	}
	return nil
}
func (v SendBulkMessagesRequestE) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSendBulkMessagesRequestE, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.MsgArr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MsgArr: %v", err)
	}
	return nil
}
func (v SendBulkMessagesRequestE) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangePublicKeyRequestI) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangePublicKeyRequestI); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.NewPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewPublicKey: %v", err)
	}
	if err = v.RotationSignature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RotationSignature: %v", err)
	}
	return nil
}
func (v ChangePublicKeyRequestI) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangePublicKeyRequestI, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.NewPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewPublicKey: %v", err)
	}
	if err = v.RotationSignature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RotationSignature: %v", err)
	}
	return nil
}
func (v ChangePublicKeyRequestI) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ChangePublicKeyRequestE) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangePublicKeyRequestE); err != nil {
		return err
	}
	if err = v.Header.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Header: %v", err)
	}
	if err = v.NewPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewPublicKey: %v", err)
	}
	if err = v.RotationSignature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RotationSignature: %v", err)
	}
	return nil
}
func (v ChangePublicKeyRequestE) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangePublicKeyRequestE, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Header.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Header: %v", err)
	}
	if err = v.NewPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewPublicKey: %v", err)
	}
	if err = v.RotationSignature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RotationSignature: %v", err)
	}
	return nil
}
func (v ChangePublicKeyRequestE) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg SignedExternalRequest) ToExternal(address ton.AccountID, init *tlb.StateInitT[*Storage]) (tlb.Message, error) {
	return ton.CreateExternalMessageTWithState(address, msg, init, tlb.VarUInteger16{})
}

func (msg SignedInternalRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*Storage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
