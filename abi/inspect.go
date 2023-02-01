package abi

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
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
	hash, err := getCodeHash(code)
	if err != nil {
		return nil, err
	}
	var restricted map[ContractInterface]struct{}
	implemented := make(map[ContractInterface]InterfaceDescription)

	if wallets, ok := walletsByHashCode[hash.Hex()]; ok {
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
		typeHint, result, err := method.InvokeFn(ctx, executor, reqAccountID)
		if err != nil {
			continue
		}
		for _, ifaceName := range method.ImplementedBy {
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

func getCodeHash(code []byte) (*tongo.Bits256, error) {
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
	return &hash, nil
}

func (ci ContractDescription) ImplementedInterfaces() []ContractInterface {
	results := make([]ContractInterface, 0, len(ci.Interfaces))
	for ifaceName := range ci.Interfaces {
		results = append(results, ifaceName)
	}
	return results
}
