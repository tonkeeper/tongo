package parser

import (
	"fmt"
	"testing"
)

var intSizes = []int{96, 264, 320, 352}
var bitsSizes = []int{96, 264, 320, 352}

func TestGenerateConstantBigInts(t *testing.T) {
	s := GenerateConstantBigInts(intSizes)
	fmt.Printf(s)
}

func TestGenerateVarUintTypes(t *testing.T) {
	GenerateVarUintTypes(32)
}

func TestGenerateConstantInts(t *testing.T) {
	GenerateConstantInts(64)
}

func TestGenerateBitsTypes(t *testing.T) {
	GenerateBitsTypes(bitsSizes)
}
