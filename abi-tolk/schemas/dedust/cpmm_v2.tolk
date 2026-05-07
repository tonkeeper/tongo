contract dedust_cpmm_v2 {
incomingMessages: (),
storage: null,
forceAbiExport: []
}

fun onInternalMessage(in: InMessage) {

}

enum PoolStatus: uint2 {
    uninitialized = 0b00,
    initializing = 0b01,
    initialized = 0b10,
}

struct Q120X120 {
    value: varuint32
}

struct PoolReward {
    remainingTime: uint40 = 0;
    remainingBudget: coins = 0;
    rewardsPerToken: Q120X120;
    lastUpdate: uint40 = 0;
}

struct Asset {
    value: address?
}

struct DedustCpmmV2GetPoolData {
    status: PoolStatus
    depositActive: bool
    swapActive: bool
    assetX: Asset
    assetY: Asset
    walletsByAssets: map<uint256, address>
    assetsByWallets: map<uint256, address>
    walletsByResolutions: map<uint256, address>
    baseFeeBPS: uint16
    reserveX: coins
    reserveY: coins
    liquidity: coins
    protocolFeeX: coins
    protocolFeeY: coins
    creatorFeeX: coins
    creatorFeeY: coins
    xLPFeePerToken: Q120X120
    yLPFeePerToken: Q120X120
    rewards: map<uint2, PoolReward>?
}

get fun get_pool_data(): DedustCpmmV2GetPoolData {
    return DedustCpmmV2GetPoolData{
        status: PoolStatus.initialized,
        depositActive: false,
        swapActive: false,
        assetX: Asset{
            value: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")
        },
        assetY: Asset{
            value: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")
        },
        walletsByAssets: createEmptyMap(),
        assetsByWallets: createEmptyMap(),
        walletsByResolutions: createEmptyMap(),
        baseFeeBPS: 0,
        reserveX: 0,
        reserveY: 0,
        liquidity: 0,
        protocolFeeX: 0,
        protocolFeeY: 0,
        creatorFeeX: 0,
        creatorFeeY: 0,
        xLPFeePerToken: Q120X120 {
            value: 0
        },
        yLPFeePerToken: Q120X120 {
            value: 0
        },
        rewards: createEmptyMap()
    }
}
