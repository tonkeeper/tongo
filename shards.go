package tongo

import (
	"github.com/tonkeeper/tongo/ton"
)

type ShardID = ton.ShardID

func ParseShardID(m int64) (ShardID, error) {
	return ton.ParseShardID(m)
}

func MustParseShardID(m int64) ShardID {
	return ton.MustParseShardID(m)
}
