tolk 1.4

// moon.cx CLMM pool: ABI stub for tongo. The pool deploys its three ranges
// (rangeNum -1, 0 and 1), keeps their codes and relays admin messages to them;
// users swap and add liquidity on the ranges (clmmRange.tolk).
// Example: EQDOAVaKQsb1CY_CT2RvalUlntZjrIF3cjpVFwG-VzjliWZ2 (GRAM/USDT), deployed
// by the factory EQD4uiDMK7T3KMVlcfIfp_LCe9eHtm6BIHCcHJ4stNQtb4R9.
// getPool and getRanges are declared so the pool can be detected; getRangeByNum
// takes an argument and is left out.

contract mooncx_clmm_pool {
    incomingMessages: MooncxClmmPoolMessage
    storage: null
}

type MooncxClmmPoolMessage =
    | ClmmPoolDeploy
    | ClmmPoolCollectProtocolFees
    | ClmmPoolUpdateRangeWallets
    | ClmmPoolUpdateRangeFees
    | ClmmPoolUpgradeRange
    | ClmmPoolUpgrade

struct ClmmPayloadCode {
    rangeCode: cell
    lpAccCode: cell
    lpPositionCode: cell
}

/// Factory deploys the pool (with StateInit) with the start bin, the LP and
/// protocol fees (basis points) and the codes of the range, LP account and LP
/// position contracts; the pool deploys its three ranges with ClmmRangeDeploy.
struct (0x6d33fdcd) ClmmPoolDeploy {
    queryId: uint64
    initBin: int64
    lpFee: uint8
    protocolFee: uint8
    payloadCode: Cell<ClmmPayloadCode>
}

/// Factory forwards ClmmFactoryCollectProtocolFees; the pool passes it to range
/// rangeNum as ClmmRangeCollectProtocolFees.
struct (0xefef86da) ClmmPoolCollectProtocolFees {
    queryId: uint64
    rangeNum: int16
    recipient: address
}

/// Factory forwards ClmmFactoryUpdateRangeWallets; the pool passes it to range
/// rangeNum as ClmmRangeUpdateWallets.
struct (0xe82e11a1) ClmmPoolUpdateRangeWallets {
    queryId: uint64
    rangeNum: int16
    xWallet: address?
    yWallet: address
}

/// Factory forwards ClmmFactoryUpdateRangeFees; the pool passes it to range
/// rangeNum as ClmmRangeUpdateFees.
struct (0xd97a5c35) ClmmPoolUpdateRangeFees {
    queryId: uint64
    rangeNum: int16
    lpFee: uint8
    protocolFee: uint8
}

/// Factory forwards ClmmFactoryUpgradeRange; the pool passes it to range
/// rangeNum as ClmmRangeUpgrade.
struct (0x7791fc81) ClmmPoolUpgradeRange {
    queryId: uint64
    rangeNum: int16
    newCode: cell?
    newData: cell?
}

/// Factory forwards ClmmFactoryUpgradePool: the pool replaces its code and/or data.
struct (0xb7f86a0a) ClmmPoolUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

/// Contract data returned by getPool.
struct ClmmPoolStorage {
    factory: address
    xMaster: address?
    yMaster: address
    binStep: uint16
    lpFee: uint8
    protocolFee: uint8
    rangeCode: cell
    lpAccCode: cell
    lpPositionCode: cell
}

fun onInternalMessage(in: InMessage) {
}

get fun getPool(): Cell<ClmmPoolStorage> {
    throw 0;
}

/// Addresses of the pool's ranges by rangeNum (-1, 0 and 1).
get fun getRanges(): map<int16, address> {
    throw 0;
}
