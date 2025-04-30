package wallet

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

var ErrBadSignature = errors.New("failed to verify msg signature")

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

type W5SendMessageAction struct {
	Magic tlb.Magic `tlb:"#0ec3c86d"`
	Mode  uint8
	Msg   *boc.Cell `tlb:"^"`
}

type W5Actions []W5SendMessageAction

// MessageV5Beta is a message format used by wallet v5 beta.
type MessageV5Beta struct {
	tlb.SumType
	// SignedInternal is an internal message authenticated by a signature.
	SignedInternal struct {
		WalletId   tlb.Bits80
		ValidUntil uint32
		Seqno      uint32
		Op         bool
		Signature  tlb.Bits512
		Actions    W5Actions `tlb:"^"`
	} `tlbSumType:"#73696e74"`
	// SignedExternal is an external message authenticated by a signature.
	SignedExternal struct {
		WalletId   tlb.Bits80
		ValidUntil uint32
		Seqno      uint32
		Op         bool
		Signature  tlb.Bits512
		Actions    W5Actions `tlb:"^"`
	} `tlbSumType:"#7369676e"`
}

// MessageV5 is a message format used by wallet v5.
type MessageV5 struct {
	tlb.SumType
	SignedInternal *struct {
		WalletId   uint32
		ValidUntil uint32
		Seqno      uint32
		Actions    *BlumActions `tlb:"maybe^"`
		Signature  tlb.Bits512
	} `tlbSumType:"#73696e74"` // op::internal_signed

	SignedExternal *struct {
		WalletId   uint32
		ValidUntil uint32
		Seqno      uint32
		Actions    *BlumActions `tlb:"maybe^"`
		Signature  tlb.Bits512
	} `tlbSumType:"#7369676e"` // op::external_signed

	AdminMessage *struct {
		WalletId   uint32
		ValidUntil uint32
		Seqno      uint32
		Op         uint32       //update_code, mass_suspend, single_suspend
		Actions    *BlumActions `tlb:"maybe^"`
		Signatures tlb.HashmapE[tlb.Bits256, tlb.Bits512]
	} `tlbSumType:"#5f1dc312"` // op::external_admin

	SessionKeyMessage *struct {
		PublicKey  tlb.Bits256
		Address    tlb.MsgAddress
		Commission tlb.Coins
		Actions    *BlumActions `tlb:"maybe^"`
	} `tlbSumType:"#9083118"` // op::internal_execute_with_session_key
}

type BlumActions []BlumSendMessageAction

type BlumSendMessageAction struct {
	Magic tlb.Magic `tlb:"#0ec3c86d"` // action_send_msg
	Mode  uint8
	Msg   *boc.Cell `tlb:"^"`
}

func (l BlumActions) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(l) == 0 {
		return nil
	}
	if err := c.WriteUint(0x0ec3c86d, 32); err != nil {
		return err
	}
	action := l[0]
	if err := c.WriteUint(uint64(action.Mode), 8); err != nil {
		return err
	}
	cell := boc.NewCell()
	next := l[1:]
	if err := encoder.Marshal(cell, next); err != nil {
		return err
	}
	if err := c.AddRef(cell); err != nil {
		return err
	}
	if err := c.AddRef(action.Msg); err != nil {
		return err
	}
	return nil
}

func (l *BlumActions) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var actions []BlumSendMessageAction
	for {
		switch c.BitsAvailableForRead() {
		case 0:
			*l = actions
			return nil
		case 40:
			next, err := c.NextRef()
			if err != nil {
				return err
			}
			var action BlumSendMessageAction
			if err := decoder.Unmarshal(c, &action); err != nil {
				return err
			}
			actions = append(actions, action)
			c = next
		default:
			return fmt.Errorf("unexpected bits available: %v", c.BitsAvailableForRead())
		}
	}
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
	return ErrBadSignature
}

func extractSignedMsgBody(msg *boc.Cell) (*SignedMsgBody, error) {
	var m tlb.Message
	if err := tlb.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	msgBody := SignedMsgBody{}
	bodyCell := boc.Cell(m.Body.Value)
	if err := tlb.Unmarshal(&bodyCell, &msgBody); err != nil {
		return nil, err
	}
	return &msgBody, nil
}

func DecodeMessageV5(msg *boc.Cell) (*MessageV5, error) {
	var m tlb.Message
	if err := tlb.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	var msgv5 MessageV5
	bodyCell := boc.Cell(m.Body.Value)
	if err := tlb.Unmarshal(&bodyCell, &msgv5); err != nil {
		return nil, err
	}
	return &msgv5, nil
}

func DecodeMessageV5Beta(msg *boc.Cell) (*MessageV5Beta, error) {
	var m tlb.Message
	if err := tlb.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	var msgv5 MessageV5Beta
	bodyCell := boc.Cell(m.Body.Value)
	if err := tlb.Unmarshal(&bodyCell, &msgv5); err != nil {
		return nil, err
	}
	return &msgv5, nil
}

func DecodeMessageV4(msg *boc.Cell) (*MessageV4, error) {
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

func DecodeMessageV3(msg *boc.Cell) (*MessageV3, error) {
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

func DecodeHighloadV2Message(msg *boc.Cell) (*HighloadV2Message, error) {
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
func ExtractRawMessages(ver Version, msg *boc.Cell) ([]RawMessage, error) {
	switch ver {
	case V5Beta:
		v5, err := DecodeMessageV5Beta(msg)
		if err != nil {
			return nil, err
		}
		return v5.RawMessages(), nil
	case V5R1:
		v5, err := DecodeMessageV5(msg)
		if err != nil {
			return nil, err
		}
		return v5.RawMessages(), nil
	case V4R1, V4R2:
		v4, err := DecodeMessageV4(msg)
		if err != nil {
			return nil, err
		}
		// TODO: check opcode
		return v4.RawMessages, nil
	case V3R1, V3R2, V3R2Lockup:
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
		return hl.RawMessages, nil
	default:
		return nil, fmt.Errorf("wallet version is not supported: %v", ver)
	}
}

// VerifySignature checks whether the given message (tlb.Message) represented as a cell
// was signed by the given public key of a wallet contract.
// On success, it returns nil.
// Otherwise, it returns an error.
func VerifySignature(ver Version, msg *boc.Cell, publicKey ed25519.PublicKey) error {
	switch ver {
	case V3R1, V3R2, V3R2Lockup, V4R1, V4R2, HighLoadV2R2:
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

func (l *W5Actions) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var actions []W5SendMessageAction
	for {
		switch c.BitsAvailableForRead() {
		case 0:
			*l = actions
			return nil
		case 40:
			next, err := c.NextRef()
			if err != nil {
				return err
			}
			var action W5SendMessageAction
			if err := decoder.Unmarshal(c, &action); err != nil {
				return err
			}
			actions = append(actions, action)
			c = next
		default:
			return fmt.Errorf("unexpected bits available: %v", c.BitsAvailableForRead())
		}
	}
}

func (l W5Actions) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(l) == 0 {
		return nil
	}
	if err := c.WriteUint(0x0ec3c86d, 32); err != nil {
		return err
	}
	action := l[0]
	if err := c.WriteUint(uint64(action.Mode), 8); err != nil {
		return err
	}
	cell := boc.NewCell()
	next := l[1:]
	if err := encoder.Marshal(cell, next); err != nil {
		return err
	}
	if err := c.AddRef(cell); err != nil {
		return err
	}
	if err := c.AddRef(action.Msg); err != nil {
		return err
	}
	return nil
}

func MessageV5VerifySignature(msgBody boc.Cell, publicKey ed25519.PublicKey) error {
	totalBits := msgBody.BitsAvailableForRead()
	if totalBits < 512 {
		return fmt.Errorf("not enough bits in the cell")
	}
	bits, err := msgBody.ReadBits(totalBits - 512)
	if err != nil {
		return err
	}
	signature, err := msgBody.ReadBytes(64)
	if err != nil {
		return err
	}
	msgCopy := boc.NewCell()
	if err := msgCopy.WriteBitString(bits); err != nil {
		return err
	}
	for i := 0; i < msgBody.RefsSize(); i++ {
		ref, err := msgBody.NextRef()
		if err != nil {
			return err
		}
		if err := msgCopy.AddRef(ref); err != nil {
			return err
		}
	}
	hash, err := msgCopy.Hash()
	if err != nil {
		return err
	}
	if ed25519.Verify(publicKey, hash, signature) {
		return nil
	}
	return ErrBadSignature
}

func (m *MessageV5Beta) RawMessages() []RawMessage {
	switch m.SumType {
	case "SignedInternal":
		msgs := make([]RawMessage, 0, len(m.SignedInternal.Actions))
		for _, action := range m.SignedInternal.Actions {
			msgs = append(msgs, RawMessage{
				Message: action.Msg,
				Mode:    action.Mode,
			})
		}
		return msgs
	case "SignedExternal":
		msgs := make([]RawMessage, 0, len(m.SignedExternal.Actions))
		for _, action := range m.SignedExternal.Actions {
			msgs = append(msgs, RawMessage{
				Message: action.Msg,
				Mode:    action.Mode,
			})
		}
		return msgs
	default:
		return nil
	}
}

func (m *MessageV5) RawMessages() []RawMessage {
	switch m.SumType {
	case "SignedInternal":
		if m.SignedInternal == nil || m.SignedInternal.Actions == nil {
			return nil
		}
		actions := m.SignedInternal.Actions
		msgs := make([]RawMessage, 0, len(*actions))
		for _, action := range *actions {
			msgs = append(msgs, RawMessage{
				Message: action.Msg,
				Mode:    action.Mode,
			})
		}
		return msgs
	case "SignedExternal":
		if m.SignedExternal == nil || m.SignedExternal.Actions == nil {
			return nil
		}
		msgs := make([]RawMessage, 0, len(*m.SignedExternal.Actions))
		for _, action := range *m.SignedExternal.Actions {
			msgs = append(msgs, RawMessage{
				Message: action.Msg,
				Mode:    action.Mode,
			})
		}
		return msgs
	default:
		return nil
	}
}

type W5ExtendedAction struct {
	SumType      tlb.SumType
	AddExtension *struct {
		Addr tlb.MsgAddress
	} `tlbSumType:"add_extension#02"`
	RemoveExtension *struct {
		Addr tlb.MsgAddress
	} `tlbSumType:"remove_extension#03"`
	SetSignatureAllowed *struct {
		Allowed bool
	} `tlbSumType:"set_signature_allowed#04"`
}

type W5ExtendedActions []W5ExtendedAction

func (extendedActions *W5ExtendedActions) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var actions []W5ExtendedAction
	for {
		var action W5ExtendedAction
		if err := decoder.Unmarshal(c, &action); err != nil {
			return err
		}
		actions = append(actions, action)
		nextRef, err := c.NextRef()
		if err != nil {
			if errors.Is(err, boc.ErrNotEnoughRefs) {
				*extendedActions = actions
				return nil
			}
			return err
		}
		c = nextRef
	}
}
func (extendedActions W5ExtendedActions) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	for i, action := range extendedActions {
		if err := encoder.Marshal(c, action); err != nil {
			return err
		}
		if i == len(extendedActions)-1 {
			break
		}
		cell := boc.NewCell()
		if err := c.AddRef(cell); err != nil {
			return err
		}
		c = cell
	}
	return nil
}
