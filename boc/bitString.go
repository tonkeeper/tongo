package boc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var ErrNotEnoughBits = errors.New("not enough bits")
var ErrBitStingOverflow = errors.New("BitString overflow")
var ErrInvalidHex = errors.New("invalid hex representation of bitstring")

type BitString struct {
	buf     []byte
	cap     int
	len     int
	rCursor int
}

func NewBitString(bitLen int) BitString {
	return BitString{
		buf:     make([]byte, ((bitLen+7)&-8)/8),
		cap:     bitLen,
		rCursor: 0,
		len:     0,
	}
}

// BitStringFromFiftHex constructs a new BitString from the given hex representation.
func BitStringFromFiftHex(hexRepr string) (*BitString, error) {
	var endingBits string
	// "_" at the end means there is a padding
	if strings.HasSuffix(hexRepr, "_") {
		if len(hexRepr) < 2 {
			return nil, ErrInvalidHex
		}
		var ok bool
		endingBits, ok = suffixToBits[hexRepr[len(hexRepr)-2:]]
		if !ok {
			return nil, ErrInvalidHex
		}
		hexRepr = hexRepr[:len(hexRepr)-2]
	}
	bs := NewBitString(len(hexRepr)*4 + len(endingBits))
	for _, x := range hexRepr {
		oct, err := hexToInt(uint8(x))
		if err != nil {
			return nil, err
		}
		if err := bs.WriteUint(uint64(oct), 4); err != nil {
			return nil, err
		}
	}
	for _, bit := range endingBits {
		if err := bs.WriteBit(bit == '1'); err != nil {
			return nil, err
		}
	}
	return &bs, nil
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
	if n >= s.cap {
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

	tu := ((ret.GetWriteCursor() + 7) & -8) - ret.GetWriteCursor()
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
	var num = big.NewInt(0)
	var b []byte
	var err error
	if bitLen%8 != 0 {
		firstByte, err := s.ReadUint(bitLen % 8)
		if err != nil {
			return nil, err
		}
		b = []byte{byte(firstByte)}
	}
	b, err = s.ReadBytes(bitLen / 8)
	if err != nil {
		return nil, err
	}
	num.SetBytes(b)
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
		var base, err = s.ReadBigUint(bitLen - 1)
		if err != nil {
			return nil, err
		}
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
	if bitLen == 0 {
		return 0, fmt.Errorf("integer can't be zero size")
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
		return int64(base - 1<<(bitLen-1)), nil
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

// ToFiftHex returns a TON hex representation of this bit string.
func (s *BitString) ToFiftHex() string {
	// the idea of FiftHex is the following:
	// 1. split a bitstring into groups of 4 bits each
	// 2. each group is converted to a hex value
	//    if the last group is full (meaning it contains 4 bits), we are done at this step
	// 3. if the last group contains less than 4 bits,
	//    we add a padding to the last group.
	//    the padding starts with "1" following as many zeros as required to make the last group full.
	//    Then we convert the last group to a hex value and add "_" to the resulting hex string.
	if s.len%4 == 0 {
		str := strings.ToUpper(hex.EncodeToString(s.buf[0 : (s.len+7)/8]))
		if s.len%8 == 0 {
			return str
		}
		return str[0 : len(str)-1]
	}
	// there is a need to add a padding
	temp := s.Copy()
	temp.Grow(4 - s.len%4)
	temp.WriteBit(true)
	for temp.len%4 != 0 {
		temp.WriteBit(false)
	}
	return temp.ToFiftHex() + "_"
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
	ln := minBitsRequired(uint64(n))
	res, err := s.ReadUint(ln)
	return uint(res), err
}

// WriteLimUint
// #<= n
func (s *BitString) WriteLimUint(val, n int) error {
	ln := minBitsRequired(uint64(n))
	err := s.WriteUint(uint64(val), ln)
	return err
}

func (s *BitString) ReadRemainingBits() BitString {
	bs, _ := s.ReadBits(s.BitsAvailableForRead())
	return bs
}

func (s *BitString) ResetCounter() {
	s.rCursor = 0
}

func (s *BitString) Grow(bitLen int) {
	s.buf = append(s.buf, make([]byte, bitLen/8+1)...)
	s.cap += bitLen
}

func (s *BitString) Append(b BitString) {
	needBits := b.len - s.BitsAvailableForWrite()
	if needBits > 0 {
		s.Grow(needBits)
	}
	_ = s.WriteBitString(b) // must fit
}

var tab64 = [64]int{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5}

func minBitsRequired(value uint64) int {
	if value == 0 {
		return 0
	}
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	return tab64[((value-(value>>1))*0x07EDD5E59A4E28C2)>>58] + 1
}

var suffixToBits = map[string]string{
	"4_": "0",
	"C_": "1",
	"c_": "1",
	"2_": "00",
	"6_": "01",
	"A_": "10",
	"a_": "10",
	"E_": "11",
	"e_": "11",
	"1_": "000",
	"3_": "001",
	"5_": "010",
	"7_": "011",
	"9_": "100",
	"B_": "101",
	"b_": "101",
	"D_": "110",
	"d_": "110",
	"F_": "111",
	"f_": "111",
}

func hexToInt(c uint8) (uint8, error) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', nil
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, nil
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, nil
	}
	return 0, ErrInvalidHex
}

func (s BitString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", s.ToFiftHex())), nil
}

func (s *BitString) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	bs, err := BitStringFromFiftHex(str)
	if err != nil {
		return err
	}
	*s = *bs
	return nil
}
