// Code generated - DO NOT EDIT.

package abiSingleNominatorPool

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixWithdraw uint64 = 0x00001000

type Withdraw struct {
	QueryId tlb.Uint64 // uint64
	Amount  tlb.Coins  // coins
}

const PrefixChangeValidatorAddress uint64 = 0x00001001

type ChangeValidatorAddress struct {
	QueryId             tlb.Uint64          // uint64
	NewValidatorAddress tlb.InternalAddress // address
}

const PrefixSendRawMsg uint64 = 0x00007702

type SendRawMsg struct {
	QueryId tlb.Uint64 // uint64
	Mode    tlb.Uint8  // uint8
	Msg     boc.Cell   // cell
}

const PrefixUpgrade uint64 = 0x00009903

type Upgrade struct {
	QueryId tlb.Uint64 // uint64
	Code    boc.Cell   // cell
}
type Roles struct {
	Owner     tlb.InternalAddress // address
	Validator tlb.InternalAddress // address
}
type PoolConfig struct {
	MinStake     tlb.Int257 // int
	DepositFee   tlb.Int257 // int
	WithdrawFee  tlb.Int257 // int
	PoolFee      tlb.Int257 // int
	ReceiptPrice tlb.Int257 // int
}
type PoolData struct {
	State                    tlb.Int257          // int
	NominatorsCount          tlb.Int257          // int
	StakeAmountSent          tlb.Int257          // int
	ValidatorAmount          tlb.Int257          // int
	PoolConfig               PoolConfig          // PoolConfig
	Nominators               tlb.Maybe[boc.Cell] // cell?
	WithdrawRequests         tlb.Maybe[boc.Cell] // cell?
	StakeAt                  tlb.Int257          // int
	SavedValidatorSetHash    tlb.Int257          // int
	ValidatorSetChangesCount tlb.Int257          // int
	ValidatorSetChangeTime   tlb.Int257          // int
	StakeHeldFor             tlb.Int257          // int
	ConfigProposalVotings    tlb.Maybe[boc.Cell] // cell?
}

const ( // errors
	ErrorWrongSetCode        = 0x2002 // 8194
	ErrorInsufficientBalance = 0x2004 // 8196
)

func DecodeGetRoles(stack *tlb.VmStack) (result Roles, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetRoles = 0x1FCA0

func GetRoles(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result Roles, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetRoles, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetRoles(&stack)
}

func DecodeGetPoolData(stack *tlb.VmStack) (result PoolData, err error) {
	if stack.Len() != 17 {
		err = fmt.Errorf("invalid stack size %d, expected 17", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPoolData = 0x13F19

func GetPoolData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result PoolData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPoolData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPoolData(&stack)
}

type singleNominatorPoolImpl struct {
	executor Executor
}

func NewSingleNominatorPool(executor Executor) SingleNominatorPool {
	return &singleNominatorPoolImpl{executor: executor}
}

func (c singleNominatorPoolImpl) WithAccountId(accountID ton.AccountID) SingleNominatorPoolWithAccount {
	return &singleNominatorPoolWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c singleNominatorPoolImpl) GetRoles(ctx context.Context, reqAccountID ton.AccountID) (Roles, error) {
	return GetRoles(ctx, c.executor, reqAccountID)
}

func (c singleNominatorPoolImpl) GetPoolData(ctx context.Context, reqAccountID ton.AccountID) (PoolData, error) {
	return GetPoolData(ctx, c.executor, reqAccountID)
}

type singleNominatorPoolWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c singleNominatorPoolWithAccountImpl) GetRoles(ctx context.Context) (Roles, error) {
	return GetRoles(ctx, c.executor, c.accountID)
}

func (c singleNominatorPoolWithAccountImpl) GetPoolData(ctx context.Context) (PoolData, error) {
	return GetPoolData(ctx, c.executor, c.accountID)
}
