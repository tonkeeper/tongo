package parser

import (
	"fmt"
	"github.com/startfellows/tongo/utils"
	"go/format"
	"strings"
)

func GenerateGolangTypes(t TLB) (string, error) {
	sumTypes := make(map[string][]CombinatorDeclaration)

	for i, c := range t.Declarations {
		if c == nil {
			return "", fmt.Errorf("declaration %v is nil", i)
		}
		if len(c.Combinator.TypeExpressions) > 0 {
			return "", fmt.Errorf("combinators with paramaters '%v' are not supported", c.Combinator.String())
		}
		sumTypes[c.Combinator.Name] = append(sumTypes[c.Combinator.Name], *c)
	}
	s := ""
	for _, v := range sumTypes {
		t, err := generateGolangType(v)
		if err != nil {
			return "", err
		}
		s += "\n" + t
	}

	b, err := format.Source([]byte(s))
	if err != nil {
		return s, err
	}
	return string(b), err
}

func generateGolangStruct(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("struct{")
	if len(declaration.FieldDefinitions) > 0 {
		builder.WriteRune('\n')
	}
	for i, field := range declaration.FieldDefinitions {
		if field == nil {
			return "", fmt.Errorf("nil field %v in %v", i, declaration.Constructor.Name)
		}
		if field.NamedField == nil && field.CellRef == nil && field.Implicit == nil {
			return "", fmt.Errorf("all types are nil in field %v in %v", i, declaration.Constructor.Name)
		}
		if field.Implicit != nil {
			continue
		}
		var name string
		var e TypeExpression
		if field.CellRef != nil {
			e = field.CellRef.TypeExpression
		} else {
			name = field.NamedField.Name
			e = field.NamedField.Expression
		}
		if name == "" || name == "_" {
			name = fmt.Sprintf("Field%v", i)
		}
		builder.WriteString(utils.ToCamelCase(name))
		builder.WriteRune('\t')
		t, err := e.ToGolangType()
		if err != nil {
			return "", err
		}
		builder.WriteString(t.name)
		builder.WriteRune('\n')
	}
	builder.WriteRune('}')
	return builder.String(), nil
}

func generateGolangSimpleType(declaration CombinatorDeclaration) (string, error) {
	s, err := generateGolangStruct(declaration)
	return fmt.Sprintf("type %v %v", declaration.Combinator.Name, s), err
}

func generateGolangSumType(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("type " + declarations[0].Combinator.Name + " struct{\ntlb.SumType\n")
	for _, d := range declarations {
		s, err := generateGolangStruct(d)
		if err != nil {
			return "", err
		}
		builder.WriteString(utils.ToCamelCase(d.Constructor.Name))
		builder.WriteRune(' ')
		builder.WriteString(s)
		builder.WriteString(fmt.Sprintf(" `tlbSumType:\"%v\"`", d.Constructor.Prefix))
		builder.WriteRune('\n')
	}
	builder.WriteRune('}')
	return builder.String(), nil

}

func generateGolangType(declarations []CombinatorDeclaration) (string, error) {
	if len(declarations) == 1 {
		return generateGolangSimpleType(declarations[0])
	} else {
		return generateGolangSumType(declarations)
	}
}

type golangType struct {
	name string
	tag  string
}

func (t *TypeExpression) ToGolangType() (golangType, error) {
	return golangType{
		name: "Temp",
		tag:  "",
	}, nil
}
