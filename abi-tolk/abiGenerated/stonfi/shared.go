// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type ForwardParams struct {
	Value          tlb.Coins           // coins
	SuccessPayload tlb.Maybe[boc.Cell] // cell?
	RejectPayload  tlb.Maybe[boc.Cell] // cell?
}

type LockForwardParams struct {
	Dest          tlb.MsgAddress // address?
	ForwardParams ForwardParams  // ForwardParams
}

type DutchSegments struct {
	Num      tlb.Uint4 // uint4
	Segments boc.Cell  // RemainingBitsAndRefs
}

type BilateralLockArgs struct {
	Resolver             tlb.MsgAddress           // address?
	ResolverTimeoutDelta tlb.Uint64               // uint64
	ResolverAskAmount    tlb.Coins                // coins
	ResolverPublicKey    tlb.Maybe[tlb.Uint256]   // uint256?
	DutchSegments        tlb.RefT[*DutchSegments] // Cell<DutchSegments>
}

type BilateralSignedMessage struct {
	ExpirationTime tlb.Uint64     // uint64
	Resolver       tlb.MsgAddress // address?
}

type BilateralSignPayload struct {
	Signature tlb.Bits512                       // bits512
	Message   tlb.RefT[*BilateralSignedMessage] // Cell<BilateralSignedMessage>
}

type BilateralUnlockArgs struct {
	MinOut             tlb.Coins                       // coins
	IgnoreRefundAmount tlb.Coins                       // coins
	Signed             tlb.Maybe[BilateralSignPayload] // BilateralSignPayload?
}

type VaultLockAdditionalDataMore struct {
	OrderOwner      tlb.MsgAddress // address?
	RefundTo        tlb.MsgAddress // address?
	AskJettonMinter tlb.MsgAddress // address?
}

const ErrorInvalidCaller = 0x80B6 // 32950

const ErrorInsufficientGas = 0x831D // 33565

const ErrorInvalidAmount = 0x8753 // 34643
