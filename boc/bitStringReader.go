package boc

import (
	"math"
	"math/big"
)

type BitStringReader struct {
	buf    []byte
	len    int
	cursor int
}

func NewBitStringReader(bitString *BitString) BitStringReader {
	var reader = BitStringReader{
		buf:    bitString.Buffer(),
		len:    0,
		cursor: 0,
	}
	return reader
}

func (s *BitStringReader) getBit(n int) bool {
	return (s.buf[(n/8)|0] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitStringReader) Skip(n int) {
	for i := 0; i < n; i++ {
		s.ReadBit()
	}
}

func (s *BitStringReader) ReadBit() bool {
	var bit = s.getBit(s.cursor)
	s.cursor++
	return bit
}

func (s *BitStringReader) ReadBigUint(bitLen int) *big.Int {
	if bitLen == 0 {
		return big.NewInt(0)
	}
	var res = ""
	for i := 0; i < bitLen; i++ {
		if s.ReadBit() {
			res += "1"
		} else {
			res += "0"
		}
	}
	var num = big.NewInt(0)
	num.SetString(res, 2)
	return num
}

func (s *BitStringReader) ReadBigInt(bitLen int) *big.Int {
	if bitLen == 0 {
		return big.NewInt(0)
	}
	if bitLen == 1 {
		if s.ReadBit() {
			return big.NewInt(-1)
		} else {
			return big.NewInt(0)
		}
	}

	if s.ReadBit() {
		var base = s.ReadBigUint(bitLen - 1)
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return base.Sub(base, nb)
	} else {
		return s.ReadBigUint(bitLen - 1)
	}
}

func (s *BitStringReader) ReadUint(bitLen int) uint {
	if bitLen == 0 {
		return 0
	}

	var res uint = 0

	for i := bitLen - 1; i >= 0; i-- {
		if s.ReadBit() {
			res |= 1 << i
		}
	}

	return res
}

func (s *BitStringReader) ReadInt(bitLen int) int {
	if bitLen == 0 {
		return 0
	}
	if bitLen == 1 {
		if s.ReadBit() {
			return -1
		} else {
			return 0
		}
	}

	if s.ReadBit() {
		base := s.ReadUint(bitLen - 1)
		return int(base - uint(math.Pow(2, float64(bitLen-1))))
	} else {
		return int(s.ReadUint(bitLen - 1))
	}
}

func (s *BitStringReader) ReadCoins() uint {
	bytes := s.ReadUint(4)
	if bytes == 0 {
		return 0
	}
	return s.ReadUint(int(bytes * 8))
}

func (s *BitStringReader) ReadByte() byte {
	return byte(s.ReadUint(8))
}

func (s *BitStringReader) ReadBytes(size int) []byte {
	res := make([]byte, size)

	for i := 0; i < size; i++ {
		res[i] = byte(s.ReadUint(8))
	}

	return res
}
