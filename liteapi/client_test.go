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
	accountId, _ := tongo.AccountIDFromRaw("0:E2D41ED396A9F1BA03839D63C5650FAFC6FCFB574FD03F2E67D6555B61A3ACD9")
	var lt uint64 = 28563297000010
	var hash tongo.Hash
	_ = hash.FromHex("3E55B1BB7B6DD1603AB950A783890C3D1E945D0FD6BC29CF1C0017C44AC91E5E")
	_, err = tongoClient.GetTransactions(context.Background(), 100, *accountId, lt, hash)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
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
	err = tongoClient.SendRawMessage(ctx, b)
	if err != nil {
		log.Fatalf("Send message error: %v", err)
	}
}

func TestRunSmcMethod(t *testing.T) {
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.ParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
	_, _, err = tongoClient.RunSmcMethod(context.Background(), 4, *accountId, "seqno", tongo.VmStack{})
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
	shards, err := api.BlocksGetShards(context.TODO(), info)
	if err != nil {
		t.Fatal(err)
	}
	if len(shards) == 0 {
		t.Fatal("at least one shard should returns")
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
	_, block, err := api.GetBlock(context.TODO(), info)
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

func TestGetShardAccount(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	accountID, _ := tongo.AccountIDFromRaw("0:5f00decb7da51881764dc3959cec60609045f6ca1b89e646bde49d492705d77f")
	acc, err := api.GetLastShardAccount(context.TODO(), *accountID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Last TX LT: %v\n", acc.LastTransLt)
}

func TestGetLastConfigAll(t *testing.T) {
	api, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	_, err = api.GetLastConfigAll(context.TODO())
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
	st, err := api.GetAccountState(context.TODO(), *accountID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Account status: %v\n", st.Status)
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
	fmt.Printf("Current block seqno : %v\n", info.Seqno)
	blockID := tongo.TonNodeBlockId{Workchain: info.Workchain, Shard: info.Shard, Seqno: info.Seqno - 1}
	bl, _, err := api.LookupBlock(context.TODO(), 1, blockID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Prev block seqno    : %v\n", bl.Seqno)
}

func TestGetOneTransaction(t *testing.T) {
	tongoClient, err := NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("-1:34517C7BDF5187C55AF4F8B61FDC321588C7AB768DEE24B006DF29106458D7CF")
	var lt uint64 = 32174166000001
	var rh, fh tongo.Hash
	_ = fh.FromUnknownString("4237A19CD87265EB16ADDD0EF2EADF07B8F362903386653ED3ABDCF34477F748")
	_ = rh.FromUnknownString("9F66A2B2CE8F80762F39E0F5340F873192C8D416E9B82C015918E98C851A23C2")
	blockID := tongo.TonNodeBlockIdExt{TonNodeBlockId: tongo.TonNodeBlockId{Workchain: -1, Shard: -9223372036854775808, Seqno: 24425406}, RootHash: rh, FileHash: fh}
	_, _, err = tongoClient.GetOneRawTransaction(context.Background(), blockID, *accountId, lt)
	if err != nil {
		log.Fatalf("Get transaction error: %v", err)
	}
}

func TestGetJettonWallet(t *testing.T) {
	tongoClient, err := NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	master, _ := tongo.ParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	owner, _ := tongo.ParseAccountID("EQAs87W4yJHlF8mt29ocA4agnMrLsOP69jC1HPyBUjJay-7l")
	wallet, err := tongoClient.GetJettonWallet(context.Background(), *master, *owner)
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
	master, _ := tongo.ParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	meta, err := tongoClient.GetJettonData(context.Background(), *master)
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
	jettonWallet, _ := tongo.ParseAccountID("kQCOSEttz9aEGXkjd1h_NJsQqOca3T-Pld5zSIPHcYZIxsyf")
	b, err := tongoClient.GetJettonBalance(context.Background(), *jettonWallet)
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
	root, _ := tongo.ParseAccountID("Ef_BimcWrQ5pmAWfRqfeVHUCNV8XgsLqeAMBivKryXrghFW3")
	m, _, err := tongoClient.DnsResolve(context.Background(), *root, "ton\u0000alice\u0000", big.NewInt(0))
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
