package tlb

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

type SumType string
type Magic uint32
type Maybe[T any] struct {
	Null  bool
	Value T
}
type Either[M, N any] struct {
	IsRight bool
	Left    M
	Right   N
}
type EitherRef[T any] struct {
	IsRight bool
	Value   T
}
type Ref[T any] struct {
	Value T
}
type Unary uint
type VarUInteger big.Int
type Any boc.Cell

func (m *Magic) UnmarshalTLB(c *boc.Cell, tag string) error {
	a := strings.Split(tag, "$")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 2, 32)
		if err != nil {
			return err
		}
		y, err := c.ReadUint(len(a[1]))
		if x != y {
			return fmt.Errorf("magic prefix: %v not found ", tag)
		}
		return nil
	}
	a = strings.Split(tag, "#")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 16, 32)
		if err != nil {
			return err
		}
		y, err := c.ReadUint(len(a[1]) * 4)
		if x != y {
			return fmt.Errorf("magic prefix: %v not found ", tag)
		}
		return nil
	}
	return fmt.Errorf("unsupported tag: %v", tag)
}

func (m Magic) MarshalTLB(c *boc.Cell, tag string) error {
	return encodeSumTag(c, tag)
}

func (m Maybe[_]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(!m.Null)
	if err != nil {
		return err
	}
	if !m.Null {
		err = Marshal(c, m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Maybe[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	exist, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.Null = !exist
	if exist {
		err = Unmarshal(c, &m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Either[_, _]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(m.IsRight)
	if err != nil {
		return err
	}
	if m.IsRight {
		err = Marshal(c, m.Right)
		if err != nil {
			return err
		}
	} else {
		err = Marshal(c, m.Left)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Either[_, _]) UnmarshalTLB(c *boc.Cell, tag string) error {
	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.IsRight = isRight
	if isRight {
		err = Unmarshal(c, &m.Right)
		if err != nil {
			return err
		}
	} else {
		err = Unmarshal(c, &m.Left)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m EitherRef[_]) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBit(m.IsRight)
	if err != nil {
		return err
	}
	if m.IsRight {
		c, err = c.NewRef()
		if err != nil {
			return err
		}
	}
	return Marshal(c, m.Value)
}

func (m *EitherRef[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.IsRight = isRight
	if isRight {
		c, err = c.NextRef()
		if err != nil {
			return err
		}
	}
	return Unmarshal(c, &m.Value)
}

func (m Ref[_]) MarshalTLB(c *boc.Cell, tag string) error {
	r := boc.NewCell()
	err := Marshal(r, m.Value)
	if err != nil {
		return err
	}
	err = c.AddRef(r)
	return err
}

func (m *Ref[_]) UnmarshalTLB(c *boc.Cell, tag string) error {
	r, err := c.NextRef()
	if err != nil {
		return err
	}
	err = Unmarshal(r, &m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (n Unary) MarshalTLB(c *boc.Cell, tag string) error {
	return c.WriteUnary(uint(n))
}

func (n *Unary) UnmarshalTLB(c *boc.Cell, tag string) error {
	a, err := c.ReadUnary()
	*n = Unary(a)
	return err
}

func (u VarUInteger) MarshalTLB(c *boc.Cell, tag string) error {
	n, err := decodeVarUIntegerTag(tag)
	if n < 1 {
		return fmt.Errorf("len of varuint must be at least one byte")
	}
	if err != nil {
		return err
	}
	i := big.Int(u)
	b := i.Bytes()
	err = c.WriteLimUint(len(b), n-1)
	if err != nil {
		return err
	}
	err = c.WriteBytes(b)
	if err != nil {
		return err
	}
	return nil
}

func (u *VarUInteger) UnmarshalTLB(c *boc.Cell, tag string) error {
	n, err := decodeVarUIntegerTag(tag)
	if err != nil {
		return err
	}
	ln, err := c.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger(*val)
	return nil
}

func (a Any) MarshalTLB(c *boc.Cell, tag string) error {
	x := boc.Cell(a)
	y := &x
	err := c.WriteBitString(y.RawBitString())
	if err != nil {
		return err
	}
	for y.RefsAvailableForRead() > 0 {
		ref, err := y.NextRef()
		if err != nil {
			return err
		}
		err = c.AddRef(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Any) UnmarshalTLB(c *boc.Cell, tag string) error {
	x := boc.NewCell()
	err := x.WriteBitString(c.ReadRemainingBits())
	if err != nil {
		return err
	}
	for c.RefsAvailableForRead() > 0 {
		ref, err := c.NextRef()
		if err != nil {
			return err
		}
		err = x.AddRef(ref)
		if err != nil {
			return err
		}
	}
	*a = Any(*x)
	return nil
}
