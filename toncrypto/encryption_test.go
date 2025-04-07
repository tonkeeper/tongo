package toncrypto

import (
	"crypto/rand"
	"slices"
	"testing"
)

func TestPrefix(t *testing.T) {
	for i := 1; i < 10_000; i++ {
		data := make([]byte, i)
		_, _ = rand.Read(data)
		newData, err := addPrefix(data)
		if err != nil {
			t.Fatal(err)
		}
		prefixLen := len(newData) - i
		if prefixLen < 16 || prefixLen > 31 {
			t.Fatal("prefix length should be between 16 and 31")
		}
		if len(newData)%16 != 0 {
			t.Fatal("new data len must be a multiple of 16")
		}
		if int(newData[0]) != prefixLen {
			t.Fatal("new data first byte should be prefixLen")
		}
		if !slices.Equal(newData[prefixLen:], data) {
			t.Fatal("data corrupted")
		}
	}
}
