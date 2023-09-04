package tongo

import (
	"github.com/tonkeeper/tongo/ton"
)

type Bits256 = ton.Bits256

func ParseHash(s string) (Bits256, error) {
	return ton.ParseHash(s)
}

func MustParseHash(s string) Bits256 {
	return ton.MustParseHash(s)
}
