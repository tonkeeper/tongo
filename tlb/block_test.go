package tlb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func Test_tlb_Unmarshal(t *testing.T) {
	type transaction struct {
		AccountAddr   string
		Lt            uint64
		PrevTransHash string
		PrevTransLt   uint64
		Now           uint32
		OutMsgCnt     Uint15
		OrigStatus    AccountStatus
		EndStatus     AccountStatus
	}
	type accountBlock struct {
		Transactions map[uint64]transaction
	}
	testCases := []struct {
		name   string
		folder string
	}{
		{
			name:   "all good",
			folder: "testdata/block-1",
		},
		{
			name:   "all good",
			folder: "testdata/block-2",
		},
		{
			name:   "all good",
			folder: "testdata/block-3",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputFilename := path.Join(tc.folder, "block.bin")
			data, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Errorf("ReadFile() failed: %v", err)
			}
			cell, err := boc.DeserializeBoc(data)
			if err != nil {
				t.Errorf("boc.DeserializeBoc() failed: %v", err)
			}
			var block Block
			err = Unmarshal(cell[0], &block)
			if err != nil {
				t.Errorf("Unmarshal() failed: %v", err)
			}
			accounts := map[string]*accountBlock{}
			var txHashes []string
			for _, account := range block.Extra.AccountBlocks.Values() {
				accBlock, ok := accounts[hex.EncodeToString(account.AccountAddr[:])]
				if !ok {
					accBlock = &accountBlock{Transactions: map[uint64]transaction{}}
					accounts[hex.EncodeToString(account.AccountAddr[:])] = accBlock
				}
				for _, txRef := range account.Transactions.Values() {
					tx := txRef.Value
					accBlock.Transactions[txRef.Value.Lt] = transaction{
						AccountAddr:   hex.EncodeToString(tx.AccountAddr[:]),
						Lt:            tx.Lt,
						PrevTransHash: hex.EncodeToString(tx.PrevTransHash[:]),
						PrevTransLt:   tx.PrevTransLt,
						Now:           tx.Now,
						OutMsgCnt:     tx.OutMsgCnt,
						OrigStatus:    tx.OrigStatus,
						EndStatus:     tx.EndStatus,
					}
					txHashes = append(txHashes, tx.Hash().Hex())
				}
			}
			sort.Slice(txHashes, func(i, j int) bool {
				return txHashes[i] < txHashes[j]
			})
			bs, err := json.MarshalIndent(accounts, " ", "  ")
			if err != nil {
				t.Errorf("json.MarshalIndent() failed: %v", err)
			}
			outputFilename := path.Join(tc.folder, "block.output.json")
			if err := os.WriteFile(outputFilename, bs, 0644); err != nil {
				t.Errorf("WriteFile() failed: %v", err)
			}
			expectedFilename := path.Join(tc.folder, "block.expected.json")
			content, err := os.ReadFile(expectedFilename)
			if err != nil {
				t.Errorf("ReadFile() failed: %v", err)
			}
			expectedAccounts := map[string]*accountBlock{}
			if err := json.Unmarshal(content, &expectedAccounts); err != nil {
				t.Errorf("json.Unmarshal() failed: %v", err)
			}
			if !reflect.DeepEqual(accounts, expectedAccounts) {
				t.Errorf("expectedAccounts differs from accounts")
			}
			bs, err = json.MarshalIndent(block.ValueFlow, " ", "  ")
			if err != nil {
				t.Fatalf("MarshalIndent() failed: %v", err)
			}
			valueFlowOutput := path.Join(tc.folder, "value-flow.output.json")
			if err := os.WriteFile(valueFlowOutput, bs, 0644); err != nil {
				t.Fatalf("WriteFile() failed: %v", err)
			}
			expectedValueFlowFilename := path.Join(tc.folder, "value-flow.expected.json")
			data, err = os.ReadFile(expectedValueFlowFilename)
			if err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			}
			if bytes.Compare(bytes.Trim(bs, " \n"), bytes.Trim(data, " \n")) != 0 {
				t.Errorf("ValueFlows differ")
			}
			// compare hashes
			bs, err = json.MarshalIndent(txHashes, " ", "  ")
			if err != nil {
				t.Fatalf("MarshalIndent() failed: %v", err)
			}
			hashesOutputFilename := path.Join(tc.folder, "tx-hashes.output.json")
			if err := os.WriteFile(hashesOutputFilename, bs, 0644); err != nil {
				t.Fatalf("WriteFile() failed: %v", err)
			}
			expectedHashesFilename := path.Join(tc.folder, "tx-hashes.expected.json")
			data, err = os.ReadFile(expectedHashesFilename)
			if err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			}
			if bytes.Compare(bytes.Trim(bs, " \n"), bytes.Trim(data, " \n")) != 0 {
				t.Errorf("tx hashes differ")
			}
		})
	}
}
