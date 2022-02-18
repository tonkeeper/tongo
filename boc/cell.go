package boc

import "errors"

type Cell struct {
	Bits     BitString
	isExotic bool
	refs     []*Cell
}

func NewCell(isExotic bool) Cell {
	return Cell{
		Bits:     NewBitString(1023),
		refs:     make([]*Cell, 0),
		isExotic: isExotic,
	}
}

func (c *Cell) BeginParse() BitStringReader {
	return NewBitStringReader(&c.Bits)
}

func (c *Cell) RefsSize() int {
	return len(c.refs)
}

func (c *Cell) AddReference(c2 *Cell) (error, *Cell) {
	if c.RefsSize() == 4 {
		return errors.New("cell references are filled"), c
	}

	c.refs = append(c.refs, c2)

	return nil, c
}
