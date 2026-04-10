package boc

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const BOCSizeLimit = 65536
const CellBits = 1023

var ErrCellRefsOverflow = errors.New("too many refs")
var ErrNotEnoughRefs = errors.New("not enough refs")
var ErrNotSingleRoot = errors.New("should be one root cell")
var ErrDepthIsTooBig = errors.New("depth is too big")

type CellType uint8

const (
	OrdinaryCell CellType = iota
	PrunedBranchCell
	LibraryCell
	MerkleProofCell
	MerkleUpdateCell
)

type Cell struct {
	bits      BitString
	refs      [4]*Cell
	refCursor int
	cellType  CellType
	mask      levelMask
	// TODO: add capacity checking
}

func NewCell() *Cell {
	return &Cell{
		bits:     NewBitString(CellBits),
		refs:     [4]*Cell{},
		cellType: OrdinaryCell,
	}
}

func NewCellExotic(cellType CellType) *Cell {
	return &Cell{
		bits:     NewBitString(CellBits),
		refs:     [4]*Cell{},
		cellType: cellType,
	}
}

func (c *Cell) RefsSize() int {
	var count int
	for i := range c.refs {
		if c.refs[i] != nil {
			count++
		}
	}
	return count
}

func (c *Cell) Refs() []*Cell {
	res := make([]*Cell, 0, 4)
	for _, ref := range c.refs {
		if ref != nil {
			res = append(res, ref)
		}
	}
	return res
}

func (c *Cell) IsExotic() bool {
	return c.cellType != OrdinaryCell
}

func (c *Cell) IsLibrary() bool {
	return c.cellType == LibraryCell
}

func (c *Cell) GetLibraryHash() ([32]byte, error) {
	if !c.IsLibrary() {
		return [32]byte{}, errors.New("not library cell")
	}
	bytes, err := c.ReadBytes(33)
	if err != nil {
		return [32]byte{}, err
	}
	var hash [32]byte
	copy(hash[:], bytes[1:])
	return hash, nil

}

func (c *Cell) CellType() CellType {
	return c.cellType
}

func (c *Cell) BitSize() int {
	return c.bits.GetWriteCursor()
}

func (c *Cell) Hash() ([]byte, error) {
	return c.hash(map[*Cell]*immutableCell{})
}

func (c *Cell) Hash256() ([32]byte, error) {
	b, err := c.hash(map[*Cell]*immutableCell{})
	if err != nil {
		return [32]byte{}, err
	}
	var h [32]byte
	copy(h[:], b)
	return h, nil
}

func (c *Cell) HashString() (string, error) {
	h, err := c.hash(map[*Cell]*immutableCell{})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h), nil
}

func (c *Cell) hash(cache map[*Cell]*immutableCell) ([]byte, error) {
	imc, err := newImmutableCell(c, cache)
	if err != nil {
		return nil, err
	}
	return imc.Hash(maxLevel), nil
}

func (c *Cell) ToBoc() ([]byte, error) {
	bag := newBagOfCells()
	return bag.serializeBoc([]*Cell{c}, false, false, false, 0)
}

func (c *Cell) ToBocString() (string, error) {
	return c.ToBocStringCustom(false, false, false, 0)
}

func (c *Cell) ToBocBase64() (string, error) {
	return c.ToBocBase64Custom(false, false, false, 0)
}

func (c *Cell) ToBocCustom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	bag := newBagOfCells()
	return bag.serializeBoc([]*Cell{c}, idx, hasCrc32, cacheBits, 0)
}

func (c *Cell) ToBocCustomWithHasher(hasher *Hasher, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	bag := &bagOfCells{
		hasher: hasher,
	}
	return bag.serializeBoc([]*Cell{c}, idx, hasCrc32, cacheBits, 0)
}

func (c *Cell) ToBocStringCustom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) (string, error) {
	boc, err := c.ToBocCustom(idx, hasCrc32, cacheBits, flags)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(boc), nil
}

func (c *Cell) ToBocBase64Custom(idx bool, hasCrc32 bool, cacheBits bool, flags uint) (string, error) {
	boc, err := c.ToBocCustom(idx, hasCrc32, cacheBits, flags)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(boc), nil
}

func (c *Cell) NewRef() (*Cell, error) {
	n := NewCell()
	return n, c.AddRef(n)
}

func (c *Cell) AddRef(c2 *Cell) error {
	for i := range c.refs {
		if c.refs[i] == nil {
			c.refs[i] = c2
			return nil
		}
	}
	return ErrCellRefsOverflow
}

// NextRef returns the next reference cell, nil is only returned with error
func (c *Cell) NextRef() (*Cell, error) {
	if c.refCursor > 3 {
		return nil, ErrNotEnoughRefs
	}
	ref := c.refs[c.refCursor]
	if ref != nil {
		c.refCursor++
		ref.ResetCounters()
		return ref, nil
	}
	return nil, ErrNotEnoughRefs
}

// NextRefV same as NextRef but returns struct instead of pointer
func (c *Cell) NextRefV() (Cell, error) {
	if c.refCursor > 3 {
		return Cell{}, ErrNotEnoughRefs
	}
	ref := c.refs[c.refCursor]
	if ref != nil {
		c.refCursor++
		ref.ResetCounters()
		return *ref, nil
	}
	return Cell{}, ErrNotEnoughRefs
}

func (c *Cell) toStringImpl(ident string, iterationsLimit *int) string {
	var s string
	if c.IsExotic() {
		s = ident + "!x{" + c.bits.ToFiftHex() + "}\n"
	} else {
		s = ident + "x{" + c.bits.ToFiftHex() + "}\n"
	}
	if *iterationsLimit == 0 {
		return s
	}
	*iterationsLimit -= 1
	for _, ref := range c.Refs() {
		s += ref.toStringImpl(ident+" ", iterationsLimit)
	}
	return s
}

func (c Cell) ToString() string {
	iter := BOCSizeLimit
	return c.toStringImpl("", &iter)
}

func (c *Cell) ReadUint(bitLen int) (uint64, error) {
	return c.bits.ReadUint(bitLen)
}

func (c *Cell) PickUint(bitLen int) (uint64, error) {
	return c.bits.PickUint(bitLen)
}

func (c *Cell) WriteUint(val uint64, bitLen int) error {
	return c.bits.WriteUint(val, bitLen)
}

func (c *Cell) WriteInt(val int64, bitLen int) error {
	return c.bits.WriteInt(val, bitLen)
}

func (c *Cell) setTopUppedArray(arr []byte, fulfilledBytes bool) error {
	err := c.bits.SetTopUppedArray(arr, fulfilledBytes)
	c.bits.cap = 1023
	return err
}

func (c *Cell) getBuffer() []byte {
	return c.bits.Buffer()
}

func (c *Cell) Skip(n int) error {
	return c.bits.Skip(n)
}

func (c *Cell) WriteBit(val bool) error {
	return c.bits.WriteBit(val)
}

func (c *Cell) ReadBit() (bool, error) {
	return c.bits.ReadBit()
}

func (c *Cell) ReadBits(n int) (BitString, error) {
	return c.bits.ReadBits(n)
}

// ReadBitsDeep reads bits from the cell
// going deep for concatenation: c.bits + c.refs[0].bits + c.refs[0].refs[0].bits
// last used cell (tail) is returned
func (c *Cell) ReadBitsDeep(bitsToRead int) (BitString, *Cell, error) {
	result := NewBitString(bitsToRead)
	var err error
	currCell := c
	for bitsToRead > 0 {
		available := currCell.BitsAvailableForRead()
		if available == 0 {
			currCell, err = currCell.NextRef()
			if err != nil {
				return result, currCell, fmt.Errorf("reading ref for rest %d bits: %v", bitsToRead, err)
			}
			continue
		}
		if bitsToRead <= available {
			bs, err := currCell.ReadBits(bitsToRead)
			if err != nil {
				return result, currCell, err
			}
			if err = result.WriteBitString(bs); err != nil {
				return result, currCell, err
			}
			bitsToRead = 0
		} else {
			if err = result.WriteBitString(currCell.ReadRemainingBits()); err != nil {
				return result, currCell, err
			}
			bitsToRead -= available
		}
	}
	return result, currCell, nil
}

func (c *Cell) RawBitString() BitString {
	return c.bits
}

func (c *Cell) ReadStringRefTail() (string, error) {
	ref, err := c.NextRef()
	if err != nil {
		return "", err
	}
	return ref.ReadStringTail()
}

func (c *Cell) ReadStringTail() (string, error) {
	bytes, err := c.ReadBytesTail()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Cell) ReadBytesTail() ([]byte, error) {
	availableBits := c.BitsAvailableForRead()
	if availableBits%8 != 0 {
		return nil, errors.New("read string failed: not aligned")
	}
	availableRefs := c.RefsAvailableForRead()
	if availableRefs > 1 {
		return nil, errors.New("read string failed: too many refs")
	}
	headBytes, err := c.ReadBytes(availableBits / 8)
	if err != nil {
		return nil, err
	}
	var tailBytes []byte
	if availableRefs == 1 {
		ref, err := c.NextRef()
		if err != nil {
			return nil, fmt.Errorf("reading ref for tail: %v", err)
		}
		tailBytes, err = ref.ReadBytesTail()
		if err != nil {
			return nil, fmt.Errorf("reading tail: %v", err)
		}
	}
	return append(headBytes, tailBytes...), nil
}

func (c *Cell) WriteUnary(n uint) error {
	return c.bits.WriteUnary(n)
}
func (c *Cell) ReadUnary() (uint, error) {
	return c.bits.ReadUnary()
}

func (c *Cell) ReadLimUint(n int) (uint, error) {
	return c.bits.ReadLimUint(n)
}

func (c *Cell) WriteLimUint(val, n int) error {
	return c.bits.WriteLimUint(val, n)
}

func (c *Cell) WriteBitString(s BitString) error {
	return c.bits.WriteBitString(s)
}

func (c *Cell) WriteBigInt(val *big.Int, bitLen int) error {
	return c.bits.WriteBigInt(val, bitLen)
}

func (c *Cell) WriteBigUint(val *big.Int, bitLen int) error {
	return c.bits.WriteBigUint(val, bitLen)
}

func (c *Cell) ReadInt(n int) (int64, error) {
	return c.bits.ReadInt(n)
}

func (c *Cell) ReadBytes(n int) ([]byte, error) {
	return c.bits.ReadBytes(n)
}

func (c *Cell) ReadBigInt(n int) (*big.Int, error) {
	return c.bits.ReadBigInt(n)
}

func (c *Cell) ReadBigUint(n int) (*big.Int, error) {
	return c.bits.ReadBigUint(n)
}

func (c *Cell) ReadRemainingBits() BitString {
	return c.bits.ReadRemainingBits()
}

func (c *Cell) CopyRemaining() *Cell {
	if c == nil {
		return nil
	}
	rCursor := c.bits.rCursor
	c2 := NewCellWithBits(c.bits.ReadRemainingBits())
	c.bits.rCursor = rCursor
	refCursor := c.refCursor
	refsNums := c.RefsAvailableForRead()
	for i := 0; i < refsNums; i++ {
		ref, err := c.NextRef()
		if err != nil {
			// this should never happen but anyway
			panic(err)
		}
		if err := c2.AddRef(ref); err != nil {
			// this should never happen but anyway
			panic(err)
		}
	}
	c.refCursor = refCursor
	return c2
}

func (c *Cell) CopyCell() *Cell {
	if c == nil {
		return nil
	}
	c2 := NewCellWithBits(c.RawBitString())
	for _, ref := range c.Refs() {
		refCopy := ref.CopyCell()
		if err := c2.AddRef(refCopy); err != nil {
			// this should never happen but anyway
			panic(err)
		}
	}
	return c2
}

func (c *Cell) WriteBytes(b []byte) error {
	return c.bits.WriteBytes(b)
}

func (c *Cell) ResetCounters() {
	c.bits.ResetCounter()
	c.refCursor = 0
}

func (c *Cell) BitsAvailableForRead() int {
	return c.bits.BitsAvailableForRead()
}

func (c *Cell) RefsAvailableForRead() int {
	return c.RefsSize() - c.refCursor
}

func (c *Cell) Sign(key ed25519.PrivateKey) ([]byte, error) {
	hash, err := c.Hash()
	if err != nil {
		return nil, err
	}
	return ed25519.Sign(key, hash[:]), nil
}

func (c *Cell) BitsAvailableForWrite() int {
	return c.bits.BitsAvailableForWrite()
}

func (c Cell) MarshalJSON() ([]byte, error) {
	str, err := c.ToBocString()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("\"%s\"", str)), nil
}

func (c *Cell) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	cells, err := DeserializeBocHex(str)
	if err != nil {
		return err
	}
	if len(cells) != 1 {
		return errors.New("multiple cells not supported")
	}
	*c = *cells[0]
	return nil
}

func NewCellWithBits(b BitString) *Cell {
	if b.len > CellBits {
		panic("bit string not fit to Cell")
	}
	return &Cell{bits: b, cellType: OrdinaryCell}
}

func d1(cell *Cell, mask levelMask) byte {
	specBit := 0
	if cell.IsExotic() {
		specBit = 8
	}
	return byte(cell.RefsSize() + specBit + 32*int(mask))
}

func d2(cell *Cell) byte {
	return byte((cell.BitSize()+7)/8 + cell.BitSize()/8)
}

func (c *Cell) bocReprWithoutRefs(mask levelMask) []byte {
	res := make([]byte, ((c.BitSize()+7)/8)+2)
	res[0] = d1(c, mask)
	res[1] = d2(c)
	copy(res[2:], c.getBuffer())
	if c.BitSize()%8 != 0 {
		res[len(res)-1] |= 1 << (7 - c.BitSize()%8)
	}
	return res
}

// Level returns a level of this cell.
func (c *Cell) Level() int {
	return c.mask.Level()
}

func (c *Cell) GetMerkleRoot() ([32]byte, error) {
	if c.CellType() != MerkleProofCell {
		return [32]byte{}, errors.New("not merkle proof cell")
	}
	bytes, err := c.ReadBytes(33)
	if err != nil {
		return [32]byte{}, err
	}
	var hash [32]byte
	copy(hash[:], bytes[1:])
	return hash, nil

}
