package abitolk

import (
	"fmt"
	"math/big"
)

func MustBigInt(s string) big.Int {
	b := big.Int{}
	res, ok := b.SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("bigint %v cannot not be parsed", s))
	}
	return *res
}

func BigIntFromUint(i uint64) big.Int {
	b := big.Int{}
	return *b.SetUint64(i)
}

func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
