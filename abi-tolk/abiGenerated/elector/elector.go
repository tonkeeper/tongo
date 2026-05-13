// Code generated - DO NOT EDIT.

package abiElector

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixNewStake uint64 = 0x4e73744b

type NewStake struct {
	QueryId         tlb.Uint64             // uint64
	ValidatorPubkey tlb.Bits256            // bits256
	StakeAt         tlb.Uint32             // uint32
	MaxFactor       tlb.Uint32             // uint32
	AdnlAddr        tlb.Bits256            // bits256
	Signature       tlb.RefT[*tlb.Bits512] // Cell<bits512>
}

const PrefixNewStakeConfirmation uint64 = 0xf374484c

type NewStakeConfirmation struct {
	QueryId tlb.Uint64 // uint64
	Comment tlb.Uint32 // uint32
}

const PrefixRecoverStakeRequest uint64 = 0x47657424

type RecoverStakeRequest struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixRecoverStakeResponse uint64 = 0xf96f7324

type RecoverStakeResponse struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixUpgradeCode uint64 = 0x4e436f64

type UpgradeCode struct {
	QueryId tlb.Uint64 // uint64
	Code    boc.Cell   // cell
}

const PrefixUpgradeCodeResponse uint64 = 0xce436f64

type UpgradeCodeResponse struct {
	QueryId tlb.Uint64 // uint64
	Op      tlb.Uint32 // uint32
}

const PrefixConfigAccepted uint64 = 0xee764f4b

type ConfigAccepted struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixConfigRejected uint64 = 0xee764f6f

type ConfigRejected struct {
	QueryId tlb.Uint64 // uint64
}

const PrefixRegisterComplaint uint64 = 0x52674370

type RegisterComplaint struct {
	QueryId    tlb.Uint64         // uint64
	ElectionId tlb.Uint32         // uint32
	Complaint  ValidatorComplaint // ValidatorComplaint
}

const PrefixVoteComplaint uint64 = 0x56744370

type VoteComplaint struct {
	QueryId       tlb.Uint64  // uint64
	Signature     tlb.Bits512 // bits512
	SignTag       tlb.Uint32  // uint32
	ValidatorIdx  tlb.Uint16  // uint16
	ElectionId    tlb.Uint32  // uint32
	ComplaintHash tlb.Uint256 // uint256
}

const PrefixReturnStake uint64 = 0xee6f454c

type ReturnStake struct {
	QueryId tlb.Uint64 // uint64
	Reason  tlb.Uint32 // uint32
}

const PrefixErrorResponse uint64 = 0xfffffffe

type ErrorResponse struct {
	QueryId tlb.Uint64 // uint64
	Op      tlb.Uint32 // uint32
}

const PrefixComplaintResponse uint64 = 0xf2676350

type ComplaintResponse struct {
	QueryId tlb.Uint64 // uint64
	Op      tlb.Uint32 // uint32
}

const PrefixVoteComplaintResponse uint64 = 0xd6745240

type VoteComplaintResponse struct {
	QueryId tlb.Uint64 // uint64
	Op      tlb.Uint32 // uint32
}
type ValidatorComplaint struct {
	Tag               tlb.Uint8   // uint8
	ValidatorPubkey   tlb.Uint256 // uint256
	Description       boc.Cell    // cell
	CreatedAt         tlb.Uint32  // uint32
	Severity          tlb.Uint8   // uint8
	RewardAddr        tlb.Uint256 // uint256
	Paid              tlb.Coins   // coins
	SuggestedFine     tlb.Coins   // coins
	SuggestedFinePart tlb.Uint32  // uint32
}
type ElectorMember struct {
	Stake     tlb.Coins   // coins
	Time      tlb.Uint32  // uint32
	MaxFactor tlb.Uint32  // uint32
	SrcAddr   tlb.Uint256 // uint256
	AdnlAddr  tlb.Uint256 // uint256
}
type ParticipantListValidatorData struct {
	Stake     tlb.Coins   // coins
	MaxFactor tlb.Uint32  // uint32
	Address   tlb.Uint256 // uint256
	AdnlAddr  tlb.Uint256 // uint256
}
type ParticipantListExtended struct {
	ElectAt    tlb.Uint32                                                    // uint32
	ElectClose tlb.Uint32                                                    // uint32
	MinStake   tlb.Coins                                                     // coins
	TotalStake tlb.Coins                                                     // coins
	Validators []tlb.ShapedTuple2[tlb.Uint256, ParticipantListValidatorData] // array<[uint256, ParticipantListValidatorData]>
	Failed     bool                                                          // bool
	Finished   bool                                                          // bool
}
type Elect struct {
	ElectAt    tlb.Uint32                               // uint32
	ElectClose tlb.Uint32                               // uint32
	MinStake   tlb.Coins                                // coins
	TotalStake tlb.Coins                                // coins
	Members    tlb.HashmapE[tlb.Uint256, ElectorMember] // map<uint256, ElectorMember>
	Failed     bool                                     // bool
	Finished   bool                                     // bool
}
type ElectorStorage struct {
	Elect         tlb.Maybe[tlb.RefT[*Elect]]          // Cell<Elect>?
	Credits       tlb.HashmapE[tlb.Uint256, tlb.Coins] // map<uint256, coins>
	PastElections tlb.HashmapE[tlb.Uint32, boc.Cell]   // map<uint32, slice>
	Grams         tlb.Coins                            // coins
	ActiveId      tlb.Uint32                           // uint32
	ActiveHash    tlb.Uint256                          // uint256
}

const ( // errors
)

func DecodeGetActiveElectionId(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetActiveElectionId = 0x15207

func GetActiveElectionId(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetActiveElectionId, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetActiveElectionId(&stack)
}

func DecodeGetParticipatesIn(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetParticipatesIn = 0x1572C

func GetParticipatesIn(ctx context.Context, executor Executor, reqAccountID ton.AccountID, validatorKey tlb.Uint256) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(validatorKey)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetParticipatesIn, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetParticipatesIn(&stack)
}

func DecodeGetComputeReturnedStake(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetComputeReturnedStake = 0x1FF80

func GetComputeReturnedStake(ctx context.Context, executor Executor, reqAccountID ton.AccountID, walletAddr tlb.Uint256) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(walletAddr)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetComputeReturnedStake, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetComputeReturnedStake(&stack)
}

func DecodeGetParticipantListExtended(stack *tlb.VmStack) (result ParticipantListExtended, err error) {
	if stack.Len() != 7 {
		err = fmt.Errorf("invalid stack size %d, expected 7", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetParticipantListExtended = 0x152AA

func GetParticipantListExtended(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result ParticipantListExtended, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetParticipantListExtended, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetParticipantListExtended(&stack)
}

func Elector_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage ElectorStorage, err error) {
	sa, err = executor.GetAccountState(ctx, accountID)
	if err != nil {
		return
	}
	acc := sa.Account
	if acc.SumType != "Account" {
		err = fmt.Errorf("account does not exist")
	} else if state := acc.Account.Storage.State; state.SumType != "AccountActive" {
		err = fmt.Errorf("account is not active")
	} else if data := state.AccountActive.StateInit.Data; !data.Exists {
		err = fmt.Errorf("account has no storage data")
	} else {
		err = storage.UnmarshalTLB(&data.Value.Value, tlb.NewDecoder())
	}
	return
}

type electorImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewElector(executor Executor, storageExecutor StorageExecutor) Elector {
	return &electorImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c electorImpl) WithAccountId(accountID ton.AccountID) ElectorWithAccount {
	return &electorWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c electorImpl) GetActiveElectionId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetActiveElectionId(ctx, c.executor, reqAccountID)
}

func (c electorImpl) GetParticipatesIn(ctx context.Context, reqAccountID ton.AccountID, validatorKey tlb.Uint256) (tlb.Int257, error) {
	return GetParticipatesIn(ctx, c.executor, reqAccountID, validatorKey)
}

func (c electorImpl) GetComputeReturnedStake(ctx context.Context, reqAccountID ton.AccountID, walletAddr tlb.Uint256) (tlb.Int257, error) {
	return GetComputeReturnedStake(ctx, c.executor, reqAccountID, walletAddr)
}

func (c electorImpl) GetParticipantListExtended(ctx context.Context, reqAccountID ton.AccountID) (ParticipantListExtended, error) {
	return GetParticipantListExtended(ctx, c.executor, reqAccountID)
}

func (c electorImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, ElectorStorage, error) {
	return Elector_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type electorWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c electorWithAccountImpl) GetActiveElectionId(ctx context.Context) (tlb.Int257, error) {
	return GetActiveElectionId(ctx, c.executor, c.accountID)
}

func (c electorWithAccountImpl) GetParticipatesIn(ctx context.Context, validatorKey tlb.Uint256) (tlb.Int257, error) {
	return GetParticipatesIn(ctx, c.executor, c.accountID, validatorKey)
}

func (c electorWithAccountImpl) GetComputeReturnedStake(ctx context.Context, walletAddr tlb.Uint256) (tlb.Int257, error) {
	return GetComputeReturnedStake(ctx, c.executor, c.accountID, walletAddr)
}

func (c electorWithAccountImpl) GetParticipantListExtended(ctx context.Context) (ParticipantListExtended, error) {
	return GetParticipantListExtended(ctx, c.executor, c.accountID)
}

func (c electorWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, ElectorStorage, error) {
	return Elector_AccountState(ctx, c.storageExecutor, c.accountID)
}
