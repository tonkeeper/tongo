import "@stdlib/gas-payments"
import "errors"
import "opcodes"
import "storage"


contract cocoon_wallet {
    incomingMessages: AllowedInternalMessage,
    incomingExternal: WalletExternalMessage,
    storage: WalletStorage
}


// === Wallet-specific Constants and Types ===

const SIZE_SIGNATURE = 512

struct WalletExternalMessage {
    signature: bits512
    message: ExternalSignedMessage
}

// External message struct (after signature)
struct ExternalSignedMessage {
    subwalletId: uint32
    validUntil: uint32
    msgSeqno: uint32
    forward: ForwardMsgs  // Contains mode + message refs
}

type ForwardMsgs = array<ForwardMsg>

fun ForwardMsgs.unpackFromSlice(mutate s: slice): ForwardMsgs {
    return []
}

fun ForwardMsgs.packToBuilder(self, mutate b: builder) {
    // stub
}

struct ForwardMsg {
    mode: uint8
    msg: cell
}

// Internal message structs
struct (0x9c69f376) OwnerWalletSendMessage {
    queryId: uint64
    mode: uint8
    body: cell
}

// Text command markers (op == 0)
struct (0x00000000) TextCommand {
    action: uint8
}

// Allowed internal messages
type AllowedInternalMessage = OwnerWalletSendMessage | TextCommand

// === External Message Handler (signed by owner) ===

fun onExternalMessage(inMsg: slice) {
    val signature = inMsg.loadBits(SIZE_SIGNATURE);
    val msg = ExternalSignedMessage.fromSlice(inMsg);

    assert (msg.validUntil > blockchain.now()) throw 1030;

    var storage = WalletStorage.load();

    assert (!(storage.status & 1)) throw 1031;  // Not blocked
    assert (msg.msgSeqno == storage.seqno) throw 1032;
    assert (msg.subwalletId == storage.subwalletId) throw 1033;
    assert (isSignatureValid(inMsg.hash(), signature, storage.publicKey)) throw 1034;
    assert (contract.getOriginalBalance() >= ton("2.0")) throw 1035;

    acceptExternalMessage();

    // Send all messages from the list
    // var remaining = msg.rest;
    // while (remaining.remainingRefsCount() > 0) {
    //     val mode = remaining.loadUint(8);
    //     sendRawMessage(remaining.loadRef(), mode);
    // }

    // Update seqno
    storage.seqno = storage.seqno + 1;
    storage.save();
}

// === Internal Message Handler ===

fun onInternalMessage(in: InMessage) {
    if (in.body.isEmpty()) {
        return;
    }

    val op = in.body.preloadUint(32);

    // Ignore these opcodes
    if (op == OP_EXCESSES || op == OP_PAYOUT || op == OP_DO_NOT_PROCESS) {
        return;
    }

    var storage = lazy WalletStorage.load();

    // Only owner can send internal messages
    assert (in.senderAddress == storage.ownerAddress) throw 1040;

    val msg = lazy AllowedInternalMessage.fromSlice(in.body);

    match (msg) {
        OwnerWalletSendMessage => {
            sendRawMessage(msg.body, msg.mode);
        },

        TextCommand => {
            match (msg.action) {
                119 => {
                    // "w" - withdraw
                    reserveToncoinsOnBalance(ton("0.01"), RESERVE_MODE_AT_MOST);  // Keep contract alive

                    createMessage({
                        bounce: BounceMode.Only256BitsOfBody,
                        dest: in.senderAddress,
                        value: 0,
                        body: TextCommand{
                            action: 119
                        }
                    }).send(SEND_MODE_CARRY_ALL_BALANCE);
                },
                98 => {
                    // "b" - block
                    storage.status = storage.status | 1;
                    storage.save();
                },
                117 => {
                    // "u" - unblock
                    storage.status = storage.status & 0xfffffffe;
                    storage.save();
                }
                else => {
                   throw 1041;
                }
            }
        }
    }
}

// === Get Methods ===

get fun seqno(): int {
    val storage = lazy WalletStorage.load();
    return storage.seqno;
}

get fun get_public_key(): int {
    val storage = lazy WalletStorage.load();
    return storage.publicKey;
}

get fun get_owner_address(): address {
    val storage = lazy WalletStorage.load();
    return storage.ownerAddress;
}
