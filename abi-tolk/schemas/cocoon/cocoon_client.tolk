import "@stdlib/gas-payments"
import "opcodes"
import "errors"
import "messages"
import "storage"


contract cocoon_client {
    incomingMessages: ClientMessage,
    storage: ClientStorage,
    forceAbiExport: [
        ClientMessage
    ]
}


// === Client-specific Message Types ===

// Incoming message types
struct (0xbb63ff93) ExtClientChargeSigned {
    rest: RemainingBitsAndRefs
}

struct (0xefd711e1) ExtClientGrantRefundSigned {
    rest: RemainingBitsAndRefs
}

type SignedClientMessage = ExtClientChargeSigned | ExtClientGrantRefundSigned

type ClientSignedPayload = ChargePayload | GrantRefundPayload

type ClientMessage =
    | SignedClientMessage
    | ExtClientTopUp
    | OwnerClientChangeSecretHashAndTopUp
    | OwnerClientRegister
    | OwnerClientChangeSecretHash
    | OwnerClientIncreaseStake
    | OwnerClientWithdraw
    | OwnerClientRequestRefund

// === Helper Functions ===

fun sendProxyRequestCarryAll(storage: ClientStorage, queryId: uint64, payload: cell?) {
    val constData = storage.constDataRef.load();

    createMessage({
        bounce: BounceMode.Only256BitsOfBody,
        dest: constData.proxyAddress,
        value: 0,
        body: ClientProxyRequest {
            queryId: queryId,
            ownerAddress: constData.ownerAddress,
            stateData: storage.toStateData().toCell(),
            payload: payload
        }
    }).send(SEND_MODE_CARRY_ALL_BALANCE);
}

fun sendProxyRequestWithExcesses(storage: ClientStorage, queryId: uint64,
    payload: cell?, sendExcessesTo: address) {
    val constData = storage.constDataRef.load();

    createMessage({
        bounce: BounceMode.Only256BitsOfBody,
        dest: constData.proxyAddress,
        value: 1,
        body: ClientProxyRequest {
            queryId: queryId,
            ownerAddress: constData.ownerAddress,
            stateData: storage.toStateData().toCell(),
            payload: payload
        }
    }).send(SEND_MODE_PAY_FEES_SEPARATELY);

    sendExcesses(sendExcessesTo, queryId);
}

// === Message Handlers ===

fun ClientStorage.handleSignedClientMessage(mutate self, inBody: slice, inValueCoins: coins,
    constData: ClientConstData, toLeaveOnBalance: coins) {
    assert (COMMISSION_ESTIMATE <= inValueCoins) throw ERROR_LOW_MSG_VALUE;
    assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;

    val signedMsg = SignedMessage.fromSlice(inBody);
    signedMsg.verify(constData.proxyPublicKey);

    // Match on inner payload type
    val payload = ClientSignedPayload.fromCell(signedMsg.signedDataCell);

    match (payload) {
        ChargePayload => {
            assert (self.state != client_state_closed) throw ERROR_CLOSED;
            assert (self.tokensUsed < payload.newTokensUsed) throw ERROR_OLD_MESSAGE;
            assert (payload.expectedMyAddress == contract.getAddress()) throw ERROR_EXPECTED_MY_ADDRESS;

            self.updateBalance(payload.newTokensUsed);

            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestWithExcesses(self, payload.queryId, null, signedMsg.sendExcessesTo);
        },
        GrantRefundPayload => {
            assert ((self.state != client_state_closed) | (self.balance != 0)) throw ERROR_CLOSED;
            assert (self.tokensUsed <= payload.newTokensUsed) throw ERROR_OLD_MESSAGE;
            assert (payload.expectedMyAddress == contract.getAddress()) throw ERROR_EXPECTED_MY_ADDRESS;

            self.updateBalance(payload.newTokensUsed);

            self.state = client_state_closed;
            self.unlockTs = 0;
            val toRefund = self.balance;
            self.balance = 0;

            reserveToncoinsOnBalance(STORAGE_ESTIMATE, 0);

            sendProxyRequestCarryAll(
                self,
                payload.queryId,
                ClientProxyRefundGranted {
                    coins: toRefund,
                    sendExcessesTo: signedMsg.sendExcessesTo
                }.toCell()
            );
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

    var storage = ClientStorage.load();
    val constData = storage.constDataRef.load();

    // Match on message type
    val msg = lazy ClientMessage.fromSlice(in.body);

    match (msg) {
        ExtClientChargeSigned => {
            storage.handleSignedClientMessage(in.body, in.valueCoins, constData, toLeaveOnBalance);
        },
        ExtClientGrantRefundSigned => {
            storage.handleSignedClientMessage(in.body, in.valueCoins, constData, toLeaveOnBalance);
        },
        ExtClientTopUp => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (msg.topUpAmount + COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (msg.topUpAmount + COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;

            storage.balance = storage.balance + msg.topUpAmount;
            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestCarryAll(
                storage,
                msg.queryId,
                ClientProxyTopUp {
                    topUpCoins: msg.topUpAmount,
                    sendExcessesTo: msg.sendExcessesTo
                }.toCell()
            );
        },

        OwnerClientChangeSecretHashAndTopUp => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (msg.topUpAmount + COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (msg.topUpAmount + COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;

            storage.balance = storage.balance + msg.topUpAmount;
            storage.secretHash = msg.newSecretHash;
            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestCarryAll(
                storage,
                msg.queryId,
                ClientProxyTopUp {
                    topUpCoins: msg.topUpAmount,
                    sendExcessesTo: msg.sendExcessesTo
                }.toCell()
            );
        },

        OwnerClientRegister => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;

            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestWithExcesses(
                storage,
                msg.queryId,
                ClientProxyRegister {
                }.toCell(),
                msg.sendExcessesTo
            );
        },

        OwnerClientChangeSecretHash => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;

            storage.secretHash = msg.newSecretHash;
            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestWithExcesses(storage, msg.queryId, null, msg.sendExcessesTo);
        },

        OwnerClientIncreaseStake => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;
            assert (msg.newStake > storage.stake) throw ERROR_LOW_MSG_VALUE;

            storage.stake = msg.newStake;
            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestWithExcesses(storage, msg.queryId, null, msg.sendExcessesTo);
        },

        OwnerClientWithdraw => {
            assert (storage.state == client_state_normal) throw ERROR_CLOSED;
            assert (COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_OWNER;
            assert (storage.balance > storage.stake) throw ERROR_LOW_SMC_BALANCE;

            val toRefund = storage.balance - storage.stake;
            storage.balance = storage.stake;
            reserveToncoinsOnBalance(toLeaveOnBalance, 0);

            sendProxyRequestCarryAll(
                storage,
                msg.queryId,
                ClientProxyRefundGranted {
                    coins: toRefund,
                    sendExcessesTo: msg.sendExcessesTo
                }.toCell()
            );
        },

        OwnerClientRequestRefund => {
            assert (storage.state != client_state_closed) throw ERROR_CLOSED;
            assert (COMMISSION_ESTIMATE <= in.valueCoins) throw ERROR_LOW_MSG_VALUE;
            assert (COMMISSION_ESTIMATE + toLeaveOnBalance <= contract.getOriginalBalance()) throw ERROR_LOW_SMC_BALANCE;
            assert (in.senderAddress == constData.ownerAddress) throw ERROR_EXPECTED_OWNER;

            if (storage.state == client_state_normal) {
                val params = storage.params.load();
                storage.state = client_state_closing;
                storage.unlockTs = blockchain.now() + params.clientDelayBeforeClose;

                reserveToncoinsOnBalance(toLeaveOnBalance, 0);

                if (storage.balance > storage.stake) {
                    val toRefund = storage.balance - storage.stake;
                    storage.balance = storage.stake;

                    sendProxyRequestCarryAll(
                        storage,
                        msg.queryId,
                        ClientProxyRefundGranted {
                            coins: toRefund,
                            sendExcessesTo: msg.sendExcessesTo
                        }.toCell()
                    );
                } else {
                    sendProxyRequestWithExcesses(storage, msg.queryId, null, msg.sendExcessesTo);
                }
            } else {
                assert (storage.unlockTs < blockchain.now()) throw ERROR_NOT_UNLOCKED_YET;

                storage.state = client_state_closed;
                val toRefund = storage.balance;
                storage.balance = 0;

                reserveToncoinsOnBalance(toLeaveOnBalance, 0);

                sendProxyRequestCarryAll(
                    storage,
                    msg.queryId,
                    ClientProxyRefundForce {
                        coins: toRefund,
                        sendExcessesTo: msg.sendExcessesTo
                    }.toCell()
                );
            }
        }
    }

    storage.save();
}

// === Get Methods ===

struct CocoonClientData {
    ownerAddress: address
    proxyAddress: address
    proxyPublicKey: uint256
    state: uint2
    balance: coins
    stake: coins
    tokensUsed: uint64
    unlockTs: uint32
    secretHash: uint256
}

get fun get_cocoon_client_data(): CocoonClientData {
    val storage = ClientStorage.load();
    val constData = storage.constDataRef.load();

    return CocoonClientData {
        ownerAddress: constData.ownerAddress,
        proxyAddress: constData.proxyAddress,
        proxyPublicKey: constData.proxyPublicKey,
        state: storage.state,
        balance: storage.balance,
        stake: storage.stake,
        tokensUsed: storage.tokensUsed,
        unlockTs: storage.unlockTs,
        secretHash: storage.secretHash
    };
}
