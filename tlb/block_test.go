package tlb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
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
			block, err := readBlock(inputFilename)
			if err != nil {
				t.Fatalf("readBlock() failed: %v", err)
			}
			if tc.folder == "testdata/block-5" && block.Magic == 0 {
				t.Fatalf("block magic not set")
			}
			if tc.folder == "testdata/block-5" {
				t.Logf("block magic: 0x%x state sumtype: %s", block.Magic, block.StateUpdate.ToRoot.SumType)
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

func Test_tlb_MarshallingRoundtrip(t *testing.T) {
	testCases := []struct {
		folder string
	}{
		{folder: "testdata/block-1"},
		{folder: "testdata/block-2"},
		{folder: "testdata/block-3"},
		{folder: "testdata/block-4"},
		{folder: "testdata/block-5"},
		{folder: "testdata/block-6"},
		{folder: "testdata/block-7"},
	}

	for _, tc := range testCases {
		t.Run(tc.folder, func(t *testing.T) {
			inputFilename := path.Join(tc.folder, "block.bin")
			block, err := readBlock(inputFilename)
			if err != nil {
				t.Fatalf("readBlock() failed: %v", err)
			}

			cell := boc.NewCell()
			err = Marshal(cell, block)
			if err != nil {
				t.Fatalf("Marshal() failed: %v", err)
			}

			bocData, err := cell.ToBoc()
			if err != nil {
				t.Fatalf("ToBoc() failed: %v", err)
			}
			cells, err := boc.DeserializeBoc(bocData)
			if err != nil {
				t.Fatalf("DeserializeBoc() failed: %v", err)
			}
			var decoded Block
			if err := Unmarshal(cells[0], &decoded); err != nil {
				dec := NewDecoder().WithDebug()
				var dbg Block
				cells[0].ResetCounters()
				if dbgErr := dec.Unmarshal(cells[0], &dbg); dbgErr != nil {
					t.Fatalf("Unmarshal() failed: %v", dbgErr)
				}
				t.Fatalf("Unmarshal() failed: %v", err)
			}
			if !blocksEqual(block, &decoded) {
				t.Fatalf("block mismatch after marshaling round-trip: %s", describeBlockDiff(block, &decoded))
			}
		})
	}
}

func describeBlockDiff(a, b *Block) string {
	switch {
	case a.Magic != b.Magic:
		return "Magic differs"
	case a.GlobalId != b.GlobalId:
		return "GlobalId differs"
	case !reflect.DeepEqual(a.Info, b.Info):
		return "Info differs"
	case !reflect.DeepEqual(a.ValueFlow, b.ValueFlow):
		return "ValueFlow differs"
	case !reflect.DeepEqual(a.StateUpdate, b.StateUpdate):
		return describeMerkleUpdateDiff("StateUpdate", a.StateUpdate, b.StateUpdate)
	case !reflect.DeepEqual(a.Extra, b.Extra):
		return "Extra differs"
	default:
		return "unknown difference"
	}
}

func describeMerkleUpdateDiff(name string, a, b MerkleUpdate[ShardState]) string {
	switch {
	case a.Magic != b.Magic:
		return fmt.Sprintf("%s magic differs", name)
	case a.FromHash != b.FromHash:
		return fmt.Sprintf("%s FromHash differs", name)
	case a.ToHash != b.ToHash:
		return fmt.Sprintf("%s ToHash differs", name)
	case a.FromDepth != b.FromDepth:
		return fmt.Sprintf("%s FromDepth differs", name)
	case a.ToDepth != b.ToDepth:
		return fmt.Sprintf("%s ToDepth differs", name)
	case !reflect.DeepEqual(a.FromRoot, b.FromRoot):
		return describeShardStateDiff(fmt.Sprintf("%s.FromRoot", name), a.FromRoot, b.FromRoot)
	case !reflect.DeepEqual(a.ToRoot, b.ToRoot):
		return describeShardStateDiff(fmt.Sprintf("%s.ToRoot", name), a.ToRoot, b.ToRoot)
	default:
		return fmt.Sprintf("%s unknown difference", name)
	}
}

func describeShardStateDiff(name string, a, b ShardState) string {
	if a.SumType != b.SumType {
		return fmt.Sprintf("%s sum type differs: %s vs %s", name, a.SumType, b.SumType)
	}
	switch a.SumType {
	case "UnsplitState":
		return describeShardStateUnsplitDiff(fmt.Sprintf("%s.UnsplitState", name), a.UnsplitState.Value, b.UnsplitState.Value)
	case "SplitState":
		switch {
		case !reflect.DeepEqual(a.SplitState.Left, b.SplitState.Left):
			return describeShardStateUnsplitDiff(fmt.Sprintf("%s.SplitState.Left", name), a.SplitState.Left, b.SplitState.Left)
		case !reflect.DeepEqual(a.SplitState.Right, b.SplitState.Right):
			return describeShardStateUnsplitDiff(fmt.Sprintf("%s.SplitState.Right", name), a.SplitState.Right, b.SplitState.Right)
		}
	default:
		return fmt.Sprintf("%s unknown sum type %s", name, a.SumType)
	}
	return fmt.Sprintf("%s unknown difference", name)
}

func describeShardStateUnsplitDiff(name string, a, b ShardStateUnsplit) string {
	switch {
	case a.Magic != b.Magic:
		return fmt.Sprintf("%s magic differs", name)
	case a.ShardStateUnsplit.GlobalID != b.ShardStateUnsplit.GlobalID:
		return fmt.Sprintf("%s GlobalID differs", name)
	case a.ShardStateUnsplit.ShardID != b.ShardStateUnsplit.ShardID:
		return fmt.Sprintf("%s ShardID differs", name)
	case a.ShardStateUnsplit.SeqNo != b.ShardStateUnsplit.SeqNo:
		return fmt.Sprintf("%s SeqNo differs", name)
	case a.ShardStateUnsplit.VertSeqNo != b.ShardStateUnsplit.VertSeqNo:
		return fmt.Sprintf("%s VertSeqNo differs", name)
	case a.ShardStateUnsplit.GenUtime != b.ShardStateUnsplit.GenUtime:
		return fmt.Sprintf("%s GenUtime differs", name)
	case a.ShardStateUnsplit.GenLt != b.ShardStateUnsplit.GenLt:
		return fmt.Sprintf("%s GenLt differs", name)
	case a.ShardStateUnsplit.MinRefMcSeqno != b.ShardStateUnsplit.MinRefMcSeqno:
		return fmt.Sprintf("%s MinRefMcSeqno differs", name)
	case !reflect.DeepEqual(a.ShardStateUnsplit.OutMsgQueueInfo, b.ShardStateUnsplit.OutMsgQueueInfo):
		return fmt.Sprintf("%s OutMsgQueueInfo differs", name)
	case a.ShardStateUnsplit.BeforeSplit != b.ShardStateUnsplit.BeforeSplit:
		return fmt.Sprintf("%s BeforeSplit differs", name)
	case !reflect.DeepEqual(a.ShardStateUnsplit.Accounts, b.ShardStateUnsplit.Accounts):
		return describeHashmapAugEDiff(fmt.Sprintf("%s.Accounts", name), a.ShardStateUnsplit.Accounts, b.ShardStateUnsplit.Accounts)
	case !reflect.DeepEqual(a.ShardStateUnsplit.Other, b.ShardStateUnsplit.Other):
		return fmt.Sprintf("%s Other differs", name)
	case !reflect.DeepEqual(a.ShardStateUnsplit.Custom, b.ShardStateUnsplit.Custom):
		return fmt.Sprintf("%s Custom differs", name)
	default:
		return fmt.Sprintf("%s unknown difference", name)
	}
}

func describeHashmapAugEDiff[keyT fixedSize, T1, T2 any](name string, a, b HashmapAugE[keyT, T1, T2]) string {
	switch {
	case len(a.m.keys) != len(b.m.keys):
		return fmt.Sprintf("%s key count differs", name)
	case !reflect.DeepEqual(a.m.keys, b.m.keys):
		return fmt.Sprintf("%s keys differ", name)
	case !reflect.DeepEqual(a.m.values, b.m.values):
		if len(a.m.values) != len(b.m.values) {
			return fmt.Sprintf("%s value count differs: %d vs %d", name, len(a.m.values), len(b.m.values))
		}
		for i := range a.m.values {
			if !reflect.DeepEqual(a.m.values[i], b.m.values[i]) {
				return fmt.Sprintf("%s value at index %d differs: %#v vs %#v", name, i, a.m.values[i], b.m.values[i])
			}
		}
		return fmt.Sprintf("%s values differ", name)
	case !reflect.DeepEqual(a.m.extra, b.m.extra):
		return fmt.Sprintf("%s inner extras differ", name)
	case !reflect.DeepEqual(a.extra, b.extra):
		return fmt.Sprintf("%s extra differs", name)
	default:
		return fmt.Sprintf("%s unknown difference", name)
	}
}

func blocksEqual(a, b *Block) bool {
	aj, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return bytes.Equal(aj, bj)
}

func readBlock(filename string) (*Block, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBoc(data)
	if err != nil {
		return nil, fmt.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	var block Block
	return &block, Unmarshal(cell[0], &block)
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
			block, err := readBlock(inputFilename)
			if err != nil {
				t.Fatalf("readBlock() failed: %v", err)
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
				rawKAddr := fmt.Sprintf("%v:%v", k.Workchain, k.Address.Hex())
				actualV, ok := tc.actual[MessageIDSimple{
					Address: rawKAddr,
					Lt:      k.Lt,
				}]
				if !ok {
					t.Errorf("Extra metadata found (%v, %v)", rawKAddr, k.Lt)
					continue
				}
				rawVAddr := fmt.Sprintf("%v:%v", v.Workchain, v.Address.Hex())
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
