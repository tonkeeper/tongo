// Code generated - DO NOT EDIT.

package abiCpmmV2

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type CpmmV2 interface {
	WithAccountId(accountID ton.AccountID) CpmmV2WithAccount
	GetPoolData(ctx context.Context, reqAccountID ton.AccountID) (DedustCpmmV2GetPoolData, error)
}

type CpmmV2WithAccount interface {
	GetPoolData(ctx context.Context) (DedustCpmmV2GetPoolData, error)
}
