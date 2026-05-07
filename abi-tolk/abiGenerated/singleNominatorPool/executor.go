// Code generated - DO NOT EDIT.

package abiSingleNominatorPool

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type SingleNominatorPool interface {
	WithAccountId(accountID ton.AccountID) SingleNominatorPoolWithAccount
	GetRoles(ctx context.Context, reqAccountID ton.AccountID) (Roles, error)
	GetPoolData(ctx context.Context, reqAccountID ton.AccountID) (PoolData, error)
}

type SingleNominatorPoolWithAccount interface {
	GetRoles(ctx context.Context) (Roles, error)
	GetPoolData(ctx context.Context) (PoolData, error)
}
