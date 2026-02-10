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
