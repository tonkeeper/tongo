import "payloads"
import "abi/generics"

contract Jetton {
    incomingMessages: ()
    storage: null
    forceAbiExport: Transfer
}

struct (0x0f8a7ea5) Transfer {
    queryId: uint64
    amount: coins
    destination: address
    responseDestination: address
    customPayload: cell?
    forwardAmount: coins
    forwardPayload: EitherRef<WrapJettonPayload>
}


// WrapJettonPayload avoid union expansion:
//     EitherRef<I> is I | Cell<I>
//     so EitherRef<JettonPayload> is Payload1 | Payload2 | PayloadN | Cell<...> which produces incorrect encoding
//        EitherRef<WrapJettonPayload> ensures of exactly two union variants: WrapJettonPayload | Cell<WrapJettonPayload>
//
//   todo maybe decode forward payload using type EitherRef<RemainingBitsAndRefs>
//      This decoding might be up to UI:
//        consider incoming message:
//          Transfer {
//              forwardPayload: EitherRef<RemainingBitsAndRefs> = Right(msg_cell_A_OUT)
//          }
//        in one of the next transactions there will be:
//          msg_cell_B_IN => BidaskDammProvideBothJettonPayload {
//             ...
//          }
//        if msg_cell_B_IN == msg_cell_A_OUT, then we can show that indeed payload was BidaskDammProvideBothJettonPayload
//
struct WrapJettonPayload {
    value: JettonPayload
}

fun onInternalMessage() {}
