package boc

import "encoding/hex"

// Hasher calculates a cell's hash more efficiently than calling Hash() directly on a cell.
//
// This is a performance optimization hack and must be used only with a read-only tree of cells.
// Hasher relies on two facts:
// * some cells can be used multiple times in different places inside the tree
// * the tree is readonly and thus any hash, once calculated, won't change.
type Hasher struct {
	cache    map[*Cell]*immutableCell
	cacheHex map[*Cell]string
}

func NewHasher() *Hasher {
	return &Hasher{
		cache:    map[*Cell]*immutableCell{},
		cacheHex: map[*Cell]string{},
	}
}

func (h *Hasher) Hash(c *Cell) ([]byte, error) {
	return c.hash(h.cache)
}

func (h *Hasher) HashString(c *Cell) (string, error) {
	if s, ok := h.cacheHex[c]; ok {
		return s, nil
	}
	hash, err := h.Hash(c)
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(hash)
	h.cacheHex[c] = s
	return s, nil
}
