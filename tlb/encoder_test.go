package tlb

import (
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

type testCase struct {
	A uint32
	B Uint11
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

func TestEnumFail(t *testing.T) {
	var cases = []testSumType{
		{
			B: testCase{4, 6},
		},
		{
			SumType: "C",
			B:       testCase{4, 6},
		},
	}
	for i, enum := range cases {
		b1 := boc.NewCell()
		err := Marshal(b1, enum)
		if err == nil {
			t.Fatalf("case %d should fail but didn't", i)
		}
	}
}

func TestRef(t *testing.T) {
	type A struct {
		A uint8
		B Ref[uint16]
	}
	var a A
	a.A = 5
	a.B.Value = 6
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
	t.Skip() //todo: это не тест а говно
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
		A Int15 `tlb:"^"`
		B Int11 `tlb:"^"`
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

func TestMaybeTags(t *testing.T) {

	a := Int15(101)
	b := Int11(8)

	type A struct {
		A *Int15 `tlb:"maybe^"`
		B *Int11 `tlb:"maybe"`
		C int32
	}

	value := A{A: &a, B: &b, C: 3}
	var value2 A
	c := boc.NewCell()
	err := Marshal(c, value)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Refs()) != 1 {
		t.Fatal("invalid refs number")
	}
	err = Unmarshal(c, &value2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, value2) {
		t.Fatal("not equal")
	}
}

func TestPointer(t *testing.T) {
	var a struct {
		A *int32
		B *struct {
			C *int32
		}
	}
	i := int32(100)
	s := struct {
		C *int32
	}{&i}
	a.A = &i
	a.B = &s
	c := boc.NewCell()
	err := Marshal(c, a)
	if err != nil {
		t.Fatal(err)
	}
	var b struct {
		A *int32
		B *struct {
			C *int32
		}
	}
	err = Unmarshal(c, &b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("not equal")
	}
}

type A struct {
	A int32
}

func (a *A) UnmarshalTLB(cell *boc.Cell, d *Decoder) error {
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
