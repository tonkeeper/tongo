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
	v.sumType = "union"
	return nil
}

func (u Union) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
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
			err = innerV.UnmarshalTolk(cell, variant.VariantTy, decoder)
			if err != nil {
				return err
			}
			unionV.Val = innerV
			err = u.SetValue(v, unionV)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("none of union prefixes matched")
}

func (u *UnionValue) UnmarshalTolk(cell *boc.Cell, ty Union, decoder *Decoder) error {
	unionV := UnionValue{}
	if len(ty.Variants) < 2 {
		return fmt.Errorf("union length must be at least 2")
	}
	prefixLen := ty.Variants[0].PrefixLen
	eatPrefix := ty.Variants[0].PrefixEatInPlace
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

	for _, variant := range ty.Variants {
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
			err = innerV.UnmarshalTolk(cell, variant.VariantTy, decoder)
			if err != nil {
				return err
			}
			unionV.Val = innerV
			*u = unionV

			return nil
		}
	}

	return fmt.Errorf("none of union prefixes matched")
}

func (u Union) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.union == nil {
	//	return fmt.Errorf("union is nil")
	//}
	//if len(u.Variants) < 2 {
	//	return fmt.Errorf("union length must be at least 2")
	//}
	//
	//if u.Variants[0].PrefixEatInPlace {
	//	err := cell.WriteUint(v.union.Prefix.Prefix, int(v.union.Prefix.Len))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//val := v.union.Val
	//err := val.valType.MarshalTolk(cell, &val)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (Union) Equal(v Value, o Value) bool {
	return false
}
