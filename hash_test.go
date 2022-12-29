package tongo

import (
	"bytes"
	"github.com/startfellows/tongo/tl"
	"reflect"
	"testing"
)

func TestHashTl(t *testing.T) {
	hash := Bits256{1, 2, 3, 4, 5, 6}
	b, err := tl.Marshal(hash)
	if err != nil {
		t.Fatal(err)
	}
	var hash1 Bits256
	err = tl.Unmarshal(bytes.NewReader(b), &hash1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hash, hash1) {
		t.Fatal("not equal")
	}
}
