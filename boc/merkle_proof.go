package boc

type MerkleProver struct {
	root *immutableCell

	pruned map[*immutableCell]struct{}
}

func NewMerkleProver(root *Cell) (*MerkleProver, error) {
	immRoot, err := newImmutableCell(root, make(map[*Cell]*immutableCell))
	if err != nil {
		return nil, err
	}
	return &MerkleProver{root: immRoot, pruned: make(map[*immutableCell]struct{})}, nil
}

type Cursor struct {
	cell   *immutableCell
	prover *MerkleProver
}

func (p *MerkleProver) Cursor() *Cursor {
	return &Cursor{cell: p.root, prover: p}
}

func (p *MerkleProver) CreateProof() ([]byte, error) {
	immRoot, err := p.root.pruneCells(p.pruned)
	if err != nil {
		return nil, err
	}
	mp := NewCell()
	mp.cellType = MerkleProofCell
	if err := mp.WriteUint(3, 8); err != nil {
		return nil, err
	}
	if err := mp.WriteBytes(p.root.Hash(0)); err != nil {
		return nil, err
	}
	if err := mp.WriteUint(uint64(p.root.Depth(0)), 16); err != nil {
		return nil, err
	}
	if err := mp.AddRef(immRoot); err != nil {
		return nil, err
	}
	mp.ResetCounters()
	return SerializeBoc(mp, false, false, false, 0)
}

func (c *Cursor) Prune() {
	c.prover.pruned[c.cell] = struct{}{}
}

func (c *Cursor) Ref(ref int) *Cursor {
	return &Cursor{cell: c.cell.refs[ref], prover: c.prover}
}
