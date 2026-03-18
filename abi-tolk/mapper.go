package abi_tolk

import (
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/tolk"
	"github.com/tonkeeper/tongo/tolk/parser"

	_ "embed"
)

// ==== TONKEEPER 2FA ====

//go:embed schemas/2fa/tonkeeper_2fa.json
var tonkeeper2Fa []byte

// ====== AFFLUENT =======

//go:embed schemas/affluent/lending_vault.json
var affluentLendingVault []byte

//go:embed schemas/affluent/multiply_vault.json
var affluentMultiplyVault []byte

//go:embed schemas/affluent/multiply_vault_v2.json
var affluentMultiplyV2Vault []byte

//go:embed schemas/affluent/pool.json
var affluentPoolVault []byte

// ======= AIRDROP =======

//go:embed schemas/airdrop/airdrop.json
var airdrop []byte

// ======== BEMO =========

//go:embed schemas/bemo/bemo.json
var bemo []byte

// ======= BIDASK ========

//go:embed schemas/bidask/bidask_damm_lp_wallet.json
var bidaskDammLpWallet []byte

//go:embed schemas/bidask/bidask_damm_pool.json
var bidaskDammPool []byte

//go:embed schemas/bidask/bidask_internal_liquidity_vault.json
var bidaskInternalLiquidityVault []byte

//go:embed schemas/bidask/bidask_lp_multitoken.json
var bidaskLpMultitoken []byte

//go:embed schemas/bidask/bidask_pool.json
var bidaskPool []byte

//go:embed schemas/bidask/bidask_range.json
var bidaskRange []byte

// ======= COCOON =======

//go:embed schemas/cocoon/client.json
var cocoonClient []byte

//go:embed schemas/cocoon/proxy.json
var cocoonProxy []byte

//go:embed schemas/cocoon/root.json
var cocoonRoot []byte

//go:embed schemas/cocoon/wallet.json
var cocoonWallet []byte

//go:embed schemas/cocoon/worker.json
var cocoonWorker []byte

// ======== CRON ========

//go:embed schemas/cron/cron.json
var cron []byte

// ======= DAOLAMA =======

//go:embed schemas/daolama/vault.json
var daolamaVault []byte

// ======= DEDUST =======

//go:embed schemas/dedust/factory.json
var dedustFactory []byte

//go:embed schemas/dedust/liquidity_deposit.json
var dedustLiquidityDeposit []byte

//go:embed schemas/dedust/pool.json
var dedustPool []byte

//go:embed schemas/dedust/vault.json
var dedustVault []byte

// ======== DNS ==========

//go:embed schemas/dns/dns.json
var dns []byte

// ======== GRAM =========

//go:embed schemas/gram/miner.json
var gramMiner []byte

// ======== HIPO =========

//go:embed schemas/hipo/hipo.json
var hipo []byte

// ======== INVOICE ========

//go:embed schemas/invoice/invoice.json
var invoice []byte

// ======== JETTON ========

//go:embed schemas/jetton/jetton_minter.json
var jettonMinter []byte

//go:embed schemas/jetton/jetton_wallet.json
var jettonWallet []byte

//go:embed schemas/jetton/jetton_wallet_governed.json
var jettonWalletGoverned []byte

//go:embed schemas/jetton/scaled_ui.json
var scaledUi []byte

// ======== KNOWN ========

//go:embed schemas/known/known.json
var known []byte

// === LIQUID STAKING ===

//go:embed schemas/liquid_staking/controller.json
var liquidStakingController []byte

//go:embed schemas/liquid_staking/pool.json
var liquidStakingPool []byte

// ===== MEGATONFI =====

//go:embed schemas/megatonfi/exchange.json
var megatonfiExchange []byte

//go:embed schemas/megatonfi/router.json
var megatonfiRouter []byte

// ====== MOONCX ======

//go:embed schemas/mooncx/booster.json
var mooncxBooster []byte

//go:embed schemas/mooncx/order.json
var mooncxOrder []byte

//go:embed schemas/mooncx/order_factory.json
var mooncxOrderFactory []byte

//go:embed schemas/mooncx/pool.json
var mooncxPool []byte

// ===== MULTISIG ======

//go:embed schemas/multisig/multisig_order_v2.json
var multisigOrderV2 []byte

//go:embed schemas/multisig/multisig_v2.json
var multisigV2 []byte

// ======= NFTS =======

//go:embed schemas/nfts/collection.json
var nftsCollection []byte

//go:embed schemas/nfts/editable.json
var nftsEditable []byte

//go:embed schemas/nfts/item.json
var nftsItem []byte

//go:embed schemas/nfts/item_simple.json
var nftsItemSimple []byte

//go:embed schemas/nfts/sbt.json
var nftsSbt []byte

// == PAYMENT CHANNEL ==

//go:embed schemas/payment_channel/payment_channel.json
var paymentChannel []byte

// ======= STONFI =======

//go:embed schemas/stonfi/lp_account_v2.json
var stonfiLpAccountV2 []byte

//go:embed schemas/stonfi/pool_v1.json
var stonfiPoolV1 []byte

//go:embed schemas/stonfi/pool_v2.json
var stonfiPoolV2 []byte

//go:embed schemas/stonfi/router_v1.json
var stonfiRouterV1 []byte

//go:embed schemas/stonfi/router_v2.json
var stonfiRouterV2 []byte

//go:embed schemas/stonfi/vault_v2.json
var stonfiVaultV2 []byte

// ====== STORAGE =======

//go:embed schemas/storages/contract.json
var storageContract []byte

//go:embed schemas/storages/provider.json
var storageProvider []byte

// ======= STORM =======

//go:embed schemas/storm_trade/executor.json
var stormExecutor []byte

//go:embed schemas/storm_trade/position_manager.json
var stormPositionManager []byte

//go:embed schemas/storm_trade/referral.json
var stormReferral []byte

//go:embed schemas/storm_trade/smart_account.json
var stormSmartAccount []byte

//go:embed schemas/storm_trade/smart_account_blank.json
var stormSmartAccountBlank []byte

//go:embed schemas/storm_trade/smart_account_factory.json
var stormSmartAccountFactory []byte

//go:embed schemas/storm_trade/vamm.json
var stormVamm []byte

//go:embed schemas/storm_trade/vamm_coinm.json
var stormVammCoinm []byte

//go:embed schemas/storm_trade/vault.json
var stormVault []byte

//go:embed schemas/storm_trade/vault_native.json
var stormVaultNative []byte

// ==== SUBSCRIPTIONS ====

//go:embed schemas/subscriptions/v1.json
var subscriptionV1 []byte

//go:embed schemas/subscriptions/v2.json
var subscriptionV2 []byte

// ===== SWAP COFFEE =====

//go:embed schemas/swap_coffee/cross_dex.json
var swapCofffeeCrossDex []byte

//go:embed schemas/swap_coffee/factory.json
var swapCofffeeFactory []byte

//go:embed schemas/swap_coffee/mev.json
var swapCofffeeMev []byte

//go:embed schemas/swap_coffee/pool.json
var swapCofffeePool []byte

//go:embed schemas/swap_coffee/staking_item.json
var swapCofffeeStakingItem []byte

//go:embed schemas/swap_coffee/staking_master.json
var swapCofffeeStakingMaster []byte

//go:embed schemas/swap_coffee/staking_vault.json
var swapCofffeeStakingVault []byte

//go:embed schemas/swap_coffee/vault_extra.json
var swapCofffeeVaultExtra []byte

//go:embed schemas/swap_coffee/vault_jetton.json
var swapCofffeeVaultJetton []byte

//go:embed schemas/swap_coffee/vault_native.json
var swapCofffeeVaultNative []byte

// ====== TELEGRAM =======

//go:embed schemas/telegram/telegram.json
var telegram []byte

// ======= TONCO ========

//go:embed schemas/tonco/account.json
var toncoAccount []byte

//go:embed schemas/tonco/pool.json
var toncoPool []byte

//go:embed schemas/tonco/router.json
var toncoRouter []byte

// === TON VALIDATORS ===

//go:embed schemas/tonvalidators/tonvalidators.json
var tonvalidators []byte

// =====-= WALLET =======

//go:embed schemas/wallet/highload_v2.json
var walletHighloadV2 []byte

//go:embed schemas/wallet/highload_v3.json
var walletHighloadV3 []byte

//go:embed schemas/wallet/wallet_preprocessed_v2.json
var walletPreprocessedV2 []byte

//go:embed schemas/wallet/wallet_v3.json
var walletV3 []byte

//go:embed schemas/wallet/wallet_v4.json
var walletV4 []byte

//go:embed schemas/wallet/wallet_v5.json
var walletV5 []byte

// ======= WHALES =======

//go:embed schemas/whales/whales.json
var whales []byte

// =======================

var interfaceToABI = map[abi.ContractInterface]parser.ABI{
	abi.Tonkeeper2Fa: parser.MustParseABI(tonkeeper2Fa),

	abi.AffluentLendingVault:    parser.MustParseABI(affluentLendingVault),
	abi.AffluentMultiplyVault:   parser.MustParseABI(affluentMultiplyVault),
	abi.AffluentMultiplyVaultV2: parser.MustParseABI(affluentMultiplyV2Vault),
	abi.AffluentPool:            parser.MustParseABI(affluentPoolVault),

	abi.AirdropInterlockerV1: parser.MustParseABI(airdrop),
	abi.AirdropInterlockerV2: parser.MustParseABI(airdrop),

	abi.BidaskDammLpWallet:           parser.MustParseABI(bidaskDammLpWallet),
	abi.BidaskDammPool:               parser.MustParseABI(bidaskDammPool),
	abi.BidaskInternalLiquidityVault: parser.MustParseABI(bidaskInternalLiquidityVault),
	abi.BidaskLpMultitoken:           parser.MustParseABI(bidaskLpMultitoken),
	abi.BidaskPool:                   parser.MustParseABI(bidaskPool),
	abi.BidaskRange:                  parser.MustParseABI(bidaskRange),

	abi.CocoonClient: parser.MustParseABI(cocoonClient),
	abi.CocoonProxy:  parser.MustParseABI(cocoonProxy),
	abi.CocoonRoot:   parser.MustParseABI(cocoonRoot),
	abi.CocoonWallet: parser.MustParseABI(cocoonWallet),
	abi.CocoonWorker: parser.MustParseABI(cocoonWorker),

	abi.Cron: parser.MustParseABI(cron),

	abi.DaolamaVault: parser.MustParseABI(daolamaVault),

	abi.DedustFactory:          parser.MustParseABI(dedustFactory),
	abi.DedustLiquidityDeposit: parser.MustParseABI(dedustLiquidityDeposit),
	abi.DedustPool:             parser.MustParseABI(dedustPool),
	abi.DedustVault:            parser.MustParseABI(dedustVault),

	abi.Dns: parser.MustParseABI(dns),

	abi.GramMiner: parser.MustParseABI(gramMiner),

	abi.JettonMaster:          parser.MustParseABI(jettonMinter),
	abi.JettonWallet:          parser.MustParseABI(jettonWallet),
	abi.JettonWalletV1:        parser.MustParseABI(jettonWallet),
	abi.JettonWalletV2:        parser.MustParseABI(jettonWallet),
	abi.JettonWalletGoverned:  parser.MustParseABI(jettonWalletGoverned),
	abi.JettonWalletRegulated: parser.MustParseABI(jettonWalletGoverned),

	abi.TonstakePool:        parser.MustParseABI(liquidStakingPool),
	abi.ValidatorController: parser.MustParseABI(liquidStakingController),

	abi.MegatonfiExchange: parser.MustParseABI(megatonfiExchange),
	abi.MegatonfiRouter:   parser.MustParseABI(megatonfiRouter),

	abi.MoonBooster:      parser.MustParseABI(mooncxBooster),
	abi.MoonOrder:        parser.MustParseABI(mooncxOrder),
	abi.MoonOrderFactory: parser.MustParseABI(mooncxOrderFactory),
	abi.MoonPool:         parser.MustParseABI(mooncxPool),

	abi.MultisigOrderV2: parser.MustParseABI(multisigOrderV2),
	abi.MultisigV2:      parser.MustParseABI(multisigV2),

	abi.NftCollection: parser.MustParseABI(nftsCollection),
	abi.Editable:      parser.MustParseABI(nftsEditable),
	abi.NftItem:       parser.MustParseABI(nftsItem),
	abi.NftItemSimple: parser.MustParseABI(nftsItemSimple),
	abi.Sbt:           parser.MustParseABI(nftsSbt),

	abi.PaymentChannel: parser.MustParseABI(paymentChannel),

	abi.StonfiLpAccountV2: parser.MustParseABI(stonfiLpAccountV2),
	abi.StonfiPool:        parser.MustParseABI(stonfiPoolV1),
	abi.StonfiPoolV2:      parser.MustParseABI(stonfiPoolV2),
	abi.StonfiRouter:      parser.MustParseABI(stonfiRouterV1),
	abi.StonfiRouterV2:    parser.MustParseABI(stonfiRouterV2),
	abi.StonfiVaultV2:     parser.MustParseABI(stonfiVaultV2),

	abi.StorageContract: parser.MustParseABI(storageContract),
	abi.StorageProvider: parser.MustParseABI(storageProvider),

	abi.StormExecutor:        parser.MustParseABI(stormExecutor),
	abi.StormPositionManager: parser.MustParseABI(stormPositionManager),
	abi.StormReferral:        parser.MustParseABI(stormReferral),
	abi.SmartAccount:         parser.MustParseABI(stormSmartAccount),
	abi.SmartAccountBlank:    parser.MustParseABI(stormSmartAccountBlank),
	abi.SmartAccountFactory:  parser.MustParseABI(stormSmartAccountFactory),
	abi.StormVamm:            parser.MustParseABI(stormVamm),
	abi.StormVammCoinm:       parser.MustParseABI(stormVammCoinm),
	abi.StormVault:           parser.MustParseABI(stormVault),
	abi.StormVaultNative:     parser.MustParseABI(stormVaultNative),

	abi.SubscriptionV1: parser.MustParseABI(subscriptionV1),
	abi.SubscriptionV2: parser.MustParseABI(subscriptionV2),

	abi.CoffeeCrossDex:      parser.MustParseABI(swapCofffeeCrossDex),
	abi.CoffeeFactory:       parser.MustParseABI(swapCofffeeFactory),
	abi.CoffeeMevProtector:  parser.MustParseABI(swapCofffeeMev),
	abi.CoffeePool:          parser.MustParseABI(swapCofffeePool),
	abi.CoffeeStakingItem:   parser.MustParseABI(swapCofffeeStakingItem),
	abi.CoffeeStakingMaster: parser.MustParseABI(swapCofffeeStakingMaster),
	abi.CoffeeVaultExtra:    parser.MustParseABI(swapCofffeeVaultExtra),
	abi.CoffeeVaultJetton:   parser.MustParseABI(swapCofffeeVaultJetton),
	abi.CoffeeVaultNative:   parser.MustParseABI(swapCofffeeVaultNative),

	abi.Teleitem: parser.MustParseABI(telegram),

	abi.ToncoAccount: parser.MustParseABI(toncoAccount),
	abi.ToncoPool:    parser.MustParseABI(toncoPool),
	abi.ToncoRouter:  parser.MustParseABI(toncoRouter),

	abi.TvPool: parser.MustParseABI(tonvalidators),

	abi.WalletHighloadV1R1:   parser.MustParseABI(walletHighloadV2),
	abi.WalletHighloadV1R2:   parser.MustParseABI(walletHighloadV2),
	abi.WalletHighloadV2R1:   parser.MustParseABI(walletHighloadV2),
	abi.WalletHighloadV2R2:   parser.MustParseABI(walletHighloadV2),
	abi.WalletHighloadV3R1:   parser.MustParseABI(walletHighloadV3),
	abi.WalletPreprocessedV2: parser.MustParseABI(walletPreprocessedV2),
	abi.WalletV1R1:           parser.MustParseABI(walletV3),
	abi.WalletV1R2:           parser.MustParseABI(walletV3),
	abi.WalletV1R3:           parser.MustParseABI(walletV3),
	abi.WalletV2R1:           parser.MustParseABI(walletV3),
	abi.WalletV2R2:           parser.MustParseABI(walletV3),
	abi.WalletV3R1:           parser.MustParseABI(walletV3),
	abi.WalletV3R2:           parser.MustParseABI(walletV3),
	abi.WalletV4R1:           parser.MustParseABI(walletV4),
	abi.WalletV4R2:           parser.MustParseABI(walletV4),
	abi.WalletV5Beta:         parser.MustParseABI(walletV5),
	abi.WalletV5R1:           parser.MustParseABI(walletV5),
	abi.WalletVesting:        parser.MustParseABI(walletV3),

	abi.WhalesPool: parser.MustParseABI(whales),
}

func GetDecoderWithInterfaces(interfaces ...abi.ContractInterface) (*tolk.Decoder, error) {
	abis := make([]parser.ABI, 0, len(interfaces)+4)
	abis = append(abis, parser.MustParseABI(bemo))
	abis = append(abis, parser.MustParseABI(hipo))
	abis = append(abis, parser.MustParseABI(known))
	abis = append(abis, parser.MustParseABI(invoice))
	for _, i := range interfaces {
		a, ok := interfaceToABI[i]
		if !ok {
			continue
		}
		abis = append(abis, a)
	}

	decoder := tolk.NewDecoder()
	err := decoder.WithABIs(abis...)
	if err != nil {
		return nil, err
	}
	return decoder, nil
}
