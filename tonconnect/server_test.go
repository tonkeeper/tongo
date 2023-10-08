package tonconnect

import (
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteapi"
)

func TestGenerateAndVerifyPayload(t *testing.T) {
	liteApiClient, _ := liteapi.NewClientWithDefaultMainnet()
	tonConnect, _ := NewTonConnect(liteApiClient, "my_secret", WithLifeTimePayload(300), WithLifeTimeProof(300))

	payload, err := tonConnect.GeneratePayload()
	if err != nil {
		t.Fatalf("failed generate payload: %v", err)
	}
	verify, err := tonConnect.CheckPayload(payload)
	if err != nil {
		t.Fatalf("failed verify payload: %v", err)
	}

	if !verify {
		t.Fatalf("failed verify payload")
	}
}

func TestExpirePayload(t *testing.T) {
	liteApiClient, _ := liteapi.NewClientWithDefaultMainnet()
	tonConnect, _ := NewTonConnect(liteApiClient, "my_secret", WithLifeTimePayload(1), WithLifeTimeProof(1)) // set little lifetime

	payload, err := tonConnect.GeneratePayload()
	if err != nil {
		t.Fatalf("failed generate payload: %v", err)
	}

	time.Sleep(2 * time.Second) // waiting expire

	verify, _ := tonConnect.CheckPayload(payload)
	if verify {
		t.Fatalf("payload not expire")
	}
}
