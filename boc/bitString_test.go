package boc

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestOn(t *testing.T) {

	fmt.Println(math.Ceil(float64(1023) / float64(8)))
	var str = NewBitString(2)
	//
	str.Print()
	//str.WriteUint(big.NewInt(255), 8)
	str.WriteInt(big.NewInt(-8), 8)
	str.Print()

	var reader = NewBitStringReader(&str)

	var num = reader.ReadInt(8)
	fmt.Println(num)

	//var n = big.NewInt(1)
	//fmt.Println("77777", n.BitLen())
	//
	//for i := n.BitLen() - 1; i != 0; i-- {
	//	fmt.Println(n.Bit(i))
	//}

	//str.WriteInt(big.NewInt(127), 8)
	//str.On(0)
	//
	//str.On(8)
	//str.Print()
	//
	//str.Off(8)
	//str.Print()
	//
	//str.Toggle(8)
	//str.Print()
	//str.Toggle(8)
	//str.Print()

}
