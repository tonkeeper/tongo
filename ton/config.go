package ton

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// BlockchainConfig represents the TON blockchain configuration stored inside key blocks.
// This struct is a compromise between
// a low-level cell-based representation of the config which is tlb.ConfigParams and something more convenient,
// BlockchainConfig uses lots of tlb types but
// 1. these tlb types are easy to navigate and
// 2. BlockchainConfig can be easily serialized to JSON.
type BlockchainConfig struct {
	ConfigParam0  *tlb.ConfigParam0  `json:",omitempty"`
	ConfigParam1  *tlb.ConfigParam1  `json:",omitempty"`
	ConfigParam2  *tlb.ConfigParam2  `json:",omitempty"`
	ConfigParam3  *tlb.ConfigParam3  `json:",omitempty"`
	ConfigParam4  *tlb.ConfigParam4  `json:",omitempty"`
	ConfigParam5  *tlb.ConfigParam5  `json:",omitempty"`
	ConfigParam6  *tlb.ConfigParam6  `json:",omitempty"`
	ConfigParam7  *tlb.ConfigParam7  `json:",omitempty"`
	ConfigParam8  *tlb.ConfigParam8  `json:",omitempty"`
	ConfigParam9  *tlb.ConfigParam9  `json:",omitempty"`
	ConfigParam10 *tlb.ConfigParam10 `json:",omitempty"`
	ConfigParam11 *tlb.ConfigParam11 `json:",omitempty"`
	ConfigParam12 *tlb.ConfigParam12 `json:",omitempty"`
	ConfigParam13 *tlb.ConfigParam13 `json:",omitempty"`
	ConfigParam14 *tlb.ConfigParam14 `json:",omitempty"`
	ConfigParam15 *tlb.ConfigParam15 `json:",omitempty"`
	ConfigParam16 *tlb.ConfigParam16 `json:",omitempty"`
	ConfigParam17 *tlb.ConfigParam17 `json:",omitempty"`
	ConfigParam18 *tlb.ConfigParam18 `json:",omitempty"`
	ConfigParam20 *tlb.ConfigParam20 `json:",omitempty"`
	ConfigParam21 *tlb.ConfigParam21 `json:",omitempty"`
	ConfigParam22 *tlb.ConfigParam22 `json:",omitempty"`
	ConfigParam23 *tlb.ConfigParam23 `json:",omitempty"`
	ConfigParam24 *tlb.ConfigParam24 `json:",omitempty"`
	ConfigParam25 *tlb.ConfigParam25 `json:",omitempty"`

	ConfigParam28 *tlb.ConfigParam28 `json:",omitempty"`
	ConfigParam29 *tlb.ConfigParam29 `json:",omitempty"`
	ConfigParam31 *tlb.ConfigParam31 `json:",omitempty"`
	ConfigParam32 *tlb.ConfigParam32 `json:",omitempty"`
	ConfigParam33 *tlb.ConfigParam33 `json:",omitempty"`
	ConfigParam34 *tlb.ConfigParam34 `json:",omitempty"`
	ConfigParam35 *tlb.ConfigParam35 `json:",omitempty"`
	ConfigParam36 *tlb.ConfigParam36 `json:",omitempty"`
	ConfigParam37 *tlb.ConfigParam37 `json:",omitempty"`
	ConfigParam39 *tlb.ConfigParam39 `json:",omitempty"`
	ConfigParam40 *tlb.ConfigParam40 `json:",omitempty"`
	ConfigParam43 *tlb.ConfigParam43 `json:",omitempty"`
	ConfigParam44 *tlb.ConfigParam44 `json:",omitempty"`

	ConfigParam71 *tlb.ConfigParam71 `json:",omitempty"`
	ConfigParam72 *tlb.ConfigParam72 `json:",omitempty"`
	ConfigParam73 *tlb.ConfigParam73 `json:",omitempty"`

	ConfigParam79 *tlb.ConfigParam79 `json:",omitempty"`
	ConfigParam81 *tlb.ConfigParam81 `json:",omitempty"`
	ConfigParam82 *tlb.ConfigParam82 `json:",omitempty"`

	// Negative keys don't have a schema,
	// so we store them as raw cells.

	ConfigParamNegative71  *boc.Cell `json:",omitempty"`
	ConfigParamNegative999 *boc.Cell `json:",omitempty"`
}

func (conf *BlockchainConfig) ConfigAddr() (AccountID, bool) {
	if conf.ConfigParam0 != nil {
		return AccountID{Workchain: -1, Address: conf.ConfigParam0.ConfigAddr}, true
	}
	return AccountID{}, false
}

func (conf *BlockchainConfig) ElectorAddr() (AccountID, bool) {
	if conf.ConfigParam1 != nil {
		return AccountID{Workchain: -1, Address: conf.ConfigParam1.ElectorAddr}, true
	}
	return AccountID{}, false
}

func (conf *BlockchainConfig) MinterAddr() (AccountID, bool) {
	if conf.ConfigParam2 != nil {
		return AccountID{Workchain: -1, Address: conf.ConfigParam2.MinterAddr}, true
	}
	return AccountID{}, false
}

func (conf *BlockchainConfig) FeeCollectorAddr() (AccountID, bool) {
	if conf.ConfigParam3 != nil {
		return AccountID{Workchain: -1, Address: conf.ConfigParam3.FeeCollectorAddr}, true
	}
	return AccountID{}, false
}

func (conf *BlockchainConfig) DnsRootAddr() (AccountID, bool) {
	if conf.ConfigParam4 != nil {
		return AccountID{Workchain: -1, Address: conf.ConfigParam4.DnsRootAddr}, true
	}
	return AccountID{}, false
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
		key := int32(item.Key)
		if key >= 0 {
			name := fmt.Sprintf("ConfigParam%d", key)
			field := confVal.FieldByName(name)
			if !field.IsValid() {
				continue
			}
			value := reflect.New(field.Type())
			if err := tlb.Unmarshal(&item.Value.Value, value.Interface()); err != nil {
				return nil, err
			}
			field.Set(value.Elem())
			continue
		}
		// negative key
		name := fmt.Sprintf("ConfigParamNegative%d", -key)
		field := confVal.FieldByName(name)
		if !field.IsValid() {
			continue
		}
		cell := boc.NewCell()
		if err := tlb.Unmarshal(&item.Value.Value, cell); err != nil {
			return nil, err
		}
		field.Set(reflect.ValueOf(cell))
	}
	return conf, nil
}

func DecodeConfigParams(b []byte) (tlb.ConfigParams, error) {
	cells, err := boc.DeserializeBoc(b)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	if len(cells) != 1 {
		return tlb.ConfigParams{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	if proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Exists {
		return proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Value.Value.Config, nil
	}
	return tlb.ConfigParams{}, fmt.Errorf("empty Custom field")
}
