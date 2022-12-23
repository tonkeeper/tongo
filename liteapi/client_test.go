package liteapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/startfellows/tongo"
)

func TestGetTransactions(t *testing.T) {
	t.Skip() //TODO: switch tests to archive node
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("-1:34517C7BDF5187C55AF4F8B61FDC321588C7AB768DEE24B006DF29106458D7CF")
	var lt uint64 = 33973842000003
	var hash tongo.Hash
	_ = hash.FromHex("8005AF92C0854B5A614427206673D120EA2914468C11C8F867F43740D6B4ACFB")
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
	_, _, err = tongoClient.RunSmcMethod(context.Background(), accountId, "seqno", tongo.VmStack{})
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
	p, err := block.Info.GetParents()
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
	accountID, _ := tongo.AccountIDFromRaw("0:5f00decb7da51881764dc3959cec60609045f6ca1b89e646bde49d492705d77f")
	st, err := api.GetAccountState(context.TODO(), accountID)
	if err != nil {
		t.Fatal(err)
	}
	ai, err := st.Account.GetInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Account status: %v\n", ai.Status)
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
	bl, err := api.LookupBlock(context.TODO(), blockID, 1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Prev block seqno    : %v\n", bl.SeqNo)
}

func TestGetOneTransaction(t *testing.T) {
	t.Skip() //todo: switch  to archive node
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("-1:34517C7BDF5187C55AF4F8B61FDC321588C7AB768DEE24B006DF29106458D7CF")
	var lt uint64 = 33973842000001
	var rh, fh tongo.Hash
	_ = fh.FromUnknownString("F497D5CE3DA3C2DAA217145A91A615188E5AD4D8D5EC58C86414DE3F627DFE8A")
	_ = rh.FromUnknownString("8215CADE3E7BAB4311230F35B5BAC218CFCB8B3706A21563556BCA29828206C9")
	blockID := tongo.BlockIDExt{BlockID: tongo.BlockID{Workchain: -1, Shard: uint64(9223372036854775808), Seqno: 26097165}, RootHash: rh, FileHash: fh}
	tx, err := tongoClient.WithBlock(blockID).GetOneTransaction(context.Background(), accountId, lt)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
	}
	fmt.Printf("TX utime: %v", tx.Now)
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
