package tolkgen

import (
	"fmt"
	"slices"
	"strings"

	//"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

// GenerateContractFile assembles the main content of a generated contract Go file:
// type declarations, get-method functions, storage helpers, and client impl.
// Returns an empty string (and no error) when there is nothing to emit.
func (gen TolkGolangGenerator) GenerateContractFile(pkgName, ifaceName string) (string, error) {
	declCode, _, err := gen.GenerateGocode()
	if err != nil {
		return "", fmt.Errorf("declarations: %w", err)
	}
	methodCode, err := gen.GenerateGetMethodsCode()
	if err != nil {
		return "", fmt.Errorf("get methods: %w", err)
	}
	clientImplCode := ""
	if ifaceName != "" {
		clientImplCode, err = gen.GenerateClientImpl(ifaceName)
		if err != nil {
			return "", fmt.Errorf("client impl: %w", err)
		}
	}

	storageCode, err := gen.GenerateStorageCode(ifaceName)
	if err != nil {
		return "", fmt.Errorf("storage: %w", err)
	}

	allCode := declCode
	if methodCode != "" {
		allCode += "\n" + methodCode
	}
	if storageCode != "" {
		allCode += "\n" + storageCode
	}
	if clientImplCode != "" {
		allCode += "\n" + clientImplCode
	}
	if allCode == "" {
		return "", nil
	}

	imports := deriveImportBlock(allCode, nil)
	return fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n\nimport (\n%s\n)\n\n", pkgName, imports) +
		allCode + "\n\n", nil
}

// GenerateMarshalFile assembles the content of the separate marshal Go file for a contract:
// UnmarshalTLB / MarshalTLB / ReadFromStack methods on all generated types, plus
// ToTLBMessage helpers for incoming external messages.
// Returns an empty string (and no error) when there is nothing to emit.
func (gen TolkGolangGenerator) GenerateMarshalFile(pkgName string) (string, error) {
	_, marshalCode, err := gen.GenerateGocode()
	if err != nil {
		return "", fmt.Errorf("declarations: %w", err)
	}
	extMsgCode, err := gen.GenerateExternalMessagesCode()
	if err != nil {
		return "", fmt.Errorf("external messages: %w", err)
	}
	intMsgCode, err := gen.GenerateInternalMessagesCode()
	if err != nil {
		return "", fmt.Errorf("internal messages: %w", err)
	}

	allCode := marshalCode
	if extMsgCode != "" {
		allCode += "\n" + extMsgCode
	}
	if intMsgCode != "" {
		allCode += "\n" + intMsgCode
	}
	if allCode == "" {
		return "", nil
	}

	imports := deriveImportBlock(allCode, nil)
	return fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n\nimport (\n%s\n)\n\n", pkgName, imports) +
		allCode + "\n", nil
}

// ExecutorEntry is one per-contract entry used when generating executor.go.
type ExecutorEntry struct {
	IfaceName string
	Gen       *TolkGolangGenerator
}

// GenerateExecutorFile assembles the full content of executor.go for a package.
func GenerateExecutorFile(pkgName string, entries []ExecutorEntry) (string, error) {
	// Determine if any contract in this group has storage (needs GetAccountState).
	needsStorage := false
	for _, e := range entries {
		if e.Gen.storageType != "" {
			needsStorage = true
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(ExecutorInterface)
	if needsStorage {
		sb.WriteString("\n")
		sb.WriteString(StorageExecutorInterface)
	}
	for _, entry := range entries {
		iface, err := entry.Gen.GenerateClientInterface(entry.IfaceName)
		if err != nil {
			return "", fmt.Errorf("client interface %s: %w", entry.IfaceName, err)
		}
		if iface != "" {
			sb.WriteString("\n")
			sb.WriteString(iface)
		}
	}

	ifaceBody := sb.String()
	imports := deriveImportBlock(ifaceBody, nil)
	return fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n\nimport (\n%s\n)\n\n%s", pkgName, imports, ifaceBody), nil
}

type TolkGolangGenerator struct {
	symbols     *symTable
	abi         parser.ContractABI
	storageType string // empty when no storage
}

func NewTolkGolangGenerator(abi parser.ContractABI) (*TolkGolangGenerator, error) {
	symbols := symTable{
		aliases:                       make(map[string]parser.ABIAlias),
		structs:                       make(map[string]parser.ABIStruct),
		enums:                         make(map[string]parser.ABIEnum),
		uniqueTypes:                   abi.UniqueTypes,
		structInstantiations:          abi.StructInstantiations,
		aliasInstantiations:           abi.AliasInstantiations,
		structIsReturnedFromGetMethod: make(map[string]bool),
		enumNeedsReadFromStack:        make(map[string]bool),
	}
	for _, decl := range abi.Declarations {
		switch decl.SumType {
		case parser.DeclarationKindAlias:
			symbols.aliases[decl.AliasDeclaration.Name] = decl.AliasDeclaration
		case parser.DeclarationKindStruct:
			symbols.structs[decl.StructDeclaration.Name] = decl.StructDeclaration
		case parser.DeclarationKindEnum:
			symbols.enums[decl.EnumDeclaration.Name] = decl.EnumDeclaration
		}
	}
	for _, method := range abi.GetMethods {
		symbols.markStructsReadFromStack(method.ReturnTyIdx)
	}
	var storageType string
	if abi.Storage.StorageTyIdx != nil {
		var err error
		storageType, err = symbols.emitGoType(*abi.Storage.StorageTyIdx)
		if err != nil {
			return nil, fmt.Errorf("storage type: %w", err)
		}
	}
	return &TolkGolangGenerator{
		abi:         abi,
		symbols:     &symbols,
		storageType: storageType,
	}, nil
}

type symTable struct {
	aliases                       map[string]parser.ABIAlias
	structs                       map[string]parser.ABIStruct
	enums                         map[string]parser.ABIEnum
	uniqueTypes                   []parser.Ty
	structInstantiations          []parser.ABIStructInstantiation
	aliasInstantiations           []parser.ABIAliasInstantiation
	structIsReturnedFromGetMethod map[string]bool
	enumNeedsReadFromStack        map[string]bool
}

func (s *symTable) markStructsReadFromStack(tyIdx int) {
	ty, err := s.tyByIdx(tyIdx)
	if err != nil {
		return
	}
	switch ty.SumType {
	case parser.TyKindArrayOf:
		s.markStructsReadFromStack(ty.ArrayOf.InnerTyIdx)
	case parser.TyKindStructRef:
		structName := ty.StructRef.StructName
		if s.structIsReturnedFromGetMethod[structName] {
			return
		}
		s.structIsReturnedFromGetMethod[structName] = true
		fields, err := s.structFieldsOf(tyIdx, true)
		if err != nil {
			return
		}
		for _, field := range fields {
			s.markStructsReadFromStack(field.TyIdx)
		}
	case parser.TyKindEnumRef:
		s.enumNeedsReadFromStack[ty.EnumRef.EnumName] = true
	case parser.TyKindTensor:
		for _, item := range ty.Tensor.ItemsTyIdx {
			s.markStructsReadFromStack(item)
		}
	case parser.TyKindShapedTuple:
		for _, item := range ty.ShapedTuple.ItemsTyIdx {
			s.markStructsReadFromStack(item)
		}
	case parser.TyKindNullable:
		s.markStructsReadFromStack(ty.Nullable.InnerTyIdx)
	case parser.TyKindAliasRef:
		targetTyIdx, _, err := s.aliasTargetOf(tyIdx)
		if err != nil {
			return
		}
		s.markStructsReadFromStack(targetTyIdx)
	case parser.TyKindUnion:
		for _, variant := range ty.Union.Variants {
			s.markStructsReadFromStack(variant.VariantTyIdx)
		}
	}
}

func (tgen TolkGolangGenerator) GenerateGocode() (declarations string, marshalers string, err error) {
	declarationsBuf := &strings.Builder{}
	marshalersBuf := &strings.Builder{}

	for idx, decl := range tgen.abi.Declarations {
		switch decl.SumType {
		case parser.DeclarationKindStruct:
			if err := tgen.structToGo(decl.StructDeclaration, declarationsBuf, marshalersBuf); err != nil {
				return "", "", fmt.Errorf("declaration[%d] struct %q: %w", idx, decl.StructDeclaration.Name, err)
			}
		case parser.DeclarationKindAlias:
			if err := tgen.aliasToGo(decl.AliasDeclaration, declarationsBuf, marshalersBuf); err != nil {
				return "", "", fmt.Errorf("declaration[%d] alias %q: %w", idx, decl.AliasDeclaration.Name, err)
			}
		case parser.DeclarationKindEnum:
			if err := tgen.enumToGo(decl.EnumDeclaration, declarationsBuf, marshalersBuf); err != nil {
				return "", "", fmt.Errorf("declaration[%d] enum %q: %w", idx, decl.EnumDeclaration.Name, err)
			}
		default:
			return "", "", fmt.Errorf("unexpected kind %v for top-level declaration", decl.SumType)
		}
	}

	if len(tgen.abi.ThrownErrors) > 0 {
		fmt.Fprintf(declarationsBuf, "const ( // errors\n")
		for _, t := range tgen.abi.ThrownErrors {
			if t.Name != "" {
				fmt.Fprintf(declarationsBuf, "\t%s = 0x%X  // %d\n", safeErrorIdent(t.Name), t.ErrCode, t.ErrCode)
			}
		}
		fmt.Fprintf(declarationsBuf, ")\n")
	}

	return declarationsBuf.String(), marshalersBuf.String(), nil
}

func (tgen TolkGolangGenerator) aliasToGo(decl parser.ABIAlias, out *strings.Builder, outMarshal *strings.Builder) error {
	targetTy, err := tgen.symbols.tyByIdx(decl.TargetTyIdx)
	if err != nil {
		return fmt.Errorf("target type: %w", err)
	}
	if targetTy.SumType == parser.TyKindUnion {
		if err := tgen.aliasUnionToGo(decl, out, outMarshal); err != nil {
			return fmt.Errorf("union alias: %w", err)
		}
		return nil
	}
	aliasName := safeGoIdent(decl.Name)
	if len(decl.TypeParams) > 0 {
		return fmt.Errorf("type params not supported for alias %q", decl.Name)
	}
	targetType, err := tgen.symbols.emitGoType(decl.TargetTyIdx)
	if err != nil {
		return fmt.Errorf("target type: %w", err)
	}
	out.WriteString(fmt.Sprintf("type %s %s\n", aliasName, targetType))

	if decl.CustomPackUnpack.UnpackFromSlice != decl.CustomPackUnpack.PackToBuilder {
		fmt.Printf("WARNING: custom pack/unpack for %s is not symmetric\n", decl.Name)
	}

	if !decl.CustomPackUnpack.UnpackFromSlice {
		expr, _, err := tgen.symbols.emitLoadExpr(decl.Name, decl.TargetTyIdx)
		if err != nil {
			return fmt.Errorf("emit unmarshal expression for %q: %w", decl.Name, err)
		}
		fmt.Fprintf(outMarshal, `func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	vx, err := %s
	if err != nil {
		return err
	}
	*v = %s(vx)
	return nil
}
`,
			aliasName, expr, aliasName)
	}
	if !decl.CustomPackUnpack.PackToBuilder {
		expr, err := tgen.symbols.emitStoreExpr("v", decl.TargetTyIdx)
		if err != nil {
			return fmt.Errorf("emit marshal expression for %q: %w", decl.Name, err)
		}
		fmt.Fprintf(outMarshal, `func (v %s) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return %s
}
`, aliasName, expr)
	}
	return nil
}

func (tgen TolkGolangGenerator) aliasUnionToGo(decl parser.ABIAlias, out *strings.Builder, outMarshal *strings.Builder) error {
	aliasName := safeGoIdent(decl.Name)
	if len(decl.TypeParams) > 0 {
		return fmt.Errorf("type params not supported for alias %q", decl.Name)
	}
	targetTy, err := tgen.symbols.tyByIdx(decl.TargetTyIdx)
	if err != nil {
		return fmt.Errorf("target type: %w", err)
	}
	variants, err := tgen.symbols.createLabelsForUnion(targetTy.Union.Variants, nil)
	if err != nil {
		return fmt.Errorf("create labels for union %q: %w", decl.Name, err)
	}

	// Type declarations go to out.
	out.WriteString(fmt.Sprintf("type %sKind uint\n", aliasName))
	fmt.Fprintf(out, "const (\n")
	for _, v := range variants {
		fmt.Fprintf(out, "\t%sKind_%s %sKind = %d\n", aliasName, safeGoIdent(v.label), aliasName, v.PrefixNum)
	}
	fmt.Fprintf(out, ")\n")
	out.WriteString(fmt.Sprintf("type %s struct { // tagged union\n", aliasName))
	out.WriteString(fmt.Sprintf("\tSumType %sKind\n", aliasName))
	for _, v := range variants {
		variantType, err := tgen.symbols.emitGoType(v.VariantTyIdx)
		if err != nil {
			return fmt.Errorf("variant %q type: %w", v.label, err)
		}
		out.WriteString(fmt.Sprintf("\t%s *%s\n", safeGoIdent(v.label), variantType))
	}
	out.WriteString("}\n")

	isEither01 := len(variants) == 2 &&
		variants[0].IsPrefixImplicit && variants[1].IsPrefixImplicit &&
		variants[0].PrefixNum == 0 && variants[0].PrefixLen == 1 &&
		variants[1].PrefixNum == 1 && variants[1].PrefixLen == 1

	// UnmarshalTLB goes to outMarshal.
	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(outMarshal, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {\n", aliasName)
		if isEither01 {
			rExpr, _, err := tgen.symbols.emitLoadExpr(aliasName, variants[0].VariantTyIdx)
			if err != nil {
				return fmt.Errorf("load expression for %q right variant: %w", decl.Name, err)
			}
			lExpr, _, err := tgen.symbols.emitLoadExpr(aliasName, variants[1].VariantTyIdx)
			if err != nil {
				return fmt.Errorf("load expression for %q left variant: %w", decl.Name, err)
			}
			fmt.Fprintf(outMarshal, `	isRight, err := c.ReadBit()
	if err != nil {
		return err
	}
	if isRight {
		return %s
	} else {
		return %s
	}
`, rExpr, lExpr)
		} else {
			isUniformPrefix := true
			uniformPrefixLen := variants[0].PrefixLen
			uniformPrefixEatInPlace := variants[0].IsPrefixImplicit
			for _, v := range variants[1:] {
				if v.PrefixLen != uniformPrefixLen || v.IsPrefixImplicit != uniformPrefixEatInPlace {
					isUniformPrefix = false
					break
				}
			}
			if isUniformPrefix {
				prefixReadFn := "ReadUint"
				if !uniformPrefixEatInPlace {
					prefixReadFn = "PickUint"
				}
				fmt.Fprintf(outMarshal, `prefix, err := c.%s(%d)
if err != nil {
	return err
}
v.SumType = %sKind(prefix)
switch v.SumType {
`, prefixReadFn, uniformPrefixLen, aliasName)
				for _, v := range variants {
					fmt.Fprintf(outMarshal, "\tcase %sKind_%s:\n", aliasName, safeGoIdent(v.label))
					fmt.Fprintf(outMarshal, "\t\tv.%s = new(%s)\n", safeGoIdent(v.label), safeGoIdent(v.label))
					fmt.Fprintf(outMarshal, "\t\treturn  v.%s.UnmarshalTLB(c, decoder)\n", safeGoIdent(v.label))
				}
				fmt.Fprintf(outMarshal, `	default:
		return fmt.Errorf("unknown prefix: %%x", prefix)
	}
`)
			} else {
				for _, v := range variants {
					prefixReadFn := "ReadUint"
					if !v.IsPrefixImplicit {
						prefixReadFn = "PickUint"
					}
					fmt.Fprintf(outMarshal, `if prefix, err := c.%s(%d); err != nil {
	return err
} else if prefix == uint64(%sKind_%s) {
	v.SumType = %sKind(prefix)
	v.%s = new(%s)
	return  v.%s.UnmarshalTLB(c, decoder)
}
`, prefixReadFn, v.PrefixLen, aliasName, safeGoIdent(v.label), aliasName, safeGoIdent(v.label), safeGoIdent(v.label), safeGoIdent(v.label))
				}
				fmt.Fprintf(outMarshal, "\treturn fmt.Errorf(\"could not find suitable prefix\")\n")
			}
			fmt.Fprintf(outMarshal, "}\n")
		}
	}

	// MarshalTLB goes to outMarshal.
	if !decl.CustomPackUnpack.PackToBuilder {
		fmt.Fprintf(outMarshal, "func (v %s) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {\n", aliasName)
		if isEither01 {
			// For Either0/1: variants[0] = bit 0, variants[1] = bit 1.
			v0, v1 := variants[0], variants[1]
			fmt.Fprintf(outMarshal, "\tswitch v.SumType {\n")
			fmt.Fprintf(outMarshal, "\tcase %sKind_%s:\n", aliasName, safeGoIdent(v0.label))
			fmt.Fprintf(outMarshal, "\t\tif err := c.WriteBit(false); err != nil { return err }\n")
			if v0.label != "" {
				fmt.Fprintf(outMarshal, "\t\tif v.%s == nil { return fmt.Errorf(\"%s.%s is nil\") }\n", v0.label, aliasName, v0.label)
				fmt.Fprintf(outMarshal, "\t\treturn v.%s.MarshalTLB(c, encoder)\n", v0.label)
			} else {
				fmt.Fprintf(outMarshal, "\t\treturn nil\n")
			}
			fmt.Fprintf(outMarshal, "\tcase %sKind_%s:\n", aliasName, safeGoIdent(v1.label))
			fmt.Fprintf(outMarshal, "\t\tif err := c.WriteBit(true); err != nil { return err }\n")
			if v1.label != "" {
				fmt.Fprintf(outMarshal, "\t\tif v.%s == nil { return fmt.Errorf(\"%s.%s is nil\") }\n", v1.label, aliasName, v1.label)
				fmt.Fprintf(outMarshal, "\t\treturn v.%s.MarshalTLB(c, encoder)\n", v1.label)
			} else {
				fmt.Fprintf(outMarshal, "\t\treturn nil\n")
			}
			fmt.Fprintf(outMarshal, "\tdefault:\n")
			fmt.Fprintf(outMarshal, "\t\treturn fmt.Errorf(\"unknown %s variant: %%v\", v.SumType)\n", aliasName)
			fmt.Fprintf(outMarshal, "\t}\n")
		} else {
			isUniformPrefix := true
			uniformPrefixEatInPlace := variants[0].IsPrefixImplicit
			uniformPrefixLen := variants[0].PrefixLen
			for _, v := range variants[1:] {
				if v.PrefixLen != uniformPrefixLen || v.IsPrefixImplicit != uniformPrefixEatInPlace {
					isUniformPrefix = false
					break
				}
			}
			fmt.Fprintf(outMarshal, "\tswitch v.SumType {\n")
			for _, v := range variants {
				fmt.Fprintf(outMarshal, "\tcase %sKind_%s:\n", aliasName, safeGoIdent(v.label))
				// When PrefixEatInPlace=true the union owns the prefix; write it before delegating.
				eatInPlace := (isUniformPrefix && uniformPrefixEatInPlace) || (!isUniformPrefix && v.IsPrefixImplicit)
				if eatInPlace {
					fmt.Fprintf(outMarshal, "\t\tif err := c.WriteUint(uint64(%sKind_%s), %d); err != nil { return err }\n",
						aliasName, safeGoIdent(v.label), v.PrefixLen)
				}
				// Delegate to variant's MarshalTLB (for PrefixEatInPlace=false it writes its own prefix).
				if v.label != "" {
					fmt.Fprintf(outMarshal, "\t\tif v.%s == nil { return fmt.Errorf(\"%s.%s is nil\") }\n", safeGoIdent(v.label), aliasName, safeGoIdent(v.label))
					fmt.Fprintf(outMarshal, "\t\treturn v.%s.MarshalTLB(c, encoder)\n", safeGoIdent(v.label))
				} else {
					fmt.Fprintf(outMarshal, "\t\treturn nil\n")
				}
			}
			fmt.Fprintf(outMarshal, "\tdefault:\n")
			fmt.Fprintf(outMarshal, "\t\treturn fmt.Errorf(\"unknown %s variant: %%v\", v.SumType)\n", aliasName)
			fmt.Fprintf(outMarshal, "\t}\n")
		}
		fmt.Fprintf(outMarshal, "}\n")
		emitToCellMethod(outMarshal, aliasName)
	}
	return nil
}

type unionVariantLabeled struct {
	label         string
	hasValueField bool
	parser.UnionVariant
}

func (st *symTable) createLabelsForUnion(variants []parser.UnionVariant, uLabelTyIdx *int) (result []unionVariantLabeled, err error) {
	var genericVariantsTyIdx []int
	if uLabelTyIdx != nil {
		labelTy, err := st.tyByIdx(*uLabelTyIdx)
		if err != nil {
			return nil, err
		}
		if labelTy.SumType == parser.TyKindUnion && len(labelTy.Union.Variants) == len(variants) {
			for _, v := range labelTy.Union.Variants {
				genericVariantsTyIdx = append(genericVariantsTyIdx, v.VariantTyIdx)
			}
		}
	}
	unique := make(map[string]struct{})
	hasDuplicates := false
	for i, variant := range variants {
		labelTyIdx := variant.VariantTyIdx
		if genericVariantsTyIdx != nil {
			labelTyIdx = genericVariantsTyIdx[i]
		}
		label, err := st.createLabelByIdx(labelTyIdx)
		if err != nil {
			return nil, err
		}
		if _, ok := unique[label]; ok {
			hasDuplicates = true
		} else {
			unique[label] = struct{}{}
		}
	}
	for i, variant := range variants {
		labelTyIdx := variant.VariantTyIdx
		if genericVariantsTyIdx != nil {
			labelTyIdx = genericVariantsTyIdx[i]
		}
		variantTy, err := st.tyByIdx(labelTyIdx)
		if err != nil {
			return nil, err
		}
		if variantTy.SumType == parser.TyKindNullLiteral {
			result = append(result, unionVariantLabeled{label: "", hasValueField: false, UnionVariant: variant})
		} else {
			labelVariant := unionVariantLabeled{UnionVariant: variant}
			if hasDuplicates {
				labelVariant.label, err = st.renderTy(labelTyIdx)
				if err != nil {
					return nil, err
				}
				labelVariant.hasValueField = true
			} else {
				labelVariant.label, err = st.createLabelByIdx(labelTyIdx)
				if err != nil {
					return nil, err
				}
				hasOwnLabel, err := st.isStructWithItsOwnLabel(labelTyIdx)
				if err != nil {
					return nil, err
				}
				labelVariant.hasValueField = !hasOwnLabel
			}
			result = append(result, labelVariant)
		}
	}
	return result, nil
}

func (st *symTable) createLabelByIdx(tyIdx int) (string, error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindStructRef:
		return safeGoIdent(ty.StructRef.StructName), nil
	case parser.TyKindAliasRef:
		targetTyIdx, _, err := st.aliasTargetOf(tyIdx)
		if err != nil {
			return "", err
		}
		return st.createLabelByIdx(targetTyIdx)
	case parser.TyKindEnumRef:
		return safeGoIdent(ty.EnumRef.EnumName), nil
	default:
		s, err := st.renderTy(tyIdx)
		if err != nil {
			return "", err
		}
		return safeGoIdent(s), nil
	}
}

func (st *symTable) isStructWithItsOwnLabel(tyIdx int) (bool, error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return false, err
	}
	switch ty.SumType {
	case parser.TyKindStructRef:
		return true, nil
	case parser.TyKindAliasRef:
		aliasTarget, _, err := st.aliasTargetOf(tyIdx)
		if err != nil {
			return false, err
		}
		return st.isStructWithItsOwnLabel(aliasTarget)
	default:
		return false, nil
	}
}

func (st *symTable) getAliasTarget(name string) (int, error) {
	alias, err := st.getAlias(name)
	if err != nil {
		return 0, err
	}
	return alias.TargetTyIdx, nil
}

func (st *symTable) getAlias(name string) (parser.ABIAlias, error) {
	if a, ok := st.aliases[name]; ok {
		return a, nil
	}
	return parser.ABIAlias{}, fmt.Errorf("alias %s not found", name)
}

// conver Tolk `enum X` to golang `type X = underlying type` with `const X1 = value` definitions
func (tgen TolkGolangGenerator) enumToGo(decl parser.ABIEnum, out *strings.Builder, outMarshal *strings.Builder) error {
	typeIdent := safeGoIdent(decl.Name)
	encodedType, err := tgen.symbols.emitGoType(decl.EncodedAsTyIdx)
	if err != nil {
		return fmt.Errorf("enum %q encoded type: %w", decl.Name, err)
	}
	fmt.Fprintf(out, "type %s %s\n", typeIdent, encodedType)
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "const (\n")
	for _, v := range decl.Members {
		name := enumItemIdent(typeIdent, v.Name)
		// handle toCell/store/fromSlice collision
		fmt.Fprintf(out, "\t%s %s = %s\n", name, typeIdent, v.Value.String())
	}
	fmt.Fprintf(out, ")\n")

	// otherwise the user should provide implementation in the same package
	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(outMarshal, "\nfunc (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {\n\treturn (*%s)(v).UnmarshalTLB(c, decoder)\n}\n",
			typeIdent, encodedType)
	}
	if !decl.CustomPackUnpack.PackToBuilder {
		expr, err := tgen.symbols.emitStoreExpr(fmt.Sprintf("%s(v)", encodedType), decl.EncodedAsTyIdx)
		if err != nil {
			return fmt.Errorf("enum %q store expression: %w", decl.Name, err)
		}
		fmt.Fprintf(outMarshal, "\nfunc (v %s) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {\n\treturn %s\n}\n", typeIdent, expr)
	}
	if tgen.symbols.enumNeedsReadFromStack[decl.Name] {
		fmt.Fprintf(outMarshal, "\nfunc (v *%s) ReadFromStack(stack *tlb.VmStack) error {\n\treturn (*%s)(v).ReadFromStack(stack)\n}\n",
			typeIdent, encodedType)
	}

	return nil
}

func (tgen TolkGolangGenerator) structToGo(decl parser.ABIStruct, out *strings.Builder, outMarshal *strings.Builder) error {
	if len(decl.TypeParams) > 0 {
		return fmt.Errorf("type params not supported for struct %q", decl.Name)
	}
	typeIdent := safeGoIdent(decl.Name)
	if decl.Prefix != nil {
		prefix, err := prefixConstValue(decl.Prefix)
		if err != nil {
			return fmt.Errorf("struct %q prefix: %w", decl.Name, err)
		}
		fmt.Fprintf(out, "const Prefix%s uint64 = %s\n", typeIdent, prefix)
	}
	fmt.Fprintf(out, "type %s struct {\n", typeIdent)
	for _, field := range decl.Fields {
		fieldType, err := tgen.symbols.emitGoType(field.TyIdx)
		if err != nil {
			return fmt.Errorf("struct %q field %q type: %w", decl.Name, field.Name, err)
		}
		renderedTy, err := tgen.symbols.renderTy(field.TyIdx)
		if err != nil {
			return fmt.Errorf("struct %q field %q render type: %w", decl.Name, field.Name, err)
		}
		fmt.Fprintf(out, "\t%s %s // %s\n", safePublicField(field.Name), fieldType, renderedTy)
	}
	fmt.Fprintf(out, "}\n")

	if decl.CustomPackUnpack.PackToBuilder != decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Printf("WARNING: custom pack/unpack for %s is not symmetric\n", decl.Name)
	}

	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(outMarshal, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {\n", typeIdent)
		if decl.Prefix != nil {
			if decl.Prefix.PrefixLen > 64 {
				return fmt.Errorf("struct %q prefix length too large: %d", decl.Name, decl.Prefix.PrefixLen)
			}
			fmt.Fprintf(outMarshal, "\tif err := c.ReadPrefix(%d, Prefix%s); err != nil {\n\t\treturn err\n\t}\n", decl.Prefix.PrefixLen, typeIdent)
		}
		for _, field := range decl.Fields {
			fieldPath := fmt.Sprintf("%s.%s", typeIdent, safePublicField(field.Name))
			expr, hasLoadMethod, err := tgen.symbols.emitLoadExpr(fieldPath, field.TyIdx)
			if err != nil {
				return fmt.Errorf("struct %q field %q unmarshal expression: %w", decl.Name, field.Name, err)
			}
			if hasLoadMethod {
				fmt.Fprintf(outMarshal, "\tif err = v.%s.UnmarshalTLB(c, decoder); err != nil {\n", safePublicField(field.Name))
			} else {
				fmt.Fprintf(outMarshal, "\tif v.%s, err = %s; err != nil {\n", safePublicField(field.Name), expr)
			}
			fmt.Fprintf(outMarshal, "\t\treturn fmt.Errorf(\"failed to read .%s: %%v\", err)\n", safePublicField(field.Name))
			fmt.Fprintf(outMarshal, "\t}\n")
		}
		fmt.Fprintf(outMarshal, "\treturn nil\n}\n")
	}

	if !decl.CustomPackUnpack.PackToBuilder {
		fmt.Fprintf(outMarshal, `func (v %s) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {`, typeIdent)
		if decl.Prefix != nil {
			if decl.Prefix.PrefixLen > 64 {
				return fmt.Errorf("struct %q prefix length too large: %d", decl.Name, decl.Prefix.PrefixLen)
			}
			fmt.Fprintf(outMarshal, `if err = c.WriteUint(Prefix%s, %d); err != nil {
	return fmt.Errorf("failed to write prefix: %%v", err)
}
`, typeIdent, decl.Prefix.PrefixLen)
		}
		for _, field := range decl.Fields {
			pFieldName := safePublicField(field.Name)
			fieldExpr := fmt.Sprintf("v.%s", pFieldName)
			if expr, err := tgen.symbols.emitStoreExpr(fieldExpr, field.TyIdx); err == nil {
				fmt.Fprintf(outMarshal, `if err = %s; err != nil {
	return fmt.Errorf("failed to .%s: %%v", err)
}
`, expr, pFieldName)
			} else {
				return fmt.Errorf("struct %q field %q store expression: %w", decl.Name, field.Name, err)
			}
		}
		fmt.Fprintf(outMarshal, "\treturn nil\n}\n")
		emitToCellMethod(outMarshal, typeIdent)
	}

	if tgen.symbols.structIsReturnedFromGetMethod[decl.Name] {
		// generate stack unrolling
		if decl.Prefix != nil {
			return fmt.Errorf("prefix not supported for get-method stack unmarshalling in struct %q", decl.Name)
		}
		fmt.Fprintf(outMarshal, "func (v *%s) ReadFromStack(stack *tlb.VmStack) (err error) {\n", typeIdent)
		fields := slices.Clone(decl.Fields)
		slices.Reverse(fields)
		for _, field := range fields {
			publicField := safePublicField(field.Name)
			expr, hasMethod, err := tgen.emitStackReadExpr(publicField, field.TyIdx, false)
			if err != nil {
				return fmt.Errorf("struct %q field %q stack expression: %w", decl.Name, field.Name, err)
			}
			if !hasMethod {
				fmt.Fprintf(outMarshal, "\tif v.%s, err = %s; err != nil {\n", publicField, expr)
			} else {
				fmt.Fprintf(outMarshal, "\tif err = v.%s.ReadFromStack(stack); err != nil {\n", publicField)
			}
			fmt.Fprintf(outMarshal, "\t\treturn fmt.Errorf(\"failed to read .%s: %%v\", err)\n", publicField)
			fmt.Fprintf(outMarshal, "\t}\n")
		}
		fmt.Fprintf(outMarshal, "\treturn nil\n}\n")
	}

	return nil
}

func emitToCellMethod(out *strings.Builder, typeName string) {
	fmt.Fprintf(out, `func (v %s) ToCell() (*boc.Cell, error) {
		c := boc.NewCell()
		if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
			return nil, err
		}
		return c, nil
	}
`, typeName)
}
