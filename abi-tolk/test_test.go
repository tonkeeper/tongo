package abi_tolk

import (
	"fmt"
	"testing"

	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
)

func TestRuntimeStonfiMessageDecoding(t *testing.T) {
	data := "b5ee9c720101030100fc0001b00f8a7ea5000000000000000040b84996d80125c28235ca8d125e676591513d520721b1fe99f7722f4c87723ce7ee0dfb73a3001232f21ee11f4703d43f63ecf098d577d2faf4cc38a5129ce1f85438043fd049c81c9c38010101e16664de2a801244183034d9fd59a236f71ec4271be377399056dda4cc3a5ebf5dc40967df641001232f21ee11f4703d43f63ecf098d577d2faf4cc38a5129ce1f85438043fd049e002465e43dc23e8e07a87ec7d9e131aaefa5f5e998714a2539c3f0a870087fa0938000000034ddcf06c002005552312926f738013f41a52de18928b7ee250cb8ecf3ca34aee70728893151e9d8c71c54670e474a00000510"
	boc1, _ := boc.DeserializeBocHex(data)

	decoder, err := GetDecoderWithInterfaces(abi.WalletV4R2, abi.JettonWallet, abi.JettonMaster, abi.StonfiVaultV2, abi.StonfiPoolV2ConstProduct, abi.WalletV5R1)
	if err != nil {
		t.Fatalf("Unable to get decoder: %v", err)
	}

	msg, err := decoder.UnmarshalMessage(boc1[0])
	if err != nil {
		t.Fatalf("Unable to unmarshal message: %v", err)
	}

	fmt.Println(msg)
}
