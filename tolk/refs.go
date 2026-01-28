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
	return nil
}

func (e EnumRef) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	if decoder.abiCtx == nil {
		return fmt.Errorf("struct has enum reference, but no abi has been given")
	}
	enum, found := decoder.abiCtx.enumRefs[e.EnumName]
	if !found {
		return fmt.Errorf("no enum with name %s was found in given abi", e.EnumName)
	}

	enumVal := Value{}
	err := enum.EncodedAs.UnmarshalTolk(cell, &enumVal, decoder)
	if err != nil {
		return err
	}
	var bigEnumVal big.Int
	switch enumVal.valType.SumType {
	case "IntN":
		if enumVal.valType.IntN.N > 64 {
			bigEnumVal = big.Int(*enumVal.bigInt)
		} else {
			bigEnumVal = *big.NewInt(int64(*enumVal.smallInt))
		}
	case "UintN":
		if enumVal.valType.UintN.N > 64 {
			bigEnumVal = big.Int(*enumVal.bigInt)
		} else {
			bigEnumVal = *new(big.Int).SetUint64(uint64(*enumVal.smallUint))
		}
	case "VarIntN", "VarUintN":
		bigEnumVal = big.Int(*enumVal.bigInt)
	default:
		return fmt.Errorf("enum encode type must be integer, got: %s", enumVal.valType.SumType)
	}

	for _, member := range enum.Members {
		val, ok := new(big.Int).SetString(member.Value, 10)
		if !ok {
			return fmt.Errorf("invalid enum %v value %v for member %s", e.EnumName, member.Value, member.Name)
		}

		if val.Cmp(&bigEnumVal) == 0 {
			err = v.SetValue(EnumValue{
				enumType: enum.EncodedAs,
				Name:     member.Name,
				Value:    *val,
			})
			if err != nil {
				return err
			}

			return nil
		}
	}
	// todo: maybe return err?
	err = v.SetValue(EnumValue{
		enumType: enum.EncodedAs,
		Name:     "Unknown",
		Value:    bigEnumVal,
	})
	if err != nil {
		return err
	}

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
	return nil
}

func (s StructRef) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	if decoder.abiCtx == nil {
		return fmt.Errorf("struct has struct reference, but no abi has been given")
	}
	strct, found := decoder.abiCtx.structRefs[s.StructName]
	if !found {
		return fmt.Errorf("no struct with name %s was found in given abi", s.StructName)
	}
	cont := Struct{
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
		cont.hasPrefix = true
		cont.prefix = TolkPrefix{
			Len:    int16(prefixLen),
			Prefix: prefix,
		}
	}

	oldGenericMap := decoder.abiCtx.genericRefs
	genericMap, err := resolveGeneric(s.TypeArgs, strct.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiCtx.genericRefs = genericMap

	for _, field := range strct.Fields {
		fieldVal := Value{}
		err = field.Ty.UnmarshalTolk(cell, &fieldVal, decoder)
		if err != nil {
			return err
		}

		if fieldVal.valType.SumType == "Nullable" {
			optVal := *fieldVal.optionalValue
			defVal := field.DefaultValue
			if !optVal.IsExists && defVal != nil {
				val := Value{}
				// todo remove decoder methods
				exists, err := defVal.UnmarshalDefaultValue(&val)
				if err != nil {
					return err
				}
				if exists {
					optVal.IsExists = true
					optVal.Val = val
				}
			}
		}

		cont.field[field.Name] = fieldVal
	}
	decoder.abiCtx.genericRefs = oldGenericMap

	err = v.SetValue(cont)
	if err != nil {
		return err
	}

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
	return fmt.Errorf("alias cannot be value")
}

func (a AliasRef) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	if decoder.abiCtx == nil {
		return fmt.Errorf("struct has alias reference, but no abi has been given")
	}
	alias, found := decoder.abiCtx.aliasRefs[a.AliasName]
	if !found {
		return fmt.Errorf("no alias with name %s was found in given abi", a.AliasName)
	}

	if alias.CustomUnpackFromSlice {
		// todo: maybe simply return error?
		fmt.Println("WARNING! alias has custom unpack method. Standard unpacking can be incorrect!")
	}

	oldGenericMap := decoder.abiCtx.genericRefs
	genericMap, err := resolveGeneric(a.TypeArgs, alias.TypeParams, decoder)
	if err != nil {
		return err
	}
	decoder.abiCtx.genericRefs = genericMap

	val := Value{}
	err = alias.TargetTy.UnmarshalTolk(cell, &val, decoder)
	if err != nil {
		return err
	}
	decoder.abiCtx.genericRefs = oldGenericMap

	innerVal, err := val.GetAny()
	if err != nil {
		return err
	}
	v.SetType(val.valType)
	err = v.SetValue(innerVal)
	if err != nil {
		return err
	}

	return nil
}

func (AliasRef) Equal(v Value, o Value) bool {
	return false
}

type Generic struct {
	NameT string `json:"nameT"`
}

func (Generic) SetValue(v *Value, val any) error {
	return fmt.Errorf("generic cannot be value")
}

func (g Generic) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	currentTy, found := decoder.abiCtx.genericRefs[g.NameT]
	if !found {
		return fmt.Errorf("cannot resolve generic type %v ", g)
	}
	val := Value{}
	err := currentTy.UnmarshalTolk(cell, &val, decoder)
	if err != nil {
		return err
	}

	innerVal, err := val.GetAny()
	if err != nil {
		return err
	}
	v.SetType(val.valType)
	err = v.SetValue(innerVal)
	if err != nil {
		return err
	}

	return nil
}

func (Generic) Equal(v Value, o Value) bool {
	return false
}

func resolveGeneric(typeArgs []Ty, typeParams []string, d *Decoder) (map[string]Ty, error) {
	genericMap := make(map[string]Ty)
	if d.abiCtx.genericRefs != nil {
		maps.Copy(genericMap, d.abiCtx.genericRefs)
	}

	for i, genericTy := range typeArgs {
		genericMap[typeParams[i]] = genericTy

		if genericTy.SumType == "Generic" {
			if d.abiCtx.genericRefs == nil {
				return nil, fmt.Errorf("cannot resolve generic type %v", genericTy.Generic.NameT)
			}

			ty, found := d.abiCtx.genericRefs[genericTy.Generic.NameT]
			if !found {
				return nil, fmt.Errorf("generic type %v not found", genericTy.Generic.NameT)
			}
			genericMap[typeParams[i]] = ty
		}
	}

	return genericMap, nil
}
