// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixExtClientChargeSigned uint64 = 0xbb63ff93

type ExtClientChargeSigned struct {
	Rest boc.Cell // RemainingBitsAndRefs
}

const PrefixExtClientGrantRefundSigned uint64 = 0xefd711e1

type ExtClientGrantRefundSigned struct {
	Rest boc.Cell // RemainingBitsAndRefs
}

const PrefixExtClientTopUp uint64 = 0xf172e6c2

type ExtClientTopUp struct {
	QueryId        tlb.Uint64          // uint64
	TopUpAmount    tlb.Coins           // coins
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientChangeSecretHashAndTopUp uint64 = 0x8473b408

type OwnerClientChangeSecretHashAndTopUp struct {
	QueryId        tlb.Uint64          // uint64
	TopUpAmount    tlb.Coins           // coins
	NewSecretHash  tlb.Uint256         // uint256
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientRegister uint64 = 0xc45f9f3b

type OwnerClientRegister struct {
	QueryId        tlb.Uint64          // uint64
	Nonce          tlb.Uint64          // uint64
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientChangeSecretHash uint64 = 0xa9357034

type OwnerClientChangeSecretHash struct {
	QueryId        tlb.Uint64          // uint64
	NewSecretHash  tlb.Uint256         // uint256
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientIncreaseStake uint64 = 0x6a1f6a60

type OwnerClientIncreaseStake struct {
	QueryId        tlb.Uint64          // uint64
	NewStake       tlb.Coins           // coins
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientWithdraw uint64 = 0xda068e78

type OwnerClientWithdraw struct {
	QueryId        tlb.Uint64          // uint64
	SendExcessesTo tlb.InternalAddress // address
}

const PrefixOwnerClientRequestRefund uint64 = 0xfafa6cc1

type OwnerClientRequestRefund struct {
	QueryId        tlb.Uint64          // uint64
	SendExcessesTo tlb.InternalAddress // address
}
type ClientConstData struct {
	OwnerAddress   tlb.InternalAddress // address
	ProxyAddress   tlb.InternalAddress // address
	ProxyPublicKey tlb.Uint256         // uint256
}
type ClientStorage struct {
	State        tlb.Uint2                  // uint2
	Balance      tlb.Coins                  // coins
	Stake        tlb.Coins                  // coins
	TokensUsed   tlb.Uint64                 // uint64
	UnlockTs     tlb.Uint32                 // uint32
	SecretHash   tlb.Uint256                // uint256
	ConstDataRef tlb.RefT[*ClientConstData] // Cell<ClientConstData>
	Params       tlb.RefT[*CocoonParams]    // Cell<CocoonParams>
}
type ClientMessageKind uint

const (
	ClientMessageKind_ExtClientChargeSigned               ClientMessageKind = 0xbb63ff93
	ClientMessageKind_ExtClientGrantRefundSigned          ClientMessageKind = 0xefd711e1
	ClientMessageKind_ExtClientTopUp                      ClientMessageKind = 0xf172e6c2
	ClientMessageKind_OwnerClientChangeSecretHashAndTopUp ClientMessageKind = 0x8473b408
	ClientMessageKind_OwnerClientRegister                 ClientMessageKind = 0xc45f9f3b
	ClientMessageKind_OwnerClientChangeSecretHash         ClientMessageKind = 0xa9357034
	ClientMessageKind_OwnerClientIncreaseStake            ClientMessageKind = 0x6a1f6a60
	ClientMessageKind_OwnerClientWithdraw                 ClientMessageKind = 0xda068e78
	ClientMessageKind_OwnerClientRequestRefund            ClientMessageKind = 0xfafa6cc1
)

type ClientMessage struct { // tagged union
	SumType                             ClientMessageKind
	ExtClientChargeSigned               *ExtClientChargeSigned
	ExtClientGrantRefundSigned          *ExtClientGrantRefundSigned
	ExtClientTopUp                      *ExtClientTopUp
	OwnerClientChangeSecretHashAndTopUp *OwnerClientChangeSecretHashAndTopUp
	OwnerClientRegister                 *OwnerClientRegister
	OwnerClientChangeSecretHash         *OwnerClientChangeSecretHash
	OwnerClientIncreaseStake            *OwnerClientIncreaseStake
	OwnerClientWithdraw                 *OwnerClientWithdraw
	OwnerClientRequestRefund            *OwnerClientRequestRefund
}
type CocoonClientData struct {
	OwnerAddress   tlb.InternalAddress // address
	ProxyAddress   tlb.InternalAddress // address
	ProxyPublicKey tlb.Uint256         // uint256
	State          tlb.Uint2           // uint2
	Balance        tlb.Coins           // coins
	Stake          tlb.Coins           // coins
	TokensUsed     tlb.Uint64          // uint64
	UnlockTs       tlb.Uint32          // uint32
	SecretHash     tlb.Uint256         // uint256
}

const ( // errors
)

func DecodeGetCocoonClientData(stack *tlb.VmStack) (result CocoonClientData, err error) {
	if stack.Len() != 9 {
		err = fmt.Errorf("invalid stack size %d, expected 9", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCocoonClientData = 0x12594

func GetCocoonClientData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CocoonClientData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCocoonClientData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCocoonClientData(&stack)
}

func CocoonClient_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage ClientStorage, err error) {
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

type cocoonClientImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewCocoonClient(executor Executor, storageExecutor StorageExecutor) CocoonClient {
	return &cocoonClientImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c cocoonClientImpl) WithAccountId(accountID ton.AccountID) CocoonClientWithAccount {
	return &cocoonClientWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c cocoonClientImpl) GetCocoonClientData(ctx context.Context, reqAccountID ton.AccountID) (CocoonClientData, error) {
	return GetCocoonClientData(ctx, c.executor, reqAccountID)
}

func (c cocoonClientImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ClientStorage, error) {
	return CocoonClient_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type cocoonClientWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c cocoonClientWithAccountImpl) GetCocoonClientData(ctx context.Context) (CocoonClientData, error) {
	return GetCocoonClientData(ctx, c.executor, c.accountID)
}

func (c cocoonClientWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, ClientStorage, error) {
	return CocoonClient_AccountState(ctx, c.storageExecutor, c.accountID)
}
