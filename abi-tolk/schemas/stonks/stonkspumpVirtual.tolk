tolk 1.4

contract StonksPumpVirtualMinter {
    incomingMessages: StonksPumpVirtualMinterMessage
    storage: null
    forceAbiExport: [
        AskToPresaleSell
    ]
}

type StonksPumpVirtualMinterMessage =
    | PresaleSellNotificationForMinter
    | BuyFromPresale

/// Wallet owner request to sell jettons back into the virtual presale.
struct (0xcc1a97af) AskToPresaleSell {
    queryId: uint64
    jettonAmount: coins
    tonsRecipient: address?
    minTonOut: coins
}

/// Wallet-to-minter notification carrying burned presale-sale jettons.
struct (0xcc1a97b0) PresaleSellNotificationForMinter {
    queryId: uint64
    jettonAmount: coins
    sellInitiator: address
    tonsRecipient: address?
    minTonOut: coins
}

/// Minter request to buy jettons from the virtual presale.
struct (0xcc1a97ae) BuyFromPresale {
    queryId: uint64
    tonAmount: coins
    minJettonsOut: coins
}

/// External log emitted by the minter after a presale buy or sell.
struct (0xcc1a97b3) PresaleTradeEvent {
    queryId: uint64
    isBuy: bool
    amountIn: coins
    amountOut: coins
    recipient: address
    leftover: coins
}

/// Return shape of get_bonding_data().
struct BondingData {
    progressBps: int
    migrationTonTarget: int
    cumulativeTonReserve: int
    cumulativeJettonReserve: int
    tonsCollected: int
    maxBuyPercent: int
    tokensSold: int
    presaleOpen: int
}

fun onInternalMessage(in: InMessage) {
    val msg = lazy StonksPumpVirtualMinterMessage.fromSlice(in.body);

    match (msg) {
        BuyFromPresale => {
            createExternalLogMessage({
                dest: createAddressNone(),
                body: PresaleTradeEvent {
                    queryId: msg.queryId,
                    isBuy: true,
                    amountIn: msg.tonAmount,
                    amountOut: msg.minJettonsOut,
                    recipient: in.senderAddress,
                    leftover: 0,
                }
            }).send(SEND_MODE_REGULAR);
        }

        PresaleSellNotificationForMinter => {
            createExternalLogMessage({
                dest: createAddressNone(),
                body: PresaleTradeEvent {
                    queryId: msg.queryId,
                    isBuy: false,
                    amountIn: msg.jettonAmount,
                    amountOut: msg.minTonOut,
                    recipient: msg.tonsRecipient == null ? msg.sellInitiator : msg.tonsRecipient!,
                    leftover: 0,
                }
            }).send(SEND_MODE_REGULAR);
        }

        else => { }
    }
}

get fun get_bonding_data(): BondingData {
    return BondingData {
        progressBps: 0,
        migrationTonTarget: 0,
        cumulativeTonReserve: 0,
        cumulativeJettonReserve: 0,
        tonsCollected: 0,
        maxBuyPercent: 0,
        tokensSold: 0,
        presaleOpen: 0,
    };
}
