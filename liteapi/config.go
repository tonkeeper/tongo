package liteapi

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type BlockchainConfig struct {
	ConfigParam0  *tlb.ConfigParam0
	ConfigParam1  *tlb.ConfigParam1
	ConfigParam2  *tlb.ConfigParam2
	ConfigParam3  *tlb.ConfigParam3
	ConfigParam4  *tlb.ConfigParam4
	ConfigParam5  *tlb.ConfigParam5
	ConfigParam6  *tlb.ConfigParam6
	ConfigParam7  *tlb.ConfigParam7
	ConfigParam8  *tlb.ConfigParam8
	ConfigParam9  *tlb.ConfigParam9
	ConfigParam10 *tlb.ConfigParam10
	ConfigParam11 *tlb.ConfigParam11
	ConfigParam12 *tlb.ConfigParam12
	ConfigParam13 *tlb.ConfigParam13
	ConfigParam14 *tlb.ConfigParam14
	ConfigParam15 *tlb.ConfigParam15
	ConfigParam16 *tlb.ConfigParam16
	ConfigParam17 *tlb.ConfigParam17
	ConfigParam18 *tlb.ConfigParam18
	ConfigParam20 *tlb.ConfigParam20
	ConfigParam21 *tlb.ConfigParam21
	ConfigParam22 *tlb.ConfigParam22
	ConfigParam23 *tlb.ConfigParam23
	ConfigParam24 *tlb.ConfigParam24
	ConfigParam25 *tlb.ConfigParam25

	ConfigParam28 *tlb.ConfigParam28
	ConfigParam29 *tlb.ConfigParam29
	ConfigParam31 *tlb.ConfigParam31
	ConfigParam32 *tlb.ConfigParam32
	ConfigParam33 *tlb.ConfigParam33
	ConfigParam34 *tlb.ConfigParam34
	ConfigParam35 *tlb.ConfigParam35
	ConfigParam36 *tlb.ConfigParam36
	ConfigParam37 *tlb.ConfigParam37
	ConfigParam39 *tlb.ConfigParam39
	ConfigParam40 *tlb.ConfigParam40
	ConfigParam43 *tlb.ConfigParam43
	ConfigParam44 *tlb.ConfigParam44

	ConfigParam71 *tlb.ConfigParam71
	ConfigParam72 *tlb.ConfigParam72
	ConfigParam73 *tlb.ConfigParam73

	ConfigParam79 *tlb.ConfigParam79
	ConfigParam81 *tlb.ConfigParam81
	ConfigParam82 *tlb.ConfigParam82
}

func (conf *BlockchainConfig) ConfigAddr() (ton.AccountID, bool) {
	if conf.ConfigParam0 != nil {
		return ton.AccountID{Workchain: -1, Address: conf.ConfigParam0.ConfigAddr}, true
	}
	return ton.AccountID{}, false
}

func (conf *BlockchainConfig) ElectorAddr() (ton.AccountID, bool) {
	if conf.ConfigParam1 != nil {
		return ton.AccountID{Workchain: -1, Address: conf.ConfigParam1.ElectorAddr}, true
	}
	return ton.AccountID{}, false
}

func (conf *BlockchainConfig) MinterAddr() (ton.AccountID, bool) {
	if conf.ConfigParam2 != nil {
		return ton.AccountID{Workchain: -1, Address: conf.ConfigParam2.MinterAddr}, true
	}
	return ton.AccountID{}, false
}

func (conf *BlockchainConfig) FeeCollectorAddr() (ton.AccountID, bool) {
	if conf.ConfigParam3 != nil {
		return ton.AccountID{Workchain: -1, Address: conf.ConfigParam3.FeeCollectorAddr}, true
	}
	return ton.AccountID{}, false
}

func (conf *BlockchainConfig) DnsRootAddr() (ton.AccountID, bool) {
	if conf.ConfigParam4 != nil {
		return ton.AccountID{Workchain: -1, Address: conf.ConfigParam4.DnsRootAddr}, true
	}
	return ton.AccountID{}, false
}

func (conf *BlockchainConfig) MandatoryParams() []int {
	if conf.ConfigParam9 != nil {
		keys := make([]int, 0, len(conf.ConfigParam9.MandatoryParams.Items()))
		for _, item := range conf.ConfigParam9.MandatoryParams.Items() {
			keys = append(keys, int(item.Key))
		}
		return keys
	}
	return nil
}

func (conf *BlockchainConfig) CriticalParams() []int {
	if conf.ConfigParam10 != nil {
		keys := make([]int, 0, len(conf.ConfigParam10.CriticalParams.Items()))
		for _, item := range conf.ConfigParam10.CriticalParams.Items() {
			keys = append(keys, int(item.Key))
		}
		return keys
	}
	return nil
}

func ConvertBlockchainConfig(params tlb.ConfigParams) (*BlockchainConfig, error) {
	conf := &BlockchainConfig{}
	confVal := reflect.ValueOf(conf).Elem()
	for _, item := range params.Config.Items() {
		name := fmt.Sprintf("ConfigParam%d", item.Key)
		field := confVal.FieldByName(name)
		if !field.IsValid() {
			continue
		}
		value := reflect.New(field.Type())
		if err := tlb.Unmarshal(&item.Value.Value, value.Interface()); err != nil {
			return nil, err
		}
		field.Set(value.Elem())
	}
	return conf, nil
}
