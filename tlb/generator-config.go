//go:build ignore

package main

import (
	"os"

	"github.com/tonkeeper/tongo/tlb/parser"
)

func main() {
	content, err := os.ReadFile("config.tlb")
	if err != nil {
		panic(err)
	}
	tlb, err := parser.Parse(string(content))
	if err != nil {
		panic(err)
	}
	g := parser.NewGenerator(nil, "")
	s, err := g.GenerateGolangTypes(tlb.Declarations, "", false)
	if err != nil {
		panic(err)
	}
	file := parser.File{
		Name:    "config.go",
		Package: "tlb",
		Imports: []string{"encoding/json", "fmt"},
		Code:    s,
	}
	if err := file.Save(); err != nil {
		panic(err)
	}
}
