package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"math/big"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/contract/jetton"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/wallet"
)

func main() {
	recipientAddr := tongo.MustParseAddress("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")

	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)
	w, err := wallet.New(privateKey, wallet.V4R2, 0, nil, client)
	if err != nil {
		panic("unable to create wallet")
	}

	master := tongo.MustParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	j := jetton.New(master, client)
	b, err := j.GetBalance(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Unable to get jetton wallet balance: %v", err)
	}
	d, err := j.GetDecimals(context.Background())
	if err != nil {
		log.Fatalf("Get decimals error: %v", err)
	}
	jw, err := j.GetJettonWallet(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Get jetton wallet error: %v", err)
	}

	log.Printf("Jetton balance: %v", b)
	log.Printf("Jetton decimals: %v", d)
	log.Printf("Jetton wallet owner address: %v", w.GetAddress().String())
	log.Printf("Jetton wallet address: %v", jw.String())

	amount := big.NewInt(1000)
	if amount.Cmp(b) == 1 {
		log.Fatalf("%v jettons needed, but only %v on balance", amount, b)
	}

	jettonTransfer := jetton.TransferMessage{
		Jetton:           j,
		JettonAmount:     amount,
		Destination:      recipientAddr.ID,
		AttachedTon:      400_000_000,
		ForwardTonAmount: 200_000_000,
	}
	err = w.Send(context.Background(), jettonTransfer)
	if err != nil {
		log.Fatalf("Unable to send transfer message: %v", err)
	}
	time.Sleep(time.Second * 15)
	b, err = j.GetBalance(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Unable to get jetton wallet balance: %v", err)
	}
	log.Printf("New Jetton balance: %v", b)
}
