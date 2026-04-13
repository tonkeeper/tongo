package tolkgen

import (
	"fmt"
	"strings"

	//"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type TolkGolangGenerator struct {
	symbols *symTable
	abi     parser.ABI
}

func NewTolkGolangGenerator(abi parser.ABI) *TolkGolangGenerator {
	symbols := symTable{
		aliases: make(map[string]parser.AliasDeclaration),
	}
	for _, decl := range abi.Declarations {
		if decl.SumType == parser.DeclarationKindAlias {
			symbols.aliases[decl.AliasDeclaration.Name] = decl.AliasDeclaration
		}
	}
	return &TolkGolangGenerator{
		abi:     abi,
		symbols: &symbols,
	}
}

type symTable struct {
	aliases map[string]parser.AliasDeclaration
}

func (tgen TolkGolangGenerator) GenerateGocode() (string, error) {
	declarationsBuf := &strings.Builder{}
	marshalersBuf := &strings.Builder{}

	for _, decl := range tgen.abi.Declarations {
		switch decl.SumType {
		case parser.DeclarationKindStruct:
			tgen.structToGo(decl.StructDeclaration, declarationsBuf, marshalersBuf)
			break
		case parser.DeclarationKindAlias:
			tgen.aliasToGo(decl.AliasDeclaration, declarationsBuf, marshalersBuf)
			break
		case parser.DeclarationKindEnum:
			tgen.enumToGo(decl.EnumDeclaration, declarationsBuf, marshalersBuf)
			break
		default:
			return "", fmt.Errorf("unexpect kind %v for top-level declaration", decl.SumType)
		}
	}

	if len(tgen.abi.ThrownErrors) > 0 {
		fmt.Fprintf(declarationsBuf, "const ( // errors\n")
		for _, t := range tgen.abi.ThrownErrors {
			if t.Name != "" {
				fmt.Fprintf(declarationsBuf, "\t%s = 0x%X  // %d\n", safeGoIdent(t.Name), t.ErrCode, t.ErrCode)
			}
		}
		fmt.Fprintf(declarationsBuf, ")\n")
	}

	return declarationsBuf.String() + "\n" + marshalersBuf.String(), nil
}

func (tgen TolkGolangGenerator) aliasToGo(decl parser.AliasDeclaration, out *strings.Builder, outMarshal *strings.Builder) {
	if decl.TargetTy.SumType == parser.TyKindUnion {
		tgen.aliasUnionToGo(decl, out)
		return
	}
	aliasName := safeGoIdent(decl.Name)
	if len(decl.TypeParams) > 0 {
		panic("type params not supported")
	}
	out.WriteString(fmt.Sprintf("type %s %s\n", aliasName, emitGoType(decl.TargetTy)))

	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(outMarshal, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {\n", aliasName)
		if expr, _, err := tgen.symbols.emitLoadExpr(decl.Name, decl.TargetTy); err == nil {
			fmt.Fprintf(outMarshal, "\tvx, err := %s\n", expr)
			fmt.Fprintf(outMarshal, "\tif err != nil {\n")
			fmt.Fprintf(outMarshal, "\t\treturn err\n")
			fmt.Fprintf(outMarshal, "\t}\n")
			fmt.Fprintf(outMarshal, "\t*v = %s(vx)\n", aliasName)
			fmt.Fprintf(outMarshal, "\treturn nil\n")
		} else {
			panic(err)
		}
		fmt.Fprintf(outMarshal, "}\n")
	}
}

func (tgen TolkGolangGenerator) aliasUnionToGo(decl parser.AliasDeclaration, out *strings.Builder) {
	aliasName := safeGoIdent(decl.Name)
	if len(decl.TypeParams) > 0 {
		panic("type params not supported")
	}
	out.WriteString(fmt.Sprintf("type %sKind uint\n", aliasName))
	fmt.Fprintf(out, "const (\n")
	for _, v := range tgen.symbols.createLabelsForUnion(decl.TargetTy.Union.Variants) {
		fmt.Fprintf(out, "\t%sKind_%s %sKind = %s\n", aliasName, safeGoIdent(v.label), aliasName, v.PrefixStr)
	}
	fmt.Fprintf(out, ")\n")
	out.WriteString(fmt.Sprintf("type %s struct { // tagged union\n", aliasName))
	out.WriteString(fmt.Sprintf("\tSumType %sKind\n", aliasName))
	for _, v := range tgen.symbols.createLabelsForUnion(decl.TargetTy.Union.Variants) {
		out.WriteString(fmt.Sprintf("\t%s *%s\n", safeGoIdent(v.label), emitGoType(v.VariantTy)))
	}
	out.WriteString("}\n")

	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(out, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {\n", aliasName)
		variants := tgen.symbols.createLabelsForUnion(decl.TargetTy.Union.Variants)
		isEither01 := len(variants) == 2 && variants[0].PrefixStr == "0b0" && variants[1].PrefixStr == "0b1"
		if isEither01 {
			if rExpr, _, err := tgen.symbols.emitLoadExpr(aliasName, variants[0].VariantTy); err == nil {
				if lExpr, _, err := tgen.symbols.emitLoadExpr(aliasName, variants[1].VariantTy); err == nil {
					fmt.Fprintf(out, "\tisRight, err := c.ReadBit()\n")
					fmt.Fprintf(out, "\tif err != nil {\n")
					fmt.Fprintf(out, "\t\treturn err\n")
					fmt.Fprintf(out, "\t}\n")
					fmt.Fprintf(out, "\tif isRight {\n")
					fmt.Fprintf(out, "\t\treturn %s\n", rExpr)
					fmt.Fprintf(out, "\t}\n")
					fmt.Fprintf(out, "\treturn %s\n", lExpr)
					fmt.Fprintf(out, "}\n")
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			isUniformPrefix := true
			uniformPrefixLen := variants[0].PrefixLen
			uniformPrefixEatInPlace := variants[0].PrefixEatInPlace
			for _, v := range variants[1:] {
				if v.PrefixLen != uniformPrefixLen || v.PrefixEatInPlace != uniformPrefixEatInPlace {
					isUniformPrefix = false
					break
				}
			}
			if isUniformPrefix {
				if uniformPrefixEatInPlace {
					fmt.Fprintf(out, "prefix, err := c.ReadUint(%d)\n", uniformPrefixLen)
				} else {
					fmt.Fprintf(out, "prefix, err := c.PickUint(%d)\n", uniformPrefixLen)
				}
				fmt.Fprintf(out, "if err != nil {\n")
				fmt.Fprintf(out, "\treturn err\n")
				fmt.Fprintf(out, "}\n")
				fmt.Fprintf(out, "v.SumType = %sKind(prefix)\n", aliasName)
				fmt.Fprintf(out, "switch v.SumType {\n")
				for _, v := range variants {
					fmt.Fprintf(out, "\tcase %sKind_%s:\n", aliasName, safeGoIdent(v.label))
					fmt.Fprintf(out, "\t\tv.%s = new(%s)\n", v.label, v.label)
					fmt.Fprintf(out, "\t\treturn  v.%s.UnmarshalTLB(c, decoder)\n", v.label)
				}
				fmt.Fprintf(out, "\tdefault:\n")
				fmt.Fprintf(out, "\t\treturn fmt.Errorf(\"unknown prefix: %%x\", prefix)\n")
				fmt.Fprintf(out, "\t}\n")
			} else {
				for _, v := range variants {
					if v.PrefixEatInPlace {
						fmt.Fprintf(out, "if prefix, err := c.ReadUint(%d); err != nil {\n", v.PrefixLen)
					} else {
						fmt.Fprintf(out, "if prefix, err := c.PickUint(%d); err != nil {\n", v.PrefixLen)
					}
					fmt.Fprintf(out, "\treturn err\n")
					fmt.Fprintf(out, "} else if prefix == uint64(%sKind_%s) {\n", aliasName, safeGoIdent(v.label))
					fmt.Fprintf(out, "\tv.SumType = %sKind(prefix)\n", aliasName)
					fmt.Fprintf(out, "\tv.%s = new(%s)\n", v.label, v.label)
					fmt.Fprintf(out, "\treturn  v.%s.UnmarshalTLB(c, decoder)\n", v.label)
					fmt.Fprintf(out, "\t}\n")
				}
				fmt.Fprintf(out, "\treturn fmt.Errorf(\"could not find suitable prefix\")\n")
			}
		}
		fmt.Fprintf(out, "}\n")
	}
}

type unionVariantLabeled struct {
	label         string
	hasValueField bool
	parser.UnionVariant
}

func (st *symTable) createLabel(ty parser.Ty) string {
	return safeGoIdent(ty.String())
}

func (st *symTable) createLabelsForUnion(variants []parser.UnionVariant) (result []unionVariantLabeled) {
	unique := make(map[string]struct{})
	hasDuplicates := false
	for _, variant := range variants {
		label := st.createLabel(variant.VariantTy)
		if _, ok := unique[label]; ok {
			hasDuplicates = true
		} else {
			unique[label] = struct{}{}
		}
	}
	for _, variant := range variants {
		variantTy := variant.VariantTy
		if variantTy.SumType == parser.TyKindNullLiteral {
			result = append(result, unionVariantLabeled{label: "", hasValueField: false, UnionVariant: variant})
		} else {
			labelVariant := unionVariantLabeled{UnionVariant: variant}
			if hasDuplicates {
				labelVariant.label = variantTy.String()
				labelVariant.hasValueField = true
			} else {
				labelVariant.label = st.createLabel(variantTy)
				labelVariant.hasValueField = !st.isStructWithItsOwnLabel(variantTy)
			}
			result = append(result, labelVariant)
		}
	}
	return result
}

func (st *symTable) isStructWithItsOwnLabel(ty parser.Ty) bool {
	switch ty.SumType {
	case parser.TyKindStructRef:
		return true
	case parser.TyKindAliasRef:
		return st.isStructWithItsOwnLabel(st.getAliasTarget(ty.AliasRef.AliasName))
	default:
		return false
	}
}

func (st *symTable) getAliasTarget(name string) parser.Ty {
	return st.getAlias(name).TargetTy
}

func (st *symTable) getAlias(name string) parser.AliasDeclaration {
	if a, ok := st.aliases[name]; ok {
		return a
	}
	panic(fmt.Sprintf("alias %s not found", name))
}

// conver Tolk `enum X` to golang `type X = underlying type` with `const X1 = value` definitions
func (tgen TolkGolangGenerator) enumToGo(decl parser.EnumDeclaration, out *strings.Builder, outMarshal *strings.Builder) {
	typeIdent := safeGoIdent(decl.Name)
	fmt.Fprintf(out, "type %s %s\n", typeIdent, emitGoType(decl.EncodedAs))
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "const (\n")
	for _, v := range decl.Members {
		name := enumItemIdent(typeIdent, v.Name)
		// handle toCell/store/fromSlice collision
		fmt.Fprintf(out, "\t%s %s = %s\n", name, typeIdent, v.Value)
	}
	fmt.Fprintf(out, ")\n")

	// otherwise the user should provide implementation in the same package
	if !decl.CustomPackUnpack.UnpackFromSlice {
		// UnmarshalTLB(cell *boc.Cell, decoder *tlb.Decoder)
		fmt.Fprintf(outMarshal, "\n")
		fmt.Fprintf(outMarshal, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {\n", typeIdent)
		if expr, _, err := tgen.symbols.emitLoadExpr(decl.Name, decl.EncodedAs); err == nil {
			fmt.Fprintf(outMarshal, "\tvx, err := %s\n", expr)
			fmt.Fprintf(outMarshal, "\t*v = %s(vx)\n", typeIdent)
		} else {
			panic(err)
		}
		fmt.Fprintf(outMarshal, "\treturn err\n")
		fmt.Fprintf(outMarshal, "}\n")
	}
	if !decl.CustomPackUnpack.PackToBuilder {
		// MarshalTLB(c *boc.Cell, encoder *Encoder) error
		fmt.Fprintf(outMarshal, "\n")
		fmt.Fprintf(outMarshal, "func (v %s) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {\n", typeIdent)
		if expr, err := emitStoreExpr(decl.Name, decl.EncodedAs); err == nil {
			fmt.Fprintf(outMarshal, "\treturn %s\n", expr)
		} else {
			panic(err)
		}
		fmt.Fprintf(outMarshal, "}\n")
	}

}

func (tgen TolkGolangGenerator) structToGo(decl parser.StructDeclaration, out *strings.Builder, outMarshal *strings.Builder) {
	if len(decl.TypeParams) > 0 {
		panic("type params not supported")
	}
	typeIdent := safeGoIdent(decl.Name)
	if decl.Prefix != nil {
		fmt.Fprintf(out, "const Prefix%s uint64 = %s\n", typeIdent, decl.Prefix.PrefixStr)
	}
	fmt.Fprintf(out, "type %s struct {\n", typeIdent)
	for _, field := range decl.Fields {
		fmt.Fprintf(out, "\t%s %s\n", safePublicField(field.Name), emitGoType(field.Ty))
	}
	fmt.Fprintf(out, "}\n")

	if !decl.CustomPackUnpack.UnpackFromSlice {
		fmt.Fprintf(outMarshal, "func (v *%s) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {\n", typeIdent)
		if decl.Prefix != nil {
			fmt.Fprintf(outMarshal, "\tif prefix, err := c.ReadUint(%d); err != nil {\n", decl.Prefix.PrefixLen)
			fmt.Fprintf(outMarshal, "\t\treturn err\n")
			fmt.Fprintf(outMarshal, "\t} else if prefix != Prefix%s {\n", typeIdent)
			fmt.Fprintf(outMarshal, "\t\treturn fmt.Errorf(\"unexpected prefix: %%x\", prefix)\n")
			fmt.Fprintf(outMarshal, "\t}\n")
		}
		for _, field := range decl.Fields {
			fieldPath := fmt.Sprintf("%s.%s", typeIdent, safePublicField(field.Name))
			if expr, hasLoadMethod, err := tgen.symbols.emitLoadExpr(fieldPath, field.Ty); err == nil {
				if hasLoadMethod {
					fmt.Fprintf(outMarshal, "\tif err = v.%s.UnmarshalTLB(c, decoder); err != nil {\n", safePublicField(field.Name))
				} else {
					fmt.Fprintf(outMarshal, "\tif v.%s, err = %s; err != nil {\n", safePublicField(field.Name), expr)
				}
				fmt.Fprintf(outMarshal, "\t\treturn err\n")
				fmt.Fprintf(outMarshal, "\t}\n")
			} else {
				panic(err)
			}
		}
		fmt.Fprintf(outMarshal, "\treturn nil\n")
		fmt.Fprintf(outMarshal, "}\n")
	}

}

func emitStoreExpr(name string, ty parser.Ty) (string, error) {
	switch ty.SumType {
	case parser.TyKindUintN:
		return fmt.Sprintf("c.WriteUint(uint64(v), %d)", ty.UintN.N), nil
	case parser.TyKindIntN:
		return fmt.Sprintf("c.WriteInt(int64(v), %d)", ty.IntN.N), nil
	}
	return "", fmt.Errorf("unknown type %v", ty)
}

func (st *symTable) emitLoadExpr(fieldPath string, ty parser.Ty) (expr string, hasLoadMethod bool, err error) {

	switch ty.SumType {
	case parser.TyKindBuilder, parser.TyKindSlice, parser.TyKindUnknown, parser.TyKindCallable:
		var hint string
		if ty.SumType == parser.TyKindBuilder || ty.SumType == parser.TyKindSlice {
			hint = " (it can be used for writing only)"
		}
		return "", false, fmt.Errorf("%s is %s%s", fieldPath, ty.String(), hint)
	case parser.TyKindRemaining:
		return `(func () (tlb.Any, error) {
	cc := c.CopyRemaining()
	return tlb.Any(*cc), nil
})()`, false, nil
	case parser.TyKindBitsN:
		return fmt.Sprintf("tlb.UnmarshalT[tlb.Bits%d](c, decoder)", ty.BitsN.N), true, nil
	case parser.TyKindCell:
		return "c.NextRefV()", false, nil
	case parser.TyKindCellOf:
		if ty.CellOf.Inner.SumType == parser.TyKindSlice {
			return "c.NextRefV()", false, nil
		}
		return fmt.Sprintf("tlb.UnmarshalT[tlb.RefT[*%s]](c, decoder)", emitGoType(ty.CellOf.Inner)), true, nil
	case parser.TyKindAliasRef:
		if len(ty.AliasRef.TypeArgs) > 0 {
			return "", false, fmt.Errorf("type arguments not supported for aliases")
		}
		return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", ty.AliasRef.AliasName), true, nil
	case parser.TyKindBool:
		return "c.ReadBit()", false, nil
	case parser.TyKindAddress, parser.TyKindAddressAny, parser.TyKindAddressOpt,
		parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindMapKV, parser.TyKindEnumRef, parser.TyKindStructRef:
		return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", emitGoType(ty)), true, nil
	case parser.TyKindNullable:
		return fmt.Sprintf("tlb.UnmarshalT[tlb.Maybe[%s]](c, decoder)", emitGoType(ty.Nullable.Inner)), true, nil
	case parser.TyKindString:
		return "c.ReadStringRefTail()", false, nil
	}

	return "", false, fmt.Errorf("unknown type %v", ty)
}

func safeGoIdent(name string) string {
	return name
}
func safePublicField(name string) string {
	// capitalize the first letter
	return strings.ToUpper(name[:1]) + name[1:]
}

func enumItemIdent(typeIdent, name string) string {
	return typeIdent + name
}

func emitGoType(ty parser.Ty) string {
	switch ty.SumType {
	case parser.TyKindInt:
		return "tlb.Int257"
	case parser.TyKindIntN:
		return fmt.Sprintf("tlb.Int%d", ty.IntN.N)
	case parser.TyKindUintN:
		return fmt.Sprintf("tlb.Uint%d", ty.UintN.N)
	case parser.TyKindVarIntN:
		return fmt.Sprintf("tlb.VarInteger%d", ty.VarIntN.N)
	case parser.TyKindVarUintN:
		return fmt.Sprintf("tlb.VarUInteger%d", ty.VarUintN.N)
	case parser.TyKindBitsN:
		switch ty.BitsN.N {
		case 80, 96, 128, 160, 256, 264, 320, 352, 512:
			return "tlb.Bits" + fmt.Sprintf("%d", ty.BitsN.N)
		}
		panic(fmt.Sprintf("tlb.Bits%d is not supported: update cmd/codegen/integers", ty.BitsN.N))
	case parser.TyKindCoins:
		return "tlb.Coins"
	case parser.TyKindCell:
		return "boc.Cell"
	case parser.TyKindBool:
		return "bool"
	case parser.TyKindCellOf:
		if ty.CellOf.Inner.SumType == parser.TyKindSlice {
			return "boc.Cell"
		}
		return fmt.Sprintf("tlb.RefT[*%s]", emitGoType(ty.CellOf.Inner))
	case parser.TyKindStructRef:
		switch ty.StructRef.StructName {
		case "CurrencyCollection":
			return "tlb.CurrencyCollection"
		default:
			name := ty.StructRef.StructName
			if len(ty.StructRef.TypeArgs) > 0 {
				tArgs := ""
				for _, t := range ty.StructRef.TypeArgs {
					tArgs = tArgs + emitGoType(t) + ", "
				}
				name = name + "[" + tArgs + "]"
			}
			return name
		}
	case parser.TyKindAliasRef:
		name := ty.AliasRef.AliasName
		if len(ty.AliasRef.TypeArgs) > 0 {
			tArgs := ""
			for _, t := range ty.AliasRef.TypeArgs {
				tArgs = tArgs + emitGoType(t) + ", "
			}
			name = name + "[" + tArgs + "]"
		}
		return name
	case parser.TyKindAddress:
		return "tlb.InternalAddress"
	case parser.TyKindAddressAny:
		return "tlb.MsgAddress"
	case parser.TyKindAddressOpt:
		return "tlb.MsgAddress" // upcast type to only carry none and internal cases (not external)
	case parser.TyKindNullable:
		return fmt.Sprintf("tlb.Maybe[%s]", emitGoType(ty.Nullable.Inner))
	case parser.TyKindMapKV:
		return fmt.Sprintf("tlb.Hashmap[%s, %s]", emitGoType(ty.MapKV.K), emitGoType(ty.MapKV.V))
	case parser.TyKindRemaining:
		return "tlb.Any"
	case parser.TyKindEnumRef:
		return safeGoIdent(ty.EnumRef.EnumName)
	case parser.TyKindSlice:
		return "boc.BitString"
	case parser.TyKindArrayOf:
		return fmt.Sprintf("[]%s", emitGoType(ty.ArrayOf.Inner))
	case parser.TyKindString:
		return "string"
	}
	panic(fmt.Sprintf("getGOType type not supported: %v", ty))
}
