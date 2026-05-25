contract XtrMaster {
    incomingMessages: (UpdateUser, UpdatePayment, UpdateContractAndProcessMessage, PushXTR),
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

struct (0xe2628539) UpdateContractAndProcessMessage {
    updateData: Cell<UpdateData>
    fromAddress: address
    fromAmount: coins
    payload: cell
}

struct UpdateData {
    code: cell
    data: cell
    version: uint32
}

get fun get_user_latest_version(): uint32 {
    return 0;
}

get fun get_payment_latest_version(): uint32 {
    return 0;
}