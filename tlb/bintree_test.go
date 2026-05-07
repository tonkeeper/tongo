package tlb

import (
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

// TestBinTree_PrunedBranchSkipped ensures decoding a BinTree that contains a
// pruned-branch subtree (as produced inside a Merkle update / proof) does not
// fail trying to interpret pruned-cell bytes as BinTree forks/leaves and as
// the underlying value type. The pruned subtree contributes zero values.
func TestBinTree_PrunedBranchSkipped(t *testing.T) {
	// Build a fork cell:
	//   bit '1' (fork) + 2 refs
	//     left  : pruned-branch cell      → must be skipped
	//     right : leaf cell (bit '0' + a uint8 value)
	root := boc.NewCell()
	if err := root.WriteBit(true); err != nil {
		t.Fatalf("write fork bit: %v", err)
	}

	pruned := boc.NewCellExotic(boc.PrunedBranchCell)
	if err := root.AddRef(pruned); err != nil {
		t.Fatalf("add pruned ref: %v", err)
	}

	leaf := boc.NewCell()
	if err := leaf.WriteBit(false); err != nil { // leaf indicator
		t.Fatalf("write leaf bit: %v", err)
	}
	if err := leaf.WriteUint(0x42, 8); err != nil { // payload value
		t.Fatalf("write leaf payload: %v", err)
	}
	if err := root.AddRef(leaf); err != nil {
		t.Fatalf("add leaf ref: %v", err)
	}

	var bt BinTree[uint8]
	if err := Unmarshal(root, &bt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(bt.Values) != 1 {
		t.Fatalf("expected 1 value (pruned subtree skipped), got %d", len(bt.Values))
	}
	if bt.Values[0] != 0x42 {
		t.Fatalf("expected leaf value 0x42, got 0x%x", bt.Values[0])
	}
}
