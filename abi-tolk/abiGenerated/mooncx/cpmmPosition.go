// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type CpmmPositionContent struct {
	Lp      tlb.Coins   // coins
	XFeesLp tlb.Uint256 // uint256
	YFeesLp tlb.Uint256 // uint256
}

const PrefixCpmmPositionDeploy uint64 = 0x7b5ceb3f

type CpmmPositionDeploy struct {
	QueryId tlb.Uint64                     // uint64
	Content tlb.RefT[*CpmmPositionContent] // Cell<CpmmPositionContent>
	Locked  bool                           // bool
}

const PrefixCpmmPositionWithdrawLiquidity uint64 = 0x84548afa

type CpmmPositionWithdrawLiquidity struct {
	QueryId tlb.Uint64 // uint64
	Amount  tlb.Coins  // coins
}

const PrefixCpmmPositionWithdrawFees uint64 = 0xf325e2f4

type CpmmPositionWithdrawFees struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPositionWithdrawFeesUpdate uint64 = 0x90ad0f98

type CpmmPositionWithdrawFeesUpdate struct {
	QueryId tlb.Uint64  // uint64
	XFeesLp tlb.Uint256 // uint256
	YFeesLp tlb.Uint256 // uint256
}

const PrefixCpmmPositionLockLiquidity uint64 = 0xb751f838

type CpmmPositionLockLiquidity struct {
	QueryId tlb.Uint64 // uint64
}
type CpmmPositionStorage struct {
	Pool   tlb.InternalAddress            // address
	User   tlb.InternalAddress            // address
	Data   tlb.RefT[*CpmmPositionContent] // Cell<CpmmPositionContent>
	Frozen bool                           // bool
	Locked bool                           // bool
}

const ( // errors
)

func DecodeGetPosition(stack *tlb.VmStack) (result tlb.RefT[*CpmmPositionStorage], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return (func() (value tlb.RefT[*CpmmPositionStorage], err error) {
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

const MethodIDGetPosition = 0x1AA08

func GetPosition(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.RefT[*CpmmPositionStorage], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPosition, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPosition(&stack)
}

func DecodeGetLocked(stack *tlb.VmStack) (result bool, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadBool()
}

const MethodIDGetLocked = 0x1A590

func GetLocked(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result bool, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetLocked, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetLocked(&stack)
}

type cpmmPositionImpl struct {
	executor Executor
}

func NewCpmmPosition(executor Executor) CpmmPosition {
	return &cpmmPositionImpl{executor: executor}
}

func (c cpmmPositionImpl) WithAccountId(accountID ton.AccountID) CpmmPositionWithAccount {
	return &cpmmPositionWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c cpmmPositionImpl) GetPosition(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*CpmmPositionStorage], error) {
	return GetPosition(ctx, c.executor, reqAccountID)
}

func (c cpmmPositionImpl) GetLocked(ctx context.Context, reqAccountID ton.AccountID) (bool, error) {
	return GetLocked(ctx, c.executor, reqAccountID)
}

type cpmmPositionWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c cpmmPositionWithAccountImpl) GetPosition(ctx context.Context) (tlb.RefT[*CpmmPositionStorage], error) {
	return GetPosition(ctx, c.executor, c.accountID)
}

func (c cpmmPositionWithAccountImpl) GetLocked(ctx context.Context) (bool, error) {
	return GetLocked(ctx, c.executor, c.accountID)
}
