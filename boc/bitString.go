package boc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

var ErrNotEnoughBits = errors.New("not enough bits")
var ErrBitStingOverflow = errors.New("BitString overflow")

type BitString struct {
	buf     []byte
	cap     int
	len     int
	rCursor int
}

func NewBitString(bitLen int) BitString {
	return BitString{
		buf:     make([]byte, int(math.Ceil(float64(bitLen)/float64(8)))),
		cap:     bitLen,
		rCursor: 0,
		len:     0,
	}
}

// Utils

func (s *BitString) BitsAvailableForWrite() int {
	return s.cap - s.len
}

func (s *BitString) BitsAvailableForRead() int {
	return s.len - s.rCursor
}

func (s *BitString) GetWriteCursor() int {
	return s.len
}

func (s *BitString) Buffer() []byte {
	return s.buf
}

func (s *BitString) checkRange(n int) error {
	if n > s.cap {
		return ErrBitStingOverflow
	}
	return nil
}

func (s *BitString) SetTopUppedArray(arr []byte, fulfilledBytes bool) error {
	s.cap = len(arr) * 8
	s.buf = make([]byte, len(arr))
	copy(s.buf, arr)
	s.len = s.cap

	if fulfilledBytes || s.cap == 0 {
		return nil
	}

	foundEndBit := false

	for i := 0; i < 7; i++ {
		s.len--
		if s.mustGetBit(s.len) {
			foundEndBit = true
			err := s.Off(s.len) //todo: check
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

	tu := int(math.Ceil(float64(ret.GetWriteCursor())/8))*8 - ret.GetWriteCursor()
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
	ret.buf = ret.buf[0 : ((ret.len+7)&-8)/8]
	return ret.buf, nil
}

func (s *BitString) Copy() BitString {
	var buf = make([]byte, len(s.buf))
	copy(buf, s.buf)

	return BitString{
		buf: buf,
		cap: s.cap,
		len: s.len,
	}
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

// Read methods

func (s *BitString) mustGetBit(n int) bool {
	return (s.buf[n/8] & (1 << (7 - (n % 8)))) > 0
}

func (s *BitString) mustReadBit() bool {
	var bit = s.mustGetBit(s.rCursor)
	s.rCursor++
	return bit
}

func (s *BitString) Skip(n int) error {
	if s.BitsAvailableForRead() < n {
		return ErrNotEnoughBits
	}
	s.rCursor += n
	return nil
}

func (s *BitString) ReadBit() (bool, error) {
	if s.BitsAvailableForRead() < 1 {
		return false, ErrNotEnoughBits
	}
	var bit = s.mustGetBit(s.rCursor)
	s.rCursor++
	return bit, nil
}

func (s *BitString) ReadBigUint(bitLen int) (*big.Int, error) {
	if s.BitsAvailableForRead() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	var res = ""
	for i := 0; i < bitLen; i++ {
		if s.mustReadBit() {
			res += "1"
		} else {
			res += "0"
		}
	}
	var num = big.NewInt(0)
	num.SetString(res, 2)
	return num, nil
}

func (s *BitString) ReadBigInt(bitLen int) (*big.Int, error) {
	if s.BitsAvailableForRead() < bitLen {
		return nil, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return big.NewInt(0), nil
	}
	if bitLen == 1 {
		if s.mustReadBit() {
			return big.NewInt(-1), nil
		}
		return big.NewInt(0), nil
	}
	if s.mustReadBit() {
		var base, _ = s.ReadBigUint(bitLen - 1)
		var b = big.NewInt(2)
		var nb = b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return base.Sub(base, nb), nil
	}
	return s.ReadBigUint(bitLen - 1)
}

func (s *BitString) ReadUint(bitLen int) (uint64, error) {
	if bitLen > 64 {
		return 0, fmt.Errorf("too much bits for uint64")
	}
	if s.BitsAvailableForRead() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	var res uint64 = 0
	for i := bitLen - 1; i >= 0; i-- {
		if s.mustReadBit() {
			res |= 1 << i
		}
	}
	return res, nil
}

func (s *BitString) PickUint(bitLen int) (uint64, error) {
	res, err := s.ReadUint(bitLen)
	if err != nil {
		return 0, err
	}
	s.rCursor -= bitLen
	return res, nil
}

func (s *BitString) ReadInt(bitLen int) (int64, error) {
	if bitLen > 64 {
		return 0, fmt.Errorf("too much bits for int64")
	}
	if s.BitsAvailableForRead() < bitLen {
		return 0, ErrNotEnoughBits
	}
	if bitLen == 0 {
		return 0, nil
	}
	if bitLen == 1 {
		if s.mustReadBit() {
			return -1, nil
		}
		return 0, nil
	}
	if s.mustReadBit() {
		base, err := s.ReadUint(bitLen - 1)
		if err != nil {
			return 0, err
		}
		return int64(base - uint64(math.Pow(2, float64(bitLen-1)))), nil
	}
	res, err := s.ReadUint(bitLen - 1)
	if err != nil {
		return 0, err
	}
	return int64(res), nil
}

func (s *BitString) ReadByte() (byte, error) {
	res, err := s.ReadUint(8)
	if err != nil {
		return 0, err
	}
	return byte(res), nil
}

func (s *BitString) ReadBytes(size int) ([]byte, error) {
	if s.BitsAvailableForRead() < size*8 {
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

func (s *BitString) ReadBits(n int) (BitString, error) {
	bitString := NewBitString(n)
	for i := 0; i < n; i++ {
		bit, err := s.ReadBit()
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

// Write methods

func (s *BitString) WriteBit(val bool) error {
	if val {
		err := s.On(s.len)
		if err != nil {
			return err
		}
	} else {
		err := s.Off(s.len)
		if err != nil {
			return err
		}
	}
	s.len++
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

func (s *BitString) WriteBitString(bs BitString) error {
	bs.rCursor = 0
	for i := 0; i < bs.len; i++ {
		bit, err := bs.ReadBit()
		if err != nil {
			return err
		}
		err = s.WriteBit(bit)
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

func (s *BitString) WriteUint(val uint64, bitLen int) error {
	for i := bitLen - 1; i >= 0; i-- {
		err := s.WriteBit(((val >> i) & 1) > 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BitString) WriteInt(val int64, bitLen int) error {
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
			err = s.WriteUint(uint64(1<<(bitLen-1)+val), bitLen-1)
			if err != nil {
				return err
			}
		} else {
			err := s.WriteBit(false)
			if err != nil {
				return err
			}
			err = s.WriteUint(uint64(val), bitLen-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BitString) WriteByte(val byte) error {
	err := s.WriteUint(uint64(val), 8)
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

// Misc

func (s *BitString) ToFiftHex() string {
	if s.len%4 == 0 {
		str := strings.ToUpper(hex.EncodeToString(s.buf[0 : (s.len+7)/8]))
		if s.len%8 == 0 {
			return str
		}
		return str[0 : len(str)-1]
	}
	temp := s.Copy()
	temp.WriteBit(true)
	for temp.len%4 != 0 {
		temp.WriteBit(false)
	}
	hex := temp.ToFiftHex()
	return hex + "_"
}

func (s *BitString) Print() {
	for _, n := range s.buf {
		fmt.Printf("% 08b\n", n)
	}
}
func (s *BitString) BinaryString() string {
	buf := strings.Builder{}
	for i, n := range s.buf {
		if (i+1)*8 <= s.len {
			buf.WriteString(fmt.Sprintf("%08b", n))
		} else if (i)*8 > s.len {
			break
		} else {
			str := fmt.Sprintf("%08b", n)
			for j := 0; buf.Len() < s.len; j++ {
				buf.WriteByte(str[j])
			}
		}

	}
	return buf.String()
}

func (s *BitString) WriteUnary(n uint) error {
	if n < 63 {
		err := s.WriteUint(1<<n-1, int(n))
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < int(n); i++ {
			err := s.WriteBit(true)
			if err != nil {
				return err
			}
		}
	}
	return s.WriteBit(false)
}

func (s *BitString) ReadUnary() (uint, error) {
	var n uint
	for {
		bit, err := s.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit {
			n += 1
		} else {
			break
		}
	}
	return n, nil
}

// ReadLimUint
// #<= n
func (s *BitString) ReadLimUint(n int) (uint, error) {
	ln := int(math.Ceil(math.Log2(float64(n + 1))))
	res, err := s.ReadUint(ln)
	return uint(res), err
}

// WriteLimUint
// #<= n
func (s *BitString) WriteLimUint(val, n int) error {
	ln := int(math.Ceil(math.Log2(float64(n + 1))))
	err := s.WriteUint(uint64(val), ln)
	return err
}

func (s *BitString) ReadRemainingBits() BitString {
	bs, _ := s.ReadBits(s.BitsAvailableForRead())
	return bs
}

func (s BitString) MarshalTLB(cell *Cell, tag string) error {
	err := cell.WriteBitString(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *BitString) UnmarshalTLB(cell *Cell, tag string) error {
	ln, err := decodeBitStringTag(tag)
	if err != nil {
		return err
	}
	bs, err := cell.ReadBits(ln)
	if err != nil {
		return err
	}
	*s = bs
	return nil
}

func decodeBitStringTag(tag string) (int, error) {
	var n int
	if tag == "" {
		return 0, fmt.Errorf("empty BitString tag")
	}
	_, err := fmt.Sscanf(tag, "%dbits", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *BitString) ResetCounter() {
	s.rCursor = 0
}
