package tlb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func Test_tlb_Unmarshal(t *testing.T) {
	type Transaction struct {
		AccountAddr   string
		Lt            uint64
		PrevTransHash string
		PrevTransLt   uint64
		Now           uint32
		OutMsgCnt     Uint15
		OrigStatus    AccountStatus
		EndStatus     AccountStatus
	}
	type AccountBlock struct {
		Transactions map[uint64]Transaction
	}
	type BlockContent struct {
		Accounts          map[string]*AccountBlock
		TxHashes          []string
		ValueFlow         ValueFlow
		InMsgDescrLength  int
		OutMsgDescrLength int
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
			accounts := map[string]*AccountBlock{}
			var txHashes []string
			for _, account := range block.Extra.AccountBlocks.Values() {
				accBlock, ok := accounts[hex.EncodeToString(account.AccountAddr[:])]
				if !ok {
					accBlock = &AccountBlock{Transactions: map[uint64]Transaction{}}
					accounts[hex.EncodeToString(account.AccountAddr[:])] = accBlock
				}
				for _, txRef := range account.Transactions.Values() {
					tx := txRef.Value
					accBlock.Transactions[txRef.Value.Lt] = Transaction{
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
			inMsgLength, err := block.Extra.InMsgDescrLength()
			if err != nil {
				t.Errorf("InMsgDescrLength() failed: %v", err)
			}
			outMsgLength, err := block.Extra.OutMsgDescrLength()
			if err != nil {
				t.Errorf("InMsgDescrLength() failed: %v", err)
			}
			blk := BlockContent{
				Accounts:          accounts,
				TxHashes:          txHashes,
				ValueFlow:         block.ValueFlow,
				InMsgDescrLength:  inMsgLength,
				OutMsgDescrLength: outMsgLength,
			}
			bs, err := json.MarshalIndent(blk, " ", "  ")
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
			if bytes.Compare(bytes.Trim(content, " \n"), bytes.Trim(bs, " \n")) != 0 {
				t.Errorf("tx hashes differ")
			}
		})
	}
}
