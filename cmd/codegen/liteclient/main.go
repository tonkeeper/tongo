package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/tonkeeper/tongo/tl/parser"
)

func main() {
	input := flag.String("input", "liteclient/lite_api.tl", "path to lite_api.tl")
	outputDir := flag.String("output", "liteclient", "directory to write generated files")
	flag.Parse()

	scheme, err := os.ReadFile(*input)
	if err != nil {
		panic(err)
	}
	parsed, err := parser.Parse(string(scheme))
	if err != nil {
		panic(err)
	}

	g := parser.NewGenerator(nil, "*Client")

	types, err := g.LoadTypes(parsed.Declarations)
	if err != nil {
		panic(err)
	}
	functions, err := g.LoadFunctions(parsed.Functions)
	if err != nil {
		panic(err)
	}

	src := `// Code generated - DO NOT EDIT.

package liteclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/tonkeeper/tongo/tl"
	"io"
)
` + types + functions

	formatted, err := format.Source([]byte(src))
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(*outputDir, "generated.go")
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		panic(err)
	}
	fmt.Println(outPath)
}
