package utils

import (
	"testing"
)

func TestCrc32String(t *testing.T) {
	tests := []struct {
		data string
		want uint32
	}{
		{
			data: "swap_refund_no_liq",
			want: 1610486421,
		},
	}
	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			if got := Crc32String(tt.data); got != tt.want {
				t.Errorf("Crc32String() = %v, want %v", got, tt.want)
			}
		})
	}
}
