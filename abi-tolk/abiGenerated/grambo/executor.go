// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type StorageExecutor interface {
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

type GramboJettonMaster interface {
	WithAccountId(accountID ton.AccountID) GramboJettonMasterWithAccount
	GetAmountOut(ctx context.Context, reqAccountID ton.AccountID, tonIn tlb.Uint64) (tlb.Int257, error)
	GetTonOut(ctx context.Context, reqAccountID ton.AccountID, tokensIn tlb.Int257) (tlb.Uint64, error)
	GetCurveState(ctx context.Context, reqAccountID ton.AccountID) (GetCurveStateResult, error)
	GetJettonData(ctx context.Context, reqAccountID ton.AccountID) (GetJettonDataResult, error)
	GetWalletAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, GramboMasterStorage, error)
}

type GramboJettonMasterWithAccount interface {
	GetAmountOut(ctx context.Context, tonIn tlb.Uint64) (tlb.Int257, error)
	GetTonOut(ctx context.Context, tokensIn tlb.Int257) (tlb.Uint64, error)
	GetCurveState(ctx context.Context) (GetCurveStateResult, error)
	GetJettonData(ctx context.Context) (GetJettonDataResult, error)
	GetWalletAddress(ctx context.Context, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, GramboMasterStorage, error)
}

type GramboJettonWallet interface {
	WithAccountId(accountID ton.AccountID) GramboJettonWalletWithAccount
	GetWalletStatus(ctx context.Context, reqAccountID ton.AccountID) (bool, error)
	GetWalletData(ctx context.Context, reqAccountID ton.AccountID) (GetWalletDataResult, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, GramboWalletStorage, error)
}

type GramboJettonWalletWithAccount interface {
	GetWalletStatus(ctx context.Context) (bool, error)
	GetWalletData(ctx context.Context) (GetWalletDataResult, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, GramboWalletStorage, error)
}
