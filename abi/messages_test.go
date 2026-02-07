package abi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestEncodeAndDecodeInMsgBody(t *testing.T) {
	value := InMsgBody{
		SumType: "JettonTransfer",
		OpCode:  pointer(uint32(260734629)),
		Value: JettonTransferMsgBody{
			QueryId:             0,
			Amount:              mustToVarUInteger16("1000000"),
			Destination:         pointer(ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d220")).ToMsgAddress(),
			ResponseDestination: pointer(ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d220")).ToMsgAddress(),
			CustomPayload:       nil,
			ForwardTonAmount:    mustToVarUInteger16("300000000"),
			ForwardPayload: tlb.EitherRef[JettonPayload]{
				IsRight: false,
			},
		},
	}

	msg := boc.NewCell()
	if err := tlb.Marshal(msg, value); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	var b InMsgBody
	if err := tlb.Unmarshal(msg, &b); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(b, value) {
		t.Fatalf("got different result")
	}
}

func TestDecodeAndEncodeUnknownInMsgBody(t *testing.T) {
	data := "b5ee9c720101010100350000650000022800000000000000005012a05f200800c46663fd6592a0ca9ff0f04ed928fe74019bac1756d06092c14464fb3ce8d373"
	boc1, _ := boc.DeserializeBocHex(data)

	var x InMsgBody
	if err := tlb.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}

// AddressKey is an address that can be used as a key in hashmap
type AddressKey struct {
	Prefix    tlb.Uint3
	Workchain tlb.Uint8
	Hash      tlb.Bits256
}

func (a AddressKey) Compare(other any) (int, bool) {
	otherAddr, ok := other.(AddressKey)
	if !ok {
		return 0, false
	}
	if !a.Prefix.Equal(otherAddr.Prefix) {
		return a.Prefix.Compare(otherAddr.Prefix)
	}
	if !a.Workchain.Equal(otherAddr.Workchain) {
		return a.Workchain.Compare(otherAddr.Workchain)
	}
	if !a.Hash.Equal(otherAddr.Hash) {
		return a.Hash.Compare(otherAddr.Hash)
	}
	return 0, true
}

func (a AddressKey) Equal(other any) bool {
	if otherAddr, ok := other.(AddressKey); ok {
		return a.Prefix.Equal(otherAddr.Prefix) && a.Workchain.Equal(otherAddr.Workchain) &&
			a.Hash.Equal(otherAddr.Hash)
	}
	return false
}

func (a AddressKey) FixedSize() int {
	return 267
}

func TestDecodeAndEncodeUnknownInMsgBody1(t *testing.T) {
	data := "b5ee9c720101050100940002058170020102004bbfc437d1533afba1a3d846ea30352a4dc7edcc414e59415e0ca4e9d897dd76cee6006061a82002012003040057bfb113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe307a120307a120e8c35040005fbf90e545323c7acb7102653c073377f7e3c67f122eb94d430a250739f109d4a57d50161978f1b50161978f1b80c35040"
	boc1, _ := boc.DeserializeBocHex(data)

	var x tlb.Hashmap[AddressKey, AssetData]
	if err := tlb.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}
