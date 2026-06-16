// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type EscrowFactory interface {
	WithAccountId(accountID ton.AccountID) EscrowFactoryWithAccount
	GetEscrowData(ctx context.Context, reqAccountID ton.AccountID) (EscrowData, error)
	GetItemAddress(ctx context.Context, reqAccountID ton.AccountID, quoteId tlb.Int257) (tlb.InternalAddress, error)
	GetVaultAddress(ctx context.Context, reqAccountID ton.AccountID, jettonWallet tlb.InternalAddress, owner tlb.InternalAddress) (tlb.InternalAddress, error)
	GetVersion(ctx context.Context, reqAccountID ton.AccountID) (GetVersionResult, error)
}

type EscrowFactoryWithAccount interface {
	GetEscrowData(ctx context.Context) (EscrowData, error)
	GetItemAddress(ctx context.Context, quoteId tlb.Int257) (tlb.InternalAddress, error)
	GetVaultAddress(ctx context.Context, jettonWallet tlb.InternalAddress, owner tlb.InternalAddress) (tlb.InternalAddress, error)
	GetVersion(ctx context.Context) (GetVersionResult, error)
}

type EscrowPosition interface {
	WithAccountId(accountID ton.AccountID) EscrowPositionWithAccount
	GetOrderData(ctx context.Context, reqAccountID ton.AccountID) (GetOrderDataResult, error)
}

type EscrowPositionWithAccount interface {
	GetOrderData(ctx context.Context) (GetOrderDataResult, error)
}

type EscrowVault interface {
	WithAccountId(accountID ton.AccountID) EscrowVaultWithAccount
	GetVaultData(ctx context.Context, reqAccountID ton.AccountID) (GetVaultDataResult, error)
}

type EscrowVaultWithAccount interface {
	GetVaultData(ctx context.Context) (GetVaultDataResult, error)
}
