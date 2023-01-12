package code

import (
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

type ContractMethod string

const (
	GetSubwalletId           ContractMethod = "get_subwallet_id"
	GetPluginList            ContractMethod = "get_plugin_list"
	GetSubscriptionData      ContractMethod = "get_subscription_data"
	GetCollectionData        ContractMethod = "get_collection_data"
	GetNftData               ContractMethod = "get_nft_data"
	GetJettonData            ContractMethod = "get_jetton_data"
	GetWalletData            ContractMethod = "get_wallet_data"
	RoyaltyParams            ContractMethod = "royalty_params"
	GetEditor                ContractMethod = "get_editor"
	GetSaleData              ContractMethod = "get_sale_data"
	DnsResolve               ContractMethod = "dnsresolve"
	GetDomain                ContractMethod = "get_domain"
	GetFullDomain            ContractMethod = "get_full_domain"
	GetAuctionInfo           ContractMethod = "get_auction_info"
	GetTelemintTokenName     ContractMethod = "get_telemint_token_name"
	GetTelemintAuctionConfig ContractMethod = "get_telemint_auction_config"
	GetLockupData            ContractMethod = "get_lockup_data"
)

// crc16 xmodem || 0x10000
var Methods = map[int64]ContractMethod{
	81467:  GetSubwalletId,
	107653: GetPluginList,
	92260:  GetSubscriptionData,
	102491: GetCollectionData,
	102351: GetNftData,
	106029: GetJettonData,
	97026:  GetWalletData,
	85719:  RoyaltyParams,
	90228:  GetEditor,
	72748:  GetSaleData,
	123660: DnsResolve,
	119378: GetDomain,
	66763:  GetFullDomain,
	80697:  GetAuctionInfo,
	69506:  GetTelemintTokenName,
	129619: GetTelemintAuctionConfig,
	107305: GetLockupData,
}

// ParseContractMethods tries to extract method names from the given code.
// It does work for most of the code compiled by FunC
// but there are smart contracts that can't be handled by ParseContractMethods.
// So if ParseContractMethods returns an error,
// you need to find another way to get method names.
func ParseContractMethods(code []byte) ([]int64, error) {
	cell, err := boc.DeserializeBoc(code)
	if err != nil {
		return nil, err
	}
	c, err := cell[0].NextRef()
	if err != nil {
		return nil, err
	}
	type GetMethods struct {
		Hashmap tlb.Hashmap[tlb.Uint19, boc.Cell]
	}
	var methods GetMethods
	err = tlb.Unmarshal(c, &methods)
	if err != nil {
		return nil, err
	}

	keys := methods.Hashmap.Keys()
	ifs := make([]int64, len(keys))
	for i := range keys {
		ifs[i] = int64(keys[i])
	}
	return ifs, nil
}
