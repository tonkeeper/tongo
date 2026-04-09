package tolkgen

import (
	"fmt"
	"strings"
	"text/template"

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

func (tgen TolkGolangGenerator) genGetMethod(method parser.GetMethod, out *strings.Builder) error {
	goFnName := MethodGoName(method.Name)
	decodeFnName := "Decode" + goFnName

	retGoType, err := stackReturnGoType(method.ReturnTy)
	if err != nil {
		return err
	}

	expectedSize, err := tgen.calcWidthOnStack(method.ReturnTy)
	if err != nil {
		return err
	}

	readExpr, hasMethod, err := tgen.emitStackReadExpr("result", method.ReturnTy, false)
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
			pushOne, err := stackPushCode(utils.ToCamelCasePrivate(p.Name), p.Ty)
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
		"MethodID":     method.TvmMethodID,
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

func stackPushCode(paramGoName string, ty parser.Ty) (string, error) {
	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins:
		return fmt.Sprintf("\tstack.Put(tlb.VmStackValue{SumType: \"VmStkInt\", VmStkInt: tlb.Int257(%s)})\n", paramGoName), nil
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

func (tgen TolkGolangGenerator) emitStackReadExpr(fieldPath string, ty parser.Ty, unTupleIfW bool) (string, bool, error) {
	if unTupleIfW { // inside `array<T>` or `[T, ...]`, if T is non-primitive, it's a sub-tuple
		wOnStack, err := tgen.calcWidthOnStack(ty)
		if err != nil {
			return "", false, err
		}
		if wOnStack != 1 {
			return "", false, fmt.Errorf("non-primitive type in array/slice: %s (width=%d)", ty.String(), wOnStack)
		}
	}

	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindStructRef, parser.TyKindAddress, parser.TyKindAddressOpt:
		goType, err := emitGoType(ty)
		if err != nil {
			return "", false, fmt.Errorf("type %s: %w", ty.String(), err)
		}
		return fmt.Sprintf("tlb.ReadFromStack[%s](stack)", goType), true, nil
	case parser.TyKindCell:
		return "stack.ReadCell()", false, nil
	case parser.TyKindBool:
		return "stack.ReadBool()", false, nil
	case parser.TyKindNullable:
		if ty.Nullable.StackWidth != 0 {
			return "", false, fmt.Errorf("nullable type with wide stack is not implemented yet: %s", ty.String())
		}
		innerRead, _, err := tgen.emitStackReadExpr(fieldPath, ty.Nullable.Inner, unTupleIfW)
		if err != nil {
			return "", false, fmt.Errorf("nullable inner: %w", err)
		}
		innerType, err := emitGoType(ty.Nullable.Inner)
		if err != nil {
			return "", false, fmt.Errorf("nullable inner type %s: %w", ty.Nullable.Inner.String(), err)
		}
		return fmt.Sprintf(`tlb.StackReadMaybeCallback(stack, func (stack *tlb.VmStack) (%s, error) {
	return %s
})`, innerType, innerRead), false, nil
	default:
		return "", false, fmt.Errorf("unexpected type in stack read: %s", ty.String())
	}
}

// paramsDecl builds the extra parameter declarations for a get method signature,
// e.g. ", hash tlb.Int257, amount tlb.Coins".
func (tgen TolkGolangGenerator) paramsDecl(params []parser.Parameter) (string, error) {
	if len(params) == 0 {
		return "", nil
	}
	var sb strings.Builder
	for _, p := range params {
		goName := utils.ToCamelCasePrivate(p.Name)
		goType, err := emitGoType(p.Ty)
		if err != nil {
			return "", fmt.Errorf("param %q type %s: %w", p.Name, p.Ty.String(), err)
		}
		fmt.Fprintf(&sb, ", %s %s", goName, goType)
	}
	return sb.String(), nil
}

func (tgen TolkGolangGenerator) paramsCall(params []parser.Parameter) string {
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
		retGoType, err := stackReturnGoType(method.ReturnTy)
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
		msgName, err := ext.GetMsgName()
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
		msgName, err := msg.GetMsgName()
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
		initParamType = "*tlb.StateInitT[tlb.Any]"
	}

	writeToInternal := func(goMsgType string) {
		fmt.Fprintf(&out, "func (msg %s) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init %s) (tlb.Message, error) {\n",
			goMsgType, initParamType)
		fmt.Fprintf(&out, "\treturn tlb.BuildInternal(msg, dest, amount, bounce, init)\n")
		fmt.Fprintf(&out, "}\n\n")
	}

	// Generate ToInternal for each direct incoming message type.
	for _, msg := range tgen.abi.IncomingMessages {
		msgName, err := msg.GetMsgName()
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
		if alias.TargetTy.SumType != parser.TyKindUnion {
			continue
		}
		allInternal := len(alias.TargetTy.Union.Variants) > 0
		for _, variant := range alias.TargetTy.Union.Variants {
			if variant.VariantTy.SumType != parser.TyKindStructRef {
				allInternal = false
				break
			}
			if _, ok := incomingMsgNames[variant.VariantTy.StructRef.StructName]; !ok {
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
