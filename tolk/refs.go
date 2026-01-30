package tolk

import (
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"golang.org/x/exp/maps"
)

type EnumRef struct {
	EnumName string `json:"enumName"`
}

func (EnumRef) SetValue(v *Value, val any) error {
	e, ok := val.(EnumValue)
	if !ok {
		return fmt.Errorf("value is not an enum value")
	}
	v.enum = &e
	v.sumType = "enum"
	return nil
}

func (e EnumRef) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	if decoder.abiRefs.enumRefs == nil {
		return fmt.Errorf("struct has enum reference, but no abi has been given")
	}
	enum, found := decoder.abiRefs.enumRefs[e.EnumName]
	if !found {
		return fmt.Errorf("no enum with name %s was found in given abi", e.EnumName)
	}

	enumVal := Value{}
	err := enumVal.UnmarshalTolk(cell, enum.EncodedAs, decoder)
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
			bigEnumVal = big.Int(*enumVal.bigInt)
		} else {
			bigEnumVal = *new(big.Int).SetUint64(uint64(*enumVal.smallUint))
		}
	case "VarIntN", "VarUintN":
		bigEnumVal = big.Int(*enumVal.bigInt)
	default:
		return fmt.Errorf("enum encode type must be integer, got: %s", enum.EncodedAs.SumType)
	}

	for _, member := range enum.Members {
		val, ok := new(big.Int).SetString(member.Value, 10)
		if !ok {
			return fmt.Errorf("invalid enum %v value %v for member %s", e.EnumName, member.Value, member.Name)
		}

		if val.Cmp(&bigEnumVal) == 0 {
			err = e.SetValue(v, EnumValue{
				val:   enumVal,
				Name:  member.Name,
				Value: *val,
			})
			if err != nil {
				return err
			}

			return nil
		}
	}
	// todo: maybe return err?
	err = e.SetValue(v, EnumValue{
		val:   enumVal,
		Name:  "Unknown",
		Value: bigEnumVal,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *EnumValue) UnmarshalTolk(cell *boc.Cell, ty EnumRef, decoder *Decoder) error {
	if decoder.abiRefs.enumRefs == nil {
		return fmt.Errorf("struct has enum reference, but no abi has been given")
	}
	enum, found := decoder.abiRefs.enumRefs[ty.EnumName]
	if !found {
		return fmt.Errorf("no enum with name %s was found in given abi", ty.EnumName)
	}

	enumVal := Value{}
	err := enumVal.UnmarshalTolk(cell, enum.EncodedAs, decoder)
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
			bigEnumVal = big.Int(*enumVal.bigInt)
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

func (EnumRef) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.enum == nil {
	//	return fmt.Errorf("enum is nil")
	//}
	//
	//valType := v.enum.enumType
	//if valType.SumType != "IntN" && valType.SumType != "UintN" && valType.SumType != "VarIntN" && valType.SumType != "VarUintN" {
	//	return fmt.Errorf("enum type must be interger, got: %s", valType.SumType)
	//}
	//err := valType.MarshalTolk(cell, &v.enum.val)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (EnumRef) Equal(v Value, o Value) bool {
	return false
}

type StructRef struct {
	StructName string `json:"structName"`
	TypeArgs   []Ty   `json:"typeArgs,omitempty"`
}

func (StructRef) SetValue(v *Value, val any) error {
	s, ok := val.(Struct)
	if !ok {
		return fmt.Errorf("value is not a struct")
	}
	v.structValue = &s
	v.sumType = "structValue"
	return nil
}

func (s StructRef) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	if decoder.abiRefs.structRefs == nil {
		return fmt.Errorf("struct has struct reference, but no abi has been given")
	}
	strct, found := decoder.abiRefs.structRefs[s.StructName]
	if !found {
		return fmt.Errorf("no struct with name %s was found in given abi", s.StructName)
	}
	tolkStruct := Struct{
		field: make(map[string]Value),
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
		tolkStruct.prefix = TolkPrefix{
			Len:    int16(prefixLen),
			Prefix: prefix,
		}
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(s.TypeArgs, strct.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	for _, field := range strct.Fields {
		fieldVal := Value{}
		err = fieldVal.UnmarshalTolk(cell, field.Ty, decoder)
		if err != nil {
			return err
		}

		if field.Ty.SumType == "Nullable" {
			optVal := *fieldVal.optionalValue
			defVal := field.DefaultValue
			if !optVal.IsExists && defVal != nil {
				val := Value{}
				exists, err := defVal.unmarshalDefaultValue(&val, field.Ty)
				if err != nil {
					return err
				}
				if exists {
					optVal.IsExists = true
					optVal.Val = val
				}
			}
		}

		tolkStruct.field[field.Name] = fieldVal
	}
	decoder.abiRefs.genericRefs = oldGenericMap

	err = s.SetValue(v, tolkStruct)
	if err != nil {
		return err
	}

	return nil
}

func (s *Struct) UnmarshalTolk(cell *boc.Cell, ty StructRef, decoder *Decoder) error {
	if decoder.abiRefs.structRefs == nil {
		return fmt.Errorf("struct has struct reference, but no abi has been given")
	}
	strct, found := decoder.abiRefs.structRefs[ty.StructName]
	if !found {
		return fmt.Errorf("no struct with name %s was found in given abi", ty.StructName)
	}
	tolkStruct := Struct{
		field: make(map[string]Value),
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
		tolkStruct.prefix = TolkPrefix{
			Len:    int16(prefixLen),
			Prefix: prefix,
		}
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, strct.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	for _, field := range strct.Fields {
		fieldVal := Value{}
		err = fieldVal.UnmarshalTolk(cell, field.Ty, decoder)
		if err != nil {
			return err
		}

		if field.Ty.SumType == "Nullable" {
			optVal := *fieldVal.optionalValue
			defVal := field.DefaultValue
			if !optVal.IsExists && defVal != nil {
				val := Value{}
				exists, err := defVal.unmarshalDefaultValue(&val, field.Ty)
				if err != nil {
					return err
				}
				if exists {
					optVal.IsExists = true
					optVal.Val = val
				}
			}
		}

		tolkStruct.field[field.Name] = fieldVal
	}

	decoder.abiRefs.genericRefs = oldGenericMap
	*s = tolkStruct

	return nil
}

func (StructRef) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.structValue == nil {
	//	return fmt.Errorf("struct is nil")
	//}
	//
	//if v.structValue.hasPrefix {
	//	err := cell.WriteUint(v.structValue.prefix.Prefix, int(v.structValue.prefix.Len))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//for _, f := range v.structValue.field {
	//	err := f.valType.MarshalTolk(cell, &f)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (StructRef) Equal(v Value, o Value) bool {
	return false
}

type AliasRef struct {
	AliasName string `json:"aliasName"`
	TypeArgs  []Ty   `json:"typeArgs,omitempty"`
}

func (AliasRef) SetValue(v *Value, val any) error {
	return fmt.Errorf("alias cannot be a value")
}

func (a AliasRef) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	if decoder.abiRefs.aliasRefs == nil {
		return fmt.Errorf("struct has alias reference, but no abi has been given")
	}
	alias, found := decoder.abiRefs.aliasRefs[a.AliasName]
	if !found {
		return fmt.Errorf("no alias with name %s was found in given abi", a.AliasName)
	}

	if alias.CustomUnpackFromSlice {
		// todo: maybe simply return error?
		fmt.Println("WARNING! alias has custom unpack method. Standard unpacking can be incorrect!")
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(a.TypeArgs, alias.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	val := Value{}
	err = val.UnmarshalTolk(cell, alias.TargetTy, decoder)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = oldGenericMap

	fmt.Println(123)
	err = a.SetValue(v, val)
	if err != nil {
		return err
	}

	return nil
}

func (a *AliasValue) UnmarshalTolk(cell *boc.Cell, ty AliasRef, decoder *Decoder) error {
	if decoder.abiRefs.aliasRefs == nil {
		return fmt.Errorf("struct has alias reference, but no abi has been given")
	}
	alias, found := decoder.abiRefs.aliasRefs[ty.AliasName]
	if !found {
		return fmt.Errorf("no alias with name %s was found in given abi", ty.AliasName)
	}

	if alias.CustomUnpackFromSlice {
		// todo: maybe simply return error?
		fmt.Println("WARNING! alias has custom unpack method. Standard unpacking can be incorrect!")
	}

	oldGenericMap := decoder.abiRefs.genericRefs
	genericMap, err := resolveGeneric(ty.TypeArgs, alias.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiRefs.genericRefs = genericMap

	val := Value{}
	err = val.UnmarshalTolk(cell, alias.TargetTy, decoder)
	if err != nil {
		return err
	}

	decoder.abiRefs.genericRefs = oldGenericMap
	*a = AliasValue(val)

	return nil
}

func (AliasRef) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("alias ref cannot be a value")
}

func (AliasRef) Equal(v Value, o Value) bool {
	return false
}

type Generic struct {
	NameT string `json:"nameT"`
}

func (Generic) SetValue(v *Value, val any) error {
	return fmt.Errorf("generic cannot be a value")
}

func (g Generic) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	currentTy, found := decoder.abiRefs.genericRefs[g.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic type %v ", g)
	}
	val := Value{}
	err := val.UnmarshalTolk(cell, currentTy, decoder)
	if err != nil {
		return err
	}

	err = g.SetValue(v, val)
	if err != nil {
		return err
	}

	return nil
}

func (g *GenericValue) UnmarshalTolk(cell *boc.Cell, ty Generic, decoder *Decoder) error {
	currentTy, found := decoder.abiRefs.genericRefs[ty.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic type %v ", ty.NameT)
	}
	val := Value{}
	err := val.UnmarshalTolk(cell, currentTy, decoder)
	if err != nil {
		return err
	}

	*g = GenericValue(val)

	return nil
}

func (Generic) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("generic cannot be a value")
}

func (Generic) Equal(v Value, o Value) bool {
	return false
}

func resolveGeneric(typeArgs []Ty, typeParams []string, d *Decoder) (map[string]Ty, error) {
	genericMap := make(map[string]Ty)
	if d.abiRefs.genericRefs != nil {
		maps.Copy(genericMap, d.abiRefs.genericRefs)
	}

	for i, genericTy := range typeArgs {
		genericMap[typeParams[i]] = genericTy

		if genericTy.SumType == "Generic" {
			if d.abiRefs.genericRefs == nil {
				return nil, fmt.Errorf("cannot resolve generic type %v", genericTy.Generic.NameT)
			}

			ty, found := d.abiRefs.genericRefs[genericTy.Generic.NameT]
			if !found {
				return nil, fmt.Errorf("generic type %v not found", genericTy.Generic.NameT)
			}
			genericMap[typeParams[i]] = ty
		}
	}

	return genericMap, nil
}
