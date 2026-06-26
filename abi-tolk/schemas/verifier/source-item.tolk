import "imports/content"

// === Storage ===

struct SourceItemStorage {
    verifierId: uint256
    verifiedCodeCellHash: uint256
    sourceItemRegistry: address
    content: Cell<SourceContent>
}

// === Contract ===

contract verifier_source_item {
    storage: SourceItemStorage
    incomingMessages: SourceContent
}

// === Handlers ===

fun onInternalMessage(inMsg: InMessage) {
}

// === Get Methods ===

struct SourceItemData {
    verifierId: uint256
    verifiedCodeCellHash: uint256
    sourceItemRegistry: address
    content: Cell<SourceContent>
}

get fun get_source_item_data(): SourceItemData {
    return SourceItemData{
        verifierId: 0,
        verifiedCodeCellHash: 0,
        sourceItemRegistry: address("0:0000000000000000000000000000000000000000000000000000000000000000"),
        content: SourceContent{version: 0, url: ""}.toCell() as Cell<SourceContent>
    };
}
