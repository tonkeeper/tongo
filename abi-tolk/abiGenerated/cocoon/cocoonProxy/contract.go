// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixTextCmd uint64 = 0x00000000

type TextCmd struct {
	Action tlb.Uint8
}

const PrefixExtProxyPayoutRequest uint64 = 0x7610e6eb

type ExtProxyPayoutRequest struct {
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixExtProxyIncreaseStake uint64 = 0x9713f187

type ExtProxyIncreaseStake struct {
	QueryId        tlb.Uint64
	Grams          tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerProxyClose uint64 = 0xb51d5a01

type OwnerProxyClose struct {
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixCloseRequestPayload uint64 = 0x636a4391

type CloseRequestPayload struct {
	QueryId           tlb.Uint64
	ExpectedMyAddress tlb.InternalAddress
}

const PrefixCloseCompleteRequestPayload uint64 = 0xe511abc7

type CloseCompleteRequestPayload struct {
	QueryId           tlb.Uint64
	ExpectedMyAddress tlb.InternalAddress
}
type SignedProxyPayloadKind uint

const (
	SignedProxyPayloadKind_CloseRequestPayload         SignedProxyPayloadKind = 0x636a4391
	SignedProxyPayloadKind_CloseCompleteRequestPayload SignedProxyPayloadKind = 0xe511abc7
)

type SignedProxyPayload struct { // tagged union
	SumType                     SignedProxyPayloadKind
	CloseRequestPayload         *CloseRequestPayload
	CloseCompleteRequestPayload *CloseCompleteRequestPayload
}

func (v *SignedProxyPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SignedProxyPayloadKind(prefix)
	switch v.SumType {
	case SignedProxyPayloadKind_CloseRequestPayload:
		v.CloseRequestPayload = new(CloseRequestPayload)
		return v.CloseRequestPayload.UnmarshalTLB(c, decoder)
	case SignedProxyPayloadKind_CloseCompleteRequestPayload:
		v.CloseCompleteRequestPayload = new(CloseCompleteRequestPayload)
		return v.CloseCompleteRequestPayload.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

const PrefixExtProxyCloseRequestSigned uint64 = 0x636a4391

type ExtProxyCloseRequestSigned struct {
	Rest tlb.Any
}

const PrefixExtProxyCloseCompleteRequestSigned uint64 = 0xe511abc7

type ExtProxyCloseCompleteRequestSigned struct {
	Rest tlb.Any
}
type ClientStateData struct {
	State      tlb.Uint2
	Balance    tlb.Coins
	Stake      tlb.Coins
	TokensUsed tlb.Uint64
	SecretHash tlb.Uint256
}

const PrefixWorkerProxyRequest uint64 = 0x4d725d2c

type WorkerProxyRequest struct {
	QueryId      tlb.Uint64
	OwnerAddress tlb.InternalAddress
	State        tlb.Uint2
	Tokens       tlb.Uint64
	Payload      tlb.Maybe[boc.Cell]
}

const PrefixWorkerProxyPayoutRequest uint64 = 0x08e7d036

type WorkerProxyPayoutRequest struct {
	WorkerPart     tlb.Coins
	ProxyPart      tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixClientProxyRequest uint64 = 0x65448ff4

type ClientProxyRequest struct {
	QueryId      tlb.Uint64
	OwnerAddress tlb.InternalAddress
	StateData    tlb.RefT[*ClientStateData]
	Payload      tlb.Maybe[boc.Cell]
}

const PrefixClientProxyTopUp uint64 = 0x5cfc6b87

type ClientProxyTopUp struct {
	TopUpCoins     tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixClientProxyRegister uint64 = 0xa35cb580

type ClientProxyRegister struct {
}

const PrefixClientProxyRefundGranted uint64 = 0xc68ebc7b

type ClientProxyRefundGranted struct {
	Coins          tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixClientProxyRefundForce uint64 = 0xf4c354c9

type ClientProxyRefundForce struct {
	Coins          tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64
}
type SignedProxyMessageKind uint

const (
	SignedProxyMessageKind_ExtProxyCloseRequestSigned         SignedProxyMessageKind = 0x636a4391
	SignedProxyMessageKind_ExtProxyCloseCompleteRequestSigned SignedProxyMessageKind = 0xe511abc7
)

type SignedProxyMessage struct { // tagged union
	SumType                            SignedProxyMessageKind
	ExtProxyCloseRequestSigned         *ExtProxyCloseRequestSigned
	ExtProxyCloseCompleteRequestSigned *ExtProxyCloseCompleteRequestSigned
}

func (v *SignedProxyMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SignedProxyMessageKind(prefix)
	switch v.SumType {
	case SignedProxyMessageKind_ExtProxyCloseRequestSigned:
		v.ExtProxyCloseRequestSigned = new(ExtProxyCloseRequestSigned)
		return v.ExtProxyCloseRequestSigned.UnmarshalTLB(c, decoder)
	case SignedProxyMessageKind_ExtProxyCloseCompleteRequestSigned:
		v.ExtProxyCloseCompleteRequestSigned = new(ExtProxyCloseCompleteRequestSigned)
		return v.ExtProxyCloseCompleteRequestSigned.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

type ClientProxyPayloadKind uint

const (
	ClientProxyPayloadKind_ClientProxyTopUp         ClientProxyPayloadKind = 0x5cfc6b87
	ClientProxyPayloadKind_ClientProxyRegister      ClientProxyPayloadKind = 0xa35cb580
	ClientProxyPayloadKind_ClientProxyRefundGranted ClientProxyPayloadKind = 0xc68ebc7b
	ClientProxyPayloadKind_ClientProxyRefundForce   ClientProxyPayloadKind = 0xf4c354c9
)

type ClientProxyPayload struct { // tagged union
	SumType                  ClientProxyPayloadKind
	ClientProxyTopUp         *ClientProxyTopUp
	ClientProxyRegister      *ClientProxyRegister
	ClientProxyRefundGranted *ClientProxyRefundGranted
	ClientProxyRefundForce   *ClientProxyRefundForce
}

func (v *ClientProxyPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = ClientProxyPayloadKind(prefix)
	switch v.SumType {
	case ClientProxyPayloadKind_ClientProxyTopUp:
		v.ClientProxyTopUp = new(ClientProxyTopUp)
		return v.ClientProxyTopUp.UnmarshalTLB(c, decoder)
	case ClientProxyPayloadKind_ClientProxyRegister:
		v.ClientProxyRegister = new(ClientProxyRegister)
		return v.ClientProxyRegister.UnmarshalTLB(c, decoder)
	case ClientProxyPayloadKind_ClientProxyRefundGranted:
		v.ClientProxyRefundGranted = new(ClientProxyRefundGranted)
		return v.ClientProxyRefundGranted.UnmarshalTLB(c, decoder)
	case ClientProxyPayloadKind_ClientProxyRefundForce:
		v.ClientProxyRefundForce = new(ClientProxyRefundForce)
		return v.ClientProxyRefundForce.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

type WorkerProxyPayload WorkerProxyPayoutRequest
type ProxyMessageKind uint

const (
	ProxyMessageKind_TextCmd                            ProxyMessageKind = 0x00000000
	ProxyMessageKind_ExtProxyCloseRequestSigned         ProxyMessageKind = 0x636a4391
	ProxyMessageKind_ExtProxyCloseCompleteRequestSigned ProxyMessageKind = 0xe511abc7
	ProxyMessageKind_ExtProxyPayoutRequest              ProxyMessageKind = 0x7610e6eb
	ProxyMessageKind_ExtProxyIncreaseStake              ProxyMessageKind = 0x9713f187
	ProxyMessageKind_OwnerProxyClose                    ProxyMessageKind = 0xb51d5a01
	ProxyMessageKind_WorkerProxyRequest                 ProxyMessageKind = 0x4d725d2c
	ProxyMessageKind_ClientProxyRequest                 ProxyMessageKind = 0x65448ff4
)

type ProxyMessage struct { // tagged union
	SumType                            ProxyMessageKind
	TextCmd                            *TextCmd
	ExtProxyCloseRequestSigned         *ExtProxyCloseRequestSigned
	ExtProxyCloseCompleteRequestSigned *ExtProxyCloseCompleteRequestSigned
	ExtProxyPayoutRequest              *ExtProxyPayoutRequest
	ExtProxyIncreaseStake              *ExtProxyIncreaseStake
	OwnerProxyClose                    *OwnerProxyClose
	WorkerProxyRequest                 *WorkerProxyRequest
	ClientProxyRequest                 *ClientProxyRequest
}

func (v *ProxyMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = ProxyMessageKind(prefix)
	switch v.SumType {
	case ProxyMessageKind_TextCmd:
		v.TextCmd = new(TextCmd)
		return v.TextCmd.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_ExtProxyCloseRequestSigned:
		v.ExtProxyCloseRequestSigned = new(ExtProxyCloseRequestSigned)
		return v.ExtProxyCloseRequestSigned.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_ExtProxyCloseCompleteRequestSigned:
		v.ExtProxyCloseCompleteRequestSigned = new(ExtProxyCloseCompleteRequestSigned)
		return v.ExtProxyCloseCompleteRequestSigned.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_ExtProxyPayoutRequest:
		v.ExtProxyPayoutRequest = new(ExtProxyPayoutRequest)
		return v.ExtProxyPayoutRequest.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_ExtProxyIncreaseStake:
		v.ExtProxyIncreaseStake = new(ExtProxyIncreaseStake)
		return v.ExtProxyIncreaseStake.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_OwnerProxyClose:
		v.OwnerProxyClose = new(OwnerProxyClose)
		return v.OwnerProxyClose.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_WorkerProxyRequest:
		v.WorkerProxyRequest = new(WorkerProxyRequest)
		return v.WorkerProxyRequest.UnmarshalTLB(c, decoder)
	case ProxyMessageKind_ClientProxyRequest:
		v.ClientProxyRequest = new(ClientProxyRequest)
		return v.ClientProxyRequest.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v *TextCmd) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixTextCmd {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Action.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ExtProxyPayoutRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtProxyPayoutRequest {
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
func (v *ExtProxyIncreaseStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtProxyIncreaseStake {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Grams.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerProxyClose) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerProxyClose {
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
func (v *CloseRequestPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixCloseRequestPayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ExpectedMyAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *CloseCompleteRequestPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixCloseCompleteRequestPayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ExpectedMyAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ExtProxyCloseRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtProxyCloseRequestSigned {
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
func (v *ExtProxyCloseCompleteRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtProxyCloseCompleteRequestSigned {
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
func (v *ClientStateData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.TokensUsed.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SecretHash.UnmarshalTLB(c, decoder); err != nil {
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
func (v *ClientProxyRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixClientProxyRequest {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.StateData.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Payload.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ClientProxyTopUp) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixClientProxyTopUp {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.TopUpCoins.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ClientProxyRegister) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixClientProxyRegister {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	return nil
}
func (v *ClientProxyRefundGranted) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixClientProxyRefundGranted {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Coins.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *ClientProxyRefundForce) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixClientProxyRefundForce {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.Coins.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
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
func (v *WorkerProxyPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	vx, err := tlb.UnmarshalT[WorkerProxyPayoutRequest](c, decoder)
	if err != nil {
		return err
	}
	*v = WorkerProxyPayload(vx)
	return nil
}
