// Code generated - DO NOT EDIT.

package abiFfVault

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixAssetDeposit uint64 = 0x534917c3

type AssetDeposit struct {
	PriceUpdateData tlb.Maybe[boc.Cell] // cell?
}

const PrefixStableDeposit uint64 = 0x3d669db7

type StableDeposit struct {
	DistributeBetweenStakers bool // bool
}

const PrefixStakeOperation uint64 = 0x7a294fdc

type StakeOperation struct {
	QueryId        tlb.Uint64          // uint64
	DepositAddress tlb.InternalAddress // address
	AssetAmount    tlb.Coins           // coins
}

const PrefixUnstakeExecute uint64 = 0x4e835078

type UnstakeExecute struct {
	QueryId         tlb.Uint64 // uint64
	PriceUpdateData boc.Cell   // cell
	AmountToUnstake tlb.Coins  // coins
	ItemIndex       tlb.Uint64 // uint64
}

const PrefixUnstakeOperation uint64 = 0x97004d72

type UnstakeOperation struct {
	QueryId         tlb.Uint64 // uint64
	AmountToUnstake tlb.Coins  // coins
	ItemIndex       tlb.Uint64 // uint64
}
type PositionData struct {
	XautPriceStaked                 tlb.Uint256 // uint256
	PreviousStablePerAssetCollected tlb.Uint256 // uint256
}

const PrefixUnstakeExecuteInternalCallback uint64 = 0x61755200

type UnstakeExecuteInternalCallback struct {
	QueryId         tlb.Uint64              // uint64
	ItemIndex       tlb.Uint64              // uint64
	AmountToUnstake tlb.Coins               // coins
	XautPriceActual tlb.Uint256             // uint256
	Owner           tlb.InternalAddress     // address
	PositionData    tlb.RefT[*PositionData] // Cell<PositionData>
}

const PrefixWithdrawJetton uint64 = 0xe44c1944

type WithdrawJetton struct {
	QueryId      tlb.Uint64          // uint64
	JettonWallet tlb.InternalAddress // address
	ToAddress    tlb.InternalAddress // address
	Amount       tlb.Coins           // coins
}
type VaultStaticData struct {
	PYTHOracle         tlb.InternalAddress // address
	AssetJettonWallet  tlb.InternalAddress // address
	StableJettonWallet tlb.InternalAddress // address
	NftItemCode        boc.Cell            // cell
	StableDecimals     tlb.Uint8           // uint8
	AssetDecimals      tlb.Uint8           // uint8
}

const PrefixVaultStorage uint64 = 0x00

type VaultStorage struct {
	AssetStaked         tlb.Coins                  // coins
	DistributedStables  tlb.Uint256                // uint256
	MaxPriceChange      tlb.Uint32                 // uint32
	NextItemIndex       tlb.Uint64                 // uint64
	Admin               tlb.InternalAddress        // address
	AssetAmountOnVault  tlb.Coins                  // coins
	StableAmountOnVault tlb.Coins                  // coins
	StaticData          tlb.RefT[*VaultStaticData] // Cell<VaultStaticData>
}
type CollectionData struct {
	Next_item_index tlb.Int257          // int
	Content         boc.Cell            // cell
	Owner           tlb.InternalAddress // address
}
type StakingData struct {
	TokensStaked       tlb.Coins   // coins
	DistributedStables tlb.Uint256 // uint256
	MaxPriceChange     tlb.Uint32  // uint32
}
type Balance struct {
	AssetAmount  tlb.Coins // coins
	StableAmount tlb.Coins // coins
}

const ( // errors
)

func DecodeGetCollectionData(stack *tlb.VmStack) (result CollectionData, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCollectionData = 0x1905B

func GetCollectionData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result CollectionData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCollectionData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCollectionData(&stack)
}

func DecodeGetNftAddressByIndex(stack *tlb.VmStack) (result tlb.InternalAddress, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetNftAddressByIndex = 0x167A3

func GetNftAddressByIndex(ctx context.Context, executor Executor, reqAccountID ton.AccountID, itemIndex tlb.Int257) (result tlb.InternalAddress, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(itemIndex)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetNftAddressByIndex, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetNftAddressByIndex(&stack)
}

func DecodeGetStakingData(stack *tlb.VmStack) (result StakingData, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetStakingData = 0x1A601

func GetStakingData(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result StakingData, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetStakingData, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetStakingData(&stack)
}

func DecodeGetBalance(stack *tlb.VmStack) (result Balance, err error) {
	if stack.Len() != 2 {
		err = fmt.Errorf("invalid stack size %d, expected 2", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetBalance = 0x1FD27

func GetBalance(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result Balance, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetBalance, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetBalance(&stack)
}

func Vault_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage VaultStorage, err error) {
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

type vaultImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewVault(executor Executor, storageExecutor StorageExecutor) Vault {
	return &vaultImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c vaultImpl) WithAccountId(accountID ton.AccountID) VaultWithAccount {
	return &vaultWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c vaultImpl) GetCollectionData(ctx context.Context, reqAccountID ton.AccountID) (CollectionData, error) {
	return GetCollectionData(ctx, c.executor, reqAccountID)
}

func (c vaultImpl) GetNftAddressByIndex(ctx context.Context, reqAccountID ton.AccountID, itemIndex tlb.Int257) (tlb.InternalAddress, error) {
	return GetNftAddressByIndex(ctx, c.executor, reqAccountID, itemIndex)
}

func (c vaultImpl) GetStakingData(ctx context.Context, reqAccountID ton.AccountID) (StakingData, error) {
	return GetStakingData(ctx, c.executor, reqAccountID)
}

func (c vaultImpl) GetBalance(ctx context.Context, reqAccountID ton.AccountID) (Balance, error) {
	return GetBalance(ctx, c.executor, reqAccountID)
}

func (c vaultImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, VaultStorage, error) {
	return Vault_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type vaultWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c vaultWithAccountImpl) GetCollectionData(ctx context.Context) (CollectionData, error) {
	return GetCollectionData(ctx, c.executor, c.accountID)
}

func (c vaultWithAccountImpl) GetNftAddressByIndex(ctx context.Context, itemIndex tlb.Int257) (tlb.InternalAddress, error) {
	return GetNftAddressByIndex(ctx, c.executor, c.accountID, itemIndex)
}

func (c vaultWithAccountImpl) GetStakingData(ctx context.Context) (StakingData, error) {
	return GetStakingData(ctx, c.executor, c.accountID)
}

func (c vaultWithAccountImpl) GetBalance(ctx context.Context) (Balance, error) {
	return GetBalance(ctx, c.executor, c.accountID)
}

func (c vaultWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, VaultStorage, error) {
	return Vault_AccountState(ctx, c.storageExecutor, c.accountID)
}
