package boc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
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
	return (s.buf[n/8] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitString) On(n int) error {
	err := s.checkRange(n)
	if err != nil {
		return err
	}
	s.buf[n/8] |= 1 << (7 - (n % 8))
	return nil
}

func (s *BitString) Off(n int) error {
	err := s.checkRange(n)
	if err != nil {
		return err
	}
	s.buf[n/8] &= ^(1 << (7 - (n % 8)))
	return nil
}

func (s *BitString) Toggle(n int) error {
	err := s.checkRange(n)
	if err != nil {
		return err
	}
	s.buf[n/8] ^= 1 << (7 - (n % 8))
	return nil
}

func (s *BitString) WriteBit(val bool) error {
	if val {
		err := s.On(s.cursor)
		if err != nil {
			return err
		}
	} else {
		err := s.Off(s.cursor)
		if err != nil {
			return err
		}
	}
	s.cursor++
	return nil
}

func (s *BitString) WriteBitArray(val []bool) error {
	for _, item := range val {
		err := s.WriteBit(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BitString) WriteBigUint(val *big.Int, bitLen int) error {
	if bitLen == 0 || val.BitLen() > bitLen {
		return errors.New("bit length is too small")
	}

	for i := bitLen - 1; i >= 0; i-- {
		err := s.WriteBit(val.Bit(i) > 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BitString) WriteBigInt(val *big.Int, bitLen int) error {
	if bitLen == 1 {
		if val.Int64() == -1 {
			err := s.WriteBit(true)
			if err != nil {
				return err
			}
			return nil
		}
		if val.Int64() == 0 {
			err := s.WriteBit(false)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("bit length is too small")
	}
	if val.Sign() == -1 {
		err := s.WriteBit(true)
		if err != nil {
			return err
		}
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		err = s.WriteBigUint(nb.Add(nb, val), bitLen-1)
		if err != nil {
			return err
		}
	} else {
		err := s.WriteBit(false)
		if err != nil {
			return err
		}
		err = s.WriteBigUint(val, bitLen-1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BitString) WriteUint(val uint, bitLen int) error {
	for i := bitLen - 1; i >= 0; i-- {
		err := s.WriteBit(((val >> i) & 1) > 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BitString) WriteInt(val int, bitLen int) error {
	if bitLen == 1 {
		if val == -1 {
			err := s.WriteBit(true)
			if err != nil {
				return err
			}
		}
		if val == 0 {
			err := s.WriteBit(false)
			if err != nil {
				return err
			}
		}
	} else {
		if val < 0 {
			err := s.WriteBit(true)
			if err != nil {
				return err
			}
			err = s.WriteUint(uint(1<<(bitLen-1)+val), bitLen-1)
			if err != nil {
				return err
			}
		} else {
			err := s.WriteBit(false)
			if err != nil {
				return err
			}
			err = s.WriteUint(uint(val), bitLen-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BitString) WriteCoins(amount uint) error {
	if amount == 0 {
		err := s.WriteUint(0, 4)
		if err != nil {
			return err
		}
	} else {

		l := (bits.Len(amount) + 7) & -8
		err := s.WriteInt(l, 4)
		if err != nil {
			return err
		}
		err = s.WriteUint(amount, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BitString) WriteByte(val byte) error {
	err := s.WriteUint(uint(val), 8)
	if err != nil {
		return err
	}
	return nil
}

func (s *BitString) WriteBytes(data []byte) error {
	for _, item := range data {
		err := s.WriteByte(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BitString) WriteAddress(address *Address) error {
	if address == nil {
		err := s.WriteUint(0, 2)
		if err != nil {
			return err
		}
	} else {
		err := s.WriteUint(2, 2)
		if err != nil {
			return err
		}
		err = s.WriteUint(0, 1)
		if err != nil {
			return err
		}
		err = s.WriteInt(address.Workchain, 8)
		if err != nil {
			return err
		}
		err = s.WriteBytes(address.Address)
		if err != nil {
			return err
		}
	}
	return nil
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
		s.cursor--
		if s.Get(s.cursor) {
			foundEndBit = true
			err := s.Off(s.cursor)
			if err != nil {
				return err
			}
			break
		}
	}
	if !foundEndBit {
		return errors.New("incorrect topUppedArray")
	}
	return nil
}

func (s *BitString) GetTopUppedArray() ([]byte, error) {
	ret := s.Copy()

	tu := int(math.Ceil(float64(ret.Cursor())/8))*8 - ret.Cursor()
	if tu > 0 {
		tu = tu - 1
		err := ret.WriteBit(true)
		if err != nil {
			return nil, err
		}
		for tu > 0 {
			tu = tu - 1
			err := ret.WriteBit(false)
			if err != nil {
				return nil, err
			}
		}
	}
	ret.buf = ret.buf[0 : ((ret.cursor+7)&-8)/8]
	return ret.buf, nil
}

func (s *BitString) Print() {
	for _, n := range s.buf {
		fmt.Printf("% 08b\n", n)
	}
}

func (s *BitString) ToFiftHex() string {
	if s.cursor%4 == 0 {
		str := strings.ToUpper(hex.EncodeToString(s.buf[0 : (s.cursor+7)/8]))
		if s.cursor%8 == 0 {
			return str
		}
		return str[0 : len(str)-1]
	}
	temp := s.Copy()
	temp.WriteBit(true)
	for temp.cursor%4 != 0 {
		temp.WriteBit(false)
	}
	hex := temp.ToFiftHex()
	return hex + "_"
}

func (s *BitString) checkRange(n int) error {
	if n > s.Length() {
		return errors.New("BitString overflow")
	}
	return nil
}
