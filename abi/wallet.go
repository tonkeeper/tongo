package abi

import (
	"errors"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// we need to write wallets payload type manually because it can be described on tlb

type WalletV1ToV4Payload []SendMessageAction

func (p *WalletV1ToV4Payload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	for {
		var msg SendMessageAction
		err := decoder.Unmarshal(c, &msg)
		if errors.Is(err, boc.ErrNotEnoughBits) {
			break
		}
		if err != nil {
			return err
		}
		*p = append(*p, msg)
	}
	return nil
}
