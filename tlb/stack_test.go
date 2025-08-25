package tlb

import (
	"math/big"
	"testing"

	"github.com/tonkeeper/tongo/boc"
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

func Test_StackUnmarshal(t *testing.T) {
	cell, err := boc.DeserializeSinglRootBase64("te6ccgEBCAEAiwACCAAABQMBBwIJBFLlYCACBAIJBBAUsCADBAESAQAAAAAAAAA+BQGVAAAAAAAAAD6AG01jbMx0HJoxmQrpmbO2pS3Gz7smSODed3s9Y5pZXhfwALPO1uMiR5RfJrdvaHAOGoJzKy7Dj+vYwtRz8gVIyWsuBwESAf//////////BgAAAA42Mi5qc29u")
	if err != nil {
		t.Fatal(err)
	}
	var stack VmStack
	err = Unmarshal(cell, &stack)
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Init              bool
		Index             Int257
		CollectionAddress MsgAddress
		OwnerAddress      MsgAddress
		IndividualContent Any
	}
	err = stack.Unmarshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	i := big.Int(data.Index)
	if data.Init != true ||
		i.Int64() != 62 ||
		data.CollectionAddress.AddrStd.Address.Hex() != "da6b1b6663a0e4d18cc8574ccd9db5296e367dd9324706f3bbd9eb1cd2caf0bf" {
		t.Fatalf("invalid decoding")
	}

}

func Test_IntTupleUnmarshal(t *testing.T) {
	type test struct {
		name       string
		hex        string
		intsResult []int64
	}

	tests := []test{
		{
			name:       "tuple with 0 items",
			hex:        "b5ee9c7201010201000b00010c000001070000010000",
			intsResult: []int64{},
		},
		{
			name:       "tuple with 1 item",
			hex:        "b5ee9c7201010301001700020c000001070001010200000012010000000000000003",
			intsResult: []int64{3},
		},
		{
			name:       "tuple with 2 items",
			hex:        "b5ee9c7201010401002300030c000001070002010203000000120100000000000000030012010000000000000004",
			intsResult: []int64{3, 4},
		},
		{
			name:       "tuple with 3 items",
			hex:        "b5ee9c7201010601003200030c000001070003010203000002000405001201000000000000000500120100000000000000030012010000000000000004",
			intsResult: []int64{3, 4, 5},
		},
		{
			name:       "tuple with 4 items",
			hex:        "b5ee9c7201010801004100030c000001070004010203000002000405001201000000000000000602000607001201000000000000000500120100000000000000030012010000000000000004",
			intsResult: []int64{3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := boc.DeserializeSinglRootHex(tt.hex)
			if err != nil {
				t.Fatal(err)
			}
			var stack VmStack
			err = Unmarshal(cell, &stack)
			if err != nil {
				t.Fatal(err)
			}
			val := stack[0]
			if val.SumType != "VmStkTuple" {
				t.Fatalf("Stack value must be tuple, got %v", val.SumType)
			}
			tuple := val.VmStkTuple
			if int(tuple.Len) != len(tt.intsResult) {
				t.Fatalf("want %v tuple len, got %v", len(tt.intsResult), tuple.Len)
			}
			if tuple.Data == nil { // for test case with 0 values in a tuple
				return
			}
			values, err := tuple.Data.RecursiveToSlice(len(tt.intsResult))
			if err != nil {
				t.Fatal(err)
			}
			for i, v := range values {
				if v.SumType != "VmStkTinyInt" {
					t.Fatalf("want values[%v] to be VmStkTinyInt, got %v", i, v.SumType)
				}
				if v.VmStkTinyInt != tt.intsResult[i] {
					t.Fatalf("want values[%v] == %v, got %v", i, tt.intsResult[i], v.VmStkTinyInt)
				}
			}
		})
	}
}
