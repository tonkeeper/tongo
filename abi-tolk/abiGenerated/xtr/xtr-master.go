// Code generated - DO NOT EDIT.

package abiXtr

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixPushXTR uint64 = 0x6f027868

type PushXTR struct {
	Seqno  tlb.Uint64 // uint64
	Amount tlb.Coins  // coins
}

func DecodeGetUserLatestVersion(stack *tlb.VmStack) (result tlb.Uint32, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUserLatestVersion = 0x18058

func GetUserLatestVersion(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Uint32, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUserLatestVersion, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUserLatestVersion(&stack)
}

func DecodeGetPaymentLatestVersion(stack *tlb.VmStack) (result tlb.Uint32, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPaymentLatestVersion = 0x1ECAF

func GetPaymentLatestVersion(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Uint32, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPaymentLatestVersion, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPaymentLatestVersion(&stack)
}

type xtrMasterImpl struct {
	executor Executor
}

func NewXtrMaster(executor Executor) XtrMaster {
	return &xtrMasterImpl{executor: executor}
}

func (c xtrMasterImpl) WithAccountId(accountID ton.AccountID) XtrMasterWithAccount {
	return &xtrMasterWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c xtrMasterImpl) GetUserLatestVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error) {
	return GetUserLatestVersion(ctx, c.executor, reqAccountID)
}

func (c xtrMasterImpl) GetPaymentLatestVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error) {
	return GetPaymentLatestVersion(ctx, c.executor, reqAccountID)
}

type xtrMasterWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c xtrMasterWithAccountImpl) GetUserLatestVersion(ctx context.Context) (tlb.Uint32, error) {
	return GetUserLatestVersion(ctx, c.executor, c.accountID)
}

func (c xtrMasterWithAccountImpl) GetPaymentLatestVersion(ctx context.Context) (tlb.Uint32, error) {
	return GetPaymentLatestVersion(ctx, c.executor, c.accountID)
}
