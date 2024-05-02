package ton

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestConvertBlockchainConfig(t *testing.T) {
	tests := []struct {
		name                 string
		configProofFilename  string
		expectedFilename     string
		expectedBrokenParams []int
	}{
		{
			name:                "all good",
			configProofFilename: "testdata/config_proof_33651872.boc",
			expectedFilename:    "testdata/config_33651872",
		},
		{
			name:                 "broken config in testnet",
			configProofFilename:  "testdata/config_proof_4324374.boc",
			expectedFilename:     "testdata/config_4324374",
			expectedBrokenParams: []int{79},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configProof, err := os.ReadFile(tt.configProofFilename)
			if err != nil {
				t.Fatalf("os.ReadFile() failed: %v", err)
			}
			params, err := DecodeConfigParams(configProof)
			if err != nil {
				t.Fatalf("DecodeConfigParams() failed: %v", err)
			}
			config, brokenParams, err := ConvertBlockchainConfig(params, true)
			if err != nil {
				t.Fatalf("ConvertBlockchainConfig() failed: %v", err)
			}
			bs, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent() failed: %v", err)
			}
			if !reflect.DeepEqual(brokenParams, tt.expectedBrokenParams) {
				t.Fatalf("expected: %v\nactual: %v", tt.expectedBrokenParams, brokenParams)
			}
			outputFilename := fmt.Sprintf("%s.output.json", tt.expectedFilename)
			if err := os.WriteFile(outputFilename, bs, 0644); err != nil {
				t.Fatalf("os.WriteFile() failed: %v", err)
			}
			expectedFilename := fmt.Sprintf("%s.json", tt.expectedFilename)
			expected, err := os.ReadFile(expectedFilename)
			if err != nil {
				t.Fatalf("os.ReadFile() failed: %v", err)
			}
			if string(expected) != string(bs) {
				t.Fatalf("expected: %s\nactual: %s", expected, bs)
			}
		})
	}
}
