package boc

import "math/big"

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

func (s *BitStringReader) GetBit(n int) bool {
	return (s.buf[(n/8)|0] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitStringReader) ReadBit() bool {
	var bit = s.GetBit(s.cursor)
	s.cursor++
	return bit
}

func (s *BitStringReader) ReadUint(bitLen int) *big.Int {
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

func (s *BitStringReader) ReadInt(bitLen int) *big.Int {
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
		var base = s.ReadUint(bitLen - 1)
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return base.Sub(base, nb)
	} else {
		return s.ReadUint(bitLen - 1)
	}
}
