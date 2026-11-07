// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type ClmmRangeInitPrice struct {
	InitBin       tlb.Int64   // int64
	InitSqrtPrice tlb.Uint256 // uint256
}

const PrefixClmmRangeDeploy uint64 = 0x724b6005

type ClmmRangeDeploy struct {
	QueryId        tlb.Uint64                    // uint64
	XMaster        tlb.MsgAddress                // address?
	YMaster        tlb.InternalAddress           // address
	BinStep        tlb.Uint16                    // uint16
	LpFee          tlb.Uint8                     // uint8
	ProtocolFee    tlb.Uint8                     // uint8
	InitPrice      tlb.RefT[*ClmmRangeInitPrice] // Cell<ClmmRangeInitPrice>
	LpAccCode      boc.Cell                      // cell
	LpPositionCode boc.Cell                      // cell
}

type ClmmTonSwapParams struct {
	SwapAmount     tlb.Coins           // coins
	Excess         tlb.InternalAddress // address
	SqrtPriceLimit tlb.Uint256         // uint256
	NextFulfill    tlb.Maybe[boc.Cell] // cell?
	NextReject     tlb.Maybe[boc.Cell] // cell?
}

const PrefixClmmRangeSwap uint64 = 0xd151210f

type ClmmRangeSwap struct {
	QueryId    tlb.Uint64        // uint64
	SwapParams ClmmTonSwapParams // ClmmTonSwapParams
}

const PrefixClmmRangeDepositLiquidity uint64 = 0xec401442

type ClmmRangeDepositLiquidity struct {
	QueryId       tlb.Uint64                              // uint64
	DepositAmount tlb.Coins                               // coins
	IsSingle      bool                                    // bool
	DepositBins   tlb.HashmapE[tlb.Int64, ClmmDepositBin] // map<int64, ClmmDepositBin>
}

const PrefixClmmRangeProvideLiquidity uint64 = 0x344629fc

type ClmmRangeProvideLiquidity struct {
	QueryId     tlb.Uint64                              // uint64
	Provider    tlb.InternalAddress                     // address
	XAmount     tlb.Coins                               // coins
	YAmount     tlb.Coins                               // coins
	DepositBins tlb.HashmapE[tlb.Int64, ClmmDepositBin] // map<int64, ClmmDepositBin>
}

const PrefixClmmRangeWithdrawFees uint64 = 0x002618d9

type ClmmRangeWithdrawFees struct {
	QueryId      tlb.Uint64                       // uint64
	Provider     tlb.InternalAddress              // address
	NftIndex     tlb.Uint64                       // uint64
	WithdrawBins tlb.HashmapE[tlb.Int64, ClmmBin] // map<int64, ClmmBin>
}

const PrefixClmmRangeWithdrawLiquidity uint64 = 0x69099bab

type ClmmRangeWithdrawLiquidity struct {
	QueryId      tlb.Uint64                       // uint64
	Provider     tlb.InternalAddress              // address
	NftIndex     tlb.Uint64                       // uint64
	WithdrawBins tlb.HashmapE[tlb.Int64, ClmmBin] // map<int64, ClmmBin>
}

const PrefixClmmRangeCancelDeposit uint64 = 0x2dcb38cb

type ClmmRangeCancelDeposit struct {
	QueryId  tlb.Uint64          // uint64
	XAmount  tlb.Coins           // coins
	YAmount  tlb.Coins           // coins
	Provider tlb.InternalAddress // address
}

const PrefixClmmRangeCollectProtocolFees uint64 = 0xdbc809b5

type ClmmRangeCollectProtocolFees struct {
	QueryId   tlb.Uint64          // uint64
	Recipient tlb.InternalAddress // address
}

const PrefixClmmRangeUpdateFees uint64 = 0x90be29eb

type ClmmRangeUpdateFees struct {
	QueryId     tlb.Uint64 // uint64
	LpFee       tlb.Uint8  // uint8
	ProtocolFee tlb.Uint8  // uint8
}

const PrefixClmmRangeUpdateWallets uint64 = 0xa17882c8

type ClmmRangeUpdateWallets struct {
	QueryId tlb.Uint64          // uint64
	XWallet tlb.MsgAddress      // address?
	YWallet tlb.InternalAddress // address
}

const PrefixClmmRangeUpgrade uint64 = 0xf0262313

type ClmmRangeUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}

type ClmmJettonSwapParams struct {
	Excess         tlb.InternalAddress // address
	SqrtPriceLimit tlb.Uint256         // uint256
	NextFulfill    tlb.Maybe[boc.Cell] // cell?
	NextReject     tlb.Maybe[boc.Cell] // cell?
}

const PrefixClmmRangeSwapPayload uint64 = 0xd151210f

type ClmmRangeSwapPayload struct {
	SwapParams ClmmJettonSwapParams // ClmmJettonSwapParams
}

const PrefixClmmRangeDepositLiquidityPayload uint64 = 0xec401442

type ClmmRangeDepositLiquidityPayload struct {
	IsSingle    bool                                    // bool
	DepositBins tlb.HashmapE[tlb.Int64, ClmmDepositBin] // map<int64, ClmmDepositBin>
}

const PrefixClmmRangeDepositLiquidityExcess uint64 = 0x89983597

type ClmmRangeDepositLiquidityExcess struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmRangeDepositLiquidityRefund uint64 = 0x77714776

type ClmmRangeDepositLiquidityRefund struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmRangeWithdrawLiquidityPayout uint64 = 0xacae4171

type ClmmRangeWithdrawLiquidityPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmRangeWithdrawFeesPayout uint64 = 0xbc49e37b

type ClmmRangeWithdrawFeesPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmRangeCancelDepositRefund uint64 = 0x5db3ae19

type ClmmRangeCancelDepositRefund struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmRangeProtocolFeesPayout uint64 = 0x9d6c9b24

type ClmmRangeProtocolFeesPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmExcess uint64 = 0x5fd5f2b3

type ClmmExcess struct {
	QueryId tlb.Uint64 // uint64
}

type ClmmRangeState struct {
	CurrentBin    tlb.Int64   // int64
	CurrentSqrtP  tlb.Uint256 // uint256
	XProtocolFees tlb.Coins   // coins
	YProtocolFees tlb.Coins   // coins
}

type ClmmRangeStorage struct {
	Pool           tlb.InternalAddress              // address
	RangeNum       tlb.Int16                        // int16
	XWallet        tlb.MsgAddress                   // address?
	YWallet        tlb.InternalAddress              // address
	BinStep        tlb.Uint16                       // uint16
	LpFee          tlb.Uint8                        // uint8
	ProtocolFee    tlb.Uint8                        // uint8
	PositionIndex  tlb.Uint64                       // uint64
	RangeState     tlb.RefT[*ClmmRangeState]        // Cell<ClmmRangeState>
	Bins           tlb.HashmapE[tlb.Int64, ClmmBin] // map<int64, ClmmBin>
	LpAccountCode  boc.Cell                         // cell
	LpPositionCode boc.Cell                         // cell
}

func DecodeGetRange(stack *tlb.VmStack) (result tlb.RefT[*ClmmRangeStorage], err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return (func() (value tlb.RefT[*ClmmRangeStorage], err error) {
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

const MethodIDGetRange = 0x19323

func GetRange(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.RefT[*ClmmRangeStorage], err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetRange, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetRange(&stack)
}

type clmmRangeImpl struct {
	executor Executor
}

func NewClmmRange(executor Executor) ClmmRange {
	return &clmmRangeImpl{executor: executor}
}

func (c clmmRangeImpl) WithAccountId(accountID ton.AccountID) ClmmRangeWithAccount {
	return &clmmRangeWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c clmmRangeImpl) GetRange(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*ClmmRangeStorage], error) {
	return GetRange(ctx, c.executor, reqAccountID)
}

type clmmRangeWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c clmmRangeWithAccountImpl) GetRange(ctx context.Context) (tlb.RefT[*ClmmRangeStorage], error) {
	return GetRange(ctx, c.executor, c.accountID)
}
