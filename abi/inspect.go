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

type InterfaceIntrospection struct {
	// GetMethods contains successfully executed methods and results of executions.
	GetMethods map[string]MethodInvocation
}

type ContractIntrospection struct {
	// Interfaces is a list of interfaces implemented by a contract.
	Interfaces map[ContractInterface]InterfaceIntrospection
}

// InspectContract tries to execute all known get method on a contract and returns a summary about it.
// The contract is not provided directly,
// instead, it is expected that "executor" has been created with knowledge of the contract's code and data.
// Executor must be ready to execute multiple different methods and must not rely on a particular order of execution.
func InspectContract(ctx context.Context, code []byte, executor Executor, reqAccountID tongo.AccountID) (*ContractIntrospection, error) {
	if len(code) == 0 {
		return &ContractIntrospection{}, nil
	}
	hash, err := getCodeHash(code)
	if err != nil {
		return nil, err
	}
	implemented := make(map[ContractInterface]InterfaceIntrospection)
	for _, iface := range walletsByHashCode[hash.Hex()] {
		implemented[iface] = InterfaceIntrospection{
			GetMethods: map[string]MethodInvocation{},
		}
	}
	for _, method := range methodInvocationOrder {
		typeHint, result, err := method.InvokeFn(ctx, executor, reqAccountID)
		if err != nil {
			continue
		}
		for _, ifaceName := range method.ImplementedBy {
			intr, ok := implemented[ifaceName]
			if !ok {
				intr = InterfaceIntrospection{
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
	return &ContractIntrospection{
		Interfaces: implemented,
	}, nil
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
