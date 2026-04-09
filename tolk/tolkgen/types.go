package tolkgen

import (
	"fmt"

	"github.com/tonkeeper/tongo/tolk/parser"
)

func emitStoreExpr(expr string, ty parser.Ty) (string, error) {
	switch ty.SumType {
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindAddress, parser.TyKindAddressAny, parser.TyKindAddressOpt,
		parser.TyKindStructRef, parser.TyKindEnumRef, parser.TyKindAliasRef, parser.TyKindMapKV, parser.TyKindBitsN:
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
		if ty.CellOf.Inner.SumType == parser.TyKindSlice {
			return fmt.Sprintf("c.AddRef(&%s)", expr), nil
		}
		return fmt.Sprintf("%s.MarshalTLB(c, encoder)", expr), nil
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
	case parser.TyKindInt, parser.TyKindIntN, parser.TyKindUintN, parser.TyKindVarIntN, parser.TyKindVarUintN, parser.TyKindCoins,
		parser.TyKindAddress, parser.TyKindAddressAny, parser.TyKindAddressOpt,
		parser.TyKindStructRef, parser.TyKindEnumRef, parser.TyKindMapKV:
		goType, err := emitGoType(ty)
		if err != nil {
			return "", false, fmt.Errorf("type %s: %w", ty.String(), err)
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
		if ty.CellOf.Inner.SumType == parser.TyKindSlice {
			return "c.NextRefV()", false, nil
		}
		innerType, err := emitGoType(ty.CellOf.Inner)
		if err != nil {
			return "", false, fmt.Errorf("cellOf inner type %s: %w", ty.CellOf.Inner.String(), err)
		}
		return fmt.Sprintf("tlb.UnmarshalT[tlb.RefT[*%s]](c, decoder)", innerType), true, nil
	case parser.TyKindAliasRef:
		if len(ty.AliasRef.TypeArgs) > 0 {
			return "", false, fmt.Errorf("type arguments not supported for aliases")
		}
		return fmt.Sprintf("tlb.UnmarshalT[%s](c, decoder)", ty.AliasRef.AliasName), true, nil
	case parser.TyKindNullable:
		innerLoadExpr, innerHasLoadMethod, err := st.emitLoadExpr(fieldPath, ty.Nullable.Inner)
		if err != nil {
			return "", false, err
		}
		if innerHasLoadMethod {
			innerType, err := emitGoType(ty.Nullable.Inner)
			if err != nil {
				return "", false, fmt.Errorf("nullable inner type %s: %w", ty.Nullable.Inner.String(), err)
			}
			return fmt.Sprintf(`(func () (result tlb.Maybe[%s], err error) {
		result.UnmarshalTLB(c, decoder)
		return
	})()`, innerType), true, nil
		} else {
			innerType, err := emitGoType(ty.Nullable.Inner)
			if err != nil {
				return "", false, fmt.Errorf("nullable inner type %s: %w", ty.Nullable.Inner.String(), err)
			}
			return fmt.Sprintf(`tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (%s, error) {
	return %s
})`, innerType, innerLoadExpr), false, nil
		}
	}
	return "", false, fmt.Errorf("unknown type %v", ty)
}

func emitGoType(ty parser.Ty) (string, error) {
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
			if len(ty.StructRef.TypeArgs) > 0 {
				tArgs := ""
				for _, t := range ty.StructRef.TypeArgs {
					goType, err := emitGoType(t)
					if err != nil {
						return "", err
					}
					tArgs += goType + ", "
				}
				name += "[" + tArgs + "]"
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
		if ty.CellOf.Inner.SumType == parser.TyKindSlice {
			return "boc.Cell", nil
		}
		innerType, err := emitGoType(ty.CellOf.Inner)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.RefT[*%s]", innerType), nil
	case parser.TyKindAliasRef:
		name := ty.AliasRef.AliasName
		if len(ty.AliasRef.TypeArgs) > 0 {
			tArgs := ""
			for _, t := range ty.AliasRef.TypeArgs {
				goType, err := emitGoType(t)
				if err != nil {
					return "", err
				}
				tArgs = tArgs + goType + ", "
			}
			name = name + "[" + tArgs + "]"
		}
		return name, nil
	case parser.TyKindNullable:
		innerType, err := emitGoType(ty.Nullable.Inner)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.Maybe[%s]", innerType), nil
	case parser.TyKindMapKV:
		kType, err := emitGoType(ty.MapKV.K)
		if err != nil {
			return "", err
		}
		vType, err := emitGoType(ty.MapKV.V)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tlb.HashmapE[%s, %s]", kType, vType), nil
	case parser.TyKindArrayOf:
		innerType, err := emitGoType(ty.ArrayOf.Inner)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[]%s", innerType), nil
	case parser.TyKindVoid:
		return "struct{}", nil
	}
	return "", fmt.Errorf("go type not supported: %v", ty)
}
