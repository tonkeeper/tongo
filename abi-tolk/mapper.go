package abi_tolk

import (
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/tolk/parser"

	_ "embed"
)

//go:embed schemas/jetton_minter.json
var jettonMinterABI []byte

//go:embed schemas/jetton_wallet.json
var jettonWalletABI []byte

var InterfaceToABI = map[abi.ContractInterface]parser.ABI{
	abi.JettonMaster: parser.MustParseABI(jettonMinterABI),
	abi.JettonWallet: parser.MustParseABI(jettonWalletABI),
}
