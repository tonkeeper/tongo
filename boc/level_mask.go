package boc

import (
	"fmt"
	"math/bits"
)

// levelMask is a tricky way to keep track of two things simultaneously:
// a cell level and a number of different hashes that makes sense to calculate for a cell.
type levelMask uint32

func (m levelMask) Level() int {
	return 32 - bits.LeadingZeros32(uint32(m))
}

func (m levelMask) HashIndex() int {
	return bits.OnesCount32(uint32(m))
}

func (m levelMask) HashesCount() int {
	return m.HashIndex() + 1
}

func (m levelMask) Apply(level int) levelMask {
	return levelMask(uint32(m) & ((1 << uint32(level)) - 1))
}

func (m levelMask) IsSignificant(level uint32) bool {
	if level == 0 {
		return true
	}
	return (m>>(level-1))%2 != 0
}

// maxLevelMask is the largest level mask that fits into the 3 mask bits of a d1 byte.
const maxLevelMask = levelMask(0b111)

// levelMaskOf derives a cell's level mask from its type, its data and the level masks of
// its refs, following ton-blockchain's DataCell::create:
// https://github.com/ton-blockchain/ton/blob/master/crypto/vm/cells/CellBuilder.cpp
//
// A cell's level mask is never stored on the cell itself: it is a function of the cell's
// content, so deriving it is the only way to keep it in sync with the tree. In particular
// an ordinary parent inherits the union of its children's masks, while a merkle cell
// shifts them right - a plain union would produce a wrong level for merkle proofs.
func levelMaskOf(cellType CellType, bitsBuf []byte, bitsLen int, refs []levelMask) (levelMask, error) {
	switch cellType {
	case OrdinaryCell:
		var mask levelMask
		for _, ref := range refs {
			mask |= ref
		}
		return mask, nil
	case PrunedBranchCell:
		if len(refs) != 0 {
			return 0, fmt.Errorf("%w: pruned branch cell must have no refs, got %v", ErrMalformedExoticCell, len(refs))
		}
		if bitsLen < 16 || len(bitsBuf) < 2 {
			return 0, fmt.Errorf("%w: pruned branch cell is too short: %v bits", ErrMalformedExoticCell, bitsLen)
		}
		// the second byte of a pruned branch cell stores the level mask of the cell it replaces.
		mask := levelMask(bitsBuf[1])
		if mask == 0 || mask > maxLevelMask {
			return 0, fmt.Errorf("%w: invalid pruned branch level mask %v", ErrMalformedExoticCell, mask)
		}
		// a pruned branch cell stores a hash and a depth for every significant level
		// of the original cell, on top of the two header bytes.
		if want := 16 + (hashSize+depthSize)*8*mask.HashIndex(); bitsLen != want {
			return 0, fmt.Errorf("%w: pruned branch cell with mask %v must have %v bits, got %v", ErrMalformedExoticCell, mask, want, bitsLen)
		}
		return mask, nil
	case LibraryCell:
		return 0, nil
	case MerkleProofCell:
		if len(refs) != 1 {
			return 0, fmt.Errorf("%w: merkle proof cell must have exactly one ref, got %v", ErrMalformedExoticCell, len(refs))
		}
		return refs[0] >> 1, nil
	case MerkleUpdateCell:
		if len(refs) != 2 {
			return 0, fmt.Errorf("%w: merkle update cell must have exactly two refs, got %v", ErrMalformedExoticCell, len(refs))
		}
		return (refs[0] | refs[1]) >> 1, nil
	}
	return 0, fmt.Errorf("%w: unknown cell type %v", ErrMalformedExoticCell, uint8(cellType))
}
