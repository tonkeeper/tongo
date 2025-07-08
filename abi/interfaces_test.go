package abi

import (
	"encoding/json"
	"testing"
)

func TestContractInterface_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		contract ContractInterface
		expected string
	}{
		{
			name:     "JettonMaster",
			contract: JettonMaster,
			expected: `"jetton_master"`,
		},
		{
			name:     "JettonWallet",
			contract: JettonWallet,
			expected: `"jetton_wallet"`,
		},
		{
			name:     "NftItem",
			contract: NftItem,
			expected: `"nft_item"`,
		},
		{
			name:     "WalletV4R1",
			contract: WalletV4R1,
			expected: `"wallet_v4r1"`,
		},
		{
			name:     "Unknown",
			contract: ContractInterface(999999),
			expected: `"unknown"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.contract)
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(data), tt.expected)
			}
		})
	}
}

func TestContractInterface_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected ContractInterface
	}{
		{
			name:     "JettonMaster",
			json:     `"jetton_master"`,
			expected: JettonMaster,
		},
		{
			name:     "JettonWallet",
			json:     `"jetton_wallet"`,
			expected: JettonWallet,
		},
		{
			name:     "NftItem",
			json:     `"nft_item"`,
			expected: NftItem,
		},
		{
			name:     "WalletV4R1",
			json:     `"wallet_v4r1"`,
			expected: WalletV4R1,
		},
		{
			name:     "Unknown",
			json:     `"unknown_interface"`,
			expected: IUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var contract ContractInterface
			err := json.Unmarshal([]byte(tt.json), &contract)
			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if contract != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", contract, tt.expected)
			}
		})
	}
}

func TestContractInterface_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		contract ContractInterface
	}{
		{
			name:     "JettonMaster",
			contract: JettonMaster,
		},
		{
			name:     "JettonWallet",
			contract: JettonWallet,
		},
		{
			name:     "NftItem",
			contract: NftItem,
		},
		{
			name:     "WalletV4R1",
			contract: WalletV4R1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.contract)
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			// Unmarshal from JSON
			var contract ContractInterface
			err = json.Unmarshal(data, &contract)
			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			// Check that we got back the same value
			if contract != tt.contract {
				t.Errorf("Round trip failed: got %v, want %v", contract, tt.contract)
			}
		})
	}
}

func TestContractInterface_JSON(t *testing.T) {
	contract := JettonMaster
	data, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	expected := `"jetton_master"`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	var result ContractInterface
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result != contract {
		t.Errorf("Expected %v, got %v", contract, result)
	}

	unknownJSON := `"unknown_interface"`
	var unknown ContractInterface
	err = json.Unmarshal([]byte(unknownJSON), &unknown)
	if err != nil {
		t.Fatalf("Failed to unmarshal unknown: %v", err)
	}

	if unknown != IUnknown {
		t.Errorf("Expected IUnknown, got %v", unknown)
	}
}
