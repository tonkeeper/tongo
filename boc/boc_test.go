package boc

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestKek(t *testing.T) {
	s := "b5ee9c72c10101010003000000028058c23e9f"

	data, _ := hex.DecodeString(s)

	cells, err := DeserializeBoc(data)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	//cells[0].Bits.Print()
	parse := cells[0].BeginParse()

	fmt.Println(parse.ReadUint(8))
}
