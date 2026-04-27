import "@stdlib/gas-payments"
import "opcodes"
import "errors"
import "messages"
import "storage"

// === Worker-specific Message Types ===

// Signed message markers (for opcode matching - same opcodes as inner payloads)
struct (0xa040ad28) ExtWorkerPayoutRequestSigned {
    rest: RemainingBitsAndRefs
}

struct (0xf5f26a36) ExtWorkerLastPayoutRequestSigned {
    rest: RemainingBitsAndRefs
}

type SignedWorkerMessage = ExtWorkerPayoutRequestSigned | ExtWorkerLastPayoutRequestSigned

type SignedPayload = PayoutPayload | LastPayoutPayload

type WorkerMessage = SignedWorkerMessage | OwnerWorkerRegister

// === Message Handlers ===

fun handleWorkerPayout(storage: WorkerStorage, signedMsg: SignedMessage, toLeaveOnBalance: coins, data: PayoutPayloadData, isLast: bool) {
    assert (data.expectedMyAddress == contract.getAddress()) throw ERROR_EXPECTED_MY_ADDRESS;

    if (isLast) {
        assert (storage.tokens <= data.newTokens) throw ERROR_OLD_MESSAGE;
        reserveToncoinsOnBalance(STORAGE_ESTIMATE, RESERVE_MODE_EXACT_AMOUNT);
        storage.state = worker_state_closed;
    } else {
        assert (storage.tokens < data.newTokens) throw ERROR_OLD_MESSAGE;
        reserveToncoinsOnBalance(toLeaveOnBalance, RESERVE_MODE_EXACT_AMOUNT);
    }

    val params = lazy storage.params.load();
    val (workerPart, proxyPart) = params.calculateEarnings(data.newTokens - storage.tokens);

    createMessage({
        bounce: BounceMode.Only256BitsOfBody,
        dest: storage.proxyAddress,
        value: 0,
        body: WorkerProxyRequest {
            queryId: signedMsg.queryId,
            ownerAddress: storage.ownerAddress,
            state: storage.state,
            tokens: data.newTokens,
            payload: WorkerProxyPayoutRequest {
                workerPart: workerPart,
                proxyPart: proxyPart,
                sendExcessesTo: signedMsg.sendExcessesTo,
            }.toCell(),
        }
    }).send(SEND_MODE_CARRY_ALL_BALANCE);

    storage.tokens = data.newTokens;
    storage.save();
}

fun handleSignedWorkerMessage(inBody: slice, storage: WorkerStorage, toLeaveOnBalance: coins) {
    val signedMsg = SignedMessage.fromSlice(inBody);
    signedMsg.verify(storage.proxyPublicKey);

    // Match on inner payload type
    val payload = lazy SignedPayload.fromCell(signedMsg.signedDataCell);

    match (payload) {
        PayoutPayload => {
            handleWorkerPayout(storage, signedMsg, toLeaveOnBalance, payload.data, false)
        },
        LastPayoutPayload => {
            handleWorkerPayout(storage, signedMsg, toLeaveOnBalance, payload.data, true)
        }
    }
}

// === Main Message Handler ===

fun onInternalMessage(in: InMessage) {
    if (in.body.isEmpty()) {
        return;
    }

    val toLeaveOnBalance = max(
        contract.getOriginalBalance() - in.valueCoins + contract.getStorageDuePayment(),
        COMMISSION_ESTIMATE + STORAGE_ESTIMATE
    );

    val op = in.body.preloadUint(32);
    if (op == OP_DO_NOT_PROCESS) {
        return;
    }

    var storage = lazy WorkerStorage.load();
    assert (storage.state != worker_state_closed) throw ERROR_CLOSED;
    assert (in.valueCoins > COMMISSION_ESTIMATE) throw ERROR_LOW_MSG_VALUE;

    val msg = lazy WorkerMessage.fromSlice(in.body);

    match (msg) {
        ExtWorkerPayoutRequestSigned => {
            handleSignedWorkerMessage(in.body, storage, toLeaveOnBalance);
        },
        ExtWorkerLastPayoutRequestSigned => {
            handleSignedWorkerMessage(in.body, storage, toLeaveOnBalance);
        },
        OwnerWorkerRegister => {
            assert (in.senderAddress == storage.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;

            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            // Send registration to proxy
            createMessage({
                bounce: BounceMode.Only256BitsOfBody,
                dest: storage.proxyAddress,
                value: 1,
                body: WorkerProxyRequest {
                    queryId: msg.queryId,
                    ownerAddress: storage.ownerAddress,
                    state: storage.state,
                    tokens: storage.tokens,
                    payload: null,
                }
            }).send(SEND_MODE_PAY_FEES_SEPARATELY);

            sendExcesses(msg.sendExcessesTo, msg.queryId);
        }
    }
}

// === Get Methods ===

struct CocoonWorkerData {
    ownerAddress: address
    proxyAddress: address
    proxyPublicKey: uint256
    state: uint2
    tokens: uint64
}

get fun get_cocoon_worker_data(): CocoonWorkerData {
    var storage = WorkerStorage.load();
    return CocoonWorkerData {
        ownerAddress: storage.ownerAddress,
        proxyAddress: storage.proxyAddress,
        proxyPublicKey: storage.proxyPublicKey,
        state: storage.state,
        tokens: storage.tokens
    };
}
