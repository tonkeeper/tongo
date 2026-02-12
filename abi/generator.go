//go:build ignore

package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonkeeper/tongo/abi/parser"
)

const schemasPath = "schemas/"

func main() {
	var xmlFiles []string
	filepath.Walk(schemasPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".xml") {
			xmlFiles = append(xmlFiles, path)
		}
		return nil
	})

	files, err := parser.Generate(xmlFiles, parser.GenerateOptions{})
	if err != nil {
		panic(err)
	}

	for filename, content := range files {
		if err := os.WriteFile(filename, content, 0644); err != nil {
			panic(err)
		}
	}
}
