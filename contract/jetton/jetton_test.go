package jetton

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

func initDefaultWallet(blockchain *liteapi.Client) wallet.Wallet {
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)
	w, err := wallet.New(privateKey, wallet.V4R2, 0, nil, blockchain)
	if err != nil {
		panic("unable to create wallet")
	}
	fmt.Printf("Wallet address: %v\n", w.GetAddress())
	return w
}

func TestSendJetton(t *testing.T) {
	t.Skip()
	recipientAddr, _ := ton.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")

	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	w := initDefaultWallet(client)

	master, _ := ton.ParseAccountID("kQCKt2WPGX-fh0cIAz38Ljd_OKQjoZE_cqk7QrYGsNP6wfP0")
	j := New(master, client)
	b, err := j.GetBalance(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Unable to get jetton wallet balance: %v", err)
	}
	amount := big.NewInt(1000)
	if amount.Cmp(b) == 1 {
		log.Fatalf("%v jettons needed, but only %v on balance", amount, b)
	}

	log.Printf("Prev balance: %v", b)
	jettonTransfer := TransferMessage{
		Jetton:           j,
		JettonAmount:     amount,
		Destination:      recipientAddr,
		AttachedTon:      ton.OneTON / 2,
		ForwardTonAmount: 200_000_000,
	}
	err = w.Send(context.Background(), jettonTransfer)
	if err != nil {
		t.Fatalf("Unable to send transfer message: %v", err)
	}
	time.Sleep(time.Second * 15)
	b, err = j.GetBalance(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Unable to get jetton wallet balance: %v", err)
	}
	log.Printf("New balance: %v", b)
}
