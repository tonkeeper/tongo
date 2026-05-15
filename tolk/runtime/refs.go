package runtime

import (
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Prefix struct {
	Len    int16  `json:"len"`
	Prefix uint64 `json:"prefix"`
}

type Struct struct {
	hasPrefix  bool
	name       string
	prefix     Prefix
	fieldNames []string
	fields     map[string]Value
}

func (s *Struct) UnmarshalTyIdx(cell *boc.Cell, tyIdx int, decoder *Decoder) error {
	if decoder.abiIndex.Structs == nil {
		return fmt.Errorf("struct has struct reference but no abi has been given")
	}
	ty, err := decoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	if ty.SumType != parser.TyKindStructRef {
		return fmt.Errorf("expected StructRef at ty_idx=%d", tyIdx)
	}
	structDef, found := decoder.abiIndex.Structs[ty.StructRef.StructName]
	if !found {
		return fmt.Errorf("struct with name %s was not found in given abi", ty.StructRef.StructName)
	}
	structVal := Struct{
		fieldNames: make([]string, 0),
		fields:     make(map[string]Value),
		name:       ty.StructRef.StructName,
	}
	if structDef.Prefix != nil {
		prefixLen := structDef.Prefix.PrefixLen
		if prefixLen > 64 {
			return fmt.Errorf("struct %v prefix length must be lower than 64", structDef.Name)
		}

		prefix, err := cell.ReadUint(prefixLen)
		if err != nil {
			return fmt.Errorf("failed to read struct's %v-bit length prefix: %w", prefixLen, err)
		}
		actualPrefix := uint64(structDef.Prefix.PrefixNum)

		if prefix != actualPrefix {
			return fmt.Errorf("struct %v prefix does not match actual prefix %v", structDef.Name, actualPrefix)
		}
		structVal.hasPrefix = true
		structVal.prefix = Prefix{
			Len:    int16(prefixLen),
			Prefix: prefix,
		}
	}

	fields, err := decoder.abiIndex.StructFieldsOf(tyIdx, false)
	if err != nil {
		return fmt.Errorf("failed to resolve struct fields: %w", err)
	}
	for _, field := range fields {
		fieldVal := Value{}

		err = fieldVal.UnmarshalTyIdx(cell, field.TyIdx, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct's field %s: %w", field.Name, err)
		}

		structVal.fieldNames = append(structVal.fieldNames, field.Name)
		structVal.fields[field.Name] = fieldVal
	}

	*s = structVal

	return nil
}

// try to resolve payload from all known abi in trace
func (s *Struct) resolvePayload(cell *boc.Cell, tyIdx int, decoder *Decoder) (Value, bool, error) {
	ty, err := decoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return Value{}, false, err
	}
	switch ty.SumType {
	case parser.TyKindRemaining:
		isRef, err := cell.ReadBit()
		if err != nil {
			return Value{}, false, fmt.Errorf("failed to read isRef prefix: %w", err)
		}
		payload := cell.CopyRemaining()
		if isRef {
			payload, err = payload.NextRef()
			if err != nil {
				return Value{}, false, fmt.Errorf("failed to read payload ref: %w", err)
			}
		}
		v, isResolved, err := decoder.resolvePayload(payload)
		if err != nil {
			return Value{}, false, fmt.Errorf("failed to resolve payload: %w", err)
		}
		if isResolved {
			return v, true, nil
		}
	case parser.TyKindCell:
		payload, err := cell.NextRef()
		if err != nil {
			return Value{}, false, fmt.Errorf("failed to read payload ref: %w", err)
		}
		payload = payload.CopyRemaining()
		v, isResolved, err := decoder.resolvePayload(payload)
		if err != nil {
			return Value{}, false, fmt.Errorf("failed to resolve payload: %w", err)
		}
		if isResolved {
			return v, true, nil
		}
	case parser.TyKindAliasRef:
		if decoder.abiIndex.Aliases == nil {
			return Value{}, false, fmt.Errorf("struct has alias reference but no abi has been given")
		}
		targetTyIdx, _, err := decoder.abiIndex.AliasTargetOf(tyIdx)
		if err != nil {
			return Value{}, false, fmt.Errorf("failed to resolve alias target: %w", err)
		}
		return s.resolvePayload(cell, targetTyIdx, decoder)
	case parser.TyKindGenericT:
		currentTyIdx, found := decoder.genericRefs[ty.GenericT.NameT]
		if !found {
			return Value{}, false, fmt.Errorf("cannot resolve generic's type %v ", ty.GenericT.NameT)
		}

		return s.resolvePayload(cell, currentTyIdx, decoder)
	}

	return Value{}, false, nil
}

func (s *Struct) MarshalTyIdx(cell *boc.Cell, tyIdx int, encoder *Encoder) error {
	if encoder.abiIndex.Structs == nil {
		return fmt.Errorf("struct has struct reference but no abi has been given")
	}
	ty, err := encoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	if ty.SumType != parser.TyKindStructRef {
		return fmt.Errorf("expected StructRef at ty_idx=%d", tyIdx)
	}
	strct, found := encoder.abiIndex.Structs[ty.StructRef.StructName]
	if !found {
		return fmt.Errorf("struct with name %s was not found in given abi", ty.StructRef.StructName)
	}

	if strct.Prefix != nil {
		actualPrefix := uint64(strct.Prefix.PrefixNum)

		if s.prefix.Prefix != actualPrefix {
			return fmt.Errorf("struct %v prefix does not match actual prefix %v", strct.Name, actualPrefix)
		}

		err = cell.WriteUint(s.prefix.Prefix, int(s.prefix.Len))
		if err != nil {
			return fmt.Errorf("failed to write struct's %v-bit prefix %v: %w", s.prefix.Len, s.prefix.Prefix, err)
		}
	}

	fields, err := encoder.abiIndex.StructFieldsOf(tyIdx, false)
	if err != nil {
		return fmt.Errorf("failed to resolve struct fields: %w", err)
	}
	for _, field := range fields {
		val, ok := s.GetField(field.Name)
		if !ok {
			return fmt.Errorf("struct %v has no field %v", strct.Name, field.Name)
		}

		err = val.MarshalTyIdx(cell, field.TyIdx, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal struct's field %s: %w", field.Name, err)
		}
	}

	return nil
}

func (s *Struct) GetField(field string) (Value, bool) {
	if s.fields == nil {
		return Value{}, false
	}
	v, ok := s.fields[field]
	return v, ok
}

func (s *Struct) MustGetField(field string) Value {
	if s.fields != nil {
		if v, ok := s.fields[field]; ok {
			return v
		}
	}
	panic("field with name " + field + " is not found")
}

func (s *Struct) SetField(field string, v Value) bool {
	for _, name := range s.fieldNames {
		if name == field {
			if s.fields == nil {
				s.fields = make(map[string]Value)
			}
			s.fields[field] = v
			return true
		}
	}
	return false
}

func (s *Struct) RemoveField(field string) {
	for i, name := range s.fieldNames {
		if name == field {
			s.fieldNames = append(s.fieldNames[:i], s.fieldNames[i+1:]...)
			delete(s.fields, field)
			return
		}
	}
}

func (s *Struct) GetPrefix() (Prefix, bool) {
	if !s.hasPrefix {
		return Prefix{}, false
	}

	return s.prefix, true
}

func (s *Struct) GetName() string {
	return s.name
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
	if !slices.Equal(s.fieldNames, otherStruct.fieldNames) {
		return false
	}
	if len(s.fields) != len(otherStruct.fields) {
		return false
	}
	for name, value := range s.fields {
		otherValue, ok := otherStruct.fields[name]
		if !ok || !value.Equal(otherValue) {
			return false
		}
	}
	return true
}

func (s Struct) MarshalJSON() ([]byte, error) {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for i, name := range s.fieldNames {
		if i != 0 {
			builder.WriteRune(',')
		}
		builder.WriteString(fmt.Sprintf("\"%s\":", name))
		val, ok := s.fields[name]
		if !ok {
			return nil, fmt.Errorf("struct field %s not found", name)
		}
		data, err := json.Marshal(&val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal struct's field %s: %w", name, err)
		}
		builder.Write(data)
	}
	builder.WriteRune('}')
	return []byte(builder.String()), nil
}

type EnumValue struct {
	ActualValue Value
	Name        string
	Value       big.Int
}

func (e *EnumValue) Unmarshal(cell *boc.Cell, ty parser.EnumRef, decoder *Decoder) error {
	if decoder.abiIndex.Enums == nil {
		return fmt.Errorf("struct has enum reference but no abi has been given")
	}
	enum, found := decoder.abiIndex.Enums[ty.EnumName]
	if !found {
		return fmt.Errorf("enum with name %s was not found in given abi", ty.EnumName)
	}

	enumVal := Value{}
	err := enumVal.UnmarshalTyIdx(cell, enum.EncodedAsTyIdx, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal enum's value: %w", err)
	}
	encodedAs, err := decoder.abiIndex.TyByIdx(enum.EncodedAsTyIdx)
	if err != nil {
		return fmt.Errorf("failed to resolve enum encoding type: %w", err)
	}
	var bigEnumVal big.Int
	switch encodedAs.SumType {
	case parser.TyKindIntN:
		if encodedAs.IntN.N > 64 {
			bigEnumVal = big.Int(*enumVal.BigInt)
		} else {
			bigEnumVal = *big.NewInt(int64(*enumVal.SmallInt))
		}
	case parser.TyKindUintN:
		if encodedAs.UintN.N > 64 {
			bigEnumVal = big.Int(*enumVal.BigUint)
		} else {
			bigEnumVal = *new(big.Int).SetUint64(uint64(*enumVal.SmallUint))
		}
	case parser.TyKindVarIntN:
		bigEnumVal = big.Int(*enumVal.VarInt)
	case parser.TyKindVarUintN:
		bigEnumVal = big.Int(*enumVal.VarUint)
	default:
		return fmt.Errorf("enum encode type must be integer, got: %s", encodedAs.SumType)
	}

	for _, member := range enum.Members {
		val := member.Value.Int
		if val.Cmp(&bigEnumVal) == 0 {
			*e = EnumValue{
				ActualValue: enumVal,
				Name:        member.Name,
				Value:       val,
			}

			return nil
		}
	}

	return fmt.Errorf("enum value didn't match any values")
}

func (e EnumValue) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(e.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enum's value: %w", err)
	}
	return data, nil
}

func (e *EnumValue) Marshal(cell *boc.Cell, ty parser.EnumRef, encoder *Encoder) error {
	if encoder.abiIndex.Enums == nil {
		return fmt.Errorf("struct has enum reference but no abi has been given")
	}
	enum, found := encoder.abiIndex.Enums[ty.EnumName]
	if !found {
		return fmt.Errorf("enum with name %s was not found in given abi", ty.EnumName)
	}

	err := e.ActualValue.MarshalTyIdx(cell, enum.EncodedAsTyIdx, encoder)
	if err != nil {
		return fmt.Errorf("failed to marshal enum's value: %w", err)
	}

	for _, member := range enum.Members {
		if member.Value.Cmp(&e.Value) == 0 {
			return nil
		}
	}
	return fmt.Errorf("enum value not matched, got: %s", e.Value.String())
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
	return e.ActualValue.Equal(otherEnumValue.ActualValue)
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

func (a *AliasValue) UnmarshalTyIdx(cell *boc.Cell, tyIdx int, decoder *Decoder) error {
	if decoder.abiIndex.Aliases == nil {
		return fmt.Errorf("struct has alias reference but no abi has been given")
	}
	ty, err := decoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	if ty.SumType != parser.TyKindAliasRef {
		return fmt.Errorf("expected AliasRef at ty_idx=%d", tyIdx)
	}

	alias, err := decoder.abiIndex.GetAlias(ty.AliasRef.AliasName)
	if err != nil {
		return err
	}
	if alias.CustomPackUnpack.UnpackFromSlice {
		if decoder.customUnpackResolver == nil {
			return fmt.Errorf("custom unpack resolver for alias %q is not configured", ty.AliasRef.AliasName)
		}
		return decoder.customUnpackResolver(*ty.AliasRef, cell, a)
	}

	targetTyIdx, _, err := decoder.abiIndex.AliasTargetOf(tyIdx)
	if err != nil {
		return fmt.Errorf("failed to resolve alias target: %w", err)
	}

	val := Value{}
	err = val.UnmarshalTyIdx(cell, targetTyIdx, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal alias value: %w", err)
	}

	*a = AliasValue(val)

	return nil
}

func (a *AliasValue) MarshalTyIdx(cell *boc.Cell, tyIdx int, encoder *Encoder) error {
	if encoder.abiIndex.Aliases == nil {
		return fmt.Errorf("struct has alias reference but no abi has been given")
	}
	ty, err := encoder.abiIndex.TyByIdx(tyIdx)
	if err != nil {
		return err
	}
	if ty.SumType != parser.TyKindAliasRef {
		return fmt.Errorf("expected AliasRef at ty_idx=%d", tyIdx)
	}

	alias, err := encoder.abiIndex.GetAlias(ty.AliasRef.AliasName)
	if err != nil {
		return err
	}
	if alias.CustomPackUnpack.PackToBuilder {
		if encoder.customPackResolver == nil {
			return fmt.Errorf("custom pack resolver for alias %q is not configured", ty.AliasRef.AliasName)
		}
		return encoder.customPackResolver(*ty.AliasRef, cell, a)
	}

	targetTyIdx, _, err := encoder.abiIndex.AliasTargetOf(tyIdx)
	if err != nil {
		return fmt.Errorf("failed to resolve alias target: %w", err)
	}

	val := Value(*a)
	err = val.MarshalTyIdx(cell, targetTyIdx, encoder)
	if err != nil {
		return fmt.Errorf("failed to marshal alias value: %w", err)
	}

	return nil
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

func (g *GenericValue) Unmarshal(cell *boc.Cell, ty parser.GenericT, decoder *Decoder) error {
	currentTyIdx, found := decoder.genericRefs[ty.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic's type %v ", ty.NameT)
	}
	val := Value{}
	err := val.UnmarshalTyIdx(cell, currentTyIdx, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal generic's value: %w", err)
	}

	*g = GenericValue(val)

	return nil
}

func (g *GenericValue) Marshal(cell *boc.Cell, ty parser.GenericT, encoder *Encoder) error {
	currentTyIdx, found := encoder.genericRefs[ty.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic's type %v ", ty.NameT)
	}
	val := Value(*g)
	err := val.MarshalTyIdx(cell, currentTyIdx, encoder)
	if err != nil {
		return fmt.Errorf("failed to marshal generic's value: %w", err)
	}

	return nil
}
