import "payloads"
import "abi/generics"

contract Jetton {
    incomingMessages: ()
    storage: null
    forceAbiExport: [
        Transfer,
        JettonPayload,
    ]
}

struct (0x0f8a7ea5) Transfer {
    queryId: uint64
    amount: coins
    destination: address
    responseDestination: address
    customPayload: cell?
    forwardAmount: coins
    forwardPayload: EitherRef<RemainingBitsAndRefs>
}

fun onInternalMessage() {}
