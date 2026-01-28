package tolk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
)

func TestRuntimeUnmarshal(t *testing.T) {
	type cases struct {
		name     string
		filename string
		cell     string
		t        Ty
		check    func(Value)
	}
	for _, c := range []cases{
		{
			name:     "unmarshal small int",
			filename: "simple",
			cell:     "b5ee9c72410101010005000006ff76c41616db06",
			t: Ty{
				SumType: "IntN",
				IntN: &IntN{
					N: 24,
				},
			},
			check: func(v Value) {
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
			name:     "unmarshal big int",
			filename: "simple",
			cell:     "b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c",
			t: Ty{
				SumType: "IntN",
				IntN: &IntN{
					N: 183,
				},
			},
			check: func(v Value) {
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
			name:     "unmarshal small uint",
			filename: "simple",
			cell:     "b5ee9c7241010101000900000d00000000001d34e435eafd",
			t: Ty{
				SumType: "UintN",
				UintN: &UintN{
					N: 53,
				},
			},
			check: func(v Value) {
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
			name:     "unmarshal big uint",
			filename: "simple",
			cell:     "b5ee9c7241010101002300004100000000000000000000000000000000000000000000000000009fc4212a38ba40b11cce12",
			t: Ty{
				SumType: "UintN",
				UintN: &UintN{
					N: 257,
				},
			},
			check: func(v Value) {
				val, ok := v.GetBigInt()
				if !ok {
					t.Errorf("v.GetBigInt() not successeded")
				}
				if val.Cmp(big.NewInt(351329437643124)) != 0 {
					t.Errorf("val != 351329437643124, got %v", val)
				}
			},
		},
		{
			name:     "unmarshal varint",
			filename: "simple",
			cell:     "b5ee9c7241010101000600000730c98588449b6923",
			t: Ty{
				SumType: "VarIntN",
				VarIntN: &VarIntN{
					N: 16,
				},
			},
			check: func(v Value) {
				val, ok := v.GetBigInt()
				if !ok {
					t.Errorf("v.GetBigInt() not successeded")
				}
				if val.Cmp(big.NewInt(825432)) != 0 {
					t.Errorf("val != 825432, got %v", val)
				}
			},
		},
		{
			name:     "unmarshal varuint",
			filename: "simple",
			cell:     "b5ee9c7241010101000800000b28119ab36b44d3a86c0f",
			t: Ty{
				SumType: "VarUintN",
				VarUintN: &VarUintN{
					N: 32,
				},
			},
			check: func(v Value) {
				val, ok := v.GetBigInt()
				if !ok {
					t.Errorf("v.GetBigInt() not successeded")
				}
				if val.Cmp(big.NewInt(9451236712)) != 0 {
					t.Errorf("val != 9451236712, got %v", val)
				}
			},
		},
		{
			name:     "unmarshal bits",
			filename: "simple",
			cell:     "b5ee9c7241010101000500000631323318854035",
			t: Ty{
				SumType: "BitsN",
				BitsN: &BitsN{
					N: 24,
				},
			},
			check: func(v Value) {
				val, ok := v.GetBits()
				if !ok {
					t.Errorf("v.GetBits() not successeded")
				}
				if bytes.Equal(val.Buffer(), []byte{55, 56, 57}) {
					t.Errorf("val != {55, 56, 57}, got %v", val)
				}
			},
		},
		{
			name:     "unmarshal coins",
			filename: "simple",
			cell:     "b5ee9c72410101010007000009436ec6e0189ebbd7f4",
			t: Ty{
				SumType: "Coins",
				Coins:   &Coins{},
			},
			check: func(v Value) {
				val, ok := v.GetBigInt()
				if !ok {
					t.Errorf("v.GetBigInt() not successeded")
				}
				if val.Cmp(big.NewInt(921464321)) != 0 {
					t.Errorf("val != 921464321, got %v", val)
				}
			},
		},
		{
			name:     "unmarshal bool",
			filename: "simple",
			cell:     "b5ee9c7241010101000300000140f6d24034",
			t: Ty{
				SumType: "Bool",
				Bool:    &Bool{},
			},
			check: func(v Value) {
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
			name:     "unmarshal cell",
			filename: "simple",
			cell:     "b5ee9c724101020100090001000100080000007ba52a3292",
			t: Ty{
				SumType: "Cell",
				Cell:    &Cell{},
			},
			check: func(v Value) {
				val, ok := v.GetCell()
				if !ok {
					t.Errorf("v.GetCell() not successeded")
				}
				hs, err := val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "644e68a539c5107401d194bc82169cbf0ad1635796891551e0750705ab2d74ae" {
					t.Errorf("val.Hash() != 644e68a539c5107401d194bc82169cbf0ad1635796891551e0750705ab2d74ae, got %v", hs)
				}
			},
		},
		{
			name:     "unmarshal remaining",
			filename: "simple",
			cell:     "b5ee9c7241010101000900000dc0800000000ab8d04726e4",
			t: Ty{
				SumType:   "Remaining",
				Remaining: &Remaining{},
			},
			check: func(v Value) {
				val, ok := v.GetCell()
				if !ok {
					t.Errorf("v.GetCell() not successeded")
				}
				hs, err := val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "f1c4e07fbd1786411c2caa9ac9f5d7240aa2007a2a1d5e5ac44f8a168cd4e36b" {
					t.Errorf("val.Hash() != f1c4e07fbd1786411c2caa9ac9f5d7240aa2007a2a1d5e5ac44f8a168cd4e36b, got %v", hs)
				}
			},
		},
		{
			name:     "unmarshal address",
			filename: "simple",
			cell:     "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			t: Ty{
				SumType: "Address",
				Address: &Address{},
			},
			check: func(v Value) {
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
			name:     "unmarshal not exists optional address",
			filename: "simple",
			cell:     "b5ee9c724101010100030000012094418655",
			t: Ty{
				SumType:    "AddressOpt",
				AddressOpt: &AddressOpt{},
			},
			check: func(v Value) {
				val, ok := v.GetOptionalAddress()
				if !ok {
					t.Errorf("v.GetOptionalAddress() not successeded")
				}

				if val.SumType != "NoneAddress" {
					t.Errorf("val.GetAddress() != none address")
				}
			},
		},
		{
			name:     "unmarshal exists optional address",
			filename: "simple",
			cell:     "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			t: Ty{
				SumType:    "AddressOpt",
				AddressOpt: &AddressOpt{},
			},
			check: func(v Value) {
				val, ok := v.GetOptionalAddress()
				if !ok {
					t.Errorf("v.GetOptionalAddress() not successeded")
				}

				if val.SumType == "InternalAddress" && val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
				}
			},
		},
		{
			name:     "unmarshal external address",
			filename: "simple",
			cell:     "b5ee9c7241010101000600000742082850fcbd94fd",
			t: Ty{
				SumType:    "AddressExt",
				AddressExt: &AddressExt{},
			},
			check: func(v Value) {
				val, ok := v.GetExternalAddress()
				if !ok {
					t.Errorf("v.GetExternalAddress() not successeded")
				}
				addressPart := boc.NewBitString(16)
				err := addressPart.WriteBytes([]byte{97, 98})
				if err != nil {
					t.Fatal(err)
				}
				if val.Len != 8 && bytes.Equal(val.Address.Buffer(), []byte{97, 98}) {
					t.Errorf("val.GetExternalAddress() != {97, 98}, got %v", val.Address.Buffer())
				}
			},
		},
		{
			name:     "unmarshal none any address",
			filename: "simple",
			cell:     "b5ee9c724101010100030000012094418655",
			t: Ty{
				SumType:    "AddressAny",
				AddressAny: &AddressAny{},
			},
			check: func(v Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != "NoneAddress" {
					t.Errorf("val.GetAddress() != none address")
				}
			},
		},
		{
			name:     "unmarshal internal any address",
			filename: "simple",
			cell:     "b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6",
			t: Ty{
				SumType:    "AddressAny",
				AddressAny: &AddressAny{},
			},
			check: func(v Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType == "InternalAddress" && val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
				}
			},
		},
		{
			name:     "unmarshal external any address",
			filename: "simple",
			cell:     "b5ee9c7241010101000600000742082850fcbd94fd",
			t: Ty{
				SumType:    "AddressAny",
				AddressAny: &AddressAny{},
			},
			check: func(v Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				addressPart := boc.NewBitString(16)
				err := addressPart.WriteBytes([]byte{97, 98})
				if err != nil {
					t.Fatal(err)
				}
				if val.SumType == "ExternalAddress" && val.ExternalAddress.Len != 8 && bytes.Equal(val.ExternalAddress.Address.Buffer(), []byte{97, 98}) {
					t.Errorf("val.GetExternalAddress() != {97, 98}, got %v", val.ExternalAddress.Address.Buffer())
				}
			},
		},
		{
			name:     "unmarshal var any address",
			filename: "simple",
			cell:     "b5ee9c7241010101000900000dc0800000000ab8d04726e4",
			t: Ty{
				SumType:    "AddressAny",
				AddressAny: &AddressAny{},
			},
			check: func(v Value) {
				val, ok := v.GetAnyAddress()
				if !ok {
					t.Errorf("v.GetAnyAddress() not successeded")
				}
				if val.SumType != "VarAddress" {
					t.Errorf("val.GetAddress() != VarAddress")
				}
				if val.VarAddress.Len != 8 {
					t.Errorf("val.VarAddress.Len != 8, got %v", val.VarAddress.Len)
				}
				if val.VarAddress.Workchain != 0 {
					t.Errorf("val.VarAddress.Workchain != 0, got %v", val.VarAddress.Workchain)
				}
				if bytes.Equal(val.VarAddress.Address.Buffer(), []byte{97, 98}) {
					t.Errorf("val.GetExternalAddress() != {97, 98}, got %v", val.ExternalAddress.Address.Buffer())
				}
			},
		},
		{
			name:     "unmarshal not exists nullable",
			filename: "simple",
			cell:     "b5ee9c7241010101000300000140f6d24034",
			t: Ty{
				SumType: "Nullable",
				Nullable: &Nullable{
					Inner: Ty{
						SumType:   "Remaining",
						Remaining: &Remaining{},
					},
				},
			},
			check: func(v Value) {
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
			name:     "unmarshal exists nullable",
			filename: "simple",
			cell:     "b5ee9c7241010201000b000101c001000900000c0ae007880db9",
			t: Ty{
				SumType: "Nullable",
				Nullable: &Nullable{
					Inner: Ty{
						SumType: "Cell",
						Cell:    &Cell{},
					},
				},
			},
			check: func(v Value) {
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
				hs, err := innerVal.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "df05386a55563049a4834a4cc1ec0dc22f3dcb63c04f7258ae475c5d28981773" {
					t.Errorf("v.GetOptionalValue().GetCell() != df05386a55563049a4834a4cc1ec0dc22f3dcb63c04f7258ae475c5d28981773, got %v", hs)
				}
			},
		},
		{
			name:     "unmarshal ref",
			filename: "simple",
			cell:     "b5ee9c7241010201000e000100010011000000000009689e40e150b4c5",
			t: Ty{
				SumType: "CellOf",
				CellOf: &CellOf{
					Inner: Ty{
						SumType: "IntN",
						IntN: &IntN{
							N: 65,
						},
					},
				},
			},
			check: func(v Value) {
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
			name:     "unmarshal empty tensor",
			filename: "simple",
			cell:     "b5ee9c724101010100020000004cacb9cd",
			t: Ty{
				SumType: "Tensor",
				Tensor:  &Tensor{},
			},
			check: func(v Value) {
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
			name:     "unmarshal tensor",
			filename: "simple",
			cell:     "b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117",
			t: Ty{
				SumType: "Tensor",
				Tensor: &Tensor{
					Items: []Ty{
						{
							SumType: "UintN",
							UintN: &UintN{
								N: 123,
							},
						},
						{
							SumType: "Bool",
							Bool:    &Bool{},
						},
						{
							SumType: "Coins",
							Coins:   &Coins{},
						},
						{
							SumType: "Tensor",
							Tensor: &Tensor{
								Items: []Ty{
									{
										SumType: "IntN",
										IntN: &IntN{
											N: 23,
										},
									},
									{
										SumType: "Nullable",
										Nullable: &Nullable{
											Inner: Ty{
												SumType: "IntN",
												IntN: &IntN{
													N: 2,
												},
											},
										},
									},
								},
							},
						},
						{
							SumType: "VarIntN",
							VarIntN: &VarIntN{
								N: 32,
							},
						},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetTensor()
				if !ok {
					t.Errorf("v.GetTensor() not successeded")
				}

				val0, ok := val[0].GetBigInt()
				if !ok {
					t.Errorf("val[0].GetBigInt() not successeded")
				}
				if val0.Cmp(big.NewInt(4325)) != 0 {
					t.Errorf("val[0].GetBigInt() != 4325, got %v", val0.String())
				}

				val1, ok := val[1].GetBool()
				if !ok {
					t.Errorf("val[1].GetBigInt() not successeded")
				}
				if !val1 {
					t.Error("val[1].GetBool() is false")
				}

				val2, ok := val[2].GetBigInt()
				if !ok {
					t.Errorf("val[2].GetBigInt() not successeded")
				}
				if val2.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Errorf("val[2].GetBigInt() != 1000000000, got %v", val2.String())
				}

				val3, ok := val[3].GetTensor()
				if !ok {
					t.Errorf("val[3].GetTensor() not successeded")
				}

				val30, ok := val3[0].GetSmallInt()
				if !ok {
					t.Errorf("val[3][0].GetSmallInt() not successeded")
				}
				if val30 != -342 {
					t.Errorf("val[3][0].GetSmallInt() != -342, got %v", val30)
				}

				optVal31, ok := val3[1].GetOptionalValue()
				if !ok {
					t.Errorf("val[3][1].GetOptionalValue() not successeded")
				}
				if !optVal31.IsExists {
					t.Errorf("val[3][1].GetOptionalValue() != exists")
				}
				val31, ok := optVal31.Val.GetSmallInt()
				if !ok {
					t.Errorf("val[3][1].GetOptionalValue().GetSmallInt() not successeded")
				}
				if val31 != 0 {
					t.Errorf("val[3][1].GetOptionalValue().GetSmallInt() != 0, got %v", val31)
				}

				val4, ok := val[4].GetBigInt()
				if !ok {
					t.Errorf("val[4].GetBigInt() not successeded")
				}
				if val4.Cmp(big.NewInt(-9_304_000_000)) != 0 {
					t.Errorf("val[4].GetBigInt() != -9304000000, got %v", val4.String())
				}
			},
		},
		//{
		//	name:     "unmarshal empty tuple",
		//	filename: "simple",
		//	t: Ty{
		//		SumType:   "TupleWith",
		//		TupleWith: &TupleWith{},
		//	},
		//	check: func(v Value) {
		//		val, ok := v.GetTupleValues()
		//		if !ok {
		//			t.Errorf("v.GetTupleValues() not successeded")
		//		}
		//
		//		if len(val) != 0 {
		//			t.Errorf("v.GetTupleValues() != empty")
		//		}
		//	},
		//},
		//{
		//	name:     "unmarshal tuple",
		//	filename: "simple",
		//	t: Ty{
		//		SumType: "TupleWith",
		//		TupleWith: &TupleWith{
		//			Items: []Ty{
		//				{
		//					SumType: "Nullable",
		//					Nullable: &Nullable{
		//						Inner: Ty{
		//							SumType: "CellOf",
		//							CellOf: &CellOf{
		//								Inner: Ty{
		//									SumType: "IntN",
		//									IntN: &IntN{
		//										N: 1,
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//				{
		//					SumType: "TupleWith",
		//					TupleWith: &TupleWith{
		//						Items: []Ty{
		//							{
		//								SumType: "UintN",
		//								UintN: &UintN{
		//									N: 1,
		//								},
		//							},
		//							{
		//								SumType: "Tensor",
		//								Tensor: &Tensor{
		//									Items: []Ty{
		//										{
		//											SumType: "CellOf",
		//											CellOf: &CellOf{
		//												Inner: Ty{
		//													SumType: "Cell",
		//													Cell:    &Cell{},
		//												},
		//											},
		//										},
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//		},
		//	},
		//	check: func(v Value) {
		//		val, ok := v.GetTupleValues()
		//		if !ok {
		//			t.Errorf("v.GetTupleValues() not successeded")
		//		}
		//
		//		val0, ok := val[0].GetOptionalValue()
		//		if !ok {
		//			t.Errorf("val[0].GetOptionalValue() not successeded")
		//		}
		//		if !val0.IsExists {
		//			t.Errorf("val[0].GetOptionalValue() != exists")
		//		}
		//		val0Ref, ok := val0.Val.GetRefValue()
		//		if !ok {
		//			t.Errorf("val[0].GetOptionalValue().GetRefValue() not successeded")
		//		}
		//		val0RefVal, ok := val0Ref.GetSmallInt()
		//		if !ok {
		//			t.Errorf("val[0].GetOptionalValue().GetRefValue().GetSmallInt() not successeded")
		//		}
		//		if val0RefVal != -1 {
		//			t.Errorf("val[0].GetOptionalValue().GetRefValue().GetSmallInt() != -1, got %v", val0RefVal)
		//		}
		//
		//		val1, ok := val[1].GetTupleValues()
		//		if !ok {
		//			t.Errorf("val[1].GetTupleValues() not successeded")
		//		}
		//
		//		val10, ok := val1[0].GetSmallUInt()
		//		if !ok {
		//			t.Errorf("val[1][0].GetSmallUInt() not successeded")
		//		}
		//		if val10 != 1 {
		//			t.Errorf("val[1][0].GetSmallUInt() != 1, got %v", val0RefVal)
		//		}
		//
		//		val11, ok := val1[1].GetTensor()
		//		if !ok {
		//			t.Errorf("val[1][1].GetTensor() not successeded")
		//		}
		//
		//		val110Ref, ok := val11[0].GetRefValue()
		//		if !ok {
		//			t.Errorf("val[1][1][0].GetRefValue() not successeded")
		//		}
		//		val110, ok := val110Ref.GetCell()
		//		if !ok {
		//			t.Errorf("val[1][1][0].GetRefValue().GetCell() not successeded")
		//		}
		//		hs, err := val110.HashString()
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		if hs != "123" {
		//			t.Errorf("val[1][1][0].GetHash() != 123, got %v", hs)
		//		}
		//	},
		//},
		{
			name:     "unmarshal int-key map",
			filename: "simple",
			cell:     "b5ee9c7241010201000c000101c001000ba00000007bc09a662c32",
			t: Ty{
				SumType: "Map",
				Map: &Map{
					K: Ty{
						SumType: "IntN",
						IntN: &IntN{
							N: 32,
						},
					},
					V: Ty{
						SumType: "Bool",
						Bool:    &Bool{},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetMap()
				if !ok {
					t.Errorf("v.GetMap() not successeded")
				}
				val123, ok := val.GetBySmallInt(Int64(123))
				if !ok {
					t.Errorf("val[123] not found")
				}
				val123Val, ok := val123.GetBool()
				if !ok {
					t.Errorf("val[123].GetBool() not successeded")
				}
				if !val123Val {
					t.Errorf("val[123] is false")
				}

				_, ok = val.GetBySmallInt(Int64(0))
				if ok {
					t.Errorf("val[0] was found")
				}
			},
		},
		{
			name:     "unmarshal uint-key map",
			filename: "simple",
			cell:     "b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8",
			t: Ty{
				SumType: "Map",
				Map: &Map{
					K: Ty{
						SumType: "UintN",
						UintN: &UintN{
							N: 16,
						},
					},
					V: Ty{
						SumType: "Address",
						Address: &Address{},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetMap()
				if !ok {
					t.Errorf("v.GetMap() not successeded")
				}
				val23, ok := val.GetBySmallUInt(UInt64(23))
				if !ok {
					t.Errorf("val[23] not found")
				}
				val23Val, ok := val23.GetAddress()
				if !ok {
					t.Errorf("val[23].GetAddress() not successeded")
				}
				if val23Val.ToRaw() != "-1:5555555555555555555555555555555555555555555555555555555555555555" {
					t.Errorf("val[23] != -1:5555555555555555555555555555555555555555555555555555555555555555, got %v", val23Val.ToRaw())
				}

				val14, ok := val.GetBySmallUInt(UInt64(14))
				if !ok {
					t.Errorf("val[14] not found")
				}
				val14Val, ok := val14.GetAddress()
				if !ok {
					t.Errorf("val[14].GetAddress() not successeded")
				}
				if val14Val.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Errorf("val[14] != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", val23Val.ToRaw())
				}

				_, ok = val.GetBySmallInt(Int64(0))
				if ok {
					t.Errorf("val[0] was found")
				}
			},
		},
		{
			name:     "unmarshal bigint-key map",
			filename: "simple",
			cell:     "b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9",
			t: Ty{
				SumType: "Map",
				Map: &Map{
					K: Ty{
						SumType: "UintN",
						UintN: &UintN{
							N: 78,
						},
					},
					V: Ty{
						SumType: "Cell",
						Cell:    &Cell{},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetMap()
				if !ok {
					t.Errorf("v.GetMap() not successeded")
				}
				val1, ok := val.GetByBigInt(BigInt(*big.NewInt(2337412)))
				if !ok {
					t.Errorf("val[2337412] not found")
				}
				val1Val, ok := val1.GetCell()
				if !ok {
					t.Errorf("val[2337412].GetCell() not successeded")
				}
				hs1, err := val1Val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs1 != "8be375797c46a090b06973ee57e96b1d1ae127609c400ceba7194e77e41c5150" {
					t.Errorf("val[2337412].GetCell().GetHashString() != 8be375797c46a090b06973ee57e96b1d1ae127609c400ceba7194e77e41c5150, got %v", hs1)
				}

				_, ok = val.GetByBigInt(BigInt(*big.NewInt(34)))
				if ok {
					t.Errorf("val[34] was found")
				}
			},
		},
		{
			name:     "unmarshal bits-key map",
			filename: "simple",
			cell:     "b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86",
			t: Ty{
				SumType: "Map",
				Map: &Map{
					K: Ty{
						SumType: "BitsN",
						BitsN: &BitsN{
							N: 16,
						},
					},
					V: Ty{
						SumType: "Map",
						Map: &Map{
							K: Ty{
								SumType: "IntN",
								IntN: &IntN{
									N: 64,
								},
							},
							V: Ty{
								SumType: "Tensor",
								Tensor: &Tensor{
									Items: []Ty{
										{
											SumType: "Address",
											Address: &Address{},
										},
										{
											SumType: "Coins",
											Coins:   &Coins{},
										},
									},
								},
							},
						},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetMap()
				if !ok {
					t.Errorf("v.GetMap() not successeded")
				}
				key1 := boc.NewBitString(16)
				err := key1.WriteBytes([]byte{65, 66})
				if err != nil {
					t.Fatal(err)
				}
				val1, ok := val.GetByBits(Bits(key1))
				if !ok {
					t.Errorf("val[{65, 66}] not found")
				}

				mp, ok := val1.GetMap()
				if !ok {
					t.Errorf("val[{65, 66}].GetMap() not successeded")
				}
				val1_124, ok := mp.GetBySmallInt(124)
				if !ok {
					t.Errorf("val[{65, 66}][124] not found")
				}
				val1_124Val, ok := val1_124.GetTensor()
				if !ok {
					t.Errorf("val[{65, 66}][124].GetTensor() not successeded")
				}
				val1_124Val0, ok := val1_124Val[0].GetAddress()
				if !ok {
					t.Errorf("val[{65, 66}][124][0].GetAddress() not successeded")
				}
				if val1_124Val0.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Errorf("val[{65, 66}][124][0].GetAddress() != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", val1_124Val0.ToRaw())
				}

				val1_124Val1, ok := val1_124Val[1].GetBigInt()
				if !ok {
					t.Errorf("val[{97, 98}][124][1].GetBigInt() not successeded")
				}
				if val1_124Val1.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Errorf("val[{97, 98}][124][1].GetBigInt() != 1_000_000_000, got %v", val1_124Val1.String())
				}

				key2 := boc.NewBitString(16)
				err = key2.WriteBytes([]byte{98, 99})
				if err != nil {
					t.Fatal(err)
				}
				_, ok = val.GetByBits(Bits(key2))
				if ok {
					t.Errorf("val[{98, 99}] was found")
				}
			},
		},
		{
			name:     "unmarshal address-key map",
			filename: "simple",
			cell:     "b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e",
			t: Ty{
				SumType: "Map",
				Map: &Map{
					K: Ty{
						SumType: "Address",
						Address: &Address{},
					},
					V: Ty{
						SumType: "Coins",
						Coins:   &Coins{},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetMap()
				if !ok {
					t.Errorf("v.GetMap() not successeded")
				}
				// todo: create converter
				addr := tongo.MustParseAddress("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")
				val1, ok := val.GetByInternalAddress(InternalAddress{
					Workchain: int8(addr.ID.Workchain),
					Address:   addr.ID.Address,
				})
				if !ok {
					t.Errorf("val[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"] not found")
				}
				val1Val, ok := val1.GetBigInt()
				if !ok {
					t.Errorf("val[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"].GetCoins() not successeded")
				}
				if val1Val.Cmp(big.NewInt(10_000_000_000)) != 0 {
					t.Errorf("val[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"].GetCoins() != 10_000_000_000, got %v", val1Val)
				}

				addr = tongo.MustParseAddress("UQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqEBI")
				_, ok = val.GetByInternalAddress(InternalAddress{
					Workchain: int8(addr.ID.Workchain),
					Address:   addr.ID.Address,
				})
				if ok {
					t.Errorf("val[\"UQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqEBI\"] was found")
				}
			},
		},
		{
			name:     "unmarshal union with dec prefix",
			filename: "simple",
			cell:     "b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029",
			t: Ty{
				SumType: "Union",
				Union: &Union{
					Variants: []UnionVariant{
						{
							PrefixStr:        "0",
							PrefixLen:        1,
							PrefixEatInPlace: true,
							VariantTy: Ty{
								SumType: "IntN",
								IntN: &IntN{
									N: 16,
								},
							},
						},
						{
							PrefixStr:        "1",
							PrefixLen:        1,
							PrefixEatInPlace: true,
							VariantTy: Ty{
								SumType: "IntN",
								IntN: &IntN{
									N: 128,
								},
							},
						},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetUnion()
				if !ok {
					t.Errorf("v.GetUnion() not successeded")
				}
				if val.Prefix.Len != 1 {
					t.Errorf("val.Prefix.Len != 1")
				}
				if val.Prefix.Prefix != 1 {
					t.Errorf("val.Prefix != 1, got %v", val.Prefix.Prefix)
				}

				unionVal, ok := val.Val.GetBigInt()
				if !ok {
					t.Errorf("val.Val.GetBigInt() not successeded")
				}
				if unionVal.Cmp(big.NewInt(124432123)) != 0 {
					t.Errorf("val.Val.GetBigInt() != 124432123, got %v", unionVal.String())
				}
			},
		},
		{
			name:     "unmarshal union with bin prefix",
			filename: "bin_union",
			cell:     "b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650",
			t: Ty{
				SumType: "Union",
				Union: &Union{
					Variants: []UnionVariant{
						{
							PrefixStr: "0b001",
							PrefixLen: 3,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "AddressWithPrefix",
								},
							},
						},
						{
							PrefixStr: "0b011",
							PrefixLen: 3,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "MapWithPrefix",
								},
							},
						},
						{
							PrefixStr: "0b111",
							PrefixLen: 3,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "CellWithPrefix",
								},
							},
						},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetUnion()
				if !ok {
					t.Errorf("v.GetUnion() not successeded")
				}
				if val.Prefix.Len != 3 {
					t.Errorf("val.Prefix.Len != 3, got %v", val.Prefix.Len)
				}
				if val.Prefix.Prefix != 3 {
					t.Errorf("val.Prefix.Prefix != 3, got %v", val.Prefix.Prefix)
				}

				mapStruct, ok := val.Val.GetStruct()
				if !ok {
					t.Errorf("val.GetStruct() not successeded")
				}
				mapStructVal, ok := mapStruct.GetField("v")
				if !ok {
					t.Errorf("val[v] not successeded")
				}
				unionVal, ok := mapStructVal.GetMap()
				if !ok {
					t.Errorf("val[v].GetMap() not successeded")
				}
				addr := tongo.MustParseAddress("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")
				mapVal, ok := unionVal.GetByInternalAddress(InternalAddress{
					Workchain: int8(addr.ID.Workchain),
					Address:   addr.ID.Address,
				})
				if !ok {
					t.Errorf("val.GetMap()[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"] not successeded")
				}
				mapCoins, ok := mapVal.GetBigInt()
				if !ok {
					t.Errorf("val.GetMap()[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"].GetBigInt() not successeded")
				}
				if mapCoins.Cmp(big.NewInt(43213412)) != 0 {
					t.Errorf("val.GetMap()[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"] != 43213412, got %v", mapCoins.String())
				}
			},
		},
		{
			name:     "unmarshal union with hex prefix",
			filename: "hex_union",
			cell:     "b5ee9c7241010101000b000011deadbeef00000000c0d75977b9",
			t: Ty{
				SumType: "Union",
				Union: &Union{
					Variants: []UnionVariant{
						{
							PrefixStr: "0x12345678",
							PrefixLen: 32,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "UInt66WithPrefix",
								},
							},
						},
						{
							PrefixStr: "0xdeadbeef",
							PrefixLen: 32,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "UInt33WithPrefix",
								},
							},
						},
						{
							PrefixStr: "0x89abcdef",
							PrefixLen: 32,
							VariantTy: Ty{
								SumType: "StructRef",
								StructRef: &StructRef{
									StructName: "UInt4WithPrefix",
								},
							},
						},
					},
				},
			},
			check: func(v Value) {
				val, ok := v.GetUnion()
				if !ok {
					t.Errorf("v.GetUnion() not successeded")
				}
				if val.Prefix.Len != 32 {
					t.Errorf("val.Prefix.Len != 32, got %v", val.Prefix.Len)
				}
				if val.Prefix.Prefix != 0xdeadbeef {
					t.Errorf("val.Prefix.Prefix != 0xdeadbeef, got %x", val.Prefix.Prefix)
				}

				structVal, ok := val.Val.GetStruct()
				if !ok {
					t.Errorf("val.Val.GetStruct() not successeded")
				}
				structV, ok := structVal.GetField("v")
				if !ok {
					t.Errorf("val.Val[v] not successeded")
				}
				unionVal, ok := structV.GetSmallUInt()
				if !ok {
					t.Errorf("val.GetSmallUInt() not successeded")
				}
				if unionVal != 1 {
					t.Errorf("val.GetSmallUInt() != 1, got %v", unionVal)
				}
			},
		},
		{
			name:     "unmarshal a lot refs from alias",
			filename: "refs",
			cell:     "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
			t: Ty{
				SumType: "AliasRef",
				AliasRef: &AliasRef{
					AliasName: "GoodNamingForMsg",
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}
				pref, ok := currStruct.GetStructPrefix()
				if !ok {
					t.Fatalf("currStruct.Prefix not found")
				}
				if pref.Len != 32 {
					t.Errorf("pref.Len != 32, got %v", pref.Len)
				}
				if pref.Prefix != 0xdeadbeef {
					t.Errorf("val.Prefix.Prefix != 0xdeadbeef, got %x", pref.Prefix)
				}

				user1, ok := currStruct.GetField("user1")
				if !ok {
					t.Fatalf("currStruct[user1] not found")
				}
				user1Val, ok := user1.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user1].GetStruct() not successeded")
				}

				user1Addr, ok := user1Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user1][addr] not found")
				}
				user1AddrVal, ok := user1Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user1][addr].GetAddress() not successeded")
				}
				if user1AddrVal.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("user1AddrVal.ToRaw() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", user1AddrVal.ToRaw())
				}

				user1Balance, ok := user1Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user1][balance] not found")
				}
				user1BalanceVal, ok := user1Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user1][balance].GetBigInt() not successeded")
				}
				if user1BalanceVal.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Errorf("currStruct[user1][balance].ToRaw() != 1000000000, got %v", user1BalanceVal.String())
				}

				user2, ok := currStruct.GetField("user2")
				if !ok {
					t.Fatalf("currStruct[user2] not found")
				}
				user2Opt, ok := user2.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[user2].GetOptionalValue() not successeded")
				}
				if !user2Opt.IsExists {
					t.Errorf("currStruct[user2] is not exists")
				}
				user2Ref, ok := user2Opt.Val.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[user2].GetRefValue() not successeded")
				}
				user2Val, ok := user2Ref.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user2].GetStruct() not successeded")
				}

				user2Addr, ok := user2Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user2][addr] not found")
				}
				user2AddrVal, ok := user2Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user2][addr].GetAddress() not successeded")
				}
				if user2AddrVal.ToRaw() != "0:086fa2a675f74347b08dd4606a549b8fdb98829cb282bc1949d3b12fbaed9dcc" {
					t.Errorf("user1AddrVal.ToRaw() != 0:086fa2a675f74347b08dd4606a549b8fdb98829cb282bc1949d3b12fbaed9dcc, got %v", user2AddrVal.ToRaw())
				}

				user2Balance, ok := user2Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user2][balance] not found")
				}
				user2BalanceVal, ok := user2Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user2][balance].GetBigInt() not successeded")
				}
				if user2BalanceVal.Cmp(big.NewInt(100_000_000)) != 0 {
					t.Errorf("currStruct[user2][balance].ToRaw() != 100000000, got %v", user2BalanceVal.String())
				}

				user3, ok := currStruct.GetField("user3")
				if !ok {
					t.Fatalf("currStruct[user3] not found")
				}
				user3Val, ok := user3.GetCell()
				if !ok {
					t.Fatalf("currStruct[user3].GetCell() not successeded")
				}
				hs, err := user3Val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "47f4b117a301111ec48d763a3cd668a246c174efd2df9ba8bd1db406f017453a" {
					t.Errorf("currStruct[user3][hashString].Hash != 47f4b117a301111ec48d763a3cd668a246c174efd2df9ba8bd1db406f017453a, got %v", hs)
				}

				user4, ok := currStruct.GetField("user4")
				if !ok {
					t.Fatalf("currStruct[user4] not found")
				}
				user4Opt, ok := user4.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[user4].GetOptionalValue() not successeded")
				}
				if user4Opt.IsExists {
					t.Errorf("currStruct[user4] exists")
				}

				user5, ok := currStruct.GetField("user5")
				if !ok {
					t.Fatalf("currStruct[user2] not found")
				}
				user5Ref, ok := user5.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[user5].GetRefValue() not successeded")
				}
				user5Val, ok := user5Ref.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user5].GetStruct() not successeded")
				}

				user5Addr, ok := user5Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user5][addr] not found")
				}
				user5AddrVal, ok := user5Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user5][addr].GetAddress() not successeded")
				}
				if user5AddrVal.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Errorf("user1AddrVal.ToRaw() != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", user5AddrVal.ToRaw())
				}

				user5Balance, ok := user5Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user5][balance] not found")
				}
				user5BalanceVal, ok := user5Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user5][balance].GetBigInt() not successeded")
				}
				if user5BalanceVal.Cmp(big.NewInt(10_000_000_000_000)) != 0 {
					t.Errorf("currStruct[user5][balance].ToRaw() != 10000000000000, got %v", user5BalanceVal.String())
				}

				role, ok := currStruct.GetField("role")
				if !ok {
					t.Fatalf("currStruct[role] not found")
				}
				roleEnum, ok := role.GetEnum()
				if !ok {
					t.Fatalf("currStruct[role].GetEnum() not successeded")
				}
				if roleEnum.Value.Cmp(big.NewInt(1)) != 0 {
					t.Errorf("currStruct[role].GetEnum().Value != 1, got %v", roleEnum.Value.String())
				}
				if roleEnum.Name != "Aboba" {
					t.Errorf("currStruct[role].GetEnum().Name != Aboba, got %v", roleEnum.Name)
				}

				oper1, ok := currStruct.GetField("oper1")
				if !ok {
					t.Fatalf("currStruct[oper1] not found")
				}
				oper1Enum, ok := oper1.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper1].GetEnum() not successeded")
				}
				if oper1Enum.Value.Cmp(big.NewInt(0)) != 0 {
					t.Errorf("currStruct[oper1].GetEnum().Value != 0, got %v", oper1Enum.Value.String())
				}
				if oper1Enum.Name != "Add" {
					t.Errorf("currStruct[oper1].GetEnum().Name != Add, got %v", oper1Enum.Name)
				}

				oper2, ok := currStruct.GetField("oper2")
				if !ok {
					t.Fatalf("currStruct[oper2] not found")
				}
				oper2Enum, ok := oper2.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper2].GetEnum() not successeded")
				}
				if oper2Enum.Value.Cmp(big.NewInt(-10000)) != 0 {
					t.Errorf("currStruct[oper2].GetEnum().Value != -10000, got %v", oper2Enum.Value.String())
				}
				if oper2Enum.Name != "TopUp" {
					t.Errorf("currStruct[oper2].GetEnum().Name != TopUp, got %v", oper2Enum.Name)
				}

				oper3, ok := currStruct.GetField("oper3")
				if !ok {
					t.Fatalf("currStruct[oper3] not found")
				}
				oper3Enum, ok := oper3.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper3].GetEnum() not successeded")
				}
				if oper3Enum.Value.Cmp(big.NewInt(1)) != 0 {
					t.Errorf("currStruct[oper3].GetEnum().Value != 1, got %v", oper3Enum.Value.String())
				}
				if oper3Enum.Name != "Something" {
					t.Errorf("currStruct[oper3].GetEnum().Name != Something, got %v", oper3Enum.Name)
				}
			},
		},
		{
			name:     "unmarshal a lot refs from struct",
			filename: "refs",
			cell:     "b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8",
			t: Ty{
				SumType: "StructRef",
				StructRef: &StructRef{
					StructName: "ManyRefsMsg",
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}
				pref, ok := currStruct.GetStructPrefix()
				if !ok {
					t.Fatalf("currStruct.Prefix not found")
				}
				if pref.Len != 32 {
					t.Errorf("pref.Len != 32, got %v", pref.Len)
				}
				if pref.Prefix != 0xdeadbeef {
					t.Errorf("val.Prefix.Prefix != 0xdeadbeef, got %x", pref.Prefix)
				}

				user1, ok := currStruct.GetField("user1")
				if !ok {
					t.Fatalf("currStruct[user1] not found")
				}
				user1Val, ok := user1.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user1].GetStruct() not successeded")
				}

				user1Addr, ok := user1Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user1][addr] not found")
				}
				user1AddrVal, ok := user1Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user1][addr].GetAddress() not successeded")
				}
				if user1AddrVal.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
					t.Errorf("user1AddrVal.ToRaw() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", user1AddrVal.ToRaw())
				}

				user1Balance, ok := user1Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user1][balance] not found")
				}
				user1BalanceVal, ok := user1Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user1][balance].GetBigInt() not successeded")
				}
				if user1BalanceVal.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Errorf("currStruct[user1][balance].ToRaw() != 1000000000, got %v", user1BalanceVal.String())
				}

				user2, ok := currStruct.GetField("user2")
				if !ok {
					t.Fatalf("currStruct[user2] not found")
				}
				user2Opt, ok := user2.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[user2].GetOptionalValue() not successeded")
				}
				if !user2Opt.IsExists {
					t.Errorf("currStruct[user2] is not exists")
				}
				user2Ref, ok := user2Opt.Val.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[user2].GetRefValue() not successeded")
				}
				user2Val, ok := user2Ref.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user2].GetStruct() not successeded")
				}

				user2Addr, ok := user2Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user2][addr] not found")
				}
				user2AddrVal, ok := user2Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user2][addr].GetAddress() not successeded")
				}
				if user2AddrVal.ToRaw() != "0:086fa2a675f74347b08dd4606a549b8fdb98829cb282bc1949d3b12fbaed9dcc" {
					t.Errorf("user1AddrVal.ToRaw() != 0:086fa2a675f74347b08dd4606a549b8fdb98829cb282bc1949d3b12fbaed9dcc, got %v", user2AddrVal.ToRaw())
				}

				user2Balance, ok := user2Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user2][balance] not found")
				}
				user2BalanceVal, ok := user2Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user2][balance].GetBigInt() not successeded")
				}
				if user2BalanceVal.Cmp(big.NewInt(100_000_000)) != 0 {
					t.Errorf("currStruct[user2][balance].ToRaw() != 100000000, got %v", user2BalanceVal.String())
				}

				user3, ok := currStruct.GetField("user3")
				if !ok {
					t.Fatalf("currStruct[user3] not found")
				}
				user3Val, ok := user3.GetCell()
				if !ok {
					t.Fatalf("currStruct[user3].GetCell() not successeded")
				}
				hs, err := user3Val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "47f4b117a301111ec48d763a3cd668a246c174efd2df9ba8bd1db406f017453a" {
					t.Errorf("currStruct[user3][hashString].Hash != 47f4b117a301111ec48d763a3cd668a246c174efd2df9ba8bd1db406f017453a, got %v", hs)
				}

				user4, ok := currStruct.GetField("user4")
				if !ok {
					t.Fatalf("currStruct[user4] not found")
				}
				user4Opt, ok := user4.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[user4].GetOptionalValue() not successeded")
				}
				if user4Opt.IsExists {
					t.Errorf("currStruct[user4] exists")
				}

				user5, ok := currStruct.GetField("user5")
				if !ok {
					t.Fatalf("currStruct[user2] not found")
				}
				user5Ref, ok := user5.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[user5].GetRefValue() not successeded")
				}
				user5Val, ok := user5Ref.GetStruct()
				if !ok {
					t.Fatalf("currStruct[user5].GetStruct() not successeded")
				}

				user5Addr, ok := user5Val.GetField("addr")
				if !ok {
					t.Fatalf("currStruct[user5][addr] not found")
				}
				user5AddrVal, ok := user5Addr.GetAddress()
				if !ok {
					t.Fatalf("currStruct[user5][addr].GetAddress() not successeded")
				}
				if user5AddrVal.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Errorf("user1AddrVal.ToRaw() != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", user5AddrVal.ToRaw())
				}

				user5Balance, ok := user5Val.GetField("balance")
				if !ok {
					t.Fatalf("currStruct[user5][balance] not found")
				}
				user5BalanceVal, ok := user5Balance.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[user5][balance].GetBigInt() not successeded")
				}
				if user5BalanceVal.Cmp(big.NewInt(10_000_000_000_000)) != 0 {
					t.Errorf("currStruct[user5][balance].ToRaw() != 10000000000000, got %v", user5BalanceVal.String())
				}

				role, ok := currStruct.GetField("role")
				if !ok {
					t.Fatalf("currStruct[role] not found")
				}
				roleEnum, ok := role.GetEnum()
				if !ok {
					t.Fatalf("currStruct[role].GetEnum() not successeded")
				}
				if roleEnum.Value.Cmp(big.NewInt(1)) != 0 {
					t.Errorf("currStruct[role].GetEnum().Value != 1, got %v", roleEnum.Value.String())
				}
				if roleEnum.Name != "Aboba" {
					t.Errorf("currStruct[role].GetEnum().Name != Aboba, got %v", roleEnum.Name)
				}

				oper1, ok := currStruct.GetField("oper1")
				if !ok {
					t.Fatalf("currStruct[oper1] not found")
				}
				oper1Enum, ok := oper1.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper1].GetEnum() not successeded")
				}
				if oper1Enum.Value.Cmp(big.NewInt(0)) != 0 {
					t.Errorf("currStruct[oper1].GetEnum().Value != 0, got %v", oper1Enum.Value.String())
				}
				if oper1Enum.Name != "Add" {
					t.Errorf("currStruct[oper1].GetEnum().Name != Add, got %v", oper1Enum.Name)
				}

				oper2, ok := currStruct.GetField("oper2")
				if !ok {
					t.Fatalf("currStruct[oper2] not found")
				}
				oper2Enum, ok := oper2.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper2].GetEnum() not successeded")
				}
				if oper2Enum.Value.Cmp(big.NewInt(-10000)) != 0 {
					t.Errorf("currStruct[oper2].GetEnum().Value != -10000, got %v", oper2Enum.Value.String())
				}
				if oper2Enum.Name != "TopUp" {
					t.Errorf("currStruct[oper2].GetEnum().Name != TopUp, got %v", oper2Enum.Name)
				}

				oper3, ok := currStruct.GetField("oper3")
				if !ok {
					t.Fatalf("currStruct[oper3] not found")
				}
				oper3Enum, ok := oper3.GetEnum()
				if !ok {
					t.Fatalf("currStruct[oper3].GetEnum() not successeded")
				}
				if oper3Enum.Value.Cmp(big.NewInt(1)) != 0 {
					t.Errorf("currStruct[oper3].GetEnum().Value != 1, got %v", oper3Enum.Value.String())
				}
				if oper3Enum.Name != "Something" {
					t.Errorf("currStruct[oper3].GetEnum().Name != Something, got %v", oper3Enum.Name)
				}
			},
		},
		{
			name:     "unmarshal a lot generics with struct",
			filename: "generics",
			cell:     "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
			t: Ty{
				SumType: "StructRef",
				StructRef: &StructRef{
					StructName: "ManyRefsMsg",
					TypeArgs: []Ty{
						{
							SumType: "UintN",
							UintN: &UintN{
								N: 16,
							},
						},
					},
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}

				payloadV, ok := currStruct.GetField("payload")
				if !ok {
					t.Fatalf("currStruct[payload] not found")
				}
				payloadStruct, ok := payloadV.GetStruct()
				if !ok {
					t.Fatalf("currStruct[payload].GetStruct() not found")
				}
				payloadRef, ok := payloadStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[payload][value] not found")
				}
				payloadUnion, ok := payloadRef.GetUnion()
				if !ok {
					t.Fatalf("currStruct[payload][value].GetUnion() not found")
				}
				payload, ok := payloadUnion.Val.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[payload].GetRefValue() not successeded")
				}
				payloadVal, ok := payload.GetSmallUInt()
				if !ok {
					t.Fatalf("currStruct[payload].GetSmallUInt() not found")
				}
				if payloadVal != 123 {
					t.Errorf("currStruct[payload].GetSmallUInt() != 123, got %v", payloadVal)
				}

				either, ok := currStruct.GetField("either")
				if !ok {
					t.Fatalf("currStruct[either] not found")
				}
				eitherUnion, ok := either.GetUnion()
				if !ok {
					t.Fatalf("currStruct[either].GetUnion() not successeded")
				}
				eitherStruct, ok := eitherUnion.Val.GetStruct()
				if !ok {
					t.Fatalf("currStruct[either].GetStruct() not successeded")
				}
				eitherV, ok := eitherStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[either][value] not successeded")
				}
				eitherVal, ok := eitherV.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[either][value].GetBigInt() not successeded")
				}
				if eitherVal.Cmp(big.NewInt(100000000)) != 0 {
					t.Fatalf("currStruct[either][value].GetBigInt() != 1000000000, got %v", eitherVal.String())
				}

				anotherEither, ok := currStruct.GetField("anotherEither")
				if !ok {
					t.Fatalf("currStruct[anotherEither] not found")
				}
				anotherEitherUnion, ok := anotherEither.GetUnion()
				if !ok {
					t.Fatalf("currStruct[anotherEither].GetUnion() not successeded")
				}
				anotherEitherStruct, ok := anotherEitherUnion.Val.GetStruct()
				if !ok {
					t.Fatalf("currStruct[anotherEither].GetStruct() not successeded")
				}
				anotherEitherV, ok := anotherEitherStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[anotherEither][value] not successeded")
				}
				anotherEitherVal, ok := anotherEitherV.GetTensor()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value].GetTensor() not successeded")
				}

				anotherEitherValBool, ok := anotherEitherVal[0].GetBool()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value][0].GetBool() not successeded")
				}
				if !anotherEitherValBool {
					t.Fatalf("currStruct[anotherEither][value][0].GetBool() is false")
				}
				anotherEitherValCoins, ok := anotherEitherVal[1].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value][0].GetBigInt() not successeded")
				}
				if anotherEitherValCoins.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Fatalf("currStruct[anotherEither][value][0].GetBigInt() != 1000000000, got %v", anotherEitherValCoins.String())
				}

				doubler, ok := currStruct.GetField("doubler")
				if !ok {
					t.Fatalf("currStruct[doubler] not found")
				}
				doublerRef, ok := doubler.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[doubler].GetRefValue() not successeded")
				}
				doublerTensor, ok := doublerRef.GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler].GetTensor() not successeded")
				}

				doublerTensor0, ok := doublerTensor[0].GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler][0] not successeded")
				}
				doublerTensor0Coins, ok := doublerTensor0[0].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[doubler][0][0].GetBigInt() not successeded")
				}
				if doublerTensor0Coins.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Fatalf("currStruct[doubler][0][0].GetBigInt() != 1000000000, got %v", doublerTensor0Coins.String())
				}

				doublerTensor0Addr, ok := doublerTensor0[1].GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[doubler][0][1].GetOptionalAddress() not successeded")
				}
				if doublerTensor0Addr.SumType != "NoneAddress" {
					t.Fatalf("currStruct[doubler][0][1].GetOptionalAddress() != NoneAddress")
				}

				doublerTensor1, ok := doublerTensor[1].GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler][1] not successeded")
				}
				doublerTensor1Coins, ok := doublerTensor1[0].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[doubler][1][0].GetBigInt() not successeded")
				}
				if doublerTensor1Coins.Cmp(big.NewInt(100_000_000)) != 0 {
					t.Fatalf("currStruct[doubler][1][0].GetBigInt() != 100000000, got %v", doublerTensor1Coins.String())
				}

				doublerTensor1Addr, ok := doublerTensor1[1].GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[doubler][1][1].GetOptionalAddress() not successeded")
				}
				if doublerTensor1Addr.SumType != "InternalAddress" {
					t.Fatalf("currStruct[doubler][1][1].GetOptionalAddress() != InternalAddress")
				}
				if doublerTensor1Addr.InternalAddress.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Fatalf("currStruct[doubler][1][1] != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", doublerTensor1Addr.InternalAddress.ToRaw())
				}

				myVal, ok := currStruct.GetField("myVal")
				if !ok {
					t.Fatalf("currStruct[myVal] not found")
				}
				myValVal, ok := myVal.GetSmallUInt()
				if !ok {
					t.Fatalf("currStruct[myVal].GetSmallUInt() not successed")
				}
				if myValVal != 16 {
					t.Fatalf("currStruct[myVal] != 16, got %v", myValVal)
				}
			},
		},
		{
			name:     "unmarshal a lot generics with alias",
			filename: "generics",
			cell:     "b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647",
			t: Ty{
				SumType: "AliasRef",
				AliasRef: &AliasRef{
					AliasName: "GoodNamingForMsg",
					TypeArgs: []Ty{
						{
							SumType: "UintN",
							UintN: &UintN{
								N: 16,
							},
						},
					},
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}

				payloadV, ok := currStruct.GetField("payload")
				if !ok {
					t.Fatalf("currStruct[payload] not found")
				}
				payloadStruct, ok := payloadV.GetStruct()
				if !ok {
					t.Fatalf("currStruct[payload].GetStruct() not found")
				}
				payloadRef, ok := payloadStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[payload][value] not found")
				}
				payloadUnion, ok := payloadRef.GetUnion()
				if !ok {
					t.Fatalf("currStruct[payload][value].GetUnion() not found")
				}
				payload, ok := payloadUnion.Val.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[payload].GetRefValue() not successeded")
				}
				payloadVal, ok := payload.GetSmallUInt()
				if !ok {
					t.Fatalf("currStruct[payload].GetSmallUInt() not found")
				}
				if payloadVal != 123 {
					t.Errorf("currStruct[payload].GetSmallUInt() != 123, got %v", payloadVal)
				}

				either, ok := currStruct.GetField("either")
				if !ok {
					t.Fatalf("currStruct[either] not found")
				}
				eitherUnion, ok := either.GetUnion()
				if !ok {
					t.Fatalf("currStruct[either].GetUnion() not successeded")
				}
				eitherStruct, ok := eitherUnion.Val.GetStruct()
				if !ok {
					t.Fatalf("currStruct[either].GetStruct() not successeded")
				}
				eitherV, ok := eitherStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[either][value] not successeded")
				}
				eitherVal, ok := eitherV.GetBigInt()
				if !ok {
					t.Fatalf("currStruct[either][value].GetBigInt() not successeded")
				}
				if eitherVal.Cmp(big.NewInt(100000000)) != 0 {
					t.Fatalf("currStruct[either][value].GetBigInt() != 1000000000, got %v", eitherVal.String())
				}

				anotherEither, ok := currStruct.GetField("anotherEither")
				if !ok {
					t.Fatalf("currStruct[anotherEither] not found")
				}
				anotherEitherUnion, ok := anotherEither.GetUnion()
				if !ok {
					t.Fatalf("currStruct[anotherEither].GetUnion() not successeded")
				}
				anotherEitherStruct, ok := anotherEitherUnion.Val.GetStruct()
				if !ok {
					t.Fatalf("currStruct[anotherEither].GetStruct() not successeded")
				}
				anotherEitherV, ok := anotherEitherStruct.GetField("value")
				if !ok {
					t.Fatalf("currStruct[anotherEither][value] not successeded")
				}
				anotherEitherVal, ok := anotherEitherV.GetTensor()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value].GetTensor() not successeded")
				}

				anotherEitherValBool, ok := anotherEitherVal[0].GetBool()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value][0].GetBool() not successeded")
				}
				if !anotherEitherValBool {
					t.Fatalf("currStruct[anotherEither][value][0].GetBool() is false")
				}
				anotherEitherValCoins, ok := anotherEitherVal[1].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[anotherEither][value][0].GetBigInt() not successeded")
				}
				if anotherEitherValCoins.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Fatalf("currStruct[anotherEither][value][0].GetBigInt() != 1000000000, got %v", anotherEitherValCoins.String())
				}

				doubler, ok := currStruct.GetField("doubler")
				if !ok {
					t.Fatalf("currStruct[doubler] not found")
				}
				doublerRef, ok := doubler.GetRefValue()
				if !ok {
					t.Fatalf("currStruct[doubler].GetRefValue() not successeded")
				}
				doublerTensor, ok := doublerRef.GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler].GetTensor() not successeded")
				}

				doublerTensor0, ok := doublerTensor[0].GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler][0] not successeded")
				}
				doublerTensor0Coins, ok := doublerTensor0[0].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[doubler][0][0].GetBigInt() not successeded")
				}
				if doublerTensor0Coins.Cmp(big.NewInt(1_000_000_000)) != 0 {
					t.Fatalf("currStruct[doubler][0][0].GetBigInt() != 1000000000, got %v", doublerTensor0Coins.String())
				}

				doublerTensor0Addr, ok := doublerTensor0[1].GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[doubler][0][1].GetOptionalAddress() not successeded")
				}
				if doublerTensor0Addr.SumType != "NoneAddress" {
					t.Fatalf("currStruct[doubler][0][1].GetOptionalAddress() != NoneAddress")
				}

				doublerTensor1, ok := doublerTensor[1].GetTensor()
				if !ok {
					t.Fatalf("currStruct[doubler][1] not successeded")
				}
				doublerTensor1Coins, ok := doublerTensor1[0].GetBigInt()
				if !ok {
					t.Fatalf("currStruct[doubler][1][0].GetBigInt() not successeded")
				}
				if doublerTensor1Coins.Cmp(big.NewInt(100_000_000)) != 0 {
					t.Fatalf("currStruct[doubler][1][0].GetBigInt() != 100000000, got %v", doublerTensor1Coins.String())
				}

				doublerTensor1Addr, ok := doublerTensor1[1].GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[doubler][1][1].GetOptionalAddress() not successeded")
				}
				if doublerTensor1Addr.SumType != "InternalAddress" {
					t.Fatalf("currStruct[doubler][1][1].GetOptionalAddress() != InternalAddress")
				}
				if doublerTensor1Addr.InternalAddress.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Fatalf("currStruct[doubler][1][1] != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", doublerTensor1Addr.InternalAddress.ToRaw())
				}

				myVal, ok := currStruct.GetField("myVal")
				if !ok {
					t.Fatalf("currStruct[myVal] not found")
				}
				myValVal, ok := myVal.GetSmallUInt()
				if !ok {
					t.Fatalf("currStruct[myVal].GetSmallUInt() not successed")
				}
				if myValVal != 16 {
					t.Fatalf("currStruct[myVal] != 16, got %v", myValVal)
				}
			},
		},
		{
			name:     "unmarshal a struct with default values",
			filename: "default_values",
			cell:     "b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde",
			t: Ty{
				SumType: "StructRef",
				StructRef: &StructRef{
					StructName: "DefaultTest",
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}

				optNum1, ok := currStruct.GetField("num1")
				if !ok {
					t.Fatalf("currStruct[num1] not found")
				}
				num1, ok := optNum1.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[num1].GetOptionalValue() not successeded")
				}
				if !num1.IsExists {
					t.Fatalf("currStruct[num1] is not exists")
				}
				num1Val, ok := num1.Val.GetSmallUInt()
				if !ok {
					t.Fatalf("currStruct[num1].GetSmallUInt() not successeded")
				}
				if num1Val != 4 {
					t.Fatalf("currStruct[num1].GetSmallUInt() != 4, got %v", num1Val)
				}

				optNum2, ok := currStruct.GetField("num2")
				if !ok {
					t.Fatalf("currStruct[num2] not found")
				}
				num2, ok := optNum2.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[num2].GetOptionalValue() not successeded")
				}
				if !num2.IsExists {
					t.Fatalf("currStruct[num2] is not exists")
				}
				num2Val, ok := num2.Val.GetSmallInt()
				if !ok {
					t.Fatalf("currStruct[num2].GetSmallInt() not successeded")
				}
				if num2Val != 5 {
					t.Fatalf("currStruct[num2].GetSmallInt() != 5, got %v", num2Val)
				}

				optSlice3, ok := currStruct.GetField("slice3")
				if !ok {
					t.Fatalf("currStruct[slice3] not found")
				}
				slice3, ok := optSlice3.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[slice3].GetOptionalValue() not successeded")
				}
				if !slice3.IsExists {
					t.Fatalf("currStruct[slice3] is not exists")
				}
				slice3Val, ok := slice3.Val.GetCell()
				if !ok {
					t.Fatalf("currStruct[slice3].GetCell() not successeded")
				}
				hs, err := slice3Val.HashString()
				if err != nil {
					t.Fatal(err)
				}
				if hs != "55e960f1409af0d7670e382c61276a559fa9330185984d91faffebf32d5fa383" {
					t.Fatalf("currStruct[slice3].GetCell().Hash() != 55e960f1409af0d7670e382c61276a559fa9330185984d91faffebf32d5fa383, got %v", hs)
				}

				optAddr4, ok := currStruct.GetField("addr4")
				if !ok {
					t.Fatalf("currStruct[addr4] not found")
				}
				addr4, ok := optAddr4.GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[addr4].GetOptionalAddress() not successeded")
				}
				if addr4.SumType != "NoneAddress" {
					t.Fatalf("currStruct[addr4].GetOptionalAddress() != NoneAddress")
				}

				optAddr5, ok := currStruct.GetField("addr5")
				if !ok {
					t.Fatalf("currStruct[addr5] not found")
				}
				addr5, ok := optAddr5.GetOptionalAddress()
				if !ok {
					t.Fatalf("currStruct[addr5].GetOptionalAddress() not successeded")
				}
				if addr5.SumType != "InternalAddress" {
					t.Fatalf("currStruct[addr5].GetOptionalAddress() != InternalAddress")
				}
				if addr5.InternalAddress.ToRaw() != "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe" {
					t.Fatalf("currStruct[addr5].GetOptionalAddress() != 0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe, got %v", addr5.InternalAddress.ToRaw())
				}

				optTensor6, ok := currStruct.GetField("tensor6")
				if !ok {
					t.Fatalf("currStruct[tensor6] not found")
				}
				tensor6, ok := optTensor6.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[tensor6].GetOptionalValue() not successeded")
				}
				if !tensor6.IsExists {
					t.Fatalf("currStruct[tensor6] is not exists")
				}
				tensor6Val, ok := tensor6.Val.GetTensor()
				if !ok {
					t.Fatalf("currStruct[tensor6].GetTensor() not successeded")
				}
				tensor6Val0, ok := tensor6Val[0].GetSmallInt()
				if !ok {
					t.Fatalf("currStruct[tensor6][0].GetSmallInt() not successed")
				}
				if tensor6Val0 != 342 {
					t.Fatalf("currStruct[tensor6][0] != 342, got %v", tensor6Val0)
				}

				tensor6Val1, ok := tensor6Val[1].GetBool()
				if !ok {
					t.Fatalf("currStruct[tensor6][1].GetBool() not successed")
				}
				if !tensor6Val1 {
					t.Fatalf("currStruct[tensor6][0] is false")
				}

				optNum7, ok := currStruct.GetField("num7")
				if !ok {
					t.Fatalf("currStruct[num7] not found")
				}
				num7, ok := optNum7.GetOptionalValue()
				if !ok {
					t.Fatalf("currStruct[num7].GetOptionalValue() not successeded")
				}
				if num7.IsExists {
					t.Fatalf("currStruct[num7] exists")
				}
			},
		},
		{
			name:     "unmarshal numbers",
			filename: "numbers",
			cell:     "b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014",
			t: Ty{
				SumType: "StructRef",
				StructRef: &StructRef{
					StructName: "Numbers",
				},
			},
			check: func(v Value) {
				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}
				num1, ok := currStruct.GetField("num1")
				if !ok {
					t.Fatalf("num1 not found")
				}
				val1, ok := num1.GetSmallUInt()
				if !ok {
					t.Fatalf("num1.GetSmallUInt() not successeded")
				}
				if val1 != 0 {
					t.Fatalf("num1 != 0, got %v", val1)
				}

				num3, ok := currStruct.GetField("num3")
				if !ok {
					t.Fatalf("num3 not found")
				}
				val3, ok := num3.GetBigInt()
				if !ok {
					t.Fatalf("num3.GetBigInt() not successeded")
				}
				if val3.Cmp(big.NewInt(241)) != 0 {
					t.Fatalf("num3 != 241, got %v", val3.String())
				}

				num4, ok := currStruct.GetField("num4")
				if !ok {
					t.Fatalf("num4 not found")
				}
				val4, ok := num4.GetBigInt()
				if !ok {
					t.Fatalf("num4.GetSmallUInt() not successeded")
				}
				if val4.Cmp(big.NewInt(3421)) != 0 {
					t.Fatalf("num4 != 3421, got %s", val4.String())
				}

				num5, ok := currStruct.GetField("num5")
				if !ok {
					t.Fatalf("num5 not found")
				}
				val5, ok := num5.GetBool()
				if !ok {
					t.Fatalf("num5.GetBool() not successeded")
				}
				if !val5 {
					t.Fatalf("num5 != true")
				}

				num7, ok := currStruct.GetField("num7")
				if !ok {
					t.Fatalf("num7 not found")
				}
				val7, ok := num7.GetBits()
				if !ok {
					t.Fatalf("num7.GetBits() not successeded")
				}
				if !bytes.Equal(val7.Buffer(), []byte{49, 50}) {
					t.Fatalf("num7 != \"12\", got %v", val7.Buffer())
				}

				num8, ok := currStruct.GetField("num8")
				if !ok {
					t.Fatalf("num8 not found")
				}
				val8, ok := num8.GetSmallInt()
				if !ok {
					t.Fatalf("num8.GetSmallInt() not successeded")
				}
				if val8 != 0 {
					t.Fatalf("num8 != 0, got %v", val8)
				}

				num9, ok := currStruct.GetField("num9")
				if !ok {
					t.Fatalf("num9 not found")
				}
				val9, ok := num9.GetBigInt()
				if !ok {
					t.Fatalf("num9.GetSmallInt() not successeded")
				}
				if val9.Cmp(big.NewInt(2342)) != 0 {
					t.Fatalf("num9 != 2342, got %s", val9.String())
				}
			},
		},
		{
			name:     "unmarshal random fields",
			filename: "random_fields",
			cell:     "b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb",
			t: Ty{
				SumType: "StructRef",
				StructRef: &StructRef{
					StructName: "RandomFields",
				},
			},
			check: func(v Value) {
				addr := ton.MustParseAccountID("UQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqEBI")

				currStruct, ok := v.GetStruct()
				if !ok {
					t.Fatalf("struct not found")
				}
				pref, ok := currStruct.GetStructPrefix()
				if !ok {
					t.Fatalf("struct prefix not found")
				}
				if pref.Len != 12 {
					t.Fatalf("pref.Len != 12, got %d", pref.Len)
				}
				if pref.Prefix != 1940 {
					t.Fatalf("struct prefix != 1940, got %d", pref)
				}

				destInt, ok := currStruct.GetField("dest_int")
				if !ok {
					t.Fatalf("dest_int not found")
				}
				destIntVal, ok := destInt.GetAddress()
				if !ok {
					t.Fatalf("num1.GetAddress() not successeded")
				}
				if destIntVal.ToRaw() != addr.ToRaw() {
					t.Fatalf("destInt != UQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqEBI")
				}

				amount, ok := currStruct.GetField("amount")
				if !ok {
					t.Fatalf("amount not found")
				}
				amountVal, ok := amount.GetBigInt()
				if !ok {
					t.Fatalf("amount.GetBigInt() not successeded")
				}
				expectedAmount, ok := big.NewInt(0).SetString("500000123400000", 10)
				if !ok {
					t.Fatalf("cannot set 500000123400000 value to big.Int")
				}
				if amountVal.Cmp(expectedAmount) != 0 {
					t.Fatalf("amount != 500000123400000, got %s", amountVal.String())
				}

				destExt, ok := currStruct.GetField("dest_ext")
				if !ok {
					t.Fatalf("num4 dest_ext found")
				}
				destExtVal, ok := destExt.GetAnyAddress()
				if !ok {
					t.Fatalf("num4.GetAnyAddress() not successeded")
				}
				if destExtVal.SumType != "NoneAddress" {
					t.Fatalf("destExt != a none address")
				}

				intVector, ok := currStruct.GetField("intVector")
				if !ok {
					t.Fatalf("intVector not found")
				}
				intVectorVal, ok := intVector.GetTensor()
				if !ok {
					t.Fatalf("num5.GetTensor() not successeded")
				}
				val1, ok := intVectorVal[0].GetSmallInt()
				if !ok {
					t.Fatalf("intVector[0].GetSmallInt() not successeded")
				}
				if val1 != 342 {
					t.Fatalf("intVector[0].GetSmallInt() != 342, got %v", val1)
				}

				optVal2, ok := intVectorVal[1].GetOptionalValue()
				if !ok {
					t.Fatalf("intVector[1].GetOptionalValue() not successeded")
				}
				if !optVal2.IsExists {
					t.Fatalf("intVector[1].GetOptionalValue() != exists")
				}
				val2, ok := optVal2.Val.GetBigInt()
				if !ok {
					t.Fatalf("intVector[1].GetOptionalValue().GetBigInt() not successeded")
				}
				if val2.Cmp(big.NewInt(1000000000)) != 0 {
					t.Fatalf("intVector[1].GetOptionalValue().GetBigInt() != 1000000000, got %v", val1)
				}

				val3, ok := intVectorVal[2].GetSmallUInt()
				if !ok {
					t.Fatalf("intVector[2].GetSmallUInt() not successeded")
				}
				if val3 != 23443 {
					t.Fatalf("intVector[2].GetSmallUInt() != 23443, got %v", val1)
				}

				needsMoreRef, ok := currStruct.GetField("needs_more")
				if !ok {
					t.Fatalf("needs_more not found")
				}
				needsMore, ok := needsMoreRef.GetRefValue()
				if !ok {
					t.Fatalf("needsMoreRef.GetRefValue() not successeded")
				}
				needsMoreVal, ok := needsMore.GetBool()
				if !ok {
					t.Fatalf("needsMore.GetBool() not successeded")
				}
				if !needsMoreVal {
					t.Fatalf("needsMore != true")
				}

				somePayload, ok := currStruct.GetField("some_payload")
				if !ok {
					t.Fatalf("some_payload not found")
				}
				somePayloadVal, ok := somePayload.GetCell()
				if !ok {
					t.Fatalf("num8.GetCell() not successeded")
				}
				somePayloadHash, err := somePayloadVal.HashString()
				if err != nil {
					t.Fatalf("somePayload.HashString() not successeded")
				}
				if somePayloadHash != "f2017ee9d429c16689ba2243d26d2a070a1e8e4a6106cee2129a049deee727d9" {
					t.Fatalf("somePayloadHash != f2017ee9d429c16689ba2243d26d2a070a1e8e4a6106cee2129a049deee727d9, got %v", somePayloadHash)
				}

				myInt, ok := currStruct.GetField("my_int")
				if !ok {
					t.Fatalf("my_int not found")
				}
				myIntVal, ok := myInt.GetSmallInt()
				if !ok {
					t.Fatalf("my_int.GetSmallInt() not successeded")
				}
				if myIntVal != 432 {
					t.Fatalf("my_int != 432, got %v", myIntVal)
				}

				someUnion, ok := currStruct.GetField("some_union")
				if !ok {
					t.Fatalf("my_int not found")
				}
				someUnionVal, ok := someUnion.GetUnion()
				if !ok {
					t.Fatalf("someUnion.GetSmallInt() not successeded")
				}
				unionVal, ok := someUnionVal.Val.GetSmallInt()
				if !ok {
					t.Fatalf("someUnion.GetSmallInt() not successeded")
				}
				if unionVal != 30000 {
					t.Fatalf("some_union != 30000, got %v", someUnionVal)
				}

				default1, ok := currStruct.GetField("default_1")
				if !ok {
					t.Fatalf("default_1 not found")
				}
				default1Val, ok := default1.GetSmallInt()
				if !ok {
					t.Fatalf("default1.GetSmallInt() not successeded")
				}
				if default1Val != 1 {
					t.Fatalf("default1 != 1, got %v", default1Val)
				}

				optDefault2, ok := currStruct.GetField("default_2")
				if !ok {
					t.Fatalf("default_2 not found")
				}
				default2, ok := optDefault2.GetOptionalValue()
				if !ok {
					t.Fatalf("default2.GetOptionalValue() not successeded")
				}
				if !default2.IsExists {
					t.Fatalf("default2.GetOptionalValue() != exists")
				}
				default2Val, ok := default2.Val.GetSmallInt()
				if !ok {
					t.Fatalf("default2.GetSmallInt() not successeded")
				}
				if default2Val != 55 {
					t.Fatalf("default2 != 55, got %v", default2Val)
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			inputFilename := fmt.Sprintf("testdata/%s.json", c.filename)
			data, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Fatal(err)
			}

			var abi ABI
			err = json.Unmarshal(data, &abi)
			if err != nil {
				t.Fatal(err)
			}

			currCell, err := boc.DeserializeBocHex(c.cell)
			if err != nil {
				t.Fatal(err)
			}
			decoder := NewDecoder()
			decoder.WithABI(abi)
			val, err := decoder.UnmarshalTolk(currCell[0], c.t)
			if err != nil {
				t.Fatal(err)
			}
			c.check(*val)
		})
	}
}
