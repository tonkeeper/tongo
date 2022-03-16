package boc

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
)

const BOCSizeLimit = 10000

type Cell struct {
	Bits     BitString
	isExotic bool
	refs     []*Cell
}

func NewCell() *Cell {
	return &Cell{
		Bits:     NewBitString(1023),
		refs:     make([]*Cell, 4),
		isExotic: false,
	}
}

func NewCellExotic() *Cell {
	return &Cell{
		Bits:     NewBitString(1023),
		refs:     make([]*Cell, 4),
		isExotic: true,
	}
}

func (c *Cell) BeginParse() BitStringReader {
	return NewBitStringReader(&c.Bits)
}

func (c *Cell) RefsSize() int {
	return len(c.Refs())
}

func (c *Cell) Refs() []*Cell {
	res := make([]*Cell, 0)
	for _, ref := range c.refs {
		if ref != nil {
			res = append(res, ref)
		}
	}
	return res
	//return c.refs
}

func (c *Cell) IsExotic() bool {
	return c.isExotic
}

func (c *Cell) BitSize() int {
	return c.Bits.Cursor()
}

func (c *Cell) Hash() ([]byte, error) {
	return hashCell(c)
}

func (c *Cell) HashString() (string, error) {
	h, err := hashCell(c) //todo: check error
	return hex.EncodeToString(h), err
}

func (c *Cell) ToBoc() ([]byte, error) {
	return SerializeBoc(c, true, true, false, 0)
}

func (c *Cell) ToBocString() (string, error) {
	return c.ToBocStringCustom(true, true, false, 0)
}

func (c *Cell) ToBocBase64() (string, error) {
	return c.ToBocBase64Custom(true, true, false, 0)
}

func (c *Cell) ToBocCustom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	return SerializeBoc(c, idx, hasCrc32, cacheBits, flags)
}

func (c *Cell) ToBocStringCustom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) (string, error) {
	boc, err := c.ToBocCustom(idx, hasCrc32, cacheBits, flags)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(boc), nil
}

func (c *Cell) ToBocBase64Custom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) (string, error) {
	boc, err := c.ToBocCustom(idx, hasCrc32, cacheBits, flags)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(boc), nil
}

func (c *Cell) AddReference(c2 *Cell) (*Cell, error) {
	if c.RefsSize() == 4 {
		return c, errors.New("cell references are filled")
	}

	c.refs = append(c.refs, c2)

	return c, nil
}

func (c *Cell) toStringImpl(ident string, iterationsLimit *int) string {
	s := ident + "x{" + c.Bits.ToFiftHex() + "}\n"
	if *iterationsLimit == 0 {
		return s
	}
	*iterationsLimit -= 1
	for _, ref := range c.Refs() {
		s += ref.toStringImpl(ident+" ", iterationsLimit)
	}
	return s
}

func (c *Cell) ToString() string {
	iter := BOCSizeLimit
	return c.toStringImpl("", &iter)
}
