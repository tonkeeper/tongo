package tlb

import (
	"encoding/hex"
	"math/big"
)

func (b Bits256) Hex() string {
	return hex.EncodeToString(b[:])
}

func NewVarUInteger16FromInt64(i int64) VarUInteger16 {
	b := big.NewInt(i)
	return VarUInteger16(*b)
}
