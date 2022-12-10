package parser

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/startfellows/tongo/utils"
	"go/format"
	"strings"
)

type DefaultType struct {
	Name          string
	IsPointerType bool
}

var (
	defaultKnownTypes = map[string]DefaultType{
		"#":      {"int32", false},
		"int":    {"int32", false},
		"int256": {"tl.Int256", false},
		"long":   {"int64", false},
		"bytes":  {"[]byte", true},
		"Bool":   {"bool", false},
		"string": {"string", false},
	}

	unmarshalerReturnErr = "if err != nil {return err}\n"
	marshalerReturnErr   = "if err != nil {return nil, err}\n"
)

type Generator struct {
	knownTypes map[string]DefaultType
	newTlTypes []string
	typeName   string
}

func NewGenerator(knownTypes map[string]DefaultType, typeName string) *Generator {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes: knownTypes,
		typeName:   typeName,
	}
}

func (g *Generator) LoadTypes(Declarations []CombinatorDeclaration) (string, error) {
	sumTypes := make(map[string][]CombinatorDeclaration)

	for _, c := range Declarations {
		sumTypes[c.Combinator] = append(sumTypes[c.Combinator], c)
	}

	s := ""
	for _, v := range sumTypes {
		typeString, err := generateGolangType(v)
		if err != nil {
			return "", err
		}
		g.newTlTypes = append(g.newTlTypes, typeString)

		unmarshaler, err := generateUnmarshalers(v)
		if err != nil {
			return "", err
		}
		marshaler, err := generateMarshalers(v)
		if err != nil {
			return "", err
		}
		s += "\n" + typeString + "\n"
		s += "\n" + marshaler + "\n"
		s += "\n" + unmarshaler + "\n"
	}

	b, err := format.Source([]byte(s))
	if err != nil {
		return s, err
	}

	return string(b), err
}

func (g *Generator) LoadFunctions(functions []CombinatorDeclaration) (string, error) {
	return "", fmt.Errorf("not implemnted")
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
	name        string
	optional    bool
	pointerType bool
}

func (t golangType) String() string {
	if !t.optional || t.pointerType {
		return t.name
	}
	return "*" + t.name
}

func mapToGoType(name string, optional bool) golangType {
	goType, ok := defaultKnownTypes[name]
	if ok {
		return golangType{
			name:        goType.Name,
			optional:    optional,
			pointerType: goType.IsPointerType,
		}
	}
	return golangType{
		name:        utils.ToCamelCase(name),
		optional:    optional,
		pointerType: false,
	}
}

func toGolangType(t TypeExpression, optional bool) (golangType, error) {
	if t.BuiltIn != nil {
		return mapToGoType(*t.BuiltIn, optional), nil
	}
	if t.NamedRef != nil {
		return mapToGoType(*t.NamedRef, optional), nil
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
			name:        "[]" + gt.String(),
			optional:    optional,
			pointerType: true,
		}, nil
	}
	return golangType{}, fmt.Errorf("invalid type expression")
}

func generateUnmarshalers(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func (t *%s) UnmarshalTL(r io.Reader) error {\n",
		utils.ToCamelCase(declarations[0].Combinator)))
	builder.WriteString("var err error\n")
	if len(declarations) == 1 {
		s, err := generateSimpleTypeUnmarshaler(declarations[0])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	} else {
		s, err := generateSumTypeUnmarshaler(declarations)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}

	builder.WriteString("return nil\n")
	builder.WriteRune('}')
	return builder.String(), nil

}

func generateSimpleTypeUnmarshaler(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	code, err := generateUnmarshalerCode(declaration, "t")
	if err != nil {
		return "", err
	}
	builder.WriteString(code)
	return builder.String(), nil
}

func generateUnmarshalerCode(declaration CombinatorDeclaration, receiverName string) (string, error) {
	builder := strings.Builder{}

	for i, field := range declaration.FieldDefinitions {
		if field == nil {
			return "", fmt.Errorf("nil field %v in %v", i, declaration.Constructor)
		}

		name := utils.ToCamelCase(field.Name)

		if field.Modificator.Name != "" { // mode.0?field
			gt, err := toGolangType(field.Expression, false)
			if err != nil {
				return "", err
			}
			typeString := gt.String()
			builder.WriteString(fmt.Sprintf("if (t.%s>>%s)&1 == 1{\n",
				utils.ToCamelCase(field.Modificator.Name), field.Modificator.Bit))
			builder.WriteString(fmt.Sprintf("var temp%s %s\n", name, typeString))
			builder.WriteString("err = tl.Unmarshal(r, &temp" + name + ")\n")
			builder.WriteString(unmarshalerReturnErr)
			link := "&"
			if gt.pointerType {
				link = ""
			}
			builder.WriteString(receiverName + "." + name + " = " + link + "temp" + name)
			builder.WriteString("}\n")
		} else {
			builder.WriteString("err = tl.Unmarshal(r, &" + receiverName + "." + name + ")\n")
			builder.WriteString(unmarshalerReturnErr)
		}

	}
	return builder.String(), nil
}

func generateSumTypeUnmarshaler(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("var b [4]byte\n")
	builder.WriteString("_, err = io.ReadFull(r, b[:])\n")
	builder.WriteString(unmarshalerReturnErr)
	builder.WriteString("tag := int(binary.LittleEndian.Uint32(b[:]))\n")
	builder.WriteString("switch tag {\n")

	for _, d := range declarations {
		tag, err := tagToInt(d.Tag)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("case 0x%x:\n", tag))
		name := utils.ToCamelCase(d.Constructor)
		builder.WriteString("t.SumType = \"" + name + "\"\n")
		code, err := generateUnmarshalerCode(d, "t."+name)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	builder.WriteString("default: return fmt.Errorf(\"invalid tag\")}\n")
	return builder.String(), nil
}

func generateMarshalers(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func (t %s) MarshalTL() ([]byte, error) {\n",
		utils.ToCamelCase(declarations[0].Combinator)))
	builder.WriteString("var (err error \n b []byte)\n")
	builder.WriteString("buf := new(bytes.Buffer)\n")
	if len(declarations) == 1 {
		s, err := generateSimpleTypeMarshaler(declarations[0])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	} else {
		s, err := generateSumTypeMarshaler(declarations)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	builder.WriteString("return buf.Bytes(), nil\n")
	builder.WriteRune('}')
	return builder.String(), nil
}

func generateSimpleTypeMarshaler(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	code, err := generateMarshalerCode(declaration, "t")
	if err != nil {
		return "", err
	}
	builder.WriteString(code)
	return builder.String(), nil
}

func generateMarshalerCode(declaration CombinatorDeclaration, receiverName string) (string, error) {
	builder := strings.Builder{}

	for i, field := range declaration.FieldDefinitions {
		if field == nil {
			return "", fmt.Errorf("nil field %v in %v", i, declaration.Constructor)
		}

		name := utils.ToCamelCase(field.Name)

		if field.Modificator.Name != "" { // mode.0?field
			builder.WriteString(fmt.Sprintf("if (t.%s>>%s)&1 == 1{\n",
				utils.ToCamelCase(field.Modificator.Name), field.Modificator.Bit))
			builder.WriteString("b, err = tl.Marshal(" + receiverName + "." + name + ")\n")
			builder.WriteString(marshalerReturnErr)
			builder.WriteString("_, err = buf.Write(b)\n")
			builder.WriteString(marshalerReturnErr)
			builder.WriteString("}\n")
		} else {
			builder.WriteString("b, err = tl.Marshal(" + receiverName + "." + name + ")\n")
			builder.WriteString(marshalerReturnErr)
			builder.WriteString("_, err = buf.Write(b)\n")
			builder.WriteString(marshalerReturnErr)
		}

	}
	return builder.String(), nil
}

func generateSumTypeMarshaler(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("switch t.SumType {\n")

	for _, d := range declarations {
		name := utils.ToCamelCase(d.Constructor)
		tag, err := tagToInt(d.Tag)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("case \"%s\":\n", name))
		builder.WriteString(fmt.Sprintf("b, err = tl.Marshal(uint32(0x%x))\n", tag))
		builder.WriteString(marshalerReturnErr)
		builder.WriteString("_, err = buf.Write(b)\n")
		code, err := generateMarshalerCode(d, "t."+name)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	builder.WriteString("default: return nil, fmt.Errorf(\"invalid sum type\")}\n")
	return builder.String(), nil
}

func tagToInt(tag string) (int, error) {
	if !strings.HasPrefix(tag, "#") {
		return 0, fmt.Errorf("invalid tag prefix")
	}
	tag = strings.TrimPrefix(tag, "#")
	b, err := hex.DecodeString(tag)
	if err != nil {
		return 0, err
	}
	b1 := make([]byte, 4)
	copy(b1[:], b)
	return int(binary.LittleEndian.Uint32(b1[:])), nil
}
