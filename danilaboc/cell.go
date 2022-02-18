package danilaboc

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

type Cell struct {
	contents  []byte
	refs      []Cell
	exotic    bool
	bitsCount int
}

func NewCell() Cell {
	return Cell{
		exotic:   false,
		contents: make([]byte, 0),
		refs:     make([]Cell, 0),
	}
}

func (c *Cell) BitSize() int {
	return c.bitsCount
}

func (c *Cell) RefsSize() int {
	return len(c.refs)
}

func (c *Cell) AppendBytes(data []byte, bits int) error {
	if bits == 0 {
		return nil
	}

	if c.BitSize()+bits >= 1023 {
		return errors.New("cell overflow error")
	}

	shift := c.BitSize() % 8
	if shift == 0 {
		c.contents = append(c.contents, 0)
	}
	appendedBits := 0
	c.contents[len(c.contents)-1] |= data[0] >> shift
	appendedBits += 8 - shift
	for i := 1; i < len(data); i++ {
		b := data[i-1] << (8 - shift)
		b |= data[i] >> shift

		c.contents = append(c.contents, b)
		appendedBits += 8
	}
	if shift != 0 && appendedBits < bits {
		c.contents = append(c.contents, data[len(data)-1]<<(8-shift))
	}

	c.bitsCount += bits

	return nil
}

func (c *Cell) AppendBits(data byte, bits int) error {
	return c.AppendBytes([]byte{data << (8 - bits)}, bits)
}

func (c *Cell) Append(data []byte) error {
	return c.AppendBytes(data, len(data)*8)
}

func (c *Cell) AppendGrams(nanograms uint64) error {
	for b := 0; b < 16; b++ {
		if nanograms < 1<<(b*8) {
			bytes := make([]byte, b)
			for j := 0; j < b; j++ {
				bytes[b-j-1] = byte(nanograms >> (j * 8))
			}
			err := c.AppendBits(byte(b), 4)
			if err != nil {
				return err
			}
			return c.Append(bytes)
		}
	}
	return errors.New("too many nanograms")
}

func (c *Cell) AppendCell(c2 Cell) error {
	err := c.AppendBytes(c2.contents, c2.BitSize())
	if err != nil {
		return err
	}
	for _, r := range c2.refs {
		err := c.AddReference(r)
		if err != nil {
			return err
		}
	}
	return nil
}

// AppendWallet appends workchain id and the wallet bytes
// from the specified url-safe base64 wallet.
func (c *Cell) AppendWallet(wallet string) error {
	wb, err := base64.URLEncoding.DecodeString(wallet)
	if err != nil {
		return err
	}

	return c.Append(wb[1:34])
}

func (c *Cell) AddReference(c2 Cell) error {
	if c.RefsSize() == 4 {
		return errors.New("Cell references are filled")
	}

	c.refs = append(c.refs, c2)

	return nil
}

func (c *Cell) getAllTree(cells *[]Cell) *[]Cell {
	*cells = append(*cells, *c)
	for i := 0; i < c.RefsSize(); i++ {
		*cells = *c.refs[i].getAllTree(cells)
	}
	return cells
}

func (c *Cell) GetAllTree() []Cell {
	cells := make([]Cell, 0)
	return *c.getAllTree(&cells)
}

func (c *Cell) GetMaxLevel() int {
	return 0
}

func (c *Cell) GetDepth() uint16 {
	var maxDepth uint16 = 0
	if c.RefsSize() > 0 {
		for _, r := range c.refs {
			rd := r.GetDepth()
			if rd > maxDepth {
				maxDepth = rd
			}
		}
		maxDepth += 1
	}
	return maxDepth
}

func (c *Cell) HashRepr() []byte {
	data := c.bocReprWithoutRefs()
	for _, r := range c.refs {
		depth := r.GetDepth()
		depthRepr := make([]byte, 2)
		binary.BigEndian.PutUint16(depthRepr, depth)
		data = append(data, depthRepr...)
	}
	for _, r := range c.refs {
		data = append(data, r.HashRepr()...)
	}
	sum := sha256.Sum256(data)
	return sum[:]
}
