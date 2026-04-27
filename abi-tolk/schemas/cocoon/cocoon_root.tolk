
import "@stdlib/gas-payments"
import "@stdlib/tvm-dicts"
import "opcodes"
import "errors"
import "messages"
import "storage"


contract cocoon_root {
    incomingMessages: RootMessage,
    storage: RootStorage
}

// === Root-specific Union Types ===

type RootMessage =
    | AddWorkerType
    | DelWorkerType
    | AddModelType
    | DelModelType
    | AddProxyType
    | DelProxyType
    | RegisterProxy
    | UnregisterProxy
    | UpdateProxy
    | ChangeFees
    | ChangeParams
    | UpgradeContracts
    | UpgradeCode
    | ResetRoot
    | UpgradeFull
    | ChangeOwner

// === Helper Functions ===

fun validateProxyAddr(addr: slice) {
    // Validate: type(1) + len(7) + addr(len*8)
    // Clone the slice to avoid consuming it
    var tempAddr = addr;
    assert (tempAddr.loadUint(1) == 0) throw ERROR_MSG_FORMAT_MISMATCH;
    val addrLen = tempAddr.loadUint(7);
    tempAddr.skipBits(addrLen * 8);
    tempAddr.assertEnd();
}

// === Message Handler ===

fun onInternalMessage(in: InMessage) {
    if (in.body.isEmpty()) {
        return;
    }

    val toLeaveOnBalance = max(
        contract.getOriginalBalance() - in.valueCoins + contract.getStorageDuePayment(),
        STORAGE_ESTIMATE
    );
    reserveToncoinsOnBalance(toLeaveOnBalance, RESERVE_MODE_EXACT_AMOUNT);

    val op = in.body.preloadUint(32);
    if (op == OP_DO_NOT_PROCESS) {
        return;
    }
    val queryHeader = lazy QueryHeader.fromSlice(in.body, {assertEndAfterReading: false});
    val queryId = queryHeader.queryId;

    var storage = lazy RootStorage.load();
    val ownerAddress = storage.ownerAddress;

    // Only owner can perform operations
    assert (in.senderAddress == storage.ownerAddress) throw ERROR_EXPECTED_MESSAGE_FROM_OWNER;

    var data = lazy storage.data.load();
    var params = lazy storage.params.load();

    // Match on message type
    val msg = lazy RootMessage.fromSlice(in.body);


    match (msg) {
        AddWorkerType => {
            // data.workerHashes.set(256, msg.workerHash, ());
        },

        DelWorkerType => {
            // data.workerHashes.uDictDelete(256, msg.workerHash);
        },

        AddModelType => {
            // data.modelHashes.uDictSet(256, msg.modelHash, createEmptySlice());
        },

        DelModelType => {
            // data.modelHashes.uDictDelete(256, msg.modelHash);
        },

        AddProxyType => {
            // data.proxyHashes.uDictSet(256, msg.proxyHash, createEmptySlice());
        },

        DelProxyType => {
            // data.proxyHashes.uDictDelete(256, msg.proxyHash);
        },

        RegisterProxy => {
            // Validate and extract hash from proxyInfo
            validateProxyAddr(msg.proxyInfo);

            data.lastProxySeqno = data.lastProxySeqno + 1;
            // data.registeredProxies.uDictSet(32, data.lastProxySeqno, msg.proxyInfo);
        },

        UnregisterProxy => {
            // data.registeredProxies.uDictDelete(32, msg.seqno);
        },

        UpdateProxy => {
            // Validate new address format
            validateProxyAddr(msg.proxyAddr);

            // Get existing entry and extract hash
            // val (value, found) = data.registeredProxies.uDictGet(32, msg.seqno);
            assert (true) throw ERROR_UNKNOWN_PROXY_TYPE;

            // Create new entry
            // data.registeredProxies.uDictSet(32, msg.seqno, msg.proxyAddr);
        },

        ChangeFees => {
            params.paramsVersion = params.paramsVersion + 1;
            params.pricePerToken = msg.pricePerToken;
            params.workerFeePerToken = msg.workerFeePerToken;
        },

        ChangeParams => {
            params.paramsVersion = params.paramsVersion + 1;
            params.pricePerToken = msg.pricePerToken;
            params.workerFeePerToken = msg.workerFeePerToken;
            params.proxyDelayBeforeClose = msg.proxyDelayBeforeClose;
            params.clientDelayBeforeClose = msg.clientDelayBeforeClose;
            params.minProxyStake = msg.minProxyStake;
            params.minClientStake = msg.minClientStake;
        },

        UpgradeContracts => {
            params.paramsVersion = params.paramsVersion + 1;
            params.proxyScCode = msg.proxyCode;
            params.workerScCode = msg.workerCode;
            params.clientScCode = msg.clientCode;
        },

        UpgradeCode => {
            contract.setCodePostponed(msg.newCode);
            return;  // Don't save storage when upgrading code
        },

        ResetRoot => {
            // Reset all dictionaries to null
            // data.proxyHashes = null;
            // data.workerHashes = null;
            // data.modelHashes = null;
            // data.registeredProxies = null;
        },

        UpgradeFull => {
            // Full upgrade: replace both data and code
            contract.setData(msg.newData);
            contract.setCodePostponed(msg.newCode);
            sendExcesses(storage.ownerAddress, msg.queryId);
            return;  // Don't save storage, we just replaced it
        },
        
        ChangeOwner => {
            storage.ownerAddress = msg.newOwnerAddress;
        }
    }

    sendExcesses(ownerAddress, queryId);

    // Save updated data
    storage.data = data.toCell();
    storage.params = params.toCell();
    storage.version = storage.version + 1;
    storage.save();
}

// === Get Methods ===

get fun last_proxy_seqno(): int {
    var storage = RootStorage.load();
    var data = storage.data.load();
    return data.lastProxySeqno;
}

struct CocoonData {
    version: uint32
    lastProxySeqno: uint32
    paramsVersion: uint32
    uniqueId: uint32
    isTest: bool
    pricePerToken: coins
    workerFeePerToken: coins
    minProxyStake: coins
    minClientStake: coins
    ownerAddress: address
}

get fun get_cocoon_data(): CocoonData {
    var storage = RootStorage.load();
    var data = storage.data.load();
    var params = storage.params.load();

    return CocoonData {
        version: storage.version,
        lastProxySeqno: data.lastProxySeqno,
        paramsVersion: params.paramsVersion,
        uniqueId: params.uniqueId,
        isTest: params.isTest,
        pricePerToken: params.pricePerToken,
        workerFeePerToken: params.workerFeePerToken,
        minProxyStake: params.minProxyStake,
        minClientStake: params.minClientStake,
        ownerAddress: storage.ownerAddress
    };
}

struct CurrentCocoonParams {
    paramsVersion: uint32
    uniqueId: uint32
    isTest: bool
    pricePerToken: coins
    workerFeePerToken: coins
    cachedTokensPriceMultiplier: uint32
    reasoningTokensPriceMultiplier: uint32
    proxyDelayBeforeClose: uint32
    clientDelayBeforeClose: uint32
    minProxyStake: coins
    minClientStake: coins
    proxyCodeHash: uint256
    workerCodeHash: uint256
    clientCodeHash: uint256
}

get fun get_cur_params(): CurrentCocoonParams {
    var storage = RootStorage.load();
    var params = storage.params.load();

    return CurrentCocoonParams {
        paramsVersion: params.paramsVersion,
        uniqueId: params.uniqueId,
        isTest: params.isTest,
        pricePerToken: params.pricePerToken,
        workerFeePerToken: params.workerFeePerToken,
        cachedTokensPriceMultiplier: params.cachedTokensPriceMultiplier,
        reasoningTokensPriceMultiplier: params.reasoningTokensPriceMultiplier,
        proxyDelayBeforeClose: params.proxyDelayBeforeClose,
        clientDelayBeforeClose: params.clientDelayBeforeClose,
        minProxyStake: params.minProxyStake,
        minClientStake: params.minClientStake,
        proxyCodeHash: params.proxyScCode != null ? params.proxyScCode!.hash() : 0,
        workerCodeHash: params.workerScCode != null ? params.workerScCode!.hash() : 0,
        clientCodeHash: params.clientScCode != null ? params.clientScCode!.hash() : 0
    };
}

get fun proxy_hash_is_valid(hash: int): int {
    var storage = RootStorage.load();
    var data = storage.data.load();
    // val (_, found) = data.proxyHashes.uDictGet(256, hash);
    return /*found ? -1 :*/ 0;
}

get fun worker_hash_is_valid(hash: int): int {
    var storage = RootStorage.load();
    var data = storage.data.load();
    // val (_, found) = data.workerHashes.uDictGet(256, hash);
    return /*found ? -1 :*/ 0;
}

get fun model_hash_is_valid(hash: int): int {
    var storage = RootStorage.load();
    var data = storage.data.load();
    // val (_, found) = data.modelHashes.uDictGet(256, hash);
    return /*found ? -1 :*/ 0;
}
