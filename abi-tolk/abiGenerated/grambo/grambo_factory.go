// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const PrefixGramboDeploy uint64 = 0x0d4295e0

type GramboDeploy struct {
	QueryId      tlb.Uint64 // uint64
	Metadata     boc.Cell   // cell
	TonTarget    tlb.Coins  // coins
	PrebuyAmount tlb.Coins  // coins
}

const PrefixGramboUpdateData uint64 = 0x20b031d9

type GramboUpdateData struct {
	NewData boc.Cell // cell
}

const PrefixGramboUpdateCode uint64 = 0x61d196c1

type GramboUpdateCode struct {
	NewCode boc.Cell // cell
}

type FactoryIncomingMessageKind uint

const (
	FactoryIncomingMessageKind_GramboDeploy     FactoryIncomingMessageKind = 222467552
	FactoryIncomingMessageKind_GramboUpdateData FactoryIncomingMessageKind = 548418009
	FactoryIncomingMessageKind_GramboUpdateCode FactoryIncomingMessageKind = 1641125569
)

type FactoryIncomingMessage struct {
	SumType          FactoryIncomingMessageKind
	GramboDeploy     *GramboDeploy
	GramboUpdateData *GramboUpdateData
	GramboUpdateCode *GramboUpdateCode
}
