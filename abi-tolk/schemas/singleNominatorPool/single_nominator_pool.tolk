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

// modes
const `MODE::SEND_MODE_REMAINING_AMOUNT` = 64

// errors
const `ERROR::WRONG_NOMINATOR_WC` = 0x2000
const `ERROR::WRONG_QUERY_ID` = 0x2001
const `ERROR::WRONG_SET_CODE` = 0x2002
const `ERROR::WRONG_VALIDATOR_WC` = 0x2003
const `ERROR::INSUFFICIENT_BALANCE` = 0x2004
const `ERROR::INSUFFICIENT_ELECTOR_FEE` = 0x2005

contract single_nominator_pool {
    incomingMessages: Message
}

// =============== storage =============================

// storage#_ owner_address:MsgAddressInt validator_address:MsgAddressInt = Storage

@inline
fun loadData(): (address, address) {
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

type Message =
    Withdraw |
    ChangeValidatorAddress |
    SendRawMsg |
    Upgrade

struct (0x00001000) Withdraw {
    queryId: uint64
    amount: coins
}

struct (0x00001001) ChangeValidatorAddress {
    queryId: uint64
    newValidatorAddress: address
}

struct (0x00007702) SendRawMsg {
    queryId: uint64
    mode: uint8
    msg: cell
}

struct (0x00009903) Upgrade {
    queryId: uint64
    code: cell
}

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
    if ( (op == 0) & (inMsgBody.bitsEqual(beginCell().storeString("w").endCell().beginParse())) & (inMsgBody.remainingRefsCount() == 0) ) {
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
            createMessage({
                bounce: BounceMode.Only256BitsOfBody,
                dest: ownerAddress,
                value: amount
            }).send(`MODE::SEND_MODE_REMAINING_AMOUNT`); // owner pays gas fees
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
            throw `ERROR::WRONG_SET_CODE`;
            // contract.setCodePostponed(code);
        }
    }

    // validator role - the wallet whose private key is on the validator node (can sign blocks but can't steal the funds used for stake)
    // if (sender.bitsEqual(validatorAddress)) {
    // }
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L68
@inline_ref
fun makeAddress(wc: int, addr: int): address {
    return beginCell()
        .storeUint(4, 3).storeInt(wc, 8).storeUint(addr, ADDRESS_SIZE).endCell().beginParse().loadAddress();
}

// taken from nominator-pool: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L78
@inline_ref
fun electorAddress(): address {
    var elector: int = blockchain.configParam(1)!.beginParse().preloadUint(ADDRESS_SIZE);
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

struct Roles {
    owner: address
    validator: address
}

get fun get_roles(): Roles {
    var (ownerAddress, validatorAddress) = loadData();
    return Roles {
        owner: ownerAddress,
        validator: validatorAddress,
    };
}

struct PoolConfig {
    min_stake: int
    deposit_fee: int
    withdraw_fee: int
    pool_fee: int
    receipt_price: int
}

struct PoolData {
    state: int
    nominators_count: int
    stake_amount_sent: int
    validator_amount: int
    pool_config: PoolConfig
    nominators: cell?
    withdraw_requests: cell?
    stake_at: int
    saved_validator_set_hash: int
    validator_set_changes_count: int
    validator_set_change_time: int
    stake_held_for: int
    config_proposal_votings: cell?
}

// nominator-pool interface with mytonctrl: https://github.com/ton-blockchain/nominator-pool/blob/2f35c36b5ad662f10fd7b01ef780c3f1949c399d/func/pool.fc#L198
// since we are relying on the existing interface between mytonctrl and nominator-pool, we return values that instruct mytonctrl
// to recover stake on every cycle, although mytonctrl samples every 10 minutes we assume its current behavior that new_stake
// and recover_stake are only sent once per cycle and don't waste gas
get fun get_pool_data(): PoolData {
    return PoolData {
        state: 2, // funds staked at elector and should be recovered by mytonctrl
        nominators_count: 1, // owner is the single nominator
        stake_amount_sent: 0, // unused, mytonctrl does not rely on this param
        validator_amount: 0, // unused, since gas is always paid by validator there is no concept of validator_amount
        pool_config: (PoolConfig {
            min_stake: 0,
            deposit_fee: 0,
            withdraw_fee: 0,
            pool_fee: 0,
            receipt_price: 0
        }), // unused, since not inviting third party nominators
        nominators: null, // unused, mytonctrl does not rely on this param
        withdraw_requests: null, // unused, not needed since owner controls the validator
        stake_at: 0, // unused, mytonctrl does not rely on this param
        saved_validator_set_hash: 0, // unused, required for maintaining validator_amount that we don't need
        validator_set_changes_count: 2, // funds staked at elector and should be recovered by mytonctrl
        validator_set_change_time: 0, // back in the past so mytonctrl will always attempt to recover stake
        stake_held_for: 0, // back in the past so mytonctrl will always attempt to recover stake
        config_proposal_votings: null // unused, not needed since owner controls the validator
    };
}
