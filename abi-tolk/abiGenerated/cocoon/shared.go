// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixClientProxyRequest uint64 = 0x65448ff4

type ClientProxyRequest struct {
	QueryId      tlb.Uint64                 // uint64
	OwnerAddress tlb.InternalAddress        // address
	StateData    tlb.RefT[*ClientStateData] // Cell<ClientStateData>
	Payload      tlb.Maybe[boc.Cell]        // cell?
}
type ClientStateData struct {
	State      tlb.Uint2   // uint2
	Balance    tlb.Coins   // coins
	Stake      tlb.Coins   // coins
	TokensUsed tlb.Uint64  // uint64
	SecretHash tlb.Uint256 // uint256
}
type CocoonParams struct {
	StructVersion                   tlb.Uint8           // uint8
	ParamsVersion                   tlb.Uint32          // uint32
	UniqueId                        tlb.Uint32          // uint32
	IsTest                          bool                // bool
	PricePerToken                   tlb.Coins           // coins
	WorkerFeePerToken               tlb.Coins           // coins
	PromptTokensPriceMultiplier     tlb.Uint32          // uint32
	CachedTokensPriceMultiplier     tlb.Uint32          // uint32
	CompletionTokensPriceMultiplier tlb.Uint32          // uint32
	ReasoningTokensPriceMultiplier  tlb.Uint32          // uint32
	ProxyDelayBeforeClose           tlb.Uint32          // uint32
	ClientDelayBeforeClose          tlb.Uint32          // uint32
	MinProxyStake                   tlb.Coins           // coins
	MinClientStake                  tlb.Coins           // coins
	ProxyScCode                     tlb.Maybe[boc.Cell] // cell?
	WorkerScCode                    tlb.Maybe[boc.Cell] // cell?
	ClientScCode                    tlb.Maybe[boc.Cell] // cell?
}

const PrefixPayout uint64 = 0xc59a7cd3

type Payout struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixReturnExcessesBack uint64 = 0x2565934c

type ReturnExcessesBack struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixWorkerProxyRequest uint64 = 0x4d725d2c

type WorkerProxyRequest struct {
	QueryId      tlb.Uint64          // uint64
	OwnerAddress tlb.InternalAddress // address
	State        tlb.Uint2           // uint2
	Tokens       tlb.Uint64          // uint64
	Payload      tlb.Maybe[boc.Cell] // cell?
}
