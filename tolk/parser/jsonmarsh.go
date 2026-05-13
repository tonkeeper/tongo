package parser

import (
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/utils"
)

// jsonKind used for discriminator decoding
type jsonKind struct {
	Kind string `json:"kind"`
}

func (d *ABIDeclaration) UnmarshalJSON(b []byte) error {
	var r jsonKind
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	switch r.Kind {
	case "struct":
		d.SumType = DeclarationKindStruct
		if err := json.Unmarshal(b, &d.StructDeclaration); err != nil {
			return err
		}
	case "alias":
		d.SumType = DeclarationKindAlias
		if err := json.Unmarshal(b, &d.AliasDeclaration); err != nil {
			return err
		}
	case "enum":
		d.SumType = DeclarationKindEnum
		if err := json.Unmarshal(b, &d.EnumDeclaration); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown declaration type %q", d.SumType)
	}

	return nil
}

func (d ABIDeclaration) MarshalJSON() ([]byte, error) {
	var kind jsonKind
	kind.Kind = string(d.SumType)

	var payload []byte
	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}

	switch d.SumType {
	case DeclarationKindStruct:
		payload, err = json.Marshal(d.StructDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	case DeclarationKindAlias:
		payload, err = json.Marshal(d.AliasDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	case DeclarationKindEnum:
		payload, err = json.Marshal(d.EnumDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	default:
		return nil, fmt.Errorf("unknown declaration type %q", d.SumType)
	}
}

func (t *Ty) UnmarshalJSON(b []byte) error {
	var kind jsonKind
	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}
	switch kind.Kind {
	case TyKindIntN:
		t.SumType = TyKindIntN
		return json.Unmarshal(b, &t.IntN)
	case TyKindUintN:
		t.SumType = TyKindUintN
		return json.Unmarshal(b, &t.UintN)
	case TyKindVarIntN:
		t.SumType = TyKindVarIntN
		return json.Unmarshal(b, &t.VarIntN)
	case TyKindVarUintN:
		t.SumType = TyKindVarUintN
		return json.Unmarshal(b, &t.VarUintN)
	case TyKindBitsN:
		t.SumType = TyKindBitsN
		return json.Unmarshal(b, &t.BitsN)
	case TyKindNullable:
		t.SumType = TyKindNullable
		return json.Unmarshal(b, &t.Nullable)
	case TyKindCellOf:
		t.SumType = TyKindCellOf
		return json.Unmarshal(b, &t.CellOf)
	case TyKindArrayOf:
		t.SumType = TyKindArrayOf
		return json.Unmarshal(b, &t.ArrayOf)
	case TyKindLispListOf:
		t.SumType = TyKindLispListOf
		return json.Unmarshal(b, &t.LispListOf)
	case TyKindTensor:
		t.SumType = TyKindTensor
		return json.Unmarshal(b, &t.Tensor)
	case TyKindShapedTuple:
		t.SumType = TyKindShapedTuple
		return json.Unmarshal(b, &t.ShapedTuple)
	case TyKindMapKV:
		t.SumType = TyKindMapKV
		return json.Unmarshal(b, &t.MapKV)
	case TyKindEnumRef:
		t.SumType = TyKindEnumRef
		return json.Unmarshal(b, &t.EnumRef)
	case TyKindStructRef:
		t.SumType = TyKindStructRef
		return json.Unmarshal(b, &t.StructRef)
	case TyKindAliasRef:
		t.SumType = TyKindAliasRef
		return json.Unmarshal(b, &t.AliasRef)
	case TyKindGenericT:
		t.SumType = TyKindGenericT
		return json.Unmarshal(b, &t.GenericT)
	case TyKindUnion:
		t.SumType = TyKindUnion
		return json.Unmarshal(b, &t.Union)
	case TyKindInt, TyKindCoins, TyKindBool, TyKindCell, TyKindBuilder, TyKindSlice, TyKindString, TyKindRemaining, TyKindAddress:
		t.SumType = kind.Kind
	case TyKindAddressOpt, TyKindAddressExt, TyKindAddressAny, TyKindNullLiteral, TyKindCallable, TyKindVoid, TyKindUnknown:
		t.SumType = kind.Kind
		return nil
	default:
		return fmt.Errorf("unknown type kind %q", kind.Kind)
	}

	return nil
}

func (t *Ty) MarshalJSON() ([]byte, error) {
	var payload []byte
	var err error
	switch t.SumType {
	case TyKindIntN:
		payload, err = json.Marshal(t.IntN)
	case TyKindUintN:
		payload, err = json.Marshal(t.UintN)
	case TyKindVarIntN:
		payload, err = json.Marshal(t.VarIntN)
	case TyKindVarUintN:
		payload, err = json.Marshal(t.VarUintN)
	case TyKindBitsN:
		payload, err = json.Marshal(t.BitsN)
	case TyKindNullable:
		payload, err = json.Marshal(t.Nullable)
	case TyKindCellOf:
		payload, err = json.Marshal(t.CellOf)
	case TyKindArrayOf:
		payload, err = json.Marshal(t.ArrayOf)
	case TyKindLispListOf:
		payload, err = json.Marshal(t.LispListOf)
	case TyKindTensor:
		payload, err = json.Marshal(t.Tensor)
	case TyKindShapedTuple:
		payload, err = json.Marshal(t.ShapedTuple)
	case TyKindMapKV:
		payload, err = json.Marshal(t.MapKV)
	case TyKindEnumRef:
		payload, err = json.Marshal(t.EnumRef)
	case TyKindStructRef:
		payload, err = json.Marshal(t.StructRef)
	case TyKindAliasRef:
		payload, err = json.Marshal(t.AliasRef)
	case TyKindGenericT:
		payload, err = json.Marshal(t.GenericT)
	case TyKindUnion:
		payload, err = json.Marshal(t.Union)
	default:
	}
	if err != nil {
		return nil, err
	}
	return utils.ConcatPrefixAndSuffixIfExists([]byte(fmt.Sprintf(`{"kind": "%s"}`, t.SumType)), payload), nil
}
