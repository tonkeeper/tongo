tolk 1.4

// moon.cx CLMM range: ABI stub for tongo. Each pool has three ranges (rangeNum
// -1, 0 and 1); range r holds the bins 2048*r - 1024 .. 2048*r + 1023 with all
// their liquidity and protocol fees. A swap runs inside one range, so clients
// send it to the range that holds the current bin.
// Example: EQCEDqD6beyufWkmtCEnkW490ZcOZVrpkqRF7fkqn9OZyLOs (GRAM/USDT, range 0).
// Jettons arrive as transfer_notification; the range reads the forward payload
// only from a ref (an inline payload gets the jettons refunded).
// Swap results have the same layout as in the moon.cx CPMM pool and are declared
// in cpmmPool.tolk: SwapFulfill / SwapReject when paid in GRAM,
// SwapFulfillPayload / SwapRejectPayload inside the jetton
// transfer otherwise. The range also accepts SwapFulfill and
// SwapFulfillPayload as the next hop of a route, with ClmmTonSwapParams /
// ClmmJettonSwapParams in their payload.
// getRange is declared so the range can be detected; getUserLpAccount takes an
// argument and is left out.

contract mooncx_clmm_range {
    incomingMessages: MooncxClmmRangeMessage
    storage: null
    forceAbiExport: [
        ClmmRangeDepositLiquidityExcess,
        ClmmRangeDepositLiquidityRefund,
        ClmmRangeWithdrawLiquidityPayout,
        ClmmRangeWithdrawFeesPayout,
        ClmmRangeCancelDepositRefund,
        ClmmRangeProtocolFeesPayout,
        ClmmExcess
    ]
}

// Registers the jetton forward payloads of transfers to the range with tongo.
struct (0x7362d09c) JettonTransferNotification<T> {
    payload: T
}

type MooncxClmmRangeMessage =
    | ClmmRangeDeploy
    | ClmmRangeSwap
    | ClmmRangeDepositLiquidity
    | ClmmRangeProvideLiquidity
    | ClmmRangeWithdrawFees
    | ClmmRangeWithdrawLiquidity
    | ClmmRangeCancelDeposit
    | ClmmRangeCollectProtocolFees
    | ClmmRangeUpdateFees
    | ClmmRangeUpdateWallets
    | ClmmRangeUpgrade
    | JettonTransferNotification<ClmmRangeSwapPayload | ClmmRangeDepositLiquidityPayload>

struct ClmmRangeInitPrice {
    initBin: int64
    initSqrtPrice: uint256
}

/// Pool deploys the range (with StateInit and 1 GRAM) with the pair, bin step,
/// LP and protocol fees (basis points), start price and the LP account and LP
/// position codes; the range then asks the jetton masters for its wallets.
struct (0x724b6005) ClmmRangeDeploy {
    queryId: uint64
    xMaster: address?
    yMaster: address
    binStep: uint16
    lpFee: uint8
    protocolFee: uint8
    initPrice: Cell<ClmmRangeInitPrice>
    lpAccCode: cell
    lpPositionCode: cell
}

/// GRAM swap parameters; sqrtPriceLimit 0 means no price limit. excess is the
/// response destination of the jetton transfer that pays the result.
/// nextFulfill / nextReject, when set, name the next hop: a recipient address,
/// optionally followed by a Maybe ^Cell payload that is forwarded to it.
struct ClmmTonSwapParams {
    swapAmount: coins
    excess: address
    sqrtPriceLimit: uint256
    nextFulfill: cell?
    nextReject: cell?
}

/// Trader swaps the GRAM attached to the message (ranges whose x side is GRAM);
/// swapAmount 0 swaps all GRAM left after gas. The jettons bought arrive as a
/// jetton transfer with SwapFulfillPayload; a rejected swap returns the
/// GRAM with SwapReject.
struct (0xd151210f) ClmmRangeSwap {
    queryId: uint64
    swapParams: ClmmTonSwapParams
}

struct ClmmDepositBin {
    x: coins
    y: coins
}

/// Provider adds GRAM liquidity (ranges whose x side is GRAM): depositAmount is
/// spread over depositBins (bin index -> amounts). With isSingle the range adds
/// it at once and deploys a new LP position; otherwise the range keeps the GRAM
/// and records the amount in the provider's LP account, which pairs it with the
/// jetton leg (ClmmRangeProvideLiquidity).
struct (0xec401442) ClmmRangeDepositLiquidity {
    queryId: uint64
    depositAmount: coins
    isSingle: bool
    depositBins: map<int64, ClmmDepositBin>
}

/// LP account of provider sends the amounts of both legs of a two-sided deposit,
/// which the range already holds. The range adds the liquidity and deploys a new
/// LP position, or refunds both legs (the GRAM leg with
/// ClmmRangeDepositLiquidityRefund).
struct (0x344629fc) ClmmRangeProvideLiquidity {
    queryId: uint64
    provider: address
    xAmount: coins
    yAmount: coins
    depositBins: map<int64, ClmmDepositBin>
}

struct ClmmBin {
    l: uint256
    xFee: uint256
    yFee: uint256
}

/// LP position nftIndex of provider asks the range to pay the fees accrued on its
/// bins (bin index -> liquidity and fee checkpoints); the range pays them and
/// answers with ClmmPositionWithdrawFeesUpdate.
struct (0x002618d9) ClmmRangeWithdrawFees {
    queryId: uint64
    provider: address
    nftIndex: uint64
    withdrawBins: map<int64, ClmmBin>
}

/// LP position nftIndex of provider withdraws the liquidity of the given bins;
/// the range pays out both assets plus fees and answers with
/// ClmmPositionWithdrawLiquidityUpdate.
struct (0x69099bab) ClmmRangeWithdrawLiquidity {
    queryId: uint64
    provider: address
    nftIndex: uint64
    withdrawBins: map<int64, ClmmBin>
}

/// LP account of provider asks the range to refund a pending deposit leg that a
/// new leg of the same side replaced.
struct (0x2dcb38cb) ClmmRangeCancelDeposit {
    queryId: uint64
    xAmount: coins
    yAmount: coins
    provider: address
}

/// Pool forwards ClmmPoolCollectProtocolFees: the range pays its protocol fees to
/// recipient.
struct (0xdbc809b5) ClmmRangeCollectProtocolFees {
    queryId: uint64
    recipient: address
}

/// Pool forwards ClmmPoolUpdateRangeFees: new LP and protocol fees (basis points).
struct (0x90be29eb) ClmmRangeUpdateFees {
    queryId: uint64
    lpFee: uint8
    protocolFee: uint8
}

/// Pool forwards ClmmPoolUpdateRangeWallets: the range replaces its jetton wallets.
struct (0xa17882c8) ClmmRangeUpdateWallets {
    queryId: uint64
    xWallet: address?
    yWallet: address
}

/// Pool forwards ClmmPoolUpgradeRange: the range replaces its code and/or data.
struct (0xf0262313) ClmmRangeUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

/// Jetton swap parameters: as ClmmTonSwapParams, without swapAmount (the whole
/// transferred amount is swapped).
struct ClmmJettonSwapParams {
    excess: address
    sqrtPriceLimit: uint256
    nextFulfill: cell?
    nextReject: cell?
}

/// Forward payload of a jetton transfer to the range: swap the jettons for the
/// other asset. The query id is the transfer's.
struct (0xd151210f) ClmmRangeSwapPayload {
    swapParams: ClmmJettonSwapParams
}

/// Forward payload of a jetton transfer to the range: a jetton liquidity deposit,
/// as ClmmRangeDepositLiquidity. With isSingle the range adds it at once and
/// deploys a new LP position; otherwise it is one leg of a two-sided deposit.
struct (0xec401442) ClmmRangeDepositLiquidityPayload {
    isSingle: bool
    depositBins: map<int64, ClmmDepositBin>
}

/// Range returns the GRAM of a deposit that did not go into the bins: the amount
/// above the bin amounts and the part the bins could not take.
struct (0x89983597) ClmmRangeDepositLiquidityExcess {
    queryId: uint64
}

/// Range refunds the GRAM of a failed deposit (a single GRAM deposit or the GRAM
/// leg of ClmmRangeProvideLiquidity); a jetton leg comes back as a plain jetton
/// transfer.
struct (0x77714776) ClmmRangeDepositLiquidityRefund {
    queryId: uint64
}

/// Range pays the GRAM part of a ClmmRangeWithdrawLiquidity (liquidity plus fees)
/// to the provider; the jetton part is a plain jetton transfer.
struct (0xacae4171) ClmmRangeWithdrawLiquidityPayout {
    queryId: uint64
}

/// Range pays the GRAM part of the fees of a ClmmRangeWithdrawFees to the
/// provider; the jetton part is a plain jetton transfer.
struct (0xbc49e37b) ClmmRangeWithdrawFeesPayout {
    queryId: uint64
}

/// Range refunds the GRAM leg of a deposit cancelled by ClmmRangeCancelDeposit.
struct (0x5db3ae19) ClmmRangeCancelDepositRefund {
    queryId: uint64
}

/// Range pays the GRAM part of its protocol fees to the recipient of
/// ClmmRangeCollectProtocolFees; the jetton part is a plain jetton transfer.
struct (0x9d6c9b24) ClmmRangeProtocolFeesPayout {
    queryId: uint64
}

/// Rest of the GRAM sent back by a moon.cx CLMM range (after a payout or refund)
/// or LP account (after recording a deposit leg).
struct (0x5fd5f2b3) ClmmExcess {
    queryId: uint64
}

struct ClmmRangeState {
    currentBin: int64
    currentSqrtP: uint256
    xProtocolFees: coins
    yProtocolFees: coins
}

/// Contract data returned by getRange once the range knows its jetton wallets.
struct ClmmRangeStorage {
    pool: address
    rangeNum: int16
    xWallet: address?
    yWallet: address
    binStep: uint16
    lpFee: uint8
    protocolFee: uint8
    positionIndex: uint64
    rangeState: Cell<ClmmRangeState>
    bins: map<int64, ClmmBin>
    lpAccountCode: cell
    lpPositionCode: cell
}

fun onInternalMessage(in: InMessage) {
}

get fun getRange(): Cell<ClmmRangeStorage> {
    throw 0;
}
