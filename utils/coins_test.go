package utils

import (
	"testing"
)

func Test_HumanFriendlyCoinsRepr(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		want   string
	}{
		{
			amount: 50,
			want:   "50 nanoTon",
		},
		{
			amount: 1_000,
			want:   "1 microTon",
		},
		{
			amount: 999_950,
			want:   "999.95 microTon",
		},
		{
			amount: 3_000_000,
			want:   "3 milliTon",
		},
		{
			amount: 994_000_745,
			want:   "994.000745 milliTon",
		},
		{
			amount: 995_000_745_000,
			want:   "995.000745 Ton",
		},
		{
			amount: 1_000_000_000,
			want:   "1 Ton",
		},
		{
			amount: 1_000_000_050,
			want:   "1.00000005 Ton",
		},
		{
			amount: 999_000_000_050,
			want:   "999.00000005 Ton",
		},
		{
			amount: 1_000_000_000_000,
			want:   "1 kiloTon",
		},
		{
			amount: 999_000_350_000_000,
			want:   "999.00035 kiloTon",
		},
		{
			amount: 9_500_000_000_000,
			want:   "9.5 kiloTon",
		},
		{
			amount: 1_000_000_000_000_000,
			want:   "1 megaTon",
		},
		{
			amount: 990_000_000_000_000_333,
			want:   "990.000000000000333 megaTon",
		},
		{
			amount: 8_950_000_000_000_000,
			want:   "8.95 megaTon",
		},
		{
			amount: 9_999_999_990_000_000,
			want:   "9.99999999 megaTon",
		},
		{
			amount: 1_000_000_000_000_000_000,
			want:   "1 gigaTon",
		},
		{
			amount: 1_600_000_000_000_000_000,
			want:   "1.6 gigaTon",
		},
		{
			amount: 2_950_100_000_000_000_000,
			want:   "2.9501 gigaTon",
		},
		{
			amount: 9_000_000_000_000_000_000,
			want:   "9 gigaTon",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repr := HumanFriendlyCoinsRepr(tt.amount)
			if tt.want != repr {
				t.Errorf("expected: %v, got: %v", tt.want, repr)
			}
		})
	}
}
