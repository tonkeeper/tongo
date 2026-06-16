// Grambo Factory (Registry) Contract ABI schema

import "@stdlib/common"

struct (0x0d4295e0) GramboDeploy {
    query_id: uint64
    metadata: cell
    ton_target: coins
    prebuy_amount: coins
}

struct (0x20b031d9) GramboUpdateData {
    new_data: cell
}

struct (0x61d196c1) GramboUpdateCode {
    new_code: cell
}

type FactoryIncomingMessage =
    | GramboDeploy
    | GramboUpdateData
    | GramboUpdateCode

struct (0x372f1ae3) GramboLaunch {
    query_id: uint64
    buyer: address
    virtual_ton_reserve: coins
    ton_target: coins
    prebuy_amount: coins
}

contract grambo_factory {
    incomingMessages: FactoryIncomingMessage
    storage: null
    outgoingMessages: GramboLaunch
    forceAbiExport: [FactoryIncomingMessage]
}

fun onInternalMessage(in: InMessage) {
    throw 65299; // Invalid curve parameters
    throw 65347; // Unauthorized
    throw 65360; // Insufficient message value
}
