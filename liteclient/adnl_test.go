package liteclient

import (
	"encoding/hex"
	"testing"
)

func TestParseAdnlAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantHex string
		wantErr bool
	}{
		{
			name:    "empty address",
			address: "",
			wantErr: true,
		},
		{
			name:    "valid .adnl",
			address: "v5xqa7wy3iv7dhphiubqx6iyh6nnav32lwgca62s6zh66xcjfvritic.adnl",
			wantHex: "7b7803f6c6d15f8cef3a28185fc8c1fcd682bbd2ec6103da97b27f7ae2496b14",
		},
		{
			name:    "valid",
			address: "v5xqa7wy3iv7dhphiubqx6iyh6nnav32lwgca62s6zh66xcjfvritic",
			wantHex: "7b7803f6c6d15f8cef3a28185fc8c1fcd682bbd2ec6103da97b27f7ae2496b14",
		},
		{
			name:    "hex",
			address: "7b7803f6c6d15f8cef3a28185fc8c1fcd682bbd2ec6103da97b27f7ae2496b14",
			wantErr: true,
		},
		{
			name:    "invalid",
			address: "v5xqa7wy3iv7dhphiubqx6iyh6nnav32lwgca62s6zh66xcjfvritid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := ParseADNLAddress(tt.address)
			if tt.wantErr && err != nil {
				return
			} else if tt.wantErr {
				t.Errorf("want error but return nil")
			}
			h, _ := hex.DecodeString(tt.wantHex)
			var res [32]byte
			copy(res[:], h)
			if res != a {
				t.Errorf("invalid address")
			}
		})
	}
}
