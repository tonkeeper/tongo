package wallet

import (
	"context"
	"fmt"
	"github.com/startfellows/tongo"
)

// SimpleMockBlockchain
// Very simple mock. It does not provide blockchain logic for calculating state and seqno for different addresses.
// Only for internal tests and demonstration purposes.
type SimpleMockBlockchain struct {
	seqno    uint32
	state    tongo.AccountInfo
	messages chan []byte
}

func NewMockBlockchain(seqno uint32, state tongo.AccountInfo) (*SimpleMockBlockchain, chan []byte) {
	c := make(chan []byte, 100)
	return &SimpleMockBlockchain{
		seqno:    seqno,
		state:    state,
		messages: c,
	}, c
}

func (b *SimpleMockBlockchain) GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error) {
	return b.seqno, nil
}

func (b *SimpleMockBlockchain) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	b.messages <- payload
	b.seqno++ // it does not check message, address and seqno logic
	// it does not modify account state
	return 0, nil
}

func (b *SimpleMockBlockchain) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tongo.ShardAccount, error) {
	// TODO: fix
	return tongo.ShardAccount{}, fmt.Errorf("not implemnted")
}
