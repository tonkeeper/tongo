package parser

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"

	"github.com/tonkeeper/tongo/tlb"
	tlbParser "github.com/tonkeeper/tongo/tlb/parser"
	"github.com/tonkeeper/tongo/utils"
)

var defaultKnownTypes = map[string]string{
	"accountid":     "tongo.AccountID",
	"cell":          "boc.Cell",
	"int8":          "int8",
	"int64":         "int64",
	"bool":          "bool",
	"uint16":        "uint16",
	"uint32":        "uint32",
	"uint64":        "uint64",
	"int256":        "tlb.Int256",
	"int257":        "tlb.Int257",
	"any":           "tlb.Any",
	"[]byte":        "[]byte",
	"coins":         "tlb.Grams",
	"big.int":       "big.Int",
	"dnsrecord":     "tlb.DNSRecord",
	"dns_recordset": "tlb.DNSRecordSet",
	"msgaddress":    "tlb.MsgAddress",
}

var (
	msgDecoderReturnErr = "if err != nil {return \"\", nil, err}\n"
	returnInvalidStack  = "{return \"\", nil, fmt.Errorf(\"invalid stack format\")}\n"
	returnStrNilErr     = "if err != nil {return \"\", nil, err}\n"
	//go:embed invocation_order.tmpl
	invocationOrderTemplate string
)

type TLBMsgBody struct {
	TypePrefix string
	TypeName   string
	Tag        uint64
	Code       string
}

type Generator struct {
	knownTypes        map[string]string
	interfaces        map[string]Interface
	newTlbTypes       map[string]struct{}
	loadedTlbTypes    []string
	loadedTlbMsgTypes map[uint32]TLBMsgBody
	typeName          string
}

func NewGenerator(knownTypes map[string]string, typeName string) *Generator {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	return &Generator{
		knownTypes:        knownTypes,
		interfaces:        map[string]Interface{},
		typeName:          typeName,
		loadedTlbMsgTypes: make(map[uint32]TLBMsgBody),
		newTlbTypes:       make(map[string]struct{}),
	}
}

func (g *Generator) GetMethods() (string, error) {
	var builder, methodMap, resultMap strings.Builder
	builder.WriteString(`
type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID tongo.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

	`)

	methods := make(map[int][]string)
	var resultTypes []string

	usedNames := map[string]struct{}{}

	for _, i := range g.interfaces {
		for _, m := range i.Methods {
			methodName := m.GolangFunctionName()
			var methodID int
			if m.ID != 0 {
				methodID = m.ID
			} else {
				methodID = utils.MethodIdFromName(m.Name)
			}
			if _, ok := usedNames[methodName]; ok {
				continue
			}
			usedNames[methodName] = struct{}{}

			if len(m.Input.StackValues) == 0 {
				methods[methodID] = append(methods[methodID], methodName)
			}

			for _, o := range m.Output {
				resultTypeName := utils.ToCamelCase(o.Version) + methodName + "Result"
				resultTypes = append(resultTypes, resultTypeName)
				r, err := g.buildResultType(resultTypeName, o.Stack)
				if err != nil {
					return "", err
				}
				builder.WriteString(r)
				builder.WriteRune('\n')
			}

			s, err := g.getMethod(m, methodID, methodName)
			if err != nil {
				return "", err
			}
			builder.WriteString(s)
			builder.WriteRune('\n')
		}
	}

	methodMap.WriteString("var KnownSimpleGetMethods = map[int][]func(ctx context.Context, executor Executor, reqAccountID tongo.AccountID) (string, any, error){\n")
	for _, k := range utils.GetOrderedKeys(methods) {
		methodMap.WriteString(fmt.Sprintf("%d: {", k))
		for _, m := range methods[k] {
			methodMap.WriteString(fmt.Sprintf("%s, ", m))
		}
		methodMap.WriteString("},\n")
	}
	methodMap.WriteString("}\n\n")

	resultMap.WriteString("var ResultTypes = []interface{}{\n")
	slices.Sort(resultTypes)
	for _, r := range resultTypes {
		resultMap.WriteString(fmt.Sprintf("&%s{}, \n", r))
	}
	resultMap.WriteString("}\n\n")

	b, err := format.Source([]byte(methodMap.String() + resultMap.String() + builder.String()))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *Generator) RegisterInterfaces(interfaces []Interface) error {
	for _, i := range interfaces {
		g.interfaces[i.Name] = i
		if i.Types != "" {
			err := g.registerType(i.Types)
			if err != nil {
				return err
			}
		}

		for _, internal := range i.Internals {
			err := g.registerMsgType(i.Name, internal.Input)
			if err != nil {
				return err
			}
			for _, out := range internal.Outputs {
				err := g.registerMsgType(i.Name, out)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *Generator) registerType(s string) error {
	tlbData, err := tlbParser.Parse(s)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}
	if len(tlbData.Declarations) == 0 {
		return fmt.Errorf("can't parse type %v", s)
	}

	gen := tlbParser.NewGenerator(nil, "")
	_, err = gen.LoadTypes(tlbData.Declarations, "", false)
	if err != nil {
		return fmt.Errorf("load types error: %v", err)
	}
	types := gen.GetTlbTypes()
	for _, t := range types {
		g.newTlbTypes[t.Name] = struct{}{}
		g.loadedTlbTypes = append(g.loadedTlbTypes, t.Definition)
	}
	return nil
}

func (g *Generator) registerMsgType(name, s string) error {
	parsed, err := tlbParser.Parse(s)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}
	if len(parsed.Declarations) != 1 {
		return fmt.Errorf("must be only one declaration for MsgBody %v", s)
	}

	gen := tlbParser.NewGenerator(nil, name)

	tag, err := tlb.ParseTag(parsed.Declarations[0].Constructor.Prefix)
	if err != nil {
		return fmt.Errorf("can't decode tag error %w", err)
	}
	if tag.Len != 32 {
		return fmt.Errorf("message %s body tag must be 32 bit lenght", parsed.Declarations[0].Constructor.Name)
	}

	var typePrefix string
	_, ok := g.loadedTlbMsgTypes[uint32(tag.Val)]
	if ok {
		typePrefix = utils.ToCamelCase(parsed.Declarations[0].Constructor.Name)
	} else {
		typePrefix = fmt.Sprintf("%s%s", utils.ToCamelCase(name), utils.ToCamelCase(parsed.Declarations[0].Constructor.Name))
	}

	t, err := gen.LoadTypes(parsed.Declarations, typePrefix, true)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}

	g.loadedTlbMsgTypes[uint32(tag.Val)] = TLBMsgBody{
		TypePrefix: typePrefix,
		TypeName:   utils.ToCamelCase(parsed.Declarations[0].Combinator.Name),
		Tag:        tag.Val,
		Code:       t,
	}

	return nil
}

func (g *Generator) checkType(s string) (string, error) {
	if typeName, prs := g.knownTypes[strings.ToLower(s)]; prs {
		return typeName, nil
	}
	_, ok := g.newTlbTypes[s]
	if !ok {
		return "", fmt.Errorf("not defined type: %s", s)
	}
	return s, nil
}

func (g *Generator) getMethod(m GetMethod, methodID int, methodName string) (string, error) {
	var builder strings.Builder
	var args []string

	builder.WriteString(fmt.Sprintf("func %v(ctx context.Context, executor Executor, reqAccountID tongo.AccountID, ", m.GolangFunctionName()))

	for _, s := range m.Input.StackValues {
		t, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		args = append(args, fmt.Sprintf("%v %v", utils.ToCamelCasePrivate(s.Name), t))
	}
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(")")

	if len(m.Output) == 0 {
		return "", fmt.Errorf("empty output for get method")
	}

	builder.WriteString(" (string, any, error) {\n")

	builder.WriteString(buildInputStackValues(m.Input.StackValues))
	builder.WriteRune('\n')

	builder.WriteString(fmt.Sprintf("// MethodID = %d for \"%s\" method\n", methodID, methodName))
	builder.WriteString(fmt.Sprintf("errCode, stack, err := executor.RunSmcMethodByID(ctx, reqAccountID, %d, stack)\n", methodID))
	builder.WriteString(returnStrNilErr)
	builder.WriteString("if errCode != 0 && errCode != 1 {return \"\", nil, fmt.Errorf(\"method execution failed with code: %v\", errCode)}\n")

	if len(m.Output) == 1 {
		name := utils.ToCamelCase(m.Output[0].Version) + methodName + "Result"
		resDecoder, err := g.buildOutputDecoder(name, m.Output[0].Stack, m.Output[0].FixedLength)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("return decode%s(stack)\n", name))
		builder.WriteString("}\n\n")
		builder.WriteString(resDecoder)
		builder.WriteRune('\n')
	} else {
		decoders := ""
		builder.WriteString("for _, f := range []func(tlb.VmStack)(string, any, error){")
		for _, o := range m.Output {
			name := utils.ToCamelCase(o.Version) + methodName + "Result"
			builder.WriteString(fmt.Sprintf("decode%s, ", name))
			resDecoder, err := g.buildOutputDecoder(name, o.Stack, o.FixedLength)
			if err != nil {
				return "", err
			}
			decoders = decoders + "\n\n" + resDecoder
		}
		builder.WriteString("} {\n")
		builder.WriteString("s, r, err := f(stack)\n")
		builder.WriteString("if err == nil {return s, r, nil}\n")
		builder.WriteString("}\n")
		builder.WriteString("return \"\", nil, fmt.Errorf(\"can not decode outputs\")\n}\n")
		builder.WriteString(decoders)
		builder.WriteRune('\n')
	}
	return builder.String(), nil
}

func buildInputStackValues(r []StackRecord) string {
	var builder strings.Builder
	builder.WriteString("stack := tlb.VmStack{}\n")

	if len(r) > 0 {
		builder.WriteString("var (val tlb.VmStackValue\n err error)\n")
	}

	for _, s := range r {
		switch s.XMLName.Local {
		case "tinyint":
			builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkTinyInt\", VmStkTinyInt: int64(%s)}\n",
				utils.ToCamelCasePrivate(s.Name)))
		case "int":
			builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkInt\", VmStkInt: %s}\n",
				utils.ToCamelCasePrivate(s.Name)))
		case "slice":
			builder.WriteString(fmt.Sprintf("val, err = tlb.TlbStructToVmCellSlice(%s)\n",
				utils.ToCamelCasePrivate(s.Name)))
			builder.WriteString(returnStrNilErr)
		case "cell":
			builder.WriteString(fmt.Sprintf("val, err = tlb.TlbStructToVmCell(%s)\n",
				utils.ToCamelCasePrivate(s.Name)))
			builder.WriteString(returnStrNilErr)
		}
		builder.WriteString("stack.Put(val)\n")
	}

	return builder.String()
}

func buildOutputStackCheck(r []StackRecord, isFixed bool) string {
	var builder strings.Builder
	if isFixed {
		builder.WriteString(fmt.Sprintf("if len(stack) != %d ", len(r)))
	} else {
		builder.WriteString(fmt.Sprintf("if len(stack) < %d ", len(r)))
	}
	for i, s := range r {
		nullableCheck := ""
		stackType := fmt.Sprintf("stack[%d].SumType", i)
		if s.Nullable {
			nullableCheck = fmt.Sprintf(" && %s != \"VmStkNull\"", stackType)
		}
		switch s.XMLName.Local {
		case "tinyint":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkTinyInt\"%s) ", stackType, nullableCheck))
		case "int":
			builder.WriteString(fmt.Sprintf("|| ((%s != \"VmStkTinyInt\" && %s != \"VmStkInt\")%s) ", stackType, stackType, nullableCheck))
		case "slice":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkSlice\"%s) ", stackType, nullableCheck))
		case "cell":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkCell\"%s) ", stackType, nullableCheck))
		}
	}
	builder.WriteString(returnInvalidStack)
	return builder.String()
}

func (g *Generator) buildOutputDecoder(name string, r []StackRecord, isFixed bool) (string, error) {
	var builder, resBuilder strings.Builder

	builder.WriteString(fmt.Sprintf("func decode%s(stack tlb.VmStack) (resultType string, resultAny any, err error) {\n", name))

	builder.WriteString(buildOutputStackCheck(r, isFixed))

	resBuilder.WriteString(fmt.Sprintf("%s{\n", name))

	for i, s := range r {
		varType, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		varName := utils.ToCamelCasePrivate(s.Name)
		recName := varName
		if s.Nullable {
			builder.WriteString(fmt.Sprintf("var %s *%s\n", varName, varType))
			builder.WriteString(fmt.Sprintf("if stack[%d].SumType != \"VmStkNull\" {", i))
			builder.WriteString(fmt.Sprintf("var temp %s\n", varType))
			recName = "temp"
		} else {
			builder.WriteString(fmt.Sprintf("var %s %s\n", varName, varType))
		}

		resBuilder.WriteString(fmt.Sprintf("%s: %s,\n", utils.ToCamelCase(s.Name), varName))

		switch s.XMLName.Local {
		case "tinyint":
			if varType == "bool" {
				builder.WriteString(fmt.Sprintf("%s = stack[%d].Int64() != 0\n", recName, i))
			} else {
				builder.WriteString(fmt.Sprintf("%s = %s(stack[%d].Int64())\n", recName, varType, i))
			}
		case "int":
			builder.WriteString(fmt.Sprintf("%s = stack[%d].Int257()\n", recName, i))
		case "slice":
			builder.WriteString(fmt.Sprintf("err = stack[%d].VmStkSlice.UnmarshalToTlbStruct(&%s)\n", i, recName))
			builder.WriteString(returnStrNilErr)
		case "cell":
			builder.WriteString(fmt.Sprintf("%sCell := &stack[%d].VmStkCell.Value\n", recName, i))
			builder.WriteString(fmt.Sprintf("err = tlb.Unmarshal(%sCell, &%s)\n", recName, recName))
			builder.WriteString(returnStrNilErr)
		}

		if s.Nullable {
			builder.WriteString(fmt.Sprintf("%s = &temp\n", varName))
			builder.WriteString("}\n")
		}
	}

	resBuilder.WriteString("}")
	builder.WriteString(fmt.Sprintf("return \"%s\",%s, nil", name, resBuilder.String()))

	builder.WriteString("}")
	return builder.String(), nil
}

func (g *Generator) buildResultType(name string, s []StackRecord) (string, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("type %s struct {\n", name))
	for _, c := range s {
		t, err := g.checkType(c.Type)
		if err != nil {
			return "", err
		}
		if c.Nullable {
			builder.WriteString(fmt.Sprintf("%s *%s\n", utils.ToCamelCase(c.Name), t))
		} else {
			builder.WriteString(fmt.Sprintf("%s %s\n", utils.ToCamelCase(c.Name), t))
		}
	}
	builder.WriteString("}\n")
	return builder.String(), nil
}

func (g *Generator) CollectedTypes() string {
	var builder strings.Builder
	builder.WriteString(strings.Join(g.loadedTlbTypes, "\n\n"))

	for _, k := range utils.GetOrderedKeys(g.loadedTlbMsgTypes) {
		builder.WriteString(g.loadedTlbMsgTypes[k].Code)
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (g *Generator) GenerateMsgDecoder() string {
	var builder strings.Builder

	builder.WriteString("func MessageDecoder(cell *boc.Cell) (string, any, error) {\n")

	builder.WriteString("tag, err := cell.ReadUint(32)\n")
	builder.WriteString(msgDecoderReturnErr)

	builder.WriteString("switch tag {\n")

	for _, k := range utils.GetOrderedKeys(g.loadedTlbMsgTypes) {
		builder.WriteString(fmt.Sprintf("case 0x%x:\n", g.loadedTlbMsgTypes[k].Tag))
		builder.WriteString(fmt.Sprintf("var res %s%s\n", g.loadedTlbMsgTypes[k].TypePrefix, g.loadedTlbMsgTypes[k].TypeName))
		builder.WriteString("err = tlb.Unmarshal(cell, &res)\n")
		builder.WriteString(fmt.Sprintf("return \"%s\", res, err\n", g.loadedTlbMsgTypes[k].TypePrefix))
	}

	builder.WriteString("}\n")
	builder.WriteString("return \"\", nil, fmt.Errorf(\"invalid message tag\")\n")
	builder.WriteString("}\n")
	builder.WriteRune('\n')

	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	return string(b)
}

type templateContext struct {
	Interfaces      map[string]Interface
	InvocationOrder []methodDescription
}

type methodDescription struct {
	Name         string
	InvokeFnName string
	Interfaces   []string
}

func (g *Generator) RenderInvocationOrderList() (string, error) {
	context := templateContext{
		Interfaces: map[string]Interface{},
	}
	descriptions := map[string]methodDescription{}
	for ifaceName, iface := range g.interfaces {
		context.Interfaces[utils.ToCamelCase(ifaceName)] = iface
		for _, method := range iface.Methods {
			if !method.UsedByIntrospection() {
				continue
			}
			invokeFnName := method.GolangFunctionName()
			desc, ok := descriptions[invokeFnName]
			if !ok {
				desc = methodDescription{
					Name:         method.Name,
					InvokeFnName: invokeFnName,
					Interfaces:   []string{},
				}
			}
			desc.Interfaces = append(desc.Interfaces, utils.ToCamelCase(ifaceName))
			descriptions[invokeFnName] = desc
			sort.Strings(desc.Interfaces)
		}
	}

	for _, desc := range descriptions {
		context.InvocationOrder = append(context.InvocationOrder, desc)
	}
	sort.Slice(context.InvocationOrder, func(i, j int) bool {
		return context.InvocationOrder[i].Name < context.InvocationOrder[j].Name
	})
	tmpl, err := template.New("invocationOrder").Parse(invocationOrderTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", err
	}
	return buf.String(), nil
}
