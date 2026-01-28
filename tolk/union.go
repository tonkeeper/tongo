package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type Union struct {
	Variants []UnionVariant `json:"variants"`
}

func (Union) SetValue(v *Value, val any) error {
	u, ok := val.(UnionValue)
	if !ok {
		return fmt.Errorf("value is not an union")
	}
	v.union = &u
	return nil
}

func (u Union) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	unionV := UnionValue{}
	if len(u.Variants) < 2 {
		return fmt.Errorf("union length must be at least 2")
	}
	prefixLen := u.Variants[0].PrefixLen
	eatPrefix := u.Variants[0].PrefixEatInPlace
	if prefixLen > 64 {
		// todo: maybe prefix len can be bigger than 64?
		return fmt.Errorf("union prefix length must be less than 64")
	}

	var prefix uint64
	var err error
	if !eatPrefix {
		copyCell := cell.CopyRemaining()
		prefix, err = copyCell.ReadUint(prefixLen)
		if err != nil {
			return err
		}
	} else {
		prefix, err = cell.ReadUint(prefixLen)
		if err != nil {
			return err
		}
	}

	for _, variant := range u.Variants {
		variantPrefix, err := PrefixToUint(variant.PrefixStr)
		if err != nil {
			return err
		}

		if prefix == variantPrefix {
			unionV.Prefix = TolkPrefix{
				Len:    int16(variant.PrefixLen),
				Prefix: prefix,
			}
			innerV := Value{}
			err = variant.VariantTy.UnmarshalTolk(cell, &innerV, decoder)
			if err != nil {
				return err
			}
			unionV.Val = innerV
			err = v.SetValue(unionV)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("none of union prefixes matched")
}

func (Union) Equal(v Value, o Value) bool {
	return false
}
