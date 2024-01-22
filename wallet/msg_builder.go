package wallet

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (w *Wallet) CreateMessage(lifetime time.Duration, messages ...Sendable) (*tlb.Message, error) {
	var msgArray []tlb.RawMessage
	for _, m := range messages {
		intMsg, mode, err := m.ToInternal()
		if err != nil {
			return nil, err
		}
		cell := boc.NewCell()
		err = tlb.Marshal(cell, intMsg)
		if err != nil {
			return nil, err
		}
		msgArray = append(msgArray, tlb.RawMessage{Message: cell, Mode: mode})
	}
	err := checkMessagesLimit(len(msgArray), w.ver)
	if err != nil {
		return nil, err
	}
	bodyCell := boc.NewCell()
	switch w.ver {
	case HighLoadV2R2:
		boundedID := uint64(time.Now().Add(lifetime).UTC().Unix()<<32) + uint64(rand.Uint32())
		body := HighloadV2Message{
			SubWalletId:    uint32(w.subWalletId),
			BoundedQueryID: boundedID,
			RawMessages:    tlb.PayloadHighload(msgArray),
		}
		err = tlb.Marshal(bodyCell, body)
	default:
		return nil, fmt.Errorf("message body generation for this wallet is not supported: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("can not marshal wallet message body: %v", err)
	}

	signBytes, err := bodyCell.Sign(w.key)
	if err != nil {
		return nil, fmt.Errorf("can not sign wallet message body: %v", err)
	}
	bits512 := tlb.Bits512{}
	copy(bits512[:], signBytes[:])
	signedBody := SignedMsgBody{
		Sign:    bits512,
		Message: tlb.Any(*bodyCell),
	}
	signedBodyCell := boc.NewCell()
	if err = tlb.Marshal(signedBodyCell, signedBody); err != nil {
		return nil, fmt.Errorf("can not marshal signed body: %v", err)
	}
	extMsg, err := ton.CreateExternalMessage(w.address, signedBodyCell, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("can not create external message: %v", err)
	}
	return &extMsg, nil
}
