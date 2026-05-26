// Code generated - DO NOT EDIT.

package abiEvaa

import (
	"context"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type Master interface {
	WithAccountId(accountID ton.AccountID) MasterWithAccount
	GetAssetSbRate(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (SbRate, error)
	GetAssetRates(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (SbRate, error)
	GetAssetReserves(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetTotals(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTotals, error)
	GetUpdatedRates(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, assetId tlb.Int257, timeElapsed tlb.Int257) (SbRate, error)
	GetUpdatedRatesForAllAssets(ctx context.Context, reqAccountID ton.AccountID, timeElapsed tlb.Int257) (boc.Cell, error)
	GetCollateralQuote(ctx context.Context, reqAccountID ton.AccountID, borrowAssetId tlb.Int257, borrowLiquidateAmount tlb.Int257, collateralAssetId tlb.Int257, pricesPacked boc.Cell) (tlb.Int257, error)
	GetUserAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error)
	GetUserSubaccountAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress, subaccountId tlb.Int32) (tlb.InternalAddress, error)
	GetActive(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetTokensKeys(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error)
	GetLastUserScVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetUpgradeConfig(ctx context.Context, reqAccountID ton.AccountID) (UpgradeConfigResult, error)
	GetAssetTrackingInfo(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTrackingInfo, error)
	GetSupervisor(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error)
	GetAssetTotalPrincipals(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTotals, error)
	GetAssetBalance(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetLiquidityById(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetLiquidityMinusReservesById(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error)
	GetStore(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error)
	GetClaimAssetReservesMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257) (tlb.Int257, error)
	GetSupplyMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, supplyUserMessage boc.Cell) (tlb.Int257, error)
	GetWithdrawMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, withdrawUserMessage boc.Cell) (tlb.Int257, error)
	GetLiquidateMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, liquidateUserMessage boc.Cell) (tlb.Int257, error)
}

type MasterWithAccount interface {
	GetAssetSbRate(ctx context.Context, assetId tlb.Int257) (SbRate, error)
	GetAssetRates(ctx context.Context, assetId tlb.Int257) (SbRate, error)
	GetAssetReserves(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetTotals(ctx context.Context, assetId tlb.Int257) (AssetTotals, error)
	GetUpdatedRates(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, assetId tlb.Int257, timeElapsed tlb.Int257) (SbRate, error)
	GetUpdatedRatesForAllAssets(ctx context.Context, timeElapsed tlb.Int257) (boc.Cell, error)
	GetCollateralQuote(ctx context.Context, borrowAssetId tlb.Int257, borrowLiquidateAmount tlb.Int257, collateralAssetId tlb.Int257, pricesPacked boc.Cell) (tlb.Int257, error)
	GetUserAddress(ctx context.Context, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error)
	GetUserSubaccountAddress(ctx context.Context, ownerAddress tlb.InternalAddress, subaccountId tlb.Int32) (tlb.InternalAddress, error)
	GetActive(ctx context.Context) (tlb.Int257, error)
	GetTokensKeys(ctx context.Context) (boc.Cell, error)
	GetLastUserScVersion(ctx context.Context) (tlb.Int257, error)
	GetUpgradeConfig(ctx context.Context) (UpgradeConfigResult, error)
	GetAssetTrackingInfo(ctx context.Context, assetId tlb.Int257) (AssetTrackingInfo, error)
	GetSupervisor(ctx context.Context) (tlb.InternalAddress, error)
	GetAssetTotalPrincipals(ctx context.Context, assetId tlb.Int257) (AssetTotals, error)
	GetAssetBalance(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetLiquidityById(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error)
	GetAssetLiquidityMinusReservesById(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error)
	GetStore(ctx context.Context) (boc.Cell, error)
	GetClaimAssetReservesMinAttachment(ctx context.Context, fwdFee tlb.Int257) (tlb.Int257, error)
	GetSupplyMinAttachment(ctx context.Context, fwdFee tlb.Int257, supplyUserMessage boc.Cell) (tlb.Int257, error)
	GetWithdrawMinAttachment(ctx context.Context, fwdFee tlb.Int257, withdrawUserMessage boc.Cell) (tlb.Int257, error)
	GetLiquidateMinAttachment(ctx context.Context, fwdFee tlb.Int257, liquidateUserMessage boc.Cell) (tlb.Int257, error)
}

type User interface {
	WithAccountId(accountID ton.AccountID) UserWithAccount
	GetCodeVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetIsUserSc(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetAccountAssetBalance(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257, sRate tlb.Int257, bRate tlb.Int257) (tlb.Int257, error)
	GetAccountBalances(ctx context.Context, reqAccountID ton.AccountID, assetDynamicsCollection boc.Cell) (boc.Cell, error)
	GetAccountHealth(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetAvailableToBorrow(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetIsLiquidable(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetAggregatedBalances(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (AggregatedBalances, error)
	GetAssetPrincipal(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error)
	GetPrincipals(ctx context.Context, reqAccountID ton.AccountID) (tlb.Hashmap[tlb.Bits256, tlb.Int64], error)
	GetRewards(ctx context.Context, reqAccountID ton.AccountID) (tlb.Maybe[boc.Cell], error)
	GetAllUserScData(ctx context.Context, reqAccountID ton.AccountID) (UserScData, error)
	GetMaximumWithdrawAmount(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257, pricesPacked boc.Cell, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell) (tlb.Int257, error)
}

type UserWithAccount interface {
	GetCodeVersion(ctx context.Context) (tlb.Int257, error)
	GetIsUserSc(ctx context.Context) (tlb.Int257, error)
	GetAccountAssetBalance(ctx context.Context, assetId tlb.Int257, sRate tlb.Int257, bRate tlb.Int257) (tlb.Int257, error)
	GetAccountBalances(ctx context.Context, assetDynamicsCollection boc.Cell) (boc.Cell, error)
	GetAccountHealth(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetAvailableToBorrow(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetIsLiquidable(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error)
	GetAggregatedBalances(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (AggregatedBalances, error)
	GetAssetPrincipal(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error)
	GetPrincipals(ctx context.Context) (tlb.Hashmap[tlb.Bits256, tlb.Int64], error)
	GetRewards(ctx context.Context) (tlb.Maybe[boc.Cell], error)
	GetAllUserScData(ctx context.Context) (UserScData, error)
	GetMaximumWithdrawAmount(ctx context.Context, assetId tlb.Int257, pricesPacked boc.Cell, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell) (tlb.Int257, error)
}
