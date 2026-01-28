package tolkParser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/tolk"
	"github.com/tonkeeper/tongo/utils"
)

type DefaultType struct {
	Name  string
	IsRef bool
}

var (
	defaultKnownTypes = map[string]DefaultType{
		"Int":        {"tlb.Int257", false},
		"Coins":      {"tlb.VarUInteger16", false},
		"Bool":       {"bool", false},
		"Cell":       {"tlb.Any", true},
		"Slice":      {"tlb.Any", false},
		"Builder":    {"tlb.Any", false},
		"Remaining":  {"tlb.Any", false},
		"Address":    {"tlb.MsgAddress", false},
		"AddressOpt": {"tlb.MsgAddress", false},
		"AddressExt": {"tlb.MsgAddress", false},
		"AddressAny": {"tlb.MsgAddress", false},
		"Void":       {"tlb.Void", false},
	}
)

type DeclrResult struct {
	Tag  Tag
	Name string
	Code string
}

type MsgResult struct {
	Code string
}

type Tag struct {
	Len int
	Val uint64
}

func ParseStructMsg(ty tolk.Ty, msgName, contractName string) (*MsgResult, error) {
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(msgName)
	res.WriteRune(' ')
	code, err := createGolangAlias(ty.StructRef.StructName, ty.StructRef.TypeArgs, contractName)
	if err != nil {
		return nil, err
	}
	res.WriteString(code)

	return &MsgResult{
		Code: res.String(),
	}, nil
}

func ParseAliasMsg(ty tolk.Ty, msgName, contractName string) (*MsgResult, error) {
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(msgName)
	res.WriteRune(' ')
	code, err := createGolangAlias(ty.AliasRef.AliasName, ty.AliasRef.TypeArgs, contractName)
	if err != nil {
		return nil, err
	}
	res.WriteString(code)

	return &MsgResult{
		Code: res.String(),
	}, nil
}

func ParseGetMethodCode(ty tolk.Ty, contractName string) (string, error) {
	switch ty.SumType {
	case "StructRef":
		name := utils.ToCamelCase(ty.StructRef.StructName)
		return createGolangAlias(name, ty.StructRef.TypeArgs, contractName)
	case "AliasRef":
		name := utils.ToCamelCase(ty.AliasRef.AliasName)
		return createGolangAlias(name, ty.AliasRef.TypeArgs, contractName)
	case "Tensor":
		var res strings.Builder
		res.WriteString("struct {\nField ")
		tlbType, _, err := parseTy(ty, contractName)
		if err != nil {
			return "", err
		}
		res.WriteString(tlbType)
		res.WriteString("`vmStackHint:\"tensor\"`\n}\n")

		return res.String(), nil
	case "TupleWith":
		var res strings.Builder
		res.WriteString("struct {\nField ")
		tlbType, _, err := parseTy(ty, contractName)
		if err != nil {
			return "", err
		}
		res.WriteString(tlbType)
		res.WriteString("\n}\n")

		return res.String(), nil
	default:
		var res strings.Builder
		res.WriteString("struct {\nValue ")
		tlbType, tlbTag, err := parseTy(ty, contractName)
		if err != nil {
			return "", err
		}
		wrappedType := wrapTlbTypeIntoTlbTag(tlbType, tlbTag)
		res.WriteString(wrappedType)
		res.WriteString("\n}\n")

		return res.String(), nil
	}
}

func createGolangAlias(aliasType string, typeArgs []tolk.Ty, contractName string) (string, error) {
	var res strings.Builder
	res.WriteString(" = ")
	res.WriteString(contractName + aliasType)

	err := writeGenericTypesPart(&res, typeArgs, contractName)
	if err != nil {
		return "", err
	}

	res.WriteRune('\n')
	return res.String(), nil
}

func ParseTag(
	ty tolk.Ty,
	structRefs map[string]tolk.StructDeclaration,
	aliasRefs map[string]tolk.AliasDeclaration,
	enumRefs map[string]tolk.EnumDeclaration,
) (*Tag, error) {
	switch ty.SumType {
	case "StructRef":
		prefix := structRefs[ty.StructRef.StructName].Prefix
		if prefix == nil {
			return nil, fmt.Errorf("%v tag is nil", ty.StructRef.StructName)
		}
		convertedTag, err := convertPrefixTagToInt(*prefix)
		if err != nil {
			return nil, err
		}
		return convertedTag, nil
	case "AliasRef":
		aliasRef := aliasRefs[ty.AliasRef.AliasName]
		return ParseTag(aliasRef.TargetTy, structRefs, aliasRefs, enumRefs)
	case "EnumRef":
		enumRef := enumRefs[ty.EnumRef.EnumName]
		return ParseTag(enumRef.EncodedAs, structRefs, aliasRefs, enumRefs)
	default:
		return nil, fmt.Errorf("cannot get tag from %s type", ty.SumType)
	}
}

func ParseStructDeclr(declr tolk.StructDeclaration, contractName string) (*DeclrResult, error) {
	result := DeclrResult{}

	unionStructs := make([]string, 0)
	var res strings.Builder
	if declr.Prefix != nil {
		convertedTag, err := convertPrefixTagToInt(*declr.Prefix)
		if err != nil {
			return nil, err
		}
		result.Tag = *convertedTag
	}
	typeName := contractName + utils.ToCamelCase(declr.Name)
	res.WriteString("type ")
	res.WriteString(typeName)
	writeGenericStructPart(&res, declr.TypeParams)
	res.WriteString(" struct {\n")
	for _, field := range declr.Fields {
		tlbField, unionStruct, err := parseField(typeName, field, contractName)
		if err != nil {
			return nil, err
		}
		if unionStruct != "" {
			unionStructs = append(unionStructs, unionStruct)
		}
		res.WriteString(tlbField)
		res.WriteRune('\n')
	}
	res.WriteString("}\n\n")

	for _, unionStruct := range unionStructs {
		res.WriteString(unionStruct)
		res.WriteRune('\n')
	}

	result.Code = res.String()
	result.Name = typeName

	return &result, nil
}

func ParseAliasDeclr(declr tolk.AliasDeclaration, contractName string) (*DeclrResult, error) {
	result := DeclrResult{}

	typeName := contractName + utils.ToCamelCase(declr.Name)
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(typeName)
	writeGenericStructPart(&res, declr.TypeParams)
	res.WriteRune(' ')
	tlbType, tlbTag, err := parseTy(declr.TargetTy, contractName)
	if err != nil {
		return nil, err
	}
	wrappedType := wrapTlbTypeIntoTlbTag(tlbType, tlbTag)

	res.WriteString(wrappedType)
	res.WriteRune('\n')

	if declr.TargetTy.SumType == "Union" {
		res.WriteString(generateJSONMarshalForUnion(declr.TargetTy, typeName))
	}

	result.Code = res.String()
	result.Name = typeName

	return &result, nil
}

func ParseEnumDeclr(declr tolk.EnumDeclaration, contractName string) (*DeclrResult, error) {
	result := DeclrResult{}

	isBigInt := false
	// enum is always an integer
	// if number has >64 bits or its varint/varuint then type should be bigInt
	switch declr.EncodedAs.SumType {
	case "IntN":
		if declr.EncodedAs.IntN.N > 64 {
			isBigInt = true
		}
	case "UintN":
		if declr.EncodedAs.UintN.N > 64 {
			isBigInt = true
		}
	case "VarIntN", "VarUintN":
		isBigInt = true
	default:
		return nil, fmt.Errorf("enum encode type must be integer, got: %s", declr.EncodedAs.SumType)
	}

	var res strings.Builder
	enumName := contractName + utils.ToCamelCase(declr.Name)
	res.WriteString("type ")
	res.WriteString(enumName)
	res.WriteString(" = ")

	tlbType, tlbTag, err := parseTy(declr.EncodedAs, contractName)
	if err != nil {
		return nil, err
	}
	wrappedType := wrapTlbTypeIntoTlbTag(tlbType, tlbTag)

	res.WriteString(wrappedType)
	if isBigInt {
		res.WriteString("\n\n  var (\n")
	} else {
		res.WriteString("\n\n  const (\n")
	}
	for _, member := range declr.Members {
		res.WriteString(enumName)
		res.WriteRune('_')
		res.WriteString(utils.ToCamelCase(member.Name))
		res.WriteString(" = ")
		res.WriteString(enumName)
		res.WriteRune('(')
		if isBigInt {
			res.WriteString("MustBigInt(\"")
			res.WriteString(member.Value)
			res.WriteString("\")")
		} else {
			res.WriteString(member.Value)
		}
		res.WriteString(")\n")
	}
	res.WriteString(")\n\n")

	result.Code = res.String()
	result.Name = enumName

	return &result, nil
}

func parseField(structName string, field tolk.Field, contractName string) (string, string, error) {
	var res strings.Builder
	var unionStruct strings.Builder
	name := utils.ToCamelCase(field.Name)
	res.WriteString(name)
	res.WriteRune(' ')

	if field.IsPayload != nil && *field.IsPayload == true {
		res.WriteString("tlb.EitherRef[Payload]")

		return res.String(), "", nil
	}

	tlbTypePart, tlbTagPart, err := parseTy(field.Ty, contractName)
	if err != nil {
		return "", "", err
	}
	// union and not tlb.Either
	if field.Ty.SumType == "Union" && len(field.Ty.Union.Variants) > 2 {
		res.WriteString(structName + name)

		unionStruct.WriteString("type ")
		unionStruct.WriteString(structName + name)
		unionStruct.WriteString(" ")
		unionStruct.WriteString(tlbTypePart)
		unionStruct.WriteString("\n\n")

		unionStruct.WriteString(generateJSONMarshalForUnion(field.Ty, name))
	} else {
		res.WriteString(tlbTypePart)
	}
	if tlbTagPart != "" {
		res.WriteRune(' ')
		res.WriteString(fmt.Sprintf("`tlb:\"%v\"`", tlbTagPart))
	}

	return res.String(), unionStruct.String(), nil
}

func generateJSONMarshalForUnion(union tolk.Ty, structName string) string {
	if len(union.Union.Variants) == 2 { // tlb.Either
		return ""
	}
	var res strings.Builder
	res.WriteString(fmt.Sprintf("func (t *%v) MarshalJSON() ([]byte, error) {\n", structName))
	res.WriteString("    switch t.SumType {")
	for i := range union.Union.Variants {
		name := fmt.Sprintf("UnionVariant%v", i)
		res.WriteString(fmt.Sprintf(`case "%v": `, name))
		res.WriteString(fmt.Sprintf(`bytes, err := json.Marshal(t.%v)`, name))
		res.WriteString("\nif err != nil {\n")
		res.WriteString("\treturn nil, err\n")
		res.WriteString("}\n")
		res.WriteString("return []byte(fmt.Sprintf(`{\"SumType\": \"" + name + "\",\"" + name + "\":%v}`, string(bytes))), nil\n")
	}
	res.WriteString("default: ")
	res.WriteString(`return nil, fmt.Errorf("unknown sum type %v", t.SumType)`)
	res.WriteString("}\n")
	res.WriteString("}\n")
	return res.String()
}

func ParseType(ty tolk.Ty, contractName string) (string, error) {
	tlbType, tlbTag, err := parseTy(ty, contractName)
	if err != nil {
		return "", err
	}

	return wrapTlbTypeIntoTlbTag(tlbType, tlbTag), nil
}

// todo: избавится от огромного switch case
func parseTy(ty tolk.Ty, contractName string) (string, string, error) {
	var tlbType strings.Builder
	var tlbTag strings.Builder
	defaultType, ok := defaultKnownTypes[ty.SumType]
	if ok {
		if defaultType.IsRef {
			tlbTag.WriteRune('^')
		}
		tlbType.WriteString(defaultType.Name)
		return tlbType.String(), tlbTag.String(), nil
	}

	switch ty.SumType {
	case "IntN":
		n := ty.IntN.N
		tlbIntType := fmt.Sprintf("tlb.Int%v", n)
		if n == 8 {
			tlbIntType = "int8"
		} else if n == 16 {
			tlbIntType = "int16"
		} else if n == 32 {
			tlbIntType = "int32"
		} else if n == 64 {
			tlbIntType = "int64"
		}
		tlbType.WriteString(tlbIntType)
	case "UintN":
		n := ty.UintN.N
		tlbUIntType := fmt.Sprintf("tlb.Uint%v", n)
		if n == 8 {
			tlbUIntType = "uint8"
		} else if n == 16 {
			tlbUIntType = "uint16"
		} else if n == 32 {
			tlbUIntType = "uint32"
		} else if n == 64 {
			tlbUIntType = "uint64"
		}
		tlbType.WriteString(tlbUIntType)
	case "VarIntN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.VarIntN.N)
		tlbType.WriteString(tlbVarUIntType)
	case "VarUintN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.VarUintN.N)
		tlbType.WriteString(tlbVarUIntType)
	case "BitsN":
		tlbVarUIntType := fmt.Sprintf("tlb.Bits%v", ty.BitsN.N)
		tlbType.WriteString(tlbVarUIntType)
	case "Nullable":
		tlbMaybeType, tlbMaybeTag, err := parseTy(ty.Nullable.Inner, contractName)
		if err != nil {
			return "", "", err
		}

		tlbTag.WriteString("maybe")
		tlbTag.WriteString(tlbMaybeTag)

		tlbType.WriteRune('*')
		tlbType.WriteString(tlbMaybeType)
		tlbType.WriteRune(' ')
	case "CellOf":
		tlbRefType, tlbRefTag, err := parseTy(ty.CellOf.Inner, contractName)
		if err != nil {
			return "", "", err
		}
		tlbTag.WriteRune('^')
		tlbTag.WriteString(tlbRefTag)

		tlbType.WriteString(tlbRefType)
		tlbType.WriteRune(' ')
	case "Tensor":
		tlbType.WriteString("struct{\n")
		for i, innerTy := range ty.Tensor.Items {
			innerTlbType, innerTlbTag, err := parseTy(innerTy, contractName)
			if err != nil {
				return "", "", err
			}

			wrappedType := wrapTlbTypeIntoTlbTag(innerTlbType, innerTlbTag)

			tlbType.WriteString("Field")
			tlbType.WriteString(strconv.Itoa(i))
			tlbType.WriteRune(' ')
			tlbType.WriteString(wrappedType)
			tlbType.WriteRune('\n')
		}
		tlbType.WriteRune('}')
	case "TupleWith":
		tlbType.WriteString("struct{\n")
		for i, innerTy := range ty.TupleWith.Items {
			innerTlbType, innerTlbTag, err := parseTy(innerTy, contractName)
			if err != nil {
				return "", "", err
			}

			wrappedType := wrapTlbTypeIntoTlbTag(innerTlbType, innerTlbTag)

			tlbType.WriteString("Field")
			tlbType.WriteString(strconv.Itoa(i))
			tlbType.WriteRune(' ')
			tlbType.WriteString(wrappedType)
			tlbType.WriteRune('\n')
		}
		tlbType.WriteRune('}')
	case "Map":
		tlbType.WriteString("tlb.HashmapE[")
		fixedType, err := getFixedSizeTypeForMap(ty.Map.K)
		if err != nil {
			return "", "", err
		}
		valueType, valueTypeTag, err := parseTy(ty.Map.V, contractName)
		wrappedType := wrapTlbTypeIntoTlbTag(valueType, valueTypeTag)

		tlbType.WriteString(fixedType)
		tlbType.WriteRune(',')
		tlbType.WriteString(wrappedType)
		tlbType.WriteRune(']')
	case "EnumRef":
		name := contractName + utils.ToCamelCase(ty.EnumRef.EnumName)
		tlbType.WriteString(name)
	case "StructRef":
		name := contractName + utils.ToCamelCase(ty.StructRef.StructName)
		tlbType.WriteString(name)
		err := writeGenericTypesPart(&tlbType, ty.StructRef.TypeArgs, contractName)
		if err != nil {
			return "", "", err
		}
	case "AliasRef":
		name := contractName + utils.ToCamelCase(ty.AliasRef.AliasName)
		tlbType.WriteString(name)
		err := writeGenericTypesPart(&tlbType, ty.AliasRef.TypeArgs, contractName)
		if err != nil {
			return "", "", err
		}
	case "Generic":
		name := ty.Generic.NameT
		tlbType.WriteString(utils.ToCamelCase(name))
	case "Union":
		if len(ty.Union.Variants) == 2 && ty.Union.Variants[0].PrefixStr == "0" {
			tlbType.WriteString("tlb.Either[")

			for _, variant := range ty.Union.Variants {
				tlbTypeArg, tlbTagArg, err := parseTy(variant.VariantTy, contractName)
				if err != nil {
					return "", "", err
				}
				wrappedType := wrapTlbTypeIntoTlbTag(tlbTypeArg, tlbTagArg)
				tlbType.WriteString(wrappedType)

				tlbType.WriteRune(',')
			}

			tlbType.WriteRune(']')
		} else {
			tlbType.WriteString("struct{\ntlb.SumType\n")
			for i, variant := range ty.Union.Variants {
				sumTypeSymbol := "$"
				if strings.HasPrefix(variant.PrefixStr, "0x") {
					sumTypeSymbol = "#"
				}

				tlbType.WriteString(fmt.Sprintf("UnionVariant%v ", i))

				tlbTypeArg, tlbTagArg, err := parseTy(variant.VariantTy, contractName)
				if err != nil {
					return "", "", err
				}
				tlbType.WriteString(tlbTypeArg)

				tlbType.WriteString(fmt.Sprintf(" `tlbSumType:\"%v%v\"`\n", sumTypeSymbol, variant.PrefixStr[2:]))
				tlbTag.WriteString(tlbTagArg)
			}
			tlbType.WriteRune('}')
		}
	case "NullLiteral":
		tlbType.WriteString("tlb.NullLiteral")
	case "TupleAny", "Callable":
		return "", "", fmt.Errorf("cannot convert type %v", ty.SumType)
	default:
		return "", "", fmt.Errorf("unknown ty type %q", ty.SumType)
	}

	return tlbType.String(), tlbTag.String(), nil
}

func getFixedSizeTypeForMap(ty tolk.Ty) (string, error) {
	switch ty.SumType {
	case "IntN":
		n := ty.IntN.N
		tlbIntType := fmt.Sprintf("tlb.Int%v", n)
		return tlbIntType, nil
	case "UintN":
		n := ty.UintN.N
		tlbIntType := fmt.Sprintf("tlb.Uint%v", n)
		return tlbIntType, nil
	case "VarIntN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.VarIntN.N)
		return tlbVarUIntType, nil
	case "VarUintN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.VarUintN.N)
		return tlbVarUIntType, nil
	case "BitsN":
		tlbVarUIntType := fmt.Sprintf("tlb.Bits%v", ty.BitsN.N)
		return tlbVarUIntType, nil
	case "Bool":
		return "tlb.Uint1", nil
	case "Address", "AddressOpt", "AddressExt", "AddressAny":
		return "tlb.InternalAddress", nil
	default:
		return "", fmt.Errorf("%v not supported as map key", ty.SumType)
	}
}

func wrapTlbTypeIntoTlbTag(tlbTypeArg, tlbTagArg string) string {
	var tlbType strings.Builder
	cnt := 0
	i := 0
	for i < len(tlbTagArg) {
		if tlbTagArg[i] == 'm' { // maybe
			i += 5
			tlbType.WriteString("tlb.Maybe[")
		} else if tlbTagArg[i] == '^' { // ref
			i += 1
			tlbType.WriteString("tlb.Ref[")
		} else {
			panic(fmt.Sprintf("unknown tag start %v in %v", tlbTagArg[i], tlbTagArg))
		}

		cnt += 1
	}
	if tlbTypeArg[0] == '*' { // remove pointer since it wrapping into type
		tlbTypeArg = tlbTypeArg[1:]
	}
	tlbType.WriteString(tlbTypeArg)
	tlbType.WriteString(strings.Repeat("]", cnt))

	return tlbType.String()
}

func convertPrefixTagToInt(tag tolk.Prefix) (*Tag, error) {
	if tag.PrefixLen == 0 {
		return nil, fmt.Errorf("prefix tag len must be > 0")
	}

	val, err := tolk.PrefixToUint(tag.PrefixStr)
	if err != nil {
		return nil, err
	}

	return &Tag{
		Len: tag.PrefixLen,
		Val: val,
	}, nil
}

func writeGenericStructPart(builder *strings.Builder, typeParams []string) {
	if len(typeParams) > 0 {
		builder.WriteRune('[')
		for _, param := range typeParams {
			builder.WriteString(utils.ToCamelCase(param))
			builder.WriteString(" any,")
		}
		builder.WriteRune(']')
	}
}

func writeGenericTypesPart(builder *strings.Builder, typeArgs []tolk.Ty, contractName string) error {
	if len(typeArgs) > 0 {
		builder.WriteRune('[')
		for _, typeArg := range typeArgs {
			tlbType, tlbTag, err := parseTy(typeArg, contractName)
			if err != nil {
				return err
			}

			wrappedType := wrapTlbTypeIntoTlbTag(tlbType, tlbTag)
			builder.WriteString(wrappedType)
			builder.WriteRune(',')
		}
		builder.WriteRune(']')
	}

	return nil
}
