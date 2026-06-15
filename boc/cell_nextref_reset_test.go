package boc

import "testing"

func TestNextRefShallowReset(t *testing.T) {
	// root -> child -> grandchild, each carrying 8 readable bits.
	grandchild := NewCell()
	if err := grandchild.WriteUint(0xAA, 8); err != nil {
		t.Fatal(err)
	}
	child := NewCell()
	if err := child.WriteUint(0xBB, 8); err != nil {
		t.Fatal(err)
	}
	if err := child.AddRef(grandchild); err != nil {
		t.Fatal(err)
	}
	root := NewCell()
	if err := root.AddRef(child); err != nil {
		t.Fatal(err)
	}
	// Advance the read cursors of child and grandchild so a reset is observable.
	if _, err := child.ReadUint(8); err != nil {
		t.Fatal(err)
	}
	if _, err := grandchild.ReadUint(8); err != nil {
		t.Fatal(err)
	}
	if child.BitsAvailableForRead() != 0 || grandchild.BitsAvailableForRead() != 0 {
		t.Fatalf("setup: expected both cursors exhausted")
	}
	// Descend into child: its own cursor must be reset
	got, err := root.NextRef()
	if err != nil {
		t.Fatal(err)
	}
	if got != child {
		t.Fatal("NextRef returned wrong cell")
	}
	if child.BitsAvailableForRead() != 8 {
		t.Fatalf("child cursor not reset: got %d want 8", child.BitsAvailableForRead())
	}
	// but the grandchild must NOT have been touched (shallow reset).
	if grandchild.BitsAvailableForRead() != 0 {
		t.Fatalf("grandchild was reset by NextRef into child — reset is not shallow (the regression)")
	}
	// Correctness preserved: descending into the grandchild resets it then.
	gotGC, err := child.NextRef()
	if err != nil {
		t.Fatal(err)
	}
	if gotGC != grandchild {
		t.Fatal("NextRef returned wrong grandchild")
	}
	if grandchild.BitsAvailableForRead() != 8 {
		t.Fatalf("grandchild cursor not reset when descended into: got %d want 8", grandchild.BitsAvailableForRead())
	}
}
