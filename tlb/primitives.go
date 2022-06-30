package tlb

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"math/big"
)

type SumType string

type tag struct {
	Name string
	Len  int
	Val  uint64
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
	// TODO: implement
	return fmt.Errorf("not empty HashmapE marshaling not implmented")
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

type VarUInteger big.Int

func (u VarUInteger) MarshalTLB(c *boc.Cell, tag string) error {
	n, err := decodeVarUIntegerTag(tag)
	if err != nil {
		return err
	}
	i := big.Int(u)
	b := i.Bytes()
	err = c.WriteLimUint(len(b), n)
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
