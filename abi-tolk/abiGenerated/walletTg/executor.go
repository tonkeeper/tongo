// Code generated - DO NOT EDIT.

package abiWalletTg

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

type WalletTg interface {
	WithAccountId(accountID ton.AccountID) WalletTgWithAccount
	GetRevision(ctx context.Context, reqAccountID ton.AccountID) (Revision, error)
	GetSeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error)
	GetSubwalletId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error)
	GetPublicKey(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint256, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, Storage, error)
}

type WalletTgWithAccount interface {
	GetRevision(ctx context.Context) (Revision, error)
	GetSeqno(ctx context.Context) (tlb.Uint32, error)
	GetSubwalletId(ctx context.Context) (tlb.Uint32, error)
	GetPublicKey(ctx context.Context) (tlb.Uint256, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, Storage, error)
}
