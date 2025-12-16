package tlb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

//go:generate go run generator.go

type SumType string

type Magic uint32

type Void struct{}

type NullLiteral = Void

type Maybe[T any] struct {
	Exists bool
	Value  T
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

type Any boc.Cell

func (m Maybe[T]) Pointer() *T {
	if m.Exists {
		return &m.Value
	}
	return nil
}

func (m *Magic) ValidateTag(c *boc.Cell, tag string) error {
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
		*m = Magic(x)
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
		*m = Magic(x)
		return nil
	}
	return fmt.Errorf("unsupported tag: %v", tag)
}

func (m Magic) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", uint32(m)))
}

func (m *Magic) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	magic, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return err
	}
	*m = Magic(magic)
	return nil
}

func (m Magic) EncodeTag(c *boc.Cell, tag string) error {
	return encodeSumTag(c, tag)
}

func (v Void) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return nil
}

func (v *Void) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	return nil
}

func (v Void) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (v *Void) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	return fmt.Errorf("not a void value %v", string(b))
}

func (m Maybe[_]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	err := c.WriteBit(m.Exists)
	if err != nil {
		return err
	}
	if m.Exists {
		err = encoder.Marshal(c, m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Maybe[_]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	exist, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.Exists = exist
	if exist {
		err = decoder.Unmarshal(c, &m.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if m.Exists {
		return json.Marshal(m.Value)
	}
	return []byte("null"), nil
}

func (m *Maybe[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		m.Exists = false
		var defaultValue T
		m.Value = defaultValue
		return nil
	}
	m.Exists = true
	return json.Unmarshal(b, &m.Value)
}

func (m Either[_, _]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
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

func (m *Either[_, _]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	m.IsRight = isRight
	if isRight {
		err = decoder.Unmarshal(c, &m.Right)
		if err != nil {
			return err
		}
	} else {
		err = decoder.Unmarshal(c, &m.Left)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m EitherRef[_]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
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

func (m *EitherRef[_]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
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
	return decoder.Unmarshal(c, &m.Value)
}

func (m Ref[_]) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	r := boc.NewCell()
	err := Marshal(r, m.Value)
	if err != nil {
		return err
	}
	err = c.AddRef(r)
	return err
}

func (m *Ref[T]) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	r, err := c.NextRef()
	if err != nil {
		return err
	}
	if r.CellType() == boc.PrunedBranchCell {
		var value T
		m.Value = value
		return nil
	}
	err = decoder.Unmarshal(r, &m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (n Unary) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteUnary(uint(n))
}

func (n *Unary) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	a, err := c.ReadUnary()
	*n = Unary(a)
	return err
}

func (a Any) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
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

func (a *Any) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	x := c.CopyRemaining()
	*a = Any(*x)
	return nil
}

func (a Any) MarshalJSON() ([]byte, error) {
	return boc.Cell(a).MarshalJSON()
}

func (a *Any) UnmarshalJSON(b []byte) error {
	return (*boc.Cell)(a).UnmarshalJSON(b)
}
