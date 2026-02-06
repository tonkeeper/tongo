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
	if len(b.Values) == 0 {
		return fmt.Errorf("BinTree must have at least one value")
	}
	return b.marshalRecursive(c, b.Values, encoder)
}

func (b BinTree[T]) marshalRecursive(c *boc.Cell, values []T, encoder *Encoder) error {
	if len(values) == 1 {
		if err := c.WriteBit(false); err != nil {
			return err
		}
		return encoder.Marshal(c, values[0])
	}
	if err := c.WriteBit(true); err != nil {
		return err
	}
	mid := len(values) / 2
	l, err := c.NewRef()
	if err != nil {
		return err
	}
	if err := b.marshalRecursive(l, values[:mid], encoder); err != nil {
		return err
	}
	r, err := c.NewRef()
	if err != nil {
		return err
	}
	return b.marshalRecursive(r, values[mid:], encoder)
}

func (b *BinTree[T]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	dec, err := decodeRecursiveBinTree(c)
	if err != nil {
		return err
	}
	b.Values = make([]T, 0, len(dec))
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
