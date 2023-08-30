package dns

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
)

func mustDecodeAdnlHex(s string) [32]byte {
	bs, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	var b [32]byte
	copy(b[:], bs)
	return b
}

func mustToMsgAddress(s string) tlb.MsgAddress {
	a, err := tongo.AccountIDFromRaw(s)
	if err != nil {
		panic(err)
	}
	return a.ToMsgAddress()
}

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
		records []tlb.DNSRecord
		success bool
	}{
		{
			domain: "dark-matter-token.ton",
			records: []tlb.DNSRecord{
				{
					SumType: "DNSAdnlAddress",
					DNSAdnlAddress: struct {
						Address   [32]byte
						ProtoList []string
					}{
						Address: mustDecodeAdnlHex("45061c1d4ec44a937d0318589e13c73d151d1cef5d3c0e53afbcf56a6c2fe2bd"),
					},
				},
			},
			success: true,
		},
		{
			domain: "sudomuddydev.ton",
			records: []tlb.DNSRecord{
				{
					SumType: "DNSSmcAddress",
					DNSSmcAddress: struct {
						Address       tlb.MsgAddress
						SmcCapability tlb.SmcCapabilities
					}{
						Address: mustToMsgAddress("0:a99f55b54f6baba0588c55a1958a042662afa44df0fbf2fc234ed19454906103"),
					},
				},
				{
					SumType: "DNSAdnlAddress",
					DNSAdnlAddress: struct {
						Address   [32]byte
						ProtoList []string
					}{
						Address: mustDecodeAdnlHex("976cad182b618aed74cb2515de972606085bff7240654c2da74a3ac5619e733a"),
					},
				},
			},
			success: true,
		},
		{
			domain: "industries.ton",
			records: []tlb.DNSRecord{
				{
					SumType: "DNSSmcAddress",
					DNSSmcAddress: struct {
						Address       tlb.MsgAddress
						SmcCapability tlb.SmcCapabilities
					}{
						Address: mustToMsgAddress("0:ab027c8b08f5bbb529d643b64eff3b434e3d236347d697e4ffc8a6e8ba160504"),
					},
				},
			},
			success: true,
		},
		{
			domain: "oo0ili0oo.t.me",
			records: []tlb.DNSRecord{
				{
					SumType: "DNSSmcAddress",
					DNSSmcAddress: struct {
						Address       tlb.MsgAddress
						SmcCapability tlb.SmcCapabilities
					}{
						Address: mustToMsgAddress("0:afa066774812c345ff23ad53e04225e657134a087d133021b3fb8667a11efe74"),
					},
				},
			},
			success: true,
		},
		//{"thekiba.dolbaeb.t.me", "0:82683859071f85ed07d10016b19b8ecd183933a46987aed9fdc502f250d9404a", true}, //todo: fix. thekiba contract
		{
			domain:  "hfdshfkjshkjdhfklhldkfhlakjh.ton",
			success: false,
		},
	} {
		t.Run(c.domain, func(t *testing.T) {
			res, err := dns.Resolve(context.Background(), c.domain)
			if (err == nil) != c.success {
				t.Fatalf("Unable to resolve domain: %v", err)
			}
			if !c.success {
				return
			}
			fmt.Printf("Qty of DNS records: %v\n", len(res))
			if len(res) != len(c.records) {
				t.Fatalf("Wrong qty of records: %v, expected: %v", len(res), len(c.records))
			}
			for i, r := range res {
				if r.SumType != c.records[i].SumType {
					t.Fatalf("Wrong record type: %v, expected: %v", r.SumType, c.records[i].SumType)
				}
				switch r.SumType {
				case "DNSAdnlAddress":
					fmt.Printf("DNSAdnlAddress: %v\n", hex.EncodeToString(r.DNSAdnlAddress.Address[:]))
				case "DNSSmcAddress":
					addr, _ := tongo.AccountIDFromTlb(r.DNSSmcAddress.Address)
					fmt.Printf("DNSSmcAddress: %v\n", addr.ToRaw())
				}

				if !reflect.DeepEqual(r, c.records[i]) {
					t.Fatalf("Wrong record: %v, expected: %v", r, c.records[i])
				}
			}
		})

	}

}
