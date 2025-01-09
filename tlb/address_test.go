package tlb

import (
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
)

func mustHexToBits256(s string) Bits256 {
	value, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	var bits Bits256
	copy(bits[:], value)
	return bits
}

func mustToAddress(s string) AddressWithWorkchain {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		panic("invalid address format")
	}
	workchain, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		panic(err)
	}
	return AddressWithWorkchain{
		Workchain: int8(workchain),
		Address:   mustHexToBits256(parts[1]),
	}
}

func TestAddressWithWorkchain_Compare(t *testing.T) {
	tests := []struct {
		name  string
		addr  AddressWithWorkchain
		other any
		want  int
	}{
		{
			name:  "all good",
			addr:  mustToAddress("0:0000000000000000000000000000000000000000000000000000000000000000"),
			other: mustToAddress("0:0769ffdea3d8261cb8844691f963979baffcf8a57e0dcac0263cc7076bd4976a"),
			want:  -1,
		},
		{
			name:  "all good",
			addr:  mustToAddress("0:0769ffdea3d8261cb8844691f963979baffcf8a57e0dcac0263cc7076bd4976a"),
			other: mustToAddress("-1:01b573bd6dc4cc5e383d6e08af2a1e258499995903cfebfadbf6f7e39533f914"),
			want:  -1,
		},
		{
			name:  "all good",
			addr:  mustToAddress("-1:fc3d252d2b2fd4f8964348d50da8de5c56c9fd39126a4bddcbe8344cf476eca1"),
			other: mustToAddress("-1:01b573bd6dc4cc5e383d6e08af2a1e258499995903cfebfadbf6f7e39533f914"),
			want:  1,
		},
		{
			name:  "equal",
			addr:  mustToAddress("-1:fc3d252d2b2fd4f8964348d50da8de5c56c9fd39126a4bddcbe8344cf476eca1"),
			other: mustToAddress("-1:fc3d252d2b2fd4f8964348d50da8de5c56c9fd39126a4bddcbe8344cf476eca1"),
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp, ok := tt.addr.Compare(tt.other)
			if !ok {
				t.Errorf("Compare() gotOk = %v, want %v", ok, true)
			}
			if cmp != tt.want {
				t.Errorf("Compare() got = %v, want %v", cmp, tt.want)
			}

			reverseCmd, ok := tt.other.(AddressWithWorkchain).Compare(tt.addr)
			if !ok {
				t.Errorf("Compare() gotOk = %v, want %v", ok, true)
			}
			if reverseCmd != -tt.want {
				t.Errorf("Compare() got = %v, want %v", reverseCmd, -tt.want)
			}
		})
	}
}
