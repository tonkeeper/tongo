contract RuntimeRandomFields {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        MyInt32,
        RandomFields
    )
}

type MyInt32 = int32

struct (0x794) RandomFields {
    dest_int: address
    amount: coins
    dest_ext: any_address
    intVector: (int32, coins?, uint64)
    needs_more: Cell<bool>
    some_payload: cell
    my_int: MyInt32
    some_union: int8 | int16 | int256
    default_1: int16 = 1
    default_2: int16? = 55
}

fun onInternalMessage() {}
