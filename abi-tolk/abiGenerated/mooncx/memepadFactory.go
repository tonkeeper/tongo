// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixMemepadFactoryCreateCurve uint64 = 0x7e838abb

type MemepadFactoryCreateCurve struct {
	QueryId    tlb.Uint64 // uint64
	Salt       tlb.Uint64 // uint64
	FeeBps     tlb.Uint8  // uint8
	InitialBuy tlb.Coins  // coins
	Content    boc.Cell   // cell
}

const PrefixMemepadFactoryWithdraw uint64 = 0x329068a6

type MemepadFactoryWithdraw struct {
	QueryId   tlb.Uint64          // uint64
	Amount    tlb.Coins           // coins
	Recipient tlb.InternalAddress // address
}

const PrefixMemepadFactoryWithdrawPayout uint64 = 0xfce532d8

type MemepadFactoryWithdrawPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixMemepadFactoryUpgrade uint64 = 0x0eb678fe

type MemepadFactoryUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}

const PrefixMemepadFactoryUpgradeCurve uint64 = 0x9057f56b

type MemepadFactoryUpgradeCurve struct {
	QueryId tlb.Uint64          // uint64
	Curve   tlb.InternalAddress // address
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}
