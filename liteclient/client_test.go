package liteclient

import (
	"context"
	"encoding/hex"
	"github.com/startfellows/tongo"
	"log"
	"testing"
)

func TestGetTransactions(t *testing.T) {
	tongoClient, err := NewClient(nil)
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
	b, _ := hex.DecodeString("b5ee9c72010204010001700003e1880111b05b70f10022319f670ac91fa98660b3dc71a88892adbce0efcedfb15bc366119fdfc5395c5eb526485a4fa810c3d487ef036f3f8712ef3cce5c77e108fb9b6913d7f8a335a3e9a5ddee7e9ac4fa9da1be58490a5738293a1999ce6eab482de185353462ffffffffe0000000105001020300deff0020dd2082014c97ba218201339cbab19f71b0ed44d0d31fd31f31d70bffe304e0a4f2608308d71820d31fd31fd31ff82313bbf263ed44d0d31fd31fd3ffd15132baf2a15144baf2a204f901541055f910f2a3f8009320d74a96d307d402fb00e8d101a4c8cb1fcb1fcbffc9ed5400500000000029a9a317466f16a147b9b9db427d4e4763f455bc7c242757184ff564c421b371a41b705700ba62006707e00a47440d27444d3bedced2323ef6d64e68543c1736839c777d16e8309f2a098a678000000000000000000000000000000000000064636163363637332d656566342d343038662d623561652d346235363561323265643238")
	tongoClient, err := NewClient(nil)
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	err = tongoClient.SendRawMessage(context.Background(), b)
	if err != nil {
		log.Fatalf("Send message error: %v", err)
	}
}

func TestRunSmcMethod(t *testing.T) {
	tongoClient, err := NewClient(nil)
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("0:deaae6518a11fd24c1da9c53ad38aedd35f4a66d1bef4f1e3081472d9276a920")
	_, err = tongoClient.RunSmcMethod(context.Background(), 4, *accountId, "seqno", tongo.VmStack{})
	if err != nil {
		log.Fatalf("Run smc error: %v", err)
	}
}
