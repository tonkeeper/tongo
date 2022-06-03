package boc

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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

var crcTable = crc32.MakeTable(crc32.Castagnoli)

func ByteArrayEquals(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

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
	hashCrc32    bool
	hasCacheBits bool
	flags        int
	sizeBytes    int
	cellsNum     uint
	rootsNum     uint
	absentNum    uint
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

	if ByteArrayEquals(prefix, reachBocMagicPrefix) {
		var flagsByte = boc[0]
		hasIdx = (flagsByte & 128) > 0
		hashCrc32 = (flagsByte & 64) > 0
		hasCacheBits = (flagsByte & 32) > 0
		flags = int((flagsByte&16)*2 + (flagsByte & 8))
		sizeBytes = int(flagsByte % 8)
	} else if ByteArrayEquals(prefix, leanBocMagicPrefix) {
		hasIdx = true
		hashCrc32 = false
		hasCacheBits = false
		flags = 0
		sizeBytes = int(boc[0])
	} else if ByteArrayEquals(prefix, leanBocMagicPrefixCRC) {
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
	cellsNum := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	rootsNum := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	absentNum := readNBytesUIntFromArray(sizeBytes, boc)
	boc = boc[sizeBytes:]
	totCellsSize := readNBytesUIntFromArray(offsetBytes, boc)
	boc = boc[offsetBytes:]

	if len(boc) < int(rootsNum)*sizeBytes {
		return nil, errors.New("not enough bytes for encoding root cells hashes")
	}

	// Roots
	rootList := make([]uint, 0)
	for i := 0; i < int(rootsNum); i++ {
		rootList = append(rootList, readNBytesUIntFromArray(sizeBytes, boc))
		boc = boc[sizeBytes:]
	}

	// Index
	index := make([]uint, 0)
	if hasIdx {
		if len(boc) < offsetBytes*int(cellsNum) {
			return nil, errors.New("not enough bytes for index encoding")
		}
		for i := 0; i < int(cellsNum); i++ {
			index = append(index, readNBytesUIntFromArray(offsetBytes, boc))
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
		hasIdx,
		hashCrc32,
		hasCacheBits,
		flags,
		sizeBytes,
		cellsNum,
		rootsNum,
		absentNum,
		totCellsSize,
		rootList,
		index,
		cellsData,
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
	dataBytesSize := int(math.Ceil(float64(d2) / 2))
	fullfilledBytes := !((d2 % 2) > 0)

	var cell *Cell
	if isExotic {
		cell = NewCellExotic()
	} else {
		cell = NewCell()
	}
	var refs = make([]int, 0)

	if len(cellData) < dataBytesSize+referenceIndexSize*refNum {
		return nil, nil, nil, errors.New("not enough bytes to encode cell data")
	}

	err := cell.Bits.SetTopUppedArray(cellData[0:dataBytesSize], fullfilledBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	cellData = cellData[dataBytesSize:]

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
	cellsArray := make([]*Cell, 0)
	refsArray := make([][]int, 0)

	for i := 0; i < int(header.cellsNum); i++ {
		cell, refs, residue, err := deserializeCellData(cellsData, header.sizeBytes)
		if err != nil {
			return nil, err
		}
		cellsData = residue
		cellsArray = append(cellsArray, cell)
		refsArray = append(refsArray, refs)
	}

	for i := int(header.cellsNum - 1); i >= 0; i-- {
		c := refsArray[i]
		if len(c) > 4 {
			return nil, fmt.Errorf("too long refs array")
		}
		for ri, r := range c {
			if r < i {
				return nil, errors.New("topological order is broken")
			}
			if r >= len(cellsArray) {
				return nil, errors.New("index out of range for boc deseriailization")
			}
			cellsArray[i].refs[ri] = cellsArray[r]
		}
	}

	rootCells := make([]*Cell, 0)

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

func getMaxDepth(cell *Cell, iterationCounter *int) (int, error) {
	if *iterationCounter == 0 {
		return 0, errors.New("to big boc for processing")
	}
	*iterationCounter -= 1
	maxDepth := 0
	if cell.RefsSize() > 0 {
		for _, ref := range cell.Refs() {
			depth, err := getMaxDepth(ref, iterationCounter)
			if err != nil {
				return 0, err
			}
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		maxDepth++
	}
	return maxDepth, nil
}

func bocReprWithoutRefs(cell *Cell) []byte {
	d1 := byte(cell.RefsSize())
	d2 := byte((cell.BitSize()+7)/8 + cell.BitSize()/8)

	res := make([]byte, ((cell.BitSize()+7)/8)+2)
	res[0] = d1
	res[1] = d2
	copy(res[2:], cell.Bits.buf)

	if cell.BitSize()%8 != 0 {
		res[len(res)-1] |= 1 << (7 - cell.BitSize()%8)
	}

	return res
}

func hashRepr(cell *Cell) ([]byte, error) {
	res := bocReprWithoutRefs(cell)
	for _, r := range cell.Refs() {
		depthRepr := make([]byte, 2)
		depthLimit := BOCSizeLimit
		v, err := getMaxDepth(r, &depthLimit)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint16(depthRepr, uint16(v))
		res = append(res, depthRepr...)
	}
	for _, r := range cell.Refs() {
		h, err := r.Hash()
		if err != nil {
			return nil, err
		}
		res = append(res, h...)
	}
	return res, nil
}

func hashCell(cell *Cell) ([]byte, error) {
	repr, err := hashRepr(cell)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(repr)
	return hash[:], nil
}

func topologicalSortImpl(cell *Cell, seen []string) ([]*Cell, error) {
	var res = make([]*Cell, 0)

	res = append(res, cell)

	hash, err := cell.HashString()
	if err != nil {
		return nil, err
	}
	if contains(seen, hash) {
		return nil, errors.New("circular references are not allowed")
	}

	for _, ref := range cell.Refs() {
		res2, err := topologicalSortImpl(ref, append(seen, hash))
		if err != nil {
			return nil, err
		}
		res = append(res, res2...)
	}

	return res, nil
}

func topologicalSort(cell *Cell) ([]*Cell, map[string]int, error) {
	res, err := topologicalSortImpl(cell, []string{})
	if err != nil {
		return nil, nil, err
	}

	indexesMap := make(map[string]int)
	for i := 0; i < len(res); i++ {
		h, err := res[i].HashString()
		if err != nil {
			return nil, nil, err
		}
		indexesMap[h] = i
	}

	return res, indexesMap, nil
}

func bocRepr(c *Cell, indexesMap map[string]int, sBytes int) ([]byte, error) {
	res := bocReprWithoutRefs(c)

	for _, ref := range c.Refs() {
		h, err := ref.HashString()
		if err != nil {
			return nil, err
		}
		idx := uint16(indexesMap[h])
		b := make([]byte, sBytes)
		binary.BigEndian.PutUint16(b[:sBytes], idx)
		res = append(res, b...)
	}

	return res, nil
}

func SerializeBoc(cell *Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	rootCell := cell
	allCells, indexesMap, err := topologicalSort(rootCell)
	if err != nil {
		return nil, err
	}

	cellsNum := uint(len(allCells))
	sBits := bits.Len(cellsNum)
	sBytes := int(math.Max(math.Ceil(float64(sBits)/8), 1))
	fullSize := uint(0)
	sizeIndex := make([]uint, 0)
	for _, cell := range allCells {
		sizeIndex = append(sizeIndex, fullSize)
		repr, err := bocRepr(cell, indexesMap, sBytes)
		if err != nil {
			return nil, err
		}
		fullSize = fullSize + uint(len(repr))
	}

	offsetBits := bits.Len(uint(fullSize))
	offsetBytes := int(math.Max(math.Ceil(float64(offsetBits)/8), 1))

	serStr := NewBitString((1023 + 32*4 + 32*3) * int(cellsNum))

	err = serStr.WriteBytes(reachBocMagicPrefix)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteBitArray([]bool{idx, hasCrc32, cacheBits})
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(flags, 2)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteInt(sBytes, 3)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteInt(offsetBytes, 8)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(cellsNum, sBytes*8)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(1, sBytes*8)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(0, sBytes*8)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(fullSize, offsetBytes*8)
	if err != nil {
		return nil, err
	}
	err = serStr.WriteUint(0, sBytes*8)
	if err != nil {
		return nil, err
	}

	if idx {
		for i := range allCells {
			err = serStr.WriteUint(sizeIndex[i], offsetBytes*8)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, cell := range allCells {
		repr, err := bocRepr(cell, indexesMap, sBytes)
		if err != nil {
			return nil, err
		}
		err = serStr.WriteBytes(repr)
		if err != nil {
			return nil, err
		}
	}

	resBytes, err := serStr.GetTopUppedArray()
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
