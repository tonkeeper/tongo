package boc

import (
	"math"
	"testing"
)

func TestAppend(t *testing.T) {
	s1 := NewBitString(8 * 10)
	s1.WriteUint(1, 80)
	s1.ReadBit()
	s2 := NewBitString(8 * 10)
	s2.WriteUint(1, 80)
	s2.ReadBit()
	s1.Append(s2)
}

func TestMinBits(t *testing.T) {
	for i := 0; i < 1000500; i++ {
		if minBitsRequired(uint64(i)) != int(math.Ceil(math.Log2(float64(i+1)))) {
			t.Fatal(i)
		}
	}
}

func BenchmarkMinbits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = minBitsRequired(uint64(i))
	}
}

func BenchmarkOldMinbits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int(math.Ceil(math.Log2(float64(i + 1))))
	}
}
