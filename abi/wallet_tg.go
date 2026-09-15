package abi

import (
	"context"

	abiWalletTg "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/walletTg"
	"github.com/tonkeeper/tongo/ton"
)

// walletTgCodeHash telegram wallet contract code is expected to be stable for a long time because of dynamic dispatch to config param -123
// https://github.com/ton-blockchain/tg-wallet-contract So we can store this hash indefinetely
var walletTgCodeHash = ton.MustParseHash("9149ae51c1e4689710cebf7830297b16acfbadb363a920a537893e7ffeeca768")

func init() {
	knownContracts[walletTgCodeHash] = knownContractDescription{
		contractInterfaces: []ContractInterface{WalletTg},
		getMethods: []InvokeFn{
			func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiWalletTg.GetRevision(ctx, executor, id)
				return "GetRevision_WalletTgResult", r, err
			},
			func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiWalletTg.GetSeqno(ctx, executor, id)
				return "GetSeqno_WalletTgResult", r, err
			},
			func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiWalletTg.GetSubwalletId(ctx, executor, id)
				return "GetSubwalletId_WalletTgResult", r, err
			},
			func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiWalletTg.GetPublicKey(ctx, executor, id)
				return "GetPublicKey_WalletTgResult", r, err
			},
		},
	}
}
