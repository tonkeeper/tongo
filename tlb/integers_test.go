package tlb

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestVarUInteger_MarshalJSON(t *testing.T) {
	type block struct {
		Seqno VarUInteger16
	}
	b1 := block{Seqno: VarUInteger16(*big.NewInt(11))}
	body, err := json.Marshal(b1)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	b2 := block{}
	err = json.Unmarshal(body, &b2)
	if err != nil {
		t.Errorf("json.Unmarshal() failed: %v", err)
	}
	value1 := big.Int(b1.Seqno)
	value2 := big.Int(b2.Seqno)
	if value1.Int64() != value2.Int64() {
		t.Errorf("want: %v, got: %v", value1.Int64(), value2.Int64())
	}
}

func TestUInt_MarshalJSON(t *testing.T) {
	type block struct {
		Seqno Uint256
	}
	b1 := block{Seqno: Uint256(*big.NewInt(147))}
	body, err := json.Marshal(b1)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	b2 := block{}
	err = json.Unmarshal(body, &b2)
	if err != nil {
		t.Errorf("json.Unmarshal() failed: %v", err)
	}
	value1 := big.Int(b1.Seqno)
	value2 := big.Int(b2.Seqno)
	if value1.Int64() != value2.Int64() {
		t.Errorf("want: %v, got: %v", value1.Int64(), value2.Int64())
	}
}

func TestUintBitsConversions(t *testing.T) {
	t.Run("128", func(t *testing.T) {
		want := "0102030405060708090a0b0c0d0e0f10"
		i, ok := new(big.Int).SetString(want, 16)
		if !ok {
			t.Fatalf("invalid test integer")
		}
		u := Uint128(*i)
		b := u.ToBits()
		if b.HexString() != want {
			t.Fatalf("ToBits() = %v, want %v", b.HexString(), want)
		}
		if got := big.Int(b.ToUint()); got.Cmp(i) != 0 {
			t.Fatalf("ToUint() = %v, want %v", &got, i)
		}
	})

	t.Run("160", func(t *testing.T) {
		want := "0102030405060708090a0b0c0d0e0f1011121314"
		i, ok := new(big.Int).SetString(want, 16)
		if !ok {
			t.Fatalf("invalid test integer")
		}
		u := Uint160(*i)
		b := u.ToBits()
		if b.HexString() != want {
			t.Fatalf("ToBits() = %v, want %v", b.HexString(), want)
		}
		if got := big.Int(b.ToUint()); got.Cmp(i) != 0 {
			t.Fatalf("ToUint() = %v, want %v", &got, i)
		}
	})

	t.Run("256", func(t *testing.T) {
		want := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		i, ok := new(big.Int).SetString(want, 16)
		if !ok {
			t.Fatalf("invalid test integer")
		}
		u := Uint256(*i)
		b := u.ToBits()
		if b.HexString() != want {
			t.Fatalf("ToBits() = %v, want %v", b.HexString(), want)
		}
		if got := big.Int(b.ToUint()); got.Cmp(i) != 0 {
			t.Fatalf("ToUint() = %v, want %v", &got, i)
		}
	})
}

func TestInt_MarshalJSON(t *testing.T) {
	type block struct {
		Seqno Int257
	}
	b1 := block{Seqno: Int257(*big.NewInt(257))}
	body, err := json.Marshal(b1)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	b2 := block{}
	err = json.Unmarshal(body, &b2)
	if err != nil {
		t.Errorf("json.Unmarshal() failed: %v", err)
	}
	value1 := big.Int(b1.Seqno)
	value2 := big.Int(b2.Seqno)
	if value1.Int64() != value2.Int64() {
		t.Errorf("want: %v, got: %v", value1.Int64(), value2.Int64())
	}
}

func TestBits_JSON(t *testing.T) {
	tests := []struct {
		name    string
		u       any
		want    string
		wantErr bool
	}{
		{
			u:    Bits256([32]byte{255, 0, 255}),
			want: `"ff00ff0000000000000000000000000000000000000000000000000000000000"`,
		},
		{
			u:    Bits512([64]byte{255, 0, 255}),
			want: `"ff00ff00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := json.Marshal(tt.u)
			if err != nil {
				t.Fatalf("json.Marshal() failed: %v", err)
			}
			if tt.want != string(value) {
				t.Errorf("want: %v, got: %v", tt.want, value)
			}
		})
	}
}
