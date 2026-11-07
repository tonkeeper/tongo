tolk 1.4

// moon.cx CPMM factory: ABI stub for tongo. Users create pools here (a memepad
// curve does the same at graduation); the admin messages except
// CpmmFactoryUpgrade are relayed to a pool.
// Deployed factory: EQCAYGE5-ogaKfvlFb-TO_GYKQsuiVOOLbKSDvAKJ9jkvP5W
// (the older factory EQCTBiBa_B4UknFdYaeCh_zMHLjStrdzIZ357cygQFmP0xhv uses the
// same messages).
// No get methods: the deployed factory exposes none.

contract mooncx_cpmm_factory {
    incomingMessages: MooncxCpmmFactoryMessage
    storage: null
    forceAbiExport: [
        CpmmFactoryTakePoolAddress
    ]
}

type MooncxCpmmFactoryMessage =
    | CpmmFactoryCreatePool
    | CpmmFactoryProvidePoolAddress
    | CpmmFactoryCollectProtocolFees
    | CpmmFactoryUpdatePoolWallets
    | CpmmFactoryUpgradePool
    | CpmmFactoryUpgrade

/// Creates a pool for a pair; a null master stands for GRAM.
struct (0xec725c4f) CpmmFactoryCreatePool {
    queryId: uint64
    xMaster: address?
    yMaster: address?
    feeBps: uint8
}

/// Asks the factory for the address of the pool of a pair created by initiator;
/// the answer goes back to the sender.
struct (0x1d42af37) CpmmFactoryProvidePoolAddress {
    queryId: uint64
    xMaster: address?
    yMaster: address
    feeBps: uint8
    initiator: address
}

/// Factory's answer to CpmmFactoryProvidePoolAddress.
struct (0x81a58422) CpmmFactoryTakePoolAddress {
    queryId: uint64
    pool: address
    xMaster: address?
    yMaster: address
    feeBps: uint8
}

/// Admin collects the protocol fees of a pool; the factory forwards it to the
/// pool as CpmmPoolCollectProtocolFees.
struct (0xc492cce9) CpmmFactoryCollectProtocolFees {
    queryId: uint64
    pool: address
    recipient: address
}

/// Admin sets the jetton wallets of a pool (a null xWallet means the x side is
/// GRAM); the factory forwards it to the pool as CpmmPoolUpdateWallets.
struct (0x51d80e5e) CpmmFactoryUpdatePoolWallets {
    queryId: uint64
    pool: address
    xWallet: address?
    yWallet: address
}

/// Admin replaces the code and/or data of a pool; the factory forwards it to the
/// pool as CpmmPoolUpgrade.
struct (0x0d730f18) CpmmFactoryUpgradePool {
    queryId: uint64
    pool: address
    newCode: cell?
    newData: cell?
}

/// Admin replaces the factory code and/or data.
struct (0xbae2139a) CpmmFactoryUpgrade {
    queryId: uint64
    newCode: cell?
    newData: cell?
}

fun onInternalMessage(in: InMessage) {
}
