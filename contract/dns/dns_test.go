package dns

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
)

func TestResolve(t *testing.T) {
	client, err := liteapi.NewClientWithDefaultMainnet()
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
		//{"thekiba.dolbaeb.t.me", "0:82683859071f85ed07d10016b19b8ecd183933a46987aed9fdc502f250d9404a", true}, //todo: fix. thekiba contract
		{"hfdshfkjshkjdhfklhldkfhlakjh.ton", "0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504", false},
	} {
		t.Run(c.domain, func(t *testing.T) {
			res, err := dns.Resolve(context.Background(), c.domain)
			if (err == nil) != c.success {
				t.Fatalf("Unable to resolve domain: %v", err)
			}
			if c.success {
				a, _ := ton.AccountIDFromTlb(res[0].DNSSmcAddress.Address)
				if a.ToRaw() != c.wallet {
					t.Fatal("invalid wallet")
				}
			}
			fmt.Printf("Qty of DNS records: %v\n", len(res))
		})

	}

}
