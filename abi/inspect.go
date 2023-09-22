package abi

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/utils"
)

type MethodInvocation struct {
	Result   any
	TypeHint string
}

type ContractDescription struct {
	// Interfaces is a list of interfaces implemented by a contract.
	ContractInterfaces []ContractInterface
	GetMethods         []MethodInvocation
}

func (d ContractDescription) hasAllResults(results []string) bool {
	for _, r := range results {
		found := false
		for _, m := range d.GetMethods {
			if m.TypeHint == r {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type contractInspector struct {
	knownMethods    []MethodDescription
	knownInterfaces []InterfaceDescription
}

type InspectorOptions struct {
	additionalMethods []MethodDescription
	knownInterfaces   []InterfaceDescription
}

type ContractInterface string

type InvokeFn func(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (string, any, error)

// MethodDescription describes a particular method and provides a function to execute it.
type MethodDescription struct {
	Name string
	// InvokeFn executes this method on a contract and returns parsed execution results.
	InvokeFn InvokeFn
}

type InterfaceDescription struct {
	Name    ContractInterface
	Results []string
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
		knownMethods:    append(methodInvocationOrder, options.additionalMethods...),
		knownInterfaces: append(contractInterfacesOrder, options.knownInterfaces...),
	}
}

// InspectContract tries to execute all known get method on a contract and returns a summary.
// The contract is not provided directly,
// instead, it is expected that "executor" has been created with knowledge of the contract's code and data.
// Executor must be ready to execute multiple different methods and must not rely on a particular order of execution.
func (ci contractInspector) InspectContract(ctx context.Context, code []byte, executor Executor, reqAccountID ton.AccountID) (*ContractDescription, error) {
	if len(code) == 0 {
		return &ContractDescription{}, nil
	}
	desc := ContractDescription{}
	info, err := getCodeInfo(code)
	if err != nil {
		return nil, err
	}

	if contract, ok := knownContracts[info.hash]; ok { //for known contracts we just need to run get methods
		desc.ContractInterfaces = contract.contractInterfaces
		for _, method := range contract.getMethods {
			typeHint, result, err := method(ctx, executor, reqAccountID)
			if err != nil {
				return &desc, nil
			}
			desc.GetMethods = append(desc.GetMethods, MethodInvocation{
				Result:   result,
				TypeHint: typeHint,
			})
		}
		return &desc, nil
	}

	for _, method := range ci.knownMethods {
		// let's avoid running get methods that we know don't exist
		if !info.isMethodOkToTry(method.Name) {
			continue
		}
		typeHint, result, err := method.InvokeFn(ctx, executor, reqAccountID)
		if err != nil {
			continue
		}
		desc.GetMethods = append(desc.GetMethods, MethodInvocation{
			Result:   result,
			TypeHint: typeHint,
		})
	}
	for _, iface := range ci.knownInterfaces {
		if desc.hasAllResults(iface.Results) {
			desc.ContractInterfaces = append(desc.ContractInterfaces, iface.Name)
		}
	}

	return &desc, nil
}

type codeInfo struct {
	hash    ton.Bits256
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
	var hash ton.Bits256
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
