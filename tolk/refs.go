package tolk

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
	"golang.org/x/exp/maps"
)

type Prefix struct {
	Len    int16
	Prefix uint64
}

type Struct struct {
	hasPrefix   bool
	prefix      Prefix
	fieldNames  []string
	fieldValues []Value
}

func (s *Struct) Unmarshal(cell *boc.Cell, ty tolkParser.StructRef, decoder *Decoder) error {
	if decoder.abiRefs.structRefs == nil {
		return fmt.Errorf("struct has struct reference, but no abi has been given")
	}
	strct, found := decoder.abiRefs.structRefs[ty.StructName]
	if !found {
		return fmt.Errorf("no struct with name %s was found in given abi", ty.StructName)
	}
	tolkStruct := Struct{
		fieldNames:  make([]string, 0),
		fieldValues: make([]Value, 0),
	}
	if strct.Prefix != nil {
		prefixLen := strct.Prefix.PrefixLen
		if prefixLen > 64 {
			return fmt.Errorf("struct %v prefix length must be lower than 64", strct.Name)
		}

		prefix, err := cell.ReadUint(prefixLen)
		if err != nil {
			return err
		}
		actualPrefix, err := binHexToUint64(strct.Prefix.PrefixStr)
		if err != nil {
			return err
		}

		if prefix != actualPrefix {
			return fmt.Errorf("struct %v prefix does not match actual prefix %v", strct.Name, actualPrefix)
		}
		tolkStruct.hasPrefix = true
		tolkStruct.prefix = Prefix{
			Len:    int16(prefixLen),
			Prefix: prefix,
		}
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, strct.TypeParams, &decoder.abiRefs)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	for _, field := range strct.Fields {
		fieldVal := Value{}
		err = fieldVal.Unmarshal(cell, field.Ty, decoder)
		if err != nil {
			return err
		}

		if field.Ty.SumType == "Nullable" {
			optVal := *fieldVal.optionalValue
			defVal := field.DefaultValue
			if !optVal.IsExists && defVal != nil {
				val := Value{}
				exists, err := val.unmarshalDefaultValue(defVal, field.Ty)
				if err != nil {
					return err
				}
				if exists {
					optVal.IsExists = true
					optVal.Val = val
				}
			}
		}

		tolkStruct.fieldNames = append(tolkStruct.fieldNames, field.Name)
		tolkStruct.fieldValues = append(tolkStruct.fieldValues, fieldVal)
	}

	decoder.abiRefs.genericRefs = oldGenericMap
	*s = tolkStruct

	return nil
}

func (s *Struct) Marshal(cell *boc.Cell, ty tolkParser.StructRef, encoder *Encoder) error {
	if encoder.abiRefs.structRefs == nil {
		return fmt.Errorf("struct has struct reference, but no abi has been given")
	}
	strct, found := encoder.abiRefs.structRefs[ty.StructName]
	if !found {
		return fmt.Errorf("no struct with name %s was found in given abi", ty.StructName)
	}

	if strct.Prefix != nil {
		actualPrefix, err := binHexToUint64(strct.Prefix.PrefixStr)
		if err != nil {
			return err
		}

		if s.prefix.Prefix != actualPrefix {
			return fmt.Errorf("struct %v prefix does not match actual prefix %v", strct.Name, actualPrefix)
		}

		err = cell.WriteUint(s.prefix.Prefix, int(s.prefix.Len))
		if err != nil {
			return err
		}
	}

	oldGenericMap := encoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, strct.TypeParams, &encoder.abiRefs)
	if err != nil {
		return err
	}
	encoder.abiRefs.genericRefs = genericMap

	for _, field := range strct.Fields {
		val, ok := s.GetField(field.Name)
		if !ok {
			return fmt.Errorf("struct %v has no field %v", strct.Name, field.Name)
		}

		err = val.Marshal(cell, field.Ty, encoder)
		if err != nil {
			return err
		}
	}

	encoder.abiRefs.genericRefs = oldGenericMap

	return nil
}

func (s *Struct) GetField(field string) (Value, bool) {
	for i, name := range s.fieldNames {
		if name == field {
			return s.fieldValues[i], true
		}
	}
	return Value{}, false
}

func (s *Struct) SetField(field string, v Value) bool {
	for i, name := range s.fieldNames {
		if name == field {
			s.fieldValues[i] = v
			return true
		}
	}
	return false
}

func (s *Struct) RemoveField(field string) {
	for i, name := range s.fieldNames {
		if name == field {
			s.fieldNames[i] = "|" // impossible symbol for field name in Tolk language
			s.fieldValues[i] = Value{}
		}
	}
}

func (s *Struct) GetPrefix() (Prefix, bool) {
	if !s.hasPrefix {
		return Prefix{}, false
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
	if !slices.Equal(s.fieldNames, otherStruct.fieldNames) {
		return false
	}
	return slices.Equal(s.fieldValues, otherStruct.fieldValues)
}

type EnumValue struct {
	val   Value
	Name  string
	Value big.Int
}

func (e *EnumValue) Unmarshal(cell *boc.Cell, ty tolkParser.EnumRef, decoder *Decoder) error {
	if decoder.abiRefs.enumRefs == nil {
		return fmt.Errorf("struct has enum reference, but no abi has been given")
	}
	enum, found := decoder.abiRefs.enumRefs[ty.EnumName]
	if !found {
		return fmt.Errorf("no enum with name %s was found in given abi", ty.EnumName)
	}

	enumVal := Value{}
	err := enumVal.Unmarshal(cell, enum.EncodedAs, decoder)
	if err != nil {
		return err
	}
	var bigEnumVal big.Int
	switch enum.EncodedAs.SumType {
	case "IntN":
		if enum.EncodedAs.IntN.N > 64 {
			bigEnumVal = big.Int(*enumVal.bigInt)
		} else {
			bigEnumVal = *big.NewInt(int64(*enumVal.smallInt))
		}
	case "UintN":
		if enum.EncodedAs.UintN.N > 64 {
			bigEnumVal = big.Int(*enumVal.bigUint)
		} else {
			bigEnumVal = *new(big.Int).SetUint64(uint64(*enumVal.smallUint))
		}
	case "VarIntN":
		bigEnumVal = big.Int(*enumVal.varInt)
	case "VarUintN":
		bigEnumVal = big.Int(*enumVal.varUint)
	default:
		return fmt.Errorf("enum encode type must be integer, got: %s", enum.EncodedAs.SumType)
	}

	for _, member := range enum.Members {
		val, ok := new(big.Int).SetString(member.Value, 10)
		if !ok {
			return fmt.Errorf("invalid enum %v value %v for member %s", ty.EnumName, member.Value, member.Name)
		}

		if val.Cmp(&bigEnumVal) == 0 {
			*e = EnumValue{
				val:   enumVal,
				Name:  member.Name,
				Value: *val,
			}

			return nil
		}
	}
	// todo: maybe return err?
	*e = EnumValue{
		val:   enumVal,
		Name:  "Unknown",
		Value: bigEnumVal,
	}

	return nil
}

func (e *EnumValue) Marshal(cell *boc.Cell, ty tolkParser.EnumRef, encoder *Encoder) error {
	if encoder.abiRefs.enumRefs == nil {
		return fmt.Errorf("struct has enum reference, but no abi has been given")
	}
	enum, found := encoder.abiRefs.enumRefs[ty.EnumName]
	if !found {
		return fmt.Errorf("no enum with name %s was found in given abi", ty.EnumName)
	}

	err := e.val.Marshal(cell, enum.EncodedAs, encoder)
	if err != nil {
		return err
	}

	for _, member := range enum.Members {
		val, ok := new(big.Int).SetString(member.Value, 10)
		if !ok {
			return fmt.Errorf("invalid enum %v value %v for member %s", ty.EnumName, member.Value, member.Name)
		}

		if val.Cmp(&e.Value) == 0 {
			return nil
		}
	}
	return fmt.Errorf("enum value not matcher, got: %s", e.Value.String())
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

type AliasValue Value

func (a *AliasValue) Equal(other any) bool {
	otherAlias, ok := other.(AliasValue)
	if !ok {
		return false
	}
	v := Value(*a)
	return v.Equal(Value(otherAlias))
}

func (a *AliasValue) Unmarshal(cell *boc.Cell, ty tolkParser.AliasRef, decoder *Decoder) error {
	if decoder.abiRefs.aliasRefs == nil {
		return fmt.Errorf("struct has alias reference, but no abi has been given")
	}
	alias, found := decoder.abiRefs.aliasRefs[ty.AliasName]
	if !found {
		return fmt.Errorf("no alias with name %s was found in given abi", ty.AliasName)
	}

	if alias.CustomUnpackFromSlice {
		// todo: maybe simply return error?
		fmt.Println("WARNING! alias has custom unpack method. Default unpacking can be incorrect!")
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, alias.TypeParams, &decoder.abiRefs)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	val := Value{}
	err = val.Unmarshal(cell, alias.TargetTy, decoder)
	if err != nil {
		return err
	}

	decoder.abiRefs.genericRefs = oldGenericMap
	*a = AliasValue(val)

	return nil
}

func (a *AliasValue) Marshal(cell *boc.Cell, ty tolkParser.AliasRef, encoder *Encoder) error {
	if encoder.abiRefs.aliasRefs == nil {
		return fmt.Errorf("struct has alias reference, but no abi has been given")
	}
	alias, found := encoder.abiRefs.aliasRefs[ty.AliasName]
	if !found {
		return fmt.Errorf("no alias with name %s was found in given abi", ty.AliasName)
	}

	if alias.CustomPackToBuilder {
		// todo: maybe simply return error?
		fmt.Println("WARNING! alias has custom pack method. Default packing can be incorrect!")
	}

	oldGenericMap := encoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, alias.TypeParams, &encoder.abiRefs)
	if err != nil {
		return err
	}
	encoder.abiRefs.genericRefs = genericMap

	val := Value(*a)
	err = val.Marshal(cell, alias.TargetTy, encoder)
	if err != nil {
		return err
	}

	encoder.abiRefs.genericRefs = oldGenericMap

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

func (g *GenericValue) Unmarshal(cell *boc.Cell, ty tolkParser.Generic, decoder *Decoder) error {
	currentTy, found := decoder.abiRefs.genericRefs[ty.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic type %v ", ty.NameT)
	}
	val := Value{}
	err := val.Unmarshal(cell, currentTy, decoder)
	if err != nil {
		return err
	}

	*g = GenericValue(val)

	return nil
}

func (g *GenericValue) Marshal(cell *boc.Cell, ty tolkParser.Generic, encoder *Encoder) error {
	currentTy, found := encoder.abiRefs.genericRefs[ty.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic type %v ", ty.NameT)
	}
	val := Value(*g)
	err := val.Marshal(cell, currentTy, encoder)
	if err != nil {
		return err
	}

	return nil
}

func resolveGeneric(typeArgs []tolkParser.Ty, typeParams []string, abiRefs *abiRefs) (map[string]tolkParser.Ty, error) {
	genericMap := make(map[string]tolkParser.Ty)
	if abiRefs.genericRefs != nil {
		maps.Copy(genericMap, abiRefs.genericRefs)
	}

	for i, genericTy := range typeArgs {
		genericMap[typeParams[i]] = genericTy

		if genericTy.SumType == "Generic" {
			if abiRefs.genericRefs == nil {
				return nil, fmt.Errorf("cannot resolve generic type %v", genericTy.Generic.NameT)
			}

			ty, found := abiRefs.genericRefs[genericTy.Generic.NameT]
			if !found {
				return nil, fmt.Errorf("generic type %v not found", genericTy.Generic.NameT)
			}
			genericMap[typeParams[i]] = ty
		}
	}

	return genericMap, nil
}
