// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type FactoryIncomingMessageKind uint

const (
	FactoryIncomingMessageKind_MinterInitTransfer            FactoryIncomingMessageKind = 1192680726
	FactoryIncomingMessageKind_MinterRefundRequest           FactoryIncomingMessageKind = 3608432653
	FactoryIncomingMessageKind_MinterInternalLock            FactoryIncomingMessageKind = 1707821141
	FactoryIncomingMessageKind_MinterInternalUnlock          FactoryIncomingMessageKind = 821961387
	FactoryIncomingMessageKind_MinterInternalWithdrawTokens  FactoryIncomingMessageKind = 742082194
	FactoryIncomingMessageKind_MinterGiveProtocolOwnership   FactoryIncomingMessageKind = 1576598069
	FactoryIncomingMessageKind_MinterUpdateProtocolTier      FactoryIncomingMessageKind = 3727994588
	FactoryIncomingMessageKind_MinterTakeProtocolOwnership   FactoryIncomingMessageKind = 2282898964
	FactoryIncomingMessageKind_MinterCancelNextProtocolOwner FactoryIncomingMessageKind = 133727974
	FactoryIncomingMessageKind_MinterResetGas                FactoryIncomingMessageKind = 4262875806
	FactoryIncomingMessageKind_MinterLockPayload             FactoryIncomingMessageKind = 1042572932
	FactoryIncomingMessageKind_MinterUnlockPayload           FactoryIncomingMessageKind = 1784215644
	FactoryIncomingMessageKind_MinterDepositVault            FactoryIncomingMessageKind = 1881886152
)

type FactoryIncomingMessage struct {
	SumType                       FactoryIncomingMessageKind
	MinterInitTransfer            *MinterInitTransfer
	MinterRefundRequest           *MinterRefundRequest
	MinterInternalLock            *MinterInternalLock
	MinterInternalUnlock          *MinterInternalUnlock
	MinterInternalWithdrawTokens  *MinterInternalWithdrawTokens
	MinterGiveProtocolOwnership   *MinterGiveProtocolOwnership
	MinterUpdateProtocolTier      *MinterUpdateProtocolTier
	MinterTakeProtocolOwnership   *MinterTakeProtocolOwnership
	MinterCancelNextProtocolOwner *MinterCancelNextProtocolOwner
	MinterResetGas                *MinterResetGas
	MinterLockPayload             *MinterLockPayload
	MinterUnlockPayload           *MinterUnlockPayload
	MinterDepositVault            *MinterDepositVault
}

type EscrowData struct {
	GlobalId     tlb.Uint64          // uint64
	Protocol     tlb.InternalAddress // address
	ProtocolTier tlb.Uint16          // uint16
	ItemCode     boc.Cell            // cell
	VaultCode    boc.Cell            // cell
}

type GetVersionResult struct {
	Major       tlb.Uint8 // uint8
	Minor       tlb.Uint8 // uint8
	Development string    // string
}

type BidPaymentData struct {
	UsingVault      bool                                // bool
	Recipient       tlb.MsgAddress                      // address?
	BidJettonWallet tlb.MsgAddress                      // address?
	Amount          tlb.Coins                           // coins
	ForwardParams   tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

type AskRefundData struct {
	RefundTo   tlb.MsgAddress // address?
	UsingVault bool           // bool
	Amount     tlb.Coins      // coins
}

type UserPaymentData struct {
	User                           tlb.MsgAddress                      // address?
	UserReceiveJettonWallet        tlb.MsgAddress                      // address?
	UserReceiveAmount              tlb.Coins                           // coins
	UserDepositViaVault            bool                                // bool
	UserForwardParams              tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
	UserExcesses                   tlb.MsgAddress                      // address?
	UserSafeDepositAndForwardValue tlb.Uint64                          // uint64
}

type MinterRefundRequestExtraFields struct {
	JettonWallet     tlb.MsgAddress            // address?
	Excesses         tlb.MsgAddress            // address?
	ForwardTonAmount tlb.Coins                 // coins
	ForwardPayload   tlb.Maybe[boc.Cell]       // cell?
	LockRejectDest   tlb.Maybe[tlb.MsgAddress] // any_address?
}

type InternalLockExtraNested struct {
	AskJettonWallet tlb.MsgAddress         // address?
	OwnerPubkey     tlb.Maybe[tlb.Uint256] // uint256?
}

type InternalLockExtra struct {
	QuoteId                    tlb.Uint256                             // uint256
	RefFee                     tlb.MsgAddress                          // address?
	RefFeeTier                 tlb.Uint16                              // uint16
	Excesses                   tlb.MsgAddress                          // address?
	SafeDepositAndForwardValue tlb.Coins                               // coins
	RefundToVault              bool                                    // bool
	UserFillToVault            bool                                    // bool
	More                       tlb.RefT[*VaultLockAdditionalDataMore]  // Cell<VaultLockAdditionalDataMore>
	UserUnlockForwardParams    tlb.Maybe[tlb.RefT[*ForwardParams]]     // Cell<ForwardParams>?
	LockForwardParams          tlb.Maybe[tlb.RefT[*LockForwardParams]] // Cell<LockForwardParams>?
	Nested                     tlb.RefT[*InternalLockExtraNested]      // Cell<InternalLockExtraNested>
}

type UnlockAdditionalData struct {
	QuoteId       tlb.Uint256                         // uint256
	FillToVault   bool                                // bool
	RefundToVault bool                                // bool
	Recipient     tlb.MsgAddress                      // address?
	RefundTo      tlb.MsgAddress                      // address?
	UnlockArgs    tlb.RefT[*BilateralUnlockArgs]      // Cell<BilateralUnlockArgs>
	ForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

type MinterLockPayloadExtraFieldsMore struct {
	OrderOwner      tlb.MsgAddress // address?
	AskJettonMinter tlb.MsgAddress // address?
}

type MinterLockPayloadExtraFields struct {
	AskJettonWallet         tlb.MsgAddress                              // address?
	RefundTo                tlb.MsgAddress                              // address?
	OwnerPubkey             tlb.Maybe[tlb.Uint256]                      // uint256?
	More                    tlb.RefT[*MinterLockPayloadExtraFieldsMore] // Cell<MinterLockPayloadExtraFieldsMore>
	UserUnlockForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]]         // Cell<ForwardParams>?
	LockForwardParams       tlb.Maybe[tlb.RefT[*LockForwardParams]]     // Cell<LockForwardParams>?
}

type UnlockPayloadExtraFields struct {
	Recipient     tlb.MsgAddress // address?
	RefundTo      tlb.MsgAddress // address?
	FillToVault   bool           // bool
	RefundToVault bool           // bool
}

const PrefixMinterInitTransfer uint64 = 0x4716dd16

type MinterInitTransfer struct {
	QueryId         tlb.Uint64                 // uint64
	QuoteId         tlb.Uint256                // uint256
	TransferType    tlb.Uint32                 // uint32
	RefFee          tlb.MsgAddress             // address?
	RefFeeTier      tlb.Uint16                 // uint16
	Excesses        tlb.MsgAddress             // address?
	BidPaymentData  tlb.RefT[*BidPaymentData]  // Cell<BidPaymentData>
	AskRefundData   tlb.RefT[*AskRefundData]   // Cell<AskRefundData>
	UserPaymentData tlb.RefT[*UserPaymentData] // Cell<UserPaymentData>
}

const PrefixMinterRefundRequest uint64 = 0xd714500d

type MinterRefundRequest struct {
	QueryId     tlb.Uint64                                // uint64
	QuoteId     tlb.Uint256                               // uint256
	Amount      tlb.Coins                                 // coins
	UseVault    bool                                      // bool
	Recipient   tlb.MsgAddress                            // address?
	ExitCode    tlb.Uint32                                // uint32
	PrevMessage tlb.Uint32                                // uint32
	ExtraFields tlb.RefT[*MinterRefundRequestExtraFields] // Cell<MinterRefundRequestExtraFields>
}

const PrefixMinterInternalLock uint64 = 0x65cb4855

type MinterInternalLock struct {
	QueryId         tlb.Uint64                   // uint64
	Owner           tlb.MsgAddress               // address?
	JettonWallet    tlb.MsgAddress               // address?
	Amount          tlb.Coins                    // coins
	LockArgs        tlb.RefT[*BilateralLockArgs] // Cell<BilateralLockArgs>
	UnlockCondition boc.Cell                     // cell
	AdditionalData  tlb.RefT[*InternalLockExtra] // Cell<InternalLockExtra>
}

const PrefixMinterInternalUnlock uint64 = 0x30fe22ab

type MinterInternalUnlock struct {
	QueryId        tlb.Uint64                      // uint64
	Owner          tlb.MsgAddress                  // address?
	JettonWallet   tlb.MsgAddress                  // address?
	Amount         tlb.Coins                       // coins
	Excesses       tlb.MsgAddress                  // address?
	AdditionalData tlb.RefT[*UnlockAdditionalData] // Cell<UnlockAdditionalData>
}

const PrefixMinterInternalWithdrawTokens uint64 = 0x2c3b4692

type MinterInternalWithdrawTokens struct {
	QueryId       tlb.Uint64                          // uint64
	JettonWallet  tlb.MsgAddress                      // address?
	CollectedFees tlb.Coins                           // coins
	Owner         tlb.MsgAddress                      // address?
	Excesses      tlb.MsgAddress                      // address?
	ForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixMinterGiveProtocolOwnership uint64 = 0x5df8fa35

type MinterGiveProtocolOwnership struct {
	QueryId            tlb.Uint64     // uint64
	NewProtocolAddress tlb.MsgAddress // address?
	Excesses           tlb.MsgAddress // address?
}

const PrefixMinterUpdateProtocolTier uint64 = 0xde34aedc

type MinterUpdateProtocolTier struct {
	QueryId  tlb.Uint64     // uint64
	NewTier  tlb.Uint16     // uint16
	Excesses tlb.MsgAddress // address?
}

const PrefixMinterTakeProtocolOwnership uint64 = 0x88124614

type MinterTakeProtocolOwnership struct {
	QueryId  tlb.Uint64     // uint64
	Excesses tlb.MsgAddress // address?
}

const PrefixMinterCancelNextProtocolOwner uint64 = 0x07f886e6

type MinterCancelNextProtocolOwner struct {
	QueryId  tlb.Uint64     // uint64
	Excesses tlb.MsgAddress // address?
}

const PrefixMinterResetGas uint64 = 0xfe16529e

type MinterResetGas struct {
	QueryId   tlb.Uint64     // uint64
	Recipient tlb.MsgAddress // address?
}

const PrefixMinterLockPayload uint64 = 0x3e246684

type MinterLockPayload struct {
	QuoteId         tlb.Uint256                             // uint256
	RefundToVault   bool                                    // bool
	RefFee          tlb.MsgAddress                          // address?
	RefFeeTier      tlb.Uint16                              // uint16
	Excesses        tlb.MsgAddress                          // address?
	TonSafeDeposit  tlb.Coins                               // coins
	UserFillToVault bool                                    // bool
	LockArgs        tlb.RefT[*BilateralLockArgs]            // Cell<BilateralLockArgs>
	UnlockCondition boc.Cell                                // cell
	ExtraFields     tlb.RefT[*MinterLockPayloadExtraFields] // Cell<MinterLockPayloadExtraFields>
}

const PrefixMinterUnlockPayload uint64 = 0x6a58f85c

type MinterUnlockPayload struct {
	QuoteId       tlb.Uint256                         // uint256
	Excesses      tlb.MsgAddress                      // address?
	ExtraFields   tlb.RefT[*UnlockPayloadExtraFields] // Cell<UnlockPayloadExtraFields>
	UnlockArgs    tlb.RefT[*BilateralUnlockArgs]      // Cell<BilateralUnlockArgs>
	ForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixMinterDepositVault uint64 = 0x702b4dc8

type MinterDepositVault struct {
	Excesses      tlb.MsgAddress                      // address?
	ForwardParams tlb.Maybe[tlb.RefT[*ForwardParams]] // Cell<ForwardParams>?
}

const PrefixLockRejectNotification uint64 = 0xac8a7373

type LockRejectNotification struct {
	QueryId        tlb.Uint64          // uint64
	ForwardPayload tlb.Maybe[boc.Cell] // cell?
}

const ErrorInvalidCall = 0x1446 // 5190

const ErrorWrongWorkchain = 0x82FD // 33533

const ErrorInvalidTransferNotificationPayload = 0x8762 // 34658

func DecodeGetEscrowData(stack *tlb.VmStack) (result EscrowData, err error) {
	if stack.Len() != 5 {
		err = fmt.Errorf("invalid stack size %d, expected 5", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetEscrowData = 0x12DFB

func GetEscrowData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result EscrowData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetEscrowData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetEscrowData(&stack)
}

func DecodeGetItemAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetItemAddress = 0x10061

func GetItemAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, quoteId tlb.Int257) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(quoteId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetItemAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetItemAddress(&stack)
}

func DecodeGetVaultAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVaultAddress = 0x1C364

func GetVaultAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, jettonWallet tlb.InternalAddress, owner tlb.InternalAddress) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(jettonWallet)
		if err != nil {
			err = fmt.Errorf("encode param jettonWallet: %w", err)
			return
		}
		stack.Put(val)
	}
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(owner)
		if err != nil {
			err = fmt.Errorf("encode param owner: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVaultAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVaultAddress(&stack)
}

func DecodeGetVersion(stack *tlb.VmStack) (result GetVersionResult, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetVersion = 0x123E4

func GetVersion(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetVersionResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetVersion, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetVersion(&stack)
}

type escrowFactoryImpl struct {
	executor Executor
}

func NewEscrowFactory(executor Executor) EscrowFactory {
	return &escrowFactoryImpl{executor: executor}
}

func (c escrowFactoryImpl) WithAccountId(accountID ton.AccountID) EscrowFactoryWithAccount {
	return &escrowFactoryWithAccountImpl{executor: c.executor, accountID: accountID}
}

func (c escrowFactoryImpl) GetEscrowData(ctx context.Context, reqAccountID ton.AccountID) (EscrowData, error) {
	return GetEscrowData(ctx, c.executor, reqAccountID)
}

func (c escrowFactoryImpl) GetItemAddress(ctx context.Context, reqAccountID ton.AccountID, quoteId tlb.Int257) (tlb.InternalAddress, error) {
	return GetItemAddress(ctx, c.executor, reqAccountID, quoteId)
}

func (c escrowFactoryImpl) GetVaultAddress(ctx context.Context, reqAccountID ton.AccountID, jettonWallet tlb.InternalAddress, owner tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetVaultAddress(ctx, c.executor, reqAccountID, jettonWallet, owner)
}

func (c escrowFactoryImpl) GetVersion(ctx context.Context, reqAccountID ton.AccountID) (GetVersionResult, error) {
	return GetVersion(ctx, c.executor, reqAccountID)
}

type escrowFactoryWithAccountImpl struct {
	executor  Executor
	accountID ton.AccountID
}

func (c escrowFactoryWithAccountImpl) GetEscrowData(ctx context.Context) (EscrowData, error) {
	return GetEscrowData(ctx, c.executor, c.accountID)
}

func (c escrowFactoryWithAccountImpl) GetItemAddress(ctx context.Context, quoteId tlb.Int257) (tlb.InternalAddress, error) {
	return GetItemAddress(ctx, c.executor, c.accountID, quoteId)
}

func (c escrowFactoryWithAccountImpl) GetVaultAddress(ctx context.Context, jettonWallet tlb.InternalAddress, owner tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetVaultAddress(ctx, c.executor, c.accountID, jettonWallet, owner)
}

func (c escrowFactoryWithAccountImpl) GetVersion(ctx context.Context) (GetVersionResult, error) {
	return GetVersion(ctx, c.executor, c.accountID)
}
