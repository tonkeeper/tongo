package tlb

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

type fixedSize interface {
	FixedSize() int
	Equal(other any) bool
	Compare(other any) (int, bool)
}

// HashmapItem represents a key-value pair stored in HashmapE[T].
type HashmapItem[keyT fixedSize, T any] struct {
	Key   keyT
	Value T
}

type Hashmap[keyT fixedSize, T any] struct {
	keys   []keyT
	values []T
}

// NewHashmap returns a new instance of Hashmap.
// Make sure that a key at index "i" corresponds to a value at the same index.
func NewHashmap[keyT fixedSize, T any](keys []keyT, values []T) Hashmap[keyT, T] {
	return Hashmap[keyT, T]{
		keys:   keys,
		values: values,
	}
}

func (h Hashmap[keyT, T]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	// Marshal empty Hashmap
	if len(h.values) == 0 || h.values == nil {
		return nil
	}
	var s keyT
	keys := make([]boc.BitString, 0, len(h.keys))
	for _, k := range h.keys {
		cell := boc.NewCell()
		err := Marshal(cell, k)
		if err != nil {
			return err
		}
		keys = append(keys, cell.RawBitString())
	}
	err := h.encodeMap(c, keys, h.values, s.FixedSize())
	if err != nil {
		return err
	}
	return nil
}

func (h Hashmap[keyT, T]) encodeMap(c *boc.Cell, keys []boc.BitString, values []T, keySize int) error {
	if len(keys) == 0 || len(values) == 0 {
		return fmt.Errorf("keys or values are empty")
	}
	label, err := encodeLabel(c, &keys[0], &keys[len(keys)-1], keySize)
	if err != nil {
		return err
	}
	keySize = keySize - label.BitsAvailableForRead() - 1 // l = n - m - 1 // see tlb
	var leftKeys, rightKeys []boc.BitString
	var leftValues, rightValues []T
	if len(keys) > 1 {
		for i := range keys {
			_, err := keys[i].ReadBits(label.BitsAvailableForRead()) // skip common label
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
		err = h.encodeMap(l, leftKeys, leftValues, keySize)
		if err != nil {
			return err
		}
		r, err := c.NewRef()
		if err != nil {
			return err
		}
		err = h.encodeMap(r, rightKeys, rightValues, keySize)
		if err != nil {
			return err
		}
		return err
	}
	// marshal value
	err = Marshal(c, values[0])
	if err != nil {
		return err
	}
	return nil
}

func (h *Hashmap[keyT, T]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var s keyT
	keySize := s.FixedSize()
	keyPrefix := boc.NewBitString(keySize)
	err := h.mapInner(keySize, keySize, c, &keyPrefix, decoder)
	if err != nil {
		return err
	}
	return nil
}

func hashmapAugExtraCountLeafs[keyT fixedSize](c *boc.Cell) (int, error) {
	maybe, err := c.ReadBit()
	if err != nil {
		return 0, err
	}
	if !maybe {
		return 0, nil
	}
	ref, err := c.NextRef()
	if err != nil {
		return 0, err
	}
	var s keyT
	keySize := s.FixedSize()
	return countLeafs(keySize, keySize, ref)
}

func countLeafs(keySize, leftKeySize int, c *boc.Cell) (int, error) {
	if c.CellType() == boc.PrunedBranchCell {
		return 0, errors.New("can't count leafs for hashmap with pruned branch cell")
	}
	size, err := loadLabelSize(leftKeySize, c)
	if err != nil {
		return 0, err
	}
	if keySize-leftKeySize+size < keySize {
		// 0 bit branch
		left, err := c.NextRef()
		if err != nil {
			return 0, err
		}
		leftCount, err := countLeafs(keySize, leftKeySize-(1+size), left)
		if err != nil {
			return 0, err
		}
		// 1 bit branch
		right, err := c.NextRef()
		if err != nil {
			return 0, err
		}
		rightCount, err := countLeafs(keySize, leftKeySize-(1+size), right)
		if err != nil {
			return 0, err
		}
		return leftCount + rightCount, nil
	}
	return 1, nil
}

func ProveKeyInHashmap[T any](prover *boc.MerkleProver, cell *boc.Cell, key boc.BitString) (T, []byte, error) {
	keySize := key.BitsAvailableForRead()
	bitString := boc.NewBitString(keySize)
	prefix := &bitString
	cursor := prover.Cursor()
	var t T
	remaining := keySize
	for {
		var err error
		var size int
		size, prefix, err = loadLabel(remaining, cell, prefix)
		if err != nil {
			return t, nil, err
		}
		_ = prefix
		if remaining <= size {
			break
		}
		if _, err = key.ReadBits(size); err != nil {
			return t, nil, err
		}
		isRight, err := key.ReadBit()
		if err != nil {
			return t, nil, err
		}
		if err := prefix.WriteBit(isRight); err != nil {
			return t, nil, err
		}
		remaining = remaining - size - 1
		next, err := cell.NextRef()
		if err != nil {
			return t, nil, err
		}
		if isRight {
			cursor.Ref(0).Prune()
			next, err = cell.NextRef()
			if err != nil {
				return t, nil, err
			}
			cursor = cursor.Ref(1)
		} else {
			cursor.Ref(1).Prune()
			cursor = cursor.Ref(0)
		}
		cell = next
		cell.ResetCounters()
	}
	if err := Unmarshal(cell, &t); err != nil {
		return t, nil, err
	}
	constructedKey, err := prefix.ReadBits(keySize)
	if err != nil {
		return t, nil, err
	}
	if constructedKey.ToFiftHex() != key.ToFiftHex() {
		return t, nil, errors.New("key is not found")
	}
	proof, err := prover.CreateProof(cursor)
	if err != nil {
		return t, nil, err
	}
	return t, proof, nil

}

func (h *Hashmap[keyT, T]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString, decoder *Decoder) error {
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
		err = h.mapInner(keySize, leftKeySize-(1+size), left, &lp, decoder)
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
		err = h.mapInner(keySize, leftKeySize-(1+size), right, &rp, decoder)
		if err != nil {
			return err
		}
		return nil
	}
	// add node to map
	var value T
	err = decoder.Unmarshal(c, &value)
	if err != nil {
		return err
	}
	h.values = append(h.values, value)
	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return err
	}

	var k keyT
	cell := boc.NewCellWithBits(key)
	err = decoder.Unmarshal(cell, &k)
	if err != nil {
		return err
	}
	h.keys = append(h.keys, k)
	return nil
}

func (h Hashmap[keyT, T]) Values() []T {
	return h.values
}

func (h Hashmap[keyT, T]) Keys() []keyT {
	return h.keys
}

func (h Hashmap[keyT, T]) Get(key keyT) (T, bool) {
	for i, k := range h.keys {
		if k.Equal(key) {
			return h.values[i], true
		}
	}
	return *new(T), false
}
func (h *Hashmap[keyT, T]) Put(key keyT, value T) {
	for i, k := range h.keys {
		if k.Equal(key) {
			h.values[i] = value
			return
		}
	}

	index := len(h.keys)
	for idx, other := range h.keys {
		cmp, ok := key.Compare(other)
		if !ok {
			// only happens when the user implements their own
			// fixedSize interface and intentionally returns false
			panic("key type mismatch")
		}
		if cmp < 0 {
			index = idx
			break
		}
	}

	h.keys = slices.Insert(h.keys, index, key)
	h.values = slices.Insert(h.values, index, value)
}

type HashmapE[keyT fixedSize, T any] struct {
	m Hashmap[keyT, T]
}

// NewHashmapE returns a new instance of HashmapE.
// Make sure that a key at index "i" corresponds to a value at the same index.
func NewHashmapE[keyT fixedSize, T any](keys []keyT, values []T) HashmapE[keyT, T] {
	return HashmapE[keyT, T]{
		m: Hashmap[keyT, T]{
			keys:   keys,
			values: values,
		},
	}
}

func (h HashmapE[keyT, T]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	var temp Maybe[Ref[Hashmap[keyT, T]]]
	temp.Exists = len(h.m.keys) > 0
	temp.Value.Value = h.m
	return encoder.Marshal(c, temp)
}

func (h *HashmapE[keyT, T]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var temp Maybe[Ref[Hashmap[keyT, T]]]
	err := decoder.Unmarshal(c, &temp)
	h.m = temp.Value.Value
	return err
}

func (h HashmapE[keyT, T]) Values() []T {
	return h.m.values
}

func (h HashmapE[keyT, T]) Keys() []keyT {
	return h.m.keys
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

type HashmapAug[keyT fixedSize, T1, T2 any] struct {
	keys   []keyT
	values []T1
	extra  HashMapAugExtraList[T2]
}

type HashMapAugExtraList[T any] struct {
	Left  *HashMapAugExtraList[T]
	Right *HashMapAugExtraList[T]
	Data  T
}

func (h HashmapAug[keyT, T1, T2]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return fmt.Errorf("not implemented")
}

func (h *HashmapAug[keyT, T1, T2]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var t keyT
	keySize := t.FixedSize()
	keyPrefix := boc.NewBitString(keySize)
	err := h.mapInner(keySize, keySize, c, &keyPrefix, &h.extra, decoder)
	if err != nil {
		return err
	}
	return nil
}

func (h *HashmapAug[keyT, T1, T2]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString, extras *HashMapAugExtraList[T2], decoder *Decoder) error {
	var err error
	var size int
	if c.CellType() == boc.PrunedBranchCell {
		return nil
	}
	size, keyPrefix, err = loadLabel(leftKeySize, c, keyPrefix)
	if err != nil {
		return err
	}
	var extra T2
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
		var extraLeft HashMapAugExtraList[T2]
		err = h.mapInner(keySize, leftKeySize-(1+size), left, &lp, &extraLeft, decoder)
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
		var extraRight HashMapAugExtraList[T2]
		err = h.mapInner(keySize, leftKeySize-(1+size), right, &rp, &extraRight, decoder)
		if err != nil {
			return err
		}
		extras.Left = &extraLeft
		extras.Right = &extraRight
		err = decoder.Unmarshal(c, &extra)
		if err != nil {
			return err
		}
		extras.Data = extra
		return nil
	}
	err = decoder.Unmarshal(c, &extra)
	if err != nil {
		return err
	}
	extras.Data = extra
	// add node to map
	var value T1
	err = decoder.Unmarshal(c, &value)
	if err != nil {
		return err
	}
	h.values = append(h.values, value)
	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return err
	}

	var k keyT
	cell := boc.NewCellWithBits(key)
	err = decoder.Unmarshal(cell, &k)
	if err != nil {
		return err
	}
	h.keys = append(h.keys, k)
	return nil
}

type HashmapAugE[keyT fixedSize, T1, T2 any] struct {
	m     HashmapAug[keyT, T1, T2]
	extra T2
}

func (h *HashmapAugE[keyT, T1, T2]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var temp struct {
		M     Maybe[Ref[HashmapAug[keyT, T1, T2]]]
		Extra T2
	}
	err := decoder.Unmarshal(c, &temp)
	h.m = temp.M.Value.Value
	h.extra = temp.Extra
	return err
}

func (h HashmapAugE[keyT, T1, T2]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	var temp struct {
		M     Maybe[Ref[HashmapAug[keyT, T1, T2]]]
		Extra T2
	}
	temp.M.Exists = len(h.m.keys) > 0
	temp.M.Value.Value = h.m
	temp.Extra = h.extra
	return Marshal(c, temp)
}

func (h HashmapAugE[keyT, T1, T2]) Values() []T1 {
	return h.m.values
}

func (h HashmapAugE[keyT, T1, T2]) Keys() []keyT {
	return h.m.keys
}

func loadLabelSize(size int, c *boc.Cell) (int, error) {
	first, err := c.ReadBit()
	if err != nil {
		return 0, err
	}
	// hml_short$0
	if !first {
		// Unary, while 1, add to ln
		ln, err := c.ReadUnary()
		if err != nil {
			return 0, err
		}
		return int(ln), nil
	}
	second, err := c.ReadBit()
	if err != nil {
		return 0, err
	}
	// hml_long$10
	if !second {
		ln, err := c.ReadLimUint(size)
		if err != nil {
			return 0, err
		}
		return int(ln), nil
	}
	// hml_same$11
	_, err = c.ReadBit()
	if err != nil {
		return 0, err
	}
	ln, err := c.ReadLimUint(size)
	if err != nil {
		return 0, err
	}
	return int(ln), nil
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

// Values returns a list of value of this hashmap.
func (h HashmapAug[_, T1, _]) Values() []T1 {
	return h.values
}

// Keys returns a list of key of this hashmap.
func (h HashmapAug[keyT, _, _]) Keys() []keyT {
	return h.keys
}

// Items returns key-value pairs of this hashmap.
func (h HashmapE[keyT, T]) Items() []HashmapItem[keyT, T] {
	return h.m.Items()
}

func (h Hashmap[keyT, T]) Items() []HashmapItem[keyT, T] {
	items := make([]HashmapItem[keyT, T], len(h.keys))
	for i, key := range h.keys {
		items[i] = HashmapItem[keyT, T]{
			Key:   key,
			Value: h.values[i],
		}
	}
	return items
}

func (h HashmapE[keyT, T]) Get(key keyT) (T, bool) {
	return h.m.Get(key)
}

func (h *HashmapE[keyT, T]) Put(key keyT, value T) {
	h.m.Put(key, value)
}

func (h Hashmap[keyT, T]) MarshalJSON() ([]byte, error) {
	m := make(map[string]T, len(h.Keys()))
	for _, item := range h.Items() {
		key, err := json.Marshal(item.Key)
		if err != nil {
			return nil, err
		}
		m[strings.Trim(string(key), "\"")] = item.Value
	}
	return json.Marshal(m)
}

func (h *Hashmap[keyT, T]) UnmarshalJSON(data []byte) error {
	var m map[string]T
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	h.keys = make([]keyT, 0, len(m))
	h.values = make([]T, 0, len(m))
	for k, v := range m {
		var key keyT
		if err := json.Unmarshal([]byte(strconv.Quote(k)), &key); err != nil {
			return err
		}
		h.keys = append(h.keys, key)
		h.values = append(h.values, v)
	}
	return nil
}
