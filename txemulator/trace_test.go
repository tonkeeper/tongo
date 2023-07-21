package txemulator

import (
	"context"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

const SEED = "way label strategy scheme park virtual walnut illegal fringe once state defense museum bone satoshi feel diary buddy notice solve moral maple video local"

func TestSimpleEmulation(t *testing.T) {
	ctx := context.Background()
	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		t.Fatal(err)
	}
	tracer, err := NewTraceBuilder(WithAccountsSource(client))
	if err != nil {
		t.Fatal(err)
	}
	w, err := wallet.DefaultWalletFromSeed(SEED, client)
	seqno := uint32(0)

	mock, messages := wallet.NewMockBlockchain(seqno, tongo.AccountInfo{Status: "active"})
	w, err = wallet.DefaultWalletFromSeed(SEED, mock)
	if err != nil {
		t.Fatal(err)
	}
	err = w.Send(ctx, wallet.SimpleTransfer{
		Amount:  tongo.OneTON / 10,
		Address: w.GetAddress(),
	}, wallet.SimpleTransfer{
		Amount:  tongo.OneTON / 10,
		Address: w.GetAddress(),
	})
	if err != nil {
		t.Fatal(err)
	}

	c, err := boc.DeserializeBoc(<-messages)
	if err != nil {
		t.Fatal(err)
	}
	var m tlb.Message
	err = tlb.Unmarshal(c[0], &m)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := tracer.Run(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Children) != 2 {
		t.Fatal(len(tree.Children))
	}
	if tree.Children[0].TX.Msgs.InMsg.Value.Value.Info.IntMsgInfo.Value.Grams != tongo.OneTON/10 {
		t.Fatal("invalid amount")
	}
}

func TestEmulate(t *testing.T) {
	// this message is for "EQBAF7OBsy_1R8Zs33l6XMP3k1OyMv6Nv-b_-n-qf7de9qp2", which uses a public library.
	c, err := boc.DeserializeSinglRootBase64("te6ccgEBAgEAoAABz4gAgC9nA2Zf6o+M2b7y9LmH7yanZGX9G3/N//T/VP9uvewComZfYno/fswnemt9B6xfHWRtZ2vKvL8C7ZiExKR3s3vsDDRnpxb5Oaoi7ATNea26glvtLlEwEFRoyIL2ZgqIaAAAAAgcAQBmYgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEBzEtAAAAAAAAAAAAAAAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	var m tlb.Message
	if err = tlb.Unmarshal(c, &m); err != nil {
		t.Fatal(err)
	}
	client, err := liteapi.NewClient(liteapi.Mainnet(), liteapi.FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	emulator, err := NewTraceBuilder(WithAccountsSource(client))
	if err != nil {
		t.Fatalf("NewTraceBuilder() failed: %v", err)
	}
	tree, err := emulator.Run(context.Background(), m)
	if err != nil {
		t.Fatalf("Run() failed: %v", err)
	}
	if !tree.TX.IsSuccess() {
		t.Fatalf("tx failed")
	}
	if len(tree.Children) != 1 {
		t.Fatalf("expected tx to has 1 child")
	}
	second := tree.Children[0].TX
	if !second.IsSuccess() {
		t.Fatalf("second tx failed")
	}
}
