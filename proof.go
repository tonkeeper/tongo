package tongo

import (
	"fmt"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// MerkleProof
// !merkle_proof#03 {X:Type} virtual_hash:bits256 depth:uint16 virtual_root:^X = MERKLE_PROOF X;
type MerkleProof[T any] struct {
	Magic       tlb.Magic `tlb:"!merkle_proof#03"`
	VirtualHash Hash
	Depth       boc.BitString `tlb:"16bits"`
	VirtualRoot T             `tlb:"^"`
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
	Magic             tlb.Magic `tlb:"shard_state#9023afe2"`
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
	Accounts        ShardAccounts          `tlb:"^"`
	Other           ShardStateUnsplitOther `tlb:"^"`
	Custom          tlb.Maybe[tlb.Ref[McStateExtra]]
}

// ShardState
// _ ShardStateUnsplit = ShardState;
// split_state#5f327da5 left:^ShardStateUnsplit right:^ShardStateUnsplit = ShardState;
type ShardState struct { // only manual decoding
	tlb.SumType
	UnsplitState struct {
		Value ShardStateUnsplit
	} `tlbSumType:"_"`
	SplitState struct {
		Left  ShardStateUnsplit `tlb:"^"` // ^ but decodes manually
		Right ShardStateUnsplit `tlb:"^"` // ^ but decodes manually
	} `tlbSumType:"split_state#5f327da5"`
}

func (s *ShardState) UnmarshalTLB(c *boc.Cell, tag string) error {
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
		err = tlb.Unmarshal(c1, &s.SplitState.Left)
		if err != nil {
			return err
		}
		c1, err = c.NextRef()
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c1, &s.SplitState.Right)
		if err != nil {
			return err
		}
		break
	case 0x9023afe2:
		var shardUnsplitData ShardStateUnsplitData
		err = tlb.Unmarshal(c, &shardUnsplitData)
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
	Magic        tlb.Magic `tlb:"shardident$00"`
	ShardPfxBits uint64    `tlb:"6bits"`
	WorkchainID  int32
	ShardPrefix  uint64
}

// sig_pair$_ node_id_short:bits256 sign:CryptoSignature = CryptoSignaturePair;  // 256+x ~ 772 bits
type CryptoSignaturePair struct {
	NodeIdShort Hash
	Sign        CryptoSignature
}

// ed25519_signature#5 R:bits256 s:bits256 = CryptoSignatureSimple;  // 516 bits
// _ CryptoSignatureSimple = CryptoSignature;
// chained_signature#f signed_cert:^SignedCertificate temp_key_signature:CryptoSignatureSimple
//   = CryptoSignature;   // 4+(356+516)+516 = 520 bits+ref (1392 bits total)
type CryptoSignatureSimple struct {
	Magic                 tlb.Magic `tlb:"ed25519_signature#5"`
	CryptoSignatureSimple CryptoSignatureSimpleData
}
type CryptoSignature struct {
	tlb.SumType
	CryptoSignatureSimple CryptoSignatureSimpleData `tlbSumType:"ed25519_signature#5"`
	CryptoSignature       struct {
		SignedCert       *SignedSertificate `tlb:"^"`
		TempKeySignature CryptoSignatureSimple
	} `tlbSumType:"chained_signature#f"`
}

type CryptoSignatureSimpleData struct {
	R Hash
	S Hash
}

func (cr *CryptoSignature) UnmarshalTLB(c *boc.Cell, tag string) error {
	sumType, err := c.ReadUint(4)
	if err != nil {
		return err
	}
	if sumType == 0x5 {
		cr.SumType = "CryptoSignatureSimple"
		err = tlb.Unmarshal(c, &cr.CryptoSignatureSimple)
		if err != nil {
			return err
		}
	} else if sumType == 0xf {
		cr.SumType = "CryptoSignature"
		c1, err := c.NextRef()
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c1, &cr.CryptoSignature.SignedCert)
		if err != nil {
			return err
		}
		err = tlb.Unmarshal(c, &cr.CryptoSignature.TempKeySignature)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid tag")
	}
	return nil
}

// signed_certificate$_ certificate:Certificate certificate_signature:CryptoSignature
//   = SignedCertificate;  // 356+516 = 872 bits
type SignedSertificate struct {
	Certificate          Certificate
	CertificateSignature CryptoSignature
}

// certificate#4 temp_key:SigPubKey valid_since:uint32
// valid_until:uint32 = Certificate;  // 356 bits
type Certificate struct {
	Magic      tlb.Magic `tlb:"certificate#4"`
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
	Hashmap tlb.HashmapAugE[ShardFeeCreated, ShardFeeCreated] `tlb:"96bits"`
}

// _ (HashmapAugE 256 InMsg ImportFees) = InMsgDescr;
type InMsgDescr struct {
	Hashmap tlb.HashmapAugE[InMsg, ImportFees] `tlb:"256bits"`
}

// _ (HashmapAugE 256 AccountBlock CurrencyCollection) = ShardAccountBlocks
type ShardAccountBlocks struct {
	Hashmap tlb.HashmapAugE[AccountBlock, CurrencyCollection] `tlb:"256bits"`
}

// acc_trans#5 account_addr:bits256
//             transactions:(HashmapAug 64 ^Transaction CurrencyCollection)
//             state_update:^(HASH_UPDATE Account)
//           = AccountBlock;
type AccountBlock struct {
	Magic        tlb.Magic `tlb:"acc_trans#5"`
	AccountAddr  Hash
	Transactions tlb.HashmapAug[tlb.Ref[Transaction], CurrencyCollection] `tlb:"64bits"`
	StateUpdate  HashUpdate                                               `tlb:"^"`
}

// _ (HashmapE 96 ProcessedUpto) = ProcessedInfo;
type ProcessedInfo struct {
	ProcessedInfo tlb.HashmapE[tlb.Size96, ProcessedUpto]
}

// processed_upto$_ last_msg_lt:uint64 last_msg_hash:bits256 = ProcessedUpto;
type ProcessedUpto struct {
	LastMsg     uint64
	LastMsgHash Hash
}

// _ (HashmapE 320 IhrPendingSince) = IhrPendingInfo;
type IhrPendingInfo struct {
	IhrPendingInfo tlb.HashmapE[tlb.Size320, IhrPendingSince]
}

// ihr_pending$_ import_lt:uint64 = IhrPendingSince;
type IhrPendingSince struct {
	ImportLt uint64
}

// _ (HashmapAugE 256 ShardAccount DepthBalanceInfo) = ShardAccounts;
type ShardAccounts struct {
	Accounts tlb.HashmapAugE[ShardAccount, DepthBalanceInfo] `tlb:"256bits"`
}

// depth_balance$_ split_depth:(#<= 30) balance:CurrencyCollection = DepthBalanceInfo;
type DepthBalanceInfo struct {
	SplitDepth uint32 `tlb:"5bits"`
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
	Libraries          tlb.HashmapE[tlb.Size256, LibDescr]
	MasterRef          tlb.Maybe[BlkMasterInfo]
}

// shared_lib_descr$00 lib:^Cell publishers:(Hashmap 256 True)
//   = LibDescr;
type LibDescr struct {
	Magic      tlb.Magic `tlb:"shared_lib_descr$00"`
	Lib        boc.Cell  `tlb:"^"`
	Publishers tlb.Hashmap[tlb.Size256, bool]
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
	Magic         tlb.Magic `tlb:"masterchain_state_extra#cc26"`
	ShardHashes   ShardHashes
	Config        ConfigParams
	Other         McStateExtraOther `tlb:"^"`
	GlobalBalance CurrencyCollection
}

// ShardHashes
// _ (HashmapE 32 ^(BinTree ShardDescr)) = ShardHashes;
type ShardHashes struct {
	Hashes tlb.HashmapE[tlb.KeySize32, tlb.Ref[ShardInfoBinTree]]
}

type ConfigHashMap struct {
	Hashmap tlb.Hashmap[tlb.KeySize32, tlb.Ref[boc.Cell]]
}

// ConfigParams
// _ config_addr:bits256 config:^(Hashmap 32 ^Cell)
// = ConfigParams;
type ConfigParams struct {
	ConfigAddr Hash
	Config     ConfigHashMap `tlb:"^"`
}

// ^[ flags:(## 16) { flags <= 1 }
// validator_info:ValidatorInfo
// prev_blocks:OldMcBlocksInfo
// after_key_block:Bool
// last_key_block:(Maybe ExtBlkRef)
// block_create_stats:(flags . 0)?BlockCreateStats ]
type McStateExtraOther struct {
	Flags            uint32 `tlb:"16bits"`
	ValidatorInfo    ValidatorInfo
	PrevBlocks       OldMcBlocksInfo
	AfterKeyBlock    bool
	LastKeyBlock     tlb.Maybe[ExtBlkRef]
	BlockCreateStats BlockCreateStats
}

func (m *McStateExtraOther) UnmarshalTLB(c *boc.Cell, tag string) error {
	flags, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	m.Flags = uint32(flags)
	err = tlb.Unmarshal(c, &m.ValidatorInfo)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.PrevBlocks)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.AfterKeyBlock)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &m.LastKeyBlock)
	if err != nil {
		return err
	}
	if m.Flags == 1 {
		err = tlb.Unmarshal(c, &m.BlockCreateStats)
		if err != nil {
			return err
		}
	}

	return nil
}

// _ (HashmapAugE 32 KeyExtBlkRef KeyMaxLt) = OldMcBlocksInfo
type OldMcBlocksInfo struct {
	Info tlb.HashmapAugE[KeyExtBlkRef, KeyMaxLt] `tlb:"32bits"`
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
	tlb.SumType
	BlockCreateStats struct {
		Counters tlb.HashmapE[tlb.Size256, CreatorStats]
	} `tlbSumType:"block_create_stats#17"`
	BlockCreateStatsExt struct {
		Counters tlb.HashmapAugE[CreatorStats, uint32] `tlb:"256bits"`
	} `tlbSumType:"block_create_statsext#34"`
}

// creator_info#4 mc_blocks:Counters shard_blocks:Counters = CreatorStats;
type CreatorStats struct {
	Magic       tlb.Magic `tlb:"creator_info#4"`
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
