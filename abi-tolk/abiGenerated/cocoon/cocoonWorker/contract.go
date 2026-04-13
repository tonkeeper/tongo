// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixWorkerProxyPayoutRequest uint64 = 0x08e7d036

type WorkerProxyPayoutRequest struct {
	WorkerPart     tlb.Coins
	ProxyPart      tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerWorkerRegister uint64 = 0x26ed7f65

type OwnerWorkerRegister struct {
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixExtWorkerPayoutRequestSigned uint64 = 0xa040ad28

type ExtWorkerPayoutRequestSigned struct {
	Rest tlb.Any
}

const PrefixExtWorkerLastPayoutRequestSigned uint64 = 0xf5f26a36

type ExtWorkerLastPayoutRequestSigned struct {
	Rest tlb.Any
}

const PrefixWorkerProxyRequest uint64 = 0x4d725d2c

type WorkerProxyRequest struct {
	QueryId      tlb.Uint64
	OwnerAddress tlb.InternalAddress
	State        tlb.Uint2
	Tokens       tlb.Uint64
	Payload      tlb.Maybe[boc.Cell]
}
type PayoutPayloadData struct {
	QueryId           tlb.Uint64
	NewTokens         tlb.Uint64
	ExpectedMyAddress tlb.InternalAddress
}

const PrefixPayoutPayload uint64 = 0xa040ad28

type PayoutPayload struct {
	Data PayoutPayloadData
}

const PrefixLastPayoutPayload uint64 = 0xf5f26a36

type LastPayoutPayload struct {
	Data PayoutPayloadData
}

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64
}
type SignedWorkerMessageKind uint

const (
	SignedWorkerMessageKind_ExtWorkerPayoutRequestSigned     SignedWorkerMessageKind = 0xa040ad28
	SignedWorkerMessageKind_ExtWorkerLastPayoutRequestSigned SignedWorkerMessageKind = 0xf5f26a36
)

type SignedWorkerMessage struct { // tagged union
	SumType                          SignedWorkerMessageKind
	ExtWorkerPayoutRequestSigned     *ExtWorkerPayoutRequestSigned
	ExtWorkerLastPayoutRequestSigned *ExtWorkerLastPayoutRequestSigned
}

func (v *SignedWorkerMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SignedWorkerMessageKind(prefix)
	switch v.SumType {
	case SignedWorkerMessageKind_ExtWorkerPayoutRequestSigned:
		v.ExtWorkerPayoutRequestSigned = new(ExtWorkerPayoutRequestSigned)
		return v.ExtWorkerPayoutRequestSigned.UnmarshalTLB(c, decoder)
	case SignedWorkerMessageKind_ExtWorkerLastPayoutRequestSigned:
		v.ExtWorkerLastPayoutRequestSigned = new(ExtWorkerLastPayoutRequestSigned)
		return v.ExtWorkerLastPayoutRequestSigned.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

type SignedPayloadKind uint

const (
	SignedPayloadKind_PayoutPayload     SignedPayloadKind = 0xa040ad28
	SignedPayloadKind_LastPayoutPayload SignedPayloadKind = 0xf5f26a36
)

type SignedPayload struct { // tagged union
	SumType           SignedPayloadKind
	PayoutPayload     *PayoutPayload
	LastPayoutPayload *LastPayoutPayload
}

func (v *SignedPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SignedPayloadKind(prefix)
	switch v.SumType {
	case SignedPayloadKind_PayoutPayload:
		v.PayoutPayload = new(PayoutPayload)
		return v.PayoutPayload.UnmarshalTLB(c, decoder)
	case SignedPayloadKind_LastPayoutPayload:
		v.LastPayoutPayload = new(LastPayoutPayload)
		return v.LastPayoutPayload.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

type WorkerMessageKind uint

const (
	WorkerMessageKind_ExtWorkerPayoutRequestSigned     WorkerMessageKind = 0xa040ad28
	WorkerMessageKind_ExtWorkerLastPayoutRequestSigned WorkerMessageKind = 0xf5f26a36
	WorkerMessageKind_OwnerWorkerRegister              WorkerMessageKind = 0x26ed7f65
	WorkerMessageKind_WorkerProxyRequest               WorkerMessageKind = 0x4d725d2c
)

type WorkerMessage struct { // tagged union
	SumType                          WorkerMessageKind
	ExtWorkerPayoutRequestSigned     *ExtWorkerPayoutRequestSigned
	ExtWorkerLastPayoutRequestSigned *ExtWorkerLastPayoutRequestSigned
	OwnerWorkerRegister              *OwnerWorkerRegister
	WorkerProxyRequest               *WorkerProxyRequest
}

func (v *WorkerMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = WorkerMessageKind(prefix)
	switch v.SumType {
	case WorkerMessageKind_ExtWorkerPayoutRequestSigned:
		v.ExtWorkerPayoutRequestSigned = new(ExtWorkerPayoutRequestSigned)
		return v.ExtWorkerPayoutRequestSigned.UnmarshalTLB(c, decoder)
	case WorkerMessageKind_ExtWorkerLastPayoutRequestSigned:
		v.ExtWorkerLastPayoutRequestSigned = new(ExtWorkerLastPayoutRequestSigned)
		return v.ExtWorkerLastPayoutRequestSigned.UnmarshalTLB(c, decoder)
	case WorkerMessageKind_OwnerWorkerRegister:
		v.OwnerWorkerRegister = new(OwnerWorkerRegister)
		return v.OwnerWorkerRegister.UnmarshalTLB(c, decoder)
	case WorkerMessageKind_WorkerProxyRequest:
		v.WorkerProxyRequest = new(WorkerProxyRequest)
		return v.WorkerProxyRequest.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v *WorkerProxyPayoutRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixWorkerProxyPayoutRequest {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.WorkerPart.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ProxyPart.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerWorkerRegister) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerWorkerRegister {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ExtWorkerPayoutRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtWorkerPayoutRequestSigned {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.Rest, err = (func() (tlb.Any, error) {
		cc := c.CopyRemaining()
		return tlb.Any(*cc), nil
	})(); err != nil {
		return err
	}
	return nil
}
func (v *ExtWorkerLastPayoutRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtWorkerLastPayoutRequestSigned {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if v.Rest, err = (func() (tlb.Any, error) {
		cc := c.CopyRemaining()
		return tlb.Any(*cc), nil
	})(); err != nil {
		return err
	}
	return nil
}
func (v *WorkerProxyRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixWorkerProxyRequest {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Tokens.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Payload.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *PayoutPayloadData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewTokens.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ExpectedMyAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *PayoutPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixPayoutPayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Data.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *LastPayoutPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixLastPayoutPayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Data.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
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
