package tongo

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"math/big"
)

// Grams
// nanograms$_ amount:(VarUInteger 16) = Grams;
type Grams uint64 // total value fit to uint64

func (g Grams) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("grams marshaling not implemented")
}

func (g *Grams) UnmarshalTLB(c *boc.Cell, tag string) error {
	var amount struct {
		Val tlb.VarUInteger `tlb:"16bytes"`
	}
	err := tlb.Unmarshal(c, &amount)
	if err != nil {
		return err
	}
	val := big.Int(amount.Val)
	if !val.IsUint64() {
		return fmt.Errorf("grams overflow")
	}
	*g = Grams(val.Uint64())
	return nil
}

// CurrencyCollection
// currencies$_ grams:Grams other:ExtraCurrencyCollection
// = CurrencyCollection;
type CurrencyCollection struct {
	Grams Grams
	Other ExtraCurrencyCollection
}

// ExtraCurrencyCollection
// extra_currencies$_ dict:(HashmapE 32 (VarUInteger 32))
// = ExtraCurrencyCollection;
type ExtraCurrencyCollection struct {
	Dict tlb.HashmapE[struct {
		Val tlb.VarUInteger `tlb:"32bytes"`
	}] `tlb:"32bits"`
}

// HashUpdate
// update_hashes#72 {X:Type} old_hash:bits256 new_hash:bits256
// = HASH_UPDATE X;
type HashUpdate struct {
	tlb.SumType
	HashUpdate struct {
		OldHash Hash
		NewHash Hash
	} `tlbSumType:"update_hashes#72"`
}
