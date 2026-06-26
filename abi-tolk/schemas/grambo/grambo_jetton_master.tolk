// Grambo Jetton Master (Bonding Curve) Contract ABI schema

import "@stdlib/common"

struct GramboTradeFees {
    referrer_fee: coins
    fee_1: coins
    fee_2: coins
}

struct GramboSellParams {
    min_ton_out: coins
    referrer: address
    custom_payload: cell?
}

struct (0x372f1ae3) GramboLaunch {
    query_id: uint64
    buyer: address
    virtual_ton_reserve: coins
    ton_target: coins
    prebuy_amount: coins
}

struct (0x2f54494e) GramboBuy {
    query_id: uint64
    ton_amount: coins
    recipient: address
    min_tokens_out: coins
    referrer: address
    custom_payload: cell?
}

struct (0x31fc1276) GramboSellNotification {
    query_id: uint64
    amount: coins
    from: address
    response_destination: address
    sell_params: Cell<GramboSellParams>?
}

struct (0x17b5012c) GramboProvideWalletAddress {
    query_id: uint64
    owner: address
    include_address: bool
}

struct (0x7966410f) GramboUnlockWallet {
    query_id: uint64
    owner: address
}

struct (0x43e84bb5) GramboWithdraw {
    query_id: uint64
    amount: coins
}

type MasterIncomingMessage =
    | GramboLaunch
    | GramboBuy
    | GramboSellNotification
    | GramboProvideWalletAddress
    | GramboUnlockWallet
    | GramboWithdraw

struct (0x178d4519) JettonInternalTransfer {
    query_id: uint64
    amount: coins
    from: address
    response_address: address
    forward_ton_amount: coins
    forward_payload: cell
}

struct (0x6b75d3e2) GramboActivateWallet {
    query_id: uint64
}

struct (0x29cb7264) GramboTakeWalletAddress {
    query_id: uint64
    owner_address: cell?
}

struct (0xd53276db) JettonExcess {
    query_id: uint64
}

type MasterOutgoingMessage =
    | JettonInternalTransfer
    | GramboActivateWallet
    | GramboTakeWalletAddress
    | JettonExcess

struct (0x5025398e) GramboSwapEvent {
    query_id: uint64
    recipient: address
    is_buy: bool
    new_collected_ton: coins
    ton_amount: coins
    token_amount: coins
    referrer: address
    fees: Cell<GramboTradeFees>
    custom_payload: cell?
}

struct (0x0ab3d75c) GramboGraduateEvent {
    query_id: uint64
    ton_to_dex: coins
    tokens_to_dex: coins
}

type MasterExternalMessage =
    | GramboSwapEvent
    | GramboGraduateEvent

struct GramboCurveState {
    collected_ton: coins
    tokens_sold: coins
    graduated: bool
    fee_addr_1: address
    fee_addr_2: address
}

struct GramboDexConfig {
    stonfi_router: address
    pton_wallet: address
}

struct GramboMasterStorage {
    virtual_ton_reserve: coins
    ton_target: coins
    factory: address
    total_supply: coins
    seed: uint256
    curve_state: Cell<GramboCurveState>
    wallet_code: cell
    metadata: cell
    dex_config: Cell<GramboDexConfig>
}

contract grambo_jetton_master {
    incomingMessages: MasterIncomingMessage
    outgoingMessages: MasterOutgoingMessage
    incomingExternal: MasterExternalMessage
    storage: GramboMasterStorage
    forceAbiExport: [MasterIncomingMessage, MasterOutgoingMessage, MasterExternalMessage]
}

fun onInternalMessage(in: InMessage) {
    throw 74;    // Unauthorized
    throw 65296; // Already launched
    throw 65344; // Slippage error
    throw 65347; // Unauthorized
    throw 65360; // Insufficient value or balance
    throw 65397; // Not graduated
    throw 65398; // Already graduated
}

get fun get_amount_out(ton_in: uint64): int {
    return 0;
}

get fun get_ton_out(tokens_in: int): uint64 {
    return 0;
}

struct GetCurveStateResult {
    virtual_ton_reserve: uint64
    ton_target: uint64
    collected_ton: uint64
    tokens_sold: int
    graduated: bool
}

get fun get_curve_state(): GetCurveStateResult {
    return GetCurveStateResult {
        virtual_ton_reserve: 0,
        ton_target: 0,
        collected_ton: 0,
        tokens_sold: 0,
        graduated: false,
    };
}

struct GetJettonDataResult {
    total_supply: int
    mintable: int
    admin_address: address
    jetton_content: cell
    jetton_wallet_code: cell
}

get fun get_jetton_data(): GetJettonDataResult {
    return GetJettonDataResult {
        total_supply: 0,
        mintable: 0,
        admin_address: address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
        jetton_content: beginCell().endCell(),
        jetton_wallet_code: beginCell().endCell(),
    };
}

get fun get_wallet_address(owner_address: address): address {
    return address("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs");
}
