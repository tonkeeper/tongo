package tolkgen

import (
	"fmt"
	"strings"

	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

type fieldView struct {
	parser.Field
	ULabelTyIdx *int
}

func (st *symTable) tyByIdx(idx int) (parser.Ty, error) {
	if idx < 0 || idx >= len(st.uniqueTypes) {
		return parser.Ty{}, fmt.Errorf("ty_idx out of range: %d", idx)
	}
	return st.uniqueTypes[idx], nil
}

func (st *symTable) renderTy(idx int) (string, error) {
	ty, err := st.tyByIdx(idx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindInt:
		return "int", nil
	case parser.TyKindIntN:
		return fmt.Sprintf("int%d", ty.IntN.N), nil
	case parser.TyKindUintN:
		return fmt.Sprintf("uint%d", ty.UintN.N), nil
	case parser.TyKindVarIntN:
		return fmt.Sprintf("varint%d", ty.VarIntN.N), nil
	case parser.TyKindVarUintN:
		return fmt.Sprintf("varuint%d", ty.VarUintN.N), nil
	case parser.TyKindCoins:
		return "coins", nil
	case parser.TyKindBool:
		return "bool", nil
	case parser.TyKindCell:
		return "cell", nil
	case parser.TyKindBuilder:
		return "builder", nil
	case parser.TyKindSlice:
		return "slice", nil
	case parser.TyKindString:
		return "string", nil
	case parser.TyKindRemaining:
		return "RemainingBitsAndRefs", nil
	case parser.TyKindAddress:
		return "address", nil
	case parser.TyKindAddressOpt:
		return "address?", nil
	case parser.TyKindAddressExt:
		return "ext_address", nil
	case parser.TyKindAddressAny:
		return "any_address", nil
	case parser.TyKindBitsN:
		return fmt.Sprintf("bits%d", ty.BitsN.N), nil
	case parser.TyKindNullLiteral:
		return "null", nil
	case parser.TyKindCallable:
		return "continuation", nil
	case parser.TyKindVoid:
		return "void", nil
	case parser.TyKindUnknown:
		return "unknown", nil
	case parser.TyKindNullable:
		inner, err := st.renderTy(ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("nullable inner: %w", err)
		}
		return inner + "?", nil
	case parser.TyKindCellOf:
		inner, err := st.renderTy(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("cell inner: %w", err)
		}
		return fmt.Sprintf("Cell<%s>", inner), nil
	case parser.TyKindArrayOf:
		inner, err := st.renderTy(ty.ArrayOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("array inner: %w", err)
		}
		return fmt.Sprintf("array<%s>", inner), nil
	case parser.TyKindLispListOf:
		inner, err := st.renderTy(ty.LispListOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("lisp list inner: %w", err)
		}
		return fmt.Sprintf("lisp_list<%s>", inner), nil
	case parser.TyKindTensor:
		items, err := utils.MapSliceErr(ty.Tensor.ItemsTyIdx, st.renderTy)
		if err != nil {
			return "", fmt.Errorf("tensor items: %w", err)
		}
		return fmt.Sprintf("(%s)", strings.Join(items, ", ")), nil
	case parser.TyKindShapedTuple:
		items, err := utils.MapSliceErr(ty.ShapedTuple.ItemsTyIdx, st.renderTy)
		if err != nil {
			return "", fmt.Errorf("tuple items: %w", err)
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", ")), nil
	case parser.TyKindMapKV:
		key, err := st.renderTy(ty.MapKV.KeyTyIdx)
		if err != nil {
			return "", fmt.Errorf("map key: %w", err)
		}
		value, err := st.renderTy(ty.MapKV.ValueTyIdx)
		if err != nil {
			return "", fmt.Errorf("map value: %w", err)
		}
		return fmt.Sprintf("map<%s, %s>", key, value), nil
	case parser.TyKindEnumRef:
		return ty.EnumRef.EnumName, nil
	case parser.TyKindStructRef:
		args, err := st.renderTypeArgs(ty.StructRef.TypeArgsTyIdx)
		if err != nil {
			return "", fmt.Errorf("struct type args: %w", err)
		}
		return ty.StructRef.StructName + args, nil
	case parser.TyKindAliasRef:
		args, err := st.renderTypeArgs(ty.AliasRef.TypeArgsTyIdx)
		if err != nil {
			return "", fmt.Errorf("alias type args: %w", err)
		}
		return ty.AliasRef.AliasName + args, nil
	case parser.TyKindGenericT:
		return ty.GenericT.NameT, nil
	case parser.TyKindUnion:
		items, err := utils.MapSliceErr(ty.Union.Variants, func(v parser.UnionVariant) (string, error) {
			return st.renderTy(v.VariantTyIdx)
		})
		if err != nil {
			return "", fmt.Errorf("union variants: %w", err)
		}
		return strings.Join(items, " | "), nil
	default:
		return "", fmt.Errorf("unknown type kind %q", ty.SumType)
	}
}

func (st *symTable) renderTypeArgs(args []int) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	rendered, err := utils.MapSliceErr(args, st.renderTy)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("<%s>", strings.Join(rendered, ", ")), nil
}

func (st *symTable) structFieldsOf(tyIdx int, isForStack bool) ([]fieldView, error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return nil, err
	}
	if ty.SumType != parser.TyKindStructRef {
		return nil, fmt.Errorf("expected StructRef at ty_idx=%d", tyIdx)
	}
	decl, ok := st.structs[ty.StructRef.StructName]
	if !ok {
		return nil, fmt.Errorf("struct %q not found", ty.StructRef.StructName)
	}
	fields := make([]fieldView, len(decl.Fields))
	for i, f := range decl.Fields {
		fields[i] = fieldView{Field: f}
	}
	for _, inst := range st.structInstantiations {
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

func (st *symTable) aliasTargetOf(tyIdx int) (targetTyIdx int, uLabelTyIdx *int, err error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return 0, nil, err
	}
	if ty.SumType != parser.TyKindAliasRef {
		return 0, nil, fmt.Errorf("expected AliasRef at ty_idx=%d", tyIdx)
	}
	decl, err := st.getAlias(ty.AliasRef.AliasName)
	if err != nil {
		return 0, nil, err
	}
	targetTyIdx = decl.TargetTyIdx
	for _, inst := range st.aliasInstantiations {
		if inst.TyIdx != tyIdx {
			continue
		}
		old := targetTyIdx
		return inst.MonomorphicTargetTyIdx, &old, nil
	}
	return targetTyIdx, nil, nil
}

func (st *symTable) msgName(bodyTyIdx int) (string, error) {
	ty, err := st.tyByIdx(bodyTyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindStructRef:
		return ty.StructRef.StructName, nil
	case parser.TyKindAliasRef:
		return ty.AliasRef.AliasName, nil
	default:
		return "", fmt.Errorf("cannot get name for %q body", ty.SumType)
	}
}

func prefixConstValue(p *parser.Prefix) (string, error) {
	if p == nil {
		return "", nil
	}
	if p.PrefixLen > 64 {
		return "", fmt.Errorf("prefix length too large: %d", p.PrefixLen)
	}
	if p.PrefixLen%4 == 0 {
		return fmt.Sprintf("0x%0*x", p.PrefixLen/4, p.PrefixNum), nil
	}
	return fmt.Sprintf("0b%0*b", p.PrefixLen, p.PrefixNum), nil
}
