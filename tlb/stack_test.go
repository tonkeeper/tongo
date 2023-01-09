package tlb

import (
	"github.com/startfellows/tongo/boc"
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
		int257 := Int257(*v)
		c := boc.NewCell()
		err := Marshal(c, int257)
		if err != nil {
			t.Fatalf("can not marshal Int257 %v", err)
		}
		var newInt1 Int257
		err = Unmarshal(c, &newInt1)
		if err != nil {
			t.Fatalf("can not unmarshal Int257 %v", err)
		}
		x := big.Int(int257)
		if x.Cmp(v) != 0 {
			t.Fatalf("converted and original big ints not equal")
		}
		x1 := big.Int(newInt1)
		if x1.Cmp(v) != 0 {
			t.Fatalf("unmarhalled and original big ints not equal")
		}
	}
	//veryBigHex, _ := hex.DecodeString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	//veryBigInt := big.NewInt(0)
	//veryBigInt.SetBytes(veryBigHex)
	//veryBigInt.Mul(veryBigInt, big.NewInt(-1))
	//_, err := Int257(veryBigInt)
	//if err == nil {
	//	t.Fatalf("big int greater than 32 bytes must throw error")
	//}
}
