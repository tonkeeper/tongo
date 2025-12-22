package parser

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/tonkeeper/tongo/tlb"
	tolkAbi "github.com/tonkeeper/tongo/tolk/abi"
	tolkParser "github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
	"golang.org/x/exp/maps"
)

var (
	returnInvalidStack = "{return \"\", nil, fmt.Errorf(\"invalid stack format\")}\n"
	returnStrNilErr    = "if err != nil {return \"\", nil, err}\n"
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
	//go:embed errors.tmpl
	contractErrorsTmpl string
)

const (
	MsgTypeInt    MsgType = 0
	MsgTypeExtIn  MsgType = 1
	MsgTypeExtOut MsgType = 2
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
	structRefs        map[string]tolkAbi.StructDeclaration
	aliasRefs         map[string]tolkAbi.AliasDeclaration
	enumRefs          map[string]tolkAbi.EnumDeclaration
	abi               []tolkAbi.ABI
	newTlbTypes       map[string]struct{}
	loadedTlbTypes    []string
	loadedTlbMsgTypes map[tlb.Tag][]TLBMsgBody
}

type MsgType int

func NewGenerator(abi []tolkAbi.ABI) *Generator {
	g := &Generator{
		structRefs:        make(map[string]tolkAbi.StructDeclaration),
		aliasRefs:         make(map[string]tolkAbi.AliasDeclaration),
		enumRefs:          make(map[string]tolkAbi.EnumDeclaration),
		loadedTlbMsgTypes: make(map[tlb.Tag][]TLBMsgBody),
		newTlbTypes:       make(map[string]struct{}),
		abi:               abi,
	}
	err := g.registerABI()
	if err != nil {
		return g
	}

	return g
}

func (g *Generator) registerABI() error {
	if err := g.mapRefs(); err != nil {
		return err
	}

	msgsName := make(map[string]struct{})
	for _, abi := range g.abi {
		for _, declr := range abi.Declarations {
			err := g.registerType(declr)
			if err != nil {
				return err
			}
		}

		typePrefix := abi.ContractName
		if abi.IsGeneral {
			typePrefix = ""
		}
		for _, intMsg := range abi.IncomingMessages {
			err := g.registerMsgType(MsgTypeInt, intMsg.BodyTy, typePrefix, msgsName)
			if err != nil {
				return err
			}
		}
		for _, intMsg := range abi.OutgoingMessages {
			err := g.registerMsgType(MsgTypeInt, intMsg.BodyTy, typePrefix, msgsName)
			if err != nil {
				return err
			}
		}
		for _, extOutMsg := range abi.EmittedMessages {
			err := g.registerMsgType(MsgTypeExtOut, extOutMsg.BodyTy, typePrefix, msgsName)
			if err != nil {
				return err
			}
		}
		if abi.IncomingExternal != nil {
			err := g.registerMsgType(MsgTypeExtIn, abi.IncomingExternal.BodyTy, typePrefix, msgsName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) registerType(declr tolkAbi.Declaration) error {
	var result *tolkParser.DeclrResult
	var err error
	switch declr.SumType {
	case "Struct":
		result, err = tolkParser.ParseStructDeclr(declr.StructDeclaration)
		if err != nil {
			return err
		}
	case "Alias":
		result, err = tolkParser.ParseAliasDeclr(declr.AliasDeclaration)
		if err != nil {
			return err
		}
	case "Enum":
		result, err = tolkParser.ParseEnumDeclr(declr.EnumDeclaration)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown declaration type")
	}

	if _, ok := g.newTlbTypes[result.Name]; ok {
		return nil
	}

	g.newTlbTypes[result.Name] = struct{}{}
	g.loadedTlbTypes = append(g.loadedTlbTypes, result.Code)

	return nil
}

func (g *Generator) registerMsgType(mType MsgType, ty tolkAbi.Ty, abiName string, msgsName map[string]struct{}) error {
	golangAbiName := utils.ToCamelCase(abiName)

	tag, err := tolkParser.ParseTag(ty, g.structRefs, g.aliasRefs)
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
	case MsgTypeInt:
		opSuffix = "MsgOp"
		typeSuffix = "MsgBody"
	case MsgTypeExtIn:
		opSuffix = "ExtInMsgOp"
		typeSuffix = "ExtInMsgBody"
	case MsgTypeExtOut:
		opSuffix = "ExtOutMsgOp"
		typeSuffix = "ExtOutMsgBody"
	}

	var typePrefix string
	var msgName string
	var res *tolkParser.MsgResult
	switch ty.SumType {
	case "StructRef":
		typePrefix = utils.ToCamelCase(ty.StructRefTy.StructName)
		msgName = golangAbiName + typePrefix + typeSuffix
		res, err = tolkParser.ParseStructMsg(ty, msgName)
		if err != nil {
			return err
		}
	case "AliasRef":
		typePrefix = utils.ToCamelCase(ty.AliasRefTy.AliasName)
		msgName = golangAbiName + typePrefix + typeSuffix
		res, err = tolkParser.ParseAliasMsg(ty, msgName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("message type must be either struct or alias, got %v", ty.SumType)
	}

	if _, ok := msgsName[msgName]; ok {
		return nil
	}
	msgsName[msgName] = struct{}{}
	g.loadedTlbMsgTypes[key] = append(g.loadedTlbMsgTypes[key], TLBMsgBody{
		Type:             mType,
		GolangTypeName:   msgName,
		GolangOpcodeName: golangAbiName + typePrefix + opSuffix,
		OperationName:    golangAbiName + typePrefix,
		Tag:              tag.Val,
		Code:             res.Code,
	})

	return nil
}

func (g *Generator) mapRefs() error {
	for _, abi := range g.abi {
		for _, declr := range abi.Declarations {
			switch declr.SumType {
			case "Struct":
				g.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
			case "Alias":
				g.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
			case "Enum":
				g.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
			default:
				return fmt.Errorf("unknown declaration type")
			}
		}
	}
	return nil
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

type messagesContext struct {
	Operations map[tlb.Tag][]TLBMsgBody
	WhatRender string
}

func (g *Generator) GenerateMsgDecoder() string {
	s := g.generateMsgDecoder(MsgTypeInt, "MsgIn")
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
	var simpleMethods []string
	for _, abi := range g.abi {
		for _, m := range abi.GetMethods {
			methodName := m.GolangFunctionName()
			var err error
			desc := getMethodDesc{
				Name:        m.Name,
				MethodName:  methodName,
				ResultTypes: make(map[string]string),
			}
			var methodID int
			if m.TvmMethodID != 0 {
				methodID = m.TvmMethodID
			} else {
				methodID = utils.MethodIdFromName(m.Name)
			}
			if _, ok := usedNames[methodName]; ok {
				continue
			}
			usedNames[methodName] = struct{}{}

			if len(m.Parameters) == 0 {
				simpleMethods = append(simpleMethods, m.Name)
				context.SimpleMethods[methodID] = append(context.SimpleMethods[methodID], methodName)
			}

			contractName := abi.ContractName
			if abi.IsGeneral {
				contractName = ""
			}
			resultTypeName := m.FullResultName(contractName)
			desc.Decoders = append(desc.Decoders, "Decode"+resultTypeName)
			r, err := tolkParser.ParseGetMethodCode(m.ReturnTy)
			if err != nil {
				return "", nil, err
			}
			desc.ResultTypes[resultTypeName] = r
			desc.Body, err = g.getMethod(m, methodID, contractName)
			if err != nil {
				return "", nil, err
			}

			context.GetMethods = append(context.GetMethods, desc)
		}
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

func (g *Generator) getMethod(m tolkAbi.GetMethod, methodID int, contractName string) (string, error) {
	var builder strings.Builder
	var args []string

	builder.WriteString(fmt.Sprintf("func %v(ctx context.Context, executor Executor, reqAccountID ton.AccountID, ", m.GolangFunctionName()))

	for _, p := range m.Parameters {
		t, err := tolkParser.ParseType(p.Ty)
		if err != nil {
			return "", err
		}
		args = append(args, fmt.Sprintf("%v %v", utils.ToCamelCasePrivate(p.Name), t))
	}
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(") (string, any, error) {\n")

	builder.WriteString(buildInputStackValues(m.Parameters))
	builder.WriteRune('\n')

	builder.WriteString(fmt.Sprintf("// MethodID = %d for \"%s\" method\n", methodID, m.Name))
	builder.WriteString(fmt.Sprintf("errCode, stack, err := executor.RunSmcMethodByID(ctx, reqAccountID, %d, stack)\n", methodID))
	builder.WriteString(returnStrNilErr)
	builder.WriteString("if errCode != 0 && errCode != 1 {return \"\", nil, fmt.Errorf(\"method execution failed with code: %v\", errCode)}\n")

	decoders := ""
	builder.WriteString("for _, f := range []func(tlb.VmStack)(string, any, error){")

	name := m.FullResultName(contractName)
	builder.WriteString(fmt.Sprintf("Decode%s, ", name))
	resDecoder, err := g.buildOutputDecoder(name, m.ReturnTy)
	if err != nil {
		return "", err
	}
	decoders = decoders + "\n\n" + resDecoder

	builder.WriteString("} {\n")
	builder.WriteString("s, r, err := f(stack)\n")
	builder.WriteString("if err == nil {return s, r, nil}\n")
	builder.WriteString("}\n")
	builder.WriteString("return \"\", nil, fmt.Errorf(\"can not decode outputs\")\n}\n")
	builder.WriteString(decoders)
	builder.WriteRune('\n')

	return builder.String(), nil
}

func buildInputStackValues(p []tolkAbi.Parameter) string {
	var builder strings.Builder
	builder.WriteString("stack := tlb.VmStack{}\n")

	if len(p) > 0 {
		builder.WriteString("var (val tlb.VmStackValue\n err error)\n")
	}

	for _, s := range p {
		switch s.Ty.SumType {
		case "intN":
			if s.Ty.NumberTy.N <= 64 {
				builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkTinyInt\", VmStkTinyInt: int64(%s)}\n",
					utils.ToCamelCasePrivate(s.Name)))
			} else {
				builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkInt\", VmStkInt: %s}\n",
					utils.ToCamelCasePrivate(s.Name)))
			}
		case "uintN":
			if s.Ty.NumberTy.N <= 63 {
				builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkTinyInt\", VmStkTinyInt: int64(%s)}\n",
					utils.ToCamelCasePrivate(s.Name)))
			} else {
				builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkInt\", VmStkInt: %s}\n",
					utils.ToCamelCasePrivate(s.Name)))
			}
		case "bool":
			builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkTinyInt\", VmStkTinyInt: BoolToInt64(%s)}\n",
				utils.ToCamelCasePrivate(s.Name)))
		case "coins":
			builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkInt\", VmStkInt: tlb.Int257(BigIntFromUint(uint64(%s)))}\n",
				utils.ToCamelCasePrivate(s.Name)))
		case "varintN", "varuintN", "int":
			builder.WriteString(fmt.Sprintf("val = tlb.VmStackValue{SumType:  \"VmStkInt\", VmStkInt: %s}\n",
				utils.ToCamelCasePrivate(s.Name)))
		case "address", "addressExt", "addressOpt", "addressAny", "slice", "remaining", "bitsN":
			builder.WriteString(fmt.Sprintf("val, err = tlb.TlbStructToVmCellSlice(%s)\n",
				utils.ToCamelCasePrivate(s.Name)))
			builder.WriteString(returnStrNilErr)
		case "cell", "builder", "cellOf", "mapKV":
			builder.WriteString(fmt.Sprintf("val, err = tlb.TlbStructToVmCell(%s)\n",
				utils.ToCamelCasePrivate(s.Name)))
			builder.WriteString(returnStrNilErr)
		case "nullLiteral", "void":
			builder.WriteString("val = tlb.VmStackValue{SumType:  \"VmStkNull\"}\n")
		default:
			builder.WriteString(fmt.Sprintf("val, err = tlb.TlbStructToVmCell(%s)\n",
				utils.ToCamelCasePrivate(s.Name)))
			builder.WriteString(returnStrNilErr)
		}
		builder.WriteString("stack.Put(val)\n")
	}

	return builder.String()
}

func (g *Generator) buildOutputDecoder(name string, ty tolkAbi.Ty) (string, error) {
	builder := new(strings.Builder)

	builder.WriteString(fmt.Sprintf("func Decode%s(stack tlb.VmStack) (resultType string, resultAny any, err error) {\n", name))

	output, err := g.buildOutputStackCheck(ty)
	if err != nil {
		return "", err
	}
	builder.WriteString(output)
	builder.WriteString(fmt.Sprintf("var result %v\n", name))
	builder.WriteString("err = stack.Unmarshal(&result)\n")

	builder.WriteString(fmt.Sprintf("return \"%s\",result, err", name))

	builder.WriteString("}\n")

	return builder.String(), nil
}

func (g *Generator) buildOutputStackCheck(ty tolkAbi.Ty) (string, error) {
	var builder strings.Builder
	// todo: what to do with is fixed? is it really needed? if so, maybe add some annotation to tolk

	var checksBuilder strings.Builder
	res, err := g.buildOutputStackTy(ty, &checksBuilder, 0, false, make(map[string]tolkAbi.Ty))
	if err != nil {
		return "", err
	}

	builder.WriteString(fmt.Sprintf("if len(stack) < %d ", res+1))
	builder.WriteString(checksBuilder.String())
	builder.WriteString(returnInvalidStack)
	return builder.String(), nil
}

func (g *Generator) buildOutputStackTy(
	ty tolkAbi.Ty,
	builder *strings.Builder,
	stackIndex int,
	isNullable bool,
	genericsMap map[string]tolkAbi.Ty,
) (int, error) {
	stackType := fmt.Sprintf("stack[%d].SumType", stackIndex)
	nullableCheck := ""
	if isNullable {
		nullableCheck = fmt.Sprintf(" && %s != \"VmStkNull\"", stackType)
	}

	switch ty.SumType {
	case "intN", "uintN", "varintN", "varuintN", "coins", "bool", "int":
		builder.WriteString(fmt.Sprintf("|| ((%s != \"VmStkTinyInt\" && %s != \"VmStkInt\")%s) ", stackType, stackType, nullableCheck))
	case "address", "addressExt", "addressOpt", "addressAny", "slice", "remaining", "bitsN":
		builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkSlice\"%s) ", stackType, nullableCheck))
	case "cell", "builder", "cellOf", "mapKV":
		builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkCell\"%s) ", stackType, nullableCheck))
	case "tensor":
		for _, item := range ty.TensorTy.Items {
			i, err := g.buildOutputStackTy(item, builder, stackIndex, isNullable, genericsMap)
			if err != nil {
				return 0, err
			}
			stackIndex = i + 1
		}
		return stackIndex - 1, nil
	case "tupleWith":
		for _, item := range ty.TupleWithTy.Items {
			i, err := g.buildOutputStackTy(item, builder, stackIndex, isNullable, genericsMap)
			if err != nil {
				return 0, err
			}
			stackIndex = i + 1
		}
		return stackIndex - 1, nil
	case "nullLiteral", "void":
		builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkNull\") ", stackType))
	case "callable":
		builder.WriteString(fmt.Sprintf("|| (%s != \"VmStkCont\") ", stackType))
	case "nullable":
		if ty.NullableTy.Inner == nil {
			return 0, fmt.Errorf("nullable must have inner ty")
		}
		i, err := g.buildOutputStackTy(*ty.NullableTy.Inner, builder, stackIndex, true, genericsMap)
		if err != nil {
			return 0, err
		}
		return i, nil
	case "EnumRef":
		enumTy, found := g.enumRefs[ty.EnumRefTy.EnumName]
		if !found {
			return 0, fmt.Errorf("EnumRefTy %s not found in enumRefs", ty.EnumRefTy.EnumName)
		}
		if enumTy.EncodedAs == nil {
			return 0, fmt.Errorf("enum %s has no EncodedAs", enumTy.Name)
		}
		i, err := g.buildOutputStackTy(*enumTy.EncodedAs, builder, stackIndex, isNullable, genericsMap)
		if err != nil {
			return 0, err
		}
		return i, nil
	case "AliasRef":
		aliasTy, found := g.aliasRefs[ty.AliasRefTy.AliasName]
		if !found {
			return 0, fmt.Errorf("alias %s not found in aliasRefs", ty.AliasRefTy.AliasName)
		}
		if aliasTy.TargetTy == nil {
			return 0, fmt.Errorf("alias %s has no TargetTy", ty.AliasRefTy.AliasName)
		}
		genericMap := make(map[string]tolkAbi.Ty)
		for i, genericName := range aliasTy.TypeParams {
			resolvedTy, err := g.resolveGenericT(genericsMap, ty.AliasRefTy.TypeArgs[i])
			if err != nil {
				return 0, err
			}
			genericMap[genericName] = *resolvedTy
		}
		i, err := g.buildOutputStackTy(*aliasTy.TargetTy, builder, stackIndex, isNullable, genericMap)
		if err != nil {
			return 0, err
		}
		return i, nil
	case "StructRef":
		structTy, found := g.structRefs[ty.StructRefTy.StructName]
		if !found {
			return 0, fmt.Errorf("StructRefTy %s not found in structRefs", ty.StructRefTy.StructName)
		}
		genericMap := make(map[string]tolkAbi.Ty)
		for i, genericName := range structTy.TypeParams {
			resolvedTy, err := g.resolveGenericT(genericsMap, ty.StructRefTy.TypeArgs[i])
			if err != nil {
				return 0, err
			}
			genericMap[genericName] = *resolvedTy
		}
		for _, f := range structTy.Fields {
			i, err := g.buildOutputStackTy(f.Ty, builder, stackIndex, isNullable, genericMap)
			if err != nil {
				return 0, err
			}
			stackIndex = i + 1
		}
		return stackIndex - 1, nil
	case "genericT":
		currTy, ok := genericsMap[ty.GenericTy.NameT]
		if !ok {
			return 0, fmt.Errorf("type for generic %v not found", ty.GenericTy.NameT)
		}
		i, err := g.buildOutputStackTy(currTy, builder, stackIndex, isNullable, genericsMap)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("unsupported type %s", ty.SumType)
	}
	return stackIndex, nil
}

func (g *Generator) resolveGenericT(genericMap map[string]tolkAbi.Ty, ty tolkAbi.Ty) (*tolkAbi.Ty, error) {
	switch ty.SumType {
	case "genericT":
		resolvedTy, ok := genericMap[ty.GenericTy.NameT]
		if !ok {
			return nil, fmt.Errorf("type for generic %v not found", ty.GenericTy.NameT)
		}
		return &resolvedTy, nil
	}
	return &ty, nil
}

type methodDescription struct {
	Name         string
	InvokeFnName string
}

type interfaceDescription struct {
	Name       string
	Results    []string
	GetMethods []string
}

func (g *Generator) RenderInvocationOrderList(simpleMethods []string) (string, error) {
	context := struct {
		Interfaces                     map[string]string
		InvocationOrder                []methodDescription
		InterfaceOrder                 []interfaceDescription
		KnownHashes                    map[string]interfaceDescription
		Inheritance                    map[string]string
		IntMsgs, ExtInMsgs, ExtOutMsgs map[string][]string
	}{
		Interfaces:  map[string]string{},
		KnownHashes: map[string]interfaceDescription{},
		Inheritance: map[string]string{},
		IntMsgs:     map[string][]string{},
		ExtInMsgs:   map[string][]string{},
		ExtOutMsgs:  map[string][]string{},
	}
	descriptions := map[string]methodDescription{}

	inheritance := map[string]string{}               // interface name -> parent interface
	methodsByIface := map[string]map[string]string{} // interface name -> method name -> result name

	for _, abi := range g.abi {
		for _, method := range abi.GetMethods {
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

		ifaceName := utils.ToCamelCase(abi.ContractName)
		methodsByIface[ifaceName] = map[string]string{}
		if abi.InheritsContract != "" {
			inheritance[ifaceName] = utils.ToCamelCase(abi.InheritsContract)
		}
		for _, method := range abi.GetMethods {
			if !slices.Contains(simpleMethods, method.Name) {
				continue
			}
			methodName := utils.ToCamelCase(method.Name)
			var resultName string
			if abi.IsGeneral {
				resultName = methodName + "Result"
			} else {
				resultName = fmt.Sprintf("%s_%sResult", methodName, ifaceName)
			}
			if _, ok := methodsByIface[ifaceName][methodName]; ok {
				return "", fmt.Errorf("method duplicate %v, interface %v", methodName, ifaceName)
			}
			methodsByIface[ifaceName][methodName] = resultName
		}
	}

	context.Inheritance = inheritance

	for _, abi := range g.abi {
		ifaceName := utils.ToCamelCase(abi.ContractName)
		context.Interfaces[ifaceName] = abi.ContractName
		ifaceMethods := map[string]string{}
		for currentIface := ifaceName; currentIface != ""; currentIface = inheritance[currentIface] {
			currentMethods := methodsByIface[currentIface]
			for methodName, resultName := range currentMethods {
				ifaceMethods[methodName] = resultName
			}
		}
		description := interfaceDescription{
			Name: ifaceName,
		}
		methodNames := maps.Keys(ifaceMethods)
		sort.Strings(methodNames)
		for _, methodName := range methodNames {
			description.GetMethods = append(description.GetMethods, methodName)
			description.Results = append(description.Results, ifaceMethods[methodName])
		}
		for _, m := range abi.IncomingMessages {
			name, err := m.GetMsgName()
			if err != nil {
				return "", err
			}
			if !abi.IsGeneral {
				name = abi.ContractName + name
			}
			context.IntMsgs[ifaceName] = append(context.IntMsgs[ifaceName], utils.ToCamelCase(name))
		}
		for _, m := range abi.OutgoingMessages {
			name, err := m.GetMsgName()
			if err != nil {
				return "", err
			}
			if !abi.IsGeneral {
				name = abi.ContractName + name
			}
			context.IntMsgs[ifaceName] = append(context.IntMsgs[ifaceName], utils.ToCamelCase(name))
		}
		for _, m := range abi.EmittedMessages {
			name, err := m.GetMsgName()
			if err != nil {
				return "", err
			}
			if !abi.IsGeneral {
				name = abi.ContractName + name
			}
			context.ExtOutMsgs[ifaceName] = append(context.ExtOutMsgs[ifaceName], utils.ToCamelCase(name))
		}
		if abi.IncomingExternal != nil {
			name, err := abi.IncomingExternal.GetMsgName()
			if err != nil {
				return "", err
			}
			if !abi.IsGeneral {
				name = abi.ContractName + name
			}
			context.ExtInMsgs[ifaceName] = append(context.ExtInMsgs[ifaceName], utils.ToCamelCase(name))
		}
		if len(abi.CodeHashes) > 0 { //we don't need to detect interfaces with code hashes because we can them directly
			for _, hash := range abi.CodeHashes {
				context.KnownHashes[hash] = description
			}
		} else {
			context.InterfaceOrder = append(context.InterfaceOrder, description)
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

func (g *Generator) RenderContractErrors() (string, error) {
	tmpl, err := template.New("contractErrors").Parse(contractErrorsTmpl)
	if err != nil {
		return "", err
	}
	var context = struct {
		Interfaces map[string]map[int]string
	}{
		Interfaces: map[string]map[int]string{},
	}
	for _, abi := range g.abi {
		ifaceName := utils.ToCamelCase(abi.ContractName)
		context.Interfaces[ifaceName] = map[int]string{}
		for _, e := range abi.ThrownErrors {
			if e.Name == "" {
				continue // skip unnamed errors
			}
			context.Interfaces[ifaceName][e.ErrCode] = e.Name
		}
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, context)
	return buf.String(), err

}
