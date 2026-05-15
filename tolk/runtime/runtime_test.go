package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

const jsonFilesPath = "testdata/json/"

func TestRuntime_SimpleValues(t *testing.T) {
	abi := loadTestABI(t, "testdata/abi/benchmark_types.json")

	for _, tt := range []struct {
		name         string
		expectedJson string
		typeName     string
		cellHex      string
		expected     *Value
		assert       func(*testing.T, *Value)
	}{
		{
			name:         "small int",
			expectedJson: `{"value":-35132}`,
			typeName:     "BenchmarkSmallInt",
			cellHex:      "b5ee9c72410101010005000006ff76c41616db06",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetSmallInt()
				if !ok {
					t.Errorf("v.GetSmallInt() not successeded")
				}
				if val != -35132 {
					t.Errorf("val != -35132, got %v", val)
				}
			},
		},
		{
			name:         "big int",
			expectedJson: `{"value":"-3513294376431"}`,
			typeName:     "BenchmarkBigInt",
			cellHex:      "b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetBigInt()
				if !ok {
					t.Errorf("v.GetBigInt() not successeded")
				}
				if val.Cmp(big.NewInt(-3513294376431)) != 0 {
					t.Errorf("val != -3513294376431, got %v", val)
				}
			},
		},
		{
			name:         "small uint",
			expectedJson: `{"value":934}`,
			typeName:     "BenchmarkSmallUint",
			cellHex:      "b5ee9c7241010101000900000d00000000001d34e435eafd",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetSmallUInt()
				if !ok {
					t.Errorf("v.GetSmallUInt() not successeded")
				}
				if val != 934 {
					t.Errorf("val != 934, got %v", val)
				}
			},
		},
		{
			name:         "big uint",
			expectedJson: `{"value":"351329437643124"}`,
			typeName:     "BenchmarkBigUint",
			cellHex:      "b5ee9c7201010101002200004000000000000000000000000000000000000000000000000000013f8842547174",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetBigUInt()
				if !ok {
					t.Errorf("v.GetBigUInt() not successeded")
				}
				if val.Cmp(big.NewInt(351329437643124)) != 0 {
					t.Errorf("val != 351329437643124, got %v", val.String())
				}
			},
		},
		{
			name:         "var int",
			expectedJson: `{"value":"825432"}`,
			typeName:     "BenchmarkVarInt16",
			cellHex:      "b5ee9c7241010101000600000730c98588449b6923",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetVarInt()
				if !ok {
					t.Errorf("v.GetVarInt() not successeded")
				}
				if val.Cmp(big.NewInt(825432)) != 0 {
					t.Errorf("val != 825432, got %v", val.String())
				}
			},
		},
		{
			name:         "var uint",
			expectedJson: `{"value":"9451236712"}`,
			typeName:     "BenchmarkVarUint32",
			cellHex:      "b5ee9c7241010101000800000b28119ab36b44d3a86c0f",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetVarUInt()
				if !ok {
					t.Errorf("v.GetVarUInt() not successeded")
				}
				if val.Cmp(big.NewInt(9451236712)) != 0 {
					t.Errorf("val != 9451236712, got %v", val.String())
				}
			},
		},
		{
			name:         "bits",
			expectedJson: `{"value":"313233"}`,
			typeName:     "BenchmarkBits24",
			cellHex:      "b5ee9c7241010101000500000631323318854035",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetBits()
				if !ok {
					t.Errorf("v.GetBits() not successeded")
				}
				if !bytes.Equal(val.Buffer(), []byte{49, 50, 51}) {
					t.Errorf("val != {49, 50, 51}, got %v", val)
				}
			},
		},
		{
			name:         "coins",
			expectedJson: `{"value":"921464321"}`,
			typeName:     "BenchmarkCoins",
			cellHex:      "b5ee9c72410101010007000009436ec6e0189ebbd7f4",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetCoins()
				if !ok {
					t.Errorf("v.GetCoins() not successeded")
				}
				if val.Cmp(big.NewInt(921464321)) != 0 {
					t.Errorf("val != 921464321, got %v", val)
				}
			},
		},
		{
			name:         "bool",
			expectedJson: `{"value":false}`,
			typeName:     "BenchmarkBool",
			cellHex:      "b5ee9c7241010101000300000140f6d24034",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetBool()
				if !ok {
					t.Errorf("v.GetBool() not successeded")
				}
				if val {
					t.Error("val is true")
				}
			},
		},
		{
			name:         "cell",
			expectedJson: `{"value":"b5ee9c720101010100060000080000007b"}`,
			typeName:     "BenchmarkCell",
			cellHex:      "b5ee9c724101020100090001000100080000007ba52a3292",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetCell()
				if !ok {
					t.Errorf("v.GetCell() not successeded")
				}
				assertCellHash(t, &val, "644e68a539c5107401d194bc82169cbf0ad1635796891551e0750705ab2d74ae")
			},
		},
		{
			name:         "remaining",
			expectedJson: `{"value":{"isRef":false,"value":"b5ee9c7201010101000900000dc0800000000ab8"}}`,
			typeName:     "BenchmarkRemaining",
			cellHex:      "b5ee9c7241010101000900000dc0800000000ab8d04726e4",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetRemaining()
				if !ok {
					t.Errorf("v.GetCell() not successeded")
				}
				assertCellHash(t, &val.Value, "f1c4e07fbd1786411c2caa9ac9f5d7240aa2007a2a1d5e5ac44f8a168cd4e36b")
			},
		},
		{
			name:         "internal address",
			expectedJson: `{"value":"0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8"}`,
			typeName:     "BenchmarkAddress",
			cellHex:      "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAddress()
				if !ok {
					t.Errorf("v.GetAddress() not successeded")
				}
				if val.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.ToRaw())
				}
			},
		},
		{
			name:         "not exists optional address",
			expectedJson: `{"value":""}`,
			typeName:     "BenchmarkOptionalAddress",
			cellHex:      "b5ee9c724101010100030000012094418655",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetOptionalAddress()
				if !ok {
					t.Errorf("v.GetOptionalAddress() not successeded")
				}
				if val.SumType != SumTypeNoneAddress {
					t.Errorf("val.GetAddress() != none address")
				}
			},
		},
		{
			name:         "exists optional address",
			expectedJson: `{"value":"0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8"}`,
			typeName:     "BenchmarkOptionalAddress",
			cellHex:      "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetOptionalAddress()
				if !ok {
					t.Errorf("v.GetOptionalAddress() not successeded")
				}
				if val.SumType != SumTypeInternalAddress {
					t.Errorf("val.GetAddress() != InternalAddress, got %v", val.SumType)
				}
				if val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
				}
			},
		},
		{
			name:         "external address",
			expectedJson: `{"value":"4142"}`,
			typeName:     "BenchmarkExternalAddress",
			cellHex:      "b5ee9c7241010101000600000742082850fcbd94fd",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != SumTypeExternalAddress {
					t.Errorf("val.GetExternalAddress() != ExternalAddress, got %v", val.SumType)
				}
				if val.ExternalAddress.Len != 16 {
					t.Errorf("val.GetExternalAddress().Len != 16, got %v", val.ExternalAddress.Len)
				}
				if !bytes.Equal(val.ExternalAddress.Address.Buffer(), []byte{65, 66}) {
					t.Errorf("val.GetExternalAddress() != {65, 66}, got %v", val.ExternalAddress.Address.Buffer())
				}
			},
		},
		{
			name:         "any none address",
			expectedJson: `{"value":""}`,
			typeName:     "BenchmarkAnyAddress",
			cellHex:      "b5ee9c724101010100030000012094418655",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != SumTypeNoneAddress {
					t.Errorf("val.GetAddress() != none address")
				}
			},
		},
		{
			name:         "any internal address",
			expectedJson: `{"value":"0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8"}`,
			typeName:     "BenchmarkAnyAddress",
			cellHex:      "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != SumTypeInternalAddress {
					t.Errorf("val.GetAddress() != InternalAddress, got %v", val.SumType)
				}
				if val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
				}
			},
		},
		{
			name:         "any external address",
			expectedJson: `{"value":"4142"}`,
			typeName:     "BenchmarkAnyAddress",
			cellHex:      "b5ee9c7241010101000600000742082850fcbd94fd",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != SumTypeExternalAddress {
					t.Errorf("val.GetExternalAddress() != ExternalAddress, got %v", val.SumType)
				}
				if val.ExternalAddress.Len != 16 {
					t.Errorf("val.GetExternalAddress().Len != 16, got %v", val.ExternalAddress.Len)
				}
				if !bytes.Equal(val.ExternalAddress.Address.Buffer(), []byte{65, 66}) {
					t.Errorf("val.GetExternalAddress() != {65, 66}, got %v", val.ExternalAddress.Address.Buffer())
				}
			},
		},
		{
			name:         "any var address",
			expectedJson: `{"value":"0:AB"}`,
			typeName:     "BenchmarkAnyAddress",
			cellHex:      "b5ee9c7241010101000900000dc0800000000ab8d04726e4",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != SumTypeVarAddress {
					t.Errorf("val.GetAddress() != VarAddress")
				}
				if val.VarAddress.Len != 8 {
					t.Errorf("val.VarAddress.Len != 8, got %v", val.VarAddress.Len)
				}
				if val.VarAddress.Workchain != 0 {
					t.Errorf("val.VarAddress.Workchain != 0, got %v", val.VarAddress.Workchain)
				}
				if !bytes.Equal(val.VarAddress.Address.Buffer(), []byte{171}) {
					t.Errorf("val.VarAddress.Address != {171}, got %v", val.VarAddress.Address.Buffer())
				}
			},
		},
		{
			name:         "not exists nullable",
			expectedJson: `{"value":{"isExists":false}}`,
			typeName:     "BenchmarkNullableRemaining",
			cellHex:      "b5ee9c7241010101000300000140f6d24034",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetOptionalValue()
				if !ok {
					t.Errorf("v.GetOptionalValue() not successeded")
				}
				if val.IsExists {
					t.Errorf("v.GetOptionalValue() is exists")
				}
			},
		},
		{
			name:         "exists nullable",
			expectedJson: `{"value":{"isExists":true,"value":"b5ee9c7201010101000700000900000c0ae0"}}`,
			typeName:     "BenchmarkNullableCell",
			cellHex:      "b5ee9c7241010201000b000101c001000900000c0ae007880db9",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetOptionalValue()
				if !ok {
					t.Errorf("v.GetOptionalValue() not successeded")
				}
				if !val.IsExists {
					t.Errorf("v.GetOptionalValue() != exists")
				}
				innerVal, ok := val.Val.GetCell()
				if !ok {
					t.Errorf("v.GetOptionalValue().GetCell() not successeded")
				}
				assertCellHash(t, &innerVal, "df05386a55563049a4834a4cc1ec0dc22f3dcb63c04f7258ae475c5d28981773")
			},
		},
		{
			name:         "ref",
			expectedJson: `{"value":"1233212"}`,
			typeName:     "BenchmarkRefInt65",
			cellHex:      "b5ee9c7241010201000e000100010011000000000009689e40e150b4c5",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetRefValue()
				if !ok {
					t.Errorf("v.GetRefValue() not successeded")
				}
				innerVal, ok := val.GetBigInt()
				if !ok {
					t.Errorf("v.GetRefValue().GetBigInt() not successeded")
				}
				if innerVal.Cmp(big.NewInt(1233212)) != 0 {
					t.Errorf("v.GetRefValue().GetBigInt() != 1233212, got %v", innerVal.String())
				}
			},
		},
		{
			name:         "empty tensor",
			expectedJson: `{"value":[]}`,
			typeName:     "BenchmarkEmptyTensor",
			cellHex:      "b5ee9c724101010100020000004cacb9cd",
			assert: func(t *testing.T, v *Value) {
				val, ok := v.GetTensor()
				if !ok {
					t.Errorf("v.GetTensor() not successeded")
				}
				if len(val) != 0 {
					t.Errorf("v.GetTensor() != empty")
				}
			},
		},
		{
			name:         "tensor",
			expectedJson: `{"value":["4325",true,"1000000000",[-342,{"isExists":true,"value":0}],"-9304000000"]}`,
			typeName:     "BenchmarkNotEmptyTensor",
			cellHex:      "b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117",
			expected: &Value{SumType: SumTypeTensor, Tensor: &(TensorValues{
				{SumType: SumTypeBigUint, BigUint: ptr(BigUInt(*big.NewInt(4325)))},
				{SumType: SumTypeBool, Bool: ptr(BoolValue(true))},
				{SumType: SumTypeCoins, Coins: ptr(CoinsValue(*big.NewInt(1_000_000_000)))},
				{SumType: SumTypeTensor, Tensor: &TensorValues{
					{SumType: SumTypeSmallInt, SmallInt: ptr(Int64(-342))},
					{SumType: SumTypeOptionalValue, OptionalValue: ptr((OptValue{IsExists: true, Val: Value{SumType: SumTypeSmallInt, SmallInt: ptr(Int64(0))}}))},
				}},
				{SumType: SumTypeVarInt, VarInt: ptr(VarInt(*big.NewInt(-9_304_000_000)))},
			})},
		},
		{
			name:         "int key map",
			expectedJson: `{"value":{"123":true}}`,
			typeName:     "BenchmarkSmallIntKeyMap",
			cellHex:      "b5ee9c7241010201000c000101c001000ba00000007bc09a662c32",
			expected: &Value{SumType: SumTypeMap, Map: &(MapValue{
				keys:   []Value{{SumType: SumTypeSmallInt, SmallInt: ptr(Int64(123))}},
				values: []Value{{SumType: SumTypeBool, Bool: ptr(BoolValue(true))}},
				len:    1,
			})},
		},
		{
			name: "uint key map",
			expectedJson: `{"value":{
				"14":"0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe",
				"23":"-1:5555555555555555555555555555555555555555555555555555555555555555"
			}}`,
			typeName: "BenchmarkSmallUintKeyMap",
			cellHex:  "b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8",
			expected: &Value{SumType: SumTypeMap, Map: &(MapValue{
				keys: []Value{
					{SumType: SumTypeSmallUint, SmallUint: ptr(UInt64(14))},
					{SumType: SumTypeSmallUint, SmallUint: ptr(UInt64(23))},
				},
				values: []Value{
					{SumType: SumTypeInternalAddress, InternalAddress: ptr(mustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))},
					{SumType: SumTypeInternalAddress, InternalAddress: ptr(mustParseAddr("Ef9VVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVbxn"))},
				},
				len: 2,
			})},
		},
		{
			name:         "big int key map",
			expectedJson: `{"value":{"2337412":"b5ee9c7201010101000800000b000000001ab0"}}`,
			typeName:     "BenchmarkBigUintKeyMap",
			cellHex:      "b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9",
			expected: &Value{SumType: SumTypeMap, Map: &(MapValue{
				keys:   []Value{{SumType: SumTypeBigUint, BigUint: ptr(BigUInt(*big.NewInt(2337412)))}},
				values: []Value{{SumType: SumTypeCell, Cell: ptr(Any(*boc.MustDeserializeSinglRootHex("b5ee9c7201010101000800000b000000001ab0")))}},
				len:    1,
			})},
		},
		{
			name:         "bits key map",
			expectedJson: `{"value":{"4142":{"124":["0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe","1000000000"]}}}`,
			typeName:     "BenchmarkBitsKeyMap",
			cellHex:      "b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86",
			expected: &Value{SumType: SumTypeMap, Map: &(MapValue{
				keys: []Value{{SumType: SumTypeBits, Bits: ptr(Bits(makeBitString(t, 16, []byte{65, 66})))}},
				values: []Value{{SumType: SumTypeMap, Map: &(MapValue{
					keys: []Value{{SumType: SumTypeSmallInt, SmallInt: ptr(Int64(124))}},
					values: []Value{{SumType: SumTypeTensor, Tensor: &(TensorValues{
						{SumType: SumTypeInternalAddress, InternalAddress: ptr(mustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))},
						{SumType: SumTypeCoins, Coins: ptr(CoinsValue(*big.NewInt(1_000_000_000)))},
					})}},
					len: 1,
				})}},
				len: 1,
			})},
		},
		{
			name:         "address key map",
			expectedJson: `{"value":{"0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe":"10000000000"}}`,
			typeName:     "BenchmarkAddressKeyMap",
			cellHex:      "b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e",
			expected: &Value{SumType: SumTypeMap, Map: &(MapValue{
				keys: []Value{{
					SumType:         SumTypeInternalAddress,
					InternalAddress: ptr(mustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")),
				}},
				values: []Value{{SumType: SumTypeCoins, Coins: ptr(CoinsValue(*big.NewInt(10_000_000_000)))}},
				len:    1,
			})},
		},
		{
			name:         "union with dec prefix",
			expectedJson: `{"value":"124432123"}`,
			typeName:     "BenchmarkDecUnion",
			cellHex:      "b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029",
			expected: &Value{
				SumType: SumTypeUnion,
				Union: &UnionValue{
					Prefix: Prefix{Len: 1, Prefix: 1},
					Val:    Value{SumType: SumTypeBigInt, BigInt: ptr(BigInt(*big.NewInt(124432123)))},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			currCell, err := boc.DeserializeBocHex(tt.cellHex)
			if err != nil {
				t.Fatal(err)
			}
			tyIdx := testTyIdx(t, abi, tt.typeName)
			decoder := NewDecoder(abi)
			v, err := decoder.UnmarshalTyIdx(currCell[0], tyIdx)
			if err != nil {
				t.Fatal(err)
			}

			assertVal := v
			if alias, ok := v.GetAlias(); ok {
				assertVal = &alias
			}
			if tt.expected != nil {
				assertRuntimeValueEqual(t, assertVal, *tt.expected)
			} else {
				tt.assert(t, assertVal)
			}
			if err := compareExpectedInlineJson(tt.expectedJson, *assertVal); err != nil {
				t.Fatal(err)
			}

			encoder := NewEncoder(abi)
			newCell, err := encoder.MarshalTyIdx(v, tyIdx)
			if err != nil {
				t.Fatal(err)
			}
			assertSameCellHash(t, currCell[0], newCell)
		})
	}
}

func mustParseAddr(addr string) InternalAddress {
	return InternalAddressFromTLB(tongo.MustParseAddress(addr).ID)
}

func makeBitString(t *testing.T, size int, data []byte) boc.BitString {
	bitsKey := boc.NewBitString(size)
	if err := bitsKey.WriteBytes(data); err != nil {
		t.Fatal(err)
	}
	return bitsKey
}

func assertCellHash(t testing.TB, cell *boc.Cell, expected string) {
	t.Helper()
	hs, err := cell.HashString()
	if err != nil {
		t.Fatal(err)
	}
	if hs != expected {
		t.Errorf("cell hash != %s, got %v", expected, hs)
	}
}

func assertSameCellHash(t testing.TB, oldCell, newCell *boc.Cell) {
	t.Helper()
	oldHs, err := oldCell.HashString()
	if err != nil {
		t.Fatal(err)
	}
	newHs, err := newCell.HashString()
	if err != nil {
		t.Fatal(err)
	}
	if oldHs != newHs {
		t.Errorf("input and output cells are different")
	}
}

func assertRuntimeValueEqual(t testing.TB, actual *Value, expected Value) {
	t.Helper()
	if actual.Equal(expected) {
		return
	}
	actualJSON, actualErr := json.Marshal(actual)
	expectedJSON, expectedErr := json.Marshal(expected)
	if actualErr == nil && expectedErr == nil {
		t.Fatalf("unexpected decoded value:\nexpected: %s\nactual:   %s", expectedJSON, actualJSON)
	}
	t.Fatalf("unexpected decoded value: expected %#v, actual %#v", expected, *actual)
}

func TestRuntime_UnmarshalABIValues(t *testing.T) {
	for _, tt := range []struct {
		name         string
		jsonName     string
		abiFile      string
		typeName     string
		cellHex      string
		expected     *Value
		customUnpack func(testing.TB, parser.ContractABI, *Decoder)
	}{
		{
			name:     "union with bin prefix",
			jsonName: "union_with_bin_prefix",
			abiFile:  "testdata/abi/bin_union.json",
			typeName: "AddressWithPrefix | MapWithPrefix | CellWithPrefix",
			cellHex:  "b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650",
			expected: &Value{
				SumType: SumTypeUnion,
				Union: &UnionValue{
					Prefix: Prefix{Len: 3, Prefix: 3},
					Val: Value{
						SumType: SumTypeStruct,
						Struct: &Struct{
							hasPrefix:  true,
							name:       "MapWithPrefix",
							prefix:     Prefix{Len: 3, Prefix: 3},
							fieldNames: []string{"v"},
							fields: map[string]Value{
								"v": {
									SumType: SumTypeMap,
									Map: &(MapValue{
										keys: []Value{{
											SumType:         SumTypeInternalAddress,
											InternalAddress: ptr(mustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")),
										}},
										values: []Value{{
											SumType: SumTypeCoins,
											Coins:   ptr(CoinsValue(*big.NewInt(43213412))),
										}},
										len: 1,
									}),
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "union with hex prefix",
			jsonName: "union_with_hex_prefix",
			abiFile:  "testdata/abi/hex_union.json",
			typeName: "UInt66WithPrefix | UInt33WithPrefix | UInt4WithPrefix",
			cellHex:  "b5ee9c7241010101000b000011deadbeef00000000c0d75977b9",
			expected: &Value{
				SumType: SumTypeUnion,
				Union: &UnionValue{
					Prefix: Prefix{Len: 32, Prefix: 0xdeadbeef},
					Val: Value{
						SumType: SumTypeStruct,
						Struct: &Struct{
							hasPrefix:  true,
							name:       "UInt33WithPrefix",
							prefix:     Prefix{Len: 32, Prefix: 0xdeadbeef},
							fieldNames: []string{"v"},
							fields: map[string]Value{
								"v": {
									SumType:   SumTypeSmallUint,
									SmallUint: ptr(UInt64(1)),
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "refs from alias",
			jsonName: "a_lot_refs_from_alias",
			abiFile:  "testdata/abi/refs.json",
			typeName: "GoodNamingForMsg",
			cellHex:  "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
		},
		{
			name:     "refs from struct",
			jsonName: "a_lot_refs_from_struct",
			abiFile:  "testdata/abi/refs.json",
			typeName: "ManyRefsMsg",
			cellHex:  "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
		},
		{
			name:     "generics from struct",
			jsonName: "a_lot_generics_from_struct",
			abiFile:  "testdata/abi/generics.json",
			typeName: "ManyRefsMsg<uint16>",
			cellHex:  "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
		},
		{
			name:     "generics from alias",
			jsonName: "a_lot_generics_from_alias",
			abiFile:  "testdata/abi/generics.json",
			typeName: "GoodNamingForMsg<uint16>",
			cellHex:  "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
		},
		{
			name:     "default values",
			jsonName: "a_lot_generics_with_default_values",
			abiFile:  "testdata/abi/default_values.json",
			typeName: "DefaultTest",
			cellHex:  "b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde",
		},
		{
			name:     "numbers",
			jsonName: "a_lot_numbers",
			abiFile:  "testdata/abi/numbers.json",
			typeName: "Numbers",
			cellHex:  "b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014",
		},
		{
			name:     "random fields",
			jsonName: "a_lot_random_fields",
			abiFile:  "testdata/abi/random_fields.json",
			typeName: "RandomFields",
			cellHex:  "b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb",
		},
		{
			name:     "alias with custom unpack",
			jsonName: "custom_pack_unpack",
			abiFile:  "testdata/abi/custom_pack_unpack.json",
			typeName: "MyAlias",
			cellHex:  "b5ee9c724101010100470000890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043b9aca00886e91196",
			customUnpack: func(t testing.TB, abi parser.ContractABI, decoder *Decoder) {
				decoder.WithCustomUnpackResolver(func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
					if err := cell.Skip(512); err != nil {
						return fmt.Errorf("failed to 512 bits from alias")
					}
					val, err := decoder.UnmarshalTyIdx(cell, testTyIdx(t, abi, "My"))
					if err != nil {
						return fmt.Errorf("failed to unmarshal alias coins")
					}
					*value = AliasValue(*val)
					return nil
				})
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			abi := loadTestABI(t, tt.abiFile)
			tyIdx := testTyIdx(t, abi, tt.typeName)

			currCell, err := boc.DeserializeBocHex(tt.cellHex)
			if err != nil {
				t.Fatal(err)
			}
			decoder := NewDecoder(abi)
			if tt.customUnpack != nil {
				tt.customUnpack(t, abi, decoder)
			}
			v, err := decoder.UnmarshalTyIdx(currCell[0], tyIdx)
			if err != nil {
				t.Fatal(err)
			}
			if tt.expected != nil {
				assertRuntimeValueEqual(t, v, *tt.expected)
			}
			if err = compareExpectedJson(tt.jsonName, *v); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRuntime_MarshalABIValues(t *testing.T) {
	for _, tt := range []struct {
		name     string
		abiFile  string
		typeName string
		cellHex  string
	}{
		{
			name:     "union with bin prefix",
			abiFile:  "testdata/abi/bin_union.json",
			typeName: "AddressWithPrefix | MapWithPrefix | CellWithPrefix",
			cellHex:  "b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650",
		},
		{
			name:     "union with hex prefix",
			abiFile:  "testdata/abi/hex_union.json",
			typeName: "UInt66WithPrefix | UInt33WithPrefix | UInt4WithPrefix",
			cellHex:  "b5ee9c7241010101000b000011deadbeef00000000c0d75977b9",
		},
		{
			name:     "refs from alias",
			abiFile:  "testdata/abi/refs.json",
			typeName: "GoodNamingForMsg",
			cellHex:  "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
		},
		{
			name:     "refs from struct",
			abiFile:  "testdata/abi/refs.json",
			typeName: "ManyRefsMsg",
			cellHex:  "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
		},
		{
			name:     "generics from struct",
			abiFile:  "testdata/abi/generics.json",
			typeName: "ManyRefsMsg<uint16>",
			cellHex:  "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
		},
		{
			name:     "generics from alias",
			abiFile:  "testdata/abi/generics.json",
			typeName: "GoodNamingForMsg<uint16>",
			cellHex:  "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
		},
		{
			name:     "default values",
			abiFile:  "testdata/abi/default_values.json",
			typeName: "DefaultTest",
			cellHex:  "b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde",
		},
		{
			name:     "numbers",
			abiFile:  "testdata/abi/numbers.json",
			typeName: "Numbers",
			cellHex:  "b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014",
		},
		{
			name:     "random fields",
			abiFile:  "testdata/abi/random_fields.json",
			typeName: "RandomFields",
			cellHex:  "b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			abi := loadTestABI(t, tt.abiFile)
			tyIdx := testTyIdx(t, abi, tt.typeName)

			currCell, err := boc.DeserializeBocHex(tt.cellHex)
			if err != nil {
				t.Fatal(err)
			}
			decoder := NewDecoder(abi)
			v, err := decoder.UnmarshalTyIdx(currCell[0], tyIdx)
			if err != nil {
				t.Fatal(err)
			}

			encoder := NewEncoder(abi)
			newCell, err := encoder.MarshalTyIdx(v, tyIdx)
			if err != nil {
				t.Fatal(err)
			}
			assertSameCellHash(t, currCell[0], newCell)
		})
	}
}

func TestRuntime_MarshalAliasWithCustomUnpack(t *testing.T) {
	inputFilename := "testdata/abi/custom_pack_unpack.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi parser.ContractABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	tyIdx := testTyIdx(t, abi, "MyAlias")

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100470000890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043b9aca00886e91196")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder(abi)
	decoder.WithCustomUnpackResolver(func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
		err := cell.Skip(512)
		if err != nil {
			return fmt.Errorf("failed to 512 bits from alias")
		}
		val, err := decoder.UnmarshalTyIdx(cell, testTyIdx(t, abi, "My"))
		if err != nil {
			return fmt.Errorf("failed to unmarshal alias' coins")
		}
		*value = AliasValue(*val)
		return nil
	})
	v, err := decoder.UnmarshalTyIdx(currCell[0], tyIdx)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder(abi)
	encoder.WithCustomPackResolver(func(ref parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
		err := cell.WriteUint(0, 256)
		if err != nil {
			return fmt.Errorf("failed to write 256 bits to alias")
		}
		err = cell.WriteUint(0, 256)
		if err != nil {
			return fmt.Errorf("failed to write 256 bits to alias again")
		}
		val := Value(*value)
		cl, err := encoder.MarshalTyIdx(&val, testTyIdx(t, abi, "My"))
		if err != nil {
			return fmt.Errorf("failed to marshal alias' coins")
		}
		err = cell.WriteBitString(cl.ReadRemainingBits())
		if err != nil {
			return fmt.Errorf("failed to write marshalled bits to alias")
		}
		for _, clRef := range cl.Refs() {
			err = cell.AddRef(clRef)
			if err != nil {
				return fmt.Errorf("failed to add ref to alias")
			}
		}
		return nil
	})
	newCell, err := encoder.MarshalTyIdx(v, tyIdx)
	if err != nil {
		t.Error(err)
	}

	oldHs, err := currCell[0].HashString()
	if err != nil {
		t.Fatal(err)
	}
	newHs, err := newCell.HashString()
	if err != nil {
		t.Fatal(err)
	}
	if oldHs != newHs {
		t.Errorf("input and output cells are different")
	}

}

func compareExpectedJson(inputFilename string, v Value) error {
	val := struct {
		Value Value `json:"value"`
	}{
		Value: v,
	}
	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		return err
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		return err
	}
	if !bytes.Equal(actual, expected) {
		return fmt.Errorf("%s got different results", pathPrefix)
	}
	return nil
}

func compareExpectedInlineJson(expected string, v Value) error {
	val := struct {
		Value Value `json:"value"`
	}{
		Value: v,
	}
	actual, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if string(actual) != normalizeJson(expected) {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}
	return nil
}

func normalizeJson(jsonstr string) string {
	var v any
	err := json.Unmarshal([]byte(jsonstr), &v)
	if err != nil {
		panic(err)
	}
	norm, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(norm)
}

func ptr[T any](t T) *T {
	return &t
}
