package tolk

// todo: move this to some package or rename somehow.

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/ton"
)

type SumType string

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

func (v *Value) SetValue(val any, ty tolkParser.Ty) error {
	switch ty.SumType {
	case "IntN":
		bi, ok := val.(BigInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to BigInt", val)
		}
		if ty.IntN.N <= 64 {
			b := big.Int(bi)
			wVal := Int64(b.Int64())
			v.SmallInt = &wVal
		} else {
			v.BigInt = &bi
		}
	case "UintN":
		bi, ok := val.(BigUInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to BigUInt", val)
		}
		if ty.UintN.N <= 64 {
			b := big.Int(bi)
			wVal := UInt64(b.Uint64())
			v.SmallUint = &wVal
		} else {
			v.BigUint = &bi
		}
	case "VarIntN":
		vi, ok := val.(VarInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to VarInt", val)
		}
		v.VarInt = &vi
	case "VarUintN":
		vi, ok := val.(VarUInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to VarUInt", val)
		}
		v.VarUint = &vi
	case "BitsN":
		b, ok := val.(Bits)
		if !ok {
			return fmt.Errorf("cannot convert %v to Bits", val)
		}
		v.Bits = &b
	case "Nullable":
		o, ok := val.(OptValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to OptValue", val)
		}
		v.OptionalValue = &o
	case "CellOf":
		r, ok := val.(RefValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to RefValue", val)
		}
		v.RefValue = &r
	case "Tensor":
		t, ok := val.(TensorValues)
		if !ok {
			return fmt.Errorf("cannot convert %v to TensorValues", val)
		}
		v.Tensor = &t
	case "TupleWith":
		t, ok := val.(TupleValues)
		if !ok {
			return fmt.Errorf("cannot convert %v to TupleValues", val)
		}
		v.TupleWith = &t
	case "Map":
		m, ok := val.(MapValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to MapValue", val)
		}
		v.Map = &m
	case "EnumRef":
		e, ok := val.(EnumValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to EnumValue", val)
		}
		v.Enum = &e
	case "StructRef":
		s, ok := val.(Struct)
		if !ok {
			return fmt.Errorf("cannot convert %v to Struct", val)
		}
		v.Struct = &s
	case "AliasRef":
		a, ok := val.(AliasValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to AliasValue", val)
		}
		v.Alias = &a
	case "Generic":
		g, ok := val.(GenericValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to GenericValue", val)
		}
		v.Generic = &g
	case "Union":
		u, ok := val.(UnionValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to UnionValue", val)
		}
		v.Union = &u
	case "Int":
		return fmt.Errorf("int not supported")
	case "Coins":
		c, ok := val.(CoinsValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to CoinsValue", val)
		}
		v.Coins = &c
	case "Bool":
		b, ok := val.(BoolValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to BoolValue", val)
		}
		v.Bool = &b
	case "Cell":
		a, ok := val.(Any)
		if !ok {
			return fmt.Errorf("cannot convert %v to Any", val)
		}
		v.Cell = &a
	case "Slice":
		return fmt.Errorf("slice not supported")
	case "Builder":
		return fmt.Errorf("builder not supported")
	case "Callable":
		return fmt.Errorf("callable not supported")
	case "Remaining":
		r, ok := val.(RemainingValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to RemainingValue", val)
		}
		v.Remaining = &r
	case "Address":
		i, ok := val.(InternalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to InternalAddress", val)
		}
		v.InternalAddress = &i
	case "AddressOpt":
		o, ok := val.(OptionalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to OptionalAddress", val)
		}
		v.OptionalAddress = &o
	case "AddressExt":
		e, ok := val.(ExternalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to ExternalAddress", val)
		}
		v.ExternalAddress = &e
	case "AddressAny":
		a, ok := val.(AnyAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to AnyAddress", val)
		}
		v.AnyAddress = &a
	case "TupleAny":
		return fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		n, ok := val.(NullValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to NullValue", val)
		}
		v.Null = &n
	case "Void":
		vo, ok := val.(VoidValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to VoidValue", val)
		}
		v.Void = &vo
	default:
		return fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	return nil
}

func (v *Value) Unmarshal(cell *boc.Cell, ty tolkParser.Ty, decoder *Decoder) error {
	var err error
	switch ty.SumType {
	case "IntN":
		if ty.IntN.N <= 64 {
			v.SumType = "smallInt"
			def := Int64(0)
			v.SmallInt = &def
			err = v.SmallInt.Unmarshal(cell, *ty.IntN, decoder)
		} else {
			v.SumType = "bigInt"
			v.BigInt = &BigInt{}
			err = v.BigInt.Unmarshal(cell, *ty.IntN, decoder)
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			v.SumType = "smallUint"
			def := UInt64(0)
			v.SmallUint = &def
			err = v.SmallUint.Unmarshal(cell, *ty.UintN, decoder)
		} else {
			v.SumType = "bigUint"
			v.BigUint = &BigUInt{}
			err = v.BigUint.Unmarshal(cell, *ty.UintN, decoder)
		}
	case "VarIntN":
		v.SumType = "varInt"
		v.VarInt = &VarInt{}
		err = v.VarInt.Unmarshal(cell, *ty.VarIntN, decoder)
	case "VarUintN":
		v.SumType = "varUint"
		v.VarUint = &VarUInt{}
		err = v.VarUint.Unmarshal(cell, *ty.VarUintN, decoder)
	case "BitsN":
		v.SumType = "bits"
		v.Bits = &Bits{}
		err = v.Bits.Unmarshal(cell, *ty.BitsN, decoder)
	case "Nullable":
		v.SumType = "optionalValue"
		v.OptionalValue = &OptValue{}
		err = v.OptionalValue.Unmarshal(cell, *ty.Nullable, decoder)
	case "CellOf":
		v.SumType = "refValue"
		v.RefValue = &RefValue{}
		err = v.RefValue.Unmarshal(cell, *ty.CellOf, decoder)
	case "Tensor":
		v.SumType = "tensor"
		v.Tensor = &TensorValues{}
		err = v.Tensor.Unmarshal(cell, *ty.Tensor, decoder)
	case "TupleWith":
		v.SumType = "tupleWith"
		v.TupleWith = &TupleValues{}
		err = v.TupleWith.Unmarshal(cell, *ty.TupleWith, decoder)
	case "Map":
		v.SumType = "map"
		v.Map = &MapValue{}
		err = v.Map.Unmarshal(cell, *ty.Map, decoder)
	case "EnumRef":
		v.SumType = "enum"
		v.Enum = &EnumValue{}
		err = v.Enum.Unmarshal(cell, *ty.EnumRef, decoder)
	case "StructRef":
		v.SumType = "struct"
		v.Struct = &Struct{}
		err = v.Struct.Unmarshal(cell, *ty.StructRef, decoder)
	case "AliasRef":
		v.SumType = "alias"
		v.Alias = &AliasValue{}
		err = v.Alias.Unmarshal(cell, *ty.AliasRef, decoder)
	case "Generic":
		v.SumType = "generic"
		v.Generic = &GenericValue{}
		err = v.Generic.Unmarshal(cell, *ty.Generic, decoder)
	case "Union":
		v.SumType = "union"
		v.Union = &UnionValue{}
		err = v.Union.Unmarshal(cell, *ty.Union, decoder)
	case "Int":
		err = fmt.Errorf("int not supported")
	case "Coins":
		v.SumType = "coins"
		v.Coins = &CoinsValue{}
		err = v.Coins.Unmarshal(cell, *ty.Coins, decoder)
	case "Bool":
		v.SumType = "bool"
		def := BoolValue(false)
		v.Bool = &def
		err = v.Bool.Unmarshal(cell, *ty.Bool, decoder)
	case "Cell":
		v.SumType = "cell"
		v.Cell = &Any{}
		err = v.Cell.Unmarshal(cell, *ty.Cell, decoder)
	case "Slice":
		err = fmt.Errorf("slice not supported")
	case "Builder":
		err = fmt.Errorf("builder not supported")
	case "Callable":
		err = fmt.Errorf("callable not supported")
	case "Remaining":
		v.SumType = "remaining"
		v.Remaining = &RemainingValue{}
		err = v.Remaining.Unmarshal(cell, *ty.Remaining, decoder)
	case "Address":
		v.SumType = "internalAddress"
		v.InternalAddress = &InternalAddress{}
		err = v.InternalAddress.Unmarshal(cell, *ty.Address, decoder)
	case "AddressOpt":
		v.SumType = "optionalAddress"
		v.OptionalAddress = &OptionalAddress{}
		err = v.OptionalAddress.Unmarshal(cell, *ty.AddressOpt, decoder)
	case "AddressExt":
		v.SumType = "externalAddress"
		v.ExternalAddress = &ExternalAddress{}
		err = v.ExternalAddress.Unmarshal(cell, *ty.AddressExt, decoder)
	case "AddressAny":
		v.SumType = "anyAddress"
		v.AnyAddress = &AnyAddress{}
		err = v.AnyAddress.Unmarshal(cell, *ty.AddressAny, decoder)
	case "TupleAny":
		err = fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		v.SumType = "null"
		v.Null = &NullValue{}
		err = v.Null.Unmarshal(cell, *ty.NullLiteral, decoder)
	case "Void":
		v.SumType = "void"
		v.Void = &VoidValue{}
		err = v.Void.Unmarshal(cell, *ty.Void, decoder)
	default:
		return fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	if err != nil {
		return err
	}
	return nil
}

func (v *Value) Marshal(cell *boc.Cell, ty tolkParser.Ty, encoder *Encoder) error {
	var err error
	switch ty.SumType {
	case "IntN":
		if ty.IntN.N <= 64 {
			if v.SumType != "smallInt" {
				return fmt.Errorf("expected smallInt, but got %v", v.SumType)
			}
			return v.SmallInt.Marshal(cell, *ty.IntN, encoder)
		} else {
			if v.SumType != "bigInt" {
				return fmt.Errorf("expected bigInt, but got %v", v.SumType)
			}
			return v.BigInt.Marshal(cell, *ty.IntN, encoder)
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			if v.SumType != "smallUint" {
				return fmt.Errorf("expected smallUint, but got %v", v.SumType)
			}
			return v.SmallUint.Marshal(cell, *ty.UintN, encoder)
		} else {
			if v.SumType != "bigUint" {
				return fmt.Errorf("expected bigUint, but got %v", v.SumType)
			}
			return v.BigUint.Marshal(cell, *ty.UintN, encoder)
		}
	case "VarIntN":
		if v.SumType != "varInt" {
			return fmt.Errorf("expected varInt, but got %v", v.SumType)
		}
		return v.VarInt.Marshal(cell, *ty.VarIntN, encoder)
	case "VarUintN":
		if v.SumType != "varUint" {
			return fmt.Errorf("expected varUint, but got %v", v.SumType)
		}
		return v.VarUint.Marshal(cell, *ty.VarUintN, encoder)
	case "BitsN":
		if v.SumType != "bits" {
			return fmt.Errorf("expected bits, but got %v", v.SumType)
		}
		return v.Bits.Marshal(cell, *ty.BitsN, encoder)
	case "Nullable":
		if v.SumType != "optionalValue" {
			return fmt.Errorf("expected optionalValue, but got %v", v.SumType)
		}
		return v.OptionalValue.Marshal(cell, *ty.Nullable, encoder)
	case "CellOf":
		if v.SumType != "refValue" {
			return fmt.Errorf("expected refValue, but got %v", v.SumType)
		}
		return v.RefValue.Marshal(cell, *ty.CellOf, encoder)
	case "Tensor":
		if v.SumType != "tensor" {
			return fmt.Errorf("expected tensor, but got %v", v.SumType)
		}
		return v.Tensor.Marshal(cell, *ty.Tensor, encoder)
	case "TupleWith":
		if v.SumType != "tupleWith" {
			return fmt.Errorf("expected tupleWith, but got %v", v.SumType)
		}
		return v.TupleWith.Marshal(cell, *ty.TupleWith, encoder)
	case "Map":
		if v.SumType != "map" {
			return fmt.Errorf("expected map, but got %v", v.SumType)
		}
		return v.Map.Marshal(cell, *ty.Map, encoder)
	case "EnumRef":
		if v.SumType != "enum" {
			return fmt.Errorf("expected enum, but got %v", v.SumType)
		}
		return v.Enum.Marshal(cell, *ty.EnumRef, encoder)
	case "StructRef":
		if v.SumType != "struct" {
			return fmt.Errorf("expected struct, but got %v", v.SumType)
		}
		return v.Struct.Marshal(cell, *ty.StructRef, encoder)
	case "AliasRef":
		if v.SumType != "alias" {
			return fmt.Errorf("expected alias, but got %v", v.SumType)
		}
		return v.Alias.Marshal(cell, *ty.AliasRef, encoder)
	case "Generic":
		if v.SumType != "generic" {
			return fmt.Errorf("expected generic, but got %v", v.SumType)
		}
		return v.Generic.Marshal(cell, *ty.Generic, encoder)
	case "Union":
		if v.SumType != "union" {
			return fmt.Errorf("expected union, but got %v", v.SumType)
		}
		return v.Union.Marshal(cell, *ty.Union, encoder)
	case "Int":
		err = fmt.Errorf("int not supported")
	case "Coins":
		if v.SumType != "coins" {
			return fmt.Errorf("expected coins, but got %v", v.SumType)
		}
		return v.Coins.Marshal(cell, *ty.Coins, encoder)
	case "Bool":
		if v.SumType != "bool" {
			return fmt.Errorf("expected bool, but got %v", v.SumType)
		}
		return v.Bool.Marshal(cell, *ty.Bool, encoder)
	case "Cell":
		if v.SumType != "cell" {
			return fmt.Errorf("expected cell, but got %v", v.SumType)
		}
		return v.Cell.Marshal(cell, *ty.Cell, encoder)
	case "Slice":
		err = fmt.Errorf("slice not supported")
	case "Builder":
		err = fmt.Errorf("builder not supported")
	case "Callable":
		err = fmt.Errorf("callable not supported")
	case "Remaining":
		if v.SumType != "remaining" {
			return fmt.Errorf("expected remaining, but got %v", v.SumType)
		}
		return v.Remaining.Marshal(cell, *ty.Remaining, encoder)
	case "Address":
		if v.SumType != "internalAddress" {
			return fmt.Errorf("expected internalAddress, but got %v", v.SumType)
		}
		return v.InternalAddress.Marshal(cell, *ty.Address, encoder)
	case "AddressOpt":
		if v.SumType != "optionalAddress" {
			return fmt.Errorf("expected optionalAddress, but got %v", v.SumType)
		}
		return v.OptionalAddress.Marshal(cell, *ty.AddressOpt, encoder)
	case "AddressExt":
		if v.SumType != "externalAddress" {
			return fmt.Errorf("expected externalAddress, but got %v", v.SumType)
		}
		return v.ExternalAddress.Marshal(cell, *ty.AddressExt, encoder)
	case "AddressAny":
		if v.SumType != "anyAddress" {
			return fmt.Errorf("expected anyAddress, but got %v", v.SumType)
		}
		return v.AnyAddress.Marshal(cell, *ty.AddressAny, encoder)
	case "TupleAny":
		err = fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		if v.SumType != "null" {
			return fmt.Errorf("expected null, but got %v", v.SumType)
		}
		return v.Null.Marshal(cell, *ty.NullLiteral, encoder)
	case "Void":
		if v.SumType != "void" {
			return fmt.Errorf("expected void, but got %v", v.SumType)
		}
		return v.Void.Marshal(cell, *ty.Void, encoder)
	default:
		err = fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	if err != nil {
		return err
	}
	return nil
}

func (v *Value) unmarshalDefaultValue(d *tolkParser.DefaultValue, vType tolkParser.Ty) (bool, error) {
	switch d.SumType {
	case "IntDefaultValue":
		val, err := binDecHexToUint(d.IntDefaultValue.V)
		if err != nil {
			return false, err
		}
		err = v.SetValue(BigInt(*val), vType)
		if err != nil {
			return false, err
		}
	case "BoolDefaultValue":
		err := v.SetValue(BoolValue(d.BoolDefaultValue.V), vType)
		if err != nil {
			return false, err
		}
	case "SliceDefaultValue":
		val, err := hex.DecodeString(d.SliceDefaultValue.Hex)
		if err != nil {
			return false, err
		}
		bs := boc.NewBitString(hex.DecodedLen(len(val)))
		err = bs.WriteBytes(val)
		if err != nil {
			return false, err
		}
		err = v.SetValue(*boc.NewCellWithBits(bs), vType)
		if err != nil {
			return false, err
		}
	case "AddressDefaultValue":
		accountID, err := ton.ParseAccountID(d.AddressDefaultValue.Address)
		if err != nil {
			return false, err
		}
		err = v.SetValue(InternalAddress{
			Workchain: int8(accountID.Workchain),
			Address:   accountID.Address,
		}, vType)
		if err != nil {
			return false, err
		}
	case "TensorDefaultValue":
		if vType.SumType != "Tensor" {
			return false, fmt.Errorf("tensor default value type must be tensor, got %q", d.SumType)
		}
		tensor := make([]Value, len(d.TensorDefaultValue.Items))
		for i, item := range d.TensorDefaultValue.Items {
			val := Value{}
			_, err := v.unmarshalDefaultValue(&item, vType.Tensor.Items[i])
			if err != nil {
				return false, err
			}
			tensor[i] = val
		}
		err := v.SetValue(tensor, vType)
		if err != nil {
			return false, err
		}
	case "NullDefaultValue":
		return false, nil
	default:
		return false, fmt.Errorf("unknown default value type %q", d.SumType)
	}

	return true, nil
}

func (v *Value) Equal(o any) bool {
	otherValue, ok := o.(Value)
	if !ok {
		return false
	}

	switch v.SumType {
	case "bool":
		if otherValue.Bool == nil {
			return false
		}
		return v.Bool.Equal(*otherValue.Bool)
	case "smallInt":
		if otherValue.SmallInt == nil {
			return false
		}
		return v.SmallInt.Equal(*otherValue.SmallInt)
	case "smallUint":
		if otherValue.SmallUint == nil {
			return false
		}
		return v.SmallUint.Equal(*otherValue.SmallUint)
	case "bigInt":
		if otherValue.BigInt == nil {
			return false
		}
		return v.BigInt.Equal(*otherValue.BigInt)
	case "bigUint":
		if otherValue.BigUint == nil {
			return false
		}
		return v.BigUint.Equal(*otherValue.BigUint)
	case "varInt":
		if otherValue.VarInt == nil {
			return false
		}
		return v.VarInt.Equal(*otherValue.VarInt)
	case "varUint":
		if otherValue.VarUint == nil {
			return false
		}
		return v.VarUint.Equal(*otherValue.VarUint)
	case "coins":
		if otherValue.Coins == nil {
			return false
		}
		return v.Coins.Equal(*otherValue.Coins)
	case "bits":
		if otherValue.Bits == nil {
			return false
		}
		return v.Bits.Equal(*otherValue.Bits)
	case "cell":
		if otherValue.Cell == nil {
			return false
		}
		return v.Cell.Equal(*otherValue.Cell)
	case "remaining":
		if otherValue.Remaining == nil {
			return false
		}
		return v.Remaining.Equal(*otherValue.Remaining)
	case "internalAddress":
		if otherValue.InternalAddress == nil {
			return false
		}
		return v.InternalAddress.Equal(*otherValue.InternalAddress)
	case "optionalAddress":
		if otherValue.OptionalAddress == nil {
			return false
		}
		return v.OptionalAddress.Equal(*otherValue.OptionalAddress)
	case "externalAddress":
		if otherValue.ExternalAddress == nil {
			return false
		}
		return v.ExternalAddress.Equal(*otherValue.ExternalAddress)
	case "anyAddress":
		if otherValue.AnyAddress == nil {
			return false
		}
		return v.AnyAddress.Equal(*otherValue.AnyAddress)
	case "optionalValue":
		if otherValue.OptionalValue == nil {
			return false
		}
		return v.OptionalValue.Equal(*otherValue.OptionalValue)
	case "refValue":
		if otherValue.RefValue == nil {
			return false
		}
		return v.RefValue.Equal(*otherValue.RefValue)
	case "tupleWith":
		if otherValue.TupleWith == nil {
			return false
		}
		return v.TupleWith.Equal(*otherValue.TupleWith)
	case "tensor":
		if otherValue.Tensor == nil {
			return false
		}
		return v.Tensor.Equal(*otherValue.Tensor)
	case "map":
		if otherValue.Map == nil {
			return false
		}
		return v.Map.Equal(*otherValue.Map)
	case "struct":
		if otherValue.Struct == nil {
			return false
		}
		return v.Struct.Equal(*otherValue.Struct)
	case "alias":
		if otherValue.Alias == nil {
			return false
		}
		return v.Alias.Equal(*otherValue.Alias)
	case "enum":
		if otherValue.Enum == nil {
			return false
		}
		return v.Enum.Equal(*otherValue.Enum)
	case "generic":
		if otherValue.Generic == nil {
			return false
		}
		return v.Generic.Equal(*otherValue.Generic)
	case "union":
		if otherValue.Union == nil {
			return false
		}
		return v.Union.Equal(*otherValue.Union)
	case "null":
		if otherValue.Null == nil {
			return false
		}
		return v.Null.Equal(*otherValue.Null)
	case "void":
		if otherValue.Void == nil {
			return false
		}
		return v.Void.Equal(*otherValue.Void)
	default:
		return false
	}
}
