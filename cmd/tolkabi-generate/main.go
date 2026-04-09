package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	talkCodegen "github.com/tonkeeper/tongo/abi-tolk/codegen"
	"github.com/tonkeeper/tongo/tolk/parser"
)

func main() {
	schemasDir := "abi-tolk/schemas"
	outputDir := "abi-tolk/abiGenerated"

	if len(os.Args) > 1 {
		schemasDir = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	err := filepath.WalkDir(schemasDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		var abi parser.ABI
		if err := json.Unmarshal(data, &abi); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		gen := talkCodegen.NewTolkGolangGenerator(abi)
		code, err := gen.GenerateGocode()
		if err != nil {
			return fmt.Errorf("codegen %s: %w", path, err)
		}
		if code == "" {
			return nil
		}

		code = `package abiGenerated

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

` + code + "\n\n"

		rel, _ := filepath.Rel(schemasDir, path)
		outPath := filepath.Join(outputDir, strings.TrimSuffix(rel, ".json")+".go")

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
		}

		if err := os.WriteFile(outPath, []byte(code), 0644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}

		fmt.Printf("generated %s\n", outPath)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
