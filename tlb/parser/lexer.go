package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Ident`, `[a-zA-Z_][0-9a-zA-Z0-9_]*`},
		{`HexTag`, `#([0-9a-f]+_?|_)`},
		{`BinTag`, `\$[01]*_?`},
		{"BuiltIn", `(#<=|#<|##|#)`},
		{"NUMBER", `[\d]+`},
		{`Punct`, `[][={};<>^~:)(]`},
		{"comment", `//[^\n]*`},
		{"whitespace", `\s+`},
	})
	tlbParser = participle.MustBuild[TLB](
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

func (c Combinator) String() string {
	s := c.Name
	for _, e := range c.TypeExpressions {
		s += " " + e.String()
	}
	return s
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
	ParenExpression      *ParenExpression `(@@`
	AnonymousConstructor *Anon            `| @@`
	CellRef              *CellRef         `| @@`
	BuiltIn              *string          `| @BuiltIn`
	NUMBER               *string          `| @NUMBER`
	NamedRef             *string          `| @Ident)`
}

func (t TypeExpression) String() string {
	return "Temp" //todo: implement
}

func Parse(tlb string) (*TLB, error) {
	return tlbParser.ParseString("", tlb)
}
