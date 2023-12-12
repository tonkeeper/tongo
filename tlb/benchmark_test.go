package tlb

import (
	"os"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func Benchmark_Tlb_Unmarshal(b *testing.B) {
	data, err := os.ReadFile("../testdata/raw-13516764.bin")
	if err != nil {
		b.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		b.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cell[0].ResetCounters()
		var block Block
		decoder := NewDecoder()
		err = decoder.Unmarshal(cell[0], &block)
		if err != nil {
			b.Errorf("Unmarshal() failed: %v", err)
		}
	}
}

func Test_block(b *testing.T) {
	data, err := os.ReadFile("testdata/block-1/block.bin")
	if err != nil {
		b.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		b.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	var block Block
	decoder := NewDecoder()
	err = decoder.Unmarshal(cell[0], &block)
	if err != nil {
		b.Errorf("Unmarshal() failed: %v", err)
	}
}
