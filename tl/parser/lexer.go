package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`FunctionsSeparator`, `---functions---`},
		{`Ident`, `[a-zA-Z][0-9a-zA-Z0-9\_]*(\.[a-zA-Z][0-9a-zA-Z0-9\_]*)*`},
		{`HexTag`, `#([0-9a-f]+_?|_)`},
		{"BuiltIn", `(#<=|#<|##|#)`},
		{"NUMBER", `[\d]+`},
		{`Punct`, `[][={}.?;\-<>^:)(]`},
		{"comment", `//[^\n]*`},
		{"whitespace", `\s+`},
	})
	tlParser = participle.MustBuild[TL](
		participle.Lexer(iniLexer),
	)
	typesMapping = map[string]string{
		"#":      "int",
		"int":    "int",
		"int256": "tl.Int256",
		"long":   "int64",
		"bytes":  "[]byte",
		"Bool":   "bool",
	}
)

type TL struct {
	Declarations []*CombinatorDeclaration `@@*`
	Separator    string                   `FunctionsSeparator`
	Functions    []*CombinatorDeclaration "@@*"
}

type CombinatorDeclaration struct {
	Constructor      string        `@Ident`
	Tag              string        `@HexTag?`
	FieldDefinitions []*NamedField `@@*`
	Equal            string        `"="`
	Combinator       string        `@Ident`
	End              string        `";"`
}

type NamedField struct {
	Name        string         `@Ident`
	Sep         string         `":"`
	Modificator Modificator    `@@?`
	Expression  TypeExpression `@@`
}

type Modificator struct {
	Name string `@Ident "."`
	Bit  string `@NUMBER "?"`
}

type TypeExpression struct {
	BuiltIn  *string          `@BuiltIn`
	NamedRef *string          `| @Ident`
	Vector   *ParenExpression `| @@`
}
type ParenExpression struct {
	Name      string           `"(" @Ident `
	Parameter []TypeExpression `@@* ")"`
}

func (t TypeExpression) String() string {
	if t.BuiltIn != nil {
		return *t.BuiltIn
	}
	if t.NamedRef != nil {
		return *t.NamedRef
	}
	return "<nil>"
}

func Parse(tl string) (*TL, error) {
	a, err := tlParser.ParseString("", tl)
	return a, err
}
