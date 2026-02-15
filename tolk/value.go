package tolk

// todo: move this to some package or rename somehow.

import (
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

type SumType string

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
	SumType         SumType          `json:"sumType"`
	Bool            *BoolValue       `json:"bool,omitempty"`
	SmallInt        *Int64           `json:"smallInt,omitempty"`
	SmallUint       *UInt64          `json:"smallUint,omitempty"`
	BigInt          *BigInt          `json:"bigInt,omitempty"`
	BigUint         *BigUInt         `json:"bigUint,omitempty"`
	VarInt          *VarInt          `json:"varInt,omitempty"`
	VarUint         *VarUInt         `json:"varUint,omitempty"`
	Coins           *CoinsValue      `json:"coins,omitempty"`
	Bits            *Bits            `json:"bits,omitempty"`
	Cell            *Any             `json:"cell,omitempty"`
	Remaining       *RemainingValue  `json:"remaining,omitempty"`
	InternalAddress *InternalAddress `json:"internalAddress,omitempty"`
	OptionalAddress *OptionalAddress `json:"optionalAddress,omitempty"`
	ExternalAddress *ExternalAddress `json:"externalAddress,omitempty"`
	AnyAddress      *AnyAddress      `json:"anyAddress,omitempty"`
	OptionalValue   *OptValue        `json:"optionalValue,omitempty"`
	RefValue        *RefValue        `json:"refValue,omitempty"`
	TupleWith       *TupleValues     `json:"tupleWith,omitempty"`
	Tensor          *TensorValues    `json:"tensor,omitempty"`
	Map             *MapValue        `json:"map,omitempty"`
	Struct          *Struct          `json:"struct,omitempty"`
	Alias           *AliasValue      `json:"alias,omitempty"`
	Enum            *EnumValue       `json:"enum,omitempty"`
	Generic         *GenericValue    `json:"generic,omitempty"`
	Union           *UnionValue      `json:"union,omitempty"`
	Null            *NullValue       `json:"null,omitempty"`
	Void            *VoidValue       `json:"void,omitempty"`
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

func (v *Value) GetRemaining() (boc.Cell, bool) {
	if v.Remaining == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.Remaining), true
}

func (v *Value) MustGetRemaining() boc.Cell {
	if v.Remaining == nil {
		panic("value is not a remaining")
	}
	return boc.Cell(*v.Remaining)
}

func (v *Value) GetType() string {
	return string(v.SumType)
}

func (v *Value) Unmarshal(cell *boc.Cell, ty tolkParser.Ty, decoder *Decoder) error {
	var err error
	switch ty.SumType {
	case "IntN":
		if ty.IntN.N <= 64 {
			v.SumType = "SmallInt"
			def := Int64(0)
			v.SmallInt = &def
			err = v.SmallInt.Unmarshal(cell, *ty.IntN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal small int value: %w", err)
			}
		} else {
			v.SumType = "BigInt"
			v.BigInt = &BigInt{}
			err = v.BigInt.Unmarshal(cell, *ty.IntN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal big int value: %w", err)
			}
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			v.SumType = "SmallUint"
			def := UInt64(0)
			v.SmallUint = &def
			err = v.SmallUint.Unmarshal(cell, *ty.UintN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal small uint value: %w", err)
			}
		} else {
			v.SumType = "BigUint"
			v.BigUint = &BigUInt{}
			err = v.BigUint.Unmarshal(cell, *ty.UintN, decoder)
			if err != nil {
				return fmt.Errorf("failed to unmarshal big uint value: %w", err)
			}
		}
	case "VarIntN":
		v.SumType = "VarInt"
		v.VarInt = &VarInt{}
		err = v.VarInt.Unmarshal(cell, *ty.VarIntN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal var int value: %w", err)
		}
	case "VarUintN":
		v.SumType = "VarUint"
		v.VarUint = &VarUInt{}
		err = v.VarUint.Unmarshal(cell, *ty.VarUintN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal var uint value: %w", err)
		}
	case "BitsN":
		v.SumType = "Bits"
		v.Bits = &Bits{}
		err = v.Bits.Unmarshal(cell, *ty.BitsN, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal bits value: %w", err)
		}
	case "Nullable":
		v.SumType = "OptionalValue"
		v.OptionalValue = &OptValue{}
		err = v.OptionalValue.Unmarshal(cell, *ty.Nullable, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nullable value: %w", err)
		}
	case "CellOf":
		v.SumType = "RefValue"
		v.RefValue = &RefValue{}
		err = v.RefValue.Unmarshal(cell, *ty.CellOf, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cell of: %w", err)
		}
	case "Tensor":
		v.SumType = "Tensor"
		v.Tensor = &TensorValues{}
		err = v.Tensor.Unmarshal(cell, *ty.Tensor, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tensor value: %w", err)
		}
	case "TupleWith":
		v.SumType = "TupleWith"
		v.TupleWith = &TupleValues{}
		err = v.TupleWith.Unmarshal(cell, *ty.TupleWith, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tuple value: %w", err)
		}
	case "Map":
		v.SumType = "Map"
		v.Map = &MapValue{}
		err = v.Map.Unmarshal(cell, *ty.Map, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal map value: %w", err)
		}
	case "EnumRef":
		v.SumType = "Enum"
		v.Enum = &EnumValue{}
		err = v.Enum.Unmarshal(cell, *ty.EnumRef, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal enum value: %w", err)
		}
	case "StructRef":
		v.SumType = "Struct"
		v.Struct = &Struct{}
		err = v.Struct.Unmarshal(cell, *ty.StructRef, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct value: %w", err)
		}
	case "AliasRef":
		v.SumType = "Alias"
		v.Alias = &AliasValue{}
		err = v.Alias.Unmarshal(cell, *ty.AliasRef, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal alias value: %w", err)
		}
	case "Generic":
		v.SumType = "Generic"
		v.Generic = &GenericValue{}
		err = v.Generic.Unmarshal(cell, *ty.Generic, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal generic value: %w", err)
		}
	case "Union":
		v.SumType = "Union"
		v.Union = &UnionValue{}
		err = v.Union.Unmarshal(cell, *ty.Union, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal union value: %w", err)
		}
	case "Int":
		err = fmt.Errorf("failed to unmarshal int value: int is not supported")
	case "Coins":
		v.SumType = "Coins"
		v.Coins = &CoinsValue{}
		err = v.Coins.Unmarshal(cell, *ty.Coins, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal coins value: %w", err)
		}
	case "Bool":
		v.SumType = "Bool"
		def := BoolValue(false)
		v.Bool = &def
		err = v.Bool.Unmarshal(cell, *ty.Bool, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal bool value: %w", err)
		}
	case "Cell":
		v.SumType = "Cell"
		v.Cell = &Any{}
		err = v.Cell.Unmarshal(cell, *ty.Cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cell value: %w", err)
		}
	case "Slice":
		err = fmt.Errorf("failed to unmarshal slice value: slice is not supported")
	case "Builder":
		err = fmt.Errorf("failed to unmarshal builder value: builder is not supported")
	case "Callable":
		err = fmt.Errorf("failed to unmarshal callable value: callable is not supported")
	case "Remaining":
		v.SumType = "Remaining"
		v.Remaining = &RemainingValue{}
		err = v.Remaining.Unmarshal(cell, *ty.Remaining, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal remaining value: %w", err)
		}
	case "Address":
		v.SumType = "InternalAddress"
		v.InternalAddress = &InternalAddress{}
		err = v.InternalAddress.Unmarshal(cell, *ty.Address, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal internal address value: %w", err)
		}
	case "AddressOpt":
		v.SumType = "OptionalAddress"
		v.OptionalAddress = &OptionalAddress{}
		err = v.OptionalAddress.Unmarshal(cell, *ty.AddressOpt, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal optional address value: %w", err)
		}
	case "AddressExt":
		v.SumType = "ExternalAddress"
		v.ExternalAddress = &ExternalAddress{}
		err = v.ExternalAddress.Unmarshal(cell, *ty.AddressExt, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal external address value: %w", err)
		}
	case "AddressAny":
		v.SumType = "AnyAddress"
		v.AnyAddress = &AnyAddress{}
		err = v.AnyAddress.Unmarshal(cell, *ty.AddressAny, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal any address value: %w", err)
		}
	case "TupleAny":
		err = fmt.Errorf("failed to unmarshal tuple any value: tuple any is not supported")
	case "NullLiteral":
		v.SumType = "Null"
		v.Null = &NullValue{}
		err = v.Null.Unmarshal(cell, *ty.NullLiteral, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal null value: %w", err)
		}
	case "Void":
		v.SumType = "Void"
		v.Void = &VoidValue{}
		err = v.Void.Unmarshal(cell, *ty.Void, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal void value: %w", err)
		}
	default:
		return fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	return nil
}

func (v *Value) Marshal(cell *boc.Cell, ty tolkParser.Ty, encoder *Encoder) error {
	var err error
	switch ty.SumType {
	case "IntN":
		if ty.IntN.N <= 64 {
			if v.SumType != "SmallInt" {
				return fmt.Errorf("expected SmallInt, but got %v", v.SumType)
			}
			err = v.SmallInt.Marshal(cell, *ty.IntN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal small int value: %w", err)
			}
		} else {
			if v.SumType != "BigInt" {
				return fmt.Errorf("expected BigInt, but got %v", v.SumType)
			}
			err = v.BigInt.Marshal(cell, *ty.IntN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal big int value: %w", err)
			}
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			if v.SumType != "SmallUint" {
				return fmt.Errorf("expected SmallUint, but got %v", v.SumType)
			}
			err = v.SmallUint.Marshal(cell, *ty.UintN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal small uint value: %w", err)
			}
		} else {
			if v.SumType != "BigUint" {
				return fmt.Errorf("expected BigUint, but got %v", v.SumType)
			}
			err = v.BigUint.Marshal(cell, *ty.UintN, encoder)
			if err != nil {
				return fmt.Errorf("failed to marshal big uint value: %w", err)
			}
		}
	case "VarIntN":
		if v.SumType != "VarInt" {
			return fmt.Errorf("expected VarInt, but got %v", v.SumType)
		}
		err = v.VarInt.Marshal(cell, *ty.VarIntN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal var int value: %w", err)
		}
	case "VarUintN":
		if v.SumType != "VarUint" {
			return fmt.Errorf("expected VarUint, but got %v", v.SumType)
		}
		err = v.VarUint.Marshal(cell, *ty.VarUintN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal var uint value: %w", err)
		}
	case "BitsN":
		if v.SumType != "Bits" {
			return fmt.Errorf("expected Bits, but got %v", v.SumType)
		}
		err = v.Bits.Marshal(cell, *ty.BitsN, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal bits value: %w", err)
		}
	case "Nullable":
		if v.SumType != "OptionalValue" {
			return fmt.Errorf("expected OptionalValue, but got %v", v.SumType)
		}
		err = v.OptionalValue.Marshal(cell, *ty.Nullable, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal nullable value: %w", err)
		}
	case "CellOf":
		if v.SumType != "RefValue" {
			return fmt.Errorf("expected RefValue, but got %v", v.SumType)
		}
		err = v.RefValue.Marshal(cell, *ty.CellOf, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal cell of value: %w", err)
		}
	case "Tensor":
		if v.SumType != "Tensor" {
			return fmt.Errorf("expected Tensor, but got %v", v.SumType)
		}
		err = v.Tensor.Marshal(cell, *ty.Tensor, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal tensor value: %w", err)
		}
	case "TupleWith":
		if v.SumType != "TupleWith" {
			return fmt.Errorf("expected TupleWith, but got %v", v.SumType)
		}
		err = v.TupleWith.Marshal(cell, *ty.TupleWith, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal tuple value: %w", err)
		}
	case "Map":
		if v.SumType != "Map" {
			return fmt.Errorf("expected Map, but got %v", v.SumType)
		}
		err = v.Map.Marshal(cell, *ty.Map, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal map value: %w", err)
		}
	case "EnumRef":
		if v.SumType != "Enum" {
			return fmt.Errorf("expected Enum, but got %v", v.SumType)
		}
		err = v.Enum.Marshal(cell, *ty.EnumRef, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal enum value: %w", err)
		}
	case "StructRef":
		if v.SumType != "Struct" {
			return fmt.Errorf("expected Struct, but got %v", v.SumType)
		}
		err = v.Struct.Marshal(cell, *ty.StructRef, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal struct value: %w", err)
		}
	case "AliasRef":
		if v.SumType != "Alias" {
			return fmt.Errorf("expected Alias, but got %v", v.SumType)
		}
		err = v.Alias.Marshal(cell, *ty.AliasRef, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal alias value: %w", err)
		}
	case "Generic":
		if v.SumType != "Generic" {
			return fmt.Errorf("expected Generic, but got %v", v.SumType)
		}
		err = v.Generic.Marshal(cell, *ty.Generic, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal generic value: %w", err)
		}
	case "Union":
		if v.SumType != "Union" {
			return fmt.Errorf("expected Union, but got %v", v.SumType)
		}
		err = v.Union.Marshal(cell, *ty.Union, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal union value: %w", err)
		}
	case "Int":
		err = fmt.Errorf("failed to marshal int value: int is not supported")
	case "Coins":
		if v.SumType != "Coins" {
			return fmt.Errorf("expected Coins, but got %v", v.SumType)
		}
		err = v.Coins.Marshal(cell, *ty.Coins, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal coins value: %w", err)
		}
	case "Bool":
		if v.SumType != "Bool" {
			return fmt.Errorf("expected Bool, but got %v", v.SumType)
		}
		err = v.Bool.Marshal(cell, *ty.Bool, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal bool value: %w", err)
		}
	case "Cell":
		if v.SumType != "Cell" {
			return fmt.Errorf("expected Cell, but got %v", v.SumType)
		}
		err = v.Cell.Marshal(cell, *ty.Cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal cell value: %w", err)
		}
	case "Slice":
		err = fmt.Errorf("failed to marshal slice value: slice is not supported")
	case "Builder":
		err = fmt.Errorf("failed to marshal builder value: builder is not supported")
	case "Callable":
		err = fmt.Errorf("failed to marshal int callable: callable is not supported")
	case "Remaining":
		if v.SumType != "Remaining" {
			return fmt.Errorf("expected Remaining, but got %v", v.SumType)
		}
		err = v.Remaining.Marshal(cell, *ty.Remaining, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal remaining value: %w", err)
		}
	case "Address":
		if v.SumType != "InternalAddress" {
			return fmt.Errorf("expected InternalAddress, but got %v", v.SumType)
		}
		err = v.InternalAddress.Marshal(cell, *ty.Address, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal internal address value: %w", err)
		}
	case "AddressOpt":
		if v.SumType != "OptionalAddress" {
			return fmt.Errorf("expected OptionalAddress, but got %v", v.SumType)
		}
		err = v.OptionalAddress.Marshal(cell, *ty.AddressOpt, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address value: %w", err)
		}
	case "AddressExt":
		if v.SumType != "ExternalAddress" {
			return fmt.Errorf("expected ExternalAddress, but got %v", v.SumType)
		}
		err = v.ExternalAddress.Marshal(cell, *ty.AddressExt, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal external address value: %w", err)
		}
	case "AddressAny":
		if v.SumType != "AnyAddress" {
			return fmt.Errorf("expected AnyAddress, but got %v", v.SumType)
		}
		err = v.AnyAddress.Marshal(cell, *ty.AddressAny, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal any address value: %w", err)
		}
	case "TupleAny":
		err = fmt.Errorf("failed to marshal tuple any value: tuple any not supported")
	case "NullLiteral":
		if v.SumType != "Null" {
			return fmt.Errorf("expected Null, but got %v", v.SumType)
		}
		err = v.Null.Marshal(cell, *ty.NullLiteral, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal null value: %w", err)
		}
	case "Void":
		if v.SumType != "Void" {
			return fmt.Errorf("expected Void, but got %v", v.SumType)
		}
		err = v.Void.Marshal(cell, *ty.Void, encoder)
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
	case "Bool":
		if otherValue.Bool == nil {
			return false
		}
		return v.Bool.Equal(*otherValue.Bool)
	case "SmallInt":
		if otherValue.SmallInt == nil {
			return false
		}
		return v.SmallInt.Equal(*otherValue.SmallInt)
	case "SmallUint":
		if otherValue.SmallUint == nil {
			return false
		}
		return v.SmallUint.Equal(*otherValue.SmallUint)
	case "BigInt":
		if otherValue.BigInt == nil {
			return false
		}
		return v.BigInt.Equal(*otherValue.BigInt)
	case "BigUint":
		if otherValue.BigUint == nil {
			return false
		}
		return v.BigUint.Equal(*otherValue.BigUint)
	case "VarInt":
		if otherValue.VarInt == nil {
			return false
		}
		return v.VarInt.Equal(*otherValue.VarInt)
	case "VarUint":
		if otherValue.VarUint == nil {
			return false
		}
		return v.VarUint.Equal(*otherValue.VarUint)
	case "Coins":
		if otherValue.Coins == nil {
			return false
		}
		return v.Coins.Equal(*otherValue.Coins)
	case "Bits":
		if otherValue.Bits == nil {
			return false
		}
		return v.Bits.Equal(*otherValue.Bits)
	case "Cell":
		if otherValue.Cell == nil {
			return false
		}
		return v.Cell.Equal(*otherValue.Cell)
	case "Remaining":
		if otherValue.Remaining == nil {
			return false
		}
		return v.Remaining.Equal(*otherValue.Remaining)
	case "InternalAddress":
		if otherValue.InternalAddress == nil {
			return false
		}
		return v.InternalAddress.Equal(*otherValue.InternalAddress)
	case "OptionalAddress":
		if otherValue.OptionalAddress == nil {
			return false
		}
		return v.OptionalAddress.Equal(*otherValue.OptionalAddress)
	case "ExternalAddress":
		if otherValue.ExternalAddress == nil {
			return false
		}
		return v.ExternalAddress.Equal(*otherValue.ExternalAddress)
	case "AnyAddress":
		if otherValue.AnyAddress == nil {
			return false
		}
		return v.AnyAddress.Equal(*otherValue.AnyAddress)
	case "OptionalValue":
		if otherValue.OptionalValue == nil {
			return false
		}
		return v.OptionalValue.Equal(*otherValue.OptionalValue)
	case "RefValue":
		if otherValue.RefValue == nil {
			return false
		}
		return v.RefValue.Equal(*otherValue.RefValue)
	case "TupleWith":
		if otherValue.TupleWith == nil {
			return false
		}
		return v.TupleWith.Equal(*otherValue.TupleWith)
	case "Tensor":
		if otherValue.Tensor == nil {
			return false
		}
		return v.Tensor.Equal(*otherValue.Tensor)
	case "Map":
		if otherValue.Map == nil {
			return false
		}
		return v.Map.Equal(*otherValue.Map)
	case "Struct":
		if otherValue.Struct == nil {
			return false
		}
		return v.Struct.Equal(*otherValue.Struct)
	case "Alias":
		if otherValue.Alias == nil {
			return false
		}
		return v.Alias.Equal(*otherValue.Alias)
	case "Enum":
		if otherValue.Enum == nil {
			return false
		}
		return v.Enum.Equal(*otherValue.Enum)
	case "Generic":
		if otherValue.Generic == nil {
			return false
		}
		return v.Generic.Equal(*otherValue.Generic)
	case "Union":
		if otherValue.Union == nil {
			return false
		}
		return v.Union.Equal(*otherValue.Union)
	case "Null":
		if otherValue.Null == nil {
			return false
		}
		return v.Null.Equal(*otherValue.Null)
	case "Void":
		if otherValue.Void == nil {
			return false
		}
		return v.Void.Equal(*otherValue.Void)
	default:
		return false
	}
}
