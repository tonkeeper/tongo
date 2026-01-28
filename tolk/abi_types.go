package tolk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
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

type Declaration struct {
	SumType           string  `json:"kind"`
	PayloadType       *string `json:"payloadType,omitempty"` // todo: think abt naming
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

	d.SumType = r.Kind
	d.PayloadType = r.PayloadType
	switch d.SumType {
	case "Struct":
		if err := json.Unmarshal(b, &d.StructDeclaration); err != nil {
			return err
		}
	case "Alias":
		if err := json.Unmarshal(b, &d.AliasDeclaration); err != nil {
			return err
		}
	case "Enum":
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
	kind.Kind = d.SumType

	var payload []byte
	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}

	switch d.SumType {
	case "Struct":
		payload, err = json.Marshal(d.StructDeclaration)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "Alias":
		payload, err = json.Marshal(d.AliasDeclaration)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "Enum":
		payload, err = json.Marshal(d.EnumDeclaration)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	default:
		return nil, fmt.Errorf("unknown declaration type %q", d.SumType)
	}
}

type StructDeclaration struct {
	Name       string   `json:"name"`
	TypeParams []string `json:"typeParams,omitempty"`
	Prefix     *Prefix  `json:"prefix,omitempty"`
	Fields     []Field  `json:"fields"`
}

type Prefix struct {
	PrefixStr string `json:"prefixStr"`
	PrefixLen int    `json:"prefixLen"`
}

type Field struct {
	Name         string        `json:"name"`
	IsPayload    *bool         `json:"isPayload,omitempty"`
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
	return concatPrefixAndPayload(prefix, payload), nil
}

func (d *DefaultValue) UnmarshalDefaultValue(v TolkValue) (bool, error) {
	vType := v.GetType()
	switch d.SumType {
	case "IntDefaultValue":
		val, err := binDecHexToUint(d.IntDefaultValue.V)
		if err != nil {
			return false, err
		}
		err = v.SetValue(val)
		if err != nil {
			return false, err
		}
	case "BoolDefaultValue":
		err := v.SetValue(d.BoolDefaultValue.V)
		if err != nil {
			return false, err
		}
	case "SliceDefaultValue":
		if vType.SumType != "Remaining" && vType.SumType != "Slice" {
			return false, fmt.Errorf("default type %v and field type %v mismatch", d.SumType, vType.SumType)
		}
		val, err := hex.DecodeString(d.SliceDefaultValue.Hex)
		if err != nil {
			return false, err
		}
		bs := boc.NewBitString(hex.DecodedLen(len(val)))
		err = bs.WriteBytes(val)
		if err != nil {
			return false, err
		}
		anySlice := Any(*boc.NewCellWithBits(bs))
		err = v.SetValue(anySlice)
		if err != nil {
			return false, err
		}
	case "AddressDefaultValue":
		if vType.SumType != "Address" {
			return false, fmt.Errorf("default type %v and field type %v mismatch", d.SumType, vType.SumType)
		}
		accountID, err := ton.ParseAccountID(d.AddressDefaultValue.Address)
		if err != nil {
			return false, err
		}
		err = v.SetValue(InternalAddress{
			Workchain: int8(accountID.Workchain),
			Address:   accountID.Address,
		})
		if err != nil {
			return false, err
		}
	case "TensorDefaultValue":
		tensor := make([]Value, len(d.TensorDefaultValue.Items))
		for i, item := range d.TensorDefaultValue.Items {
			val := Value{}
			_, err := item.UnmarshalDefaultValue(&val)
			if err != nil {
				return false, err
			}
			tensor[i] = val
		}
		err := v.SetValue(tensor)
		if err != nil {
			return false, err
		}
	case "NullDefaultValue":
		return false, nil
	default:
		return false, fmt.Errorf("unknown default value type %q", d.SumType)
	}

	return true, nil
}

type AliasDeclaration struct {
	Name                  string   `json:"name"`
	TargetTy              Ty       `json:"targetTy"`
	TypeParams            []string `json:"typeParams,omitempty"`
	CustomPackToBuilder   bool     `json:"customPackToBuilder,omitempty"`
	CustomUnpackFromSlice bool     `json:"customUnpackFromSlice,omitempty"`
}

type EnumDeclaration struct {
	Name      string       `json:"name"`
	EncodedAs Ty           `json:"encodedAs"`
	Members   []EnumMember `json:"members"`
}

type EnumMember struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Ty struct {
	SumType     string `json:"kind"`
	Int         *Int
	IntN        *IntN
	UintN       *UintN
	VarIntN     *VarIntN
	VarUintN    *VarUintN
	BitsN       *BitsN
	Coins       *Coins
	Bool        *Bool
	Cell        *Cell
	Slice       *Slice
	Builder     *Builder
	Callable    *Callable
	Remaining   *Remaining
	Address     *Address
	AddressOpt  *AddressOpt
	AddressExt  *AddressExt
	AddressAny  *AddressAny
	Nullable    *Nullable
	CellOf      *CellOf
	Tensor      *Tensor
	TupleWith   *TupleWith
	Map         *Map
	EnumRef     *EnumRef
	AliasRef    *AliasRef
	StructRef   *StructRef
	Generic     *Generic
	Union       *Union
	TupleAny    *TupleAny
	NullLiteral *NullLiteral
	Void        *Void
}

func (t *Ty) UnmarshalJSON(b []byte) error {
	var kind Kind
	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}
	switch kind.Kind {
	case "intN":
		t.SumType = "IntN"
		if err := json.Unmarshal(b, &t.IntN); err != nil {
			return err
		}
	case "uintN":
		t.SumType = "UintN"
		if err := json.Unmarshal(b, &t.UintN); err != nil {
			return err
		}
	case "varintN":
		t.SumType = "VarIntN"
		if err := json.Unmarshal(b, &t.VarIntN); err != nil {
			return err
		}
	case "varuintN":
		t.SumType = "VarUintN"
		if err := json.Unmarshal(b, &t.VarUintN); err != nil {
			return err
		}
	case "bitsN":
		t.SumType = "BitsN"
		if err := json.Unmarshal(b, &t.BitsN); err != nil {
			return err
		}
	case "nullable":
		t.SumType = "Nullable"
		if err := json.Unmarshal(b, &t.Nullable); err != nil {
			return err
		}
	case "cellOf":
		t.SumType = "CellOf"
		if err := json.Unmarshal(b, &t.CellOf); err != nil {
			return err
		}
	case "tensor":
		t.SumType = "Tensor"
		if err := json.Unmarshal(b, &t.Tensor); err != nil {
			return err
		}
	case "tupleWith":
		t.SumType = "TupleWith"
		if err := json.Unmarshal(b, &t.TupleWith); err != nil {
			return err
		}
	case "mapKV":
		t.SumType = "Map"
		if err := json.Unmarshal(b, &t.Map); err != nil {
			return err
		}
	case "EnumRef":
		t.SumType = "EnumRef"
		if err := json.Unmarshal(b, &t.EnumRef); err != nil {
			return err
		}
	case "StructRef":
		t.SumType = "StructRef"
		if err := json.Unmarshal(b, &t.StructRef); err != nil {
			return err
		}
	case "AliasRef":
		t.SumType = "AliasRef"
		if err := json.Unmarshal(b, &t.AliasRef); err != nil {
			return err
		}
	case "genericT":
		t.SumType = "Generic"
		if err := json.Unmarshal(b, &t.Generic); err != nil {
			return err
		}
	case "union":
		t.SumType = "Union"
		if err := json.Unmarshal(b, &t.Union); err != nil {
			return err
		}
	case "int":
		t.SumType = "Int"
		t.IntN = &IntN{}
	case "coins":
		t.SumType = "Coins"
		t.Coins = &Coins{}
	case "bool":
		t.SumType = "Bool"
		t.Bool = &Bool{}
	case "cell":
		t.SumType = "Cell"
		t.Cell = &Cell{}
	case "slice":
		t.SumType = "Slice"
		t.Slice = &Slice{}
	case "builder":
		t.SumType = "Builder"
		t.Builder = &Builder{}
	case "remaining":
		t.SumType = "Remaining"
		t.Remaining = &Remaining{}
	case "address":
		t.SumType = "Address"
		t.Address = &Address{}
	case "addressOpt":
		t.SumType = "AddressOpt"
		t.AddressOpt = &AddressOpt{}
	case "addressExt":
		t.SumType = "AddressExt"
		t.AddressExt = &AddressExt{}
	case "addressAny":
		t.SumType = "AddressAny"
		t.AddressAny = &AddressAny{}
	case "tupleAny":
		t.SumType = "TupleAny"
		t.TupleAny = &TupleAny{}
	case "nullLiteral":
		t.SumType = "NullLiteral"
		t.NullLiteral = &NullLiteral{}
	case "callable":
		t.SumType = "Callable"
		t.Callable = &Callable{}
	case "void":
		t.SumType = "Void"
		t.Void = &Void{}
	default:
		return fmt.Errorf("unknown ty type %q", kind.Kind)
	}

	return nil
}

func (t *Ty) MarshalJSON() ([]byte, error) {
	var kind Kind
	var prefix []byte
	var payload []byte
	var err error

	switch t.SumType {
	case "IntN":
		kind.Kind = "intN"
		payload, err = json.Marshal(t.IntN)
		if err != nil {
			return nil, err
		}
	case "UintN":
		kind.Kind = "uintN"
		payload, err = json.Marshal(t.UintN)
		if err != nil {
			return nil, err
		}
	case "VarIntN":
		kind.Kind = "varintN"
		payload, err = json.Marshal(t.VarIntN)
		if err != nil {
			return nil, err
		}
	case "VarUintN":
		kind.Kind = "varuintN"
		payload, err = json.Marshal(t.VarUintN)
		if err != nil {
			return nil, err
		}
	case "BitsN":
		kind.Kind = "bitsN"
		payload, err = json.Marshal(t.BitsN)
		if err != nil {
			return nil, err
		}
	case "Nullable":
		kind.Kind = "nullable"
		payload, err = json.Marshal(t.Nullable)
		if err != nil {
			return nil, err
		}
	case "CellOf":
		kind.Kind = "cellOf"
		payload, err = json.Marshal(t.CellOf)
		if err != nil {
			return nil, err
		}
	case "Tensor":
		kind.Kind = "tensor"
		payload, err = json.Marshal(t.Tensor)
		if err != nil {
			return nil, err
		}
	case "TupleWith":
		kind.Kind = "tupleWith"
		payload, err = json.Marshal(t.TupleWith)
		if err != nil {
			return nil, err
		}
	case "Map":
		kind.Kind = "mapKV"
		payload, err = json.Marshal(t.Map)
		if err != nil {
			return nil, err
		}
	case "EnumRef":
		kind.Kind = "EnumRef"
		payload, err = json.Marshal(t.EnumRef)
		if err != nil {
			return nil, err
		}
	case "StructRef":
		kind.Kind = "StructRef"
		payload, err = json.Marshal(t.StructRef)
		if err != nil {
			return nil, err
		}
	case "AliasRef":
		kind.Kind = "AliasRef"
		payload, err = json.Marshal(t.AliasRef)
		if err != nil {
			return nil, err
		}
	case "Generic":
		kind.Kind = "genericT"
		payload, err = json.Marshal(t.Generic)
		if err != nil {
			return nil, err
		}
	case "Union":
		kind.Kind = "union"
		payload, err = json.Marshal(t.Union)
		if err != nil {
			return nil, err
		}
	case "Int":
		kind.Kind = "int"
	case "Coins":
		kind.Kind = "coins"
	case "Bool":
		kind.Kind = "bool"
	case "Cell":
		kind.Kind = "cell"
	case "Slice":
		kind.Kind = "slice"
	case "Builder":
		kind.Kind = "builder"
	case "Remaining":
		kind.Kind = "remaining"
	case "Address":
		kind.Kind = "address"
	case "AddressOpt":
		kind.Kind = "addressOpt"
	case "AddressExt":
		kind.Kind = "addressExt"
	case "AddressAny":
		kind.Kind = "addressAny"
	case "TupleAny":
		kind.Kind = "tupleAny"
	case "NullLiteral":
		kind.Kind = "nullLiteral"
	case "Callable":
		kind.Kind = "callable"
	case "Void":
		kind.Kind = "void"
	default:
		return nil, fmt.Errorf("unknown ty type %q", t.SumType)
	}

	prefix, err = json.Marshal(kind)
	if err != nil {
		return nil, err
	}
	return concatPrefixAndPayload(prefix, payload), nil
}

func (t *Ty) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	v.SetType(*t)
	switch t.SumType {
	case "Int":
		return t.Int.UnmarshalTolk(cell, v, decoder)
	case "NullLiteral":
		return t.NullLiteral.UnmarshalTolk(cell, v, decoder)
	case "Void":
		return t.Void.UnmarshalTolk(cell, v, decoder)
	case "IntN":
		return t.IntN.UnmarshalTolk(cell, v, decoder)
	case "UintN":
		return t.UintN.UnmarshalTolk(cell, v, decoder)
	case "VarIntN":
		return t.VarIntN.UnmarshalTolk(cell, v, decoder)
	case "VarUintN":
		return t.VarUintN.UnmarshalTolk(cell, v, decoder)
	case "BitsN":
		return t.BitsN.UnmarshalTolk(cell, v, decoder)
	case "Coins":
		return t.Coins.UnmarshalTolk(cell, v, decoder)
	case "Bool":
		return t.Bool.UnmarshalTolk(cell, v, decoder)
	case "Cell":
		return t.Cell.UnmarshalTolk(cell, v, decoder)
	case "Slice":
		return t.Slice.UnmarshalTolk(cell, v, decoder)
	case "Builder":
		return t.Builder.UnmarshalTolk(cell, v, decoder)
	case "Callable":
		return t.Callable.UnmarshalTolk(cell, v, decoder)
	case "Remaining":
		return t.Remaining.UnmarshalTolk(cell, v, decoder)
	case "Address":
		return t.Address.UnmarshalTolk(cell, v, decoder)
	case "AddressOpt":
		return t.AddressOpt.UnmarshalTolk(cell, v, decoder)
	case "AddressExt":
		return t.AddressExt.UnmarshalTolk(cell, v, decoder)
	case "AddressAny":
		return t.AddressAny.UnmarshalTolk(cell, v, decoder)
	case "Nullable":
		return t.Nullable.UnmarshalTolk(cell, v, decoder)
	case "CellOf":
		return t.CellOf.UnmarshalTolk(cell, v, decoder)
	case "Tensor":
		return t.Tensor.UnmarshalTolk(cell, v, decoder)
	case "TupleWith":
		return t.TupleWith.UnmarshalTolk(cell, v, decoder)
	case "TupleAny":
		return t.TupleAny.UnmarshalTolk(cell, v, decoder)
	case "Map":
		return t.Map.UnmarshalTolk(cell, v, decoder)
	case "EnumRef":
		return t.EnumRef.UnmarshalTolk(cell, v, decoder)
	case "StructRef":
		return t.StructRef.UnmarshalTolk(cell, v, decoder)
	case "AliasRef":
		return t.AliasRef.UnmarshalTolk(cell, v, decoder)
	case "Generic":
		return t.Generic.UnmarshalTolk(cell, v, decoder)
	case "Union":
		return t.Union.UnmarshalTolk(cell, v, decoder)
	default:
		return fmt.Errorf("unknown t type %q", t.SumType)
	}
}

func (t *Ty) GetComparableType() (TolkComparableType, bool) {
	var tolkType TolkType
	switch t.SumType {
	case "Int":
		tolkType = t.Int
	case "NullLiteral":
		tolkType = t.NullLiteral
	case "Void":
		tolkType = t.Void
	case "IntN":
		tolkType = t.IntN
	case "UintN":
		tolkType = t.UintN
	case "VarIntN":
		tolkType = t.VarIntN
	case "VarUintN":
		tolkType = t.VarUintN
	case "BitsN":
		tolkType = t.BitsN
	case "Coins":
		tolkType = t.Coins
	case "Bool":
		tolkType = t.Bool
	case "Cell":
		tolkType = t.Cell
	case "Slice":
		tolkType = t.Slice
	case "Builder":
		tolkType = t.Builder
	case "Callable":
		tolkType = t.Callable
	case "Remaining":
		tolkType = t.Remaining
	case "Address":
		tolkType = t.Address
	case "AddressOpt":
		tolkType = t.AddressOpt
	case "AddressExt":
		tolkType = t.AddressExt
	case "AddressAny":
		tolkType = t.AddressAny
	case "Nullable":
		tolkType = t.Nullable
	case "CellOf":
		tolkType = t.CellOf
	case "Tensor":
		tolkType = t.Tensor
	case "TupleWith":
		tolkType = t.TupleWith
	case "TupleAny":
		tolkType = t.TupleAny
	case "Map":
		tolkType = t.Map
	case "EnumRef":
		tolkType = t.EnumRef
	case "StructRef":
		tolkType = t.StructRef
	case "AliasRef":
		tolkType = t.AliasRef
	case "Generic":
		tolkType = t.Generic
	case "Union":
		tolkType = t.Union
	default:
		return nil, false
	}

	var comparableType TolkComparableType
	var ok bool
	comparableType, ok = tolkType.(TolkComparableType)
	return comparableType, ok
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

func concatPrefixAndPayload(prefix, payload []byte) []byte {
	if len(payload) == 0 {
		return prefix
	}
	prefix = prefix[:len(prefix)-1] // remove '}'
	payload[0] = ','                // replace '{' with ','
	result := make([]byte, 0, len(prefix)+len(payload))
	result = append(result, prefix...)
	result = append(result, payload...)
	return result
}
