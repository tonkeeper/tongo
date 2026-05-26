// EVAA Master Contract ABI
// Generated from FunC sources: https://github.com/evaafi/contracts

contract evaa_master {
    incomingMessages: (),
    storage: null,
}

fun onInternalMessage(in: InMessage) {

}


// === Get Methods ===

struct UpgradeConfigResult {
    masterCodeVersion: int
    userCodeVersion: int
    timeout: uint32
    updateTime: uint64
    freezeTime: uint64
    userCode: cell
    newMasterCode: cell?
    newUserCode: cell?
}

struct AssetTrackingInfo {
    trackingSupplyIndex: int
    trackingBorrowIndex: int
    lastAccrual: int
}

struct SbRate {
    sRate: int
    bRate: int
}

struct AssetTotals {
    totalSupply: int
    totalBorrow: int
}

get fun get_asset_sb_rate(assetId: int): SbRate {
    return SbRate{sRate: 0, bRate: 0};
}

get fun getAssetRates(assetId: int): SbRate {
    return SbRate{sRate: 0, bRate: 0};
}

get fun getAssetReserves(assetId: int): int {
    return 0;
}

get fun getAssetTotals(assetId: int): AssetTotals {
    return AssetTotals{totalSupply: 0, totalBorrow: 0};
}

get fun getUpdatedRates(assetConfigCollection: cell, assetDynamicsCollection: cell, assetId: int, timeElapsed: int): SbRate {
    return SbRate{sRate: 0, bRate: 0};
}

get fun getUpdatedRatesForAllAssets(timeElapsed: int): cell {
    return beginCell().endCell();
}

get fun getCollateralQuote(borrowAssetId: int, borrowLiquidateAmount: int, collateralAssetId: int, pricesPacked: cell): int {
    return 0;
}

get fun get_user_address(ownerAddress: address): address {
    return get_user_subaccount_address(ownerAddress, 0);
}

get fun get_user_subaccount_address(ownerAddress: address, subaccountId: int32): address {
    throw 0;
}

get fun get_active(): int {
    return 0;
}

get fun getTokensKeys(): cell {
    return beginCell().endCell();
}

get fun getLastUserScVersion(): int {
    return 0;
}

get fun getUpgradeConfig(): UpgradeConfigResult {
    return UpgradeConfigResult{
        masterCodeVersion: 0,
        userCodeVersion: 0,
        timeout: 0,
        updateTime: 0,
        freezeTime: 0,
        userCode: beginCell().endCell(),
        newMasterCode: null,
        newUserCode: null
    };
}

get fun get_asset_tracking_info(assetId: int): AssetTrackingInfo {
    return AssetTrackingInfo{trackingSupplyIndex: 0, trackingBorrowIndex: 0, lastAccrual: 0};
}

get fun get_supervisor(): address {
    return address("EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr");
}

get fun get_asset_total_principals(assetId: int): AssetTotals {
    return AssetTotals{totalSupply: 0, totalBorrow: 0};
}

get fun get_asset_balance(assetId: int): int {
    return 0;
}

get fun get_asset_liquidity_by_id(assetId: int): int {
    return 0;
}

get fun get_asset_liquidity_minus_reserves_by_id(assetId: int): int {
    return 0;
}

get fun getStore(): cell {
    return beginCell().endCell();
}

get fun claim_asset_reserves_min_attachment(fwdFee: int): int {
    return 0;
}

get fun supply_min_attachment(fwdFee: int, supplyUserMessage: cell): int {
    return 0;
}

get fun withdraw_min_attachment(fwdFee: int, withdrawUserMessage: cell): int {
    return 0;
}

get fun liquidate_min_attachment(fwdFee: int, liquidateUserMessage: cell): int {
    return 0;
}
