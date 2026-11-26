// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixCpmmAccountDepositLiquidityRecord uint64 = 0xfda9374a

type CpmmAccountDepositLiquidityRecord struct {
	QueryId  tlb.Uint64                    // uint64
	XIn      bool                          // bool
	Amount   tlb.Coins                     // coins
	MinLpOut tlb.Coins                     // coins
	Locked   bool                          // bool
	Split    tlb.Maybe[tlb.RefT[*LpSplit]] // Cell<LpSplit>?
}
