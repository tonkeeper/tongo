package tychoclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/tlb"
)

const (
	masterchainWorkchain = -1
	masterchainShard     = 0x8000000000000000 // Full shard range
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if client.conn == nil {
		t.Fatal("Connection is nil")
	}
	if client.client == nil {
		t.Fatal("Client is nil")
	}
}

func TestGetStatus(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.Version == 0 {
		t.Error("Version should not be 0")
	}
	if status.Timestamp == 0 {
		t.Error("Timestamp should not be 0")
	}
}

func TestGetShardAccount(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with a dummy address (32 bytes)
	dummyAddress := make([]byte, 32)
	copy(dummyAddress, []byte("test_address_for_tycho_demo_12"))

	// This will likely fail since the address doesn't exist, but tests the API
	_, err = client.GetShardAccount(ctx, 0, dummyAddress, false)
	if err != nil {
		t.Logf("GetShardAccount failed as expected: %v", err)
		// This is expected for a non-existent account
	} else {
		t.Logf("GetShardAccount succeeded unexpectedly")
	}
}

func TestGetShardAccountAtSeqno(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get status to find a valid seqno
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.McStateInfo == nil || status.McStateInfo.McSeqno == 0 {
		t.Skip("No valid seqno available")
	}

	// Test with a dummy address
	dummyAddress := make([]byte, 32)
	copy(dummyAddress, []byte("test_address_for_tycho_demo_12"))

	targetSeqno := status.McStateInfo.McSeqno - 1
	_, err = client.GetShardAccountAtSeqno(ctx, 0, dummyAddress, false, masterchainWorkchain, masterchainShard, targetSeqno)
	if err != nil {
		t.Logf("GetShardAccountAtSeqno failed as expected: %v", err)
		// This is expected for a non-existent account
	} else {
		t.Logf("GetShardAccountAtSeqno succeeded unexpectedly")
	}
}

func TestGetRawBlockData(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get status to find a valid seqno
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.McStateInfo == nil {
		t.Skip("Node not ready")
	}

	// Get raw block data
	bocData, err := client.GetRawBlockData(ctx, masterchainWorkchain, masterchainShard, status.McStateInfo.McSeqno)
	if err != nil {
		t.Fatalf("Failed to get raw block data: %v", err)
	}

	if len(bocData) == 0 {
		t.Error("Empty BOC data received")
	}

	t.Logf("Successfully got raw block data: %d bytes", len(bocData))
}

// TestParseTychoBlock verifies that we can parse Tycho blocks and extract Tycho-specific fields.
// This test demonstrates the TLB parser working with Tycho's modified block structure:
// - Block magic: 0x11ef55bb (vs TON's 0x11ef55aa)
// - BlockInfo with gen_utime_ms field (millisecond precision timestamp)
// - BlockExtra magic: 0x4a33f6fc (vs TON's 0x4a33f6fd)
func TestParseTychoBlock(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get status first
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.McStateInfo == nil {
		t.Skip("Node not ready")
	}

	// Get raw block data
	bocData, err := client.GetRawBlockData(ctx, masterchainWorkchain, masterchainShard, status.McStateInfo.McSeqno)
	if err != nil {
		t.Fatalf("Failed to get raw block data: %v", err)
	}

	// Parse as Tycho block
	block, err := ParseTychoBlock(bocData)
	if err != nil {
		t.Fatalf("Failed to parse Tycho block: %v", err)
	}

	// Verify block structure
	if block.GlobalId == 0 {
		t.Error("GlobalId should not be 0")
	}

	// Verify Tycho-specific fields
	if block.Info.GenUtimeMs == 0 {
		t.Error("GenUtimeMs should not be 0 (Tycho-specific field)")
	}

	if block.Info.SeqNo != status.McStateInfo.McSeqno {
		t.Errorf("Block seqno mismatch: got %d, want %d", block.Info.SeqNo, status.McStateInfo.McSeqno)
	}

	t.Logf("Successfully parsed Tycho block:")
	t.Logf("  GlobalId: %d", block.GlobalId)
	t.Logf("  SeqNo: %d", block.Info.SeqNo)
	t.Logf("  GenUtime: %d", block.Info.GenUtime)
	t.Logf("  GenUtimeMs: %d (Tycho-specific)", block.Info.GenUtimeMs)
	t.Logf("  StartLt: %d", block.Info.StartLt)
	t.Logf("  EndLt: %d", block.Info.EndLt)

	// Verify ValueFlow
	t.Logf("  ValueFlow.ToNextBlk.Grams: %d", block.ValueFlow.ToNextBlk.Grams)
	t.Logf("  ValueFlow.Exported.Grams: %d", block.ValueFlow.Exported.Grams)

	// Verify message descriptors
	inMsgDescr, err := block.Extra.InMsgDescr()
	if err != nil {
		t.Logf("  InMsgDescr: parsing failed: %v", err)
	} else {
		t.Logf("  InMsgDescr: %d messages", len(inMsgDescr.Keys()))
	}

	outMsgDescr, err := block.Extra.OutMsgDescr()
	if err != nil {
		t.Logf("  OutMsgDescr: parsing failed: %v", err)
	} else {
		t.Logf("  OutMsgDescr: %d messages", len(outMsgDescr.Keys()))
	}
}

// TestParseTychoBlockFromFixture tests parsing using a saved block fixture.
// This test is faster and more deterministic than fetching from the API.
// To generate a new fixture, run: go run cmd/fetch_test_block/main.go
func TestParseTychoBlockFromFixture(t *testing.T) {
	// Read the test fixture
	data, err := os.ReadFile("testdata/tycho_block.json")
	if err != nil {
		t.Skipf("Test fixture not found (run: go run cmd/fetch_test_block/main.go): %v", err)
	}

	var fixture struct {
		Seqno      uint32 `json:"seqno"`
		Magic      string `json:"magic"`
		GenUtimeMs uint16 `json:"gen_utime_ms"`
		BlockData  string `json:"block_data"`
	}
	err = json.Unmarshal(data, &fixture)
	if err != nil {
		t.Fatalf("Failed to parse fixture: %v", err)
	}

	// Decode BOC
	blockData, err := base64.StdEncoding.DecodeString(fixture.BlockData)
	if err != nil {
		t.Fatalf("Failed to decode block data: %v", err)
	}

	t.Logf("Testing with fixture: seqno=%d, size=%d bytes", fixture.Seqno, len(blockData))

	// Parse the block
	block, err := ParseTychoBlock(blockData)
	if err != nil {
		t.Fatalf("Failed to parse block: %v", err)
	}

	// === VERIFY BLOCK HEADER ===
	if block.Magic != 0x11ef55bb {
		t.Errorf("Expected magic 0x11ef55bb, got 0x%x", block.Magic)
	}

	if block.GlobalId == 0 {
		t.Error("GlobalId should not be 0")
	}

	// === VERIFY BLOCKINFO ===
	// Note: Magic field is not exposed in TychoBlockInfo (same as TON's BlockInfo)
	// It's only used during unmarshaling

	if block.Info.SeqNo != fixture.Seqno {
		t.Errorf("Expected seqno %d, got %d", fixture.Seqno, block.Info.SeqNo)
	}

	if block.Info.GenUtimeMs != fixture.GenUtimeMs {
		t.Errorf("Expected gen_utime_ms %d, got %d", fixture.GenUtimeMs, block.Info.GenUtimeMs)
	}

	if block.Info.GenUtime == 0 {
		t.Error("GenUtime should not be 0")
	}

	if block.Info.StartLt == 0 {
		t.Error("StartLt should not be 0")
	}

	if block.Info.EndLt == 0 {
		t.Error("EndLt should not be 0")
	}

	if block.Info.EndLt < block.Info.StartLt {
		t.Errorf("EndLt (%d) should be >= StartLt (%d)", block.Info.EndLt, block.Info.StartLt)
	}

	// === VERIFY VALUEFLOW ===
	t.Logf("✅ ValueFlow:")
	t.Logf("   ToNextBlk.Grams: %d", block.ValueFlow.ToNextBlk.Grams)
	t.Logf("   Exported.Grams: %d", block.ValueFlow.Exported.Grams)
	t.Logf("   FeesCollected.Grams: %d", block.ValueFlow.FeesCollected.Grams)

	// === VERIFY OUTMSGQUEUEUPDATES ===
	if block.Other.OutMsgQueueUpdates.Magic != 0x1 {
		t.Errorf("Expected OutMsgQueueUpdates magic 0x1, got 0x%x", block.Other.OutMsgQueueUpdates.Magic)
	}

	emptyHash := tlb.Bits256{}
	if block.Other.OutMsgQueueUpdates.DiffHash == emptyHash {
		t.Error("OutMsgQueueUpdates.DiffHash should not be empty")
	}

	// TailLen can be 0 for empty queues, so don't check it

	// === VERIFY BLOCKEXTRA ===
	if block.Extra.RandSeed == emptyHash {
		t.Error("Extra.RandSeed should not be empty")
	}

	if block.Extra.CreatedBy == emptyHash {
		t.Error("Extra.CreatedBy should not be empty")
	}

	// AccountBlocks can be empty for blocks with no transactions
	numAccounts := len(block.Extra.AccountBlocks.Keys())
	if numAccounts == 0 {
		t.Log("Note: Block has no account blocks (no transactions)")
	}

	// === LOGGING: Detailed verification results ===
	t.Logf("✅ Block Header:")
	t.Logf("   Magic: 0x%x (Tycho)", block.Magic)
	t.Logf("   GlobalId: %d", block.GlobalId)

	t.Logf("✅ BlockInfo:")
	// Note: Magic field is not exposed (same as TON's BlockInfo pattern)
	t.Logf("   SeqNo: %d", block.Info.SeqNo)
	t.Logf("   GenUtime: %d", block.Info.GenUtime)
	t.Logf("   GenUtimeMs: %d (Tycho-specific!)", block.Info.GenUtimeMs)
	t.Logf("   NotMaster: %v, KeyBlock: %v", block.Info.NotMaster, block.Info.KeyBlock)
	t.Logf("   StartLt: %d, EndLt: %d", block.Info.StartLt, block.Info.EndLt)

	t.Logf("✅ OutMsgQueueUpdates:")
	t.Logf("   Magic: 0x%x", block.Other.OutMsgQueueUpdates.Magic)
	t.Logf("   TailLen: %d", block.Other.OutMsgQueueUpdates.TailLen)
	t.Logf("   DiffHash: %x... (first 8 bytes)", block.Other.OutMsgQueueUpdates.DiffHash[:8])

	t.Logf("✅ BlockExtra:")
	t.Logf("   RandSeed: %x... (first 8 bytes)", block.Extra.RandSeed[:8])
	t.Logf("   CreatedBy: %x... (first 8 bytes)", block.Extra.CreatedBy[:8])
	t.Logf("   AccountBlocks: %d accounts", numAccounts)

	// Verify message descriptors
	inMsgDescr, err := block.Extra.InMsgDescr()
	if err != nil {
		t.Logf("   InMsgDescr: parsing failed: %v", err)
	} else {
		t.Logf("   InMsgDescr: %d messages", len(inMsgDescr.Keys()))
	}

	outMsgDescr, err := block.Extra.OutMsgDescr()
	if err != nil {
		t.Logf("   OutMsgDescr: parsing failed: %v", err)
	} else {
		t.Logf("   OutMsgDescr: %d messages", len(outMsgDescr.Keys()))
	}

	if block.Extra.Custom.Exists {
		t.Logf("   McBlockExtra: exists")
	} else {
		t.Logf("   McBlockExtra: not present")
	}

	// === VERIFY BLOCKINFO CONDITIONAL FIELDS ===
	if !block.Info.NotMaster {
		t.Logf("ℹ️  Block is masterchain block")
		if block.Info.MasterRef != nil {
			t.Error("MasterRef should be nil for masterchain blocks")
		}
	} else {
		t.Logf("ℹ️  Block is shardchain block")
		if block.Info.MasterRef == nil {
			t.Error("MasterRef should not be nil for shardchain blocks")
		}
	}

	t.Logf("\n🎉 Complete block data validation passed!")
}

// TestParseTychoBlockErrorCases tests error handling in the parser
func TestParseTychoBlockErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		bocData []byte
		wantErr bool
	}{
		{
			name:    "empty data",
			bocData: []byte{},
			wantErr: true,
		},
		{
			name:    "invalid BOC header",
			bocData: []byte{0x00, 0x01, 0x02, 0x03, 0x04},
			wantErr: true,
		},
		{
			name:    "nil data",
			bocData: nil,
			wantErr: true,
		},
		{
			name:    "too short data",
			bocData: []byte{0xb5, 0xee},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTychoBlock(tt.bocData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTychoBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Logf("Got expected error: %v", err)
			}
		})
	}
}

// TestParseTychoShardAccount tests parsing of shard account data fetched from Tycho.
func TestParseTychoShardAccount(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with a zero address (most likely to return some data)
	testAddress := make([]byte, 32)

	accountInfo, err := client.GetShardAccount(ctx, 0, testAddress, false)
	if err != nil {
		t.Fatalf("GetShardAccount failed, account might not exists: %v", err)
	}

	if len(accountInfo.AccountState) == 0 {
		t.Skip("No account data to test parsing")
	}

	// Verify that ParseShardAccount works
	if accountInfo.ParsedAccountState == nil {
		t.Error("ParsedAccountState should not be nil when AccountState exists")
	}

	t.Logf("Successfully parsed account data: %d bytes -> cell with %d bits, %d refs",
		len(accountInfo.AccountState),
		accountInfo.ParsedAccountState.BitsAvailableForRead(),
		accountInfo.ParsedAccountState.RefsAvailableForRead())
}

// TestParseTychoShardAccountAtSeqno tests fetching and parsing account data at a specific seqno
func TestParseTychoShardAccountAtSeqno(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get status to find a valid seqno
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.McStateInfo == nil || status.McStateInfo.McSeqno < 2 {
		t.Skip("Need at least seqno 2 for historical testing")
	}

	// Test with a historical seqno
	targetSeqno := status.McStateInfo.McSeqno - 1
	testAddress := make([]byte, 32) // Zero address

	accountInfo, err := client.GetShardAccountAtSeqno(ctx, 0, testAddress, false, masterchainWorkchain, masterchainShard, targetSeqno)
	if err != nil {
		t.Logf("GetShardAccountAtSeqno failed: %v", err)
		return // Expected for non-existent accounts
	}

	// If we got data, verify parsing works
	if len(accountInfo.AccountState) > 0 && accountInfo.ParsedAccountState == nil {
		t.Error("ParsedAccountState should not be nil when AccountState exists")
	}

	t.Logf("Historical account parsing test completed")
}
