// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixClmmPositionDeploy uint64 = 0x4e5504a9

type ClmmPositionDeploy struct {
	QueryId      tlb.Uint64                       // uint64
	Provider     tlb.InternalAddress              // address
	PositionBins tlb.HashmapE[tlb.Int64, ClmmBin] // map<int64, ClmmBin>
}

const PrefixClmmPositionWithdrawLiquidity uint64 = 0xc84175bb

type ClmmPositionWithdrawLiquidity struct {
	QueryId      tlb.Uint64                    // uint64
	WithdrawBins tlb.HashmapE[tlb.Int64, bool] // map<int64, bool>
}

const PrefixClmmPositionWithdrawFees uint64 = 0x716c13d8

type ClmmPositionWithdrawFees struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmPositionWithdrawFeesUpdate uint64 = 0xfe49ce3f

type ClmmPositionWithdrawFeesUpdate struct {
	QueryId     tlb.Uint64                       // uint64
	UpdatedBins tlb.HashmapE[tlb.Int64, ClmmBin] // map<int64, ClmmBin>
}

const PrefixClmmPositionWithdrawLiquidityUpdate uint64 = 0x18256917

type ClmmPositionWithdrawLiquidityUpdate struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmPositionLockLiquidity uint64 = 0x850bdd2a

type ClmmPositionLockLiquidity struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixClmmPositionClosed uint64 = 0x80381933

type ClmmPositionClosed struct {
	QueryId tlb.Uint64 // uint64
}
