package account

import (
	"encoding/json"
	"testing"
)

func TestAccountIDJsonUnmarshal(t *testing.T) {
	input := []byte(`{"A": "-1:7014a79eb7a81cf37542a62b75defa99427580e6612f956d47caa0fe0ec5d05e"}`)
	var a struct {
		A ID
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

func TestToMsgAddress(t *testing.T) {
	ma := (*ID)(nil).ToMsgAddress()
	if ma.SumType != "AddrNone" {
		t.Fatal(ma.SumType)
	}
}
