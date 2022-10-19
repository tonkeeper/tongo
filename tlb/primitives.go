package tlb

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/startfellows/tongo/boc"
)

type SumType string

type Magic uint32

func (m *Magic) UnmarshalTLB(c *boc.Cell, tag string) error {
	a := strings.Split(tag, "$")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 2, 32)
		if err != nil {
			return err
		}
		y, err := c.PickUint(len(a[1]))
		if x == y {
			_ = c.Skip(len(a[1])) // already checked
			return nil
		}
		return fmt.Errorf("magic prefix: %v not found ", tag)
	}
	a = strings.Split(tag, "#")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 16, 32)
		if err != nil {
			return err
		}
		y, err := c.PickUint(len(a[1]) * 4)
		if x == y {
			_ = c.Skip(len(a[1]) * 4) // already checked
			return nil
		}
		return fmt.Errorf("magic prefix: %v not found ", tag)
	}

	return fmt.Errorf("unsupported tag: %v", tag)
}

type Maybe[T any] struct {
	Null  bool
	Value T
}

func (m Maybe[_]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(!m.Null)
	if err != nil {
		return err
	}
	if !m.Null {
		err = Marshal(c, m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Maybe[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	exist, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.Null = !exist
	if exist {
		err = Unmarshal(c, &m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

type Either[M, N any] struct {
	IsRight bool
	Left    M
	Right   N
}

func (m Either[_, _]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(m.IsRight)
	if err != nil {
		return err
	}
	if m.IsRight {
		err = Marshal(c, m.Right)
		if err != nil {
			return err
		}
	} else {
		err = Marshal(c, m.Left)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Either[_, _]) UnmarshalTLB(c *boc.Cell, tag string) error {
	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.IsRight = isRight
	if isRight {
		err = Unmarshal(c, &m.Right)
		if err != nil {
			return err
		}
	} else {
		err = Unmarshal(c, &m.Left)
		if err != nil {
			return err
		}
	}
	return nil
}

type EitherRef[T any] struct {
	IsRight bool
	Value   T
}

func (m EitherRef[_]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(m.IsRight)
	if err != nil {
		return err
	}
	if m.IsRight {
		c, err = c.NewRef()
		if err != nil {
			return err
		}
	}
	return Marshal(c, m.Value)
}

func (m *EitherRef[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.IsRight = isRight
	if isRight {
		c, err = c.NextRef()
		if err != nil {
			return err
		}
	}
	return Unmarshal(c, &m.Value)
}

type Ref[T any] struct {
	Value T
}

func (m Ref[_]) MarshalTLB(c *boc.Cell, tag string) error {
	r := boc.NewCell()
	err := Marshal(r, m.Value)
	if err != nil {
		return err
	}
	err = c.AddRef(r)
	return err
}

func (m *Ref[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	r, err := c.NextRef()
	if err != nil {
		return err
	}
	err = Unmarshal(r, &m.Value)
	if err != nil {
		return err
	}
	return nil
}

type Unary uint

func (n Unary) MarshalTLB(c *boc.Cell, tag string) error {
	return c.WriteUnary(uint(n))
}

func (n *Unary) UnmarshalTLB(c *boc.Cell, tag string) error {
	a, err := c.ReadUnary()
	*n = Unary(a)
	return err
}

type Hashmap[T any] struct {
	keys    []boc.BitString
	keySize int
	values  []T
}

func (h Hashmap[T]) MarshalTLB(c *boc.Cell, tag string) error {
	// Marshal empty Hashmap
	if len(h.values) == 0 || h.values == nil {
		return nil
	}
	err := h.encodeMap(c, h.keys, h.values, h.keySize)
	if err != nil {
		return err
	}
	return nil
}

func (h Hashmap[T]) encodeMap(c *boc.Cell, keys []boc.BitString, values []T, size int) error {
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
			return nil
		}
		err = h.encodeMap(l, leftKeys, leftValues, size)
		if err != nil {
			return err
		}
		r, err := c.NewRef()
		if err != nil {
			return nil
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

func (h *Hashmap[T]) UnmarshalTLB(c *boc.Cell, tag string) error {
	keySize, err := decodeHashmapTag(tag)
	if err != nil {
		return err
	}
	h.keySize = keySize
	keyPrefix := boc.NewBitString(keySize)
	err = h.mapInner(keySize, keySize, c, &keyPrefix)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hashmap[T]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString) error {
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
			return nil
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

func (h Hashmap[T]) Values() []T {
	return h.values
}

func (h Hashmap[T]) Keys() []boc.BitString {
	return h.keys
}

type HashmapE[T any] struct {
	keys    []boc.BitString
	keySize int
	values  []T
}

func (h HashmapE[T]) MarshalTLB(c *boc.Cell, tag string) error {
	// Marshal empty Hashmap
	if len(h.values) == 0 || h.values == nil {
		err := c.WriteBit(false)
		if err != nil {
			return err
		}
		return nil
	}
	err := c.WriteBit(true)
	if err != nil {
		return err
	}
	r, err := c.NewRef()
	if err != nil {
		return err
	}
	err = h.encodeMap(r, h.keys, h.values, h.keySize)
	if err != nil {
		return err
	}
	return nil
}

func (h HashmapE[T]) encodeMap(c *boc.Cell, keys []boc.BitString, values []T, size int) error {
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
			return nil
		}
		err = h.encodeMap(l, leftKeys, leftValues, size)
		if err != nil {
			return err
		}
		r, err := c.NewRef()
		if err != nil {
			return nil
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

func (h *HashmapE[T]) UnmarshalTLB(c *boc.Cell, tag string) error {
	keySize, err := decodeHashmapTag(tag)
	if err != nil {
		return err
	}
	h.keySize = keySize
	isExists, err := c.ReadBit()
	if !isExists {
		return nil
	} // hme_empty$0 {n:#} {X:Type} = HashmapE n X;
	r, err := c.NextRef() // hme_root$1 {n:#} {X:Type} root:^(Hashmap n X) = HashmapE n X;
	if err != nil {
		return err
	}
	keyPrefix := boc.NewBitString(keySize)
	err = h.mapInner(keySize, keySize, r, &keyPrefix)
	if err != nil {
		return err
	}
	return nil
}

func decodeHashmapTag(tag string) (int, error) {
	var ln int
	if tag == "" {
		return 0, fmt.Errorf("empty hashmap tag")
	}
	_, err := fmt.Sscanf(tag, "%dbits", &ln)
	if err != nil {
		return 0, err
	}
	return ln, nil
}

func (h *HashmapE[T]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString) error {
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
			return nil
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

func (h HashmapE[T]) Values() []T {
	return h.values
}

func (h HashmapE[T]) Keys() []boc.BitString {
	return h.keys
}

type HashmapAug[T1, T2 any] struct {
	keys      []boc.BitString
	keySize   int
	values    []T1
	extra     HashMapAugExtraList[T2]
	rootExtra T2
}

type HashMapAugExtraList[T any] struct {
	Left  *HashMapAugExtraList[T]
	Right *HashMapAugExtraList[T]
	Data  T
}

func (h *HashmapAug[T1, T2]) UnmarshalTLB(c *boc.Cell, tag string) error {
	keySize, err := decodeHashmapTag(tag)
	if err != nil {
		return err
	}
	h.keySize = keySize

	keyPrefix := boc.NewBitString(keySize)
	err = h.mapInner(keySize, keySize, c, &keyPrefix, &h.extra)
	if err != nil {
		return err
	}

	return nil
}

func (h *HashmapAug[T1, T2]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString, extras *HashMapAugExtraList[T2]) error {
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
			return nil
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

type HashmapAugE[T1, T2 any] struct {
	keys      []boc.BitString
	keySize   int
	values    []T1
	extra     HashMapAugExtraList[T2]
	rootExtra T2
}

func (h *HashmapAugE[T1, T2]) UnmarshalTLB(c *boc.Cell, tag string) error {
	keySize, err := decodeHashmapTag(tag)
	if err != nil {
		return err
	}
	h.keySize = keySize
	isExists, err := c.ReadBit()
	if err != nil {
		return err
	}
	// hme_empty$0 {n:#} {X:Type} {Y:Type} extra:Y = HashmapAugE n X Y;
	// ahme_root$1 {n:#} {X:Type} {Y:Type} root:^(HashmapAug n X Y) extra:Y = HashmapAugE n X Y;
	if !isExists {
		return nil
	}
	r, err := c.NextRef()
	if err != nil {
		return err
	}
	keyPrefix := boc.NewBitString(keySize)
	err = h.mapInner(keySize, keySize, r, &keyPrefix, &h.extra)
	if err != nil {
		return err
	}
	var extra T2
	err = Unmarshal(c, &extra)
	if err != nil {
		return err
	}
	h.rootExtra = extra

	return nil
}

func (h *HashmapAugE[T1, T2]) mapInner(keySize, leftKeySize int, c *boc.Cell, keyPrefix *boc.BitString, extras *HashMapAugExtraList[T2]) error {
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
			return nil
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
	// add node to map
	err = Unmarshal(c, &extra)
	if err != nil {
		return err
	}
	extras.Data = extra
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

func (h HashmapAugE[T1, T2]) Values() []T1 {
	return h.values
}

func (h HashmapAugE[T1, T2]) Keys() []boc.BitString {
	return h.keys
}

type VarUInteger big.Int

func (u VarUInteger) MarshalTLB(c *boc.Cell, tag string) error {
	n, err := decodeVarUIntegerTag(tag)
	if n < 1 {
		return fmt.Errorf("len of varuint must be at least one byte")
	}
	if err != nil {
		return err
	}
	i := big.Int(u)
	b := i.Bytes()
	err = c.WriteLimUint(len(b), n-1)
	if err != nil {
		return err
	}
	err = c.WriteBytes(b)
	if err != nil {
		return err
	}
	return nil
}

func (u *VarUInteger) UnmarshalTLB(c *boc.Cell, tag string) error {
	n, err := decodeVarUIntegerTag(tag)
	if err != nil {
		return err
	}
	ln, err := c.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger(*val)
	return nil
}

func decodeVarUIntegerTag(tag string) (int, error) {
	var n int
	if tag == "" {
		return 0, fmt.Errorf("empty varuint tag")
	}
	_, err := fmt.Sscanf(tag, "%dbytes", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// TODO: replace bitstring with Cell
type Any boc.Cell

func (a Any) MarshalTLB(c *boc.Cell, tag string) error {
	x := boc.Cell(a)
	y := &x
	err := c.WriteBitString(y.RawBitString())
	if err != nil {
		return err
	}
	for y.RefsAvailableForRead() > 0 {
		ref, err := y.NextRef()
		if err != nil {
			return err
		}
		err = c.AddRef(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Any) UnmarshalTLB(c *boc.Cell, tag string) error {
	x := boc.NewCell()
	err := x.WriteBitString(c.ReadRemainingBits())
	if err != nil {
		return err
	}
	for c.RefsAvailableForRead() > 0 {
		ref, err := c.NextRef()
		if err != nil {
			return err
		}
		err = x.AddRef(ref)
		if err != nil {
			return err
		}
	}
	*a = Any(*x)
	return nil
}

type BinTree[T any] struct {
	Values []T
}

func decodeRecursiveBinTree(c *boc.Cell) ([]*boc.Cell, error) {
	var cellAr []*boc.Cell
	isBranch, err := c.ReadBit()
	if err != nil {
		return nil, err
	}
	if !isBranch {
		return append(cellAr, c), nil
	}

	l, err := c.NextRef()
	if err != nil {
		return nil, err
	}
	rec, err := decodeRecursiveBinTree(l)
	if err != nil {
		return nil, err
	}
	cellAr = append(cellAr, rec...)
	r, err := c.NextRef()
	if err != nil {
		return nil, err
	}
	rec, err = decodeRecursiveBinTree(r)
	if err != nil {
		return nil, err
	}
	cellAr = append(cellAr, rec...)

	return cellAr, nil
}

func (b BinTree[T]) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("BinTree marshaling not implmented")
}

func (b *BinTree[T]) UnmarshalTLB(c *boc.Cell, tag string) error {
	dec, err := decodeRecursiveBinTree(c)
	if err != nil {
		return err
	}
	for _, i := range dec {
		var t T
		err := Unmarshal(i, &t)
		if err != nil {
			return err
		}
		b.Values = append(b.Values, t)
	}
	return nil
}
