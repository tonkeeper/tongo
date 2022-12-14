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

type tlType struct {
	name       string
	tags       []uint32
	definition string
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
	functionReturnErr    = "if err != nil {return %s{}, %s}\n"
)

type Generator struct {
	knownTypes map[string]DefaultType
	newTlTypes map[string]tlType
	typeName   string
}

func NewGenerator(knownTypes map[string]DefaultType, typeName string) *Generator {
	tlTypes := make(map[string]tlType)
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes: knownTypes,
		newTlTypes: tlTypes,
		typeName:   typeName,
	}
}

func (g *Generator) LoadTypes(declarations []CombinatorDeclaration) (string, error) {
	sumTypes := make(map[string][]CombinatorDeclaration)

	for _, c := range declarations {
		sumTypes[c.Combinator] = append(sumTypes[c.Combinator], c)
	}

	s := ""
	for _, v := range sumTypes {
		t, err := g.generateGolangType(v)
		if err != nil {
			return "", err
		}
		g.newTlTypes[t.name] = t
		unmarshaler, err := g.generateUnmarshalers(v, t.name)
		if err != nil {
			return "", err
		}
		marshaler, err := generateMarshalers(v, t.name)
		if err != nil {
			return "", err
		}
		s += "\n" + t.definition + "\n"
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
	s := ""
	for _, c := range functions {
		if len(c.FieldDefinitions) > 0 {
			name, requestType, err := g.generateGolangMethodRequestType(c)
			if err != nil {
				return "", err
			}
			marshaler, err := generateMarshalers([]CombinatorDeclaration{c}, name)
			if err != nil {
				return "", err
			}
			s += "\n" + requestType + "\n"
			s += "\n" + marshaler + "\n"
		}
		method, err := g.generateGolangMethod(g.typeName, c)
		if err != nil {
			return "", err
		}
		s += "\n" + method + "\n"
	}
	b, err := format.Source([]byte(s))
	if err != nil {
		return s, err
	}
	return string(b), err
}

func (g *Generator) generateGolangType(declarations []CombinatorDeclaration) (tlType, error) {
	if len(declarations) == 1 {
		return g.generateGolangSimpleType(declarations[0])
	} else {
		return g.generateGolangSumType(declarations)
	}
}

func (g *Generator) generateGolangSimpleType(declaration CombinatorDeclaration) (tlType, error) {
	name := utils.ToCamelCase(declaration.Combinator)
	s, err := g.generateGolangStruct(declaration)
	if err != nil {
		return tlType{}, err
	}
	tag, err := tagToUint32(declaration.Tag)
	if err != nil {
		return tlType{}, err
	}
	return tlType{
		name:       name,
		tags:       []uint32{tag},
		definition: fmt.Sprintf("type %v %v", name, s),
	}, nil
}

func (g *Generator) generateGolangSumType(declarations []CombinatorDeclaration) (tlType, error) {
	name := utils.ToCamelCase(declarations[0].Combinator)
	var tags []uint32
	builder := strings.Builder{}
	builder.WriteString("type " + name + " struct{\ntl.SumType\n")
	for _, d := range declarations {
		tag, err := tagToUint32(d.Tag)
		if err != nil {
			return tlType{}, err
		}
		tags = append(tags, tag)
		s, err := g.generateGolangStruct(d)
		if err != nil {
			return tlType{}, err
		}
		builder.WriteString(utils.ToCamelCase(d.Constructor))
		builder.WriteRune(' ')
		builder.WriteString(s)
		builder.WriteRune('\n')
	}
	builder.WriteRune('}')
	return tlType{
		name:       name,
		tags:       tags,
		definition: builder.String(),
	}, nil
}

func (g *Generator) generateGolangStruct(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("struct{")
	if len(declaration.FieldDefinitions) > 0 {
		builder.WriteRune('\n')
	}
	for i, field := range declaration.FieldDefinitions {
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
		t, err := toGolangType(e, optional, g.knownTypes)
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

func mapToGoType(name string, optional bool, defaultTypes map[string]DefaultType) golangType {
	goType, ok := defaultTypes[name]
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

func toGolangType(t TypeExpression, optional bool, defaultTypes map[string]DefaultType) (golangType, error) {
	if t.BuiltIn != nil {
		return mapToGoType(*t.BuiltIn, optional, defaultTypes), nil
	}
	if t.NamedRef != nil {
		return mapToGoType(*t.NamedRef, optional, defaultTypes), nil
	}

	if t.Vector != nil {
		if len(t.Vector.Parameter) != 1 {
			return golangType{}, fmt.Errorf("vector must contains only one parameter")
		}
		gt, err := toGolangType(t.Vector.Parameter[0], false, defaultTypes) // can not be pointer type under vector
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

func (g *Generator) generateUnmarshalers(declarations []CombinatorDeclaration, receiverType string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func (t *%s) UnmarshalTL(r io.Reader) error {\n", receiverType))
	builder.WriteString("var err error\n")
	if len(declarations) == 1 {
		s, err := g.generateSimpleTypeUnmarshaler(declarations[0])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	} else {
		s, err := g.generateSumTypeUnmarshaler(declarations)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}

	builder.WriteString("return nil\n")
	builder.WriteRune('}')
	return builder.String(), nil

}

func (g *Generator) generateSimpleTypeUnmarshaler(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	code, err := g.generateUnmarshalerCode(declaration, "t")
	if err != nil {
		return "", err
	}
	builder.WriteString(code)
	return builder.String(), nil
}

func (g *Generator) generateUnmarshalerCode(declaration CombinatorDeclaration, receiverName string) (string, error) {
	builder := strings.Builder{}

	for _, field := range declaration.FieldDefinitions {
		name := utils.ToCamelCase(field.Name)
		if field.Modificator.Name != "" { // mode.0?field
			gt, err := toGolangType(field.Expression, false, g.knownTypes)
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

func (g *Generator) generateSumTypeUnmarshaler(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("var b [4]byte\n")
	builder.WriteString("_, err = io.ReadFull(r, b[:])\n")
	builder.WriteString(unmarshalerReturnErr)
	builder.WriteString("tag := int(binary.LittleEndian.Uint32(b[:]))\n")
	builder.WriteString("switch tag {\n")

	for _, d := range declarations {
		tag, err := tagToUint32(d.Tag)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("case 0x%x:\n", tag))
		name := utils.ToCamelCase(d.Constructor)
		builder.WriteString("t.SumType = \"" + name + "\"\n")
		code, err := g.generateUnmarshalerCode(d, "t."+name)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	builder.WriteString("default: return fmt.Errorf(\"invalid tag\")}\n")
	return builder.String(), nil
}

func generateMarshalers(declarations []CombinatorDeclaration, receiverType string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func (t %s) MarshalTL() ([]byte, error) {\n", receiverType))
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

	for _, field := range declaration.FieldDefinitions {
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
		tag, err := tagToUint32(d.Tag)
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

func tagToUint32(tag string) (uint32, error) {
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
	return binary.BigEndian.Uint32(b1[:]), nil
}

func (g *Generator) generateGolangMethod(typeName string, c CombinatorDeclaration) (string, error) {
	methodName := utils.ToCamelCase(c.Constructor)
	responseName := utils.ToCamelCase(c.Combinator)
	tag, err := tagToUint32(c.Tag)
	if err != nil {
		return "", err
	}

	errType, ok := g.newTlTypes["LiteServerError"]
	if !ok {
		return "", fmt.Errorf("LiteServerError not parsed")
	}
	if len(errType.tags) != 1 {
		return "", fmt.Errorf("invalid error tag")
	}

	respType, ok := g.newTlTypes[responseName]
	if !ok {
		return "", fmt.Errorf("response type %s not parsed", responseName)
	}

	builder := strings.Builder{}

	// func signature
	builder.WriteString(fmt.Sprintf("func (c %s) %s(ctx context.Context", typeName, methodName))
	if len(c.FieldDefinitions) > 0 {
		builder.WriteString(fmt.Sprintf(", request %sRequest", methodName))
		builder.WriteString(fmt.Sprintf(") (%s ,error) {\n", responseName))

		// request marshaling
		builder.WriteString(fmt.Sprintf("payload, err := tl.Marshal(struct{tl.SumType \n Req %sRequest", methodName))
		builder.WriteString(fmt.Sprintf(" `tlSumType:\"%x\"`}{Req: request})\n", tag))
		builder.WriteString(fmt.Sprintf(functionReturnErr, responseName, "err"))
	} else {
		builder.WriteString(fmt.Sprintf(") (%s ,error) {\n", responseName))
		builder.WriteString("payload := make([]byte, 4)\n")
		builder.WriteString(fmt.Sprintf("binary.BigEndian.PutUint32(payload, 0x%x)\n", tag))
	}

	builder.WriteString("req := makeLiteServerQueryRequest(payload)\n")
	builder.WriteString("server := c.getMasterchainServer()\n")
	builder.WriteString("resp, err := server.Request(ctx, req)\n")
	builder.WriteString(fmt.Sprintf(functionReturnErr, responseName, "err"))

	builder.WriteString(fmt.Sprintf("if len(resp) < 4 {return %s{}, fmt.Errorf(\"not enought bytes for tag\")}\n",
		responseName))
	builder.WriteString("tag := binary.BigEndian.Uint32(resp[:4])\n")

	// lite server error processing
	builder.WriteString(fmt.Sprintf("if tag == 0x%x {\n", errType.tags[0]))
	builder.WriteString("var errRes LiteServerError\n")
	builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)\n")
	builder.WriteString(fmt.Sprintf(functionReturnErr, responseName, "err"))
	builder.WriteString(fmt.Sprintf("return %s{}, errRes\n", responseName))
	builder.WriteString("}\n")

	if len(respType.tags) == 0 {
		return "", fmt.Errorf("invalid response type %s tag", responseName)
	}

	if len(respType.tags) == 1 {
		// simple type response
		builder.WriteString(fmt.Sprintf("if tag == 0x%x {\n", respType.tags[0]))
		builder.WriteString(fmt.Sprintf("var res %s\n", responseName))
		builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)\n")
		builder.WriteString("return res, err\n}\n")
		builder.WriteString(fmt.Sprintf("return %s{}, fmt.Errorf(\"invalid tag\")\n", responseName))
	} else if len(respType.tags) > 1 {
		builder.WriteString(fmt.Sprintf("var res %s\n", responseName))
		builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp), &res)\n")
		builder.WriteString("return res, err\n")
	}

	builder.WriteRune('}')
	return builder.String(), nil
}

func (g *Generator) generateGolangMethodRequestType(c CombinatorDeclaration) (string, string, error) {
	s, err := g.generateGolangStruct(c)
	name := utils.ToCamelCase(c.Constructor) + "Request"
	return name, fmt.Sprintf("type %s %s", name, s), err
}

// TODO: add custom method
//func (t LiteServerError) Error() string {
//	return fmt.Sprintf("error code: %d message: %s", t.Code, t.Message)
//}
