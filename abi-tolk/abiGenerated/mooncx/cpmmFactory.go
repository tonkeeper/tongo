// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixCpmmFactoryCreatePool uint64 = 0xec725c4f

type CpmmFactoryCreatePool struct {
	QueryId tlb.Uint64     // uint64
	XMaster tlb.MsgAddress // address?
	YMaster tlb.MsgAddress // address?
	FeeBps  tlb.Uint8      // uint8
}

const PrefixCpmmFactoryProvidePoolAddress uint64 = 0x1d42af37

type CpmmFactoryProvidePoolAddress struct {
	QueryId   tlb.Uint64          // uint64
	XMaster   tlb.MsgAddress      // address?
	YMaster   tlb.InternalAddress // address
	FeeBps    tlb.Uint8           // uint8
	Initiator tlb.InternalAddress // address
}

const PrefixCpmmFactoryTakePoolAddress uint64 = 0x81a58422

type CpmmFactoryTakePoolAddress struct {
	QueryId tlb.Uint64          // uint64
	Pool    tlb.InternalAddress // address
	XMaster tlb.MsgAddress      // address?
	YMaster tlb.InternalAddress // address
	FeeBps  tlb.Uint8           // uint8
}

const PrefixCpmmFactoryCollectProtocolFees uint64 = 0xc492cce9

type CpmmFactoryCollectProtocolFees struct {
	QueryId   tlb.Uint64          // uint64
	Pool      tlb.InternalAddress // address
	Recipient tlb.InternalAddress // address
}

const PrefixCpmmFactoryUpdatePoolWallets uint64 = 0x51d80e5e

type CpmmFactoryUpdatePoolWallets struct {
	QueryId tlb.Uint64          // uint64
	Pool    tlb.InternalAddress // address
	XWallet tlb.MsgAddress      // address?
	YWallet tlb.InternalAddress // address
}

const PrefixCpmmFactoryUpgradePool uint64 = 0x0d730f18

type CpmmFactoryUpgradePool struct {
	QueryId tlb.Uint64          // uint64
	Pool    tlb.InternalAddress // address
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}

const PrefixCpmmFactoryUpgrade uint64 = 0xbae2139a

type CpmmFactoryUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}
