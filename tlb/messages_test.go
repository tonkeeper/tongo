package tlb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func byte32FromHex(x string) [32]byte {
	var address []byte
	fmt.Sscanf(x, "%x", &address)
	var result [32]byte
	copy(result[:], address)
	return result
}

func TestMsgAddress_JSON(t *testing.T) {

	bitstr := boc.NewBitString(16)
	err := bitstr.WriteBytes([]byte{255, 44})
	if err != nil {
		t.Fatalf("WriteBytes() failed: %v", err)
	}

	tests := []struct {
		name    string
		addr    MsgAddress
		want    []byte
		wantErr bool
	}{
		{
			name: "AddrNone",
			addr: MsgAddress{
				SumType:  "AddrNone",
				AddrNone: struct{}{},
			},
			want: []byte(`""`),
		},
		{
			name: "AddrStd",
			addr: MsgAddress{
				SumType: "AddrStd",
				AddrStd: struct {
					Anycast     Maybe[Anycast]
					WorkchainId int8
					Address     Bits256
				}{
					WorkchainId: -1,
					Address:     byte32FromHex("adfd5f1d28db13e50591d5c76a976c15d8ab6cad90554748ab254871390d9334"),
				},
			},
			want: []byte(`"-1:adfd5f1d28db13e50591d5c76a976c15d8ab6cad90554748ab254871390d9334"`),
		},
		{
			name: "AddrStd with Anycast",
			addr: MsgAddress{
				SumType: "AddrStd",
				AddrStd: struct {
					Anycast     Maybe[Anycast]
					WorkchainId int8
					Address     Bits256
				}{
					Anycast: Maybe[Anycast]{
						Exists: true,
						Value: Anycast{
							Depth:      16,
							RewritePfx: 9,
						},
					},
					WorkchainId: -1,
					Address:     byte32FromHex("adfd5f1d28db13e50591d5c76a976c15d8ab6cad90554748ab254871390d9334"),
				},
			},
			want: []byte(`"-1:adfd5f1d28db13e50591d5c76a976c15d8ab6cad90554748ab254871390d9334:Anycast(16,9)"`),
		},
		{
			name: "AddrExtern",
			addr: MsgAddress{
				SumType: "AddrExtern",
				AddrExtern: &struct {
					Len             Uint9
					ExternalAddress boc.BitString
				}{
					Len:             16,
					ExternalAddress: bitstr,
				},
			},
			want: []byte(`"FF2C"`),
		},
		{
			name: "AddrVar",
			addr: MsgAddress{
				SumType: "AddrVar",
				AddrVar: &struct {
					Anycast     Maybe[Anycast]
					AddrLen     Uint9
					WorkchainId int32
					Address     boc.BitString
				}{
					WorkchainId: 0,
					AddrLen:     16,
					Address:     bitstr,
				},
			},
			want: []byte(`"0:FF2C"`),
		},
		{
			name: "AddrVar with Anycast",
			addr: MsgAddress{
				SumType: "AddrVar",
				AddrVar: &struct {
					Anycast     Maybe[Anycast]
					AddrLen     Uint9
					WorkchainId int32
					Address     boc.BitString
				}{
					Anycast: Maybe[Anycast]{
						Exists: true,
						Value: Anycast{
							Depth:      16,
							RewritePfx: 8,
						},
					},
					WorkchainId: 0,
					AddrLen:     16,
					Address:     bitstr,
				},
			},
			want: []byte(`"0:FF2C:Anycast(16,8)"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.addr.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() failed: %v", err)
			}
			t.Logf("got: %v", string(got))
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want: %v got: %v", tt.want, got)
			}
			var x MsgAddress
			if err := x.UnmarshalJSON(got); err != nil {
				t.Fatalf("UnmarshalJSON() failed: %v", err)
			}
			if !reflect.DeepEqual(tt.addr, x) {
				t.Fatalf("want: \n%v got: \n%v", tt.addr, x)
			}

		})
	}
}

func mustFromFiftHex(hexStr string) boc.BitString {
	value, err := boc.BitStringFromFiftHex(hexStr)
	if err != nil {
		panic(err)
	}
	return *value
}

func TestMsgAddress_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		addressStr string
		wantAddr   *MsgAddress
	}{
		{
			name:       "AddrVar - all good",
			addressStr: "1968328842:0D4ED212D5D996B95A9FC3_:Anycast(10,374)",
			wantAddr: &MsgAddress{
				SumType: "AddrVar",
				AddrVar: &struct {
					Anycast     Maybe[Anycast]
					AddrLen     Uint9
					WorkchainId int32
					Address     boc.BitString
				}{
					Anycast: Maybe[Anycast]{
						Exists: true,
						Value: Anycast{
							Depth:      10,
							RewritePfx: 374,
						},
					},
					AddrLen:     87,
					WorkchainId: 1968328842,
					Address:     mustFromFiftHex("0D4ED212D5D996B95A9FC3_"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &MsgAddress{}
			err := addr.UnmarshalJSON([]byte(tt.addressStr))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(addr, tt.wantAddr) {
				t.Fatalf("want: %v, got: %#v\n", tt.wantAddr, addr)
			}
		})
	}
}

func TestMessage_Marshal_and_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		boc      string
		filename string
		wantHash string
	}{
		{
			name:     "ExtInMsg with body",
			boc:      "te6ccgEBAgEAqgAB4YgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEAFG8W4Jpf7AeOqfzL9vZ79mX3eM6UEBxZvN6+QmpYwXBq32QOBIrP4lF5ijGgQmZbC6KDeiiptxmTNwl5f59OAGU1NGLsixYlYAAAA2AAcAQBoYgBZQOG7qXmeA/2Tw1pLX2IkcQ5h5fxWzzcBskMJbVVRsKNaTpAAAAAAAAAAAAAAAAAAAA==",
			filename: "testdata/message-1",
			wantHash: "d5376cf6e9de8813d0640016545000e17bcc399bd654826f4fd7a3000b2fad68",
		},
		{
			name:     "IntMsg with body",
			boc:      "te6ccgEBAgEAjAABsUgALXKEDiSWLCdVuhCWy/hYz3hnzF93uwd93pYymUX+v88AGzNMlqFjDeWk/rivKqxwpBoMaThCmw7tE7othW8odIgQBycOAAYdyRAAAEQo20NHEsixYeDAAQBbBRONkQAAAAAAAAAAgBZQOG7qXmeA/2Tw1pLX2IkcQ5h5fxWzzcBskMJbVVRsKA==",
			filename: "testdata/message-2",
			wantHash: "b55e0995ab2428b7ccffa4d417ff78caca62dc4d33bc0e33b2d9bcf0c396f08c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgCell, err := boc.DeserializeSinglRootBase64(tt.boc)
			if err != nil {
				t.Fatalf("DeserializeSinglRootBase64() failed: %v", err)
			}
			var msg Message
			if err := Unmarshal(msgCell, &msg); err != nil {
				t.Fatalf("Unmarshal(() failed: %v", err)
			}
			hash := fmt.Sprintf("%x", msg.Hash())
			if hash != tt.wantHash {
				t.Fatalf("want hash: %v, got: %v", tt.wantHash, hash)
			}
			jsonMsg, err := json.MarshalIndent(msg, " ", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent(() failed: %v", err)
			}
			outputFilename := fmt.Sprintf("%v.output.json", tt.filename)
			if err := os.WriteFile(outputFilename, jsonMsg, 0644); err != nil {
				t.Fatalf("os.WriteFile() failed: %v", err)
			}
			inputFilename := fmt.Sprintf("%v.json", tt.filename)
			expected, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Fatalf("os.ReadFile() failed: %v", err)
			}
			if !bytes.Equal(expected, jsonMsg) {
				t.Fatalf("got different results")
			}
			c := boc.NewCell()
			if err := Marshal(c, msg); err != nil {
				t.Fatalf("tlb.Marshal() failed: %v", err)
			}
			base64, err := c.ToBocBase64()
			if err != nil {
				t.Fatalf("Cell.ToBocBase64() failed: %v", err)
			}
			if tt.boc != base64 {
				t.Fatalf("tlb.Marshal yields a different string")
			}
		})
	}
}
