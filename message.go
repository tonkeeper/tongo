package tongo

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func decodeBoc(bocStr string) ([]byte, error) {
	bocData, err := base64.StdEncoding.DecodeString(bocStr)
	if err != nil {
		return hex.DecodeString(bocStr)
	}
	return bocData, nil
}

// Message contains a tlb.Message, its boc representation and a Cell.
type Message struct {
	Boc    []byte
	TlbMsg *tlb.Message
	Cell   *boc.Cell
}

// ParseTlbMessage returns a Message unmarshalled from the given boc string.
// The boc string can be either in base64 or hex format.
func ParseTlbMessage(bocStr string) (*Message, error) {
	b, err := decodeBoc(bocStr)
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(b)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid message boc")
	}
	var msg tlb.Message
	if err := tlb.Unmarshal(cells[0], &msg); err != nil {
		return nil, err
	}
	cells[0].ResetCounters()
	return &Message{Boc: b, TlbMsg: &msg, Cell: cells[0]}, nil
}

func (m *Message) DestinationAccountID() (ton.AccountID, error) {
	var dest tlb.MsgAddress
	switch m.TlbMsg.Info.SumType {
	case "IntMsgInfo":
		dest = m.TlbMsg.Info.IntMsgInfo.Dest
	case "ExtInMsgInfo":
		dest = m.TlbMsg.Info.ExtInMsgInfo.Dest
	case "ExtOutMsgInfo":
		dest = m.TlbMsg.Info.ExtOutMsgInfo.Dest
	}
	accountID, err := ton.AccountIDFromTlb(dest)
	if err != nil {
		return ton.AccountID{}, err
	}
	if accountID == nil {
		return ton.AccountID{}, fmt.Errorf("failed to extract the destination address")
	}
	return *accountID, nil
}
