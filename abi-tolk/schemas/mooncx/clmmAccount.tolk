tolk 1.4

// moon.cx CLMM LP account: ABI stub for tongo. There is one per (range, user).
// It only records the amount of the first leg of a two-sided liquidity deposit;
// the GRAM or jettons of that leg stay on the range. When the other leg is
// recorded, it sends both amounts to the range (ClmmRangeProvideLiquidity),
// which adds the liquidity or refunds both legs. A new leg of the same side
// makes the range refund the pending one (ClmmRangeCancelDeposit). After
// recording a leg it returns the gas with ClmmExcess.
// No get methods: the deployed LP account exposes none.

contract mooncx_clmm_account {
    incomingMessages: ClmmAccountDepositLiquidityRecord
    storage: null
}

struct ClmmDepositBin {
    x: coins
    y: coins
}

/// Range records one leg of a two-sided deposit in the provider's LP account
/// (deploying it with StateInit on first use); xIn tells which side amount is.
struct (0xed8efa06) ClmmAccountDepositLiquidityRecord {
    queryId: uint64
    amount: coins
    xIn: bool
    depositBins: map<int64, ClmmDepositBin>
}

fun onInternalMessage(in: InMessage) {
}
