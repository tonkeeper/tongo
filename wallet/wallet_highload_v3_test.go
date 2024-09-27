package wallet

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"testing"
	"time"
)

func Test_HighloadV3_Send(t *testing.T) {
	t.Skip("Only for manual testing")
	tests := []struct {
		name        string
		msgQty      int
		specialMode string
	}{
		{
			name:        "single message",
			msgQty:      1,
			specialMode: "",
		},
		{
			name:        "single message with init",
			msgQty:      1,
			specialMode: "with_init",
		},
		{
			name:        "one batch",
			msgQty:      100,
			specialMode: "",
		},
		{
			name:        "few batches",
			msgQty:      300,
			specialMode: "",
		},
		{
			name:        "single ext out message",
			msgQty:      1,
			specialMode: "ext_out",
		},
		{
			name:        "ext out message few batches",
			msgQty:      300,
			specialMode: "ext_out",
		},
	}

	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		t.Fatalf("Unable to create lite client: %v", err)
	}
	pk, _ := SeedToPrivateKey("birth pattern then forest walnut then phrase walnut fan pumpkin pattern then cluster blossom verify then forest velvet pond fiction pattern collect then then")
	opts := []Option{WithSubWalletID(DefaultSubWallet)}
	w, err := New(pk, HighLoadV3R1, client, opts...)
	if err != nil {
		t.Fatalf("Unable to create wallet: %v", err)
	}
	recipient := tongo.MustParseAccountID("kQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIX8f")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfers := make([]Sendable, tt.msgQty)
			for i := 0; i < tt.msgQty; i++ {
				switch tt.specialMode {
				case "with_init":
					transfers[i] = Message{
						Amount:  1_000_000,
						Address: recipient,
						Code:    boc.NewCell(),
						Data:    boc.NewCell(),
					}
				case "ext_out":
					transfers[i] = LogMessage{
						Comment: "123",
					}
				default:
					transfers[i] = SimpleTransfer{
						Amount:     100_000,
						Address:    recipient,
						Comment:    fmt.Sprintf("%d", i+1),
						Bounceable: false,
					}
				}
			}
			err = w.Send(context.TODO(), transfers...)
			if err != nil {
				t.Fatalf("Sending err: %v", err)
			}
			time.Sleep(time.Second)
		})
	}
}

func Test_HighloadV3_generateAddress(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		opts       []Option
		want       ton.AccountID
	}{
		{
			name:       "workchain 0",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithSubWalletID(DefaultSubWallet),
			},
			want: ton.MustParseAccountID("0:41ad70d17c024e9b4e2e3a5948d3f5ff855e339c6d9d504c679abc4eb08c4b7c"), // TODO: depends of DefaultTimeoutHighloadV3
		},
		{
			name:       "workchain -1",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(-1),
				WithSubWalletID(DefaultSubWallet),
			},
			want: ton.MustParseAccountID("-1:41ad70d17c024e9b4e2e3a5948d3f5ff855e339c6d9d504c679abc4eb08c4b7c"), // TODO: depends of DefaultTimeoutHighloadV3
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := hex.DecodeString(tt.privateKey)
			if err != nil {
				t.Fatalf("hex.DecodeString() error = %v", err)
			}
			publicKey := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
			w := newWalletHighloadV3(HighLoadV3R1, publicKey, applyOptions(tt.opts...))
			address, err := w.generateAddress()
			if err != nil {
				t.Fatalf("generateAddress() error = %v", err)
			}
			if address.ToRaw() != tt.want.ToRaw() {
				t.Errorf("generateAddress() got = %v, want %v", address, tt.want)
			}
		})
	}
}

func Test_DecodeHighloadV3Message(t *testing.T) {
	tests := []struct {
		name   string
		msgQty int
	}{
		{
			name:   "single message",
			msgQty: 1,
		},
		{
			name:   "one batch",
			msgQty: 200,
		},
		{
			name:   "few batches",
			msgQty: 600,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfers := make([]SimpleTransfer, tt.msgQty)
			raws := make([]RawMessage, tt.msgQty)
			for i := 0; i < tt.msgQty; i++ {
				transfers[i] = SimpleTransfer{
					Amount:     tlb.Grams(i * 100),
					Address:    tongo.MustParseAccountID("kQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIX8f"),
					Bounceable: false,
				}
				msg, _, _ := transfers[i].ToInternal()
				c := boc.NewCell()
				_ = tlb.Marshal(c, msg)
				raws[i].Message = c
				raws[i].Mode = DefaultMessageMode
			}
			pubKey, privKey, _ := ed25519.GenerateKey(nil)
			w, _ := newWallet(pubKey, HighLoadV3R1, Options{})
			signedBodyCell, err := w.createSignedMsgBodyCell(privKey, raws, MessageConfig{})
			if err != nil {
				t.Fatalf("Unable to createSignedMsgBodyCell: %v", err)
			}
			addr, _ := w.generateAddress()
			extMsg, _ := ton.CreateExternalMessage(addr, signedBodyCell, nil, tlb.VarUInteger16{})
			extMsgCell := boc.NewCell()
			_ = tlb.Marshal(extMsgCell, extMsg)
			decodedMsg, err := DecodeHighloadV3Message(extMsgCell)
			if err != nil {
				t.Fatalf("Unable to docode external message: %v", err)
			}
			for i, m := range raws {
				if m.Mode != decodedMsg.Messages[i].Mode {
					t.Fatalf("Invalid mode at step: %d", i)
				}
				h1, _ := m.Message.Hash256()
				h2, _ := decodedMsg.Messages[i].Message.Hash256()
				if h1 != h2 {
					t.Fatalf("Invalid message hash at step: %d", i)
				}
			}
		})
	}
}
