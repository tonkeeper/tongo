package boc

import (
	"fmt"
	"math"
	"math/big"
)

type CellReader struct {
	Cell
	bitCursor int
	refCursor int
}

func NewCellReader(cell Cell) CellReader {
	return CellReader{Cell: cell, bitCursor: 0, refCursor: 0}
}

func (c *CellReader) bitsAvailable() int {
	return c.Bits.len - c.bitCursor
}

func (c *CellReader) readBit() bool {
	var bit = c.Bits.Get(c.bitCursor)
	c.bitCursor++
	return bit
}

func (c *CellReader) GetRef() (CellReader, error) {
	c.refCursor += 1
	if c.refCursor > len(c.refs) {
		return CellReader{}, fmt.Errorf("not enought refs")
	}
	return NewCellReader(*c.refs[c.refCursor-1]), nil
}

func (c *CellReader) ReadBit() (bool, error) {
	if c.bitsAvailable() < 1 {
		return false, ErrNotEnoughBits
	}
	return c.readBit(), nil
}

func (c *CellReader) ReadUint(bitLen int) (uint64, error) {
	if bitLen > 64 {
		return 0, fmt.Errorf("len more than 64 bits, need to use bigint")
	}
	if c.bitsAvailable() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	var res uint64 = 0
	for i := bitLen - 1; i >= 0; i-- {
		if c.readBit() {
			res |= 1 << i
		}
	}
	return res, nil
}

func (c *CellReader) ReadInt(bitLen int) (int64, error) {
	if bitLen > 64 {
		return 0, fmt.Errorf("len more than 64 bits, need to use bigint")
	}
	if c.bitsAvailable() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	if bitLen == 1 {
		if c.readBit() {
			return -1, nil
		}
		return 0, nil
	}
	if c.readBit() {
		base, err := c.ReadUint(bitLen - 1)
		if err != nil {
			return 0, err
		}
		return int64(base - uint64(math.Pow(2, float64(bitLen-1)))), nil
	}
	res, err := c.ReadUint(bitLen - 1)
	if err != nil {
		return 0, err
	}
	return int64(res), nil
}

func (c *CellReader) ReadBigUint(bitLen int) (*big.Int, error) {
	if c.bitsAvailable() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	var res = ""
	for i := 0; i < bitLen; i++ {
		if c.readBit() {
			res += "1"
		} else {
			res += "0"
		}
	}
	var num = big.NewInt(0)
	num.SetString(res, 2)
	return num, nil
}

func (c *CellReader) ReadBigInt(bitLen int) (*big.Int, error) {
	if c.bitsAvailable() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	if bitLen == 1 {
		if c.readBit() {
			return big.NewInt(-1), nil
		}
		return big.NewInt(0), nil
	}
	if c.readBit() {
		var base, _ = c.ReadBigUint(bitLen - 1)
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return base.Sub(base, nb), nil
	}
	return c.ReadBigUint(bitLen - 1)
}

func (c *CellReader) ReadBytes(size int) ([]byte, error) {
	if c.bitsAvailable() < size*8 {
		return nil, ErrNotEnoughBits
	}
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		b, err := c.ReadUint(8)
		if err != nil {
			return nil, err
		}
		res[i] = byte(b)
	}
	return res, nil
}

// ReadVarUint
// TL-B: var_uint$_ {n:#} len:(#< n) value:(uint (len * 8)) = VarUInteger n;
func (c *CellReader) ReadVarUint(byteLen int) (*big.Int, error) {
	if byteLen < 2 {
		return nil, fmt.Errorf("invalid varuint length")
	}
	lenBits := int(math.Ceil(math.Log2(float64(byteLen))))
	uintLen, err := c.ReadUint(lenBits)
	if err != nil {
		return nil, err
	}
	value, err := c.ReadBigUint(int(uintLen) * 8)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// ReadUnary
// TL-B:
// unary_zero$0 = Unary ~0;
// unary_succ$1 {n:#} x:(Unary ~n) = Unary ~(n + 1);
func (c *CellReader) ReadUnary() (uint64, error) {
	var n uint64
	for {
		succ, err := c.ReadBit()
		if err != nil {
			return 0, err
		}
		if succ {
			n++
		} else {
			return n, nil
		}
	}
}

// ReadMaybe
// TL-B:
// nothing$0 {X:Type} = Maybe X;
// just$1 {X:Type} value:X = Maybe X;
func (c *CellReader) ReadMaybe() (bool, error) {
	return c.ReadBit()
}

// ReadBits
// {n:#} (n * Bit)
func (c *CellReader) ReadBits(n int) (BitString, error) {
	bitString := NewBitString(n)
	for i := 0; i < n; i++ {
		bit, err := c.ReadBit()
		if err != nil {
			return BitString{}, err
		}
		err = bitString.WriteBit(bit)
		if err != nil {
			return BitString{}, err
		}
	}
	return bitString, nil
}
