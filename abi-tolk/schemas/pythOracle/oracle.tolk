tolk 1.2

import "@stdlib/tvm-dicts"
type SkippedBytes = uint8

fun SkippedBytes.unpackFromSlice(mutate s: slice): SkippedBytes {
    val byteCount = s.loadUint(8);
    s.skipBits(byteCount * 8);
    return byteCount;
}

fun SkippedBytes.packToBuilder(self, mutate b: builder) {
    b.storeAny(self);
    b.storeUint(0, self * 8);
}

type WormholeProofBytes = array<uint8>

fun WormholeProofBytes.unpackFromSlice(mutate s: slice): WormholeProofBytes {
    return [];
}

fun WormholeProofBytes.packToBuilder(self, mutate b: builder) {
}

contract pyth_oracle {
    incomingMessages: MainMessage,
    storage: MainStorage
    forceAbiExport: (
        MainStorage,
        PriceFeedsSection,
        StoredPriceFeed,
        PricePoint,
        DataSourcesSection,
        DataSource,
        GuardianSetsSection,
        GuardianSetRecord,
        GovernanceSection,
        AccumulatorUpdateHeader,
        PriceFeedUpdateData,
        PriceFeedIdList,
        PriceFeedUpdateResponse,
        PriceFeedResponseEntry,
        ErrorResponse,
        SuccessResponse,
        UpdateGuardianSetMessage,
        UpdatePriceFeedsMessage,
        ExecuteGovernanceActionMessage,
        UpgradeContractMessage,
        ParsePriceFeedUpdatesMessage,
        ParseUniquePriceFeedUpdatesMessage
    )
}

type MainMessage =
    UpdateGuardianSetMessage |
    UpdatePriceFeedsMessage |
    ExecuteGovernanceActionMessage |
    UpgradeContractMessage |
    ParsePriceFeedUpdatesMessage |
    ParseUniquePriceFeedUpdatesMessage

struct MainStorage {
    priceFeeds: Cell<PriceFeedsSection>
    dataSources: Cell<DataSourcesSection>
    guardianSets: Cell<GuardianSetsSection>
    governance: Cell<GovernanceSection>
}

struct PriceFeedsSection {
    latestPriceFeeds: map<uint256, Cell<StoredPriceFeed>>
    singleUpdateFee: uint256
}

struct StoredPriceFeed {
    price: Cell<PricePoint>
    emaPrice: Cell<PricePoint>
}

struct PricePoint {
    price: int64
    conf: uint64
    expo: int32
    publishTime: uint64
}

struct AccumulatorUpdateHeader {
    magic: uint32
    majorVersion: uint8
    minorVersion: uint8
    trailingHeader: SkippedBytes
    updateType: uint8
}

struct PriceFeedUpdateData {
    header: AccumulatorUpdateHeader
    wormholeProof: WormholeProofBytes
    numUpdates: uint8
}

type PriceFeedIdList = array<uint256>

fun PriceFeedIdList.unpackFromSlice(mutate s: slice): PriceFeedIdList {
    return []
}

struct DataSourcesSection {
    isValidDataSource: map<uint256, bool>
}

struct DataSource {
    emitterChainId: uint16
    emitterAddress: uint256
}

struct GuardianSetsSection {
    currentGuardianSetIndex: uint32
    guardianSets: map<uint32, Cell<GuardianSetRecord>>
}

struct GuardianSetRecord {
    expirationTime: uint64
    guardianKeys: map<uint8, bits160>
}

struct GuardianSetInfo {
    expirationTime: int
    keysDict: cell
    keyCount: int
}

struct GovernanceSection {
    chainId: uint16
    governanceChainId: uint16
    governanceContract: uint256
    consumedGovernanceActions: map<uint256, bool>
    governanceDataSource: Cell<DataSource>
    lastExecutedGovernanceSequence: uint64
    governanceDataSourceIndex: uint32
    upgradeCodeHash: uint256
}

struct PriceFeedResponseEntry {
    priceId: uint256
    priceFeed: Cell<StoredPriceFeed>
    next: Cell<PriceFeedResponseEntry>?
}

struct PriceFeedUpdateResponse {
    op: uint32
    priceFeedCount: uint8
    priceFeeds: Cell<PriceFeedResponseEntry>
    originalSender: any_address
    customPayload: Cell<slice>
}

struct (0x10002) ErrorResponse {
    errorCode: uint32
    operation: uint32
    customPayload: Cell<slice>
}

struct (0x10001) SuccessResponse {
    result: cell
    customPayload: cell
}

struct (0x00000001) UpdateGuardianSetMessage {
    wormholeMessage: cell
}

struct (0x00000002) UpdatePriceFeedsMessage {
    updateData: cell
}

struct (0x00000003) ExecuteGovernanceActionMessage {
    governanceVm: cell
}

struct (0x00000004) UpgradeContractMessage {
    newCode: cell
}

struct (0x00000005) ParsePriceFeedUpdatesMessage {
    updateData: Cell<PriceFeedUpdateData>
    priceIds: Cell<PriceFeedIdList>
    minPublishTime: uint64
    maxPublishTime: uint64
    targetAddress: any_address
    customPayload: cell
}

struct (0x00000006) ParseUniquePriceFeedUpdatesMessage {
    updateData: Cell<PriceFeedUpdateData>
    priceIds: Cell<PriceFeedIdList>
    publishTime: uint64
    maxStaleness: uint64
    targetAddress: any_address
    customPayload: cell
}

struct PriceData {
    price: int64
    conf: uint64
    expo: int32
    timestamp: uint64
}

struct PriceFeesCell {
    assetID: uint256
    priceData: Cell<Cell<PriceData>> // здесь нет ошибки, там есть другая информация, но она бесполезна
}

struct (0x00000005) OracleResponseSuccess {
    someNum: uint8
    priceFeedsCell: Cell<PriceFeesCell>
    initialSender: address
    AfterOperation: cell
}

fun onInternalMessage(inMsg: InMessage) {

    createMessage({
        bounce: BounceMode.NoBounce,
        dest: inMsg.senderAddress,
        value: 0,
        body: OracleResponseSuccess{
            someNum: 0,
            priceFeedsCell: createEmptyCell() as Cell<PriceFeesCell>,
            initialSender: inMsg.senderAddress,
            AfterOperation: createEmptyCell(),
        }
    }).send(SEND_MODE_IGNORE_ERRORS + SEND_MODE_CARRY_ALL_REMAINING_MESSAGE_VALUE);
}

fun onBouncedMessage(msg: InMessageBounced) {
}

get fun get_update_fee(data: slice): int {
    return 0;
}

get fun get_single_update_fee(): int {
    return 0;
}

get fun get_governance_data_source_index(): int {
    return 0;
}

get fun get_governance_data_source(): cell {
    return beginCell().endCell();
}

get fun get_last_executed_governance_sequence(): int {
    return 0;
}

get fun get_is_valid_data_source(data_source: cell): int {
    return 0;
}

get fun get_price_unsafe(price_feed_id: uint256): PricePoint {
    return PricePoint { price: 0, conf: 0, expo: 0, publishTime: 0 };
}

get fun get_price_no_older_than(time_period: int, price_feed_id: uint256): PricePoint {
    return PricePoint { price: 0, conf: 0, expo: 0, publishTime: 0 };
}

get fun get_ema_price_unsafe(price_feed_id: uint256): PricePoint {
    return PricePoint { price: 0, conf: 0, expo: 0, publishTime: 0 };
}

get fun get_ema_price_no_older_than(time_period: int, price_feed_id: uint256): PricePoint {
    return PricePoint { price: 0, conf: 0, expo: 0, publishTime: 0 };
}

get fun get_chain_id(): int {
    return 0;
}

get fun get_current_guardian_set_index(): int {
    return 0;
}

get fun get_guardian_set(index: int): GuardianSetInfo {
    return GuardianSetInfo { expirationTime: 0, keysDict: beginCell().endCell(), keyCount: 0 };
}

get fun get_governance_chain_id(): int {
    return 0;
}

get fun get_governance_contract(): int {
    return 0;
}

get fun governance_action_is_consumed(hash: int): int {
    return 0;
}
