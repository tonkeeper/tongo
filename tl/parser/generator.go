package parser

import (
	"fmt"
	"github.com/startfellows/tongo/utils"
	"go/format"
	"strings"
)

func GenerateGolangTypes(t TL) (string, error) {
	sumTypes := make(map[string][]CombinatorDeclaration)

	for i, c := range t.Declarations {
		if c == nil {
			return "", fmt.Errorf("declaration %v is nil", i)
		}
		sumTypes[c.Combinator] = append(sumTypes[c.Combinator], *c)
	}

	s := ""
	for _, v := range sumTypes {
		t, err := generateGolangType(v)
		if err != nil {
			return "", err
		}
		s += "\n" + t + "\n"
	}

	b, err := format.Source([]byte(s))
	if err != nil {
		return s, err
	}
	return string(b), err
}

func generateGolangType(declarations []CombinatorDeclaration) (string, error) {
	if len(declarations) == 1 {
		return generateGolangSimpleType(declarations[0])
	} else {
		return generateGolangSumType(declarations)
	}
}

func generateGolangSimpleType(declaration CombinatorDeclaration) (string, error) {
	s, err := generateGolangStruct(declaration)
	return fmt.Sprintf("type %v %v", utils.ToCamelCase(declaration.Combinator), s), err
}

func generateGolangSumType(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("type " + utils.ToCamelCase(declarations[0].Combinator) + " struct{\ntl.SumType\n")
	for _, d := range declarations {
		s, err := generateGolangStruct(d)
		if err != nil {
			return "", err
		}
		builder.WriteString(utils.ToCamelCase(d.Constructor))
		builder.WriteRune(' ')
		builder.WriteString(s)
		builder.WriteRune('\n')
	}
	builder.WriteRune('}')
	return builder.String(), nil
}

func generateGolangStruct(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("struct{")
	if len(declaration.FieldDefinitions) > 0 {
		builder.WriteRune('\n')
	}
	for i, field := range declaration.FieldDefinitions {
		if field == nil {
			return "", fmt.Errorf("nil field %v in %v", i, declaration.Constructor)
		}

		var name string
		var e TypeExpression
		name = field.Name
		e = field.Expression

		if name == "" || name == "_" {
			name = fmt.Sprintf("Field%v", i)
		}

		optional := false
		if field.Modificator.Name == "mode" { // mode.0?field
			optional = true
		}

		builder.WriteString(utils.ToCamelCase(name))
		builder.WriteRune('\t')
		t, err := toGolangType(e, optional)
		if err != nil {
			return "", err
		}
		builder.WriteString(t.String())
		builder.WriteRune('\n')
	}
	builder.WriteRune('}')
	return builder.String(), nil
}

type golangType struct {
	name     string
	optional bool
}

func (t golangType) String() string {
	if !t.optional {
		return t.name
	}
	if strings.HasPrefix(t.name, "[]") {
		return t.name
	}
	return "*" + t.name
}

func mapToGoType(tlType string) string {
	goType, ok := typesMapping[tlType]
	if ok {
		return goType
	}
	return utils.ToCamelCase(tlType)
}

func toGolangType(t TypeExpression, optional bool) (golangType, error) {
	if t.BuiltIn != nil {
		return golangType{
			name:     mapToGoType(*t.BuiltIn),
			optional: optional,
		}, nil
	}
	if t.NamedRef != nil {
		return golangType{
			name:     mapToGoType(*t.NamedRef),
			optional: optional,
		}, nil
	}

	if t.Vector != nil {
		if len(t.Vector.Parameter) != 1 {
			return golangType{}, fmt.Errorf("vector must contains only one parameter")
		}
		gt, err := toGolangType(t.Vector.Parameter[0], false) // can not be pointer type under vector
		if err != nil {
			return golangType{}, err
		}
		return golangType{
			name:     "[]" + mapToGoType(gt.name),
			optional: optional,
		}, nil
	}
	return golangType{}, fmt.Errorf("invalid type expression")
}
