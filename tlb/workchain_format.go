package tlb

import (
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type WorkchainFormat struct {
	SumType
	WorkchainFormat0 WorkchainFormat0
	WorkchainFormat1 WorkchainFormat1
}

func (t *WorkchainFormat) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	tag, err := c.PickUint(4)
	if err != nil {
		return err
	}
	switch tag {
	case 0:
		t.SumType = "WorkchainFormat0"
		return decoder.Unmarshal(c, &t.WorkchainFormat0)
	case 1:
		t.SumType = "WorkchainFormat1"
		return decoder.Unmarshal(c, &t.WorkchainFormat1)
	default:
		return fmt.Errorf("can not decode sumtype WorkchainFormat")
	}
}

func (t WorkchainFormat) MarshalJSON() ([]byte, error) {
	switch t.SumType {
	case "WorkchainFormat0":
		bytes, err := json.Marshal(t.WorkchainFormat0)
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf(`{"SumType": "WorkchainFormat0","WorkchainFormat0":%v}`, string(bytes))), nil
	case "WorkchainFormat1":
		bytes, err := json.Marshal(t.WorkchainFormat1)
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf(`{"SumType": "WorkchainFormat1","WorkchainFormat1":%v}`, string(bytes))), nil
	case "":
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("unknown sum type %v", t.SumType)
	}
}
