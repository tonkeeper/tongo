
contract ff_vault_position {
    incomingMessages: PositionMessage,
    forceAbiExport: (
    )
}

type PositionMessage = UnstakeRequest | UnstakeExecuteInternal | UnstakeExecuteCancel;

struct (0x8a77cf39) UnstakeRequest {
    queryId: uint64
    amountWantedToUnstake: coins
}

struct (0xd1029c99) UnstakeExecuteInternal {
    queryId: uint64
    amountToUnstake: coins
    xautPriceActual: uint256
}

struct (0x0b30393c) UnstakeExecuteCancel {
    queryId: uint64
    excessGetter: address
    amountToAdd: coins
}


struct NftData {
    init: bool
    index: int
    collection_address: address
    owner_address: address?
    content: cell?
}

struct StakePositionInfo {
    lockedAssetAmount: int,
    amountWantedToUnstake: int,
    previousStablePerAssetCollected: int,
    priceStake: int
}

fun onInternalMessage(in: PositionMessage) {
    throw 0;
}

get fun get_nft_data(): NftData {
    throw 0;
}
get fun get_stake_position_info(): StakePositionInfo {
    throw 0;
}
