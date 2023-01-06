package tongo

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

func Benchmark_Tlb_Unmarshal(b *testing.B) {
	data, err := os.ReadFile("testdata/raw-block.bin")
	if err != nil {
		b.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		b.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cell[0].ResetCounters()
		var block Block
		err = tlb.Unmarshal(cell[0], &block)
		if err != nil {
			b.Errorf("Unmarshal() failed: %v", err)
		}
	}
}

func Test_tlb_Unmarshal(t *testing.T) {
	type transaction struct {
		AccountAddr   string
		Lt            uint64
		PrevTransHash string
		PrevTransLt   uint64
		Now           uint32
		OutMsgCnt     tlb.Uint15
		OrigStatus    AccountStatus
		EndStatus     AccountStatus
	}
	type accountBlock struct {
		Transactions map[uint64]transaction
	}
	data, err := os.ReadFile("testdata/raw-block.bin")
	if err != nil {
		t.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		t.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	var block Block
	err = tlb.Unmarshal(cell[0], &block)
	if err != nil {
		t.Errorf("Unmarshal() failed: %v", err)
	}
	accounts := map[string]*accountBlock{}
	for _, account := range block.Extra.AccountBlocks.Values() {
		accBlock, ok := accounts[account.AccountAddr.Hex()]
		if !ok {
			accBlock = &accountBlock{Transactions: map[uint64]transaction{}}
			accounts[account.AccountAddr.Hex()] = accBlock
		}
		for _, txRef := range account.Transactions.Values() {
			tx := txRef.Value
			accBlock.Transactions[txRef.Value.Lt] = transaction{
				AccountAddr:   tx.AccountAddr.Hex(),
				Lt:            tx.Lt,
				PrevTransHash: tx.PrevTransHash.Hex(),
				PrevTransLt:   tx.PrevTransLt,
				Now:           tx.Now,
				OutMsgCnt:     tx.OutMsgCnt,
				OrigStatus:    tx.OrigStatus,
				EndStatus:     tx.EndStatus,
			}
		}
	}

	bs, err := json.MarshalIndent(accounts, " ", "  ")
	if err != nil {
		t.Errorf("json.MarshalIndent() failed: %v", err)
	}
	if err := os.WriteFile("testdata/raw-block.output.json", bs, 0644); err != nil {
		t.Errorf("WriteFile() failed: %v", err)
	}
	content, err := os.ReadFile("testdata/raw-block.expected.json")
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
}
