package parser

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/tonkeeper/tongo/tlb"
	tlbParser "github.com/tonkeeper/tongo/tlb/parser"
	"github.com/tonkeeper/tongo/utils"
)

var defaultKnownTypes = map[string]tlbParser.DefaultType{
	"cell":          {"boc.Cell", false},
	"int8":          {"int8", false},
	"int32":         {"int32", false},
	"int64":         {"int64", false},
	"bool":          {"bool", false},
	"uint8":         {"uint8", false},
	"uint16":        {"uint16", false},
	"uint32":        {"uint32", false},
	"uint64":        {"uint64", false},
	"uint128":       {"tlb.Uint128", false},
	"int256":        {"tlb.Int256", false},
	"int257":        {"tlb.Int257", false},
	"bits256":       {"tlb.Bits256", false},
	"any":           {"tlb.Any", false},
	"[]byte":        {"[]byte", false},
	"string":        {"string", false},
	"coins":         {"tlb.Grams", false},
	"big.int":       {"big.Int", false},
	"dnsrecord":     {"tlb.DNSRecord", false},
	"dns_recordset": {"tlb.DNSRecordSet", false},
	"msgaddress":    {"tlb.MsgAddress", false},
	"text":          {"tlb.Text", false},
	"fullcontent":   {"tlb.FullContent", false},
	"contentdata ":  {"tlb.ContentData", false},
	"StateInit":     {"tlb.StateInit", false},
}

var (
	msgDecoderReturnErr = "if err != nil {return \"\", nil, err}\n"
	returnInvalidStack  = "{return \"\", nil, fmt.Errorf(\"invalid stack format\")}\n"
	returnStrNilErr     = "if err != nil {return \"\", nil, err}\n"
	//go:embed messages.md.tmpl
	messagesMDTemplate string
	//go:embed interfaces.tmpl
	invocationOrderTemplate string
	//go:embed get_methods.tmpl
	getMethodsTemplate string
	//go:embed messages.tmpl
	messagesTemplate string
	//go:embed payloads.tmpl
	payloadTmpl string
)

type MsgType int

const (
	MsgTypeIn            MsgType = 0
	MsgTypeExtIn         MsgType = 1
	MsgTypeExtOut        MsgType = 2
	MsgTypeJettonPayload MsgType = 3
	MsgTypeNFTPayload    MsgType = 4
)

type TLBMsgBody struct {
	Type             MsgType
	GolangTypeName   string
	GolangOpcodeName string
	OperationName    string
	Tag              uint64
	Code             string
	FixedLength      bool
}

type Generator struct {
	knownTypes        map[string]tlbParser.DefaultType
	abi               ABI
	newTlbTypes       map[string]struct{}
	loadedTlbTypes    []string
	loadedTlbMsgTypes map[tlb.Tag][]TLBMsgBody
	typeName          string
}

func NewGenerator(knownTypes map[string]tlbParser.DefaultType, abi ABI) (*Generator, error) {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	g := &Generator{
		knownTypes:        knownTypes,
		abi:               abi,
		loadedTlbMsgTypes: make(map[tlb.Tag][]TLBMsgBody),
		newTlbTypes:       make(map[string]struct{}),
	}
	err := g.registerABI()
	return g, err
}

type getMethodContext struct {
	GetMethods    []getMethodDesc
	SimpleMethods map[int][]string
}
type getMethodDesc struct {
	Name        string
	MethodName  string
	ID          int
	Body        string
	Decoders    []string
	ResultTypes map[string]string
}

func (g *Generator) GetMethods() (string, []string, error) {
	context := getMethodContext{
		SimpleMethods: map[int][]string{},
		GetMethods:    []getMethodDesc{},
	}

	usedNames := map[string]struct{}{}
	simpleMethods := []string{}
	for _, m := range g.abi.Methods {
		methodName := m.GolangFunctionName()
		var err error
		desc := getMethodDesc{
			Name:        m.Name,
			MethodName:  methodName,
			ResultTypes: make(map[string]string),
		}
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
			simpleMethods = append(simpleMethods, m.Name)
			context.SimpleMethods[methodID] = append(context.SimpleMethods[methodID], methodName)
		}
		for _, o := range m.Output {
			resultTypeName := o.FullResultName(methodName)
			desc.Decoders = append(desc.Decoders, "Decode"+resultTypeName)
			r, err := g.buildStackStruct(o.Stack)
			if err != nil {
				return "", nil, err
			}
			desc.ResultTypes[resultTypeName] = r
		}
		desc.Body, err = g.getMethod(m, methodID, methodName)
		if err != nil {
			return "", nil, err
		}

		context.GetMethods = append(context.GetMethods, desc)
	}

	tmpl, err := template.New("getMethods").Parse(getMethodsTemplate)
	if err != nil {
		return "", nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", nil, err
	}

	b, err := format.Source([]byte(buf.String()))
	if err != nil {
		return "", nil, err
	}
	return string(b), simpleMethods, nil
}

func (g *Generator) registerABI() error {
	for _, t := range g.abi.Types {
		if t != "" {
			err := g.registerType(t)
			if err != nil {
				return err
			}
		}
	}
	for _, internal := range g.abi.Internals {
		err := registerMsgType(g.loadedTlbMsgTypes, MsgTypeIn, internal.Name, internal.Input, internal.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, m := range g.abi.ExtIn {
		err := registerMsgType(g.loadedTlbMsgTypes, MsgTypeExtIn, m.Name, m.Input, m.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, m := range g.abi.ExtOut {
		err := registerMsgType(g.loadedTlbMsgTypes, MsgTypeExtOut, m.Name, m.Input, m.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, jetton := range g.abi.JettonPayloads {
		err := registerMsgType(g.loadedTlbMsgTypes, MsgTypeJettonPayload, jetton.Name, jetton.Input, jetton.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, nft := range g.abi.NFTPayloads {
		err := registerMsgType(g.loadedTlbMsgTypes, MsgTypeNFTPayload, nft.Name, nft.Input, nft.FixedLength)
		if err != nil {
			return err
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

	gen := tlbParser.NewGenerator(tlbParser.WithDefaultTypes(defaultKnownTypes, false))
	_, err = gen.GenerateGolangTypes(tlbData.Declarations, "", false)
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

func registerMsgType(known map[tlb.Tag][]TLBMsgBody, mType MsgType, name, s string, fixedLength bool) error {
	parsed, err := tlbParser.Parse(s)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}
	if len(parsed.Declarations) != 1 {
		return fmt.Errorf("must be only one declaration for MsgBody %v", s)
	}

	gen := tlbParser.NewGenerator(tlbParser.WithDefaultTypes(defaultKnownTypes, false))

	tag, err := tlb.ParseTag(parsed.Declarations[0].Constructor.Prefix)
	if err != nil {
		return fmt.Errorf("can't decode tag error %w", err)
	}
	key := tlb.Tag{
		Len: tag.Len,
		Val: tag.Val,
	}
	var typeSuffix string
	var opSuffix string
	switch mType {
	case MsgTypeIn:
		opSuffix = "MsgOp"
		typeSuffix = "MsgBody"
	case MsgTypeExtIn:
		opSuffix = "ExtInMsgOp"
		typeSuffix = "ExtInMsgBody"
	case MsgTypeExtOut:
		opSuffix = "ExtOutMsgOp"
		typeSuffix = "ExtOutMsgBody"
	case MsgTypeJettonPayload:
		opSuffix = "JettonOp"
		typeSuffix = "JettonPayload"
	case MsgTypeNFTPayload:
		opSuffix = "NFTOp"
		typeSuffix = "NFTPayload"
	}
	typePrefix := utils.ToCamelCase(name)

	t, err := gen.GenerateGolangTypes(parsed.Declarations, typePrefix+typeSuffix, true)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}

	known[key] = append(known[key], TLBMsgBody{
		Type:             mType,
		GolangTypeName:   typePrefix + typeSuffix,
		GolangOpcodeName: utils.ToCamelCase(name) + opSuffix,
		OperationName:    utils.ToCamelCase(name),
		Tag:              tag.Val,
		Code:             t,
		FixedLength:      fixedLength,
	})

	return nil
}

func (g *Generator) checkType(s string) (string, error) {
	if typeName, prs := g.knownTypes[strings.ToLower(s)]; prs {
		return typeName.Name, nil
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

	builder.WriteString(fmt.Sprintf("func %v(ctx context.Context, executor Executor, reqAccountID ton.AccountID, ", m.GolangFunctionName()))

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

	builder.WriteString(fmt.Sprintf("// MethodID = %d for \"%s\" method\n", methodID, m.Name))
	builder.WriteString(fmt.Sprintf("errCode, stack, err := executor.RunSmcMethodByID(ctx, reqAccountID, %d, stack)\n", methodID))
	builder.WriteString(returnStrNilErr)
	builder.WriteString("if errCode != 0 && errCode != 1 {return \"\", nil, fmt.Errorf(\"method execution failed with code: %v\", errCode)}\n")

	decoders := ""
	builder.WriteString("for _, f := range []func(tlb.VmStack)(string, any, error){")
	for _, o := range m.Output {
		name := o.FullResultName(methodName)
		builder.WriteString(fmt.Sprintf("Decode%s, ", name))
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
		if s.Nullable || (s.XMLName.Local == "tuple" && s.List) {
			nullableCheck = fmt.Sprintf(" && %s != \"VmStkNull\"", stackType)
		}
		switch s.XMLName.Local {
		case "int":
			builder.WriteString(fmt.Sprintf("|| ((%s != \"VmStkTinyInt\" && %s != \"VmStkInt\")%s) ", stackType, stackType, nullableCheck))
		case "slice":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkSlice\"%s) ", stackType, nullableCheck))
		case "cell":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkCell\"%s) ", stackType, nullableCheck))
		case "tuple":
			builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkTuple\"%s) ", stackType, nullableCheck))
		}
	}
	builder.WriteString(returnInvalidStack)
	return builder.String()
}

func (g *Generator) buildOutputDecoder(name string, r []StackRecord, isFixed bool) (string, error) {
	builder := new(strings.Builder)

	builder.WriteString(fmt.Sprintf("func Decode%s(stack tlb.VmStack) (resultType string, resultAny any, err error) {\n", name))

	builder.WriteString(buildOutputStackCheck(r, isFixed))
	builder.WriteString(fmt.Sprintf("var result %v\n", name))
	builder.WriteString("err = stack.Unmarshal(&result)\n")

	builder.WriteString(fmt.Sprintf("return \"%s\",result, err", name))

	builder.WriteString("}\n")

	return builder.String(), nil
}

func (g *Generator) buildResultType(name string, s []StackRecord) (string, error) {
	str, err := g.buildStackStruct(s)
	return "type " + name + " " + str, err
}

func (g *Generator) buildStackStruct(s []StackRecord) (string, error) {
	var builder strings.Builder
	builder.WriteString("struct {\n")
	for _, c := range s {
		if c.XMLName.Local == "tuple" {
			list := ""
			pointer := ""
			if c.List {
				list = "[]"
			} else if c.Nullable {
				pointer = "*"
			}
			str, err := g.buildStackStruct(c.SubStack)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&builder, "%s %s%s%s\n", utils.ToCamelCase(c.Name), pointer, list, str)
			continue
		}
		t, err := g.checkType(c.Type)
		if err != nil {
			return "", err
		}
		pointer := ""
		if c.Nullable {
			pointer = "*"
		}
		fmt.Fprintf(&builder, "%s %s%s\n", utils.ToCamelCase(c.Name), pointer, t)

	}
	builder.WriteString("}\n")
	return builder.String(), nil

}

func (g *Generator) CollectedTypes() string {
	var builder strings.Builder
	builder.WriteString(strings.Join(g.loadedTlbTypes, "\n\n"))

	builder.WriteRune('\n')
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	return string(b)
}

type opCode struct {
	OperationName string
	Tag           uint64
}

type messagesContext struct {
	Operations map[tlb.Tag][]TLBMsgBody
	WhatRender string
}

func (g *Generator) GenerateMsgDecoder() string {
	s := g.generateMsgDecoder(MsgTypeIn, "MsgIn")
	s += g.generateMsgDecoder(MsgTypeExtIn, "MsgExtIn")
	s += g.generateMsgDecoder(MsgTypeExtOut, "MsgExtOut")
	return s
}
func (g *Generator) generateMsgDecoder(msgType MsgType, what string) string {

	context := messagesContext{Operations: make(map[tlb.Tag][]TLBMsgBody), WhatRender: what}
	for tag, operation := range g.loadedTlbMsgTypes {
		filtered := make([]TLBMsgBody, 0, len(operation))
		for _, body := range operation {
			if body.Type == msgType {
				filtered = append(filtered, body)
			}
		}
		if len(filtered) > 0 {
			context.Operations[tag] = filtered
		}
	}
	tmpl, err := template.New("messages").Parse(messagesTemplate)
	if err != nil {
		panic(err)
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		panic(err)
		return ""
	}
	return buf.String()
}

type messageOperation struct {
	Name      string
	OpCode    string
	OpCodeLen int
}

type messagesMDContext struct {
	Operations []messageOperation
}

// RenderMessagesMD renders messages.md file with messages and their names + opcodes.
func (g *Generator) RenderMessagesMD() (string, error) {
	context := messagesMDContext{}
	for opcode, bodies := range g.loadedTlbMsgTypes {
		for _, body := range bodies {
			operation := messageOperation{
				Name:      body.OperationName,
				OpCode:    fmt.Sprintf("0x%08x", opcode.Val),
				OpCodeLen: opcode.Len,
			}
			context.Operations = append(context.Operations, operation)
		}
	}
	sort.Slice(context.Operations, func(i, j int) bool {
		if context.Operations[i].Name == context.Operations[j].Name {
			return context.Operations[i].OpCode < context.Operations[j].OpCode
		}
		return context.Operations[i].Name < context.Operations[j].Name
	})
	tmpl, err := template.New("messagesMD").Parse(messagesMDTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type methodDescription struct {
	Name         string
	InvokeFnName string
}

type interfacDescripion struct {
	Name       string
	Results    []string
	GetMethods []string
}

func (g *Generator) RenderInvocationOrderList(simpleMethods []string) (string, error) {
	context := struct {
		Interfaces                     map[string]string
		InvocationOrder                []methodDescription
		InterfaceOrder                 []interfacDescripion
		KnownHashes                    map[string]interfacDescripion
		Inheritance                    map[string]string
		IntMsgs, ExtInMsgs, ExtOutMsgs map[string][]string
	}{
		Interfaces:  map[string]string{},
		KnownHashes: map[string]interfacDescripion{},
		Inheritance: map[string]string{},
		IntMsgs:     map[string][]string{},
		ExtInMsgs:   map[string][]string{},
		ExtOutMsgs:  map[string][]string{},
	}
	descriptions := map[string]methodDescription{}

	for _, method := range g.abi.Methods {
		if !method.UsedByIntrospection() {
			continue
		}

		invokeFnName := method.GolangFunctionName()
		desc, ok := descriptions[invokeFnName]
		if ok {
			return "", fmt.Errorf("method duplicate %v", invokeFnName)
		}

		desc = methodDescription{
			Name:         method.Name,
			InvokeFnName: invokeFnName,
		}
		descriptions[invokeFnName] = desc
	}
	for _, iface := range g.abi.Interfaces {
		ifaceName := utils.ToCamelCase(iface.Name)
		context.Interfaces[ifaceName] = iface.Name
		descripion := interfacDescripion{
			Name: ifaceName,
		}
		if iface.Inherits != "" {
			context.Inheritance[ifaceName] = utils.ToCamelCase(iface.Inherits)
		}
		for _, method := range iface.Methods {
			if !slices.Contains(simpleMethods, method.Name) {
				continue
			}
			resultName := utils.ToCamelCase(method.Name) + "Result"
			descripion.GetMethods = append(descripion.GetMethods, utils.ToCamelCase(method.Name))
			if method.Version != "" {
				resultName = fmt.Sprintf("%s_%sResult", utils.ToCamelCase(method.Name), utils.ToCamelCase(method.Version))
			}
			descripion.Results = append(descripion.Results, resultName)
		}
		for _, m := range iface.Input.Internals {
			context.IntMsgs[ifaceName] = append(context.IntMsgs[ifaceName], utils.ToCamelCase(m.Name))
		}
		for _, m := range iface.Input.Externals {
			context.ExtInMsgs[ifaceName] = append(context.ExtInMsgs[ifaceName], utils.ToCamelCase(m.Name))
		}
		for _, m := range iface.Output.Externals {
			context.ExtOutMsgs[ifaceName] = append(context.ExtOutMsgs[ifaceName], utils.ToCamelCase(m.Name))
		}
		if len(iface.CodeHashes) > 0 { //we dont' need to detect interfaces with code hashes because we can them directly
			for _, hash := range iface.CodeHashes {

				context.KnownHashes[hash] = descripion
			}
		} else {
			context.InterfaceOrder = append(context.InterfaceOrder, descripion)
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

func (g *Generator) RenderJetton() (string, error) {
	return g.renderPayload(MsgTypeJettonPayload)
}

func (g *Generator) RenderNFT() (string, error) {
	return g.renderPayload(MsgTypeNFTPayload)
}

func (g *Generator) renderPayload(mType MsgType) (string, error) {

	context := messagesContext{Operations: make(map[tlb.Tag][]TLBMsgBody)}
	switch mType {
	case MsgTypeNFTPayload:
		context.WhatRender = "NFT"
	case MsgTypeJettonPayload:
		context.WhatRender = "Jetton"
	}
	for tag, operation := range g.loadedTlbMsgTypes {
		filtered := make([]TLBMsgBody, 0, len(operation))
		for _, body := range operation {
			if body.Type == mType {
				filtered = append(filtered, body)
			}
		}
		if len(filtered) > 0 {
			context.Operations[tag] = filtered
		}
	}
	tmpl, err := template.New("jettons").Parse(payloadTmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getOrderedKeys[M ~map[tlb.Tag]V, V any](m M) []tlb.Tag {
	keys := maps.Keys(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Val < keys[j].Val
	})
	return keys
}
