package boc

import (
	"errors"
	"testing"
)

// newPrunedBranch builds a well-formed level-1 pruned branch cell
// standing in for a cell with the given hash and depth.
func newPrunedBranch(t *testing.T, hash []byte, depth int) *Cell {
	t.Helper()
	c := NewCellExotic(PrunedBranchCell)
	if err := c.WriteUint(uint64(PrunedBranchCell), 8); err != nil {
		t.Fatalf("write type: %v", err)
	}
	if err := c.WriteUint(1, 8); err != nil { // level mask of the replaced cell
		t.Fatalf("write mask: %v", err)
	}
	if err := c.WriteBytes(hash); err != nil {
		t.Fatalf("write hash: %v", err)
	}
	if err := c.WriteUint(uint64(depth), 16); err != nil {
		t.Fatalf("write depth: %v", err)
	}
	c.ResetCounters()
	return c
}

// declaredMasks returns the level mask each cell of the given boc declares in its d1 byte,
// in serialization order (the root comes first).
func declaredMasks(t *testing.T, bocBytes []byte) []levelMask {
	t.Helper()
	header, err := parseBocHeader(bocBytes)
	if err != nil {
		t.Fatalf("parseBocHeader() failed: %v", err)
	}
	data := header.cellsData
	masks := make([]levelMask, 0, header.cellCount)
	for i := 0; i < int(header.cellCount); i++ {
		d1, d2 := data[0], data[1]
		if d1&0b10000 != 0 {
			t.Fatalf("cell %v unexpectedly serialized with hashes", i)
		}
		masks = append(masks, levelMask(d1>>5))
		dataBytes := int(d2>>1) + int(d2%2)
		data = data[2+dataBytes+header.sizeBytes*int(d1%8):]
	}
	return masks
}

// TestSerializeBoc_LevelMaskOfBuiltTree makes sure a cell tree built in memory
// around a leveled cell declares the right level mask. Getting this wrong produces
// a boc that ton rejects with "level mask mismatch".
func TestSerializeBoc_LevelMaskOfBuiltTree(t *testing.T) {
	pruned := newPrunedBranch(t, make([]byte, hashSize), 7)

	root := NewCell()
	if err := root.WriteUint(0xdeadbeef, 32); err != nil {
		t.Fatalf("WriteUint() failed: %v", err)
	}
	// an intermediate ordinary cell must propagate the level up to the root
	middle := NewCell()
	if err := middle.AddRef(pruned); err != nil {
		t.Fatalf("AddRef() failed: %v", err)
	}
	if err := root.AddRef(middle); err != nil {
		t.Fatalf("AddRef() failed: %v", err)
	}

	if level := root.Level(); level != 1 {
		t.Fatalf("root.Level() = %v, want 1", level)
	}
	bocBytes, err := root.ToBoc()
	if err != nil {
		t.Fatalf("ToBoc() failed: %v", err)
	}
	for i, mask := range declaredMasks(t, bocBytes) {
		if mask != 1 {
			t.Errorf("cell %v declares level mask %v, want 1", i, mask)
		}
	}
	if _, err := DeserializeBoc(bocBytes); err != nil {
		t.Fatalf("DeserializeBoc() failed: %v", err)
	}
}

// TestCreateProof_MerkleProofLevelMask makes sure a merkle proof root shifts
// its child's level mask right instead of inheriting it.
func TestCreateProof_MerkleProofLevelMask(t *testing.T) {
	root := NewCell()
	for i := 0; i < 2; i++ {
		ref, err := root.NewRef()
		if err != nil {
			t.Fatalf("NewRef() failed: %v", err)
		}
		if err := ref.WriteUint(uint64(i), 32); err != nil {
			t.Fatalf("WriteUint() failed: %v", err)
		}
	}
	rootHash, err := root.Hash256()
	if err != nil {
		t.Fatalf("Hash256() failed: %v", err)
	}

	prover, err := NewMerkleProver(root)
	if err != nil {
		t.Fatalf("NewMerkleProver() failed: %v", err)
	}
	cursor := prover.Cursor()
	cursor.Ref(1).Prune()
	proof, err := prover.CreateProof(cursor)
	if err != nil {
		t.Fatalf("CreateProof() failed: %v", err)
	}

	masks := declaredMasks(t, proof)
	if masks[0] != 0 {
		t.Errorf("merkle proof root declares level mask %v, want 0", masks[0])
	}

	proofRoot, err := DeserializeSingleRootBoc(proof)
	if err != nil {
		t.Fatalf("DeserializeSingleRootBoc() failed: %v", err)
	}
	if proofRoot.CellType() != MerkleProofCell {
		t.Fatalf("proof root type = %v, want MerkleProofCell", proofRoot.CellType())
	}
	if level := proofRoot.Level(); level != 0 {
		t.Errorf("proof root level = %v, want 0", level)
	}
	virtualRoot := proofRoot.Refs()[0]
	if level := virtualRoot.Level(); level != 1 {
		t.Errorf("virtual root level = %v, want 1", level)
	}
	// the proof must still be about the original tree
	merkleRoot, err := proofRoot.GetMerkleRoot()
	if err != nil {
		t.Fatalf("GetMerkleRoot() failed: %v", err)
	}
	if merkleRoot != rootHash {
		t.Errorf("merkle root = %x, want %x", merkleRoot, rootHash)
	}
	if hash, err := virtualRoot.HashAtLevel(0); err != nil || hash != rootHash {
		t.Errorf("virtual root hash at level 0 = %x (err %v), want %x", hash, err, rootHash)
	}
}

func TestLevelMaskOf(t *testing.T) {
	prunedBits := 16 + (hashSize+depthSize)*8
	testCases := []struct {
		name     string
		cellType CellType
		bitsBuf  []byte
		bitsLen  int
		refs     []levelMask
		want     levelMask
		wantErr  bool
	}{
		{name: "ordinary without refs", cellType: OrdinaryCell, want: 0},
		{name: "ordinary is a union of its refs", cellType: OrdinaryCell, refs: []levelMask{0b001, 0b100}, want: 0b101},
		{name: "library is always level 0", cellType: LibraryCell, bitsBuf: []byte{2}, bitsLen: 264, want: 0},
		{name: "merkle proof shifts its ref right", cellType: MerkleProofCell, refs: []levelMask{0b011}, want: 0b001},
		{name: "merkle proof over a level 1 ref is level 0", cellType: MerkleProofCell, refs: []levelMask{0b001}, want: 0},
		{name: "merkle proof without a ref", cellType: MerkleProofCell, wantErr: true},
		{name: "merkle update shifts the union right", cellType: MerkleUpdateCell, refs: []levelMask{0b010, 0b101}, want: 0b011},
		{name: "merkle update with a single ref", cellType: MerkleUpdateCell, refs: []levelMask{0b010}, wantErr: true},
		{name: "pruned branch takes its mask from the data", cellType: PrunedBranchCell, bitsBuf: []byte{1, 1}, bitsLen: prunedBits, want: 1},
		{name: "pruned branch with a zero mask", cellType: PrunedBranchCell, bitsBuf: []byte{1, 0}, bitsLen: prunedBits, wantErr: true},
		{name: "pruned branch with an out of range mask", cellType: PrunedBranchCell, bitsBuf: []byte{1, 0b1000}, bitsLen: prunedBits, wantErr: true},
		{name: "pruned branch with a truncated body", cellType: PrunedBranchCell, bitsBuf: []byte{1, 1}, bitsLen: prunedBits - 8, wantErr: true},
		{name: "empty pruned branch", cellType: PrunedBranchCell, wantErr: true},
		{name: "pruned branch with refs", cellType: PrunedBranchCell, bitsBuf: []byte{1, 1}, bitsLen: prunedBits, refs: []levelMask{0}, wantErr: true},
		{name: "unknown cell type", cellType: CellType(42), wantErr: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mask, err := levelMaskOf(tc.cellType, tc.bitsBuf, tc.bitsLen, tc.refs)
			if tc.wantErr {
				if !errors.Is(err, ErrMalformedExoticCell) {
					t.Fatalf("levelMaskOf() error = %v, want ErrMalformedExoticCell", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("levelMaskOf() failed: %v", err)
			}
			if mask != tc.want {
				t.Errorf("levelMaskOf() = %b, want %b", mask, tc.want)
			}
		})
	}
}

// TestMalformedExoticCell makes sure a broken exotic cell is reported
// instead of producing a bogus hash or panicking on an out of range read.
func TestMalformedExoticCell(t *testing.T) {
	truncated := NewCellExotic(PrunedBranchCell)
	if err := truncated.WriteUint(0x0101, 16); err != nil {
		t.Fatalf("WriteUint() failed: %v", err)
	}
	truncated.ResetCounters()

	if _, err := truncated.Hash256(); !errors.Is(err, ErrMalformedExoticCell) {
		t.Errorf("Hash256() error = %v, want ErrMalformedExoticCell", err)
	}
	if _, err := truncated.ToBoc(); !errors.Is(err, ErrMalformedExoticCell) {
		t.Errorf("ToBoc() error = %v, want ErrMalformedExoticCell", err)
	}
	if _, err := truncated.LevelMask(); !errors.Is(err, ErrMalformedExoticCell) {
		t.Errorf("LevelMask() error = %v, want ErrMalformedExoticCell", err)
	}

	// and it must be reported through its parents too
	parent := NewCell()
	if err := parent.AddRef(truncated); err != nil {
		t.Fatalf("AddRef() failed: %v", err)
	}
	if _, err := parent.ToBoc(); !errors.Is(err, ErrMalformedExoticCell) {
		t.Errorf("ToBoc() error = %v, want ErrMalformedExoticCell", err)
	}
}
