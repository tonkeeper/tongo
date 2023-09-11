package abi

import (
	"encoding/json"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"reflect"
	"testing"
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestJettonPayloadJSMarshaling(t *testing.T) {
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
		{
			name: "StonfiSwap",
			boc:  "b5ee9c7201010101006f0000d92593856180022a16a3164c4d5aa3133f3110ff10496e00ca8ac8abeffc5027e024d33480c3e8227f36ff000b3ced6e32247945f26b76f68700e1a82732b2ec38febd8c2d473f20548c96b2f001b23bbd527c00fd3d0874026bb7941f8fdd4d599ed8e7a63f426d5723f0388f2e",
			want: `{"SumType":"StonfiSwap","OpCode":630424929,"Value":{"TokenWallet":"0:1150b518b2626ad51899f98887f8824b70065456455f7fe2813f012699a4061f","MinOut":"289381247","ToAddress":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","ReferralAddress":"0:6c8eef549f003f4f421d009aede507e3f7535667b639e98fd09b55c8fc0e23cb"}}`,
		},
		{
			name: "incomplete cell",
			boc:  "b5ee9c72010101010004000004053d",
			want: `{"SumType":"Cell","Value":"b5ee9c72010101010004000004053d"}`,
		},
		{
			name: "unknown payload",
			boc:  "b5ee9c720101010100120000206f6c6f6c6f6c6f6c6f74726f6c6f6c6f",
			want: `{"SumType":"Cell","OpCode":1869377388,"Value":"b5ee9c720101010100120000206f6c6f6c6f6c6f6c6f74726f6c6f6c6f"}`,
		},
		{
			name: "dedust swap",
			boc:  "b5ee9c72010102010034000153e3a0d482801c2bfd75b91e3b5a9402036fdb733e5dc096f4c773a328dd33586324e0cdb4bd2307b6a7400100090000000002",
			want: `{"SumType":"DedustSwap","OpCode":3818968194}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := boc.DeserializeBocHex(tt.boc)
			assert(t, err)
			var j, j2 JettonPayload
			err = tlb.Unmarshal(c[0], &j)
			assert(t, err)
			b, err := json.Marshal(j)
			assert(t, err)
			if string(b) != tt.want {
				t.Errorf("JettonPayload.MarshalJSON() = %v, want %v", string(b), tt.want)
			}
			err = json.Unmarshal(b, &j2)
			assert(t, err)
			if j.SumType != UnknownJettonOp && !reflect.DeepEqual(j, j2) {
				t.Errorf("JettonPayload.UnmarshalJSON() = %v, want %v", j2, j)
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
