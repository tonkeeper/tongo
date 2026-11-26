// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixClmmAccountDepositLiquidityRecord uint64 = 0xed8efa06

type ClmmAccountDepositLiquidityRecord struct {
	QueryId     tlb.Uint64                              // uint64
	Amount      tlb.Coins                               // coins
	XIn         bool                                    // bool
	DepositBins tlb.HashmapE[tlb.Int64, ClmmDepositBin] // map<int64, ClmmDepositBin>
}
