package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/liteclient"
	"github.com/startfellows/tongo/wallet"
	"log"
)

func main() {
	recipientAddr, _ := tongo.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)

	w, err := wallet.NewWallet(privateKey, wallet.V4R2, 0, nil)
	if err != nil {
		log.Fatalf("Unable to create wallet: %v", err)
	}

	log.Printf("Wallet address: %v\n", w.GetAddress().ToRaw())

	tonTransfer := wallet.TonTransfer{
		Recipient: *recipientAddr,
		Amount:    10000,
		Comment:   "hello",
		Bounce:    false,
		Mode:      1,
	}

	client, err := liteclient.NewClient(nil)
	if err != nil {
		log.Fatalf("Unable to create lite client: %v", err)
	}

	// Need to set seqno. Get seqno method will be later.
	seqno := uint32(15)

	msg, err := w.GenerateTonTransferMessage(seqno, 0xFFFFFFFF, []wallet.TonTransfer{tonTransfer})
	if err != nil {
		log.Fatalf("Unable to generate transfer message: %v", err)
	}

	err = client.SendRawMessage(context.Background(), msg)
	if err != nil {
		log.Fatalf("Send message error: %v", err)
	}
	log.Printf("The message was sent successfully")
}
