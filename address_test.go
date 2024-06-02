package tongo

import (
	"context"
	"testing"

	"github.com/tonkeeper/tongo/contract/dns"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
)

func TestParseAddress(t *testing.T) {
	cli, err := liteapi.NewClient(liteapi.Mainnet(), liteapi.FromEnvs())
	if err != nil {
		t.Fatalf("failed to create liteapi client: %v", err)
	}
	resolver := dns.NewDNS(ton.MustParseAccountID(DefaultRoot), cli)
	parser := NewAccountAddressParser(resolver)

	const (
		parseToHumanAddress = iota
		parseToRawAddress
		parseDnsToRawAddress
		parseTestnetOnlyAddress
	)

	type testCase struct {
		name                string
		typeParse           int
		request             string
		response            string
		responseTestnetOnly bool
	}

	for _, test := range []testCase{
		{
			name:                "Parse to raw address",
			typeParse:           parseToHumanAddress,
			request:             "0:91d73056e035232f09aaf8242a1d51eea98b6a5bebbf8ac0c9e521d02a1a4bdb",
			response:            "EQCR1zBW4DUjLwmq-CQqHVHuqYtqW-u_isDJ5SHQKhpL2wQV",
			responseTestnetOnly: false,
		},
		{
			name:                "Parse to human address",
			typeParse:           parseToRawAddress,
			request:             "EQCR1zBW4DUjLwmq-CQqHVHuqYtqW-u_isDJ5SHQKhpL2wQV",
			response:            "0:91d73056e035232f09aaf8242a1d51eea98b6a5bebbf8ac0c9e521d02a1a4bdb",
			responseTestnetOnly: false,
		},
		{
			name:                "Parse dns to raw address",
			typeParse:           parseDnsToRawAddress,
			request:             "blackpepper.ton",
			response:            "0:44556b55c15052eb44c6b75a9eccbc6280d32d598d12e975f435195795bb11d5",
			responseTestnetOnly: false,
		},
		{
			name:                "Parse dns to raw address",
			typeParse:           parseDnsToRawAddress,
			request:             "subbotin.ton",
			response:            "0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb",
			responseTestnetOnly: false,
		},
		{
			name:                "Parse only testnet",
			typeParse:           parseTestnetOnlyAddress,
			request:             "kQCR1zBW4DUjLwmq-CQqHVHuqYtqW-u_isDJ5SHQKhpL27-f",
			response:            "0:91d73056e035232f09aaf8242a1d51eea98b6a5bebbf8ac0c9e521d02a1a4bdb",
			responseTestnetOnly: true,
		},
		{
			name:      "url-unsafe",
			typeParse: parseToRawAddress,
			request:   "UQBDWnKuRx7eqYtr5Kr9HdFGHnBsUyX_jsPGC/RO/K4BaVdu",
			response:  "0:435a72ae471edea98b6be4aafd1dd1461e706c5325ff8ec3c60bf44efcae0169",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			account, err := parser.ParseAddress(context.Background(), test.request)
			if err != nil {
				t.Fatalf("failed parse %v address: %v", test.request, err)
			}
			switch test.typeParse {
			case parseToHumanAddress:
				if account.ID.ToHuman(true, false) != test.response && account.TestnetOnly != test.responseTestnetOnly {
					t.Fatalf("not equal address")
				}
			case parseToRawAddress:
				if account.ID.ToRaw() != test.response && account.TestnetOnly != test.responseTestnetOnly {
					t.Fatalf("not equal address")
				}
			case parseDnsToRawAddress:
				if account.ID.ToRaw() != test.response && account.TestnetOnly != test.responseTestnetOnly {
					t.Fatalf("not equal address")
				}
			case parseTestnetOnlyAddress:
				if account.ID.ToRaw() != test.response && account.TestnetOnly != test.responseTestnetOnly {
					t.Fatalf("not equal address")
				}
			}
		})
	}
}
