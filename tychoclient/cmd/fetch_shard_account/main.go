package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tychoclient"
)

// ShardAccountFixture represents a serialized test fixture for a TON shard account
type ShardAccountFixture struct {
	Description   string                 `json:"description"`
	Workchain     int32                  `json:"workchain"`
	Address       string                 `json:"address"`
	WithProof     bool                   `json:"with_proof"`
	McStateInfo   map[string]interface{} `json:"mc_state_info,omitempty"`
	AccountFound  bool                   `json:"account_found"`
	AccountState  string                 `json:"account_state,omitempty"`
	AccountType   string                 `json:"account_type,omitempty"`
	Proof         string                 `json:"proof,omitempty"`
	Balance       string                 `json:"balance,omitempty"`
	LastTransHash string                 `json:"last_trans_hash,omitempty"`
	LastTransLt   uint64                 `json:"last_trans_lt,omitempty"`
	Debug         *DebugInfo             `json:"debug,omitempty"`
}

// DebugInfo contains debugging information about account parsing
type DebugInfo struct {
	AccountStateLength int      `json:"account_state_length"`
	BOCCells           []string `json:"boc_cells,omitempty"`
	BOCError           string   `json:"boc_error,omitempty"`
	ParseError         string   `json:"parse_error,omitempty"`
}

func main() {
	var (
		address   = flag.String("address", "", "Account address (user-friendly or 64 hex chars)")
		withProof = flag.Bool("with-proof", false, "Include proof in the response")
		output    = flag.String("output", "", "Output file path (optional, prints to stdout if not specified)")
		debug     = flag.Bool("debug", false, "Include debug information")
	)
	flag.Parse()

	if *address == "" {
		log.Fatal("Address is required")
	}

	// Parse address (support both user-friendly and hex formats)
	var addressBytes []byte
	var workchain int32
	var err error

	if len(*address) == 64 {
		// Assume raw hex format
		addressBytes, err = hex.DecodeString(*address)
		if err != nil {
			log.Fatalf("Failed to decode hex address: %v", err)
		}
		workchain = 0
	} else {
		// Try to parse as user-friendly address
		accountID, err := ton.ParseAccountID(*address)
		if err != nil {
			log.Fatalf("Failed to parse address: %v", err)
		}
		addressBytes = accountID.Address[:]
		workchain = accountID.Workchain
	}

	if len(addressBytes) != 32 {
		log.Fatalf("Address must be 32 bytes, got %d", len(addressBytes))
	}

	// Connect to Tycho client
	client, err := tychoclient.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Get shard account info
	accountInfo, err := client.GetShardAccountRaw(context.Background(), workchain, addressBytes, *withProof)
	if err != nil {
		log.Fatalf("Failed to get shard account: %v", err)
	}

	// Create fixture
	fixture := ShardAccountFixture{
		Description: "Tycho testnet shard account fixture",
		Workchain:   workchain,
		Address:     hex.EncodeToString(addressBytes),
		WithProof:   *withProof,
		McStateInfo: map[string]interface{}{
			"mc_seqno": accountInfo.McStateInfo.McSeqno,
			"lt":       accountInfo.McStateInfo.Lt,
			"utime":    accountInfo.McStateInfo.Utime,
		},
	}

	// Add proof if requested
	if *withProof && len(accountInfo.Proof) > 0 {
		fixture.Proof = base64.StdEncoding.EncodeToString(accountInfo.Proof)
	}

	// Process account state
	fixture.AccountFound = len(accountInfo.AccountState) > 0
	if fixture.AccountFound {
		fixture.AccountState = base64.StdEncoding.EncodeToString(accountInfo.AccountState)
		parseAccountData(&fixture, accountInfo.AccountState, *debug)
	} else {
		fixture.AccountType = "none"
	}

	// Output
	data, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if *output != "" {
		if err := os.MkdirAll("testdata", 0755); err != nil {
			log.Fatalf("Failed to create testdata directory: %v", err)
		}
		if err := os.WriteFile(*output, data, 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("✅ Saved to %s\n", *output)
	} else {
		fmt.Println(string(data))
	}
}

// parseAccountData parses account state and extracts information
func parseAccountData(fixture *ShardAccountFixture, accountState []byte, includeDebug bool) {
	var debugInfo *DebugInfo
	if includeDebug {
		debugInfo = &DebugInfo{
			AccountStateLength: len(accountState),
		}
		fixture.Debug = debugInfo

		// Try BOC parsing for debug
		cells, err := boc.DeserializeBoc(accountState)
		if err != nil {
			debugInfo.BOCError = err.Error()
		} else {
			debugInfo.BOCCells = make([]string, len(cells))
			for i, cell := range cells {
				bitSize := cell.BitsAvailableForRead()
				debugInfo.BOCCells[i] = fmt.Sprintf("Cell %d: %d bits, %d refs", i, bitSize, cell.RefsSize())
			}
		}
	}

	// Try to parse account
	parsedAccount, err := tychoclient.ParseShardAccount(accountState)
	if err != nil {
		if debugInfo != nil {
			debugInfo.ParseError = err.Error()
		}
		fixture.AccountType = "parse_error"
		return
	}

	// Extract account information
	fixture.LastTransLt = parsedAccount.LastTransLt
	fixture.LastTransHash = hex.EncodeToString(parsedAccount.LastTransHash[:])

	switch parsedAccount.Account.SumType {
	case "AccountNone":
		fixture.AccountType = "none"
	case "Account":
		account := parsedAccount.Account.Account
		fixture.Balance = fmt.Sprintf("%d", uint64(account.Storage.Balance.Grams))

		switch account.Storage.State.SumType {
		case "AccountUninit":
			fixture.AccountType = "uninitialized"
		case "AccountActive":
			fixture.AccountType = "active"
		case "AccountFrozen":
			fixture.AccountType = "frozen"
		default:
			fixture.AccountType = "unknown"
		}
	default:
		fixture.AccountType = "unknown"
	}
}
