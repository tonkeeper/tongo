package tlb

import (
	"encoding/hex"
)

func (b Bits256) Hex() string {
	return hex.EncodeToString(b[:])
}