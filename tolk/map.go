package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type TolkType interface {
	UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error
}

type TolkComparableType interface {
	UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error
	GetFixedSize() int
}

type TolkComparableValue interface {
	Equal(other any) bool
	Compare(other any) (int, bool)
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
	return nil
}

func (m Map) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
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
		keyType: m.K,
		valType: m.V,
		keys:    make([]Value, 0),
		values:  make([]Value, 0),
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
	err = v.SetValue(mp)
	if err != nil {
		return err
	}

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
	err = vt.UnmarshalTolk(c, &v, decoder)
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
	err = kt.UnmarshalTolk(cell, &k, decoder)
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

func (Map) Equal(v Value, o Value) bool {
	return false
}
