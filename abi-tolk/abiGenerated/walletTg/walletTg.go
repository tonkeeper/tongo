// Code generated - DO NOT EDIT.

package abiWalletTg

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixStorage uint64 = 0x00

type Storage struct {
	Seqno       tlb.Uint32  // uint32
	SubwalletId tlb.Uint32  // uint32
	PublicKey   tlb.Uint256 // uint256
}
type Revision tlb.Uint8

const (
	RevisionRev00Initial Revision = 0
)

type MessageToSend struct {
	SendMode    tlb.Uint8 // uint8
	MessageCell boc.Cell  // cell
}
type RawArrayOfMessagesToSend []MessageToSend
type KeyRotationProofPayload struct {
	Tag             tlb.Uint96  // uint96
	WalletWorkchain tlb.Int8    // int8
	WalletAddrHash  tlb.Uint256 // uint256
}
type SignedInternalRequest struct {
	Signature tlb.Bits512        // bits512
	Request   AllowedInternalMsg // AllowedInternalMsg
}
type SignedExternalRequest struct {
	Signature tlb.Bits512        // bits512
	Request   AllowedExternalMsg // AllowedExternalMsg
}
type AllowedInternalMsgKind uint

const (
	AllowedInternalMsgKind_SendOneMessageRequestI   AllowedInternalMsgKind = 1669951092
	AllowedInternalMsgKind_SendBulkMessagesRequestI AllowedInternalMsgKind = 1938386548
	AllowedInternalMsgKind_ChangePublicKeyRequestI  AllowedInternalMsgKind = 4223310279
)

type AllowedInternalMsg struct { // tagged union
	SumType                  AllowedInternalMsgKind
	SendOneMessageRequestI   *SendOneMessageRequestI
	SendBulkMessagesRequestI *SendBulkMessagesRequestI
	ChangePublicKeyRequestI  *ChangePublicKeyRequestI
}
type AllowedExternalMsgKind uint

const (
	AllowedExternalMsgKind_SendOneMessageRequestE   AllowedExternalMsgKind = 1669951093
	AllowedExternalMsgKind_SendBulkMessagesRequestE AllowedExternalMsgKind = 1938386549
	AllowedExternalMsgKind_ChangePublicKeyRequestE  AllowedExternalMsgKind = 4223310280
)

type AllowedExternalMsg struct { // tagged union
	SumType                  AllowedExternalMsgKind
	SendOneMessageRequestE   *SendOneMessageRequestE
	SendBulkMessagesRequestE *SendBulkMessagesRequestE
	ChangePublicKeyRequestE  *ChangePublicKeyRequestE
}
type SeqnoHeader struct {
	SubwalletId tlb.Uint32 // uint32
	ValidUntil  tlb.Uint32 // uint32
	Seqno       tlb.Uint32 // uint32
}

const PrefixSendOneMessageRequestI uint64 = 0x63896e74

type SendOneMessageRequestI struct {
	Header SeqnoHeader   // SeqnoHeader
	Msg    MessageToSend // MessageToSend
}

const PrefixSendOneMessageRequestE uint64 = 0x63896e75

type SendOneMessageRequestE struct {
	Header SeqnoHeader   // SeqnoHeader
	Msg    MessageToSend // MessageToSend
}

const PrefixSendBulkMessagesRequestI uint64 = 0x73896e74

type SendBulkMessagesRequestI struct {
	Header SeqnoHeader              // SeqnoHeader
	MsgArr RawArrayOfMessagesToSend // RawArrayOfMessagesToSend
}

const PrefixSendBulkMessagesRequestE uint64 = 0x73896e75

type SendBulkMessagesRequestE struct {
	Header SeqnoHeader              // SeqnoHeader
	MsgArr RawArrayOfMessagesToSend // RawArrayOfMessagesToSend
}

const PrefixChangePublicKeyRequestI uint64 = 0xfbba99c7

type ChangePublicKeyRequestI struct {
	Header            SeqnoHeader            // SeqnoHeader
	NewPublicKey      tlb.Uint256            // uint256
	RotationSignature tlb.RefT[*tlb.Bits512] // Cell<bits512>
}

const PrefixChangePublicKeyRequestE uint64 = 0xfbba99c8

type ChangePublicKeyRequestE struct {
	Header            SeqnoHeader            // SeqnoHeader
	NewPublicKey      tlb.Uint256            // uint256
	RotationSignature tlb.RefT[*tlb.Bits512] // Cell<bits512>
}

const ( // errors
	ErrorsBeginParseOfNull                                = 0x7    // 7
	ErrorsInvalidSeqno                                    = 0x85   // 133
	ErrorsInvalidSubwalletId                              = 0x86   // 134
	ErrorsInvalidSignature                                = 0x87   // 135
	ErrorsExpired                                         = 0x88   // 136
	ErrorsExternalSendMessageMustHaveIgnoreErrorsSendMode = 0x89   // 137
	ErrorsInvalidOutArrayLen                              = 0x93   // 147
	ErrorsNewKeyMustDiffer                                = 0x94   // 148
	ErrorsNewKeySignatureIsInvalid                        = 0x95   // 149
	ErrorsInvalidOpcode                                   = 0xFFFF // 65535
)

func DecodeGetRevision(stack *tlb.VmStack) (result Revision, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetRevision = 0x1C6E5

func GetRevision(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result Revision, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetRevision, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetRevision(&stack)
}

func DecodeGetSeqno(stack *tlb.VmStack) (result tlb.Uint32, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSeqno = 0x14C97

func GetSeqno(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Uint32, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSeqno, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSeqno(&stack)
}

func DecodeGetSubwalletId(stack *tlb.VmStack) (result tlb.Uint32, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSubwalletId = 0x13E3B

func GetSubwalletId(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Uint32, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSubwalletId, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSubwalletId(&stack)
}

func DecodeGetPublicKey(stack *tlb.VmStack) (result tlb.Uint256, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPublicKey = 0x1339C

func GetPublicKey(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Uint256, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPublicKey, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPublicKey(&stack)
}

func WalletTg_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage Storage, err error) {
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

type walletTgImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewWalletTg(executor Executor, storageExecutor StorageExecutor) WalletTg {
	return &walletTgImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c walletTgImpl) WithAccountId(accountID ton.AccountID) WalletTgWithAccount {
	return &walletTgWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c walletTgImpl) GetRevision(ctx context.Context, reqAccountID ton.AccountID) (Revision, error) {
	return GetRevision(ctx, c.executor, reqAccountID)
}

func (c walletTgImpl) GetSeqno(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error) {
	return GetSeqno(ctx, c.executor, reqAccountID)
}

func (c walletTgImpl) GetSubwalletId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint32, error) {
	return GetSubwalletId(ctx, c.executor, reqAccountID)
}

func (c walletTgImpl) GetPublicKey(ctx context.Context, reqAccountID ton.AccountID) (tlb.Uint256, error) {
	return GetPublicKey(ctx, c.executor, reqAccountID)
}

func (c walletTgImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, Storage, error) {
	return WalletTg_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type walletTgWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c walletTgWithAccountImpl) GetRevision(ctx context.Context) (Revision, error) {
	return GetRevision(ctx, c.executor, c.accountID)
}

func (c walletTgWithAccountImpl) GetSeqno(ctx context.Context) (tlb.Uint32, error) {
	return GetSeqno(ctx, c.executor, c.accountID)
}

func (c walletTgWithAccountImpl) GetSubwalletId(ctx context.Context) (tlb.Uint32, error) {
	return GetSubwalletId(ctx, c.executor, c.accountID)
}

func (c walletTgWithAccountImpl) GetPublicKey(ctx context.Context) (tlb.Uint256, error) {
	return GetPublicKey(ctx, c.executor, c.accountID)
}

func (c walletTgWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, Storage, error) {
	return WalletTg_AccountState(ctx, c.storageExecutor, c.accountID)
}
