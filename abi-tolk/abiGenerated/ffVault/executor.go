// Code generated - DO NOT EDIT.

package abiFfVault

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

type Position interface {
	WithAccountId(accountID ton.AccountID) PositionWithAccount
	GetNftData(ctx context.Context, reqAccountID ton.AccountID) (NftData, error)
	GetStakePositionInfo(ctx context.Context, reqAccountID ton.AccountID) (StakePositionInfo, error)
}

type PositionWithAccount interface {
	GetNftData(ctx context.Context) (NftData, error)
	GetStakePositionInfo(ctx context.Context) (StakePositionInfo, error)
}

type Vault interface {
	WithAccountId(accountID ton.AccountID) VaultWithAccount
	GetCollectionData(ctx context.Context, reqAccountID ton.AccountID) (CollectionData, error)
	GetNftAddressByIndex(ctx context.Context, reqAccountID ton.AccountID, itemIndex tlb.Int257) (tlb.InternalAddress, error)
	GetStakingData(ctx context.Context, reqAccountID ton.AccountID) (StakingData, error)
	GetBalance(ctx context.Context, reqAccountID ton.AccountID) (Balance, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, VaultStorage, error)
}

type VaultWithAccount interface {
	GetCollectionData(ctx context.Context) (CollectionData, error)
	GetNftAddressByIndex(ctx context.Context, itemIndex tlb.Int257) (tlb.InternalAddress, error)
	GetStakingData(ctx context.Context) (StakingData, error)
	GetBalance(ctx context.Context) (Balance, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, VaultStorage, error)
}
