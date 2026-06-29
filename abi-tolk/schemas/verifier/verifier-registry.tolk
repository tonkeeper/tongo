import "@stdlib/tvm-dicts"

// === Constants ===

const MAX_VERIFIERS = 20

// === Types ===

// SignatureItem represents one entry in the signature linked list.
// Wire format: 512 bits signature + 256 bits pub_key, optional ref to next.
// The full list is stored as a cell (custom marshaling required for the linked list).
struct SignatureItem {
    signature: bits512
    pubKey: uint256
}

struct MessageDescription {
    verifierId: uint256
    validUntil: uint32
    sourceAddr: any_address
    targetAddr: any_address
    msg: cell
}

struct VerifierSettings {
    multiSigThreshold: uint8
    pubKeyEndpoints: map<uint256, uint32>
    name: string
    marketingUrl: string
}

struct Verifier {
    admin: any_address
    settings: VerifierSettings
}

// === Storage ===

struct VerifierRegistryStorage {
    verifiers: map<uint256, Verifier>
}

// === Internal Messages ===

struct (0x75217758) ForwardMessage {
    queryId: uint64
    msg: Cell<MessageDescription>
    signatures: cell  // linked list of SignatureItem (custom marshaling)
}

struct (0x6002d61a) UpdateVerifier {
    queryId: uint64
    verifierId: uint256
    settings: VerifierSettings
}

struct (0x19fa5637) RemoveVerifier {
    queryId: uint64
    id: uint256
}

type VerifierRegistryInternalMessage = ForwardMessage | UpdateVerifier | RemoveVerifier

// === Contract ===

contract verifier_registry {
    incomingMessages: VerifierRegistryInternalMessage,
    storage: VerifierRegistryStorage
    forceAbiExport: [VerifierRegistryInternalMessage],
}

// === Handlers ===

fun onInternalMessage(inMsg: InMessage) {
    // stub to declare outgoing message types

    // Text comment message (sent on verifier remove/update)
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: inMsg.senderAddress,
        value: 0,
        body: beginCell().endCell()
    }).send(SEND_MODE_CARRY_ALL_BALANCE);

    // Forwarded message to target (sent on forward_message processing)
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: inMsg.senderAddress,
        value: 0,
        body: beginCell().endCell()
    }).send(64);
}

fun onBouncedMessage(msg: InMessageBounced) {
}

// === Get Methods ===

struct VerifierInfo {
    admin: address
    settings: cell
    found: bool
}

get fun get_verifier(id: int): VerifierInfo {
    return VerifierInfo{
        admin: address("0:0000000000000000000000000000000000000000000000000000000000000000"),
        settings: beginCell().endCell(),
        found: false
    };
}

get fun get_verifiers_num(): int {
    return 0;
}

get fun get_verifiers(): Cell<VerifierRegistryStorage> {
    throw 0;;
}
