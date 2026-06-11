// Code generated - DO NOT EDIT.

package abiStonkspump

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type StonksPumpVirtualFactoryMessage CreateVirtualLiquidityJetton

const PrefixCreateVirtualLiquidityJetton uint64 = 0xcc1a9801

type CreateVirtualLiquidityJetton struct {
	QueryId         tlb.Uint64          // uint64
	Content         boc.Cell            // cell
	DevBuyTonAmount tlb.Coins           // coins
	MaxBuyPercent   tlb.Uint16          // uint16
	DevWallet       tlb.InternalAddress // address
}
