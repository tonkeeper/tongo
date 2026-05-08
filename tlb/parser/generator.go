package parser

import (
	"fmt"
	"go/format"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/utils"
	"golang.org/x/exp/maps"
)

func pkgRef(pkg, name string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

type DefaultType struct {
	Name          string
	IsPointerType bool
}

type TlbType struct {
	Name       string
	Definition string
}

type Generator struct {
	knownTypes  map[string]DefaultType
	newTlbTypes map[string]TlbType
	tlbPkg      string // package prefix for tlb types, e.g. "tlb"; empty when generating for the tlb package itself
}

var (
	defaultKnownTypes = map[string]DefaultType{
		"#":                    {"uint32", false},
		"int8":                 {"int8", false},
		"int16":                {"int16", false},
		"int32":                {"int32", false},
		"int64":                {"int64", false},
		"uint8":                {"uint8", false},
		"uint16":               {"uint16", false},
		"uint32":               {"uint32", false},
		"uint64":               {"uint64", false},
		"Bool":                 {"bool", false},
		"True":                 {"struct{}", false},
		"Unit":                 {"struct{}", false},
		"Cell":                 {"tlb.Any", false},
		"MsgAddress":           {"tlb.MsgAddress", false},
		"MsgAddressInt":        {"tlb.MsgAddress", false}, //todo: replace with MsgAddressInt after adding to tlb package
		"MsgAddressExt":        {"tlb.MsgAddress", false}, //todo: replace with MsgAddressExt after adding to tlb package
		"AddressWithWorkchain": {"tlb.AddressWithWorkchain", false},
		"Coins":                {"tlb.Grams", false},
		"Grams":                {"tlb.Grams", false},
		"Text":                 {"tlb.Text", false},
		"Bytes":                {"tlb.Bytes", false},
		"FixedLengthText":      {"tlb.FixedLengthText", false},
		"SnakeData":            {"tlb.SnakeData", false},
		"ChunkedData":          {"tlb.ChunkedData", false},
		"DNSRecord":            {"tlb.DNSRecord", false},
		"DNS_RecordSet":        {"tlb.DNSRecordSet", false},
		"CurrencyCollection":   {"tlb.CurrencyCollection", false},
	}
)

type Option func(*Generator)

func WithDefaultTypes(types map[string]DefaultType, replace bool) Option {
	return func(g *Generator) {
		if replace {
			g.knownTypes = make(map[string]DefaultType)
		}
		for k, v := range types {
			g.knownTypes[k] = v
		}
	}
}

func WithTlbPackage(pkg string) Option {
	return func(g *Generator) {
		if pkg == "" {
			for k, v := range g.knownTypes {
				v.Name = strings.TrimPrefix(v.Name, "tlb.")
				g.knownTypes[k] = v
			}
		}
		g.tlbPkg = pkg
	}
}

func NewGenerator(options ...Option) *Generator {

	g := &Generator{
		knownTypes:  maps.Clone(defaultKnownTypes),
		newTlbTypes: make(map[string]TlbType),
		tlbPkg:      "tlb",
	}
	for _, o := range options {
		o(g)
	}
	return g
}

func (g *Generator) GetTlbTypes() []TlbType {
	var res []TlbType

	for _, k := range utils.GetOrderedKeys(g.newTlbTypes) {
		res = append(res, g.newTlbTypes[k])
	}
	return res
}

func (g *Generator) GenerateGolangTypes(declarations []CombinatorDeclaration, typePrefix string, skipMagic bool) (string, error) {
	dec := make([][]CombinatorDeclaration, 0)
	for _, c := range declarations {
		if len(c.Combinator.TypeExpressions) == 1 && c.Combinator.TypeExpressions[0].Number != nil {
			c = CombinatorDeclaration{
				Constructor:      c.Constructor,
				FieldDefinitions: c.FieldDefinitions,
				Equal:            c.Equal,
				Combinator: Combinator{
					Name: fmt.Sprintf("%v%v", c.Combinator.Name, *c.Combinator.TypeExpressions[0].Number),
				},
				End: c.End,
			}
		} else if len(c.Combinator.TypeExpressions) > 0 {
			return "", fmt.Errorf("combinators with parameters '%v' are not supported", c.Combinator.Name)
		}
		f := false
		for i, c1 := range dec {
			if c1[0].Combinator.Name == c.Combinator.Name {
				dec[i] = append(dec[i], c)
				f = true
				break
			}
		}
		if !f {
			dec = append(dec, []CombinatorDeclaration{c})
		}
	}
	s := ""

	for _, v := range dec {
		name := v[0].Combinator.Name
		if typePrefix != "" {
			name = utils.ToCamelCase(typePrefix)
		}
		t, err := g.generateGolangType(v, name, skipMagic)
		if err != nil {
			return "", err
		}

		b, err := format.Source([]byte(t))
		if err != nil {
			return "", err
		}

		g.newTlbTypes[name] = TlbType{
			Name:       name,
			Definition: string(b),
		}
		s += "\n\n" + t
	}

	return s, nil
}

func (g *Generator) generateGolangStruct(declaration CombinatorDeclaration, skipMagic bool, enclosingType string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("struct{")
	if len(declaration.FieldDefinitions) > 0 {
		builder.WriteRune('\n')
	}

	if !skipMagic && declaration.Constructor.Prefix != "" && declaration.Constructor.Prefix != "#_" && declaration.Constructor.Prefix != "$_" {
		builder.WriteString(fmt.Sprintf("Magic %s `tlb:\"%v\"`\n", pkgRef(g.tlbPkg, "Magic"), declaration.Constructor.Prefix))
	}
	s, err := g.fieldDefinitionsToStruct(declaration.FieldDefinitions, enclosingType)
	if err != nil {
		return "", err
	}
	builder.WriteString(s)

	builder.WriteRune('}')
	return builder.String(), nil
}

// fieldDefinitionsToStruct generates the Go struct body for a list of TLB field definitions.
// enclosingType is the name of the sum type currently being generated; fields whose resolved
// type matches enclosingType are automatically wrapped in a pointer to break recursive cycles.
func (g *Generator) fieldDefinitionsToStruct(definitions []FieldDefinition, enclosingType string) (string, error) {
	var builder strings.Builder
	for i, field := range definitions {
		if field.IsEmpty() {
			return "", fmt.Errorf("all types are nil in field %v ", i)
		}
		if field.Implicit != nil {
			continue
		}
		var name string
		var e TypeExpression
		if field.CellRef != nil {
			e = field.CellRef.TypeExpression
		} else if field.NamedField != nil {
			name = field.NamedField.Name
			e = field.NamedField.Expression
		} else if field.TypeRef != nil {
			builder.WriteString(fmt.Sprintf("%s %s\n", field.TypeRef.Name, field.TypeRef.Name))
			continue
		}
		if field.Anon != nil {
			t, err := field.Anon.toGolangType(g)
			if err != nil {
				return "", err
			}
			value := fmt.Sprintf("Value %s\n", t.String())
			builder.WriteString(value)
			continue
		}
		if name == "" || name == "_" {
			name = fmt.Sprintf("Field%v", i)
		}
		t, err := e.toGolangType(g)
		if err != nil {
			return "", err
		}
		typeName := t.String()
		if enclosingType != "" && typeName == enclosingType {
			t = golangType{name: "*" + typeName, tag: t.tag}
		}
		builder.WriteString(utils.ToCamelCase(name))
		builder.WriteRune('\t')
		builder.WriteString(t.String())
		if len(t.tag) > 0 {
			builder.WriteString(fmt.Sprintf("`tlb:\"%s\"`", t.tag))
		} else if field.CellRef != nil {
			builder.WriteString("`tlb:\"^\"`")
		}
		builder.WriteRune('\n')
	}
	return builder.String(), nil
}

func (g *Generator) generateGolangSimpleType(declaration CombinatorDeclaration, typeName string, skipMagic bool) (string, error) {
	s, err := g.generateGolangStruct(declaration, skipMagic, "")
	return fmt.Sprintf("type %s %v", typeName, s), err
}

func (g *Generator) generateGolangSumType(declarations []CombinatorDeclaration, typeName string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("type " + typeName + " struct{\n" + pkgRef(g.tlbPkg, "SumType") + "\n")
	for _, d := range declarations {
		s, err := g.generateGolangStruct(d, true, typeName)
		if err != nil {
			return "", err
		}
		builder.WriteString(utils.ToCamelCase(d.Constructor.Name))
		builder.WriteRune(' ')
		builder.WriteString(s)
		builder.WriteString(fmt.Sprintf(" `tlbSumType:\"%v\"`", d.Constructor.Prefix))
		builder.WriteRune('\n')
	}
	builder.WriteString("}\n")

	builder.WriteString(fmt.Sprintf(`func (t *%v) MarshalJSON() ([]byte, error) {`, typeName))
	builder.WriteString(`    switch t.SumType {`)
	for _, d := range declarations {
		name := utils.ToCamelCase(d.Constructor.Name)
		builder.WriteString(fmt.Sprintf(`case "%v": `, name))
		builder.WriteString(fmt.Sprintf(`bytes, err := json.Marshal(t.%v)`+"\n", name))
		builder.WriteString("if err != nil {\n")
		builder.WriteString("return nil, err\n")
		builder.WriteString("}\n")
		//builder.WriteString("return []byte(fmt.Sprintf(`{\"SumType\": %v}`, string(bytes))), nil")
		builder.WriteString("return []byte(fmt.Sprintf(`{\"SumType\": \"" + name + "\",\"" + name + "\":%v}`, string(bytes))), nil\n")
	}

	builder.WriteString("default: ")
	builder.WriteString(`return nil, fmt.Errorf("unknown sum type %v", t.SumType)`)
	builder.WriteString("}\n")
	builder.WriteString("}\n")
	return builder.String(), nil

}

func (g *Generator) generateGolangType(declarations []CombinatorDeclaration, typeName string, skipMagic bool) (string, error) {
	if len(declarations) == 1 {
		return g.generateGolangSimpleType(declarations[0], typeName, skipMagic)
	} else {
		return g.generateGolangSumType(declarations, typeName)
	}
}

type golangType struct {
	name   string
	tag    string
	params []golangType
}

func (t TypeExpression) toGolangType(g *Generator) (golangType, error) {
	if t.ParenExpression != nil {
		return t.ParenExpression.toGolangType(g)
	}
	if t.NamedRef != nil {
		return mapToGoType(*t.NamedRef, false, g.knownTypes, g.tlbPkg), nil
	}
	if t.BuiltIn != nil {
		return mapToGoType(*t.BuiltIn, false, g.knownTypes, g.tlbPkg), nil
	}
	if t.Number != nil {
		return mapToGoType(fmt.Sprintf("%d", *t.Number), false, g.knownTypes, g.tlbPkg), nil
	}
	if t.CellRef != nil {
		gt, err := t.CellRef.TypeExpression.toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		return golangType{
			name: fmt.Sprintf("%s", gt.String()),
			tag:  "^",
		}, nil
	}
	if t.AnonymousConstructor != nil {
		s, err := g.fieldDefinitionsToStruct(t.AnonymousConstructor.Values, "")
		if err != nil {
			return golangType{}, err
		}
		return golangType{
			name:   fmt.Sprintf("struct {\n%s\n}", s),
			tag:    "",
			params: nil,
		}, nil
	}

	return golangType{
		name: "UnknownType",
		tag:  "",
	}, nil
}

func (t *ParenExpression) toGolangType(g *Generator) (golangType, error) {
	var res golangType
	name, err := t.Name.toGolangType(g)
	if err != nil {
		return golangType{}, err
	}
	res.name = name.String()
	switch name.String() {
	case "Either":
		if len(t.Parameter) != 2 {
			return golangType{}, fmt.Errorf("invalid parameters qty for Either")
		}
		p1, err := t.Parameter[0].toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		p2, err := t.Parameter[1].toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		if p1.name == p2.name && p2.tag == "^" {
			// todo: compare tags?
			res.name = fmt.Sprintf("%s[%s]", pkgRef(g.tlbPkg, "EitherRef"), p1.String())
			return res, nil
		}
		res.name = fmt.Sprintf("%s[%s, %s]", pkgRef(g.tlbPkg, "Either"), p1.String(), p2.String())
		return res, nil
	case "HashmapE", "Hashmap":
		if len(t.Parameter) != 2 {
			return golangType{}, fmt.Errorf("invalid parameters qty for HashmapE")
		}
		if t.Parameter[0].Number == nil && t.Parameter[0].NamedRef == nil {
			return golangType{}, fmt.Errorf("invalid bitsize type for HashmapE")
		}
		p, err := t.Parameter[1].toGolangType(g)
		if p.tag == "^" {
			p.name = fmt.Sprintf("%s[%s]", pkgRef(g.tlbPkg, "Ref"), p.String())
		}
		if err != nil {
			return golangType{}, err
		}
		if t.Parameter[0].Number != nil {
			size := mapBitsSizeToType(*t.Parameter[0].Number, g.tlbPkg)
			res.name = fmt.Sprintf("%s[%s, %s]", pkgRef(g.tlbPkg, name.String()), size.String(), p.String())
			return res, nil
		}

		param0Type, ok := g.knownTypes[*t.Parameter[0].NamedRef]
		if !ok {
			return golangType{}, fmt.Errorf("unknown type %v", *t.Parameter[0].BuiltIn)
		}
		res.name = fmt.Sprintf("%s[%s, %s]", pkgRef(g.tlbPkg, name.String()), param0Type.Name, p.String())
		return res, nil
	case "Maybe":
		if len(t.Parameter) != 1 {
			return golangType{}, fmt.Errorf("invalid parameters qty for Maybe")
		}
		tag := "maybe"
		param := t.Parameter[0]
		if t.Parameter[0].CellRef != nil {
			tag = "maybe^"
			param = t.Parameter[0].CellRef.TypeExpression
		}
		p, err := param.toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		if len(p.tag) > 0 {
			return golangType{}, fmt.Errorf("can't combine tags: %v and %v", tag, p.tag)
		}
		res.name = fmt.Sprintf("*%s", p.String())
		res.tag = tag
		return res, nil
	case "VarUInteger":
		if len(t.Parameter) != 1 {
			return golangType{}, fmt.Errorf("invalid parameters qty for VarUInteger")
		}
		p, err := t.Parameter[0].toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		res.name = fmt.Sprintf("%s%s", pkgRef(g.tlbPkg, "VarUInteger"), p.String())
		return res, nil
	case "##":
		if len(t.Parameter) != 1 {
			return golangType{}, fmt.Errorf("invalid parameters qty for ##")
		}
		p, err := t.Parameter[0].toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		size := p.String()
		if size == "8" || size == "16" || size == "32" || size == "64" {
			res.name = fmt.Sprintf("uint%s", p.String())
		} else {
			res.name = fmt.Sprintf("%s%s", pkgRef(g.tlbPkg, "Uint"), p.String())
		}
		return res, nil
	}

	for _, p := range t.Parameter {
		param, err := p.toGolangType(g)
		if err != nil {
			return golangType{}, err
		}
		res.params = append(res.params, param)
	}
	return res, nil
}

func mapBitsSizeToType(bits int, tlbPkg string) golangType {
	if bits <= 64 {
		return golangType{
			name: fmt.Sprintf("%s%d", pkgRef(tlbPkg, "Uint"), bits),
		}
	}
	return golangType{
		name: fmt.Sprintf("%s%d", pkgRef(tlbPkg, "Bits"), bits),
	}
}

func mapToGoType(name string, optional bool, knownTypes map[string]DefaultType, tlbPkg string) golangType {
	goType, ok := knownTypes[name]
	if ok {
		return golangType{
			name: goType.Name,
		}
	}
	t, ok := parseBuildInInt(name, tlbPkg)
	if ok {
		return t
	}
	if name == "##" {
		return golangType{
			name: name,
		}
	}
	return golangType{
		name: utils.ToCamelCase(name),
	}
}

func parseBuildInInt(s string, tlbPkg string) (golangType, bool) {
	if strings.HasPrefix(s, "int") {
		last := strings.TrimPrefix(s, "int")
		bits, err := strconv.Atoi(last)
		if err != nil {
			return golangType{}, false
		}
		return golangType{
			name: fmt.Sprintf("%s%d", pkgRef(tlbPkg, "Int"), bits),
		}, true
	}

	if strings.HasPrefix(s, "uint") {
		last := strings.TrimPrefix(s, "uint")
		bits, err := strconv.Atoi(last)
		if err != nil {
			return golangType{}, false
		}
		return golangType{
			name: fmt.Sprintf("%s%d", pkgRef(tlbPkg, "Uint"), bits),
		}, true
	}

	if strings.HasPrefix(s, "bits") {
		last := strings.TrimPrefix(s, "bits")
		bits, err := strconv.Atoi(last)
		if err != nil {
			return golangType{}, false
		}
		return golangType{
			name: fmt.Sprintf("%s%d", pkgRef(tlbPkg, "Bits"), bits),
		}, true
	}

	return golangType{}, false
}

func (t golangType) String() string {
	return t.name
}
