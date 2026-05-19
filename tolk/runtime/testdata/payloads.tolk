
type JettonPayload =
        CoffeeStakingLockJettonPayload |
        TextCommentJettonPayload |
        TegroJettonSwapJettonPayload |
        EncryptedTextCommentJettonPayload |
        StonfiSwapJettonPayload |
        TegroAddLiquidityJettonPayload |
        StonfiProvideLpV2JettonPayload |
        BidaskProvideBothJettonPayload |
        DedustDepositLiquidityJettonPayload |
        PoolFundAccountJettonPayload |
        StonfiSwapOkRefJettonPayload |
        CoffeeCrossDexResendJettonPayload |
        BidaskDammProvideJettonPayload |
        StonfiSwapV2JettonPayload |
        BidaskDammProvideOneSideJettonPayload |
        StormDepositJettonJettonPayload |
        InvoicePayloadJettonPayload |
        TonkeeperRelayerFeeJettonPayload |
        BidaskSwapV2JettonPayload |
        MoonBoostPoolJettonPayload |
        BidaskProvideJettonPayload |
        MoonFillOrderJettonPayload |
        BidaskDammProvideBothJettonPayload |
        MoonDepositLiquidityJettonPayload |
        MoonSwapJettonPayload |
        CoffeeCrossDexFailureJettonPayload |
        CoffeeSwapJettonPayload |
        CoffeeCreatePoolJettonPayload |
        CoffeeDepositLiquidityJettonPayload |
        CoffeeNotificationJettonPayload |
        MoonSwapFailedJettonPayload |
        StonfiSwapOkJettonPayload |
        StormStakeJettonPayload |
        WithdrawPayloadJettonPayload |
        MoonSwapSucceedJettonPayload |
        MoonCreateOrderJettonPayload |
        BidaskDammSwapJettonPayload |
        DedustSwapJettonPayload |
        CoffeeMevProtectFailedSwapJettonPayload |
        BidaskSwapJettonPayload |
        DepositPayloadJettonPayload |
        StonfiProvideLiquidityJettonPayload;

struct (0x0c0ffede) CoffeeStakingLockJettonPayload {
    periodID: uint32
}

struct (0x00000000) TextCommentJettonPayload {
    text: slice
}

struct (0x01fb7a25) TegroJettonSwapJettonPayload {
    extract: bool
    maxIn: coins
    minOut: coins
    destination: address
    errorDestination: address
    payload: cell
}

struct (0x2167da4b) EncryptedTextCommentJettonPayload {
    cipherText: slice
}

struct (0x25938561) StonfiSwapJettonPayload {
    tokenWallet: address
    minOut: coins
    toAddress: address
    referralAddress: address?
}

struct (0x287e167a) TegroAddLiquidityJettonPayload {
    amountA: coins
    amountB: coins
}

struct (0x37c096df) StonfiProvideLpV2JettonPayload {
    tokenWallet1: address
    refundAddress: address
    excessesAddress: address
    txDeadline: uint64
    crossProvideLpBody: Cell<CrossProvideLpBody>
}

struct CrossProvideLpBody {
    minLpOut: coins
    toAddress: address
    bothPositive: bool
    fwdAmount: coins
    customPayload: cell?
}

struct (0x3ea0bafc) BidaskProvideBothJettonPayload {
    tonAmount: coins
    depositType: uint4
    liquidityDict: map<uint32, int32>
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0b0) DedustVolatile {}

struct (0b1) DedustStable {}

struct (0b0000) DedustAssetNative {}

struct (0b0001) DedustAssetJetton {
    workchain: int8
    hash: bits256
}

struct (0b0010) DedustAssetExtraCurrency {
    currencyId: int32
}

type DedustPoolType = DedustVolatile | DedustStable

type DedustAsset = DedustAssetNative | DedustAssetJetton | DedustAssetExtraCurrency

struct (0x40e108d6) DedustDepositLiquidityJettonPayload {
    poolType: DedustPoolType
    asset0: DedustAsset
    asset1: DedustAsset
    asset0TargetBalance: coins
    asset1TargetBalance: coins
    fulfillPayload: cell?
    rejectPayload: cell?
}

struct (0x4468de77) PoolFundAccountJettonPayload {
    jettonTarget: address
    enough0: coins
    enough1: coins
    liquidity: uint128
    tickLower: int24
    tickUpper: int24
}

struct (0x45078540) StonfiSwapOkRefJettonPayload {}

struct (0x4ee9b106) CoffeeCrossDexResendJettonPayload {
    next: RemainingBitsAndRefs
}

struct (0x63ec24ae) BidaskDammProvideJettonPayload {
    receiver: address
    lockLiquidity: bool
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0x6664de2a) StonfiSwapV2JettonPayload {
    tokenWallet1: address
    refundAddress: address
    excessesAddress: address
    txDeadline: uint64
    crossSwapBody: Cell<CrossSwapBody>
}

struct CrossSwapBody {
    minOut: coins
    receiver: address
    fwdGas: coins
    customPayload: cell?
    refundFwdGas: coins
    refundPayload: cell?
    refFee: uint16
    refAddress: address
}

struct (0x729c04c8) BidaskDammProvideOneSideJettonPayload {
    receiver: address
    lockLiquidity: bool
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0b0) StormNoKeyInit {}

struct (0b1) StormNeedKeyInit {
    userPublicKeys: map<uint256, cell>
}

type InitializationRequest = StormNoKeyInit | StormNeedKeyInit

struct (0x76840119) StormDepositJettonJettonPayload {
    queryId: uint64
    receiverAddress: address
    init: bool
    keyInit: InitializationRequest
}

struct (0x00) InvoiceUrlNone {}

struct (0x01) InvoiceUrlTonsite {
    address: bits256
}

type InvoiceUrl = InvoiceUrlNone | InvoiceUrlTonsite

struct (0x7aa23eb5) InvoicePayloadJettonPayload {
    id: bits128
    url: InvoiceUrl
}

struct (0x878da6e3) TonkeeperRelayerFeeJettonPayload {}

struct AdditionalData {
    fromAddress: address
    refAddress: address
}

struct (0x87d36990) BidaskSwapV2JettonPayload {
    toAddress: address
    slippage: coins | uint256
    exactOut: coins
    additionalData: Cell<AdditionalData>?
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0x96aa1586) MoonBoostPoolJettonPayload {}

struct (0x96feef7b) BidaskProvideJettonPayload {
    depositType: uint4
    liquidityDict: map<uint32, int32>
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0x99b49842) MoonFillOrderJettonPayload {
    recipient: address
    recipientPayload: cell?
    rejectAddress: address
}

struct (0xa8904134) BidaskDammProvideBothJettonPayload {
    nativeAmount: coins
    reciever: address
    lockLiquidity: bool
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0xb31db781) MoonDepositLiquidityJettonPayload {
    minLpOut: coins
}

struct MoonNextPayload {
    recipient: address
    payload: cell?
}

struct MoonSwapParams {
    minOut: coins
    deadline: uint64
    excess: address
    referral: address
    nextFulfill: Cell<MoonNextPayload>?
    nextReject: Cell<MoonNextPayload>?
}

struct (0xb37a900b) MoonSwapJettonPayload {
    moonSwap: MoonSwapParams
}

struct (0xb902e61a) CoffeeCrossDexFailureJettonPayload {
    queryId: uint64
    recipient: address
}

struct CoffeeSwapStepParams {
    poolAddressHash: uint256
    minOutAmount: coins
    next: Cell<CoffeeSwapStepParams>?
}

struct CoffeeNotificationDataSingle {
    reciever: address
    fwdGas: coins
    payload: cell
}

struct CoffeeNotificationData {
    onSuccess: Cell<CoffeeNotificationDataSingle>?
    onFailure: Cell<CoffeeNotificationDataSingle>?
}

struct CoffeeSwapParams {
    deadline: uint32
    recipient: address
    referral: address
    notificationData: Cell<CoffeeNotificationData>?
}

struct (0xc0ffee10) CoffeeSwapJettonPayload {
    step: CoffeeSwapStepParams
    params: Cell<CoffeeSwapParams>
}

struct (0b00) CoffeeAssetNative {}

struct (0b01) CoffeeAssetJetton {
    chain: uint8
    hash: uint256
}

struct (0b10) CoffeAssetExtra {
    id: uint32
}

type CoffeeAsset = CoffeeAssetNative | CoffeeAssetJetton | CoffeAssetExtra

struct (0b000) CoffeeAmmConstProd {}

struct (0b001) CoffeeAmmCurveFiStable {}

type CoffeeAmm = CoffeeAmmConstProd | CoffeeAmmCurveFiStable

struct CoffeePoolParams {
    first: CoffeeAsset
    second: CoffeeAsset
    amm: CoffeeAmm
    ammSettings: cell?
}

struct CoffeePublicPoolCreationParams {
    recipient: address
    useRecipientOnFailure: bool
    notificationData: cell?
}

struct CoffeePrivatePoolCreationParams {
    isActive: bool
    extraSettings: cell?
}

struct CoffeePoolCreationParams {
    public: CoffeePublicPoolCreationParams
    privateP: CoffeePrivatePoolCreationParams
}

struct (0xc0ffee11) CoffeeCreatePoolJettonPayload {
    params: CoffeePoolParams
    creationParams: CoffeePoolCreationParams
}

struct (0b00) CoffeeDepositLiquidityConditionNone {}

struct (0b01) CoffeeDepositLiquidityConditionLpQuantity {
    minLpAmount: coins
}

struct (0b10) CoffeeDepositLiquidityConditionReservedRatio {
    denominator: coins
    minNominator: coins
    maxNominator: coins
}

struct (0b11) CoffeeDepositLiquidityConditionComplex {
    minLpAmoint: coins
    denominator: coins
    minNominator: coins
    maxNominator: coins
}

type CoffeeDepositLiquidityCondition = CoffeeDepositLiquidityConditionNone | CoffeeDepositLiquidityConditionLpQuantity | CoffeeDepositLiquidityConditionReservedRatio | CoffeeDepositLiquidityConditionComplex

struct CoffeeDepositLiquidityParamsTrimmed {
    recipient: address
    useRecipientOnFailure: bool
    referral: address?
    deadline: uint32
    condition: CoffeeDepositLiquidityCondition
    extraSettings: cell?
    notificationData: Cell<CoffeeNotificationData>?
}

struct CoffeeDepositLiquidityParams {
    params: CoffeeDepositLiquidityParamsTrimmed
    poolParmas: CoffeePoolParams
}

struct (0xc0ffee12) CoffeeDepositLiquidityJettonPayload {
    params: CoffeeDepositLiquidityParams
}

struct (0xc0ffee36) CoffeeNotificationJettonPayload {
    queryId: uint64
    body: cell
}

struct (0xc47c1f57) MoonSwapFailedJettonPayload {}

struct (0xc64370e5) StonfiSwapOkJettonPayload {}

struct (0xc89a3ee4) StormStakeJettonPayload {}

struct (0xcb03bfaf) WithdrawPayloadJettonPayload {
    assetAddress: address
    oracleParams: cell?
    forwardTonAmount: coins
    forwardPayload: cell?
}

struct (0xcb7f38d6) MoonSwapSucceedJettonPayload {}

struct MoonOrderParams {
    rate: uint256
    lock: uint2
    vestingTime: uint64
}

struct (0xda067c19) MoonCreateOrderJettonPayload {
    asset1: address
    asset2: address
    orderData: MoonOrderParams
}

struct (0xdd79732c) BidaskDammSwapJettonPayload {
    toAddress: address
    slippage: coins
    fromAddress: address
    exactOut: coins
    additionalData: cell?
    rejectPayload: cell?
    forwardPayload: cell?
}

struct DedustSwapStepParams {
    kindOut: bool
    limit: coins
    next: Cell<DedustSwapStep>?
}

struct DedustSwapStep {
    poolAddr: address
    params: DedustSwapStepParams
}

struct DedustSwapParams {
    deadline: uint32
    recipientAddr: address
    referralAddr: address
    fulfillPayload: cell?
    rejectPayload: cell?
}

struct (0xe3a0d482) DedustSwapJettonPayload {
    step: DedustSwapStep
    swapParams: Cell<DedustSwapParams>
}

struct (0xee51ce51) CoffeeMevProtectFailedSwapJettonPayload {
    queryId: uint64
    recipient: address
}

struct (0xf2ef6c1b) BidaskSwapJettonPayload {
    toAddress: address
    slippage: coins | uint256
    exactOut: coins
    refAddress: address
    additionalData: cell?
    rejectPayload: cell?
    forwardPayload: cell?
}

struct (0xf9471134) DepositPayloadJettonPayload {
    oracleParams: cell?
    forwardTonAmount: coins
    forwardPayload: cell?
}

struct (0xfcf9e58f) StonfiProvideLiquidityJettonPayload {
    tokenWallet: address
    minLpOut: coins
}
