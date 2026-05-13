// Common message structures for all Cocoon contracts

import "errors"

// === Common Messages ===

struct (0x2565934c) ReturnExcessesBack {
    queryId: uint64
}

struct (0xc59a7cd3) Payout {
    queryId: uint64
}

// === Worker <-> Proxy Messages ===

struct (0x4d725d2c) WorkerProxyRequest {
    queryId: uint64
    ownerAddress: address
    state: uint2
    tokens: uint64
    payload: cell?
}

struct (0x08e7d036) WorkerProxyPayoutRequest {
    workerPart: coins
    proxyPart: coins
    sendExcessesTo: address
}

// === Client <-> Proxy Messages ===

struct ClientStateData {
    state: uint2
    balance: coins
    stake: coins
    tokensUsed: uint64
    secretHash: uint256
}

struct (0x65448ff4) ClientProxyRequest {
    queryId: uint64
    ownerAddress: address
    stateData: Cell<ClientStateData>
    payload: cell?
}

struct (0x5cfc6b87) ClientProxyTopUp {
    topUpCoins: coins
    sendExcessesTo: address
}

struct (0xa35cb580) ClientProxyRegister {
    // empty payload
}

struct (0xc68ebc7b) ClientProxyRefundGranted {
    coins: coins
    sendExcessesTo: address
}

struct (0xf4c354c9) ClientProxyRefundForce {
    coins: coins
    sendExcessesTo: address
}

// === Signed Message Structures ===

// Common header inside all signed data
struct QueryHeader {
    op: uint32
    queryId: uint64
}

// Generic signed message wrapper (common for all signed messages)
struct SignedMessage {
    op: uint32
    queryId: uint64
    sendExcessesTo: address
    signature: bits512
    signedDataCell: cell  // Will match on contents after signature verification
}

// Verify signed message and check QueryHeader matches
fun SignedMessage.verify(self, publicKey: int) {
    val hash = self.signedDataCell.hash();
    assert (isSignatureValid(hash, self.signature as slice, publicKey)) throw ERROR_BAD_SIGNATURE;

    var payloadSlice = self.signedDataCell.beginParse();
    val header = QueryHeader.fromSlice(payloadSlice, {
        assertEndAfterReading: false
    });
    assert (header.op == self.op) throw ERROR_SIGNED_MSG_FORMAT_MISMATCH;
    assert (header.queryId == self.queryId) throw ERROR_SIGNED_MSG_FORMAT_MISMATCH;
}

// Common payload data (same for all payout types)
struct PayoutPayloadData {
    queryId: uint64
    newTokens: uint64
    expectedMyAddress: address
}

// Minimal payload types - just for opcode matching
struct (0xa040ad28) PayoutPayload {
    data: PayoutPayloadData
}

struct (0xf5f26a36) LastPayoutPayload {
    data: PayoutPayloadData
}

// === Root Messages ===

struct (0xe34b1c60) AddWorkerType {
    queryId: uint64
    workerHash: uint256
}

struct (0x8d94a79a) DelWorkerType {
    queryId: uint64
    workerHash: uint256
}

struct (0xc146134d) AddModelType {
    queryId: uint64
    modelHash: uint256
}

struct (0x92b11c18) DelModelType {
    queryId: uint64
    modelHash: uint256
}

struct (0x71860e80) AddProxyType {
    queryId: uint64
    proxyHash: uint256
}

struct (0x3c41d0b2) DelProxyType {
    queryId: uint64
    proxyHash: uint256
}

struct (0x927c7cb5) RegisterProxy {
    queryId: uint64
    proxyInfo: RemainingBitsAndRefs  // type + len + addr
}

struct (0x6d49eaf2) UnregisterProxy {
    queryId: uint64
    seqno: uint32
}

struct (0x9c7924ba) UpdateProxy {
    queryId: uint64
    seqno: uint32
    proxyAddr: RemainingBitsAndRefs  // type + len + addr
}

struct (0xc52ed8d4) ChangeFees {
    queryId: uint64
    pricePerToken: coins
    workerFeePerToken: coins
}

struct (0x022fa189) ChangeParams {
    queryId: uint64
    pricePerToken: coins
    workerFeePerToken: coins
    proxyDelayBeforeClose: uint32
    clientDelayBeforeClose: uint32
    minProxyStake: coins
    minClientStake: coins
}

struct (0xa2370f61) UpgradeContracts {
    queryId: uint64
    proxyCode: cell
    workerCode: cell
    clientCode: cell
}

struct (0x11aefd51) UpgradeCode {
    queryId: uint64
    newCode: cell
}

struct (0x563c1d96) ResetRoot {
    queryId: uint64
}

struct (0x4f7c5789) UpgradeFull {
    queryId: uint64
    newData: cell
    newCode: cell
}

struct (0xc4a1ae54) ChangeOwner {
    queryId: uint64
    newOwnerAddress: address
}

// === Worker Messages ===

struct (0x26ed7f65) OwnerWorkerRegister {
    queryId: uint64
    sendExcessesTo: address
}

// === Proxy Messages ===

struct (0x7610e6eb) ExtProxyPayoutRequest {
    queryId: uint64
    sendExcessesTo: address
}

struct (0x9713f187) ExtProxyIncreaseStake {
    queryId: uint64
    grams: coins
    sendExcessesTo: address
}

struct (0xb51d5a01) OwnerProxyClose {
    queryId: uint64
    sendExcessesTo: address
}

struct (0x636a4391) CloseRequestPayload {
    queryId: uint64
    expectedMyAddress: address
}

struct (0xe511abc7) CloseCompleteRequestPayload {
    queryId: uint64
    expectedMyAddress: address
}

// === Client Messages ===

struct (0xf172e6c2) ExtClientTopUp {
    queryId: uint64
    topUpAmount: coins
    sendExcessesTo: address
}

struct (0x8473b408) OwnerClientChangeSecretHashAndTopUp {
    queryId: uint64
    topUpAmount: coins
    newSecretHash: uint256
    sendExcessesTo: address
}

struct (0xc45f9f3b) OwnerClientRegister {
    queryId: uint64
    nonce: uint64
    sendExcessesTo: address
}

struct (0xa9357034) OwnerClientChangeSecretHash {
    queryId: uint64
    newSecretHash: uint256
    sendExcessesTo: address
}

struct (0x6a1f6a60) OwnerClientIncreaseStake {
    queryId: uint64
    newStake: coins
    sendExcessesTo: address
}

struct (0xda068e78) OwnerClientWithdraw {
    queryId: uint64
    sendExcessesTo: address
}

struct (0xfafa6cc1) OwnerClientRequestRefund {
    queryId: uint64
    sendExcessesTo: address
}

struct (0xbb63ff93) ChargePayload {
    queryId: uint64
    newTokensUsed: uint64
    expectedMyAddress: address
}

struct (0xefd711e1) GrantRefundPayload {
    queryId: uint64
    newTokensUsed: uint64
    expectedMyAddress: address
}

// === Helper Functions ===

@inline
fun sendExcesses(dest: address, queryId: int) {
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: dest,
        value: 0,
        body: ReturnExcessesBack {
            queryId: queryId
        }
    }).send(SEND_MODE_CARRY_ALL_BALANCE);
}

@inline
fun sendPayout(dest: address, coins: coins, queryId: uint64, mode: int) {
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: dest,
        value: coins,
        body: Payout{
            queryId: queryId
        }
    }).send(mode);
}

