// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type CocoonWorkerData struct {
	OwnerAddress   tlb.InternalAddress // address
	ProxyAddress   tlb.InternalAddress // address
	ProxyPublicKey tlb.Uint256         // uint256
	State          tlb.Uint2           // uint2
	Tokens         tlb.Uint64          // uint64
}

const ( // errors
)

func DecodeGetCocoonWorkerData(stack *tlb.VmStack) (result CocoonWorkerData, err error) {
	if stack.Len() != 5 {
		err = fmt.Errorf("invalid stack size %d, expected 5", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCocoonWorkerData = 0x19FBB

func GetCocoonWorkerData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CocoonWorkerData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCocoonWorkerData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCocoonWorkerData(&stack)
}

type cocoonWorkerImpl struct {
	executor Executor
}

func NewCocoonWorker(executor Executor) CocoonWorker {
	return &cocoonWorkerImpl{executor: executor}
}

func (c cocoonWorkerImpl) WithAccountId(accountID ton.AccountID) CocoonWorkerWithAccount {
	return &cocoonWorkerWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c cocoonWorkerImpl) GetCocoonWorkerData(ctx context.Context, reqAccountID ton.AccountID) (CocoonWorkerData, error) {
	return GetCocoonWorkerData(ctx, c.executor, reqAccountID)
}

type cocoonWorkerWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c cocoonWorkerWithAccountImpl) GetCocoonWorkerData(ctx context.Context) (CocoonWorkerData, error) {
	return GetCocoonWorkerData(ctx, c.executor, c.accountID)
}
