// Code generated - DO NOT EDIT.

package abi

import (
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
)

const (
	FfVaultStableDepositJettonOp     JettonOpName = "FfVaultStableDeposit"
	FfVaultAssetDepositJettonOp      JettonOpName = "FfVaultAssetDeposit"
	FfVaultStableDepositJettonOpCode JettonOpCode = 0x3d669db7
	FfVaultAssetDepositJettonOpCode  JettonOpCode = 0x534917c3
)

var tolkKnownJettonTypes = map[JettonOpName]any{
	FfVaultStableDepositJettonOp: abiFfVault.StableDeposit{},
	FfVaultAssetDepositJettonOp:  abiFfVault.AssetDeposit{},
}

var tolkJettonOpCodes = map[JettonOpName]JettonOpCode{
	FfVaultStableDepositJettonOp: FfVaultStableDepositJettonOpCode,
	FfVaultAssetDepositJettonOp:  FfVaultAssetDepositJettonOpCode,
}

var tolkJettonDecodersMapping = map[JettonOpCode]jettonDecoder{
	FfVaultStableDepositJettonOpCode: decodeJettonPayload[abiFfVault.StableDeposit](FfVaultStableDepositJettonOp, FfVaultStableDepositJettonOpCode, false, false),
	FfVaultAssetDepositJettonOpCode:  decodeJettonPayload[abiFfVault.AssetDeposit](FfVaultAssetDepositJettonOp, FfVaultAssetDepositJettonOpCode, false, false),
}
