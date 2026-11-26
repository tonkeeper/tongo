// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixClmmFactoryCreatePool uint64 = 0xede56ae4

type ClmmFactoryCreatePool struct {
	QueryId tlb.Uint64     // uint64
	XMaster tlb.MsgAddress // address?
	YMaster tlb.MsgAddress // address?
	BinStep tlb.Uint16     // uint16
	InitBin tlb.Int64      // int64
}

const PrefixClmmFactoryUpgradePool uint64 = 0x4c76e24e

type ClmmFactoryUpgradePool struct {
	QueryId tlb.Uint64          // uint64
	Pool    tlb.InternalAddress // address
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}

const PrefixClmmFactoryUpdateRangeWallets uint64 = 0xe55e6dc7

type ClmmFactoryUpdateRangeWallets struct {
	QueryId  tlb.Uint64          // uint64
	Pool     tlb.InternalAddress // address
	RangeNum tlb.Int16           // int16
	XWallet  tlb.MsgAddress      // address?
	YWallet  tlb.InternalAddress // address
}

const PrefixClmmFactoryUpdateRangeFees uint64 = 0x40073294

type ClmmFactoryUpdateRangeFees struct {
	QueryId     tlb.Uint64          // uint64
	Pool        tlb.InternalAddress // address
	RangeNum    tlb.Int16           // int16
	LpFee       tlb.Uint8           // uint8
	ProtocolFee tlb.Uint8           // uint8
}

const PrefixClmmFactoryUpgradeRange uint64 = 0x06dbf780

type ClmmFactoryUpgradeRange struct {
	QueryId  tlb.Uint64          // uint64
	Pool     tlb.InternalAddress // address
	RangeNum tlb.Int16           // int16
	NewCode  tlb.Maybe[boc.Cell] // cell?
	NewData  tlb.Maybe[boc.Cell] // cell?
}

const PrefixClmmFactoryCollectProtocolFees uint64 = 0x4b3332cb

type ClmmFactoryCollectProtocolFees struct {
	QueryId   tlb.Uint64          // uint64
	Pool      tlb.InternalAddress // address
	RangeNum  tlb.Int16           // int16
	Recipient tlb.InternalAddress // address
}

const PrefixClmmFactoryUpgrade uint64 = 0x496938f7

type ClmmFactoryUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}
