package tlb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
		Libraries         map[string]map[string]struct{}
	}
	testCases := []struct {
		name   string
		folder string
	}{
		{
			name:   "block (0,8000000000000000,30816553)",
			folder: "testdata/block-1",
		},
		{
			name:   "block (0,8000000000000000,40484416)",
			folder: "testdata/block-2",
		},
		{
			name:   "block (0,8000000000000000,40484438)",
			folder: "testdata/block-3",
		},
		{
			name:   "block (0,D83800000000000,4168601)",
			folder: "testdata/block-4",
		},
		{
			name:   "block (0,D83800000000000,(-1,8000000000000000,17734191)",
			folder: "testdata/block-5",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			inputFilename := path.Join(tc.folder, "block.bin")
			data, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			}
			cell, err := boc.DeserializeBoc(data)
			if err != nil {
				t.Fatalf("boc.DeserializeBoc() failed: %v", err)
			}
			var block Block
			err = Unmarshal(cell[0], &block)
			if err != nil {
				t.Fatalf("Unmarshal() failed: %v", err)
			}
			accounts := map[string]*AccountBlock{}
			var txHashes []string
			for _, account := range block.Extra.AccountBlocks.Values() {
				accBlock, ok := accounts[hex.EncodeToString(account.AccountAddr[:])]
				if !ok {
					accBlock = &AccountBlock{
						Transactions: map[uint64]Transaction{},
					}
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
				Libraries:         libraries(&block.StateUpdate.ToRoot),
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
				t.Errorf("block content mismatch")
			}
		})
	}
}

func libraries(s *ShardState) map[string]map[string]struct{} {
	libs := map[string]map[string]struct{}{}
	for _, item := range s.UnsplitState.Value.ShardStateUnsplit.Other.Libraries.Items() {
		lib := fmt.Sprintf("%x", item.Key)
		if _, ok := libs[lib]; !ok {
			libs[lib] = map[string]struct{}{}
		}
		for _, pub := range item.Value.Publishers.Keys() {
			libs[lib][fmt.Sprintf("%x", pub)] = struct{}{}
		}
	}
	return libs
}

func Test_GetInMsgsMetadata(t *testing.T) {
	type MessageIDSimple struct {
		Address string
		Lt      uint64
	}
	testCases := []struct {
		name   string
		folder string
		actual map[MessageIDSimple]MessageIDSimple
	}{
		{
			name:   "block (0,a000000000000000,54302208)",
			folder: "testdata/block-6",
			actual: map[MessageIDSimple]MessageIDSimple{
				{
					Address: "0:5691bb16d415f3c11a0b2864415c605eb0f12a17456ca9bbedfce51008aaf2b9",
					Lt:      58969747000002,
				}: {
					Address: "0:5691bb16d415f3c11a0b2864415c605eb0f12a17456ca9bbedfce51008aaf2b9",
					Lt:      58969747000001,
				},
				{
					Address: "0:01cc5bf05b36701c6758a02bc501ac328dc253ae06ff9efce770b2e7b012656b",
					Lt:      58969748000002,
				}: {
					Address: "0:01cc5bf05b36701c6758a02bc501ac328dc253ae06ff9efce770b2e7b012656b",
					Lt:      58969748000001,
				},
				{
					Address: "0:7321efca606adece6492dee6bd6ee3343fcbf31273050ebb2fd69ac59f03f8de",
					Lt:      58969747000012,
				}: {
					Address: "0:8c6c53da461a65ee09e6dccb012e2ad4e4870913eb74bdac124bd9f9062ae69c",
					Lt:      58969744000001,
				},
				{
					Address: "0:140beaba8ef50e82edbb5f74ee421638461ff634fd88acec628eafdc31c91614",
					Lt:      58969748000002,
				}: {
					Address: "0:ebb3f08c9e7e66035e54b1b3abf4d177f889a055d6b40b2c70d728040d8a62f6",
					Lt:      58969744000001,
				},
				{
					Address: "0:52eea2f3c4eb0f2335518fa86303934cbc0d75a23eaa4b7a1636a08a274e2d19",
					Lt:      58969747000003,
				}: {
					Address: "0:cdce58745d265d6f9fd0ae5c79423c991d500eef8bf9a0c79556bb45ff956dd6",
					Lt:      58969743000001,
				},
				{
					Address: "0:01cc5bf05b36701c6758a02bc501ac328dc253ae06ff9efce770b2e7b012656b",
					Lt:      58969748000003,
				}: {
					Address: "0:01cc5bf05b36701c6758a02bc501ac328dc253ae06ff9efce770b2e7b012656b",
					Lt:      58969748000001,
				},
				{
					Address: "0:61b2678a5d4c2db1e24382d875429662560d29d1988ea31147ccdda0fb2c253b",
					Lt:      58969747000003,
				}: {
					Address: "0:8c6c53da461a65ee09e6dccb012e2ad4e4870913eb74bdac124bd9f9062ae69c",
					Lt:      58969744000001,
				},
				{
					Address: "0:d9fff84e4578537a82b56114eaccbfb9d8714ccae54e81c5182e8402066f6193",
					Lt:      58969747000002,
				}: {
					Address: "0:d9fff84e4578537a82b56114eaccbfb9d8714ccae54e81c5182e8402066f6193",
					Lt:      58969747000001,
				},
			},
		},
		{
			name:   "block (0,a000000000000000,54302202)",
			folder: "testdata/block-7",
			actual: map[MessageIDSimple]MessageIDSimple{
				{
					Address: "0:5f33baaf2bb55f85d5c3ee6a4375a4857f09aa0154121b9a39ae854bf615f709",
					Lt:      58969742000002,
				}: {
					Address: "0:5f33baaf2bb55f85d5c3ee6a4375a4857f09aa0154121b9a39ae854bf615f709",
					Lt:      58969742000001,
				},
				{
					Address: "0:6181aebf31392c2c373d513b2b265de5cc6b4269f9853a5e80caf3f40c4cb56b",
					Lt:      58969742000003,
				}: {
					Address: "0:86a00101e0a1b9ec5b9a135be16d66b71ef32f17ea028895240ebecfe8f3aa23",
					Lt:      58969737000001,
				},
				{
					Address: "0:53602c7b4c681581a22bdedb89debffdd44e855010562ed4fbde54d8626ac873",
					Lt:      58969742000008,
				}: {
					Address: "0:8f0e6d2cc08fb952d58e5e1d438107bec7868943d87a84ba03212ac119b102b1",
					Lt:      58969737000001,
				},
				{
					Address: "0:39b1a920cd03f7ef701c80527de2ade5a0aa077392f41e62079d4a54b59606a2",
					Lt:      58969743000011,
				}: {
					Address: "0:6850c5ff38020184f7bfd667864c99d4e18ded1a19d692b33262f2a8afdc0fc5",
					Lt:      58969721000001,
				},
				{
					Address: "0:010433a2f450726099588bb155c56b9d058d3cfc64868e05bca84a95b91814a9",
					Lt:      58969743000006,
				}: {
					Address: "0:6850c5ff38020184f7bfd667864c99d4e18ded1a19d692b33262f2a8afdc0fc5",
					Lt:      58969721000001,
				},
				{
					Address: "0:06653502ddafd029af877a6ce0aa35cc4f5ff2c0b8914751c4ce8651ca430359",
					Lt:      58969743000004,
				}: {
					Address: "0:76db58f0da229ddcfc4ca3239a492f48617b1b7b5c49e831d5d94c41c4dfaa26",
					Lt:      58969739000001,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputFilename := path.Join(tc.folder, "block.bin")
			data, err := os.ReadFile(inputFilename)
			if err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			}
			cell, err := boc.DeserializeBoc(data)
			if err != nil {
				t.Fatalf("boc.DeserializeBoc() failed: %v", err)
			}
			var block Block
			err = Unmarshal(cell[0], &block)
			if err != nil {
				t.Fatalf("Unmarshal() failed: %v", err)
			}
			inMsgsMetadata, err := block.GetInMsgsMetadata()
			if err != nil {
				t.Errorf("GetInMsgsMetadata() failed: %v", err)
			}
			if len(inMsgsMetadata) != len(tc.actual) {
				t.Errorf("Length haven't match: expected %v, got %v", len(tc.actual), len(inMsgsMetadata))
				return
			}
			for k, v := range inMsgsMetadata {
				rawKAddr := "0:" + k.Address.AddrStd.Address.Hex()
				actualV, ok := tc.actual[MessageIDSimple{
					Address: rawKAddr,
					Lt:      k.Lt,
				}]
				if !ok {
					t.Errorf("Extra metadata found (%v, %v)", rawKAddr, k.Lt)
					continue
				}
				rawVAddr := "0:" + v.Address.AddrStd.Address.Hex()
				if rawVAddr != actualV.Address {
					t.Errorf("Addresses haven't matched: expected %v, got %v", actualV.Address, rawVAddr)
				}
				if v.Lt != actualV.Lt {
					t.Errorf("Lt haven't matched: expected %v, got %v", actualV.Lt, v.Lt)
				}
			}
		})
	}
}
