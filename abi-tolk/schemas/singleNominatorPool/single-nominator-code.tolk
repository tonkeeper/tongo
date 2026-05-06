// build by `acton func2tolk ./single-nominator-code.fc`

// https://github.com/orbs-network/single-nominator/blob/main/contracts/single-nominator.fc
// this contract is very similar to https://github.com/ton-blockchain/nominator-pool but much simpler since it only supports a single nominator
// frankly speaking, we tried using nominator-pool but it's so complicated that we couldn't be sure there were no bugs hiding around
// this contract is very simple and easy to review, it is laser focused on protecting your stake and nothing else!

// =============== consts =============================

const ADDRESS_SIZE = 256
const MIN_TON_FOR_STORAGE = 1000000000 // 10% from nominator-pool since we have a single cell
const MIN_TON_FOR_SEND_MSG = 1200000000

// owner ops
const `OP::WITHDRAW` = 0x1000
const `OP::CHANGE_VALIDATOR_ADDRESS` = 0x1001
const `OP::SEND_RAW_MSG` = 0x7702
const `OP::UPGRADE` = 0x9903

// elector ops
const `OP::NEW_STAKE` = 0x4e73744b
const `OP::RECOVER_STAKE` = 0x47657424

// modes
const `MODE::SEND_MODE_REMAINING_AMOUNT` = 64

// errors
const `ERROR::WRONG_NOMINATOR_WC` = 0x2000
const `ERROR::WRONG_QUERY_ID` = 0x2001
const `ERROR::WRONG_SET_CODE` = 0x2002
const `ERROR::WRONG_VALIDATOR_WC` = 0x2003
const `ERROR::INSUFFICIENT_BALANCE` = 0x2004
const `ERROR::INSUFFICIENT_ELECTOR_FEE` = 0x2005

// =============== storage =============================

// storage#_ owner_address:MsgAddressInt validator_address:MsgAddressInt = Storage

@inline
fun loadData(): (slice, slice) {
    var ds: slice = contract.getData().beginParse();

    var ownerAddress: address = ds.loadAddress();
    var validatorAddress: address = ds.loadAddress();

    ds.assertEnd();

    return (ownerAddress, validatorAddress);
}

@inline
fun saveData(ownerAddress: address, validatorAddress: address) {
    contract.setData(beginCell()
        .storeAddress(ownerAddress)
        .storeAddress(validatorAddress)
        .endCell());
}

// =============== messages =============================

// defined below

// main entry point for receiving messages
// my_balance contains balance after adding msg_value
// in the future, use: fun onInternalMessage(in: InMessage) {
fun onInternalMessage(myBalance: int, msgValue: int, inMsgFull: cell, inMsgBody: slice) {
    var (ownerAddress, validatorAddress) = loadData();

    if (inMsgBody.isEmpty()) {
        return;
    }

    var cs: slice = inMsgFull.beginParse();
    var flags: int = cs.loadUint(4); // int_msg_info$0 ihr_disabled:Bool bounce:Bool bounced:Bool
    if (flags & 1) {
        // ignore all bounced messages
        return;
    }
    var sender: address = cs.loadAddress();

    var op: int = inMsgBody.loadUint(32);
    var queryId: int = 0; // filled below

    // if just a message with comment "w" - means withdraw
    if ( (op == 0) & (inMsgBody.bitsEqual("w")) & (inMsgBody.remainingRefsCount() == 0) ) {
        op = `OP::WITHDRAW`;
        // in few lines, amount to withdraw is loaded from in_msg_body,
        // so we set it to the maximum avaliable balance
        inMsgBody = beginCell().storeCoins(myBalance - msgValue - MIN_TON_FOR_STORAGE).endCell().beginParse();
    } else {
        queryId = inMsgBody.loadUint(64);
    }

    // owner role - cold wallet (private key that is not connected to the internet) that owns the funds used for staking and acts as the single nominator
    if (sender.bitsEqual(ownerAddress)) {

        // allow owner to withdraw funds - take the money home and stop validating with it
        if (op == `OP::WITHDRAW`) {
            var amount: int = inMsgBody.loadCoins();
            amount = min(amount, myBalance - msgValue - MIN_TON_FOR_STORAGE);
            assert(amount > 0) throw `ERROR::INSUFFICIENT_BALANCE`;
            sendMsg(ownerAddress, amount, null, BOUNCEABLE, `MODE::SEND_MODE_REMAINING_AMOUNT`); // owner pays gas fees
        }

        // mainly used when the validator was compromised to prevent validator from entering new election cycles
        if (op == `OP::CHANGE_VALIDATOR_ADDRESS`) {
            var newValidatorAddress: address = inMsgBody.loadAddress();
            saveData(ownerAddress, newValidatorAddress);
        }

        // emergency safeguard to allow owner to send arbitrary messages as the nominator contract
        if (op == `OP::SEND_RAW_MSG`) {
            var mode: int = inMsgBody.loadUint(8);
            var msg: cell = inMsgBody.loadRef();
            sendRawMessage(msg, mode);
        }

        // second emergency safeguard to allow owner to replace nominator logic - you should never need to use this
        if (op == `OP::UPGRADE`) {
            var code: cell = inMsgBody.loadRef();
            assert(!((code == null))) throw `ERROR::WRONG_SET_CODE`;
            contract.setCodePostponed(code);
        }
    }

    // validator role - the wallet whose private key is on the validator node (can sign blocks but can't steal the funds used for stake)
    if (sender.bitsEqual(validatorAddress)) {

        // send stake to the elector for the next validation cycle (sent every period ~18 hours)
        if (op == `OP::NEW_STAKE`) {
            var (senderWc: int, _) = sender.getWorkchainAndHash();
            var (myWc: int, _) = contract.getAddress().getWorkchainAndHash();
            assert(senderWc == -1) throw `ERROR::WRONG_VALIDATOR_WC`; // for voting purpose
            assert(myWc == -1) throw `ERROR::WRONG_NOMINATOR_WC`; // must be deployed on masterchain
            assert(queryId) throw `ERROR::WRONG_QUERY_ID`; // query_id must be greater then 0 to receive confirmation message from elector
            assert(msgValue >= MIN_TON_FOR_SEND_MSG) throw `ERROR::INSUFFICIENT_ELECTOR_FEE`; // must be greater then new_stake sending to elector fee
            var stakeAmount: int = inMsgBody.loadCoins();
            var msg: slice = inMsgBody;
            checkNewStakeMsg(inMsgBody);
            assert(stakeAmount <= myBalance - msgValue - MIN_TON_FOR_STORAGE) throw `ERROR::INSUFFICIENT_BALANCE`;

            sendMsg(electorAddress(), stakeAmount, beginCell().storeUint(`OP::NEW_STAKE`, 32).storeUint(queryId, 64).storeSlice(msg).endCell(), BOUNCEABLE, `MODE::SEND_MODE_REMAINING_AMOUNT`); // bounceable, validator pays gas fees
        }

        // recover stake from elector of previous validation cycle (sent every period ~18 hours)
        if (op == `OP::RECOVER_STAKE`) {
            var payload: cell = beginCell().storeUint(`OP::RECOVER_STAKE`, 32).storeUint(queryId, 64).endCell();
            sendMsg(electorAddress(), 0, payload, BOUNCEABLE, `MODE::SEND_MODE_REMAINING_AMOUNT`); // bounceable, validator pays gas fees
        }
    }
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L217
@inline_ref
fun sendMsg(toAddress: address, amount: int, payload: cell, flags: int, sendMode: int) {
    var hasPayload: int = ~ (payload == null);

    var msg: builder = beginCell()
        .storeUint(flags, 6)
        .storeAddress(toAddress)
        .storeCoins(amount)
        .storeUint(hasPayload ? 1 : 0, 1 + 4 + 4 + 64 + 32 + 1 + 1);

    if (hasPayload) {
        msg = msg.storeRef(payload);
    }

    sendRawMessage(msg.endCell(), sendMode);
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L68
@inline_ref
fun makeAddress(wc: int, addr: int): address {
    return beginCell()
        .storeUint(4, 3).storeInt(wc, 8).storeUint(addr, ADDRESS_SIZE).endCell().beginParse();
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L78
@inline_ref
fun electorAddress(): address {
    var elector: int = blockchain.configParam(1).beginParse().preloadUint(ADDRESS_SIZE);
    return makeAddress(-1, elector);
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L139
// check the validity of the new_stake message
// https://github.com/ton-blockchain/ton/blob/b38d227a469666d83ac535ad2eea80cb49d911b8/crypto/smartcont/elector-code.fc#L208
@inline_ref
fun checkNewStakeMsg(cs: slice): int {
    var validatorPubkey = cs.loadUint(256);
    var stakeAt = cs.loadUint(32);
    var maxFactor = cs.loadUint(32);
    var adnlAddr = cs.loadUint(256);
    var signature = cs.loadRef().beginParse().preloadBits(512);
    cs.assertEnd();
    return stakeAt; // supposed start of next validation round (utime_since)
}

// =============== getters =============================

get fun get_roles(): (slice, slice) {
    var (ownerAddress, validatorAddress) = loadData();
    return (ownerAddress, validatorAddress);
}

// nominator-pool interface with mytonctrl: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L198
// since we are relying on the existing interface between mytonctrl and nominator-pool, we return values that instruct mytonctrl
// to recover stake on every cycle, although mytonctrl samples every 10 minutes we assume its current behavior that new_stake
// and recover_stake are only sent once per cycle and don't waste gas
get fun get_pool_data(): (int, int, int, int, (int, int, int, int, int), cell, cell, int, int, int, int, int, cell) {
    return (
        2, // state - funds staked at elector and should be recovered by mytonctrl
        1, // nominators_count - owner is the single nominator
        0, // stake_amount_sent - unused, mytonctrl does not rely on this param
        0, // validator_amount - unused, since gas is always paid by validator there is no concept of validator_amount
        (0, 0, 0, 0, 0), // pool config - unused, since not inviting third party nominators
        null, // nominators - unused, mytonctrl does not rely on this param
        null, // withdraw_requests - unused, not needed since owner controls the validator
        0, // stake_at - unused, mytonctrl does not rely on this param
        0, // saved_validator_set_hash - unused, required for maintaining validator_amount that we don't need
        2, // validator_set_changes_count - funds staked at elector and should be recovered by mytonctrl
        0, // validator_set_change_time - back in the past so mytonctrl will always attempt to recover stake
        0, // stake_held_for - back in the past so mytonctrl will always attempt to recover stake
        null // config_proposal_votings - unused, not needed since owner controls the validator
    );
}
