package abi

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"reflect"
	"testing"
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
