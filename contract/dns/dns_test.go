package dns

import (
	"context"
	"fmt"
	"github.com/startfellows/tongo/liteapi"
	"log"
	"testing"
)

func TestResolve(t *testing.T) {
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}
	//root, _ := tongo.ParseAccountID("Ef_BimcWrQ5pmAWfRqfeVHUCNV8XgsLqeAMBivKryXrghFW3")
	dns, err := NewDNS(nil, client)
	if err != nil {
		log.Fatalf("Unable to create DNS: %v", err)
	}
	res, err := dns.Resolve(context.Background(), "industries.ton")
	if err != nil {
		log.Fatalf("Unable to resolve domain: %v", err)
	}
	fmt.Printf("Qty of DNS records: %v\n", len(res))
}
