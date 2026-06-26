package abiVerifier

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// StringTail matches FunC's storeStringTail/loadStringTail — snake-encoded string stored
// ordinary tolk string is like Cell<StringTail>

func (t *StringTail) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	s, err := c.ReadStringTail()
	if err != nil {
		return fmt.Errorf("StringTail: %w", err)
	}
	*t = StringTail(s)
	return nil
}

func (t StringTail) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	_, err := c.WriteStringTail(string(t))
	return err
}
