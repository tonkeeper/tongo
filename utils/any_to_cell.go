package utils

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func AnyToCell(i any) (*boc.Cell, error) {
	if i == nil {
		return nil, nil
	}
	switch v := i.(type) {
	case string:
		b, err := hex.DecodeString(v)
		if err != nil {
			b, err = base64.StdEncoding.DecodeString(v)
			if err != nil {
				return nil, err
			}
		}
		cell, err := boc.DeserializeSingleRootBoc(b)
		if err != nil {
			return nil, err
		}
		return cell, nil
	case []byte:
		cell, err := boc.DeserializeSingleRootBoc(v)
		if err != nil {
			return nil, err
		}
		return cell, nil
	case *boc.Cell:
		return v, nil
	case boc.Cell:
		return &v, nil
	default:
		c := boc.NewCell()
		err := tlb.Marshal(c, i)
		return c, err
	}
}
