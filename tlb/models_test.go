package tlb

import (
	"math"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func TestGrams_UnmarshalTLB(t *testing.T) {
	tests := []struct {
		name      string
		bocHex    string
		wantGrams Grams
	}{
		{
			bocHex:    "b5ee9c720101010100040000031fa8",
			wantGrams: Grams(250),
		},
		{
			bocHex:    "b5ee9c7201010101000300000108",
			wantGrams: Grams(0),
		},
		{
			bocHex:    "b5ee9c7201010101000b00001187fffffffffffffff8",
			wantGrams: Grams(math.MaxInt64),
		},
		{
			bocHex:    "b5ee9c7201010101000700000947fffffff8",
			wantGrams: Grams(math.MaxInt32),
		},
		{
			bocHex:    "b5ee9c720101010100070000094800000008",
			wantGrams: Grams(math.MaxInt32 + 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := boc.DeserializeBocHex(tt.bocHex)
			if err != nil {
				t.Fatalf("DeserializeBocHex() failed: %v", err)
			}
			var g Grams
			err = g.UnmarshalTLB(cell[0], &Decoder{})
			if err != nil {
				t.Fatalf("UnmarshalTLB() error = %v", err)
			}
			if tt.wantGrams != g {
				t.Fatalf("want: %v, got: %v", tt.wantGrams, g)
			}
		})
	}
}
