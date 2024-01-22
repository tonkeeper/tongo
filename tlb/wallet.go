package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

// RawMessage when received by a wallet contract will be resent as an outgoing message.
type RawMessage struct {
	Message *boc.Cell
	Mode    byte
}

type WalletPayloadV1toV4 []RawMessage
type PayloadHighload []RawMessage

func (p WalletPayloadV1toV4) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	if len(p) > 4 {
		return fmt.Errorf("WalletPayloadV1toV4 supports only up to 4 messages")
	}
	for _, msg := range p {
		err := c.WriteUint(uint64(msg.Mode), 8)
		if err != nil {
			return err
		}
		err = c.AddRef(msg.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *WalletPayloadV1toV4) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	for {
		ref, err := c.NextRef()
		if err != nil {
			break
		}
		mode, err := c.ReadUint(8)
		if err != nil {
			return err
		}
		msg := RawMessage{
			Message: ref,
			Mode:    byte(mode),
		}
		*p = append(*p, msg)
	}
	return nil
}

func (p PayloadHighload) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	if len(p) > 254 {
		return fmt.Errorf("PayloadHighload supports only up to 254 messages")
	}
	var keys []Uint16
	var values []Any
	for i, msg := range p {
		cell := boc.NewCell()
		if err := cell.WriteUint(uint64(msg.Mode), 8); err != nil {
			return err
		}
		if err := cell.AddRef(msg.Message); err != nil {
			return err
		}
		keys = append(keys, Uint16(i))
		values = append(values, Any(*cell))
	}
	hashmap := NewHashmap[Uint16, Any](keys, values)
	dict := boc.NewCell()
	if err := Marshal(dict, hashmap); err != nil {
		return err
	}
	if err := c.WriteBit(true); err != nil {
		return err
	}
	if err := c.AddRef(dict); err != nil {
		return err
	}
	return nil
}

func (p *PayloadHighload) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	var m HashmapE[Uint16, Any]
	if err := decoder.Unmarshal(c, &m); err != nil {
		return err
	}
	rawMessages := make([]RawMessage, 0, len(m.Values()))
	for _, item := range m.Items() {
		cell := boc.Cell(item.Value)
		mode, err := cell.ReadUint(8)
		if err != nil {
			return fmt.Errorf("failed to read msg mode: %v", err)
		}
		msg, err := cell.NextRef()
		if err != nil {
			return fmt.Errorf("failed to read msg: %v", err)
		}
		rawMsg := RawMessage{
			Message: msg,
			Mode:    byte(mode),
		}
		rawMessages = append(rawMessages, rawMsg)
	}
	*p = rawMessages
	return nil
}
