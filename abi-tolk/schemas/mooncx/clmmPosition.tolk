tolk 1.4

// moon.cx CLMM LP position: ABI stub for tongo. Every successful deposit deploys
// a new position (range, nftIndex) that holds the liquidity per bin; it is not
// an NFT.
// Its messages to the range use rich bounces; tongo already decodes those as BounceV2.
// getPosition is left out: the CPMM LP position has a get method with the same
// name and a different result, and the Go code of both shares one package.

contract mooncx_clmm_position {
    incomingMessages: MooncxClmmPositionMessage
    storage: null
    forceAbiExport: [
        ClmmPositionClosed
    ]
}

type MooncxClmmPositionMessage =
    | ClmmPositionDeploy
    | ClmmPositionWithdrawLiquidity
    | ClmmPositionWithdrawFees
    | ClmmPositionWithdrawFeesUpdate
    | ClmmPositionWithdrawLiquidityUpdate
    | ClmmPositionLockLiquidity

struct ClmmBin {
    l: uint256
    xFee: uint256
    yFee: uint256
}

/// Range deploys a new position (with StateInit) for provider with the bins of a
/// deposit (bin index -> liquidity and fee checkpoints).
struct (0x4e5504a9) ClmmPositionDeploy {
    queryId: uint64
    provider: address
    positionBins: map<int64, ClmmBin>
}

/// Owner withdraws the liquidity of the given bins (only the keys are used); the
/// position sends ClmmRangeWithdrawLiquidity to its range. Refused while the
/// position is locked or waiting for the range.
struct (0xc84175bb) ClmmPositionWithdrawLiquidity {
    queryId: uint64
    withdrawBins: map<int64, bool>
}

/// Owner collects the fees of all bins; the position sends ClmmRangeWithdrawFees
/// to its range.
struct (0x716c13d8) ClmmPositionWithdrawFees {
    queryId: uint64
}

/// Range's answer to ClmmRangeWithdrawFees: the bins with new fee checkpoints.
struct (0xfe49ce3f) ClmmPositionWithdrawFeesUpdate {
    queryId: uint64
    updatedBins: map<int64, ClmmBin>
}

/// Range's answer to ClmmRangeWithdrawLiquidity. A position left without bins
/// then returns its balance to the owner with ClmmPositionClosed.
struct (0x18256917) ClmmPositionWithdrawLiquidityUpdate {
    queryId: uint64
}

/// Owner locks the position's liquidity for good (there is no unlock); fees can
/// still be collected.
struct (0x850bdd2a) ClmmPositionLockLiquidity {
    queryId: uint64
}

/// Position left without bins returns its whole balance to the owner and is
/// destroyed.
struct (0x80381933) ClmmPositionClosed {
    queryId: uint64
}

fun onInternalMessage(in: InMessage) {
}
