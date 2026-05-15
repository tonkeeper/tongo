package runtime

// todo: move this to some package or rename somehow.

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

type SumType string

const (
	SumTypeBool            SumType = "Bool"
	SumTypeSmallInt        SumType = "SmallInt"
	SumTypeSmallUint       SumType = "SmallUint"
	SumTypeBigInt          SumType = "BigInt"
	SumTypeBigUint         SumType = "BigUint"
	SumTypeVarInt          SumType = "VarInt"
	SumTypeVarUint         SumType = "VarUint"
	SumTypeCoins           SumType = "Coins"
	SumTypeBits            SumType = "Bits"
	SumTypeCell            SumType = "Cell"
	SumTypeRemaining       SumType = "Remaining"
	SumTypeInternalAddress SumType = "InternalAddress"
	SumTypeOptionalAddress SumType = "OptionalAddress"
	SumTypeExternalAddress SumType = "ExternalAddress"
	SumTypeAnyAddress      SumType = "AnyAddress"
	SumTypeNoneAddress     SumType = "NoneAddress"
	SumTypeVarAddress      SumType = "VarAddress"
	SumTypeOptionalValue   SumType = "OptionalValue"
	SumTypeRefValue        SumType = "RefValue"
	SumTypeTupleWith       SumType = "TupleWith"
	SumTypeTensor          SumType = "Tensor"
	SumTypeMap             SumType = "Map"
	SumTypeStruct          SumType = "Struct"
	SumTypeAlias           SumType = "Alias"
	SumTypeEnum            SumType = "Enum"
	SumTypeGenericT        SumType = "GenericT"
	SumTypeUnion           SumType = "Union"
	SumTypeNull            SumType = "Null"
	SumTypeVoid            SumType = "Void"
)

func (s SumType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + utils.ToCamelCasePrivate(string(s)) + `"`), nil
}

func (s *SumType) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid SumType value: %s", string(bytes))
	}
	*s = SumType(bytes[1 : len(bytes)-1])
	return nil
}

type Value struct {
	SumType         SumType
	Bool            *BoolValue
	SmallInt        *Int64
	SmallUint       *UInt64
	BigInt          *BigInt
	BigUint         *BigUInt
	VarInt          *VarInt
	VarUint         *VarUInt
	Coins           *CoinsValue
	Bits            *Bits
	Cell            *Any
	Remaining       *RemainingValue
	InternalAddress *InternalAddress
	OptionalAddress *OptionalAddress
	ExternalAddress *ExternalAddress
	AnyAddress      *AnyAddress
	OptionalValue   *OptValue
	RefValue        *RefValue
	TupleWith       *TupleValues
	Tensor          *TensorValues
	Map             *MapValue
	Struct          *Struct
	Alias           *AliasValue
	Enum            *EnumValue
	Generic         *GenericValue
	Union           *UnionValue
	Null            *NullValue
	Void            *VoidValue
}

func (v *Value) GetBool() (bool, bool) {
	if v.Bool == nil {
		return false, false
	}
	return bool(*v.Bool), true
}

func (v *Value) MustGetBool() bool {
	if v.Bool == nil {
		panic("value is not a bool")
	}
	return bool(*v.Bool)
}

func (v *Value) GetSmallInt() (int64, bool) {
	if v.SmallInt == nil {
		return 0, false
	}
	return int64(*v.SmallInt), true
}

func (v *Value) MustGetSmallInt() int64 {
	if v.SmallInt == nil {
		panic("value is not a small int")
	}
	return int64(*v.SmallInt)
}

func (v *Value) GetBigInt() (big.Int, bool) {
	if v.BigInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.BigInt), true
}

func (v *Value) MustGetBigInt() big.Int {
	if v.BigInt == nil {
		panic("value is not a big int")
	}
	return big.Int(*v.BigInt)
}

func (v *Value) GetSmallUInt() (uint64, bool) {
	if v.SmallUint == nil {
		return 0, false
	}
	return uint64(*v.SmallUint), true
}

func (v *Value) MustGetSmallUInt() uint64 {
	if v.SmallUint == nil {
		panic("value is not a small uint")
	}
	return uint64(*v.SmallUint)
}

func (v *Value) GetBigUInt() (big.Int, bool) {
	if v.BigUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.BigUint), true
}

func (v *Value) MustGetBigUInt() big.Int {
	if v.BigUint == nil {
		panic("value is not a big uint")
	}
	return big.Int(*v.BigUint)
}

func (v *Value) GetVarInt() (big.Int, bool) {
	if v.VarInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.VarInt), true
}

func (v *Value) MustGetVarInt() big.Int {
	if v.VarInt == nil {
		panic("value is not a var int")
	}
	return big.Int(*v.VarInt)
}

func (v *Value) GetVarUInt() (big.Int, bool) {
	if v.VarUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.VarUint), true
}

func (v *Value) MustGetVarUInt() big.Int {
	if v.VarUint == nil {
		panic("value is not a var uint")
	}
	return big.Int(*v.VarUint)
}

func (v *Value) GetCoins() (big.Int, bool) {
	if v.Coins == nil {
		return big.Int{}, false
	}
	return big.Int(*v.Coins), true
}

func (v *Value) MustGetCoins() big.Int {
	if v.Coins == nil {
		panic("value is not a coins")
	}
	return big.Int(*v.Coins)
}

func (v *Value) GetBits() (boc.BitString, bool) {
	if v.Bits == nil {
		return boc.BitString{}, false
	}
	return boc.BitString(*v.Bits), true
}

func (v *Value) MustGetBits() boc.BitString {
	if v.Bits == nil {
		panic("value is not a bits")
	}
	return boc.BitString(*v.Bits)
}

func (v *Value) GetAddress() (InternalAddress, bool) {
	if v.InternalAddress == nil {
		return InternalAddress{}, false
	}
	return *v.InternalAddress, true
}

func (v *Value) MustGetAddress() InternalAddress {
	if v.InternalAddress == nil {
		panic("value is not an internal address")
	}
	return *v.InternalAddress
}

func (v *Value) GetOptionalAddress() (OptionalAddress, bool) {
	if v.OptionalAddress == nil {
		return OptionalAddress{}, false
	}
	return *v.OptionalAddress, true
}

func (v *Value) MustGetOptionalAddress() OptionalAddress {
	if v.OptionalAddress == nil {
		panic("value is not an optional address")
	}
	return *v.OptionalAddress
}

func (v *Value) GetExternalAddress() (ExternalAddress, bool) {
	if v.ExternalAddress == nil {
		return ExternalAddress{}, false
	}
	return *v.ExternalAddress, true
}

func (v *Value) MustGetExternalAddress() ExternalAddress {
	if v.ExternalAddress == nil {
		panic("value is not an external address")
	}
	return *v.ExternalAddress
}

func (v *Value) GetAnyAddress() (AnyAddress, bool) {
	if v.AnyAddress == nil {
		return AnyAddress{}, false
	}
	return *v.AnyAddress, true
}

func (v *Value) MustGetAnyAddress() AnyAddress {
	if v.AnyAddress == nil {
		panic("value is not an any address")
	}
	return *v.AnyAddress
}

func (v *Value) GetOptionalValue() (OptValue, bool) {
	if v.OptionalValue == nil {
		return OptValue{}, false
	}
	return *v.OptionalValue, true
}

func (v *Value) MustGetOptionalValue() OptValue {
	if v.OptionalValue == nil {
		panic("value is not an optional value")
	}
	return *v.OptionalValue
}

func (v *Value) GetRefValue() (Value, bool) {
	if v.RefValue == nil {
		return Value{}, false
	}
	return Value(*v.RefValue), true
}

func (v *Value) MustGetRefValue() Value {
	if v.RefValue == nil {
		panic("value is not a reference")
	}
	return Value(*v.RefValue)
}

func (v *Value) GetTensor() ([]Value, bool) {
	if v.Tensor == nil {
		return TensorValues{}, false
	}
	return *v.Tensor, true
}

func (v *Value) MustGetTensor() []Value {
	if v.Tensor == nil {
		panic("value is not a tensor")
	}
	return *v.Tensor
}

func (v *Value) GetMap() (MapValue, bool) {
	if v.Map == nil {
		return MapValue{}, false
	}
	return *v.Map, true
}

func (v *Value) MustGetMap() MapValue {
	if v.Map == nil {
		panic("value is not a map")
	}
	return *v.Map
}

func (v *Value) GetStruct() (Struct, bool) {
	if v.Struct == nil {
		return Struct{}, false
	}
	return *v.Struct, true
}

func (v *Value) MustGetStruct() Struct {
	if v.Struct == nil {
		panic("value is not a struct")
	}
	return *v.Struct
}

func (v *Value) GetAlias() (Value, bool) {
	if v.Alias == nil {
		return Value{}, false
	}
	return Value(*v.Alias), true
}

func (v *Value) MustGetAlias() Value {
	if v.Alias == nil {
		panic("value is not an alias")
	}
	return Value(*v.Alias)
}

func (v *Value) GetGeneric() (Value, bool) {
	if v.Generic == nil {
		return Value{}, false
	}
	return Value(*v.Generic), true
}

func (v *Value) MustGetGeneric() Value {
	if v.Generic == nil {
		panic("value is not a generic")
	}
	return Value(*v.Generic)
}

func (v *Value) GetEnum() (EnumValue, bool) {
	if v.Enum == nil {
		return EnumValue{}, false
	}
	return *v.Enum, true
}

func (v *Value) MustGetEnum() EnumValue {
	if v.Enum == nil {
		panic("value is not an enum")
	}
	return *v.Enum
}

func (v *Value) GetUnion() (UnionValue, bool) {
	if v.Union == nil {
		return UnionValue{}, false
	}
	return *v.Union, true
}

func (v *Value) MustGetUnion() UnionValue {
	if v.Union == nil {
		panic("value is not an union")
	}
	return *v.Union
}

func (v *Value) GetTupleValues() ([]Value, bool) {
	if v.TupleWith == nil {
		return TupleValues{}, false
	}
	return *v.TupleWith, true
}

func (v *Value) MustGetTupleValues() []Value {
	if v.TupleWith == nil {
		panic("value is not a tuple")
	}
	return *v.TupleWith
}

func (v *Value) GetCell() (boc.Cell, bool) {
	if v.Cell == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.Cell), true
}

func (v *Value) MustGetCell() boc.Cell {
	if v.Cell == nil {
		panic("value is not a cell")
	}
	return boc.Cell(*v.Cell)
}

func (v *Value) GetRemaining() (RemainingValue, bool) {
	if v.Remaining == nil {
		return RemainingValue{}, false
	}
	return *v.Remaining, true
}

func (v *Value) MustGetRemaining() RemainingValue {
	if v.Remaining == nil {
		panic("value is not a remaining")
	}
	return *v.Remaining
}

func (v *Value) GetType() string {
	return string(v.SumType)
}

func (v *Value) UnmarshalTyIdx(cell *boc.Cell, tyIdx int, decoder *Decoder) error {
	ty, err := decoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	return v.unmarshal(cell, ty, &tyIdx, decoder)
}

func (v *Value) Unmarshal(cell *boc.Cell, ty parser.Ty, decoder *Decoder) error {
	var tyIdx *int
	if decoder.abiIndex != nil {
		if idx, ok := decoder.abiIndex.TyIdxOf(ty); ok {
			tyIdx = &idx
		}
	}
	return v.unmarshal(cell, ty, tyIdx, decoder)
}

func (v *Value) unmarshal(cell *boc.Cell, ty parser.Ty, tyIdx *int, decoder *Decoder) error {
	var err error
	switch ty.SumType {
	case parser.TyKindIntN:
		if ty.IntN.N <= 64 {
			v.SumType = SumTypeSmallInt
			def := Int64(0)
			v.SmallInt = &def
			err = v.SmallInt.Unmarshal(cell, *ty.IntN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal small int value: %w", err)
			}
		} else {
			v.SumType = SumTypeBigInt
			v.BigInt = &BigInt{}
			err = v.BigInt.Unmarshal(cell, *ty.IntN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal big int value: %w", err)
			}
		}
	case parser.TyKindUintN:
		if ty.UintN.N <= 64 {
			v.SumType = SumTypeSmallUint
			def := UInt64(0)
			v.SmallUint = &def
			err = v.SmallUint.Unmarshal(cell, *ty.UintN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal small uint value: %w", err)
			}
		} else {
			v.SumType = SumTypeBigUint
			v.BigUint = &BigUInt{}
			err = v.BigUint.Unmarshal(cell, *ty.UintN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal big uint value: %w", err)
			}
		}
	case parser.TyKindVarIntN:
		v.SumType = SumTypeVarInt
		v.VarInt = &VarInt{}
		err = v.VarInt.Unmarshal(cell, *ty.VarIntN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal var int value: %w", err)
		}
	case parser.TyKindVarUintN:
		v.SumType = SumTypeVarUint
		v.VarUint = &VarUInt{}
		err = v.VarUint.Unmarshal(cell, *ty.VarUintN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal var uint value: %w", err)
		}
	case parser.TyKindBitsN:
		v.SumType = SumTypeBits
		v.Bits = &Bits{}
		err = v.Bits.Unmarshal(cell, *ty.BitsN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal bits value: %w", err)
		}
	case parser.TyKindNullable:
		v.SumType = SumTypeOptionalValue
		v.OptionalValue = &OptValue{}
		err = v.OptionalValue.Unmarshal(cell, *ty.Nullable, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nullable value: %w", err)
		}
	case parser.TyKindCellOf:
		v.SumType = SumTypeRefValue
		v.RefValue = &RefValue{}
		err = v.RefValue.Unmarshal(cell, *ty.CellOf, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cell of: %w", err)
		}
	case parser.TyKindTensor:
		v.SumType = SumTypeTensor
		v.Tensor = &TensorValues{}
		err = v.Tensor.Unmarshal(cell, *ty.Tensor, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tensor value: %w", err)
		}
	case parser.TyKindShapedTuple:
		v.SumType = SumTypeTupleWith
		v.TupleWith = &TupleValues{}
		err = v.TupleWith.Unmarshal(cell, *ty.ShapedTuple, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tuple value: %w", err)
		}
	case parser.TyKindMapKV:
		v.SumType = SumTypeMap
		v.Map = &MapValue{}
		err = v.Map.Unmarshal(cell, *ty.MapKV, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal map value: %w", err)
		}
	case parser.TyKindEnumRef:
		v.SumType = SumTypeEnum
		v.Enum = &EnumValue{}
		err = v.Enum.Unmarshal(cell, *ty.EnumRef, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal enum value: %w", err)
		}
	case parser.TyKindStructRef:
		v.SumType = SumTypeStruct
		v.Struct = &Struct{}
		if tyIdx == nil {
			return fmt.Errorf("struct %q requires ty_idx context", ty.StructRef.StructName)
		}
		err = v.Struct.UnmarshalTyIdx(cell, *tyIdx, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct value: %w", err)
		}
	case parser.TyKindAliasRef:
		v.SumType = SumTypeAlias
		v.Alias = &AliasValue{}
		if tyIdx == nil {
			return fmt.Errorf("alias %q requires ty_idx context", ty.AliasRef.AliasName)
		}
		err = v.Alias.UnmarshalTyIdx(cell, *tyIdx, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal alias value: %w", err)
		}
	case parser.TyKindGenericT:
		v.SumType = SumTypeGenericT
		v.Generic = &GenericValue{}
		err = v.Generic.Unmarshal(cell, *ty.GenericT, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal generic value: %w", err)
		}
	case parser.TyKindUnion:
		v.SumType = SumTypeUnion
		v.Union = &UnionValue{}
		err = v.Union.Unmarshal(cell, *ty.Union, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal union value: %w", err)
		}
	case parser.TyKindInt:
		err = fmt.Errorf("failed to unmarshal int value: int is not supported")
	case parser.TyKindCoins:
		v.SumType = SumTypeCoins
		v.Coins = &CoinsValue{}
		err = v.Coins.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal coins value: %w", err)
		}
	case parser.TyKindBool:
		v.SumType = SumTypeBool
		def := BoolValue(false)
		v.Bool = &def
		err = v.Bool.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal bool value: %w", err)
		}
	case parser.TyKindCell:
		v.SumType = SumTypeCell
		v.Cell = &Any{}
		err = v.Cell.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cell value: %w", err)
		}
	case parser.TyKindSlice:
		err = fmt.Errorf("failed to unmarshal slice value: slice is not supported")
	case parser.TyKindBuilder:
		err = fmt.Errorf("failed to unmarshal builder value: builder is not supported")
	case parser.TyKindCallable:
		err = fmt.Errorf("failed to unmarshal callable value: callable is not supported")
	case parser.TyKindRemaining:
		v.SumType = SumTypeRemaining
		v.Remaining = &RemainingValue{}
		err = v.Remaining.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal remaining value: %w", err)
		}
	case parser.TyKindAddress:
		v.SumType = SumTypeInternalAddress
		v.InternalAddress = &InternalAddress{}
		err = v.InternalAddress.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal internal address value: %w", err)
		}
	case parser.TyKindAddressOpt:
		v.SumType = SumTypeOptionalAddress
		v.OptionalAddress = &OptionalAddress{}
		err = v.OptionalAddress.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal optional address value: %w", err)
		}
	case parser.TyKindAddressExt:
		v.SumType = SumTypeExternalAddress
		v.ExternalAddress = &ExternalAddress{}
		err = v.ExternalAddress.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal external address value: %w", err)
		}
	case parser.TyKindAddressAny:
		v.SumType = SumTypeAnyAddress
		v.AnyAddress = &AnyAddress{}
		err = v.AnyAddress.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal any address value: %w", err)
		}
	case parser.TyKindNullLiteral:
		v.SumType = SumTypeNull
		v.Null = &NullValue{}
		err = v.Null.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal null value: %w", err)
		}
	case parser.TyKindVoid:
		v.SumType = SumTypeVoid
		v.Void = &VoidValue{}
		err = v.Void.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal void value: %w", err)
		}
	default:
		return fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	return nil
}

func (v *Value) MarshalTyIdx(cell *boc.Cell, tyIdx int, encoder *Encoder) error {
	ty, err := encoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	return v.marshal(cell, ty, &tyIdx, encoder)
}

func (v *Value) Marshal(cell *boc.Cell, ty parser.Ty, encoder *Encoder) error {
	var tyIdx *int
	if encoder.abiIndex != nil {
		if idx, ok := encoder.abiIndex.TyIdxOf(ty); ok {
			tyIdx = &idx
		}
	}
	return v.marshal(cell, ty, tyIdx, encoder)
}

func (v *Value) marshal(cell *boc.Cell, ty parser.Ty, tyIdx *int, encoder *Encoder) error {
	var err error
	switch ty.SumType {
	case parser.TyKindIntN:
		if ty.IntN.N <= 64 {
			if v.SumType != SumTypeSmallInt {
				return fmt.Errorf("expected SmallInt, but got %v", v.SumType)
			}
			err = v.SmallInt.Marshal(cell, *ty.IntN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal small int value: %w", err)
			}
		} else {
			if v.SumType != SumTypeBigInt {
				return fmt.Errorf("expected BigInt, but got %v", v.SumType)
			}
			err = v.BigInt.Marshal(cell, *ty.IntN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal big int value: %w", err)
			}
		}
	case parser.TyKindUintN:
		if ty.UintN.N <= 64 {
			if v.SumType != SumTypeSmallUint {
				return fmt.Errorf("expected SmallUint, but got %v", v.SumType)
			}
			err = v.SmallUint.Marshal(cell, *ty.UintN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal small uint value: %w", err)
			}
		} else {
			if v.SumType != SumTypeBigUint {
				return fmt.Errorf("expected BigUint, but got %v", v.SumType)
			}
			err = v.BigUint.Marshal(cell, *ty.UintN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal big uint value: %w", err)
			}
		}
	case parser.TyKindVarIntN:
		if v.SumType != SumTypeVarInt {
			return fmt.Errorf("expected VarInt, but got %v", v.SumType)
		}
		err = v.VarInt.Marshal(cell, *ty.VarIntN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal var int value: %w", err)
		}
	case parser.TyKindVarUintN:
		if v.SumType != SumTypeVarUint {
			return fmt.Errorf("expected VarUint, but got %v", v.SumType)
		}
		err = v.VarUint.Marshal(cell, *ty.VarUintN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal var uint value: %w", err)
		}
	case parser.TyKindBitsN:
		if v.SumType != SumTypeBits {
			return fmt.Errorf("expected Bits, but got %v", v.SumType)
		}
		err = v.Bits.Marshal(cell, *ty.BitsN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal bits value: %w", err)
		}
	case parser.TyKindNullable:
		if v.SumType != SumTypeOptionalValue {
			return fmt.Errorf("expected OptionalValue, but got %v", v.SumType)
		}
		err = v.OptionalValue.Marshal(cell, *ty.Nullable, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal nullable value: %w", err)
		}
	case parser.TyKindCellOf:
		if v.SumType != SumTypeRefValue {
			return fmt.Errorf("expected RefValue, but got %v", v.SumType)
		}
		err = v.RefValue.Marshal(cell, *ty.CellOf, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal cell of value: %w", err)
		}
	case parser.TyKindTensor:
		if v.SumType != SumTypeTensor {
			return fmt.Errorf("expected Tensor, but got %v", v.SumType)
		}
		err = v.Tensor.Marshal(cell, *ty.Tensor, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal tensor value: %w", err)
		}
	case parser.TyKindShapedTuple:
		if v.SumType != SumTypeTupleWith {
			return fmt.Errorf("expected TupleWith, but got %v", v.SumType)
		}
		err = v.TupleWith.Marshal(cell, *ty.ShapedTuple, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal tuple value: %w", err)
		}
	case parser.TyKindMapKV:
		if v.SumType != SumTypeMap {
			return fmt.Errorf("expected Map, but got %v", v.SumType)
		}
		err = v.Map.Marshal(cell, *ty.MapKV, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal map value: %w", err)
		}
	case parser.TyKindEnumRef:
		if v.SumType != SumTypeEnum {
			return fmt.Errorf("expected Enum, but got %v", v.SumType)
		}
		err = v.Enum.Marshal(cell, *ty.EnumRef, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal enum value: %w", err)
		}
	case parser.TyKindStructRef:
		if v.SumType != SumTypeStruct {
			return fmt.Errorf("expected Struct, but got %v", v.SumType)
		}
		if tyIdx == nil {
			return fmt.Errorf("struct %q requires ty_idx context", ty.StructRef.StructName)
		}
		err = v.Struct.MarshalTyIdx(cell, *tyIdx, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal struct value: %w", err)
		}
	case parser.TyKindAliasRef:
		if v.SumType != SumTypeAlias {
			return fmt.Errorf("expected Alias, but got %v", v.SumType)
		}
		if tyIdx == nil {
			return fmt.Errorf("alias %q requires ty_idx context", ty.AliasRef.AliasName)
		}
		err = v.Alias.MarshalTyIdx(cell, *tyIdx, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal alias value: %w", err)
		}
	case parser.TyKindGenericT:
		if v.SumType != SumTypeGenericT {
			return fmt.Errorf("expected GenericT, but got %v", v.SumType)
		}
		err = v.Generic.Marshal(cell, *ty.GenericT, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal generic value: %w", err)
		}
	case parser.TyKindUnion:
		if v.SumType != SumTypeUnion {
			return fmt.Errorf("expected Union, but got %v", v.SumType)
		}
		err = v.Union.Marshal(cell, *ty.Union, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal union value: %w", err)
		}
	case parser.TyKindInt:
		err = fmt.Errorf("failed to marshal int value: int is not supported")
	case parser.TyKindCoins:
		if v.SumType != SumTypeCoins {
			return fmt.Errorf("expected Coins, but got %v", v.SumType)
		}
		err = v.Coins.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal coins value: %w", err)
		}
	case parser.TyKindBool:
		if v.SumType != SumTypeBool {
			return fmt.Errorf("expected Bool, but got %v", v.SumType)
		}
		err = v.Bool.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal bool value: %w", err)
		}
	case parser.TyKindCell:
		if v.SumType != SumTypeCell {
			return fmt.Errorf("expected Cell, but got %v", v.SumType)
		}
		err = v.Cell.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal cell value: %w", err)
		}
	case parser.TyKindSlice:
		err = fmt.Errorf("failed to marshal slice value: slice is not supported")
	case parser.TyKindBuilder:
		err = fmt.Errorf("failed to marshal builder value: builder is not supported")
	case parser.TyKindCallable:
		err = fmt.Errorf("failed to marshal int callable: callable is not supported")
	case parser.TyKindRemaining:
		if v.SumType != SumTypeRemaining {
			return fmt.Errorf("expected Remaining, but got %v", v.SumType)
		}
		err = v.Remaining.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal remaining value: %w", err)
		}
	case parser.TyKindAddress:
		if v.SumType != SumTypeInternalAddress {
			return fmt.Errorf("expected InternalAddress, but got %v", v.SumType)
		}
		err = v.InternalAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal internal address value: %w", err)
		}
	case parser.TyKindAddressOpt:
		if v.SumType != SumTypeOptionalAddress {
			return fmt.Errorf("expected OptionalAddress, but got %v", v.SumType)
		}
		err = v.OptionalAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address value: %w", err)
		}
	case parser.TyKindAddressExt:
		if v.SumType != SumTypeExternalAddress {
			return fmt.Errorf("expected ExternalAddress, but got %v", v.SumType)
		}
		err = v.ExternalAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal external address value: %w", err)
		}
	case parser.TyKindAddressAny:
		if v.SumType != SumTypeAnyAddress {
			return fmt.Errorf("expected AnyAddress, but got %v", v.SumType)
		}
		err = v.AnyAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal any address value: %w", err)
		}
	case parser.TyKindNullLiteral:
		if v.SumType != SumTypeNull {
			return fmt.Errorf("expected Null, but got %v", v.SumType)
		}
		err = v.Null.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal null value: %w", err)
		}
	case parser.TyKindVoid:
		if v.SumType != SumTypeVoid {
			return fmt.Errorf("expected Void, but got %v", v.SumType)
		}
		err = v.Void.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal void value: %w", err)
		}
	default:
		err = fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	return nil
}

func (v *Value) Equal(o any) bool {
	otherValue, ok := o.(Value)
	if !ok {
		return false
	}

	switch v.SumType {
	case SumTypeBool:
		if otherValue.Bool == nil {
			return false
		}
		return v.Bool.Equal(*otherValue.Bool)
	case SumTypeSmallInt:
		if otherValue.SmallInt == nil {
			return false
		}
		return v.SmallInt.Equal(*otherValue.SmallInt)
	case SumTypeSmallUint:
		if otherValue.SmallUint == nil {
			return false
		}
		return v.SmallUint.Equal(*otherValue.SmallUint)
	case SumTypeBigInt:
		if otherValue.BigInt == nil {
			return false
		}
		return v.BigInt.Equal(*otherValue.BigInt)
	case SumTypeBigUint:
		if otherValue.BigUint == nil {
			return false
		}
		return v.BigUint.Equal(*otherValue.BigUint)
	case SumTypeVarInt:
		if otherValue.VarInt == nil {
			return false
		}
		return v.VarInt.Equal(*otherValue.VarInt)
	case SumTypeVarUint:
		if otherValue.VarUint == nil {
			return false
		}
		return v.VarUint.Equal(*otherValue.VarUint)
	case SumTypeCoins:
		if otherValue.Coins == nil {
			return false
		}
		return v.Coins.Equal(*otherValue.Coins)
	case SumTypeBits:
		if otherValue.Bits == nil {
			return false
		}
		return v.Bits.Equal(*otherValue.Bits)
	case SumTypeCell:
		if otherValue.Cell == nil {
			return false
		}
		return v.Cell.Equal(*otherValue.Cell)
	case SumTypeRemaining:
		if otherValue.Remaining == nil {
			return false
		}
		return v.Remaining.Equal(*otherValue.Remaining)
	case SumTypeInternalAddress:
		if otherValue.InternalAddress == nil {
			return false
		}
		return v.InternalAddress.Equal(*otherValue.InternalAddress)
	case SumTypeOptionalAddress:
		if otherValue.OptionalAddress == nil {
			return false
		}
		return v.OptionalAddress.Equal(*otherValue.OptionalAddress)
	case SumTypeExternalAddress:
		if otherValue.ExternalAddress == nil {
			return false
		}
		return v.ExternalAddress.Equal(*otherValue.ExternalAddress)
	case SumTypeAnyAddress:
		if otherValue.AnyAddress == nil {
			return false
		}
		return v.AnyAddress.Equal(*otherValue.AnyAddress)
	case SumTypeOptionalValue:
		if otherValue.OptionalValue == nil {
			return false
		}
		return v.OptionalValue.Equal(*otherValue.OptionalValue)
	case SumTypeRefValue:
		if otherValue.RefValue == nil {
			return false
		}
		return v.RefValue.Equal(*otherValue.RefValue)
	case SumTypeTupleWith:
		if otherValue.TupleWith == nil {
			return false
		}
		return v.TupleWith.Equal(*otherValue.TupleWith)
	case SumTypeTensor:
		if otherValue.Tensor == nil {
			return false
		}
		return v.Tensor.Equal(*otherValue.Tensor)
	case SumTypeMap:
		if otherValue.Map == nil {
			return false
		}
		return v.Map.Equal(*otherValue.Map)
	case SumTypeStruct:
		if otherValue.Struct == nil {
			return false
		}
		return v.Struct.Equal(*otherValue.Struct)
	case SumTypeAlias:
		if otherValue.Alias == nil {
			return false
		}
		return v.Alias.Equal(*otherValue.Alias)
	case SumTypeEnum:
		if otherValue.Enum == nil {
			return false
		}
		return v.Enum.Equal(*otherValue.Enum)
	case SumTypeGenericT:
		if otherValue.Generic == nil {
			return false
		}
		return v.Generic.Equal(*otherValue.Generic)
	case SumTypeUnion:
		if otherValue.Union == nil {
			return false
		}
		return v.Union.Equal(*otherValue.Union)
	case SumTypeNull:
		if otherValue.Null == nil {
			return false
		}
		return v.Null.Equal(*otherValue.Null)
	case SumTypeVoid:
		if otherValue.Void == nil {
			return false
		}
		return v.Void.Equal(*otherValue.Void)
	default:
		return false
	}
}

func (v Value) MarshalJSON() ([]byte, error) {
	var s strings.Builder
	var data []byte
	var err error
	switch v.SumType {
	case SumTypeBool:
		data, err = json.Marshal(v.Bool)
	case SumTypeSmallInt:
		data, err = json.Marshal(v.SmallInt)
	case SumTypeSmallUint:
		data, err = json.Marshal(v.SmallUint)
	case SumTypeBigInt:
		data, err = json.Marshal(v.BigInt)
	case SumTypeBigUint:
		data, err = json.Marshal(v.BigUint)
	case SumTypeVarInt:
		data, err = json.Marshal(v.VarInt)
	case SumTypeVarUint:
		data, err = json.Marshal(v.VarUint)
	case SumTypeCoins:
		data, err = json.Marshal(v.Coins)
	case SumTypeBits:
		data, err = json.Marshal(v.Bits)
	case SumTypeCell:
		data, err = json.Marshal(v.Cell)
	case SumTypeRemaining:
		data, err = json.Marshal(v.Remaining)
	case SumTypeInternalAddress:
		data, err = json.Marshal(v.InternalAddress)
	case SumTypeOptionalAddress:
		data, err = json.Marshal(v.OptionalAddress)
	case SumTypeExternalAddress:
		data, err = json.Marshal(v.ExternalAddress)
	case SumTypeAnyAddress:
		data, err = json.Marshal(v.AnyAddress)
	case SumTypeOptionalValue:
		data, err = json.Marshal(v.OptionalValue)
	case SumTypeRefValue:
		data, err = json.Marshal(v.RefValue)
	case SumTypeTupleWith:
		data, err = json.Marshal(v.TupleWith)
	case SumTypeTensor:
		data, err = json.Marshal(v.Tensor)
	case SumTypeMap:
		data, err = json.Marshal(v.Map)
	case SumTypeStruct:
		data, err = json.Marshal(v.Struct)
	case SumTypeAlias:
		val := Value(*v.Alias)
		data, err = json.Marshal(val)
	case SumTypeEnum:
		data, err = json.Marshal(v.Enum)
	case SumTypeGenericT:
		val := Value(*v.Generic)
		data, err = json.Marshal(val)
	case SumTypeUnion:
		data, err = json.Marshal(v.Union)
	case SumTypeNull:
		data, err = json.Marshal(v.Null)
	case SumTypeVoid:
		data, err = json.Marshal(v.Void)
	default:
		err = fmt.Errorf("unknown value type: %s", v.SumType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	s.Write(data)
	return []byte(s.String()), err
}
