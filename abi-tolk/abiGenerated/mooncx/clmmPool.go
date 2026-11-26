// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type ClmmPayloadCode struct {
	RangeCode      boc.Cell // cell
	LpAccCode      boc.Cell // cell
	LpPositionCode boc.Cell // cell
}

const PrefixClmmPoolDeploy uint64 = 0x6d33fdcd

type ClmmPoolDeploy struct {
	QueryId     tlb.Uint64                 // uint64
	InitBin     tlb.Int64                  // int64
	LpFee       tlb.Uint8                  // uint8
	ProtocolFee tlb.Uint8                  // uint8
	PayloadCode tlb.RefT[*ClmmPayloadCode] // Cell<ClmmPayloadCode>
}

const PrefixClmmPoolCollectProtocolFees uint64 = 0xefef86da

type ClmmPoolCollectProtocolFees struct {
	QueryId   tlb.Uint64          // uint64
	RangeNum  tlb.Int16           // int16
	Recipient tlb.InternalAddress // address
}

const PrefixClmmPoolUpdateRangeWallets uint64 = 0xe82e11a1

type ClmmPoolUpdateRangeWallets struct {
	QueryId  tlb.Uint64          // uint64
	RangeNum tlb.Int16           // int16
	XWallet  tlb.MsgAddress      // address?
	YWallet  tlb.InternalAddress // address
}

const PrefixClmmPoolUpdateRangeFees uint64 = 0xd97a5c35

type ClmmPoolUpdateRangeFees struct {
	QueryId     tlb.Uint64 // uint64
	RangeNum    tlb.Int16  // int16
	LpFee       tlb.Uint8  // uint8
	ProtocolFee tlb.Uint8  // uint8
}

const PrefixClmmPoolUpgradeRange uint64 = 0x7791fc81

type ClmmPoolUpgradeRange struct {
	QueryId  tlb.Uint64          // uint64
	RangeNum tlb.Int16           // int16
	NewCode  tlb.Maybe[boc.Cell] // cell?
	NewData  tlb.Maybe[boc.Cell] // cell?
}

const PrefixClmmPoolUpgrade uint64 = 0xb7f86a0a

type ClmmPoolUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}
type ClmmPoolStorage struct {
	Factory        tlb.InternalAddress // address
	XMaster        tlb.MsgAddress      // address?
	YMaster        tlb.InternalAddress // address
	BinStep        tlb.Uint16          // uint16
	LpFee          tlb.Uint8           // uint8
	ProtocolFee    tlb.Uint8           // uint8
	RangeCode      boc.Cell            // cell
	LpAccCode      boc.Cell            // cell
	LpPositionCode boc.Cell            // cell
}

const ( // errors
)

func DecodeGetPool(stack *tlb.VmStack) (result tlb.RefT[*ClmmPoolStorage], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return (func() (value tlb.RefT[*ClmmPoolStorage], err error) {
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

const MethodIDGetPool = 0x15698

func GetPool(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.RefT[*ClmmPoolStorage], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPool, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPool(&stack)
}

func DecodeGetRanges(stack *tlb.VmStack) (result tlb.Hashmap[tlb.Int16, tlb.InternalAddress], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return tlb.ReadHashmapFromStack[tlb.Int16, tlb.InternalAddress](stack)
}

const MethodIDGetRanges = 0x1DE2E

func GetRanges(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Hashmap[tlb.Int16, tlb.InternalAddress], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetRanges, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetRanges(&stack)
}

type clmmPoolImpl struct {
	executor Executor
}

func NewClmmPool(executor Executor) ClmmPool {
	return &clmmPoolImpl{executor: executor}
}

func (c clmmPoolImpl) WithAccountId(accountID ton.AccountID) ClmmPoolWithAccount {
	return &clmmPoolWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c clmmPoolImpl) GetPool(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*ClmmPoolStorage], error) {
	return GetPool(ctx, c.executor, reqAccountID)
}

func (c clmmPoolImpl) GetRanges(ctx context.Context, reqAccountID ton.AccountID) (tlb.Hashmap[tlb.Int16, tlb.InternalAddress], error) {
	return GetRanges(ctx, c.executor, reqAccountID)
}

type clmmPoolWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c clmmPoolWithAccountImpl) GetPool(ctx context.Context) (tlb.RefT[*ClmmPoolStorage], error) {
	return GetPool(ctx, c.executor, c.accountID)
}

func (c clmmPoolWithAccountImpl) GetRanges(ctx context.Context) (tlb.Hashmap[tlb.Int16, tlb.InternalAddress], error) {
	return GetRanges(ctx, c.executor, c.accountID)
}
