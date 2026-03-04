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
	templateStr := `
type Uint{{.NameIndex}} uint{{.P}}

func (u Uint{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteUint(uint64(u), {{.NameIndex}})
}

func (u *Uint{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadUint({{.NameIndex}})
	*u = Uint{{.NameIndex}}(v)
	return err
}

func (u Uint{{.NameIndex}}) FixedSize() int {
	return {{.NameIndex}}
}

func (u Uint{{.NameIndex}}) Equal(other any) bool {
    otherInt, ok := other.(Uint{{.NameIndex}})
	if !ok {
		return false
	}
	return u == otherInt
}

func (u Uint{{.NameIndex}}) Compare(other any) (int, bool) {
    otherInt, ok := other.(Uint{{.NameIndex}})
	if !ok {
		return 0, false
	}
	if u == otherInt {
		return 0, true
	}
	if u < otherInt {
		return -1, true
	}
	return 1, true
}

{{- if lt .NameIndex 57 }}
func (u Uint{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("%d", u)), nil
}
{{- else }}
func (u Uint{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("\"%d\"", u)), nil
}
{{- end }}

func (u *Uint{{.NameIndex}}) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseUint(strings.Trim(string(p), "\""), 10, {{.NameIndex}})
    if err != nil {
		return err
    }
    *u = Uint{{.NameIndex}}(value)
    return nil
}

type Int{{.NameIndex}} int{{.P}}

func (u Int{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	return c.WriteInt(int64(u), {{.NameIndex}})
}

func (u *Int{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadInt({{.NameIndex}})
	*u = Int{{.NameIndex}}(v)
	return err
}

func (u Int{{.NameIndex}}) FixedSize() int {
	return {{.NameIndex}}
}

func (u Int{{.NameIndex}}) Equal(other any) bool {
    otherInt, ok := other.(Int{{.NameIndex}})
	if !ok {
		return false
	}
	return u == otherInt
}

func (u Int{{.NameIndex}}) Compare(other any) (int, bool) {
    otherInt, ok := other.(Int{{.NameIndex}})
	if !ok {
		return 0, false
	}
	if u == otherInt {
		return 0, true
	}
	if u < otherInt {
		return -1, true
	}
	return 1, true
}

{{- if lt .NameIndex 57 }}
func (u Int{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("%d", u)), nil
}
{{- else }}
func (u Int{{.NameIndex}}) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("\"%d\"", u)), nil
}
{{- end }}

func (u *Int{{.NameIndex}}) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseInt(strings.Trim(string(p), "\""), 10, {{.NameIndex}})
    if err != nil {
		return err
    }
    *u = Int{{.NameIndex}}(value)
    return nil
}
`
	tpl, err := template.New("smallInts").Parse(templateStr)
	if err != nil {
		panic(err)
	}
	type context struct {
		NameIndex int
		P         int
	}
	for i := 1; i <= max; i++ {
		p := nearestPow(i)
		ctx := context{NameIndex: i, P: p}
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

func GenerateConstantBigInts(sizes []int) string {
	var b bytes.Buffer

	templateStr := `
type Int{{.NameIndex}} big.Int

func (u Int{{.NameIndex}}) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	x := big.Int(u)
	return c.WriteBigInt(&x, {{.NameIndex}})
}

func (u *Int{{.NameIndex}}) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	v, err := c.ReadBigInt({{.NameIndex}})
	if err != nil {
		return err
	}
	*u = Int{{.NameIndex}}(*v)
	return err
}

func (u Int{{.NameIndex}}) FixedSize() int {
	return {{.NameIndex}}
}

func (u Int{{.NameIndex}}) Equal(other any) bool {
	otherInt, ok := other.(Int{{.NameIndex}})
	if !ok {
		return false
	}
	bigU := big.Int(u)
	otherBigInt := big.Int(otherInt)
	return bigU.Cmp(&otherBigInt) == 0
}

func (u Int{{.NameIndex}}) Compare(other any) (int, bool) {
	otherInt, ok := other.(Int{{.NameIndex}})
	if !ok {
		return 0, false
	}
	bigU := big.Int(u)
	otherBigInt := big.Int(otherInt)
	return bigU.Cmp(&otherBigInt), true
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

func GenerateConstantBigUints(sizes []int) string {
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

func (u Uint{{.NameIndex}}) Equal(other any) bool {
	otherUint, ok := other.(Uint{{.NameIndex}})
	if !ok {
		return false
	}
	bigU := big.Int(u)
	otherBigUint := big.Int(otherUint)
	return bigU.Cmp(&otherBigUint) == 0
}

func (u Uint{{.NameIndex}}) Compare(other any) (int, bool) {
	otherUint, ok := other.(Uint{{.NameIndex}})
	if !ok {
		return 0, false
	}
	bigU := big.Int(u)
	otherBigUint := big.Int(otherUint)
	return bigU.Cmp(&otherBigUint), true
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

func (u Bits%v) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("\"%%x\"", u[:])), nil
}

func (u *Bits%v) UnmarshalJSON(b []byte) error {
	bs, err := hex.DecodeString(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	if len(bs) != %v {
		return fmt.Errorf("can't parse Bits%v %%v", string(b))
	}
	copy(u[:], bs)
    return nil
}

func (u Bits%v) Equal(other any) bool {
    otherBits, ok := other.(Bits%v)
	if !ok {
		return false
	}
	return u == otherBits
}

func (u Bits%v) Compare(other any) (int, bool) {
    otherBits, ok := other.(Bits%v)
	if !ok {
		return 0, false
	}
	return bytes.Compare(u[:], otherBits[:]), true
}
	`, i, i/8, i, i, i, i, i/8, i, i, i, i, i)
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
