package dns

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/tonkeeper/tongo/liteapi"
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
		domain  string
		wallet  string
		success bool
	}{
		{"industries.ton", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", true},
		{"oo0ili0oo.t.me", "0:afa066774812c345ff23ad53e04225e657134a087d133021b3fb8667a11efe74", true},
		{"packages.ton", "0:afa066774812c345ff23ad53e04225e657134a087d133021b3fb8667a11efe74", false},
		{"shibdev.dolboeb.t.me", "0:2c979f6a3b2f42d972f9e39da519b1d4f5b23797896145aa2978d67cd8d44af2", true},
		{"hfdshfkjshkjdhfklhldkfhlakjh.ton", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", false},
		{"ololo.png", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", false},
	} {
		t.Run(c.domain, func(t *testing.T) {
			res, err := dns.Resolve(context.Background(), c.domain)
			if (err == nil) != c.success {
				t.Fatalf("Unable to resolve domain: %v", err)
			}
			if err != nil && !errors.Is(err, ErrNotResolved) && !errors.Is(err, liteapi.ErrAccountNotFound) {
				t.Fatal(err)
			}
			if c.success {
				a, err := ton.AccountIDFromTlb(res[0].DNSSmcAddress.Address)
				if err != nil {
					t.Fatal(err)
				}
				if a.ToRaw() != c.wallet {
					t.Fatal("invalid wallet")
				}
			}
			fmt.Printf("Qty of DNS records: %v\n", len(res))
		})

	}

}
