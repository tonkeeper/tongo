package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type TolkType interface {
	UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error
}

type TolkComparableType interface {
	GetFixedSize() int
}

type Map struct {
	K Ty `json:"k"`
	V Ty `json:"v"`
}

func (Map) SetValue(v *Value, val any) error {
	m, ok := val.(MapValue)
	if !ok {
		return fmt.Errorf("value is not a map")
	}
	v.mp = &m
	v.sumType = "mp"
	return nil
}

func (m Map) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	cmp, ok := m.K.GetComparableType()
	if !ok {
		return fmt.Errorf("%v is not comparable", m.K.SumType)
	}
	keySize := cmp.GetFixedSize()
	keyPrefix := boc.NewBitString(cmp.GetFixedSize())

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
		err = mapInner(keySize, keySize, mpCell, &keyPrefix, m.K, m.V, &mp.keys, &mp.values, decoder)
		if err != nil {
			return err
		}
	}

	mp.len = len(mp.keys)
	err = m.SetValue(v, mp)
	if err != nil {
		return err
	}

	return nil
}

func (m *MapValue) UnmarshalTolk(cell *boc.Cell, ty Map, decoder *Decoder) error {
	cmp, ok := ty.K.GetComparableType()
	if !ok {
		return fmt.Errorf("%v is not comparable", ty.K.SumType)
	}
	keySize := cmp.GetFixedSize()
	keyPrefix := boc.NewBitString(cmp.GetFixedSize())

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
	kt, vt Ty,
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
	err = v.UnmarshalTolk(c, vt, decoder)
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
	err = k.UnmarshalTolk(cell, kt, decoder)
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

func (Map) MarshalTolk(cell *boc.Cell, v *Value) error {
	return nil
}

//	if v.mp == nil {
//		return fmt.Errorf("map is nil")
//	}
//
//	compKey, ok := v.mp.keyType.GetComparableType()
//	if !ok {
//		return fmt.Errorf("map key is not a comparable type, got %v", v.mp.keyType.SumType)
//	}
//
//	if len(v.mp.values) == 0 {
//		err := cell.WriteBit(false)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	err := cell.WriteBit(true)
//	if err != nil {
//		return err
//	}
//
//	keys := make([]boc.BitString, 0, len(v.mp.keys))
//	for _, k := range v.mp.keys {
//		keyCell := boc.NewCell()
//		err = k.valType.MarshalTolk(keyCell, &k)
//		if err != nil {
//			return err
//		}
//		keys = append(keys, keyCell.RawBitString())
//	}
//
//	ref := boc.NewCell()
//	err = encodeMap(ref, keys, v.mp.values, compKey.GetFixedSize(), v.mp.valType)
//	if err != nil {
//		return err
//	}
//
//	err = cell.AddRef(ref)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func encodeMap(c *boc.Cell, keys []boc.BitString, values []Value, keySize int, vt Ty) error {
//	if len(keys) == 0 || len(values) == 0 {
//		return fmt.Errorf("keys or values are empty")
//	}
//	label, err := encodeLabel(c, &keys[0], &keys[len(keys)-1], keySize)
//	if err != nil {
//		return err
//	}
//	keySize = keySize - label.BitsAvailableForRead() - 1 // l = n - m - 1 // see tlb
//	var leftKeys, rightKeys []boc.BitString
//	var leftValues, rightValues []Value
//	if len(keys) > 1 {
//		for i := range keys {
//			_, err := keys[i].ReadBits(label.BitsAvailableForRead()) // skip common label
//			if err != nil {
//				return err
//			}
//			isRight, err := keys[i].ReadBit()
//			if err != nil {
//				return err
//			}
//			if isRight {
//				rightKeys = append(rightKeys, keys[i].ReadRemainingBits())
//				rightValues = append(rightValues, values[i])
//			} else {
//				leftKeys = append(leftKeys, keys[i].ReadRemainingBits())
//				leftValues = append(leftValues, values[i])
//			}
//		}
//		l, err := c.NewRef()
//		if err != nil {
//			return err
//		}
//		err = encodeMap(l, leftKeys, leftValues, keySize, vt)
//		if err != nil {
//			return err
//		}
//		r, err := c.NewRef()
//		if err != nil {
//			return err
//		}
//		err = encodeMap(r, rightKeys, rightValues, keySize, vt)
//		if err != nil {
//			return err
//		}
//		return err
//	}
//	// marshal value
//	err = vt.MarshalTolk(c, &values[0])
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func encodeLabel(c *boc.Cell, keyFirst, keyLast *boc.BitString, keySize int) (boc.BitString, error) {
//	label := boc.NewBitString(keySize)
//	if keyFirst != keyLast {
//		bitLeft, err := keyFirst.ReadBit()
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		for keyFirst.BitsAvailableForRead() > 0 {
//			bitRight, err := keyLast.ReadBit()
//			if err != nil {
//				return boc.BitString{}, err
//			}
//			if bitLeft != bitRight {
//				break
//			}
//			if err := label.WriteBit(bitLeft); err != nil {
//				return boc.BitString{}, err
//			}
//			bitLeft, err = keyFirst.ReadBit()
//			if err != nil {
//				return boc.BitString{}, err
//			}
//		}
//	} else {
//		label = keyFirst.Copy()
//	}
//	keyFirst.ResetCounter()
//	keyLast.ResetCounter()
//	if label.BitsAvailableForRead() < 8 {
//		//hml_short$0 {m:#} {n:#} len:(Unary ~n) {n <= m} s:(n * Bit) = HmLabel ~n m;
//		err := c.WriteBit(false)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		// todo pack label
//		err = c.WriteUnary(uint(label.BitsAvailableForRead()))
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		err = c.WriteBitString(label)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//
//	} else {
//		// hml_long$10 {m:#} n:(#<= m) s:(n * Bit) = HmLabel ~n m;
//		err := c.WriteBit(true)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		err = c.WriteBit(false)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		// todo pack label
//		err = c.WriteLimUint(label.BitsAvailableForRead(), keySize)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//		err = c.WriteBitString(label)
//		if err != nil {
//			return boc.BitString{}, err
//		}
//	}
//	return label, nil
//}

func (Map) Equal(v Value, o Value) bool {
	return false
}
