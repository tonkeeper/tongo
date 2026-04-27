// Code generated - DO NOT EDIT.

package abiCocoon

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

type CocoonClient interface {
	WithAccountId(accountID ton.AccountID) CocoonClientWithAccount
	GetCocoonClientData(ctx context.Context, reqAccountID ton.AccountID) (CocoonClientData, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ClientStorage, error)
}

type CocoonClientWithAccount interface {
	GetCocoonClientData(ctx context.Context) (CocoonClientData, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, ClientStorage, error)
}

type CocoonProxy interface {
	WithAccountId(accountID ton.AccountID) CocoonProxyWithAccount
	GetCocoonProxyData(ctx context.Context, reqAccountID ton.AccountID) (CocoonProxyData, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ProxyStorage, error)
}

type CocoonProxyWithAccount interface {
	GetCocoonProxyData(ctx context.Context) (CocoonProxyData, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, ProxyStorage, error)
}

type CocoonRoot interface {
	WithAccountId(accountID ton.AccountID) CocoonRootWithAccount
	GetLastProxySeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetCocoonData(ctx context.Context, reqAccountID ton.AccountID) (CocoonData, error)
	GetCurParams(ctx context.Context, reqAccountID ton.AccountID) (CurrentCocoonParams, error)
	GetProxyHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error)
	GetWorkerHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error)
	GetModelHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, RootStorage, error)
}

type CocoonRootWithAccount interface {
	GetLastProxySeqno(ctx context.Context) (tlb.Int257, error)
	GetCocoonData(ctx context.Context) (CocoonData, error)
	GetCurParams(ctx context.Context) (CurrentCocoonParams, error)
	GetProxyHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error)
	GetWorkerHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error)
	GetModelHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, RootStorage, error)
}

type CocoonWallet interface {
	WithAccountId(accountID ton.AccountID) CocoonWalletWithAccount
	GetSeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetPublicKey(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetOwnerAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, WalletStorage, error)
}

type CocoonWalletWithAccount interface {
	GetSeqno(ctx context.Context) (tlb.Int257, error)
	GetPublicKey(ctx context.Context) (tlb.Int257, error)
	GetOwnerAddress(ctx context.Context) (tlb.InternalAddress, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, WalletStorage, error)
}

type CocoonWorker interface {
	WithAccountId(accountID ton.AccountID) CocoonWorkerWithAccount
	GetCocoonWorkerData(ctx context.Context, reqAccountID ton.AccountID) (CocoonWorkerData, error)
}

type CocoonWorkerWithAccount interface {
	GetCocoonWorkerData(ctx context.Context) (CocoonWorkerData, error)
}
