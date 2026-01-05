package tolkAbi

import (
	"encoding/json"
	"fmt"
	"math/big"

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
	SumType           string `json:"kind"`
	StructDeclaration StructDeclaration
	AliasDeclaration  AliasDeclaration
	EnumDeclaration   EnumDeclaration
}

func (d *Declaration) UnmarshalJSON(b []byte) error {
	var kind Kind

	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}

	d.SumType = kind.Kind
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
}

func (d *DefaultValue) UnmarshalJSON(b []byte) error {
	var kind Kind

	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}

	d.SumType = kind.Kind
	switch d.SumType {
	case "int":
		if err := json.Unmarshal(b, &d.IntDefaultValue); err != nil {
			return err
		}
	case "bool":
		if err := json.Unmarshal(b, &d.BoolDefaultValue); err != nil {
			return err
		}
	case "address":
		if err := json.Unmarshal(b, &d.AddressDefaultValue); err != nil {
			return err
		}
	case "tensor":
		if err := json.Unmarshal(b, &d.TensorDefaultValue); err != nil {
			return err
		}
	case "null": // do nothing since null value have no additional fields
	default:
		return fmt.Errorf("unknown default value type %q", d.SumType)
	}

	return nil
}

func (d *DefaultValue) MarshalJSON() ([]byte, error) {
	var kind Kind
	kind.Kind = d.SumType

	var payload []byte
	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}

	switch d.SumType {
	case "int":
		payload, err = json.Marshal(d.IntDefaultValue)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "bool":
		payload, err = json.Marshal(d.BoolDefaultValue)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "address":
		payload, err = json.Marshal(d.AddressDefaultValue)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "tensor":
		payload, err = json.Marshal(d.TensorDefaultValue)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "null":
		return prefix, nil
	default:
		return nil, fmt.Errorf("unknown default value type %q", d.SumType)
	}
}

type AliasDeclaration struct {
	Name                  string   `json:"name"`
	TargetTy              *Ty      `json:"targetTy"`
	TypeParams            []string `json:"typeParams,omitempty"`
	CustomPackToBuilder   bool     `json:"customPackToBuilder,omitempty"`
	CustomUnpackFromSlice bool     `json:"customUnpackFromSlice,omitempty"`
}

type EnumDeclaration struct {
	Name      string       `json:"name"`
	EncodedAs *Ty          `json:"encodedAs"`
	Members   []EnumMember `json:"members"`
}

type EnumMember struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Ty struct {
	SumType  string `json:"kind"`
	NumberTy struct {
		N int `json:"n"`
	}
	BitsTy struct {
		N int `json:"n"`
	}
	NullableTy struct {
		Inner *Ty `json:"inner"`
	}
	CellOf struct {
		Inner *Ty `json:"inner"`
	}
	TensorTy struct {
		Items []Ty `json:"items"`
	}
	TupleWithTy struct {
		Items []Ty `json:"items"`
	}
	MapTy struct {
		K *Ty `json:"k"`
		V *Ty `json:"v"`
	}
	EnumRefTy struct {
		EnumName string `json:"enumName"`
	}
	StructRefTy struct {
		StructName string `json:"structName"`
		TypeArgs   []Ty   `json:"typeArgs,omitempty"`
	}
	AliasRefTy struct {
		AliasName string `json:"aliasName"`
		TypeArgs  []Ty   `json:"typeArgs,omitempty"`
	}
	GenericTy struct {
		NameT string `json:"nameT"`
	}
	Union struct {
		Variants []UnionVariant `json:"variants"`
	}
}

func (t *Ty) UnmarshalJSON(b []byte) error {
	var kind Kind
	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}

	t.SumType = kind.Kind
	switch t.SumType {
	case "intN", "uintN", "varintN", "varuintN":
		if err := json.Unmarshal(b, &t.NumberTy); err != nil {
			return err
		}
	case "bitsN":
		if err := json.Unmarshal(b, &t.BitsTy); err != nil {
			return err
		}
	case "nullable":
		if err := json.Unmarshal(b, &t.NullableTy); err != nil {
			return err
		}
	case "cellOf":
		if err := json.Unmarshal(b, &t.CellOf); err != nil {
			return err
		}
	case "tensor":
		if err := json.Unmarshal(b, &t.TensorTy); err != nil {
			return err
		}
	case "tupleWith":
		if err := json.Unmarshal(b, &t.TupleWithTy); err != nil {
			return err
		}
	case "mapKV":
		if err := json.Unmarshal(b, &t.MapTy); err != nil {
			return err
		}
	case "EnumRef":
		if err := json.Unmarshal(b, &t.EnumRefTy); err != nil {
			return err
		}
	case "StructRef":
		if err := json.Unmarshal(b, &t.StructRefTy); err != nil {
			return err
		}
	case "AliasRef":
		if err := json.Unmarshal(b, &t.AliasRefTy); err != nil {
			return err
		}
	case "genericT":
		if err := json.Unmarshal(b, &t.GenericTy); err != nil {
			return err
		}
	case "union":
		if err := json.Unmarshal(b, &t.Union); err != nil {
			return err
		}
	case "int", "coins", "bool", "cell", "slice", "builder", "remaining", "address", "addressOpt", "addressExt",
		"addressAny", "tupleAny", "nullLiteral", "callable", "void":
		return nil
	default:
		return fmt.Errorf("unknown ty type %q", t.SumType)
	}

	return nil
}

func (t *Ty) MarshalJSON() ([]byte, error) {
	var kind Kind
	kind.Kind = t.SumType

	var payload []byte
	prefix, err := json.Marshal(kind)
	if err != nil {
		return nil, err
	}

	switch t.SumType {
	case "intN", "uintN", "varintN", "varuintN":
		payload, err = json.Marshal(t.NumberTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "bitsN":
		payload, err = json.Marshal(t.BitsTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "nullable":
		payload, err = json.Marshal(t.NullableTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "cellOf":
		payload, err = json.Marshal(t.CellOf)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "tensor":
		payload, err = json.Marshal(t.TensorTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "tupleWith":
		payload, err = json.Marshal(t.TupleWithTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "mapKV":
		payload, err = json.Marshal(t.MapTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "EnumRef":
		payload, err = json.Marshal(t.EnumRefTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "StructRef":
		payload, err = json.Marshal(t.StructRefTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "AliasRef":
		payload, err = json.Marshal(t.AliasRefTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "genericT":
		payload, err = json.Marshal(t.GenericTy)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "union":
		payload, err = json.Marshal(t.Union)
		if err != nil {
			return nil, err
		}
		return concatPrefixAndPayload(prefix, payload), nil
	case "int", "coins", "bool", "cell", "slice", "builder", "remaining", "address", "addressOpt", "addressExt",
		"addressAny", "tupleAny", "nullLiteral", "callable", "void":
		return prefix, nil
	default:
		return nil, fmt.Errorf("unknown ty type %q", t.SumType)
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
		return ty.StructRefTy.StructName, nil
	case "AliasRef":
		return ty.AliasRefTy.AliasName, nil
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
	prefix = prefix[:len(prefix)-1] // remove '}'
	payload[0] = ','                // replace '{' with ','
	result := make([]byte, 0, len(prefix)+len(payload))
	result = append(result, prefix...)
	result = append(result, payload...)
	return result
}
