package parser

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

func GenerateVarUintTypes(max int) string {
	var b bytes.Buffer
	templateStr := ` 
type VarUInteger{{.NameIndex}} big.Int

func (u VarUInteger{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	i := big.Int(u)
    b := i.Bytes()
	err := c.WriteLimUint(len(b), {{.BitsLimit}})
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	ln, err := c.ReadLimUint({{.BitsLimit}})
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger{{.NameIndex}}(*val)
	return nil
}

func (u VarUInteger{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    i := big.Int(u)
    return []byte(fmt.Sprintf("\"%s\"", i.String())), nil
}

func (u *VarUInteger{{.NameIndex}}) UnmarshalJSON(p []byte) error {
    var z big.Int
    _, ok := z.SetString(strings.Trim(string(p), "\""), 10)
    if !ok {
        return fmt.Errorf("invalid integer: %s", p)
    }
    *u = VarUInteger{{.NameIndex}}(z)
    return nil
}
`
	tpl, err := template.New("varUInteger").Parse(templateStr)
	if err != nil {
		panic(err)
	}
	type context struct {
		NameIndex int
		BitsLimit int
	}
	for i := 1; i <= max; i++ {
		ctx := context{NameIndex: i, BitsLimit: i - 1}
		if err := tpl.Execute(&b, ctx); err != nil {
			panic(err)
		}
	}
	bytes, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func GenerateConstantInts(max int) string {
	var b bytes.Buffer
	for i := 1; i <= max; i++ {
		p := nearestPow(i)
		fmt.Fprintf(&b, `
type Uint%v uint%v

func (u Uint%v) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteUint(uint64(u), %v)
}

func (u *Uint%v) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadUint(%v)
	*u = Uint%v(v)
	return err
}

func (u Uint%v) FixedSize() int {
	return %v
}

type Int%v int%v

func (u Int%v) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteInt(int64(u), %v)
}

func (u *Int%v) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadInt(%v)
	*u = Int%v(v)
	return err
}

func (u Int%v) FixedSize() int {
	return %v
}
`, i, p, i, i, i, i, i, i, i, i, p, i, i, i, i, i, i, i)
	}
	bytes, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func GenerateConstantBigInts(sizes []int) string {
	var b bytes.Buffer

	templateStr := `
type Uint{{.NameIndex}} big.Int

func (u Uint{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	x := big.Int(u)
	return c.WriteBigUint(&x, {{.NameIndex}})
}

func (u *Uint{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadBigUint({{.NameIndex}})
	if err != nil {
		return err
	}
	*u = Uint{{.NameIndex}}(*v)
	return err
}

func (u Uint{{.NameIndex}}) FixedSize() int {
	return {{.NameIndex}}
}

func (u Uint{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    i := big.Int(u)
    return []byte(fmt.Sprintf("\"%s\"", i.String())), nil
}

func (u *Uint{{.NameIndex}}) UnmarshalJSON(p []byte) error {
    var z big.Int
    _, ok := z.SetString(strings.Trim(string(p), "\""), 10)
    if !ok {
        return fmt.Errorf("invalid integer: %s", p)
    }
    *u = Uint{{.NameIndex}}(z)
    return nil
}

type Int{{.NameIndex}} big.Int

func (u Int{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	x := big.Int(u)
	return c.WriteBigInt(&x, {{.NameIndex}})
}

func (u *Int{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadBigInt({{.NameIndex}})
	*u = Int{{.NameIndex}}(*v)
	return err
}

func (u Int{{.NameIndex}}) FixedSize() int {
	return {{.NameIndex}}
}

func (u Int{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    i := big.Int(u)
    return []byte(fmt.Sprintf("\"%s\"", i.String())), nil
}

func (u *Int{{.NameIndex}}) UnmarshalJSON(p []byte) error {
    var z big.Int
    _, ok := z.SetString(strings.Trim(string(p), "\""), 10)
    if !ok {
        return fmt.Errorf("invalid integer: %s", p)
    }
    *u = Int{{.NameIndex}}(z)
    return nil
}
`
	tpl, err := template.New("bigInts").Parse(templateStr)
	if err != nil {
		panic(err)
	}
	type context struct {
		NameIndex int
	}
	for _, i := range sizes {
		ctx := context{NameIndex: i}
		if err := tpl.Execute(&b, ctx); err != nil {
			panic(err)
		}
	}
	bytes, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func GenerateBitsTypes(sizes []int) string {
	var b bytes.Buffer
	for _, i := range sizes {
		if i%8 == 0 {
			fmt.Fprintf(&b, `
type Bits%v [%v]byte

func (u Bits%v) FixedSize() int {
	return %v
}
	`, i, i/8, i, i)
		} else {
			fmt.Fprintf(&b, `
type Bits%v boc.BitString
	
func (u Bits%v) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteBitString(boc.BitString(u))
}
	
func (u *Bits%v) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadBits(%v)
	*u = Bits%v(v)
	return err
}

func (u Bits%v) FixedSize() int {
	return %v
}
	`, i, i, i, i, i, i, i)
		}
	}
	bytes, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func nearestPow(i int) int {
	switch {
	case i <= 8:
		return 8
	case i <= 16:
		return 16
	case i <= 32:
		return 32
	case i <= 64:
		return 64
	default:
		panic(i)
	}
}
