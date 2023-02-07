package tongo

import "github.com/tonkeeper/tongo/tlb"

type Transaction struct {
	tlb.Transaction
	BlockID BlockIDExt
}
