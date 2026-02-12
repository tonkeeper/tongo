package dns

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestResolve(t *testing.T) {
	client, err := liteapi.NewClient(liteapi.Mainnet(), liteapi.FromEnvs())
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	root, err := client.GetRootDNS(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	dns := NewDNS(root, client)

	for _, c := range []struct {
		domain       string
		wallet       string
		pictureURL   string
		pictureBagID string
		success      bool
	}{
		{"industries.ton", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", "", "", true},
		{"oo0ili0oo.t.me", "0:afa066774812c345ff23ad53e04225e657134a087d133021b3fb8667a11efe74", "", "", true},
		{"metadata.ton", "0:afa066774812c345ff23ad53e04225e657134a087d133021b3fb8667a11efe74", "", "", false},
		{"shibdev.dolboeb.t.me", "0:2c979f6a3b2f42d972f9e39da519b1d4f5b23797896145aa2978d67cd8d44af2", "", "", true},
		{"hfdshfkjshkjdhfklhldkfhlakjh.ton", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", "", "", false},
		{"ololo.png", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", "", "", false},
		{"hasselbach.ton", "0:85b1573c32727a3f053a2e2f04157d6fa65cdcd1fe06444b7d52cc8a35277146", "", "3496f55fcad2007d1ef9fcf18639493e25109c353c0583d4a3f694dfbdf06dad", true},
		{"kontext.ton", "0:85b1573c32727a3f053a2e2f04157d6fa65cdcd1fe06444b7d52cc8a35277146", "https://hasselbach.ru/hydratingkitten.png", "", true},
	} {
		t.Run(c.domain, func(t *testing.T) {
			res, err := dns.Resolve(context.Background(), c.domain)
			if err != nil {
				if c.success {
					t.Fatalf("Expected to resolve domain: %v", err)
				}
				return
			}
			if !c.success {
				t.Fatalf("Not expected to resolve domain: %v", err)
			}

			r, ok := res[DNSCategoryWallet]
			if c.wallet != "" {
				if !ok {
					t.Error("wallet record not found")
				} else {
					assertDNSSmcAddress(t, r, c.wallet)
				}
			} else if ok {
				t.Error("unexpected wallet record")
			}

			r, ok = res[DNSCategoryPicture]
			if c.pictureURL != "" || c.pictureBagID != "" {
				if !ok {
					t.Error("picture record not found")
				} else {
					if c.pictureURL != "" {
						assertDNSText(t, r, c.pictureURL)
					}
					if c.pictureBagID != "" {
						assertDNSStorageAddress(t, r, c.pictureBagID)
					}
				}
			} else if ok {
				t.Error("unexpected picture record")
			}

			fmt.Printf("Qty of DNS records: %v\n", len(res))
		})

	}

}

func assertDNSSmcAddress(t *testing.T, r tlb.DNSRecord, expected string) {
	t.Helper()
	if r.SumType != "DNSSmcAddress" {
		t.Errorf("expected DNSSmcAddress, got %v", r.SumType)
		return
	}
	a, err := ton.AccountIDFromTlb(r.DNSSmcAddress.Address)
	if err != nil {
		t.Errorf("AccountIDFromTlb: %v", err)
		return
	}
	if a == nil {
		t.Error("expected non-nil account ID")
		return
	}
	if a.ToRaw() != expected {
		t.Errorf("expected wallet %v, got %v", expected, a.ToRaw())
	}
}

func assertDNSText(t *testing.T, r tlb.DNSRecord, expected string) {
	t.Helper()
	if r.SumType != "DNSText" {
		t.Errorf("expected DNSText, got %v", r.SumType)
		return
	}
	if string(r.DNSText) != expected {
		t.Errorf("expected %v, got %v", expected, string(r.DNSText))
	}
}

func assertDNSStorageAddress(t *testing.T, r tlb.DNSRecord, expected string) {
	t.Helper()
	if r.SumType != "DNSStorageAddress" {
		t.Errorf("expected DNSStorageAddress, got %v", r.SumType)
		return
	}
	if hex.EncodeToString(r.DNSStorageAddress[:]) != expected {
		t.Errorf("expected bag ID %v, got %v", expected, hex.EncodeToString(r.DNSStorageAddress[:]))
	}
}
