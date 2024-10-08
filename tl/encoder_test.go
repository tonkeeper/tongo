package tl

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type testCase struct {
	A int32
	B int64
	C uint32
	D uint64
	E []byte
}

func TestMarshal(t *testing.T) {
	longSlice := make([]byte, 10000)
	for i := range longSlice {
		longSlice[i] = byte(i % 256)
	}
	var cases = []testCase{
		{
			A: 1,
			B: 2,
			C: 3,
			D: 4,
			E: []byte{1},
		},
		{
			A: 1<<31 - 1,
			B: 1<<63 - 1,
			C: 1<<32 - 1,
			D: 1<<64 - 1,
			E: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249},
		},
		{
			A: -10,
			B: -100,
			C: 94823094,
			D: 4124124124,
			E: longSlice,
		},
	}
	for _, c := range cases {
		b1, err := Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		var unmarshaled testCase
		err = Unmarshal(bytes.NewReader(b1), &unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(c, unmarshaled) {
			t.Fatal("not equal")
		}
		b2, err := Marshal(unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b1, b2) {
			t.Fatal("not equal")
		}
	}
}

type custom struct {
	A int32
}

type A struct {
	A int32
	B custom
	C int32
}

func (c custom) MarshalTL() ([]byte, error) {
	return []byte{1, 1, 1, 1}, nil
}

func (c *custom) UnmarshalTL(r io.Reader) error {
	b := []byte{0, 0, 0, 0}
	r.Read(b)
	c.A = 15
	return nil
}

func TestMarshalCustomDecode(t *testing.T) {
	a := A{
		A: 12,
		B: custom{A: 15},
		C: 13,
	}
	b, err := Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 A
	err = Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}

type testSumType struct {
	SumType
	A testCase `tlSumType:"0a0a0a0a"`
	B testCase `tlSumType:"0b0b0b0b"`
}

func TestSumType(t *testing.T) {
	a := testSumType{
		SumType: "B",
		A:       testCase{},
		B:       testCase{A: 1, E: []byte{1}},
	}
	b, err := Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 testSumType
	err = Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}
