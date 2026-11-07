tolk 1.4

// moon.cx CPMM LP account: ABI stub for tongo. There is one per (pool, user).
// It only records the amount of the first leg of a two-sided liquidity deposit;
// the GRAM or jettons of that leg stay on the pool. When the other leg is
// recorded, it sends both amounts to the pool (CpmmPoolProvideLiquidity). A new
// leg of the same side makes the pool refund the pending one
// (CpmmPoolCancelDeposit). After recording a leg it returns the rest of the GRAM
// with CpmmExcess.
// No get methods: the deployed LP account exposes none.

contract mooncx_cpmm_account {
    incomingMessages: CpmmAccountDepositLiquidityRecord
    storage: null
}

struct LpSplit {
    aRecipient: address
    aBps: uint16
    bRecipient: address
}

/// Pool records one leg of a liquidity deposit in the provider's LP account
/// (deploying it with StateInit on first use); xIn tells which side amount is.
struct (0xfda9374a) CpmmAccountDepositLiquidityRecord {
    queryId: uint64
    xIn: bool
    amount: coins
    minLpOut: coins
    locked: bool
    split: Cell<LpSplit>?
}

fun onInternalMessage(in: InMessage) {
}
