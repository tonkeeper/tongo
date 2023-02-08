package tlb

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func (b Bits256) Hex() string {
	return hex.EncodeToString(b[:])
}

func (b Bits256) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%x\"", b), nil
}

func (b *Bits256) UnmarshalJSON(buf []byte) error {
	var sl []byte
	_, err := fmt.Fscanf(bytes.NewReader(buf), "\"%x\"", &sl)
	if len(sl) != 32 {
		return fmt.Errorf("can't parse 256bits %v", string(buf))
	}
	copy(b[:], sl)
	return err
}
