package tolk

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type MapValue struct {
	keys   []Value
	values []Value
	len    int
}

func (m *MapValue) Unmarshal(cell *boc.Cell, ty tolkParser.Map, decoder *Decoder) error {
	keySize, ok := ty.K.GetFixedSize()
	if !ok {
		return fmt.Errorf("%v is not comparable", ty.K.SumType)
	}
	keyPrefix := boc.NewBitString(keySize)

	isNotEmpty, err := cell.ReadBit()
	if err != nil {
		return err
	}
	mp := MapValue{
		keys:   make([]Value, 0),
		values: make([]Value, 0),
	}
	if isNotEmpty {
		mpCell, err := cell.NextRef()
		if err != nil {
			return err
		}
		err = mapInner(keySize, keySize, mpCell, &keyPrefix, ty.K, ty.V, &mp.keys, &mp.values, decoder)
		if err != nil {
			return err
		}
	}

	m.len = len(mp.keys)
	m.keys = mp.keys
	m.values = mp.values

	return nil
}

func mapInner(
	keySize, leftKeySize int,
	c *boc.Cell,
	keyPrefix *boc.BitString,
	kt, vt tolkParser.Ty,
	keys, values *[]Value,
	decoder *Decoder,
) error {
	var err error
	var size int
	if c.CellType() == boc.PrunedBranchCell {
		return nil
	}
	size, keyPrefix, err = loadLabel(leftKeySize, c, keyPrefix)
	if err != nil {
		return err
	}
	// until key size is not equals we go deeper
	if keyPrefix.BitsAvailableForRead() < keySize {
		// 0 bit branch
		left, err := c.NextRef()
		if err != nil {
			return err
		}
		lp := keyPrefix.Copy()
		err = lp.WriteBit(false)
		if err != nil {
			return err
		}
		err = mapInner(keySize, leftKeySize-(1+size), left, &lp, kt, vt, keys, values, decoder)
		if err != nil {
			return err
		}
		// 1 bit branch
		right, err := c.NextRef()
		if err != nil {
			return err
		}
		rp := keyPrefix.Copy()
		err = rp.WriteBit(true)
		if err != nil {
			return err
		}
		err = mapInner(keySize, leftKeySize-(1+size), right, &rp, kt, vt, keys, values, decoder)
		if err != nil {
			return err
		}
		return nil
	}
	// add node to map
	v := Value{}
	err = v.Unmarshal(c, vt, decoder)
	if err != nil {
		return err
	}
	*values = append(*values, v)

	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return err
	}
	k := Value{}
	cell := boc.NewCellWithBits(key)
	err = k.Unmarshal(cell, kt, decoder)
	if err != nil {
		return err
	}
	*keys = append(*keys, k)

	return nil
}

func loadLabel(size int, c *boc.Cell, key *boc.BitString) (int, *boc.BitString, error) {
	first, err := c.ReadBit()
	if err != nil {
		return 0, nil, err
	}
	// hml_short$0
	if !first {
		// Unary, while 1, add to ln
		ln, err := c.ReadUnary()
		if err != nil {
			return 0, nil, err
		}
		// add bits to key
		for i := 0; i < int(ln); i++ {
			bit, err := c.ReadBit()
			if err != nil {
				return 0, nil, err
			}
			err = key.WriteBit(bit)
			if err != nil {
				return 0, nil, err
			}
		}
		return int(ln), key, nil
	}
	second, err := c.ReadBit()
	if err != nil {
		return 0, nil, err
	}
	// hml_long$10
	if !second {
		ln, err := c.ReadLimUint(size)
		if err != nil {
			return 0, nil, err
		}
		for i := 0; i < int(ln); i++ {
			bit, err := c.ReadBit()
			if err != nil {
				return 0, nil, err
			}
			err = key.WriteBit(bit)
			if err != nil {
				return 0, nil, err
			}
		}
		return int(ln), key, nil
	}
	// hml_same$11
	bitType, err := c.ReadBit()
	if err != nil {
		return 0, nil, err
	}
	ln, err := c.ReadLimUint(size)
	if err != nil {
		return 0, nil, err
	}
	for i := 0; i < int(ln); i++ {
		err = key.WriteBit(bitType)
		if err != nil {
			return 0, nil, err
		}
	}
	return int(ln), key, nil
}

func (m *MapValue) Marshal(cell *boc.Cell, ty tolkParser.Map, encoder *Encoder) error {
	keySize, ok := ty.K.GetFixedSize()
	if !ok {
		return fmt.Errorf("map key is not a comparable type, got %v", ty.K.SumType)
	}

	if len(m.keys) != len(m.values) {
		return fmt.Errorf("map keys and values lengths do not match")
	}

	if len(m.values) == 0 {
		err := cell.WriteBit(false)
		if err != nil {
			return err
		}
		return nil
	}

	err := cell.WriteBit(true)
	if err != nil {
		return err
	}

	keys := make([]boc.BitString, len(m.keys))
	for i, k := range m.keys {
		keyCell := boc.NewCell()
		err = k.Marshal(keyCell, ty.K, encoder)
		if err != nil {
			return err
		}
		keys[i] = keyCell.RawBitString()
	}

	ref := boc.NewCell()
	err = encodeMap(ref, keys, m.values, keySize, ty.V, encoder)
	if err != nil {
		return err
	}

	err = cell.AddRef(ref)
	if err != nil {
		return err
	}

	return nil
}

func encodeMap(c *boc.Cell, keys []boc.BitString, values []Value, keySize int, vt tolkParser.Ty, encoder *Encoder) error {
	if len(keys) == 0 || len(values) == 0 {
		return fmt.Errorf("keys or values are empty")
	}
	label, err := encodeLabel(c, &keys[0], &keys[len(keys)-1], keySize)
	if err != nil {
		return err
	}
	keySize = keySize - label.BitsAvailableForRead() - 1 // l = n - m - 1 // see tlb
	var leftKeys, rightKeys []boc.BitString
	var leftValues, rightValues []Value
	if len(keys) > 1 {
		for i := range keys {
			_, err = keys[i].ReadBits(label.BitsAvailableForRead()) // skip common label
			if err != nil {
				return err
			}
			isRight, err := keys[i].ReadBit()
			if err != nil {
				return err
			}
			if isRight {
				rightKeys = append(rightKeys, keys[i].ReadRemainingBits())
				rightValues = append(rightValues, values[i])
			} else {
				leftKeys = append(leftKeys, keys[i].ReadRemainingBits())
				leftValues = append(leftValues, values[i])
			}
		}
		l, err := c.NewRef()
		if err != nil {
			return err
		}
		err = encodeMap(l, leftKeys, leftValues, keySize, vt, encoder)
		if err != nil {
			return err
		}
		r, err := c.NewRef()
		if err != nil {
			return err
		}
		err = encodeMap(r, rightKeys, rightValues, keySize, vt, encoder)
		if err != nil {
			return err
		}
		return err
	}
	// marshal value
	err = values[0].Marshal(c, vt, encoder)
	if err != nil {
		return err
	}
	return nil
}

func encodeLabel(c *boc.Cell, keyFirst, keyLast *boc.BitString, keySize int) (boc.BitString, error) {
	label := boc.NewBitString(keySize)
	if keyFirst != keyLast {
		bitLeft, err := keyFirst.ReadBit()
		if err != nil {
			return boc.BitString{}, err
		}
		for keyFirst.BitsAvailableForRead() > 0 {
			bitRight, err := keyLast.ReadBit()
			if err != nil {
				return boc.BitString{}, err
			}
			if bitLeft != bitRight {
				break
			}
			if err := label.WriteBit(bitLeft); err != nil {
				return boc.BitString{}, err
			}
			bitLeft, err = keyFirst.ReadBit()
			if err != nil {
				return boc.BitString{}, err
			}
		}
	} else {
		label = keyFirst.Copy()
	}
	keyFirst.ResetCounter()
	keyLast.ResetCounter()
	labelLen := label.BitsAvailableForRead()

	// We must find the most compact way to serialize key
	hmlShortSize := 2*labelLen + 2
	hmlLongSize := 2 + int(math.Ceil(math.Log2(float64(keySize)+1))) + labelLen
	var encodeFunc func(*boc.Cell, int, boc.BitString) error
	isShort := false
	if hmlShortSize <= hmlLongSize {
		isShort = true
		encodeFunc = encodeShortLabel
	} else {
		encodeFunc = encodeLongLabel
	}

	// If all bits in label are the same then we can use hml_same
	isAllZero := true
	isAllOne := true
	for label.BitsAvailableForRead() > 0 {
		bit, err := label.ReadBit()
		if err != nil {
			return boc.BitString{}, err
		}
		if bit {
			isAllZero = false
		} else {
			isAllOne = false
		}
	}
	label.ResetCounter()
	if isAllZero || isAllOne {
		hmlSameSize := 2 + 1 + int(math.Ceil(math.Log2(float64(keySize)+1)))
		if isShort && hmlSameSize < hmlShortSize {
			encodeFunc = encodeSameLabel
		} else if !isShort && hmlSameSize < hmlLongSize {
			encodeFunc = encodeSameLabel
		}
	}

	err := encodeFunc(c, keySize, label)
	if err != nil {
		return boc.BitString{}, err
	}
	return label, nil
}

func encodeShortLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	//hml_short$0 {m:#} {n:#} len:(Unary ~n) {n <= m} s:(n * Bit) = HmLabel ~n m;
	err := c.WriteBit(false)
	if err != nil {
		return err
	}
	err = c.WriteUnary(uint(label.BitsAvailableForRead()))
	if err != nil {
		return err
	}
	err = c.WriteBitString(label)
	if err != nil {
		return err
	}
	return nil
}

func encodeLongLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	// hml_long$10 {m:#} n:(#<= m) s:(n * Bit) = HmLabel ~n m;
	err := c.WriteBit(true)
	if err != nil {
		return err
	}
	err = c.WriteBit(false)
	if err != nil {
		return err
	}
	err = c.WriteLimUint(label.BitsAvailableForRead(), keySize)
	if err != nil {
		return err
	}
	err = c.WriteBitString(label)
	if err != nil {
		return err
	}
	return nil
}

func encodeSameLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	//hml_same$11 {m:#} v:Bit n:(#<= m) = HmLabel ~n m;
	err := c.WriteUint(0b11, 2)
	if err != nil {
		return err
	}
	err = c.WriteBit(label.BinaryString()[0] == '1')
	if err != nil {
		return err
	}
	err = c.WriteLimUint(label.BitsAvailableForRead(), keySize)
	if err != nil {
		return err
	}
	return nil
}

func (m *MapValue) Equal(other any) bool {
	otherMapValue, ok := other.(MapValue)
	if !ok {
		return false
	}
	if m.len != otherMapValue.len {
		return false
	}
	for i := range m.keys {
		if !m.keys[i].Equal(otherMapValue.keys[i]) {
			return false
		}
		if !m.values[i].Equal(otherMapValue.values[i]) {
			return false
		}
	}
	return true
}

func (m *MapValue) Get(key Value) (Value, bool) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			return m.values[i], true
		}
	}

	return Value{}, false
}

func (m *MapValue) GetBySmallInt(v Int64) (Value, bool) {
	key := Value{
		sumType:  "smallInt",
		smallInt: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetBySmallUInt(v UInt64) (Value, bool) {
	key := Value{
		sumType:   "smallUint",
		smallUint: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBigInt(v BigInt) (Value, bool) {
	key := Value{
		sumType: "bigInt",
		bigInt:  &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBigUInt(v BigUInt) (Value, bool) {
	key := Value{
		sumType: "bigUint",
		bigUint: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBits(v Bits) (Value, bool) {
	key := Value{
		sumType: "bits",
		bits:    &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByInternalAddress(v InternalAddress) (Value, bool) {
	key := Value{
		sumType:         "internalAddress",
		internalAddress: &v,
	}
	return m.Get(key)
}

func (m *MapValue) Set(key Value, value Value) (bool, error) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			m.values[i] = value
			return true, nil
		}
	}

	m.keys = append(m.keys, Value(key))
	m.values = append(m.values, value)
	m.len++
	return true, nil
}

func (m *MapValue) SetBySmallInt(k Int64, value Value) (bool, error) {
	key := Value{
		sumType:  "smallInt",
		smallInt: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetBySmallUInt(k UInt64, value Value) (bool, error) {
	key := Value{
		sumType:   "smallUint",
		smallUint: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBigInt(k BigInt, value Value) (bool, error) {
	key := Value{
		sumType: "bigInt",
		bigInt:  &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBigUInt(k BigUInt, value Value) (bool, error) {
	key := Value{
		sumType: "bigUint",
		bigUint: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBits(k Bits, value Value) (bool, error) {
	key := Value{
		sumType: "bits",
		bits:    &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByInternalAddress(k InternalAddress, value Value) (bool, error) {
	key := Value{
		sumType:         "internalAddress",
		internalAddress: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) Delete(key Value) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			m.keys[i] = Value{}
			m.values[i] = Value{}
			m.len--
		}
	}
}

func (m *MapValue) DeleteBySmallInt(k Int64) {
	key := Value{
		sumType:  "smallInt",
		smallInt: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteBySmallUInt(k UInt64) {
	key := Value{
		sumType:   "smallUint",
		smallUint: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBigInt(k BigInt) {
	key := Value{
		sumType: "bigInt",
		bigInt:  &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBigUInt(k BigUInt) {
	key := Value{
		sumType: "bigUint",
		bigUint: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBits(k Bits) {
	key := Value{
		sumType: "bits",
		bits:    &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByInternalAddress(k InternalAddress) {
	key := Value{
		sumType:         "internalAddress",
		internalAddress: &k,
	}
	m.Delete(key)
}

func (m *MapValue) Len() int {
	return m.len
}

func (m *MapValue) MarshalJSON() ([]byte, error) {
	s := strings.Builder{}
	s.WriteString("{\n")
	for i, k := range m.keys {
		s.WriteString(fmt.Sprintf("\"%v\":", k))
		val, err := json.Marshal(m.values[i])
		if err != nil {
			return nil, err
		}
		s.WriteString(string(val))
	}
	s.WriteString("}")

	return []byte(s.String()), nil
}
