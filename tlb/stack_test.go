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

func Test_IntMatrixTupleUnmarshal(t *testing.T) {
	type test struct {
		name             string
		hex              string
		intsResultMatrix [][]int64
	}

	tests := []test{
		{
			name:             "tuple 0x0",
			hex:              "b5ee9c7201010301001100020c000001070001010200000006070000",
			intsResultMatrix: [][]int64{{}},
		},
		{
			name:             "tuple 0x1",
			hex:              "b5ee9c7201010401001d00020c000001070001010200000106070001030012010000000000000001",
			intsResultMatrix: [][]int64{{1}},
		},
		{
			name:             "tuple 0x2",
			hex:              "b5ee9c7201010501002900020c000001070001010200000206070002030400120100000000000000010012010000000000000002",
			intsResultMatrix: [][]int64{{1, 2}},
		},
		{
			name:             "tuple 0x3",
			hex:              "b5ee9c7201010701003800020c000001070001010200000206070003030402000506001201000000000000000300120100000000000000010012010000000000000002",
			intsResultMatrix: [][]int64{{1, 2, 3}},
		},
		{
			name:             "tuple 1x0",
			hex:              "b5ee9c7201010301001200030c00000107000201020200000006070000",
			intsResultMatrix: [][]int64{{}, {}},
		},
		{
			name:             "tuple 1x1",
			hex:              "b5ee9c7201010601002f00030c000001070002010203000001060700010401060700010500120100000000000000010012010000000000000002",
			intsResultMatrix: [][]int64{{1}, {2}},
		},
		{
			name:             "tuple 1x2",
			hex:              "b5ee9c7201010801004700030c000001070002010203000002060700020405020607000206070012010000000000000001001201000000000000000200120100000000000000030012010000000000000004",
			intsResultMatrix: [][]int64{{1, 2}, {3, 4}},
		},
		{
			name:             "tuple 1x3",
			hex:              "b5ee9c7201010c01006500030c0000010700020102030000020607000304050206070003060702000809001201000000000000000302000a0b00120100000000000000060012010000000000000001001201000000000000000200120100000000000000040012010000000000000005",
			intsResultMatrix: [][]int64{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name:             "tuple 2x0",
			hex:              "b5ee9c7201010401001600030c0000010700030102030000020003030006070000",
			intsResultMatrix: [][]int64{{}, {}, {}},
		},
		{
			name:             "tuple 2x1",
			hex:              "b5ee9c7201010901004400030c000001070003010203000002000405010607000106010607000107010607000108001201000000000000000300120100000000000000010012010000000000000002",
			intsResultMatrix: [][]int64{{1}, {2}, {3}},
		},
		{
			name:             "tuple 2x2",
			hex:              "b5ee9c7201010c01006800030c000001070003010203000002000405020607000206070206070002080902060700020a0b001201000000000000000500120100000000000000060012010000000000000001001201000000000000000200120100000000000000030012010000000000000004",
			intsResultMatrix: [][]int64{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:             "tuple 2x3",
			hex:              "b5ee9c7201011201009500030c000001070003010203000002000405020607000306070206070003080902060700030a0b02000c0d001201000000000000000902000e0f0012010000000000000003020010110012010000000000000006001201000000000000000700120100000000000000080012010000000000000001001201000000000000000200120100000000000000040012010000000000000005",
			intsResultMatrix: [][]int64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{
			name:             "tuple 3x0",
			hex:              "b5ee9c7201010501001a00030c000001070004010204000002000304020004040006070000",
			intsResultMatrix: [][]int64{{}, {}, {}, {}},
		},
		{
			name:             "tuple 3x1",
			hex:              "b5ee9c7201010c01005900030c00000107000401020300000200040501060700010602000708010607000109001201000000000000000401060700010a01060700010b001201000000000000000300120100000000000000010012010000000000000002",
			intsResultMatrix: [][]int64{{1}, {2}, {3}, {4}},
		},
		{
			name:             "tuple 3x2",
			hex:              "b5ee9c7201011001008900030c000001070004010203000002000405020607000206070200080902060700020a0b0012010000000000000007001201000000000000000802060700020c0d02060700020e0f001201000000000000000500120100000000000000060012010000000000000001001201000000000000000200120100000000000000030012010000000000000004",
			intsResultMatrix: [][]int64{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
		{
			name:             "tuple 3x3",
			hex:              "b5ee9c720101180100c500030c000001070004010203000002000405020607000306070200080902060700030a0b02000c0d001201000000000000000c02060700030e0f02060700031011020012130012010000000000000009001201000000000000000a001201000000000000000b020014150012010000000000000003020016170012010000000000000006001201000000000000000700120100000000000000080012010000000000000001001201000000000000000200120100000000000000040012010000000000000005",
			intsResultMatrix: [][]int64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
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
				t.Errorf("Stack value must be tuple, got %v", val.SumType)
			}
			tuple := val.VmStkTuple
			if int(tuple.Len) != len(tt.intsResultMatrix) {
				t.Errorf("want %v tuple len, got %v", len(tt.intsResultMatrix), tuple.Len)
			}
			values, err := tuple.Data.RecursiveToSlice(len(tt.intsResultMatrix))
			if err != nil {
				t.Fatal(err)
			}
			for i, v := range values {
				if v.SumType != "VmStkTuple" {
					t.Errorf("want values[%v] to be VmStkTuple, got %v", i, v.SumType)
				}
				if int(v.VmStkTuple.Len) != len(tt.intsResultMatrix[i]) {
					t.Errorf("want %v tuple[%v] len, got %v", len(tt.intsResultMatrix[i]), i, v.VmStkTuple.Len)
				}
				if v.VmStkTuple.Data == nil { // for test case with 0 values in a tuple
					return
				}
				innerValues, err := v.VmStkTuple.Data.RecursiveToSlice(len(tt.intsResultMatrix[i]))
				if err != nil {
					t.Fatal(err)
				}
				for j, iV := range innerValues {
					if iV.SumType != "VmStkTinyInt" {
						t.Errorf("want values[%v][%v] to be VmStkTinyInt, got %v", i, j, iV.SumType)
					}
					if iV.VmStkTinyInt != tt.intsResultMatrix[i][j] {
						t.Errorf("want values[%v][%v] == %v, got %v", i, j, tt.intsResultMatrix[i][j], iV.VmStkTinyInt)
					}
				}
			}
		})
	}
}
