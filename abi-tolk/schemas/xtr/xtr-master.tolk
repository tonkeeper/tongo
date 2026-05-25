contract XtrMaster {
    incomingMessages: (UpdateUser, UpdatePayment, PushXTR),
    storage: null,
    forceAbiExport: []
}

fun onInternalMessage(in: InMessage) {

}

struct (0x6f027868) PushXTR {
    seqno: uint64
    amount: coins
}

struct (0x59ff81d6) UpdateUser {
    destAddress: address
    payload: cell
}

struct (0x2a999df4) UpdatePayment {
    destAddress: address
    payload: cell
}

get fun get_user_latest_version(): uint32 {
    return 0;
}

get fun get_payment_latest_version(): uint32 {
    return 0;
}