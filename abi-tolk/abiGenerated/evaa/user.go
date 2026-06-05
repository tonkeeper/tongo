// Code generated - DO NOT EDIT.

package abiEvaa

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type UserScData struct {
	CodeVersion    tlb.Int257          // int
	MasterAddress  tlb.InternalAddress // address
	OwnerAddress   tlb.InternalAddress // address
	UserPrincipals boc.Cell            // cell
	State          tlb.Int64           // int64
	UserRewards    tlb.Maybe[boc.Cell] // cell?
	BackupCell1    tlb.Maybe[boc.Cell] // cell?
	BackupCell2    tlb.Maybe[boc.Cell] // cell?
}
type AggregatedBalances struct {
	Supply tlb.Int257 // int
	Borrow tlb.Int257 // int
}

const ( // errors
)

func DecodeGetCodeVersion(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCodeVersion = 0x16EBE

func GetCodeVersion(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCodeVersion, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCodeVersion(&stack)
}

func DecodeGetIsUserSc(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetIsUserSc = 0x11C45

func GetIsUserSc(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetIsUserSc, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetIsUserSc(&stack)
}

func DecodeGetAccountAssetBalance(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAccountAssetBalance = 0x14EFC

func GetAccountAssetBalance(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257, sRate tlb.Int257, bRate tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(sRate)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(bRate)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAccountAssetBalance, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAccountAssetBalance(&stack)
}

func DecodeGetAccountBalances(stack *tlb.VmStack) (result boc.Cell, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadCell()
}

const MethodIDGetAccountBalances = 0x1468D

func GetAccountBalances(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetDynamicsCollection boc.Cell) (result boc.Cell, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(assetDynamicsCollection)
		if err != nil {
			err = fmt.Errorf("encode param assetDynamicsCollection: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAccountBalances, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAccountBalances(&stack)
}

func DecodeGetAccountHealth(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAccountHealth = 0x1B072

func GetAccountHealth(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (result tlb.Int257, err error) {
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
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAccountHealth, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAccountHealth(&stack)
}

func DecodeGetAvailableToBorrow(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAvailableToBorrow = 0x16E2C

func GetAvailableToBorrow(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (result tlb.Int257, err error) {
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
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAvailableToBorrow, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAvailableToBorrow(&stack)
}

func DecodeGetIsLiquidable(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetIsLiquidable = 0x117E7

func GetIsLiquidable(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (result tlb.Int257, err error) {
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
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetIsLiquidable, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetIsLiquidable(&stack)
}

func DecodeGetAggregatedBalances(stack *tlb.VmStack) (result AggregatedBalances, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAggregatedBalances = 0x14994

func GetAggregatedBalances(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (result AggregatedBalances, err error) {
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
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAggregatedBalances, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAggregatedBalances(&stack)
}

func DecodeGetAssetPrincipal(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAssetPrincipal = 0x134FB

func GetAssetPrincipal(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAssetPrincipal, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAssetPrincipal(&stack)
}

func DecodeGetPrincipals(stack *tlb.VmStack) (result tlb.Hashmap[tlb.Bits256, tlb.Int64], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return tlb.ReadHashmapFromStack[tlb.Bits256, tlb.Int64](stack)
}

const MethodIDGetPrincipals = 0x1FAF2

func GetPrincipals(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Hashmap[tlb.Bits256, tlb.Int64], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPrincipals, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPrincipals(&stack)
}

func DecodeGetRewards(stack *tlb.VmStack) (result tlb.Maybe[boc.Cell], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return tlb.StackReadMaybeCallback(stack, func(stack *tlb.VmStack) (value boc.Cell, err error) {
		return stack.ReadCell()

	})
}

const MethodIDGetRewards = 0x19B4E

func GetRewards(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Maybe[boc.Cell], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetRewards, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetRewards(&stack)
}

func DecodeGetAllUserScData(stack *tlb.VmStack) (result UserScData, err error) {
	if stack.Len() != 8 {
		err = fmt.Errorf("invalid stack size %d, expected 8", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAllUserScData = 0x1702A

func GetAllUserScData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result UserScData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAllUserScData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAllUserScData(&stack)
}

func DecodeGetMaximumWithdrawAmount(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetMaximumWithdrawAmount = 0x1D090

func GetMaximumWithdrawAmount(ctx context.Context, executor Executor, reqAccountID ton.AccountID, assetId tlb.Int257, pricesPacked boc.Cell, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(assetId)})
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(pricesPacked)
		if err != nil {
			err = fmt.Errorf("encode param pricesPacked: %w", err)
			return
		}
		stack.Put(val)
	}
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
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetMaximumWithdrawAmount, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetMaximumWithdrawAmount(&stack)
}

type userImpl struct {
	executor Executor
}

func NewUser(executor Executor) User {
	return &userImpl{executor: executor}
}

func (c userImpl) WithAccountId(accountID ton.AccountID) UserWithAccount {
	return &userWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c userImpl) GetCodeVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetCodeVersion(ctx, c.executor, reqAccountID)
}

func (c userImpl) GetIsUserSc(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetIsUserSc(ctx, c.executor, reqAccountID)
}

func (c userImpl) GetAccountAssetBalance(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257, sRate tlb.Int257, bRate tlb.Int257) (tlb.Int257, error) {
	return GetAccountAssetBalance(ctx, c.executor, reqAccountID, assetId, sRate, bRate)
}

func (c userImpl) GetAccountBalances(ctx context.Context, reqAccountID ton.AccountID, assetDynamicsCollection boc.Cell) (boc.Cell, error) {
	return GetAccountBalances(ctx, c.executor, reqAccountID, assetDynamicsCollection)
}

func (c userImpl) GetAccountHealth(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetAccountHealth(ctx, c.executor, reqAccountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userImpl) GetAvailableToBorrow(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetAvailableToBorrow(ctx, c.executor, reqAccountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userImpl) GetIsLiquidable(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetIsLiquidable(ctx, c.executor, reqAccountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userImpl) GetAggregatedBalances(ctx context.Context, reqAccountID ton.AccountID, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (AggregatedBalances, error) {
	return GetAggregatedBalances(ctx, c.executor, reqAccountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userImpl) GetAssetPrincipal(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetPrincipal(ctx, c.executor, reqAccountID, assetId)
}

func (c userImpl) GetPrincipals(ctx context.Context, reqAccountID ton.AccountID) (tlb.Hashmap[tlb.Bits256, tlb.Int64], error) {
	return GetPrincipals(ctx, c.executor, reqAccountID)
}

func (c userImpl) GetRewards(ctx context.Context, reqAccountID ton.AccountID) (tlb.Maybe[boc.Cell], error) {
	return GetRewards(ctx, c.executor, reqAccountID)
}

func (c userImpl) GetAllUserScData(ctx context.Context, reqAccountID ton.AccountID) (UserScData, error) {
	return GetAllUserScData(ctx, c.executor, reqAccountID)
}

func (c userImpl) GetMaximumWithdrawAmount(ctx context.Context, reqAccountID ton.AccountID, assetId tlb.Int257, pricesPacked boc.Cell, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell) (tlb.Int257, error) {
	return GetMaximumWithdrawAmount(ctx, c.executor, reqAccountID, assetId, pricesPacked, assetConfigCollection, assetDynamicsCollection)
}

type userWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c userWithAccountImpl) GetCodeVersion(ctx context.Context) (tlb.Int257, error) {
	return GetCodeVersion(ctx, c.executor, c.accountID)
}

func (c userWithAccountImpl) GetIsUserSc(ctx context.Context) (tlb.Int257, error) {
	return GetIsUserSc(ctx, c.executor, c.accountID)
}

func (c userWithAccountImpl) GetAccountAssetBalance(ctx context.Context, assetId tlb.Int257, sRate tlb.Int257, bRate tlb.Int257) (tlb.Int257, error) {
	return GetAccountAssetBalance(ctx, c.executor, c.accountID, assetId, sRate, bRate)
}

func (c userWithAccountImpl) GetAccountBalances(ctx context.Context, assetDynamicsCollection boc.Cell) (boc.Cell, error) {
	return GetAccountBalances(ctx, c.executor, c.accountID, assetDynamicsCollection)
}

func (c userWithAccountImpl) GetAccountHealth(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetAccountHealth(ctx, c.executor, c.accountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userWithAccountImpl) GetAvailableToBorrow(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetAvailableToBorrow(ctx, c.executor, c.accountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userWithAccountImpl) GetIsLiquidable(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (tlb.Int257, error) {
	return GetIsLiquidable(ctx, c.executor, c.accountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userWithAccountImpl) GetAggregatedBalances(ctx context.Context, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell, pricesPacked boc.Cell) (AggregatedBalances, error) {
	return GetAggregatedBalances(ctx, c.executor, c.accountID, assetConfigCollection, assetDynamicsCollection, pricesPacked)
}

func (c userWithAccountImpl) GetAssetPrincipal(ctx context.Context, assetId tlb.Int257) (tlb.Int257, error) {
	return GetAssetPrincipal(ctx, c.executor, c.accountID, assetId)
}

func (c userWithAccountImpl) GetPrincipals(ctx context.Context) (tlb.Hashmap[tlb.Bits256, tlb.Int64], error) {
	return GetPrincipals(ctx, c.executor, c.accountID)
}

func (c userWithAccountImpl) GetRewards(ctx context.Context) (tlb.Maybe[boc.Cell], error) {
	return GetRewards(ctx, c.executor, c.accountID)
}

func (c userWithAccountImpl) GetAllUserScData(ctx context.Context) (UserScData, error) {
	return GetAllUserScData(ctx, c.executor, c.accountID)
}

func (c userWithAccountImpl) GetMaximumWithdrawAmount(ctx context.Context, assetId tlb.Int257, pricesPacked boc.Cell, assetConfigCollection boc.Cell, assetDynamicsCollection boc.Cell) (tlb.Int257, error) {
	return GetMaximumWithdrawAmount(ctx, c.executor, c.accountID, assetId, pricesPacked, assetConfigCollection, assetDynamicsCollection)
}
