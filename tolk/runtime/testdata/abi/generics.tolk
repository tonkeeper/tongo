contract RuntimeGenerics {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        Wrapper<uint16>,
        EitherLeft<address>,
        EitherRight<coins>,
        Either<address, coins>,
        EitherAlias<bool, (bool, coins)>,
        EitherRef<uint32>,
        Doubler<(coins, address?)>,
        ManyRefsMsg<uint16>,
        GoodNamingForMsg<uint16>
    )
}

type Wrapper<W> = W

struct (0b0) EitherLeft<R> {
    value: R
}

struct (0b1) EitherRight<T> {
    value: T
}

type Either<L, R> = EitherLeft<L> | EitherRight<R>

struct EitherRef<I> {
    value: I | Cell<I>
}

type EitherAlias<R, L> = Either<R, L>

type Doubler<U> = (U, U)

struct ManyRefsMsg<T> {
    payload: EitherRef<uint32>
    either: Either<address, coins>
    anotherEither: EitherAlias<bool, (bool, coins)>
    doubler: Cell<Doubler<(coins, address?)>>
    myVal: Wrapper<T>
}

type GoodNamingForMsg<T> = ManyRefsMsg<T>

fun onInternalMessage() {}
