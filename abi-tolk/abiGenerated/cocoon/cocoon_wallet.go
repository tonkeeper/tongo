// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixOwnerWalletSendMessage uint64 = 0x9c69f376

type OwnerWalletSendMessage struct {
	QueryId tlb.Uint64 // uint64
	Mode    tlb.Uint8  // uint8
	Body    boc.Cell   // cell
}

const PrefixTextCommand uint64 = 0x00000000

type TextCommand struct {
	Action tlb.Uint8 // uint8
}
type ForwardMsg struct {
	Mode tlb.Uint8 // uint8
	Msg  boc.Cell  // cell
}
type ForwardMsgs []ForwardMsg
type ExternalSignedMessage struct {
	SubwalletId tlb.Uint32  // uint32
	ValidUntil  tlb.Uint32  // uint32
	MsgSeqno    tlb.Uint32  // uint32
	Forward     ForwardMsgs // ForwardMsgs
}
type WalletExternalMessage struct {
	Signature tlb.Bits512           // bits512
	Message   ExternalSignedMessage // ExternalSignedMessage
}
type WalletStorage struct {
	Seqno        tlb.Uint32          // uint32
	SubwalletId  tlb.Uint32          // uint32
	PublicKey    tlb.Uint256         // uint256
	Status       tlb.Uint32          // uint32
	OwnerAddress tlb.InternalAddress // address
}

const ( // errors
)

func DecodeGetSeqno(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSeqno = 0x14C97

func GetSeqno(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSeqno, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSeqno(&stack)
}

func DecodeGetPublicKey(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPublicKey = 0x1339C

func GetPublicKey(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPublicKey, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPublicKey(&stack)
}

func DecodeGetOwnerAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetOwnerAddress = 0x1BFBB

func GetOwnerAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetOwnerAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetOwnerAddress(&stack)
}

func CocoonWallet_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage WalletStorage, err error) {
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

type cocoonWalletImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewCocoonWallet(executor Executor, storageExecutor StorageExecutor) CocoonWallet {
	return &cocoonWalletImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c cocoonWalletImpl) WithAccountId(accountID ton.AccountID) CocoonWalletWithAccount {
	return &cocoonWalletWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c cocoonWalletImpl) GetSeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetSeqno(ctx, c.executor, reqAccountID)
}

func (c cocoonWalletImpl) GetPublicKey(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetPublicKey(ctx, c.executor, reqAccountID)
}

func (c cocoonWalletImpl) GetOwnerAddress(ctx context.Context, reqAccountID ton.AccountID) (tlb.InternalAddress, error) {
	return GetOwnerAddress(ctx, c.executor, reqAccountID)
}

func (c cocoonWalletImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, WalletStorage, error) {
	return CocoonWallet_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type cocoonWalletWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c cocoonWalletWithAccountImpl) GetSeqno(ctx context.Context) (tlb.Int257, error) {
	return GetSeqno(ctx, c.executor, c.accountID)
}

func (c cocoonWalletWithAccountImpl) GetPublicKey(ctx context.Context) (tlb.Int257, error) {
	return GetPublicKey(ctx, c.executor, c.accountID)
}

func (c cocoonWalletWithAccountImpl) GetOwnerAddress(ctx context.Context) (tlb.InternalAddress, error) {
	return GetOwnerAddress(ctx, c.executor, c.accountID)
}

func (c cocoonWalletWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, WalletStorage, error) {
	return CocoonWallet_AccountState(ctx, c.storageExecutor, c.accountID)
}
