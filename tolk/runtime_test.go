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
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/ton"
)

const jsonFilesPath = "testdata/json/"

func TestRuntime_UnmarshalSmallInt(t *testing.T) {
	inputFilename := "small_int"
	ty := tolkParser.Ty{
		SumType: "IntN",
		IntN: &tolkParser.IntN{
			N: 24,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010005000006ff76c41616db06")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetSmallInt()
	if !ok {
		t.Errorf("v.GetSmallInt() not successeded")
	}
	if val != -35132 {
		t.Errorf("val != -35132, got %v", val)
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBigInt(t *testing.T) {
	inputFilename := "big_int"
	ty := tolkParser.Ty{
		SumType: "IntN",
		IntN: &tolkParser.IntN{
			N: 183,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetBigInt()
	if !ok {
		t.Errorf("v.GetBigInt() not successeded")
	}
	if val.Cmp(big.NewInt(-3513294376431)) != 0 {
		t.Errorf("val != -3513294376431, got %v", val)
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalSmallUInt(t *testing.T) {
	inputFilename := "small_uint"
	ty := tolkParser.Ty{
		SumType: "UintN",
		UintN: &tolkParser.UintN{
			N: 53,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000d00000000001d34e435eafd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetSmallUInt()
	if !ok {
		t.Errorf("v.GetSmallUInt() not successeded")
	}
	if val != 934 {
		t.Errorf("val != 934, got %v", val)
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBigUInt(t *testing.T) {
	inputFilename := "big_uint"
	ty := tolkParser.Ty{
		SumType: "UintN",
		UintN: &tolkParser.UintN{
			N: 257,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002300004100000000000000000000000000000000000000000000000000009fc4212a38ba40b11cce12")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetBigUInt()
	if !ok {
		t.Errorf("v.GetBigUInt() not successeded")
	}
	if val.Cmp(big.NewInt(351329437643124)) != 0 {
		t.Errorf("val != 351329437643124, got %v", val.String())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalVarInt(t *testing.T) {
	inputFilename := "var_int"
	ty := tolkParser.Ty{
		SumType: "VarIntN",
		VarIntN: &tolkParser.VarIntN{
			N: 16,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000730c98588449b6923")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetVarInt()
	if !ok {
		t.Errorf("v.GetVarInt() not successeded")
	}
	if val.Cmp(big.NewInt(825432)) != 0 {
		t.Errorf("val != 825432, got %v", val.String())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalVarUInt(t *testing.T) {
	inputFilename := "var_uint"
	ty := tolkParser.Ty{
		SumType: "VarUintN",
		VarUintN: &tolkParser.VarUintN{
			N: 32,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000800000b28119ab36b44d3a86c0f")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetVarUInt()
	if !ok {
		t.Errorf("v.GetVarUInt() not successeded")
	}
	if val.Cmp(big.NewInt(9451236712)) != 0 {
		t.Errorf("val != 9451236712, got %v", val.String())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBits(t *testing.T) {
	inputFilename := "bits"
	ty := tolkParser.Ty{
		SumType: "BitsN",
		BitsN: &tolkParser.BitsN{
			N: 24,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000500000631323318854035")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetBits()
	if !ok {
		t.Errorf("v.GetBits() not successeded")
	}
	if bytes.Equal(val.Buffer(), []byte{55, 56, 57}) {
		t.Errorf("val != {55, 56, 57}, got %v", val)
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalCoins(t *testing.T) {
	inputFilename := "coins"
	ty := tolkParser.Ty{
		SumType: "Coins",
		Coins:   &tolkParser.Coins{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010007000009436ec6e0189ebbd7f4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetCoins()
	if !ok {
		t.Errorf("v.GetCoins() not successeded")
	}
	if val.Cmp(big.NewInt(921464321)) != 0 {
		t.Errorf("val != 921464321, got %v", val)
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBool(t *testing.T) {
	inputFilename := "bool"
	ty := tolkParser.Ty{
		SumType: "Bool",
		Bool:    &tolkParser.Bool{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000300000140f6d24034")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetBool()
	if !ok {
		t.Errorf("v.GetBool() not successeded")
	}
	if val {
		t.Error("val is true")
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalCell(t *testing.T) {
	inputFilename := "cell"
	ty := tolkParser.Ty{
		SumType: "Cell",
		Cell:    &tolkParser.Cell{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101020100090001000100080000007ba52a3292")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalRemaining(t *testing.T) {
	inputFilename := "remaining"
	ty := tolkParser.Ty{
		SumType:   "Remaining",
		Remaining: &tolkParser.Remaining{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetRemaining()
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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAddress(t *testing.T) {
	inputFilename := "internal_address"
	ty := tolkParser.Ty{
		SumType: "Address",
		Address: &tolkParser.Address{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetAddress()
	if !ok {
		t.Errorf("v.GetAddress() not successeded")
	}
	if val.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
		t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.ToRaw())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalNotExitsOptionalAddress(t *testing.T) {
	inputFilename := "not_exists_optional_address"
	ty := tolkParser.Ty{
		SumType:    "AddressOpt",
		AddressOpt: &tolkParser.AddressOpt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100030000012094418655")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetOptionalAddress()
	if !ok {
		t.Errorf("v.GetOptionalAddress() not successeded")
	}

	if val.SumType != "NoneAddress" {
		t.Errorf("val.GetAddress() != none address")
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalExistsOptionalAddress(t *testing.T) {
	inputFilename := "exists_optional_address"
	ty := tolkParser.Ty{
		SumType:    "AddressOpt",
		AddressOpt: &tolkParser.AddressOpt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetOptionalAddress()
	if !ok {
		t.Errorf("v.GetOptionalAddress() not successeded")
	}

	if val.SumType == "InternalAddress" && val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
		t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalExternalAddress(t *testing.T) {
	inputFilename := "external_address"
	ty := tolkParser.Ty{
		SumType:    "AddressExt",
		AddressExt: &tolkParser.AddressExt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000742082850fcbd94fd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetExternalAddress()
	if !ok {
		t.Errorf("v.GetExternalAddress() not successeded")
	}
	addressPart := boc.NewBitString(16)
	err = addressPart.WriteBytes([]byte{97, 98})
	if err != nil {
		t.Fatal(err)
	}
	if val.Len != 8 && bytes.Equal(val.Address.Buffer(), []byte{97, 98}) {
		t.Errorf("val.GetExternalAddress() != {97, 98}, got %v", val.Address.Buffer())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAnyNoneAddress(t *testing.T) {
	inputFilename := "any_none_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100030000012094418655")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetAnyAddress()
	if !ok {
		t.Errorf("v.GetAnyAddress() not successeded")
	}
	if val.SumType != "NoneAddress" {
		t.Errorf("val.GetAddress() != none address")
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAnyInternalAddress(t *testing.T) {
	inputFilename := "any_internal_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetAnyAddress()
	if !ok {
		t.Errorf("v.GetAnyAddress() not successeded")
	}
	if val.SumType == "InternalAddress" && val.InternalAddress.ToRaw() != "0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8" {
		t.Errorf("val.GetAddress() != 0:83dfd552e63729b472fcbcc8c45ebcc6691702558b68ec7527e1ba403a0f31a8, got %v", val.InternalAddress.ToRaw())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAnyExternalAddress(t *testing.T) {
	inputFilename := "any_external_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000742082850fcbd94fd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetAnyAddress()
	if !ok {
		t.Errorf("v.GetAnyAddress() not successeded")
	}
	addressPart := boc.NewBitString(16)
	err = addressPart.WriteBytes([]byte{97, 98})
	if err != nil {
		t.Fatal(err)
	}
	if val.SumType == "ExternalAddress" && val.ExternalAddress.Len != 8 && bytes.Equal(val.ExternalAddress.Address.Buffer(), []byte{97, 98}) {
		t.Errorf("val.GetExternalAddress() != {97, 98}, got %v", val.ExternalAddress.Address.Buffer())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAnyVarAddress(t *testing.T) {
	inputFilename := "any_var_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalNotExistsNullable(t *testing.T) {
	inputFilename := "not_exists_nullable"
	ty := tolkParser.Ty{
		SumType: "Nullable",
		Nullable: &tolkParser.Nullable{
			Inner: tolkParser.Ty{
				SumType:   "Remaining",
				Remaining: &tolkParser.Remaining{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000300000140f6d24034")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetOptionalValue()
	if !ok {
		t.Errorf("v.GetOptionalValue() not successeded")
	}
	if val.IsExists {
		t.Errorf("v.GetOptionalValue() is exists")
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalExistsNullable(t *testing.T) {
	inputFilename := "exists_nullable"
	ty := tolkParser.Ty{
		SumType: "Nullable",
		Nullable: &tolkParser.Nullable{
			Inner: tolkParser.Ty{
				SumType: "Cell",
				Cell:    &tolkParser.Cell{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000b000101c001000900000c0ae007880db9")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalRef(t *testing.T) {
	inputFilename := "ref"
	ty := tolkParser.Ty{
		SumType: "CellOf",
		CellOf: &tolkParser.CellOf{
			Inner: tolkParser.Ty{
				SumType: "IntN",
				IntN: &tolkParser.IntN{
					N: 65,
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000e000100010011000000000009689e40e150b4c5")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalEmptyTensor(t *testing.T) {
	inputFilename := "empty_tensor"
	ty := tolkParser.Ty{
		SumType: "Tensor",
		Tensor:  &tolkParser.Tensor{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100020000004cacb9cd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetTensor()
	if !ok {
		t.Errorf("v.GetTensor() not successeded")
	}

	if len(val) != 0 {
		t.Errorf("v.GetTensor() != empty")
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalTensor(t *testing.T) {
	inputFilename := "tensor"
	ty := tolkParser.Ty{
		SumType: "Tensor",
		Tensor: &tolkParser.Tensor{
			Items: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 123,
					},
				},
				{
					SumType: "Bool",
					Bool:    &tolkParser.Bool{},
				},
				{
					SumType: "Coins",
					Coins:   &tolkParser.Coins{},
				},
				{
					SumType: "Tensor",
					Tensor: &tolkParser.Tensor{
						Items: []tolkParser.Ty{
							{
								SumType: "IntN",
								IntN: &tolkParser.IntN{
									N: 23,
								},
							},
							{
								SumType: "Nullable",
								Nullable: &tolkParser.Nullable{
									Inner: tolkParser.Ty{
										SumType: "IntN",
										IntN: &tolkParser.IntN{
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
					VarIntN: &tolkParser.VarIntN{
						N: 32,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetTensor()
	if !ok {
		t.Errorf("v.GetTensor() not successeded")
	}

	val0, ok := val[0].GetBigUInt()
	if !ok {
		t.Errorf("val[0].GetBigUInt() not successeded")
	}
	if val0.Cmp(big.NewInt(4325)) != 0 {
		t.Errorf("val[0].GetBigUInt() != 4325, got %v", val0.String())
	}

	val1, ok := val[1].GetBool()
	if !ok {
		t.Errorf("val[1].GetBool() not successeded")
	}
	if !val1 {
		t.Error("val[1].GetBool() is false")
	}

	val2, ok := val[2].GetCoins()
	if !ok {
		t.Errorf("val[2].GetCoins() not successeded")
	}
	if val2.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Errorf("val[2].GetCoins() != 1000000000, got %v", val2.String())
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

	val4, ok := val[4].GetVarInt()
	if !ok {
		t.Errorf("val[4].GetVarInt() not successeded")
	}
	if val4.Cmp(big.NewInt(-9_304_000_000)) != 0 {
		t.Errorf("val[4].GetVarInt() != -9304000000, got %v", val4.String())
	}

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalIntKeyMap(t *testing.T) {
	inputFilename := "int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "IntN",
				IntN: &tolkParser.IntN{
					N: 32,
				},
			},
			V: tolkParser.Ty{
				SumType: "Bool",
				Bool:    &tolkParser.Bool{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000c000101c001000ba00000007bc09a662c32")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalUIntKeyMap(t *testing.T) {
	inputFilename := "uint_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "UintN",
				UintN: &tolkParser.UintN{
					N: 16,
				},
			},
			V: tolkParser.Ty{
				SumType: "Address",
				Address: &tolkParser.Address{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBigIntKeyMap(t *testing.T) {
	inputFilename := "big_int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "UintN",
				UintN: &tolkParser.UintN{
					N: 78,
				},
			},
			V: tolkParser.Ty{
				SumType: "Cell",
				Cell:    &tolkParser.Cell{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetMap()
	if !ok {
		t.Errorf("v.GetMap() not successeded")
	}
	val1, ok := val.GetByBigUInt(BigUInt(*big.NewInt(2337412)))
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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalBitsKeyMap(t *testing.T) {
	inputFilename := "bits_int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "BitsN",
				BitsN: &tolkParser.BitsN{
					N: 16,
				},
			},
			V: tolkParser.Ty{
				SumType: "Map",
				Map: &tolkParser.Map{
					K: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 64,
						},
					},
					V: tolkParser.Ty{
						SumType: "Tensor",
						Tensor: &tolkParser.Tensor{
							Items: []tolkParser.Ty{
								{
									SumType: "Address",
									Address: &tolkParser.Address{},
								},
								{
									SumType: "Coins",
									Coins:   &tolkParser.Coins{},
								},
							},
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := v.GetMap()
	if !ok {
		t.Errorf("v.GetMap() not successeded")
	}
	key1 := boc.NewBitString(16)
	err = key1.WriteBytes([]byte{65, 66})
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

	val1_124Val1, ok := val1_124Val[1].GetCoins()
	if !ok {
		t.Errorf("val[{97, 98}][124][1].GetCoins() not successeded")
	}
	if val1_124Val1.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Errorf("val[{97, 98}][124][1].GetCoins() != 1_000_000_000, got %v", val1_124Val1.String())
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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalAddressKeyMap(t *testing.T) {
	inputFilename := "address_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "Address",
				Address: &tolkParser.Address{},
			},
			V: tolkParser.Ty{
				SumType: "Coins",
				Coins:   &tolkParser.Coins{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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
	val1Val, ok := val1.GetCoins()
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

	pathPrefix := jsonFilesPath + inputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalUnionWithDecPrefix(t *testing.T) {
	jsonInputFilename := "union_with_dec_prefix"
	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr:        "0",
					PrefixLen:        1,
					PrefixEatInPlace: true,
					VariantTy: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 16,
						},
					},
				},
				{
					PrefixStr:        "1",
					PrefixLen:        1,
					PrefixEatInPlace: true,
					VariantTy: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 128,
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalUnionWithBinPrefix(t *testing.T) {
	jsonInputFilename := "union_with_bin_prefix"
	inputFilename := "testdata/bin_union.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr: "0b001",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "AddressWithPrefix",
						},
					},
				},
				{
					PrefixStr: "0b011",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "MapWithPrefix",
						},
					},
				},
				{
					PrefixStr: "0b111",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "CellWithPrefix",
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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
	mapCoins, ok := mapVal.GetCoins()
	if !ok {
		t.Errorf("val.GetMap()[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"].GetCoins() not successeded")
	}
	if mapCoins.Cmp(big.NewInt(43213412)) != 0 {
		t.Errorf("val.GetMap()[\"EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs\"].GetCoins() != 43213412, got %v", mapCoins.String())
	}

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalUnionWithHexPrefix(t *testing.T) {
	jsonInputFilename := "union_with_hex_prefix"
	inputFilename := "testdata/hex_union.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr: "0x12345678",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt66WithPrefix",
						},
					},
				},
				{
					PrefixStr: "0xdeadbeef",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt33WithPrefix",
						},
					},
				},
				{
					PrefixStr: "0x89abcdef",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt4WithPrefix",
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000b000011deadbeef00000000c0d75977b9")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotRefsFromAlias(t *testing.T) {
	jsonInputFilename := "a_lot_refs_from_alias"
	inputFilename := "testdata/refs.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "AliasRef",
		AliasRef: &tolkParser.AliasRef{
			AliasName: "GoodNamingForMsg",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	currAlias, ok := v.GetAlias()
	if !ok {
		t.Errorf("v.GetAlias() not successeded")
	}
	currStruct, ok := currAlias.GetStruct()
	if !ok {
		t.Fatalf("struct not found")
	}
	pref, ok := currStruct.GetPrefix()
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
	user1Alias, ok := user1.GetAlias()
	if !ok {
		t.Fatalf("currStruct[user1].GetAlias() not found")
	}
	user1Val, ok := user1Alias.GetStruct()
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
	user1BalanceVal, ok := user1Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user1][balance].GetCoins() not successeded")
	}
	if user1BalanceVal.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Errorf("currStruct[user1][balance].GetCoins() != 1000000000, got %v", user1BalanceVal.String())
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
	user2Alias, ok := user2Ref.GetAlias()
	if !ok {
		t.Fatalf("currStruct[user2].GetAlias() not found")
	}
	user2Val, ok := user2Alias.GetStruct()
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
	user2BalanceVal, ok := user2Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user2][balance].GetCoins() not successeded")
	}
	if user2BalanceVal.Cmp(big.NewInt(100_000_000)) != 0 {
		t.Errorf("currStruct[user2][balance].GetCoins() != 100000000, got %v", user2BalanceVal.String())
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
	user5BalanceVal, ok := user5Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user5][balance].GetCoins() not successeded")
	}
	if user5BalanceVal.Cmp(big.NewInt(10_000_000_000_000)) != 0 {
		t.Errorf("currStruct[user5][balance].GetCoins() != 10000000000000, got %v", user5BalanceVal.String())
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotRefsFromStruct(t *testing.T) {
	jsonInputFilename := "a_lot_refs_from_struct"
	inputFilename := "testdata/refs.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "ManyRefsMsg",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	currStruct, ok := v.GetStruct()
	if !ok {
		t.Fatalf("struct not found")
	}
	pref, ok := currStruct.GetPrefix()
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
	user1Alias, ok := user1.GetAlias()
	if !ok {
		t.Fatalf("currStruct[user1].GetAlias() not found")
	}
	user1Val, ok := user1Alias.GetStruct()
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
	user1BalanceVal, ok := user1Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user1][balance].GetCoins() not successeded")
	}
	if user1BalanceVal.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Errorf("currStruct[user1][balance].GetCoins() != 1000000000, got %v", user1BalanceVal.String())
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
	user2Alias, ok := user2Ref.GetAlias()
	if !ok {
		t.Fatalf("currStruct[user2].GetAlias() not found")
	}
	user2Val, ok := user2Alias.GetStruct()
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
	user2BalanceVal, ok := user2Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user2][balance].GetCoins() not successeded")
	}
	if user2BalanceVal.Cmp(big.NewInt(100_000_000)) != 0 {
		t.Errorf("currStruct[user2][balance].GetCoins() != 100000000, got %v", user2BalanceVal.String())
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
	user5BalanceVal, ok := user5Balance.GetCoins()
	if !ok {
		t.Fatalf("currStruct[user5][balance].GetCoins() not successeded")
	}
	if user5BalanceVal.Cmp(big.NewInt(10_000_000_000_000)) != 0 {
		t.Errorf("currStruct[user5][balance].GetCoins() != 10000000000000, got %v", user5BalanceVal.String())
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotGenericsFromStruct(t *testing.T) {
	jsonInputFilename := "a_lot_generics_from_struct"
	inputFilename := "testdata/generics.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "ManyRefsMsg",
			TypeArgs: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 16,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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
	// todo: remove GetGeneric because its not convenient
	payloadGeneric, ok := payload.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[payload].GetGeneric() not found")
	}
	payloadVal, ok := payloadGeneric.GetSmallUInt()
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
	eitherAlias, ok := either.GetAlias()
	if !ok {
		t.Fatalf("currStruct[either].GetAlias() not found")
	}
	eitherUnion, ok := eitherAlias.GetUnion()
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
	eitherValGeneric, ok := eitherV.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[either][value].GetGeneric() not found")
	}
	eitherVal, ok := eitherValGeneric.GetCoins()
	if !ok {
		t.Fatalf("currStruct[either][value].GetCoins() not successeded")
	}
	if eitherVal.Cmp(big.NewInt(100000000)) != 0 {
		t.Fatalf("currStruct[either][value].GetCoins() != 1000000000, got %v", eitherVal.String())
	}

	anotherEither, ok := currStruct.GetField("anotherEither")
	if !ok {
		t.Fatalf("currStruct[anotherEither] not found")
	}
	anotherEitherAlias, ok := anotherEither.GetAlias()
	if !ok {
		t.Fatalf("currStruct[anotherEither].GetAlias() not found")
	}
	anotherEitherAliasAlias, ok := anotherEitherAlias.GetAlias()
	if !ok {
		t.Fatalf("currStruct[anotherEither].GetAlias().GetAlias() not found")
	}
	anotherEitherUnion, ok := anotherEitherAliasAlias.GetUnion()
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
	anotherEitherVGeneric, ok := anotherEitherV.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[anotherEither][value].GetGeneric() not found")
	}
	anotherEitherVal, ok := anotherEitherVGeneric.GetTensor()
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
	anotherEitherValCoins, ok := anotherEitherVal[1].GetCoins()
	if !ok {
		t.Fatalf("currStruct[anotherEither][value][0].GetCoins() not successeded")
	}
	if anotherEitherValCoins.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Fatalf("currStruct[anotherEither][value][0].GetCoins() != 1000000000, got %v", anotherEitherValCoins.String())
	}

	doubler, ok := currStruct.GetField("doubler")
	if !ok {
		t.Fatalf("currStruct[doubler] not found")
	}
	doublerRef, ok := doubler.GetRefValue()
	if !ok {
		t.Fatalf("currStruct[doubler].GetRefValue() not successeded")
	}
	doublerRefAlias, ok := doublerRef.GetAlias()
	if !ok {
		t.Fatalf("currStruct[doubler].GetAlias() not found")
	}
	doublerTensor, ok := doublerRefAlias.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler].GetTensor() not successeded")
	}

	doublerGenric0, ok := doublerTensor[0].GetGeneric()
	if !ok {
		t.Fatalf("currStruct[doubler][0].GetGeneric() not successeded")
	}
	doublerTensor0, ok := doublerGenric0.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler][0].GetTensor() not successeded")
	}
	doublerTensor0Coins, ok := doublerTensor0[0].GetCoins()
	if !ok {
		t.Fatalf("currStruct[doubler][0][0].GetCoins() not successeded")
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

	doublerUnion1, ok := doublerTensor[1].GetGeneric()
	if !ok {
		t.Fatalf("currStruct[doubler][1].GetGeneric() not successeded")
	}
	doublerTensor1, ok := doublerUnion1.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler][1] not successeded")
	}
	doublerTensor1Coins, ok := doublerTensor1[0].GetCoins()
	if !ok {
		t.Fatalf("currStruct[doubler][1][0].GetCoins() not successeded")
	}
	if doublerTensor1Coins.Cmp(big.NewInt(100_000_000)) != 0 {
		t.Fatalf("currStruct[doubler][1][0].GetCoins() != 100000000, got %v", doublerTensor1Coins.String())
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
	myValAlias, ok := myVal.GetAlias()
	if !ok {
		t.Fatalf("currStruct[myVal].GetAlias() not found")
	}
	myValGeneric, ok := myValAlias.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[myVal].GetGeneric() not found")
	}
	myValVal, ok := myValGeneric.GetSmallUInt()
	if !ok {
		t.Fatalf("currStruct[myVal].GetSmallUInt() not successed")
	}
	if myValVal != 16 {
		t.Fatalf("currStruct[myVal] != 16, got %v", myValVal)
	}

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotGenericsFromAlias(t *testing.T) {
	jsonInputFilename := "a_lot_generics_from_alias"
	inputFilename := "testdata/generics.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "AliasRef",
		AliasRef: &tolkParser.AliasRef{
			AliasName: "GoodNamingForMsg",
			TypeArgs: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 16,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	currAlias, ok := v.GetAlias()
	if !ok {
		t.Fatalf("v.GetAlias() not found")
	}
	currStruct, ok := currAlias.GetStruct()
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
	// todo: remove GetGeneric because its not convenient
	payloadGeneric, ok := payload.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[payload].GetGeneric() not found")
	}
	payloadVal, ok := payloadGeneric.GetSmallUInt()
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
	eitherAlias, ok := either.GetAlias()
	if !ok {
		t.Fatalf("currStruct[either].GetAlias() not found")
	}
	eitherUnion, ok := eitherAlias.GetUnion()
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
	eitherValGeneric, ok := eitherV.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[either][value].GetGeneric() not found")
	}
	eitherVal, ok := eitherValGeneric.GetCoins()
	if !ok {
		t.Fatalf("currStruct[either][value].GetCoins() not successeded")
	}
	if eitherVal.Cmp(big.NewInt(100000000)) != 0 {
		t.Fatalf("currStruct[either][value].GetCoins() != 1000000000, got %v", eitherVal.String())
	}

	anotherEither, ok := currStruct.GetField("anotherEither")
	if !ok {
		t.Fatalf("currStruct[anotherEither] not found")
	}
	anotherEitherAlias, ok := anotherEither.GetAlias()
	if !ok {
		t.Fatalf("currStruct[anotherEither].GetAlias() not found")
	}
	anotherEitherAliasAlias, ok := anotherEitherAlias.GetAlias()
	if !ok {
		t.Fatalf("currStruct[anotherEither].GetAlias().GetAlias() not found")
	}
	anotherEitherUnion, ok := anotherEitherAliasAlias.GetUnion()
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
	anotherEitherVGeneric, ok := anotherEitherV.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[anotherEither][value].GetGeneric() not found")
	}
	anotherEitherVal, ok := anotherEitherVGeneric.GetTensor()
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
	anotherEitherValCoins, ok := anotherEitherVal[1].GetCoins()
	if !ok {
		t.Fatalf("currStruct[anotherEither][value][0].GetCoins() not successeded")
	}
	if anotherEitherValCoins.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Fatalf("currStruct[anotherEither][value][0].GetCoins() != 1000000000, got %v", anotherEitherValCoins.String())
	}

	doubler, ok := currStruct.GetField("doubler")
	if !ok {
		t.Fatalf("currStruct[doubler] not found")
	}
	doublerRef, ok := doubler.GetRefValue()
	if !ok {
		t.Fatalf("currStruct[doubler].GetRefValue() not successeded")
	}
	doublerRefAlias, ok := doublerRef.GetAlias()
	if !ok {
		t.Fatalf("currStruct[doubler].GetAlias() not found")
	}
	doublerTensor, ok := doublerRefAlias.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler].GetTensor() not successeded")
	}

	doublerGenric0, ok := doublerTensor[0].GetGeneric()
	if !ok {
		t.Fatalf("currStruct[doubler][0].GetGeneric() not successeded")
	}
	doublerTensor0, ok := doublerGenric0.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler][0].GetTensor() not successeded")
	}
	doublerTensor0Coins, ok := doublerTensor0[0].GetCoins()
	if !ok {
		t.Fatalf("currStruct[doubler][0][0].GetCoins() not successeded")
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

	doublerUnion1, ok := doublerTensor[1].GetGeneric()
	if !ok {
		t.Fatalf("currStruct[doubler][1].GetGeneric() not successeded")
	}
	doublerTensor1, ok := doublerUnion1.GetTensor()
	if !ok {
		t.Fatalf("currStruct[doubler][1] not successeded")
	}
	doublerTensor1Coins, ok := doublerTensor1[0].GetCoins()
	if !ok {
		t.Fatalf("currStruct[doubler][1][0].GetCoins() not successeded")
	}
	if doublerTensor1Coins.Cmp(big.NewInt(100_000_000)) != 0 {
		t.Fatalf("currStruct[doubler][1][0].GetCoins() != 100000000, got %v", doublerTensor1Coins.String())
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
	myValAlias, ok := myVal.GetAlias()
	if !ok {
		t.Fatalf("currStruct[myVal].GetAlias() not found")
	}
	myValGeneric, ok := myValAlias.GetGeneric()
	if !ok {
		t.Fatalf("currStruct[myVal].GetGeneric() not found")
	}
	myValVal, ok := myValGeneric.GetSmallUInt()
	if !ok {
		t.Fatalf("currStruct[myVal].GetSmallUInt() not successed")
	}
	if myValVal != 16 {
		t.Fatalf("currStruct[myVal] != 16, got %v", myValVal)
	}

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalStructWithDefaultValues(t *testing.T) {
	jsonInputFilename := "a_lot_generics_with_default_values"
	inputFilename := "testdata/default_values.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "DefaultTest",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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
	num2Alias, ok := num2.Val.GetAlias()
	if !ok {
		t.Fatalf("currStruct[num2].GetAlias() not found")
	}
	num2Val, ok := num2Alias.GetSmallInt()
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
	slice3Val, ok := slice3.Val.GetRemaining()
	if !ok {
		t.Fatalf("currStruct[slice3].GetRemaining() not successeded")
	}
	hs, err := slice3Val.HashString()
	if err != nil {
		t.Fatal(err)
	}
	if hs != "55e960f1409af0d7670e382c61276a559fa9330185984d91faffebf32d5fa383" {
		t.Fatalf("currStruct[slice3].GetRemaining().Hash() != 55e960f1409af0d7670e382c61276a559fa9330185984d91faffebf32d5fa383, got %v", hs)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotNumbers(t *testing.T) {
	jsonInputFilename := "a_lot_numbers"
	inputFilename := "testdata/numbers.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "Numbers",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

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
	val3, ok := num3.GetBigUInt()
	if !ok {
		t.Fatalf("num3.GetBigUInt() not successeded")
	}
	if val3.Cmp(big.NewInt(241)) != 0 {
		t.Fatalf("num3 != 241, got %v", val3.String())
	}

	num4, ok := currStruct.GetField("num4")
	if !ok {
		t.Fatalf("num4 not found")
	}
	val4, ok := num4.GetVarUInt()
	if !ok {
		t.Fatalf("num4.GetVarUInt() not successeded")
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
	val9, ok := num9.GetVarInt()
	if !ok {
		t.Fatalf("num9.GetVarInt() not successeded")
	}
	if val9.Cmp(big.NewInt(2342)) != 0 {
		t.Fatalf("num9 != 2342, got %s", val9.String())
	}

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_UnmarshalALotRandomFields(t *testing.T) {
	jsonInputFilename := "a_lot_random_fields"
	inputFilename := "testdata/random_fields.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "RandomFields",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	addr := ton.MustParseAccountID("UQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqEBI")

	currStruct, ok := v.GetStruct()
	if !ok {
		t.Fatalf("struct not found")
	}
	pref, ok := currStruct.GetPrefix()
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
	amountVal, ok := amount.GetCoins()
	if !ok {
		t.Fatalf("amount.GetCoins() not successeded")
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
	val2, ok := optVal2.Val.GetCoins()
	if !ok {
		t.Fatalf("intVector[1].GetOptionalValue().GetCoins() not successeded")
	}
	if val2.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("intVector[1].GetOptionalValue().GetCoins() != 1000000000, got %v", val1)
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
	myIntAlias, ok := myInt.GetAlias()
	if !ok {
		t.Fatalf("my_int.GetAlias() not successeded")
	}
	myIntVal, ok := myIntAlias.GetSmallInt()
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actual, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(pathPrefix+".output.json", actual, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalSmallInt(t *testing.T) {
	jsonInputFilename := "small_int"
	ty := tolkParser.Ty{
		SumType: "IntN",
		IntN: &tolkParser.IntN{
			N: 24,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010005000006ff76c41616db06")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalSmallUInt(t *testing.T) {
	jsonInputFilename := "small_uint"
	ty := tolkParser.Ty{
		SumType: "UintN",
		UintN: &tolkParser.UintN{
			N: 53,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000d00000000001d34e435eafd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBigInt(t *testing.T) {
	jsonInputFilename := "big_int"
	ty := tolkParser.Ty{
		SumType: "IntN",
		IntN: &tolkParser.IntN{
			N: 183,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBigUInt(t *testing.T) {
	jsonInputFilename := "big_uint"
	ty := tolkParser.Ty{
		SumType: "UintN",
		UintN: &tolkParser.UintN{
			N: 257,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002300004100000000000000000000000000000000000000000000000000009fc4212a38ba40b11cce12")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalVarInt(t *testing.T) {
	jsonInputFilename := "var_int"
	ty := tolkParser.Ty{
		SumType: "VarIntN",
		VarIntN: &tolkParser.VarIntN{
			N: 16,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000730c98588449b6923")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalVarUInt(t *testing.T) {
	jsonInputFilename := "var_uint"
	ty := tolkParser.Ty{
		SumType: "VarUintN",
		VarUintN: &tolkParser.VarUintN{
			N: 32,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000800000b28119ab36b44d3a86c0f")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBits(t *testing.T) {
	jsonInputFilename := "bits"
	ty := tolkParser.Ty{
		SumType: "BitsN",
		BitsN: &tolkParser.BitsN{
			N: 24,
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000500000631323318854035")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalCoins(t *testing.T) {
	jsonInputFilename := "coins"
	ty := tolkParser.Ty{
		SumType: "Coins",
		Coins:   &tolkParser.Coins{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010007000009436ec6e0189ebbd7f4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBool(t *testing.T) {
	jsonInputFilename := "bool"
	ty := tolkParser.Ty{
		SumType: "Bool",
		Bool:    &tolkParser.Bool{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000300000140f6d24034")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalCell(t *testing.T) {
	jsonInputFilename := "cell"
	ty := tolkParser.Ty{
		SumType: "Cell",
		Cell:    &tolkParser.Cell{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101020100090001000100080000007ba52a3292")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalRemaining(t *testing.T) {
	jsonInputFilename := "remaining"
	ty := tolkParser.Ty{
		SumType:   "Remaining",
		Remaining: &tolkParser.Remaining{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAddress(t *testing.T) {
	jsonInputFilename := "internal_address"
	ty := tolkParser.Ty{
		SumType: "Address",
		Address: &tolkParser.Address{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalNotExitsOptionalAddress(t *testing.T) {
	jsonInputFilename := "not_exists_optional_address"
	ty := tolkParser.Ty{
		SumType:    "AddressOpt",
		AddressOpt: &tolkParser.AddressOpt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100030000012094418655")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalExistsOptionalAddress(t *testing.T) {
	jsonInputFilename := "exists_optional_address"
	ty := tolkParser.Ty{
		SumType:    "AddressOpt",
		AddressOpt: &tolkParser.AddressOpt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalExternalAddress(t *testing.T) {
	jsonInputFilename := "external_address"
	ty := tolkParser.Ty{
		SumType:    "AddressExt",
		AddressExt: &tolkParser.AddressExt{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000742082850fcbd94fd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAnyNoneAddress(t *testing.T) {
	jsonInputFilename := "any_none_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100030000012094418655")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAnyInternalAddress(t *testing.T) {
	jsonInputFilename := "any_internal_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAnyExternalAddress(t *testing.T) {
	jsonInputFilename := "any_external_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000600000742082850fcbd94fd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAnyVarAddress(t *testing.T) {
	jsonInputFilename := "any_var_address"
	ty := tolkParser.Ty{
		SumType:    "AddressAny",
		AddressAny: &tolkParser.AddressAny{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalNotExistsNullable(t *testing.T) {
	jsonInputFilename := "not_exists_nullable"
	ty := tolkParser.Ty{
		SumType: "Nullable",
		Nullable: &tolkParser.Nullable{
			Inner: tolkParser.Ty{
				SumType:   "Remaining",
				Remaining: &tolkParser.Remaining{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000300000140f6d24034")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalExistsNullable(t *testing.T) {
	jsonInputFilename := "exists_nullable"
	ty := tolkParser.Ty{
		SumType: "Nullable",
		Nullable: &tolkParser.Nullable{
			Inner: tolkParser.Ty{
				SumType: "Cell",
				Cell:    &tolkParser.Cell{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000b000101c001000900000c0ae007880db9")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalRef(t *testing.T) {
	jsonInputFilename := "ref"
	ty := tolkParser.Ty{
		SumType: "CellOf",
		CellOf: &tolkParser.CellOf{
			Inner: tolkParser.Ty{
				SumType: "IntN",
				IntN: &tolkParser.IntN{
					N: 65,
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000e000100010011000000000009689e40e150b4c5")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalEmptyTensor(t *testing.T) {
	jsonInputFilename := "empty_tensor"
	ty := tolkParser.Ty{
		SumType: "Tensor",
		Tensor:  &tolkParser.Tensor{},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101010100020000004cacb9cd")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalTensor(t *testing.T) {
	jsonInputFilename := "tensor"
	ty := tolkParser.Ty{
		SumType: "Tensor",
		Tensor: &tolkParser.Tensor{
			Items: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 123,
					},
				},
				{
					SumType: "Bool",
					Bool:    &tolkParser.Bool{},
				},
				{
					SumType: "Coins",
					Coins:   &tolkParser.Coins{},
				},
				{
					SumType: "Tensor",
					Tensor: &tolkParser.Tensor{
						Items: []tolkParser.Ty{
							{
								SumType: "IntN",
								IntN: &tolkParser.IntN{
									N: 23,
								},
							},
							{
								SumType: "Nullable",
								Nullable: &tolkParser.Nullable{
									Inner: tolkParser.Ty{
										SumType: "IntN",
										IntN: &tolkParser.IntN{
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
					VarIntN: &tolkParser.VarIntN{
						N: 32,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalIntKeyMap(t *testing.T) {
	jsonInputFilename := "int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "IntN",
				IntN: &tolkParser.IntN{
					N: 32,
				},
			},
			V: tolkParser.Ty{
				SumType: "Bool",
				Bool:    &tolkParser.Bool{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201000c000101c001000ba00000007bc09a662c32")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalUIntKeyMap(t *testing.T) {
	jsonInputFilename := "uint_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "UintN",
				UintN: &tolkParser.UintN{
					N: 16,
				},
			},
			V: tolkParser.Ty{
				SumType: "Address",
				Address: &tolkParser.Address{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
	if err != nil {
		t.Fatal(err)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBigIntKeyMap(t *testing.T) {
	jsonInputFilename := "big_int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "UintN",
				UintN: &tolkParser.UintN{
					N: 78,
				},
			},
			V: tolkParser.Ty{
				SumType: "Cell",
				Cell:    &tolkParser.Cell{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalBitsKeyMap(t *testing.T) {
	jsonInputFilename := "bits_int_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "BitsN",
				BitsN: &tolkParser.BitsN{
					N: 16,
				},
			},
			V: tolkParser.Ty{
				SumType: "Map",
				Map: &tolkParser.Map{
					K: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 64,
						},
					},
					V: tolkParser.Ty{
						SumType: "Tensor",
						Tensor: &tolkParser.Tensor{
							Items: []tolkParser.Ty{
								{
									SumType: "Address",
									Address: &tolkParser.Address{},
								},
								{
									SumType: "Coins",
									Coins:   &tolkParser.Coins{},
								},
							},
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalAddressKeyMap(t *testing.T) {
	jsonInputFilename := "address_key_map"
	ty := tolkParser.Ty{
		SumType: "Map",
		Map: &tolkParser.Map{
			K: tolkParser.Ty{
				SumType: "Address",
				Address: &tolkParser.Address{},
			},
			V: tolkParser.Ty{
				SumType: "Coins",
				Coins:   &tolkParser.Coins{},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalUnionWithDecPrefix(t *testing.T) {
	jsonInputFilename := "union_with_dec_prefix"
	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr:        "0",
					PrefixLen:        1,
					PrefixEatInPlace: true,
					VariantTy: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 16,
						},
					},
				},
				{
					PrefixStr:        "1",
					PrefixLen:        1,
					PrefixEatInPlace: true,
					VariantTy: tolkParser.Ty{
						SumType: "IntN",
						IntN: &tolkParser.IntN{
							N: 128,
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	newCell, err := Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalUnionWithBinPrefix(t *testing.T) {
	jsonInputFilename := "union_with_bin_prefix"
	inputFilename := "testdata/bin_union.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr: "0b001",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "AddressWithPrefix",
						},
					},
				},
				{
					PrefixStr: "0b011",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "MapWithPrefix",
						},
					},
				},
				{
					PrefixStr: "0b111",
					PrefixLen: 3,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "CellWithPrefix",
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalUnionWithHexPrefix(t *testing.T) {
	jsonInputFilename := "union_with_hex_prefix"
	inputFilename := "testdata/hex_union.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "Union",
		Union: &tolkParser.Union{
			Variants: []tolkParser.UnionVariant{
				{
					PrefixStr: "0x12345678",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt66WithPrefix",
						},
					},
				},
				{
					PrefixStr: "0xdeadbeef",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt33WithPrefix",
						},
					},
				},
				{
					PrefixStr: "0x89abcdef",
					PrefixLen: 32,
					VariantTy: tolkParser.Ty{
						SumType: "StructRef",
						StructRef: &tolkParser.StructRef{
							StructName: "UInt4WithPrefix",
						},
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101000b000011deadbeef00000000c0d75977b9")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotRefsFromAlias(t *testing.T) {
	jsonInputFilename := "a_lot_refs_from_alias"
	inputFilename := "testdata/refs.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "AliasRef",
		AliasRef: &tolkParser.AliasRef{
			AliasName: "GoodNamingForMsg",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
	if err != nil {
		t.Fatal(err)
	}

	cb := currCell[0].Refs()[0].RawBitString()
	fmt.Println(cb.BinaryString())

	nb := newCell.Refs()[0].RawBitString()
	fmt.Println(nb.BinaryString())

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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotRefsFromStruct(t *testing.T) {
	jsonInputFilename := "a_lot_refs_from_struct"
	inputFilename := "testdata/refs.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "ManyRefsMsg",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotGenericsFromStruct(t *testing.T) {
	jsonInputFilename := "a_lot_generics_from_struct"
	inputFilename := "testdata/generics.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "ManyRefsMsg",
			TypeArgs: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 16,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotGenericsFromAlias(t *testing.T) {
	jsonInputFilename := "a_lot_generics_from_alias"
	inputFilename := "testdata/generics.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "AliasRef",
		AliasRef: &tolkParser.AliasRef{
			AliasName: "GoodNamingForMsg",
			TypeArgs: []tolkParser.Ty{
				{
					SumType: "UintN",
					UintN: &tolkParser.UintN{
						N: 16,
					},
				},
			},
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalStructWithDefaultValues(t *testing.T) {
	jsonInputFilename := "a_lot_generics_with_default_values"
	inputFilename := "testdata/default_values.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "DefaultTest",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotNumbers(t *testing.T) {
	jsonInputFilename := "a_lot_numbers"
	inputFilename := "testdata/numbers.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "Numbers",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}

func TestRuntime_MarshalALotRandomFields(t *testing.T) {
	jsonInputFilename := "a_lot_random_fields"
	inputFilename := "testdata/random_fields.json"
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		t.Fatal(err)
	}

	var abi tolkParser.ABI
	err = json.Unmarshal(data, &abi)
	if err != nil {
		t.Fatal(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "RandomFields",
		},
	}

	currCell, err := boc.DeserializeBocHex("b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder()
	decoder.WithABI(abi)
	v, err := decoder.Unmarshal(currCell[0], ty)
	if err != nil {
		t.Fatal(err)
	}

	encoder := NewEncoder()
	encoder.WithABI(abi)
	newCell, err := encoder.Marshal(v, ty)
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

	pathPrefix := jsonFilesPath + jsonInputFilename
	actualJson, err := os.ReadFile(pathPrefix + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonV Value
	if err := json.Unmarshal(actualJson, &jsonV); err != nil {
		t.Fatal(err)
	}
	if !v.Equal(jsonV) {
		t.Errorf("%s got different results", pathPrefix)
	}
}
