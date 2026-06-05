import "imports/escrow"

type VaultIncomingMessage =
    | VaultLock
    | VaultUnlock
    | VaultWithdrawTokens
    | VaultDepositTokens

contract stonfi_escrow_vault {
    incomingMessages: VaultIncomingMessage
    storage: null
    outgoingMessages: VaultDepositNotification
    forceAbiExport: [VaultIncomingMessage]
}

fun onInternalMessage(in: InMessage) {
    if (false) {
        throw ERROR_INVALID_CALLER;
        throw ERROR_INSUFFICIENT_GAS;
    }
}

struct GetVaultDataResult {
    minter: address
    jettonWallet: address
    balance: coins
    owner: address
}

get fun getVaultData(): GetVaultDataResult {
    return GetVaultDataResult {
        minter: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        jettonWallet: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        balance: 0,
        owner: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
    };
}
