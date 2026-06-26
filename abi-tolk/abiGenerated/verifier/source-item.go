// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type SourceItemStorage struct {
	VerifierId           tlb.Uint256              // uint256
	VerifiedCodeCellHash tlb.Uint256              // uint256
	SourceItemRegistry   tlb.InternalAddress      // address
	Content              tlb.RefT[*SourceContent] // Cell<SourceContent>
}

type SourceItemData struct {
	VerifierId           tlb.Uint256              // uint256
	VerifiedCodeCellHash tlb.Uint256              // uint256
	SourceItemRegistry   tlb.InternalAddress      // address
	Content              tlb.RefT[*SourceContent] // Cell<SourceContent>
}

func DecodeGetSourceItemData(stack *tlb.VmStack) (result SourceItemData, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSourceItemData = 0x177B6

func GetSourceItemData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result SourceItemData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSourceItemData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSourceItemData(&stack)
}

func SourceItem_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage SourceItemStorage, err error) {
	sa, err = executor.GetAccountState(ctx, accountID)
	if err != nil {
		return
	}
	acc := sa.Account
	if acc.SumType != "Account" {
		err = fmt.Errorf("account does not exist")
	} else if state := acc.Account.Storage.State; state.SumType != "AccountActive" {
		err = fmt.Errorf("account is not active")
	} else if data := state.AccountActive.StateInit.Data; !data.Exists {
		err = fmt.Errorf("account has no storage data")
	} else {
		err = storage.UnmarshalTLB(&data.Value.Value, tlb.NewDecoder())
	}
	return
}

type sourceItemImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewSourceItem(executor Executor, storageExecutor StorageExecutor) SourceItem {
	return &sourceItemImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c sourceItemImpl) WithAccountId(accountID ton.AccountID) SourceItemWithAccount {
	return &sourceItemWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c sourceItemImpl) GetSourceItemData(ctx context.Context, reqAccountID ton.AccountID) (SourceItemData, error) {
	return GetSourceItemData(ctx, c.executor, reqAccountID)
}

func (c sourceItemImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, SourceItemStorage, error) {
	return SourceItem_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type sourceItemWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c sourceItemWithAccountImpl) GetSourceItemData(ctx context.Context) (SourceItemData, error) {
	return GetSourceItemData(ctx, c.executor, c.accountID)
}

func (c sourceItemWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, SourceItemStorage, error) {
	return SourceItem_AccountState(ctx, c.storageExecutor, c.accountID)
}
