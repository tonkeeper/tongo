package parser_tolk

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseABY(t *testing.T) {
	type Case struct {
		name           string
		filenamePrefix string
	}
	for _, c := range []Case{
		{
			name:           "simple abi",
			filenamePrefix: "simple",
		},
		{
			name:           "a lot of wrappers in abi",
			filenamePrefix: "alot-wrappers",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			inputFilename := fmt.Sprintf("testdata/%v.json", c.filenamePrefix)
			outputFilename := fmt.Sprintf("testdata/%v.output.json", c.filenamePrefix)
			expected, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ParseABI(expected)
			if err != nil {
				t.Errorf("failed to unmarshall abi: %v", err)
			}

			bs, err := json.MarshalIndent(parsed, " ", "  ")
			if err != nil {
				t.Fatal(err)
			}
			err = os.WriteFile(outputFilename, bs, 0644)
			if err != nil {
				t.Errorf("failed to write output: %v", err)
			}
			if !bytes.Equal(bs, expected) {
				t.Errorf("output does not match expected output")
			}
		})
	}
}
