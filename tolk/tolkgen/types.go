package tolkgen

import (
	"fmt"
	"strings"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

func (st *symTable) emitStoreExpr(expr string, tyIdx int) (string, error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindAddress, parser.TyKindAddressAny, parser.TyKindAddressOpt,
		parser.TyKindStructRef, parser.TyKindEnumRef, parser.TyKindAliasRef, parser.TyKindMapKV, parser.TyKindBitsN, parser.TyKindTensor:
		return fmt.Sprintf("%s.MarshalTLB(c, encoder)", expr), nil
	case parser.TyKindBool:
		return fmt.Sprintf("c.WriteBit(%s)", expr), nil
	case parser.TyKindCell, parser.TyKindRemaining:
		return fmt.Sprintf("c.AddRef(&%s)", expr), nil
	case parser.TyKindSlice:
		return fmt.Sprintf("c.AddRef(%s)", expr), nil
	case parser.TyKindString:
		return fmt.Sprintf("c.WriteStringRefTail(%s)", expr), nil
	case parser.TyKindNullable:
		return fmt.Sprintf("%s.MarshalTLB(c, encoder)", expr), nil
	case parser.TyKindCellOf:
		inner, err := st.tyByIdx(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("cellOf inner: %w", err)
		}
		if inner.SumType == parser.TyKindSlice {
			return fmt.Sprintf("c.AddRef(&%s)", expr), nil
		}
		return fmt.Sprintf("%s.MarshalTLB(c, encoder)", expr), nil
	}
	return "", fmt.Errorf("unknown type %v", ty)
}

func (st *symTable) emitLoadExpr(fieldPath string, tyIdx int) (expr string, hasLoadMethod bool, err error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return "", false, err
	}

	switch ty.SumType {
	case parser.TyKindBuilder, parser.TyKindSlice, parser.TyKindUnknown, parser.TyKindCallable:
		renderedTy, err := st.renderTy(tyIdx)
		if err != nil {
			return "", false, err
		}
		var hint string
		if ty.SumType == parser.TyKindBuilder || ty.SumType == parser.TyKindSlice {
			hint = " (it can be used for writing only)"
		}
		return "", false, fmt.Errorf("%s is %s%s", fieldPath, renderedTy, hint)
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindAddress, parser.TyKindAddressAny, parser.TyKindAddressOpt,
		parser.TyKindStructRef, parser.TyKindEnumRef, parser.TyKindMapKV:
		goType, err := st.emitGoType(tyIdx)
		if err != nil {
			return "", false, fmt.Errorf("type: %w", err)
		}
		return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", goType), true, nil
	case parser.TyKindBitsN:
		return fmt.Sprintf("tlb.UnmarshalT[tlb.Bits%d](c, decoder)", ty.BitsN.N), true, nil
	case parser.TyKindBool:
		return "c.ReadBit()", false, nil
	case parser.TyKindCell:
		return "c.NextRefV()", false, nil
	case parser.TyKindString:
		return "c.ReadStringRefTail()", false, nil
	case parser.TyKindRemaining:
		return `(func () (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})()`, false, nil
	case parser.TyKindCellOf:
		inner, err := st.tyByIdx(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("cellOf inner: %w", err)
		}
		if inner.SumType == parser.TyKindSlice {
			return "c.NextRefV()", false, nil
		}
		innerType, err := st.emitGoType(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", false, fmt.Errorf("cellOf inner type: %w", err)
		}
		return fmt.Sprintf("tlb.UnmarshalT[tlb.RefT[*%s]](c, decoder)", innerType), true, nil
	case parser.TyKindAliasRef:
		if len(ty.AliasRef.TypeArgsTyIdx) > 0 {
			return "", false, fmt.Errorf("type arguments not supported for aliases")
		}
		return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", ty.AliasRef.AliasName), true, nil
	case parser.TyKindNullable:
		innerLoadExpr, innerHasLoadMethod, err := st.emitLoadExpr(fieldPath, ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", false, err
		}
		if innerHasLoadMethod {
			innerType, err := st.emitGoType(ty.Nullable.InnerTyIdx)
			if err != nil {
				return "", false, fmt.Errorf("nullable inner type: %w", err)
			}
			return fmt.Sprintf(`(func () (result tlb.Maybe[%s], err error) {
		result.UnmarshalTLB(c, decoder)
		return
	})()`, innerType), true, nil
		} else {
			innerType, err := st.emitGoType(ty.Nullable.InnerTyIdx)
			if err != nil {
				return "", false, fmt.Errorf("nullable inner type: %w", err)
			}
			return fmt.Sprintf(`tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (%s, error) {
	return %s
})`, innerType, innerLoadExpr), false, nil
		}
	case parser.TyKindTensor:
		if len(ty.Tensor.ItemsTyIdx) > tlb.MaxTensorSize {
			return "", false, fmt.Errorf("tensor size %d is too big, update tlb/tensor.go", len(ty.Tensor.ItemsTyIdx))
		}
		loaders := make([]string, len(ty.Tensor.ItemsTyIdx))
		allHaveMethods := true
		for i, itemTyIdx := range ty.Tensor.ItemsTyIdx {
			expr, hasMethod, err := st.emitLoadExpr(fmt.Sprintf("%s[%d]", fieldPath, i), itemTyIdx)
			allHaveMethods = allHaveMethods && hasMethod
			if err != nil {
				return "", false, fmt.Errorf("tensor item %d: %w", i, err)
			}
			loaders[i] = expr
		}
		if allHaveMethods {
			typ, err := st.emitGoType(tyIdx)
			if err != nil {
				return "", false, fmt.Errorf("tensor type: %w", err)
			}
			return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", typ), true, nil
		} else {
			return "", false, fmt.Errorf("tensor items have no methods")
		}
	}
	return "", false, fmt.Errorf("unknown type %v", ty)
}

func (st *symTable) emitGoType(tyIdx int) (string, error) {
	ty, err := st.tyByIdx(tyIdx)
	if err != nil {
		return "", err
	}
	switch ty.SumType {
	case parser.TyKindInt:
		return "tlb.Int257", nil
	case parser.TyKindIntN:
		return fmt.Sprintf("tlb.Int%d", ty.IntN.N), nil
	case parser.TyKindUintN:
		return fmt.Sprintf("tlb.Uint%d", ty.UintN.N), nil
	case parser.TyKindVarIntN:
		return fmt.Sprintf("tlb.VarInteger%d", ty.VarIntN.N), nil
	case parser.TyKindVarUintN:
		return fmt.Sprintf("tlb.VarUInteger%d", ty.VarUintN.N), nil
	case parser.TyKindCoins:
		return "tlb.Coins", nil
	case parser.TyKindBool:
		return "bool", nil
	case parser.TyKindCell, parser.TyKindSlice, parser.TyKindRemaining:
		return "boc.Cell", nil
	case parser.TyKindString:
		return "string", nil
	case parser.TyKindAddress:
		return "tlb.InternalAddress", nil
	case parser.TyKindAddressAny, parser.TyKindAddressOpt:
		return "tlb.MsgAddress", nil
	case parser.TyKindStructRef:
		switch ty.StructRef.StructName {
		case "CurrencyCollection":
			return "tlb.CurrencyCollection", nil
		default:
			name := ty.StructRef.StructName
			if len(ty.StructRef.TypeArgsTyIdx) > 0 {
				tArgs := make([]string, 0, len(ty.StructRef.TypeArgsTyIdx))
				for _, t := range ty.StructRef.TypeArgsTyIdx {
					goType, err := st.emitGoType(t)
					if err != nil {
						return "", err
					}
					tArgs = append(tArgs, goType)
				}
				name += "[" + strings.Join(tArgs, ", ") + "]"
			}
			return name, nil
		}
	case parser.TyKindEnumRef:
		return safeGoIdent(ty.EnumRef.EnumName), nil
	case parser.TyKindBitsN:
		switch ty.BitsN.N {
		case 80, 96, 128, 160, 256, 264, 320, 352, 512:
			return "tlb.Bits" + fmt.Sprintf("%d", ty.BitsN.N), nil
		}
		return "", fmt.Errorf("tlb.Bits%d is not supported: update cmd/codegen/integers", ty.BitsN.N)
	case parser.TyKindCellOf:
		inner, err := st.tyByIdx(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", fmt.Errorf("cellOf inner: %w", err)
		}
		if inner.SumType == parser.TyKindSlice {
			return "boc.Cell", nil
		}
		innerType, err := st.emitGoType(ty.CellOf.InnerTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.RefT[*%s]", innerType), nil
	case parser.TyKindAliasRef:
		name := ty.AliasRef.AliasName
		if len(ty.AliasRef.TypeArgsTyIdx) > 0 {
			tArgs := make([]string, 0, len(ty.AliasRef.TypeArgsTyIdx))
			for _, t := range ty.AliasRef.TypeArgsTyIdx {
				goType, err := st.emitGoType(t)
				if err != nil {
					return "", err
				}
				tArgs = append(tArgs, goType)
			}
			name = name + "[" + strings.Join(tArgs, ", ") + "]"
		}
		return name, nil
	case parser.TyKindNullable:
		innerType, err := st.emitGoType(ty.Nullable.InnerTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.Maybe[%s]", innerType), nil
	case parser.TyKindMapKV:
		kType, err := st.emitGoType(ty.MapKV.KeyTyIdx)
		if err != nil {
			return "", err
		}
		vType, err := st.emitGoType(ty.MapKV.ValueTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.HashmapE[%s, %s]", kType, vType), nil
	case parser.TyKindArrayOf:
		innerType, err := st.emitGoType(ty.ArrayOf.InnerTyIdx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[]%s", innerType), nil
	case parser.TyKindVoid:
		return "struct{}", nil
	case parser.TyKindTensor:
		items := ty.Tensor.ItemsTyIdx
		if len(items) > tlb.MaxTensorSize {
			return "", fmt.Errorf("tensor size %d is too big, update tlb/tensor.go", len(items))
		}
		params, err := utils.MapSliceErr(items, st.emitGoType)
		if err != nil {
			return "", fmt.Errorf("tensor items: %w", err)
		}
		if len(items) > 0 {
			return fmt.Sprintf("tlb.Tensor%d[%s]", len(items), strings.Join(params, ", ")), nil
		} else {
			return "tlb.Tensor0", nil
		}
	case parser.TyKindShapedTuple:
		items := ty.ShapedTuple.ItemsTyIdx
		if len(items) > tlb.MaxTensorSize {
			return "", fmt.Errorf("tensor size %d is too big, update tlb/tensor.go", len(items))
		}
		params, err := utils.MapSliceErr(items, st.emitGoType)
		if err != nil {
			return "", fmt.Errorf("tensor items: %w", err)
		}
		if len(items) > 0 {
			return fmt.Sprintf("tlb.ShapedTuple%d[%s]", len(items), strings.Join(params, ", ")), nil
		} else {
			return "tlb.ShapedTuple0", nil
		}
	}
	return "", fmt.Errorf("go type not supported: %v", ty)
}
