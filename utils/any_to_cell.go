package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
		cells, err := boc.DeserializeBoc(b)
		if err != nil {
			return nil, err
		}
		if len(cells) != 1 {
			return nil, fmt.Errorf("invalid number of cells: %v", len(cells))
		}
		return cells[0], nil
	case []byte:
		cells, err := boc.DeserializeBoc(v)
		if err != nil {
			return nil, err
		}
		if len(cells) != 1 {
			return nil, fmt.Errorf("invalid number of cells: %v", len(cells))
		}
		return cells[0], nil
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
