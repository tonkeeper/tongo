package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/tonkeeper/tongo/tlb/parser"
)

var bitsSizes = []int{80, 96, 128, 256, 264, 320, 352, 512}
var intSizes = []int{128, 256, 257}
var uintSizes = []int{128, 160, 220, 256}

func main() {
	outputDir := flag.String("output", "tlb", "directory to write generated files")
	flag.Parse()

	src := `// Code generated - DO NOT EDIT.

package tlb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)
` + parser.GenerateVarUintTypes(32) +
		parser.GenerateConstantInts(64) +
		parser.GenerateConstantBigInts(intSizes) +
		parser.GenerateConstantBigUints(uintSizes) +
		parser.GenerateBitsTypes(bitsSizes)

	formatted, err := format.Source([]byte(src))
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(*outputDir, "integers.go")
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		panic(err)
	}
	fmt.Println(outPath)
}
