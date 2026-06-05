import "imports/escrow"

type FactoryIncomingMessage =
    | MinterInitTransfer
    | MinterRefundRequest
    | MinterInternalLock
    | MinterInternalUnlock
    | MinterInternalWithdrawTokens
    | MinterGiveProtocolOwnership
    | MinterUpdateProtocolTier
    | MinterTakeProtocolOwnership
    | MinterCancelNextProtocolOwner
    | MinterResetGas
    | MinterLockPayload
    | MinterUnlockPayload
    | MinterDepositVault

contract stonfi_escrow_factory {
    incomingMessages: FactoryIncomingMessage
    storage: null
    outgoingMessages: LockRejectNotification
    forceAbiExport: [FactoryIncomingMessage]
}

fun onInternalMessage(in: InMessage) {
    throw ERROR_INVALID_CALL;
    throw ERROR_INVALID_CALLER;
    throw ERROR_WRONG_WORKCHAIN;
    throw ERROR_INSUFFICIENT_GAS;
    throw ERROR_INVALID_AMOUNT;
    throw ERROR_INVALID_TRANSFER_NOTIFICATION_PAYLOAD;
}

struct EscrowData {
    globalId: uint64
    protocol: address
    protocolTier: uint16
    itemCode: cell
    vaultCode: cell
}

get fun getEscrowData(): EscrowData {
    return EscrowData {
        globalId: 0,
        protocol: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        protocolTier: 0,
        itemCode: beginCell().endCell(),
        vaultCode: beginCell().endCell(),
    };
}

get fun getItemAddress(quoteId: int): address {
    return address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs");
}

get fun getVaultAddress(jettonWallet: address, owner: address): address {
    return address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs");
}

struct GetVersionResult {
    major: uint8
    minor: uint8
    development: string
}

get fun getVersion(): GetVersionResult {
    return GetVersionResult {
        major: 0,
        minor: 0,
        development: "",
    };
}
