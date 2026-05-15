package tolkgen

import (
	"fmt"

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
		goType, err := tgen.symbols.emitGoType(tyIdx)
		if err != nil {
			return "", fmt.Errorf("return type: %w", err)
		}
		return goType, nil
	}
}
