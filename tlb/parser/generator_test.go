package parser

import (
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"testing"
)

//go:embed testdata/_config.tlb
var config string

//go:embed testdata/_big.tlb
var big string

func TestGenerateGolangTypes(t *testing.T) {

	tests := []struct {
		source           string
		expectedFilename string
	}{
		{
			source:           config,
			expectedFilename: "testdata/config.go.out",
		},
		{
			source:           big,
			expectedFilename: "testdata/big.go.out",
		},
	}

	for _, testCase := range tests {
		parsed, err := Parse(testCase.source)
		if err != nil {
			t.Fatalf("failed to parse tlb: %s", err)
		}
		g := NewGenerator()
		s, err := g.GenerateGolangTypes(parsed.Declarations, "", false)
		if err != nil {
			t.Fatalf("failed to generate golang types: %s", err)
		}
		sourceCode, err := format.Source([]byte(s))
		if err != nil {
			t.Fatalf("failed to format source code: %s", err)
		}
		filename := fmt.Sprintf("%vput.out", testCase.expectedFilename)
		if err := os.WriteFile(filename, sourceCode, 0644); err != nil {
			t.Fatalf("failed to write output file %s: %s", filename, err)
		}
		content, err := os.ReadFile(testCase.expectedFilename)
		if err != nil {
			t.Fatalf("failed to read expected file %s: %s", testCase.expectedFilename, err)
		}
		if string(content) != string(sourceCode) {
			t.Fatalf("expected file %s does not match generated code", testCase.expectedFilename)
		}
	}
}
