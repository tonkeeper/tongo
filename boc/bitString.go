package boc

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

type BitString struct {
	buf    []byte
	len    int
	cursor int
}

func NewBitString(bitLen int) BitString {
	return BitString{
		buf:    make([]byte, int(math.Ceil(float64(bitLen)/float64(8)))),
		len:    bitLen,
		cursor: 0,
	}
}

func (s *BitString) Available() int {
	return s.len - s.cursor
}

func (s *BitString) Length() int {
	return s.len
}

func (s *BitString) Cursor() int {
	return s.cursor
}

func (s *BitString) Buffer() []byte {
	return s.buf
}

func (s *BitString) Get(n int) bool {

	return (s.buf[(n/8)|0] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitString) On(n int) {
	s.buf[(n/8)|0] |= 1 << (7 - (n % 8))
}

func (s *BitString) Off(n int) {
	s.buf[(n/8)|0] &= ^(1 << (7 - (n % 8)))
}

func (s *BitString) Toggle(n int) {
	s.buf[(n/8)|0] ^= 1 << (7 - (n % 8))
}

func (s *BitString) WriteBit(val bool) {
	if val {
		s.On(s.cursor)
	} else {
		s.Off(s.cursor)
	}
	s.cursor++
}

func (s *BitString) WriteBitArray(val []bool) {
	for _, item := range val {
		s.WriteBit(item)
	}
}

func (s *BitString) WriteUint(val *big.Int, bitLen int) error {
	if bitLen == 0 || val.BitLen() > bitLen {
		return errors.New("bit length is too small")
	}

	for i := bitLen - 1; i >= 0; i-- {
		s.WriteBit(val.Bit(i) > 0)
	}

	return nil
}

func (s *BitString) WriteInt(val *big.Int, bitLen int) error {
	if bitLen == 1 {
		if val.Cmp(big.NewInt(-1)) == 0 {
			s.WriteBit(true)
			return nil
		}
		if val.Cmp(big.NewInt(0)) == 0 {
			s.WriteBit(false)
			return nil
		}
		return errors.New("bit length is too small")
	} else {
		if val.Sign() == -1 {
			s.WriteBit(true)
			var b = big.NewInt(2)
			var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
			err := s.WriteUint(nb.Add(nb, val), bitLen-1)
			if err != nil {
				return err
			}
		} else {
			s.WriteBit(false)
			err := s.WriteUint(val, bitLen-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BitString) SetTopUppedArray(arr []byte, fullfilledBytes bool) error {
	s.len = len(arr) * 8
	s.buf = make([]byte, len(arr))
	copy(s.buf, arr)
	s.cursor = s.len

	if fullfilledBytes || s.len == 0 {
		return nil
	}

	foundEndBit := false

	for i := 0; i < 7; i++ {
		s.cursor -= 1
		if s.Get(s.cursor) {
			foundEndBit = true
			s.Off(s.cursor)
			break
		}
	}
	if !foundEndBit {
		return errors.New("incorrect topUppedArray")
	}
	return nil
}

func (s *BitString) Print() {
	for _, n := range s.buf {
		fmt.Printf("% 08b", n)
	}
	fmt.Println()
}
