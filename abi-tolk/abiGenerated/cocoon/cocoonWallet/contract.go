// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64
}
type QueryHeader struct {
	Op      tlb.Uint32
	QueryId tlb.Uint64
}
type SignedMessage struct {
	Op             tlb.Uint32
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
	Signature      tlb.Bits512
	SignedDataCell boc.Cell
}
type ExternalSignedMessage struct {
	SubwalletId tlb.Uint32
	ValidUntil  tlb.Uint32
	MsgSeqno    tlb.Uint32
	Rest        tlb.Any
}

const PrefixOwnerWalletSendMessage uint64 = 0x9c69f376

type OwnerWalletSendMessage struct {
	QueryId tlb.Uint64
	Mode    tlb.Uint8
	Body    boc.Cell
}

const PrefixTextCommand uint64 = 0x00000000

type TextCommand struct {
	Action tlb.Uint8
}
type AllowedInternalMessageKind uint

const (
	AllowedInternalMessageKind_OwnerWalletSendMessage AllowedInternalMessageKind = 0x9c69f376
	AllowedInternalMessageKind_TextCommand            AllowedInternalMessageKind = 0x00000000
)

type AllowedInternalMessage struct { // tagged union
	SumType                AllowedInternalMessageKind
	OwnerWalletSendMessage *OwnerWalletSendMessage
	TextCommand            *TextCommand
}

func (v *AllowedInternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = AllowedInternalMessageKind(prefix)
	switch v.SumType {
	case AllowedInternalMessageKind_OwnerWalletSendMessage:
		v.OwnerWalletSendMessage = new(OwnerWalletSendMessage)
		return v.OwnerWalletSendMessage.UnmarshalTLB(c, decoder)
	case AllowedInternalMessageKind_TextCommand:
		v.TextCommand = new(TextCommand)
		return v.TextCommand.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v *ReturnExcessesBack) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixReturnExcessesBack {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *Payout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixPayout {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *QueryHeader) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *SignedMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Op.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.SignedDataCell, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *ExternalSignedMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.SubwalletId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ValidUntil.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.MsgSeqno.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.Rest, err = (func() (tlb.Any, error) {
		cc := c.CopyRemaining()
		return tlb.Any(*cc), nil
	})(); err != nil {
		return err
	}
	return nil
}
func (v *OwnerWalletSendMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerWalletSendMessage {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Mode.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if v.Body, err = c.NextRefV(); err != nil {
		return err
	}
	return nil
}
func (v *TextCommand) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixTextCommand {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Action.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
