// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type VaultIncomingMessageKind uint

const (
	VaultIncomingMessageKind_VaultLock           VaultIncomingMessageKind = 612781780
	VaultIncomingMessageKind_VaultUnlock         VaultIncomingMessageKind = 1272347135
	VaultIncomingMessageKind_VaultWithdrawTokens VaultIncomingMessageKind = 1477043330
	VaultIncomingMessageKind_VaultDepositTokens  VaultIncomingMessageKind = 1432280907
)

type VaultIncomingMessage struct { // tagged union
	SumType             VaultIncomingMessageKind
	VaultLock           *VaultLock
	VaultUnlock         *VaultUnlock
	VaultWithdrawTokens *VaultWithdrawTokens
	VaultDepositTokens  *VaultDepositTokens
}
type GetVaultDataResult struct {
	Minter       tlb.InternalAddress // address
	JettonWallet tlb.InternalAddress // address
	Balance      tlb.Coins           // coins
	Owner        tlb.InternalAddress // address
}
type VaultLockAdditionalData struct {
	Excesses                tlb.MsgAddress                          // any_address
	TonSafeDeposit          tlb.Coins                               // coins
	UserFillToVault         bool                                    // bool
	AskJettonWallet         tlb.MsgAddress                          // any_address
	OwnerPubkey             tlb.Maybe[tlb.Uint256]                  // uint256?
	More                    tlb.RefT[*VaultLockAdditionalDataMore]  // Cell<VaultLockAdditionalDataMore>
	UserUnlockForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]]     // Cell<ForwardParams>?
	LockForwardParams       tlb.Maybe[tlb.RefT[*LockForwardParams]] // Cell<LockForwardParams>?
}
type VaultUnlockExtraFields struct {
	Recipient tlb.MsgAddress // any_address
	RefundTo  tlb.MsgAddress // any_address
}

const PrefixVaultLock uint64 = 0x24864ed4

type VaultLock struct {
	QueryId         tlb.Uint64                         // uint64
	Amount          tlb.Coins                          // coins
	QuoteId         tlb.Uint256                        // uint256
	RefundToVault   bool                               // bool
	RefFee          tlb.MsgAddress                     // any_address
	RefFeeTier      tlb.Uint16                         // uint16
	AdditionalData  tlb.RefT[*VaultLockAdditionalData] // Cell<VaultLockAdditionalData>
	LockArgs        tlb.RefT[*BilateralLockArgs]       // Cell<BilateralLockArgs>
	UnlockCondition boc.Cell                           // cell
}

const PrefixVaultUnlock uint64 = 0x4bd679ff

type VaultUnlock struct {
	QueryId             tlb.Uint64                          // uint64
	Amount              tlb.Coins                           // coins
	QuoteId             tlb.Uint256                         // uint256
	FillToVault         bool                                // bool
	RefundToVault       bool                                // bool
	UnlockArgs          tlb.RefT[*BilateralUnlockArgs]      // Cell<BilateralUnlockArgs>
	Excesses            tlb.MsgAddress                      // any_address
	ExtraFields         tlb.RefT[*VaultUnlockExtraFields]   // Cell<VaultUnlockExtraFields>
	UnlockForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixVaultWithdrawTokens uint64 = 0x5809e482

type VaultWithdrawTokens struct {
	QueryId               tlb.Uint64                          // uint64
	Amount                tlb.Coins                           // coins
	Excesses              tlb.MsgAddress                      // any_address
	WithdrawForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixVaultDepositTokens uint64 = 0x555edf4b

type VaultDepositTokens struct {
	QueryId          tlb.Uint64          // uint64
	Amount           tlb.Coins           // coins
	Excesses         tlb.MsgAddress      // any_address
	ForwardTonAmount tlb.Coins           // coins
	ForwardPayload   tlb.Maybe[boc.Cell] // cell?
}

func DecodeGetVaultData(stack *tlb.VmStack) (result GetVaultDataResult, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVaultData = 0x1BFEB

func GetVaultData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetVaultDataResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVaultData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVaultData(&stack)
}

type escrowVaultImpl struct {
	executor Executor
}

func NewEscrowVault(executor Executor) EscrowVault {
	return &escrowVaultImpl{executor: executor}
}

func (c escrowVaultImpl) WithAccountId(accountID ton.AccountID) EscrowVaultWithAccount {
	return &escrowVaultWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c escrowVaultImpl) GetVaultData(ctx context.Context, reqAccountID ton.AccountID) (GetVaultDataResult, error) {
	return GetVaultData(ctx, c.executor, reqAccountID)
}

type escrowVaultWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c escrowVaultWithAccountImpl) GetVaultData(ctx context.Context) (GetVaultDataResult, error) {
	return GetVaultData(ctx, c.executor, c.accountID)
}
