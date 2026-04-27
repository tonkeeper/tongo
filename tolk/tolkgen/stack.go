package tolkgen

import (
	"fmt"

	"github.com/tonkeeper/tongo/tolk/parser"
)

func (tgen TolkGolangGenerator) calcWidthOnStack(ty parser.Ty) (int, error) {
	switch ty.SumType {
	case parser.TyKindVoid:
		return 0, nil
	case parser.TyKindTensor:
		totalWidth := 0
		for i, item := range ty.Tensor.Items {
			itemW, err := tgen.calcWidthOnStack(item)
			if err != nil {
				return 0, fmt.Errorf("tensor item %d: %w", i, err)
			}
			totalWidth += itemW
		}
		return totalWidth, nil
	case parser.TyKindStructRef:
		sTy, ok := tgen.symbols.structs[ty.StructRef.StructName]
		if !ok {
			return 0, fmt.Errorf("struct %q not found in declarations", ty.StructRef.StructName)
		}
		totalWidth := 0
		for _, f := range sTy.Fields {
			fTy := f.Ty
			if len(ty.StructRef.TypeArgs) > 0 {
				fTy = fTy.InstantiateGenerics(sTy.TypeParams, ty.StructRef.TypeArgs)
			}
			n, err := tgen.calcWidthOnStack(fTy)
			if err != nil {
				return 0, fmt.Errorf("field %q: %w", f.Name, err)
			}
			totalWidth += n
		}
		return totalWidth, nil
	case parser.TyKindAliasRef:
		aliasRef, err := tgen.symbols.getAlias(ty.AliasRef.AliasName)
		if err != nil {
			return 0, fmt.Errorf("alias %q: %w", ty.AliasRef.AliasName, err)
		}
		targetTy := aliasRef.TargetTy
		if len(aliasRef.TypeParams) > 0 {
			targetTy = targetTy.InstantiateGenerics(aliasRef.TypeParams, ty.AliasRef.TypeArgs)
		}
		return tgen.calcWidthOnStack(targetTy)
	case parser.TyKindNullable:
		return max(ty.Nullable.StackWidth, 1), nil
	case parser.TyKindUnion:
		return ty.Union.StackWidth, nil
	case parser.TyKindGenericT:
		return 0, fmt.Errorf("unexpected genericT=%sTy in calcWidthOnStack", ty.Generic.NameT)
	default:
		return 1, nil
	}
}

func stackReturnGoType(ty parser.Ty) (string, error) {
	switch ty.SumType {
	case parser.TyKindTensor:
		return "", fmt.Errorf("tensor return type is not supported: use a struct instead")
	default:
		goType, err := emitGoType(ty)
		if err != nil {
			return "", fmt.Errorf("return type %s: %w", ty.String(), err)
		}
		return goType, nil
	}
}
