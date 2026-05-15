contract RuntimeDefaultValues {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        Int8,
        Addr,
        DefaultTest
    )
}

type Int8 = int8
type Addr = address

struct DefaultTest {
    num1: uint32? = 4
    num2: Int8? = 5
    addr4: address? = null
    addr5: address? = address("0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe")
    tensor6: (int32, bool)? = (342, true)
    num7: int128? = null
    slice3: RemainingBitsAndRefs? = "616263".hexToSlice()
}

fun onInternalMessage() {}
