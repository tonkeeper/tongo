package boc

// Hasher calculates a cell's hash more efficiently than calling Hash() directly on a cell.
//
// This is a performance optimization hack and must be used only with a read-only tree of cells.
// Hasher relies on two facts:
// * some cells can be used multiple times in different places inside the tree
// * the tree is readonly and thus any hash, once calculated, won't change.
type Hasher struct {
	cache map[*Cell]*immutableCell
}

func NewHasher() *Hasher {
	return &Hasher{
		cache: map[*Cell]*immutableCell{},
	}
}

func (h *Hasher) Hash(c *Cell) ([]byte, error) {
	return hashCell(c, h.cache)
}
