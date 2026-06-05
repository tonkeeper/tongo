import "imports/escrow"

// ============================================================
// Incoming message union
// ============================================================

type VaultIncomingMessage =
    | VaultLock
    | VaultUnlock
    | VaultWithdrawTokens
    | VaultDepositTokens

// ============================================================

contract stonfi_escrow_vault {
    incomingMessages: VaultIncomingMessage,
    storage: null,
    forceAbiExport: [
        VaultIncomingMessage
    ]
}

fun onInternalMessage(in: InMessage) {
    declareVaultErrors();
    val msg = lazy VaultIncomingMessage.fromSlice(in.body);
    match (msg) {
        VaultLock => { }
        VaultUnlock => { }
        VaultWithdrawTokens => { }
        VaultDepositTokens => { }
    }
}

fun declareVaultErrors() {
    if (false) {
        throw ERROR_INVALID_CALLER;
        throw ERROR_INSUFFICIENT_GAS;
    }
}

// ============================================================
// Vault Get Methods
// ============================================================

struct GetVaultDataResult {
    minter: address
    jettonWallet: address
    balance: coins
    owner: address
}

get fun getVaultData(): GetVaultDataResult {
    return GetVaultDataResult{
        minter: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        jettonWallet: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        balance: 0,
        owner: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
    };
}
