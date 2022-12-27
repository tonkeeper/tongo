package tongo

import (
	"io/ioutil"
	"testing"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

func Benchmark_Tlb_Unmarshal(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/raw-block.bin")
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
		err = tlb.Unmarshal(cell[0], &block)
		if err != nil {
			b.Errorf("Unmarshal() failed: %v", err)
		}
	}
}
