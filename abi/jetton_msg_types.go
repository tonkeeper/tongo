package abi

// Code autogenerated. DO NOT EDIT.

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func decodeCoffeeStakingLockJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeStakingLockJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeStakingLockJettonOp
		j.Value = res
		return nil
	}
	return err
}

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

func decodeBidaskProvideBothJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res BidaskProvideBothJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = BidaskProvideBothJettonOp
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

func decodeCoffeeCrossDexResendJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeCrossDexResendJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeCrossDexResendJettonOp
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

func decodeBidaskProvideJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res BidaskProvideJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = BidaskProvideJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeCoffeeCrossDexFailureJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeCrossDexFailureJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeCrossDexFailureJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeCoffeeSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeSwapJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeCoffeeCreatePoolJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeCreatePoolJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeCreatePoolJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeCoffeeDepositLiquidityJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeDepositLiquidityJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeDepositLiquidityJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeCoffeeNotificationJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeNotificationJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeNotificationJettonOp
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

func decodeCoffeeMevProtectFailedSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res CoffeeMevProtectFailedSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = CoffeeMevProtectFailedSwapJettonOp
		j.Value = res
		return nil
	}
	return err
}

func decodeBidaskSwapJettonOpJetton(j *JettonPayload, c *boc.Cell) error {
	var res BidaskSwapJettonPayload
	err := tlb.Unmarshal(c, &res)
	if err == nil {
		j.SumType = BidaskSwapJettonOp
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
	CoffeeStakingLockJettonOp          JettonOpName = "CoffeeStakingLock"
	TextCommentJettonOp                JettonOpName = "TextComment"
	TegroJettonSwapJettonOp            JettonOpName = "TegroJettonSwap"
	EncryptedTextCommentJettonOp       JettonOpName = "EncryptedTextComment"
	StonfiSwapJettonOp                 JettonOpName = "StonfiSwap"
	TegroAddLiquidityJettonOp          JettonOpName = "TegroAddLiquidity"
	StonfiProvideLpV2JettonOp          JettonOpName = "StonfiProvideLpV2"
	BidaskProvideBothJettonOp          JettonOpName = "BidaskProvideBoth"
	DedustDepositLiquidityJettonOp     JettonOpName = "DedustDepositLiquidity"
	StonfiSwapOkRefJettonOp            JettonOpName = "StonfiSwapOkRef"
	CoffeeCrossDexResendJettonOp       JettonOpName = "CoffeeCrossDexResend"
	StonfiSwapV2JettonOp               JettonOpName = "StonfiSwapV2"
	InvoicePayloadJettonOp             JettonOpName = "InvoicePayload"
	TonkeeperRelayerFeeJettonOp        JettonOpName = "TonkeeperRelayerFee"
	BidaskProvideJettonOp              JettonOpName = "BidaskProvide"
	CoffeeCrossDexFailureJettonOp      JettonOpName = "CoffeeCrossDexFailure"
	CoffeeSwapJettonOp                 JettonOpName = "CoffeeSwap"
	CoffeeCreatePoolJettonOp           JettonOpName = "CoffeeCreatePool"
	CoffeeDepositLiquidityJettonOp     JettonOpName = "CoffeeDepositLiquidity"
	CoffeeNotificationJettonOp         JettonOpName = "CoffeeNotification"
	StonfiSwapOkJettonOp               JettonOpName = "StonfiSwapOk"
	DedustSwapJettonOp                 JettonOpName = "DedustSwap"
	CoffeeMevProtectFailedSwapJettonOp JettonOpName = "CoffeeMevProtectFailedSwap"
	BidaskSwapJettonOp                 JettonOpName = "BidaskSwap"
	StonfiProvideLiquidityJettonOp     JettonOpName = "StonfiProvideLiquidity"

	CoffeeStakingLockJettonOpCode          JettonOpCode = 0x0c0ffede
	TextCommentJettonOpCode                JettonOpCode = 0x00000000
	TegroJettonSwapJettonOpCode            JettonOpCode = 0x01fb7a25
	EncryptedTextCommentJettonOpCode       JettonOpCode = 0x2167da4b
	StonfiSwapJettonOpCode                 JettonOpCode = 0x25938561
	TegroAddLiquidityJettonOpCode          JettonOpCode = 0x287e167a
	StonfiProvideLpV2JettonOpCode          JettonOpCode = 0x37c096df
	BidaskProvideBothJettonOpCode          JettonOpCode = 0x3ea0bafc
	DedustDepositLiquidityJettonOpCode     JettonOpCode = 0x40e108d6
	StonfiSwapOkRefJettonOpCode            JettonOpCode = 0x45078540
	CoffeeCrossDexResendJettonOpCode       JettonOpCode = 0x4ee9b106
	StonfiSwapV2JettonOpCode               JettonOpCode = 0x6664de2a
	InvoicePayloadJettonOpCode             JettonOpCode = 0x7aa23eb5
	TonkeeperRelayerFeeJettonOpCode        JettonOpCode = 0x878da6e3
	BidaskProvideJettonOpCode              JettonOpCode = 0x96feef7b
	CoffeeCrossDexFailureJettonOpCode      JettonOpCode = 0xb902e61a
	CoffeeSwapJettonOpCode                 JettonOpCode = 0xc0ffee10
	CoffeeCreatePoolJettonOpCode           JettonOpCode = 0xc0ffee11
	CoffeeDepositLiquidityJettonOpCode     JettonOpCode = 0xc0ffee12
	CoffeeNotificationJettonOpCode         JettonOpCode = 0xc0ffee36
	StonfiSwapOkJettonOpCode               JettonOpCode = 0xc64370e5
	DedustSwapJettonOpCode                 JettonOpCode = 0xe3a0d482
	CoffeeMevProtectFailedSwapJettonOpCode JettonOpCode = 0xee51ce51
	BidaskSwapJettonOpCode                 JettonOpCode = 0xf2ef6c1b
	StonfiProvideLiquidityJettonOpCode     JettonOpCode = 0xfcf9e58f
)

var KnownJettonTypes = map[string]any{
	CoffeeStakingLockJettonOp:          CoffeeStakingLockJettonPayload{},
	TextCommentJettonOp:                TextCommentJettonPayload{},
	TegroJettonSwapJettonOp:            TegroJettonSwapJettonPayload{},
	EncryptedTextCommentJettonOp:       EncryptedTextCommentJettonPayload{},
	StonfiSwapJettonOp:                 StonfiSwapJettonPayload{},
	TegroAddLiquidityJettonOp:          TegroAddLiquidityJettonPayload{},
	StonfiProvideLpV2JettonOp:          StonfiProvideLpV2JettonPayload{},
	BidaskProvideBothJettonOp:          BidaskProvideBothJettonPayload{},
	DedustDepositLiquidityJettonOp:     DedustDepositLiquidityJettonPayload{},
	StonfiSwapOkRefJettonOp:            StonfiSwapOkRefJettonPayload{},
	CoffeeCrossDexResendJettonOp:       CoffeeCrossDexResendJettonPayload{},
	StonfiSwapV2JettonOp:               StonfiSwapV2JettonPayload{},
	InvoicePayloadJettonOp:             InvoicePayloadJettonPayload{},
	TonkeeperRelayerFeeJettonOp:        TonkeeperRelayerFeeJettonPayload{},
	BidaskProvideJettonOp:              BidaskProvideJettonPayload{},
	CoffeeCrossDexFailureJettonOp:      CoffeeCrossDexFailureJettonPayload{},
	CoffeeSwapJettonOp:                 CoffeeSwapJettonPayload{},
	CoffeeCreatePoolJettonOp:           CoffeeCreatePoolJettonPayload{},
	CoffeeDepositLiquidityJettonOp:     CoffeeDepositLiquidityJettonPayload{},
	CoffeeNotificationJettonOp:         CoffeeNotificationJettonPayload{},
	StonfiSwapOkJettonOp:               StonfiSwapOkJettonPayload{},
	DedustSwapJettonOp:                 DedustSwapJettonPayload{},
	CoffeeMevProtectFailedSwapJettonOp: CoffeeMevProtectFailedSwapJettonPayload{},
	BidaskSwapJettonOp:                 BidaskSwapJettonPayload{},
	StonfiProvideLiquidityJettonOp:     StonfiProvideLiquidityJettonPayload{},
}
var JettonOpCodes = map[JettonOpName]JettonOpCode{
	CoffeeStakingLockJettonOp:          CoffeeStakingLockJettonOpCode,
	TextCommentJettonOp:                TextCommentJettonOpCode,
	TegroJettonSwapJettonOp:            TegroJettonSwapJettonOpCode,
	EncryptedTextCommentJettonOp:       EncryptedTextCommentJettonOpCode,
	StonfiSwapJettonOp:                 StonfiSwapJettonOpCode,
	TegroAddLiquidityJettonOp:          TegroAddLiquidityJettonOpCode,
	StonfiProvideLpV2JettonOp:          StonfiProvideLpV2JettonOpCode,
	BidaskProvideBothJettonOp:          BidaskProvideBothJettonOpCode,
	DedustDepositLiquidityJettonOp:     DedustDepositLiquidityJettonOpCode,
	StonfiSwapOkRefJettonOp:            StonfiSwapOkRefJettonOpCode,
	CoffeeCrossDexResendJettonOp:       CoffeeCrossDexResendJettonOpCode,
	StonfiSwapV2JettonOp:               StonfiSwapV2JettonOpCode,
	InvoicePayloadJettonOp:             InvoicePayloadJettonOpCode,
	TonkeeperRelayerFeeJettonOp:        TonkeeperRelayerFeeJettonOpCode,
	BidaskProvideJettonOp:              BidaskProvideJettonOpCode,
	CoffeeCrossDexFailureJettonOp:      CoffeeCrossDexFailureJettonOpCode,
	CoffeeSwapJettonOp:                 CoffeeSwapJettonOpCode,
	CoffeeCreatePoolJettonOp:           CoffeeCreatePoolJettonOpCode,
	CoffeeDepositLiquidityJettonOp:     CoffeeDepositLiquidityJettonOpCode,
	CoffeeNotificationJettonOp:         CoffeeNotificationJettonOpCode,
	StonfiSwapOkJettonOp:               StonfiSwapOkJettonOpCode,
	DedustSwapJettonOp:                 DedustSwapJettonOpCode,
	CoffeeMevProtectFailedSwapJettonOp: CoffeeMevProtectFailedSwapJettonOpCode,
	BidaskSwapJettonOp:                 BidaskSwapJettonOpCode,
	StonfiProvideLiquidityJettonOp:     StonfiProvideLiquidityJettonOpCode,
}

var funcJettonDecodersMapping = map[JettonOpCode]func(*JettonPayload, *boc.Cell) error{
	CoffeeStakingLockJettonOpCode:          decodeCoffeeStakingLockJettonOpJetton,
	TextCommentJettonOpCode:                decodeTextCommentJettonOpJetton,
	TegroJettonSwapJettonOpCode:            decodeTegroJettonSwapJettonOpJetton,
	EncryptedTextCommentJettonOpCode:       decodeEncryptedTextCommentJettonOpJetton,
	StonfiSwapJettonOpCode:                 decodeStonfiSwapJettonOpJetton,
	TegroAddLiquidityJettonOpCode:          decodeTegroAddLiquidityJettonOpJetton,
	StonfiProvideLpV2JettonOpCode:          decodeStonfiProvideLpV2JettonOpJetton,
	BidaskProvideBothJettonOpCode:          decodeBidaskProvideBothJettonOpJetton,
	DedustDepositLiquidityJettonOpCode:     decodeDedustDepositLiquidityJettonOpJetton,
	StonfiSwapOkRefJettonOpCode:            decodeStonfiSwapOkRefJettonOpJetton,
	CoffeeCrossDexResendJettonOpCode:       decodeCoffeeCrossDexResendJettonOpJetton,
	StonfiSwapV2JettonOpCode:               decodeStonfiSwapV2JettonOpJetton,
	InvoicePayloadJettonOpCode:             decodeInvoicePayloadJettonOpJetton,
	TonkeeperRelayerFeeJettonOpCode:        decodeTonkeeperRelayerFeeJettonOpJetton,
	BidaskProvideJettonOpCode:              decodeBidaskProvideJettonOpJetton,
	CoffeeCrossDexFailureJettonOpCode:      decodeCoffeeCrossDexFailureJettonOpJetton,
	CoffeeSwapJettonOpCode:                 decodeCoffeeSwapJettonOpJetton,
	CoffeeCreatePoolJettonOpCode:           decodeCoffeeCreatePoolJettonOpJetton,
	CoffeeDepositLiquidityJettonOpCode:     decodeCoffeeDepositLiquidityJettonOpJetton,
	CoffeeNotificationJettonOpCode:         decodeCoffeeNotificationJettonOpJetton,
	StonfiSwapOkJettonOpCode:               decodeStonfiSwapOkJettonOpJetton,
	DedustSwapJettonOpCode:                 decodeDedustSwapJettonOpJetton,
	CoffeeMevProtectFailedSwapJettonOpCode: decodeCoffeeMevProtectFailedSwapJettonOpJetton,
	BidaskSwapJettonOpCode:                 decodeBidaskSwapJettonOpJetton,
	StonfiProvideLiquidityJettonOpCode:     decodeStonfiProvideLiquidityJettonOpJetton,
}

type CoffeeStakingLockJettonPayload struct {
	PeriodId uint32
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

type BidaskProvideBothJettonPayload struct {
	QueryId        uint64
	TonAmount      tlb.Grams
	DepositType    tlb.Uint4
	LiquidityDict  tlb.HashmapE[tlb.Uint32, int32]
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
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

type CoffeeCrossDexResendJettonPayload struct {
	Next tlb.Any `tlb:"^"`
}

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

type BidaskProvideJettonPayload struct {
	QueryId        uint64
	DepositType    tlb.Uint4
	LiquidityDict  tlb.HashmapE[tlb.Uint32, int32]
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type CoffeeCrossDexFailureJettonPayload struct {
	QueryId   uint64
	Recipient tlb.MsgAddress
}

type CoffeeSwapJettonPayload struct {
	Field0 CoffeeSwapStepParams
	Params CoffeeSwapParams `tlb:"^"`
}

type CoffeeCreatePoolJettonPayload struct {
	Params         CoffeePoolParams
	CreationParams CoffeePoolCreationParams
}

type CoffeeDepositLiquidityJettonPayload struct {
	Params CoffeeDepositLiquidityParams
}

type CoffeeNotificationJettonPayload struct {
	QueryId uint64
	Body    tlb.Any `tlb:"^"`
}

type StonfiSwapOkJettonPayload struct{}

type DedustSwapJettonPayload struct {
	Step       DedustSwapStep
	SwapParams DedustSwapParams `tlb:"^"`
}

type CoffeeMevProtectFailedSwapJettonPayload struct {
	QueryId   uint64
	Recipient tlb.MsgAddress
}

type BidaskSwapJettonPayload struct {
	QueryId        uint64
	ToAddress      tlb.MsgAddress
	Slippage       tlb.Either[tlb.Grams, tlb.Uint256]
	ExactOut       tlb.Grams
	RefAddress     tlb.MsgAddress
	AdditionalData *tlb.Any `tlb:"maybe^"`
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type StonfiProvideLiquidityJettonPayload struct {
	TokenWallet tlb.MsgAddress
	MinLpOut    tlb.VarUInteger16
}
