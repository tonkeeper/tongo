package boc

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// immutableCell provides a convenient way to calculate a cell's hash and depth.
type immutableCell struct {
	mask     levelMask
	cellType CellType
	refs     []*immutableCell
	hashes   [][]byte
	depths   []int

	// bitsBuf points to the content of this cell.
	// Used if this cell is a pruned branch cell to extract hashes and depth of an original cell.
	bitsBuf []byte
	// bitsLen is a length of the content of this cell in bits.
	// it's not always equal to len(bitsBuf) because bitsBuf can contain more bytes than the actual content.
	bitsLen int
}

// newImmutableCell returns a new instance of immutable cell.
// cache can't be nil because it helps to avoid an endless loop in case of the given cell contains a fork bomb.
func newImmutableCell(c *Cell, cache map[*Cell]*immutableCell) (*immutableCell, error) {
	if imm, ok := cache[c]; ok {
		return imm, nil
	}
	imm := &immutableCell{
		mask:     c.mask,
		cellType: c.cellType,
		refs:     make([]*immutableCell, 0, c.RefsSize()),
		hashes:   make([][]byte, 0, c.mask.HashesCount()),
		depths:   make([]int, 0, c.mask.HashesCount()),
		bitsBuf:  c.bits.buf,
		bitsLen:  c.bits.len,
	}
	for _, ref := range c.refs {
		if ref == nil {
			break
		}
		immRef, err := newImmutableCell(ref, cache)
		if err != nil {
			return nil, err
		}
		imm.refs = append(imm.refs, immRef)
	}

	mask := c.mask
	level := mask.Level()

	// the algorithm below is taken directly from ton-blockchain/ton
	// https://github.com/ton-blockchain/ton/blob/v2022.12/crypto/vm/cells/DataCell.cpp#L215

	// "offset" is a workaround for a pruned branch cell.
	// When the TON runtime creates a merkle proof,
	// it can replace some cells inside a cell tree with pruned branch cells.
	// From the perspective of a merkle tree,
	// a pruned branch cell must behave similarly to a cell it replaces,
	// providing the same hashes and depths for all significant levels of the original cell.
	// Thus, the runtime stores all hashes and depths of the original cell inside a pruned branch cell.
	//
	// This means two things:
	// 1. when we request a hash(depth) of a pruned branch cell,
	//    we return either the pruned branch cell's hash or the original cell's hash depending on the requested level.
	// 2. when we calculate hashes of a pruned branch cell,
	//    we skip several first levels because they were stored inside the pruned branch cell.
	//
	// An example.
	// A cell has a level mask "10" (binary), meaning the cell has two hashes and its level is 2.
	// A pruned branch cell that replaces it will have a level mask "110"(binary),
	// meaning it has three hashes and its level is 3.
	// For hashes of levels from 0 to 2, we return the original cell's hashes stored inside the pruned branch cell.
	// A hash of level 3 is the only hash of the pruned branch cell.
	offset := 0
	if c.cellType == PrunedBranchCell {
		offset = mask.HashIndex()
	}

	hashIndex := -1
	for i := 0; i <= level; i++ {
		if !mask.IsSignificant(uint32(i)) {
			// we care about significant hashes only
			continue
		}
		hashIndex += 1
		if hashIndex < offset {
			// happens with pruned cells only
			continue
		}

		x := sha256.New()
		if hashIndex == offset {
			// either i=0 or cellType=PrunedBranchCell
			cellRepr := c.bocReprWithoutRefs(c.mask.Apply(i))
			x.Write(cellRepr)
		} else {
			// i>0
			x.Write([]byte{d1(c, c.mask.Apply(i)), d2(c)})
			x.Write(imm.hashes[hashIndex-offset-1])
		}

		childLevelIndex := i
		if c.cellType == MerkleProofCell || c.cellType == MerkleUpdateCell {
			childLevelIndex = i + 1
		}
		depth := 0
		for _, ref := range imm.refs {
			childDepth := ref.Depth(childLevelIndex)
			var depthRepr [2]byte
			binary.BigEndian.PutUint16(depthRepr[:], uint16(childDepth))
			x.Write(depthRepr[:])
			if depth < childDepth {
				depth = childDepth
			}
		}
		if len(imm.refs) > 0 {
			if depth >= maxDepth {
				return nil, ErrDepthIsTooBig
			}
			depth += 1
		}
		imm.depths = append(imm.depths, depth)
		for _, ref := range imm.refs {
			x.Write(ref.Hash(childLevelIndex))
		}

		imm.hashes = append(imm.hashes, x.Sum(nil))
	}
	cache[c] = imm
	return imm, nil
}

// Hash return a hash as a byte slice for the given level.
func (ic *immutableCell) Hash(level int) []byte {
	index := ic.mask.Apply(level).HashIndex()
	if ic.cellType == PrunedBranchCell {
		//offset := ic.mask.HashIndex()
		//if index != offset {
		return ic.bitsBuf[2+(index)*32 : 2+(index+1)*32]
		//}
		//index = 0
	}
	return ic.hashes[index]
}

// Depth returns a depth for the given level.
func (ic *immutableCell) Depth(level int) int {
	index := ic.mask.Apply(level).HashIndex()
	if ic.cellType == PrunedBranchCell {
		offset := ic.mask.HashIndex()
		if index != offset {
			depth := readNBytesUIntFromArray(2, ic.bitsBuf[2+32*offset+index*2:])
			return int(depth)
		}
		index = 0
	}
	return ic.depths[index]
}

// pruneCells return the current subtree (which this cell represents) with pruned cells.
// if this cell is pruned, "pruneCells" returns a new pruned branch cell instead of this cell.
// As of now, this function doesn't work with MerkleProofCell and MerkleUpdateCell.
func (ic *immutableCell) pruneCells(pruned map[*immutableCell]struct{}) (*Cell, error) {
	if ic.cellType == MerkleProofCell || ic.cellType == MerkleUpdateCell {
		return nil, fmt.Errorf("unsupported cell type: %v", ic.cellType)
	}
	if _, ok := pruned[ic]; ok {
		// we are pruned
		// let's replace this cell with a pruned branch cell
		prunedCell := NewCell()
		prunedCell.mask = 1
		prunedCell.cellType = PrunedBranchCell
		if err := prunedCell.WriteUint(1, 8); err != nil {
			return nil, err
		}
		if err := prunedCell.WriteUint(1, 8); err != nil {
			return nil, err
		}
		if err := prunedCell.WriteBytes(ic.Hash(0)); err != nil {
			return nil, err
		}
		if err := prunedCell.WriteUint(uint64(ic.Depth(0)), 16); err != nil {
			return nil, err
		}
		prunedCell.ResetCounters()
		return prunedCell, nil
	}
	// all good,
	// going down the tree
	bits := BitString{
		buf: ic.bitsBuf,
		cap: ic.bitsLen,
		len: ic.bitsLen,
	}
	res := Cell{
		bits:     bits,
		refs:     [4]*Cell{},
		cellType: ic.cellType,
		mask:     ic.mask,
	}
	mask := ic.mask
	for i, ref := range ic.refs {
		cell, err := ref.pruneCells(pruned)
		if err != nil {
			return nil, err
		}
		if cell.mask > 0 {
			mask |= cell.mask
		}
		res.refs[i] = cell
	}
	res.mask = mask
	return &res, nil
}
