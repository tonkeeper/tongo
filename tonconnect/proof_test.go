package tonconnect

import (
	"context"
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/wallet"
)

func TestCreateSignedProof(t *testing.T) {
	cli, err := liteapi.NewClient(liteapi.Testnet())
	if err != nil {
		t.Fatalf("liteapi.NewClient() failed: %v", err)
	}
	tests := []struct {
		name    string
		version wallet.Version
		secret  string
	}{
		{
			name:    "v4r2",
			version: wallet.V4R2,
			secret:  "some-random-secret",
		},
		{
			name:    "v4r1",
			version: wallet.V4R1,
			secret:  "another-random-secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := NewTonConnect(cli, tt.secret)
			if err != nil {
				t.Fatalf("NewTonConnect() failed: %v", err)
			}
			payload, err := srv.GeneratePayload()
			if err != nil {
				t.Fatalf("GeneratePayload() failed: %v", err)
			}
			seed := wallet.RandomSeed()
			privateKey, err := wallet.SeedToPrivateKey(seed)
			if err != nil {
				t.Fatalf("SeedToPrivateKey() failed: %v", err)
			}
			publicKey := privateKey.Public().(ed25519.PublicKey)
			stateInit, err := wallet.GenerateStateInit(publicKey, tt.version, 0, nil)
			if err != nil {
				t.Fatalf("GenerateStateInit() failed: %v", err)
			}
			accountID, err := wallet.GenerateWalletAddress(publicKey, tt.version, 0, nil)
			if err != nil {
				t.Fatalf("GenerateWalletAddress() failed: %v", err)
			}
			signedProof, err := CreateSignedProof(payload, accountID, privateKey, stateInit, ProofOptions{Timestamp: time.Now(), Domain: "web"})
			if err != nil {
				t.Fatalf("CreateSignedProof() failed: %v", err)
			}
			verified, _, err := srv.CheckProof(context.Background(), signedProof, srv.checkPayload)
			if err != nil {
				t.Fatalf("CheckProof() failed: %v", err)
			}
			if verified != true {
				t.Fatalf("proof is invalid")
			}
		})
	}
}
