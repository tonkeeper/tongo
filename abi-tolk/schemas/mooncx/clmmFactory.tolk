tolk 1.4

// moon.cx CLMM (bin-based concentrated liquidity) factory: ABI stub for tongo.
// It deploys one pool per (pair, bin step); the pool splits its bins between
// three range contracts that hold the liquidity. Admin messages go
// factory -> pool -> range.
// Deployed factory: EQD4uiDMK7T3KMVlcfIfp_LCe9eHtm6BIHCcHJ4stNQtb4R9
// (the older factory EQDugULuWQbVGoa2ChPgWiJSRPwr0DkU04NfcXacgmCOGbWR uses the
// same messages).
// getPoolAddress is left out: it takes arguments, so it cannot detect the
// factory.

contract mooncx_clmm_factory {
    incomingMessages: MooncxClmmFactoryMessage
    storage: null
}

type MooncxClmmFactoryMessage =
    | ClmmFactoryCreatePool
    | ClmmFactoryUpgradePool
    | ClmmFactoryUpdateRangeWallets
    | ClmmFactoryUpdateRangeFees
    | ClmmFactoryUpgradeRange
    | ClmmFactoryCollectProtocolFees
    | ClmmFactoryUpgrade

/// Creates a pool for a pair with a bin step (25, 50, 100 or 200) starting at
/// initBin; a null master stands for GRAM, and the factory orders the masters
/// itself. Needs at least 4 GRAM.
struct (0xede56ae4) ClmmFactoryCreatePool {
    queryId: uint64
    xMaster: address?
    yMaster: address?
    binStep: uint16
    initBin: int64
}

/// Admin replaces the code and/or data of a pool; the factory forwards it to the
/// pool as ClmmPoolUpgrade.
struct (0x4c76e24e) ClmmFactoryUpgradePool {
    queryId: uint64
    pool: address
    newCode: cell?
    newData: cell?
}

/// Admin sets the jetton wallets of range rangeNum of a pool (a null xWallet
/// means the x side is GRAM); forwarded as ClmmPoolUpdateRangeWallets.
struct (0xe55e6dc7) ClmmFactoryUpdateRangeWallets {
    queryId: uint64
    pool: address
    rangeNum: int16
    xWallet: address?
    yWallet: address
}

/// Admin sets the LP and protocol fees (basis points) of range rangeNum of a
/// pool; forwarded as ClmmPoolUpdateRangeFees.
struct (0x40073294) ClmmFactoryUpdateRangeFees {
    queryId: uint64
    pool: address
    rangeNum: int16
    lpFee: uint8
    protocolFee: uint8
}

/// Admin replaces the code and/or data of range rangeNum of a pool; forwarded
/// as ClmmPoolUpgradeRange.
struct (0x06dbf780) ClmmFactoryUpgradeRange {
    queryId: uint64
    pool: address
    rangeNum: int16
    newCode: cell?
    newData: cell?
}

/// Admin collects the protocol fees of range rangeNum of a pool; forwarded as
/// ClmmPoolCollectProtocolFees.
struct (0x4b3332cb) ClmmFactoryCollectProtocolFees {
    queryId: uint64
    pool: address
    rangeNum: int16
    recipient: address
}

/// Admin replaces the factory code and/or data.
struct (0x496938f7) ClmmFactoryUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

fun onInternalMessage(in: InMessage) {
}
