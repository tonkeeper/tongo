// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type ClmmPool interface {
	WithAccountId(accountID ton.AccountID) ClmmPoolWithAccount
	GetPool(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*ClmmPoolStorage], error)
	GetRanges(ctx context.Context, reqAccountID ton.AccountID) (tlb.Hashmap[tlb.Int16, tlb.InternalAddress], error)
}

type ClmmPoolWithAccount interface {
	GetPool(ctx context.Context) (tlb.RefT[*ClmmPoolStorage], error)
	GetRanges(ctx context.Context) (tlb.Hashmap[tlb.Int16, tlb.InternalAddress], error)
}

type ClmmRange interface {
	WithAccountId(accountID ton.AccountID) ClmmRangeWithAccount
	GetRange(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*ClmmRangeStorage], error)
}

type ClmmRangeWithAccount interface {
	GetRange(ctx context.Context) (tlb.RefT[*ClmmRangeStorage], error)
}

type CpmmPosition interface {
	WithAccountId(accountID ton.AccountID) CpmmPositionWithAccount
	GetPosition(ctx context.Context, reqAccountID ton.AccountID) (tlb.RefT[*CpmmPositionStorage], error)
	GetLocked(ctx context.Context, reqAccountID ton.AccountID) (bool, error)
}

type CpmmPositionWithAccount interface {
	GetPosition(ctx context.Context) (tlb.RefT[*CpmmPositionStorage], error)
	GetLocked(ctx context.Context) (bool, error)
}
