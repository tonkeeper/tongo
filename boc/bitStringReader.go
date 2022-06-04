package boc

import (
	"fmt"
	"math"
	"math/big"
)

//var ErrNotEnoughBits = errors.New("not enough bits")

// BitStringReader DEPRECATED
type BitStringReader struct {
	buf    []byte
	len    int
	cursor int
}

func NewBitStringReader(bitString *BitString) BitStringReader {
	var reader = BitStringReader{
		buf:    bitString.Buffer(),
		len:    bitString.cap,
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

// ReadGrams
// TL-B: nanograms$_ amount:(VarUInteger 16) = Grams;
func (s *BitStringReader) ReadGrams() (uint64, error) {
	grams, err := s.ReadVarUint(16)
	if err != nil {
		return 0, err
	}
	if !grams.IsUint64() {
		return 0, fmt.Errorf("grams uint64 overflow")
	}
	return grams.Uint64(), nil
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

// ReadVarUint
// TL-B: var_uint$_ {n:#} len:(#< n) value:(uint (len * 8)) = VarUInteger n;
func (s *BitStringReader) ReadVarUint(byteLen int) (*big.Int, error) {
	if byteLen < 2 {
		return nil, fmt.Errorf("invalid varuint length")
	}
	lenBits := int(math.Ceil(math.Log2(float64(byteLen))))
	uintLen, err := s.ReadUint(lenBits)
	if err != nil {
		return nil, err
	}
	value, err := s.ReadBigUint(int(uintLen) * 8)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// ReadBits
// {n:#} (n * Bit)
func (s *BitStringReader) ReadBits(n uint) (BitString, error) {
	var bitString BitString
	for i := uint(0); i < n; i++ {
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
