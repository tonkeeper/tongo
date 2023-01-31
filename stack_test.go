package tongo

import (
	"encoding/hex"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"math/big"
	"testing"
)

func TestInt257(t *testing.T) {
	vals := []*big.Int{
		big.NewInt(-1256),
		big.NewInt(0),
		big.NewInt(546665),
	}

	for _, v := range vals {
		int257, err := Int257FromBigInt(v)
		if err != nil {
			t.Fatalf("can not encode Int257 %v", err)
		}
		c := boc.NewCell()
		err = tlb.Marshal(c, int257)
		if err != nil {
			t.Fatalf("can not marshal Int257 %v", err)
		}
		var newInt1 Int257
		err = tlb.Unmarshal(c, &newInt1)
		if err != nil {
			t.Fatalf("can not unmarshal Int257 %v", err)
		}
		if int257.BigInt().Cmp(v) != 0 {
			t.Fatalf("converted and original big ints not equal")
		}
		if newInt1.BigInt().Cmp(v) != 0 {
			t.Fatalf("unmarhalled and original big ints not equal")
		}
	}
	veryBigHex, _ := hex.DecodeString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	veryBigInt := big.NewInt(0)
	veryBigInt.SetBytes(veryBigHex)
	veryBigInt.Mul(veryBigInt, big.NewInt(-1))
	_, err := Int257FromBigInt(veryBigInt)
	if err == nil {
		t.Fatalf("big int greater than 32 bytes must throw error")
	}
}
