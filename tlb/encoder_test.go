package tlb

import (
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

type testCase struct {
	A uint32 `tlb:"32bits"`
	B uint32 `tlb:"11bits"`
}

type testSumType struct {
	SumType
	A testCase `tlbSumType:"a$0"`
	B testCase `tlbSumType:"b$1"`
}

func TestStruct(t *testing.T) {
	var cases = []testCase{
		{3, 0},
		{1<<32 - 1, 1},
		{94823094, 100},
	}
	for _, c := range cases {
		b1 := boc.NewCell()
		err := Marshal(b1, c)
		if err != nil {
			t.Fatal(err)
		}
		var unmarshaled testCase
		err = Unmarshal(b1, &unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(c, unmarshaled) {
			t.Fatalf("not equal. left: %v, right: %v \n", c, unmarshaled)
		}
		b2 := boc.NewCell()
		err = Marshal(b2, unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if b1.ToString() != b2.ToString() {
			t.Fatal("not equal")
		}
	}
}

func TestEnum(t *testing.T) {
	var cases = []testSumType{
		{
			SumType: "A",
			A:       testCase{3, 5},
		},
		{
			SumType: "B",
			B:       testCase{4, 6},
		},
	}
	for _, enum := range cases {
		b1 := boc.NewCell()
		err := Marshal(b1, enum)
		if err != nil {
			t.Fatal(err)
		}
		var unmarshaled testSumType
		err = Unmarshal(b1, &unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(enum, unmarshaled) {
			t.Fatalf("not equal. left: %v, right: %v \n", enum, unmarshaled)
		}
		b2 := boc.NewCell()
		err = Marshal(b2, unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if b1.ToString() != b2.ToString() {
			t.Fatal("not equal")
		}
	}
}

func TestRef(t *testing.T) {
	type A struct {
		A uint32 `tlb:"8bits"`
		B Ref[struct {
			C uint32 `tlb:"16bits"`
		}]
	}
	var a A
	a.A = 5
	a.B.Value.C = 6
	c := boc.NewCell()
	err := Marshal(c, a)
	if err != nil {
		t.Fatal(err)
	}
	var a2 A
	err = Unmarshal(c, &a2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a2) {
		t.Fatal("not equal")
	}
}

func TestCell(t *testing.T) {
	type A struct {
		A Ref[boc.Cell]
	}
	var a A
	a.A.Value = *boc.NewCell()
	a.A.Value.WriteBit(true)
	c := boc.NewCell()
	err := Marshal(c, a)
	if err != nil {
		t.Fatal(err)
	}
	var a2 A
	err = Unmarshal(c, &a2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a2) {
		t.Fatal("not equal")
	}
}

func TestBitString(t *testing.T) {
	type test struct {
		A boc.BitString `tlb:"100bits"`
	}
	var a test
	a.A = boc.NewBitString(100)
	_ = a.A.WriteUint(100, 100)
	c := boc.NewCell()
	err := Marshal(c, a)
	if err != nil {
		t.Fatal(err)
	}
	var a2 test
	err = Unmarshal(c, &a2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a2) {
		t.Fatal("not equal")
	}
	b2 := boc.NewCell()
	err = Marshal(b2, a2)
	if err != nil {
		t.Fatal(err)
	}
	if c.ToString() != b2.ToString() {
		t.Fatal("not equal")
	}
}

func TestRefTag(t *testing.T) {
	type A struct {
		A int32 `tlb:"^ 15bits"`
		B int32 `tlb:"^11bits"`
		C int32
	}
	a := A{1, 2, 3}
	var a2 A
	c := boc.NewCell()
	err := Marshal(c, a)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Refs()) != 2 {
		t.Fatal("invalid cells number")
	}
	err = Unmarshal(c, &a2)
	if err != nil {
		t.Fatal(err)
	}
	if a != a2 {
		t.Fatal("not equal")
	}
}

type A struct {
	A int32
}

func (a *A) UnmarshalTLB(cell *boc.Cell, tag string) error {
	a.A = 100
	return nil
}

func TestDecoderInterface(t *testing.T) {
	var a A
	c := boc.NewCell()
	err := Unmarshal(c, &a)
	if err != nil {
		t.Fatal(err)
	}
	if a.A != 100 {
		t.Fatal(a.A)
	}
	var b struct {
		A A
	}
	err = Unmarshal(c, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b.A.A != 100 {
		t.Fatal(b.A.A)
	}
}