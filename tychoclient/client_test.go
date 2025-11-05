package tychoclient

import (
	"context"
	"testing"
	"time"
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
	_, err = client.GetShardAccountAtSeqno(ctx, 0, dummyAddress, false, -1, 0x8000000000000000, targetSeqno)
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
	bocData, err := client.GetRawBlockData(ctx, -1, 0x8000000000000000, status.McStateInfo.McSeqno)
	if err != nil {
		t.Fatalf("Failed to get raw block data: %v", err)
	}

	if len(bocData) == 0 {
		t.Error("Empty BOC data received")
	}

	t.Logf("Successfully got raw block data: %d bytes", len(bocData))
}
