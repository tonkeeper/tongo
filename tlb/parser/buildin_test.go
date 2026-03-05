package parser

import (
	"fmt"
	"os"
	"testing"
)

var intSizes = []int{96, 264, 320, 352}
var bitsSizes = []int{96, 264, 320, 352}

func TestGenerateConstantBigInts(t *testing.T) {
	if os.Getenv("TEST_CI") == "1" {
		t.SkipNow()
	}
	s := GenerateConstantBigInts(intSizes)
	fmt.Println(s)
}

func TestGenerateVarUintTypes(t *testing.T) {
	if os.Getenv("TEST_CI") == "1" {
		t.SkipNow()
	}
	GenerateVarUintTypes(32)
}

func TestGenerateConstantInts(t *testing.T) {
	if os.Getenv("TEST_CI") == "1" {
		t.SkipNow()
	}
	GenerateConstantInts(64)
}

func TestGenerateBitsTypes(t *testing.T) {
	if os.Getenv("TEST_CI") == "1" {
		t.SkipNow()
	}
	GenerateBitsTypes(bitsSizes)
}
