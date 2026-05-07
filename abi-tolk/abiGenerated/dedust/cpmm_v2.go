// Code generated - DO NOT EDIT.

package abiDedust

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type PoolStatus tlb.Uint2

const (
	PoolStatusuninitialized PoolStatus = 0
	PoolStatusinitializing  PoolStatus = 1
	PoolStatusinitialized   PoolStatus = 2
)

type Asset struct {
	Value tlb.MsgAddress // address?
}
type Q120X120 struct {
	Value tlb.VarUInteger32 // varuint32
}
type PoolReward struct {
	RemainingTime   tlb.Uint40 // uint40
	RemainingBudget tlb.Coins  // coins
	RewardsPerToken Q120X120   // Q120X120
	LastUpdate      tlb.Uint40 // uint40
}
type DedustCpmmV2GetPoolData struct {
	Status               PoolStatus                                     // PoolStatus
	DepositActive        bool                                           // bool
	SwapActive           bool                                           // bool
	AssetX               Asset                                          // Asset
	AssetY               Asset                                          // Asset
	WalletsByAssets      tlb.HashmapE[tlb.Uint256, tlb.InternalAddress] // map<uint256, address>
	AssetsByWallets      tlb.HashmapE[tlb.Uint256, tlb.InternalAddress] // map<uint256, address>
	WalletsByResolutions tlb.HashmapE[tlb.Uint256, tlb.InternalAddress] // map<uint256, address>
	BaseFeeBPS           tlb.Uint16                                     // uint16
	ReserveX             tlb.Coins                                      // coins
	ReserveY             tlb.Coins                                      // coins
	Liquidity            tlb.Coins                                      // coins
	ProtocolFeeX         tlb.Coins                                      // coins
	ProtocolFeeY         tlb.Coins                                      // coins
	CreatorFeeX          tlb.Coins                                      // coins
	CreatorFeeY          tlb.Coins                                      // coins
	XLPFeePerToken       Q120X120                                       // Q120X120
	YLPFeePerToken       Q120X120                                       // Q120X120
	Rewards              tlb.Maybe[tlb.HashmapE[tlb.Uint2, PoolReward]] // map<uint2, PoolReward>?
}

func DecodeGetPoolData(stack *tlb.VmStack) (result DedustCpmmV2GetPoolData, err error) {
	if stack.Len() != 20 {
		err = fmt.Errorf("invalid stack size %d, expected 20", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPoolData = 0x13F19

func GetPoolData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result DedustCpmmV2GetPoolData, err error) {
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

type cpmmV2Impl struct {
	executor Executor
}

func NewCpmmV2(executor Executor) CpmmV2 {
	return &cpmmV2Impl{executor: executor}
}

func (c cpmmV2Impl) WithAccountId(accountID ton.AccountID) CpmmV2WithAccount {
	return &cpmmV2WithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c cpmmV2Impl) GetPoolData(ctx context.Context, reqAccountID ton.AccountID) (DedustCpmmV2GetPoolData, error) {
	return GetPoolData(ctx, c.executor, reqAccountID)
}

type cpmmV2WithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c cpmmV2WithAccountImpl) GetPoolData(ctx context.Context) (DedustCpmmV2GetPoolData, error) {
	return GetPoolData(ctx, c.executor, c.accountID)
}
