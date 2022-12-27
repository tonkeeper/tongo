package tlb

import (
	"fmt"

	"github.com/startfellows/tongo/boc"
)

type fixedSize interface {
	FixedSize() int
}

// HashmapItem represents a key-value pair stored in HashmapE[T].
type HashmapItem[T any] struct {
	Key   boc.BitString
	Value T
}

type Hashmap[sizeT fixedSize, T any] struct {
	keys   []boc.BitString
	values []T
}

func (h Hashmap[sizeT, T]) MarshalTLB(c *boc.Cell, tag string) error {
	// Marshal empty Hashmap
	if len(h.values) == 0 || h.values == nil {
		return nil
	}
	var s sizeT
	err := h.encodeMap(c, h.keys, h.values, s.FixedSize())
	if err != nil {
		return err
	}
	return nil
}

func (h Hashmap[sizeT, T]) encodeMap(c *boc.Cell, keys []boc.BitString, values []T, size int) error {
	if len(keys) == 0 || len(values) == 0 {
		return fmt.Errorf("keys or values are empty")
	}

	label, err := encodeLabel(c, &keys[0], &keys[len(keys)-1], size)
	if err != nil {
		return err
	}

	size = size - label.BitsAvailableForRead() - 1 // l = n - m - 1 // see tlb
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
		err = h.encodeMap(l, leftKeys, leftValues, size)
		if err != nil {
			return err
		}
		r, err := c.NewRef()
		if err != nil {
			return err
		}
		err = h.encodeMap(r, rightKeys, rightValues, size)
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

func (h *Hashmap[sizeT, T]) UnmarshalTLB(c *boc.Cell, tag string) error {
	var s sizeT
	keySize := s.FixedSize()
	keyPrefix := boc.NewBitString(keySize)
	err := h.mapInner(keySize, keySize, c, &keyPrefix)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hashmap[sizeT, T]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString) error {
	var err error
	var size int
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
		err = h.mapInner(keySize, leftKeySize-(1+size), left, &lp)
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
		err = h.mapInner(keySize, leftKeySize-(1+size), right, &rp)
		if err != nil {
			return err
		}
		return nil
	}
	// add node to map
	var value T
	err = Unmarshal(c, &value)
	if err != nil {
		return err
	}
	h.values = append(h.values, value)
	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return err
	}
	h.keys = append(h.keys, key)
	return nil
}

func (h Hashmap[sizeT, T]) Values() []T {
	return h.values
}

func (h Hashmap[sizeT, T]) Keys() []boc.BitString {
	return h.keys
}

type HashmapE[sizeT fixedSize, T any] struct {
	m Hashmap[sizeT, T]
}

func (h HashmapE[sizeT, T]) MarshalTLB(c *boc.Cell, tag string) error {
	var temp Maybe[Ref[Hashmap[sizeT, T]]]
	temp.Null = len(h.m.keys) == 0
	temp.Value.Value = h.m
	return Marshal(c, temp)
}

func (h *HashmapE[sizeT, T]) UnmarshalTLB(c *boc.Cell, tag string) error {
	var temp Maybe[Ref[Hashmap[sizeT, T]]]
	err := Unmarshal(c, &temp)
	h.m = temp.Value.Value
	return err
}

func (h HashmapE[sizeT, T]) Values() []T {
	return h.m.values
}

func (h HashmapE[sizeT, T]) Keys() []boc.BitString {
	return h.m.keys
}
func encodeLabel(c *boc.Cell, keyFirst, keyLast *boc.BitString, size int) (boc.BitString, error) {
	label := boc.NewBitString(size)
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
			label.WriteBit(bitLeft)
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
	if label.BitsAvailableForRead() < 8 {
		//hml_short$0 {m:#} {n:#} len:(Unary ~n) {n <= m} s:(n * Bit) = HmLabel ~n m;
		err := c.WriteBit(false)
		if err != nil {
			return boc.BitString{}, err
		}
		// todo pack label
		err = c.WriteUnary(uint(label.BitsAvailableForRead()))
		if err != nil {
			return boc.BitString{}, err
		}
		err = c.WriteBitString(label)
		if err != nil {
			return boc.BitString{}, err
		}

	} else {
		// hml_long$10 {m:#} n:(#<= m) s:(n * Bit) = HmLabel ~n m;
		err := c.WriteBit(true)
		if err != nil {
			return boc.BitString{}, err
		}
		err = c.WriteBit(false)
		if err != nil {
			return boc.BitString{}, err
		}
		// todo pack label
		err = c.WriteLimUint(label.BitsAvailableForRead(), size)
		if err != nil {
			return boc.BitString{}, err
		}
		err = c.WriteBitString(label)
		if err != nil {
			return boc.BitString{}, err
		}
	}
	return label, nil
}

type HashmapAug[sizeT fixedSize, T1, T2 any] struct {
	keys   []boc.BitString
	values []T1
	extra  HashMapAugExtraList[T2]
}

type HashMapAugExtraList[T any] struct {
	Left  *HashMapAugExtraList[T]
	Right *HashMapAugExtraList[T]
	Data  T
}

func (h HashmapAug[sizeT, T1, T2]) MarshalTLB(c *boc.Cell, tag string) error {
	return fmt.Errorf("not implemented")
}

func (h *HashmapAug[sizeT, T1, T2]) UnmarshalTLB(c *boc.Cell, tag string) error {
	var t sizeT
	keySize := t.FixedSize()
	keyPrefix := boc.NewBitString(keySize)
	err := h.mapInner(keySize, keySize, c, &keyPrefix, &h.extra)
	if err != nil {
		return err
	}
	return nil
}

func (h *HashmapAug[sizeT, T1, T2]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString, extras *HashMapAugExtraList[T2]) error {
	var err error
	var size int
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
		err = h.mapInner(keySize, leftKeySize-(1+size), left, &lp, &extraLeft)
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
		err = h.mapInner(keySize, leftKeySize-(1+size), right, &rp, &extraRight)
		if err != nil {
			return err
		}
		extras.Left = &extraLeft
		extras.Right = &extraRight
		err = Unmarshal(c, &extra)
		if err != nil {
			return err
		}
		extras.Data = extra
		return nil
	}
	err = Unmarshal(c, &extra)
	if err != nil {
		return err
	}
	extras.Data = extra
	// add node to map
	var value T1
	err = Unmarshal(c, &value)
	if err != nil {
		return err
	}
	h.values = append(h.values, value)
	key, err := keyPrefix.ReadBits(keySize)
	if err != nil {
		return err
	}
	h.keys = append(h.keys, key)

	return nil
}

type HashmapAugE[sizeT fixedSize, T1, T2 any] struct {
	m     HashmapAug[sizeT, T1, T2]
	extra T2
}

func (h *HashmapAugE[sizeT, T1, T2]) UnmarshalTLB(c *boc.Cell, tag string) error {
	var temp struct {
		M     Maybe[Ref[HashmapAug[sizeT, T1, T2]]]
		Extra T2
	}
	err := Unmarshal(c, &temp)
	h.m = temp.M.Value.Value
	h.extra = temp.Extra
	return err
}

func (h HashmapAugE[sizeT, T1, T2]) MarshalTLB(c *boc.Cell, tag string) error {
	var temp struct {
		M     Maybe[Ref[HashmapAug[sizeT, T1, T2]]]
		Extra T2
	}
	temp.M.Null = len(h.m.keys) == 0
	temp.M.Value.Value = h.m
	temp.Extra = h.extra
	return Marshal(c, temp)
}

func (h HashmapAugE[sizeT, T1, T2]) Values() []T1 {
	return h.m.values
}

func (h HashmapAugE[sizeT, T1, T2]) Keys() []boc.BitString {
	return h.m.keys
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
		keyBits, err := c.ReadBits(int(ln))
		if err != nil {
			return 0, nil, err
		}
		// add bits to key
		err = key.WriteBitString(keyBits)
		if err != nil {
			return 0, nil, err
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
		keyBits, err := c.ReadBits(int(ln))
		if err != nil {
			return 0, nil, err
		}
		// add bits to key
		err = key.WriteBitString(keyBits)
		if err != nil {
			return 0, nil, err
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

// Items returns key-value pairs of this hashmap.
func (h HashmapE[_, T]) Items() []HashmapItem[T] {
	return h.m.Items()
}

func (h Hashmap[_, T]) Items() []HashmapItem[T] {
	items := make([]HashmapItem[T], len(h.keys))
	for i, key := range h.keys {
		items[i] = HashmapItem[T]{
			Key:   key,
			Value: h.values[i],
		}
	}
	return items
}
