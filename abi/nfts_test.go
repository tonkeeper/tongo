package abi

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
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
			assert(t, err)
			var j, j2 NFTPayload
			err = tlb.Unmarshal(c[0], &j)
			assert(t, err)
			b, err := json.Marshal(j)
			assert(t, err)
			if string(b) != tt.want {
				t.Errorf("NftPayload.MarshalJSON() = %v, want %v", string(b), tt.want)
			}
			err = json.Unmarshal(b, &j2)
			assert(t, err)
			if j.SumType != UnknownNFTOp && !reflect.DeepEqual(j, j2) {
				t.Errorf("NftPayload.UnmarshalJSON() = %v, want %v", j2, j)
			}
			c2 := boc.NewCell()
			err = tlb.Marshal(c2, j2)
			assert(t, err)
			s, err := c2.ToBocString()
			assert(t, err)
			if s != tt.boc {
				t.Errorf("JettonPayload.MarshalTLB() = %v, want %v", s, tt.boc)
			}
		})
	}
}
