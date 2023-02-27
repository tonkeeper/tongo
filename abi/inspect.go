package abi

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/utils"
)

type MethodInvocation struct {
	Result   any
	TypeHint string
}

type InterfaceDescription struct {
	// GetMethods contains successfully executed methods and results of executions.
	GetMethods map[string]MethodInvocation
}

type ContractDescription struct {
	// Interfaces is a list of interfaces implemented by a contract.
	Interfaces map[ContractInterface]InterfaceDescription
}

type contractInspector struct {
	knownMethods []MethodDescription
}

type InspectorOptions struct {
	additionalMethods []MethodDescription
}

type InspectorOption func(o *InspectorOptions)

func InspectWithAdditionalMethods(list []MethodDescription) InspectorOption {
	return func(o *InspectorOptions) {
		o.additionalMethods = list
	}
}

func NewContractInspector(opts ...InspectorOption) *contractInspector {
	options := &InspectorOptions{}
	for _, o := range opts {
		o(options)
	}
	return &contractInspector{
		knownMethods: append(methodInvocationOrder, options.additionalMethods...),
	}
}

// InspectContract tries to execute all known get method on a contract and returns a summary.
// The contract is not provided directly,
// instead, it is expected that "executor" has been created with knowledge of the contract's code and data.
// Executor must be ready to execute multiple different methods and must not rely on a particular order of execution.
func (ci contractInspector) InspectContract(ctx context.Context, code []byte, executor Executor, reqAccountID tongo.AccountID) (*ContractDescription, error) {
	if len(code) == 0 {
		return &ContractDescription{}, nil
	}
	info, err := getCodeInfo(code)
	if err != nil {
		return nil, err
	}
	var restricted map[ContractInterface]struct{}
	implemented := make(map[ContractInterface]InterfaceDescription)

	if wallets, ok := walletsByHashCode[info.hash.Hex()]; ok {
		// ok, this is a wallet,
		restricted = make(map[ContractInterface]struct{}, len(wallets))
		for _, iface := range wallets {
			restricted[iface] = struct{}{}
			implemented[iface] = InterfaceDescription{
				GetMethods: map[string]MethodInvocation{},
			}
		}
	}

	for _, method := range ci.knownMethods {
		// if this is a wallet, we execute only wallet methods
		if !isMethodAllowed(method.ImplementedBy, restricted) {
			continue
		}
		// let's avoid running get methods that we know don't exist
		if !info.isMethodOkToTry(method.Name) {
			continue
		}
		typeHint, result, err := method.InvokeFn(ctx, executor, reqAccountID)
		if err != nil {
			continue
		}
		for _, ifaceName := range method.implementedInterfaces(typeHint) {
			if !isInterfaceAllowed(ifaceName, restricted) {
				continue
			}
			intr, ok := implemented[ifaceName]
			if !ok {
				intr = InterfaceDescription{
					GetMethods: map[string]MethodInvocation{},
				}
				implemented[ifaceName] = intr
			}
			intr.GetMethods[method.Name] = MethodInvocation{
				Result:   result,
				TypeHint: typeHint,
			}
		}
	}
	for ifaceName := range implemented {
		if IsWallet(ifaceName) {
			implemented[Wallet] = InterfaceDescription{
				GetMethods: map[string]MethodInvocation{},
			}
			break
		}
	}
	return &ContractDescription{
		Interfaces: implemented,
	}, nil
}

func (m MethodDescription) implementedInterfaces(typeHint string) []ContractInterface {
	if m.ImplementedByFn != nil {
		// implementedByFn is optional,
		// if it's defined,
		// we use typeHint to get the exact contract interface.
		iface := m.ImplementedByFn(typeHint)
		if len(iface) > 0 {
			return []ContractInterface{iface}
		}
		return nil
	}
	return m.ImplementedBy
}

func isInterfaceAllowed(name ContractInterface, restricted map[ContractInterface]struct{}) bool {
	if restricted == nil {
		return true
	}
	_, ok := restricted[name]
	return ok
}

func isMethodAllowed(implementedBy []ContractInterface, restricted map[ContractInterface]struct{}) bool {
	if restricted == nil {
		return true
	}
	for _, iface := range implementedBy {
		if _, ok := restricted[iface]; ok {
			return true
		}
	}
	return false
}

func (ci ContractDescription) ImplementedInterfaces() []ContractInterface {
	results := make([]ContractInterface, 0, len(ci.Interfaces))
	for ifaceName := range ci.Interfaces {
		results = append(results, ifaceName)
	}
	return results
}

type codeInfo struct {
	hash    tongo.Bits256
	methods map[int64]struct{}
}

func (i codeInfo) isMethodOkToTry(name string) bool {
	if i.methods == nil {
		return true
	}
	methodID := utils.MethodIdFromName(name)
	_, ok := i.methods[int64(methodID)]
	return ok
}

func getCodeInfo(code []byte) (*codeInfo, error) {
	cells, err := boc.DeserializeBoc(code)
	if err != nil {
		return nil, err
	}
	if len(cells) == 0 {
		return nil, fmt.Errorf("failed to find a root cell")
	}
	h, err := cells[0].Hash()
	if err != nil {
		return nil, err
	}
	var hash tongo.Bits256
	if err = hash.FromBytes(h); err != nil {
		return nil, err
	}
	cells[0].ResetCounters()
	c, err := cells[0].NextRef()
	if err != nil {
		// we are OK, if there is no information about get methods
		return &codeInfo{hash: hash}, nil
	}
	type GetMethods struct {
		Hashmap tlb.Hashmap[tlb.Uint19, boc.Cell]
	}
	var getMethods GetMethods
	err = tlb.Unmarshal(c, &getMethods)
	if err != nil {
		// we are OK, if there is no information about get methods
		return &codeInfo{hash: hash}, nil
	}
	keys := getMethods.Hashmap.Keys()
	methods := make(map[int64]struct{}, len(keys))
	for _, key := range keys {
		methods[int64(key)] = struct{}{}
	}
	return &codeInfo{hash: hash, methods: methods}, nil
}
