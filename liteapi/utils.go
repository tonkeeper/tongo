package liteapi

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// VerifySendMessagePayload verifies that the given payload is an external message ready to be sent to the blockchain.
func VerifySendMessagePayload(payload []byte) error {
	roots, err := boc.DeserializeBoc(payload)
	if err != nil {
		return err
	}
	if len(roots) != 1 {
		return fmt.Errorf("external message is not a valid bag of cells")
	}
	root := roots[0]
	if root.Level() != 0 {
		return fmt.Errorf("external message must have zero level")
	}
	var msg tlb.Message
	if err := tlb.Unmarshal(root, &msg); err != nil {
		return fmt.Errorf("external message is not a tlb.Message")
	}
	if msg.Info.SumType != "ExtInMsgInfo" {
		return fmt.Errorf("external message must begin with ext_in_msg_info$10")
	}
	return nil
}
