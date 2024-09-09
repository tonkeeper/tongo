package abi

import (
	"errors"
	"fmt"

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

type W5SendMessageAction struct {
	Magic tlb.Magic `tlb:"#0ec3c86d"`
	Mode  uint8
	Msg   MessageRelaxed `tlb:"^"`
}

type W5Actions []W5SendMessageAction

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
