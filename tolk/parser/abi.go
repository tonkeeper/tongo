package parser

// ContractABI is a final result — the ABI of a TON smart contract.
//
// Partially, its properties may be specified by a user manually:
// > contract MyName {
// >    author: "Dima"
// >    incomingMessages: SomeUnion
// > }
//
// Partially, its properties are automatically calculated by the compiler:
// - outgoingMessages, via the calls to `createMessage`
// - thrownErrors, via `throw` and `assert` statements
// - getMethods, essentially `get fun`
//
// Doc comments from Tolk code come as descriptions:
// > /// Desc of struct
// > struct SomeOutgoingMessage {
// >     /// Desc of field
// >     field: int32
// > }
//
// While collecting messages/getters/etc., the compiler gathers unique_types
// and emits them separately. They are referred to as `body_ty_idx`, `return_ty_idx`, etc.
type ContractABI struct {
	ABISchemaVersion string `json:"abi_schema_version,omitempty"`
	ContractName     string `json:"contract_name"`
	Author           string `json:"author,omitempty"`
	Version          string `json:"version,omitempty"`
	Description      string `json:"description,omitempty"`

	UniqueTypes          []Ty                     `json:"unique_types"`
	StructInstantiations []ABIStructInstantiation `json:"struct_instantiations"`
	AliasInstantiations  []ABIAliasInstantiation  `json:"alias_instantiations"`
	Declarations         []ABIDeclaration         `json:"declarations"`

	Storage          ABIStorage           `json:"storage"`
	IncomingMessages []ABIInternalMessage `json:"incoming_messages"`
	IncomingExternal []ABIExternalMessage `json:"incoming_external"`
	OutgoingMessages []ABIOutgoingMessage `json:"outgoing_messages"`
	EmittedEvents    []ABIOutgoingMessage `json:"emitted_events"`
	GetMethods       []ABIGetMethod       `json:"get_methods"`
	ThrownErrors     []ABIThrownError     `json:"thrown_errors"`

	CompilerName    string `json:"compiler_name"`
	CompilerVersion string `json:"compiler_version"`
}

//func (a *ABI) GetGolangNamespace() string {
//	return utils.ToCamelCase(a.Namespace)
//}
//
//func (a *ABI) GetGolangContractName() string {
//	return a.GetGolangNamespace() + utils.ToCamelCase(a.ContractName)
//}

// ABIStorage defines shape of a storage.
// Most often, it's a regular struct, serializable into a cell.
// In the case of NFT, when a storage changes its shape (several fields appear after deployment),
// the "initial storage" can also be expressed: it's called "storage at deployment".
// The storage is used to visualize current contract state and to calculate its address.
// Storage descriptions live on corresponding declarations
type ABIStorage struct {
	StorageTyIdx             *int `json:"storage_ty_idx,omitempty"`
	StorageAtDeploymentTyIdx *int `json:"storage_at_deployment_ty_idx,omitempty"`
}

// ABIInternalMessage is "an incoming internal message" (handled by `onInternalMessage`).
// In practice, a user describes each message as a struct:
// > struct (0x12345678) Increment { ... }
// > struct (0x23456789) Reset { ... }
// Then, ABI of a contract will contain those two messages:
// * body_ty = { kind: 'StructRef', struct_name: 'Increment' }
// * body_ty = { kind: 'StructRef', struct_name: 'Reset' }
// Theoretically, body_ty can be something else: e.g., instantiation `Transfer<ForwardPayload>`.
// Message descriptions live on corresponding declarations.
type ABIInternalMessage struct {
	BodyTyIdx int `json:"body_ty_idx"`
}

// ABIExternalMessage is "an incoming external message" (handled by `onExternalMessage`).
// It's either a 'slice' or some struct.
// Message descriptions live on corresponding declarations.
type ABIExternalMessage struct {
	BodyTyIdx int `json:"body_ty_idx"`
}

// ABIOutgoingMessage is "an outgoing internal/external message".
// In Tolk code, those are calls to `createMessage`.
// Message descriptions live on corresponding declarations.
type ABIOutgoingMessage struct {
	BodyTyIdx int `json:"body_ty_idx"`
}

// ABIGetMethod is a "get method" (aka "contract getter").
// In Tolk code, getters are created with `get fun`.
// Example:
// > get fun calcData(owner: address): SomeStruct { ... }
// It has one parameter with ty = { kind: 'address' },
// its return_ty is { kind: 'StructRef', struct_name: 'SomeStruct' }.
// Note, that getters are called off-chain — via the TVM stack, not via serialization.
// (For instance, they can return 'int' or 'slice', although they are not serializable)
type ABIGetMethod struct {
	TVMMethodID int                     `json:"tvm_method_id"`
	Name        string                  `json:"name"`
	Parameters  []ABIGetMethodParameter `json:"parameters"`
	ReturnTyIdx int                     `json:"return_ty_idx"`
	Description string                  `json:"description,omitempty"`
}

//
//func (g GetMethod) GolangFunctionName() string {
//	return utils.ToCamelCase(g.Name)
//}
//
//func (g GetMethod) FullResultName(contractName string) string {
//	res := ""
//	if contractName != "" {
//		res = contractName + "_"
//	}
//	res += utils.ToCamelCase(g.Name)
//
//	return res + "Result"
//}
//
//func (g GetMethod) UsedByIntrospection() bool {
//	return len(g.Parameters) == 0
//}

type ABIGetMethodParameter struct {
	Name        string `json:"name"`
	TyIdx       int    `json:"ty_idx"`
	Description string `json:"description,omitempty"`
	// todo: default value
}

type ABIThrownErrorKind string

const (
	ABIThrownError_PlainInt   ABIThrownErrorKind = "plain_int"
	ABIThrownError_Constant   ABIThrownErrorKind = "constant"
	ABIThrownError_EnumMember ABIThrownErrorKind = "enum_member"
)

// ABIThrownError is an errCode fired by `throw` or `assert` in Tolk code.
// Example:
// > assert (valid) throw ErrCodes.NoAccess;
// Then { kind: 'enum_member', name: 'ErrCodes.NoAccess', err_code: 123 } exists.
// Description is `///` comment in Tolk code above a constant or an enum member.
type ABIThrownError struct {
	Kind        ABIThrownErrorKind `json:"kind"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	ErrCode     int                `json:"err_code"`
}
