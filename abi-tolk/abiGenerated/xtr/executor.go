// Code generated - DO NOT EDIT.

package abiXtr

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type XtrMaster interface {
	WithAccountId(accountID ton.AccountID) XtrMasterWithAccount
	GetUserLatestVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error)
	GetPaymentLatestVersion(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error)
}

type XtrMasterWithAccount interface {
	GetUserLatestVersion(ctx context.Context) (tlb.Uint32, error)
	GetPaymentLatestVersion(ctx context.Context) (tlb.Uint32, error)
}
