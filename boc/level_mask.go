package boc

import "math/bits"

// levelMask is a tricky way to keep track of two things simultaneously:
// a cell level and a number of different hashes that makes sense to calculate for a cell.
type levelMask uint32

func (m levelMask) Level() int {
	return 32 - bits.LeadingZeros32(uint32(m))
}

func (m levelMask) HashIndex() int {
	return bits.OnesCount32(uint32(m))
}

func (m levelMask) HashesCount() int {
	return m.HashIndex() + 1
}

func (m levelMask) Apply(level int) levelMask {
	return levelMask(uint32(m) & ((1 << uint32(level)) - 1))
}

func (m levelMask) IsSignificant(level uint32) bool {
	if level == 0 {
		return true
	}
	return (m>>(level-1))%2 != 0
}
