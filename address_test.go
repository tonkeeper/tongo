package tongo

import (
	"context"
	"testing"
)

func TestParseAddress(t *testing.T) {
	parser := DefaultAddressParser()

	const (
		parseToHumanAddress = iota
		parseToRawAddress
		parseDnsToRawAddress
	)

	type testCase struct {
		name      string
		typeParse int
		request   string
		response  string
	}

	for _, test := range []testCase{
		{
			name:      "Parse to raw address",
			typeParse: parseToHumanAddress,
			request:   "0:91d73056e035232f09aaf8242a1d51eea98b6a5bebbf8ac0c9e521d02a1a4bdb",
			response:  "EQCR1zBW4DUjLwmq-CQqHVHuqYtqW-u_isDJ5SHQKhpL2wQV",
		},
		{
			name:      "Parse to human address",
			typeParse: parseToRawAddress,
			request:   "EQCR1zBW4DUjLwmq-CQqHVHuqYtqW-u_isDJ5SHQKhpL2wQV",
			response:  "0:91d73056e035232f09aaf8242a1d51eea98b6a5bebbf8ac0c9e521d02a1a4bdb",
		},
		{
			name:      "Parse dns to raw address",
			typeParse: parseDnsToRawAddress,
			request:   "blackpepper.ton",
			response:  "0:44556b55c15052eb44c6b75a9eccbc6280d32d598d12e975f435195795bb11d5",
		},
		{
			name:      "Parse dns to raw address",
			typeParse: parseDnsToRawAddress,
			request:   "subbotin.ton",
			response:  "0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			account, err := parser.ParseAddress(context.Background(), test.request)
			if err != nil {
				t.Fatalf("failed parse %v address: %v", test.request, err)
			}
			switch test.typeParse {
			case parseToHumanAddress:
				if account.ID.ToHuman(true, false) != test.response {
					t.Fatalf("not equal address")
				}
			case parseToRawAddress:
				if account.ID.ToRaw() != test.response {
					t.Fatalf("not equal address")
				}
			case parseDnsToRawAddress:
				if account.ID.ToRaw() != test.response {
					t.Fatalf("not equal address")
				}
			}
		})
	}
}
