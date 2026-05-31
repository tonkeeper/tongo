tolk 1.4

contract StonksPumpVirtualFactory {
    incomingMessages: StonksPumpVirtualFactoryMessage
    storage: null
}

type StonksPumpVirtualFactoryMessage =
    | CreateVirtualLiquidityJetton

/// Factory request that deploys and initializes a virtual-liquidity jetton.
struct (0xcc1a9801) CreateVirtualLiquidityJetton {
    queryId: uint64
    content: cell
    devBuyTonAmount: coins
    maxBuyPercent: uint16
    devWallet: address
}

fun onInternalMessage(in: InMessage) {
    val msg = lazy StonksPumpVirtualFactoryMessage.fromSlice(in.body);

    match (msg) {
        CreateVirtualLiquidityJetton => { }
        else => { }
    }
}
