package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tonkeeper/tongo/tychoclient"
)

func main() {
	seqno := flag.Uint("seqno", 0, "Block seqno to fetch (0 = latest)")
	output := flag.String("output", "testdata/tycho_block.json", "Output file path")
	flag.Parse()

	client, err := tychoclient.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Get latest seqno if not specified
	if *seqno == 0 {
		status, err := client.GetStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
		if status.McStateInfo == nil {
			log.Fatalf("Node not ready (McStateInfo is nil)")
		}
		// Use McSeqno from McStateInfo
		*seqno = uint(status.McStateInfo.McSeqno)
		fmt.Printf("Using latest masterchain seqno: %d\n", *seqno)
	} // Fetch the block
	blockData, err := client.GetRawBlockData(ctx, -1, 0x8000000000000000, uint32(*seqno))
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	// Parse it to verify it works
	block, err := tychoclient.ParseTychoBlock(blockData)
	if err != nil {
		log.Fatalf("Failed to parse block: %v", err)
	}

	// Create test fixture
	fixture := struct {
		Seqno      uint32 `json:"seqno"`
		Magic      string `json:"magic"`
		GenUtimeMs uint16 `json:"gen_utime_ms"`
		BlockData  string `json:"block_data"` // base64 encoded BOC
	}{
		Seqno:      block.Info.SeqNo,
		Magic:      fmt.Sprintf("0x%x", block.Magic),
		GenUtimeMs: block.Info.GenUtimeMs,
		BlockData:  base64.StdEncoding.EncodeToString(blockData),
	}

	// Save to file
	data, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create testdata directory if it doesn't exist
	os.MkdirAll("testdata", 0755)

	err = os.WriteFile(*output, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("✅ Block saved to %s\n", *output)
	fmt.Printf("   Seqno: %d\n", fixture.Seqno)
	fmt.Printf("   Magic: %s\n", fixture.Magic)
	fmt.Printf("   GenUtimeMs: %d\n", fixture.GenUtimeMs)
	fmt.Printf("   Size: %d bytes\n", len(blockData))
}
