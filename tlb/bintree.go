package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type BinTree[T any] struct {
	Values []T
}

func decodeRecursiveBinTree(c *boc.Cell) ([]*boc.Cell, error) {
	var cellAr []*boc.Cell
	isBranch, err := c.ReadBit()
	if err != nil {
		return nil, err
	}
	if !isBranch {
		return append(cellAr, c), nil
	}

	l, err := c.NextRef()
	if err != nil {
		return nil, err
	}
	rec, err := decodeRecursiveBinTree(l)
	if err != nil {
		return nil, err
	}
	cellAr = append(cellAr, rec...)
	r, err := c.NextRef()
	if err != nil {
		return nil, err
	}
	rec, err = decodeRecursiveBinTree(r)
	if err != nil {
		return nil, err
	}
	cellAr = append(cellAr, rec...)

	return cellAr, nil
}

func (b BinTree[T]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	// TODO: implement
	return fmt.Errorf("BinTree marshaling not implmented")
}

func (b *BinTree[T]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	dec, err := decodeRecursiveBinTree(c)
	if err != nil {
		return err
	}
	for _, i := range dec {
		var t T
		err := decoder.Unmarshal(i, &t)
		if err != nil {
			return err
		}
		b.Values = append(b.Values, t)
	}
	return nil
}
