package boc

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"math/bits"
)

var reachBocMagicPrefix = []byte{
	0xb5, 0xee, 0x9c, 0x72,
}

var leanBocMagicPrefix = []byte{
	0x68, 0xff, 0x65, 0xf3,
}

var leanBocMagicPrefixCRC = []byte{
	0xac, 0xc3, 0xa7, 0x28,
}

const hashSize = 32
const depthSize = 2
const maxLevel = 3
const maxDepth = 1024
const maxCellWhs = 64

var crcTable = crc32.MakeTable(crc32.Castagnoli)

func readNBytesUIntFromArray(n int, arr []byte) uint {
	var res uint = 0
	for i := 0; i < n; i++ {
		res *= 256
		res += uint(arr[i])
	}
	return res
}

type bocHeader struct {
	hasIdx       bool
	hasCrc32     bool
	hasCacheBits bool
	flags        int
	sizeBytes    int
	cellCount    uint
	rootCount    uint
	absentCount  uint
	totCellsSize uint
	rootList     []uint
	index        []uint
	cellsData    []byte
}

func parseBocHeader(boc []byte) (*bocHeader, error) {
	if len(boc) < 4+1 {
		return nil, errors.New("not enough bytes for magic prefix")
	}
	checkSum := crc32.Checksum(boc[0:len(boc)-4], crcTable)

	var prefix = boc[0:4]
	boc = boc[4:]

	var (
		hasIdx       bool
		hashCrc32    bool
		hasCacheBits bool
		flags        int
		sizeBytes    int
	)

	if bytes.Equal(prefix, reachBocMagicPrefix) {
		var flagsByte = boc[0]
		hasIdx = (flagsByte & 128) > 0
		hashCrc32 = (flagsByte & 64) > 0
		hasCacheBits = (flagsByte & 32) > 0
		flags = int((flagsByte&16)*2 + (flagsByte & 8))
		sizeBytes = int(flagsByte % 8)
	} else if bytes.Equal(prefix, leanBocMagicPrefix) {
		hasIdx = true
		hashCrc32 = false
		hasCacheBits = false
		flags = 0
		sizeBytes = int(boc[0])
	} else if bytes.Equal(prefix, leanBocMagicPrefixCRC) {
		hasIdx = true
		hashCrc32 = true
		hasCacheBits = false
		flags = 0
		sizeBytes = int(boc[0])
	} else {
		return nil, errors.New("unknown magic prefix")
	}

	boc = boc[1:]
	if len(boc) < 1+5*sizeBytes {
		return nil, errors.New("not enough bytes for encoding cells counters")
	}

	offsetBytes := int(boc[0])
	boc = boc[1:]
	cellsCount := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	rootsCount := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	absentNum := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	totCellsSize := readNBytesUIntFromArray(offsetBytes, boc)
	boc = boc[offsetBytes:]

	if len(boc) < int(rootsCount)*sizeBytes {
		return nil, errors.New("not enough bytes for encoding root cells hashes")
	}

	// Roots
	rootList := make([]uint, 0, rootsCount)
	for i := 0; i < int(rootsCount); i++ {
		rootList = append(rootList, readNBytesUIntFromArray(sizeBytes, boc))
		boc = boc[sizeBytes:]
	}

	// Index
	index := make([]uint, 0, cellsCount)
	if hasIdx {
		if len(boc) < offsetBytes*int(cellsCount) {
			return nil, errors.New("not enough bytes for index encoding")
		}
		for i := 0; i < int(cellsCount); i++ {
			val := readNBytesUIntFromArray(offsetBytes, boc)
			if hasCacheBits {
				val /= 2
			}
			index = append(index, val)
			boc = boc[offsetBytes:]
		}
	}

	// Cells
	if len(boc) < int(totCellsSize) {
		return nil, errors.New("not enough bytes for cells data")
	}

	cellsData := boc[0:totCellsSize]
	boc = boc[totCellsSize:]

	if hashCrc32 {
		if len(boc) < 4 {
			return nil, errors.New("not enough bytes for crc32c hashsum")
		}
		if binary.LittleEndian.Uint32(boc[0:4]) != checkSum {
			return nil, errors.New("crc32c hashsum mismatch")
		}
		boc = boc[4:]
	}

	if len(boc) > 0 {
		return nil, errors.New("too much bytes in provided boc")
	}

	return &bocHeader{
		hasIdx:       hasIdx,
		hasCrc32:     hashCrc32,
		hasCacheBits: hasCacheBits,
		flags:        flags,
		sizeBytes:    sizeBytes,
		cellCount:    cellsCount,
		rootCount:    rootsCount,
		absentCount:  absentNum,
		totCellsSize: totCellsSize,
		rootList:     rootList,
		index:        index,
		cellsData:    cellsData,
	}, nil
}

func deserializeCellData(cellData []byte, referenceIndexSize int) (*Cell, []int, []byte, error) {
	if len(cellData) < 2 {
		return nil, nil, nil, errors.New("not enough bytes to encode cell descriptors")
	}

	d1 := cellData[0]
	d2 := cellData[1]
	cellData = cellData[2:]

	isExotic := (d1 & 8) > 0
	refNum := int(d1 % 8)

	dataBytesSize := int(d2>>1) + int(d2%2)
	fullfilledBytes := !((d2 % 2) > 0)
	withHashes := (d1 & 0b10000) != 0
	mask := levelMask(d1 >> 5)

	if withHashes {
		offset := mask.HashesCount() * (hashSize + depthSize)
		cellData = cellData[offset:]
	}
	var cell *Cell
	if isExotic {
		// the first byte of an exotic cell stores the cell's type.
		exoticType := CellType(readNBytesUIntFromArray(1, cellData))
		cell = NewCellExotic(exoticType)
		cell.mask = mask
	} else {
		cell = NewCell()
		cell.mask = mask
	}

	if len(cellData) < dataBytesSize+referenceIndexSize*refNum {
		return nil, nil, nil, errors.New("not enough bytes to encode cell data")
	}

	err := cell.setTopUppedArray(cellData[0:dataBytesSize], fullfilledBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	cellData = cellData[dataBytesSize:]

	refs := make([]int, 0, refNum)
	for i := 0; i < refNum; i++ {
		refs = append(refs, int(readNBytesUIntFromArray(referenceIndexSize, cellData)))
		cellData = cellData[referenceIndexSize:]
	}

	return cell, refs, cellData, nil
}

func DeserializeBoc(boc []byte) ([]*Cell, error) {
	header, err := parseBocHeader(boc)
	if err != nil {
		return nil, err
	}
	cellsData := header.cellsData
	cellsArray := make([]*Cell, 0, header.cellCount)
	refsArray := make([][]int, 0, header.cellCount)

	for i := 0; i < int(header.cellCount); i++ {
		cell, refs, residue, err := deserializeCellData(cellsData, header.sizeBytes)
		if err != nil {
			return nil, err
		}
		cellsData = residue
		cellsArray = append(cellsArray, cell)
		refsArray = append(refsArray, refs)
	}
	for i := int(header.cellCount - 1); i >= 0; i-- {
		c := refsArray[i]
		if len(c) > 4 {
			return nil, fmt.Errorf("too long refs array")
		}
		for ri, r := range c {
			if r < i {
				return nil, errors.New("topological order is broken")
			}
			if r >= len(cellsArray) {
				return nil, errors.New("index out of range for boc deserialization")
			}
			cellsArray[i].refs[ri] = cellsArray[r]
		}
	}

	rootCells := make([]*Cell, 0, len(header.rootList))
	for _, item := range header.rootList {
		rootCells = append(rootCells, cellsArray[item])
	}
	return rootCells, nil
}

func DeserializeBocBase64(boc string) ([]*Cell, error) {
	bocData, err := base64.StdEncoding.DecodeString(boc)
	if err != nil {
		return nil, err
	}
	return DeserializeBoc(bocData)
}

func DeserializeSinglRootBase64(boc string) (*Cell, error) {
	cells, err := DeserializeBocBase64(boc)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid boc roots number %v", len(cells))
	}
	return cells[0], nil
}

func DeserializeBocHex(boc string) ([]*Cell, error) {
	bocData, err := hex.DecodeString(boc)
	if err != nil {
		return nil, err
	}
	return DeserializeBoc(bocData)
}

func SerializeBoc(cell *Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	bag := newBagOfCells()
	return bag.serializeBoc([]*Cell{cell}, idx, hasCrc32, cacheBits, flags)
}

// bagOfCells serializes cells to a boc.
//
// the serialization algorithms is a golang version of
// https://github.com/ton-blockchain/ton/blob/master/crypto/vm/boc.cpp#
type bagOfCells struct {
	hasher *Hasher
}

func newBagOfCells() *bagOfCells {
	return &bagOfCells{
		hasher: NewHasher(),
	}
}

// serializeBoc converts the given list of root cells to a byte representation.
//
//	serialized_boc#672fb0ac has_idx:(## 1) has_crc32c:(## 1)
//	has_cache_bits:(## 1) flags:(## 2) { flags = 0 }
//	size:(## 3) { size <= 4 }
//	off_bytes:(## 8) { off_bytes <= 8 }
//	cells:(##(size * 8))
//	roots:(##(size * 8))
//	absent:(##(size * 8)) { roots + absent <= cells }
//	tot_cells_size:(##(off_bytes * 8))
//	index:(cells * ##(off_bytes * 8))
//	cell_data:(tot_cells_size * [ uint8 ])
//	= BagOfCells;
func (boc *bagOfCells) serializeBoc(rootCells []*Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	roots, cellInfos, err := boc.importRoots(rootCells)
	if err != nil {
		return nil, err
	}
	cellCount := len(cellInfos)
	refBitSize := bits.Len(uint(cellCount))
	refByteSize := int(math.Max(math.Ceil(float64(refBitSize)/8), 1))
	offset := uint(0)
	offsets := make([]uint, cellCount)
	reps := make([][]byte, cellCount)
	for i := cellCount - 1; i >= 0; i-- {
		ci := cellInfos[i]
		repr := ci.cell.bocReprWithoutRefs(ci.cell.mask)
		for j := 0; j < ci.refsNumber; j++ {
			var b [8]byte
			k := cellCount - 1 - ci.refsIndex[j]
			binary.BigEndian.PutUint64(b[:], uint64(k))
			repr = append(repr, b[8-refByteSize:]...)
		}
		reps[i] = repr
		offset += uint(len(repr))
		fixedOffset := offset
		if cacheBits {
			fixedOffset = offset * 2
			if ci.shouldCache {
				fixedOffset += 1
			}
		}
		offsets[i] = fixedOffset
	}

	offsetBitSize := bits.Len(offset)
	offsetByteSize := int(math.Max(math.Ceil(float64(offsetBitSize)/8), 1))

	bitString := NewBitString((CellBits + 32*4 + 32*3) * int(cellCount))
	err = bitString.WriteBytes(reachBocMagicPrefix)
	if err != nil {
		return nil, err
	}
	err = bitString.WriteBitArray([]bool{idx, hasCrc32, cacheBits})
	if err != nil {
		return nil, err
	}
	err = bitString.WriteUint(uint64(flags), 2)
	if err != nil {
		return nil, err
	}
	err = bitString.WriteInt(int64(refByteSize), 3)
	if err != nil {
		return nil, err
	}
	err = bitString.WriteInt(int64(offsetByteSize), 8)
	if err != nil {
		return nil, err
	}
	err = bitString.WriteUint(uint64(cellCount), refByteSize*8)
	if err != nil {
		return nil, err
	}
	err = bitString.WriteUint(uint64(len(roots)), refByteSize*8)
	if err != nil {
		return nil, err
	}
	// absent count
	err = bitString.WriteUint(0, refByteSize*8)
	if err != nil {
		return nil, err
	}
	// total cells size
	err = bitString.WriteUint(uint64(offset), offsetByteSize*8)
	if err != nil {
		return nil, err
	}
	for _, root := range roots {
		k := len(cellInfos) - 1 - root.index
		err = bitString.WriteUint(uint64(k), refByteSize*8)
		if err != nil {
			return nil, err
		}
	}

	if idx {
		for i := cellCount - 1; i >= 0; i-- {
			err = bitString.WriteUint(uint64(offsets[i]), offsetByteSize*8)
			if err != nil {
				return nil, err
			}
		}
	}

	for i := cellCount - 1; i >= 0; i-- {
		err = bitString.WriteBytes(reps[i])
		if err != nil {
			return nil, err
		}
	}
	resBytes, err := bitString.GetTopUppedArray()
	if err != nil {
		return nil, err
	}
	if hasCrc32 {
		checksum := make([]byte, 4)
		binary.LittleEndian.PutUint32(checksum, crc32.Checksum(resBytes, crc32.MakeTable(crc32.Castagnoli)))
		resBytes = append(resBytes, checksum...)
	}
	return resBytes, nil
}

func (boc *bagOfCells) importRoots(rootCells []*Cell) ([]*rootInfo, []*cellInfo, error) {
	roots := make([]*rootInfo, 0, len(rootCells))
	state := &orderState{
		cells: map[string]int{},
	}
	for _, root := range rootCells {
		pos, err := boc.importCell(state, root, 0)
		if err != nil {
			return nil, nil, err
		}
		info := rootInfo{
			cell:  root,
			index: pos,
		}
		roots = append(roots, &info)
	}
	cellInfos := boc.reorderCells(roots, state)
	return roots, cellInfos, nil
}

func (boc *bagOfCells) importCell(state *orderState, cell *Cell, depth int) (int, error) {
	if depth > maxDepth {
		return 0, ErrDepthIsTooBig
	}
	if cell == nil {
		return 0, errors.New("failed to import nil cell")
	}
	hash, err := boc.hasher.HashString(cell)
	if err != nil {
		return 0, err
	}
	if pos, ok := state.cells[hash]; ok {
		state.cellList[pos].shouldCache = true
		return pos, nil
	}
	var refs [4]int
	sumChildWt := 1
	refsNumber := 0
	for i, ref := range cell.refs {
		if ref == nil {
			break
		}
		refsNumber += 1
		refPos, err := boc.importCell(state, ref, depth+1)
		if err != nil {
			return 0, err
		}
		refs[i] = refPos
		sumChildWt += state.cellList[refPos].wt
		state.intRefs += 1
	}
	state.cells[hash] = len(state.cellList)
	if sumChildWt > 255 {
		sumChildWt = 255
	}
	info := cellInfo{
		cell:        cell,
		shouldCache: false,
		wt:          sumChildWt,
		refsIndex:   refs,
		refsNumber:  refsNumber,
		hashCount:   cell.mask.HashesCount(),
		newIndex:    -1,
		isRootCell:  false,
	}
	return state.add(&info), nil
}

// cellInfo holds information required to come up with an order to serialize boc.
type cellInfo struct {
	cell        *Cell
	shouldCache bool
	wt          int
	refsIndex   [4]int
	refsNumber  int
	hashCount   int
	newIndex    int
	isRootCell  bool
}

func (ci *cellInfo) isSpecial() bool {
	return ci.wt == 0
}

type rootInfo struct {
	cell  *Cell
	index int
}

type orderState struct {
	cellList []*cellInfo
	cells    map[string]int
	intRefs  int
}

func (bag *orderState) add(ci *cellInfo) int {
	bag.cellList = append(bag.cellList, ci)
	return len(bag.cellList) - 1
}

type force int

const (
	previsit force = iota
	visit
	allocate
)

func (boc *bagOfCells) revisit(state, newState *orderState, cellIndex int, force force) int {
	dci := state.cellList[cellIndex]
	if dci.newIndex >= 0 {
		return dci.newIndex
	}
	if force == previsit {
		if dci.newIndex != -1 {
			// already previsited or visited
			return dci.newIndex
		}
		for j := dci.refsNumber - 1; j >= 0; j-- {
			childIndex := dci.refsIndex[j]
			childCell := state.cellList[childIndex]
			if childCell.isSpecial() {
				boc.revisit(state, newState, childIndex, visit)
			} else {
				boc.revisit(state, newState, childIndex, previsit)
			}
		}
		dci.newIndex = -2
		return dci.newIndex
	}
	if force == allocate {
		dci.newIndex = newState.add(dci)
		return dci.newIndex
	}

	if dci.newIndex == -3 {
		return dci.newIndex
	}
	if dci.isSpecial() {
		boc.revisit(state, newState, cellIndex, previsit)
	}

	for j := dci.refsNumber - 1; j >= 0; j-- {
		boc.revisit(state, newState, dci.refsIndex[j], visit)
	}
	for j := dci.refsNumber - 1; j >= 0; j-- {
		dci.refsIndex[j] = boc.revisit(state, newState, dci.refsIndex[j], allocate)
	}
	dci.newIndex = -3
	return dci.newIndex

}
func (boc *bagOfCells) reorderCells(roots []*rootInfo, state *orderState) []*cellInfo {
	for i := len(state.cellList) - 1; i >= 0; i-- {
		dci := state.cellList[i]
		c := dci.refsNumber
		sum := maxCellWhs - 1
		mask := 0
		for j := 0; j < dci.refsNumber; j++ {
			dcj := state.cellList[dci.refsIndex[j]]
			limit := (maxCellWhs - 1 + j) / dci.refsNumber
			if dcj.wt <= limit {
				sum -= dcj.wt
				c--
				mask |= 1 << j
			}
		}
		if c > 0 {
			for j := 0; j < dci.refsNumber; j++ {
				if mask&(1<<j) == 0 {
					dcj := state.cellList[dci.refsIndex[j]]
					sum += 1
					limit := sum / c
					if dcj.wt > limit {
						dcj.wt = limit
					}
				}
			}
		}
	}
	// intHashes will be used to preallocate all required space to store a boc.
	// currently, not used.
	intHashes := 0
	for _, dci := range state.cellList {
		sum := 1
		for j := 0; j < dci.refsNumber; j++ {
			sum += state.cellList[dci.refsIndex[j]].wt
		}
		if sum <= dci.wt {
			dci.wt = sum
		} else {
			dci.wt = 0
			intHashes += dci.hashCount
		}
	}
	// topHashes will be used to preallocate space required to store a boc.
	// currently, not used.
	topHashes := 0
	for _, root := range roots {
		ci := state.cellList[root.index]
		ci.isRootCell = true
		if ci.wt > 0 {
			topHashes += ci.hashCount
		}
	}

	newState := &orderState{}
	if len(state.cellList) > 0 {
		for _, root := range roots {
			boc.revisit(state, newState, root.index, previsit)
			boc.revisit(state, newState, root.index, visit)
		}
		for _, root := range roots {
			boc.revisit(state, newState, root.index, allocate)
		}
		for _, root := range roots {
			root.index = state.cellList[root.index].newIndex
		}
	}
	return newState.cellList
}
