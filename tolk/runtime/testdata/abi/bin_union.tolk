contract RuntimeBinUnion {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        AddressWithPrefix,
        MapWithPrefix,
        CellWithPrefix,
        AddressWithPrefix | MapWithPrefix | CellWithPrefix
    )
}

struct (0b001) AddressWithPrefix {
    v: address
}

struct (0b011) MapWithPrefix {
    v: map<address, coins>
}

struct (0b111) CellWithPrefix {
    v: cell
}

fun onInternalMessage() {}
