package wallet

import (
	"context"
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
			extensions, err := GetW5ExtensionsList(state)
			if err != nil {
				t.Fatalf("GetW5ExtensionsList() error = %v", err)
			}
			if !reflect.DeepEqual(extensions, tt.wantExtensions) {
				t.Errorf("GetW5ExtensionsList() = %v, want %v", extensions, tt.wantExtensions)
			}
		})
	}
}
