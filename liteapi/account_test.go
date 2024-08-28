package liteapi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/ton"
	"testing"
)

func TestGetAccountWithProof(t *testing.T) {
	api, err := NewClient(Testnet(), FromEnvs())
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name      string
		accountID string
	}{
		{
			name:      "account from masterchain",
			accountID: "-1:34517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf",
		},
		{
			name:      "active account from basechain",
			accountID: "0:e33ed33a42eb2032059f97d90c706f8400bb256d32139ca707f1564ad699c7dd",
		},
		{
			name:      "nonexisted from basechain",
			accountID: "0:5f00decb7da51881764dc3959cec60609045f6ca1b89e646bde49d492705d77c",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accountID, err := ton.AccountIDFromRaw(tt.accountID)
			if err != nil {
				t.Fatal("AccountIDFromRaw() failed: %w", err)
			}
			acc, st, err := api.GetAccountWithProof(context.TODO(), accountID)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Account status: %v\n", acc.Account.Status())
			fmt.Printf("Last proof utime: %v\n", st.ShardStateUnsplit.GenUtime)
		})
	}
}
