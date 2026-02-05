package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk"
	"github.com/tonkeeper/tongo/tolk/parser"
)
import _ "embed"

//go:embed abi.json
var abiData []byte

func main() {
	var abi tolkParser.ABI
	err := json.Unmarshal(abiData, &abi)
	if err != nil {
		panic(err)
	}

	ty := tolkParser.Ty{
		SumType: "StructRef",
		StructRef: &tolkParser.StructRef{
			StructName: "Transfer",
		},
	}
	b, err := hex.DecodeString("b5ee9c72c10101010056000000a75fcc3d140000000000000000800c0674dd00e3a7231084788441cc873e60eb8681f44901cba3a9107c5c322dc4500034a37c6673343b360e10d4e438483b555805a20e5f056742b6a42ba35311994c802625a008a90c976e")
	if err != nil {
		panic(err)
	}
	cell, err := boc.DeserializeBoc(b)
	if err != nil {
		panic(err)
	}

	decoder := tolk.NewDecoder()
	decoder.WithABI(abi)
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

	queryId, ok := tolkStruct.GetField("queryId")
	if !ok {
		panic("transfer.queryId not found")
	}
	queryIdValue, ok := queryId.GetSmallUInt()
	if !ok {
		panic("cannot get transfer.queryId value")
	}

	newOwner, ok := tolkStruct.GetField("newOwner")
	if !ok {
		panic("transfer.newOwner not found")
	}
	newOwnerValue, ok := newOwner.GetAddress()
	if !ok {
		panic("cannot get transfer.newOwner value")
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

	forwardPayload, ok := tolkStruct.GetField("forwardPayload")
	if !ok {
		panic("transfer.forwardPayload not found")
	}
	forwardPayloadValue, ok := forwardPayload.GetRemaining()
	if !ok {
		panic("cannot get transfer.forwardPayload value")
	}

	fmt.Printf("Transfer prefix: 0x%x\n", prefix.Prefix)
	fmt.Printf("Transfer query id: %v\n", queryIdValue)
	fmt.Printf("Transfer new owner: %v\n", newOwnerValue.ToRaw())
	fmt.Printf("Transfer response destination: %v\n", responseDestinationValue.ToRaw())
	fmt.Printf("Transfer is custom payload exists: %v\n", customPayloadValue.IsExists)
	fmt.Printf("Transfer forward amount: %x\n", forwardAmountValue.String())
	fmt.Printf("Transfer forward value: %x\n", forwardPayloadValue.ReadRemainingBits())
}
