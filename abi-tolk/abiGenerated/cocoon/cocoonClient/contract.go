// Code generated - DO NOT EDIT.

package abiGenerated

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type ClientStateData struct {
	State      tlb.Uint2
	Balance    tlb.Coins
	Stake      tlb.Coins
	TokensUsed tlb.Uint64
	SecretHash tlb.Uint256
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

const PrefixExtClientChargeSigned uint64 = 0xbb63ff93

type ExtClientChargeSigned struct {
	Rest tlb.Any
}

const PrefixExtClientGrantRefundSigned uint64 = 0xefd711e1

type ExtClientGrantRefundSigned struct {
	Rest tlb.Any
}

const PrefixExtClientTopUp uint64 = 0xf172e6c2

type ExtClientTopUp struct {
	QueryId        tlb.Uint64
	TopUpAmount    tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientChangeSecretHashAndTopUp uint64 = 0x8473b408

type OwnerClientChangeSecretHashAndTopUp struct {
	QueryId        tlb.Uint64
	TopUpAmount    tlb.Coins
	NewSecretHash  tlb.Uint256
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientRegister uint64 = 0xc45f9f3b

type OwnerClientRegister struct {
	QueryId        tlb.Uint64
	Nonce          tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientChangeSecretHash uint64 = 0xa9357034

type OwnerClientChangeSecretHash struct {
	QueryId        tlb.Uint64
	NewSecretHash  tlb.Uint256
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientIncreaseStake uint64 = 0x6a1f6a60

type OwnerClientIncreaseStake struct {
	QueryId        tlb.Uint64
	NewStake       tlb.Coins
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientWithdraw uint64 = 0xda068e78

type OwnerClientWithdraw struct {
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixOwnerClientRequestRefund uint64 = 0xfafa6cc1

type OwnerClientRequestRefund struct {
	QueryId        tlb.Uint64
	SendExcessesTo tlb.InternalAddress
}

const PrefixChargePayload uint64 = 0xbb63ff93

type ChargePayload struct {
	QueryId           tlb.Uint64
	NewTokensUsed     tlb.Uint64
	ExpectedMyAddress tlb.InternalAddress
}

const PrefixGrantRefundPayload uint64 = 0xefd711e1

type GrantRefundPayload struct {
	QueryId           tlb.Uint64
	NewTokensUsed     tlb.Uint64
	ExpectedMyAddress tlb.InternalAddress
}
type SignedClientMessageKind uint

const (
	SignedClientMessageKind_ExtClientChargeSigned      SignedClientMessageKind = 0xbb63ff93
	SignedClientMessageKind_ExtClientGrantRefundSigned SignedClientMessageKind = 0xefd711e1
)

type SignedClientMessage struct { // tagged union
	SumType                    SignedClientMessageKind
	ExtClientChargeSigned      *ExtClientChargeSigned
	ExtClientGrantRefundSigned *ExtClientGrantRefundSigned
}

func (v *SignedClientMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SignedClientMessageKind(prefix)
	switch v.SumType {
	case SignedClientMessageKind_ExtClientChargeSigned:
		v.ExtClientChargeSigned = new(ExtClientChargeSigned)
		return v.ExtClientChargeSigned.UnmarshalTLB(c, decoder)
	case SignedClientMessageKind_ExtClientGrantRefundSigned:
		v.ExtClientGrantRefundSigned = new(ExtClientGrantRefundSigned)
		return v.ExtClientGrantRefundSigned.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

type ClientSignedPayloadKind uint

const (
	ClientSignedPayloadKind_ChargePayload      ClientSignedPayloadKind = 0xbb63ff93
	ClientSignedPayloadKind_GrantRefundPayload ClientSignedPayloadKind = 0xefd711e1
)

type ClientSignedPayload struct { // tagged union
	SumType            ClientSignedPayloadKind
	ChargePayload      *ChargePayload
	GrantRefundPayload *GrantRefundPayload
}

func (v *ClientSignedPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = ClientSignedPayloadKind(prefix)
	switch v.SumType {
	case ClientSignedPayloadKind_ChargePayload:
		v.ChargePayload = new(ChargePayload)
		return v.ChargePayload.UnmarshalTLB(c, decoder)
	case ClientSignedPayloadKind_GrantRefundPayload:
		v.GrantRefundPayload = new(GrantRefundPayload)
		return v.GrantRefundPayload.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64
}
type ClientMessageKind uint

const (
	ClientMessageKind_ExtClientChargeSigned               ClientMessageKind = 0xbb63ff93
	ClientMessageKind_ExtClientGrantRefundSigned          ClientMessageKind = 0xefd711e1
	ClientMessageKind_ExtClientTopUp                      ClientMessageKind = 0xf172e6c2
	ClientMessageKind_OwnerClientChangeSecretHashAndTopUp ClientMessageKind = 0x8473b408
	ClientMessageKind_OwnerClientRegister                 ClientMessageKind = 0xc45f9f3b
	ClientMessageKind_OwnerClientChangeSecretHash         ClientMessageKind = 0xa9357034
	ClientMessageKind_OwnerClientIncreaseStake            ClientMessageKind = 0x6a1f6a60
	ClientMessageKind_OwnerClientWithdraw                 ClientMessageKind = 0xda068e78
	ClientMessageKind_OwnerClientRequestRefund            ClientMessageKind = 0xfafa6cc1
)

type ClientMessage struct { // tagged union
	SumType                             ClientMessageKind
	ExtClientChargeSigned               *ExtClientChargeSigned
	ExtClientGrantRefundSigned          *ExtClientGrantRefundSigned
	ExtClientTopUp                      *ExtClientTopUp
	OwnerClientChangeSecretHashAndTopUp *OwnerClientChangeSecretHashAndTopUp
	OwnerClientRegister                 *OwnerClientRegister
	OwnerClientChangeSecretHash         *OwnerClientChangeSecretHash
	OwnerClientIncreaseStake            *OwnerClientIncreaseStake
	OwnerClientWithdraw                 *OwnerClientWithdraw
	OwnerClientRequestRefund            *OwnerClientRequestRefund
}

func (v *ClientMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = ClientMessageKind(prefix)
	switch v.SumType {
	case ClientMessageKind_ExtClientChargeSigned:
		v.ExtClientChargeSigned = new(ExtClientChargeSigned)
		return v.ExtClientChargeSigned.UnmarshalTLB(c, decoder)
	case ClientMessageKind_ExtClientGrantRefundSigned:
		v.ExtClientGrantRefundSigned = new(ExtClientGrantRefundSigned)
		return v.ExtClientGrantRefundSigned.UnmarshalTLB(c, decoder)
	case ClientMessageKind_ExtClientTopUp:
		v.ExtClientTopUp = new(ExtClientTopUp)
		return v.ExtClientTopUp.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientChangeSecretHashAndTopUp:
		v.OwnerClientChangeSecretHashAndTopUp = new(OwnerClientChangeSecretHashAndTopUp)
		return v.OwnerClientChangeSecretHashAndTopUp.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientRegister:
		v.OwnerClientRegister = new(OwnerClientRegister)
		return v.OwnerClientRegister.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientChangeSecretHash:
		v.OwnerClientChangeSecretHash = new(OwnerClientChangeSecretHash)
		return v.OwnerClientChangeSecretHash.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientIncreaseStake:
		v.OwnerClientIncreaseStake = new(OwnerClientIncreaseStake)
		return v.OwnerClientIncreaseStake.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientWithdraw:
		v.OwnerClientWithdraw = new(OwnerClientWithdraw)
		return v.OwnerClientWithdraw.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientRequestRefund:
		v.OwnerClientRequestRefund = new(OwnerClientRequestRefund)
		return v.OwnerClientRequestRefund.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
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
func (v *ExtClientChargeSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtClientChargeSigned {
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
func (v *ExtClientGrantRefundSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtClientGrantRefundSigned {
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
func (v *ExtClientTopUp) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixExtClientTopUp {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.TopUpAmount.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerClientChangeSecretHashAndTopUp) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientChangeSecretHashAndTopUp {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.TopUpAmount.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewSecretHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerClientRegister) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientRegister {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.Nonce.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerClientChangeSecretHash) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientChangeSecretHash {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewSecretHash.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerClientIncreaseStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientIncreaseStake {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewStake.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *OwnerClientWithdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientWithdraw {
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
func (v *OwnerClientRequestRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixOwnerClientRequestRefund {
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
func (v *ChargePayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixChargePayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewTokensUsed.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ExpectedMyAddress.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	return nil
}
func (v *GrantRefundPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if prefix, err := c.ReadUint(32); err != nil {
		return err
	} else if prefix != PrefixGrantRefundPayload {
		return fmt.Errorf("unexpected prefix: %x", prefix)
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.NewTokensUsed.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	if err = v.ExpectedMyAddress.UnmarshalTLB(c, decoder); err != nil {
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
