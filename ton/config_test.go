package ton

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			name:                "mainnet config at 2026 May 8",
			configProofFilename: "testdata/config_proof_65466080.boc",
			expectedFilename:    "testdata/config_65466080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configProof, err := os.ReadFile(tt.configProofFilename)
			require.NoError(t, err, "failed to read config proof from %s", tt.configProofFilename)
			params, err := DecodeConfigParams(configProof)
			require.NoError(t, err, "DecodeConfigParams() failed")
			config, brokenParams, err := ConvertBlockchainConfig(params, true)
			require.NoError(t, err, "ConvertBlockchainConfig() failed")
			bs, err := json.MarshalIndent(config, "", "  ")
			require.NoError(t, err, "MarshalIndent() failed")
			if !reflect.DeepEqual(brokenParams, tt.expectedBrokenParams) {
				t.Fatalf("expected: %v\nactual: %v", tt.expectedBrokenParams, brokenParams)
			}
			outputFilename := fmt.Sprintf("%s.output.json", tt.expectedFilename)
			require.NoError(t, os.WriteFile(outputFilename, bs, 0644))
			expectedFilename := fmt.Sprintf("%s.json", tt.expectedFilename)
			expected, err := os.ReadFile(expectedFilename)
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(bs))
		})
	}
}
