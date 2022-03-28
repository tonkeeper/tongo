package boc

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

var ErrNotEnoughBits = errors.New("not enough bits")

type BitStringReader struct {
	buf    []byte
	len    int
	cursor int
}

func NewBitStringReader(bitString *BitString) BitStringReader {
	var reader = BitStringReader{
		buf:    bitString.Buffer(),
		len:    bitString.len,
		cursor: 0,
	}
	return reader
}

func (s *BitStringReader) available() int {
	return s.len - s.cursor
}

func (s *BitStringReader) getBit(n int) bool {
	return (s.buf[n/8] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitStringReader) readBit() bool {
	var bit = s.getBit(s.cursor)
	s.cursor++
	return bit
}

func (s *BitStringReader) Skip(n int) error {
	if s.available() < n {
		return ErrNotEnoughBits
	}
	s.cursor += n
	return nil
}

func (s *BitStringReader) ReadBit() (bool, error) {
	if s.available() < 1 {
		return false, ErrNotEnoughBits
	}
	var bit = s.getBit(s.cursor)
	s.cursor++
	return bit, nil
}

func (s *BitStringReader) ReadBigUint(bitLen int) (*big.Int, error) {
	if s.available() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	var res = ""
	for i := 0; i < bitLen; i++ {
		if s.readBit() {
			res += "1"
		} else {
			res += "0"
		}
	}
	var num = big.NewInt(0)
	num.SetString(res, 2)
	return num, nil
}

func (s *BitStringReader) ReadBigInt(bitLen int) (*big.Int, error) {
	if s.available() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	if bitLen == 1 {
		if s.readBit() {
			return big.NewInt(-1), nil
		}
		return big.NewInt(0), nil
	}
	if s.readBit() {
		var base, _ = s.ReadBigUint(bitLen - 1)
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return base.Sub(base, nb), nil
	}
	return s.ReadBigUint(bitLen - 1)
}

func (s *BitStringReader) ReadUint(bitLen int) (uint, error) {
	if s.available() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	var res uint = 0
	for i := bitLen - 1; i >= 0; i-- {
		if s.readBit() {
			res |= 1 << i
		}
	}
	return res, nil
}

func (s *BitStringReader) ReadInt(bitLen int) (int, error) {
	if s.available() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	if bitLen == 1 {
		if s.readBit() {
			return -1, nil
		}
		return 0, nil
	}
	if s.readBit() {
		base, err := s.ReadUint(bitLen - 1)
		if err != nil {
			return 0, err
		}
		return int(base - uint(math.Pow(2, float64(bitLen-1)))), nil
	}
	res, err := s.ReadUint(bitLen - 1)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (s *BitStringReader) ReadCoins() (uint, error) {
	bytes, err := s.ReadUint(4)
	if err != nil {
		return 0, err
	}
	if bytes == 0 {
		return 0, nil
	}
	return s.ReadUint(int(bytes * 8))
}

func (s *BitStringReader) ReadByte() (byte, error) {
	res, err := s.ReadUint(8)
	if err != nil {
		return 0, err
	}
	return byte(res), nil
}

func (s *BitStringReader) ReadBytes(size int) ([]byte, error) {
	if s.available() < size*8 {
		return nil, ErrNotEnoughBits
	}
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		b, err := s.ReadUint(8)
		if err != nil {
			return nil, err
		}
		res[i] = byte(b)
	}
	return res, nil
}

func (s *BitStringReader) ReadAddress() (*Address, error) {
	prefix, err := s.ReadUint(2)
	if err != nil {
		return nil, err
	}
	if prefix == 0 { // adr_none prefix
		return nil, nil
	}
	if prefix != 2 { // not adr_std prefix
		return nil, fmt.Errorf("not std address")
	}
	maybe, err := s.ReadBit()
	if err != nil {
		return nil, err
	}
	if maybe == true {
		return nil, fmt.Errorf("anycast not being processed") //TODO: add anycast processing
	}
	workchain, err := s.ReadInt(8)
	if err != nil {
		return nil, err
	}
	addr, err := s.ReadBytes(32)
	if err != nil {
		return nil, err
	}
	var address Address
	address.Workchain = int32(workchain)
	address.Address = addr
	return &address, nil
}
