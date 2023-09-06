package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

// MerkleProof
// !merkle_proof#03 {X:Type} virtual_hash:bits256 depth:uint16 virtual_root:^X = MERKLE_PROOF X;
type MerkleProof[T any] struct {
	Magic       Magic `tlb:"!merkle_proof#03"`
	VirtualHash Bits256
	Depth       uint16
	VirtualRoot T `tlb:"^"`
}

// ShardStateUnsplit
// shard_state#9023afe2 global_id:int32
// shard_id:ShardIdent
// seq_no:uint32 vert_seq_no:#
// gen_utime:uint32 gen_lt:uint64
// min_ref_mc_seqno:uint32
// out_msg_queue_info:^OutMsgQueueInfo
// before_split:(## 1)
// accounts:^ShardAccounts
// ^[ overload_history:uint64 underload_history:uint64
// total_balance:CurrencyCollection
// total_validator_fees:CurrencyCollection
// libraries:(HashmapE 256 LibDescr)
// master_ref:(Maybe BlkMasterInfo) ]
// custom:(Maybe ^McStateExtra)
// = ShardStateUnsplit;
type ShardStateUnsplit struct {
	Magic             Magic `tlb:"shard_state#9023afe2"`
	ShardStateUnsplit ShardStateUnsplitData
}

type ShardStateUnsplitData struct {
	GlobalID        int32
	ShardID         ShardIdent
	SeqNo           uint32
	VertSeqNo       uint32
	GenUtime        uint32
	GenLt           uint64
	MinRefMcSeqno   uint32
	OutMsgQueueInfo OutMsgQueueInfo `tlb:"^"`
	BeforeSplit     bool
	Accounts        HashmapAugE[Bits256, ShardAccount, DepthBalanceInfo] `tlb:"^"`
	Other           ShardStateUnsplitOther                               `tlb:"^"`
	Custom          Maybe[Ref[McStateExtra]]
}

// ShardState
// _ ShardStateUnsplit = ShardState;
// split_state#5f327da5 left:^ShardStateUnsplit right:^ShardStateUnsplit = ShardState;
type ShardState struct { // only manual decoding
	SumType
	UnsplitState struct {
		Value ShardStateUnsplit
	} `tlbSumType:"_"`
	SplitState struct {
		Left  ShardStateUnsplit `tlb:"^"` // ^ but decodes manually
		Right ShardStateUnsplit `tlb:"^"` // ^ but decodes manually
	} `tlbSumType:"split_state#5f327da5"`
}

func (s *ShardState) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	sumType, err := c.ReadUint(32)
	if err != nil {
		return err
	}
	switch sumType {
	case 0x5f327da5:
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		err = decoder.Unmarshal(c1, &s.SplitState.Left)
		if err != nil {
			return err
		}
		c1, err = c.NextRef()
		if err != nil {
			return err
		}
		err = decoder.Unmarshal(c1, &s.SplitState.Right)
		if err != nil {
			return err
		}
		break
	case 0x9023afe2:
		var shardUnsplitData ShardStateUnsplitData
		err = decoder.Unmarshal(c, &shardUnsplitData)
		if err != nil {
			return err
		}
		s.UnsplitState.Value.ShardStateUnsplit = shardUnsplitData
		break
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

// ShardIdent
// shard_ident$00 shard_pfx_bits:(#<= 60)
// workchain_id:int32 shard_prefix:uint64 = ShardIdent;
type ShardIdent struct {
	Magic        Magic `tlb:"shardident$00"`
	ShardPfxBits Uint6
	WorkchainID  int32
	ShardPrefix  uint64
}

// sig_pair$_ node_id_short:bits256 sign:CryptoSignature = CryptoSignaturePair;  // 256+x ~ 772 bits
type CryptoSignaturePair struct {
	NodeIdShort Bits256
	Sign        CryptoSignature
}

// ed25519_signature#5 R:bits256 s:bits256 = CryptoSignatureSimple;  // 516 bits
// _ CryptoSignatureSimple = CryptoSignature;
// chained_signature#f signed_cert:^SignedCertificate temp_key_signature:CryptoSignatureSimple
//
//	= CryptoSignature;   // 4+(356+516)+516 = 520 bits+ref (1392 bits total)
type CryptoSignatureSimple struct {
	Magic                 Magic `tlb:"ed25519_signature#5"`
	CryptoSignatureSimple CryptoSignatureSimpleData
}
type CryptoSignature struct {
	SumType
	CryptoSignatureSimple CryptoSignatureSimpleData `tlbSumType:"ed25519_signature#5"`
	CryptoSignature       struct {
		SignedCert       *SignedSertificate `tlb:"^"`
		TempKeySignature CryptoSignatureSimple
	} `tlbSumType:"chained_signature#f"`
}

type CryptoSignatureSimpleData struct {
	R Bits256
	S Bits256
}

func (cr *CryptoSignature) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	sumType, err := c.ReadUint(4)
	if err != nil {
		return err
	}
	if sumType == 0x5 {
		cr.SumType = "CryptoSignatureSimple"
		err = decoder.Unmarshal(c, &cr.CryptoSignatureSimple)
		if err != nil {
			return err
		}
	} else if sumType == 0xf {
		cr.SumType = "CryptoSignature"
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		err = decoder.Unmarshal(c1, &cr.CryptoSignature.SignedCert)
		if err != nil {
			return err
		}
		err = decoder.Unmarshal(c, &cr.CryptoSignature.TempKeySignature)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid tag")
	}
	return nil
}

// signed_certificate$_ certificate:Certificate certificate_signature:CryptoSignature
//
//	= SignedCertificate;  // 356+516 = 872 bits
type SignedSertificate struct {
	Certificate          Certificate
	CertificateSignature CryptoSignature
}

// certificate#4 temp_key:SigPubKey valid_since:uint32
// valid_until:uint32 = Certificate;  // 356 bits
type Certificate struct {
	Magic      Magic `tlb:"certificate#4"`
	TempKey    SigPubKey
	ValidSince uint32
	ValidUntil uint32
}

// _ fees:CurrencyCollection create:CurrencyCollection = ShardFeeCreated;
type ShardFeeCreated struct {
	Fees   CurrencyCollection
	Create CurrencyCollection
}

// _ (HashmapAugE 96 ShardFeeCreated ShardFeeCreated) = ShardFees;
type ShardFees struct {
	Hashmap HashmapAugE[Bits96, ShardFeeCreated, ShardFeeCreated]
}

// acc_trans#5 account_addr:bits256
//
//	  transactions:(HashmapAug 64 ^Transaction CurrencyCollection)
//	  state_update:^(HASH_UPDATE Account)
//	= AccountBlock;
type AccountBlock struct {
	Magic        Magic `tlb:"acc_trans#5"`
	AccountAddr  Bits256
	Transactions HashmapAug[Uint64, Ref[Transaction], CurrencyCollection]
	StateUpdate  HashUpdate `tlb:"^"`
}

// processed_upto$_ last_msg_lt:uint64 last_msg_hash:bits256 = ProcessedUpto;
type ProcessedUpto struct {
	LastMsgLt   uint64
	LastMsgHash Bits256
}

// ihr_pending$_ import_lt:uint64 = IhrPendingSince;
type IhrPendingSince struct {
	ImportLt uint64
}

// depth_balance$_ split_depth:(#<= 30) balance:CurrencyCollection = DepthBalanceInfo;
type DepthBalanceInfo struct {
	SplitDepth Uint5
	Balance    CurrencyCollection
}

// ^[ overload_history:uint64 underload_history:uint64
// total_balance:CurrencyCollection
// total_validator_fees:CurrencyCollection
// libraries:(HashmapE 256 LibDescr)
// master_ref:(Maybe BlkMasterInfo) ]
type ShardStateUnsplitOther struct {
	OverloadHistory    uint64
	UnderloadHistory   uint64
	TotalBalance       CurrencyCollection
	TotalValidatorFees CurrencyCollection
	Libraries          HashmapE[Bits256, LibDescr]
	MasterRef          Maybe[BlkMasterInfo]
}

// shared_lib_descr$00 lib:^Cell publishers:(Hashmap 256 True)
//
//	= LibDescr;
type LibDescr struct {
	Magic      Magic    `tlb:"shared_lib_descr$00"`
	Lib        boc.Cell `tlb:"^"`
	Publishers Hashmap[Bits256, bool]
}

// McStateExtra
// masterchain_state_extra#cc26
// shard_hashes:ShardHashes
// config:ConfigParams
// ^[ flags:(## 16) { flags <= 1 }
// validator_info:ValidatorInfo
// prev_blocks:OldMcBlocksInfo
// after_key_block:Bool
// last_key_block:(Maybe ExtBlkRef)
// block_create_stats:(flags . 0)?BlockCreateStats ]
// global_balance:CurrencyCollection
// = McStateExtra;
type McStateExtra struct {
	Magic         Magic `tlb:"masterchain_state_extra#cc26"`
	ShardHashes   HashmapE[Uint32, Ref[ShardInfoBinTree]]
	Config        ConfigParams
	Other         McStateExtraOther `tlb:"^"`
	GlobalBalance CurrencyCollection
}

// ConfigParams
// _ config_addr:bits256 config:^(Hashmap 32 ^Cell)
// = ConfigParams;
type ConfigParams struct {
	ConfigAddr Bits256
	Config     Hashmap[Uint32, Ref[boc.Cell]] `tlb:"^"`
}

// CloneKeepingSubsetOfKeys returns a new ConfigParams with only the keys specified in the keys parameter.
func (params *ConfigParams) CloneKeepingSubsetOfKeys(keys []uint32) ConfigParams {
	keysMap := make(map[uint32]struct{})
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}
	var newKeys []Uint32
	var newValues []Ref[boc.Cell]
	for _, item := range params.Config.Items() {
		if _, ok := keysMap[uint32(item.Key)]; ok {
			newKeys = append(newKeys, item.Key)
			newValues = append(newValues, item.Value)
		}
	}
	return ConfigParams{
		ConfigAddr: params.ConfigAddr,
		Config: Hashmap[Uint32, Ref[boc.Cell]]{
			keys:   newKeys,
			values: newValues,
		},
	}
}

// ^[ flags:(## 16) { flags <= 1 }
// validator_info:ValidatorInfo
// prev_blocks:OldMcBlocksInfo
// after_key_block:Bool
// last_key_block:(Maybe ExtBlkRef)
// block_create_stats:(flags . 0)?BlockCreateStats ]
type McStateExtraOther struct {
	Flags            uint16
	ValidatorInfo    ValidatorInfo
	PrevBlocks       HashmapAugE[Uint32, KeyExtBlkRef, KeyMaxLt]
	AfterKeyBlock    bool
	LastKeyBlock     Maybe[ExtBlkRef]
	BlockCreateStats BlockCreateStats
}

func (m *McStateExtraOther) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	flags, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	m.Flags = uint16(flags)
	err = decoder.Unmarshal(c, &m.ValidatorInfo)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.PrevBlocks)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.AfterKeyBlock)
	if err != nil {
		return err
	}
	err = decoder.Unmarshal(c, &m.LastKeyBlock)
	if err != nil {
		return err
	}
	if m.Flags == 1 {
		err = decoder.Unmarshal(c, &m.BlockCreateStats)
		if err != nil {
			return err
		}
	}

	return nil
}

// _ key:Bool max_end_lt:uint64 = KeyMaxLt;
type KeyMaxLt struct {
	Key      bool
	MaxEndLt uint64
}

// _ key:Bool blk_ref:ExtBlkRef = KeyExtBlkRef;
type KeyExtBlkRef struct {
	Key    bool
	BlkRef ExtBlkRef
}

// block_create_stats#17 counters:(HashmapE 256 CreatorStats) = BlockCreateStats;
// block_create_stats_ext#34 counters:(HashmapAugE 256 CreatorStats uint32) = BlockCreateStats;
type BlockCreateStats struct {
	SumType
	BlockCreateStats struct {
		Counters HashmapE[Bits256, CreatorStats]
	} `tlbSumType:"block_create_stats#17"`
	BlockCreateStatsExt struct {
		Counters HashmapAugE[Bits256, CreatorStats, uint32]
	} `tlbSumType:"block_create_stats_ext#34"`
}

// creator_info#4 mc_blocks:Counters shard_blocks:Counters = CreatorStats;
type CreatorStats struct {
	Magic       Magic `tlb:"creator_info#4"`
	McBlocks    Counters
	ShardBlocks Counters
}

// counters#_ last_updated:uint32 total:uint64 cnt2048:uint64 cnt65536:uint64 = Counters;

type Counters struct {
	LastUpdated uint32
	Total       uint64
	Cnt2048     uint64
	Cnt65536    uint64
}
