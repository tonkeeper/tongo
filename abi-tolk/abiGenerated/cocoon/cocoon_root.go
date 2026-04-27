// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixAddWorkerType uint64 = 0xe34b1c60

type AddWorkerType struct {
	QueryId    tlb.Uint64  // uint64
	WorkerHash tlb.Uint256 // uint256
}

const PrefixDelWorkerType uint64 = 0x8d94a79a

type DelWorkerType struct {
	QueryId    tlb.Uint64  // uint64
	WorkerHash tlb.Uint256 // uint256
}

const PrefixAddModelType uint64 = 0xc146134d

type AddModelType struct {
	QueryId   tlb.Uint64  // uint64
	ModelHash tlb.Uint256 // uint256
}

const PrefixDelModelType uint64 = 0x92b11c18

type DelModelType struct {
	QueryId   tlb.Uint64  // uint64
	ModelHash tlb.Uint256 // uint256
}

const PrefixAddProxyType uint64 = 0x71860e80

type AddProxyType struct {
	QueryId   tlb.Uint64  // uint64
	ProxyHash tlb.Uint256 // uint256
}

const PrefixDelProxyType uint64 = 0x3c41d0b2

type DelProxyType struct {
	QueryId   tlb.Uint64  // uint64
	ProxyHash tlb.Uint256 // uint256
}

const PrefixRegisterProxy uint64 = 0x927c7cb5

type RegisterProxy struct {
	QueryId   tlb.Uint64 // uint64
	ProxyInfo boc.Cell   // RemainingBitsAndRefs
}

const PrefixUnregisterProxy uint64 = 0x6d49eaf2

type UnregisterProxy struct {
	QueryId tlb.Uint64 // uint64
	Seqno   tlb.Uint32 // uint32
}

const PrefixUpdateProxy uint64 = 0x9c7924ba

type UpdateProxy struct {
	QueryId   tlb.Uint64 // uint64
	Seqno     tlb.Uint32 // uint32
	ProxyAddr boc.Cell   // RemainingBitsAndRefs
}

const PrefixChangeFees uint64 = 0xc52ed8d4

type ChangeFees struct {
	QueryId           tlb.Uint64 // uint64
	PricePerToken     tlb.Coins  // coins
	WorkerFeePerToken tlb.Coins  // coins
}

const PrefixChangeParams uint64 = 0x022fa189

type ChangeParams struct {
	QueryId                tlb.Uint64 // uint64
	PricePerToken          tlb.Coins  // coins
	WorkerFeePerToken      tlb.Coins  // coins
	ProxyDelayBeforeClose  tlb.Uint32 // uint32
	ClientDelayBeforeClose tlb.Uint32 // uint32
	MinProxyStake          tlb.Coins  // coins
	MinClientStake         tlb.Coins  // coins
}

const PrefixUpgradeContracts uint64 = 0xa2370f61

type UpgradeContracts struct {
	QueryId    tlb.Uint64 // uint64
	ProxyCode  boc.Cell   // cell
	WorkerCode boc.Cell   // cell
	ClientCode boc.Cell   // cell
}

const PrefixUpgradeCode uint64 = 0x11aefd51

type UpgradeCode struct {
	QueryId tlb.Uint64 // uint64
	NewCode boc.Cell   // cell
}

const PrefixResetRoot uint64 = 0x563c1d96

type ResetRoot struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixUpgradeFull uint64 = 0x4f7c5789

type UpgradeFull struct {
	QueryId tlb.Uint64 // uint64
	NewData boc.Cell   // cell
	NewCode boc.Cell   // cell
}

const PrefixChangeOwner uint64 = 0xc4a1ae54

type ChangeOwner struct {
	QueryId         tlb.Uint64          // uint64
	NewOwnerAddress tlb.InternalAddress // address
}
type RegisteredProxy struct {
	Kind    tlb.Uint1 // uint1
	Address string    // string
}
type RootData struct {
	ProxyHashes       tlb.HashmapE[tlb.Uint256, boc.Cell]       // map<uint256, slice>
	RegisteredProxies tlb.HashmapE[tlb.Uint32, RegisteredProxy] // map<uint32, RegisteredProxy>
	LastProxySeqno    tlb.Uint32                                // uint32
	WorkerHashes      tlb.HashmapE[tlb.Uint256, boc.Cell]       // map<uint256, slice>
	ModelHashes       tlb.HashmapE[tlb.Uint256, boc.Cell]       // map<uint256, slice>
}
type RootStorage struct {
	OwnerAddress tlb.InternalAddress     // address
	Data         tlb.RefT[*RootData]     // Cell<RootData>
	Params       tlb.RefT[*CocoonParams] // Cell<CocoonParams>
	Version      tlb.Uint32              // uint32
}
type CocoonData struct {
	Version           tlb.Uint32          // uint32
	LastProxySeqno    tlb.Uint32          // uint32
	ParamsVersion     tlb.Uint32          // uint32
	UniqueId          tlb.Uint32          // uint32
	IsTest            bool                // bool
	PricePerToken     tlb.Coins           // coins
	WorkerFeePerToken tlb.Coins           // coins
	MinProxyStake     tlb.Coins           // coins
	MinClientStake    tlb.Coins           // coins
	OwnerAddress      tlb.InternalAddress // address
}
type CurrentCocoonParams struct {
	ParamsVersion                  tlb.Uint32  // uint32
	UniqueId                       tlb.Uint32  // uint32
	IsTest                         bool        // bool
	PricePerToken                  tlb.Coins   // coins
	WorkerFeePerToken              tlb.Coins   // coins
	CachedTokensPriceMultiplier    tlb.Uint32  // uint32
	ReasoningTokensPriceMultiplier tlb.Uint32  // uint32
	ProxyDelayBeforeClose          tlb.Uint32  // uint32
	ClientDelayBeforeClose         tlb.Uint32  // uint32
	MinProxyStake                  tlb.Coins   // coins
	MinClientStake                 tlb.Coins   // coins
	ProxyCodeHash                  tlb.Uint256 // uint256
	WorkerCodeHash                 tlb.Uint256 // uint256
	ClientCodeHash                 tlb.Uint256 // uint256
}

const ( // errors
)

func DecodeGetLastProxySeqno(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetLastProxySeqno = 0x1006F

func GetLastProxySeqno(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetLastProxySeqno, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetLastProxySeqno(&stack)
}

func DecodeGetCocoonData(stack *tlb.VmStack) (result CocoonData, err error) {
	if stack.Len() != 10 {
		err = fmt.Errorf("invalid stack size %d, expected 10", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCocoonData = 0x17965

func GetCocoonData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CocoonData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCocoonData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCocoonData(&stack)
}

func DecodeGetCurParams(stack *tlb.VmStack) (result CurrentCocoonParams, err error) {
	if stack.Len() != 14 {
		err = fmt.Errorf("invalid stack size %d, expected 14", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCurParams = 0x15D71

func GetCurParams(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CurrentCocoonParams, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCurParams, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCurParams(&stack)
}

func DecodeGetProxyHashIsValid(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetProxyHashIsValid = 0x1F965

func GetProxyHashIsValid(ctx context.Context, executor Executor, reqAccountID ton.AccountID, hash tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(hash)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetProxyHashIsValid, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetProxyHashIsValid(&stack)
}

func DecodeGetWorkerHashIsValid(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetWorkerHashIsValid = 0x17609

func GetWorkerHashIsValid(ctx context.Context, executor Executor, reqAccountID ton.AccountID, hash tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(hash)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetWorkerHashIsValid, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetWorkerHashIsValid(&stack)
}

func DecodeGetModelHashIsValid(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetModelHashIsValid = 0x11B8B

func GetModelHashIsValid(ctx context.Context, executor Executor, reqAccountID ton.AccountID, hash tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(hash)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetModelHashIsValid, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetModelHashIsValid(&stack)
}

func CocoonRoot_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage RootStorage, err error) {
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

type cocoonRootImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewCocoonRoot(executor Executor, storageExecutor StorageExecutor) CocoonRoot {
	return &cocoonRootImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c cocoonRootImpl) WithAccountId(accountID ton.AccountID) CocoonRootWithAccount {
	return &cocoonRootWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c cocoonRootImpl) GetLastProxySeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetLastProxySeqno(ctx, c.executor, reqAccountID)
}

func (c cocoonRootImpl) GetCocoonData(ctx context.Context, reqAccountID ton.AccountID) (CocoonData, error) {
	return GetCocoonData(ctx, c.executor, reqAccountID)
}

func (c cocoonRootImpl) GetCurParams(ctx context.Context, reqAccountID ton.AccountID) (CurrentCocoonParams, error) {
	return GetCurParams(ctx, c.executor, reqAccountID)
}

func (c cocoonRootImpl) GetProxyHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error) {
	return GetProxyHashIsValid(ctx, c.executor, reqAccountID, hash)
}

func (c cocoonRootImpl) GetWorkerHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error) {
	return GetWorkerHashIsValid(ctx, c.executor, reqAccountID, hash)
}

func (c cocoonRootImpl) GetModelHashIsValid(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error) {
	return GetModelHashIsValid(ctx, c.executor, reqAccountID, hash)
}

func (c cocoonRootImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, RootStorage, error) {
	return CocoonRoot_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type cocoonRootWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c cocoonRootWithAccountImpl) GetLastProxySeqno(ctx context.Context) (tlb.Int257, error) {
	return GetLastProxySeqno(ctx, c.executor, c.accountID)
}

func (c cocoonRootWithAccountImpl) GetCocoonData(ctx context.Context) (CocoonData, error) {
	return GetCocoonData(ctx, c.executor, c.accountID)
}

func (c cocoonRootWithAccountImpl) GetCurParams(ctx context.Context) (CurrentCocoonParams, error) {
	return GetCurParams(ctx, c.executor, c.accountID)
}

func (c cocoonRootWithAccountImpl) GetProxyHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error) {
	return GetProxyHashIsValid(ctx, c.executor, c.accountID, hash)
}

func (c cocoonRootWithAccountImpl) GetWorkerHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error) {
	return GetWorkerHashIsValid(ctx, c.executor, c.accountID, hash)
}

func (c cocoonRootWithAccountImpl) GetModelHashIsValid(ctx context.Context, hash tlb.Int257) (tlb.Int257, error) {
	return GetModelHashIsValid(ctx, c.executor, c.accountID, hash)
}

func (c cocoonRootWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, RootStorage, error) {
	return CocoonRoot_AccountState(ctx, c.storageExecutor, c.accountID)
}
