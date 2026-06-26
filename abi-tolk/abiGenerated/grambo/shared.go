// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixGramboLaunch uint64 = 0x372f1ae3

type GramboLaunch struct {
	QueryId           tlb.Uint64          // uint64
	Buyer             tlb.InternalAddress // address
	VirtualTonReserve tlb.Coins           // coins
	TonTarget         tlb.Coins           // coins
	PrebuyAmount      tlb.Coins           // coins
}

type GramboSellParams struct {
	MinTonOut     tlb.Coins           // coins
	Referrer      tlb.InternalAddress // address
	CustomPayload tlb.Maybe[boc.Cell] // cell?
}

const PrefixGramboSellNotification uint64 = 0x31fc1276

type GramboSellNotification struct {
	QueryId             tlb.Uint64                             // uint64
	Amount              tlb.Coins                              // coins
	From                tlb.InternalAddress                    // address
	ResponseDestination tlb.InternalAddress                    // address
	SellParams          tlb.Maybe[tlb.RefT[*GramboSellParams]] // Cell<GramboSellParams>?
}

const PrefixGramboActivateWallet uint64 = 0x6b75d3e2

type GramboActivateWallet struct {
	QueryId tlb.Uint64 // uint64
}
