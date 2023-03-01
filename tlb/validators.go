package tlb

// validator_info$_
//
//	validator_list_hash_short:uint32
//	catchain_seqno:uint32
//	nx_cc_updated:Bool
//
// = ValidatorInfo;
type ValidatorInfo struct {
	ValidatorListHashShort uint32
	CatchainSeqno          uint32
	NxCcUpdated            bool
}

// validator_base_info$_
//
//	validator_list_hash_short:uint32
//	catchain_seqno:uint32
//
// = ValidatorBaseInfo;
type ValidatorBaseInfo struct {
	ValidatorListHashShort uint32
	CatchainSeqno          uint32
}

type ValidatorsSet struct {
	SumType
	// validators#11 utime_since:uint32 utime_until:uint32
	//   total:(## 16) main:(## 16) { main <= total } { main >= 1 }
	//   list:(Hashmap 16 ValidatorDescr) = ValidatorSet;
	Validators *struct {
		UtimeSince uint32
		UtimeUntil uint32
		Total      uint16
		Main       uint16
		List       Hashmap[Uint16, ValidatorDescr]
	} `tlbSumType:"validators#11"`
	// validators_ext#12 utime_since:uint32 utime_until:uint32
	//   total:(## 16) main:(## 16) { main <= total } { main >= 1 }
	//   total_weight:uint64 list:(HashmapE 16 ValidatorDescr) = ValidatorSet;
	ValidatorsExt *struct {
		UtimeSince  uint32
		UtimeUntil  uint32
		Total       uint16
		Main        uint16
		TotalWeight uint64
		List        HashmapE[Uint16, ValidatorDescr]
	} `tlbSumType:"validatorsext#12"`
}

type ValidatorDescr struct {
	SumType
	// validator#53 public_key:SigPubKey weight:uint64 = ValidatorDescr;
	Validator *struct {
		PublicKey SigPubKey
		Weight    uint64
	} `tlbSumType:"validator#53"`
	// validator_addr#73 public_key:SigPubKey weight:uint64 adnl_addr:bits256 = ValidatorDescr;
	ValidatorAddr *struct {
		PublicKey SigPubKey
		Weight    uint64
		AdnlAddr  Bits256
	} `tlbSumType:"validatoraddr#73"`
}

type SigPubKey struct {
	Magic  Magic `tlb:"pubkey#8e81278a"`
	PubKey Bits256
}

//ed25519_pubkey#8e81278a pubkey:bits256 = SigPubKey;
