package tychoclient

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

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

	ConfigParam28 *ConfigParam28     `json:",omitempty"`
	ConfigParam29 *ConfigParam29     `json:",omitempty"`
	ConfigParam31 *tlb.ConfigParam31 `json:",omitempty"`
	ConfigParam32 *tlb.ConfigParam32 `json:",omitempty"`
	ConfigParam33 *tlb.ConfigParam33 `json:",omitempty"`
	ConfigParam34 *tlb.ConfigParam34 `json:",omitempty"`
	ConfigParam35 *tlb.ConfigParam35 `json:",omitempty"`
	ConfigParam36 *tlb.ConfigParam36 `json:",omitempty"`
	ConfigParam37 *tlb.ConfigParam37 `json:",omitempty"`
	ConfigParam39 *tlb.ConfigParam39 `json:",omitempty"`
	ConfigParam40 *tlb.ConfigParam40 `json:",omitempty"`
	//ConfigParam43 *tlb.ConfigParam43 `json:",omitempty"` // TODO: replace
	ConfigParam44 *tlb.ConfigParam44 `json:",omitempty"`
	ConfigParam45 *tlb.ConfigParam45 `json:",omitempty"`

	ConfigParam71 *tlb.ConfigParam71 `json:",omitempty"`
	ConfigParam72 *tlb.ConfigParam72 `json:",omitempty"`
	ConfigParam73 *tlb.ConfigParam73 `json:",omitempty"`

	ConfigParam79 *tlb.ConfigParam79 `json:",omitempty"`
	ConfigParam81 *tlb.ConfigParam81 `json:",omitempty"`
	ConfigParam82 *tlb.ConfigParam82 `json:",omitempty"`

	// TODO: ConfigParam50
	// TODO: ConfigParam100

	// Negative keys don't have a schema,
	// so we store them as raw cells.

	ConfigParamNegative71  *boc.Cell `json:",omitempty"`
	ConfigParamNegative999 *boc.Cell `json:",omitempty"`
}

type ConfigParam28 = CollationConfig

// TODO: 00 tags?

// CollationConfig (only v2)
// collation_config_tycho_v2#a7
// shuffle_mc_validators:Bool
// mc_block_min_interval_ms:uint32
// mc_block_max_interval_ms:uint32
// empty_sc_block_interval_ms:uint32
// max_uncommitted_chain_length:uint8
// wu_used_to_import_next_anchor:uint64
// msgs_exec_params:MsgsExecutionParams
// work_units_params:WorkUnitsParams
// = CollationConfig;
type CollationConfig struct {
	Magic                     tlb.Magic `tlb:"#a7"`
	ShuffleMcValidators       bool
	McBlockMinIntervalMs      uint32
	McBlockMaxIntervalMs      uint32
	EmptyScBlockIntervalMs    uint32
	MaxUncommittedChainLength uint8
	WuUsedToImportNextAnchor  uint64
	MsgsExecParams            MsgsExecutionParams // TODO: clarify tags
	WorkUnitsParams           WorkUnitsParams
}

// WorkUnitsParams
// work_units_params_tycho#00
// prepare:WorkUnitParamsPrepare // TODO: missprint : WorkUnitsParamsPrepare
// execute:WorkUnitParamsExecute // TODO: missprint : WorkUnitsParamsExecute
// finalize:WorkUnitParamsFinalize // TODO: missprint : WorkUnitsParamsFinalize
// = WorkUnitsParams;
type WorkUnitsParams struct {
	Prepare  WorkUnitsParamsPrepare
	Execute  WorkUnitsParamsExecute
	Finalize WorkUnitsParamsFinalize
}

// WorkUnitsParamsPrepare
// work_units_params_prepare_tycho#00
// fixed:uint32
// msgs_stats:uint16
// remaning_msgs_stats:uint16
// read_ext_msgs:uint16
// read_int_msgs:uint16
// read_new_msgs:uint16
// add_to_msg_groups:uint16
// = WorkUnitsParamsPrepare;
type WorkUnitsParamsPrepare struct {
	Fixed             uint32
	MsgsStats         uint16
	RemaningMsgsStats uint16
	ReadExtMsgs       uint16
	ReadIntMsgs       uint16
	ReadNewMsgs       uint16
	AddToMsgGroups    uint16
}

// WorkUnitsParamsExecute
// work_units_params_execute_tycho#00
// prepare:uint32
// execute:uint16
// execute_err:uint16
// execute_delimiter:uint32
// serialize_enqueue:uint16
// serialize_dequeue:uint16
// insert_new_msgs:uint16
// subgroup_size:uint16
// = WorkUnitsParamsExecute;
type WorkUnitsParamsExecute struct {
	Prepare          uint32
	Execute          uint16
	ExecuteErr       uint16
	ExecuteDelimiter uint32
	SerializeEnqueue uint16
	SerializeDequeue uint16
	InsertNewMsgs    uint16
	SubgroupSize     uint16
}

// WorkUnitsParamsFinalize
// work_units_params_finalize_tycho#00
// build_transactions:uint16
// build_accounts:uint16
// build_in_msg:uint16
// build_out_msg:uint16
// serialize_min:uint32
// serialize_accounts:uint16
// serialize_msg:uint16
// state_update_min:uint32
// state_update_accounts:uint16
// state_update_msg:uint16
// create_diff:uint16
// serialize_diff:uint16
// apply_diff:uint16
// diff_tail_len:uint16
// = WorkUnitsParamsFinalize;
type WorkUnitsParamsFinalize struct {
	BuildTransactions   uint16
	BuildAccounts       uint16
	BuildInMsg          uint16
	BuildOutMsg         uint16
	SerializeMin        uint32
	SerializeAccounts   uint16
	SerializeMsg        uint16
	StateUpdateMin      uint32
	StateUpdateAccounts uint16
	StateUpdateMsg      uint16
	CreateDiff          uint16
	SerializeDiff       uint16
	ApplyDiff           uint16
	DiffTailLen         uint16
}

// consensus_config_tycho#d8
// clock_skew_millis:uint16 { clock_skew_millis != 0 }
// payload_batch_bytes:uint32 { payload_batch_bytes != 0 }
// commit_history_rounds:uint8 { commit_history_rounds != 0 }
// deduplicate_rounds:uint16
// max_consensus_lag_rounds:uint16 { max_consensus_lag_rounds != 0 }
// payload_buffer_bytes:uint32 { payload_buffer_bytes != 0 }
// broadcast_retry_millis:uint8 { broadcast_retry_millis != 0 }
// download_retry_millis:uint8 { download_retry_millis != 0 }
// download_peers:uint8 { download_peers != 0 }
// min_sign_attempts:uint8 { min_sign_attempts != 0 }
// download_peer_queries:uint8 { download_peer_queries != 0 }
// sync_support_rounds:uint16 { sync_support_rounds != 0 }
// = ConsensusConfig;

type ConfigParam29 = ConsensusConfig
type ConsensusConfig struct {
	Magic                 tlb.Magic `tlb:"#d8"`
	ClockSkewMillis       uint16
	PayloadBatchBytes     uint32
	CommitHistoryRounds   uint8
	DeduplicateRounds     uint16
	MaxConsensusLagRounds uint16
	PayloadBufferBytes    uint32
	BroadcastRetryMillis  uint8
	DownloadRetryMillis   uint8
	DownloadPeers         uint8
	MinSignAttempts       uint8
	DownloadPeerQueries   uint8
	SyncSupportRounds     uint16
}

func ConvertBlockchainConfig(params tlb.ConfigParams, ignoreBrokenParams bool) (*BlockchainConfig, []int, error) {
	var brokenParams []int
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
				if ignoreBrokenParams {
					brokenParams = append(brokenParams, int(key))
					continue
				}
				return nil, nil, err
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
			if ignoreBrokenParams {
				brokenParams = append(brokenParams, int(key))
				continue
			}
			return nil, nil, err
		}
		field.Set(reflect.ValueOf(cell))
	}
	return conf, brokenParams, nil
}

func ConvertBlockchainConfigStrict(params tlb.ConfigParams) (*BlockchainConfig, error) {
	conf, _, err := ConvertBlockchainConfig(params, false)
	if err != nil {
		return nil, err
	}
	return conf, nil
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
