

contract ff_vault {
    storage: VaultStorage
    incomingMessages: VaultMessage,
    forceAbiExport: (
    )
}

struct (0x7362d09c) JettonTransferNotification<T> {
    payload: T
}

type VaultMessage =
    JettonTransferNotification<AssetDeposit | StableDeposit> |
    StakeOperation |
    UnstakeExecute |
    UnstakeOperation |
    UnstakeExecuteInternalCallback |
    WithdrawJetton;

struct (0xe44c1944) WithdrawJetton {
    queryId: uint64
    jettonWallet: address
    toAddress: address
    amount: coins
}

struct (0x3d669db7) StableDeposit {
    distributeBetweenStakers: bool
}

struct (0x4e835078) UnstakeExecute {
    queryId: uint64
    priceUpdateData: cell
    amountToUnstake: coins
    itemIndex: uint64
}

struct (0x534917c3) AssetDeposit {
    priceUpdateData: cell?
}

struct (0x7a294fdc) StakeOperation {
    queryId: uint64
    depositAddress: address // адрес приславшего токены
    assetAmount: coins // количество присланных токенов
}

struct (0x97004d72) UnstakeOperation {
    queryId: uint64
    amountToUnstake: coins
    itemIndex: uint64
}

struct PositionData {
    xautPriceStaked: uint256
    previousStablePerAssetCollected: uint256
}

struct (0x61755200) UnstakeExecuteInternalCallback {
    queryId: uint64
    itemIndex: uint64
    amountToUnstake: coins
    xautPriceActual: uint256
    owner: address
    positionData: Cell<PositionData>
}


fun onInternalMessage(in: InMessage) {
    throw 0;
}

struct VaultStaticData {
    PYTHOracle: address
    AssetJettonWallet: address
    StableJettonWallet: address
    nftItemCode: cell
    stableDecimals: uint8
    assetDecimals: uint8
}

struct (0x00) VaultStorage {
    assetStaked: coins
    distributedStables: uint256
    maxPriceChange: uint32
    nextItemIndex: uint64
    admin: address
    assetAmountOnVault: coins
    stableAmountOnVault: coins
    staticData: Cell<VaultStaticData>
}

struct CollectionData {
    next_item_index: int
    content: cell
    owner: address
}

get fun get_collection_data(): CollectionData {
    throw 0;
}

struct StakingData {
    tokensStaked: coins // возвращает количество токенов в стейкинге
    distributedStables: uint256 // награды делить на количество токенов
    maxPriceChange: uint32
}

get fun get_nft_address_by_index(itemIndex: int): address {
    throw 0;
}
get fun get_staking_data(): StakingData {
    throw 0;
}

struct Balance {
    assetAmount: coins
    stableAmount: coins
}

get fun get_balance(): Balance {
    throw 0;
}
