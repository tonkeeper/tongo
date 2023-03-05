package tlb

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/tonkeeper/tongo/boc"
)

// Grams
// nanograms$_ amount:(VarUInteger 16) = Grams;
type Grams uint64 // total value fit to uint64

func (g Grams) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	var amount VarUInteger16
	amount = VarUInteger16(*big.NewInt(int64(g)))
	err := Marshal(c, amount)
	return err
}

func (g *Grams) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	ln, err := c.ReadLimUint(15)
	if err != nil {
		return err
	}
	if ln > 8 {
		return fmt.Errorf("grams overflow")
	}
	var amount uint64
	for i := 0; i < int(ln); i++ {
		b, err := c.ReadUint(8)
		if err != nil {
			return err
		}
		amount = uint64(b) | (amount << 8)
	}
	*g = Grams(amount)
	return nil
}

func (g Grams) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", g)), nil
}

func (g *Grams) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseUint(string(bytes.Trim(data, "\" \n")), 10, 64)
	if err != nil {
		return err
	}
	*g = Grams(val)
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
	Dict HashmapE[Uint32, VarUInteger32]
}

// HashUpdate
// update_hashes#72 {X:Type} old_hash:bits256 new_hash:bits256
// = HASH_UPDATE X;
type HashUpdate struct {
	Magic   Magic `tlb:"update_hashes#72"`
	OldHash Bits256
	NewHash Bits256
}

// SnakeData
// tail#_ {bn:#} b:(bits bn) = SnakeData ~0;
// cons#_ {bn:#} {n:#} b:(bits bn) next:^(SnakeData ~n) = SnakeData ~(n + 1);
type SnakeData boc.BitString

func (s SnakeData) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	bs := boc.BitString(s)
	if c.BitsAvailableForWrite() < bs.GetWriteCursor() {
		s, err := bs.ReadBits(c.BitsAvailableForWrite())
		if err != nil {
			return err
		}
		err = c.WriteBitString(s)
		if err != nil {
			return err
		}
		ref := boc.NewCell()
		err = Marshal(ref, SnakeData(bs.ReadRemainingBits()))
		if err != nil {
			return err
		}
		err = c.AddRef(ref)
		return err
	}
	return c.WriteBitString(bs)
}

func (s *SnakeData) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	b := c.ReadRemainingBits()
	if c.RefsAvailableForRead() > 0 {
		cell, err := c.NextRef()
		if err != nil {
			return err
		}
		var sn SnakeData
		err = decoder.Unmarshal(cell, &sn)
		if err != nil {
			return err
		}
		b.Append(boc.BitString(sn))
	}
	*s = SnakeData(b)
	return nil
}

// text#_ {n:#} data:(SnakeData ~n) = Text;
type Text string

func (t Text) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	bs := boc.NewBitString(len(t) * 8)
	err := bs.WriteBytes([]byte(t))
	if err != nil {
		return err
	}
	return Marshal(c, SnakeData(bs))
}

func (t *Text) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var sn SnakeData
	err := decoder.Unmarshal(c, &sn)
	if err != nil {
		return err
	}
	bs := boc.BitString(sn)
	if bs.BitsAvailableForRead()%8 != 0 {
		return fmt.Errorf("text data must be a multiple of 8 bits")
	}
	b, err := bs.GetTopUppedArray()
	if err != nil {
		return err
	}
	*t = Text(b)
	return nil
}

type FixedLengthText string

func (t FixedLengthText) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	l := len(t)
	err := c.WriteUint(uint64(l), 8)
	if err != nil {
		return err
	}
	return c.WriteBytes([]byte(t))
}

func (t *FixedLengthText) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	l, err := c.ReadUint(8)
	if err != nil {
		return err
	}
	b, err := c.ReadBytes(int(l))
	*t = FixedLengthText(b)
	return err
}

// FullContent
// onchain#00 data:(HashMapE 256 ^ContentData) = FullContent;
// offchain#01 uri:Text = FullContent;
// text#_ {n:#} data:(SnakeData ~n) = Text;
type FullContent struct {
	SumType
	Onchain struct {
		Data HashmapE[Bits256, Ref[ContentData]]
	} `tlbSumType:"onchain#00"`
	Offchain struct {
		Uri SnakeData // Text
	} `tlbSumType:"offchain#01"`
}

// ContentData
// snake#00 data:(SnakeData ~n) = ContentData;
// chunks#01 data:ChunkedData = ContentData;
type ContentData struct {
	SumType
	Snake struct {
		Data SnakeData
	} `tlbSumType:"snake#00"`
	Chunks struct {
		Data ChunkedData
	} `tlbSumType:"chunks#01"`
}

func (c ContentData) Bytes() ([]byte, error) {
	var bs boc.BitString
	switch c.SumType {
	case "Snake":
		bs = boc.BitString(c.Snake.Data)
	case "Chunks":
		bs = boc.BitString(c.Chunks.Data)
	default:
		return nil, fmt.Errorf("empty content data struct")
	}
	if bs.BitsAvailableForRead()%8 != 0 {
		return nil, fmt.Errorf("data is not multiple of 8 bits")
	}
	return bs.GetTopUppedArray()
}

// ChunkedData
// chunked_data#_ data:(HashMapE 32 ^(SnakeData ~0)) = ChunkedData;
type ChunkedData boc.BitString

func (d ChunkedData) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	// TODO: implement
	return fmt.Errorf("ChunkedData marshaling not implemented")
}

func (d *ChunkedData) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	type chunkedData struct {
		Data HashmapE[Uint32, Ref[SnakeData]]
	}
	var (
		cd chunkedData
	)
	b := boc.NewBitString(boc.CellBits)
	err := decoder.Unmarshal(c, &cd)
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

type ShardDesc struct {
	SumType
	Old struct {
		SeqNo              uint32
		RegMcSeqno         uint32
		StartLT            uint64
		EndLT              uint64
		RootHash           Bits256
		FileHash           Bits256
		BeforeSplit        bool
		BeforeMerge        bool
		WantSplit          bool
		WantMerge          bool
		NXCCUpdated        bool
		Flags              Uint3
		NextCatchainSeqNo  uint32
		NextValidatorShard int64
		MinRefMcSeqNo      uint32
		GenUTime           uint32
	} `tlbSumType:"old#b"`
	New struct {
		SeqNo              uint32
		RegMcSeqno         uint32
		StartLT            uint64
		EndLT              uint64
		RootHash           Bits256
		FileHash           Bits256
		BeforeSplit        bool
		BeforeMerge        bool
		WantSplit          bool
		WantMerge          bool
		NXCCUpdated        bool
		Flags              Uint3
		NextCatchainSeqNo  uint32
		NextValidatorShard int64
		MinRefMcSeqNo      uint32
		GenUTime           uint32
	} `tlbSumType:"new#a"`
}

type ShardInfoBinTree struct {
	BinTree BinTree[ShardDesc]
}

type AllShardsInfo struct {
	ShardHashes HashmapE[Uint32, Ref[ShardInfoBinTree]]
}
