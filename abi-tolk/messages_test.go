package abitolk

import (
	"fmt"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func TestDecodeAndEncodeUnknownInMsgBody(t *testing.T) {
	data := "b5ee9c720101010100570000a9178d4519000000006948345d40bab36908010e3f58e634ebf4cca5296ea1ade141563bf278c85e3ff08161e79367808987cf0021c7eb1cc69d7e9994a52dd435bc282ac77e4f190bc7fe102c3cf26cf01130f9c405"
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
