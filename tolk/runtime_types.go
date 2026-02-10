package tolk

// todo: move this to some package or rename somehow.

import (
	"bytes"
	"fmt"
	"maps"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
)

type SumType string

type TolkPrefix struct {
	Len    int16
	Prefix uint64
}

type Struct struct {
	hasPrefix bool
	prefix    TolkPrefix
	field     map[string]Value
}

func (s *Struct) GetField(field string) (Value, bool) {
	val, ok := s.field[field]
	return val, ok
}

// SetField set new value only if types are the same
func (s *Struct) SetField(field string, v Value) bool {
	s.field[field] = v
	return true
}

func (s *Struct) RemoveField(field string) {
	delete(s.field, field)
}

func (s *Struct) GetPrefix() (TolkPrefix, bool) {
	if !s.hasPrefix {
		return TolkPrefix{}, false
	}

	return s.prefix, true
}

func (s *Struct) Equal(o any) bool {
	otherStruct, ok := o.(Struct)
	if !ok {
		return false
	}
	if s.hasPrefix != otherStruct.hasPrefix {
		return false
	}
	if s.hasPrefix {
		if s.prefix != otherStruct.prefix {
			return false
		}
	}
	return maps.Equal(s.field, otherStruct.field)
}

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

func (v *Value) GetSmallInt() (int64, bool) {
	if v.smallInt == nil {
		return 0, false
	}
	return int64(*v.smallInt), true
}

func (v *Value) GetBigInt() (big.Int, bool) {
	if v.bigInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.bigInt), true
}

func (v *Value) GetSmallUInt() (uint64, bool) {
	if v.smallUint == nil {
		return 0, false
	}
	return uint64(*v.smallUint), true
}

func (v *Value) GetBigUInt() (big.Int, bool) {
	if v.bigUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.bigUint), true
}

func (v *Value) GetVarInt() (big.Int, bool) {
	if v.varInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.varInt), true
}

func (v *Value) GetVarUInt() (big.Int, bool) {
	if v.varUint == nil {
		return big.Int{}, false
	}
	return big.Int(*v.varUint), true
}

func (v *Value) GetCoins() (big.Int, bool) {
	if v.coins == nil {
		return big.Int{}, false
	}
	return big.Int(*v.coins), true
}

func (v *Value) GetBits() (boc.BitString, bool) {
	if v.bits == nil {
		return boc.BitString{}, false
	}
	return boc.BitString(*v.bits), true
}

func (v *Value) GetAddress() (InternalAddress, bool) {
	if v.internalAddress == nil {
		return InternalAddress{}, false
	}
	return *v.internalAddress, true
}

func (v *Value) GetOptionalAddress() (OptionalAddress, bool) {
	if v.optionalAddress == nil {
		return OptionalAddress{}, false
	}
	return *v.optionalAddress, true
}

func (v *Value) GetExternalAddress() (ExternalAddress, bool) {
	if v.externalAddress == nil {
		return ExternalAddress{}, false
	}
	return *v.externalAddress, true
}

func (v *Value) GetAnyAddress() (AnyAddress, bool) {
	if v.anyAddress == nil {
		return AnyAddress{}, false
	}
	return *v.anyAddress, true
}

func (v *Value) GetOptionalValue() (OptValue, bool) {
	if v.optionalValue == nil {
		return OptValue{}, false
	}
	return *v.optionalValue, true
}

func (v *Value) GetRefValue() (Value, bool) {
	if v.refValue == nil {
		return Value{}, false
	}
	return Value(*v.refValue), true
}

func (v *Value) GetTensor() ([]Value, bool) {
	if v.tensor == nil {
		return TensorValues{}, false
	}
	return *v.tensor, true
}

func (v *Value) GetMap() (MapValue, bool) {
	if v.mp == nil {
		return MapValue{}, false
	}
	return *v.mp, true
}

func (v *Value) GetStruct() (Struct, bool) {
	if v.structValue == nil {
		return Struct{}, false
	}
	return *v.structValue, true
}

func (v *Value) GetAlias() (Value, bool) {
	if v.alias == nil {
		return Value{}, false
	}
	return Value(*v.alias), true
}

func (v *Value) GetGeneric() (Value, bool) {
	if v.generic == nil {
		return Value{}, false
	}
	return Value(*v.generic), true
}

func (v *Value) GetEnum() (EnumValue, bool) {
	if v.enum == nil {
		return EnumValue{}, false
	}
	return *v.enum, true
}

func (v *Value) GetUnion() (UnionValue, bool) {
	if v.union == nil {
		return UnionValue{}, false
	}
	return *v.union, true
}

func (v *Value) GetTupleValues() ([]Value, bool) {
	if v.tupleWith == nil {
		return TupleValues{}, false
	}
	return *v.tupleWith, true
}

func (v *Value) GetCell() (boc.Cell, bool) {
	if v.cell == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.cell), true
}

func (v *Value) GetRemaining() (boc.Cell, bool) {
	if v.remaining == nil {
		return boc.Cell{}, false
	}
	return boc.Cell(*v.remaining), true
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

func (v *Value) UnmarshalTolk(cell *boc.Cell, ty Ty, decoder *Decoder) error {
	var err error
	switch ty.SumType {
	case "IntN":
		if ty.IntN.N <= 64 {
			v.sumType = "smallInt"
			def := Int64(0)
			v.smallInt = &def
			err = v.smallInt.UnmarshalTolk(cell, *ty.IntN, decoder)
		} else {
			v.sumType = "bigInt"
			v.bigInt = &BigInt{}
			err = v.bigInt.UnmarshalTolk(cell, *ty.IntN, decoder)
		}
	case "UintN":
		if ty.UintN.N <= 64 {
			v.sumType = "smallUint"
			def := UInt64(0)
			v.smallUint = &def
			err = v.smallUint.UnmarshalTolk(cell, *ty.UintN, decoder)
		} else {
			v.sumType = "bigUint"
			v.bigUint = &BigUInt{}
			err = v.bigUint.UnmarshalTolk(cell, *ty.UintN, decoder)
		}
	case "VarIntN":
		v.sumType = "varInt"
		v.varInt = &VarInt{}
		err = v.varInt.UnmarshalTolk(cell, *ty.VarIntN, decoder)
	case "VarUintN":
		v.sumType = "varUint"
		v.varUint = &VarUInt{}
		err = v.varUint.UnmarshalTolk(cell, *ty.VarUintN, decoder)
	case "BitsN":
		v.sumType = "bits"
		v.bits = &Bits{}
		err = v.bits.UnmarshalTolk(cell, *ty.BitsN, decoder)
	case "Nullable":
		v.sumType = "optionalValue"
		v.optionalValue = &OptValue{}
		err = v.optionalValue.UnmarshalTolk(cell, *ty.Nullable, decoder)
	case "CellOf":
		v.sumType = "refValue"
		v.refValue = &RefValue{}
		err = v.refValue.UnmarshalTolk(cell, *ty.CellOf, decoder)
	case "Tensor":
		v.sumType = "tensor"
		v.tensor = &TensorValues{}
		err = v.tensor.UnmarshalTolk(cell, *ty.Tensor, decoder)
	case "TupleWith":
		v.sumType = "tupleWith"
		v.tupleWith = &TupleValues{}
		err = v.tupleWith.UnmarshalTolk(cell, *ty.TupleWith, decoder)
	case "Map":
		v.sumType = "mp"
		v.mp = &MapValue{}
		err = v.mp.UnmarshalTolk(cell, *ty.Map, decoder)
	case "EnumRef":
		v.sumType = "enum"
		v.enum = &EnumValue{}
		err = v.enum.UnmarshalTolk(cell, *ty.EnumRef, decoder)
	case "StructRef":
		v.sumType = "structValue"
		v.structValue = &Struct{}
		err = v.structValue.UnmarshalTolk(cell, *ty.StructRef, decoder)
	case "AliasRef":
		v.sumType = "alias"
		v.alias = &AliasValue{}
		err = v.alias.UnmarshalTolk(cell, *ty.AliasRef, decoder)
	case "Generic":
		v.sumType = "generic"
		v.generic = &GenericValue{}
		err = v.generic.UnmarshalTolk(cell, *ty.Generic, decoder)
	case "Union":
		v.sumType = "union"
		v.union = &UnionValue{}
		err = v.union.UnmarshalTolk(cell, *ty.Union, decoder)
	case "Int":
		err = fmt.Errorf("int not supported")
	case "Coins":
		v.sumType = "coins"
		v.coins = &CoinsValue{}
		err = v.coins.UnmarshalTolk(cell, *ty.Coins, decoder)
	case "Bool":
		v.sumType = "bool"
		def := BoolValue(false)
		v.bool = &def
		err = v.bool.UnmarshalTolk(cell, *ty.Bool, decoder)
	case "Cell":
		v.sumType = "cell"
		v.cell = &Any{}
		err = v.cell.UnmarshalTolk(cell, *ty.Cell, decoder)
	case "Slice":
		err = fmt.Errorf("slice not supported")
	case "Builder":
		err = fmt.Errorf("builder not supported")
	case "Callable":
		err = fmt.Errorf("callable not supported")
	case "Remaining":
		v.sumType = "remaining"
		v.remaining = &RemainingValue{}
		err = v.remaining.UnmarshalTolk(cell, *ty.Remaining, decoder)
	case "Address":
		v.sumType = "internalAddress"
		v.internalAddress = &InternalAddress{}
		err = v.internalAddress.UnmarshalTolk(cell, *ty.Address, decoder)
	case "AddressOpt":
		v.sumType = "optionalAddress"
		v.optionalAddress = &OptionalAddress{}
		err = v.optionalAddress.UnmarshalTolk(cell, *ty.AddressOpt, decoder)
	case "AddressExt":
		v.sumType = "externalAddress"
		v.externalAddress = &ExternalAddress{}
		err = v.externalAddress.UnmarshalTolk(cell, *ty.AddressExt, decoder)
	case "AddressAny":
		v.sumType = "anyAddress"
		v.anyAddress = &AnyAddress{}
		err = v.anyAddress.UnmarshalTolk(cell, *ty.AddressAny, decoder)
	case "TupleAny":
		err = fmt.Errorf("tuple any not supported")
	case "NullLiteral":
		v.sumType = "null"
		v.null = &NullValue{}
		err = v.null.UnmarshalTolk(cell, *ty.NullLiteral, decoder)
	case "Void":
		v.sumType = "void"
		v.void = &VoidValue{}
		err = v.void.UnmarshalTolk(cell, *ty.Void, decoder)
	default:
		return fmt.Errorf("unknown ty type %q", ty.SumType)
	}
	if err != nil {
		return err
	}
	return nil
}

func (v *Value) SetValue(val any, ty Ty) error {
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

type BoolValue bool

func (b *BoolValue) Equal(o any) bool {
	otherBool, ok := o.(BoolValue)
	if !ok {
		return false
	}
	return *b == otherBool
}

type Any boc.Cell

func (a *Any) Equal(o any) bool {
	other, ok := o.(Any)
	if !ok {
		return false
	}
	cellV := boc.Cell(*a)
	vHash, err := cellV.HashString()
	if err != nil {
		return false
	}
	cellO := boc.Cell(other)
	oHash, err := cellO.HashString()
	if err != nil {
		return false
	}
	return oHash == vHash
}

type RemainingValue boc.Cell

func (r *RemainingValue) Equal(o any) bool {
	other, ok := o.(RemainingValue)
	if !ok {
		return false
	}
	cellV := boc.Cell(*r)
	vHash, err := cellV.HashString()
	if err != nil {
		return false
	}
	cellO := boc.Cell(other)
	oHash, err := cellO.HashString()
	if err != nil {
		return false
	}
	return oHash == vHash
}

type Int64 int64

func (i *Int64) Equal(other any) bool {
	otherInt, ok := other.(Int64)
	if !ok {
		return false
	}
	return *i == otherInt
}

type UInt64 uint64

func (i *UInt64) Equal(other any) bool {
	otherUint, ok := other.(UInt64)
	if !ok {
		return false
	}
	return *i == otherUint
}

type BigInt big.Int

func (b *BigInt) Equal(other any) bool {
	otherBigInt, ok := other.(BigInt)
	if !ok {
		return false
	}
	bi := big.Int(*b)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type BigUInt big.Int

func (b *BigUInt) Equal(other any) bool {
	otherBigInt, ok := other.(BigUInt)
	if !ok {
		return false
	}
	bi := big.Int(*b)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type VarInt big.Int

func (vi *VarInt) Equal(other any) bool {
	otherBigInt, ok := other.(VarInt)
	if !ok {
		return false
	}
	bi := big.Int(*vi)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type VarUInt big.Int

func (vu *VarUInt) Equal(other any) bool {
	otherBigInt, ok := other.(VarUInt)
	if !ok {
		return false
	}
	bi := big.Int(*vu)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type CoinsValue big.Int

func (c *CoinsValue) Equal(other any) bool {
	otherBigInt, ok := other.(CoinsValue)
	if !ok {
		return false
	}
	bi := big.Int(*c)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type Bits boc.BitString

func (b *Bits) Equal(other any) bool {
	otherBits, ok := other.(Bits)
	if !ok {
		return false
	}
	bs := boc.BitString(*b)
	otherBs := boc.BitString(otherBits)
	return bytes.Equal(bs.Buffer(), otherBs.Buffer())
}

type MapKey Value

type RefValue Value

func (r *RefValue) Equal(other any) bool {
	otherRefValue, ok := other.(RefValue)
	if !ok {
		return false
	}
	return r.Equal(otherRefValue)
}

type OptValue struct {
	IsExists bool
	Val      Value
}

func (o *OptValue) Equal(other any) bool {
	otherOptValue, ok := other.(OptValue)
	if !ok {
		return false
	}
	if o.IsExists != otherOptValue.IsExists {
		return false
	}
	if o.IsExists {
		return o.Val.Equal(otherOptValue.Val)
	}
	return true
}

type UnionValue struct {
	Prefix TolkPrefix
	Val    Value
}

func (u *UnionValue) Equal(other any) bool {
	otherUnionValue, ok := other.(UnionValue)
	if !ok {
		return false
	}
	if u.Prefix != otherUnionValue.Prefix {
		return false
	}
	return u.Val.Equal(otherUnionValue.Val)
}

type EnumValue struct {
	val   Value
	Name  string
	Value big.Int
}

func (e *EnumValue) Equal(other any) bool {
	otherEnumValue, ok := other.(EnumValue)
	if !ok {
		return false
	}
	if e.Name != otherEnumValue.Name {
		return false
	}
	if e.Value.Cmp(&otherEnumValue.Value) != 0 {
		return false
	}
	return e.val.Equal(otherEnumValue.val)
}

type AnyAddress struct {
	SumType
	InternalAddress *InternalAddress
	NoneAddress     *NoneAddress
	ExternalAddress *ExternalAddress
	VarAddress      *VarAddress
}

func (a *AnyAddress) Equal(other any) bool {
	otherAnyAddress, ok := other.(AnyAddress)
	if !ok {
		return false
	}
	if otherAnyAddress.SumType != a.SumType {
		return false
	}
	switch a.SumType {
	case "NoneAddress":
		return true
	case "InternalAddress":
		return a.InternalAddress.Equal(otherAnyAddress.InternalAddress)
	case "ExternalAddress":
		return a.ExternalAddress.Equal(otherAnyAddress.ExternalAddress)
	case "VarAddress":
		return a.VarAddress.Equal(otherAnyAddress.VarAddress)
	}
	return false
}

type NoneAddress struct {
}

type ExternalAddress struct {
	Len     int16
	Address boc.BitString
}

func (e *ExternalAddress) Equal(other any) bool {
	otherExternalAddress, ok := other.(ExternalAddress)
	if !ok {
		return false
	}
	if e.Len != otherExternalAddress.Len {
		return false
	}
	return bytes.Equal(e.Address.Buffer(), otherExternalAddress.Address.Buffer())
}

type InternalAddress struct {
	Workchain int8
	Address   [32]byte
}

func (i *InternalAddress) Equal(other any) bool {
	otherInternalAddress, ok := other.(InternalAddress)
	if !ok {
		return false
	}
	return *i == otherInternalAddress
}

func (i *InternalAddress) ToRaw() string {
	return fmt.Sprintf("%v:%x", i.Workchain, i.Address)
}

type VarAddress struct {
	Len       int16
	Workchain int32
	Address   boc.BitString
}

func (v *VarAddress) Equal(other any) bool {
	otherVarAddress, ok := other.(VarAddress)
	if !ok {
		return false
	}
	if v.Len != otherVarAddress.Len {
		return false
	}
	if v.Workchain != otherVarAddress.Workchain {
		return false
	}
	return bytes.Equal(v.Address.Buffer(), otherVarAddress.Address.Buffer())
}

type OptionalAddress struct {
	SumType
	NoneAddress     NoneAddress
	InternalAddress InternalAddress
}

func (o *OptionalAddress) Equal(other any) bool {
	otherOptionalAddress, ok := other.(OptionalAddress)
	if !ok {
		return false
	}
	if o.SumType != otherOptionalAddress.SumType {
		return false
	}
	if o.SumType == "InternalAddress" {
		return o.InternalAddress.Equal(otherOptionalAddress.InternalAddress)
	}
	return true
}

type TupleValues []Value

func (v *TupleValues) Equal(other any) bool {
	otherTupleValues, ok := other.(TupleValues)
	if !ok {
		return false
	}
	wV := *v
	if len(otherTupleValues) != len(wV) {
		return false
	}
	for i := range wV {
		if !wV[i].Equal(otherTupleValues[i]) {
			return false
		}
	}
	return true
}

type TensorValues []Value

func (v *TensorValues) Equal(other any) bool {
	otherTensorValues, ok := other.(TensorValues)
	if !ok {
		return false
	}
	wV := *v
	if len(otherTensorValues) != len(wV) {
		return false
	}
	for i := range wV {
		if !wV[i].Equal(otherTensorValues[i]) {
			return false
		}
	}
	return true
}

type MapValue struct {
	keys   []Value
	values []Value
	len    int
}

func (m *MapValue) Equal(other any) bool {
	otherMapValue, ok := other.(MapValue)
	if !ok {
		return false
	}
	if m.len != otherMapValue.len {
		return false
	}
	for i := range m.keys {
		if !m.keys[i].Equal(otherMapValue.keys[i]) {
			return false
		}
		if !m.values[i].Equal(otherMapValue.values[i]) {
			return false
		}
	}
	return true
}

func (m *MapValue) Get(key MapKey) (Value, bool) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			return m.values[i], true
		}
	}

	return Value{}, false
}

func (m *MapValue) GetBySmallInt(v Int64) (Value, bool) {
	key := MapKey{
		sumType:  "smallInt",
		smallInt: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetBySmallUInt(v UInt64) (Value, bool) {
	key := MapKey{
		sumType:   "smallUint",
		smallUint: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBigInt(v BigInt) (Value, bool) {
	key := MapKey{
		sumType: "bigInt",
		bigInt:  &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBigUInt(v BigUInt) (Value, bool) {
	key := MapKey{
		sumType: "bigUint",
		bigUint: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBits(v Bits) (Value, bool) {
	key := MapKey{
		sumType: "bits",
		bits:    &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByInternalAddress(v InternalAddress) (Value, bool) {
	key := MapKey{
		sumType:         "internalAddress",
		internalAddress: &v,
	}
	return m.Get(key)
}

func (m *MapValue) Set(key MapKey, value Value) (bool, error) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			m.values[i] = value
			return true, nil
		}
	}

	m.keys = append(m.keys, Value(key))
	m.values = append(m.values, value)
	m.len++
	return true, nil
}

func (m *MapValue) SetBySmallInt(k Int64, value Value) (bool, error) {
	key := MapKey{
		sumType:  "smallInt",
		smallInt: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetBySmallUInt(k UInt64, value Value) (bool, error) {
	key := MapKey{
		sumType:   "smallUint",
		smallUint: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBigInt(k BigInt, value Value) (bool, error) {
	key := MapKey{
		sumType: "bigInt",
		bigInt:  &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBigUInt(k BigUInt, value Value) (bool, error) {
	key := MapKey{
		sumType: "bigUint",
		bigUint: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBits(k Bits, value Value) (bool, error) {
	key := MapKey{
		sumType: "bits",
		bits:    &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByInternalAddress(k InternalAddress, value Value) (bool, error) {
	key := MapKey{
		sumType:         "internalAddress",
		internalAddress: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) Delete(key MapKey) {
	for i, k := range m.keys {
		if k.Equal(Value(key)) {
			m.keys[i] = Value{}
			m.values[i] = Value{}
			m.len--
		}
	}
}

func (m *MapValue) DeleteBySmallInt(k Int64) {
	key := MapKey{
		sumType:  "smallInt",
		smallInt: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteBySmallUInt(k UInt64) {
	key := MapKey{
		sumType:   "smallUint",
		smallUint: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBigInt(k BigInt) {
	key := MapKey{
		sumType: "bigInt",
		bigInt:  &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBigUInt(k BigUInt) {
	key := MapKey{
		sumType: "bigUint",
		bigUint: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBits(k Bits) {
	key := MapKey{
		sumType: "bits",
		bits:    &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByInternalAddress(k InternalAddress) {
	key := MapKey{
		sumType:         "internalAddress",
		internalAddress: &k,
	}
	m.Delete(key)
}

func (m *MapValue) Len() int {
	return m.len
}

type AliasValue Value

func (a *AliasValue) Equal(other any) bool {
	otherAlias, ok := other.(AliasValue)
	if !ok {
		return false
	}
	v := Value(*a)
	return v.Equal(Value(otherAlias))
}

type GenericValue Value

func (g *GenericValue) Equal(other any) bool {
	otherGeneric, ok := other.(GenericValue)
	if !ok {
		return false
	}
	v := Value(*g)
	return v.Equal(Value(otherGeneric))
}

type NullValue struct{}

func (n *NullValue) Equal(other any) bool {
	_, ok := other.(NullValue)
	if !ok {
		return false
	}
	return true
}

type VoidValue struct{}

func (v *VoidValue) Equal(other any) bool {
	_, ok := other.(VoidValue)
	if !ok {
		return false
	}
	return true
}
