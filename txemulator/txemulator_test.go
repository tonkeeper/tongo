package txemulator

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/liteclient"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/tvm"
	"github.com/startfellows/tongo/wallet"
	"log"
	"testing"
)

func TestExec(t *testing.T) {
	recipientAddr, _ := tongo.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)

	w, err := wallet.NewWallet(privateKey, wallet.V4R2, 0, nil)
	if err != nil {
		log.Fatalf("Unable to create wallet: %v", err)
	}
	tongoClient, err := liteclient.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	config, err := tongoClient.GetLastConfigAll(context.Background())
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	account, err := tongoClient.GetLastRawAccount(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	tonTransfer := wallet.TonTransfer{
		Recipient: *recipientAddr,
		Amount:    10000,
		Comment:   "hello",
		Bounce:    false,
		Mode:      1,
	}

	accountID := w.GetAddress()
	res, err := tvm.RunTvm(
		&account.Account.Storage.State.AccountActive.StateInit.Code.Value.Value,
		&account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value,
		"seqno", []tvm.StackEntry{}, &accountID)
	if err != nil {
		log.Fatalf("TVM run error: %v", err)
	}
	if res.ExitCode != 0 || len(res.Stack) != 1 || !res.Stack[0].IsInt() {
		log.Fatalf("TVM execution failed")
	}

	msg, err := w.GenerateTonTransferMessage(uint32(res.Stack[0].Int64()), 0xFFFFFFFF, []wallet.TonTransfer{tonTransfer})
	if err != nil {
		log.Fatalf("Unable to generate transfer message: %v", err)
	}

	var message tongo.Message[tlb.Any]
	c, err := boc.DeserializeBoc(msg)
	if err != nil {
		log.Fatalf("unable to deserialize transfer message: %v", err)
	}
	err = tlb.Unmarshal(c[0], &message)
	if err != nil {
		log.Fatalf("unable to unmarshal transfer message: %v", err)
	}

	var shardAccount tongo.ShardAccount
	shardAccount.Account.Value = account
	shardAccount.LastTransLt = account.Account.Storage.LastTransLt - 1

	e, err := NewEmulator(config)
	if err != nil {
		log.Fatalf("unable to create emulator: %v", err)
	}

	acc, tx, err := e.Emulate(shardAccount, message)
	if err != nil {
		log.Fatalf("emulator error: %v", err)
	}
	fmt.Printf("Account last transaction hash: %x\n", acc.LastTransHash)
	fmt.Printf("Prev transaction hash: %x\n", tx.Transaction.PrevTransHash)

}
