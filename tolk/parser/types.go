package parser

import (
	"math/big"
)

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

type Ty struct {
	SumType     TyKind `json:"kind"`
	IntN        *IntN
	UintN       *UintN
	VarIntN     *VarIntN
	VarUintN    *VarUintN
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
	GenericT    *GenericT
	Union       *Union
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

type ABIDeclarationKind string

const (
	DeclarationKindStruct ABIDeclarationKind = "struct"
	DeclarationKindAlias  ABIDeclarationKind = "alias"
	DeclarationKindEnum   ABIDeclarationKind = "enum"
)

type ABIDeclaration struct {
	SumType           ABIDeclarationKind `json:"kind"`
	StructDeclaration ABIStruct
	AliasDeclaration  ABIAlias
	EnumDeclaration   ABIEnum
}

// ABIStruct represents a Tolk struct.
// Examples:
//
// > struct Point { x: int, y: int }
// A simple struct, not serializable (because 'int', not 'int8' or similar)
//
// > struct Wrapper<TItem> { item: TItem }
// A generic struct (has type_params), fields[0].ty = { kind: 'genericT', name_t: 'TItem' }
//
// > struct (0x12345678) Increment { ... }
// Has a serialization prefix: prefix_num = 0x12345678, prefix_len = 32
//
// A field can have `@abi.clientType(<type>)` annotation to override how it's rendered in explorers/UI.
// Then the field has `client_ty_idx` set. It's used for wrappers (serialization), but not for stack.
type ABIStruct struct {
	Name             string               `json:"name"`
	TyIdx            int                  `json:"ty_idx"`
	TypeParams       []string             `json:"type_params,omitempty"`
	Prefix           *Prefix              `json:"prefix,omitempty"`
	Fields           []Field              `json:"fields"`
	CustomPackUnpack ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
	Description      string               `json:"description,omitempty"`
}

type Prefix struct {
	PrefixNum int `json:"prefix_num"`
	PrefixLen int `json:"prefix_len"`
}

type Field struct {
	Name        string `json:"name"`
	TyIdx       int    `json:"ty_idx"`
	ClientTyIdx *int   `json:"client_ty_idx,omitempty"`
	Description string `json:"description,omitempty"`
	// DefaultValue *ABIConstExpression // not used
}

// ABIAlias represents a Tolk type alias.
// Examples:
//
// > type UserId = int32
// A simple alias, target_ty = { kind: 'intN', n: 32 }
//
// > type Maybe<T> = MaybeNothing | MaybeJust<T>
// A generic alias (has type_params), target_ty = { kind: 'union', variants: ... }
//
// An alias is serialized as its target, unless it has custom serializers in Tolk code.
type ABIAlias struct {
	Name             string               `json:"name"`
	TyIdx            int                  `json:"ty_idx"`
	TargetTyIdx      int                  `json:"target_ty_idx"`
	TypeParams       []string             `json:"type_params,omitempty"`
	CustomPackUnpack ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
	Description      string               `json:"description,omitempty"`
}

// ABIEnum represents a Tolk enum.
// Examples:
//
// > enum Color { Red, Green, Blue }
// Has 3 members (values '0', '1', '2'), encoded as 'uint2' (auto-calculated by the compiler).
//
// > enum Mode: int8 { User = 0, Admin = 127 }
// Has 2 members, encoded as 'int8' (specified manually).
type ABIEnum struct {
	Name             string               `json:"name"`
	TyIdx            int                  `json:"ty_idx"`
	EncodedAsTyIdx   int                  `json:"encoded_as_ty_idx"`
	Members          []ABIEnumMember      `json:"members"`
	CustomPackUnpack ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
	Description      string               `json:"description,omitempty"`
}

type ABIEnumMember struct {
	Name  string  `json:"name"`
	Value big.Int `json:"value"`
}

// ABIStructInstantiation is a resolved runtime shape for a generic struct.
// For example, we have `struct Wrapper<T> { item: T }`, and somewhere `Wrapper<Point>` is used.
// Then `Wrapper<Point>` is present in `unique_types` as `StructRef "Wrapper" type_args=Point",
// and its instantiated fields types = ["Point"] exist here.
// Without monomorphization, client-side cannot reconstruct stack layout in case of generic nullables/unions.
type ABIStructInstantiation struct {
	TyIdx                  int                  `json:"ty_idx"`
	StructName             string               `json:"struct_name"`
	MonomorphicFieldsTyIdx []int                `json:"monomorphic_fields_ty_idx"`
	CustomPackUnpack       ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
}

// ABICustomSerializers is present in a struct/alias if that type has custom serializers in Tolk code:
//
// fun SomeType.packToBuilder(self, mutate b: builder) { ... }
// fun SomeType.unpackFromSlice(mutate s: slice): SomeAlias { ... }
//
// Body of these functions is not a part of ABI (it's arbitrary Tolk code).
// To make serialization work (e.g., in TypeScript wrappers),
// one should provide equivalent implementations in a language ABI is applied to.
type ABICustomSerializers struct {
	PackToBuilder   bool `json:"pack_to_builder"`
	UnpackFromSlice bool `json:"unpack_from_slice"`
}

// ABIAliasInstantiation is a resolved runtime shape for a generic alias.
// For example, we have `type Maybe<T> = None | Just<T>`, and somewhere `Maybe<int>` is used.
// Then `Maybe<int>` is present in `unique_types` as `AliasRef "Maybe" type_args=int",
// and its instantiated target = "None | Just<int>" exists here.
// Without monomorphization, client-side cannot reconstruct stack layout in case of generic nullables/unions.
type ABIAliasInstantiation struct {
	TyIdx                  int                  `json:"ty_idx"`
	AliasName              string               `json:"alias_name"`
	MonomorphicTargetTyIdx int                  `json:"monomorphic_target_ty_idx"`
	CustomPackUnpack       ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
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
	InnerTyIdx  int `json:"inner_ty_idx"`
	StackTypeId int `json:"stack_type_id,omitempty"`
	StackWidth  int `json:"stack_width,omitempty"`
}

type CellOf struct {
	InnerTyIdx int `json:"inner_ty_idx"`
}

type ArrayOf struct {
	InnerTyIdx int `json:"inner_ty_idx"`
}

type LispListOf struct {
	InnerTyIdx int `json:"inner_ty_idx"`
}

type Tensor struct {
	ItemsTyIdx []int `json:"items_ty_idx"`
}

type ShapedTuple struct {
	ItemsTyIdx []int `json:"items_ty_idx"`
}

type MapKV struct {
	KeyTyIdx   int `json:"key_ty_idx"`
	ValueTyIdx int `json:"value_ty_idx"`
}

type EnumRef struct {
	EnumName string `json:"enum_name"`
}

type StructRef struct {
	StructName    string `json:"struct_name"`
	TypeArgsTyIdx []int  `json:"type_args_ty_idx,omitempty"`
}

type AliasRef struct {
	AliasName     string `json:"alias_name"`
	TypeArgsTyIdx []int  `json:"type_args_ty_idx,omitempty"`
}

type GenericT struct {
	NameT string `json:"name_t"`
}

type Union struct {
	Variants   []UnionVariant `json:"variants"`
	StackWidth int            `json:"stack_width,omitempty"`
}

type UnionVariant struct {
	VariantTyIdx     int  `json:"variant_ty_idx"`
	PrefixNum        int  `json:"prefix_num"`
	PrefixLen        int  `json:"prefix_len"`
	IsPrefixImplicit bool `json:"is_prefix_implicit,omitempty"`
	StackTypeId      int  `json:"stack_type_id,omitempty"`
	StackWidth       int  `json:"stack_width,omitempty"`
}

func (t *Ty) GetFixedSize() (int, bool) {
	switch t.SumType {
	case TyKindIntN:
		return t.IntN.N, true
	case TyKindUintN:
		return t.UintN.N, true
	case TyKindBitsN:
		return t.BitsN.N, true
	case TyKindBool:
		return 1, true
	case TyKindAddress:
		return 267, true
	default:
		return 0, false
	}
}
