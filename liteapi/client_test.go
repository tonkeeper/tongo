package liteapi

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi/pool"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"golang.org/x/exp/maps"
)

func TestNewClient_WithMaxConnectionsNumber(t *testing.T) {
	t.Skip("when public lite servers are down, this test will fail")
	cli, err := NewClient(Mainnet())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	if cli.pool.ConnectionsNumber() != defaultMaxConnectionsNumber {
		t.Fatalf("want connections number: %v, got: %v", defaultMaxConnectionsNumber, cli.pool.ConnectionsNumber())
	}
	cli, err = NewClient(Mainnet(), FromEnvs(), WithMaxConnectionsNumber(defaultMaxConnectionsNumber+1))
	if err != nil {
		t.Fatalf("Unable to create tongo client: %v", err)
	}
	if cli.pool.ConnectionsNumber() != defaultMaxConnectionsNumber+1 {
		t.Fatalf("want connections number: %v, got: %v", defaultMaxConnectionsNumber+1, cli.pool.ConnectionsNumber())
	}
}

func TestAsyncInitialization(t *testing.T) {
	accountId := ton.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")

	cli, err := NewClient(Mainnet(),
		WithMaxConnectionsNumber(10),
		WithAsyncConnectionsInit())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	if cli.pool.ConnectionsNumber() != 0 {
		t.Fatalf("expected 0 connections")
	}

	iterations := 0
	noConnections := 0
	for {
		_, _, err = cli.RunSmcMethod(context.Background(), accountId, "seqno", tlb.VmStack{})
		if errors.Is(err, pool.ErrNoConnections) {
			noConnections += 1
			fmt.Printf("No connections\n")
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Fatalf("Run smc error: %v", err)
		}
		if cli.pool.ConnectionsNumber() > 2 {
			iterations += 1
		}
		if iterations > 3 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if noConnections == 0 {
		t.Fatalf("expected no connections error")
	}
}

func TestSyncInitialization(t *testing.T) {
	cli, err := NewClient(Mainnet(), WithMaxConnectionsNumber(2))
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	if cli.pool.ConnectionsNumber() == 0 {
		t.Fatalf("0 connections")
	}

	iterations := 0
	for {
		accountId := ton.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
		_, _, err = cli.RunSmcMethod(context.Background(), accountId, "seqno", tlb.VmStack{})
		if errors.Is(err, pool.ErrNoConnections) {
			t.Fatalf("no connections error")
		}
		if err != nil {
			log.Fatalf("Run smc error: %v", err)
		}
		if cli.pool.ConnectionsNumber() > 1 {
			iterations += 1
		}
		if iterations > 3 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func TestGetTransactions_archive(t *testing.T) {
	if len(os.Getenv("ARCHIVE_NODES_CONFIG")) == 0 {
		t.Skip("ARCHIVE_NODES_CONFIG env is not set")
	}
	value := os.Getenv("ARCHIVE_NODES_CONFIG")
	servers, err := config.ParseLiteServersEnvVar(value)
	if err != nil {
		t.Fatalf("ParseLiteServersEnvVar() failed: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected servers length: 2, got: %v", len(servers))
	}
	tongoClient, err := NewClient(WithLiteServers(servers), WithDetectArchiveNodes())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	time.Sleep(15 * time.Second)
	accountId, _ := ton.AccountIDFromRaw("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d220")
	txs, err := tongoClient.GetLastTransactions(context.Background(), accountId, 1000)
	if err != nil {
		t.Fatalf("Get transaction error: %v", err)
	}
	fmt.Printf("archive txs: %v\n", len(txs))
}

func TestGetTransactions(t *testing.T) {
	if len(os.Getenv("LITE_SERVERS")) == 0 {
		t.Skip("LITE_SERVERS env is not set")
	}
	tongoClient, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := ton.AccountIDFromRaw("0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb")
	for i := 1; i < 77; i++ {
		txs, err := tongoClient.GetLastTransactions(context.Background(), accountId, i)
		if err != nil {
			t.Fatalf("Get transaction error: %v", err)
		}
		if len(txs) != i {
			t.Fatalf("expected #txs: %v, got: %v", i, len(txs))
		}
		hashes := make(map[string]struct{}, len(txs))
		for i, tx := range txs {
			if i > 0 {
				if txs[i-1].Lt <= tx.Lt {
					log.Fatalf("wrong order")
				}
			}
			s := tx.Hash().Hex()
			if _, ok := hashes[s]; ok {
				log.Fatalf("duplicated hash: %v", s)
			}
			hashes[s] = struct{}{}
		}
	}
}

func TestSendRawMessage(t *testing.T) {
	t.Skip() //TODO: generate new valid message
	b, _ := hex.DecodeString("b5ee9c72010204010001700003e1880111b05b70f10022319f670ac91fa98660b3dc71a88892adbce0efcedfb15bc366119fdfc5395c5eb526485a4fa810c3d487ef036f3f8712ef3cce5c77e108fb9b6913d7f8a335a3e9a5ddee7e9ac4fa9da1be58490a5738293a1999ce6eab482de185353462ffffffffe0000000105001020300deff0020dd2082014c97ba218201339cbab19f71b0ed44d0d31fd31f31d70bffe304e0a4f2608308d71820d31fd31fd31ff82313bbf263ed44d0d31fd31fd3ffd15132baf2a15144baf2a204f901541055f910f2a3f8009320d74a96d307d402fb00e8d101a4c8cb1fcb1fcbffc9ed5400500000000029a9a317466f16a147b9b9db427d4e4763f455bc7c242757184ff564c421b371a41b705700ba62006707e00a47440d27444d3bedced2323ef6d64e68543c1736839c777d16e8309f2a098a678000000000000000000000000000000000000064636163363637332d656566342d343038662d623561652d346235363561323265643238")
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	code, err := tongoClient.SendMessage(ctx, b)
	if err != nil {
		log.Fatalf("Send message error: %v", err)
	}
	fmt.Printf("Send msg code: %v", code)
}

func TestRunSmcMethod(t *testing.T) {
	tongoClient, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId := ton.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
	_, _, err = tongoClient.RunSmcMethod(context.Background(), accountId, "seqno", tlb.VmStack{})
	if err != nil {
		log.Fatalf("Run smc error: %v", err)
	}
}

func TestGetAllShards(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	shards, err := api.GetAllShardsInfo(context.TODO(), info.Last.ToBlockIdExt())
	if err != nil {
		t.Fatal(err)
	}
	if len(shards) == 0 {
		t.Fatal("at least one shard should returns")
	}
	for _, s := range shards {
		fmt.Printf("Shard: %v\n", s.Shard)
	}
}

type BlockContent struct {
	Balances map[string]int64
}

func createOutputFile(api *Client, extID ton.BlockIDExt, filename string, accounts []tlb.Bits256) {
	content := BlockContent{
		Balances: map[string]int64{},
	}
	for _, accountAddr := range accounts {
		accountID := ton.AccountID{
			Workchain: extID.Workchain,
			Address:   accountAddr,
		}
		state, err := api.WithBlock(extID).GetAccountState(context.TODO(), accountID)
		if err != nil {
			panic(err)
		}
		collection, ok := state.Account.CurrencyCollection()
		if !ok {
			panic("failed to get currency collection from account state")
		}
		content.Balances[accountID.ToRaw()] = int64(collection.Grams)
	}
	bs, err := json.MarshalIndent(content, "", " ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filename, bs, 0644); err != nil {
		panic(err)
	}
}

func TestGetBlock(t *testing.T) {
	testCases := []struct {
		name    string
		blockID string
		file    string
	}{
		{
			name:    "block (0,4000000000000000,41731611)",
			blockID: "(0,4000000000000000,41731611)",
			file:    "block-41731611",
		},
		{
			name:    "block (-1,8000000000000000,34606335)",
			blockID: "(-1,8000000000000000,34606335)",
			file:    "block-34606335",
		},
		{
			name:    "block (0,8000000000000000,40429107)",
			blockID: "(0,8000000000000000,40429107)",
			file:    "block-40429107",
		},
		{
			name:    "block (0,c000000000000000,41745429)",
			blockID: "(0,c000000000000000,41745429)",
			file:    "block-41745429",
		},
	}
	api, err := NewClient(FromEnvs(), WithProofPolicy(ProofPolicyFast))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			blockID := ton.MustParseBlockID(tt.blockID)
			extID, _, err := api.LookupBlock(context.TODO(), blockID, 1, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			b, err := api.GetBlock(context.TODO(), extID)
			if err != nil {
				t.Fatal(err)
			}
			content := BlockContent{
				Balances: map[string]int64{},
			}
			balances := b.StateUpdate.ToRoot.AccountBalances()
			for _, tx := range b.AllTransactions() {
				_, ok := balances[tx.AccountAddr]
				if !ok {
					t.Fatalf("tx account not found in balances")
				}
			}
			for accountAddr, currencyCollection := range balances {
				accountID := ton.AccountID{
					Workchain: blockID.Workchain,
					Address:   accountAddr,
				}
				content.Balances[accountID.ToRaw()] = int64(currencyCollection.Grams)
			}
			outputFilename := fmt.Sprintf("testdata/%v.output.json", tt.file)
			bs, err := json.MarshalIndent(content, "", " ")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(outputFilename, bs, 0644); err != nil {
				t.Fatal(err)
			}
			expectedFilename := fmt.Sprintf("testdata/%v.expected.json", tt.file)
			if _, err := os.Stat(expectedFilename); errors.Is(err, os.ErrNotExist) {
				createOutputFile(api, extID, expectedFilename, maps.Keys(balances))
			}
			expected, err := os.ReadFile(expectedFilename)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(bytes.Trim(expected, " \n"), bytes.Trim(bs, " \n")) != 0 {
				t.Fatalf("block content mismatch")
			}
		})
	}
}

func TestGetConfigAll(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	_, err = api.GetConfigAll(context.TODO(), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConfigAllWithSafePolicy(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvsOrMainnet(), WithProofPolicy(ProofPolicyFast))
	if err != nil {
		t.Fatal(err)
	}
	_, err = api.GetConfigAll(context.TODO(), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAccountState(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name      string
		accountID string
	}{
		{
			name:      "account from masterchain",
			accountID: "-1:34517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf",
		},
		{
			name:      "account from basechain",
			accountID: "0:5f00decb7da51881764dc3959cec60609045f6ca1b89e646bde49d492705d77f",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accountID, err := ton.AccountIDFromRaw(tt.accountID)
			if err != nil {
				t.Fatal("AccountIDFromRaw() failed: %w", err)
			}
			st, err := api.GetAccountState(context.TODO(), accountID)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Account status: %v\n", st.Account.Status())
		})
	}
}

func TestLookupBlock(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Current block seqno : %v\n", info.Last.Seqno)
	blockID := ton.BlockID{
		Workchain: int32(info.Last.Workchain),
		Shard:     info.Last.Shard,
		Seqno:     info.Last.Seqno - 1,
	}
	_, bl, err := api.LookupBlock(context.TODO(), blockID, 1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Prev block seqno    : %v\n", bl.SeqNo)
}

func TestGetOneTransaction(t *testing.T) {
	tongoClient, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	ctx := context.Background()
	info, err := tongoClient.GetMasterchainInfo(ctx)
	shards, err := tongoClient.GetAllShardsInfo(ctx, info.Last.ToBlockIdExt())
	if err != nil {
		t.Fatal(err)
	}
	blockID := shards[0]
	var tx1 *tlb.Transaction
	for i := 0; i < 10; i++ {
		b, err := tongoClient.GetBlock(ctx, blockID)
		if err != nil {
			t.Fatal(err)
		}
		if len(b.AllTransactions()) == 0 {
			prev := b.Info.PrevRef.PrevBlkInfo.Prev
			blockID = ton.BlockIDExt{ton.BlockID{
				blockID.Workchain, blockID.Shard, prev.SeqNo,
			},
				ton.Bits256(prev.RootHash), ton.Bits256(prev.FileHash),
			}
			continue
		}
		tx1 = b.AllTransactions()[0]
		break
	}

	tx2, err := tongoClient.GetOneTransactionFromBlock(context.Background(), ton.AccountID{Workchain: blockID.Workchain, Address: tx1.AccountAddr}, blockID, tx1.Lt)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
	}
	if tx2.Hash() != tx1.Hash() {
		log.Fatalf("mismatch hashes")
	}
}

func TestGetLibraries(t *testing.T) {
	tongoClient, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	hash := ton.MustParseHash("587CC789EFF1C84F46EC3797E45FC809A14FF5AE24F1E0C7A6A99CC9DC9061FF")
	libs, err := tongoClient.GetLibraries(context.Background(), []ton.Bits256{hash})
	if err != nil {
		log.Fatalf("GetLibraries() failed: %v", err)
	}
	if len(libs) != 1 {
		t.Fatalf("expected libs lengths: 1, got: %v", len(libs))
	}
	cell, ok := libs[hash]
	if !ok {
		t.Fatalf("expected lib is not found")
	}
	base64, err := cell.ToBocBase64()
	if err != nil {
		t.Fatalf("ToBocBase64() failed: %v", err)
	}
	expected := "te6ccgEBAQEAXwAAuv8AIN0gggFMl7ohggEznLqxnHGw7UTQ0x/XC//jBOCk8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVA=="
	if base64 != expected {
		t.Fatalf("want: %v, got: %v", expected, base64)
	}
}

func TestGetJettonWallet(t *testing.T) {
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	master := ton.MustParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	owner := ton.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
	wallet, err := tongoClient.GetJettonWallet(context.Background(), master, owner)
	if err != nil {
		log.Fatalf("get jetton wallet error: %v", err)
	}
	fmt.Printf("jetton wallet address: %v\n", wallet.String())
}

func TestGetJettonData(t *testing.T) {
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	master := ton.MustParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	meta, err := tongoClient.GetJettonData(context.Background(), master)
	if err != nil {
		log.Fatalf("get jetton decimals error: %v", err)
	}
	fmt.Printf("jetton symbol: %v\n", meta.Symbol)
}

func TestGetJettonBalance(t *testing.T) {
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	jettonWallet := ton.MustParseAccountID("kQCOSEttz9aEGXkjd1h_NJsQqOca3T-Pld5zSIPHcYZIxsyf")
	b, err := tongoClient.GetJettonBalance(context.Background(), jettonWallet)
	if err != nil {
		log.Fatalf("get jetton decimals error: %v", err)
	}
	fmt.Printf("jetton balance: %v\n", b)
}

func TestDnsResolve(t *testing.T) {
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	root := ton.MustParseAccountID("Ef_BimcWrQ5pmAWfRqfeVHUCNV8XgsLqeAMBivKryXrghFW3")
	m, _, err := tongoClient.DnsResolve(context.Background(), root, "ton\u0000alice\u0000", big.NewInt(0))
	if err != nil {
		log.Fatalf("dns resolve error: %v", err)
	}
	fmt.Printf("Bytes resolved: %v\n", m)
}

func TestGetRootDNS(t *testing.T) {
	tongoClient, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	root, err := tongoClient.GetRootDNS(context.Background())
	if err != nil {
		log.Fatalf("get root dns error: %v", err)
	}
	fmt.Printf("Root DNS: %v\n", root.ToRaw())
}

func TestClient_GetTransactionsForUnknownAccount(t *testing.T) {
	var a ton.AccountID
	rand.Read(a.Address[:])
	client, err := NewClientWithDefaultTestnet()
	if err != nil {
		t.Error(err)
	}

	txs, err := client.GetTransactions(context.Background(), 10, a, 0, ton.Bits256{})
	if err != nil || len(txs) != 0 {
		t.Error(err, len(txs))
	}
}

func TestMappingTransactionsToBlocks(t *testing.T) {
	const limit = 100
	c, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	txs, err := c.GetLastTransactions(ctx, ton.MustParseAccountID("0:408da3b28b6c065a593e10391269baaa9c5f8caebc0c69d9f0aabbab2a99256b"), limit)
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != limit {
		t.Fatal(len(txs))
	}

	for i := 0; i < len(txs)-1; i++ {
		if txs[i].PrevTransLt != txs[i+1].Lt {
			t.Fatal("something with tx order")
		}
	}

	txsInBlockCache := make(map[ton.Bits256][]tlb.Bits256)
	for _, tx := range txs {
		if _, ok := txsInBlockCache[tx.BlockID.RootHash]; !ok {
			block, err := c.GetBlock(ctx, tx.BlockID)
			if err != nil {
				t.Fatal(err)
			}
			var txHashes []tlb.Bits256
			for _, accountBlock := range block.Extra.AccountBlocks.Values() {
				for _, txRef := range accountBlock.Transactions.Values() {
					txHashes = append(txHashes, txRef.Value.Hash())
				}
			}
			txsInBlockCache[tx.BlockID.RootHash] = txHashes
		}
		found := false
		for _, hash := range txsInBlockCache[tx.BlockID.RootHash] {
			if hash == tx.Hash() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("can't find tx %v in block %v", tx.Hash(), tx.BlockID.String())
		}
	}
}

func TestFromEnvs(t *testing.T) {

	os.Setenv("LITE_SERVERS", "some-value")
	options := Options{}
	err := FromEnvs()(&options)
	if err == nil {
		t.Fatal("expected err")
	}

	os.Setenv("LITE_SERVERS", "127.0.0.1:22095:6PGkPQSbyFp12esf1+Mp5+cAx5wtTU=")
	options = Options{}
	err = FromEnvs()(&options)
	if err != nil {
		t.Fatalf("FromEnv()() failed: %v", err)
	}
	if len(options.LiteServers) != 1 {
		t.Fatal("expected 1 lite server")
	}

	os.Unsetenv("LITE_SERVERS")
	options = Options{}
	err = FromEnvs()(&options)
	if err != nil {
		t.Fatalf("FromEnv()() failed: %v", err)
	}
	if len(options.LiteServers) != 0 {
		t.Fatal("expected 0 lite server")
	}
}

func TestWaitMasterchainBlock(t *testing.T) {
	api, err := NewClient(Mainnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Current block seqno : %v\n", info.Last.Seqno)
	nextSeqno := info.Last.Seqno + 1
	bl, err := api.WaitMasterchainBlock(context.TODO(), nextSeqno, time.Second*15)
	if err != nil {
		t.Fatal(err)
	}
	if bl.Seqno != nextSeqno {
		t.Fatal("wrong block seqno")
	}
	fmt.Printf("Next block seqno    : %v\n", bl.Seqno)
}
