// Code generated - DO NOT EDIT.

package abiFfVault

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixUnstakeRequest uint64 = 0x8a77cf39

type UnstakeRequest struct {
	QueryId               tlb.Uint64 // uint64
	AmountWantedToUnstake tlb.Coins  // coins
}

const PrefixUnstakeExecuteInternal uint64 = 0xd1029c99

type UnstakeExecuteInternal struct {
	QueryId         tlb.Uint64  // uint64
	AmountToUnstake tlb.Coins   // coins
	XautPriceActual tlb.Uint256 // uint256
}

const PrefixUnstakeExecuteCancel uint64 = 0x0b30393c

type UnstakeExecuteCancel struct {
	QueryId      tlb.Uint64          // uint64
	ExcessGetter tlb.InternalAddress // address
	AmountToAdd  tlb.Coins           // coins
}
type NftData struct {
	Init               bool                // bool
	Index              tlb.Int257          // int
	Collection_address tlb.InternalAddress // address
	Owner_address      tlb.MsgAddress      // address?
	Content            tlb.Maybe[boc.Cell] // cell?
}
type StakePositionInfo struct {
	LockedAssetAmount               tlb.Int257 // int
	AmountWantedToUnstake           tlb.Int257 // int
	PreviousStablePerAssetCollected tlb.Int257 // int
	PriceStake                      tlb.Int257 // int
}

const ( // errors
)

func DecodeGetNftData(stack *tlb.VmStack) (result NftData, err error) {
	if stack.Len() != 5 {
		err = fmt.Errorf("invalid stack size %d, expected 5", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetNftData = 0x18FCF

func GetNftData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result NftData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetNftData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetNftData(&stack)
}

func DecodeGetStakePositionInfo(stack *tlb.VmStack) (result StakePositionInfo, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetStakePositionInfo = 0x190F0

func GetStakePositionInfo(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result StakePositionInfo, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetStakePositionInfo, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetStakePositionInfo(&stack)
}

type positionImpl struct {
	executor Executor
}

func NewPosition(executor Executor) Position {
	return &positionImpl{executor: executor}
}

func (c positionImpl) WithAccountId(accountID ton.AccountID) PositionWithAccount {
	return &positionWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c positionImpl) GetNftData(ctx context.Context, reqAccountID ton.AccountID) (NftData, error) {
	return GetNftData(ctx, c.executor, reqAccountID)
}

func (c positionImpl) GetStakePositionInfo(ctx context.Context, reqAccountID ton.AccountID) (StakePositionInfo, error) {
	return GetStakePositionInfo(ctx, c.executor, reqAccountID)
}

type positionWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c positionWithAccountImpl) GetNftData(ctx context.Context) (NftData, error) {
	return GetNftData(ctx, c.executor, c.accountID)
}

func (c positionWithAccountImpl) GetStakePositionInfo(ctx context.Context) (StakePositionInfo, error) {
	return GetStakePositionInfo(ctx, c.executor, c.accountID)
}
