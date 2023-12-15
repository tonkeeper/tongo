package wallet

import (
	"crypto/ed25519"
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

type HighloadV2Message struct {
	SubWalletId    uint32
	BoundedQueryID uint64
	RawMessages    PayloadHighload
}

// SignedMsgBody represents an external message's body sent by an offchain application to a wallet contract.
// The signature is created using the wallet's private key.
// So the wallet will verify that it is the recipient of the payload and accept the payload.
type SignedMsgBody struct {
	Sign    tlb.Bits512
	Message tlb.Any
}

// RawMessage when received by a wallet contract will be resent as an outgoing message.
type RawMessage struct {
	Message *boc.Cell
	Mode    byte
}

// ToTlbMessage converts a RawMessage to a tlb.Message.
func (rm *RawMessage) ToTlbMessage() (*tlb.Message, error) {
	var msg tlb.Message
	if err := tlb.Unmarshal(rm.Message, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

type PayloadV1toV4 []RawMessage
type PayloadHighload []RawMessage

func (body *SignedMsgBody) Verify(publicKey ed25519.PublicKey) error {
	msg := boc.Cell(body.Message)
	hash, err := msg.Hash()
	if err != nil {
		return err
	}
	if ed25519.Verify(publicKey, hash, body.Sign[:]) {
		return nil
	}
	return fmt.Errorf("failed to verify msg signature")
}

func extractSignedMsgBody(msg *tlb.Message) (*SignedMsgBody, error) {
	bodyCell := boc.Cell(msg.Body.Value)
	msgBody := SignedMsgBody{}
	if err := tlb.Unmarshal(&bodyCell, &msgBody); err != nil {
		return nil, err
	}
	return &msgBody, nil
}

func DecodeMessageV4(msg *tlb.Message) (*MessageV4, error) {
	signedMsgBody, err := extractSignedMsgBody(msg)
	if err != nil {
		return nil, err
	}
	return decodeMessageV4(signedMsgBody)
}

func decodeMessageV4(body *SignedMsgBody) (*MessageV4, error) {
	msgv4 := MessageV4{}
	cell := boc.Cell(body.Message)
	if err := tlb.Unmarshal(&cell, &msgv4); err != nil {
		return nil, err
	}
	return &msgv4, nil
}

func DecodeMessageV3(msg *tlb.Message) (*MessageV3, error) {
	signedMsgBody, err := extractSignedMsgBody(msg)
	if err != nil {
		return nil, err
	}
	return decodeMessageV3(signedMsgBody)
}

func decodeMessageV3(body *SignedMsgBody) (*MessageV3, error) {
	msgv3 := MessageV3{}
	payloadCell := boc.Cell(body.Message)
	if err := tlb.Unmarshal(&payloadCell, &msgv3); err != nil {
		return nil, err
	}
	return &msgv3, nil
}

func DecodeHighloadV2Message(msg *tlb.Message) (*HighloadV2Message, error) {
	signedMsgBody, err := extractSignedMsgBody(msg)
	if err != nil {
		return nil, err
	}
	return decodeHighloadV2Message(signedMsgBody)
}

func decodeHighloadV2Message(body *SignedMsgBody) (*HighloadV2Message, error) {
	msg := HighloadV2Message{}
	payloadCell := boc.Cell(body.Message)
	if err := tlb.Unmarshal(&payloadCell, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ExtractRawMessages extracts a list of RawMessages from an external message.
func ExtractRawMessages(ver Version, msg *tlb.Message) (PayloadV1toV4, error) {
	if msg == nil {
		return nil, fmt.Errorf("msg is nil")
	}
	switch ver {
	case V4R1, V4R2:
		v4, err := DecodeMessageV4(msg)
		if err != nil {
			return nil, err
		}
		// TODO: check opcode
		return v4.RawMessages, nil
	case V3R1, V3R2:
		v3, err := DecodeMessageV3(msg)
		if err != nil {
			return nil, err
		}
		return v3.RawMessages, nil
	case HighLoadV2R2:
		hl, err := DecodeHighloadV2Message(msg)
		if err != nil {
			return nil, err
		}
		return []RawMessage(hl.RawMessages), nil
	default:
		return nil, fmt.Errorf("wallet version is not supported: %v", ver)
	}
}

// VerifySignature checks whether the given message (tlb.Message) represented as a cell
// was signed by the given public key of a wallet contract.
// On success, it returns nil.
// Otherwise, it returns an error.
func VerifySignature(ver Version, msg *tlb.Message, publicKey ed25519.PublicKey) error {
	if msg == nil {
		return fmt.Errorf("msg is nil")
	}
	switch ver {
	case V3R1, V3R2, V4R1, V4R2, HighLoadV2R2:
		signedMsgBody, err := extractSignedMsgBody(msg)
		if err != nil {
			return err
		}
		return signedMsgBody.Verify(publicKey)
	default:
		return fmt.Errorf("wallet version is not supported: %v", ver)
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

func (p PayloadHighload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(p) > 254 {
		return fmt.Errorf("PayloadHighload supports only up to 254 messages")
	}
	var keys []tlb.Uint16
	var values []tlb.Any
	for i, msg := range p {
		cell := boc.NewCell()
		if err := cell.WriteUint(uint64(msg.Mode), 8); err != nil {
			return err
		}
		if err := cell.AddRef(msg.Message); err != nil {
			return err
		}
		keys = append(keys, tlb.Uint16(i))
		values = append(values, tlb.Any(*cell))
	}
	hashmap := tlb.NewHashmap[tlb.Uint16, tlb.Any](keys, values)
	dict := boc.NewCell()
	if err := tlb.Marshal(dict, hashmap); err != nil {
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

func (p *PayloadHighload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var m tlb.HashmapE[tlb.Uint16, tlb.Any]
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
