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
	hash, err := cells[0].Hash()
	fmt.Println(hex.EncodeToString(hash))

	//fmt.Println(hex.EncodeToString(cells[0].ToBOC()))
	//fmt.Sprintf("%x", cells[0].Hash())
	//fmt.Println(cells[0].Hash())
	//cells[0].Bits.Print()
	//parse := cells[0].BeginParse()

	//fmt.Println(parse.ReadBigUint(8))
}

func TestKek2(t *testing.T) {
	s := "b5ee9c72c10101010003000000028058c23e9f"

	data, _ := hex.DecodeString(s)

	cells, err := DeserializeBoc(data)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	//fmt.Println(hex.EncodeToString(cells[0].Hash()))

	//fmt.Println("danila:", hex.EncodeToString(cells[0].ToBOC()))

	serialized, err := SerializeBoc(cells[0], true, true, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("my:", hex.EncodeToString(serialized))

	testCell := NewCell()
	err = testCell.Bits.WriteUint(128, 8)
	if err != nil {
		t.Fatal(err)
	}
	testCell2 := NewCell()
	err = testCell2.Bits.WriteInt(-55, 32)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testCell.AddReference(testCell2)
	if err != nil {
		t.Fatal(err)
	}
	serialized2, err := SerializeBoc(testCell, true, true, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("kek:", hex.EncodeToString(serialized2))

	//fmt.Sprintf("%x", cells[0].Hash())
	//fmt.Println(cells[0].Hash())
	//cells[0].Bits.Print()
	//parse := cells[0].BeginParse()

	//fmt.Println(parse.ReadBigUint(8))
}
