package tongo

import (
	"encoding/binary"
	"errors"
	"math/bits"
)

type ShardID struct {
	prefix int64
	mask   int64
}

func ParseShardID(m int64) (ShardID, error) {
	if m == 0 {
		return ShardID{}, errors.New("at least one non-zero bit required in shard id")
	}
	trailingZeros := bits.TrailingZeros64(uint64(m))
	return ShardID{
		prefix: m ^ (1 << trailingZeros),
		mask:   -1 << (trailingZeros + 1),
	}, nil
}

func MustParseShardID(m int64) ShardID {
	s, err := ParseShardID(m)
	if err != nil {
		panic(err)
	}
	return s
}

func (s ShardID) Encode() int64 {
	return s.prefix | (1 << (bits.TrailingZeros64(uint64(s.mask)) - 1))
}

func (s ShardID) MatchAccountID(a AccountID) bool {
	aPrefix := binary.BigEndian.Uint64(a.Address[:8])
	return (int64(aPrefix) & s.mask) == s.prefix
}

func (s ShardID) MatchBlockID(block TonNodeBlockId) bool {
	return true //todo: write tests and fixes
	sub, err := ParseShardID(block.Shard)
	if err != nil {
		return false
	}
	if bits.TrailingZeros64(uint64(s.mask)) < bits.TrailingZeros64(uint64(sub.mask)) {
		return false
	}
	return sub.prefix&s.mask == s.mask
}
