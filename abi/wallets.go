package abi

import (
	"fmt"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

const (
	// Wallet is an abstract interface,
	// any wallet in the blockchain has a concrete version like v1R1, v4R1 but
	// whenever a contract implements any specific wallet interface, this one will be added too.
	Wallet     ContractInterface = "wallet"
	WalletV1R1 ContractInterface = "wallet_v1R1"
	WalletV1R2 ContractInterface = "wallet_v1R2"
	WalletV1R3 ContractInterface = "wallet_v1R3"
	WalletV2R1 ContractInterface = "wallet_v2R1"
	WalletV2R2 ContractInterface = "wallet_v2R2"
	WalletV3R1 ContractInterface = "wallet_v3R1"
	WalletV4R1 ContractInterface = "wallet_v4R1"
	// WalletV4 is an abstract interface, added once a wallet implements any of v4R* versions.
	WalletV4 ContractInterface = "wallet_v4"
	// WalletHighload is an abstract interface, added once a wallet implements any of HighLoad versions.
	WalletHighload ContractInterface = "wallet_highload"
)

var walletInterfaces = map[ContractInterface]struct{}{
	Wallet:         {},
	WalletV1R1:     {},
	WalletV1R2:     {},
	WalletV1R3:     {},
	WalletV2R1:     {},
	WalletV2R2:     {},
	WalletV3R1:     {},
	WalletV3R2:     {},
	WalletV4:       {},
	WalletV4R1:     {},
	WalletV4R2:     {},
	WalletHighload: {},
}

// IsWallet returns true if the given interface is one of the wallet interfaces.
func IsWallet(i ContractInterface) bool {
	_, ok := walletInterfaces[i]
	return ok
}

func bitsToHex(code tlb.Bits256) string {
	return fmt.Sprintf("%x", code[:])
}

var walletsByHashCode = map[string][]ContractInterface{
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R1)):         {WalletV1R1},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R2)):         {WalletV1R2},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V1R3)):         {WalletV1R3},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V2R1)):         {WalletV2R1},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V2R2)):         {WalletV2R2},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V3R1)):         {WalletV3R1},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V3R2)):         {WalletV3R2},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V4R1)):         {WalletV4R1, WalletV4},
	bitsToHex(wallet.GetCodeHashByVer(wallet.V4R2)):         {WalletV4R2, WalletV4},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV1R1)): {WalletHighload},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV1R2)): {WalletHighload},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV2R1)): {WalletHighload},
	bitsToHex(wallet.GetCodeHashByVer(wallet.HighLoadV2R2)): {WalletHighload},
}
