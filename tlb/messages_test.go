package tlb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
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

func TestMessage_normalized_hash(t *testing.T) {
	info := struct {
		Src       MsgAddress
		Dest      MsgAddress
		ImportFee VarUInteger16
	}{
		Src: MsgAddress{
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
		Dest: MsgAddress{
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
		ImportFee: VarUInteger16(*big.NewInt(12364)),
	}
	msg := Message{
		Info: CommonMsgInfo{
			SumType:      "ExtInMsgInfo",
			ExtInMsgInfo: &info,
		},
	}
	msg.Init.Exists = true
	msg.Init.Value.IsRight = true
	msg.Init.Value.Value.Code.Exists = true
	code := boc.NewCell()
	_ = code.WriteUint(102, 32)
	msg.Init.Value.Value.Code.Value.Value = *code

	body := boc.NewCell()
	_ = body.WriteUint(200, 32)
	msg.Body.Value = Any(*body)

	if msg.Hash(true).Hex() != "dfacc0b48826e33a5a127ee1def710a449d8ce79def7c19f43e57b7996e870df" {
		t.Fatalf("invalid mesg hash %s", msg.Hash(true).Hex())
	}
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
				SumType:    "AddrExtern",
				AddrExtern: &bitstr,
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
		{
			name:       "AddrVar - all good",
			addressStr: "-43464703:B36192BFD8DB3EABFCC47AE5EF2E7F4ED6AEFB1C2E8F0A3549A66FFC45051CD3:Anycast(30,927184944)",
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
							Depth:      30,
							RewritePfx: 927184944,
						},
					},
					AddrLen:     256,
					WorkchainId: -43464703,
					Address:     mustFromFiftHex("B36192BFD8DB3EABFCC47AE5EF2E7F4ED6AEFB1C2E8F0A3549A66FFC45051CD3"),
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
			wantHash: "23ff6f150d573f64d5599a57813f991882b7b4d5ae0550ebd08ea658431e62f6",
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
			hash := fmt.Sprintf("%x", msg.Hash(true))
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
