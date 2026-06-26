// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type SourcesRegistryStorage struct {
	MinGram                 tlb.Coins           // coins
	MaxGram                 tlb.Coins           // coins
	AdminAddress            tlb.InternalAddress // address
	VerifierRegistryAddress tlb.InternalAddress // address
	SourceItemCode          boc.Cell            // cell
}

const PrefixDeploySourceItemPayload uint64 = 0x000003ea

type DeploySourceItemPayload struct {
	QueryId              tlb.Uint64               // uint64
	VerifierId           tlb.Uint256              // uint256
	VerifiedCodeCellHash tlb.Uint256              // uint256
	SourceContent        tlb.RefT[*SourceContent] // Cell<SourceContent>
}

const PrefixChangeVerifierRegistry uint64 = 0x000007d3

type ChangeVerifierRegistry struct {
	QueryId             tlb.Uint64          // uint64
	NewVerifierRegistry tlb.InternalAddress // address
}

const PrefixChangeAdmin uint64 = 0x00000bbc

type ChangeAdmin struct {
	QueryId  tlb.Uint64          // uint64
	NewAdmin tlb.InternalAddress // address
}

const PrefixSetSourceItemCode uint64 = 0x00000fa5

type SetSourceItemCode struct {
	QueryId           tlb.Uint64 // uint64
	NewSourceItemCode boc.Cell   // cell
}

const PrefixSetCode uint64 = 0x0000138e

type SetCode struct {
	QueryId tlb.Uint64 // uint64
	NewCode boc.Cell   // cell
}

const PrefixSetDeploymentCosts uint64 = 0x00001777

type SetDeploymentCosts struct {
	QueryId    tlb.Uint64 // uint64
	NewMinGram tlb.Coins  // coins
	NewMaxGram tlb.Coins  // coins
}

type SourcesRegistryInternalMessageKind uint

const (
	SourcesRegistryInternalMessageKind_DeploySourceItemPayload SourcesRegistryInternalMessageKind = 1002
	SourcesRegistryInternalMessageKind_ChangeVerifierRegistry  SourcesRegistryInternalMessageKind = 2003
	SourcesRegistryInternalMessageKind_ChangeAdmin             SourcesRegistryInternalMessageKind = 3004
	SourcesRegistryInternalMessageKind_SetSourceItemCode       SourcesRegistryInternalMessageKind = 4005
	SourcesRegistryInternalMessageKind_SetCode                 SourcesRegistryInternalMessageKind = 5006
	SourcesRegistryInternalMessageKind_SetDeploymentCosts      SourcesRegistryInternalMessageKind = 6007
)

type SourcesRegistryInternalMessage struct {
	SumType                 SourcesRegistryInternalMessageKind
	DeploySourceItemPayload *DeploySourceItemPayload
	ChangeVerifierRegistry  *ChangeVerifierRegistry
	ChangeAdmin             *ChangeAdmin
	SetSourceItemCode       *SetSourceItemCode
	SetCode                 *SetCode
	SetDeploymentCosts      *SetDeploymentCosts
}

type DeploymentCosts struct {
	MinGram tlb.Coins // coins
	MaxGram tlb.Coins // coins
}

func DecodeGetSourceItemAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSourceItemAddress = 0x134D1

func GetSourceItemAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, verifierId tlb.Int257, verifiedCodeCellHash tlb.Int257) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(verifierId)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(verifiedCodeCellHash)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSourceItemAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSourceItemAddress(&stack)
}

func DecodeGetVerifierRegistryAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVerifierRegistryAddress = 0x1101C

func GetVerifierRegistryAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVerifierRegistryAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVerifierRegistryAddress(&stack)
}

func DecodeGetAdminAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAdminAddress = 0x1702F

func GetAdminAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAdminAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAdminAddress(&stack)
}

func DecodeGetDeploymentCosts(stack *tlb.VmStack) (result DeploymentCosts, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetDeploymentCosts = 0x135C3

func GetDeploymentCosts(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result DeploymentCosts, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetDeploymentCosts, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetDeploymentCosts(&stack)
}

func SourcesRegistry_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage SourcesRegistryStorage, err error) {
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

type sourcesRegistryImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewSourcesRegistry(executor Executor, storageExecutor StorageExecutor) SourcesRegistry {
	return &sourcesRegistryImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c sourcesRegistryImpl) WithAccountId(accountID ton.AccountID) SourcesRegistryWithAccount {
	return &sourcesRegistryWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c sourcesRegistryImpl) GetSourceItemAddress(ctx context.Context, reqAccountID ton.AccountID, verifierId tlb.Int257, verifiedCodeCellHash tlb.Int257) (tlb.InternalAddress, error) {
	return GetSourceItemAddress(ctx, c.executor, reqAccountID, verifierId, verifiedCodeCellHash)
}

func (c sourcesRegistryImpl) GetVerifierRegistryAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error) {
	return GetVerifierRegistryAddress(ctx, c.executor, reqAccountID)
}

func (c sourcesRegistryImpl) GetAdminAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error) {
	return GetAdminAddress(ctx, c.executor, reqAccountID)
}

func (c sourcesRegistryImpl) GetDeploymentCosts(ctx context.Context, reqAccountID ton.AccountID) (DeploymentCosts, error) {
	return GetDeploymentCosts(ctx, c.executor, reqAccountID)
}

func (c sourcesRegistryImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, SourcesRegistryStorage, error) {
	return SourcesRegistry_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type sourcesRegistryWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c sourcesRegistryWithAccountImpl) GetSourceItemAddress(ctx context.Context, verifierId tlb.Int257, verifiedCodeCellHash tlb.Int257) (tlb.InternalAddress, error) {
	return GetSourceItemAddress(ctx, c.executor, c.accountID, verifierId, verifiedCodeCellHash)
}

func (c sourcesRegistryWithAccountImpl) GetVerifierRegistryAddress(ctx context.Context) (tlb.InternalAddress, error) {
	return GetVerifierRegistryAddress(ctx, c.executor, c.accountID)
}

func (c sourcesRegistryWithAccountImpl) GetAdminAddress(ctx context.Context) (tlb.InternalAddress, error) {
	return GetAdminAddress(ctx, c.executor, c.accountID)
}

func (c sourcesRegistryWithAccountImpl) GetDeploymentCosts(ctx context.Context) (DeploymentCosts, error) {
	return GetDeploymentCosts(ctx, c.executor, c.accountID)
}

func (c sourcesRegistryWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, SourcesRegistryStorage, error) {
	return SourcesRegistry_AccountState(ctx, c.storageExecutor, c.accountID)
}
