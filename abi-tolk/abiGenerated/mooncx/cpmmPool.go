// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixCpmmPoolDeploy uint64 = 0x6ca4cb66

type CpmmPoolDeploy struct {
	QueryId        tlb.Uint64 // uint64
	LpFee          tlb.Uint8  // uint8
	ProtocolFee    tlb.Uint8  // uint8
	LpAccCode      boc.Cell   // cell
	LpPositionCode boc.Cell   // cell
}

const PrefixCpmmPoolProvideReady uint64 = 0x44daf70f

type CpmmPoolProvideReady struct {
	QueryId tlb.Uint64 // uint64
	Attempt tlb.Uint8  // uint8
}

const PrefixCpmmPoolTakeReady uint64 = 0x67071f0d

type CpmmPoolTakeReady struct {
	QueryId tlb.Uint64 // uint64
	Attempt tlb.Uint8  // uint8
	IsReady bool       // bool
}

type CpmmTonSwapParams struct {
	SwapAmount  tlb.Coins           // coins
	Excess      tlb.InternalAddress // address
	MinOut      tlb.Coins           // coins
	NextFulfill tlb.Maybe[boc.Cell] // cell?
	NextReject  tlb.Maybe[boc.Cell] // cell?
}

const PrefixCpmmPoolSwap uint64 = 0xfa45056b

type CpmmPoolSwap struct {
	QueryId    tlb.Uint64        // uint64
	SwapParams CpmmTonSwapParams // CpmmTonSwapParams
}

const PrefixSwapFulfill uint64 = 0x9d142643

type SwapFulfill struct {
	QueryId tlb.Uint64          // uint64
	Amount  tlb.Coins           // coins
	Payload tlb.Maybe[boc.Cell] // cell?
}

const PrefixSwapReject uint64 = 0xff32da99

type SwapReject struct {
	QueryId tlb.Uint64          // uint64
	Amount  tlb.Coins           // coins
	Payload tlb.Maybe[boc.Cell] // cell?
}

const PrefixCpmmPoolDepositLiquidity uint64 = 0x90af0b63

type CpmmPoolDepositLiquidity struct {
	QueryId  tlb.Uint64                    // uint64
	MinLpOut tlb.Coins                     // coins
	Locked   bool                          // bool
	Split    tlb.Maybe[tlb.RefT[*LpSplit]] // Cell<LpSplit>?
}

const PrefixCpmmPoolProvideLiquidity uint64 = 0x48a936dd

type CpmmPoolProvideLiquidity struct {
	QueryId  tlb.Uint64                    // uint64
	Provider tlb.InternalAddress           // address
	XAmount  tlb.Coins                     // coins
	YAmount  tlb.Coins                     // coins
	MinLpOut tlb.Coins                     // coins
	Locked   bool                          // bool
	Split    tlb.Maybe[tlb.RefT[*LpSplit]] // Cell<LpSplit>?
}

const PrefixCpmmPoolWithdrawLiquidity uint64 = 0xcbd7f377

type CpmmPoolWithdrawLiquidity struct {
	QueryId  tlb.Uint64          // uint64
	Provider tlb.InternalAddress // address
	LpPos    tlb.Coins           // coins
	XFeesPos tlb.Uint256         // uint256
	YFeesPos tlb.Uint256         // uint256
}

const PrefixCpmmPoolWithdrawFees uint64 = 0xd5c9f8c3

type CpmmPoolWithdrawFees struct {
	QueryId  tlb.Uint64          // uint64
	Provider tlb.InternalAddress // address
	LpPos    tlb.Coins           // coins
	XFeesPos tlb.Uint256         // uint256
	YFeesPos tlb.Uint256         // uint256
}

const PrefixCpmmPoolCancelDeposit uint64 = 0x1f911dd9

type CpmmPoolCancelDeposit struct {
	QueryId  tlb.Uint64          // uint64
	Provider tlb.InternalAddress // address
	XAmount  tlb.Coins           // coins
	YAmount  tlb.Coins           // coins
}

const PrefixCpmmPoolBoost uint64 = 0x50bbec04

type CpmmPoolBoost struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolCollectProtocolFees uint64 = 0xa7335450

type CpmmPoolCollectProtocolFees struct {
	QueryId   tlb.Uint64          // uint64
	Recipient tlb.InternalAddress // address
}

const PrefixCpmmPoolUpdateWallets uint64 = 0xc29e72ba

type CpmmPoolUpdateWallets struct {
	QueryId    tlb.Uint64          // uint64
	NewXWallet tlb.MsgAddress      // address?
	NewYWallet tlb.InternalAddress // address
}

const PrefixCpmmPoolUpgrade uint64 = 0x38599428

type CpmmPoolUpgrade struct {
	QueryId tlb.Uint64          // uint64
	NewCode tlb.Maybe[boc.Cell] // cell?
	NewData tlb.Maybe[boc.Cell] // cell?
}

type CpmmJettonSwapParams struct {
	Excess      tlb.InternalAddress // address
	MinOut      tlb.Coins           // coins
	NextFulfill tlb.Maybe[boc.Cell] // cell?
	NextReject  tlb.Maybe[boc.Cell] // cell?
}

const PrefixCpmmPoolSwapPayload uint64 = 0xfa45056b

type CpmmPoolSwapPayload struct {
	SwapParams CpmmJettonSwapParams // CpmmJettonSwapParams
}

const PrefixSwapFulfillPayload uint64 = 0x9d142643

type SwapFulfillPayload struct {
	Payload tlb.Maybe[boc.Cell] // cell?
}

const PrefixSwapRejectPayload uint64 = 0xff32da99

type SwapRejectPayload struct {
	Payload tlb.Maybe[boc.Cell] // cell?
}

const PrefixCpmmPoolDepositLiquidityPayload uint64 = 0x90af0b63

type CpmmPoolDepositLiquidityPayload struct {
	MinLpOut tlb.Coins                     // coins
	Locked   bool                          // bool
	Split    tlb.Maybe[tlb.RefT[*LpSplit]] // Cell<LpSplit>?
}

const PrefixCpmmPoolBoostPayload uint64 = 0x50bbec04

type CpmmPoolBoostPayload struct {
}

const PrefixCpmmPoolProvideLiquidityRefund uint64 = 0x67568a3a

type CpmmPoolProvideLiquidityRefund struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolProvideLiquidityExcess uint64 = 0x9d9f7bad

type CpmmPoolProvideLiquidityExcess struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolWithdrawLiquidityPayout uint64 = 0x15b7ef52

type CpmmPoolWithdrawLiquidityPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolWithdrawFeesPayout uint64 = 0x33e81d59

type CpmmPoolWithdrawFeesPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolCancelDepositRefund uint64 = 0xd212503b

type CpmmPoolCancelDepositRefund struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmPoolProtocolFeesPayout uint64 = 0x15e67eab

type CpmmPoolProtocolFeesPayout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixCpmmExcess uint64 = 0x13925ff2

type CpmmExcess struct {
	QueryId tlb.Uint64 // uint64
}
