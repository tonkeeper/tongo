package main

import (
	"encoding/hex"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

func main() {
	// Deserialize boc from hex or use boc.DeserializeBocBase64("base64_encoded_boc").
	b, err := hex.DecodeString("B5EE9C72C1020301000109000000005000AD019AC8A9A995CAED6DDA5DD1FC6E4934C0D209DA376C027F6FC48264B257A31E38D2A2EAB37B7821221A748419EFFD17CCD508FE27665252ECB0CB82A19CE9CA360729A9A31762C2ECD700000705010101B4620060989BF00560E387F3F3ADAE93CA236D3C67CB9BF75DD826E471703BFB1D045F180C35000000000000000000000000000000000000F09F928E20436F6E67726174756C6174696F6E732120596F75722077616C6C657420630200B4616E206265636F6D6520612076616C696461746F722E2053656E64203130206F72206D6F726520544F4E7320746F20746869732077616C6C657420616E6420676574206261636B20646F75626C652074686520616D6F756E742E9E83C3C2")
	if err != nil {
		panic(err)
	}
	bagOfCells, err := boc.DeserializeBoc(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("First root cell from boc:\n %v\n", bagOfCells[0].ToString())

	// Build new cell
	newCell := boc.NewCell()
	fmt.Printf("Empty cell: %v\n", newCell.ToString())
	err = newCell.WriteUint(533, 15) // write 533 as 15 bit uint
	if err != nil {
		panic(err)
	}
	err = newCell.WriteBit(true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated cell: %v\n", newCell.ToString())

	// Read something from newCell
	uint15Picked, err := newCell.PickUint(15) // read without read cursor moving
	if err != nil {
		panic(err)
	}
	fmt.Printf("uint15 read from the cell: %v\n", uint15Picked)
	uint15, err := newCell.ReadUint(15) // read with cursor moving
	if err != nil {
		panic(err)
	}
	fmt.Printf("uint15 read from the cell: %v\n", uint15)
}
