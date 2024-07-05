package wallet

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
)

func TestGetW5ExtensionsList(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		wantErr        bool
		wantExtensions map[ton.AccountID]struct{}
	}{
		{
			name:      "wallet v5 with extensions",
			accountID: "0:8a189456e840670f5c09c5db75c829454cc6b1e5dd81cc20065fdb1999a9cbad",
			wantExtensions: map[ton.AccountID]struct{}{
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d220"): {},
				ton.MustParseAccountID("0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb"): {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := ton.MustParseAccountID(tt.accountID)
			cli, err := liteapi.NewClient(liteapi.Testnet())
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			state, err := cli.GetAccountState(context.Background(), accountID)
			if err != nil {
				t.Fatalf("GetAccountState() error = %v", err)
			}
			extensions, err := GetW5BetaExtensionsList(state)
			if err != nil {
				t.Fatalf("GetW5BetaExtensionsList() error = %v", err)
			}
			if !reflect.DeepEqual(extensions, tt.wantExtensions) {
				t.Errorf("GetW5BetaExtensionsList() = %v, want %v", extensions, tt.wantExtensions)
			}
		})
	}
}

func TestNewWalletV5R1(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want *walletV5R1
	}{
		{
			name: "workchain 0, testnet",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: &walletV5R1{
				workchain: 0,
				walletID:  2147483645,
			},
		},
		{
			name: "workchain 0, mainnet",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: &walletV5R1{
				workchain: 0,
				walletID:  2147483409,
			},
		},
		{
			name: "workchain -1, mainnet",
			opts: []Option{
				WithWorkchain(-1),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: &walletV5R1{
				workchain: -1,
				walletID:  8388369,
			},
		},
		{
			name: "workchain -1, testnet",
			opts: []Option{
				WithWorkchain(-1),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: &walletV5R1{
				workchain: -1,
				walletID:  8388605,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wallet := NewWalletV5R1(nil, applyOptions(tt.opts...))
			if wallet.walletID != tt.want.walletID {
				t.Errorf("NewWalletV5R1() = %v, want %v", wallet.walletID, tt.want.walletID)
			}
			if wallet.workchain != tt.want.workchain {
				t.Errorf("NewWalletV5R1() = %v, want %v", wallet.workchain, tt.want.workchain)
			}
		})
	}
}

func Test_walletV5R1_generateAddress(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		opts       []Option
		want       ton.AccountID
	}{
		{
			name:       "workchain 0, testnet",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: ton.MustParseAccountID("0:aa7bd5aa1614bc01f5460cfdb14224cb8db7a89612c6b2f21b6e043a0b75d3e6"),
		},
		{
			name:       "workchain 0, mainnet",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: ton.MustParseAccountID("0:827137ba7a1ad871a8a8605e8dba799666abb952dfe9eff6e9dfa96700ae16f4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := hex.DecodeString(tt.privateKey)
			if err != nil {
				t.Fatalf("hex.DecodeString() error = %v", err)
			}

			publicKey := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
			w := NewWalletV5R1(publicKey, applyOptions(tt.opts...))
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

func TestGetW5R1ExtensionsList(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		wantExtensions map[ton.AccountID]struct{}
	}{
		{
			name:      "wallet v5 with extensions",
			accountID: "0:d7391407f03695b3af341b29ab3137a8c269ab091e0043b4301a15a034828e1c",
			wantExtensions: map[ton.AccountID]struct{}{
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d221"): {},
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d225"): {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := ton.MustParseAccountID(tt.accountID)
			cli, err := liteapi.NewClient(liteapi.Testnet())
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			state, err := cli.GetAccountState(context.Background(), accountID)
			if err != nil {
				t.Fatalf("GetAccountState() error = %v", err)
			}
			extensions, err := GetW5R1ExtensionsList(state, 0)
			if err != nil {
				t.Fatalf("GetW5BetaExtensionsList() error = %v", err)
			}
			if !reflect.DeepEqual(extensions, tt.wantExtensions) {
				t.Errorf("GetW5BetaExtensionsList() \ngot  = %v\nwant = %v\n", extensions, tt.wantExtensions)
			}
		})
	}
}
