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
	sumType         SumType
	bool            *BoolValue
	smallInt        *Int64
	smallUint       *UInt64
	bigInt          *BigInt
	bigUint         *BigUInt
	varInt          *VarInt
	varUint         *VarUInt
	coins           *CoinsValue
	bits            *Bits
	cell            *Any
	remaining       *RemainingValue
	internalAddress *InternalAddress
	optionalAddress *OptionalAddress
	externalAddress *ExternalAddress
	anyAddress      *AnyAddress
	optionalValue   *OptValue
	refValue        *RefValue
	tupleWith       *TupleValues
	tensor          *TensorValues
	mp              *MapValue
	structValue     *Struct
	alias           *AliasValue
	enum            *EnumValue
	generic         *GenericValue
	union           *UnionValue
	null            *NullValue
	void            *VoidValue
}

func (v *Value) GetBool() (bool, bool) {
	if v.bool == nil {
		return false, false
	}
	return bool(*v.bool), true
}

func (v *Value) MustGetBool() bool {
	if v.bool == nil {
		panic("value is not a bool")
	}
	return bool(*v.bool)
}

func (v *Value) GetSmallInt() (int64, bool) {
	if v.smallInt == nil {
		return 0, false
	}
	return int64(*v.smallInt), true
}

func (v *Value) MustGetSmallInt() int64 {
	if v.smallInt == nil {
		panic("value is not a small int")
	}
	return int64(*v.smallInt)
}

func (v *Value) GetBigInt() (big.Int, bool) {
	if v.bigInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.bigInt), true
}

func (v *Value) MustGetBigInt() big.Int {
	if v.bigInt == nil {
		panic("value is not a big int")
	}
	return big.Int(*v.bigInt)
}

func (v *Value) GetSmallUInt() (uint64, bool) {
	if v.smallUint == nil {
		return 0, false
	}
	return uint64(*v.smallUint), true
}

func (v *Value) MustGetSmallUInt() uint64 {
	if v.smallUint == nil {
		panic("value is not a small uint")
	}
	return uint64(*v.smallUint)
}

func (v *Value) GetBigUInt() (big.Int, bool) {
	if v.bigUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.bigUint), true
}

func (v *Value) MustGetBigUInt() big.Int {
	if v.bigUint == nil {
		panic("value is not a big uint")
	}
	return big.Int(*v.bigUint)
}

func (v *Value) GetVarInt() (big.Int, bool) {
	if v.varInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.varInt), true
}

func (v *Value) MustGetVarInt() big.Int {
	if v.varInt == nil {
		panic("value is not a var int")
	}
	return big.Int(*v.varInt)
}

func (v *Value) GetVarUInt() (big.Int, bool) {
	if v.varUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.varUint), true
}

func (v *Value) MustGetVarUInt() big.Int {
	if v.varUint == nil {
		panic("value is not a var uint")
	}
	return big.Int(*v.varUint)
}

func (v *Value) GetCoins() (big.Int, bool) {
	if v.coins == nil {
		return big.Int{}, false
	}
	return big.Int(*v.coins), true
}

func (v *Value) MustGetCoins() big.Int {
	if v.coins == nil {
		panic("value is not a coins")
	}
	return big.Int(*v.coins)
}

func (v *Value) GetBits() (boc.BitString, bool) {
	if v.bits == nil {
		return boc.BitString{}, false
	}
	return boc.BitString(*v.bits), true
}

func (v *Value) MustGetBits() boc.BitString {
	if v.bits == nil {
		panic("value is not a bits")
	}
	return boc.BitString(*v.bits)
}

func (v *Value) GetAddress() (InternalAddress, bool) {
	if v.internalAddress == nil {
		return InternalAddress{}, false
	}
	return *v.internalAddress, true
}

func (v *Value) MustGetAddress() InternalAddress {
	if v.internalAddress == nil {
		panic("value is not an internal address")
	}
	return *v.internalAddress
}

func (v *Value) GetOptionalAddress() (OptionalAddress, bool) {
	if v.optionalAddress == nil {
		return OptionalAddress{}, false
	}
	return *v.optionalAddress, true
}

func (v *Value) MustGetOptionalAddress() OptionalAddress {
	if v.optionalAddress == nil {
		panic("value is not an optional address")
	}
	return *v.optionalAddress
}

func (v *Value) GetExternalAddress() (ExternalAddress, bool) {
	if v.externalAddress == nil {
		return ExternalAddress{}, false
	}
	return *v.externalAddress, true
}

func (v *Value) MustGetExternalAddress() ExternalAddress {
	if v.externalAddress == nil {
		panic("value is not an external address")
	}
	return *v.externalAddress
}

func (v *Value) GetAnyAddress() (AnyAddress, bool) {
	if v.anyAddress == nil {
		return AnyAddress{}, false
	}
	return *v.anyAddress, true
}

func (v *Value) MustGetAnyAddress() AnyAddress {
	if v.anyAddress == nil {
		panic("value is not an any address")
	}
	return *v.anyAddress
}

func (v *Value) GetOptionalValue() (OptValue, bool) {
	if v.optionalValue == nil {
		return OptValue{}, false
	}
	return *v.optionalValue, true
}

func (v *Value) MustGetOptionalValue() OptValue {
	if v.optionalValue == nil {
		panic("value is not an optional value")
	}
	return *v.optionalValue
}

func (v *Value) GetRefValue() (Value, bool) {
	if v.refValue == nil {
		return Value{}, false
	}
	return Value(*v.refValue), true
}

func (v *Value) MustGetRefValue() Value {
	if v.refValue == nil {
		panic("value is not a reference")
	}
	return Value(*v.refValue)
}

func (v *Value) GetTensor() ([]Value, bool) {
	if v.tensor == nil {
		return TensorValues{}, false
	}
	return *v.tensor, true
}

func (v *Value) MustGetTensor() []Value {
	if v.tensor == nil {
		panic("value is not a tensor")
	}
	return *v.tensor
}

func (v *Value) GetMap() (MapValue, bool) {
	if v.mp == nil {
		return MapValue{}, false
	}
	return *v.mp, true
}

func (v *Value) MustGetMap() MapValue {
	if v.mp == nil {
		panic("value is not a map")
	}
	return *v.mp
}

func (v *Value) GetStruct() (Struct, bool) {
	if v.structValue == nil {
		return Struct{}, false
	}
	return *v.structValue, true
}

func (v *Value) MustGetStruct() Struct {
	if v.structValue == nil {
		panic("value is not a struct")
	}
	return *v.structValue
}

func (v *Value) GetAlias() (Value, bool) {
	if v.alias == nil {
		return Value{}, false
	}
	return Value(*v.alias), true
}

func (v *Value) MustGetAlias() Value {
	if v.alias == nil {
		panic("value is not an alias")
	}
	return Value(*v.alias)
}

func (v *Value) GetGeneric() (Value, bool) {
	if v.generic == nil {
		return Value{}, false
	}
	return Value(*v.generic), true
}

func (v *Value) MustGetGeneric() Value {
	if v.generic == nil {
		panic("value is not a generic")
	}
	return Value(*v.generic)
}

func (v *Value) GetEnum() (EnumValue, bool) {
	if v.enum == nil {
		return EnumValue{}, false
	}
	return *v.enum, true
}

func (v *Value) MustGetEnum() EnumValue {
	if v.enum == nil {
		panic("value is not an enum")
	}
	return *v.enum
}

func (v *Value) GetUnion() (UnionValue, bool) {
	if v.union == nil {
		return UnionValue{}, false
	}
	return *v.union, true
}

func (v *Value) MustGetUnion() UnionValue {
	if v.union == nil {
		panic("value is not an union")
	}
	return *v.union
}

func (v *Value) GetTupleValues() ([]Value, bool) {
	if v.tupleWith == nil {
		return TupleValues{}, false
	}
	return *v.tupleWith, true
}

func (v *Value) MustGetTupleValues() []Value {
	if v.tupleWith == nil {
		panic("value is not a tuple")
	}
	return *v.tupleWith
}

func (v *Value) GetCell() (boc.Cell, bool) {
	if v.cell == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.cell), true
}

func (v *Value) MustGetCell() boc.Cell {
	if v.cell == nil {
		panic("value is not a cell")
	}
	return boc.Cell(*v.cell)
}

func (v *Value) GetRemaining() (boc.Cell, bool) {
	if v.remaining == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.remaining), true
}

func (v *Value) MustGetRemaining() boc.Cell {
	if v.remaining == nil {
		panic("value is not a remaining")
	}
	return boc.Cell(*v.remaining)
}

func (v *Value) GetType() string {
	return string(v.sumType)
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
			v.smallInt = &wVal
		} else {
			v.bigInt = &bi
		}
	case "UintN":
		bi, ok := val.(BigUInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to BigUInt", val)
		}
		if ty.UintN.N <= 64 {
			b := big.Int(bi)
			wVal := UInt64(b.Uint64())
			v.smallUint = &wVal
		} else {
			v.bigUint = &bi
		}
	case "VarIntN":
		vi, ok := val.(VarInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to VarInt", val)
		}
		v.varInt = &vi
	case "VarUintN":
		vi, ok := val.(VarUInt)
		if !ok {
			return fmt.Errorf("cannot convert %v to VarUInt", val)
		}
		v.varUint = &vi
	case "BitsN":
		b, ok := val.(Bits)
		if !ok {
			return fmt.Errorf("cannot convert %v to Bits", val)
		}
		v.bits = &b
	case "Nullable":
		o, ok := val.(OptValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to OptValue", val)
		}
		v.optionalValue = &o
	case "CellOf":
		r, ok := val.(RefValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to RefValue", val)
		}
		v.refValue = &r
	case "Tensor":
		t, ok := val.(TensorValues)
		if !ok {
			return fmt.Errorf("cannot convert %v to TensorValues", val)
		}
		v.tensor = &t
	case "TupleWith":
		t, ok := val.(TupleValues)
		if !ok {
			return fmt.Errorf("cannot convert %v to TupleValues", val)
		}
		v.tupleWith = &t
	case "Map":
		m, ok := val.(MapValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to MapValue", val)
		}
		v.mp = &m
	case "EnumRef":
		e, ok := val.(EnumValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to EnumValue", val)
		}
		v.enum = &e
	case "StructRef":
		s, ok := val.(Struct)
		if !ok {
			return fmt.Errorf("cannot convert %v to Struct", val)
		}
		v.structValue = &s
	case "AliasRef":
		a, ok := val.(AliasValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to AliasValue", val)
		}
		v.alias = &a
	case "Generic":
		g, ok := val.(GenericValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to GenericValue", val)
		}
		v.generic = &g
	case "Union":
		u, ok := val.(UnionValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to UnionValue", val)
		}
		v.union = &u
	case "Int":
		return fmt.Errorf("int not supported")
	case "Coins":
		c, ok := val.(CoinsValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to CoinsValue", val)
		}
		v.coins = &c
	case "Bool":
		b, ok := val.(BoolValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to BoolValue", val)
		}
		v.bool = &b
	case "Cell":
		a, ok := val.(Any)
		if !ok {
			return fmt.Errorf("cannot convert %v to Any", val)
		}
		v.cell = &a
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
		v.remaining = &r
	case "Address":
		i, ok := val.(InternalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to InternalAddress", val)
		}
		v.internalAddress = &i
	case "AddressOpt":
		o, ok := val.(OptionalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to OptionalAddress", val)
		}
		v.optionalAddress = &o
	case "AddressExt":
		e, ok := val.(ExternalAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to ExternalAddress", val)
		}
		v.externalAddress = &e
	case "AddressAny":
		a, ok := val.(AnyAddress)
		if !ok {
			return fmt.Errorf("cannot convert %v to AnyAddress", val)
		}
		v.anyAddress = &a
	case "TupleAny":
		return fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		n, ok := val.(NullValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to NullValue", val)
		}
		v.null = &n
	case "Void":
		vo, ok := val.(VoidValue)
		if !ok {
			return fmt.Errorf("cannot convert %v to VoidValue", val)
		}
		v.void = &vo
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
			v.sumType = "smallInt"
			def := Int64(0)
			v.smallInt = &def
			err = v.smallInt.Unmarshal(cell, *ty.IntN, decoder)
		} else {
			v.sumType = "bigInt"
			v.bigInt = &BigInt{}
			err = v.bigInt.Unmarshal(cell, *ty.IntN, decoder)
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			v.sumType = "smallUint"
			def := UInt64(0)
			v.smallUint = &def
			err = v.smallUint.Unmarshal(cell, *ty.UintN, decoder)
		} else {
			v.sumType = "bigUint"
			v.bigUint = &BigUInt{}
			err = v.bigUint.Unmarshal(cell, *ty.UintN, decoder)
		}
	case "VarIntN":
		v.sumType = "varInt"
		v.varInt = &VarInt{}
		err = v.varInt.Unmarshal(cell, *ty.VarIntN, decoder)
	case "VarUintN":
		v.sumType = "varUint"
		v.varUint = &VarUInt{}
		err = v.varUint.Unmarshal(cell, *ty.VarUintN, decoder)
	case "BitsN":
		v.sumType = "bits"
		v.bits = &Bits{}
		err = v.bits.Unmarshal(cell, *ty.BitsN, decoder)
	case "Nullable":
		v.sumType = "optionalValue"
		v.optionalValue = &OptValue{}
		err = v.optionalValue.Unmarshal(cell, *ty.Nullable, decoder)
	case "CellOf":
		v.sumType = "refValue"
		v.refValue = &RefValue{}
		err = v.refValue.Unmarshal(cell, *ty.CellOf, decoder)
	case "Tensor":
		v.sumType = "tensor"
		v.tensor = &TensorValues{}
		err = v.tensor.Unmarshal(cell, *ty.Tensor, decoder)
	case "TupleWith":
		v.sumType = "tupleWith"
		v.tupleWith = &TupleValues{}
		err = v.tupleWith.Unmarshal(cell, *ty.TupleWith, decoder)
	case "Map":
		v.sumType = "mp"
		v.mp = &MapValue{}
		err = v.mp.Unmarshal(cell, *ty.Map, decoder)
	case "EnumRef":
		v.sumType = "enum"
		v.enum = &EnumValue{}
		err = v.enum.Unmarshal(cell, *ty.EnumRef, decoder)
	case "StructRef":
		v.sumType = "structValue"
		v.structValue = &Struct{}
		err = v.structValue.Unmarshal(cell, *ty.StructRef, decoder)
	case "AliasRef":
		v.sumType = "alias"
		v.alias = &AliasValue{}
		err = v.alias.Unmarshal(cell, *ty.AliasRef, decoder)
	case "Generic":
		v.sumType = "generic"
		v.generic = &GenericValue{}
		err = v.generic.Unmarshal(cell, *ty.Generic, decoder)
	case "Union":
		v.sumType = "union"
		v.union = &UnionValue{}
		err = v.union.Unmarshal(cell, *ty.Union, decoder)
	case "Int":
		err = fmt.Errorf("int not supported")
	case "Coins":
		v.sumType = "coins"
		v.coins = &CoinsValue{}
		err = v.coins.Unmarshal(cell, *ty.Coins, decoder)
	case "Bool":
		v.sumType = "bool"
		def := BoolValue(false)
		v.bool = &def
		err = v.bool.Unmarshal(cell, *ty.Bool, decoder)
	case "Cell":
		v.sumType = "cell"
		v.cell = &Any{}
		err = v.cell.Unmarshal(cell, *ty.Cell, decoder)
	case "Slice":
		err = fmt.Errorf("slice not supported")
	case "Builder":
		err = fmt.Errorf("builder not supported")
	case "Callable":
		err = fmt.Errorf("callable not supported")
	case "Remaining":
		v.sumType = "remaining"
		v.remaining = &RemainingValue{}
		err = v.remaining.Unmarshal(cell, *ty.Remaining, decoder)
	case "Address":
		v.sumType = "internalAddress"
		v.internalAddress = &InternalAddress{}
		err = v.internalAddress.Unmarshal(cell, *ty.Address, decoder)
	case "AddressOpt":
		v.sumType = "optionalAddress"
		v.optionalAddress = &OptionalAddress{}
		err = v.optionalAddress.Unmarshal(cell, *ty.AddressOpt, decoder)
	case "AddressExt":
		v.sumType = "externalAddress"
		v.externalAddress = &ExternalAddress{}
		err = v.externalAddress.Unmarshal(cell, *ty.AddressExt, decoder)
	case "AddressAny":
		v.sumType = "anyAddress"
		v.anyAddress = &AnyAddress{}
		err = v.anyAddress.Unmarshal(cell, *ty.AddressAny, decoder)
	case "TupleAny":
		err = fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		v.sumType = "null"
		v.null = &NullValue{}
		err = v.null.Unmarshal(cell, *ty.NullLiteral, decoder)
	case "Void":
		v.sumType = "void"
		v.void = &VoidValue{}
		err = v.void.Unmarshal(cell, *ty.Void, decoder)
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
			if v.sumType != "smallInt" {
				return fmt.Errorf("expected smallInt, but got %v", v.sumType)
			}
			return v.smallInt.Marshal(cell, *ty.IntN, encoder)
		} else {
			if v.sumType != "bigInt" {
				return fmt.Errorf("expected bigInt, but got %v", v.sumType)
			}
			return v.bigInt.Marshal(cell, *ty.IntN, encoder)
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			if v.sumType != "smallUint" {
				return fmt.Errorf("expected smallUint, but got %v", v.sumType)
			}
			return v.smallUint.Marshal(cell, *ty.UintN, encoder)
		} else {
			if v.sumType != "bigUint" {
				return fmt.Errorf("expected bigUint, but got %v", v.sumType)
			}
			return v.bigUint.Marshal(cell, *ty.UintN, encoder)
		}
	case "VarIntN":
		if v.sumType != "varInt" {
			return fmt.Errorf("expected varInt, but got %v", v.sumType)
		}
		return v.varInt.Marshal(cell, *ty.VarIntN, encoder)
	case "VarUintN":
		if v.sumType != "varUint" {
			return fmt.Errorf("expected varUint, but got %v", v.sumType)
		}
		return v.varUint.Marshal(cell, *ty.VarUintN, encoder)
	case "BitsN":
		if v.sumType != "bits" {
			return fmt.Errorf("expected bits, but got %v", v.sumType)
		}
		return v.bits.Marshal(cell, *ty.BitsN, encoder)
	case "Nullable":
		if v.sumType != "optionalValue" {
			return fmt.Errorf("expected optionalValue, but got %v", v.sumType)
		}
		return v.optionalValue.Marshal(cell, *ty.Nullable, encoder)
	case "CellOf":
		if v.sumType != "refValue" {
			return fmt.Errorf("expected refValue, but got %v", v.sumType)
		}
		return v.refValue.Marshal(cell, *ty.CellOf, encoder)
	case "Tensor":
		if v.sumType != "tensor" {
			return fmt.Errorf("expected tensor, but got %v", v.sumType)
		}
		return v.tensor.Marshal(cell, *ty.Tensor, encoder)
	case "TupleWith":
		if v.sumType != "tupleWith" {
			return fmt.Errorf("expected tupleWith, but got %v", v.sumType)
		}
		return v.tupleWith.Marshal(cell, *ty.TupleWith, encoder)
	case "Map":
		if v.sumType != "mp" {
			return fmt.Errorf("expected mp, but got %v", v.sumType)
		}
		return v.mp.Marshal(cell, *ty.Map, encoder)
	case "EnumRef":
		if v.sumType != "enum" {
			return fmt.Errorf("expected enum, but got %v", v.sumType)
		}
		return v.enum.Marshal(cell, *ty.EnumRef, encoder)
	case "StructRef":
		if v.sumType != "structValue" {
			return fmt.Errorf("expected structValue, but got %v", v.sumType)
		}
		return v.structValue.Marshal(cell, *ty.StructRef, encoder)
	case "AliasRef":
		if v.sumType != "alias" {
			return fmt.Errorf("expected alias, but got %v", v.sumType)
		}
		return v.alias.Marshal(cell, *ty.AliasRef, encoder)
	case "Generic":
		if v.sumType != "generic" {
			return fmt.Errorf("expected generic, but got %v", v.sumType)
		}
		return v.generic.Marshal(cell, *ty.Generic, encoder)
	case "Union":
		if v.sumType != "union" {
			return fmt.Errorf("expected union, but got %v", v.sumType)
		}
		return v.union.Marshal(cell, *ty.Union, encoder)
	case "Int":
		err = fmt.Errorf("int not supported")
	case "Coins":
		if v.sumType != "coins" {
			return fmt.Errorf("expected coins, but got %v", v.sumType)
		}
		return v.coins.Marshal(cell, *ty.Coins, encoder)
	case "Bool":
		if v.sumType != "bool" {
			return fmt.Errorf("expected bool, but got %v", v.sumType)
		}
		return v.bool.Marshal(cell, *ty.Bool, encoder)
	case "Cell":
		if v.sumType != "cell" {
			return fmt.Errorf("expected cell, but got %v", v.sumType)
		}
		return v.cell.Marshal(cell, *ty.Cell, encoder)
	case "Slice":
		err = fmt.Errorf("slice not supported")
	case "Builder":
		err = fmt.Errorf("builder not supported")
	case "Callable":
		err = fmt.Errorf("callable not supported")
	case "Remaining":
		if v.sumType != "remaining" {
			return fmt.Errorf("expected remaining, but got %v", v.sumType)
		}
		return v.remaining.Marshal(cell, *ty.Remaining, encoder)
	case "Address":
		if v.sumType != "internalAddress" {
			return fmt.Errorf("expected internalAddress, but got %v", v.sumType)
		}
		return v.internalAddress.Marshal(cell, *ty.Address, encoder)
	case "AddressOpt":
		if v.sumType != "optionalAddress" {
			return fmt.Errorf("expected optionalAddress, but got %v", v.sumType)
		}
		return v.optionalAddress.Marshal(cell, *ty.AddressOpt, encoder)
	case "AddressExt":
		if v.sumType != "externalAddress" {
			return fmt.Errorf("expected externalAddress, but got %v", v.sumType)
		}
		return v.externalAddress.Marshal(cell, *ty.AddressExt, encoder)
	case "AddressAny":
		if v.sumType != "anyAddress" {
			return fmt.Errorf("expected anyAddress, but got %v", v.sumType)
		}
		return v.anyAddress.Marshal(cell, *ty.AddressAny, encoder)
	case "TupleAny":
		err = fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		if v.sumType != "null" {
			return fmt.Errorf("expected null, but got %v", v.sumType)
		}
		return v.null.Marshal(cell, *ty.NullLiteral, encoder)
	case "Void":
		if v.sumType != "void" {
			return fmt.Errorf("expected void, but got %v", v.sumType)
		}
		return v.void.Marshal(cell, *ty.Void, encoder)
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

	switch v.sumType {
	case "bool":
		if otherValue.bool == nil {
			return false
		}
		return v.bool.Equal(*otherValue.bool)
	case "smallInt":
		if otherValue.smallInt == nil {
			return false
		}
		return v.smallInt.Equal(*otherValue.smallInt)
	case "smallUint":
		if otherValue.smallUint == nil {
			return false
		}
		return v.smallUint.Equal(*otherValue.smallUint)
	case "bigInt":
		if otherValue.bigInt == nil {
			return false
		}
		return v.bigInt.Equal(*otherValue.bigInt)
	case "bigUint":
		if otherValue.bigUint == nil {
			return false
		}
		return v.bigUint.Equal(*otherValue.bigUint)
	case "varInt":
		if otherValue.varInt == nil {
			return false
		}
		return v.varInt.Equal(*otherValue.varInt)
	case "varUint":
		if otherValue.varUint == nil {
			return false
		}
		return v.varUint.Equal(*otherValue.varUint)
	case "coins":
		if otherValue.coins == nil {
			return false
		}
		return v.coins.Equal(*otherValue.coins)
	case "bits":
		if otherValue.bits == nil {
			return false
		}
		return v.bits.Equal(*otherValue.bits)
	case "cell":
		if otherValue.cell == nil {
			return false
		}
		return v.cell.Equal(*otherValue.cell)
	case "remaining":
		if otherValue.remaining == nil {
			return false
		}
		return v.remaining.Equal(*otherValue.remaining)
	case "internalAddress":
		if otherValue.internalAddress == nil {
			return false
		}
		return v.internalAddress.Equal(*otherValue.internalAddress)
	case "optionalAddress":
		if otherValue.optionalAddress == nil {
			return false
		}
		return v.optionalAddress.Equal(*otherValue.optionalAddress)
	case "externalAddress":
		if otherValue.externalAddress == nil {
			return false
		}
		return v.externalAddress.Equal(*otherValue.externalAddress)
	case "anyAddress":
		if otherValue.anyAddress == nil {
			return false
		}
		return v.anyAddress.Equal(*otherValue.anyAddress)
	case "optionalValue":
		if otherValue.optionalValue == nil {
			return false
		}
		return v.optionalValue.Equal(*otherValue.optionalValue)
	case "refValue":
		if otherValue.refValue == nil {
			return false
		}
		return v.refValue.Equal(*otherValue.refValue)
	case "tupleWith":
		if otherValue.tupleWith == nil {
			return false
		}
		return v.tupleWith.Equal(*otherValue.tupleWith)
	case "tensor":
		if otherValue.tensor == nil {
			return false
		}
		return v.tensor.Equal(*otherValue.tensor)
	case "mp":
		if otherValue.mp == nil {
			return false
		}
		return v.mp.Equal(*otherValue.mp)
	case "structValue":
		if otherValue.structValue == nil {
			return false
		}
		return v.structValue.Equal(*otherValue.structValue)
	case "alias":
		if otherValue.alias == nil {
			return false
		}
		return v.alias.Equal(*otherValue.alias)
	case "enum":
		if otherValue.enum == nil {
			return false
		}
		return v.enum.Equal(*otherValue.enum)
	case "generic":
		if otherValue.generic == nil {
			return false
		}
		return v.generic.Equal(*otherValue.generic)
	case "union":
		if otherValue.union == nil {
			return false
		}
		return v.union.Equal(*otherValue.union)
	case "null":
		if otherValue.null == nil {
			return false
		}
		return v.null.Equal(*otherValue.null)
	case "void":
		if otherValue.void == nil {
			return false
		}
		return v.void.Equal(*otherValue.void)
	default:
		return false
	}
}
