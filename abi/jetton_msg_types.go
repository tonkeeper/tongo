package abi

// Code autogenerated. DO NOT EDIT.

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func decodeTextCommentJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res TextCommentJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = TextCommentJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeTegroJettonSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res TegroJettonSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = TegroJettonSwapJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeEncryptedTextCommentJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res EncryptedTextCommentJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = EncryptedTextCommentJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = StonfiSwapJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeTegroAddLiquidityJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res TegroAddLiquidityJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = TegroAddLiquidityJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiProvideLpV2JettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiProvideLpV2JettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = StonfiProvideLpV2JettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeDedustDepositLiquidityJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res DedustDepositLiquidityJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = DedustDepositLiquidityJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiSwapOkRefJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiSwapOkRefJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil && completedRead(c) {
		j.SumType = StonfiSwapOkRefJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiSwapV2JettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiSwapV2JettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = StonfiSwapV2JettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeInvoicePayloadJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res InvoicePayloadJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = InvoicePayloadJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeTonkeeperRelayerFeeJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res TonkeeperRelayerFeeJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = TonkeeperRelayerFeeJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiSwapOkJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiSwapOkJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil && completedRead(c) {
		j.SumType = StonfiSwapOkJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeDedustSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res DedustSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = DedustSwapJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeStonfiProvideLiquidityJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res StonfiProvideLiquidityJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = StonfiProvideLiquidityJettonOp
		j.Value = res
		return nil
	}
	return err
}

const (
	TextCommentJettonOp            JettonOpName = "TextComment"
	TegroJettonSwapJettonOp        JettonOpName = "TegroJettonSwap"
	EncryptedTextCommentJettonOp   JettonOpName = "EncryptedTextComment"
	StonfiSwapJettonOp             JettonOpName = "StonfiSwap"
	TegroAddLiquidityJettonOp      JettonOpName = "TegroAddLiquidity"
	StonfiProvideLpV2JettonOp      JettonOpName = "StonfiProvideLpV2"
	DedustDepositLiquidityJettonOp JettonOpName = "DedustDepositLiquidity"
	StonfiSwapOkRefJettonOp        JettonOpName = "StonfiSwapOkRef"
	StonfiSwapV2JettonOp           JettonOpName = "StonfiSwapV2"
	InvoicePayloadJettonOp         JettonOpName = "InvoicePayload"
	TonkeeperRelayerFeeJettonOp    JettonOpName = "TonkeeperRelayerFee"
	StonfiSwapOkJettonOp           JettonOpName = "StonfiSwapOk"
	DedustSwapJettonOp             JettonOpName = "DedustSwap"
	StonfiProvideLiquidityJettonOp JettonOpName = "StonfiProvideLiquidity"

	TextCommentJettonOpCode            JettonOpCode = 0x00000000
	TegroJettonSwapJettonOpCode        JettonOpCode = 0x01fb7a25
	EncryptedTextCommentJettonOpCode   JettonOpCode = 0x2167da4b
	StonfiSwapJettonOpCode             JettonOpCode = 0x25938561
	TegroAddLiquidityJettonOpCode      JettonOpCode = 0x287e167a
	StonfiProvideLpV2JettonOpCode      JettonOpCode = 0x37c096df
	DedustDepositLiquidityJettonOpCode JettonOpCode = 0x40e108d6
	StonfiSwapOkRefJettonOpCode        JettonOpCode = 0x45078540
	StonfiSwapV2JettonOpCode           JettonOpCode = 0x6664de2a
	InvoicePayloadJettonOpCode         JettonOpCode = 0x7aa23eb5
	TonkeeperRelayerFeeJettonOpCode    JettonOpCode = 0x878da6e3
	StonfiSwapOkJettonOpCode           JettonOpCode = 0xc64370e5
	DedustSwapJettonOpCode             JettonOpCode = 0xe3a0d482
	StonfiProvideLiquidityJettonOpCode JettonOpCode = 0xfcf9e58f
)

var KnownJettonTypes = map[string]any{
	TextCommentJettonOp:            TextCommentJettonPayload{},
	TegroJettonSwapJettonOp:        TegroJettonSwapJettonPayload{},
	EncryptedTextCommentJettonOp:   EncryptedTextCommentJettonPayload{},
	StonfiSwapJettonOp:             StonfiSwapJettonPayload{},
	TegroAddLiquidityJettonOp:      TegroAddLiquidityJettonPayload{},
	StonfiProvideLpV2JettonOp:      StonfiProvideLpV2JettonPayload{},
	DedustDepositLiquidityJettonOp: DedustDepositLiquidityJettonPayload{},
	StonfiSwapOkRefJettonOp:        StonfiSwapOkRefJettonPayload{},
	StonfiSwapV2JettonOp:           StonfiSwapV2JettonPayload{},
	InvoicePayloadJettonOp:         InvoicePayloadJettonPayload{},
	TonkeeperRelayerFeeJettonOp:    TonkeeperRelayerFeeJettonPayload{},
	StonfiSwapOkJettonOp:           StonfiSwapOkJettonPayload{},
	DedustSwapJettonOp:             DedustSwapJettonPayload{},
	StonfiProvideLiquidityJettonOp: StonfiProvideLiquidityJettonPayload{},
}
var JettonOpCodes = map[JettonOpName]JettonOpCode{
	TextCommentJettonOp:            TextCommentJettonOpCode,
	TegroJettonSwapJettonOp:        TegroJettonSwapJettonOpCode,
	EncryptedTextCommentJettonOp:   EncryptedTextCommentJettonOpCode,
	StonfiSwapJettonOp:             StonfiSwapJettonOpCode,
	TegroAddLiquidityJettonOp:      TegroAddLiquidityJettonOpCode,
	StonfiProvideLpV2JettonOp:      StonfiProvideLpV2JettonOpCode,
	DedustDepositLiquidityJettonOp: DedustDepositLiquidityJettonOpCode,
	StonfiSwapOkRefJettonOp:        StonfiSwapOkRefJettonOpCode,
	StonfiSwapV2JettonOp:           StonfiSwapV2JettonOpCode,
	InvoicePayloadJettonOp:         InvoicePayloadJettonOpCode,
	TonkeeperRelayerFeeJettonOp:    TonkeeperRelayerFeeJettonOpCode,
	StonfiSwapOkJettonOp:           StonfiSwapOkJettonOpCode,
	DedustSwapJettonOp:             DedustSwapJettonOpCode,
	StonfiProvideLiquidityJettonOp: StonfiProvideLiquidityJettonOpCode,
}

var funcJettonDecodersMapping = map[JettonOpCode]func(*JettonPayload, *boc.Cell) error{
	TextCommentJettonOpCode:            decodeTextCommentJettonOpJetton,
	TegroJettonSwapJettonOpCode:        decodeTegroJettonSwapJettonOpJetton,
	EncryptedTextCommentJettonOpCode:   decodeEncryptedTextCommentJettonOpJetton,
	StonfiSwapJettonOpCode:             decodeStonfiSwapJettonOpJetton,
	TegroAddLiquidityJettonOpCode:      decodeTegroAddLiquidityJettonOpJetton,
	StonfiProvideLpV2JettonOpCode:      decodeStonfiProvideLpV2JettonOpJetton,
	DedustDepositLiquidityJettonOpCode: decodeDedustDepositLiquidityJettonOpJetton,
	StonfiSwapOkRefJettonOpCode:        decodeStonfiSwapOkRefJettonOpJetton,
	StonfiSwapV2JettonOpCode:           decodeStonfiSwapV2JettonOpJetton,
	InvoicePayloadJettonOpCode:         decodeInvoicePayloadJettonOpJetton,
	TonkeeperRelayerFeeJettonOpCode:    decodeTonkeeperRelayerFeeJettonOpJetton,
	StonfiSwapOkJettonOpCode:           decodeStonfiSwapOkJettonOpJetton,
	DedustSwapJettonOpCode:             decodeDedustSwapJettonOpJetton,
	StonfiProvideLiquidityJettonOpCode: decodeStonfiProvideLiquidityJettonOpJetton,
}

type TextCommentJettonPayload struct {
	Text tlb.Text
}

type TegroJettonSwapJettonPayload struct {
	Extract          bool
	MaxIn            tlb.VarUInteger16
	MinOut           tlb.VarUInteger16
	Destination      tlb.MsgAddress
	ErrorDestination tlb.MsgAddress
	Payload          *tlb.Any `tlb:"maybe^"`
}

type EncryptedTextCommentJettonPayload struct {
	CipherText tlb.Bytes
}

type StonfiSwapJettonPayload struct {
	TokenWallet     tlb.MsgAddress
	MinOut          tlb.VarUInteger16
	ToAddress       tlb.MsgAddress
	ReferralAddress *tlb.MsgAddress `tlb:"maybe"`
}

type TegroAddLiquidityJettonPayload struct {
	AmountA tlb.VarUInteger16
	AbountB tlb.VarUInteger16
}

type StonfiProvideLpV2JettonPayload struct {
	TokenWallet1       tlb.MsgAddress
	RefundAddress      tlb.MsgAddress
	ExcessesAddress    tlb.MsgAddress
	TxDeadline         uint64
	CrossProvideLpBody struct {
		MinLpOut      tlb.VarUInteger16
		ToAddress     tlb.MsgAddress
		BothPositive  tlb.Uint1
		FwdAmount     tlb.Grams
		CustomPayload *tlb.Any `tlb:"maybe^"`
	} `tlb:"^"`
}

type DedustDepositLiquidityJettonPayload struct {
	PoolParams          DedustPoolParams
	MinLpAmount         tlb.Grams
	Asset0TargetBalance tlb.Grams
	Asset1TargetBalance tlb.Grams
	FulfillPayload      *tlb.Any `tlb:"maybe^"`
	RejectPayload       *tlb.Any `tlb:"maybe^"`
}

type StonfiSwapOkRefJettonPayload struct{}

type StonfiSwapV2JettonPayload struct {
	TokenWallet1    tlb.MsgAddress
	RefundAddress   tlb.MsgAddress
	ExcessesAddress tlb.MsgAddress
	TxDeadline      uint64
	CrossSwapBody   struct {
		MinOut        tlb.VarUInteger16
		Receiver      tlb.MsgAddress
		FwdGas        tlb.Grams
		CustomPayload *tlb.Any `tlb:"maybe^"`
		RefundFwdGas  tlb.Grams
		RefundPayload *tlb.Any `tlb:"maybe^"`
		RefFee        uint16
		RefAddress    tlb.MsgAddress
	} `tlb:"^"`
}

type InvoicePayloadJettonPayload struct {
	Id  tlb.Bits128
	Url PaymentProviderUrl
}

type TonkeeperRelayerFeeJettonPayload struct{}

type StonfiSwapOkJettonPayload struct{}

type DedustSwapJettonPayload struct {
	Step       DedustSwapStep
	SwapParams DedustSwapParams `tlb:"^"`
}

type StonfiProvideLiquidityJettonPayload struct {
	TokenWallet tlb.MsgAddress
	MinLpOut    tlb.VarUInteger16
}
