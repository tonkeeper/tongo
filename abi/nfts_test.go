package abi

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"reflect"
	"testing"
)

func TestNftPayloadJSMarshaling(t *testing.T) {
	tests := []struct {
		name string
		boc  string
		want string
	}{
		{
			name: "Empty",
			boc:  "b5ee9c72010101010002000000",
			want: `{}`,
		},
		{
			name: "TextComment",
			boc:  "b5ee9c7201010101000b0000120000000048656c6c6f",
			want: `{"SumType":"TextComment","OpCode":0,"Value":{"Text":"Hello"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := boc.DeserializeBocHex(tt.boc)
			require.NoError(t, err)
			var j, j2 NFTPayload
			err = tlb.Unmarshal(c[0], &j)
			require.NoError(t, err)
			b, err := json.Marshal(j)
			require.NoError(t, err)
			if string(b) != tt.want {
				t.Errorf("NftPayload.MarshalJSON() = %v, want %v", string(b), tt.want)
			}
			err = json.Unmarshal(b, &j2)
			require.NoError(t, err)
			if j.SumType != UnknownNFTOp && !reflect.DeepEqual(j, j2) {
				t.Errorf("NftPayload.UnmarshalJSON() = %v, want %v", j2, j)
			}
			c2 := boc.NewCell()
			err = tlb.Marshal(c2, j2)
			require.NoError(t, err)
			s, err := c2.ToBocString()
			require.NoError(t, err)
			if s != tt.boc {
				t.Errorf("JettonPayload.MarshalTLB() = %v, want %v", s, tt.boc)
			}
		})
	}
}
