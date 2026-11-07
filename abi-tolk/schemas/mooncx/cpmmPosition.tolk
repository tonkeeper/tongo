tolk 1.4

// moon.cx CPMM LP position: ABI stub for tongo. There is one per (pool, owner);
// it holds the owner's LP and fee checkpoints and is not an NFT.
// Its messages to the pool use rich bounces; tongo already decodes those as BounceV2.
// getPosition and getLocked (camelCase, specific to moon.cx) are declared so the
// position can be detected: both must answer, and getPosition must parse as
// CpmmPositionStorage.

contract mooncx_cpmm_position {
    incomingMessages: MooncxCpmmPositionMessage
    storage: null
}

type MooncxCpmmPositionMessage =
    | CpmmPositionDeploy
    | CpmmPositionWithdrawLiquidity
    | CpmmPositionWithdrawFees
    | CpmmPositionWithdrawFeesUpdate
    | CpmmPositionLockLiquidity

struct CpmmPositionContent {
    lp: coins
    xFeesLp: uint256
    yFeesLp: uint256
}

/// Pool mints LP into the owner's position, deploying it with StateInit on the
/// first deposit; a later deposit is merged in. The position returns the GRAM
/// above 0.05 to the owner with CpmmExcess.
struct (0x7b5ceb3f) CpmmPositionDeploy {
    queryId: uint64
    content: Cell<CpmmPositionContent>
    locked: bool
}

/// Owner withdraws amount LP; the position sends CpmmPoolWithdrawLiquidity to the
/// pool. Refused while the position is locked or a fee withdrawal is pending.
struct (0x84548afa) CpmmPositionWithdrawLiquidity {
    queryId: uint64
    amount: coins
}

/// Owner collects the accrued fees; the position sends CpmmPoolWithdrawFees to the
/// pool and stays frozen until CpmmPositionWithdrawFeesUpdate arrives.
struct (0xf325e2f4) CpmmPositionWithdrawFees {
    queryId: uint64
}

/// Pool's answer to CpmmPoolWithdrawFees: the new fee checkpoints; unfreezes the position.
struct (0x90ad0f98) CpmmPositionWithdrawFeesUpdate {
    queryId: uint64
    xFeesLp: uint256
    yFeesLp: uint256
}

/// Owner locks the position's liquidity for good (there is no unlock); fees can
/// still be collected.
struct (0xb751f838) CpmmPositionLockLiquidity {
    queryId: uint64
}

/// Contract data returned by getPosition.
struct CpmmPositionStorage {
    pool: address
    user: address
    data: Cell<CpmmPositionContent>
    frozen: bool
    locked: bool
}

fun onInternalMessage(in: InMessage) {
}

get fun getPosition(): Cell<CpmmPositionStorage> {
    throw 0;
}

get fun getLocked(): bool {
    throw 0;
}
