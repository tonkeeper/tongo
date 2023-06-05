package tlb

import (
	"fmt"
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
