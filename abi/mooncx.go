package abi

import (
	"encoding/json"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (m *MoonNextPayload) UnmarshalJSON(data []byte) error {
	if data[0] == '"' { // MoonNextPayload has been marshalled to json as cell
		data = data[1 : len(data)-1]
		c, err := boc.DeserializeBocHex(string(data))
		if err != nil {
			return err
		}
		msgAddr := tlb.MsgAddress{}
		if err := tlb.Unmarshal(c[0], &msgAddr); err != nil {
			return err
		}
		m.Recipient = msgAddr
		hasRef, err := c[0].ReadBit()
		if err != nil {
			return err
		}
		if hasRef {
			ref, err := c[0].NextRef()
			if err != nil {
				return err
			}
			m.Payload = (*tlb.Any)(ref)
		}
		return nil
	}
	// otherwise MoonNextPayload has been marshalled to json as struct
	var r struct {
		Recipient string
		Payload   string
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	id := ton.MustParseAccountID(r.Recipient)
	m.Recipient = id.ToMsgAddress()
	if r.Payload != "" {
		c, err := boc.DeserializeBocHex(r.Payload)
		if err != nil {
			return err
		}
		m.Payload = (*tlb.Any)(c[0])
	}

	return nil
}
