package boc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
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

func (s *BitString) Copy() BitString {
	var buf = make([]byte, len(s.buf))
	copy(buf, s.buf)

	return BitString{
		buf:    buf,
		len:    len(s.buf) * 8,
		cursor: s.cursor,
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

func (s *BitString) WriteBigUint(val *big.Int, bitLen int) error {
	if bitLen == 0 || val.BitLen() > bitLen {
		return errors.New("bit length is too small")
	}

	for i := bitLen - 1; i >= 0; i-- {
		s.WriteBit(val.Bit(i) > 0)
	}

	return nil
}

func (s *BitString) WriteBigInt(val *big.Int, bitLen int) error {
	if bitLen == 1 {
		if val.Int64() == -1 {
			s.WriteBit(true)
			return nil
		}
		if val.Int64() == 0 {
			s.WriteBit(false)
			return nil
		}
		return errors.New("bit length is too small")
	} else {
		if val.Sign() == -1 {
			s.WriteBit(true)
			var b = big.NewInt(2)
			var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
			err := s.WriteBigUint(nb.Add(nb, val), bitLen-1)
			if err != nil {
				return err
			}
		} else {
			s.WriteBit(false)
			err := s.WriteBigUint(val, bitLen-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BitString) WriteUint(val int, bitLen int) {
	for i := bitLen - 1; i >= 0; i-- {
		s.WriteBit(((val >> i) & 1) > 0)
	}
}

func (s *BitString) WriteInt(val int, bitLen int) {
	if bitLen == 1 {
		if val == -1 {
			s.WriteBit(true)
		}
		if val == 0 {
			s.WriteBit(false)
		}
	} else {
		if val < 0 {
			s.WriteBit(true)
			s.WriteUint(int(math.Pow(2, float64(bitLen-1)))+val, bitLen-1)
		} else {
			s.WriteBit(false)
			s.WriteUint(val, bitLen-1)
		}
	}
}

func (s *BitString) WriteCoins(amount int) {
	if amount == 0 {
		s.WriteUint(0, 4)
	} else {
		l := int(math.Ceil(float64(len(strconv.FormatInt(int64(amount), 16)) / 2)))
		s.WriteUint(l, 4)
		s.WriteUint(amount, l*8)
	}
}

func (s *BitString) WriteByte(val byte) {
	s.WriteUint(int(val), 8)
}

func (s *BitString) WriteBytes(data []byte) {
	for _, item := range data {
		s.WriteByte(item)
	}
}

func (s *BitString) WriteAddress(address *Address) {
	if address == nil {
		s.WriteUint(0, 2)
	} else {
		s.WriteUint(2, 2)
		s.WriteUint(0, 1)
		s.WriteInt(address.Workchain, 8)
		s.WriteBytes(address.Address)
	}
}

func (s *BitString) SetTopUppedArray(arr []byte, fulfilledBytes bool) error {
	s.len = len(arr) * 8
	s.buf = make([]byte, len(arr))
	copy(s.buf, arr)
	s.cursor = s.len

	if fulfilledBytes || s.len == 0 {
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

func (s *BitString) GetTopUppedArray() []byte {
	ret := s.Copy()

	tu := int(math.Ceil(float64(ret.Cursor())/8))*8 - ret.Cursor()
	if tu > 0 {
		tu = tu - 1
		ret.WriteBit(true)
		for tu > 0 {
			tu = tu - 1
			ret.WriteBit(false)
		}
	}
	ret.buf = ret.buf[0:int(math.Ceil(float64(ret.cursor/8)))]
	return ret.buf
}

func (s *BitString) Print() {
	for _, n := range s.buf {
		fmt.Printf("% 08b", n)
	}
	fmt.Println()
}

func (s *BitString) ToFiftHex() string {
	if s.cursor%4 == 0 {
		str := strings.ToUpper(hex.EncodeToString(s.buf[0 : (s.cursor+7)/8]))
		if s.cursor%8 == 0 {
			return str
		} else {
			return str[0 : len(str)-1]
		}
	} else {
		temp := s.Copy()
		temp.WriteBit(true)
		for temp.cursor%4 != 0 {
			temp.WriteBit(false)
		}
		hex := temp.ToFiftHex()
		return hex + "_"
	}
}
