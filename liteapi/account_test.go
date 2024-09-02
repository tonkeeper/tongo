package liteapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
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

func TestUnmarshallingProofWithPrunedResolver(t *testing.T) {
	testCases := []struct {
		name      string
		accountID string
		state     string
		proof     string
	}{
		{
			name:      "account from masterchain",
			accountID: "-1:34517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf",
			state:     "te6ccgEBAQEANwAAac/zRRfHvfUYfFWvT4th/cMhWIx6t2je4ksAbfKRBkWNfPICV8MiQ7WQAAAnwcjbgQjD0JAE",
			proof:     "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",
		},
		{
			name:      "active account from basechain",
			accountID: "0:e33ed33a42eb2032059f97d90c706f8400bb256d32139ca707f1564ad699c7dd",
			state:     "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",
			proof:     "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",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accountID, err := ton.AccountIDFromRaw(tt.accountID)
			if err != nil {
				t.Fatal("AccountIDFromRaw() failed: %w", err)
			}
			state, err := base64.StdEncoding.DecodeString(tt.state)
			if err != nil {
				t.Fatal("base64 decoding failed: %w", err)
			}
			proof, err := base64.StdEncoding.DecodeString(tt.proof)
			if err != nil {
				t.Fatal("base64 decoding failed: %w", err)
			}
			stateCells, err := boc.DeserializeBoc(state)
			if err != nil {
				t.Fatal("DeserializeBoc() failed: %w", err)
			}
			proofCells, err := boc.DeserializeBoc(proof)
			if err != nil {
				t.Fatal("DeserializeBoc() failed: %w", err)
			}
			cellMap, err := stateCells[0].NonPrunedCells()
			if err != nil {
				t.Fatal("Get NonPrunedCells() failed: %w", err)
			}
			decoder := tlb.NewDecoder().WithDebug().WithPrunedResolver(func(hash tlb.Bits256) (*boc.Cell, error) {
				if cellMap == nil {
					return nil, fmt.Errorf("failed to fetch library: no resolver provided")
				}
				cell, ok := cellMap[hash]
				if ok {
					return cell, nil
				}
				return nil, errors.New("not found")
			})
			var stateProof struct {
				Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
			}
			err = decoder.Unmarshal(proofCells[1], &stateProof)
			if err != nil {
				t.Fatal("proof unmarshalling failed: %w", err)
			}
			values := stateProof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Values()
			keys := stateProof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Keys()
			for i, k := range keys {
				if bytes.Equal(k[:], accountID.Address[:]) {
					fmt.Printf("Account status: %v\n", values[i].Account.Status())
				}
			}
		})
	}
}
