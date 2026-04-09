package parser

import "slices"

// InstantiateGenerics replaces all generic Ts (typeParams) with instantiation (typeArgs) recursively.
// Example: `(int, T, Wrapper<T?>)` and T=coins => `(int, coins, Wrapper<coins?>)`
// A new struct is returned
func (ty Ty) InstantiateGenerics(typeParams []string, typeArgs []Ty) Ty {
	switch ty.SumType {
	case TyKindNullable:
		ty.Nullable = &Nullable{
			Inner: ty.Nullable.Inner.InstantiateGenerics(typeParams, typeArgs),
		}
	case TyKindCellOf:
		ty.CellOf = &CellOf{
			Inner: ty.CellOf.Inner.InstantiateGenerics(typeParams, typeArgs),
		}
	case TyKindArrayOf:
		ty.ArrayOf = &ArrayOf{
			Inner: ty.ArrayOf.Inner.InstantiateGenerics(typeParams, typeArgs),
		}
	case TyKindLispListOf:
		ty.LispListOf = &LispListOf{
			Inner: ty.LispListOf.Inner.InstantiateGenerics(typeParams, typeArgs),
		}
	case TyKindTensor:
		updatedItems := make([]Ty, len(ty.Tensor.Items))
		for i := range ty.Tensor.Items {
			updatedItems[i] = ty.Tensor.Items[i].InstantiateGenerics(typeParams, typeArgs)
		}
		ty.Tensor = &Tensor{
			Items: updatedItems,
		}
	case TyKindShapedTuple:
		updatedItems := make([]Ty, len(ty.ShapedTuple.Items))
		for i := range ty.ShapedTuple.Items {
			updatedItems[i] = ty.ShapedTuple.Items[i].InstantiateGenerics(typeParams, typeArgs)
		}
		ty.ShapedTuple = &ShapedTuple{
			Items: updatedItems,
		}
	case TyKindMapKV:
		ty.MapKV = &MapKV{
			K: ty.MapKV.K.InstantiateGenerics(typeParams, typeArgs),
			V: ty.MapKV.V.InstantiateGenerics(typeParams, typeArgs),
		}
	case TyKindStructRef:
		if len(ty.StructRef.TypeArgs) > 0 {
			updatedTypeArgs := make([]Ty, len(ty.StructRef.TypeArgs))
			for i := range ty.StructRef.TypeArgs {
				updatedTypeArgs[i] = ty.StructRef.TypeArgs[i].InstantiateGenerics(typeParams, typeArgs)
			}
			ty.StructRef = &StructRef{
				StructName: ty.StructRef.StructName,
				TypeArgs:   updatedTypeArgs,
			}
		}
	case TyKindAliasRef:
		if len(ty.AliasRef.TypeArgs) > 0 {
			updatedTypeArgs := make([]Ty, len(ty.AliasRef.TypeArgs))
			for i := range ty.AliasRef.TypeArgs {
				updatedTypeArgs[i] = ty.AliasRef.TypeArgs[i].InstantiateGenerics(typeParams, typeArgs)
			}
			ty.AliasRef = &AliasRef{
				AliasName: ty.AliasRef.AliasName,
				TypeArgs:  updatedTypeArgs,
			}
		}
	case TyKindUnion:
		newVariants := make([]UnionVariant, len(ty.Union.Variants))
		for i := range ty.Union.Variants {
			newVariants[i] = ty.Union.Variants[i]
			newVariants[i].VariantTy = newVariants[i].VariantTy.InstantiateGenerics(typeParams, typeArgs)
		}
		ty.Union = &Union{
			Variants: make([]UnionVariant, len(ty.Union.Variants)),
		}
	case TyKindGenericT:
		if i := slices.Index(typeParams, ty.Generic.NameT); 0 <= i && i < len(typeArgs) {
			return typeArgs[i]
		}
		panic("inconsistent generics: could not find type argument for " + ty.Generic.NameT)
	default:
	}
	return ty
}
