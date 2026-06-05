// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type PositionIncomingMessageKind uint

const (
	PositionIncomingMessageKind_ItemInternalLock   PositionIncomingMessageKind = 2058379649
	PositionIncomingMessageKind_ItemInternalUnlock PositionIncomingMessageKind = 1768923864
	PositionIncomingMessageKind_ItemWithdraw       PositionIncomingMessageKind = 314091911
)

type PositionIncomingMessage struct {
	SumType            PositionIncomingMessageKind
	ItemInternalLock   *ItemInternalLock
	ItemInternalUnlock *ItemInternalUnlock
	ItemWithdraw       *ItemWithdraw
}

type PositionExternalMessageKind uint

const (
	PositionExternalMessageKind_ExternalItemWithdraw PositionExternalMessageKind = 2331607777
	PositionExternalMessageKind_ExternalCronTrigger  PositionExternalMessageKind = 554987565
)

type PositionExternalMessage struct {
	SumType              PositionExternalMessageKind
	ExternalItemWithdraw *ExternalItemWithdraw
	ExternalCronTrigger  *ExternalCronTrigger
}

type GetOrderDataResult struct {
	Status                          tlb.Uint8            // uint8
	QuoteId                         tlb.Int257           // int
	Minter                          tlb.InternalAddress  // address
	BidJettonWallet                 tlb.InternalAddress  // address
	BidJettonAmount                 tlb.Coins            // coins
	BidFromUser                     tlb.InternalAddress  // address
	RefFee                          tlb.InternalAddress  // address
	RefFeeTier                      tlb.Uint16           // uint16
	PendingUserReceiveAmount        tlb.Coins            // coins
	AskJettonWallet                 tlb.InternalAddress  // address
	AskJettonMinter                 tlb.InternalAddress  // address
	ConditionState                  boc.Cell             // cell
	ConditionCode                   boc.Cell             // cell
	UserUnlockForwardParams         tlb.Maybe[tlb.Coins] // coins?
	UserUnlockForwardSuccessPayload tlb.Maybe[boc.Cell]  // cell?
	UserFillToVault                 bool                 // bool
	SafeDepositAndForwardValue      tlb.Uint64           // uint64
	ExcessesReceiver                tlb.InternalAddress  // address
}

type GetCronInfoResult struct {
	NextCallTime        tlb.Uint32 // uint32
	Reward              tlb.Coins  // coins
	BalanceMinusAmounts tlb.Int64  // int64
	RepeatEvery         tlb.Uint32 // uint32
}

type ItemAdditionalFieldMore struct {
	AskJettonMinter tlb.MsgAddress // address?
	OrderOwner      tlb.MsgAddress // address?
	RefundTo        tlb.MsgAddress // address?
}

type ItemAdditionalField struct {
	AskJettonWallet            tlb.MsgAddress                          // address?
	RefFee                     tlb.MsgAddress                          // address?
	RefFeeTier                 tlb.Uint16                              // uint16
	SafeDepositAndForwardValue tlb.Coins                               // coins
	UserFillToVault            bool                                    // bool
	OwnerPubkey                tlb.Maybe[tlb.Uint256]                  // uint256?
	More                       tlb.RefT[*ItemAdditionalFieldMore]      // Cell<ItemAdditionalFieldMore>
	UserUnlockForwardParams    tlb.Maybe[tlb.RefT[*ForwardParams]]     // Cell<ForwardParams>?
	LockForwardParams          tlb.Maybe[tlb.RefT[*LockForwardParams]] // Cell<LockForwardParams>?
}

type ItemInternalUnlockExtraFields struct {
	Recipient tlb.MsgAddress // address?
	RefundTo  tlb.MsgAddress // address?
	Excesses  tlb.MsgAddress // address?
}

type ExternalItemWithdrawPayload struct {
	EscrowItem   tlb.MsgAddress // address?
	SignatureTTL tlb.Uint64     // uint64
	WithdrawArgs boc.Cell       // cell
}

const PrefixEscrowWithdrawSignMessage uint64 = 0x75569022

type EscrowWithdrawSignMessage struct {
	SchemaHash        tlb.Uint32                             // uint32
	Timestamp         tlb.Uint64                             // uint64
	UserWalletAddress tlb.MsgAddress                         // address?
	Domain            boc.Cell                               // cell
	Payload           tlb.RefT[*ExternalItemWithdrawPayload] // Cell<ExternalItemWithdrawPayload>
}

const PrefixItemInternalLock uint64 = 0x7ab06181

type ItemInternalLock struct {
	QueryId         tlb.Uint64                     // uint64
	TokenAddress    tlb.MsgAddress                 // address?
	Amount          tlb.Coins                      // coins
	Excesses        tlb.MsgAddress                 // address?
	RefundToVault   bool                           // bool
	AdditionalField tlb.RefT[*ItemAdditionalField] // Cell<ItemAdditionalField>
	LockArgs        tlb.RefT[*BilateralLockArgs]   // Cell<BilateralLockArgs>
	UnlockCondition boc.Cell                       // cell
}

const PrefixItemInternalUnlock uint64 = 0x696fa2d8

type ItemInternalUnlock struct {
	QueryId                  tlb.Uint64                               // uint64
	FillToVault              bool                                     // bool
	RefundToVault            bool                                     // bool
	Resolver                 tlb.MsgAddress                           // address?
	ResolverSentJettonWallet tlb.MsgAddress                           // address?
	ResolverSentAmount       tlb.Coins                                // coins
	UnlockArgs               tlb.RefT[*BilateralUnlockArgs]           // Cell<BilateralUnlockArgs>
	ExtraFields              tlb.RefT[*ItemInternalUnlockExtraFields] // Cell<ItemInternalUnlockExtraFields>
	ForwardParams            tlb.Maybe[tlb.RefT[*ForwardParams]]      // Cell<ForwardParams>?
}

const PrefixItemWithdraw uint64 = 0x12b8a987

type ItemWithdraw struct {
	QueryId       tlb.Uint64                          // uint64
	WithdrawArgs  boc.Cell                            // cell
	Excesses      tlb.MsgAddress                      // address?
	ForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixExternalItemWithdraw uint64 = 0x8af982e1

type ExternalItemWithdraw struct {
	QueryId     tlb.Uint64                           // uint64
	Signature   tlb.Bits512                          // bits512
	SignMessage tlb.RefT[*EscrowWithdrawSignMessage] // Cell<EscrowWithdrawSignMessage>
}

const PrefixExternalCronTrigger uint64 = 0x2114702d

type ExternalCronTrigger struct {
	RewardAddress tlb.MsgAddress // address?
	Salt          tlb.Uint32     // uint32
}

const ErrorWrongStatus = 0xB4AE // 46254

func DecodeGetOrderData(stack *tlb.VmStack) (result GetOrderDataResult, err error) {
	if stack.Len() != 18 {
		err = fmt.Errorf("invalid stack size %d, expected 18", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetOrderData = 0x1D175

func GetOrderData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetOrderDataResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetOrderData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetOrderData(&stack)
}

func DecodeGetCronInfo(stack *tlb.VmStack) (result GetCronInfoResult, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCronInfo = 0x1305B

func GetCronInfo(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetCronInfoResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCronInfo, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCronInfo(&stack)
}

type escrowPositionImpl struct {
	executor Executor
}

func NewEscrowPosition(executor Executor) EscrowPosition {
	return &escrowPositionImpl{executor: executor}
}

func (c escrowPositionImpl) WithAccountId(accountID ton.AccountID) EscrowPositionWithAccount {
	return &escrowPositionWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c escrowPositionImpl) GetOrderData(ctx context.Context, reqAccountID ton.AccountID) (GetOrderDataResult, error) {
	return GetOrderData(ctx, c.executor, reqAccountID)
}

func (c escrowPositionImpl) GetCronInfo(ctx context.Context, reqAccountID ton.AccountID) (GetCronInfoResult, error) {
	return GetCronInfo(ctx, c.executor, reqAccountID)
}

type escrowPositionWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c escrowPositionWithAccountImpl) GetOrderData(ctx context.Context) (GetOrderDataResult, error) {
	return GetOrderData(ctx, c.executor, c.accountID)
}

func (c escrowPositionWithAccountImpl) GetCronInfo(ctx context.Context) (GetCronInfoResult, error) {
	return GetCronInfo(ctx, c.executor, c.accountID)
}
