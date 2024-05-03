package wallet

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tontest"
)

func TestGetCodeByVer(t *testing.T) {
	for ver := V1R1; ver < HighLoadV2; ver++ {
		_ = GetCodeByVer(ver)
	}
}

func TestVersionToString(t *testing.T) {
	testData := map[Version]string{
		V1R1:       "v1R1",
		V3R1:       "v3R1",
		V3R2:       "v3R2",
		V4R1:       "v4R1",
		V4R2:       "v4R2",
		HighLoadV2: "highload_v2",
	}
	for ver, name := range testData {
		if ver.ToString() != name {
			t.Fatal("invalid mapping version to string")
		}
	}
}

func TestGenerateWalletAddress(t *testing.T) {
	type walletData struct {
		Address   string
		PublicKey string
	}
	testData := map[Version]walletData{
		// TODO: add other versions
		V3R2: {"0:f3a069b7fc4631da4401de03eddd7cd30caca618c6ad0e3ac3fa454370b73a96",
			"f96db56e72de2e84e0aef780428e439a6c84e0b27bc2b2591075785479f2e9c3"},
		V4R1: {"0:17afeaaa61cb575e3e340a296da6bf55bc6b996cfab1d9f87840b2b6dc4cf613",
			"6f58b9fecb87e847825a7ecf3ae1f32b5578eee156ac10b398e2f1d67c12ca05"},
		V4R2: {"0:8f2983152d1480ba6af25e087d672232080b294dc8992525e35e4ff6d601f405",
			"7843fd9de6cd858154d9a914b8c3cd0bf1dc5af3a0c1dd273586568fc4d1c002"},
	}
	for ver, data := range testData {
		key, _ := hex.DecodeString(data.PublicKey)
		publicKey := ed25519.PublicKey(key)
		address, err := GenerateWalletAddress(publicKey, ver, 0, nil, nil)
		if err != nil {
			t.Fatalf("address generation failed: %v", err)
		}
		if address.ToRaw() != data.Address {
			t.Fatal("addresses mismatch")
		}
	}
}

func TestLongCommentSerialization(t *testing.T) {
	// TODO: add real serialized data
	longComment := `
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog`

	cell := boc.NewCell()
	err := tlb.Marshal(cell, TextComment(longComment))
	if err != nil {
		t.Fatalf("long comment serialization error: %v", err)
	}
	var text TextComment
	err = tlb.Unmarshal(cell, &text)
	if err != nil {
		t.Fatalf("long comment deserialization error: %v", err)
	}
	if string(text) != longComment {
		t.Fatal("TextComment invalid serialization/deserialization")
	}
}

func TestSimpleSend(t *testing.T) {
	t.Skip()
	recipientAddr, _ := ton.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	w := initDefaultWallet(client)
	tonTransfer := SimpleTransfer{
		Amount:  10000,
		Address: recipientAddr,
		Comment: "hello",
	}
	err = w.Send(context.Background(), tonTransfer)
	if err != nil {
		t.Fatalf("Unable to generate transfer message: %v", err)
	}
}

func TestGetSeqno(t *testing.T) {
	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	w := initDefaultWallet(client)
	seqno, err := client.GetSeqno(context.Background(), w.GetAddress())
	if err != nil {
		t.Fatalf("Unable to get wallet seqno: %v", err)
	}
	fmt.Printf("Seqno: %v\n", seqno)
}

func TestMockBlockchain(t *testing.T) {
	recipientAddr, _ := ton.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	client, c := NewMockBlockchain(1, tontest.Account().Balance(10000).Address(ton.AccountID{}).MustShardAccount())
	w := initDefaultWallet(client)
	tonTransfer := SimpleTransfer{
		Amount:  10000,
		Address: recipientAddr,
		Comment: "hello",
		// Body:    *boc.Cell, // empty
		// Init:    *tongo.StateInit, // empty
		// Bounce: *bool, // default
		// Mode:       *byte, // default
	}
	err := w.Send(context.Background(), tonTransfer)
	if err != nil {
		t.Fatalf("Unable to send message: %v", err)
	}
	res := <-c
	fmt.Printf("Transfer message: %x\n", res)
	b, _ := w.GetBalance(context.Background())
	fmt.Printf("Wallet balance: %v\n", b)
}

func initDefaultWallet(blockchain blockchain) Wallet {
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)
	w, err := New(privateKey, V4R2, 0, nil, blockchain)
	if err != nil {
		panic("unable to create wallet")
	}
	fmt.Printf("Wallet address: %v\n", w.GetAddress())
	return w
}

func TestDeserializeMessage(t *testing.T) {
	raw := "te6ccgECBgEAAXgAAZyTKd6JrxqOdZVwzAQwORTSmzEWraYc1v/B1trCjeQrl2o3reaq5d5fnMSdsStxpuQiWtclR30WuOqT4B5llBUBKamjF2NhSmsAAAAhAAMBAWpiAGyOLO7S4cYY8vrxnyGjmQy4fbFHsBTVru7/HlYASPPFKAlQL5AAAAAAAAAAAAAAAAAAAQIDtUY3KJoNoPRt29Q+0/k+b+3eMh/+pICK4Bdz8Cp7zpaXiLqHOGifBStWBqk5GdQn02DkhGQ+wOAtjn9w4gNcdEaIjZ0HAAAnE2NhShBjYUqmCXRlc3RlZWRlMcADBAUAWgFodHRwczovL25mdC5zdGVsLmNvbS91c2VybmFtZS90ZXN0ZWVkZTEuanNvbgBrgBCxVofvaKaMhDVkiORUzmkG5xwA0tZBmR3uDe2eRLI7agJUC+QAoCVAvkAAoAABwgAAACWQAEsABQBkgBCxVofvaKaMhDVkiORUzmkG5xwA0tZBmR3uDe2eRLI7cA=="
	cells, err := boc.DeserializeBocBase64(raw)
	if err != nil {
		panic(err)
	}
	var msg MessageV4
	cells[0].Skip(512)
	err = tlb.Unmarshal(cells[0], &msg)
	if err != nil {
		panic(err)
	}
}
func pointer[T any](t T) *T {
	return &t
}

func TestGenerateWalletAddress1(t *testing.T) {
	tests := []struct {
		name            string
		key             ed25519.PublicKey
		ver             Version
		workchain       int
		subWalletId     *int
		networkGlobalID *int
		want            ton.AccountID
		wantErr         bool
	}{
		{
			name:            "V5R1",
			ver:             V5R1,
			workchain:       0,
			networkGlobalID: pointer(-3),
			key:             mustPubkeyFromHex("cfa50eeb1c3293c92bd33d5aa672c1717accd8a21b96033debb6d30b5bb230df"),
			want:            ton.MustParseAccountID("EQCsa9xhVJCw2BRL07dhxwkOoAjNHRPLN2iPggZG_ZauRYt-"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateWalletAddress(tt.key, tt.ver, tt.workchain, tt.subWalletId, tt.networkGlobalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateWalletAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateWalletAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
