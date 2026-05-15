contract RuntimeCustomPackUnpack {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        My,
        MyAlias
    )
}

struct My {
    amount: coins
}

type MyAlias = My

fun MyAlias.packToBuilder(self, mutate b: builder) {
    b.storeUint(0, 256);
    b.storeUint(0, 256);
    b.storeCoins(self.amount);
}

fun MyAlias.unpackFromSlice(mutate s: slice): MyAlias {
    s.skipBits(512);
    return { amount: s.loadCoins() };
}

fun onInternalMessage() {}
