package liteapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/tlb"
)

func TestNewClient_WithMaxConnectionsNumber(t *testing.T) {
	cli, err := NewClient(Mainnet())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	if cli.pool.ConnectionsNumber() != defaultMaxConnectionsNumber {
		t.Fatalf("want connections number: %v, got: %v", defaultMaxConnectionsNumber, cli.pool.ConnectionsNumber())
	}
	cli, err = NewClient(Mainnet(), WithMaxConnectionsNumber(defaultMaxConnectionsNumber+1))
	if err != nil {
		t.Fatalf("Unable to create tongo client: %v", err)
	}
	if cli.pool.ConnectionsNumber() != defaultMaxConnectionsNumber+1 {
		t.Fatalf("want connections number: %v, got: %v", defaultMaxConnectionsNumber+1, cli.pool.ConnectionsNumber())
	}
}

func TestGetTransactions(t *testing.T) {
	t.Skip() //TODO: switch tests to archive node
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("-1:34517C7BDF5187C55AF4F8B61FDC321588C7AB768DEE24B006DF29106458D7CF")
	var lt uint64 = 33973842000003
	hash := tongo.MustParseHash("8005AF92C0854B5A614427206673D120EA2914468C11C8F867F43740D6B4ACFB")
	tx, err := tongoClient.GetTransactions(context.Background(), 100, accountId, lt, hash)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
	}
	fmt.Printf("Tx qty: %v\n", len(tx))
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
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId := tongo.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
	_, _, err = tongoClient.RunSmcMethod(context.Background(), accountId, "seqno", tlb.VmStack{})
	if err != nil {
		log.Fatalf("Run smc error: %v", err)
	}
}

func TestGetAllShards(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
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

func TestGetBlock(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	block, err := api.GetBlock(context.TODO(), info.Last.ToBlockIdExt())
	if err != nil {
		t.Fatal(err)
	}
	p, err := tongo.GetParents(block.Info)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Block seqno: %v\n", block.Info.SeqNo)
	fmt.Printf("1st parent block seqno: %v\n", p[0].Seqno)
}

func TestGetConfigAll(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	_, err = api.GetConfigAll(context.TODO(), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAccountState(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
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
			accountID, err := tongo.AccountIDFromRaw(tt.accountID)
			if err != nil {
				t.Fatal("AccountIDFromRaw() failed: %w", err)
			}
			st, err := api.GetAccountState(context.TODO(), accountID)
			if err != nil {
				t.Fatal(err)
			}
			ai, err := tongo.GetAccountInfo(st.Account)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Account status: %v\n", ai.Status)
		})
	}
}

func TestLookupBlock(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Current block seqno : %v\n", info.Last.Seqno)
	blockID := tongo.BlockID{
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
	tongoClient, err := NewClientWithDefaultMainnet()
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
		block, err := tongoClient.GetBlock(ctx, blockID)
		if err != nil {
			t.Fatal(err)
		}
		if len(block.AllTransactions()) == 0 {
			prev := block.Info.PrevRef.PrevBlkInfo.Prev
			blockID = tongo.BlockIDExt{tongo.BlockID{
				blockID.Workchain, blockID.Shard, prev.SeqNo,
			},
				tongo.Bits256(prev.RootHash), tongo.Bits256(prev.FileHash),
			}
			continue
		}
		tx1 = block.AllTransactions()[0]
		break
	}

	tx2, err := tongoClient.GetOneTransactionFromBlock(context.Background(), tongo.AccountID{Workchain: blockID.Workchain, Address: tx1.AccountAddr}, blockID, tx1.Lt)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
	}
	if tx2.Hash() != tx1.Hash() {
		log.Fatalf("mismatch hashes")
	}
}

func TestGetLibraries(t *testing.T) {
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	hash := tongo.MustParseHash("587CC789EFF1C84F46EC3797E45FC809A14FF5AE24F1E0C7A6A99CC9DC9061FF")
	libs, err := tongoClient.GetLibraries(context.Background(), []tongo.Bits256{hash})
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
	master := tongo.MustParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	owner := tongo.MustParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
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
	master := tongo.MustParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
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
	jettonWallet := tongo.MustParseAccountID("kQCOSEttz9aEGXkjd1h_NJsQqOca3T-Pld5zSIPHcYZIxsyf")
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
	root := tongo.MustParseAccountID("Ef_BimcWrQ5pmAWfRqfeVHUCNV8XgsLqeAMBivKryXrghFW3")
	m, _, err := tongoClient.DnsResolve(context.Background(), root, "ton\u0000alice\u0000", big.NewInt(0))
	if err != nil {
		log.Fatalf("dns resolve error: %v", err)
	}
	fmt.Printf("Bytes resolved: %v\n", m)
}

func TestGetRootDNS(t *testing.T) {
	tongoClient, err := NewClientWithDefaultMainnet()
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
	var a tongo.AccountID
	rand.Read(a.Address[:])
	client, err := NewClientWithDefaultTestnet()
	if err != nil {
		t.Error(err)
	}

	txs, err := client.GetTransactions(context.Background(), 10, a, 0, tongo.Bits256{})
	if err != nil || len(txs) != 0 {
		t.Error(err, len(txs))
	}
}

func TestMappingTransactionsToBlocks(t *testing.T) {
	const limit = 100
	c, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	txs, err := c.GetLastTransactions(ctx, tongo.MustParseAccountID("0:408da3b28b6c065a593e10391269baaa9c5f8caebc0c69d9f0aabbab2a99256b"), limit)
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

	txsInBlockCache := make(map[tongo.Bits256][]tlb.Bits256)
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
