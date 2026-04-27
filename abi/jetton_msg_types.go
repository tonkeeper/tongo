// Code generated - DO NOT EDIT.

package abi

import (
	"github.com/tonkeeper/tongo/tlb"
)

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
	PoolFundAccountJettonOp            JettonOpName = "PoolFundAccount"
	StonfiSwapOkRefJettonOp            JettonOpName = "StonfiSwapOkRef"
	CoffeeCrossDexResendJettonOp       JettonOpName = "CoffeeCrossDexResend"
	BidaskDammProvideJettonOp          JettonOpName = "BidaskDammProvide"
	StonfiSwapV2JettonOp               JettonOpName = "StonfiSwapV2"
	BidaskDammProvideOneSideJettonOp   JettonOpName = "BidaskDammProvideOneSide"
	StormDepositJettonJettonOp         JettonOpName = "StormDepositJetton"
	InvoicePayloadJettonOp             JettonOpName = "InvoicePayload"
	TonkeeperRelayerFeeJettonOp        JettonOpName = "TonkeeperRelayerFee"
	BidaskSwapV2JettonOp               JettonOpName = "BidaskSwapV2"
	MoonBoostPoolJettonOp              JettonOpName = "MoonBoostPool"
	BidaskProvideJettonOp              JettonOpName = "BidaskProvide"
	MoonFillOrderJettonOp              JettonOpName = "MoonFillOrder"
	BidaskDammProvideBothJettonOp      JettonOpName = "BidaskDammProvideBoth"
	MoonDepositLiquidityJettonOp       JettonOpName = "MoonDepositLiquidity"
	MoonSwapJettonOp                   JettonOpName = "MoonSwap"
	CoffeeCrossDexFailureJettonOp      JettonOpName = "CoffeeCrossDexFailure"
	CoffeeSwapJettonOp                 JettonOpName = "CoffeeSwap"
	CoffeeCreatePoolJettonOp           JettonOpName = "CoffeeCreatePool"
	CoffeeDepositLiquidityJettonOp     JettonOpName = "CoffeeDepositLiquidity"
	CoffeeNotificationJettonOp         JettonOpName = "CoffeeNotification"
	MoonSwapFailedJettonOp             JettonOpName = "MoonSwapFailed"
	StonfiSwapOkJettonOp               JettonOpName = "StonfiSwapOk"
	StormStakeJettonOp                 JettonOpName = "StormStake"
	WithdrawPayloadJettonOp            JettonOpName = "WithdrawPayload"
	MoonSwapSucceedJettonOp            JettonOpName = "MoonSwapSucceed"
	MoonCreateOrderJettonOp            JettonOpName = "MoonCreateOrder"
	BidaskDammSwapJettonOp             JettonOpName = "BidaskDammSwap"
	DedustSwapJettonOp                 JettonOpName = "DedustSwap"
	CoffeeMevProtectFailedSwapJettonOp JettonOpName = "CoffeeMevProtectFailedSwap"
	BidaskSwapJettonOp                 JettonOpName = "BidaskSwap"
	DepositPayloadJettonOp             JettonOpName = "DepositPayload"
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
	PoolFundAccountJettonOpCode            JettonOpCode = 0x4468de77
	StonfiSwapOkRefJettonOpCode            JettonOpCode = 0x45078540
	CoffeeCrossDexResendJettonOpCode       JettonOpCode = 0x4ee9b106
	BidaskDammProvideJettonOpCode          JettonOpCode = 0x63ec24ae
	StonfiSwapV2JettonOpCode               JettonOpCode = 0x6664de2a
	BidaskDammProvideOneSideJettonOpCode   JettonOpCode = 0x729c04c8
	StormDepositJettonJettonOpCode         JettonOpCode = 0x76840119
	InvoicePayloadJettonOpCode             JettonOpCode = 0x7aa23eb5
	TonkeeperRelayerFeeJettonOpCode        JettonOpCode = 0x878da6e3
	BidaskSwapV2JettonOpCode               JettonOpCode = 0x87d36990
	MoonBoostPoolJettonOpCode              JettonOpCode = 0x96aa1586
	BidaskProvideJettonOpCode              JettonOpCode = 0x96feef7b
	MoonFillOrderJettonOpCode              JettonOpCode = 0x99b49842
	BidaskDammProvideBothJettonOpCode      JettonOpCode = 0xa8904134
	MoonDepositLiquidityJettonOpCode       JettonOpCode = 0xb31db781
	MoonSwapJettonOpCode                   JettonOpCode = 0xb37a900b
	CoffeeCrossDexFailureJettonOpCode      JettonOpCode = 0xb902e61a
	CoffeeSwapJettonOpCode                 JettonOpCode = 0xc0ffee10
	CoffeeCreatePoolJettonOpCode           JettonOpCode = 0xc0ffee11
	CoffeeDepositLiquidityJettonOpCode     JettonOpCode = 0xc0ffee12
	CoffeeNotificationJettonOpCode         JettonOpCode = 0xc0ffee36
	MoonSwapFailedJettonOpCode             JettonOpCode = 0xc47c1f57
	StonfiSwapOkJettonOpCode               JettonOpCode = 0xc64370e5
	StormStakeJettonOpCode                 JettonOpCode = 0xc89a3ee4
	WithdrawPayloadJettonOpCode            JettonOpCode = 0xcb03bfaf
	MoonSwapSucceedJettonOpCode            JettonOpCode = 0xcb7f38d6
	MoonCreateOrderJettonOpCode            JettonOpCode = 0xda067c19
	BidaskDammSwapJettonOpCode             JettonOpCode = 0xdd79732c
	DedustSwapJettonOpCode                 JettonOpCode = 0xe3a0d482
	CoffeeMevProtectFailedSwapJettonOpCode JettonOpCode = 0xee51ce51
	BidaskSwapJettonOpCode                 JettonOpCode = 0xf2ef6c1b
	DepositPayloadJettonOpCode             JettonOpCode = 0xf9471134
	StonfiProvideLiquidityJettonOpCode     JettonOpCode = 0xfcf9e58f
)

var xmlKnownJettonTypes = map[JettonOpName]any{
	CoffeeStakingLockJettonOp:          CoffeeStakingLockJettonPayload{},
	TextCommentJettonOp:                TextCommentJettonPayload{},
	TegroJettonSwapJettonOp:            TegroJettonSwapJettonPayload{},
	EncryptedTextCommentJettonOp:       EncryptedTextCommentJettonPayload{},
	StonfiSwapJettonOp:                 StonfiSwapJettonPayload{},
	TegroAddLiquidityJettonOp:          TegroAddLiquidityJettonPayload{},
	StonfiProvideLpV2JettonOp:          StonfiProvideLpV2JettonPayload{},
	BidaskProvideBothJettonOp:          BidaskProvideBothJettonPayload{},
	DedustDepositLiquidityJettonOp:     DedustDepositLiquidityJettonPayload{},
	PoolFundAccountJettonOp:            PoolFundAccountJettonPayload{},
	StonfiSwapOkRefJettonOp:            StonfiSwapOkRefJettonPayload{},
	CoffeeCrossDexResendJettonOp:       CoffeeCrossDexResendJettonPayload{},
	BidaskDammProvideJettonOp:          BidaskDammProvideJettonPayload{},
	StonfiSwapV2JettonOp:               StonfiSwapV2JettonPayload{},
	BidaskDammProvideOneSideJettonOp:   BidaskDammProvideOneSideJettonPayload{},
	StormDepositJettonJettonOp:         StormDepositJettonJettonPayload{},
	InvoicePayloadJettonOp:             InvoicePayloadJettonPayload{},
	TonkeeperRelayerFeeJettonOp:        TonkeeperRelayerFeeJettonPayload{},
	BidaskSwapV2JettonOp:               BidaskSwapV2JettonPayload{},
	MoonBoostPoolJettonOp:              MoonBoostPoolJettonPayload{},
	BidaskProvideJettonOp:              BidaskProvideJettonPayload{},
	MoonFillOrderJettonOp:              MoonFillOrderJettonPayload{},
	BidaskDammProvideBothJettonOp:      BidaskDammProvideBothJettonPayload{},
	MoonDepositLiquidityJettonOp:       MoonDepositLiquidityJettonPayload{},
	MoonSwapJettonOp:                   MoonSwapJettonPayload{},
	CoffeeCrossDexFailureJettonOp:      CoffeeCrossDexFailureJettonPayload{},
	CoffeeSwapJettonOp:                 CoffeeSwapJettonPayload{},
	CoffeeCreatePoolJettonOp:           CoffeeCreatePoolJettonPayload{},
	CoffeeDepositLiquidityJettonOp:     CoffeeDepositLiquidityJettonPayload{},
	CoffeeNotificationJettonOp:         CoffeeNotificationJettonPayload{},
	MoonSwapFailedJettonOp:             MoonSwapFailedJettonPayload{},
	StonfiSwapOkJettonOp:               StonfiSwapOkJettonPayload{},
	StormStakeJettonOp:                 StormStakeJettonPayload{},
	WithdrawPayloadJettonOp:            WithdrawPayloadJettonPayload{},
	MoonSwapSucceedJettonOp:            MoonSwapSucceedJettonPayload{},
	MoonCreateOrderJettonOp:            MoonCreateOrderJettonPayload{},
	BidaskDammSwapJettonOp:             BidaskDammSwapJettonPayload{},
	DedustSwapJettonOp:                 DedustSwapJettonPayload{},
	CoffeeMevProtectFailedSwapJettonOp: CoffeeMevProtectFailedSwapJettonPayload{},
	BidaskSwapJettonOp:                 BidaskSwapJettonPayload{},
	DepositPayloadJettonOp:             DepositPayloadJettonPayload{},
	StonfiProvideLiquidityJettonOp:     StonfiProvideLiquidityJettonPayload{},
}

var xmlJettonOpCodes = map[JettonOpName]JettonOpCode{
	CoffeeStakingLockJettonOp:          CoffeeStakingLockJettonOpCode,
	TextCommentJettonOp:                TextCommentJettonOpCode,
	TegroJettonSwapJettonOp:            TegroJettonSwapJettonOpCode,
	EncryptedTextCommentJettonOp:       EncryptedTextCommentJettonOpCode,
	StonfiSwapJettonOp:                 StonfiSwapJettonOpCode,
	TegroAddLiquidityJettonOp:          TegroAddLiquidityJettonOpCode,
	StonfiProvideLpV2JettonOp:          StonfiProvideLpV2JettonOpCode,
	BidaskProvideBothJettonOp:          BidaskProvideBothJettonOpCode,
	DedustDepositLiquidityJettonOp:     DedustDepositLiquidityJettonOpCode,
	PoolFundAccountJettonOp:            PoolFundAccountJettonOpCode,
	StonfiSwapOkRefJettonOp:            StonfiSwapOkRefJettonOpCode,
	CoffeeCrossDexResendJettonOp:       CoffeeCrossDexResendJettonOpCode,
	BidaskDammProvideJettonOp:          BidaskDammProvideJettonOpCode,
	StonfiSwapV2JettonOp:               StonfiSwapV2JettonOpCode,
	BidaskDammProvideOneSideJettonOp:   BidaskDammProvideOneSideJettonOpCode,
	StormDepositJettonJettonOp:         StormDepositJettonJettonOpCode,
	InvoicePayloadJettonOp:             InvoicePayloadJettonOpCode,
	TonkeeperRelayerFeeJettonOp:        TonkeeperRelayerFeeJettonOpCode,
	BidaskSwapV2JettonOp:               BidaskSwapV2JettonOpCode,
	MoonBoostPoolJettonOp:              MoonBoostPoolJettonOpCode,
	BidaskProvideJettonOp:              BidaskProvideJettonOpCode,
	MoonFillOrderJettonOp:              MoonFillOrderJettonOpCode,
	BidaskDammProvideBothJettonOp:      BidaskDammProvideBothJettonOpCode,
	MoonDepositLiquidityJettonOp:       MoonDepositLiquidityJettonOpCode,
	MoonSwapJettonOp:                   MoonSwapJettonOpCode,
	CoffeeCrossDexFailureJettonOp:      CoffeeCrossDexFailureJettonOpCode,
	CoffeeSwapJettonOp:                 CoffeeSwapJettonOpCode,
	CoffeeCreatePoolJettonOp:           CoffeeCreatePoolJettonOpCode,
	CoffeeDepositLiquidityJettonOp:     CoffeeDepositLiquidityJettonOpCode,
	CoffeeNotificationJettonOp:         CoffeeNotificationJettonOpCode,
	MoonSwapFailedJettonOp:             MoonSwapFailedJettonOpCode,
	StonfiSwapOkJettonOp:               StonfiSwapOkJettonOpCode,
	StormStakeJettonOp:                 StormStakeJettonOpCode,
	WithdrawPayloadJettonOp:            WithdrawPayloadJettonOpCode,
	MoonSwapSucceedJettonOp:            MoonSwapSucceedJettonOpCode,
	MoonCreateOrderJettonOp:            MoonCreateOrderJettonOpCode,
	BidaskDammSwapJettonOp:             BidaskDammSwapJettonOpCode,
	DedustSwapJettonOp:                 DedustSwapJettonOpCode,
	CoffeeMevProtectFailedSwapJettonOp: CoffeeMevProtectFailedSwapJettonOpCode,
	BidaskSwapJettonOp:                 BidaskSwapJettonOpCode,
	DepositPayloadJettonOp:             DepositPayloadJettonOpCode,
	StonfiProvideLiquidityJettonOp:     StonfiProvideLiquidityJettonOpCode,
}

var xmlJettonDecodersMapping = map[JettonOpCode]jettonDecoder{
	CoffeeStakingLockJettonOpCode:          decodeJettonPayload[CoffeeStakingLockJettonPayload](CoffeeStakingLockJettonOp, CoffeeStakingLockJettonOpCode, true, false),
	TextCommentJettonOpCode:                decodeJettonPayload[TextCommentJettonPayload](TextCommentJettonOp, TextCommentJettonOpCode, true, false),
	TegroJettonSwapJettonOpCode:            decodeJettonPayload[TegroJettonSwapJettonPayload](TegroJettonSwapJettonOp, TegroJettonSwapJettonOpCode, true, false),
	EncryptedTextCommentJettonOpCode:       decodeJettonPayload[EncryptedTextCommentJettonPayload](EncryptedTextCommentJettonOp, EncryptedTextCommentJettonOpCode, true, false),
	StonfiSwapJettonOpCode:                 decodeJettonPayload[StonfiSwapJettonPayload](StonfiSwapJettonOp, StonfiSwapJettonOpCode, true, false),
	TegroAddLiquidityJettonOpCode:          decodeJettonPayload[TegroAddLiquidityJettonPayload](TegroAddLiquidityJettonOp, TegroAddLiquidityJettonOpCode, true, false),
	StonfiProvideLpV2JettonOpCode:          decodeJettonPayload[StonfiProvideLpV2JettonPayload](StonfiProvideLpV2JettonOp, StonfiProvideLpV2JettonOpCode, true, false),
	BidaskProvideBothJettonOpCode:          decodeJettonPayload[BidaskProvideBothJettonPayload](BidaskProvideBothJettonOp, BidaskProvideBothJettonOpCode, true, false),
	DedustDepositLiquidityJettonOpCode:     decodeJettonPayload[DedustDepositLiquidityJettonPayload](DedustDepositLiquidityJettonOp, DedustDepositLiquidityJettonOpCode, true, false),
	PoolFundAccountJettonOpCode:            decodeJettonPayload[PoolFundAccountJettonPayload](PoolFundAccountJettonOp, PoolFundAccountJettonOpCode, true, true),
	StonfiSwapOkRefJettonOpCode:            decodeJettonPayload[StonfiSwapOkRefJettonPayload](StonfiSwapOkRefJettonOp, StonfiSwapOkRefJettonOpCode, true, true),
	CoffeeCrossDexResendJettonOpCode:       decodeJettonPayload[CoffeeCrossDexResendJettonPayload](CoffeeCrossDexResendJettonOp, CoffeeCrossDexResendJettonOpCode, true, false),
	BidaskDammProvideJettonOpCode:          decodeJettonPayload[BidaskDammProvideJettonPayload](BidaskDammProvideJettonOp, BidaskDammProvideJettonOpCode, true, false),
	StonfiSwapV2JettonOpCode:               decodeJettonPayload[StonfiSwapV2JettonPayload](StonfiSwapV2JettonOp, StonfiSwapV2JettonOpCode, true, false),
	BidaskDammProvideOneSideJettonOpCode:   decodeJettonPayload[BidaskDammProvideOneSideJettonPayload](BidaskDammProvideOneSideJettonOp, BidaskDammProvideOneSideJettonOpCode, true, false),
	StormDepositJettonJettonOpCode:         decodeJettonPayload[StormDepositJettonJettonPayload](StormDepositJettonJettonOp, StormDepositJettonJettonOpCode, true, false),
	InvoicePayloadJettonOpCode:             decodeJettonPayload[InvoicePayloadJettonPayload](InvoicePayloadJettonOp, InvoicePayloadJettonOpCode, true, false),
	TonkeeperRelayerFeeJettonOpCode:        decodeJettonPayload[TonkeeperRelayerFeeJettonPayload](TonkeeperRelayerFeeJettonOp, TonkeeperRelayerFeeJettonOpCode, true, false),
	BidaskSwapV2JettonOpCode:               decodeJettonPayload[BidaskSwapV2JettonPayload](BidaskSwapV2JettonOp, BidaskSwapV2JettonOpCode, true, false),
	MoonBoostPoolJettonOpCode:              decodeJettonPayload[MoonBoostPoolJettonPayload](MoonBoostPoolJettonOp, MoonBoostPoolJettonOpCode, true, false),
	BidaskProvideJettonOpCode:              decodeJettonPayload[BidaskProvideJettonPayload](BidaskProvideJettonOp, BidaskProvideJettonOpCode, true, false),
	MoonFillOrderJettonOpCode:              decodeJettonPayload[MoonFillOrderJettonPayload](MoonFillOrderJettonOp, MoonFillOrderJettonOpCode, true, false),
	BidaskDammProvideBothJettonOpCode:      decodeJettonPayload[BidaskDammProvideBothJettonPayload](BidaskDammProvideBothJettonOp, BidaskDammProvideBothJettonOpCode, true, false),
	MoonDepositLiquidityJettonOpCode:       decodeJettonPayload[MoonDepositLiquidityJettonPayload](MoonDepositLiquidityJettonOp, MoonDepositLiquidityJettonOpCode, true, false),
	MoonSwapJettonOpCode:                   decodeJettonPayload[MoonSwapJettonPayload](MoonSwapJettonOp, MoonSwapJettonOpCode, true, false),
	CoffeeCrossDexFailureJettonOpCode:      decodeJettonPayload[CoffeeCrossDexFailureJettonPayload](CoffeeCrossDexFailureJettonOp, CoffeeCrossDexFailureJettonOpCode, true, false),
	CoffeeSwapJettonOpCode:                 decodeJettonPayload[CoffeeSwapJettonPayload](CoffeeSwapJettonOp, CoffeeSwapJettonOpCode, true, false),
	CoffeeCreatePoolJettonOpCode:           decodeJettonPayload[CoffeeCreatePoolJettonPayload](CoffeeCreatePoolJettonOp, CoffeeCreatePoolJettonOpCode, true, false),
	CoffeeDepositLiquidityJettonOpCode:     decodeJettonPayload[CoffeeDepositLiquidityJettonPayload](CoffeeDepositLiquidityJettonOp, CoffeeDepositLiquidityJettonOpCode, true, false),
	CoffeeNotificationJettonOpCode:         decodeJettonPayload[CoffeeNotificationJettonPayload](CoffeeNotificationJettonOp, CoffeeNotificationJettonOpCode, true, false),
	MoonSwapFailedJettonOpCode:             decodeJettonPayload[MoonSwapFailedJettonPayload](MoonSwapFailedJettonOp, MoonSwapFailedJettonOpCode, true, false),
	StonfiSwapOkJettonOpCode:               decodeJettonPayload[StonfiSwapOkJettonPayload](StonfiSwapOkJettonOp, StonfiSwapOkJettonOpCode, true, true),
	StormStakeJettonOpCode:                 decodeJettonPayload[StormStakeJettonPayload](StormStakeJettonOp, StormStakeJettonOpCode, true, false),
	WithdrawPayloadJettonOpCode:            decodeJettonPayload[WithdrawPayloadJettonPayload](WithdrawPayloadJettonOp, WithdrawPayloadJettonOpCode, true, true),
	MoonSwapSucceedJettonOpCode:            decodeJettonPayload[MoonSwapSucceedJettonPayload](MoonSwapSucceedJettonOp, MoonSwapSucceedJettonOpCode, true, false),
	MoonCreateOrderJettonOpCode:            decodeJettonPayload[MoonCreateOrderJettonPayload](MoonCreateOrderJettonOp, MoonCreateOrderJettonOpCode, true, false),
	BidaskDammSwapJettonOpCode:             decodeJettonPayload[BidaskDammSwapJettonPayload](BidaskDammSwapJettonOp, BidaskDammSwapJettonOpCode, true, false),
	DedustSwapJettonOpCode:                 decodeJettonPayload[DedustSwapJettonPayload](DedustSwapJettonOp, DedustSwapJettonOpCode, true, false),
	CoffeeMevProtectFailedSwapJettonOpCode: decodeJettonPayload[CoffeeMevProtectFailedSwapJettonPayload](CoffeeMevProtectFailedSwapJettonOp, CoffeeMevProtectFailedSwapJettonOpCode, true, false),
	BidaskSwapJettonOpCode:                 decodeJettonPayload[BidaskSwapJettonPayload](BidaskSwapJettonOp, BidaskSwapJettonOpCode, true, false),
	DepositPayloadJettonOpCode:             decodeJettonPayload[DepositPayloadJettonPayload](DepositPayloadJettonOp, DepositPayloadJettonOpCode, true, true),
	StonfiProvideLiquidityJettonOpCode:     decodeJettonPayload[StonfiProvideLiquidityJettonPayload](StonfiProvideLiquidityJettonOp, StonfiProvideLiquidityJettonOpCode, true, false),
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

type PoolFundAccountJettonPayload struct {
	JettonTarget tlb.MsgAddress
	Enough0      tlb.VarUInteger16
	Enough1      tlb.VarUInteger16
	Liquidity    tlb.Uint128
	TickLower    tlb.Int24
	TickUpper    tlb.Int24
}

type StonfiSwapOkRefJettonPayload struct{}

type CoffeeCrossDexResendJettonPayload struct {
	Next tlb.Any `tlb:"^"`
}

type BidaskDammProvideJettonPayload struct {
	Receiver       tlb.MsgAddress
	LockLiquidity  bool
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
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

type BidaskDammProvideOneSideJettonPayload struct {
	Receiver       tlb.MsgAddress
	LockLiquidity  bool
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type StormDepositJettonJettonPayload struct {
	QueryId         uint64
	ReceiverAddress tlb.MsgAddress
	Init            bool
	KeyInit         InitializationRequest
}

type InvoicePayloadJettonPayload struct {
	Id  tlb.Bits128
	Url PaymentProviderUrl
}

type TonkeeperRelayerFeeJettonPayload struct{}

type BidaskSwapV2JettonPayload struct {
	ToAddress      tlb.MsgAddress
	Slippage       tlb.Either[tlb.Grams, tlb.Uint256]
	ExactOut       tlb.Grams
	AdditionalData *AdditionalData `tlb:"maybe^"`
	RejectPayload  *tlb.Any        `tlb:"maybe^"`
	ForwardPayload *tlb.Any        `tlb:"maybe^"`
}

type MoonBoostPoolJettonPayload struct{}

type BidaskProvideJettonPayload struct {
	DepositType    tlb.Uint4
	LiquidityDict  tlb.HashmapE[tlb.Uint32, int32]
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type MoonFillOrderJettonPayload struct {
	Recipient        tlb.MsgAddress
	RecipientPayload *tlb.Any `tlb:"maybe^"`
	RejectAddress    tlb.MsgAddress
}

type BidaskDammProvideBothJettonPayload struct {
	NativeAmount   tlb.Grams
	Receiver       tlb.MsgAddress
	LockLiquidity  bool
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type MoonDepositLiquidityJettonPayload struct {
	MinLpOut tlb.VarUInteger16
}

type MoonSwapJettonPayload struct {
	SwapParams MoonSwapParams
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

type MoonSwapFailedJettonPayload struct{}

type StonfiSwapOkJettonPayload struct{}

type StormStakeJettonPayload struct{}

type WithdrawPayloadJettonPayload struct {
	AssetAddress     tlb.MsgAddress
	OracleParams     *tlb.Any `tlb:"maybe^"`
	ForwardTonAmount tlb.Grams
	ForwardPayload   *tlb.Any `tlb:"maybe^"`
}

type MoonSwapSucceedJettonPayload struct{}

type MoonCreateOrderJettonPayload struct {
	Asset1Id  tlb.MsgAddress
	Asset2Id  tlb.MsgAddress
	OrderData MoonOrderParams
}

type BidaskDammSwapJettonPayload struct {
	ToAddress      tlb.MsgAddress
	Slippage       tlb.Grams
	FromAddress    tlb.MsgAddress
	ExactOut       tlb.Grams
	AdditionalData *tlb.Any `tlb:"maybe^"`
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type DedustSwapJettonPayload struct {
	Step       DedustSwapStep
	SwapParams DedustSwapParams `tlb:"^"`
}

type CoffeeMevProtectFailedSwapJettonPayload struct {
	QueryId   uint64
	Recipient tlb.MsgAddress
}

type BidaskSwapJettonPayload struct {
	ToAddress      tlb.MsgAddress
	Slippage       tlb.Either[tlb.Grams, tlb.Uint256]
	ExactOut       tlb.Grams
	RefAddress     tlb.MsgAddress
	AdditionalData *tlb.Any `tlb:"maybe^"`
	RejectPayload  *tlb.Any `tlb:"maybe^"`
	ForwardPayload *tlb.Any `tlb:"maybe^"`
}

type DepositPayloadJettonPayload struct {
	OracleParams     *tlb.Any `tlb:"maybe^"`
	ForwardTonAmount tlb.Grams
	ForwardPayload   *tlb.Any `tlb:"maybe^"`
}

type StonfiProvideLiquidityJettonPayload struct {
	TokenWallet tlb.MsgAddress
	MinLpOut    tlb.VarUInteger16
}
