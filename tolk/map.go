package tolk

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

type MapValue struct {
	keys   []Value
	values []Value
	len    int
}

func (m *MapValue) Unmarshal(cell *boc.Cell, ty tolkParser.Map, decoder *Decoder) error {
	keySize, ok := ty.K.GetFixedSize()
	if !ok {
		return fmt.Errorf("%v type is not comparable", ty.K.SumType)
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
			return fmt.Errorf("failed to get map's next ref: %w", err)
		}
		err = mapInner(keySize, keySize, mpCell, &keyPrefix, ty.K, ty.V, &mp.keys, &mp.values, decoder)
		if err != nil {
			return fmt.Errorf("failed to parse map value: %w", err)
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
		return fmt.Errorf("failed to load map's label: %w", err)
	}
	// until key size is not equals we go deeper
	if keyPrefix.BitsAvailableForRead() < keySize {
		// 0 bit branch
		left, err := c.NextRef()
		if err != nil {
			return fmt.Errorf("failed to get map's left branch: %w", err)
		}
		lp := keyPrefix.Copy()
		err = lp.WriteBit(false)
		if err != nil {
			return fmt.Errorf("failed to write map's left branch key prefix: %w", err)
		}
		err = mapInner(keySize, leftKeySize-(1+size), left, &lp, kt, vt, keys, values, decoder)
		if err != nil {
			return fmt.Errorf("failed to get map's left value: %w", err)
		}
		// 1 bit branch
		right, err := c.NextRef()
		if err != nil {
			return fmt.Errorf("failed to get map's right branch: %w", err)
		}
		rp := keyPrefix.Copy()
		err = rp.WriteBit(true)
		if err != nil {
			return fmt.Errorf("failed to write map's right branch key prefix: %w", err)
		}
		err = mapInner(keySize, leftKeySize-(1+size), right, &rp, kt, vt, keys, values, decoder)
		if err != nil {
			return fmt.Errorf("failed to get map's right value: %w", err)
		}
		return nil
	}
	// add node to map
	v := Value{}
	err = v.Unmarshal(c, vt, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal map's value: %w", err)
	}
	*values = append(*values, v)

	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return fmt.Errorf("failed to get map's key: %w", err)
	}
	k := Value{}
	cell := boc.NewCellWithBits(key)
	err = k.Unmarshal(cell, kt, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal map's key: %w", err)
	}
	*keys = append(*keys, k)

	return nil
}

func loadLabel(size int, c *boc.Cell, key *boc.BitString) (int, *boc.BitString, error) {
	first, err := c.ReadBit()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read map's first bit of label's type: %w", err)
	}
	// hml_short$0
	if !first {
		// Unary, while 1, add to ln
		ln, err := c.ReadUnary()
		if err != nil {
			return 0, nil, fmt.Errorf("failed to read map's short label's length: %w", err)
		}
		// add bits to key
		for i := 0; i < int(ln); i++ {
			bit, err := c.ReadBit()
			if err != nil {
				return 0, nil, fmt.Errorf("failed to read map's short label's %v-bit: %w", i, err)
			}
			err = key.WriteBit(bit)
			if err != nil {
				return 0, nil, fmt.Errorf("failed to write map's short label's %v-bit: %w", i, err)
			}
		}
		return int(ln), key, nil
	}
	second, err := c.ReadBit()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read map's second bit of label's type: %w", err)
	}
	// hml_long$10
	if !second {
		ln, err := c.ReadLimUint(size)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to read map's long label's length: %w", err)
		}
		for i := 0; i < int(ln); i++ {
			bit, err := c.ReadBit()
			if err != nil {
				return 0, nil, fmt.Errorf("failed to read map's long label's %v-bit: %w", i, err)
			}
			err = key.WriteBit(bit)
			if err != nil {
				return 0, nil, fmt.Errorf("failed to write map's long label's %v-bit: %w", i, err)
			}
		}
		return int(ln), key, nil
	}
	// hml_same$11
	bitType, err := c.ReadBit()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read map's same label's bit: %w", err)
	}
	ln, err := c.ReadLimUint(size)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read map's same label's length: %w", err)
	}
	for i := 0; i < int(ln); i++ {
		err = key.WriteBit(bitType)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to write map's same label's %v-bit: %w", i, err)
		}
	}
	return int(ln), key, nil
}

func (m *MapValue) Marshal(cell *boc.Cell, ty tolkParser.Map, encoder *Encoder) error {
	keySize, ok := ty.K.GetFixedSize()
	if !ok {
		return fmt.Errorf("%s type is not comparable", ty.K.SumType)
	}

	if len(m.keys) != len(m.values) {
		return fmt.Errorf("map keys and values lengths do not match")
	}

	if len(m.values) == 0 {
		err := cell.WriteBit(false)
		if err != nil {
			return fmt.Errorf("failed to write map's emptiness bit: %w", err)
		}
		return nil
	}

	err := cell.WriteBit(true)
	if err != nil {
		return fmt.Errorf("failed to write map's not-emptiness bit: %w", err)
	}

	keys := make([]boc.BitString, len(m.keys))
	for i, k := range m.keys {
		keyCell := boc.NewCell()
		err = k.Marshal(keyCell, ty.K, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal map's %v key: %w", i, err)
		}
		keys[i] = keyCell.RawBitString()
	}

	ref := boc.NewCell()
	err = encodeMap(ref, keys, m.values, keySize, ty.V, encoder)
	if err != nil {
		return fmt.Errorf("failed to encode map: %w", err)
	}

	err = cell.AddRef(ref)
	if err != nil {
		return fmt.Errorf("failed to add map's ref value: %w", err)
	}

	return nil
}

func encodeMap(c *boc.Cell, keys []boc.BitString, values []Value, keySize int, vt tolkParser.Ty, encoder *Encoder) error {
	if len(keys) == 0 || len(values) == 0 {
		return fmt.Errorf("keys or values are empty")
	}
	label, err := encodeLabel(c, &keys[0], &keys[len(keys)-1], keySize)
	if err != nil {
		return fmt.Errorf("failed to encode map's label: %w", err)
	}
	keySize = keySize - label.BitsAvailableForRead() - 1 // l = n - m - 1 // see tlb
	var leftKeys, rightKeys []boc.BitString
	var leftValues, rightValues []Value
	if len(keys) > 1 {
		for i := range keys {
			err = keys[i].Skip(label.BitsAvailableForRead()) // skip common label
			if err != nil {
				return fmt.Errorf("failed to skip map's key common label: %w", err)
			}
			isRight, err := keys[i].ReadBit()
			if err != nil {
				return fmt.Errorf("failed to read map's is right bit: %w", err)
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
			return fmt.Errorf("failed to create map's left value: %w", err)
		}
		err = encodeMap(l, leftKeys, leftValues, keySize, vt, encoder)
		if err != nil {
			return fmt.Errorf("failed to encode map's left value: %w", err)
		}
		r, err := c.NewRef()
		if err != nil {
			return fmt.Errorf("failed to create map's right value: %w", err)
		}
		err = encodeMap(r, rightKeys, rightValues, keySize, vt, encoder)
		if err != nil {
			return fmt.Errorf("failed to encode map's right value: %w", err)
		}
		return nil
	}
	// marshal value
	err = values[0].Marshal(c, vt, encoder)
	if err != nil {
		return fmt.Errorf("failed to marshal map's values: %w", err)
	}
	return nil
}

func encodeLabel(c *boc.Cell, keyFirst, keyLast *boc.BitString, keySize int) (boc.BitString, error) {
	label := boc.NewBitString(keySize)
	if keyFirst != keyLast {
		bitLeft, err := keyFirst.ReadBit()
		if err != nil {
			return boc.BitString{}, fmt.Errorf("failed to read map's first label's bit: %w", err)
		}
		for keyFirst.BitsAvailableForRead() > 0 {
			bitRight, err := keyLast.ReadBit()
			if err != nil {
				return boc.BitString{}, fmt.Errorf("failed to read map's last label's bit: %w", err)
			}
			if bitLeft != bitRight {
				break
			}
			if err := label.WriteBit(bitLeft); err != nil {
				return boc.BitString{}, fmt.Errorf("failed to write map's label bit: %w", err)
			}
			bitLeft, err = keyFirst.ReadBit()
			if err != nil {
				return boc.BitString{}, fmt.Errorf("failed to read map's first label's bit: %w", err)
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
			return boc.BitString{}, fmt.Errorf("failed to read map's label's bit: %w", err)
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
		return boc.BitString{}, fmt.Errorf("failed to encode map's label: %w", err)
	}
	return label, nil
}

func encodeShortLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	//hml_short$0 {m:#} {n:#} len:(Unary ~n) {n <= m} s:(n * Bit) = HmLabel ~n m;
	err := c.WriteBit(false)
	if err != nil {
		return fmt.Errorf("failed to write map's short label type: %w", err)
	}
	err = c.WriteUnary(uint(label.BitsAvailableForRead()))
	if err != nil {
		return fmt.Errorf("failed to write map's short label length: %w", err)
	}
	err = c.WriteBitString(label)
	if err != nil {
		return fmt.Errorf("failed to write map's short label value: %w", err)
	}
	return nil
}

func encodeLongLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	// hml_long$10 {m:#} n:(#<= m) s:(n * Bit) = HmLabel ~n m;
	err := c.WriteUint(0b10, 2)
	if err != nil {
		return fmt.Errorf("failed to write map's long label type: %w", err)
	}
	err = c.WriteLimUint(label.BitsAvailableForRead(), keySize)
	if err != nil {
		return fmt.Errorf("failed to write map's long label length: %w", err)
	}
	err = c.WriteBitString(label)
	if err != nil {
		return fmt.Errorf("failed to write map's long label value: %w", err)
	}
	return nil
}

func encodeSameLabel(c *boc.Cell, keySize int, label boc.BitString) error {
	//hml_same$11 {m:#} v:Bit n:(#<= m) = HmLabel ~n m;
	err := c.WriteUint(0b11, 2)
	if err != nil {
		return fmt.Errorf("failed to write map's same label type: %w", err)
	}
	err = c.WriteBit(label.BinaryString()[0] == '1')
	if err != nil {
		return fmt.Errorf("failed to write map's same label bit: %w", err)
	}
	err = c.WriteLimUint(label.BitsAvailableForRead(), keySize)
	if err != nil {
		return fmt.Errorf("failed to write map's same label length: %w", err)
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

func (m *MapValue) get(key Value) (Value, bool) {
	for i, k := range m.keys {
		if k.Equal(key) {
			return m.values[i], true
		}
	}

	return Value{}, false
}

func (m *MapValue) GetBySmallInt(v Int64) (Value, bool) {
	key := Value{
		SumType:  "smallInt",
		SmallInt: &v,
	}
	return m.get(key)
}

func (m *MapValue) GetBySmallUInt(v UInt64) (Value, bool) {
	key := Value{
		SumType:   "smallUint",
		SmallUint: &v,
	}
	return m.get(key)
}

func (m *MapValue) GetByBigInt(v BigInt) (Value, bool) {
	key := Value{
		SumType: "bigInt",
		BigInt:  &v,
	}
	return m.get(key)
}

func (m *MapValue) GetByBigUInt(v BigUInt) (Value, bool) {
	key := Value{
		SumType: "bigUint",
		BigUint: &v,
	}
	return m.get(key)
}

func (m *MapValue) GetByBits(v Bits) (Value, bool) {
	key := Value{
		SumType: "bits",
		Bits:    &v,
	}
	return m.get(key)
}

func (m *MapValue) GetByInternalAddress(v InternalAddress) (Value, bool) {
	key := Value{
		SumType:         "internalAddress",
		InternalAddress: &v,
	}
	return m.get(key)
}

func (m *MapValue) set(key Value, value Value) (bool, error) {
	for i, k := range m.keys {
		if k.Equal(key) {
			m.values[i] = value
			return true, nil
		}
	}

	m.keys = append(m.keys, key)
	m.values = append(m.values, value)
	m.len++
	return true, nil
}

func (m *MapValue) SetBySmallInt(k Int64, value Value) (bool, error) {
	key := Value{
		SumType:  "smallInt",
		SmallInt: &k,
	}
	return m.set(key, value)
}

func (m *MapValue) SetBySmallUInt(k UInt64, value Value) (bool, error) {
	key := Value{
		SumType:   "smallUint",
		SmallUint: &k,
	}
	return m.set(key, value)
}

func (m *MapValue) SetByBigInt(k BigInt, value Value) (bool, error) {
	key := Value{
		SumType: "bigInt",
		BigInt:  &k,
	}
	return m.set(key, value)
}

func (m *MapValue) SetByBigUInt(k BigUInt, value Value) (bool, error) {
	key := Value{
		SumType: "bigUint",
		BigUint: &k,
	}
	return m.set(key, value)
}

func (m *MapValue) SetByBits(k Bits, value Value) (bool, error) {
	key := Value{
		SumType: "bits",
		Bits:    &k,
	}
	return m.set(key, value)
}

func (m *MapValue) SetByInternalAddress(k InternalAddress, value Value) (bool, error) {
	key := Value{
		SumType:         "internalAddress",
		InternalAddress: &k,
	}
	return m.set(key, value)
}

func (m *MapValue) delete(key Value) {
	for i, k := range m.keys {
		if k.Equal(key) {
			m.keys[i] = Value{}
			m.values[i] = Value{}
			m.len--
		}
	}
}

func (m *MapValue) DeleteBySmallInt(k Int64) {
	key := Value{
		SumType:  "smallInt",
		SmallInt: &k,
	}
	m.delete(key)
}

func (m *MapValue) DeleteBySmallUInt(k UInt64) {
	key := Value{
		SumType:   "smallUint",
		SmallUint: &k,
	}
	m.delete(key)
}

func (m *MapValue) DeleteByBigInt(k BigInt) {
	key := Value{
		SumType: "bigInt",
		BigInt:  &k,
	}
	m.delete(key)
}

func (m *MapValue) DeleteByBigUInt(k BigUInt) {
	key := Value{
		SumType: "bigUint",
		BigUint: &k,
	}
	m.delete(key)
}

func (m *MapValue) DeleteByBits(k Bits) {
	key := Value{
		SumType: "bits",
		Bits:    &k,
	}
	m.delete(key)
}

func (m *MapValue) DeleteByInternalAddress(k InternalAddress) {
	key := Value{
		SumType:         "internalAddress",
		InternalAddress: &k,
	}
	m.delete(key)
}

func (m *MapValue) Len() int {
	return m.len
}

func (m *MapValue) MarshalJSON() ([]byte, error) {
	if len(m.keys) != len(m.values) {
		return nil, errors.New("map values and keys must contain equal length")
	}
	s := strings.Builder{}
	s.WriteString("{")
	if len(m.keys) > 0 {
		s.WriteString("\"keySumType\":")
		s.WriteString(fmt.Sprintf("\"%s\",", utils.ToCamelCasePrivate(string(m.keys[0].SumType))))
		for i, k := range m.keys {
			if i != 0 {
				s.WriteString(",")
			}
			key, err := json.Marshal(k)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map's key: %w", err)
			}
			validKey, err := getJustValue(string(key))
			if err != nil {
				return nil, fmt.Errorf("failed to get map's key value: %w", err)
			}
			// there is no case that key's value may start with '[' since tensors and tuples are not fixed sizes
			if validKey[0] == '{' {
				return nil, fmt.Errorf("cannot use %v as a map key", key)
			}
			if validKey[0] != '"' {
				validKey = fmt.Sprintf("\"%s\"", validKey)
			}
			s.WriteString(fmt.Sprintf("%v:", validKey))
			val, err := json.Marshal(m.values[i])
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map's value: %w", err)
			}
			s.Write(val)
		}
	}
	s.WriteString("}")

	return []byte(s.String()), nil
}

func getJustValue(key string) (string, error) {
	foundComma := false
	start := -1
	for i, v := range key {
		if v == ',' {
			foundComma = true
		}
		if v == ':' && foundComma {
			start = i
			break
		}
	}
	if start == -1 {
		return "", fmt.Errorf("invalid key: %v", key)
	}
	return strings.ReplaceAll(key[start+1:len(key)-1], " ", ""), nil
}

func (m *MapValue) UnmarshalJSON(bytes []byte) error {
	decoder := json.NewDecoder(strings.NewReader(string(bytes)))
	_, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to unmarshal map: %w", err)
	}

	keyTypeDeclr, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to unmarshal map key type declaration: %w", err)
	}
	stringKeyTypeDeclr, ok := keyTypeDeclr.(string)
	if !ok {
		return fmt.Errorf("expected key type as a string")
	}
	if stringKeyTypeDeclr != "keySumType" {
		return fmt.Errorf("map does not have key sum type")
	}

	keyType, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to unmarshal map key type: %w", err)
	}
	stringKeyType, ok := keyType.(string)
	if !ok {
		return fmt.Errorf("expected map key type value as a string")
	}
	keyTemplate := strings.Builder{}
	keyTemplate.WriteString("{\"sumType\":\"")
	keyTemplate.WriteString(stringKeyType)
	keyTemplate.WriteString("\",\"")
	keyTemplate.WriteString(stringKeyType)
	keyTemplate.WriteString("\":%s}")
	keyTmpl := keyTemplate.String()

	for decoder.More() {
		keyValue, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to unmarshal map key's value: %w", err)
		}
		stringKeyValue, ok := keyValue.(string)
		if !ok {
			return fmt.Errorf("expected map key as a string")
		}
		keyValueJson := fmt.Sprintf(keyTmpl, wrapValue(stringKeyType, stringKeyValue))
		key := Value{}
		if err = json.Unmarshal([]byte(keyValueJson), &key); err != nil {
			return fmt.Errorf("failed to unmarshal map key value: %w", err)
		}

		var value Value
		if err = decoder.Decode(&value); err != nil {
			return fmt.Errorf("failed to unmarshal map value: %w", err)
		}

		m.keys = append(m.keys, key)
		m.values = append(m.values, value)
		m.len++
	}

	return nil
}

func wrapValue(sumType string, v any) any {
	switch sumType {
	case "smallInt", "smallUint", "bool", "optionalValue", "tupleWith", "tensor", "map", "struct", "enum", "union", "refValue", "alias", "generic":
		return v
	}
	return fmt.Sprintf("\"%v\"", v)
}
