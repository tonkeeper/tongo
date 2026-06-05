// StonFi Escrow — shared data types and message structs.
// Contracts: escrow_factory, escrow_position, escrow_vault.

// ============================================================
// Shared data types (no opcode prefix)
// ============================================================

struct ForwardParams {
    value: coins
    successPayload: cell?
    rejectPayload: cell?
}

struct LockForwardParams {
    dest: any_address
    forwardParams: ForwardParams
}

struct DutchSegment {
    askEnd: coins
    duration: uint16
}

struct DutchSegments {
    num: uint4
    segments: RemainingBitsAndRefs
}

struct DutchAuctionData {
    dutchRemaining: coins
    resolverPublicKey: uint256?
    pendingUserReceiveAmount: coins
    dutchSegments: Cell<DutchSegments>
}

struct BilateralConditionState {
    initBidAmount: coins
    resolver: any_address
    resolverTimeout: uint64
    resolverAskAmount: coins
    auctionData: Cell<DutchAuctionData>
}

struct BilateralLockArgs {
    resolver: any_address
    resolverTimeoutDelta: uint64
    resolverAskAmount: coins
    resolverPublicKey: uint256?
    dutchSegments: Cell<DutchSegments>
}

struct BilateralSignedMessage {
    expirationTime: uint64
    resolver: any_address
}

struct BilateralSignPayload {
    signature: bits512
    message: Cell<BilateralSignedMessage>
}

struct BilateralUnlockArgs {
    minOut: coins
    ignoreRefundAmount: coins
    signed: BilateralSignPayload?
}

struct DutchPriceAtArg {
    timestamp: uint64?
}

struct BidPaymentData {
    usingVault: bool
    recipient: any_address
    bidJettonWallet: any_address
    amount: coins
    forwardParams: Cell<ForwardParams>?
}

struct AskRefundData {
    refundTo: any_address
    usingVault: bool
    amount: coins
}

struct UserPaymentData {
    user: any_address
    userReceiveJettonWallet: any_address
    userReceiveAmount: coins
    userDepositViaVault: bool
    userForwardParams: Cell<ForwardParams>?
    userExcesses: any_address
    userSafeDepositAndForwardValue: uint64
}

struct MinterRefundRequestExtraFields {
    jettonWallet: any_address
    excesses: any_address
    forwardTonAmount: coins
    forwardPayload: cell?
    lockRejectDest: any_address?
}

struct VaultLockAdditionalDataMore {
    orderOwner: any_address
    refundTo: any_address
    askJettonMinter: any_address
}

struct VaultLockAdditionalData {
    excesses: any_address
    tonSafeDeposit: coins
    userFillToVault: bool
    askJettonWallet: any_address
    ownerPubkey: uint256?
    more: Cell<VaultLockAdditionalDataMore>
    userUnlockForwardParams: Cell<ForwardParams>?
    lockForwardParams: Cell<LockForwardParams>?
}

struct InternalLockExtraNested {
    askJettonWallet: any_address
    ownerPubkey: uint256?
}

struct InternalLockExtra {
    quoteId: uint256
    refFee: any_address
    refFeeTier: uint16
    excesses: any_address
    safeDepositAndForwardValue: coins
    refundToVault: bool
    userFillToVault: bool
    more: Cell<VaultLockAdditionalDataMore>
    userUnlockForwardParams: Cell<ForwardParams>?
    lockForwardParams: Cell<LockForwardParams>?
    nested: Cell<InternalLockExtraNested>
}

struct UnlockAdditionalData {
    quoteId: uint256
    fillToVault: bool
    refundToVault: bool
    recipient: any_address
    refundTo: any_address
    unlockArgs: Cell<BilateralUnlockArgs>
    forwardParams: Cell<ForwardParams>?
}

struct MinterLockPayloadExtraFieldsMore {
    orderOwner: any_address
    askJettonMinter: any_address
}

struct MinterLockPayloadExtraFields {
    askJettonWallet: any_address
    refundTo: any_address
    ownerPubkey: uint256?
    more: Cell<MinterLockPayloadExtraFieldsMore>
    userUnlockForwardParams: Cell<ForwardParams>?
    lockForwardParams: Cell<LockForwardParams>?
}

struct UnlockPayloadExtraFields {
    recipient: any_address
    refundTo: any_address
    fillToVault: bool
    refundToVault: bool
}

struct ItemAdditionalFieldMore {
    askJettonMinter: any_address
    orderOwner: any_address
    refundTo: any_address
}

struct ItemAdditionalField {
    askJettonWallet: any_address
    refFee: any_address
    refFeeTier: uint16
    safeDepositAndForwardValue: coins
    userFillToVault: bool
    ownerPubkey: uint256?
    more: Cell<ItemAdditionalFieldMore>
    userUnlockForwardParams: Cell<ForwardParams>?
    lockForwardParams: Cell<LockForwardParams>?
}

struct ItemInternalUnlockExtraFields {
    recipient: any_address
    refundTo: any_address
    excesses: any_address
}

struct ExternalItemWithdrawPayload {
    escrowItem: any_address
    signatureTTL: uint64
    withdrawArgs: cell
}

struct EscrowWithdrawSignMessage {
    schemaHash: uint32
    timestamp: uint64
    userWalletAddress: any_address
    domain: cell
    payload: Cell<ExternalItemWithdrawPayload>
}

struct VaultUnlockExtraFields {
    recipient: any_address
    refundTo: any_address
}

// ============================================================
// Message structs — Internal (shared by several contracts)
// ============================================================

// ——— Factory incoming ———

struct (0x4716dd16) MinterInitTransfer {
    queryId: uint64
    quoteId: uint256
    transferType: uint32
    refFee: address?
    refFeeTier: uint16
    excesses: address?
    bidPaymentData: Cell<BidPaymentData>
    askRefundData: Cell<AskRefundData>
    userPaymentData: Cell<UserPaymentData>
}

struct (0xd714500d) MinterRefundRequest {
    queryId: uint64
    quoteId: uint256
    amount: coins
    useVault: bool
    recipient: address?
    exitCode: uint32
    prevMessage: uint32
    extraFields: Cell<MinterRefundRequestExtraFields>
}

struct (0x65cb4855) MinterInternalLock {
    queryId: uint64
    owner: address?
    jettonWallet: address?
    amount: coins
    lockArgs: Cell<BilateralLockArgs>
    unlockCondition: cell
    additionalData: Cell<InternalLockExtra>
}

struct (0x30fe22ab) MinterInternalUnlock {
    queryId: uint64
    owner: address?
    jettonWallet: address?
    amount: coins
    excesses: address?
    additionalData: Cell<UnlockAdditionalData>
}

struct (0x2c3b4692) MinterInternalWithdrawTokens {
    queryId: uint64
    jettonWallet: address?
    collectedFees: coins
    owner: address?
    excesses: address?
    forwardParams: Cell<ForwardParams>?
}

struct (0x5df8fa35) MinterGiveProtocolOwnership {
    queryId: uint64
    newProtocolAddress: address?
    excesses: address?
}

struct (0xde34aedc) MinterUpdateProtocolTier {
    queryId: uint64
    newTier: uint16
    excesses: any_address
}

struct (0x88124614) MinterTakeProtocolOwnership {
    queryId: uint64
    excesses: any_address
}

struct (0x07f886e6) MinterCancelNextProtocolOwner {
    queryId: uint64
    excesses: any_address
}

struct (0xfe16529e) MinterResetGas {
    queryId: uint64
    recipient: any_address
}

// ——— Jetton payloads ———

struct (0x3e246684) MinterLockPayload {
    quoteId: uint256
    refundToVault: bool
    refFee: address?
    refFeeTier: uint16
    excesses: address?
    tonSafeDeposit: coins
    userFillToVault: bool
    lockArgs: Cell<BilateralLockArgs>
    unlockCondition: cell
    extraFields: Cell<MinterLockPayloadExtraFields>
}

struct (0x6a58f85c) MinterUnlockPayload {
    quoteId: uint256
    excesses: any_address
    extraFields: Cell<UnlockPayloadExtraFields>
    unlockArgs: Cell<BilateralUnlockArgs>
    forwardParams: Cell<ForwardParams>?
}

struct (0x702b4dc8) MinterDepositVault {
    excesses: any_address
    forwardParams: Cell<ForwardParams>?
}

// ——— Position incoming ———

struct (0x7ab06181) ItemInternalLock {
    queryId: uint64
    tokenAddress: any_address
    amount: coins
    excesses: any_address
    refundToVault: bool
    additionalField: Cell<ItemAdditionalField>
    lockArgs: Cell<BilateralLockArgs>
    unlockCondition: cell
}

struct (0x696fa2d8) ItemInternalUnlock {
    queryId: uint64
    fillToVault: bool
    refundToVault: bool
    resolver: any_address
    resolverSentJettonWallet: any_address
    resolverSentAmount: coins
    unlockArgs: Cell<BilateralUnlockArgs>
    extraFields: Cell<ItemInternalUnlockExtraFields>
    forwardParams: Cell<ForwardParams>?
}

struct (0x12b8a987) ItemWithdraw {
    queryId: uint64
    withdrawArgs: cell
    excesses: any_address
    forwardParams: Cell<ForwardParams>?
}

// ——— Position external ———

struct (0x8af982e1) ExternalItemWithdraw {
    queryId: uint64
    signature: bits512
    signMessage: Cell<EscrowWithdrawSignMessage>
}

struct (0x2114702d) ExternalCronTrigger {
    rewardAddress: any_address
    salt: uint32
}

// ——— Vault incoming ———

struct (0x24864ed4) VaultLock {
    queryId: uint64
    amount: coins
    quoteId: uint256
    refundToVault: bool
    refFee: any_address
    refFeeTier: uint16
    additionalData: Cell<VaultLockAdditionalData>
    lockArgs: Cell<BilateralLockArgs>
    unlockCondition: cell
}

struct (0x4bd679ff) VaultUnlock {
    queryId: uint64
    amount: coins
    quoteId: uint256
    fillToVault: bool
    refundToVault: bool
    unlockArgs: Cell<BilateralUnlockArgs>
    excesses: any_address
    extraFields: Cell<VaultUnlockExtraFields>
    unlockForwardParams: Cell<ForwardParams>?
}

struct (0x5809e482) VaultWithdrawTokens {
    queryId: uint64
    amount: coins
    excesses: any_address
    withdrawForwardParams: Cell<ForwardParams>?
}

struct (0x555edf4b) VaultDepositTokens {
    queryId: uint64
    amount: coins
    excesses: any_address
    forwardTonAmount: coins
    forwardPayload: cell?
}

// ——— Outgoing messages ———

struct (0x14b8bb09) ItemLockSuccess {
    queryId: uint64
    forwardPayload: cell?
}

struct (0xac8a7373) LockRejectNotification {
    queryId: uint64
    forwardPayload: cell?
}

struct (0x44aceac4) VaultDepositNotification {
    queryId: uint64
    amount: coins
    forwardPayload: cell
}
