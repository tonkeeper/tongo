package abi

import (
	"fmt"
	tlbParser "github.com/startfellows/tongo/tlb/parser"
	"go/format"
	"strings"
)

var defaultKnownTypes = map[string]string{
	"accountid": "tongo.AccountID",
	"cell":      "boc.Cell",
	"int8":      "int8",
	"int257":    "tongo.Int257",
	"any":       "boc.Any",
}

type Generator struct {
	knownTypes  map[string]string
	newTlbTypes []string
	typeName    string
}

func NewGenerator(knownTypes map[string]string, typeName string) *Generator {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes: knownTypes,
		typeName:   typeName,
	}
}

func (g *Generator) GetMethods(methods []GetMethod) (string, error) {
	var builder strings.Builder

	for i := range methods {
		s, err := g.GetMethod(methods[i])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
		builder.WriteRune('\n')
	}
	return builder.String(), nil
}

func (g *Generator) checkType(s string) (string, error) {
	if typeName, prs := g.knownTypes[strings.ToLower(s)]; prs {
		return typeName, nil
	}
	tlb, err := tlbParser.Parse(s)
	if err != nil {
		return "", fmt.Errorf("can't decoder %v error %w", s, err)
	}
	if len(tlb.Declarations) == 0 {
		return "", fmt.Errorf("can't parse type %v", s)
	}
	g.newTlbTypes = append(g.newTlbTypes, s)
	return tlb.Declarations[len(tlb.Declarations)-1].Combinator.Name, nil
}

func (g *Generator) GetMethod(m GetMethod) (string, error) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "func (c %v) %v(", g.typeName, m.Name)
	var args, result []string

	for _, s := range m.Input.StackValues {
		t, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		args = append(args, fmt.Sprintf("%v %v", s.Name, t))
	}
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(") (")
	for _, s := range m.Stack {
		t, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		result = append(result, t)
	}
	result = append(result, "error")
	builder.WriteString(strings.Join(result, ", "))
	builder.WriteString(") {\n}\n")
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *Generator) CollectedTypes() (string, error) {
	types, err := tlbParser.Parse(strings.Join(g.newTlbTypes, "\n"))
	if err != nil {
		return "", err
	}
	return tlbParser.GenerateGolangTypes(*types)

}
