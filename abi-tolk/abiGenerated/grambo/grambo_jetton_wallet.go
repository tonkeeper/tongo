// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type WalletIncomingMessage GramboActivateWallet

type GramboWalletStorage struct {
	Activated bool                // bool
	Balance   tlb.Coins           // coins
	Owner     tlb.InternalAddress // address
	Master    tlb.InternalAddress // address
}

type GetWalletDataResult struct {
	Balance             tlb.Int257          // int
	OwnerAddress        tlb.InternalAddress // address
	JettonMasterAddress tlb.InternalAddress // address
	JettonWalletCode    boc.Cell            // cell
}

func DecodeGetWalletStatus(stack *tlb.VmStack) (result bool, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadBool()
}

const MethodIDGetWalletStatus = 0x1B700

func GetWalletStatus(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result bool, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetWalletStatus, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetWalletStatus(&stack)
}

func DecodeGetWalletData(stack *tlb.VmStack) (result GetWalletDataResult, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetWalletData = 0x17B02

func GetWalletData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetWalletDataResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetWalletData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetWalletData(&stack)
}

func GramboJettonWallet_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage GramboWalletStorage, err error) {
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

type gramboJettonWalletImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewGramboJettonWallet(executor Executor, storageExecutor StorageExecutor) GramboJettonWallet {
	return &gramboJettonWalletImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c gramboJettonWalletImpl) WithAccountId(accountID ton.AccountID) GramboJettonWalletWithAccount {
	return &gramboJettonWalletWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c gramboJettonWalletImpl) GetWalletStatus(ctx context.Context, reqAccountID ton.AccountID) (bool, error) {
	return GetWalletStatus(ctx, c.executor, reqAccountID)
}

func (c gramboJettonWalletImpl) GetWalletData(ctx context.Context, reqAccountID ton.AccountID) (GetWalletDataResult, error) {
	return GetWalletData(ctx, c.executor, reqAccountID)
}

func (c gramboJettonWalletImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, GramboWalletStorage, error) {
	return GramboJettonWallet_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type gramboJettonWalletWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c gramboJettonWalletWithAccountImpl) GetWalletStatus(ctx context.Context) (bool, error) {
	return GetWalletStatus(ctx, c.executor, c.accountID)
}

func (c gramboJettonWalletWithAccountImpl) GetWalletData(ctx context.Context) (GetWalletDataResult, error) {
	return GetWalletData(ctx, c.executor, c.accountID)
}

func (c gramboJettonWalletWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, GramboWalletStorage, error) {
	return GramboJettonWallet_AccountState(ctx, c.storageExecutor, c.accountID)
}
