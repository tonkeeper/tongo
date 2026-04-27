// Code generated - DO NOT EDIT.

package abiPythOracle

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const PrefixUpdateGuardianSetMessage uint64 = 0x00000001

type UpdateGuardianSetMessage struct {
	WormholeMessage boc.Cell // cell
}

const PrefixUpdatePriceFeedsMessage uint64 = 0x00000002

type UpdatePriceFeedsMessage struct {
	UpdateData boc.Cell // cell
}

const PrefixExecuteGovernanceActionMessage uint64 = 0x00000003

type ExecuteGovernanceActionMessage struct {
	GovernanceVm boc.Cell // cell
}

const PrefixUpgradeContractMessage uint64 = 0x00000004

type UpgradeContractMessage struct {
	NewCode boc.Cell // cell
}
type TrailingHeaderBytes []tlb.Uint8
type WormholeProofBytes []tlb.Uint8
type PriceFeedMessage struct {
	MessageType     tlb.Uint8   // uint8
	PriceId         tlb.Uint256 // uint256
	Price           tlb.Int64   // int64
	Conf            tlb.Uint64  // uint64
	Expo            tlb.Int32   // int32
	PublishTime     tlb.Uint64  // uint64
	PrevPublishTime tlb.Uint64  // uint64
	EmaPrice        tlb.Int64   // int64
	EmaConf         tlb.Uint64  // uint64
}
type PriceFeedProof []tlb.Bits160
type PriceFeedUpdate struct {
	Message PriceFeedMessage // PriceFeedMessage
	Proof   PriceFeedProof   // PriceFeedProof
}
type WormholeUpdateSection struct {
	Proof   WormholeProofBytes // WormholeProofBytes
	Updates []PriceFeedUpdate  // array<PriceFeedUpdate>
}
type AccumulatorUpdatePayload struct {
	TrailingHeader  TrailingHeaderBytes   // TrailingHeaderBytes
	UpdateType      tlb.Uint8             // uint8
	WormholeSection WormholeUpdateSection // WormholeUpdateSection
}
type PriceFeedUpdateData struct {
	Magic        tlb.Uint32               // uint32
	MajorVersion tlb.Uint8                // uint8
	MinorVersion tlb.Uint8                // uint8
	Accumulator  AccumulatorUpdatePayload // AccumulatorUpdatePayload
}
type PriceFeedIdList []tlb.Uint256

const PrefixParsePriceFeedUpdatesMessage uint64 = 0x00000005

type ParsePriceFeedUpdatesMessage struct {
	UpdateData     tlb.RefT[*PriceFeedUpdateData] // Cell<PriceFeedUpdateData>
	PriceIds       tlb.RefT[*PriceFeedIdList]     // Cell<PriceFeedIdList>
	MinPublishTime tlb.Uint64                     // uint64
	MaxPublishTime tlb.Uint64                     // uint64
	TargetAddress  tlb.MsgAddress                 // any_address
	CustomPayload  boc.Cell                       // cell
}

const PrefixParseUniquePriceFeedUpdatesMessage uint64 = 0x00000006

type ParseUniquePriceFeedUpdatesMessage struct {
	UpdateData    tlb.RefT[*PriceFeedUpdateData] // Cell<PriceFeedUpdateData>
	PriceIds      tlb.RefT[*PriceFeedIdList]     // Cell<PriceFeedIdList>
	PublishTime   tlb.Uint64                     // uint64
	MaxStaleness  tlb.Uint64                     // uint64
	TargetAddress tlb.MsgAddress                 // any_address
	CustomPayload boc.Cell                       // cell
}
type PricePoint struct {
	Price       tlb.Int64  // int64
	Conf        tlb.Uint64 // uint64
	Expo        tlb.Int32  // int32
	PublishTime tlb.Uint64 // uint64
}
type StoredPriceFeed struct {
	Price    tlb.RefT[*PricePoint] // Cell<PricePoint>
	EmaPrice tlb.RefT[*PricePoint] // Cell<PricePoint>
}
type PriceFeedsSection struct {
	LatestPriceFeeds tlb.HashmapE[tlb.Uint256, tlb.RefT[*StoredPriceFeed]] // map<uint256, Cell<StoredPriceFeed>>
	SingleUpdateFee  tlb.Uint256                                           // uint256
}
type DataSourcesSection struct {
	IsValidDataSource tlb.HashmapE[tlb.Uint256, bool] // map<uint256, bool>
}
type GuardianSetRecord struct {
	ExpirationTime tlb.Uint64                           // uint64
	GuardianKeys   tlb.HashmapE[tlb.Uint8, tlb.Bits160] // map<uint8, bits160>
}
type GuardianSetsSection struct {
	CurrentGuardianSetIndex tlb.Uint32                                             // uint32
	GuardianSets            tlb.HashmapE[tlb.Uint32, tlb.RefT[*GuardianSetRecord]] // map<uint32, Cell<GuardianSetRecord>>
}
type DataSource struct {
	EmitterChainId tlb.Uint16  // uint16
	EmitterAddress tlb.Uint256 // uint256
}
type GovernanceSection struct {
	ChainId                        tlb.Uint16                      // uint16
	GovernanceChainId              tlb.Uint16                      // uint16
	GovernanceContract             tlb.Uint256                     // uint256
	ConsumedGovernanceActions      tlb.HashmapE[tlb.Uint256, bool] // map<uint256, bool>
	GovernanceDataSource           tlb.RefT[*DataSource]           // Cell<DataSource>
	LastExecutedGovernanceSequence tlb.Uint64                      // uint64
	GovernanceDataSourceIndex      tlb.Uint32                      // uint32
	UpgradeCodeHash                tlb.Uint256                     // uint256
}
type MainStorage struct {
	PriceFeeds   tlb.RefT[*PriceFeedsSection]   // Cell<PriceFeedsSection>
	DataSources  tlb.RefT[*DataSourcesSection]  // Cell<DataSourcesSection>
	GuardianSets tlb.RefT[*GuardianSetsSection] // Cell<GuardianSetsSection>
	Governance   tlb.RefT[*GovernanceSection]   // Cell<GovernanceSection>
}
type PriceFeedResponseEntry struct {
	PriceId   tlb.Uint256                                  // uint256
	PriceFeed tlb.RefT[*StoredPriceFeed]                   // Cell<StoredPriceFeed>
	Next      tlb.Maybe[tlb.RefT[*PriceFeedResponseEntry]] // Cell<PriceFeedResponseEntry>?
}
type PriceFeedUpdateResponse struct {
	Op             tlb.Uint32                        // uint32
	PriceFeedCount tlb.Uint8                         // uint8
	PriceFeeds     tlb.RefT[*PriceFeedResponseEntry] // Cell<PriceFeedResponseEntry>
	OriginalSender tlb.MsgAddress                    // any_address
	CustomPayload  boc.Cell                          // Cell<slice>
}

const PrefixErrorResponse uint64 = 0x10002

type ErrorResponse struct {
	ErrorCode     tlb.Uint32 // uint32
	Operation     tlb.Uint32 // uint32
	CustomPayload boc.Cell   // Cell<slice>
}

const PrefixSuccessResponse uint64 = 0x10001

type SuccessResponse struct {
	Result        boc.Cell // cell
	CustomPayload boc.Cell // cell
}
type PriceData struct {
	Price     tlb.Int64  // int64
	Conf      tlb.Uint64 // uint64
	Expo      tlb.Int32  // int32
	Timestamp tlb.Uint64 // uint64
}
type PriceFeesCell struct {
	AssetID   tlb.Uint256                     // uint256
	PriceData tlb.RefT[*tlb.RefT[*PriceData]] // Cell<Cell<PriceData>>
}

const PrefixOracleResponseSuccess uint64 = 0x00000005

type OracleResponseSuccess struct {
	SomeNum        tlb.Uint8                // uint8
	PriceFeedsCell tlb.RefT[*PriceFeesCell] // Cell<PriceFeesCell>
	InitialSender  tlb.InternalAddress      // address
	AfterOperation boc.Cell                 // cell
}
type GuardianSetInfo struct {
	ExpirationTime tlb.Int257 // int
	KeysDict       boc.Cell   // cell
	KeyCount       tlb.Int257 // int
}

func DecodeGetUpdateFee(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetUpdateFee = 0x1E44F

func GetUpdateFee(ctx context.Context, executor Executor, reqAccountID ton.AccountID, data boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.CellToVmCellSlice(&data)
		if err != nil {
			err = fmt.Errorf("encode param data: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetUpdateFee, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetUpdateFee(&stack)
}

func DecodeGetSingleUpdateFee(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetSingleUpdateFee = 0x18673

func GetSingleUpdateFee(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetSingleUpdateFee, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetSingleUpdateFee(&stack)
}

func DecodeGetGovernanceDataSourceIndex(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetGovernanceDataSourceIndex = 0x17FBE

func GetGovernanceDataSourceIndex(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGovernanceDataSourceIndex, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGovernanceDataSourceIndex(&stack)
}

func DecodeGetGovernanceDataSource(stack *tlb.VmStack) (result boc.Cell, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	return stack.ReadCell()
}

const MethodIDGetGovernanceDataSource = 0x1B157

func GetGovernanceDataSource(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result boc.Cell, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGovernanceDataSource, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGovernanceDataSource(&stack)
}

func DecodeGetLastExecutedGovernanceSequence(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetLastExecutedGovernanceSequence = 0x11234

func GetLastExecutedGovernanceSequence(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetLastExecutedGovernanceSequence, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetLastExecutedGovernanceSequence(&stack)
}

func DecodeGetIsValidDataSource(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetIsValidDataSource = 0x18D00

func GetIsValidDataSource(ctx context.Context, executor Executor, reqAccountID ton.AccountID, dataSource boc.Cell) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	{
		var val tlb.VmStackValue
		val, err = tlb.TlbStructToVmCell(dataSource)
		if err != nil {
			err = fmt.Errorf("encode param dataSource: %w", err)
			return
		}
		stack.Put(val)
	}
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetIsValidDataSource, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetIsValidDataSource(&stack)
}

func DecodeGetPriceUnsafe(stack *tlb.VmStack) (result PricePoint, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPriceUnsafe = 0x199E1

func GetPriceUnsafe(ctx context.Context, executor Executor, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (result PricePoint, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(priceFeedId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPriceUnsafe, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPriceUnsafe(&stack)
}

func DecodeGetPriceNoOlderThan(stack *tlb.VmStack) (result PricePoint, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetPriceNoOlderThan = 0x11B73

func GetPriceNoOlderThan(ctx context.Context, executor Executor, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (result PricePoint, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(timePeriod)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(priceFeedId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetPriceNoOlderThan, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetPriceNoOlderThan(&stack)
}

func DecodeGetEmaPriceUnsafe(stack *tlb.VmStack) (result PricePoint, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetEmaPriceUnsafe = 0x1F541

func GetEmaPriceUnsafe(ctx context.Context, executor Executor, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (result PricePoint, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(priceFeedId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetEmaPriceUnsafe, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetEmaPriceUnsafe(&stack)
}

func DecodeGetEmaPriceNoOlderThan(stack *tlb.VmStack) (result PricePoint, err error) {
	if stack.Len() != 4 {
		err = fmt.Errorf("invalid stack size %d, expected 4", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetEmaPriceNoOlderThan = 0x11B2A

func GetEmaPriceNoOlderThan(ctx context.Context, executor Executor, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (result PricePoint, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(timePeriod)})
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(priceFeedId)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetEmaPriceNoOlderThan, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetEmaPriceNoOlderThan(&stack)
}

func DecodeGetChainId(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetChainId = 0x1E048

func GetChainId(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetChainId, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetChainId(&stack)
}

func DecodeGetCurrentGuardianSetIndex(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetCurrentGuardianSetIndex = 0x1BFC4

func GetCurrentGuardianSetIndex(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetCurrentGuardianSetIndex, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetCurrentGuardianSetIndex(&stack)
}

func DecodeGetGuardianSet(stack *tlb.VmStack) (result GuardianSetInfo, err error) {
	if stack.Len() != 3 {
		err = fmt.Errorf("invalid stack size %d, expected 3", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetGuardianSet = 0x1DCB3

func GetGuardianSet(ctx context.Context, executor Executor, reqAccountID ton.AccountID, index tlb.Int257) (result GuardianSetInfo, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(index)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGuardianSet, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGuardianSet(&stack)
}

func DecodeGetGovernanceChainId(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetGovernanceChainId = 0x18F9E

func GetGovernanceChainId(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGovernanceChainId, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGovernanceChainId(&stack)
}

func DecodeGetGovernanceContract(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetGovernanceContract = 0x10132

func GetGovernanceContract(ctx context.Context, executor Executor, reqAccountID ton.AccountID) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGovernanceContract, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGovernanceContract(&stack)
}

func DecodeGetGovernanceActionIsConsumed(stack *tlb.VmStack) (result tlb.Int257, err error) {
	if stack.Len() != 1 {
		err = fmt.Errorf("invalid stack size %d, expected 1", stack.Len())
		return
	}
	err = result.ReadFromStack(stack)
	return
}

const MethodIDGetGovernanceActionIsConsumed = 0x13721

func GetGovernanceActionIsConsumed(ctx context.Context, executor Executor, reqAccountID ton.AccountID, hash tlb.Int257) (result tlb.Int257, err error) {
	var errCode uint32
	var stack tlb.VmStack
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(hash)})
	errCode, stack, err = executor.RunSmcMethodByID(ctx, reqAccountID, MethodIDGetGovernanceActionIsConsumed, stack)
	if err != nil {
		return
	}
	if errCode != 0 && errCode != 1 {
		err = fmt.Errorf("method execution failed with code: %v", errCode)
		return
	}
	return DecodeGetGovernanceActionIsConsumed(&stack)
}

func Oracle_AccountState(ctx context.Context, executor StorageExecutor, accountID ton.AccountID) (sa tlb.ShardAccount, storage MainStorage, err error) {
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

type oracleImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
}

func NewOracle(executor Executor, storageExecutor StorageExecutor) Oracle {
	return &oracleImpl{executor: executor, storageExecutor: storageExecutor}
}

func (c oracleImpl) WithAccountId(accountID ton.AccountID) OracleWithAccount {
	return &oracleWithAccountImpl{executor: c.executor, storageExecutor: c.storageExecutor, accountID: accountID}
}

func (c oracleImpl) GetUpdateFee(ctx context.Context, reqAccountID ton.AccountID, data boc.Cell) (tlb.Int257, error) {
	return GetUpdateFee(ctx, c.executor, reqAccountID, data)
}

func (c oracleImpl) GetSingleUpdateFee(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetSingleUpdateFee(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetGovernanceDataSourceIndex(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetGovernanceDataSourceIndex(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetGovernanceDataSource(ctx context.Context, reqAccountID ton.AccountID) (boc.Cell, error) {
	return GetGovernanceDataSource(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetLastExecutedGovernanceSequence(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetLastExecutedGovernanceSequence(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetIsValidDataSource(ctx context.Context, reqAccountID ton.AccountID, dataSource boc.Cell) (tlb.Int257, error) {
	return GetIsValidDataSource(ctx, c.executor, reqAccountID, dataSource)
}

func (c oracleImpl) GetPriceUnsafe(ctx context.Context, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetPriceUnsafe(ctx, c.executor, reqAccountID, priceFeedId)
}

func (c oracleImpl) GetPriceNoOlderThan(ctx context.Context, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetPriceNoOlderThan(ctx, c.executor, reqAccountID, timePeriod, priceFeedId)
}

func (c oracleImpl) GetEmaPriceUnsafe(ctx context.Context, reqAccountID ton.AccountID, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetEmaPriceUnsafe(ctx, c.executor, reqAccountID, priceFeedId)
}

func (c oracleImpl) GetEmaPriceNoOlderThan(ctx context.Context, reqAccountID ton.AccountID, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetEmaPriceNoOlderThan(ctx, c.executor, reqAccountID, timePeriod, priceFeedId)
}

func (c oracleImpl) GetChainId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetChainId(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetCurrentGuardianSetIndex(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetCurrentGuardianSetIndex(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetGuardianSet(ctx context.Context, reqAccountID ton.AccountID, index tlb.Int257) (GuardianSetInfo, error) {
	return GetGuardianSet(ctx, c.executor, reqAccountID, index)
}

func (c oracleImpl) GetGovernanceChainId(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetGovernanceChainId(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetGovernanceContract(ctx context.Context, reqAccountID ton.AccountID) (tlb.Int257, error) {
	return GetGovernanceContract(ctx, c.executor, reqAccountID)
}

func (c oracleImpl) GetGovernanceActionIsConsumed(ctx context.Context, reqAccountID ton.AccountID, hash tlb.Int257) (tlb.Int257, error) {
	return GetGovernanceActionIsConsumed(ctx, c.executor, reqAccountID, hash)
}

func (c oracleImpl) AccountState(ctx context.Context, reqAccountID ton.AccountID) (tlb.ShardAccount, MainStorage, error) {
	return Oracle_AccountState(ctx, c.storageExecutor, reqAccountID)
}

type oracleWithAccountImpl struct {
	executor        Executor
	storageExecutor StorageExecutor
	accountID       ton.AccountID
}

func (c oracleWithAccountImpl) GetUpdateFee(ctx context.Context, data boc.Cell) (tlb.Int257, error) {
	return GetUpdateFee(ctx, c.executor, c.accountID, data)
}

func (c oracleWithAccountImpl) GetSingleUpdateFee(ctx context.Context) (tlb.Int257, error) {
	return GetSingleUpdateFee(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetGovernanceDataSourceIndex(ctx context.Context) (tlb.Int257, error) {
	return GetGovernanceDataSourceIndex(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetGovernanceDataSource(ctx context.Context) (boc.Cell, error) {
	return GetGovernanceDataSource(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetLastExecutedGovernanceSequence(ctx context.Context) (tlb.Int257, error) {
	return GetLastExecutedGovernanceSequence(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetIsValidDataSource(ctx context.Context, dataSource boc.Cell) (tlb.Int257, error) {
	return GetIsValidDataSource(ctx, c.executor, c.accountID, dataSource)
}

func (c oracleWithAccountImpl) GetPriceUnsafe(ctx context.Context, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetPriceUnsafe(ctx, c.executor, c.accountID, priceFeedId)
}

func (c oracleWithAccountImpl) GetPriceNoOlderThan(ctx context.Context, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetPriceNoOlderThan(ctx, c.executor, c.accountID, timePeriod, priceFeedId)
}

func (c oracleWithAccountImpl) GetEmaPriceUnsafe(ctx context.Context, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetEmaPriceUnsafe(ctx, c.executor, c.accountID, priceFeedId)
}

func (c oracleWithAccountImpl) GetEmaPriceNoOlderThan(ctx context.Context, timePeriod tlb.Int257, priceFeedId tlb.Uint256) (PricePoint, error) {
	return GetEmaPriceNoOlderThan(ctx, c.executor, c.accountID, timePeriod, priceFeedId)
}

func (c oracleWithAccountImpl) GetChainId(ctx context.Context) (tlb.Int257, error) {
	return GetChainId(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetCurrentGuardianSetIndex(ctx context.Context) (tlb.Int257, error) {
	return GetCurrentGuardianSetIndex(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetGuardianSet(ctx context.Context, index tlb.Int257) (GuardianSetInfo, error) {
	return GetGuardianSet(ctx, c.executor, c.accountID, index)
}

func (c oracleWithAccountImpl) GetGovernanceChainId(ctx context.Context) (tlb.Int257, error) {
	return GetGovernanceChainId(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetGovernanceContract(ctx context.Context) (tlb.Int257, error) {
	return GetGovernanceContract(ctx, c.executor, c.accountID)
}

func (c oracleWithAccountImpl) GetGovernanceActionIsConsumed(ctx context.Context, hash tlb.Int257) (tlb.Int257, error) {
	return GetGovernanceActionIsConsumed(ctx, c.executor, c.accountID, hash)
}

func (c oracleWithAccountImpl) AccountState(ctx context.Context) (tlb.ShardAccount, MainStorage, error) {
	return Oracle_AccountState(ctx, c.storageExecutor, c.accountID)
}
