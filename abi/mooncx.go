package abi

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (m *MoonSwapParams) UnmarshalJSON(data []byte) error {
	var r struct {
		MinOut      string
		Deadline    uint64
		Excess      string
		Referral    string
		NextFulfill *MoonNextPayload
		NextReject  *MoonNextPayload
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	minOut, ok := new(big.Int).SetString(r.MinOut, 10)
	if !ok {
		return fmt.Errorf("invalid minOut: %v", minOut)
	}
	m.MinOut = tlb.VarUInteger16(*minOut)
	m.Deadline = r.Deadline
	var excess tlb.MsgAddress
	if err := excess.UnmarshalJSON([]byte(r.Excess)); err != nil {
		return err
	}
	m.Excess = excess
	var referral tlb.MsgAddress
	if err := referral.UnmarshalJSON([]byte(r.Referral)); err != nil {
		return err
	}
	m.Referral = referral
	m.NextFulfill = r.NextFulfill
	m.NextReject = r.NextReject
	return nil
}

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
