// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type GramboTradeFees struct {
	ReferrerFee tlb.Coins // coins
	Fee1        tlb.Coins // coins
	Fee2        tlb.Coins // coins
}

const PrefixGramboBuy uint64 = 0x2f54494e

type GramboBuy struct {
	QueryId       tlb.Uint64          // uint64
	TonAmount     tlb.Coins           // coins
	Recipient     tlb.InternalAddress // address
	MinTokensOut  tlb.Coins           // coins
	Referrer      tlb.InternalAddress // address
	CustomPayload tlb.Maybe[boc.Cell] // cell?
}

const PrefixGramboProvideWalletAddress uint64 = 0x17b5012c

type GramboProvideWalletAddress struct {
	QueryId        tlb.Uint64          // uint64
	Owner          tlb.InternalAddress // address
	IncludeAddress bool                // bool
}

const PrefixGramboUnlockWallet uint64 = 0x7966410f

type GramboUnlockWallet struct {
	QueryId tlb.Uint64          // uint64
	Owner   tlb.InternalAddress // address
}

const PrefixGramboWithdraw uint64 = 0x43e84bb5

type GramboWithdraw struct {
	QueryId tlb.Uint64 // uint64
	Amount  tlb.Coins  // coins
}

type MasterIncomingMessageKind uint

const (
	MasterIncomingMessageKind_GramboLaunch               MasterIncomingMessageKind = 925833955
	MasterIncomingMessageKind_GramboBuy                  MasterIncomingMessageKind = 794052942
	MasterIncomingMessageKind_GramboSellNotification     MasterIncomingMessageKind = 838603382
	MasterIncomingMessageKind_GramboProvideWalletAddress MasterIncomingMessageKind = 397738284
	MasterIncomingMessageKind_GramboUnlockWallet         MasterIncomingMessageKind = 2036744463
	MasterIncomingMessageKind_GramboWithdraw             MasterIncomingMessageKind = 1139297205
)

type MasterIncomingMessage struct {
	SumType                    MasterIncomingMessageKind
	GramboLaunch               *GramboLaunch
	GramboBuy                  *GramboBuy
	GramboSellNotification     *GramboSellNotification
	GramboProvideWalletAddress *GramboProvideWalletAddress
	GramboUnlockWallet         *GramboUnlockWallet
	GramboWithdraw             *GramboWithdraw
}

const PrefixJettonInternalTransfer uint64 = 0x178d4519

type JettonInternalTransfer struct {
	QueryId          tlb.Uint64          // uint64
	Amount           tlb.Coins           // coins
	From             tlb.InternalAddress // address
	ResponseAddress  tlb.InternalAddress // address
	ForwardTonAmount tlb.Coins           // coins
	ForwardPayload   boc.Cell            // cell
}

const PrefixGramboTakeWalletAddress uint64 = 0x29cb7264

type GramboTakeWalletAddress struct {
	QueryId      tlb.Uint64          // uint64
	OwnerAddress tlb.Maybe[boc.Cell] // cell?
}

const PrefixJettonExcess uint64 = 0xd53276db

type JettonExcess struct {
	QueryId tlb.Uint64 // uint64
}

type MasterOutgoingMessageKind uint

const (
	MasterOutgoingMessageKind_JettonInternalTransfer  MasterOutgoingMessageKind = 395134233
	MasterOutgoingMessageKind_GramboActivateWallet    MasterOutgoingMessageKind = 1802884066
	MasterOutgoingMessageKind_GramboTakeWalletAddress MasterOutgoingMessageKind = 701198948
	MasterOutgoingMessageKind_JettonExcess            MasterOutgoingMessageKind = 3576854235
)

type MasterOutgoingMessage struct {
	SumType                 MasterOutgoingMessageKind
	JettonInternalTransfer  *JettonInternalTransfer
	GramboActivateWallet    *GramboActivateWallet
	GramboTakeWalletAddress *GramboTakeWalletAddress
	JettonExcess            *JettonExcess
}

const PrefixGramboSwapEvent uint64 = 0x5025398e

type GramboSwapEvent struct {
	QueryId         tlb.Uint64                 // uint64
	Recipient       tlb.InternalAddress        // address
	IsBuy           bool                       // bool
	NewCollectedTon tlb.Coins                  // coins
	TonAmount       tlb.Coins                  // coins
	TokenAmount     tlb.Coins                  // coins
	Referrer        tlb.InternalAddress        // address
	Fees            tlb.RefT[*GramboTradeFees] // Cell<GramboTradeFees>
	CustomPayload   tlb.Maybe[boc.Cell]        // cell?
}

const PrefixGramboGraduateEvent uint64 = 0x0ab3d75c

type GramboGraduateEvent struct {
	QueryId     tlb.Uint64 // uint64
	TonToDex    tlb.Coins  // coins
	TokensToDex tlb.Coins  // coins
}

type MasterExternalMessageKind uint

const (
	MasterExternalMessageKind_GramboSwapEvent     MasterExternalMessageKind = 1344616846
	MasterExternalMessageKind_GramboGraduateEvent MasterExternalMessageKind = 179558236
)

type MasterExternalMessage struct {
	SumType             MasterExternalMessageKind
	GramboSwapEvent     *GramboSwapEvent
	GramboGraduateEvent *GramboGraduateEvent
}

type GramboCurveState struct {
	CollectedTon tlb.Coins           // coins
	TokensSold   tlb.Coins           // coins
	Graduated    bool                // bool
	FeeAddr1     tlb.InternalAddress // address
	FeeAddr2     tlb.InternalAddress // address
}

type GramboDexConfig struct {
	StonfiRouter tlb.InternalAddress // address
	PtonWallet   tlb.InternalAddress // address
}

type GramboMasterStorage struct {
	VirtualTonReserve tlb.Coins                   // coins
	TonTarget         tlb.Coins                   // coins
	Factory           tlb.InternalAddress         // address
	TotalSupply       tlb.Coins                   // coins
	Seed              tlb.Uint256                 // uint256
	CurveState        tlb.RefT[*GramboCurveState] // Cell<GramboCurveState>
	WalletCode        boc.Cell                    // cell
	Metadata          boc.Cell                    // cell
	DexConfig         tlb.RefT[*GramboDexConfig]  // Cell<GramboDexConfig>
}

type GetCurveStateResult struct {
	VirtualTonReserve tlb.Uint64 // uint64
	TonTarget         tlb.Uint64 // uint64
	CollectedTon      tlb.Uint64 // uint64
	TokensSold        tlb.Int257 // int
	Graduated         bool       // bool
}

type GetJettonDataResult struct {
	TotalSupply      tlb.Int257          // int
	Mintable         tlb.Int257          // int
	AdminAddress     tlb.InternalAddress // address
	JettonContent    boc.Cell            // cell
	JettonWalletCode boc.Cell            // cell
}

func DecodeGetAmountOut(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetAmountOut = 0x15D4C

func GetAmountOut(ctx context.Context, executor Executor, reqAccountID ton.AccountID, tonIn tlb.Uint64) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257FromInt64(int64(tonIn))})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetAmountOut, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetAmountOut(&stack)
}

func DecodeGetTonOut(stack *tlb.VmStack) (result tlb.Uint64, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetTonOut = 0x1BD1F

func GetTonOut(ctx context.Context, executor Executor, reqAccountID ton.AccountID, tokensIn tlb.Int257) (result tlb.Uint64, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(tokensIn)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetTonOut, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetTonOut(&stack)
}

func DecodeGetCurveState(stack *tlb.VmStack) (result GetCurveStateResult, err error) {
	if stack.Len() != 5 {
		err = fmt.Errorf("invalid stack size %d, expected 5", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCurveState = 0x1E48C

func GetCurveState(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetCurveStateResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCurveState, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCurveState(&stack)
}

func DecodeGetJettonData(stack *tlb.VmStack) (result GetJettonDataResult, err error) {
	if stack.Len() != 5 {
		err = fmt.Errorf("invalid stack size %d, expected 5", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetJettonData = 0x19E2D

func GetJettonData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result GetJettonDataResult, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetJettonData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetJettonData(&stack)
}

func DecodeGetWalletAddress(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetWalletAddress = 0x19379

func GetWalletAddress(ctx context.Context, executor Executor, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCellSlice(ownerAddress)
		if err != nil {
			err = fmt.Errorf("encode param ownerAddress: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetWalletAddress, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetWalletAddress(&stack)
}

func GramboJettonMaster_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage GramboMasterStorage, err error) {
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

type gramboJettonMasterImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewGramboJettonMaster(executor Executor, storageExecutor StorageExecutor) GramboJettonMaster {
	return &gramboJettonMasterImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c gramboJettonMasterImpl) WithAccountId(accountID ton.AccountID) GramboJettonMasterWithAccount {
	return &gramboJettonMasterWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c gramboJettonMasterImpl) GetAmountOut(ctx context.Context, reqAccountID ton.AccountID, tonIn tlb.Uint64) (tlb.Int257, error) {
	return GetAmountOut(ctx, c.executor, reqAccountID, tonIn)
}

func (c gramboJettonMasterImpl) GetTonOut(ctx context.Context, reqAccountID ton.AccountID, tokensIn tlb.Int257) (tlb.Uint64, error) {
	return GetTonOut(ctx, c.executor, reqAccountID, tokensIn)
}

func (c gramboJettonMasterImpl) GetCurveState(ctx context.Context, reqAccountID ton.AccountID) (GetCurveStateResult, error) {
	return GetCurveState(ctx, c.executor, reqAccountID)
}

func (c gramboJettonMasterImpl) GetJettonData(ctx context.Context, reqAccountID ton.AccountID) (GetJettonDataResult, error) {
	return GetJettonData(ctx, c.executor, reqAccountID)
}

func (c gramboJettonMasterImpl) GetWalletAddress(ctx context.Context, reqAccountID ton.AccountID, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetWalletAddress(ctx, c.executor, reqAccountID, ownerAddress)
}

func (c gramboJettonMasterImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, GramboMasterStorage, error) {
	return GramboJettonMaster_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type gramboJettonMasterWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c gramboJettonMasterWithAccountImpl) GetAmountOut(ctx context.Context, tonIn tlb.Uint64) (tlb.Int257, error) {
	return GetAmountOut(ctx, c.executor, c.accountID, tonIn)
}

func (c gramboJettonMasterWithAccountImpl) GetTonOut(ctx context.Context, tokensIn tlb.Int257) (tlb.Uint64, error) {
	return GetTonOut(ctx, c.executor, c.accountID, tokensIn)
}

func (c gramboJettonMasterWithAccountImpl) GetCurveState(ctx context.Context) (GetCurveStateResult, error) {
	return GetCurveState(ctx, c.executor, c.accountID)
}

func (c gramboJettonMasterWithAccountImpl) GetJettonData(ctx context.Context) (GetJettonDataResult, error) {
	return GetJettonData(ctx, c.executor, c.accountID)
}

func (c gramboJettonMasterWithAccountImpl) GetWalletAddress(ctx context.Context, ownerAddress tlb.InternalAddress) (tlb.InternalAddress, error) {
	return GetWalletAddress(ctx, c.executor, c.accountID, ownerAddress)
}

func (c gramboJettonMasterWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, GramboMasterStorage, error) {
	return GramboJettonMaster_AccountState(ctx, c.storageExecutor, c.accountID)
}
