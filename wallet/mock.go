package wallet

import (
	"context"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/tlb"
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
	select {
	case b.messages <- payload:
	default:
	}

	b.seqno++ // it does not check message, address and seqno logic
	// it does not modify account state
	return 0, nil
}

func (b *SimpleMockBlockchain) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tlb.ShardAccount, error) {
	return tlb.ShardAccount{
		Account: tlb.Account{
			SumType: "Account",
			Account: struct {
				Addr        tlb.MsgAddress
				StorageStat tlb.StorageInfo
				Storage     tlb.AccountStorage
			}{
				Addr: accountID.ToMsgAddress(),
				Storage: tlb.AccountStorage{
					LastTransLt: 0,
					Balance: tlb.CurrencyCollection{
						Grams: tlb.Grams(b.state.Balance),
					},
					State: tlb.AccountState{
						SumType: "AccountActive",
					},
				},
			},
		},
		LastTransHash: tlb.Bits256{},
		LastTransLt:   0,
	}, nil
}
