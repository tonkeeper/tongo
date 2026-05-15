package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldView struct {
	Field
	ULabelTyIdx *int
}

type ABIIndex struct {
	UniqueTypes          []Ty
	StructInstantiations []ABIStructInstantiation
	AliasInstantiations  []ABIAliasInstantiation
	Structs              map[string]ABIStruct
	Aliases              map[string]ABIAlias
	Enums                map[string]ABIEnum
}

func NewABIIndex(abi ContractABI) *ABIIndex {
	idx := &ABIIndex{
		UniqueTypes:          abi.UniqueTypes,
		StructInstantiations: abi.StructInstantiations,
		AliasInstantiations:  abi.AliasInstantiations,
		Structs:              make(map[string]ABIStruct),
		Aliases:              make(map[string]ABIAlias),
		Enums:                make(map[string]ABIEnum),
	}
	for _, decl := range abi.Declarations {
		switch decl.SumType {
		case DeclarationKindStruct:
			idx.Structs[decl.StructDeclaration.Name] = decl.StructDeclaration
		case DeclarationKindAlias:
			idx.Aliases[decl.AliasDeclaration.Name] = decl.AliasDeclaration
		case DeclarationKindEnum:
			idx.Enums[decl.EnumDeclaration.Name] = decl.EnumDeclaration
		}
	}
	return idx
}

func (idx *ABIIndex) TyByIdx(tyIdx int) (Ty, error) {
	if tyIdx < 0 || tyIdx >= len(idx.UniqueTypes) {
		return Ty{}, fmt.Errorf("ty_idx out of range: %d", tyIdx)
	}
	return idx.UniqueTypes[tyIdx], nil
}

func (idx *ABIIndex) TyIdxOf(ty Ty) (int, bool) {
	for i, candidate := range idx.UniqueTypes {
		if reflect.DeepEqual(candidate, ty) {
			return i, true
		}
	}
	return 0, false
}

func (idx *ABIIndex) StructFieldsOf(tyIdx int, isForStack bool) ([]FieldView, error) {
	ty, err := idx.TyByIdx(tyIdx)
	if err != nil {
		return nil, err
	}
	if ty.SumType != TyKindStructRef {
		return nil, fmt.Errorf("expected StructRef at ty_idx=%d", tyIdx)
	}
	decl, ok := idx.Structs[ty.StructRef.StructName]
	if !ok {
		return nil, fmt.Errorf("struct %q not found", ty.StructRef.StructName)
	}

	fields := make([]FieldView, len(decl.Fields))
	for i, f := range decl.Fields {
		fields[i] = FieldView{Field: f}
	}
	for _, inst := range idx.StructInstantiations {
		if inst.TyIdx != tyIdx {
			continue
		}
		if len(fields) != len(inst.MonomorphicFieldsTyIdx) {
			return nil, fmt.Errorf("mismatch monomorphic fields size for %q", inst.StructName)
		}
		for i := range fields {
			old := fields[i].TyIdx
			fields[i].TyIdx = inst.MonomorphicFieldsTyIdx[i]
			fields[i].ULabelTyIdx = &old
		}
		break
	}
	if !isForStack {
		for i := range fields {
			if fields[i].ClientTyIdx != nil {
				fields[i].TyIdx = *fields[i].ClientTyIdx
				fields[i].ULabelTyIdx = nil
			}
		}
	}
	return fields, nil
}

func (idx *ABIIndex) AliasTargetOf(tyIdx int) (targetTyIdx int, uLabelTyIdx *int, err error) {
	ty, err := idx.TyByIdx(tyIdx)
	if err != nil {
		return 0, nil, err
	}
	if ty.SumType != TyKindAliasRef {
		return 0, nil, fmt.Errorf("expected AliasRef at ty_idx=%d", tyIdx)
	}
	decl, err := idx.GetAlias(ty.AliasRef.AliasName)
	if err != nil {
		return 0, nil, err
	}
	targetTyIdx = decl.TargetTyIdx
	for _, inst := range idx.AliasInstantiations {
		if inst.TyIdx != tyIdx {
			continue
		}
		old := targetTyIdx
		return inst.MonomorphicTargetTyIdx, &old, nil
	}
	return targetTyIdx, nil, nil
}

func (idx *ABIIndex) MsgName(bodyTyIdx int) (string, error) {
	ty, err := idx.TyByIdx(bodyTyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case TyKindStructRef:
		return ty.StructRef.StructName, nil
	case TyKindAliasRef:
		return ty.AliasRef.AliasName, nil
	default:
		return "", fmt.Errorf("cannot get name for %q body", ty.SumType)
	}
}

func (idx *ABIIndex) GetAlias(name string) (ABIAlias, error) {
	if a, ok := idx.Aliases[name]; ok {
		return a, nil
	}
	return ABIAlias{}, fmt.Errorf("alias %s not found", name)
}

func (idx *ABIIndex) GetAliasTarget(name string) (int, error) {
	alias, err := idx.GetAlias(name)
	if err != nil {
		return 0, err
	}
	return alias.TargetTyIdx, nil
}

func (idx *ABIIndex) RenderTy(tyIdx int) (string, error) {
	ty, err := idx.TyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case TyKindInt:
		return "int", nil
	case TyKindIntN:
		return fmt.Sprintf("int%d", ty.IntN.N), nil
	case TyKindUintN:
		return fmt.Sprintf("uint%d", ty.UintN.N), nil
	case TyKindVarIntN:
		return fmt.Sprintf("varint%d", ty.VarIntN.N), nil
	case TyKindVarUintN:
		return fmt.Sprintf("varuint%d", ty.VarUintN.N), nil
	case TyKindCoins:
		return "coins", nil
	case TyKindBool:
		return "bool", nil
	case TyKindCell:
		return "cell", nil
	case TyKindBuilder:
		return "builder", nil
	case TyKindSlice:
		return "slice", nil
	case TyKindString:
		return "string", nil
	case TyKindRemaining:
		return "RemainingBitsAndRefs", nil
	case TyKindAddress:
		return "address", nil
	case TyKindAddressOpt:
		return "address?", nil
	case TyKindAddressExt:
		return "ext_address", nil
	case TyKindAddressAny:
		return "any_address", nil
	case TyKindBitsN:
		return fmt.Sprintf("bits%d", ty.BitsN.N), nil
	case TyKindNullLiteral:
		return "null", nil
	case TyKindCallable:
		return "continuation", nil
	case TyKindVoid:
		return "void", nil
	case TyKindUnknown:
		return "unknown", nil
	case TyKindNullable:
		inner, err := idx.RenderTy(ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("nullable inner: %w", err)
		}
		return inner + "?", nil
	case TyKindCellOf:
		inner, err := idx.RenderTy(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("cell inner: %w", err)
		}
		return fmt.Sprintf("Cell<%s>", inner), nil
	case TyKindArrayOf:
		inner, err := idx.RenderTy(ty.ArrayOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("array inner: %w", err)
		}
		return fmt.Sprintf("array<%s>", inner), nil
	case TyKindLispListOf:
		inner, err := idx.RenderTy(ty.LispListOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("lisp list inner: %w", err)
		}
		return fmt.Sprintf("lisp_list<%s>", inner), nil
	case TyKindTensor:
		items, err := mapTyIdx(ty.Tensor.ItemsTyIdx, idx.RenderTy)
		if err != nil {
			return "", fmt.Errorf("tensor items: %w", err)
		}
		return fmt.Sprintf("(%s)", strings.Join(items, ", ")), nil
	case TyKindShapedTuple:
		items, err := mapTyIdx(ty.ShapedTuple.ItemsTyIdx, idx.RenderTy)
		if err != nil {
			return "", fmt.Errorf("tuple items: %w", err)
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", ")), nil
	case TyKindMapKV:
		key, err := idx.RenderTy(ty.MapKV.KeyTyIdx)
		if err != nil {
			return "", fmt.Errorf("map key: %w", err)
		}
		value, err := idx.RenderTy(ty.MapKV.ValueTyIdx)
		if err != nil {
			return "", fmt.Errorf("map value: %w", err)
		}
		return fmt.Sprintf("map<%s, %s>", key, value), nil
	case TyKindEnumRef:
		return ty.EnumRef.EnumName, nil
	case TyKindStructRef:
		args, err := idx.renderTypeArgs(ty.StructRef.TypeArgsTyIdx)
		if err != nil {
			return "", fmt.Errorf("struct type args: %w", err)
		}
		return ty.StructRef.StructName + args, nil
	case TyKindAliasRef:
		args, err := idx.renderTypeArgs(ty.AliasRef.TypeArgsTyIdx)
		if err != nil {
			return "", fmt.Errorf("alias type args: %w", err)
		}
		return ty.AliasRef.AliasName + args, nil
	case TyKindGenericT:
		return ty.GenericT.NameT, nil
	case TyKindUnion:
		items := make([]string, len(ty.Union.Variants))
		for i, v := range ty.Union.Variants {
			rendered, err := idx.RenderTy(v.VariantTyIdx)
			if err != nil {
				return "", fmt.Errorf("union variants: %w", err)
			}
			items[i] = rendered
		}
		return strings.Join(items, " | "), nil
	default:
		return "", fmt.Errorf("unknown type kind %q", ty.SumType)
	}
}

func (idx *ABIIndex) renderTypeArgs(args []int) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	rendered, err := mapTyIdx(args, idx.RenderTy)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("<%s>", strings.Join(rendered, ", ")), nil
}

func mapTyIdx(values []int, f func(int) (string, error)) ([]string, error) {
	res := make([]string, len(values))
	for i, value := range values {
		mapped, err := f(value)
		if err != nil {
			return nil, err
		}
		res[i] = mapped
	}
	return res, nil
}
