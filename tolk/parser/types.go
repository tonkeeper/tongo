package parser

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/tonkeeper/tongo/utils"
)

type Kind struct {
	Kind string `json:"kind"`
}

type ABI struct {
	Namespace        string            `json:"namespace"`
	ContractName     string            `json:"contractName"`
	InheritsContract string            `json:"inheritsContract,omitempty"`
	Author           string            `json:"author,omitempty"`
	Version          string            `json:"version,omitempty"`
	Description      string            `json:"description,omitempty"`
	Declarations     []Declaration     `json:"declarations"`
	IncomingMessages []IncomingMessage `json:"incomingMessages"`
	IncomingExternal *IncomingExternal `json:"incomingExternal,omitempty"`
	OutgoingMessages []OutgoingMessage `json:"outgoingMessages"`
	EmittedMessages  []OutgoingMessage `json:"emittedEvents"`
	GetMethods       []GetMethod       `json:"getMethods"`
	ThrownErrors     []ThrownError     `json:"thrownErrors"`
	CompilerName     string            `json:"compilerName"`
	CompilerVersion  string            `json:"compilerVersion"`
	CodeBoc64        string            `json:"codeBoc64"`
	CodeHashes       []string          `json:"codeHashes,omitempty"`
}

func (a *ABI) GetGolangNamespace() string {
	return utils.ToCamelCase(a.Namespace)
}

func (a *ABI) GetGolangContractName() string {
	return a.GetGolangNamespace() + utils.ToCamelCase(a.ContractName)
}

type DeclarationKind string

const (
	DeclarationKindStruct DeclarationKind = "Struct"
	DeclarationKindAlias  DeclarationKind = "Alias"
	DeclarationKindEnum   DeclarationKind = "Enum"
)

type Declaration struct {
	SumType           DeclarationKind `json:"kind"`
	PayloadType       *string         `json:"payloadType,omitempty"` // todo: think abt naming
	StructDeclaration StructDeclaration
	AliasDeclaration  AliasDeclaration
	EnumDeclaration   EnumDeclaration
}

func (d *Declaration) UnmarshalJSON(b []byte) error {
	var r struct {
		Kind        string  `json:"kind"`
		PayloadType *string `json:"payloadType,omitempty"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	d.PayloadType = r.PayloadType
	switch r.Kind {
	case "Struct":
		d.SumType = DeclarationKindStruct
		if err := json.Unmarshal(b, &d.StructDeclaration); err != nil {
			return err
		}
	case "Alias":
		d.SumType = DeclarationKindAlias
		if err := json.Unmarshal(b, &d.AliasDeclaration); err != nil {
			return err
		}
	case "Enum":
		d.SumType = DeclarationKindEnum
		if err := json.Unmarshal(b, &d.EnumDeclaration); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown declaration type %q", d.SumType)
	}

	return nil
}

func (d Declaration) MarshalJSON() ([]byte, error) {
	var kind Kind
	kind.Kind = string(d.SumType)

	var payload []byte
	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}

	switch d.SumType {
	case DeclarationKindStruct:
		payload, err = json.Marshal(d.StructDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	case DeclarationKindAlias:
		payload, err = json.Marshal(d.AliasDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	case DeclarationKindEnum:
		payload, err = json.Marshal(d.EnumDeclaration)
		if err != nil {
			return nil, err
		}
		return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
	default:
		return nil, fmt.Errorf("unknown declaration type %q", d.SumType)
	}
}

type StructDeclaration struct {
	Name             string               `json:"name"`
	TypeParams       []string             `json:"typeParams,omitempty"`
	Prefix           *Prefix              `json:"prefix,omitempty"`
	Fields           []Field              `json:"fields"`
	CustomPackUnpack ABICustomSerializers `json:"customPackUnpack"`
}

type Prefix struct {
	PrefixStr string `json:"prefixStr"`
	PrefixLen int    `json:"prefixLen"`
}

type Field struct {
	Name string `json:"name"`
	//IsPayload    *bool         `json:"isPayload,omitempty"`
	Ty           Ty            `json:"ty"`
	DefaultValue *DefaultValue `json:"defaultValue,omitempty"`
	Description  string        `json:"description,omitempty"`
}

type DefaultValue struct {
	SumType         string `json:"kind"`
	IntDefaultValue struct {
		V string `json:"v"`
	}
	BoolDefaultValue struct {
		V bool `json:"v"`
	}
	SliceDefaultValue struct {
		Hex string `json:"hex"`
	}
	AddressDefaultValue struct {
		Address string `json:"addr"`
	}
	TensorDefaultValue struct {
		Items []DefaultValue `json:"items"`
	}
	NullDefaultValue struct{}
}

func (d *DefaultValue) UnmarshalJSON(b []byte) error {
	var kind Kind

	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}

	switch kind.Kind {
	case "int":
		d.SumType = "IntDefaultValue"
		if err := json.Unmarshal(b, &d.IntDefaultValue); err != nil {
			return err
		}
	case "bool":
		d.SumType = "BoolDefaultValue"
		if err := json.Unmarshal(b, &d.BoolDefaultValue); err != nil {
			return err
		}
	case "slice":
		d.SumType = "SliceDefaultValue"
		if err := json.Unmarshal(b, &d.SliceDefaultValue); err != nil {
			return err
		}
	case "address":
		d.SumType = "AddressDefaultValue"
		if err := json.Unmarshal(b, &d.AddressDefaultValue); err != nil {
			return err
		}
	case "tensor":
		d.SumType = "TensorDefaultValue"
		if err := json.Unmarshal(b, &d.TensorDefaultValue); err != nil {
			return err
		}
	case "null":
		d.SumType = "NullDefaultValue"
	default:
		return fmt.Errorf("unknown default value type %q", d.SumType)
	}

	return nil
}

func (d *DefaultValue) MarshalJSON() ([]byte, error) {
	var kind Kind
	var payload []byte
	var err error

	switch d.SumType {
	case "IntDefaultValue":
		kind.Kind = "int"
		payload, err = json.Marshal(d.IntDefaultValue)
		if err != nil {
			return nil, err
		}
	case "BoolDefaultValue":
		kind.Kind = "bool"
		payload, err = json.Marshal(d.BoolDefaultValue)
		if err != nil {
			return nil, err
		}
	case "SliceDefaultValue":
		kind.Kind = "slice"
		payload, err = json.Marshal(d.SliceDefaultValue)
		if err != nil {
			return nil, err
		}
	case "AddressDefaultValue":
		kind.Kind = "address"
		payload, err = json.Marshal(d.AddressDefaultValue)
		if err != nil {
			return nil, err
		}
	case "TensorDefaultValue":
		kind.Kind = "tensor"
		payload, err = json.Marshal(d.TensorDefaultValue)
		if err != nil {
			return nil, err
		}
	case "NullDefaultValue":
		kind.Kind = "null"
	default:
		return nil, fmt.Errorf("unknown default value type %q", d.SumType)
	}

	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}
	return utils.ConcatPrefixAndSuffixIfExists(prefix, payload), nil
}

type AliasDeclaration struct {
	Name             string               `json:"name"`
	TargetTy         Ty                   `json:"targetTy"`
	TypeParams       []string             `json:"typeParams,omitempty"`
	CustomPackUnpack ABICustomSerializers `json:"customPackUnpack"`
}

type EnumDeclaration struct {
	Name             string               `json:"name"`
	EncodedAs        Ty                   `json:"encodedAs"`
	Members          []EnumMember         `json:"members"`
	CustomPackUnpack ABICustomSerializers `json:"customPackUnpack"`
}

type ABICustomSerializers struct {
	PackToBuilder   bool `json:"packToBuilder"`
	UnpackFromSlice bool `json:"unpackFromSlice"`
}

type EnumMember struct {
	Name string `json:"name"`
	// Value is a number string
	Value string `json:"value"`
}

/*
export type Ty =
    | { kind: 'int' }
    | { kind: 'intN'; n: number }
    | { kind: 'uintN'; n: number }
    | { kind: 'varintN', n: number }
    | { kind: 'varuintN', n: number }
    | { kind: 'coins' }
    | { kind: 'bool' }
    | { kind: 'cell' }
    | { kind: 'builder' }
    | { kind: 'slice' }
    | { kind: 'string' }
    | { kind: 'remaining' }
    | { kind: 'address' }
    | { kind: 'addressOpt' }
    | { kind: 'addressExt' }
    | { kind: 'addressAny' }
    | { kind: 'bitsN'; n: number }
    | { kind: 'nullLiteral' }
    | { kind: 'callable' }
    | { kind: 'void' }
    | { kind: 'unknown' }
    | { kind: 'nullable'; inner: Ty; stackTypeId?: number; stackWidth?: number }
    | { kind: 'cellOf'; inner: Ty }
    | { kind: 'arrayOf'; inner: Ty }
    | { kind: 'lispListOf'; inner: Ty }
    | { kind: 'tensor'; items: Ty[] }
    | { kind: 'shapedTuple'; items: Ty[] }
    | { kind: 'mapKV'; k: Ty; v: Ty }
    | { kind: 'EnumRef'; enumName: string }
    | { kind: 'StructRef'; structName: string; typeArgs?: Ty[] }
    | { kind: 'AliasRef'; aliasName: string; typeArgs?: Ty[] }
    | { kind: 'genericT'; nameT: string }
    | { kind: 'union'; variants: UnionVariant[]; stackWidth: number }

*/

type TyKind = string

const (
	TyKindInt         TyKind = "int"
	TyKindIntN        TyKind = "intN"
	TyKindUintN       TyKind = "uintN"
	TyKindVarIntN     TyKind = "varintN"
	TyKindVarUintN    TyKind = "varuintN"
	TyKindCoins       TyKind = "coins"
	TyKindBool        TyKind = "bool"
	TyKindCell        TyKind = "cell"
	TyKindBuilder     TyKind = "builder"
	TyKindSlice       TyKind = "slice"
	TyKindString      TyKind = "string"
	TyKindRemaining   TyKind = "remaining"
	TyKindAddress     TyKind = "address"
	TyKindAddressOpt  TyKind = "addressOpt"
	TyKindAddressExt  TyKind = "addressExt"
	TyKindAddressAny  TyKind = "addressAny"
	TyKindBitsN       TyKind = "bitsN"
	TyKindNullLiteral TyKind = "nullLiteral"
	TyKindCallable    TyKind = "callable"
	TyKindVoid        TyKind = "void"
	TyKindUnknown     TyKind = "unknown"
	TyKindNullable    TyKind = "nullable"
	TyKindCellOf      TyKind = "cellOf"
	TyKindArrayOf     TyKind = "arrayOf"
	TyKindLispListOf  TyKind = "lispListOf"
	TyKindTensor      TyKind = "tensor"
	TyKindShapedTuple TyKind = "shapedTuple"
	TyKindMapKV       TyKind = "mapKV"
	TyKindEnumRef     TyKind = "EnumRef"
	TyKindStructRef   TyKind = "StructRef"
	TyKindAliasRef    TyKind = "AliasRef"
	TyKindGenericT    TyKind = "genericT"
	TyKindUnion       TyKind = "union"
)

type String struct{}

type Unknown struct{}

type LispListOf struct {
	Inner Ty `json:"inner"`
}

type Ty struct {
	SumType     string `json:"kind"`
	IntN        *IntN
	UintN       *UintN
	VarIntN     *VarIntN
	BitsN       *BitsN
	Nullable    *Nullable
	CellOf      *CellOf
	ArrayOf     *ArrayOf
	LispListOf  *LispListOf
	Tensor      *Tensor
	ShapedTuple *ShapedTuple
	MapKV       *MapKV
	EnumRef     *EnumRef
	StructRef   *StructRef
	AliasRef    *AliasRef
	Generic     *Generic
	Union       *Union
	VarUintN    *VarUintN
	// types without parameters
	//    Int         *Int
	//    Coins       *Coins
	//    Bool        *Bool
	//    Cell        *Cell
	//    Builder     *Builder
	//    Slice       *Slice
	//    String      *String
	//    Remaining   *Remaining
	//    Address     *Address
	//    AddressOpt  *AddressOpt
	//    AddressExt  *AddressExt
	//    AddressAny  *AddressAny
	//    NullLiteral *NullLiteral
	//    Callable    *Callable
	//    Void        *Void
	//    Unknown     *Unknown
}

/**
export function renderTy(ty: Ty): string {
    switch (ty.kind) {
        case 'int':         return `int`;
        case 'intN':        return `int${ty.n}`;
        case 'uintN':       return `uint${ty.n}`;
        case 'varintN':     return `varint${ty.n}`;
        case 'varuintN':    return `varuint${ty.n}`;
        case 'coins':       return `coins`;
        case 'bool':        return `bool`;
        case 'cell':        return `cell`;
        case 'builder':     return `builder`;
        case 'slice':       return `slice`;
        case 'string':      return `string`;
        case 'remaining':   return `RemainingBitsAndRefs`;
        case 'address':     return `address`;
        case 'addressOpt':  return `address?`;
        case 'addressExt':  return `ext_address`;
        case 'addressAny':  return `any_address`;
        case 'bitsN':       return `bits${ty.n}`;
        case 'nullLiteral': return `null`;
        case 'callable':    return `continuation`;
        case 'void':        return `void`;
        case 'unknown':     return `unknown`;
        case 'nullable':    return `${renderTy(ty.inner)}?`;
        case 'cellOf':      return `Cell<${renderTy(ty.inner)}>`;
        case 'arrayOf':     return `array<${renderTy(ty.inner)}>`;
        case 'lispListOf':  return `lisp_list<${renderTy(ty.inner)}>`;
        case 'tensor':      return `(${ty.items.map(renderTy).join(', ')})`;
        case 'shapedTuple': return `[${ty.items.map(renderTy).join(', ')}]`;
        case 'mapKV':       return `map<${renderTy(ty.k)}, ${renderTy(ty.v)}>`;
        case 'EnumRef':     return ty.enumName;
        case 'StructRef':   return ty.structName + (ty.typeArgs ? `<${ty.typeArgs.map(renderTy).join(', ')}>` : '');
        case 'AliasRef':    return ty.aliasName + (ty.typeArgs ? `<${ty.typeArgs.map(renderTy).join(', ')}>` : '');
        case 'genericT':    return ty.nameT;
        case 'union':       return ty.variants.map(v => renderTy(v.variantTy)).join(' | ');
    }
}
*/

func (ty Ty) String() string {
	switch ty.SumType {
	case TyKindInt:
		return "int"
	case TyKindIntN:
		return fmt.Sprintf("int%d", ty.IntN.N)
	case TyKindUintN:
		return fmt.Sprintf("uint%d", ty.UintN.N)
	case TyKindVarIntN:
		return fmt.Sprintf("varint%d", ty.VarIntN.N)
	case TyKindVarUintN:
		return fmt.Sprintf("varuint%d", ty.VarUintN.N)
	case TyKindCoins:
		return "coins"
	case TyKindBool:
		return "bool"
	case TyKindCell:
		return "cell"
	case TyKindBuilder:
		return "builder"
	case TyKindSlice:
		return "slice"
	case TyKindString:
		return "string"
	case TyKindRemaining:
		return "RemainingBitsAndRefs"
	case TyKindAddress:
		return "address"
	case TyKindAddressOpt:
		return "address?"
	case TyKindAddressExt:
		return "ext_address"
	case TyKindAddressAny:
		return "any_address"
	case TyKindBitsN:
		return fmt.Sprintf("bits%d", ty.BitsN.N)
	case TyKindNullLiteral:
		return "null"
	case TyKindCallable:
		return "continuation"
	case TyKindVoid:
		return "void"
	case TyKindUnknown:
		return "unknown"
	case TyKindNullable:
		return fmt.Sprintf("%s?", ty.Nullable.Inner.String())
	case TyKindCellOf:
		return fmt.Sprintf("Cell<%s>", ty.CellOf.Inner.String())
	case TyKindArrayOf:
		return fmt.Sprintf("array<%s>", ty.ArrayOf.Inner.String())
	case TyKindLispListOf:
		return fmt.Sprintf("lisp_list<%s>", ty.LispListOf.Inner.String())
	case TyKindTensor:
		items := make([]string, len(ty.Tensor.Items))
		for i, item := range ty.Tensor.Items {
			items[i] = item.String()
		}
		return fmt.Sprintf("(%s)", strings.Join(items, ", "))
	case TyKindShapedTuple:
		items := make([]string, len(ty.ShapedTuple.Items))
		for i, item := range ty.ShapedTuple.Items {
			items[i] = item.String()
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case TyKindMapKV:
		return fmt.Sprintf("map<%s, %s>", ty.MapKV.K.String(), ty.MapKV.V.String())
	case TyKindEnumRef:
		return ty.EnumRef.EnumName
	case TyKindStructRef:
		if len(ty.StructRef.TypeArgs) > 0 {
			args := make([]string, len(ty.StructRef.TypeArgs))
			for i, arg := range ty.StructRef.TypeArgs {
				args[i] = arg.String()
			}
			return fmt.Sprintf("%s<%s>", ty.StructRef.StructName, strings.Join(args, ", "))
		}
		return ty.StructRef.StructName
	case TyKindAliasRef:
		if len(ty.AliasRef.TypeArgs) > 0 {
			args := make([]string, len(ty.AliasRef.TypeArgs))
			for i, arg := range ty.AliasRef.TypeArgs {
				args[i] = arg.String()
			}
			return fmt.Sprintf("%s<%s>", ty.AliasRef.AliasName, strings.Join(args, ", "))
		}
		return ty.AliasRef.AliasName
	case TyKindGenericT:
		return ty.Generic.NameT
	case TyKindUnion:
		variants := make([]string, len(ty.Union.Variants))
		for i, v := range ty.Union.Variants {
			variants[i] = v.VariantTy.String()
		}
		return strings.Join(variants, " | ")
	default:
		return "unknown"
	}
}

type IntN struct {
	N int `json:"n"`
}

type UintN struct {
	N int `json:"n"`
}

type VarIntN struct {
	N int `json:"n"`
}

type VarUintN struct {
	N int `json:"n"`
}

type BitsN struct {
	N int `json:"n"`
}

type Nullable struct {
	Inner       Ty  `json:"inner"`
	StackTypeId int `json:"stackTypeId,omitempty"`
	StackWidth  int `json:"stackWidth,omitempty"`
}

type CellOf struct {
	Inner Ty `json:"inner"`
}

type ArrayOf struct {
	Inner Ty `json:"inner"`
}

type MapKV struct {
	K Ty `json:"k"`
	V Ty `json:"v"`
}

type EnumRef struct {
	EnumName string `json:"enumName"`
}

type StructRef struct {
	StructName string `json:"structName"`
	TypeArgs   []Ty   `json:"typeArgs,omitempty"`
}

type AliasRef struct {
	AliasName string `json:"aliasName"`
	TypeArgs  []Ty   `json:"typeArgs,omitempty"`
}

type Generic struct {
	NameT string `json:"nameT"`
}

type Tensor struct {
	Items []Ty `json:"items"`
}

type ShapedTuple struct {
	Items []Ty `json:"items"`
}

type Union struct {
	Variants   []UnionVariant `json:"variants"`
	StackWidth int            `json:"stackWidth"`
}

func (t *Ty) UnmarshalJSON(b []byte) error {
	var kind Kind
	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}
	switch kind.Kind {
	case TyKindIntN:
		t.SumType = TyKindIntN
		return json.Unmarshal(b, &t.IntN)
	case TyKindUintN:
		t.SumType = TyKindUintN
		return json.Unmarshal(b, &t.UintN)
	case TyKindVarIntN:
		t.SumType = TyKindVarIntN
		return json.Unmarshal(b, &t.VarIntN)
	case TyKindVarUintN:
		t.SumType = TyKindVarUintN
		return json.Unmarshal(b, &t.VarUintN)
	case TyKindBitsN:
		t.SumType = TyKindBitsN
		return json.Unmarshal(b, &t.BitsN)
	case TyKindNullable:
		t.SumType = TyKindNullable
		return json.Unmarshal(b, &t.Nullable)
	case TyKindCellOf:
		t.SumType = TyKindCellOf
		return json.Unmarshal(b, &t.CellOf)
	case TyKindArrayOf:
		t.SumType = TyKindArrayOf
		return json.Unmarshal(b, &t.ArrayOf)
	case TyKindLispListOf:
		t.SumType = TyKindLispListOf
		return json.Unmarshal(b, &t.LispListOf)
	case TyKindTensor:
		t.SumType = TyKindTensor
		return json.Unmarshal(b, &t.Tensor)
	case TyKindShapedTuple:
		t.SumType = TyKindShapedTuple
		return json.Unmarshal(b, &t.ShapedTuple)
	case TyKindMapKV:
		t.SumType = TyKindMapKV
		return json.Unmarshal(b, &t.MapKV)
	case TyKindEnumRef:
		t.SumType = TyKindEnumRef
		return json.Unmarshal(b, &t.EnumRef)
	case TyKindStructRef:
		t.SumType = TyKindStructRef
		return json.Unmarshal(b, &t.StructRef)
	case TyKindAliasRef:
		t.SumType = TyKindAliasRef
		return json.Unmarshal(b, &t.AliasRef)
	case TyKindGenericT:
		t.SumType = TyKindGenericT
		return json.Unmarshal(b, &t.Generic)
	case TyKindUnion:
		t.SumType = TyKindUnion
		return json.Unmarshal(b, &t.Union)
	case TyKindInt, TyKindCoins, TyKindBool, TyKindCell, TyKindBuilder, TyKindSlice, TyKindString, TyKindRemaining, TyKindAddress:
		t.SumType = kind.Kind
	case TyKindAddressOpt, TyKindAddressExt, TyKindAddressAny, TyKindNullLiteral, TyKindCallable, TyKindVoid, TyKindUnknown:
		t.SumType = kind.Kind
		return nil
	default:
		return fmt.Errorf("unknown type kind %q", kind.Kind)
	}

	return nil
}

func (t *Ty) MarshalJSON() ([]byte, error) {
	var payload []byte
	var err error
	switch t.SumType {
	case TyKindIntN:
		payload, err = json.Marshal(t.IntN)
	case TyKindUintN:
		payload, err = json.Marshal(t.UintN)
	case TyKindVarIntN:
		payload, err = json.Marshal(t.VarIntN)
	case TyKindVarUintN:
		payload, err = json.Marshal(t.VarUintN)
	case TyKindBitsN:
		payload, err = json.Marshal(t.BitsN)
	case TyKindNullable:
		payload, err = json.Marshal(t.Nullable)
	case TyKindCellOf:
		payload, err = json.Marshal(t.CellOf)
	case TyKindArrayOf:
		payload, err = json.Marshal(t.ArrayOf)
	case TyKindLispListOf:
		payload, err = json.Marshal(t.LispListOf)
	case TyKindTensor:
		payload, err = json.Marshal(t.Tensor)
	case TyKindShapedTuple:
		payload, err = json.Marshal(t.ShapedTuple)
	case TyKindMapKV:
		payload, err = json.Marshal(t.MapKV)
	case TyKindEnumRef:
		payload, err = json.Marshal(t.EnumRef)
	case TyKindStructRef:
		payload, err = json.Marshal(t.StructRef)
	case TyKindAliasRef:
		payload, err = json.Marshal(t.AliasRef)
	case TyKindGenericT:
		payload, err = json.Marshal(t.Generic)
	case TyKindUnion:
		payload, err = json.Marshal(t.Union)
	default:
	}
	if err != nil {
		return nil, err
	}
	return utils.ConcatPrefixAndSuffixIfExists([]byte(fmt.Sprintf(`{"kind": "%s"}`, t.SumType)), payload), nil
}

func (t *Ty) GetFixedSize() (int, bool) {
	switch t.SumType {
	case "IntN":
		return t.IntN.N, true
	case "UintN":
		return t.UintN.N, true
	case "BitsN":
		return t.BitsN.N, true
	case "Bool":
		return 1, true
	case "Address":
		return 267, true
	default:
		return 0, false
	}
}

type UnionVariant struct {
	PrefixStr        string `json:"prefixStr"`
	PrefixLen        int    `json:"prefixLen"`
	PrefixEatInPlace bool   `json:"prefixEatInPlace,omitempty"`
	VariantTy        Ty     `json:"variantTy"`
}

type IncomingMessage struct {
	BodyTy            Ty       `json:"bodyTy"`
	MinimalMsgValue   *big.Int `json:"minimalMsgValue,omitempty"`
	Description       string   `json:"description,omitempty"`
	PreferredSendMode int16    `json:"preferredSendMode,omitempty"`
}

func (m *IncomingMessage) GetMsgName() (string, error) {
	return getMsgName(m.BodyTy)
}

type IncomingExternal struct {
	BodyTy      Ty     `json:"bodyTy"`
	Description string `json:"description,omitempty"`
}

func (m *IncomingExternal) GetMsgName() (string, error) {
	return getMsgName(m.BodyTy)
}

type OutgoingMessage struct {
	BodyTy      Ty     `json:"bodyTy"`
	Description string `json:"description,omitempty"`
}

func (m *OutgoingMessage) GetMsgName() (string, error) {
	return getMsgName(m.BodyTy)
}

func getMsgName(ty Ty) (string, error) {
	switch ty.SumType {
	case "StructRef":
		return ty.StructRef.StructName, nil
	case "AliasRef":
		return ty.AliasRef.AliasName, nil
	default:
		return "", fmt.Errorf("cannot get name for %q body", ty.SumType)
	}
}

type GetMethod struct {
	TvmMethodID int         `json:"tvmMethodId"`
	Name        string      `json:"name"`
	Parameters  []Parameter `json:"parameters"`
	ReturnTy    Ty          `json:"returnTy"`
	Description string      `json:"description,omitempty"`
}

func (g GetMethod) GolangFunctionName() string {
	return utils.ToCamelCase(g.Name)
}

func (g GetMethod) FullResultName(contractName string) string {
	res := ""
	if contractName != "" {
		res = contractName + "_"
	}
	res += utils.ToCamelCase(g.Name)

	return res + "Result"
}

func (g GetMethod) UsedByIntrospection() bool {
	return len(g.Parameters) == 0
}

type Parameter struct {
	Name string `json:"name"`
	Ty   Ty     `json:"ty"`
}

type ThrownError struct {
	Name    string `json:"constName"`
	ErrCode int    `json:"errCode"`
}

func MustParseABI(data []byte) ABI {
	var abi ABI
	if err := json.Unmarshal(data, &abi); err != nil {
		panic(err)
	}
	return abi
}
