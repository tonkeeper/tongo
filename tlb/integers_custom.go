package tlb

import (
	"encoding/hex"
	"math/big"
)

func (b Bits256) Hex() string {
	return hex.EncodeToString(b[:])
}

func VarUInteger16FromInt64(i int64) VarUInteger16 {
	b := big.NewInt(i)
	return VarUInteger16(*b)
}

func Int257FromInt64(i int64) Int257 {
	b := big.NewInt(i)
	return Int257(*b)
}
