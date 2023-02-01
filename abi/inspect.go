package abi

import (
	"context"

	"github.com/tonkeeper/tongo"
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
	// For an interface to be implemented means there is at least one successfully executed get method.
	Interfaces map[ContractInterface]InterfaceIntrospection
}

// InspectContract tries to execute all known get method on a contract and returns a summary about it.
// The contract is not provided directly,
// instead, it is expected that "executor" has been created with knowledge of the contract's code and data.
// Executor must be ready to execute multiple different methods and must not rely on a particular order of execution.
func InspectContract(ctx context.Context, executor Executor, reqAccountID tongo.AccountID) (*ContractIntrospection, error) {
	results := make(map[ContractInterface]InterfaceIntrospection)
	for _, method := range methodInvocationOrder {
		typeHint, result, err := method.InvokeFn(ctx, executor, reqAccountID)
		if err != nil {
			continue
		}
		for _, ifaceName := range method.ImplementedBy {
			intr, ok := results[ifaceName]
			if !ok {
				intr = InterfaceIntrospection{
					GetMethods: map[string]MethodInvocation{},
				}
				results[ifaceName] = intr
			}
			intr.GetMethods[method.Name] = MethodInvocation{
				Result:   result,
				TypeHint: typeHint,
			}
		}
	}
	return &ContractIntrospection{
		Interfaces: results,
	}, nil
}
