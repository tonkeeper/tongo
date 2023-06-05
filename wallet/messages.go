package wallet

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type MessageV3 struct {
	SubWalletId uint32
	ValidUntil  uint32
	Seqno       uint32
	RawMessages PayloadV1toV4
}

type MessageV4 struct {
	// Op: 0 - simple send, 1 - deploy and install plugin, 2 - install plugin, 3 - remove plugin
	SubWalletId uint32
	ValidUntil  uint32
	Seqno       uint32
	Op          int8
	RawMessages PayloadV1toV4
}

// SingedMsgBody represents an external message's body sent by an offchain application to a wallet contract.
// The signature is created using the wallet's private key.
// So the wallet will verify that it is the recipient of the payload and accept the payload.
type SingedMsgBody struct {
	Sign    tlb.Bits512
	Message tlb.Any
}

// RawMessage when received by a wallet contract will be resent as an outgoing message.
type RawMessage struct {
	Message *boc.Cell
	Mode    byte
}

type PayloadV1toV4 []RawMessage

func decodeMessageV4(msg *boc.Cell) (*MessageV4, error) {
	var m tlb.Message
	if err := tlb.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	msgBody := SingedMsgBody{}
	bodyCell := boc.Cell(m.Body.Value)
	if err := tlb.Unmarshal(&bodyCell, &msgBody); err != nil {
		return nil, err
	}
	msgv4 := MessageV4{}
	cell := boc.Cell(msgBody.Message)
	if err := tlb.Unmarshal(&cell, &msgv4); err != nil {
		return nil, err
	}
	return &msgv4, nil
}

func decodeMessageV3(msg *boc.Cell) (*MessageV3, error) {
	var m tlb.Message
	if err := tlb.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	msgBody := SingedMsgBody{}
	bodyCell := boc.Cell(m.Body.Value)
	if err := tlb.Unmarshal(&bodyCell, &msgBody); err != nil {
		return nil, err
	}
	msgv3 := MessageV3{}
	payloadCell := boc.Cell(msgBody.Message)
	if err := tlb.Unmarshal(&payloadCell, &msgv3); err != nil {
		return nil, err
	}
	return &msgv3, nil
}

// ExtractRawMessages extracts a list of RawMessages from an external message.
func ExtractRawMessages(ver Version, msg *boc.Cell) (PayloadV1toV4, error) {
	switch ver {
	case V4R1, V4R2:
		v4, err := decodeMessageV4(msg)
		if err != nil {
			return nil, err
		}
		// TODO: check opcode
		return v4.RawMessages, nil
	case V3R1, V3R2:
		v3, err := decodeMessageV3(msg)
		if err != nil {
			return nil, err
		}
		return v3.RawMessages, nil
	default:
		return nil, fmt.Errorf("wallet version is not supported: %v", ver)
	}
}

func (p PayloadV1toV4) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(p) > 4 {
		return fmt.Errorf("PayloadV1toV4 supports only up to 4 messages")
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

func (p *PayloadV1toV4) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
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
