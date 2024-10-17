package parser

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo/utils"
	"go/format"
	"strings"
	"text/template"
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
		"#":      {"uint32", false},
		"uint":   {"uint32", false},
		"uint8":  {"uint8", false},
		"uint16": {"uint16", false},
		"uint32": {"uint32", false},
		"uint64": {"uint64", false},
		"int":    {"uint32", false}, // this feels so wrong
		"int8":   {"int8", false},
		"int32":  {"int32", false},
		"int64":  {"int64", false},
		"int256": {"tl.Int256", false},
		"long":   {"uint64", false},
		"bytes":  {"[]byte", true},
		"Bool":   {"bool", false},
		"string": {"string", false},
	}

	unmarshalerReturnErr = "if err != nil {return err}\n"
	marshalerReturnErr   = "if err != nil {return nil, err}\n"
	functionReturnErr    = "if err != nil {return res, %s}\n"
)

var (
	//go:embed decoding.tmpl
	decodingTemplate string
)

type Generator struct {
	knownTypes      map[string]DefaultType
	newTlTypes      map[string]tlType
	typeName        string
	newRequestTypes map[string]tlType
}

func NewGenerator(knownTypes map[string]DefaultType, typeName string) *Generator {
	tlTypes := make(map[string]tlType)
	requestTypes := make(map[string]tlType)
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes:      knownTypes,
		newTlTypes:      tlTypes,
		typeName:        typeName,
		newRequestTypes: requestTypes,
	}
}

func (g *Generator) LoadTypes(declarations []CombinatorDeclaration) (string, error) {
	dec := make([][]CombinatorDeclaration, 0)
	for _, c := range declarations {
		f := false
		for i, c1 := range dec {
			if c1[0].Combinator == c.Combinator {
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
		t, err := g.generateGolangType(v)
		if err != nil {
			return "", err
		}
		if len(v) == 1 {
			g.newTlTypes[v[0].Constructor] = t
		} else {
			g.newTlTypes[v[0].Combinator] = t
		}
		s += "\n" + t.definition + "\n"
		unmarshaler, err := g.generateUnmarshalers(v, t.name)
		if err != nil {
			return "", err
		}
		marshaler, err := g.generateMarshalers(v, t.name)
		if err != nil {
			return "", err
		}
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
		requestType, err := g.generateGolangMethodRequestType(c)
		if err != nil {
			return "", err
		}
		g.newRequestTypes[c.Constructor] = requestType
		s += "\n" + requestType.definition + "\n"
		if len(c.FieldDefinitions) > 0 {
			marshaler, err := g.generateMarshalers([]CombinatorDeclaration{c}, requestType.name)
			if err != nil {
				return "", err
			}
			s += "\n" + marshaler + "\n"
		}
		unmarshler, err := g.generateUnmarshalers([]CombinatorDeclaration{c}, requestType.name)
		if err != nil {
			return "", err
		}
		s += "\n" + unmarshler + "\n"
		method, err := g.generateGolangMethod(g.typeName, c)
		if err != nil {
			return "", err
		}
		s += "\n" + method + "\n"
	}
	decoder := g.generateRequestDecoder()
	s += "\n" + decoder + "\n"

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
	name := utils.ToCamelCase(declaration.Constructor)
	s, err := g.generateGolangStruct(declaration)
	if err != nil {
		return tlType{}, err
	}
	tag, err := tagToUint32(declaration.Tag)
	if err != nil {
		return tlType{}, err
	}
	return tlType{
		name:       name + "C",
		tags:       []uint32{tag},
		definition: fmt.Sprintf("type %vC %v", name, s),
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

		t, err := toGolangType(e, optional, g.knownTypes, g.newTlTypes)
		if err != nil {
			return "", err
		}

		if t.name == "True" {
			continue
		}

		builder.WriteString(utils.ToCamelCase(name))
		builder.WriteRune('\t')
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

func mapToGoType(name string, optional bool, defaultTypes map[string]DefaultType, newTlTypes map[string]tlType) golangType {
	goType, ok := defaultTypes[name]
	if ok {
		return golangType{
			name:        goType.Name,
			optional:    optional,
			pointerType: goType.IsPointerType,
		}
	}
	knType, ok := newTlTypes[name]
	if ok {
		return golangType{
			name:        knType.name,
			optional:    optional,
			pointerType: false,
		}
	}
	return golangType{
		name:        utils.ToCamelCase(name),
		optional:    optional,
		pointerType: false,
	}
}

func toGolangType(t TypeExpression, optional bool, defaultTypes map[string]DefaultType, newTlTypes map[string]tlType) (golangType, error) {
	if t.BuiltIn != nil {
		return mapToGoType(*t.BuiltIn, optional, defaultTypes, newTlTypes), nil
	}
	if t.NamedRef != nil {
		return mapToGoType(*t.NamedRef, optional, defaultTypes, newTlTypes), nil
	}

	if t.Vector != nil {
		if len(t.Vector.Parameter) != 1 {
			return golangType{}, fmt.Errorf("vector must contains only one parameter")
		}
		gt, err := toGolangType(t.Vector.Parameter[0], false, defaultTypes, newTlTypes) // can not be pointer type under vector
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
	if len(declarations) == 1 && len(declarations[0].FieldDefinitions) > 0 {
		builder.WriteString("var err error\n")
		s, err := g.generateSimpleTypeUnmarshaler(declarations[0])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	} else if len(declarations) > 1 {
		builder.WriteString("var err error\n")
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
			gt, err := toGolangType(field.Expression, false, g.knownTypes, g.newTlTypes)
			if err != nil {
				return "", err
			}
			typeString := gt.String()
			builder.WriteString(fmt.Sprintf("if (t.%s>>%s)&1 == 1{\n",
				utils.ToCamelCase(field.Modificator.Name), field.Modificator.Bit))
			if typeString == "True" {
				// TODO: add field?
				builder.WriteString("var tTrue bool\n")
				builder.WriteString("err = tl.Unmarshal(r, &tTrue)\n")
				builder.WriteString(unmarshalerReturnErr)
				builder.WriteString("if tTrue != true {return fmt.Errorf(\"not a True type\")}\n")
				builder.WriteString("}\n")
				continue
			}
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
		builder.WriteString(fmt.Sprintf("case %#x:\n", tag))
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

func (g *Generator) generateMarshalers(declarations []CombinatorDeclaration, receiverType string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func (t %s) MarshalTL() ([]byte, error) {\n", receiverType))
	builder.WriteString("var (err error \n b []byte)\n")
	builder.WriteString("buf := new(bytes.Buffer)\n")
	if len(declarations) == 1 {
		s, err := g.generateSimpleTypeMarshaler(declarations[0])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	} else {
		s, err := g.generateSumTypeMarshaler(declarations)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	builder.WriteString("return buf.Bytes(), nil\n")
	builder.WriteRune('}')
	return builder.String(), nil
}

func (g *Generator) generateSimpleTypeMarshaler(declaration CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}
	code, err := g.generateMarshalerCode(declaration, "t")
	if err != nil {
		return "", err
	}
	builder.WriteString(code)
	return builder.String(), nil
}

func (g *Generator) generateMarshalerCode(declaration CombinatorDeclaration, receiverName string) (string, error) {
	builder := strings.Builder{}

	for _, field := range declaration.FieldDefinitions {
		name := utils.ToCamelCase(field.Name)
		if field.Modificator.Name != "" { // mode.0?field
			t, err := toGolangType(field.Expression, true, g.knownTypes, g.newTlTypes)
			if err != nil {
				return "", err
			}
			if t.name == "True" {
				continue // TODO: check
			}
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

func (g *Generator) generateSumTypeMarshaler(declarations []CombinatorDeclaration) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("switch t.SumType {\n")

	for _, d := range declarations {
		name := utils.ToCamelCase(d.Constructor)
		tag, err := tagToUint32(d.Tag)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("case \"%s\":\n", name))
		builder.WriteString(fmt.Sprintf("b, err = tl.Marshal(uint32(%#x))\n", tag))
		builder.WriteString(marshalerReturnErr)
		builder.WriteString("_, err = buf.Write(b)\n")
		code, err := g.generateMarshalerCode(d, "t."+name)
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
		return 0, fmt.Errorf("ivalid tag %s err: %s", tag, err.Error())
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

	errType, ok := g.newTlTypes["liteServer.error"]
	if !ok {
		return "", fmt.Errorf("LiteServerError not parsed")
	}
	if len(errType.tags) != 1 {
		return "", fmt.Errorf("invalid error tag")
	}

	var respType tlType
	ok = false
	for k, v := range g.newTlTypes {
		// TODO: valid if only one constructor OR type
		if strings.ToLower(c.Combinator) == strings.ToLower(k) {
			respType = v
			ok = true
			break
		}
	}
	if !ok {
		return "", fmt.Errorf("response type %s not parsed", responseName)
	}

	builder := strings.Builder{}

	// func signature
	builder.WriteString(fmt.Sprintf("func (c %s) %s(ctx context.Context", typeName, methodName))
	if len(c.FieldDefinitions) > 0 {
		builder.WriteString(fmt.Sprintf(", request %sRequest", methodName))
		builder.WriteString(fmt.Sprintf(") (res %s ,err error) {\n", respType.name))

		// request marshaling
		builder.WriteString(fmt.Sprintf("payload, err := tl.Marshal(struct{tl.SumType \n Req %sRequest", methodName))
		builder.WriteString(fmt.Sprintf(" `tlSumType:\"%08x\"`}{SumType: \"Req\", Req: request})\n", tag))
		builder.WriteString(fmt.Sprintf(functionReturnErr, "err"))
	} else {
		builder.WriteString(fmt.Sprintf(") (res %s ,err error) {\n", respType.name))
		builder.WriteString("payload := make([]byte, 4)\n")
		builder.WriteString(fmt.Sprintf("binary.LittleEndian.PutUint32(payload, %#x)\n", tag))
	}

	builder.WriteString("resp, err := c.liteServerRequest(ctx, payload)\n")
	builder.WriteString(fmt.Sprintf(functionReturnErr, "err"))

	builder.WriteString("if len(resp) < 4 {return res, fmt.Errorf(\"not enough bytes for tag\")}\n")
	builder.WriteString("tag := binary.LittleEndian.Uint32(resp[:4])\n")

	// lite server error processing
	builder.WriteString(fmt.Sprintf("if tag == %#x {\n", errType.tags[0]))
	builder.WriteString("var errRes LiteServerErrorC\n")
	builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)\n")
	builder.WriteString(fmt.Sprintf(functionReturnErr, "err"))
	builder.WriteString("return res, errRes\n")
	builder.WriteString("}\n")

	if len(respType.tags) == 0 {
		return "", fmt.Errorf("invalid response type %s tag", responseName)
	}

	if len(respType.tags) == 1 {
		// simple type response
		builder.WriteString(fmt.Sprintf("if tag == %#x {\n", respType.tags[0]))
		builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)\n")
		builder.WriteString("return res, err\n}\n")
		builder.WriteString("return res, fmt.Errorf(\"invalid tag\")\n")
	} else if len(respType.tags) > 1 {
		builder.WriteString("err = tl.Unmarshal(bytes.NewReader(resp), &res)\n")
		builder.WriteString("return res, err\n")
	}

	builder.WriteRune('}')
	return builder.String(), nil
}

func (g *Generator) generateGolangMethodRequestType(c CombinatorDeclaration) (tlType, error) {
	name := utils.ToCamelCase(c.Constructor) + "Request"
	s, err := g.generateGolangStruct(c)
	if err != nil {
		return tlType{}, err
	}
	tag, err := tagToUint32(c.Tag)
	if err != nil {
		return tlType{}, err
	}
	return tlType{
		name:       name,
		tags:       []uint32{tag},
		definition: fmt.Sprintf("type %s %s", name, s),
	}, nil
}

func (g *Generator) generateRequestDecoder() string {
	type msg struct {
		TlName string
		Name   string
		Tag    uint32
	}
	type messagesContext struct {
		Types      map[string]msg
		WhatRender string
	}
	reqContext := messagesContext{Types: map[string]msg{}, WhatRender: "Request"}
	for tlName, val := range g.newRequestTypes {
		if len(val.tags) != 1 {
			panic("sum type is not supported in the request")
		}
		reqContext.Types[val.name] = msg{
			Name:   val.name,
			Tag:    val.tags[0],
			TlName: tlName,
		}
	}
	tmpl, err := template.New("decoding").Parse(decodingTemplate)
	if err != nil {
		panic(err)
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, reqContext); err != nil {
		panic(err)
		return ""
	}
	fmt.Printf("%s", buf.String())
	return buf.String()
}
