package boc

import "math"

type HashmapE struct {
	ConstructorName string
	Root            *Hashmap
	N               uint32
	X               any
}

type Hashmap struct {
	ConstructorName string
	Label           HmLabel
	Node            HashmapNode
	N               uint32
	X               any
}

// HmLabel
// hml_short (len, s)
// hml_long (n, s)
// hml_same (v, n)
type HmLabel struct {
	ConstructorName string
	Len             uint64
	S               *BitString
	N, M            uint32
	V               *bool
}

type HashmapNode struct {
	ConstructorName string
	NPlus           uint32
	Value           *TLBType
	Left            *Hashmap
	Right           *Hashmap
}

// ReadHashmapE
// TL-B:
// hme_empty$0 {n:#} {X:Type} = HashmapE n X;
// hme_root$1 {n:#} {X:Type} root:^(Hashmap n X) = HashmapE n X;
func ReadHashmapE(c *CellReader, n uint32, isRef bool, x TLBType) (HashmapE, error) {
	root, err := c.ReadBit()
	if err != nil {
		return HashmapE{}, err
	}
	if root {
		rootCell, err := c.GetRef()
		if err != nil {
			return HashmapE{}, err
		}
		root, err := ReadHashmap(&rootCell, n, isRef, x)
		if err != nil {
			return HashmapE{}, err
		}
		return HashmapE{ConstructorName: "hme_root", Root: &root, N: n, X: x}, nil
	}
	return HashmapE{ConstructorName: "hme_empty", N: n, X: x}, nil
}

// ReadHashmap
// TL-B:
// hm_edge#_ {n:#} {X:Type} {l:#} {m:#} label:(HmLabel ~l n) {n = (~m) + l} node:(HashmapNode m X) = Hashmap n X;
func ReadHashmap(c *CellReader, n uint32, isRef bool, x TLBType) (Hashmap, error) {
	label, err := ReadHmLabel(c, n)
	if err != nil {
		return Hashmap{}, err
	}
	m := n - label.N
	node, err := ReadHashmapNode(c, m, isRef, x)
	if err != nil {
		return Hashmap{}, err
	}
	return Hashmap{ConstructorName: "hm_edge", Label: label, Node: node, N: n, X: x}, nil
}

// ReadHmLabel
// TL-B:
// hml_short$0 {m:#} {n:#} len:(Unary ~n) {n <= m} s:(n * Bit) = HmLabel ~n m;
// hml_long$10 {m:#} n:(#<= m) s:(n * Bit) = HmLabel ~n m;
// hml_same$11 {m:#} v:Bit n:(#<= m) = HmLabel ~n m;
func ReadHmLabel(c *CellReader, m uint32) (HmLabel, error) {
	notShort, err := c.ReadBit()
	if err != nil {
		return HmLabel{}, err
	}
	if !notShort {
		same, err := c.ReadBit()
		if err != nil {
			return HmLabel{}, err
		}
		if same {
			// decode hml_same
			v, err := c.ReadBit()
			if err != nil {
				return HmLabel{}, err
			}
			nLen := int(math.Ceil(math.Log2(float64(m + 1))))
			n, err := c.ReadUint(nLen)
			if err != nil {
				return HmLabel{}, err
			}
			return HmLabel{ConstructorName: "hml_same", V: &v, N: uint32(n), M: m}, nil
		}
		// decode hml_long
		nLen := int(math.Ceil(math.Log2(float64(m + 1))))
		n, err := c.ReadUint(nLen)
		if err != nil {
			return HmLabel{}, err
		}
		bits, err := c.ReadBits(int(n))
		if err != nil {
			return HmLabel{}, err
		}
		return HmLabel{ConstructorName: "hml_long", S: &bits, N: uint32(n), M: m}, nil
	}
	// decode hml_short
	ln, err := c.ReadUnary()
	if err != nil {
		return HmLabel{}, err
	}
	bits, err := c.ReadBits(int(ln))
	if err != nil {
		return HmLabel{}, err
	}
	return HmLabel{ConstructorName: "hml_short", Len: ln, S: &bits, N: uint32(ln), M: m}, nil
}

// ReadHashmapNode
// TL-B:
// hmn_leaf#_ {X:Type} value:X = HashmapNode 0 X;
// hmn_fork#_ {n:#} {X:Type} left:^(Hashmap n X) right:^(Hashmap n X) = HashmapNode (n + 1) X;
func ReadHashmapNode(c *CellReader, nPlus uint32, isRef bool, x TLBType) (HashmapNode, error) {
	if nPlus == 0 {
		if isRef {
			ref, err := c.GetRef()
			if err != nil {
				return HashmapNode{}, err
			}
			err = x.UnmarshalTLB(&ref)
			if err != nil {
				return HashmapNode{}, err
			}
		} else {
			err := x.UnmarshalTLB(c)
			if err != nil {
				return HashmapNode{}, err
			}
		}
		return HashmapNode{ConstructorName: "hmn_leaf", Value: &x, NPlus: nPlus}, nil
	}
	n := nPlus - 1
	leftCell, err := c.GetRef()
	if err != nil {
		return HashmapNode{}, err
	}
	left, err := ReadHashmap(&leftCell, n, isRef, x)
	if err != nil {
		return HashmapNode{}, err
	}
	rightCell, err := c.GetRef()
	if err != nil {
		return HashmapNode{}, err
	}
	right, err := ReadHashmap(&rightCell, n, isRef, x)
	if err != nil {
		return HashmapNode{}, err
	}
	return HashmapNode{ConstructorName: "hmn_fork", NPlus: nPlus, Left: &left, Right: &right}, nil
}
