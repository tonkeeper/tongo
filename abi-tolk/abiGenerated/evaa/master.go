// Code generated - DO NOT EDIT.

package abiEvaa

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type UpgradeConfigResult struct {
	MasterCodeVersion tlb.Int257          // int
	UserCodeVersion   tlb.Int257          // int
	Timeout           tlb.Uint32          // uint32
	UpdateTime        tlb.Uint64          // uint64
	FreezeTime        tlb.Uint64          // uint64
	UserCode          boc.Cell            // cell
	NewMasterCode     tlb.Maybe[boc.Cell] // cell?
	NewUserCode       tlb.Maybe[boc.Cell] // cell?
}
type AssetTrackingInfo struct {
	TrackingSupplyIndex tlb.Int257 // int
	TrackingBorrowIndex tlb.Int257 // int
	LastAccrual         tlb.Int257 // int
}
type SbRate struct {
	SRate tlb.Int257 // int
	BRate tlb.Int257 // int
}
type AssetTotals struct {
	TotalSupply tlb.Int257 // int
	TotalBorrow tlb.Int257 // int
}

const ( // errors
)

func DecodeGetAssetSbRate(stack *tlb.VmStack) (result SbRate, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetSbRate = 0x11FA3

func GetAssetSbRate(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result SbRate, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetSbRate, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetSbRate(&stack)
}

func DecodeGetAssetRates(stack *tlb.VmStack) (result SbRate, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetRates = 0x14029

func GetAssetRates(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result SbRate, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetRates, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetRates(&stack)
}

func DecodeGetAssetReserves(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetReserves = 0x10E63

func GetAssetReserves(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetReserves, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetReserves(&stack)
}

func DecodeGetAssetTotals(stack *tlb.VmStack) (result AssetTotals, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetTotals = 0x1A2F4

func GetAssetTotals(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result AssetTotals, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetTotals, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetTotals(&stack)
}

func DecodeGetUpdatedRates(stack *tlb.VmStack) (result SbRate, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUpdatedRates = 0x1DDC9

func GetUpdatedRates(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, assetId tlb.Int257, timeElapsed tlb.Int257) (result SbRate, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(assetConfigCollection)
		if err != nil {
			err = fmt.Errorf("encode param assetConfigCollection: %w", err)
			return
		}
		stack.Put(val)
	}
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(assetDynamicsCollection)
		if err != nil {
			err = fmt.Errorf("encode param assetDynamicsCollection: %w", err)
			return
		}
		stack.Put(val)
	}
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(timeElapsed)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUpdatedRates, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUpdatedRates(&stack)
}

func DecodeGetUpdatedRatesForAllAssets(stack *tlb.VmStack) (result boc.Cell, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadCell()
}

const MethodIDGetUpdatedRatesForAllAssets = 0x166D2

func GetUpdatedRatesForAllAssets(ctx context.Context, executor Executor, reqAccountID ton.AccountID, timeElapsed tlb.Int257) (result boc.Cell, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(timeElapsed)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUpdatedRatesForAllAssets, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUpdatedRatesForAllAssets(&stack)
}

func DecodeGetCollateralQuote(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCollateralQuote = 0x11DF9

func GetCollateralQuote(ctx context.Context, executor Executor, reqAccountID ton.AccountID, borrowAssetId tlb.Int257, borrowLiquidateAmount tlb.Int257, collateralAssetId tlb.Int257, pricesPacked boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(borrowAssetId)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(borrowLiquidateAmount)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(collateralAssetId)})
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCollateralQuote, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCollateralQuote(&stack)
}

func DecodeGetUserAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUserAddress = 0x19180

func GetUserAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(ownerAddress)
		if err != nil {
			err = fmt.Errorf("encode param ownerAddress: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUserAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUserAddress(&stack)
}

func DecodeGetUserSubaccountAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUserSubaccountAddress = 0x1D6D3

func GetUserSubaccountAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress, subaccountId tlb.Int32) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(ownerAddress)
		if err != nil {
			err = fmt.Errorf("encode param ownerAddress: %w", err)
			return
		}
		stack.Put(val)
	}
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257FromInt64(int64(subaccountId))})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUserSubaccountAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUserSubaccountAddress(&stack)
}

func DecodeGetActive(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetActive = 0x1F92F

func GetActive(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetActive, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetActive(&stack)
}

func DecodeGetTokensKeys(stack *tlb.VmStack) (result boc.Cell, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadCell()
}

const MethodIDGetTokensKeys = 0x18084

func GetTokensKeys(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result boc.Cell, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetTokensKeys, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetTokensKeys(&stack)
}

func DecodeGetLastUserScVersion(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetLastUserScVersion = 0x15A10

func GetLastUserScVersion(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetLastUserScVersion, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetLastUserScVersion(&stack)
}

func DecodeGetUpgradeConfig(stack *tlb.VmStack) (result UpgradeConfigResult, err error) {
	if stack.Len() != 8 {
		err = fmt.Errorf("invalid stack size %d, expected 8", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUpgradeConfig = 0x11FDA

func GetUpgradeConfig(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result UpgradeConfigResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUpgradeConfig, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUpgradeConfig(&stack)
}

func DecodeGetAssetTrackingInfo(stack *tlb.VmStack) (result AssetTrackingInfo, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetTrackingInfo = 0x18F78

func GetAssetTrackingInfo(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result AssetTrackingInfo, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetTrackingInfo, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetTrackingInfo(&stack)
}

func DecodeGetSupervisor(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSupervisor = 0x15AB4

func GetSupervisor(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSupervisor, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSupervisor(&stack)
}

func DecodeGetAssetTotalPrincipals(stack *tlb.VmStack) (result AssetTotals, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetTotalPrincipals = 0x1180F

func GetAssetTotalPrincipals(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result AssetTotals, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetTotalPrincipals, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetTotalPrincipals(&stack)
}

func DecodeGetAssetBalance(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetBalance = 0x18ACB

func GetAssetBalance(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetBalance, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetBalance(&stack)
}

func DecodeGetAssetLiquidityById(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetLiquidityById = 0x10A18

func GetAssetLiquidityById(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetLiquidityById, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetLiquidityById(&stack)
}

func DecodeGetAssetLiquidityMinusReservesById(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetLiquidityMinusReservesById = 0x19714

func GetAssetLiquidityMinusReservesById(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetLiquidityMinusReservesById, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetLiquidityMinusReservesById(&stack)
}

func DecodeGetStore(stack *tlb.VmStack) (result boc.Cell, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadCell()
}

const MethodIDGetStore = 0x15526

func GetStore(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result boc.Cell, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetStore, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetStore(&stack)
}

func DecodeGetClaimAssetReservesMinAttachment(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetClaimAssetReservesMinAttachment = 0x1AF9D

func GetClaimAssetReservesMinAttachment(ctx context.Context, executor Executor, reqAccountID ton.AccountID, fwdFee tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(fwdFee)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetClaimAssetReservesMinAttachment, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetClaimAssetReservesMinAttachment(&stack)
}

func DecodeGetSupplyMinAttachment(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSupplyMinAttachment = 0x163D7

func GetSupplyMinAttachment(ctx context.Context, executor Executor, reqAccountID ton.AccountID, fwdFee tlb.Int257, supplyUserMessage boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(fwdFee)})
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(supplyUserMessage)
		if err != nil {
			err = fmt.Errorf("encode param supplyUserMessage: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSupplyMinAttachment, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSupplyMinAttachment(&stack)
}

func DecodeGetWithdrawMinAttachment(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetWithdrawMinAttachment = 0x1E459

func GetWithdrawMinAttachment(ctx context.Context, executor Executor, reqAccountID ton.AccountID, fwdFee tlb.Int257, withdrawUserMessage boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(fwdFee)})
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(withdrawUserMessage)
		if err != nil {
			err = fmt.Errorf("encode param withdrawUserMessage: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetWithdrawMinAttachment, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetWithdrawMinAttachment(&stack)
}

func DecodeGetLiquidateMinAttachment(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetLiquidateMinAttachment = 0x1B8D4

func GetLiquidateMinAttachment(ctx context.Context, executor Executor, reqAccountID ton.AccountID, fwdFee tlb.Int257, liquidateUserMessage boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(fwdFee)})
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(liquidateUserMessage)
		if err != nil {
			err = fmt.Errorf("encode param liquidateUserMessage: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetLiquidateMinAttachment, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetLiquidateMinAttachment(&stack)
}

type masterImpl struct {
	executor Executor
}

func NewMaster(executor Executor) Master {
	return &masterImpl{executor: executor}
}

func (c masterImpl) WithAccountId(accountID ton.AccountID) MasterWithAccount {
	return &masterWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c masterImpl) GetAssetSbRate(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (SbRate, error) {
	return GetAssetSbRate(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetRates(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (SbRate, error) {
	return GetAssetRates(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetReserves(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetReserves(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetTotals(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTotals, error) {
	return GetAssetTotals(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetUpdatedRates(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, assetId tlb.Int257, timeElapsed tlb.Int257) (SbRate, error) {
	return GetUpdatedRates(ctx, c.executor, reqAccountID, assetConfigCollection, assetDynamicsCollection, assetId, timeElapsed)
}

func (c masterImpl) GetUpdatedRatesForAllAssets(ctx context.Context, reqAccountID ton.AccountID, timeElapsed tlb.Int257) (boc.Cell, error) {
	return GetUpdatedRatesForAllAssets(ctx, c.executor, reqAccountID, timeElapsed)
}

func (c masterImpl) GetCollateralQuote(ctx context.Context, reqAccountID ton.AccountID, borrowAssetId tlb.Int257, borrowLiquidateAmount tlb.Int257, collateralAssetId tlb.Int257, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetCollateralQuote(ctx, c.executor, reqAccountID, borrowAssetId, borrowLiquidateAmount, collateralAssetId, pricesPacked)
}

func (c masterImpl) GetUserAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetUserAddress(ctx, c.executor, reqAccountID, ownerAddress)
}

func (c masterImpl) GetUserSubaccountAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress, subaccountId tlb.Int32) (tlb.InternalAddress, error) {
	return GetUserSubaccountAddress(ctx, c.executor, reqAccountID, ownerAddress, subaccountId)
}

func (c masterImpl) GetActive(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetActive(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetTokensKeys(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error) {
	return GetTokensKeys(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetLastUserScVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetLastUserScVersion(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetUpgradeConfig(ctx context.Context, reqAccountID ton.AccountID) (UpgradeConfigResult, error) {
	return GetUpgradeConfig(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetAssetTrackingInfo(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTrackingInfo, error) {
	return GetAssetTrackingInfo(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetSupervisor(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error) {
	return GetSupervisor(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetAssetTotalPrincipals(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (AssetTotals, error) {
	return GetAssetTotalPrincipals(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetBalance(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetBalance(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetLiquidityById(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetLiquidityById(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetAssetLiquidityMinusReservesById(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetLiquidityMinusReservesById(ctx, c.executor, reqAccountID, assetId)
}

func (c masterImpl) GetStore(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error) {
	return GetStore(ctx, c.executor, reqAccountID)
}

func (c masterImpl) GetClaimAssetReservesMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257) (tlb.Int257, error) {
	return GetClaimAssetReservesMinAttachment(ctx, c.executor, reqAccountID, fwdFee)
}

func (c masterImpl) GetSupplyMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, supplyUserMessage boc.Cell) (tlb.Int257, error) {
	return GetSupplyMinAttachment(ctx, c.executor, reqAccountID, fwdFee, supplyUserMessage)
}

func (c masterImpl) GetWithdrawMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, withdrawUserMessage boc.Cell) (tlb.Int257, error) {
	return GetWithdrawMinAttachment(ctx, c.executor, reqAccountID, fwdFee, withdrawUserMessage)
}

func (c masterImpl) GetLiquidateMinAttachment(ctx context.Context, reqAccountID ton.AccountID, fwdFee tlb.Int257, liquidateUserMessage boc.Cell) (tlb.Int257, error) {
	return GetLiquidateMinAttachment(ctx, c.executor, reqAccountID, fwdFee, liquidateUserMessage)
}

type masterWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c masterWithAccountImpl) GetAssetSbRate(ctx context.Context, assetId tlb.Int257) (SbRate, error) {
	return GetAssetSbRate(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetRates(ctx context.Context, assetId tlb.Int257) (SbRate, error) {
	return GetAssetRates(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetReserves(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetReserves(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetTotals(ctx context.Context, assetId tlb.Int257) (AssetTotals, error) {
	return GetAssetTotals(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetUpdatedRates(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, assetId tlb.Int257, timeElapsed tlb.Int257) (SbRate, error) {
	return GetUpdatedRates(ctx, c.executor, c.accountID, assetConfigCollection, assetDynamicsCollection, assetId, timeElapsed)
}

func (c masterWithAccountImpl) GetUpdatedRatesForAllAssets(ctx context.Context, timeElapsed tlb.Int257) (boc.Cell, error) {
	return GetUpdatedRatesForAllAssets(ctx, c.executor, c.accountID, timeElapsed)
}

func (c masterWithAccountImpl) GetCollateralQuote(ctx context.Context, borrowAssetId tlb.Int257, borrowLiquidateAmount tlb.Int257, collateralAssetId tlb.Int257, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetCollateralQuote(ctx, c.executor, c.accountID, borrowAssetId, borrowLiquidateAmount, collateralAssetId, pricesPacked)
}

func (c masterWithAccountImpl) GetUserAddress(ctx context.Context, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetUserAddress(ctx, c.executor, c.accountID, ownerAddress)
}

func (c masterWithAccountImpl) GetUserSubaccountAddress(ctx context.Context, ownerAddress tlb.InternalAddress, subaccountId tlb.Int32) (tlb.InternalAddress, error) {
	return GetUserSubaccountAddress(ctx, c.executor, c.accountID, ownerAddress, subaccountId)
}

func (c masterWithAccountImpl) GetActive(ctx context.Context) (tlb.Int257, error) {
	return GetActive(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetTokensKeys(ctx context.Context) (boc.Cell, error) {
	return GetTokensKeys(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetLastUserScVersion(ctx context.Context) (tlb.Int257, error) {
	return GetLastUserScVersion(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetUpgradeConfig(ctx context.Context) (UpgradeConfigResult, error) {
	return GetUpgradeConfig(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetAssetTrackingInfo(ctx context.Context, assetId tlb.Int257) (AssetTrackingInfo, error) {
	return GetAssetTrackingInfo(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetSupervisor(ctx context.Context) (tlb.InternalAddress, error) {
	return GetSupervisor(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetAssetTotalPrincipals(ctx context.Context, assetId tlb.Int257) (AssetTotals, error) {
	return GetAssetTotalPrincipals(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetBalance(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetBalance(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetLiquidityById(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetLiquidityById(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetAssetLiquidityMinusReservesById(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetLiquidityMinusReservesById(ctx, c.executor, c.accountID, assetId)
}

func (c masterWithAccountImpl) GetStore(ctx context.Context) (boc.Cell, error) {
	return GetStore(ctx, c.executor, c.accountID)
}

func (c masterWithAccountImpl) GetClaimAssetReservesMinAttachment(ctx context.Context, fwdFee tlb.Int257) (tlb.Int257, error) {
	return GetClaimAssetReservesMinAttachment(ctx, c.executor, c.accountID, fwdFee)
}

func (c masterWithAccountImpl) GetSupplyMinAttachment(ctx context.Context, fwdFee tlb.Int257, supplyUserMessage boc.Cell) (tlb.Int257, error) {
	return GetSupplyMinAttachment(ctx, c.executor, c.accountID, fwdFee, supplyUserMessage)
}

func (c masterWithAccountImpl) GetWithdrawMinAttachment(ctx context.Context, fwdFee tlb.Int257, withdrawUserMessage boc.Cell) (tlb.Int257, error) {
	return GetWithdrawMinAttachment(ctx, c.executor, c.accountID, fwdFee, withdrawUserMessage)
}

func (c masterWithAccountImpl) GetLiquidateMinAttachment(ctx context.Context, fwdFee tlb.Int257, liquidateUserMessage boc.Cell) (tlb.Int257, error) {
	return GetLiquidateMinAttachment(ctx, c.executor, c.accountID, fwdFee, liquidateUserMessage)
}
