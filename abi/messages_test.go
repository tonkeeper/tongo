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
	data := "b5ee9c7201010101005c0000b3178d45190000000000000000a0165ceda7f1e7eda8000801b26778d9535dc370f56d2e34225d4619deb1e0f4080a7bdde114618639d1574b001643971209f95b121f5600b0fd68e03bdb7b0994f0dfdd4a9550a2a3e56e23f382"
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

func TestDecodeAndEncodeUnknownInMsgBody1(t *testing.T) {
	data := "b5ee9c7241010b0100f4000202c8010a0201620209020120030802012004070201580506004327fd889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff400432002aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac0044d4ff53355555555555555555555555555555555555555555555555555555555555550045bd002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877fa00045a3cff555555555555555555555555555555555555555555555555555555555555555580045bd5800aa6aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab8a1b898e"
	boc1, _ := boc.DeserializeBocHex(data)

	var x tlb.Hashmap[tlb.Uint16, tlb.MsgAddress]
	if err := tlb.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	curr1 := boc1[0]
	curr2 := boc2
	fmt.Println(equal(*curr1, *curr2))

	hs1, _ := boc1[0].HashString()
	hs2, _ := boc2.HashString()
	if hs1 != hs2 {
		t.Fatalf("got different result")
	}
}

func equal(curr1 boc.Cell, curr2 boc.Cell) bool {
	curr1.ResetCounters()
	curr2.ResetCounters()
	bits1 := curr1.ReadRemainingBits()
	bits2 := curr2.ReadRemainingBits()
	if bits1.BinaryString() != bits2.BinaryString() {
		//fmt.Println(bits1.BinaryString())
		//fmt.Println(bits2.BinaryString())
		return false
	}

	if len(curr1.Refs()) != len(curr2.Refs()) {
		fmt.Println(len(curr1.Refs()), len(curr2.Refs()))
		return false
	}

	for i := range curr1.Refs() {
		if !equal(*curr1.Refs()[i], *curr2.Refs()[i]) {
			curr1.Refs()[1].ResetCounters()
			curr2.Refs()[1].ResetCounters()
			bits3 := curr1.Refs()[1].ReadRemainingBits()
			bits4 := curr2.Refs()[1].ReadRemainingBits()
			if bits3.BitsAvailableForRead() == 269 {
				fmt.Println(bits3.BinaryString())
				fmt.Println(bits4.BinaryString())
			}
			return false
		}
	}

	return true
}
