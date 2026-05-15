contract RuntimeNumbers {
    incomingMessages: ()
    storage: null
    forceAbiExport: Numbers
}

struct Numbers {
    num1: uint64
    num3: uint256
    num4: varuint32
    num5: bool
    num7: bits16
    num8: int14
    num9: varint16
}

fun onInternalMessage() {}
