package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk"
	"github.com/tonkeeper/tongo/tolk/parser"
)
import _ "embed"

//go:embed abi/jetton_wallet.json
var jettonWalletAbiData []byte

func main() {
	var jettonWalletABI tolkParser.ABI
	err := json.Unmarshal(jettonWalletAbiData, &jettonWalletABI)
	if err != nil {
		panic(err)
	}

	ty := tolkParser.NewStructType("Transfer")
	b, err := hex.DecodeString("b5ee9c720102030100011b0001ae0f8a7ea5bf4e20320df05dd0318537180125c28235ca8d125e676591513d520721b1fe99f7722f4c87723ce7ee0dfb73a3001f8cc4cadf3be6b14e892d7aa31f5be355c1eab776f0e1d61a46644ddc17e68b881908b1010101e16664de2a801244183034d9fd59a236f71ec4271be377399056dda4cc3a5ebf5dc40967df641001f8cc4cadf3be6b14e892d7aa31f5be355c1eab776f0e1d61a46644ddc17e68ba003f198995be77cd629d125af5463eb7c6ab83d56eede1c3ac348cc89bb82fcd17000000007fffffffc0020095446cf7101800fc662656f9df358a74496bd518fadf1aae0f55bbb7870eb0d233226ee0bf345c00000540095e99c8dc6a438526df4961936ff51209f307a28c37c6c78310ce140ab78ab658")
	if err != nil {
		panic(err)
	}
	cell, err := boc.DeserializeBoc(b)
	if err != nil {
		panic(err)
	}

	decoder := tolk.NewDecoder()
	err = decoder.WithABIs(jettonWalletABI)
	if err != nil {
		panic(err)
	}
	res, err := decoder.Unmarshal(cell[0], ty)
	if err != nil {
		panic(err)
	}
	tolkStruct, ok := res.GetStruct()
	if !ok {
		panic("Struct Transfer not found")
	}
	prefix, exists := tolkStruct.GetPrefix()
	if !exists {
		panic("Transfer prefix not found")
	}

	queryId := tolkStruct.MustGetField("queryId")
	queryIdValue := queryId.MustGetSmallUInt()

	newOwner, ok := tolkStruct.GetField("destination")
	if !ok {
		panic("transfer.destination not found")
	}
	newOwnerValue, ok := newOwner.GetAddress()
	if !ok {
		panic("cannot get transfer.destination value")
	}

	responseDestination, ok := tolkStruct.GetField("responseDestination")
	if !ok {
		panic("transfer.responseDestination not found")
	}
	responseDestinationValue, ok := responseDestination.GetAddress()
	if !ok {
		panic("cannot get transfer.responseDestination value")
	}

	customPayload, ok := tolkStruct.GetField("customPayload")
	if !ok {
		panic("transfer.customPayload not found")
	}
	customPayloadValue, ok := customPayload.GetOptionalValue()
	if !ok {
		panic("cannot get transfer.customPayload value")
	}

	forwardAmount, ok := tolkStruct.GetField("forwardAmount")
	if !ok {
		panic("transfer.forwardAmount not found")
	}
	forwardAmountValue, ok := forwardAmount.GetCoins()
	if !ok {
		panic("cannot get transfer.forwardAmount value")
	}

	_, ok = tolkStruct.GetField("forwardPayload")
	if !ok {
		panic("transfer.forwardPayload not found")
	}

	val, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	inputFilename := "examples/tolk/simple-unmarshal/output.json"
	err = os.WriteFile(inputFilename, val, os.ModePerm)
	if err != nil {
		panic(err)
	}

	tolkValue := tolk.Value{}
	if err := json.Unmarshal(val, &tolkValue); err != nil {
		panic(err)
	}

	fmt.Printf("Transfer prefix: 0x%x\n", prefix.Prefix)
	fmt.Printf("Transfer query id: %v\n", queryIdValue)
	fmt.Printf("Transfer new owner: %v\n", newOwnerValue.ToRaw())
	fmt.Printf("Transfer response destination: %v\n", responseDestinationValue.ToRaw())
	fmt.Printf("Transfer is custom payload exists: %v\n", customPayloadValue.IsExists)
	fmt.Printf("Transfer forward amount: %x\n", forwardAmountValue.String())
}
