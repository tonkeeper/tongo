// Code generated - DO NOT EDIT.

package abiPythOracle

import (
	"context"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type StorageExecutor interface {
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

type Oracle interface {
	WithAccountId(accountID ton.AccountID) OracleWithAccount
	GetUpdateFee(ctx context.Context, reqAccountID ton.AccountID, data boc.Cell) (tlb.Int257, error)
	GetSingleUpdateFee(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetGovernanceDataSourceIndex(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetGovernanceDataSource(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error)
	GetLastExecutedGovernanceSequence(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetIsValidDataSource(ctx context.Context, reqAccountID ton.AccountID, dataSource boc.Cell) (tlb.Int257, error)
	GetPriceUnsafe(ctx context.Context, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (PricePoint, error)
	GetPriceNoOlderThan(ctx context.Context, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error)
	GetEmaPriceUnsafe(ctx context.Context, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (PricePoint, error)
	GetEmaPriceNoOlderThan(ctx context.Context, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error)
	GetChainId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetCurrentGuardianSetIndex(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetGuardianSet(ctx context.Context, reqAccountID ton.AccountID, index tlb.Int257) (GuardianSetInfo, error)
	GetGovernanceChainId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetGovernanceContract(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetGovernanceActionIsConsumed(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, MainStorage, error)
}

type OracleWithAccount interface {
	GetUpdateFee(ctx context.Context, data boc.Cell) (tlb.Int257, error)
	GetSingleUpdateFee(ctx context.Context) (tlb.Int257, error)
	GetGovernanceDataSourceIndex(ctx context.Context) (tlb.Int257, error)
	GetGovernanceDataSource(ctx context.Context) (boc.Cell, error)
	GetLastExecutedGovernanceSequence(ctx context.Context) (tlb.Int257, error)
	GetIsValidDataSource(ctx context.Context, dataSource boc.Cell) (tlb.Int257, error)
	GetPriceUnsafe(ctx context.Context, priceFeedId tlb.Uint256) (PricePoint, error)
	GetPriceNoOlderThan(ctx context.Context, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error)
	GetEmaPriceUnsafe(ctx context.Context, priceFeedId tlb.Uint256) (PricePoint, error)
	GetEmaPriceNoOlderThan(ctx context.Context, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error)
	GetChainId(ctx context.Context) (tlb.Int257, error)
	GetCurrentGuardianSetIndex(ctx context.Context) (tlb.Int257, error)
	GetGuardianSet(ctx context.Context, index tlb.Int257) (GuardianSetInfo, error)
	GetGovernanceChainId(ctx context.Context) (tlb.Int257, error)
	GetGovernanceContract(ctx context.Context) (tlb.Int257, error)
	GetGovernanceActionIsConsumed(ctx context.Context, hash tlb.Int257) (tlb.Int257, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, MainStorage, error)
}
