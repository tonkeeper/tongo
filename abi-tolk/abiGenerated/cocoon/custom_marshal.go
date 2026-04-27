package abiCocoon

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (m *ForwardMsgs) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	/*
		FunC code:
		while (cs.slice_refs()) {
			var mode = cs~load_uint(8);
			send_raw_message(cs~load_ref(), mode);
		}
	*/
	for c.RefsAvailableForRead() > 0 {
		mode, err := c.ReadUint(8)
		if err != nil {
			return fmt.Errorf("failed to read fwd msg mode: %v", err)
		}
		ref, err := c.NextRef()
		if err != nil || ref == nil {
			return fmt.Errorf("failed to read fwd msg cell: %v", err)
		}
		*m = append(*m, ForwardMsg{
			Mode: tlb.Uint8(mode),
			Msg:  *ref,
		})
	}
	return nil
}

func (m *ForwardMsgs) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	for _, fwdMsg := range *m {
		if err := c.WriteUint(uint64(fwdMsg.Mode), 8); err != nil {
			return fmt.Errorf("failed to write fwd msg mode: %v", err)
		}
		if err := c.AddRef(&fwdMsg.Msg); err != nil {
			return fmt.Errorf("failed to write fwd msg cell: %v", err)
		}
	}
	return nil
}
