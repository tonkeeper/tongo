package parser

import (
	"bytes"
	"fmt"
	"go/format"
)

func GenerateVarUintTypes(max int) string {
	var b bytes.Buffer
	for i := 1; i <= max; i++ {
		fmt.Fprintf(&b,
			`
type VarUInteger%v big.Int

func (u VarUInteger%v) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
    b := i.Bytes()
	err := c.WriteLimUint(len(b), %v)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger%v) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(%v)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger%v(*val)
	return nil
}
`,
			i, i, i-1, i, i-1, i)
	}
	bytes, _ := format.Source(b.Bytes())
	return string(bytes)
}
