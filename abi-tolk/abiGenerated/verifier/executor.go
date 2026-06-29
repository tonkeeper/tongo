// Code generated - DO NOT EDIT.

package abiVerifier

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

type SourceItem interface {
	WithAccountId(accountID ton.AccountID) SourceItemWithAccount
	GetSourceItemData(ctx context.Context, reqAccountID ton.AccountID) (SourceItemData, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, SourceItemStorage, error)
}

type SourceItemWithAccount interface {
	GetSourceItemData(ctx context.Context) (SourceItemData, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, SourceItemStorage, error)
}

type SourcesRegistry interface {
	WithAccountId(accountID ton.AccountID) SourcesRegistryWithAccount
	GetSourceItemAddress(ctx context.Context, reqAccountID ton.AccountID, verifierId tlb.Int257, verifiedCodeCellHash tlb.Int257) (tlb.InternalAddress, error)
	GetVerifierRegistryAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error)
	GetAdminAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error)
	GetDeploymentCosts(ctx context.Context, reqAccountID ton.AccountID) (DeploymentCosts, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, SourcesRegistryStorage, error)
}

type SourcesRegistryWithAccount interface {
	GetSourceItemAddress(ctx context.Context, verifierId tlb.Int257, verifiedCodeCellHash tlb.Int257) (tlb.InternalAddress, error)
	GetVerifierRegistryAddress(ctx context.Context) (tlb.InternalAddress, error)
	GetAdminAddress(ctx context.Context) (tlb.InternalAddress, error)
	GetDeploymentCosts(ctx context.Context) (DeploymentCosts, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, SourcesRegistryStorage, error)
}

type VerifierRegistry interface {
	WithAccountId(accountID ton.AccountID) VerifierRegistryWithAccount
	GetVerifier(ctx context.Context, reqAccountID ton.AccountID, id tlb.Int257) (VerifierInfo, error)
	GetVerifiersNum(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetVerifiers(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*VerifierRegistryStorage], error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, VerifierRegistryStorage, error)
}

type VerifierRegistryWithAccount interface {
	GetVerifier(ctx context.Context, id tlb.Int257) (VerifierInfo, error)
	GetVerifiersNum(ctx context.Context) (tlb.Int257, error)
	GetVerifiers(ctx context.Context) (tlb.RefT[*VerifierRegistryStorage], error)
	AccountState(ctx context.Context) (tlb.ShardAccount, VerifierRegistryStorage, error)
}
