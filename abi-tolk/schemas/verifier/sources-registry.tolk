import "imports/content"

// === Constants ===

const MIN_GRAMS_LOWER_BOUND = 65000000 // 0.065 TON

// === Storage ===

struct SourcesRegistryStorage {
    minGram: coins
    maxGram: coins
    adminAddress: address
    verifierRegistryAddress: address
    sourceItemCode: cell
}

// === Internal Messages ===

struct (0x000003EA) DeploySourceItemPayload {
    queryId: uint64
    verifierId: uint256
    verifiedCodeCellHash: uint256
    sourceContent: Cell<SourceContent>
}

struct (0x000007D3) ChangeVerifierRegistry {
    queryId: uint64
    newVerifierRegistry: address
}

struct (0x00000BBC) ChangeAdmin {
    queryId: uint64
    newAdmin: address
}

struct (0x00000FA5) SetSourceItemCode {
    queryId: uint64
    newSourceItemCode: cell
}

struct (0x0000138E) SetCode {
    queryId: uint64
    newCode: cell
}

struct (0x00001777) SetDeploymentCosts {
    queryId: uint64
    newMinGram: coins
    newMaxGram: coins
}

type SourcesRegistryInternalMessage = DeploySourceItemPayload | ChangeVerifierRegistry | ChangeAdmin | SetSourceItemCode | SetCode | SetDeploymentCosts

// === Contract ===

contract verifier_sources_registry {
    incomingMessages: SourcesRegistryInternalMessage,
    storage: SourcesRegistryStorage
    forceAbiExport: [SourcesRegistryInternalMessage]
}

// === Handlers ===

fun onInternalMessage(inMsg: InMessage) {
    // stub to declare outgoing message types

    // Deploy source item message (sent to the calculated source item address)
    createMessage({
        bounce: BounceMode.NoBounce,
        dest: inMsg.senderAddress,
        value: 0,
        body: DeploySourceItemPayload{
            queryId: 0,
            verifierId: 0,
            verifiedCodeCellHash: 0,
            sourceContent: SourceContent{version: 0, url: ""}.toCell() as Cell<SourceContent>
        }
    }).send(64);
}

fun onBouncedMessage(msg: InMessageBounced) {
}

// === Get Methods ===

struct DeploymentCosts {
    minGram: coins
    maxGram: coins
}

get fun get_source_item_address(verifierId: int, verifiedCodeCellHash: int): address {
    return address("0:0000000000000000000000000000000000000000000000000000000000000000");
}

get fun get_verifier_registry_address(): address {
    return address("0:0000000000000000000000000000000000000000000000000000000000000000");
}

get fun get_admin_address(): address {
    return address("0:0000000000000000000000000000000000000000000000000000000000000000");
}

get fun get_deployment_costs(): DeploymentCosts {
    return DeploymentCosts{minGram: 0, maxGram: 0};
}
