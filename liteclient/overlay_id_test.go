package liteclient

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/tonkeeper/tongo/tlb"
)

func TestOverlayID_FullID(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("XplPz01CXAps5qeSWUtxcyBfdAo5zVb1N979KLSKD24=")
	if err != nil {
		panic(err)
	}
	var zeroStateFileHash tlb.Bits256
	copy(zeroStateFileHash[:], bytes)

	tests := []struct {
		name      string
		overlayID OverlayID
		want      string
	}{
		{
			name: "masterchain",
			overlayID: OverlayID{
				Workchain:         -1,
				Shard:             -9223372036854775808,
				ZeroStateFileHash: zeroStateFileHash,
			},
			want: "c684cd30e81e3ad7159bbef689daea0021dae2b90dd1a65d14fe8cc11f3523b1",
		},
		{
			name: "basechain",
			overlayID: OverlayID{
				Workchain:         0,
				Shard:             -9223372036854775808,
				ZeroStateFileHash: zeroStateFileHash,
			},
			want: "9435c212dc0ec51dac686410e9ba98f4b6fc7d5f08aeb9164109178eb950ddec",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := tt.overlayID.FullID()
			if err != nil {
				t.Fatal(err)
			}
			hexStr := hex.EncodeToString(shortID)
			if hexStr != tt.want {
				t.Errorf("got %s, want %s", hexStr, tt.want)
			}
		})
	}
}
