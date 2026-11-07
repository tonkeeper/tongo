tolk 1.4

// moon.cx CPMM pool (constant product): ABI stub for tongo.
// Pools are deployed by the factory EQCAYGE5-ogaKfvlFb-TO_GYKQsuiVOOLbKSDvAKJ9jkvP5W,
// e.g. EQDYaknWzA_oe_c2okEJmEf_7WWXZO6pKhoXuRoXFgPbD0HM.
// Jettons arrive as transfer_notification; the pool reads the forward payload
// only from a ref (an inline payload gets the jettons refunded). The payloads
// are registered through JettonTransferNotification below.
// Shared messages are declared only here and carry no contract in their names:
// the swap results (SwapFulfill, SwapReject and their jetton payloads), which
// moon.cx CLMM ranges send and accept with the same layout, and CpmmExcess,
// which the CPMM LP account and LP position send as well.
// No get methods: the deployed pool exposes none.

contract mooncx_cpmm_pool {
    incomingMessages: MooncxCpmmPoolMessage
    storage: null
    forceAbiExport: [
        CpmmPoolTakeReady,
        SwapReject,
        CpmmPoolProvideLiquidityRefund,
        CpmmPoolProvideLiquidityExcess,
        CpmmPoolWithdrawLiquidityPayout,
        CpmmPoolWithdrawFeesPayout,
        CpmmPoolCancelDepositRefund,
        CpmmPoolProtocolFeesPayout,
        CpmmExcess
    ]
}

// Registers the jetton forward payloads with tongo. SwapFulfillPayload and
// SwapRejectPayload travel in the jetton transfers the pool sends, and
// SwapFulfillPayload also in a transfer to the next-hop pool; the others travel
// in transfers to the pool.
struct (0x7362d09c) JettonTransferNotification<T> {
    payload: T
}

type MooncxCpmmPoolMessage =
    | CpmmPoolDeploy
    | CpmmPoolProvideReady
    | CpmmPoolSwap
    | SwapFulfill
    | CpmmPoolDepositLiquidity
    | CpmmPoolProvideLiquidity
    | CpmmPoolWithdrawLiquidity
    | CpmmPoolWithdrawFees
    | CpmmPoolCancelDeposit
    | CpmmPoolBoost
    | CpmmPoolCollectProtocolFees
    | CpmmPoolUpdateWallets
    | CpmmPoolUpgrade
    | JettonTransferNotification<CpmmPoolSwapPayload | SwapFulfillPayload | CpmmPoolDepositLiquidityPayload | CpmmPoolBoostPayload | SwapRejectPayload>

/// Factory deploys the pool (with StateInit). lpFee (the fee tier) and
/// protocolFee are charged on the swap input, in basis points; the codes are
/// those of the LP account and LP position contracts. The pool then asks the
/// jetton masters for its wallets.
struct (0x6ca4cb66) CpmmPoolDeploy {
    queryId: uint64
    lpFee: uint8
    protocolFee: uint8
    lpAccCode: cell
    lpPositionCode: cell
}

/// Asks the pool whether it is ready (its jetton wallets are known); attempt counts retries.
struct (0x44daf70f) CpmmPoolProvideReady {
    queryId: uint64
    attempt: uint8
}

/// Pool's answer to CpmmPoolProvideReady.
struct (0x67071f0d) CpmmPoolTakeReady {
    queryId: uint64
    attempt: uint8
    isReady: bool
}

/// GRAM swap parameters. excess is the response destination of the jetton
/// transfer that pays the result. nextFulfill / nextReject, when set, name the
/// next hop: a recipient address, optionally followed by a Maybe ^Cell payload
/// that is forwarded to it.
struct CpmmTonSwapParams {
    swapAmount: coins
    excess: address
    minOut: coins
    nextFulfill: cell?
    nextReject: cell?
}

/// Trader swaps the GRAM attached to the message (GRAM pools only); swapAmount 0
/// swaps all GRAM left after gas. The jettons bought arrive as a jetton transfer
/// with SwapFulfillPayload; a rejected swap returns the GRAM with
/// SwapReject.
struct (0xfa45056b) CpmmPoolSwap {
    queryId: uint64
    swapParams: CpmmTonSwapParams
}

/// Swap result paid in GRAM by a moon.cx CPMM pool (a moon.cx CLMM range sends
/// the same message) to the trader or to the nextFulfill recipient: amount is the
/// GRAM bought, payload is the next hop's payload. A pool or range that receives
/// it as the next hop of a route swaps the attached GRAM; payload then holds its
/// GRAM swap parameters (CpmmTonSwapParams, ClmmTonSwapParams for a range).
struct (0x9d142643) SwapFulfill {
    queryId: uint64
    amount: coins
    payload: cell?
}

/// Rejected swap: a moon.cx CPMM pool (a moon.cx CLMM range sends the same
/// message) returns the GRAM to the trader or to the nextReject recipient.
/// amount is the rejected swap amount (a CLMM range sends 0 when it rejects
/// before swapping); payload is the next hop's payload.
struct (0xff32da99) SwapReject {
    queryId: uint64
    amount: coins
    payload: cell?
}

struct LpSplit {
    aRecipient: address
    aBps: uint16
    bRecipient: address
}

/// Deposits the GRAM attached to the message as liquidity into a GRAM pool;
/// split optionally shares the LP between two recipients.
struct (0x90af0b63) CpmmPoolDepositLiquidity {
    queryId: uint64
    minLpOut: coins
    locked: bool
    split: Cell<LpSplit>?
}

/// LP account of provider sends the amounts of both legs of a deposit, which the
/// pool already holds, once the second one is recorded. The pool mints LP into
/// one LP position (two with split, allowed only on the first deposit) and
/// returns the rest with CpmmPoolProvideLiquidityExcess, or refunds both legs.
struct (0x48a936dd) CpmmPoolProvideLiquidity {
    queryId: uint64
    provider: address
    xAmount: coins
    yAmount: coins
    minLpOut: coins
    locked: bool
    split: Cell<LpSplit>?
}

/// LP position of provider asks the pool to burn lpPos LP and pay out both assets
/// plus the fees accrued since its fee checkpoints xFeesPos / yFeesPos.
struct (0xcbd7f377) CpmmPoolWithdrawLiquidity {
    queryId: uint64
    provider: address
    lpPos: coins
    xFeesPos: uint256
    yFeesPos: uint256
}

/// LP position of provider asks the pool to pay the fees accrued on lpPos LP since
/// its fee checkpoints; the pool answers with CpmmPositionWithdrawFeesUpdate.
struct (0xd5c9f8c3) CpmmPoolWithdrawFees {
    queryId: uint64
    provider: address
    lpPos: coins
    xFeesPos: uint256
    yFeesPos: uint256
}

/// LP account of provider asks the pool to refund a pending deposit leg that a
/// new leg of the same side replaced.
struct (0x1f911dd9) CpmmPoolCancelDeposit {
    queryId: uint64
    provider: address
    xAmount: coins
    yAmount: coins
}

/// Anyone donates the attached GRAM (minus 0.05) to the GRAM reserve of a GRAM
/// pool that has liquidity.
struct (0x50bbec04) CpmmPoolBoost {
    queryId: uint64
}

/// Factory forwards CpmmFactoryCollectProtocolFees: the pool pays its protocol
/// fees to recipient.
struct (0xa7335450) CpmmPoolCollectProtocolFees {
    queryId: uint64
    recipient: address
}

/// Factory forwards CpmmFactoryUpdatePoolWallets: the pool replaces its jetton wallets.
struct (0xc29e72ba) CpmmPoolUpdateWallets {
    queryId: uint64
    newXWallet: address?
    newYWallet: address
}

/// Factory forwards CpmmFactoryUpgradePool: the pool replaces its code and/or data.
struct (0x38599428) CpmmPoolUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

/// Jetton swap parameters: as CpmmTonSwapParams, without swapAmount (the whole
/// transferred amount is swapped).
struct CpmmJettonSwapParams {
    excess: address
    minOut: coins
    nextFulfill: cell?
    nextReject: cell?
}

/// Forward payload of a jetton transfer to the pool: swap the jettons for the
/// other asset. The query id is the transfer's.
struct (0xfa45056b) CpmmPoolSwapPayload {
    swapParams: CpmmJettonSwapParams
}

/// Forward payload of a swap result paid in jettons by a moon.cx CPMM pool (a
/// moon.cx CLMM range sends the same payload) to the trader or to the nextFulfill
/// recipient; payload is the next hop's payload. A pool or range that receives
/// it as the next hop of a route swaps the jettons; payload then holds its jetton
/// swap parameters (CpmmJettonSwapParams, ClmmJettonSwapParams for a range).
struct (0x9d142643) SwapFulfillPayload {
    payload: cell?
}

/// Forward payload of the jetton refund of a rejected jetton swap, sent by a
/// moon.cx CPMM pool (a moon.cx CLMM range sends the same payload) to the trader
/// or to the nextReject recipient; payload is the next hop's payload.
struct (0xff32da99) SwapRejectPayload {
    payload: cell?
}

/// Forward payload of a jetton transfer to the pool: the jetton leg of a liquidity
/// deposit (CpmmPoolDepositLiquidity is the GRAM leg).
struct (0x90af0b63) CpmmPoolDepositLiquidityPayload {
    minLpOut: coins
    locked: bool
    split: Cell<LpSplit>?
}

/// Forward payload of a jetton transfer to the pool: donate the jettons to the
/// reserves of a pool that has liquidity.
struct (0x50bbec04) CpmmPoolBoostPayload {
}

/// Pool refunds the GRAM leg of a CpmmPoolProvideLiquidity that failed (minimum
/// LP not met, zero LP out, an amount above 2^120 - 1, first deposit not made by
/// the pool initiator, or a split on a pool that already has liquidity); a jetton
/// leg comes back as a plain jetton transfer.
struct (0x67568a3a) CpmmPoolProvideLiquidityRefund {
    queryId: uint64
}

/// Pool returns the GRAM left after a successful CpmmPoolProvideLiquidity.
struct (0x9d9f7bad) CpmmPoolProvideLiquidityExcess {
    queryId: uint64
}

/// Pool pays the GRAM leg of a CpmmPoolWithdrawLiquidity to the provider; the
/// jetton leg is a plain jetton transfer.
struct (0x15b7ef52) CpmmPoolWithdrawLiquidityPayout {
    queryId: uint64
}

/// Pool pays the GRAM leg of a CpmmPoolWithdrawFees to the provider; the jetton
/// leg is a plain jetton transfer.
struct (0x33e81d59) CpmmPoolWithdrawFeesPayout {
    queryId: uint64
}

/// Pool refunds the GRAM leg of a deposit cancelled by CpmmPoolCancelDeposit.
struct (0xd212503b) CpmmPoolCancelDepositRefund {
    queryId: uint64
}

/// Pool pays the GRAM part of its protocol fees to the recipient of
/// CpmmPoolCollectProtocolFees; the jetton part is a plain jetton transfer.
struct (0x15e67eab) CpmmPoolProtocolFeesPayout {
    queryId: uint64
}

/// Rest of the GRAM sent back by a moon.cx CPMM pool (after a payout or refund),
/// LP account (after recording a deposit leg) or LP position (after a deploy).
struct (0x13925ff2) CpmmExcess {
    queryId: uint64
}

fun onInternalMessage(in: InMessage) {
}
