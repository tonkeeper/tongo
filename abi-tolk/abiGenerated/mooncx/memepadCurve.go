// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type EconState struct {
	TokensForSale tlb.Coins  // coins
	TokensForLp   tlb.Coins  // coins
	VirtualGram   tlb.Coins  // coins
	VirtualTokens tlb.Coins  // coins
	TargetGram    tlb.Coins  // coins
	TradeFeeBps   tlb.Uint16 // uint16
	CreatorFeeBps tlb.Uint16 // uint16
	PoolFeeBps    tlb.Uint8  // uint8
}

const PrefixMemepadCurveInit uint64 = 0x61c9641a

type MemepadCurveInit struct {
	QueryId           tlb.Uint64           // uint64
	CpmmFactory       tlb.InternalAddress  // address
	ProtocolRecipient tlb.InternalAddress  // address
	InitialBuy        tlb.Coins            // coins
	Econ              tlb.RefT[*EconState] // Cell<EconState>
	JwCode            boc.Cell             // cell
	Metadata          boc.Cell             // cell
}

const PrefixMemepadCurveDeployError uint64 = 0x46f6a58e

type MemepadCurveDeployError struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadCurveBuy uint64 = 0x29d1f1c1

type MemepadCurveBuy struct {
	QueryId tlb.Uint64 // uint64
	MinOut  tlb.Coins  // coins
}

const PrefixMemepadCurveSellPayload uint64 = 0xb97a94b7

type MemepadCurveSellPayload struct {
	MinOut tlb.Coins // coins
}

const PrefixMemepadCurveSellPayout uint64 = 0xdee7a887

type MemepadCurveSellPayout struct {
	QueryId tlb.Uint64 // uint64
	GramOut tlb.Coins  // coins
}

const PrefixMemepadCurveExcess uint64 = 0xca6c5778

type MemepadCurveExcess struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadCurveCreatorClaimFees uint64 = 0x6ec846e7

type MemepadCurveCreatorClaimFees struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadCurveCreatorFeesPayout uint64 = 0xaa871c43

type MemepadCurveCreatorFeesPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadCurveProtocolPayout uint64 = 0x0b8cca52

type MemepadCurveProtocolPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadCurveUpgrade uint64 = 0x24293ccc

type MemepadCurveUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}
