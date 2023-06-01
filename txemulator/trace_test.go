package txemulator

import (
	"context"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
	"testing"
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
