package tolkParser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tolkAbi "github.com/tonkeeper/tongo/tolk/abi"
	"github.com/tonkeeper/tongo/utils"
)

type DefaultType struct {
	Name  string
	IsRef bool
}

var (
	defaultKnownTypes = map[string]DefaultType{
		"int":        {"tlb.Int257", false},
		"coins":      {"tlb.VarUInteger16", false},
		"bool":       {"bool", false},
		"cell":       {"tlb.Any", true},
		"slice":      {"tlb.Any", false},
		"builder":    {"tlb.Any", false},
		"remaining":  {"tlb.Any", false},
		"address":    {"tlb.MsgAddress", false},
		"addressOpt": {"tlb.MsgAddress", false},
		"addressExt": {"tlb.MsgAddress", false},
		"addressAny": {"tlb.MsgAddress", false},
		"void":       {"tlb.Void", false},
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

func ParseStructMsg(ty tolkAbi.Ty, msgName, contractName string) (*MsgResult, error) {
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(msgName)
	res.WriteRune(' ')
	code, err := createGolangAlias(ty.StructRefTy.StructName, ty.StructRefTy.TypeArgs, contractName)
	if err != nil {
		return nil, err
	}
	res.WriteString(code)

	return &MsgResult{
		Code: res.String(),
	}, nil
}

func ParseAliasMsg(ty tolkAbi.Ty, msgName, contractName string) (*MsgResult, error) {
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(msgName)
	res.WriteRune(' ')
	code, err := createGolangAlias(ty.AliasRefTy.AliasName, ty.AliasRefTy.TypeArgs, contractName)
	if err != nil {
		return nil, err
	}
	res.WriteString(code)

	return &MsgResult{
		Code: res.String(),
	}, nil
}

func ParseGetMethodCode(ty tolkAbi.Ty, contractName string) (string, error) {
	switch ty.SumType {
	case "StructRef":
		name := utils.ToCamelCase(ty.StructRefTy.StructName)
		return createGolangAlias(name, ty.StructRefTy.TypeArgs, contractName)
	case "AliasRef":
		name := utils.ToCamelCase(ty.AliasRefTy.AliasName)
		return createGolangAlias(name, ty.AliasRefTy.TypeArgs, contractName)
	case "tensor":
		var res strings.Builder
		res.WriteString("struct {\nField ")
		tlbType, _, err := parseTy(ty, contractName)
		if err != nil {
			return "", err
		}
		res.WriteString(tlbType)
		res.WriteString("`vmStackHint:\"tensor\"`\n}\n")

		return res.String(), nil
	case "tupleWith":
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

func createGolangAlias(aliasType string, typeArgs []tolkAbi.Ty, contractName string) (string, error) {
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
	ty tolkAbi.Ty,
	structRefs map[string]tolkAbi.StructDeclaration,
	aliasRefs map[string]tolkAbi.AliasDeclaration,
	enumRefs map[string]tolkAbi.EnumDeclaration,
) (*Tag, error) {
	switch ty.SumType {
	case "StructRef":
		prefix := structRefs[ty.StructRefTy.StructName].Prefix
		if prefix == nil {
			return nil, fmt.Errorf("%v tag is nil", ty.StructRefTy.StructName)
		}
		convertedTag, err := convertPrefixTagToInt(*prefix)
		if err != nil {
			return nil, err
		}
		return convertedTag, nil
	case "AliasRef":
		aliasRef := aliasRefs[ty.AliasRefTy.AliasName]
		if aliasRef.TargetTy == nil {
			return nil, fmt.Errorf("alias ref ty cannot be nil")
		}
		return ParseTag(*aliasRef.TargetTy, structRefs, aliasRefs, enumRefs)
	case "EnumRef":
		enumRef := enumRefs[ty.EnumRefTy.EnumName]
		if enumRef.EncodedAs == nil {
			return nil, fmt.Errorf("enum ref ty cannot be nil")
		}
		return ParseTag(*enumRef.EncodedAs, structRefs, aliasRefs, enumRefs)
	default:
		return nil, fmt.Errorf("cannot get tag from %s type", ty.SumType)
	}
}

func ParseStructDeclr(declr tolkAbi.StructDeclaration, contractName string) (*DeclrResult, error) {
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

func ParseAliasDeclr(declr tolkAbi.AliasDeclaration, contractName string) (*DeclrResult, error) {
	result := DeclrResult{}

	if declr.TargetTy == nil {
		return nil, errors.New("target ty cannot be null")
	}

	typeName := contractName + utils.ToCamelCase(declr.Name)
	var res strings.Builder
	res.WriteString("type ")
	res.WriteString(typeName)
	writeGenericStructPart(&res, declr.TypeParams)
	res.WriteRune(' ')
	tlbType, tlbTag, err := parseTy(*declr.TargetTy, contractName)
	if err != nil {
		return nil, err
	}
	wrappedType := wrapTlbTypeIntoTlbTag(tlbType, tlbTag)

	res.WriteString(wrappedType)
	res.WriteRune('\n')

	if declr.TargetTy.SumType == "union" {
		res.WriteString(generateJSONMarshalForUnion(*declr.TargetTy, typeName))
	}

	result.Code = res.String()
	result.Name = typeName

	return &result, nil
}

func ParseEnumDeclr(declr tolkAbi.EnumDeclaration, contractName string) (*DeclrResult, error) {
	result := DeclrResult{}

	if declr.EncodedAs == nil {
		return nil, errors.New("target ty cannot be null")
	}

	isBigInt := false
	// enum is always an integer
	// if number has >64 bits or its varint/varuint then type should be bigInt
	if declr.EncodedAs.NumberTy.N > 64 || strings.HasPrefix(declr.EncodedAs.SumType, "var") {
		isBigInt = true
	}

	var res strings.Builder
	enumName := contractName + utils.ToCamelCase(declr.Name)
	res.WriteString("type ")
	res.WriteString(enumName)
	res.WriteString(" = ")

	tlbType, tlbTag, err := parseTy(*declr.EncodedAs, contractName)
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

func parseField(structName string, field tolkAbi.Field, contractName string) (string, string, error) {
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
	if field.Ty.SumType == "union" && len(field.Ty.Union.Variants) > 2 {
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

func generateJSONMarshalForUnion(union tolkAbi.Ty, structName string) string {
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

func ParseType(ty tolkAbi.Ty, contractName string) (string, error) {
	tlbType, tlbTag, err := parseTy(ty, contractName)
	if err != nil {
		return "", err
	}

	return wrapTlbTypeIntoTlbTag(tlbType, tlbTag), nil
}

func parseTy(ty tolkAbi.Ty, contractName string) (string, string, error) {
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
	case "intN":
		n := ty.NumberTy.N
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
	case "uintN":
		n := ty.NumberTy.N
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
	case "varintN", "varuintN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.NumberTy.N)
		tlbType.WriteString(tlbVarUIntType)
	case "bitsN":
		tlbVarUIntType := fmt.Sprintf("tlb.Bits%v", ty.BitsTy.N)
		tlbType.WriteString(tlbVarUIntType)
	case "nullable":
		if ty.NullableTy.Inner == nil {
			return "", "", fmt.Errorf("inner nullable type cannot be null")
		}
		tlbMaybeType, tlbMaybeTag, err := parseTy(*ty.NullableTy.Inner, contractName)
		if err != nil {
			return "", "", err
		}

		tlbTag.WriteString("maybe")
		tlbTag.WriteString(tlbMaybeTag)

		tlbType.WriteRune('*')
		tlbType.WriteString(tlbMaybeType)
		tlbType.WriteRune(' ')
	case "cellOf":
		if ty.CellOf.Inner == nil {
			return "", "", fmt.Errorf("inner cell type cannot be null")
		}
		tlbRefType, tlbRefTag, err := parseTy(*ty.CellOf.Inner, contractName)
		if err != nil {
			return "", "", err
		}
		tlbTag.WriteRune('^')
		tlbTag.WriteString(tlbRefTag)

		tlbType.WriteString(tlbRefType)
		tlbType.WriteRune(' ')
	case "tensor":
		tlbType.WriteString("struct{\n")
		for i, innerTy := range ty.TensorTy.Items {
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
	case "tupleWith":
		tlbType.WriteString("struct{\n")
		for i, innerTy := range ty.TupleWithTy.Items {
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
	case "mapKV":
		if ty.MapTy.K == nil {
			return "", "", fmt.Errorf("map key cannot be null")
		}
		if ty.MapTy.V == nil {
			return "", "", fmt.Errorf("map value cannot be null")
		}
		tlbType.WriteString("tlb.HashmapE[")
		fixedType, err := getFixedSizeTypeForMap(*ty.MapTy.K)
		if err != nil {
			return "", "", err
		}
		valueType, valueTypeTag, err := parseTy(*ty.MapTy.V, contractName)
		wrappedType := wrapTlbTypeIntoTlbTag(valueType, valueTypeTag)

		tlbType.WriteString(fixedType)
		tlbType.WriteRune(',')
		tlbType.WriteString(wrappedType)
		tlbType.WriteRune(']')
	case "EnumRef":
		name := contractName + utils.ToCamelCase(ty.EnumRefTy.EnumName)
		tlbType.WriteString(name)
	case "StructRef":
		name := contractName + utils.ToCamelCase(ty.StructRefTy.StructName)
		tlbType.WriteString(name)
		err := writeGenericTypesPart(&tlbType, ty.StructRefTy.TypeArgs, contractName)
		if err != nil {
			return "", "", err
		}
	case "AliasRef":
		name := contractName + utils.ToCamelCase(ty.AliasRefTy.AliasName)
		tlbType.WriteString(name)
		err := writeGenericTypesPart(&tlbType, ty.AliasRefTy.TypeArgs, contractName)
		if err != nil {
			return "", "", err
		}
	case "genericT":
		name := ty.GenericTy.NameT
		tlbType.WriteString(utils.ToCamelCase(name))
	case "union":
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
	case "nullLiteral":
		tlbType.WriteString("tlb.NullLiteral")
	case "tupleAny", "callable":
		return "", "", fmt.Errorf("cannot convert type %v", ty.SumType)
	default:
		return "", "", fmt.Errorf("unknown ty type %q", ty.SumType)
	}

	return tlbType.String(), tlbTag.String(), nil
}

func getFixedSizeTypeForMap(ty tolkAbi.Ty) (string, error) {
	switch ty.SumType {
	case "intN":
		n := ty.NumberTy.N
		tlbIntType := fmt.Sprintf("tlb.Int%v", n)
		return tlbIntType, nil
	case "uintN":
		n := ty.NumberTy.N
		tlbIntType := fmt.Sprintf("tlb.Uint%v", n)
		return tlbIntType, nil
	case "varintN", "varuintN":
		tlbVarUIntType := fmt.Sprintf("tlb.VarUInteger%v", ty.NumberTy.N)
		return tlbVarUIntType, nil
	case "bitsN":
		tlbVarUIntType := fmt.Sprintf("tlb.Bits%v", ty.BitsTy.N)
		return tlbVarUIntType, nil
	case "bool":
		return "tlb.Uint1", nil
	case "address", "addressOpt", "addressExt", "addressAny":
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

func convertPrefixTagToInt(tag tolkAbi.Prefix) (*Tag, error) {
	if tag.PrefixLen == 0 {
		return nil, fmt.Errorf("prefix tag len must be > 0")
	}

	var val uint64
	var err error
	if len(tag.PrefixStr) == 1 {
		val, err = strconv.ParseUint(tag.PrefixStr, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		if tag.PrefixStr[1] == 'b' {
			val, err = strconv.ParseUint(tag.PrefixStr[2:], 2, 64)
			if err != nil {
				return nil, err
			}
		} else if tag.PrefixStr[1] == 'x' {
			val, err = strconv.ParseUint(tag.PrefixStr[2:], 16, 64)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("prefix tag must be either binary or hex format")
		}
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

func writeGenericTypesPart(builder *strings.Builder, typeArgs []tolkAbi.Ty, contractName string) error {
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
