contract RuntimeHexUnion {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        UInt66WithPrefix,
        UInt33WithPrefix,
        UInt4WithPrefix,
        UInt66WithPrefix | UInt33WithPrefix | UInt4WithPrefix
    )
}

struct (0x12345678) UInt66WithPrefix {
    v: uint66
}

struct (0xdeadbeef) UInt33WithPrefix {
    v: uint33
}

struct (0x89abcdef) UInt4WithPrefix {
    v: uint4
}

fun onInternalMessage() {}
