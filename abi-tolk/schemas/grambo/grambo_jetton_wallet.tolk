// Grambo Jetton Wallet Contract ABI schema

import "@stdlib/common"

struct GramboSellParams {
    min_ton_out: coins
    referrer: address
    custom_payload: cell?
}

struct (0x31fc1276) GramboSellNotification {
    query_id: uint64
    amount: coins
    from: address
    response_destination: address
    sell_params: Cell<GramboSellParams>?
}

struct (0x6b75d3e2) GramboActivateWallet {
    query_id: uint64
}

type WalletIncomingMessage =
    | GramboActivateWallet

struct GramboWalletStorage {
    activated: bool
    balance: coins
    owner: address
    master: address
}

contract grambo_jetton_wallet {
    incomingMessages: WalletIncomingMessage
    outgoingMessages: GramboSellNotification
    storage: GramboWalletStorage
    forceAbiExport: [WalletIncomingMessage]
}

fun onInternalMessage(in: InMessage) {
    throw 65397; // Wallet not activated
}

get fun get_wallet_status(): bool {
    return false;
}

struct GetWalletDataResult {
    balance: int
    owner_address: address
    jetton_master_address: address
    jetton_wallet_code: cell
}

get fun get_wallet_data(): GetWalletDataResult {
    return GetWalletDataResult {
        balance: 0,
        owner_address: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        jetton_master_address: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        jetton_wallet_code: beginCell().endCell(),
    };
}
