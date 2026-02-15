package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type UnionValue struct {
	Prefix Prefix `json:"prefix"`
	Val    Value  `json:"value"`
}

func (u *UnionValue) Unmarshal(cell *boc.Cell, ty tolkParser.Union, decoder *Decoder) error {
	unionV := UnionValue{}
	if len(ty.Variants) < 2 {
		return fmt.Errorf("union length must be at least 2")
	}
	prefixLen := ty.Variants[0].PrefixLen
	eatPrefix := ty.Variants[0].PrefixEatInPlace
	if prefixLen > 64 {
		// todo: maybe prefix len can be bigger than 64?
		return fmt.Errorf("union prefix length must be less than or equal to 64")
	}

	var prefix uint64
	var err error
	if !eatPrefix {
		copyCell := cell.CopyRemaining()
		prefix, err = copyCell.ReadUint(prefixLen)
	} else {
		prefix, err = cell.ReadUint(prefixLen)
	}
	if err != nil {
		return fmt.Errorf("failed to read union's %v-length prefix: %w", prefixLen, err)
	}

	for _, variant := range ty.Variants {
		variantPrefix, err := PrefixToUint(variant.PrefixStr)
		if err != nil {
			return fmt.Errorf("failed to read union's variant prefi %v: %w", variant.PrefixStr, err)
		}

		if prefix == variantPrefix {
			unionV.Prefix = Prefix{
				Len:    int16(variant.PrefixLen),
				Prefix: prefix,
			}
			innerV := Value{}
			err = innerV.Unmarshal(cell, variant.VariantTy, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal union value: %w", err)
			}
			unionV.Val = innerV
			*u = unionV

			return nil
		}
	}

	return fmt.Errorf("none of union prefixes matched")
}

func (u *UnionValue) Marshal(cell *boc.Cell, ty tolkParser.Union, encoder *Encoder) error {
	if len(ty.Variants) < 2 {
		return fmt.Errorf("union length must be at least 2")
	}
	if u.Prefix.Len > 64 {
		return fmt.Errorf("union prefix length must be less then or equals to 64")
	}

	if ty.Variants[0].PrefixEatInPlace {
		err := cell.WriteUint(u.Prefix.Prefix, int(u.Prefix.Len))
		if err != nil {
			return fmt.Errorf("failed to write union's %v-length prefix: %w", u.Prefix.Len, err)
		}
	}

	for _, variant := range ty.Variants {
		variantPrefix, err := PrefixToUint(variant.PrefixStr)
		if err != nil {
			return fmt.Errorf("failed to parse union's variant prefix %v: %w", variant.PrefixStr, err)
		}

		if u.Prefix.Prefix == variantPrefix {
			err = u.Val.Marshal(cell, variant.VariantTy, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal union value: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("none of union prefixes matched")
}

func (u *UnionValue) Equal(other any) bool {
	otherUnionValue, ok := other.(UnionValue)
	if !ok {
		return false
	}
	if u.Prefix != otherUnionValue.Prefix {
		return false
	}
	return u.Val.Equal(otherUnionValue.Val)
}
