package parser

import (
	"fmt"
	"github.com/startfellows/tongo/tlb"
	tlbParser "github.com/startfellows/tongo/tlb/parser"
	"github.com/startfellows/tongo/utils"
	"go/format"
	"sort"
	"strings"
)

var defaultKnownTypes = map[string]string{
	"accountid":     "tongo.AccountID",
	"cell":          "boc.Cell",
	"int8":          "int8",
	"int64":         "int64",
	"bool":          "bool",
	"uint16":        "uint16",
	"uint64":        "uint64",
	"int256":        "tlb.Int256",
	"int257":        "tlb.Int257",
	"any":           "tlb.Any",
	"[]byte":        "[]byte",
	"big.int":       "big.Int",
	"dnsrecord":     "tlb.DNSRecord",
	"dns_recordset": "tlb.DNSRecordSet",
}

var (
	msgDecoderReturnErr = "if err != nil {return \"\", nil, err}\n"
	returnNilErr        = "if err != nil {return nil, err}\n"
	returnInvalidStack  = "{return \"\", nil, fmt.Errorf(\"invalid stack format\")}\n"
	returnStrNilErr     = "if err != nil {return \"\", nil, err}\n"
)

type TLBMsgBody struct {
	TypePrefix string
	TypeName   string
	Tag        uint64
	Code       string
}

type Generator struct {
	knownTypes        map[string]string
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
		typeName:          typeName,
		loadedTlbMsgTypes: make(map[uint32]TLBMsgBody),
		newTlbTypes:       make(map[string]struct{}),
	}
}

func (g *Generator) GetMethods(interfaces []Interface) (string, error) {
	var builder strings.Builder
	builder.WriteString(`
type Executor interface {
	RunSmcMethod(ctx context.Context, accountID tongo.AccountID, method string, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

	`)
	for _, i := range interfaces {
		for _, m := range i.Methods {
			s, err := g.GetMethod(m)
			if err != nil {
				return "", err
			}
			builder.WriteString(s)
			builder.WriteRune('\n')
		}
	}
	return builder.String(), nil
}

func (g *Generator) RegisterTypes(interfaces []Interface) error {
	for _, i := range interfaces {
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

func (g *Generator) GetMethod(m GetMethod) (string, error) {
	var builder strings.Builder
	var args []string

	for _, o := range m.Output {
		r, err := g.buildResultType(utils.ToCamelCase(o.Version)+utils.ToCamelCase(m.Name)+"Result", o.Stack)
		if err != nil {
			return "", err
		}
		builder.WriteString(r)
		builder.WriteRune('\n')
	}

	builder.WriteString(fmt.Sprintf("func %v(ctx context.Context, executor Executor, reqAccountID tongo.AccountID, ", utils.ToCamelCase(m.Name)))

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

	builder.WriteString(fmt.Sprintf("errCode, stack, err := executor.RunSmcMethod(ctx, reqAccountID, \"%s\", stack)\n", m.Name))
	builder.WriteString(returnStrNilErr)
	builder.WriteString("if errCode != 0 && errCode != 1 {return \"\", nil, fmt.Errorf(\"method execution failed with code: %v\", errCode)}\n")

	if len(m.Output) == 1 {
		name := utils.ToCamelCase(m.Output[0].Version) + utils.ToCamelCase(m.Name) + "Result"
		resDecoder, err := g.buildOutputDecoder(name, m.Output[0].Stack)
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
			name := utils.ToCamelCase(o.Version) + utils.ToCamelCase(m.Name) + "Result"
			builder.WriteString(fmt.Sprintf("decode%s, ", name))
			resDecoder, err := g.buildOutputDecoder(name, o.Stack)
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

	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}
	return string(b), nil
	//return builder.String(), nil
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

func buildOutputStackCheck(r []StackRecord) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("if len(stack) != %d ", len(r)))
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

func (g *Generator) buildOutputDecoder(name string, r []StackRecord) (string, error) {
	var builder, resBuilder strings.Builder

	builder.WriteString(fmt.Sprintf("func decode%s(stack tlb.VmStack) (resultType string, resultAny any, err error) {\n", name))

	builder.WriteString(buildOutputStackCheck(r))

	resBuilder.WriteString(fmt.Sprintf("%s{\n", name))

	for i, s := range r {
		varType, err := g.checkType(s.Type)
		if err != nil {
			return "", err
		}
		varName := utils.ToCamelCasePrivate(s.Name)

		if s.Nullable {
			builder.WriteString(fmt.Sprintf("var %s *%s\n", varName, varType))
			builder.WriteString(fmt.Sprintf("if stack[%d].SumType != \"VmStkNull\" {", i))
		} else {
			builder.WriteString(fmt.Sprintf("var %s %s\n", varName, varType))
		}

		resBuilder.WriteString(fmt.Sprintf("%s: %s,\n", utils.ToCamelCase(s.Name), varName))
		// TODO: add nullable values decoding
		switch s.XMLName.Local {
		case "tinyint":
			if varType == "bool" {
				builder.WriteString(fmt.Sprintf("%s = stack[%d].Int64() != 0\n", varName, i))
			} else {
				builder.WriteString(fmt.Sprintf("%s = %s(stack[%d].Int64())\n", varName, varType, i))
			}
		case "int":
			builder.WriteString(fmt.Sprintf("%s = stack[%d].Int257()\n", varName, i))
		case "slice":
			builder.WriteString(fmt.Sprintf("err = stack[%d].VmStkSlice.UnmarshalToTlbStruct(&%s)\n", i, varName))
			builder.WriteString(returnStrNilErr)
		case "cell":
			builder.WriteString(fmt.Sprintf("%sCell := &stack[%d].VmStkCell.Value\n", varName, i))
			builder.WriteString(fmt.Sprintf("err = tlb.Unmarshal(%sCell, &%s)\n", varName, varName))
			builder.WriteString(returnStrNilErr)
		}

		if s.Nullable {
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

func (g *Generator) sortedMessages() []TLBMsgBody {
	var keys []int
	for k := range g.loadedTlbMsgTypes {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var res []TLBMsgBody
	for _, k := range keys {
		res = append(res, g.loadedTlbMsgTypes[uint32(k)])
	}
	return res
}

func (g *Generator) CollectedTypes() string {
	var builder strings.Builder
	builder.WriteString(strings.Join(g.loadedTlbTypes, "\n\n"))

	for _, m := range g.sortedMessages() {
		builder.WriteString(m.Code)
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

	for _, m := range g.sortedMessages() {
		builder.WriteString(fmt.Sprintf("case 0x%x:\n", m.Tag))
		builder.WriteString(fmt.Sprintf("var res %s%s\n", m.TypePrefix, m.TypeName))
		builder.WriteString("err = tlb.Unmarshal(cell, &res)\n")
		builder.WriteString(fmt.Sprintf("return \"%s\", res, err\n", m.TypePrefix))
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
