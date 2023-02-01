package abi

import (
	"fmt"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

const (
	WalletInterface         ContractInterface = "wallet"
	WalletV1R1Interface     ContractInterface = "wallet_v1R1"
	WalletV1R2Interface     ContractInterface = "wallet_v1R2"
	WalletV1R3Interface     ContractInterface = "wallet_v1R3"
	WalletV2R1Interface     ContractInterface = "wallet_v2R1"
	WalletV2R2Interface     ContractInterface = "wallet_v2R2"
	WalletV3R1Interface     ContractInterface = "wallet_v3R1"
	WalletV3R2Interface     ContractInterface = "wallet_v3R2"
	WalletV4R1Interface     ContractInterface = "wallet_v4R1"
	WalletV4R2Interface     ContractInterface = "wallet_v4R2"
	WalletV4Interface       ContractInterface = "wallet_v4"
	WalletHighloadInterface ContractInterface = "wallet_highload"
)

func bitsToHex(code tlb.Bits256) string {
	return fmt.Sprintf("%x", code[:])
}

var walletsByHashCode = map[string][]ContractInterface{
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R1)):         {WalletInterface, WalletV1R1Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R2)):         {WalletInterface, WalletV1R2Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R3)):         {WalletInterface, WalletV1R3Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V2R1)):         {WalletInterface, WalletV2R1Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V2R2)):         {WalletInterface, WalletV2R2Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V3R1)):         {WalletInterface, WalletV3R1Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V3R2)):         {WalletInterface, WalletV3R2Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V4R1)):         {WalletInterface, WalletV4R1Interface, WalletV4Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V4R2)):         {WalletInterface, WalletV4R2Interface, WalletV4Interface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV1R1)): {WalletInterface, WalletHighloadInterface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV1R2)): {WalletInterface, WalletHighloadInterface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV2R1)): {WalletInterface, WalletHighloadInterface},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV2R2)): {WalletInterface, WalletHighloadInterface},
}
