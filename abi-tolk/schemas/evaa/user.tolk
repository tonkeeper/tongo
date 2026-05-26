// EVAA User (Position) Contract ABI
// Generated from FunC sources: https://github.com/evaafi/contracts

contract evaa_user {
    incomingMessages: (),
    storage: null,
}

fun onInternalMessage(in: InMessage) {

}


// === Get Methods ===

struct UserScData {
    codeVersion: int
    masterAddress: address
    ownerAddress: address
    userPrincipals: cell
    state: int64
    userRewards: cell?
    backupCell1: cell?
    backupCell2: cell?
}

struct AggregatedBalances {
    supply: int
    borrow: int
}

get fun codeVersion(): int {
    return 0;
}

get fun isUserSc(): int {
    return -1;
}

get fun getAccountAssetBalance(assetId: int, sRate: int, bRate: int): int {
    return 0;
}

get fun getAccountBalances(assetDynamicsCollection: cell): cell {
    return beginCell().endCell();
}

get fun getAccountHealth(assetConfigCollection: cell, assetDynamicsCollection: cell, pricesPacked: cell): int {
    return 0;
}

get fun getAvailableToBorrow(assetConfigCollection: cell, assetDynamicsCollection: cell, pricesPacked: cell): int {
    return 0;
}

get fun getIsLiquidable(assetConfigCollection: cell, assetDynamicsCollection: cell, pricesPacked: cell): int {
    return 0;
}

get fun getAggregatedBalances(assetConfigCollection: cell, assetDynamicsCollection: cell, pricesPacked: cell): AggregatedBalances {
    return AggregatedBalances{supply: 0, borrow: 0};
}

get fun get_asset_principal(assetId: int): int {
    return 0;
}

get fun getPrincipals(): map<bits256, int64> {
    throw 0;
}

get fun getRewards(): cell? {
    return beginCell().endCell();
}

get fun getAllUserScData(): UserScData {
    return UserScData{
        codeVersion: 0,
        masterAddress: address("EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr"),
        ownerAddress: address("EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr"),
        userPrincipals: beginCell().endCell(),
        state: 0,
        userRewards: beginCell().endCell(),
        backupCell1: null,
        backupCell2: null
    };
}

get fun get_maximum_withdraw_amount(assetId: int, pricesPacked: cell, assetConfigCollection: cell, assetDynamicsCollection: cell): int {
    return 0;
}
