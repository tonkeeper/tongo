package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2/lexer"
)

var SOURCE = `
int_msg_info$0 ihr_disabled:Bool bounce:Bool bounced:Bool
  src:MsgAddressInt dest:MsgAddressInt 
  value:CurrencyCollection ihr_fee:Grams fwd_fee:Grams
  created_lt:uint64 created_at:uint32 = CommonMsgInfo;
ext_in_msg_info$10 src:MsgAddressExt dest:MsgAddressInt 
  import_fee:Grams = CommonMsgInfo;
ext_out_msg_info$11 src:MsgAddressInt dest:MsgAddressExt
  created_lt:uint64 created_at:uint32 = CommonMsgInfo;

int_msg_info$0 ihr_disabled:Bool bounce:Bool bounced:Bool
  src:MsgAddress dest:MsgAddressInt 
  value:CurrencyCollection ihr_fee:Grams fwd_fee:Grams
  created_lt:uint64 created_at:uint32 = CommonMsgInfoRelaxed;
ext_out_msg_info$11 src:MsgAddress dest:MsgAddressExt
  created_lt:uint64 created_at:uint32 = CommonMsgInfoRelaxed;

tick_tock$_ tick:Bool tock:Bool = TickTock;

_ split_depth:(Maybe (## 5)) special:(Maybe TickTock)
  code:(Maybe ^Cell) data:(Maybe ^Cell)
  library:(HashmapE 256 SimpleLib) = StateInit;
`

var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Ident`, `[a-zA-Z_][0-9a-zA-Z0-9_]*`},
		{`HexTag`, `#([0-9a-f]+_?|_)`},
		{`BinTag`, `\$[01]*_?`},
		{"BuiltIn", `(#<=|#<|##|#)`},
		{"NUMBER", `[\d]+`},
		{`Punct`, `[][={};<>^:)(]`},
		{"comment", `//[^\n]*`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[TLB](
		participle.Lexer(iniLexer),
	)
)

type TLB struct {
	Declarations []*CombinatorDeclaration `@@*`
}
type CombinatorDeclaration struct {
	Constructor      Constructor        `@@`
	FieldDefinitions []*FieldDefinition `@@*`
	Equal            string             `"="`
	Combinator       Combinator         `@@`
	End              string             `";"`
}
type Constructor struct {
	Name   string `@Ident`
	Prefix string `@(HexTag|BinTag)?`
}

type Combinator struct {
	Name            string           `@Ident`
	TypeExpressions []TypeExpression `@@*`
}

type FieldDefinition struct {
	Implicit   *ImplicitDefinition `@@`
	NamedField *NamedField         `| @@`
	CellRef    *CellRef            `| @@`
}

type NamedField struct {
	Name       string         `@(Ident|"_")`
	Sep        string         `":"`
	Expression TypeExpression `@@`
}

type ImplicitDefinition struct {
	Start      string          `"{"`
	Implicit   *ImplicitField  `(@@`
	Expression *TypeExpression `| @@)`
	End        string          `"}"`
}
type ImplicitField struct {
	Name string `@Ident`
	Sep  string `":"`
	Type string `@("#"|"Type")`
}

type ParenExpression struct {
	Name      TypeExpression   `"(" @@ `
	Parameter []TypeExpression `@@* ")"`
}
type CellRef struct {
	TypeExpression TypeExpression `"^" @@`
}

type Anon struct {
	Values []FieldDefinition `"[" @@ "]"`
}

type TypeExpression struct {
	Tilda                string           `@"~"?`
	ParenExpression      *ParenExpression `@@`
	AnonymousConstructor *Anon            `| @@`
	CellRef              *CellRef         `| @@`
	BuiltIn              *string          `| @BuiltIn`
	NUMBER               *string          `| @NUMBER`
	NamedRef             *string          `| @Ident`
}

func main() {
	ini, err := parser.ParseString("", SOURCE)
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
