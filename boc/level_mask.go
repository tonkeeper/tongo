package boc

import "math/bits"

type LevelMask uint32

func (m LevelMask) Level() int {
	return 32 - bits.LeadingZeros32(uint32(m))
}

func (m LevelMask) HashIndex() int {
	return bits.OnesCount32(uint32(m))
}

func (m LevelMask) HashesCount() int {
	return m.HashIndex() + 1
}

func (m LevelMask) Apply(level uint32) LevelMask {
	return LevelMask(uint32(m) & ((1 << level) - 1))
}

func (m LevelMask) ApplyOr(other LevelMask) LevelMask {
	return m | other
}
func (m LevelMask) ShiftRight() LevelMask {
	return m >> 1
}

func (m LevelMask) IsSignificant(level uint32) bool {
	if level == 0 {
		return true
	}
	return (m>>(level-1))%2 != 0
}