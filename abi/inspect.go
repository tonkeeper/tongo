package abi

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	codePkg "github.com/tonkeeper/tongo/code"
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
	MethodsInspected   int
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

type libResolver interface {
	GetLibraries(ctx context.Context, libraryList []ton.Bits256) (map[ton.Bits256]*boc.Cell, error)
}

type contractInspector struct {
	knownMethods    []MethodDescription
	knownInterfaces []InterfaceDescription
	scanAllMethods  bool
	libResolver     libResolver
}

type InspectorOptions struct {
	additionalMethods []MethodDescription
	knownInterfaces   []InterfaceDescription
	scanAllMethods    bool
	libResolver       libResolver
}

type ContractInterface uint32

func (c ContractInterface) Implements(other ContractInterface) bool {
	if c == other {
		return true
	}
	return c.recursiveImplements(other)
}

type InvokeFn func(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (string, any, error)

// MethodDescription describes a particular method and provides a function to execute it.
type MethodDescription struct {
	Name string
	// InvokeFn executes this method on a contract and returns parsed execution results.
	InvokeFn InvokeFn
}

type knownContractDescription struct {
	contractInterfaces []ContractInterface
	getMethods         []InvokeFn
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

func InspectWithAdditionalInterfaces(list []InterfaceDescription) InspectorOption {
	return func(o *InspectorOptions) {
		o.knownInterfaces = list
	}
}
func InspectWithAllMethods() InspectorOption {
	return func(o *InspectorOptions) {
		o.scanAllMethods = true
	}
}

func InspectWithLibraryResolver(resolver libResolver) InspectorOption {
	return func(o *InspectorOptions) {
		o.libResolver = resolver
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
		scanAllMethods:  options.scanAllMethods,
		libResolver:     options.libResolver,
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
	info, err := GetCodeInfo(ctx, code, ci.libResolver)
	if err != nil {
		return nil, err
	}

	if contract, ok := knownContracts[info.Hash]; ok { //for known contracts we just need to run get methods
		desc.ContractInterfaces = contract.contractInterfaces
		for _, method := range contract.getMethods {
			desc.MethodsInspected += 1
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
		if !ci.scanAllMethods && !info.isMethodOkToTry(method.Name) {
			continue
		}
		desc.MethodsInspected += 1
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

type CodeInfo struct {
	Hash    ton.Bits256
	Methods map[int64]struct{}
}

func (i CodeInfo) isMethodOkToTry(name string) bool {
	if i.Methods == nil {
		return false
	}
	methodID := utils.MethodIdFromName(name)
	_, ok := i.Methods[int64(methodID)]
	return ok
}

func GetCodeInfo(ctx context.Context, code []byte, resolver libResolver) (*CodeInfo, error) {
	cells, err := boc.DeserializeBoc(code)
	if err != nil {
		return nil, err
	}
	if len(cells) == 0 {
		return nil, fmt.Errorf("failed to find a root cell")
	}
	root := cells[0]
	libHashes, err := codePkg.FindLibraries(root)
	if err != nil {
		return nil, fmt.Errorf("failed while looking for libraries inside cell: %w", err)
	}
	var libs map[ton.Bits256]*boc.Cell
	if len(libHashes) > 0 {
		if resolver == nil {
			return nil, fmt.Errorf("found libraries in cell, but no resolver provided")
		}
		libs, err = resolver.GetLibraries(ctx, libHashes)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch libraries: %w", err)
		}
	}
	h, err := root.Hash256()
	if err != nil {
		return nil, err
	}
	root.ResetCounters()
	if root.IsLibrary() {
		hash, err := root.GetLibraryHash()
		if err != nil {
			return nil, err
		}
		cell, ok := libs[ton.Bits256(hash)]
		if !ok {
			return nil, fmt.Errorf("library not found")
		}
		root = cell
	}
	c, err := root.NextRef()
	if err != nil {
		// we are OK, if there is no information about get methods
		return &CodeInfo{Hash: h}, nil
	}
	if c.IsLibrary() {
		hash, err := c.GetLibraryHash()
		if err != nil {
			return nil, err
		}
		cell, ok := libs[ton.Bits256(hash)]
		if !ok {
			return nil, fmt.Errorf("library not found")
		}
		c = cell
	}
	type GetMethods struct {
		Hashmap tlb.Hashmap[tlb.Uint19, boc.Cell]
	}
	var getMethods GetMethods
	decoder := tlb.NewDecoder().WithLibraryResolver(func(hash tlb.Bits256) (*boc.Cell, error) {
		if resolver == nil {
			return nil, fmt.Errorf("failed to fetch library: no resolver provided")
		}
		cell, ok := libs[ton.Bits256(hash)]
		if ok {
			return cell, nil
		}
		localLibs, err := resolver.GetLibraries(ctx, []ton.Bits256{ton.Bits256(hash)})
		if err != nil {
			return nil, err
		}
		if len(localLibs) == 0 {
			return nil, fmt.Errorf("library not found")
		}
		return localLibs[ton.Bits256(hash)], nil
	})

	err = decoder.Unmarshal(c, &getMethods)
	if err != nil {
		// we are OK, if there is no information about get methods
		return &CodeInfo{Hash: h}, nil
	}
	keys := getMethods.Hashmap.Keys()
	methods := make(map[int64]struct{}, len(keys))
	for _, key := range keys {
		methods[int64(key)] = struct{}{}
	}
	return &CodeInfo{Hash: h, Methods: methods}, nil
}
