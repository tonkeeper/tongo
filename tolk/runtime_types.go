package tolk

// todo: move this to some package or rename somehow.

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
)

type TolkValue interface {
	GetType() Ty
	SetType(Ty)
	SetValue(any) error
}

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

func (s *Struct) AddField(field string, v Value) bool {
	if _, found := s.field[field]; found {
		return false
	}

	s.field[field] = v
	return true
}

// SetField set new value only if types are the same
func (s *Struct) SetField(field string, v Value) bool {
	old, found := s.field[field]
	if !found {
		return false
	}
	if old.valType != v.valType {
		return false
	}
	s.field[field] = v
	return true
}

// UpdateField set new value even if types mismatched
func (s *Struct) UpdateField(field string, v Value) bool {
	_, found := s.field[field]
	if !found {
		return false
	}
	s.field[field] = v
	return true
}

func (s *Struct) RemoveField(field string) {
	delete(s.field, field)
}

func (s *Struct) GetStructPrefix() (TolkPrefix, bool) {
	if !s.hasPrefix {
		return TolkPrefix{}, false
	}

	return s.prefix, true
}

type Value struct {
	valType         Ty
	bool            *BoolValue
	smallInt        *Int64
	smallUint       *UInt64
	bigInt          *BigInt
	bits            *Bits
	cell            *Any
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
	generic         *GenericValue
	enum            *EnumValue
	union           *UnionValue
}

func (v *Value) GetType() Ty {
	return v.valType
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

func (v *Value) GetSmallUInt() (uint64, bool) {
	if v.smallUint == nil {
		return 0, false
	}
	return uint64(*v.smallUint), true
}

func (v *Value) GetBigInt() (big.Int, bool) {
	if v.bigInt == nil {
		return big.Int{}, false
	}
	return big.Int(*v.bigInt), true
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

func (v *Value) GetGeneric() (GenericValue, bool) {
	if v.generic == nil {
		return GenericValue{}, false
	}
	return *v.generic, true
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

// todo: is good idea to use this method
func (v *Value) GetAny() (any, error) {
	switch v.valType.SumType {
	case "Int":
		return v, fmt.Errorf("int is not supported")
	case "NullLiteral":
		return nil, nil
	case "Void":
		return nil, nil
	case "IntN":
		if i64, ok := v.GetSmallInt(); ok {
			return Int64(i64), nil
		}
		if bi, ok := v.GetBigInt(); ok {
			return BigInt(bi), nil
		}
		return nil, fmt.Errorf("value is not a BigInt or Int64")
	case "UintN":
		if ui64, ok := v.GetSmallUInt(); ok {
			return UInt64(ui64), nil
		}
		if bi, ok := v.GetBigInt(); ok {
			return BigInt(bi), nil
		}
		return nil, fmt.Errorf("value is not a BigInt or UInt64")
	case "VarIntN":
		if bi, ok := v.GetBigInt(); ok {
			return BigInt(bi), nil
		}
		return nil, fmt.Errorf("value is not a BigInt")
	case "VarUintN":
		if bi, ok := v.GetBigInt(); ok {
			return BigInt(bi), nil
		}
		return nil, fmt.Errorf("value is not a BigInt")
	case "BitsN":
		if bits, ok := v.GetBits(); ok {
			return Bits(bits), nil
		}
		return nil, fmt.Errorf("value is not a Bits")
	case "Coins":
		if bi, ok := v.GetBigInt(); ok {
			return BigInt(bi), nil
		}
		return nil, fmt.Errorf("value is not a BigInt")
	case "Bool":
		if b, ok := v.GetBool(); ok {
			return BoolValue(b), nil
		}
		return nil, fmt.Errorf("value is not a Bool")
	case "Cell":
		if b, ok := v.GetCell(); ok {
			return Any(b), nil
		}
		return nil, fmt.Errorf("value is not an Any")
	case "Slice":
		if b, ok := v.GetCell(); ok {
			return Any(b), nil
		}
		return nil, fmt.Errorf("value is not an Any")
	case "Builder":
		if b, ok := v.GetCell(); ok {
			return Any(b), nil
		}
		return nil, fmt.Errorf("value is not an Any")
	case "Callable":
		if b, ok := v.GetCell(); ok {
			return Any(b), nil
		}
		return nil, fmt.Errorf("value is not an Any")
	case "Remaining":
		if b, ok := v.GetCell(); ok {
			return Any(b), nil
		}
		return nil, fmt.Errorf("value is not an Any")
	case "Address":
		if a, ok := v.GetAddress(); ok {
			return a, nil
		}
		return nil, fmt.Errorf("value is not an InternalAddress")
	case "AddressOpt":
		if a, ok := v.GetOptionalAddress(); ok {
			return a, nil
		}
		return nil, fmt.Errorf("value is not an OptionalAddress")
	case "AddressExt":
		if a, ok := v.GetExternalAddress(); ok {
			return a, nil
		}
		return nil, fmt.Errorf("value is not an ExternalAddress")
	case "AddressAny":
		if a, ok := v.GetAnyAddress(); ok {
			return a, nil
		}
		return nil, fmt.Errorf("value is not an AnyAddress")
	case "Nullable":
		if o, ok := v.GetOptionalValue(); ok {
			return o, nil
		}
		return nil, fmt.Errorf("value is not an OptionalValue")
	case "CellOf":
		if r, ok := v.GetRefValue(); ok {
			return RefValue(r), nil
		}
		return nil, fmt.Errorf("value is not a RefValue")
	case "Tensor":
		if t, ok := v.GetTensor(); ok {
			return TensorValues(t), nil
		}
		return nil, fmt.Errorf("value is not a Tensor")
	case "TupleWith":
		if t, ok := v.GetTupleValues(); ok {
			return TupleValues(t), nil
		}
		return nil, fmt.Errorf("value is not a Tuple")
	case "TupleAny":
		return nil, fmt.Errorf("tuple any is not supported")
	case "Map":
		if m, ok := v.GetMap(); ok {
			return m, nil
		}
		return nil, fmt.Errorf("value is not a Map")
	case "EnumRef":
		if e, ok := v.GetEnum(); ok {
			return e, nil
		}
		return nil, fmt.Errorf("value is not an Enum")
	case "StructRef":
		if s, ok := v.GetStruct(); ok {
			return s, nil
		}
		return nil, fmt.Errorf("value is not a Struct")
	case "Union":
		if u, ok := v.GetUnion(); ok {
			return u, nil
		}
		return nil, fmt.Errorf("value is not an Union")
	default:
		return nil, fmt.Errorf("unknown value type %v", v.valType.SumType)
	}
}

func (v *Value) SetType(t Ty) {
	v.valType = t
}

func (v *Value) SetValue(val any) error {
	switch v.valType.SumType {
	case "Int":
		return v.valType.Int.SetValue(v, val)
	case "NullLiteral":
		return v.valType.NullLiteral.SetValue(v, val)
	case "Void":
		return v.valType.Void.SetValue(v, val)
	case "IntN":
		return v.valType.IntN.SetValue(v, val)
	case "UintN":
		return v.valType.UintN.SetValue(v, val)
	case "VarIntN":
		return v.valType.VarIntN.SetValue(v, val)
	case "VarUintN":
		return v.valType.VarUintN.SetValue(v, val)
	case "BitsN":
		return v.valType.BitsN.SetValue(v, val)
	case "Coins":
		return v.valType.Coins.SetValue(v, val)
	case "Bool":
		return v.valType.Bool.SetValue(v, val)
	case "Cell":
		return v.valType.Cell.SetValue(v, val)
	case "Slice":
		return v.valType.Slice.SetValue(v, val)
	case "Builder":
		return v.valType.Builder.SetValue(v, val)
	case "Callable":
		return v.valType.Callable.SetValue(v, val)
	case "Remaining":
		return v.valType.Remaining.SetValue(v, val)
	case "Address":
		return v.valType.Address.SetValue(v, val)
	case "AddressOpt":
		return v.valType.AddressOpt.SetValue(v, val)
	case "AddressExt":
		return v.valType.AddressExt.SetValue(v, val)
	case "AddressAny":
		return v.valType.AddressAny.SetValue(v, val)
	case "Nullable":
		return v.valType.Nullable.SetValue(v, val)
	case "CellOf":
		return v.valType.CellOf.SetValue(v, val)
	case "Tensor":
		return v.valType.Tensor.SetValue(v, val)
	case "TupleWith":
		return v.valType.TupleWith.SetValue(v, val)
	case "TupleAny":
		return v.valType.TupleAny.SetValue(v, val)
	case "Map":
		return v.valType.Map.SetValue(v, val)
	case "EnumRef":
		return v.valType.EnumRef.SetValue(v, val)
	case "StructRef":
		return v.valType.StructRef.SetValue(v, val)
	case "AliasRef":
		return v.valType.AliasRef.SetValue(v, val)
	case "Generic":
		return v.valType.Generic.SetValue(v, val)
	case "Union":
		return v.valType.Union.SetValue(v, val)
	default:
		return fmt.Errorf("unknown t type %q", v.valType.SumType)
	}
}

func (v *Value) Equal(o Value) bool {
	if v.valType != o.valType {
		return false
	}
	switch v.valType.SumType {
	case "Int":
		return v.valType.Int.Equal(*v, o)
	case "NullLiteral":
		return v.valType.NullLiteral.Equal(*v, o)
	case "Void":
		return v.valType.Void.Equal(*v, o)
	case "IntN":
		return v.valType.IntN.Equal(*v, o)
	case "UintN":
		return v.valType.UintN.Equal(*v, o)
	case "VarIntN":
		return v.valType.VarIntN.Equal(*v, o)
	case "VarUintN":
		return v.valType.VarUintN.Equal(*v, o)
	case "BitsN":
		return v.valType.BitsN.Equal(*v, o)
	case "Coins":
		return v.valType.Coins.Equal(*v, o)
	case "Bool":
		return v.valType.Bool.Equal(*v, o)
	case "Cell":
		return v.valType.Cell.Equal(*v, o)
	case "Slice":
		return v.valType.Slice.Equal(*v, o)
	case "Builder":
		return v.valType.Builder.Equal(*v, o)
	case "Callable":
		return v.valType.Callable.Equal(*v, o)
	case "Remaining":
		return v.valType.Remaining.Equal(*v, o)
	case "Address":
		return v.valType.Address.Equal(*v, o)
	case "AddressOpt":
		return v.valType.AddressOpt.Equal(*v, o)
	case "AddressExt":
		return v.valType.AddressExt.Equal(*v, o)
	case "AddressAny":
		return v.valType.AddressAny.Equal(*v, o)
	case "Nullable":
		return v.valType.Nullable.Equal(*v, o)
	case "CellOf":
		return v.valType.CellOf.Equal(*v, o)
	case "Tensor":
		return v.valType.Tensor.Equal(*v, o)
	case "TupleWith":
		return v.valType.TupleWith.Equal(*v, o)
	case "TupleAny":
		return v.valType.TupleAny.Equal(*v, o)
	case "Map":
		return v.valType.Map.Equal(*v, o)
	case "EnumRef":
		return v.valType.EnumRef.Equal(*v, o)
	case "StructRef":
		return v.valType.StructRef.Equal(*v, o)
	case "AliasRef":
		return v.valType.AliasRef.Equal(*v, o)
	case "Generic":
		return v.valType.Generic.Equal(*v, o)
	case "Union":
		return v.valType.Union.Equal(*v, o)
	default:
		return false
	}
}

type BoolValue bool

type Any boc.Cell

type Int64 int64

func (i Int64) Equal(other any) bool {
	otherInt, ok := other.(Int64)
	if !ok {
		return false
	}
	return i == otherInt
}

type UInt64 uint64

func (i UInt64) Equal(other any) bool {
	otherUint, ok := other.(UInt64)
	if !ok {
		return false
	}
	return i == otherUint
}

type BigInt big.Int

func (b BigInt) Equal(other any) bool {
	otherBigInt, ok := other.(big.Int)
	if !ok {
		return false
	}
	bi := big.Int(b)
	return bi.Cmp(&otherBigInt) == 0
}

type Bits boc.BitString

func (b Bits) Equal(other any) bool {
	otherBits, ok := other.(Bits)
	if !ok {
		return false
	}
	bs := boc.BitString(b)
	otherBs := boc.BitString(otherBits)
	return bytes.Equal(bs.Buffer(), otherBs.Buffer())
}

type MapKey Value

type RefValue Value

type OptValue struct {
	IsExists bool
	Val      Value
}

type UnionValue struct {
	Prefix TolkPrefix
	Val    Value
}

type EnumValue struct {
	enumType Ty
	Name     string
	Value    big.Int
}

type AnyAddress struct {
	SumType
	InternalAddress *InternalAddress
	NoneAddress     *NoneAddress
	ExternalAddress *ExternalAddress
	VarAddress      *VarAddress
}

type NoneAddress struct {
}

type ExternalAddress struct {
	Len     int16
	Address boc.BitString
}

type InternalAddress struct {
	Workchain int8
	Address   [32]byte
}

func (i InternalAddress) Equal(other any) bool {
	otherInternalAddress, ok := other.(InternalAddress)
	if !ok {
		return false
	}
	return i == otherInternalAddress
}

func (i InternalAddress) ToRaw() string {
	return fmt.Sprintf("%v:%x", i.Workchain, i.Address)
}

type VarAddress struct {
	Len       int16
	Workchain int32
	Address   boc.BitString
}

type OptionalAddress struct {
	SumType
	NoneAddress     NoneAddress
	InternalAddress InternalAddress
}

type TupleValues []Value

type TensorValues []Value

type MapValue struct {
	keyType Ty
	valType Ty
	keys    []Value
	values  []Value
	len     int
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
		valType:  m.keyType,
		smallInt: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetBySmallUInt(v UInt64) (Value, bool) {
	key := MapKey{
		valType:   m.keyType,
		smallUint: &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBigInt(v BigInt) (Value, bool) {
	key := MapKey{
		valType: m.keyType,
		bigInt:  &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByBits(v Bits) (Value, bool) {
	key := MapKey{
		valType: m.keyType,
		bits:    &v,
	}
	return m.Get(key)
}

func (m *MapValue) GetByInternalAddress(v InternalAddress) (Value, bool) {
	key := MapKey{
		valType:         m.keyType,
		internalAddress: &v,
	}
	return m.Get(key)
}

func (m *MapValue) Set(key MapKey, value Value) (bool, error) {
	if key.valType != m.keyType {
		return false, fmt.Errorf("map key has type %v, got %v", m.keyType, key.valType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
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
	if m.keyType.SumType != "IntN" {
		return false, fmt.Errorf("map key has type %v, got int64", m.keyType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
	key := MapKey{
		valType:  m.valType,
		smallInt: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetBySmallUInt(k UInt64, value Value) (bool, error) {
	if m.keyType.SumType != "UintN" {
		return false, fmt.Errorf("map key has type %v, got int64", m.keyType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
	key := MapKey{
		valType:   m.valType,
		smallUint: &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBigInt(k BigInt, value Value) (bool, error) {
	if m.keyType.SumType != "IntN" && m.keyType.SumType != "UintN" && m.keyType.SumType != "VarIntN" && m.keyType.SumType != "VarUintN" {
		return false, fmt.Errorf("map key has type %v, got BigInt", m.keyType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
	key := MapKey{
		valType: m.valType,
		bigInt:  &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByBits(k Bits, value Value) (bool, error) {
	if m.keyType.SumType != "BitsN" {
		return false, fmt.Errorf("map key has type %v, got BitsN", m.keyType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
	key := MapKey{
		valType: m.valType,
		bits:    &k,
	}
	return m.Set(key, value)
}

func (m *MapValue) SetByInternalAddress(k InternalAddress, value Value) (bool, error) {
	if m.keyType.SumType != "Address" {
		return false, fmt.Errorf("map key has type %v, got Address", m.keyType)
	}
	if value.valType != m.valType {
		return false, fmt.Errorf("map value has type %v, got %v", m.valType, value.valType)
	}
	key := MapKey{
		valType:         m.valType,
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
		valType:  m.valType,
		smallInt: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteBySmallUInt(k UInt64) {
	key := MapKey{
		valType:   m.valType,
		smallUint: &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBigInt(k BigInt) {
	key := MapKey{
		valType: m.valType,
		bigInt:  &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByBits(k Bits) {
	key := MapKey{
		valType: m.valType,
		bits:    &k,
	}
	m.Delete(key)
}

func (m *MapValue) DeleteByInternalAddress(k InternalAddress) {
	key := MapKey{
		valType:         m.valType,
		internalAddress: &k,
	}
	m.Delete(key)
}

func (m *MapValue) Len() int {
	return m.len
}

type AliasValue Value

type GenericValue Value
