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

//go:embed abi/swap_coffee.json
var swapCoffeeAbiData []byte

func main() {
	var jettonWalletABI parser.ABI
	err := json.Unmarshal(jettonWalletAbiData, &jettonWalletABI)
	if err != nil {
		panic(err)
	}

	var swapCoffeeABI parser.ABI
	err = json.Unmarshal(swapCoffeeAbiData, &swapCoffeeABI)
	if err != nil {
		panic(err)
	}

	ty := parser.NewStructType("Transfer")
	b, err := hex.DecodeString("b5ee9c72010205010001310001ae0f8a7ea5003c0fe80d6334813ef1895801d9d00dc43cb19c9bfa417f874e7be3670825205a6e50d3490a0aec1194fe2e8f003ffffbc78bb7c4bb37210cfa7f64397ff32df53ce36c46934e21a1e315c93137482ed968810101084ee9b1060201ae0f8a7ea5003c0fe80d6334813ef189d801fe7933ccc21619bc6f1800071d12d4e7aae8d68f637fcdb8facdca753b08c859003ffffbc78bb7c4bb37210cfa7f64397ff32df53ce36c46934e21a1e315c93137481823cf41030153c0ffee10c44cf0b64eb9dedfb010c7b24a508140d86f6e07789f9383bb9b2141ea9c08ea502d0bbc060404008d69a178e5801ffffde3c5dbe25d9b90867d3fb21cbff996fa9e71b62349a710d0f18ae4989bb00021706f0b21587534d776380d79ce6fcbef538e7aacc3f949f91149c80db932bd")
	if err != nil {
		panic(err)
	}
	cell, err := boc.DeserializeBoc(b)
	if err != nil {
		panic(err)
	}

	decoder := tolk.NewDecoder()
	err = decoder.WithABIs(jettonWalletABI, swapCoffeeABI)
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
	inputFilename := "examples/tolk/payload-unmarshal/output.json"
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
