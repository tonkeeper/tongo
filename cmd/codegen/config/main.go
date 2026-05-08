package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tonkeeper/tongo/tlb/parser"
)

func main() {
	input := flag.String("input", "tlb/config.tlb", "path to config.tlb")
	outputDir := flag.String("output", "tlb", "directory to write generated files")
	flag.Parse()

	content, err := os.ReadFile(*input)
	if err != nil {
		panic(err)
	}
	tlb, err := parser.Parse(string(content))
	if err != nil {
		panic(err)
	}
	// Filter declarations already defined elsewhere in the tlb package.
	var filtered []parser.CombinatorDeclaration
	for _, decl := range tlb.Declarations {
		if decl.Combinator.Name == "GlobalVersion" {
			continue
		}
		filtered = append(filtered, decl)
	}

	g := parser.NewGenerator(
		parser.WithTlbPackage(""),
	)
	s, err := g.GenerateGolangTypes(filtered, "", false)
	if err != nil {
		panic(err)
	}
	file := parser.File{
		Name:    filepath.Join(*outputDir, "config.go"),
		Package: "tlb",
		Imports: []string{"encoding/json", "fmt"},
		Code:    s,
	}
	if err := file.Save(); err != nil {
		panic(err)
	}
	fmt.Println(file.Name)
}
