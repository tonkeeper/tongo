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
	"accountid":     "ton.AccountID",
	"cell":          "boc.Cell",
	"int8":          "int8",
	"int32":         "int32",
	"int64":         "int64",
	"bool":          "bool",
	"uint8":         "uint8",
	"uint16":        "uint16",
	"uint32":        "uint32",
	"uint64":        "uint64",
	"uint128":       "tlb.Uint128",
	"int256":        "tlb.Int256",
	"int257":        "tlb.Int257",
	"bits256":       "tlb.Bits256",
	"any":           "tlb.Any",
	"[]byte":        "[]byte",
	"string":        "string",
	"coins":         "tlb.Grams",
	"big.int":       "big.Int",
	"dnsrecord":     "tlb.DNSRecord",
	"dns_recordset": "tlb.DNSRecordSet",
	"msgaddress":    "tlb.MsgAddress",
	"text":          "tlb.Text",
	"fullcontent":   "tlb.FullContent",
	"contentdata ":  "tlb.ContentData",
}

var (
	msgDecoderReturnErr = "if err != nil {return \"\", nil, err}\n"
	returnInvalidStack  = "{return \"\", nil, fmt.Errorf(\"invalid stack format\")}\n"
	returnStrNilErr     = "if err != nil {return \"\", nil, err}\n"
	//go:embed invocation_order.tmpl
	invocationOrderTemplate string
	//go:embed get_methods.tmpl
	getMethodsTemplate string
	//go:embed messages.tmpl
	messagesTemplate string
)

type TLBMsgBody struct {
	TypeName      string
	OperationName string
	Tag           uint64
	Code          string
	FixedLength   bool
}

type Generator struct {
	knownTypes            map[string]string
	abi                   ABI
	newTlbTypes           map[string]struct{}
	loadedTlbTypes        []string
	loadedTlbMsgTypes     map[uint32][]TLBMsgBody
	loadedJettonsMsgTypes map[uint32][]TLBMsgBody
	loadedNFTsMsgTypes    map[uint32][]TLBMsgBody
	typeName              string
}

func NewGenerator(knownTypes map[string]string, abi ABI) (*Generator, error) {
	if knownTypes == nil {
		knownTypes = defaultKnownTypes
	}
	g := &Generator{
		knownTypes:            knownTypes,
		abi:                   abi,
		loadedTlbMsgTypes:     make(map[uint32][]TLBMsgBody),
		loadedJettonsMsgTypes: make(map[uint32][]TLBMsgBody),
		loadedNFTsMsgTypes:    make(map[uint32][]TLBMsgBody),
		newTlbTypes:           make(map[string]struct{}),
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
		err := registerMstType(g.loadedTlbMsgTypes, "MsgBody", internal.Name, internal.Input, internal.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, jetton := range g.abi.JettonPayloads {
		err := registerMstType(g.loadedJettonsMsgTypes, "JettonPayload", jetton.Name, jetton.Input, jetton.FixedLength)
		if err != nil {
			return err
		}
	}
	for _, nft := range g.abi.NFTPayloads {
		err := registerMstType(g.loadedNFTsMsgTypes, "NFTPayload", nft.Name, nft.Input, nft.FixedLength)
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

	gen := tlbParser.NewGenerator(nil, "")
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

func registerMstType(known map[uint32][]TLBMsgBody, postfix, name, s string, fixedLength bool) error {
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

	typePrefix := utils.ToCamelCase(name) + postfix

	t, err := gen.GenerateGolangTypes(parsed.Declarations, typePrefix, true)
	if err != nil {
		return fmt.Errorf("can't decode %v error %w", s, err)
	}

	known[uint32(tag.Val)] = append(known[uint32(tag.Val)], TLBMsgBody{
		TypeName:      utils.ToCamelCase(name) + postfix,
		OperationName: utils.ToCamelCase(name),
		Tag:           tag.Val,
		Code:          t,
		FixedLength:   fixedLength,
	})

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

	for _, k := range utils.GetOrderedKeys(g.loadedTlbMsgTypes) {
		for _, v := range g.loadedTlbMsgTypes[k] {
			builder.WriteString(v.Code)
			builder.WriteRune('\n')
		}
	}
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
	Operations map[uint32][]TLBMsgBody
}

func (g *Generator) GenerateMsgDecoder() string {

	context := messagesContext{g.loadedTlbMsgTypes}

	tmpl, err := template.New("mesages").Parse(messagesTemplate)
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

type templateContext struct {
	Interfaces      map[string]string
	InvocationOrder []methodDescription
	InterfaceOrder  []interfacDescripion
	KnownHashes     map[string]interfacDescripion
	Inheritance     map[string]string
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
	context := templateContext{
		Interfaces:  map[string]string{},
		KnownHashes: map[string]interfacDescripion{},
		Inheritance: map[string]string{},
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
	return g.renderPayload("Jetton", g.loadedJettonsMsgTypes)
}

func (g *Generator) RenderNFT() (string, error) {
	return g.renderPayload("NFT", g.loadedNFTsMsgTypes)
}

func (g *Generator) renderPayload(payloadName string, msgTypes map[uint32][]TLBMsgBody) (string, error) {
	var builder strings.Builder

	fmt.Fprintf(&builder, `

func (j *%sPayload) UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder) error  {
	if cell.BitsAvailableForRead() == 0 && cell.RefsAvailableForRead() == 0 {
		return nil
	}
	tempCell := cell.CopyRemaining()
	op64, err := tempCell.ReadUint(32)
	if errors.Is(err, boc.ErrNotEnoughBits) {
		j.SumType = Unknown%sOp
		j.Value = cell.CopyRemaining()
		return nil
	}
	op := uint32(op64)
	j.OpCode = &op
	switch  op {
`, payloadName, payloadName)

	var knownTypes [][2]string
	var knownOpcodes []opCode
	for _, k := range utils.GetOrderedKeys(msgTypes) {
		tlbType := msgTypes[k][0] //todo: iterate
		fmt.Fprintf(&builder, "case %s%sOpCode:  // 0x%08x\n", tlbType.OperationName, payloadName, tlbType.Tag)
		fmt.Fprintf(&builder, "var res %s\n", tlbType.TypeName)
		fmt.Fprintf(&builder, `if err := tlb.Unmarshal(tempCell, &res); err == nil { 
	j.SumType = %v%sOp
	j.Value = res
	return  nil 
}`, tlbType.OperationName, payloadName)
		builder.WriteString("\n")
		if tlbType.FixedLength {
			builder.WriteString(fmt.Sprintf(` panic! if cell.RefsAvailableForRead() > 0 || cell.BitsAvailableForRead() > 0 { return %sMsgOp, nil, ErrStructSizeMismatch }`, tlbType.OperationName))
			builder.WriteString("\n")
		}
		knownTypes = append(knownTypes, [2]string{tlbType.OperationName, tlbType.TypeName})
		knownOpcodes = append(knownOpcodes, opCode{OperationName: tlbType.OperationName, Tag: tlbType.Tag})
	}

	fmt.Fprintf(&builder, `
}
		j.SumType = Unknown%sOp
		j.Value = cell.CopyRemaining()
	
	return nil
}
const (
`, payloadName)
	for _, v := range knownTypes {
		fmt.Fprintf(&builder, `%v%sOp %sOpName = "%v"`+"\n", v[0], payloadName, payloadName, v[0])
	}
	builder.WriteString(`

`)
	for _, v := range knownOpcodes {
		fmt.Fprintf(&builder, `%s%sOpCode %sOpCode = 0x%08x`+"\n", v.OperationName, payloadName, payloadName, v.Tag)
	}
	builder.WriteString(")\n")
	fmt.Fprintf(&builder, "var Known%sTypes = map[string]any{\n", payloadName)
	for _, v := range knownTypes {
		fmt.Fprintf(&builder, "%v%sOp: %v{},\n", v[0], payloadName, v[1])
	}
	fmt.Fprintf(&builder, `}
var %sOpCodes = map[%sOpName]%sOpCode{
`, strings.ToLower(payloadName), payloadName, payloadName)
	for _, v := range knownOpcodes {
		fmt.Fprintf(&builder, "%s%sOp: %s%sOpCode,\n", v.OperationName, payloadName, v.OperationName, payloadName)
	}
	builder.WriteString("\n}\n")
	for _, k := range utils.GetOrderedKeys(msgTypes) {
		builder.WriteString(msgTypes[k][0].Code) //todo: iterate
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')
	return builder.String(), nil
}
