import "imports/escrow"

// ============================================================
// Incoming message unions
// ============================================================

type PositionIncomingMessage =
    | ItemInternalLock
    | ItemInternalUnlock
    | ItemWithdraw

type PositionExternalMessage =
    | ExternalItemWithdraw
    | ExternalCronTrigger

// ============================================================

contract stonfi_escrow_position {
    incomingMessages: PositionIncomingMessage,
    incomingExternal: PositionExternalMessage,
    storage: null,
    forceAbiExport: [
        PositionIncomingMessage,
        PositionExternalMessage
    ]
}

fun onInternalMessage(in: InMessage) {
    declarePositionErrors();
    val msg = lazy PositionIncomingMessage.fromSlice(in.body);
    match (msg) {
        ItemInternalLock => { }
        ItemInternalUnlock => { }
        ItemWithdraw => { }
    }
}

fun declarePositionErrors() {
    if (false) {
        throw ERROR_WRONG_STATUS;
        throw ERROR_INSUFFICIENT_GAS;
        throw ERROR_INVALID_AMOUNT;
    }
}

fun onCreateMessageStubs() {
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        value: 0,
        body: ItemLockSuccess { queryId: 0, forwardPayload: null }
    });
}

// ============================================================
// Position Get Methods
// ============================================================

struct GetOrderDataResult {
    status: uint8
    quoteId: int
    minter: address
    bidJettonWallet: address
    bidJettonAmount: coins
    bidFromUser: address
    refFee: address
    refFeeTier: uint16
    pendingUserReceiveAmount: coins
    askJettonWallet: address
    askJettonMinter: address
    conditionState: cell
    conditionCode: cell
    userUnlockForwardParams: coins?
    userUnlockForwardSuccessPayload: cell?
    userFillToVault: bool
    safeDepositAndForwardValue: uint64
    excessesReceiver: address
}

get fun getOrderData(): GetOrderDataResult {
    return GetOrderDataResult{
        status: 0,
        quoteId: 0,
        minter: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        bidJettonWallet: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        bidJettonAmount: 0,
        bidFromUser: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        refFee: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        refFeeTier: 0,
        pendingUserReceiveAmount: 0,
        askJettonWallet: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        askJettonMinter: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        conditionState: beginCell().endCell(),
        conditionCode: beginCell().endCell(),
        userUnlockForwardParams: 0,
        userUnlockForwardSuccessPayload: null,
        userFillToVault: false,
        safeDepositAndForwardValue: 0,
        excessesReceiver: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
    };
}

struct GetCronInfoResult {
    nextCallTime: uint32
    reward: coins
    balanceMinusAmounts: int64
    repeatEvery: uint32
}

get fun get_cron_info(): GetCronInfoResult {
    return GetCronInfoResult{
        nextCallTime: 0,
        reward: 0,
        balanceMinusAmounts: 0,
        repeatEvery: 0,
    };
}
