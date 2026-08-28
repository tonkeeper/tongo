package tvm

import (
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

// TestEmulatorAcceptsBuiltLeveledBoc feeds a cell tree built in memory around a pruned
// branch to ton's own boc deserializer (through libemulator), which validates the level
// mask every cell declares. A parent claiming level 0 over a level 1 child is rejected
// there with "level mask mismatch".
func TestEmulatorAcceptsBuiltLeveledBoc(t *testing.T) {
	pruned := boc.NewCellExotic(boc.PrunedBranchCell)
	if err := pruned.WriteUint(uint64(boc.PrunedBranchCell), 8); err != nil {
		t.Fatalf("WriteUint() failed: %v", err)
	}
	if err := pruned.WriteUint(1, 8); err != nil { // level mask of the replaced cell
		t.Fatalf("WriteUint() failed: %v", err)
	}
	if err := pruned.WriteBytes(make([]byte, 32)); err != nil { // hash of the replaced cell
		t.Fatalf("WriteBytes() failed: %v", err)
	}
	if err := pruned.WriteUint(0, 16); err != nil { // depth of the replaced cell
		t.Fatalf("WriteUint() failed: %v", err)
	}
	pruned.ResetCounters()

	// an ordinary cell over a pruned branch: this is the one that used to claim level 0
	virtualRoot := boc.NewCell()
	if err := virtualRoot.AddRef(pruned); err != nil {
		t.Fatalf("AddRef() failed: %v", err)
	}
	virtualRootHash, err := virtualRoot.HashAtLevel(0)
	if err != nil {
		t.Fatalf("HashAtLevel() failed: %v", err)
	}

	// a merkle proof root keeps the tree at level 0 overall, which is what the
	// deserializer requires of a root cell.
	root := boc.NewCellExotic(boc.MerkleProofCell)
	if err := root.WriteUint(uint64(boc.MerkleProofCell), 8); err != nil {
		t.Fatalf("WriteUint() failed: %v", err)
	}
	if err := root.WriteBytes(virtualRootHash[:]); err != nil {
		t.Fatalf("WriteBytes() failed: %v", err)
	}
	if err := root.WriteUint(1, 16); err != nil { // depth of the virtual root
		t.Fatalf("WriteUint() failed: %v", err)
	}
	if err := root.AddRef(virtualRoot); err != nil {
		t.Fatalf("AddRef() failed: %v", err)
	}
	root.ResetCounters()

	dataBoc, err := root.ToBocBase64()
	if err != nil {
		t.Fatalf("ToBocBase64() failed: %v", err)
	}
	codeBoc, err := boc.NewCell().ToBocBase64()
	if err != nil {
		t.Fatalf("ToBocBase64() failed: %v", err)
	}
	if _, err := NewEmulatorFromBOCsBase64(codeBoc, dataBoc, ""); err != nil {
		t.Fatalf("ton refused the boc we produced: %v", err)
	}

	// the very same tree, but with the ordinary parent claiming level 0 over its
	// level 1 child - byte for byte what tongo used to emit before level masks
	// were derived. ton rejects it with "level mask mismatch".
	const brokenBoc = "te6ccgEBAwEATwAJRgMRvBFcAiYynJo1p/7h0pwfUkjg94QdV4KbUCeCYxSREwABAQEAAihIAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	if brokenBoc == dataBoc {
		t.Fatal("the broken boc must differ from what we emit now")
	}
	if _, err := NewEmulatorFromBOCsBase64(codeBoc, brokenBoc, ""); err == nil {
		t.Error("ton accepted a boc with a level mask mismatch, the check is not doing anything")
	}
}
