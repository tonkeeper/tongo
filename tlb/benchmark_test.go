package tlb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func Benchmark_Tlb_Unmarshal(b *testing.B) {
	data, err := os.ReadFile("../testdata/raw-13516764.bin")
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
		decoder := NewDecoder()
		err = decoder.Unmarshal(cell[0], &block)
		if err != nil {
			b.Errorf("Unmarshal() failed: %v", err)
		}
	}
}

func Test_block(b *testing.T) {
	data, err := os.ReadFile("../testdata/raw-13516764.bin")
	if err != nil {
		b.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		b.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	var block Block
	decoder := NewDecoder()
	err = decoder.Unmarshal(cell[0], &block)
	if err != nil {
		b.Errorf("Unmarshal() failed: %v", err)
	}
}

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
	data, err := os.ReadFile("testdata/raw-block.bin")
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
	bs, err = json.MarshalIndent(block.ValueFlow, " ", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() failed: %v", err)
	}
	if err := os.WriteFile("testdata/value-flow.output.json", bs, 0644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	data, err = os.ReadFile("testdata/value-flow.expected.json")
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}
	if bytes.Compare(bytes.Trim(bs, " \n"), bytes.Trim(data, " \n")) != 0 {
		t.Errorf("ValueFlows differ")

	}

}
