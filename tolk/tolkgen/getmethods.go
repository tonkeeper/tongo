package tolkgen

import (
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

const ExecutorInterface = `
type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}
`

// StorageExecutorInterface is appended to executor.go when any contract in the group has storage.
const StorageExecutorInterface = `
type StorageExecutor interface {
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}
`

// GenerateGetMethodsCode generate separate decode and execute functions
func (tgen TolkGolangGenerator) GenerateGetMethodsCode() (string, error) {
	if len(tgen.abi.GetMethods) == 0 {
		return "", nil
	}
	out := &strings.Builder{}
	for _, method := range tgen.abi.GetMethods {
		if err := tgen.genGetMethod(method, out); err != nil {
			return "", fmt.Errorf("get method %q: %w", method.Name, err)
		}
	}
	return out.String(), nil
}

func (tgen TolkGolangGenerator) genGetMethod(method parser.ABIGetMethod, out *strings.Builder) error {
	goFnName := MethodGoName(method.Name)
	decodeFnName := "Decode" + goFnName

	retGoType, err := tgen.stackReturnGoType(method.ReturnTyIdx)
	if err != nil {
		return err
	}

	expectedSize, err := tgen.calcWidthOnStack(method.ReturnTyIdx)
	if err != nil {
		return err
	}

	readExpr, hasMethod, err := tgen.emitStackReadExpr("result", method.ReturnTyIdx, false)
	if err != nil {
		return fmt.Errorf("return type expr generation failed: %w", err)
	}
	paramsDecl, err := tgen.paramsDecl(method.Parameters)
	if err != nil {
		return err
	}

	pushCode := ""
	if len(method.Parameters) > 0 {
		var sb strings.Builder
		for _, p := range method.Parameters {
			pushOne, err := tgen.stackPushCode(utils.ToCamelCasePrivate(p.Name), p.TyIdx)
			if err != nil {
				return fmt.Errorf("param %q: %w", p.Name, err)
			}
			sb.WriteString(pushOne)
		}
		pushCode = sb.String()
	}
	decodeRead := "\treturn " + readExpr + "\n"
	if hasMethod {
		decodeRead = "\terr = result.ReadFromStack(stack)\n\treturn\n"
	}
	constName := "MethodID" + goFnName
	code, err := renderGetMethodsTemplate(getMethodTpl, map[string]any{
		"DecodeFnName": decodeFnName,
		"RetGoType":    retGoType,
		"ExpectedSize": expectedSize,
		"DecodeRead":   decodeRead,
		"ConstName":    constName,
		"MethodID":     method.TVMMethodID,
		"GoFnName":     goFnName,
		"ParamsDecl":   paramsDecl,
		"PushCode":     pushCode,
	})
	if err != nil {
		return fmt.Errorf("render get method %q: %w", method.Name, err)
	}
	out.WriteString(code)
	return nil
}

// int257Conversion builds an expression converting an integer parameter to
// tlb.Int257 for pushing onto the VM stack. tlb.Int257 is backed by big.Int, so
// big.Int-backed tlb types (Int257, Var(U)IntegerN, and Int/UintN with N > 64)
// convert directly, while the fixed-width primitive types (Int/UintN with
// N <= 64 and Coins) must go through Int257FromInt64.
func int257Conversion(ty parser.Ty, paramGoName string) string {
	bigIntBacked := false
	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindVarIntN, parser.TyKindVarUintN:
		bigIntBacked = true
	case parser.TyKindIntN:
		bigIntBacked = ty.IntN.N > 64
	case parser.TyKindUintN:
		bigIntBacked = ty.UintN.N > 64
	}
	if bigIntBacked {
		return fmt.Sprintf("tlb.Int257(%s)", paramGoName)
	}
	return fmt.Sprintf("tlb.Int257FromInt64(int64(%s))", paramGoName)
}

func (tgen TolkGolangGenerator) stackPushCode(paramGoName string, tyIdx int) (string, error) {
	ty, err := tgen.symbols.TyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins:
		expr := int257Conversion(ty, paramGoName)
		return fmt.Sprintf("\tstack.Put(tlb.VmStackValue{SumType: \"VmStkInt\", VmStkInt: %s})\n", expr), nil
	case parser.TyKindBool:
		return fmt.Sprintf(`	if %s {
		stack.Put(tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: 1})
	} else {
		stack.Put(tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: 0})
	}
`, paramGoName), nil
	case parser.TyKindCell:
		return fmt.Sprintf(`	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(%s)
	if err != nil {
		err = fmt.Errorf("encode param %s: %%w", err)
		return
	}
	stack.Put(val)
	}
`, paramGoName, paramGoName), nil
	case parser.TyKindSlice:
		return fmt.Sprintf(`	{
		var val tlb.VmStackValue
		val, err = tlb.CellToVmCellSlice(&%s)
	if err != nil {
		err = fmt.Errorf("encode param %s: %%w", err)
		return
	}
	stack.Put(val)
	}
`, paramGoName, paramGoName), nil
	case parser.TyKindVoid:
		return "", fmt.Errorf("unsupported void get-method parameter")
	default:
		return fmt.Sprintf(`	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(%s)
	if err != nil {
		err = fmt.Errorf("encode param %s: %%w", err)
		return
	}
	stack.Put(val)
	}
`, paramGoName, paramGoName), nil
	}
}

func (tgen TolkGolangGenerator) emitStackReadExpr(fieldPath string, tyIdx int, unTupleIfW bool) (string, bool, error) {
	ty, err := tgen.symbols.TyByIdx(tyIdx)
	if err != nil {
		return "", false, err
	}
	if unTupleIfW { // inside `array<T>` or `[T, ...]`, if T is non-primitive, it's a sub-tuple
		wOnStack, err := tgen.calcWidthOnStack(tyIdx)
		if err != nil {
			return "", false, fmt.Errorf("calc width on stack: %w", err)
		}
		if wOnStack != 1 {
			expr, hasMethod, err := tgen.emitStackReadExpr(fieldPath, tyIdx, false)
			if err != nil {
				return "", false, fmt.Errorf("tuple inner: %w", err)
			}
			goType, err := tgen.emitStackGoType(tyIdx)
			if err != nil {
				return "", false, fmt.Errorf("tuple inner type: %w", err)
			}
			if hasMethod {
				return fmt.Sprintf(`tlb.ReadTupleFromStack(stack, func(stack *tlb.VmStack) (result %s, err error) {
	err = result.ReadFromStack(stack)
	return
})`, goType), false, nil
			} else {
				return fmt.Sprintf(`tlb.ReadTupleFromStack(stack, func(stack *tlb.VmStack) (%s, error) {
	return %s
})`, goType, expr), false, nil
			}
		}
	}

	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindStructRef, parser.TyKindAddress, parser.TyKindAddressOpt:
		goType, err := tgen.symbols.emitGoType(tyIdx)
		if err != nil {
			return "", false, fmt.Errorf("type: %w", err)
		}
		return fmt.Sprintf("tlb.ReadFromStack[%s](stack)", goType), true, nil
	case parser.TyKindCell:
		return "stack.ReadCell()", false, nil
	case parser.TyKindBool:
		return "stack.ReadBool()", false, nil
	case parser.TyKindArrayOf:
		innerExpr, hasMethod, err := tgen.emitStackReadExpr(fieldPath, ty.ArrayOf.InnerTyIdx, true)
		if err != nil {
			return "", false, fmt.Errorf("array inner: %w", err)
		}
		goType, err := tgen.emitStackGoType(ty.ArrayOf.InnerTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("array inner type: %w", err)
		}
		if hasMethod {
			return fmt.Sprintf(
				`tlb.ReadArrayFromStack[%s](stack, func(stack *tlb.VmStack) (value %s, err error) {
		err = value.ReadFromStack(stack)
		return
	})`, goType, goType), false, nil
		} else {
			return fmt.Sprintf(
				`tlb.ReadArrayFromStack[%s](stack, func(stack *tlb.VmStack) (value %s, err error) {
		return %s
	})`, goType, goType, innerExpr), false, nil
		}
	case parser.TyKindNullable:
		innerRead, _, err := tgen.emitStackReadExpr(fieldPath, ty.Nullable.InnerTyIdx, unTupleIfW)
		if err != nil {
			return "", false, fmt.Errorf("nullable inner: %w", err)
		}
		innerType, err := tgen.emitStackGoType(ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("nullable inner type: %w", err)
		}
		if ty.Nullable.StackTypeId != 0 {
			if ty.Nullable.StackWidth <= 0 {
				return "", false, fmt.Errorf("nullable type with invalid stack width %d", ty.Nullable.StackWidth)
			}
			return fmt.Sprintf(`tlb.StackReadWideMaybeCallback(stack, %d, func (stack *tlb.VmStack) (%s, error) {
	return %s
})`, ty.Nullable.StackWidth, innerType, innerRead), false, nil
		}
		return fmt.Sprintf(`tlb.StackReadMaybeCallback(stack, func (stack *tlb.VmStack) (%s, error) {
	return %s
})`, innerType, innerRead), false, nil
	case parser.TyKindTensor:
		items := ty.Tensor.ItemsTyIdx
		if len(items) > tlb.MaxTensorSize {
			return "", false, fmt.Errorf("tensor size %d is too big, update tlb/tensor.go", len(items))
		}
		loaders := make([]string, len(items))
		for i, itemTyIdx := range items {
			expr, hasMethod, err := tgen.emitStackReadExpr(fmt.Sprintf("result.V%d", i), itemTyIdx, false)
			if err != nil {
				return "", false, fmt.Errorf("tensor item %d: %w", i, err)
			}
			if hasMethod {
				loaders[i] = fmt.Sprintf(`
	if err = result.V%d.ReadFromStack(stack); err != nil {
		return
	}`, i)
			} else {
				loaders[i] = fmt.Sprintf(`
	if result.V%d, err = %s; err != nil {
	return
}`, i, expr)
			}
		}
		slices.Reverse(loaders)
		typ, err := tgen.emitStackGoType(tyIdx)
		if err != nil {
			return "", false, fmt.Errorf("tensor type: %w", err)
		}
		return fmt.Sprintf(`(func () (result %s, err error) {
		%s
		return
	})()`, typ, strings.Join(loaders, "")), false, nil
	case parser.TyKindShapedTuple:
		items := ty.ShapedTuple.ItemsTyIdx
		if len(items) > tlb.MaxTensorSize {
			return "", false, fmt.Errorf("tensor size %d is too big, update tlb/tensor.go", len(items))
		}
		loaders := make([]string, len(items))
		for i, itemTyIdx := range items {
			expr, hasMethod, err := tgen.emitStackReadExpr(fmt.Sprintf("result.V%d", i), itemTyIdx, true)
			if err != nil {
				return "", false, fmt.Errorf("tensor item %d: %w", i, err)
			}
			if hasMethod {
				loaders[i] = fmt.Sprintf(`
	if err = result.V%d.ReadFromStack(stack); err != nil {
		return
	}`, i)
			} else {
				loaders[i] = fmt.Sprintf(`
	if result.V%d, err = %s; err != nil {
	return
}`, i, expr)
			}
		}
		slices.Reverse(loaders)
		typ, err := tgen.emitStackGoType(tyIdx)
		if err != nil {
			return "", false, fmt.Errorf("tensor type: %w", err)
		}
		return fmt.Sprintf(`(func () (result %s, err error) {
	var tuple tlb.VmStkTuple
	if tuple, err = stack.ReadTuple(); err != nil {
		err = fmt.Errorf("read tensor: %%w", err)
		return
	}
	var stack *tlb.VmStack
	stack, err = tuple.AsStack()
	if err != nil {
		err = fmt.Errorf("read tensor items: %%w", err)
		return
	}
		%s
		return
	})()`, typ, strings.Join(loaders, "")), false, nil
	case parser.TyKindEnumRef:
		enum := tgen.symbols.Enums[ty.EnumRef.EnumName]
		return tgen.emitStackReadExpr(fieldPath, enum.EncodedAsTyIdx, unTupleIfW)
	case parser.TyKindMapKV:
		// A dictionary on the stack is null (empty) or a cell holding the dictionary
		// root loaded directly: a Hashmap, not the in-cell HashmapE (which has a Maybe
		// prefix and stores the root behind a ref).
		kType, err := tgen.emitStackGoType(ty.MapKV.KeyTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("map key type: %w", err)
		}
		vType, err := tgen.emitStackGoType(ty.MapKV.ValueTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("map value type: %w", err)
		}
		return fmt.Sprintf("tlb.ReadHashmapFromStack[%s, %s](stack)", kType, vType), false, nil
	default:
		renderedTy, _ := tgen.symbols.RenderTy(tyIdx)
		return "", false, fmt.Errorf("unexpected type in stack read: %s", renderedTy)
	}
}

// paramsDecl builds the extra parameter declarations for a get method signature,
// e.g. ", hash tlb.Int257, amount tlb.Coins".
func (tgen TolkGolangGenerator) paramsDecl(params []parser.ABIGetMethodParameter) (string, error) {
	if len(params) == 0 {
		return "", nil
	}
	var sb strings.Builder
	for _, p := range params {
		goName := utils.ToCamelCasePrivate(p.Name)
		goType, err := tgen.symbols.emitGoType(p.TyIdx)
		if err != nil {
			return "", fmt.Errorf("param %q type: %w", p.Name, err)
		}
		fmt.Fprintf(&sb, ", %s %s", goName, goType)
	}
	return sb.String(), nil
}

func (tgen TolkGolangGenerator) paramsCall(params []parser.ABIGetMethodParameter) string {
	if len(params) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, p := range params {
		fmt.Fprintf(&sb, ", %s", utils.ToCamelCasePrivate(p.Name))
	}
	return sb.String()
}

type clientMethodSpec struct {
	GoFnName   string
	RetGoType  string
	ParamsDecl string
	ParamsCall string
}

func (tgen TolkGolangGenerator) clientMethodSpecs() ([]clientMethodSpec, error) {
	specs := make([]clientMethodSpec, 0, len(tgen.abi.GetMethods))
	for _, method := range tgen.abi.GetMethods {
		goFnName := MethodGoName(method.Name)
		retGoType, err := tgen.stackReturnGoType(method.ReturnTyIdx)
		if err != nil {
			return nil, fmt.Errorf("get method %q return type: %w", method.Name, err)
		}
		paramsDecl, err := tgen.paramsDecl(method.Parameters)
		if err != nil {
			return nil, fmt.Errorf("get method %q params: %w", method.Name, err)
		}
		specs = append(specs, clientMethodSpec{
			GoFnName:   goFnName,
			RetGoType:  retGoType,
			ParamsDecl: paramsDecl,
			ParamsCall: tgen.paramsCall(method.Parameters),
		})
	}
	return specs, nil
}

// GenerateClientInterface generates an interface that carries all get method signatures,
// plus a WithAccount variant interface (accountID curried in).
func (tgen TolkGolangGenerator) GenerateClientInterface(ifaceName string) (string, error) {
	if len(tgen.abi.GetMethods) == 0 {
		return "", nil
	}
	withAccountName := ifaceName + "WithAccount"
	specs, err := tgen.clientMethodSpecs()
	if err != nil {
		return "", err
	}

	ifaceCode, err := renderGetMethodsTemplate(clientInterfacesTpl, clientInterfaceTemplateData{
		IfaceName:       ifaceName,
		WithAccountName: withAccountName,
		Methods:         specs,
		HasStorage:      tgen.storageType != "",
		StorageType:     tgen.storageType,
	})
	if err != nil {
		return "", fmt.Errorf("render client interfaces: %w", err)
	}
	return ifaceCode, nil
}

// GenerateClientImpl generates:
//   - the main impl struct + constructor + methods (including WithAccountId)
//   - the WithAccount impl struct + methods (accountID curried in)
func (tgen TolkGolangGenerator) GenerateClientImpl(ifaceName string) (string, error) {
	if len(tgen.abi.GetMethods) == 0 {
		return "", nil
	}
	implName := strings.ToLower(ifaceName[:1]) + ifaceName[1:] + "Impl"
	withAccountName := ifaceName + "WithAccount"
	withAccountImplName := strings.ToLower(ifaceName[:1]) + ifaceName[1:] + "WithAccountImpl"

	specs, err := tgen.clientMethodSpecs()
	if err != nil {
		return "", err
	}
	storageFnName := ifaceName + "_AccountState"

	implCode, err := renderGetMethodsTemplate(clientImplTpl, clientImplTemplateData{
		IfaceName:       ifaceName,
		WithAccountName: withAccountName,
		ImplName:        implName,
		WithAccountImpl: withAccountImplName,
		Methods:         specs,
		HasStorage:      tgen.storageType != "",
		StorageType:     tgen.storageType,
		StorageFnName:   storageFnName,
	})
	if err != nil {
		return "", fmt.Errorf("render client impl: %w", err)
	}
	return implCode, nil
}

// GenerateStorageCode generates the free AccountState function for a contract that has storage.
// The function is named {ifaceName}_AccountState to avoid collisions across contracts.
func (tgen TolkGolangGenerator) GenerateStorageCode(ifaceName string) (string, error) {
	if tgen.storageType == "" {
		return "", nil
	}
	fnName := ifaceName + "_AccountState"
	return fmt.Sprintf(`func %s(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage %s, err error) {
	sa, err = executor.GetAccountState(ctx, accountID)
	if err != nil {
		return
	}
	acc := sa.Account
	if acc.SumType != "Account" {
		err = fmt.Errorf("account does not exist")
	} else if state := acc.Account.Storage.State; state.SumType != "AccountActive" {
		err = fmt.Errorf("account is not active")
	} else if data := state.AccountActive.StateInit.Data; !data.Exists {
		err = fmt.Errorf("account has no storage data")
	} else {
		err = storage.UnmarshalTLB(&data.Value.Value, tlb.NewDecoder())
	}
	return
}
`, fnName, tgen.storageType), nil
}

// GenerateExternalMessagesCode generates a ToTLBMessage method for each IncomingExternal message type.
// The method marshals the message body, optionally builds a tlb.StateInit from a typed StateInitT,
// calls ton.CreateExternalMessage and returns the resulting tlb.Message.
func (tgen TolkGolangGenerator) GenerateExternalMessagesCode() (string, error) {
	if len(tgen.abi.IncomingExternal) == 0 {
		return "", nil
	}

	var out strings.Builder
	for _, ext := range tgen.abi.IncomingExternal {
		msgName, err := tgen.symbols.MsgName(ext.BodyTyIdx)
		if err != nil {
			continue
		}
		goMsgType := safeGoIdent(msgName)

		if tgen.storageType != "" {
			fmt.Fprintf(&out, `func (msg %s) ToExternal(address ton.AccountID, init *tlb.StateInitT[*%s]) (tlb.Message, error) {
	return ton.CreateExternalMessageTWithState(address, msg, init, tlb.VarUInteger16{})
}

`, goMsgType, tgen.storageType)
		} else {
			fmt.Fprintf(&out, `func (msg %s) ToExternal(address ton.AccountID, init *tlb.StateInit) (tlb.Message, error) {
	return ton.CreateExternalMessageT(address, msg, init, tlb.VarUInteger16{})
}

`, goMsgType)
		}
	}

	return out.String(), nil
}

// GenerateInternalMessagesCode generates a ToInternal method for each IncomingMessage type
// and for any union alias whose variants are all IncomingMessage types.
// Each method is a thin wrapper around tlb.BuildInternal.
func (tgen TolkGolangGenerator) GenerateInternalMessagesCode() (string, error) {
	if len(tgen.abi.IncomingMessages) == 0 {
		return "", nil
	}
	var out strings.Builder

	// Build set of struct names referenced in incomingMessages.
	incomingMsgNames := make(map[string]struct{}, len(tgen.abi.IncomingMessages))
	for _, msg := range tgen.abi.IncomingMessages {
		msgName, err := tgen.symbols.MsgName(msg.BodyTyIdx)
		if err != nil {
			continue
		}
		incomingMsgNames[msgName] = struct{}{}
	}

	// Determine the state-init type parameter: typed storage if defined, raw Any otherwise.
	var initParamType string
	if tgen.storageType != "" {
		initParamType = fmt.Sprintf("*tlb.StateInitT[*%s]", tgen.storageType)
	} else {
		initParamType = "*tlb.StateInitT[*tlb.Any]"
	}

	writeToInternal := func(goMsgType string) {
		fmt.Fprintf(&out, "func (msg %s) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init %s) (tlb.Message, error) {\n",
			goMsgType, initParamType)
		fmt.Fprintf(&out, "\treturn tlb.BuildInternal(&msg, dest, amount, bounce, init)\n")
		fmt.Fprintf(&out, "}\n\n")
	}

	// Generate ToInternal for each direct incoming message type.
	for _, msg := range tgen.abi.IncomingMessages {
		msgName, err := tgen.symbols.MsgName(msg.BodyTyIdx)
		if err != nil {
			continue
		}
		writeToInternal(safeGoIdent(msgName))
	}

	// Generate ToInternal for union aliases whose every variant is an incoming message type.
	for _, decl := range tgen.abi.Declarations {
		if decl.SumType != parser.DeclarationKindAlias {
			continue
		}
		alias := decl.AliasDeclaration
		targetTy, err := tgen.symbols.TyByIdx(alias.TargetTyIdx)
		if err != nil || targetTy.SumType != parser.TyKindUnion {
			continue
		}
		allInternal := len(targetTy.Union.Variants) > 0
		for _, variant := range targetTy.Union.Variants {
			variantTy, err := tgen.symbols.TyByIdx(variant.VariantTyIdx)
			if err != nil || variantTy.SumType != parser.TyKindStructRef {
				allInternal = false
				break
			}
			if _, ok := incomingMsgNames[variantTy.StructRef.StructName]; !ok {
				allInternal = false
				break
			}
		}
		if !allInternal {
			continue
		}
		writeToInternal(safeGoIdent(alias.Name))
	}

	return out.String(), nil
}

// MethodGoName converts a Tolk get-method name to its Go exported function name.
// E.g. "get_cocoon_data" → "GetCocoonData", "last_proxy_seqno" → "GetLastProxySeqno".
func MethodGoName(name string) string {
	camel := utils.ToCamelCase("get_" + name)
	if strings.HasPrefix(camel, "GetGet") {
		camel = "Get" + camel[6:]
	}
	return camel
}

type clientInterfaceTemplateData struct {
	IfaceName       string
	WithAccountName string
	Methods         []clientMethodSpec
	HasStorage      bool
	StorageType     string
}

type clientImplTemplateData struct {
	IfaceName       string
	WithAccountName string
	ImplName        string
	WithAccountImpl string
	Methods         []clientMethodSpec
	HasStorage      bool
	StorageType     string
	StorageFnName   string
}

const getMethodTpl = `func {{.DecodeFnName}}(stack *tlb.VmStack) (result {{.RetGoType}}, err error) {
	if stack.Len() != {{.ExpectedSize}} {
		err = fmt.Errorf("invalid stack size %d, expected {{.ExpectedSize}}", stack.Len())
		return
	}
{{.DecodeRead}}}

const {{.ConstName}} = 0x{{printf "%X" .MethodID}}

func {{.GoFnName}}(ctx context.Context, executor Executor, reqAccountID ton.AccountID{{.ParamsDecl}}) (result {{.RetGoType}}, err error) {
	var errCode uint32
	var stack tlb.VmStack
{{.PushCode}}	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, {{.ConstName}}, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return {{.DecodeFnName}}(&stack)
}

`

const clientInterfacesTpl = `type {{.IfaceName}} interface {
	WithAccountId(accountID ton.AccountID) {{.WithAccountName}}{{range .Methods}}
	{{.GoFnName}}(ctx context.Context, reqAccountID ton.AccountID{{.ParamsDecl}}) ({{.RetGoType}}, error){{end}}{{if .HasStorage}}
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, {{.StorageType}}, error){{end}}
}

type {{.WithAccountName}} interface {
{{range .Methods}}	{{.GoFnName}}(ctx context.Context{{.ParamsDecl}}) ({{.RetGoType}}, error)
{{end}}{{if .HasStorage}}	AccountState(ctx context.Context) (tlb.ShardAccount, {{.StorageType}}, error)
{{end}}}
`

const clientImplTpl = `{{if .HasStorage}}type {{.ImplName}} struct {
	executor Executor
	storageExecutor StorageExecutor
}

func New{{.IfaceName}}(executor Executor, storageExecutor StorageExecutor) {{.IfaceName}} {
	return &{{.ImplName}}{executor: executor, storageExecutor: storageExecutor}
}

func (c {{.ImplName}}) WithAccountId(accountID ton.AccountID) {{.WithAccountName}} {
	return &{{.WithAccountImpl}}{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}
{{else}}type {{.ImplName}} struct {
	executor Executor
}

func New{{.IfaceName}}(executor Executor) {{.IfaceName}} {
	return &{{.ImplName}}{executor: executor}
}

func (c {{.ImplName}}) WithAccountId(accountID ton.AccountID) {{.WithAccountName}} {
	return &{{.WithAccountImpl}}{executor: c.executor, accountID: accountID}
}
{{end}}{{range .Methods}}
func (c {{$.ImplName}}) {{.GoFnName}}(ctx context.Context, reqAccountID ton.AccountID{{.ParamsDecl}}) ({{.RetGoType}}, error) {
	return {{.GoFnName}}(ctx, c.executor, reqAccountID{{.ParamsCall}})
}
{{end}}{{if .HasStorage}}
func (c {{.ImplName}}) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, {{.StorageType}}, error) {
	return {{.StorageFnName}}(ctx, c.storageExecutor, reqAccountID)
}
{{end}}{{if .HasStorage}}
type {{.WithAccountImpl}} struct {
	executor Executor
	storageExecutor StorageExecutor
	accountID ton.AccountID
}
{{else}}
type {{.WithAccountImpl}} struct {
	executor Executor
	accountID ton.AccountID
}
{{end}}{{range .Methods}}
func (c {{$.WithAccountImpl}}) {{.GoFnName}}(ctx context.Context{{.ParamsDecl}}) ({{.RetGoType}}, error) {
	return {{.GoFnName}}(ctx, c.executor, c.accountID{{.ParamsCall}})
}
{{end}}{{if .HasStorage}}
func (c {{.WithAccountImpl}}) AccountState(ctx context.Context) (tlb.ShardAccount, {{.StorageType}}, error) {
	return {{.StorageFnName}}(ctx, c.storageExecutor, c.accountID)
}
{{end}}`

func renderGetMethodsTemplate(tpl string, data any) (string, error) {
	t, err := template.New("getmethods").Parse(tpl)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := t.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}
