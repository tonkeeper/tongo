tolk 1.4

// moon.cx memepad factory: ABI stub for tongo (Tonviewer / TonAPI decoding).
// Messages and opcodes match the deployed factory
// EQDOEt2Hvs6nOFA22KHUEa3BEbS7Hy2ed-9yllPy5nk8Hjhs.
// No get methods: the deployed factory exposes none.

contract mooncx_memepad_factory {
    incomingMessages: MooncxMemepadFactoryMessage
    storage: null
    forceAbiExport: [
        MemepadFactoryWithdrawPayout
    ]
}

type MooncxMemepadFactoryMessage =
    | MemepadFactoryCreateCurve
    | MemepadFactoryWithdraw
    | MemepadFactoryUpgrade
    | MemepadFactoryUpgradeCurve

/// Creator asks the factory to deploy a new curve (bonding curve + jetton).
/// content holds name, symbol and image as text refs plus an optional ref with extra fields.
struct (0x7e838abb) MemepadFactoryCreateCurve {
    queryId: uint64
    salt: uint64
    feeBps: uint8
    initialBuy: coins
    content: cell
}

/// Admin withdraws GRAM from the factory.
struct (0x329068a6) MemepadFactoryWithdraw {
    queryId: uint64
    amount: coins
    recipient: address
}

/// Factory pays the withdrawn GRAM to the recipient of MemepadFactoryWithdraw.
struct (0xfce532d8) MemepadFactoryWithdrawPayout {
    queryId: uint64
}

/// Admin replaces factory code and/or data.
struct (0x0eb678fe) MemepadFactoryUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

/// Admin asks the factory to upgrade one curve.
struct (0x9057f56b) MemepadFactoryUpgradeCurve {
    queryId: uint64
    curve: address
    newCode: cell?
    newData: cell?
}

fun onInternalMessage(in: InMessage) {
}
