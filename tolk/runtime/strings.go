package runtime

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type StringValue string

func (s *StringValue) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	value, err := cell.ReadStringRefTail()
	if err != nil {
		return fmt.Errorf("failed to read string value: %w", err)
	}
	*s = StringValue(value)
	return nil
}

func (s *StringValue) Marshal(cell *boc.Cell, encoder *Encoder) error {
	if _, err := cell.WriteStringRefTail(string(*s)); err != nil {
		return fmt.Errorf("failed to write string value: %w", err)
	}
	return nil
}

func (s *StringValue) Equal(o any) bool {
	otherString, ok := o.(StringValue)
	if !ok {
		return false
	}
	return *s == otherString
}
