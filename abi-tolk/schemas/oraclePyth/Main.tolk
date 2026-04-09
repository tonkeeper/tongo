tolk 1.3

import "@stdlib/tvm-dicts"

const OP_UPDATE_GUARDIAN_SET: uint32 = 1
const OP_UPDATE_PRICE_FEEDS: uint32 = 2
const OP_EXECUTE_GOVERNANCE_ACTION: uint32 = 3
const OP_UPGRADE_CONTRACT: uint32 = 4
const OP_PARSE_PRICE_FEED_UPDATES: uint32 = 5
const OP_PARSE_UNIQUE_PRICE_FEED_UPDATES: uint32 = 6

const OP_RESPONSE_SUCCESS: uint32 = 0x10001
const OP_RESPONSE_ERROR: uint32 = 0x10002

contract Main {
    incomingMessages: MainMessage
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
    updateData: cell
    priceIds: cell
    minPublishTime: uint64
    maxPublishTime: uint64
    targetAddress: any_address
    customPayload: cell
}

struct (0x00000006) ParseUniquePriceFeedUpdatesMessage {
    updateData: cell
    priceIds: cell
    publishTime: uint64
    maxStaleness: uint64
    targetAddress: any_address
    customPayload: cell
}

fun onInternalMessage(inMsg: InMessage) {
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

get fun get_price_unsafe(price_feed_id: uint256): (int, int, int, int) {
    return (0, 0, 0, 0);
}

get fun get_price_no_older_than(time_period: int, price_feed_id: uint256): (int, int, int, int) {
    return (0, 0, 0, 0);
}

get fun get_ema_price_unsafe(price_feed_id: uint256): (int, int, int, int) {
    return (0, 0, 0, 0);
}

get fun get_ema_price_no_older_than(time_period: int, price_feed_id: uint256): (int, int, int, int) {
    return (0, 0, 0, 0);
}

get fun get_chain_id(): int {
    return 0;
}

get fun get_current_guardian_set_index(): int {
    return 0;
}

get fun get_guardian_set(index: int): (int, cell, int) {
    return (0, beginCell().endCell(), 0);
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
