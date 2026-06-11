// Code generated - DO NOT EDIT.

package abiStonkspump

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixAskToPresaleSell uint64 = 0xcc1a97af

type AskToPresaleSell struct {
	QueryId       tlb.Uint64     // uint64
	JettonAmount  tlb.Coins      // coins
	TonsRecipient tlb.MsgAddress // address?
	MinTonOut     tlb.Coins      // coins
}

const PrefixPresaleSellNotificationForMinter uint64 = 0xcc1a97b0

type PresaleSellNotificationForMinter struct {
	QueryId       tlb.Uint64          // uint64
	JettonAmount  tlb.Coins           // coins
	SellInitiator tlb.InternalAddress // address
	TonsRecipient tlb.MsgAddress      // address?
	MinTonOut     tlb.Coins           // coins
}

const PrefixBuyFromPresale uint64 = 0xcc1a97ae

type BuyFromPresale struct {
	QueryId       tlb.Uint64 // uint64
	TonAmount     tlb.Coins  // coins
	MinJettonsOut tlb.Coins  // coins
}

const PrefixPresaleTradeEvent uint64 = 0xcc1a97b3

type PresaleTradeEvent struct {
	QueryId   tlb.Uint64          // uint64
	IsBuy     bool                // bool
	AmountIn  tlb.Coins           // coins
	AmountOut tlb.Coins           // coins
	Recipient tlb.InternalAddress // address
	Leftover  tlb.Coins           // coins
}
type BondingData struct {
	ProgressBps             tlb.Int257 // int
	MigrationTonTarget      tlb.Int257 // int
	CumulativeTonReserve    tlb.Int257 // int
	CumulativeJettonReserve tlb.Int257 // int
	TonsCollected           tlb.Int257 // int
	MaxBuyPercent           tlb.Int257 // int
	TokensSold              tlb.Int257 // int
	PresaleOpen             tlb.Int257 // int
}

func DecodeGetBondingData(stack *tlb.VmStack) (result BondingData, err error) {
	if stack.Len() != 8 {
		err = fmt.Errorf("invalid stack size %d, expected 8", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetBondingData = 0x1D388

func GetBondingData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result BondingData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetBondingData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetBondingData(&stack)
}

type stonkspumpVirtualImpl struct {
	executor Executor
}

func NewStonkspumpVirtual(executor Executor) StonkspumpVirtual {
	return &stonkspumpVirtualImpl{executor: executor}
}

func (c stonkspumpVirtualImpl) WithAccountId(accountID ton.AccountID) StonkspumpVirtualWithAccount {
	return &stonkspumpVirtualWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c stonkspumpVirtualImpl) GetBondingData(ctx context.Context, reqAccountID ton.AccountID) (BondingData, error) {
	return GetBondingData(ctx, c.executor, reqAccountID)
}

type stonkspumpVirtualWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c stonkspumpVirtualWithAccountImpl) GetBondingData(ctx context.Context) (BondingData, error) {
	return GetBondingData(ctx, c.executor, c.accountID)
}
