package tongo

import (
	"github.com/tonkeeper/tongo/tlb"
)

// validator_info$_
//   validator_list_hash_short:uint32
//   catchain_seqno:uint32
//   nx_cc_updated:Bool
// = ValidatorInfo;
type ValidatorInfo struct {
	ValidatorListHashShort uint32
	CatchainSeqno          uint32
	NxCcUpdated            bool
}

// validator_base_info$_
//   validator_list_hash_short:uint32
//   catchain_seqno:uint32
// = ValidatorBaseInfo;
type ValidatorBaseInfo struct {
	ValidatorListHashShort uint32
	CatchainSeqno          uint32
}

type ValidatorsSet struct {
	tlb.SumType
	// validators#11 utime_since:uint32 utime_until:uint32
	//   total:(## 16) main:(## 16) { main <= total } { main >= 1 }
	//   list:(Hashmap 16 ValidatorDescr) = ValidatorSet;
	Validators struct {
		UtimeSince uint32
		UtimeUntil uint32
		Total      uint32                      `tlb:"16bits"`
		Main       uint32                      `tlb:"16bits"`
		List       tlb.Hashmap[ValidatorDescr] `tlb:"16bits"`
	} `tlbSumType:"validators#11"`
	// validators_ext#12 utime_since:uint32 utime_until:uint32
	//   total:(## 16) main:(## 16) { main <= total } { main >= 1 }
	//   total_weight:uint64 list:(HashmapE 16 ValidatorDescr) = ValidatorSet;
	ValidatorsExt struct {
		UtimeSince  uint32
		UtimeUntil  uint32
		Total       uint32 `tlb:"16bits"`
		Main        uint32 `tlb:"16bits"`
		TotalWeight uint64
		List        tlb.HashmapE[ValidatorDescr] `tlb:"16bits"`
	} `tlbSumType:"validatorsext#12"`
}

type ValidatorDescr struct {
	tlb.SumType
	// validator#53 public_key:SigPubKey weight:uint64 = ValidatorDescr;
	Validator struct {
		PublicKey SigPubKey
		Weight    uint64
	} `tlbSumType:"validator#53"`
	// validator_addr#73 public_key:SigPubKey weight:uint64 adnl_addr:bits256 = ValidatorDescr;
	ValidatorAddr struct {
		PublicKey SigPubKey
		Weight    uint64
		AdnlAddr  Hash
	} `tlbSumType:"validatoraddr#73"`
}

type SigPubKey struct {
	Magic  tlb.Magic `tlb:"pubkey#8e81278a"`
	PubKey Hash
}

//ed25519_pubkey#8e81278a pubkey:bits256 = SigPubKey;
