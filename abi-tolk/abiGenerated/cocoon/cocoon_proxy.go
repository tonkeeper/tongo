// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixTextCmd uint64 = 0x00000000

type TextCmd struct {
	Action tlb.Uint8 // uint8
}

const PrefixExtProxyCloseRequestSigned uint64 = 0x636a4391

type ExtProxyCloseRequestSigned struct {
	Rest boc.Cell // RemainingBitsAndRefs
}

const PrefixExtProxyCloseCompleteRequestSigned uint64 = 0xe511abc7

type ExtProxyCloseCompleteRequestSigned struct {
	Rest boc.Cell // RemainingBitsAndRefs
}

const PrefixExtProxyPayoutRequest uint64 = 0x7610e6eb

type ExtProxyPayoutRequest struct {
	QueryId        tlb.Uint64          // uint64
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixExtProxyIncreaseStake uint64 = 0x9713f187

type ExtProxyIncreaseStake struct {
	QueryId        tlb.Uint64          // uint64
	Grams          tlb.Coins           // coins
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerProxyClose uint64 = 0xb51d5a01

type OwnerProxyClose struct {
	QueryId        tlb.Uint64          // uint64
	SendExcessesTo tlb.InternalAddress // address
}
type ProxyStorage struct {
	OwnerAddress   tlb.InternalAddress     // address
	ProxyPublicKey tlb.Uint256             // uint256
	RootAddress    tlb.InternalAddress     // address
	State          tlb.Uint2               // uint2
	Balance        tlb.Coins               // coins
	Stake          tlb.Coins               // coins
	UnlockTs       tlb.Uint32              // uint32
	Params         tlb.RefT[*CocoonParams] // Cell<CocoonParams>
}
type CocoonProxyData struct {
	OwnerAddress      tlb.InternalAddress // address
	ProxyPublicKey    tlb.Uint256         // uint256
	RootAddress       tlb.InternalAddress // address
	State             tlb.Uint2           // uint2
	Balance           tlb.Coins           // coins
	Stake             tlb.Coins           // coins
	UnlockTs          tlb.Uint32          // uint32
	PricePerToken     tlb.Coins           // coins
	WorkerFeePerToken tlb.Coins           // coins
	MinProxyStake     tlb.Coins           // coins
	MinClientStake    tlb.Coins           // coins
	ParamsVersion     tlb.Uint32          // uint32
}

const ( // errors
)

func DecodeGetCocoonProxyData(stack *tlb.VmStack) (result CocoonProxyData, err error) {
	if stack.Len() != 12 {
		err = fmt.Errorf("invalid stack size %d, expected 12", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCocoonProxyData = 0x17D97

func GetCocoonProxyData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CocoonProxyData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCocoonProxyData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCocoonProxyData(&stack)
}

func CocoonProxy_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage ProxyStorage, err error) {
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

type cocoonProxyImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewCocoonProxy(executor Executor, storageExecutor StorageExecutor) CocoonProxy {
	return &cocoonProxyImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c cocoonProxyImpl) WithAccountId(accountID ton.AccountID) CocoonProxyWithAccount {
	return &cocoonProxyWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c cocoonProxyImpl) GetCocoonProxyData(ctx context.Context, reqAccountID ton.AccountID) (CocoonProxyData, error) {
	return GetCocoonProxyData(ctx, c.executor, reqAccountID)
}

func (c cocoonProxyImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ProxyStorage, error) {
	return CocoonProxy_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type cocoonProxyWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c cocoonProxyWithAccountImpl) GetCocoonProxyData(ctx context.Context) (CocoonProxyData, error) {
	return GetCocoonProxyData(ctx, c.executor, c.accountID)
}

func (c cocoonProxyWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, ProxyStorage, error) {
	return CocoonProxy_AccountState(ctx, c.storageExecutor, c.accountID)
}
