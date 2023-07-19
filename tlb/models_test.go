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

func TestSignedCoinsMarshalTLB(t *testing.T) {
	tests := []struct {
		name            string
		signedCoins     SignedCoins
		wantSignedCoins SignedCoins
	}{
		{
			signedCoins:     SignedCoins(0),
			wantSignedCoins: SignedCoins(0),
		},
		{
			signedCoins:     SignedCoins(-1),
			wantSignedCoins: SignedCoins(-1),
		},
		{
			signedCoins:     SignedCoins(1),
			wantSignedCoins: SignedCoins(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cell := boc.NewCell()
			err := test.signedCoins.MarshalTLB(cell, &Encoder{})
			if err != nil {
				t.Fatalf("MarshalTLB() error = %v", err)
			}
			var sc SignedCoins
			err = sc.UnmarshalTLB(cell, NewDecoder())
			if err != nil {
				t.Fatalf("UnmarshalTLB() error = %v", err)
			}
			if sc != test.wantSignedCoins {
				t.Fatalf("want: %v, got: %v", test.wantSignedCoins, sc)
			}
		})
	}
}

func TestSignedCoinsFailedTLB(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		bitLen    int
		wantError error
	}{
		{
			value:     1,
			bitLen:    1,
			wantError: boc.ErrNotEnoughBits,
		},
		{
			value:     2,
			bitLen:    2,
			wantError: ErrGramsOverflow,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sc SignedCoins
			cell := boc.NewCell()
			err := cell.WriteBit(false)
			if err != nil {
				t.Fatalf("WriteBit() failed: %v", err)
			}
			err = cell.WriteLimUint(test.value, test.bitLen)
			if err != nil {
				t.Fatalf("WriteLimUint() failed: %v", err)
			}
			err = cell.WriteUint(uint64(test.value), test.bitLen)
			if err != nil {
				t.Fatalf("WriteUint() failed: %v", err)
			}
			err = sc.UnmarshalTLB(cell, NewDecoder())
			if err != test.wantError {
				t.Fatalf("want: %v, got: %v", test.wantError, err)
			}
		})
	}
}

func TestSignedCoins_UnmarshalTLB(t *testing.T) {
	tests := []struct {
		name            string
		value           uint64
		wantSignedCoins SignedCoins
	}{
		{
			value:           0,
			wantSignedCoins: SignedCoins(0),
		},
		{
			value:           1,
			wantSignedCoins: SignedCoins(1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sc SignedCoins
			cell := boc.NewCell()
			err := cell.WriteBit(false)
			if err != nil {
				t.Fatalf("WriteBit() failed: %v", err)
			}
			err = cell.WriteLimUint(3, 15)
			if err != nil {
				t.Fatalf("WriteLimUint() failed: %v", err)
			}
			err = cell.WriteUint(test.value, 24)
			if err != nil {
				t.Fatalf("WriteUint() failed: %v", err)
			}
			err = sc.UnmarshalTLB(cell, NewDecoder())
			if err != nil {
				t.Fatalf("UnmarshalTLB() error = %v", err)
			}

			if sc != test.wantSignedCoins {
				t.Fatalf("want: %v, got: %v", test.wantSignedCoins, sc)
			}
		})
	}
}
