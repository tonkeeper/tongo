package ton

import (
	"encoding/json"
	"testing"
)

func TestAccountIDJsonUnmarshal(t *testing.T) {
	input := []byte(`{"A": "-1:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e"}`)
	var a struct {
		A AccountID
	}
	err := json.Unmarshal(input, &a)
	if err != nil {
		t.Fatal(err)
	}
	if a.A.Workchain != -1 {
		t.Fatal("invalid workchain")
	}
	for i, b := range []byte{112, 20, 167, 158, 183, 168, 28, 243, 117, 66, 166, 43, 117, 222, 250, 153, 66, 117, 128, 230, 97, 47, 149, 109, 71, 202, 160, 254, 14, 197, 208, 94} {
		if a.A.Address[i] != b {
			t.Fatal("invalid address")
		}
	}
}

func TestParseRawAccountID(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		output  string
		success bool
	}{
		{
			name:    "master",
			input:   "-1:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			output:  "-1:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: true,
		},
		{
			name:    "base",
			input:   "0:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			output:  "0:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: true,
		},
		{
			name:    "upper case",
			input:   "0:7014A79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			output:  "0:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: true,
		},
		{
			name:    "zfill1",
			input:   "0:014A79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			output:  "0:0014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: true,
		},
		{
			name:    "zfill2",
			input:   "0:14A79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			output:  "0:0014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: true,
		},
		{
			name:    "invalid",
			input:   "0:14A79eb7a8ZZZZ",
			output:  "0:0014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: false,
		},
		{
			name:    "invalid2",
			input:   "0:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e:",
			output:  "0:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e",
			success: false,
		},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			a, err := AccountIDFromRaw(c.input)
			if (err == nil) != c.success {
				t.Fatal(err)
			}
			if c.success && a.ToRaw() != c.output {
				t.Fatal(a.ToRaw())
			}
		})
	}
}

func TestToMsgAddress(t *testing.T) {
	ma := (*AccountID)(nil).ToMsgAddress()
	if ma.SumType != "AddrNone" {
		t.Fatal(ma.SumType)
	}
}
