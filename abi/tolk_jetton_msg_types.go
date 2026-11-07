// Code generated - DO NOT EDIT.

package abi

import (
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
	abiMooncx "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/mooncx"
)

const (
	FfVaultStableDepositJettonOp                       JettonOpName = "FfVaultStableDeposit"
	MooncxCpmmPoolBoostPayloadJettonOp                 JettonOpName = "MooncxCpmmPoolBoostPayload"
	FfVaultAssetDepositJettonOp                        JettonOpName = "FfVaultAssetDeposit"
	MooncxCpmmPoolDepositLiquidityPayloadJettonOp      JettonOpName = "MooncxCpmmPoolDepositLiquidityPayload"
	MooncxSwapFulfillPayloadJettonOp                   JettonOpName = "MooncxSwapFulfillPayload"
	MooncxMemepadCurveSellPayloadJettonOp              JettonOpName = "MooncxMemepadCurveSellPayload"
	MooncxClmmRangeSwapPayloadJettonOp                 JettonOpName = "MooncxClmmRangeSwapPayload"
	MooncxClmmRangeDepositLiquidityPayloadJettonOp     JettonOpName = "MooncxClmmRangeDepositLiquidityPayload"
	MooncxCpmmPoolSwapPayloadJettonOp                  JettonOpName = "MooncxCpmmPoolSwapPayload"
	MooncxSwapRejectPayloadJettonOp                    JettonOpName = "MooncxSwapRejectPayload"
	FfVaultStableDepositJettonOpCode                   JettonOpCode = 0x3d669db7
	MooncxCpmmPoolBoostPayloadJettonOpCode             JettonOpCode = 0x50bbec04
	FfVaultAssetDepositJettonOpCode                    JettonOpCode = 0x534917c3
	MooncxCpmmPoolDepositLiquidityPayloadJettonOpCode  JettonOpCode = 0x90af0b63
	MooncxSwapFulfillPayloadJettonOpCode               JettonOpCode = 0x9d142643
	MooncxMemepadCurveSellPayloadJettonOpCode          JettonOpCode = 0xb97a94b7
	MooncxClmmRangeSwapPayloadJettonOpCode             JettonOpCode = 0xd151210f
	MooncxClmmRangeDepositLiquidityPayloadJettonOpCode JettonOpCode = 0xec401442
	MooncxCpmmPoolSwapPayloadJettonOpCode              JettonOpCode = 0xfa45056b
	MooncxSwapRejectPayloadJettonOpCode                JettonOpCode = 0xff32da99
)

var tolkKnownJettonTypes = map[JettonOpName]any{
	FfVaultStableDepositJettonOp:                   abiFfVault.StableDeposit{},
	MooncxCpmmPoolBoostPayloadJettonOp:             abiMooncx.CpmmPoolBoostPayload{},
	FfVaultAssetDepositJettonOp:                    abiFfVault.AssetDeposit{},
	MooncxCpmmPoolDepositLiquidityPayloadJettonOp:  abiMooncx.CpmmPoolDepositLiquidityPayload{},
	MooncxSwapFulfillPayloadJettonOp:               abiMooncx.SwapFulfillPayload{},
	MooncxMemepadCurveSellPayloadJettonOp:          abiMooncx.MemepadCurveSellPayload{},
	MooncxClmmRangeSwapPayloadJettonOp:             abiMooncx.ClmmRangeSwapPayload{},
	MooncxClmmRangeDepositLiquidityPayloadJettonOp: abiMooncx.ClmmRangeDepositLiquidityPayload{},
	MooncxCpmmPoolSwapPayloadJettonOp:              abiMooncx.CpmmPoolSwapPayload{},
	MooncxSwapRejectPayloadJettonOp:                abiMooncx.SwapRejectPayload{},
}

var tolkJettonOpCodes = map[JettonOpName]JettonOpCode{
	FfVaultStableDepositJettonOp:                   FfVaultStableDepositJettonOpCode,
	MooncxCpmmPoolBoostPayloadJettonOp:             MooncxCpmmPoolBoostPayloadJettonOpCode,
	FfVaultAssetDepositJettonOp:                    FfVaultAssetDepositJettonOpCode,
	MooncxCpmmPoolDepositLiquidityPayloadJettonOp:  MooncxCpmmPoolDepositLiquidityPayloadJettonOpCode,
	MooncxSwapFulfillPayloadJettonOp:               MooncxSwapFulfillPayloadJettonOpCode,
	MooncxMemepadCurveSellPayloadJettonOp:          MooncxMemepadCurveSellPayloadJettonOpCode,
	MooncxClmmRangeSwapPayloadJettonOp:             MooncxClmmRangeSwapPayloadJettonOpCode,
	MooncxClmmRangeDepositLiquidityPayloadJettonOp: MooncxClmmRangeDepositLiquidityPayloadJettonOpCode,
	MooncxCpmmPoolSwapPayloadJettonOp:              MooncxCpmmPoolSwapPayloadJettonOpCode,
	MooncxSwapRejectPayloadJettonOp:                MooncxSwapRejectPayloadJettonOpCode,
}

var tolkJettonDecodersMapping = map[JettonOpCode]jettonDecoder{
	FfVaultStableDepositJettonOpCode:                   decodeJettonPayload[abiFfVault.StableDeposit](FfVaultStableDepositJettonOp, FfVaultStableDepositJettonOpCode, false, false),
	MooncxCpmmPoolBoostPayloadJettonOpCode:             decodeJettonPayload[abiMooncx.CpmmPoolBoostPayload](MooncxCpmmPoolBoostPayloadJettonOp, MooncxCpmmPoolBoostPayloadJettonOpCode, false, false),
	FfVaultAssetDepositJettonOpCode:                    decodeJettonPayload[abiFfVault.AssetDeposit](FfVaultAssetDepositJettonOp, FfVaultAssetDepositJettonOpCode, false, false),
	MooncxCpmmPoolDepositLiquidityPayloadJettonOpCode:  decodeJettonPayload[abiMooncx.CpmmPoolDepositLiquidityPayload](MooncxCpmmPoolDepositLiquidityPayloadJettonOp, MooncxCpmmPoolDepositLiquidityPayloadJettonOpCode, false, false),
	MooncxSwapFulfillPayloadJettonOpCode:               decodeJettonPayload[abiMooncx.SwapFulfillPayload](MooncxSwapFulfillPayloadJettonOp, MooncxSwapFulfillPayloadJettonOpCode, false, false),
	MooncxMemepadCurveSellPayloadJettonOpCode:          decodeJettonPayload[abiMooncx.MemepadCurveSellPayload](MooncxMemepadCurveSellPayloadJettonOp, MooncxMemepadCurveSellPayloadJettonOpCode, false, false),
	MooncxClmmRangeSwapPayloadJettonOpCode:             decodeJettonPayload[abiMooncx.ClmmRangeSwapPayload](MooncxClmmRangeSwapPayloadJettonOp, MooncxClmmRangeSwapPayloadJettonOpCode, false, false),
	MooncxClmmRangeDepositLiquidityPayloadJettonOpCode: decodeJettonPayload[abiMooncx.ClmmRangeDepositLiquidityPayload](MooncxClmmRangeDepositLiquidityPayloadJettonOp, MooncxClmmRangeDepositLiquidityPayloadJettonOpCode, false, false),
	MooncxCpmmPoolSwapPayloadJettonOpCode:              decodeJettonPayload[abiMooncx.CpmmPoolSwapPayload](MooncxCpmmPoolSwapPayloadJettonOp, MooncxCpmmPoolSwapPayloadJettonOpCode, false, false),
	MooncxSwapRejectPayloadJettonOpCode:                decodeJettonPayload[abiMooncx.SwapRejectPayload](MooncxSwapRejectPayloadJettonOp, MooncxSwapRejectPayloadJettonOpCode, false, false),
}
