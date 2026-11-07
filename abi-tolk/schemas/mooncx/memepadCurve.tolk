tolk 1.4

// moon.cx memepad bonding curve (also the jetton master of its token):
// ABI stub for tongo.
// Sells are jetton burns (0x595f07bc) whose custom_payload is
// MemepadCurveSellPayload; the jetton wallet forwards it to the curve as the
// extra ref of burn_notification (0x7bdd97de). The payload is registered through
// JettonTransferNotification below; the TEP-74 burn and burn_notification
// themselves are not redeclared.
// No get methods on purpose: the curve exposes only generic jetton getters
// (get_jetton_data, get_wallet_address, get_next_admin_address), and declaring
// them would make every jetton master look like a moon.cx curve.
// Graduation messages to the moon.cx CPMM are in cpmmFactory.tolk and cpmmPool.tolk.

contract mooncx_memepad_curve {
    incomingMessages: MooncxMemepadCurveMessage
    storage: null
    forceAbiExport: [
        MemepadCurveSellPayout,
        MemepadCurveExcess,
        MemepadCurveCreatorFeesPayout,
        MemepadCurveProtocolPayout,
        MemepadCurveDeployError
    ]
}

// Registers the sell payload with tongo, whose jetton payload table also decodes
// the custom_payload of jetton_burn.
struct (0x7362d09c) JettonTransferNotification<T> {
    payload: T
}

type MooncxMemepadCurveMessage =
    | MemepadCurveInit
    | MemepadCurveBuy
    | MemepadCurveCreatorClaimFees
    | MemepadCurveUpgrade
    | JettonTransferNotification<MemepadCurveSellPayload>

struct EconState {
    tokensForSale: coins
    tokensForLp: coins
    virtualGram: coins
    virtualTokens: coins
    targetGram: coins
    tradeFeeBps: uint16
    creatorFeeBps: uint16
    poolFeeBps: uint8
}

/// Factory initializes a freshly deployed curve and runs the creator's initial buy, if any.
struct (0x61c9641a) MemepadCurveInit {
    queryId: uint64
    cpmmFactory: address
    protocolRecipient: address
    initialBuy: coins
    econ: Cell<EconState>
    jwCode: cell
    metadata: cell
}

/// Curve refunds the GRAM attached to MemepadCurveInit to the creator when the curve is already initialized.
struct (0x46f6a58e) MemepadCurveDeployError {
    queryId: uint64
}

/// Trader buys the curve's tokens for the GRAM attached to the message.
struct (0x29d1f1c1) MemepadCurveBuy {
    queryId: uint64
    minOut: coins
}

/// Custom payload of the jetton burn that sells tokens back to the curve (the
/// jetton wallet forwards it to the curve as the extra ref of burn_notification);
/// minOut is the minimum GRAM out.
struct (0xb97a94b7) MemepadCurveSellPayload {
    minOut: coins
}

/// Curve pays GRAM to the seller after a sell (jetton burn).
struct (0xdee7a887) MemepadCurveSellPayout {
    queryId: uint64
    gramOut: coins
}

/// Curve returns leftover GRAM: the change of a buy that overshoots the graduation
/// target, or the GRAM attached to a burn after the curve has graduated.
struct (0xca6c5778) MemepadCurveExcess {
    queryId: uint64
}

/// Creator collects accumulated creator fees.
struct (0x6ec846e7) MemepadCurveCreatorClaimFees {
    queryId: uint64
}

/// Curve pays accumulated creator fees (GRAM attached) to the creator.
struct (0xaa871c43) MemepadCurveCreatorFeesPayout {
    queryId: uint64
}

/// Curve sends accumulated platform fees (GRAM attached) to the factory,
/// on a creator fee claim and at graduation.
struct (0x0b8cca52) MemepadCurveProtocolPayout {
    queryId: uint64
}

/// Factory forwards a code/data upgrade to the curve.
struct (0x24293ccc) MemepadCurveUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

fun onInternalMessage(in: InMessage) {
}
