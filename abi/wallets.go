package abi

import (
	"github.com/tonkeeper/tongo/ton"
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
	WalletV3R2 ContractInterface = "wallet_v3R2"
	WalletV4R1 ContractInterface = "wallet_v4R1"
	WalletV4R2 ContractInterface = "wallet_v4R2"
	// WalletHighload is an abstract interface, added once a wallet implements any of HighLoad versions.
	WalletHighload ContractInterface = "wallet_highload"
)

// IsWallet returns true if the given interface is one of the wallet interfaces.
func (c ContractInterface) IsWallet() bool {
	switch c {
	case WalletV1R1,
		WalletV1R2,
		WalletV1R3,
		WalletV2R1,
		WalletV2R2,
		WalletV3R1,
		WalletV3R2,
		WalletV4R1,
		WalletV4R2,
		WalletHighload:
		return true
	}
	return false
}

type knownContractDescription struct {
	contractInterfaces []ContractInterface
	getMethods         []InvokeFn
}

var knownContracts = map[ton.Bits256]knownContractDescription{ //todo: add getmethods and popular contacts
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V1R1)):         {contractInterfaces: []ContractInterface{Wallet, WalletV1R1}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V1R2)):         {contractInterfaces: []ContractInterface{Wallet, WalletV1R2}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V1R3)):         {contractInterfaces: []ContractInterface{Wallet, WalletV1R3}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V2R1)):         {contractInterfaces: []ContractInterface{Wallet, WalletV2R1}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V2R2)):         {contractInterfaces: []ContractInterface{Wallet, WalletV2R2}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V3R1)):         {contractInterfaces: []ContractInterface{Wallet, WalletV3R1}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V3R2)):         {contractInterfaces: []ContractInterface{Wallet, WalletV3R2}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V4R1)):         {contractInterfaces: []ContractInterface{Wallet, WalletV4R1}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.V4R2)):         {contractInterfaces: []ContractInterface{Wallet, WalletV4R2}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.HighLoadV1R1)): {contractInterfaces: []ContractInterface{Wallet, WalletHighload}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.HighLoadV1R2)): {contractInterfaces: []ContractInterface{Wallet, WalletHighload}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.HighLoadV2R1)): {contractInterfaces: []ContractInterface{Wallet, WalletHighload}},
	ton.Bits256(wallet.GetCodeHashByVer(wallet.HighLoadV2R2)): {contractInterfaces: []ContractInterface{Wallet, WalletHighload}},
}
