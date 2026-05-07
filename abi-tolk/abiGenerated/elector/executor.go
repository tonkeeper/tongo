// Code generated - DO NOT EDIT.

package abiElector

import (
	"context"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Executor interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type StorageExecutor interface {
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

type Elector interface {
	WithAccountId(accountID ton.AccountID) ElectorWithAccount
	GetActiveElectionId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error)
	GetParticipatesIn(ctx context.Context, reqAccountID ton.AccountID, validatorKey tlb.Uint256) (tlb.Int257, error)
	GetComputeReturnedStake(ctx context.Context, reqAccountID ton.AccountID, walletAddr tlb.Uint256) (tlb.Int257, error)
	GetParticipantListExtended(ctx context.Context, reqAccountID ton.AccountID) (ParticipantListExtended, error)
	AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ElectorStorage, error)
}

type ElectorWithAccount interface {
	GetActiveElectionId(ctx context.Context) (tlb.Int257, error)
	GetParticipatesIn(ctx context.Context, validatorKey tlb.Uint256) (tlb.Int257, error)
	GetComputeReturnedStake(ctx context.Context, walletAddr tlb.Uint256) (tlb.Int257, error)
	GetParticipantListExtended(ctx context.Context) (ParticipantListExtended, error)
	AccountState(ctx context.Context) (tlb.ShardAccount, ElectorStorage, error)
}
