package tolkgen

import (
	"fmt"
	"strings"

	"github.com/tonkeeper/tongo/tolk/parser"
)

func (tgen TolkGolangGenerator) calcWidthOnStack(tyIdx int) (int, error) {
	ty, err := tgen.symbols.TyByIdx(tyIdx)
	if err != nil {
		return 0, err
	}
	switch ty.SumType {
	case parser.TyKindVoid:
		return 0, nil
	case parser.TyKindTensor:
		totalWidth := 0
		for i, item := range ty.Tensor.ItemsTyIdx {
			itemW, err := tgen.calcWidthOnStack(item)
			if err != nil {
				return 0, fmt.Errorf("tensor item %d: %w", i, err)
			}
			totalWidth += itemW
		}
		return totalWidth, nil
	case parser.TyKindStructRef:
		fields, err := tgen.symbols.StructFieldsOf(tyIdx, true)
		if err != nil {
			return 0, err
		}
		totalWidth := 0
		for _, f := range fields {
			n, err := tgen.calcWidthOnStack(f.TyIdx)
			if err != nil {
				return 0, fmt.Errorf("field %q: %w", f.Name, err)
			}
			totalWidth += n
		}
		return totalWidth, nil
	case parser.TyKindAliasRef:
		targetTyIdx, _, err := tgen.symbols.AliasTargetOf(tyIdx)
		if err != nil {
			return 0, fmt.Errorf("alias %q: %w", ty.AliasRef.AliasName, err)
		}
		return tgen.calcWidthOnStack(targetTyIdx)
	case parser.TyKindNullable:
		return max(ty.Nullable.StackWidth, 1), nil
	case parser.TyKindUnion:
		return ty.Union.StackWidth, nil
	case parser.TyKindGenericT:
		return 0, fmt.Errorf("unexpected genericT=%sTy in calcWidthOnStack", ty.GenericT.NameT)
	default:
		return 1, nil
	}
}

func (tgen TolkGolangGenerator) stackReturnGoType(tyIdx int) (string, error) {
	ty, err := tgen.symbols.TyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindTensor:
		return "", fmt.Errorf("tensor return type is not supported: use a struct instead")
	default:
		goType, err := tgen.emitStackGoType(tyIdx)
		if err != nil {
			return "", fmt.Errorf("return type: %w", err)
		}
		return goType, nil
	}
}

// emitStackGoType returns the Go type used when decoding tyIdx from a VM stack.
// It matches emitGoType except for dictionaries: a map returned on the stack is a
// directly-loaded dictionary root (tlb.Hashmap), not the in-cell tlb.HashmapE which
// carries a Maybe prefix. Container kinds recurse so nested maps get the same
// treatment; everything else delegates to emitGoType.
func (tgen TolkGolangGenerator) emitStackGoType(tyIdx int) (string, error) {
	ty, err := tgen.symbols.TyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindMapKV:
		kType, err := tgen.emitStackGoType(ty.MapKV.KeyTyIdx)
		if err != nil {
			return "", err
		}
		vType, err := tgen.emitStackGoType(ty.MapKV.ValueTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.Hashmap[%s, %s]", kType, vType), nil
	case parser.TyKindArrayOf:
		inner, err := tgen.emitStackGoType(ty.ArrayOf.InnerTyIdx)
		if err != nil {
			return "", err
		}
		return "[]" + inner, nil
	case parser.TyKindNullable:
		inner, err := tgen.emitStackGoType(ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.Maybe[%s]", inner), nil
	case parser.TyKindTensor, parser.TyKindShapedTuple:
		var items []int
		if ty.SumType == parser.TyKindTensor {
			items = ty.Tensor.ItemsTyIdx
		} else {
			items = ty.ShapedTuple.ItemsTyIdx
		}
		params := make([]string, len(items))
		for i, item := range items {
			params[i], err = tgen.emitStackGoType(item)
			if err != nil {
				return "", err
			}
		}
		prefix := "tlb.Tensor"
		if ty.SumType == parser.TyKindShapedTuple {
			prefix = "tlb.ShapedTuple"
		}
		if len(items) == 0 {
			return prefix + "0", nil
		}
		return fmt.Sprintf("%s%d[%s]", prefix, len(items), strings.Join(params, ", ")), nil
	default:
		return tgen.symbols.emitGoType(tyIdx)
	}
}
