contract XtrMaster {
    incomingMessages: (PushXTR),
    storage: null,
    forceAbiExport: []
}

fun onInternalMessage(in: InMessage) {

}

struct (0x6f027868) PushXTR {
    seqno: uint64
    amount: coins
}

get fun get_user_latest_version(): uint32 {
    return 0;
}

get fun get_payment_latest_version(): uint32 {
    return 0;
}