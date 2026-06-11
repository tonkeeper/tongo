// Code generated - DO NOT EDIT.

package abiStonkspump

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type StonkspumpVirtual interface {
	WithAccountId(accountID ton.AccountID) StonkspumpVirtualWithAccount
	GetBondingData(ctx context.Context, reqAccountID ton.AccountID) (BondingData, error)
}

type StonkspumpVirtualWithAccount interface {
	GetBondingData(ctx context.Context) (BondingData, error)
}
