package tongo

import (
	"encoding/json"
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

func (g Grams) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d", g))
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

// SnakeData
// tail#_ {bn:#} b:(bits bn) = SnakeData ~0;
// cons#_ {bn:#} {n:#} b:(bits bn) next:^(SnakeData ~n) = SnakeData ~(n + 1);
type SnakeData boc.BitString

func (s SnakeData) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("SnakeData marshaling not implemented")
}

func (s *SnakeData) UnmarshalTLB(c *boc.Cell, tag string) error {
	b := c.ReadRemainingBits()
	if c.RefsAvailableForRead() > 0 {
		cell, err := c.NextRef()
		if err != nil {
			return err
		}
		var sn SnakeData
		err = tlb.Unmarshal(cell, &sn)
		if err != nil {
			return err
		}
		b.Append(boc.BitString(sn))
	}
	*s = SnakeData(b)
	return nil
}

// FullContent
// onchain#00 data:(HashMapE 256 ^ContentData) = FullContent;
// offchain#01 uri:Text = FullContent;
// text#_ {n:#} data:(SnakeData ~n) = Text;
type FullContent struct {
	tlb.SumType
	Onchain struct {
		Data tlb.HashmapE[tlb.Ref[ContentData]] `tlb:"256bits"`
	} `tlbSumType:"onchain#00"`
	Offchain struct {
		Uri SnakeData // Text
	} `tlbSumType:"offchain#01"`
}

// ContentData
// snake#00 data:(SnakeData ~n) = ContentData;
// chunks#01 data:ChunkedData = ContentData;
type ContentData struct {
	tlb.SumType
	Snake struct {
		Data SnakeData
	} `tlbSumType:"snake#00"`
	Chunks struct {
		Data ChunkedData
	} `tlbSumType:"chunks#01"`
}

// ChunkedData
// chunked_data#_ data:(HashMapE 32 ^(SnakeData ~0)) = ChunkedData;
type ChunkedData boc.BitString

func (d ChunkedData) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("ChunkedData marshaling not implemented")
}

func (d *ChunkedData) UnmarshalTLB(c *boc.Cell, tag string) error {
	type chunkedData struct {
		Data tlb.HashmapE[tlb.Ref[SnakeData]] `tlb:"32bits"`
	}
	var (
		cd chunkedData
	)
	b := boc.NewBitString(boc.CellBits)
	err := tlb.Unmarshal(c, &cd)
	if err != nil {
		return err
	}
	// TODO: check keys sort
	for _, x := range cd.Data.Values() {
		b.Append(boc.BitString(x.Value))
	}
	*d = ChunkedData(b)
	return nil
}
