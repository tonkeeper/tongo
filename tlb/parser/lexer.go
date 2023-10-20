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
		{"Comparison", `(<=|<|>=|>|!=|=)`},
		{"NUMBER", `[\d]+`},
		{`Punct`, `[][={};<>^~:.?)(]`},
		{"comment", `//[^\n]*`},
		{"whitespace", `\s+`},
	})
	tlbParser = participle.MustBuild[TLB](
		participle.Lexer(iniLexer),
	)
)

type TLB struct {
	Declarations []CombinatorDeclaration `@@*`
}

type CombinatorDeclaration struct {
	Constructor      Constructor       `@@`
	FieldDefinitions []FieldDefinition `@@*`
	Equal            string            `"="`
	Combinator       Combinator        `@@`
	End              string            `";"`
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
	Implicit   *CurlyExpression `"{" @@ "}"`
	NamedField *NamedField      `| @@`
	Anon       *ParenExpression `| @@`
	CellRef    *CellRef         `| @@`
	TypeRef    *TypeRef         `| @@`
}

func (fd FieldDefinition) IsEmpty() bool {
	return fd.NamedField == nil && fd.CellRef == nil && fd.Implicit == nil && fd.Anon == nil && fd.TypeRef == nil
}

type NamedField struct {
	Name       string         `@(Ident|"_")`
	Sep        string         `":"`
	Expression TypeExpression `@@`
}

type TypeRef struct {
	Name string `@Ident`
}

type CurlyExpression struct {
	CompareExpr    *CompareOperatorExpr `@@`
	TypeDefinition *TypeDefinition      `| @@`
}

type TypeDefinition struct {
	Implicit   *ImplicitField  `(@@`
	Expression *TypeExpression `| @@)`
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
	Values []FieldDefinition `"[" @@* "]"`
}

type Optional struct {
	Modificator string `@Ident`
	Dot         string `"."`
	Int         int    `@NUMBER`
	QMark       string `"?"`
	Ident       string `@Ident`
}

type TypeExpression struct {
	Tilda                string           `@"~"?`
	ParenExpression      *ParenExpression `(@@`
	AnonymousConstructor *Anon            `| @@`
	CellRef              *CellRef         `| @@`
	Optional             *Optional        `| @@`
	BuiltIn              *string          `| @BuiltIn`
	Number               *int             `| @NUMBER`
	NamedRef             *string          `| @Ident)`
}

type CompareOperatorExpr struct {
	Left      *RefInner `@@`
	Operation string    `@Comparison`
	Right     *RefInner `@@`
}

type RefInner struct {
	Ident  *string `@Ident`
	Number *int    `| @NUMBER`
}

func Parse(tlb string) (*TLB, error) {
	return tlbParser.ParseString("", tlb)
}
