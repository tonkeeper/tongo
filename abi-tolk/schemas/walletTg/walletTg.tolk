tolk 1.4

/*
    ABI of the Telegram wallet contract (WalletTg),
    https://github.com/ton-blockchain/tg-wallet-contract.

    Every Telegram wallet account holds an immutable trampoline whose only job is
    `JMP to config[-123]`, where the real WalletTg bytecode lives, shared by all
    Telegram users and updated by validator voting.

    This is a stub: types, the contract declaration and get method signatures only.
    The bodies are mocked, the sources of truth are the upstream contract sources.
*/

contract walletTg {
    storage: Storage,
    incomingMessages: SignedInternalRequest,
    incomingExternal: SignedExternalRequest,
    forceAbiExport: [
        KeyRotationProofPayload
    ]
}

/// Current storage schema of the contract.
/// A new contract revision bumps its 0xXX prefix and may append new fields;
/// c4 of an account keeps the previous shape until it is migrated
/// on the first valid incoming message after a bytecode upgrade.
struct (0x00) Storage {
    /// Seqno protects against replay attacks: every incoming message contains the expected seqno.
    seqno: uint32
    /// SubwalletId allows creating multiple wallets for one private key.
    /// Default for mainnet: 0x7FFF7F11; for testnet: 0x7FFF7FFD.
    subwalletId: uint32
    /// Public key used to check signatures of all incoming messages.
    publicKey: uint256
}

/// Every time the contract logic is updated, the 0xXX prefix of `Storage` is bumped.
enum Revision : uint8 {
    Rev00_Initial = 0,
}

/// WalletTg can send up to 255 out messages.
/// Note: in contrast with w5, which accepts a well-formed c5 register, WalletTg accepts messages.
struct MessageToSend {
    sendMode: uint8
    messageCell: cell
}

/// Bulk requests carry `array<MessageToSend>`, a standard-encoded Tolk array:
/// `len: uint8` followed by a maybe-ref to the first chunk. Every chunk starts with
/// a maybe-ref to the next one, so 3 items fit in an intermediate chunk and 4 in the last one.
type RawArrayOfMessagesToSend = array<MessageToSend>

/// A payload signed by a new private key for the "change key" request.
struct KeyRotationProofPayload {
    tag: uint96 = 0x4B45595F524F544154494F4E
    walletWorkchain: int8
    walletAddrHash: uint256
}

/// Every internal/external message is signed by a private key and checked against `storage.publicKey`.
/// The client-side signing process has three steps.
/// 1. Compose some valid request (body).
/// 2. Sign it with their private key — resulting in `signature`.
/// 3. Prepend the signature to the serialized request and send.
///
/// Note that signing DIFFERS FROM WALLET V5, where the signature lies at the end.
///
/// Note: upstream this is one generic `SignedRequest<T>`;
/// here it is specialized, because ABI codegen needs concrete message types.
struct SignedInternalRequest {
    signature: bits512
    request: AllowedInternalMsg
}

struct SignedExternalRequest {
    signature: bits512
    request: AllowedExternalMsg
}

type AllowedInternalMsg =
    | SendOneMessageRequestI
    | SendBulkMessagesRequestI
    | ChangePublicKeyRequestI

type AllowedExternalMsg =
    | SendOneMessageRequestE
    | SendBulkMessagesRequestE
    | ChangePublicKeyRequestE

/// Every internal/external request starts with a header.
struct SeqnoHeader {
    /// Must match `Storage.subwalletId`, so one key can control multiple wallets.
    subwalletId: uint32
    /// Unix timestamp until which the message is valid. Must be greater than `blockchain.now`.
    validUntil: uint32
    /// Must match `Storage.seqno`.
    seqno: uint32
}

/// Check the signature and forward ONE message.
/// `MessageToSend` is stored inline, reducing import fee and gas.
/// Internal (not external) is used for gasless.
struct (0x63896E74) SendOneMessageRequestI {
    header: SeqnoHeader
    msg: MessageToSend
}

/// Differs from the "internal" struct only in opcode to prevent cross-channel replay.
/// External requests must set the ignore-errors send mode.
struct (0x63896E75) SendOneMessageRequestE {
    header: SeqnoHeader
    msg: MessageToSend
}

/// Check the signature and forward MANY messages (1-255).
struct (0x73896E74) SendBulkMessagesRequestI {
    header: SeqnoHeader
    msgArr: RawArrayOfMessagesToSend
}

struct (0x73896E75) SendBulkMessagesRequestE {
    header: SeqnoHeader
    msgArr: RawArrayOfMessagesToSend
}

/// Change the wallet signing key. The request is signed twice:
/// by the current key (as any request) and by the new key (`rotationSignature`,
/// a signature of `KeyRotationProofPayload`).
struct (0xFBBA99C7) ChangePublicKeyRequestI {
    header: SeqnoHeader
    newPublicKey: uint256
    rotationSignature: Cell<bits512>
}

struct (0xFBBA99C8) ChangePublicKeyRequestE {
    header: SeqnoHeader
    newPublicKey: uint256
    rotationSignature: Cell<bits512>
}

/// Error codes are identical to w5 for common failures.
enum Errors {
    BeginParseOfNull = 7,
    InvalidSeqno = 133,
    InvalidSubwalletId = 134,
    InvalidSignature = 135,
    Expired = 136,
    ExternalSendMessageMustHaveIgnoreErrorsSendMode = 137,
    InvalidOutArrayLen = 147,
    NewKeyMustDiffer = 148,
    NewKeySignatureIsInvalid = 149,
    InvalidOpcode = 0xFFFF,
}

/// Mocked entrypoints: WalletTg sends only the messages a signed request carries,
/// so there is nothing to describe here besides the incoming ones.
fun onInternalMessage(in: InMessage) {
    // mocked: every error code is referenced, so that all of them get into the ABI
    val n = in.body.remainingBitsCount();
    if (n == 0) { throw Errors.BeginParseOfNull; }
    if (n == 1) { throw Errors.InvalidSeqno; }
    if (n == 2) { throw Errors.InvalidSubwalletId; }
    if (n == 3) { throw Errors.InvalidSignature; }
    if (n == 4) { throw Errors.Expired; }
    if (n == 5) { throw Errors.ExternalSendMessageMustHaveIgnoreErrorsSendMode; }
    if (n == 6) { throw Errors.InvalidOutArrayLen; }
    if (n == 7) { throw Errors.NewKeyMustDiffer; }
    if (n == 8) { throw Errors.NewKeySignatureIsInvalid; }
    if (n == 9) { throw Errors.InvalidOpcode; }
}

fun onExternalMessage(inMsgBody: slice) {
}

/// Returns current storage revision. It does not trigger migration.
get fun revision(): Revision {
    throw 0;
}

/// Returns current seqno from the storage. Bumped on every valid incoming message.
get fun seqno(): uint32 {
    throw 0;
}

/// Returns the subwallet id from the storage.
get fun get_subwallet_id(): uint32 {
    throw 0;
}

/// Returns the public key corresponding to the signer's private key.
get fun get_public_key(): uint256 {
    throw 0;
}
