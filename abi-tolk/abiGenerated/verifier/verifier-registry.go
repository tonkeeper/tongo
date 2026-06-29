// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type MessageDescription struct {
	VerifierId tlb.Uint256    // uint256
	ValidUntil tlb.Uint32     // uint32
	SourceAddr tlb.MsgAddress // any_address
	TargetAddr tlb.MsgAddress // any_address
	Msg        boc.Cell       // cell
}
type VerifierSettings struct {
	MultiSigThreshold tlb.Uint8                             // uint8
	PubKeyEndpoints   tlb.HashmapE[tlb.Uint256, tlb.Uint32] // map<uint256, uint32>
	Name              string                                // string
	MarketingUrl      string                                // string
}
type Verifier struct {
	Admin    tlb.MsgAddress   // any_address
	Settings VerifierSettings // VerifierSettings
}
type VerifierRegistryStorage struct {
	Verifiers tlb.HashmapE[tlb.Uint256, Verifier] // map<uint256, Verifier>
}

const PrefixForwardMessage uint64 = 0x75217758

type ForwardMessage struct {
	QueryId    tlb.Uint64                    // uint64
	Msg        tlb.RefT[*MessageDescription] // Cell<MessageDescription>
	Signatures boc.Cell                      // cell
}

const PrefixUpdateVerifier uint64 = 0x6002d61a

type UpdateVerifier struct {
	QueryId    tlb.Uint64       // uint64
	VerifierId tlb.Uint256      // uint256
	Settings   VerifierSettings // VerifierSettings
}

const PrefixRemoveVerifier uint64 = 0x19fa5637

type RemoveVerifier struct {
	QueryId tlb.Uint64  // uint64
	Id      tlb.Uint256 // uint256
}
type VerifierRegistryInternalMessageKind uint

const (
	VerifierRegistryInternalMessageKind_ForwardMessage VerifierRegistryInternalMessageKind = 1965127512
	VerifierRegistryInternalMessageKind_UpdateVerifier VerifierRegistryInternalMessageKind = 1610798618
	VerifierRegistryInternalMessageKind_RemoveVerifier VerifierRegistryInternalMessageKind = 435836471
)

type VerifierRegistryInternalMessage struct { // tagged union
	SumType        VerifierRegistryInternalMessageKind
	ForwardMessage *ForwardMessage
	UpdateVerifier *UpdateVerifier
	RemoveVerifier *RemoveVerifier
}
type VerifierInfo struct {
	Admin    tlb.InternalAddress // address
	Settings boc.Cell            // cell
	Found    bool                // bool
}

const ( // errors
)

func DecodeGetVerifier(stack *tlb.VmStack) (result VerifierInfo, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVerifier = 0x16E70

func GetVerifier(ctx context.Context, executor Executor, reqAccountID ton.AccountID, id tlb.Int257) (result VerifierInfo, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(id)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVerifier, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVerifier(&stack)
}

func DecodeGetVerifiersNum(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVerifiersNum = 0x10C4D

func GetVerifiersNum(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVerifiersNum, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVerifiersNum(&stack)
}

func DecodeGetVerifiers(stack *tlb.VmStack) (result tlb.RefT[*VerifierRegistryStorage], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return (func() (value tlb.RefT[*VerifierRegistryStorage], err error) {
		var cIn boc.Cell
		cIn, err = stack.ReadCell()
		if err != nil {
			return
		}
		c := boc.NewCell()
		_ = c.AddRef(&cIn)
		err = value.UnmarshalTLB(c, tlb.NewDecoder())
		return
	})()
}

const MethodIDGetVerifiers = 0x1B39C

func GetVerifiers(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.RefT[*VerifierRegistryStorage], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVerifiers, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVerifiers(&stack)
}

func VerifierRegistry_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage VerifierRegistryStorage, err error) {
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

type verifierRegistryImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewVerifierRegistry(executor Executor, storageExecutor StorageExecutor) VerifierRegistry {
	return &verifierRegistryImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c verifierRegistryImpl) WithAccountId(accountID ton.AccountID) VerifierRegistryWithAccount {
	return &verifierRegistryWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c verifierRegistryImpl) GetVerifier(ctx context.Context, reqAccountID ton.AccountID, id tlb.Int257) (VerifierInfo, error) {
	return GetVerifier(ctx, c.executor, reqAccountID, id)
}

func (c verifierRegistryImpl) GetVerifiersNum(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetVerifiersNum(ctx, c.executor, reqAccountID)
}

func (c verifierRegistryImpl) GetVerifiers(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*VerifierRegistryStorage], error) {
	return GetVerifiers(ctx, c.executor, reqAccountID)
}

func (c verifierRegistryImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, VerifierRegistryStorage, error) {
	return VerifierRegistry_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type verifierRegistryWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c verifierRegistryWithAccountImpl) GetVerifier(ctx context.Context, id tlb.Int257) (VerifierInfo, error) {
	return GetVerifier(ctx, c.executor, c.accountID, id)
}

func (c verifierRegistryWithAccountImpl) GetVerifiersNum(ctx context.Context) (tlb.Int257, error) {
	return GetVerifiersNum(ctx, c.executor, c.accountID)
}

func (c verifierRegistryWithAccountImpl) GetVerifiers(ctx context.Context) (tlb.RefT[*VerifierRegistryStorage], error) {
	return GetVerifiers(ctx, c.executor, c.accountID)
}

func (c verifierRegistryWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, VerifierRegistryStorage, error) {
	return VerifierRegistry_AccountState(ctx, c.storageExecutor, c.accountID)
}
