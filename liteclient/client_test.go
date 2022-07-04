package liteclient

import (
	"context"
	"encoding/base64"
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

func TestSendMessage(t *testing.T) {
	tongoConfig, err := config.ParseConfigFile("tmp/global.config.json")
	if err != nil {
		log.Fatalf("Unable to read config json: %v", err)
	}
	tongoClient, err := NewClient(*tongoConfig)
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	payload, _ := base64.StdEncoding.DecodeString("qrvM3Q==")
	err = tongoClient.SendMessage(context.Background(), payload)
	if err != nil {
		log.Fatalf("Send message error: %v", err)
	}
}
