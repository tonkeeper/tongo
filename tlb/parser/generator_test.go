package parser

import (
	"fmt"
	"testing"
)

var SOURCE = `
addr_std$10 anycast:(Maybe Anycast) 
   workchain_id:int8 address:bits256  = MsgAddressInt;
addr_var$11 anycast:(Maybe Anycast) addr_len:(## 9) 
   workchain_id:int32 address:(bits addr_len) = MsgAddressInt;
extra_currencies$_ dict:(HashmapE 32 (VarUInteger 32)) 
                 = ExtraCurrencyCollection;
`

func TestGenerateGolangTypes(t *testing.T) {
	parsed, err := Parse(SOURCE)
	if err != nil {
		panic(err)
	}
	g := NewGenerator(nil, "LightClient")

	s, err := g.LoadTypes(parsed.Declarations)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", s)
	//_ = s
}
