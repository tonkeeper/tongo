package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/liteapi"
	"github.com/startfellows/tongo/wallet"
	"log"
)

func main() {
	recipientAddr, _ := tongo.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)

	client, err := liteclient.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create lite client: %v", err)
	}

	w, err := wallet.NewWallet(privateKey, wallet.V4R2, 0, nil, client)
	if err != nil {
		log.Fatalf("Unable to create wallet: %v", err)
	}

	log.Printf("Wallet address: %v\n", w.GetAddress().ToRaw())

	comment := "hello"
	tonTransfer := wallet.Message{
		Amount:  10000,
		Address: *recipientAddr,
		Comment: &comment,
		// Body:    *boc.Cell, // empty
		// Init:    *tongo.StateInit, // empty
		// Bounceable: *bool, // default
		// Mode:       *byte, // default
	}

	err = w.SimpleSend(context.Background(), []wallet.Message{tonTransfer})
	if err != nil {
		log.Fatalf("Unable to generate transfer message: %v", err)
	}
	log.Printf("The message was sent successfully")
}
