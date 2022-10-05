package tongo

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// MerkleProof
// !merkle_proof#03 {X:Type} virtual_hash:bits256 depth:uint16 virtual_root:^X = MERKLE_PROOF X;
type MerkleProof[T any] struct {
	tlb.SumType
	MerkleProof struct {
		VirtualHash Hash
		Depth       uint32 `tlb:"16bits"`
		VirtualRoot tlb.Ref[T]
	} `tlbSumType:"!merkle_proof#03"`
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
	tlb.SumType
	ShardStateUnsplit struct {
		GlobalID        int32
		ShardID         ShardIdent
		SeqNo           uint32
		VertSeqNo       uint32
		GenUtime        uint32
		GenLt           uint64
		MinRefMcSeqno   uint32
		OutMsgQueueInfo tlb.Ref[tlb.Any] // TODO: implement decoding OutMsgQueueInfo fields
		BeforeSplit     bool
		Accounts        tlb.Ref[struct {
			ShardAccounts tlb.HashmapE[tlb.Any] `tlb:"256bits"` // TODO: implement HashmapAugE and some problems with decoding Account
		}]
		Other  tlb.Ref[tlb.Any] // TODO: implement decoding Other fields
		Custom tlb.Maybe[tlb.Ref[McStateExtra]]
	} `tlbSumType:"shard_state#9023afe2"`
	// SplitState struct{} `tlbSumType:"split_state#5f327da5"` // rare case
}

// ShardIdent
// shard_ident$00 shard_pfx_bits:(#<= 60)
// workchain_id:int32 shard_prefix:uint64 = ShardIdent;
type ShardIdent struct {
	tlb.SumType
	ShardIdent struct {
		ShardPfxBits uint64 // TODO: implement lim uint tag
		WorkchainID  int32
		ShardPrefix  uint64
	} `tlbSumType:"shard_ident$00"`
}

func (s *ShardIdent) UnmarshalTLB(c *boc.Cell, tag string) error {
	t, err := c.ReadUint(2)
	if err != nil {
		return err
	}
	if t != 0 {
		return fmt.Errorf("invalid tag")
	}
	prefixBits, err := c.ReadLimUint(60)
	if err != nil {
		return err
	}
	workchain, err := c.ReadInt(32)
	if err != nil {
		return err
	}
	prefix, err := c.ReadUint(64)
	if err != nil {
		return err
	}
	s.SumType = "ShardIdent"
	s.ShardIdent.ShardPfxBits = uint64(prefixBits)
	s.ShardIdent.WorkchainID = int32(workchain)
	s.ShardIdent.ShardPrefix = prefix
	return nil
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
	tlb.SumType
	McStateExtra struct {
		ShardHashes   ShardHashes
		Config        ConfigParams
		Other         tlb.Ref[tlb.Any] // TODO: implement decoding other fields
		GlobalBalance CurrencyCollection
	} `tlbSumType:"masterchain_state_extra#cc26"`
}

// ShardHashes
// _ (HashmapE 32 ^(BinTree ShardDescr)) = ShardHashes;
type ShardHashes struct {
	Hashes tlb.HashmapE[tlb.Ref[tlb.Any]] `tlb:"32bits"` // TODO: implement decoding BinTree ShardDescr
}

// ConfigParams
// _ config_addr:bits256 config:^(Hashmap 32 ^Cell)
// = ConfigParams;
type ConfigParams struct {
	ConfigAddr Hash
	Config     tlb.Ref[tlb.Any] `tlb:"32bits"` // TODO: implement decoding config
}

// DepthBalanceInfo
// depth_balance$_ split_depth:(#<= 30) balance:CurrencyCollection = DepthBalanceInfo;
type DepthBalanceInfo struct {
	SplitDepth uint32
	Balance    CurrencyCollection
}

func (i *DepthBalanceInfo) UnmarshalTLB(c *boc.Cell, tag string) error {
	var balance CurrencyCollection
	splitDepth, err := c.ReadLimUint(30)
	if err != nil {
		return err
	}
	err = tlb.Unmarshal(c, &balance)
	if err != nil {
		return err
	}
	i.SplitDepth = uint32(splitDepth)
	i.Balance = balance
	return nil
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
		Left  ShardStateUnsplit // ^ but decodes manually
		Right ShardStateUnsplit // ^ but decodes manually
	} `tlbSumType:"split_state#5f327da5"`
}

func (s *ShardState) UnmarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("not implemented")
}
