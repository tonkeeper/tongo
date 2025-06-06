package liteapi

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"hash/crc32"
	"sort"
)

const (
	magicTonBlockID   = 0xc50b6e70 // crc32(ton.blockId root_cell_hash:int256 file_hash:int256 = ton.BlockId)
	magicValidatorSet = 0x901660ed // crc32(???) todo find out where this opcode came from
	magicPublicKey    = 0x4813b4c6 // crc32(pub.ed25519 key:int256 = PublicKey)
)

type signature struct {
	nodeIdShort [32]byte
	signature   []byte
}
type signatures struct {
	validatorListHashShort uint32
	catchainSeqno          uint32
	signatures             []signature
}

func (c *Client) VerifyProofChain(ctx context.Context, source, target ton.BlockIDExt) error {
	isForward := source.Seqno < target.Seqno

	for source.Seqno != target.Seqno {
		partialBlockProof, err := c.GetBlockProof(ctx, source, &target)
		if err != nil {
			return fmt.Errorf("cannot get partial block proof from liteserver: %w", err)
		}
		partialBlockSource := partialBlockProof.From.ToBlockIdExt()
		if partialBlockSource != source {
			return fmt.Errorf("incorrect source partial block: got %v, want %v", partialBlockProof.From.Seqno, source.Seqno)
		}

		for _, step := range partialBlockProof.Steps {
			switch step.SumType {
			case "LiteServerBlockLinkBack":
				if isForward {
					return fmt.Errorf("blocks can be linked forward if direction is backward")
				}
				toKeyBlock := step.LiteServerBlockLinkBack.ToKeyBlock
				sourceBlock := step.LiteServerBlockLinkBack.From.ToBlockIdExt()
				targetBlock := step.LiteServerBlockLinkBack.To.ToBlockIdExt()
				destProof, err := boc.DeserializeBoc(step.LiteServerBlockLinkBack.DestProof)
				if err != nil {
					return fmt.Errorf("unable to deserialize dest proof boc: %w", err)
				}
				stateProof, err := boc.DeserializeBoc(step.LiteServerBlockLinkBack.StateProof)
				if err != nil {
					return fmt.Errorf("unable to deserialize state proof boc: %w", err)
				}
				proof, err := boc.DeserializeBoc(step.LiteServerBlockLinkBack.Proof)
				if err != nil {
					return fmt.Errorf("unable to deserialize proof boc: %w", err)
				}
				err = VerifyBackwardProofLink(toKeyBlock, sourceBlock, targetBlock, destProof[0], stateProof[0], proof[0])
				if err != nil {
					return fmt.Errorf("failed to verify backward proof: %w", err)
				}
			case "LiteServerBlockLinkForward":
				toKeyBlock := step.LiteServerBlockLinkForward.ToKeyBlock
				sourceBlock := step.LiteServerBlockLinkForward.From.ToBlockIdExt()
				targetBlock := step.LiteServerBlockLinkForward.To.ToBlockIdExt()
				destProof, err := boc.DeserializeBoc(step.LiteServerBlockLinkForward.DestProof)
				if err != nil {
					return fmt.Errorf("unable to deserialize dest proof boc: %w", err)
				}
				configProof, err := boc.DeserializeBoc(step.LiteServerBlockLinkForward.ConfigProof)
				if err != nil {
					return fmt.Errorf("unable to deserialize state proof boc: %w", err)
				}
				signatures := signaturesMapper(step.LiteServerBlockLinkForward.Signatures)

				err = VerifyForwardProofLink(toKeyBlock, sourceBlock, targetBlock, destProof[0], configProof[0], *signatures)
				if err != nil {
					return fmt.Errorf("failed to verify forward proof: %w", err)
				}
				// Go to next step
				sourceBlock = step.LiteServerBlockLinkForward.To.ToBlockIdExt()
			}
		}
		source = partialBlockProof.To.ToBlockIdExt()
	}

	return nil
}

func VerifyBackwardProofLink(toKeyBlock bool, source, target ton.BlockIDExt, destProof, stateProof, proof *boc.Cell) error {
	if source.Workchain != -1 && target.Workchain != -1 {
		return fmt.Errorf("both blocks must be from the masterchain")
	}

	if source.Seqno <= target.Seqno {
		return fmt.Errorf("source seqno must be > target seqno for backward link")
	}

	proofsBoc, err := boc.SerializeMulitpleRootsBoc([]*boc.Cell{proof, stateProof}, false, false, false, 0)
	if err != nil {
		return fmt.Errorf("unable to serialize proof: %w", err)
	}
	stateExtra, err := checkBlockShardStateProof(proofsBoc, source.RootHash, nil)
	if err != nil {
		return fmt.Errorf("failed to check proof for shard in the masterchain: %w", err)
	}
	if !stateExtra.ShardStateUnsplit.Custom.Exists {
		return fmt.Errorf("source block is not a masterchain block")
	}
	keyBlock, ok := stateExtra.ShardStateUnsplit.Custom.Value.Value.Other.PrevBlocks.GetWithoutExtra(tlb.Uint32(target.Seqno))
	if !ok {
		return fmt.Errorf("target block not found in state proof")
	}
	if keyBlock.Key != toKeyBlock {
		return fmt.Errorf("unexpected target block in proof isKey value: got = %v, want %v", keyBlock.Key, toKeyBlock)
	}
	if keyBlock.BlkRef.RootHash.Equal(target.RootHash) {
		return fmt.Errorf("incorrect target block hash in proof")
	}

	// proof target block
	_, err = checkBlockProof(*destProof, target.RootHash)
	if err != nil {
		return fmt.Errorf("failed to check target block proof: %w", err)
	}
	var targetBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(destProof, &targetBlock)
	if err != nil {
		return fmt.Errorf("failed to unmashal target block from proof: %w", err)
	}

	if targetBlock.VirtualRoot.Info.KeyBlock != toKeyBlock {
		return fmt.Errorf("unexpected target block isKey value: got = %v, want %v", targetBlock.VirtualRoot.Info.KeyBlock, toKeyBlock)
	}

	return nil
}

func VerifyForwardProofLink(toKeyBlock bool, source, target ton.BlockIDExt, destProof, configProof *boc.Cell, signatures signatures) error {
	if source.Workchain != -1 && target.Workchain != -1 {
		return fmt.Errorf("both blocks must be source the masterchain")
	}

	if source.Seqno >= target.Seqno {
		return fmt.Errorf("source seqno must be < target seqno for forward link")
	}

	// proof source block
	_, err := checkBlockProof(*configProof, source.RootHash)
	if err != nil {
		return fmt.Errorf("failed to check source block proof: %w", err)
	}

	var sourceBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(configProof, &sourceBlock)
	if err != nil {
		return fmt.Errorf("failed to unmashal source block from proof: %w", err)
	}

	if !sourceBlock.VirtualRoot.Extra.Custom.Exists {
		return fmt.Errorf("source block is lack of extra info")
	}

	blockValidatorsCell, ok := sourceBlock.VirtualRoot.Extra.Custom.Value.Value.Config.Config.Get(tlb.Uint32(34))
	if !ok {
		return fmt.Errorf("source block is lack of 34 config param")
	}
	var blockValidators = &tlb.ValidatorsSet{}
	if err := tlb.Unmarshal(&blockValidatorsCell.Value, blockValidators); err != nil {
		return fmt.Errorf("unbale to unmashal validators set: %w", err)
	}

	catchainConfigCell, ok := sourceBlock.VirtualRoot.Extra.Custom.Value.Value.Config.Config.Get(tlb.Uint32(28))
	if !ok {
		return fmt.Errorf("source block is lack of catchain config cell")
	}
	var catchainConfig = &tlb.CatchainConfig{}
	if err := tlb.Unmarshal(&catchainConfigCell.Value, catchainConfig); err != nil {
		return fmt.Errorf("unable to unmashal catchain config: %w", err)
	}

	// proof target block
	_, err = checkBlockProof(*destProof, target.RootHash)
	if err != nil {
		return fmt.Errorf("failed to check target block proof: %w", err)
	}

	var targetBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(destProof, &targetBlock)
	if err != nil {
		return fmt.Errorf("failed to unmarshal target block from proof: %w", err)
	}

	if targetBlock.VirtualRoot.Info.KeyBlock != toKeyBlock {
		return fmt.Errorf("unexpected target block in proof isKey value: got = %v, want %v", targetBlock.VirtualRoot.Info.KeyBlock, toKeyBlock)
	}

	if targetBlock.VirtualRoot.Info.GenValidatorListHashShort != signatures.validatorListHashShort {
		return fmt.Errorf("incorrect validator list hash")
	}

	if targetBlock.VirtualRoot.Info.GenCatchainSeqno != signatures.catchainSeqno {
		return fmt.Errorf("incorrect catchain seqno")
	}
	validators, err := GetMainValidators(&target, *catchainConfig, *blockValidators, targetBlock.VirtualRoot.Info.GenCatchainSeqno)
	if err != nil {
		return fmt.Errorf("failed to get main validators: %w", err)
	}

	if err = CheckBlockSignatures(&target, signatures, validators); err != nil {
		return fmt.Errorf("failed to check block signatures: %w", err)
	}

	return nil
}

func GetMainValidators(block *ton.BlockIDExt, catchainConfig tlb.CatchainConfig, validatorsSet tlb.ValidatorsSet, catchainSeqno uint32) ([]*tlb.ValidatorAddr, error) {
	if block.Workchain != -1 {
		return nil, fmt.Errorf("block must be from the masterchain")
	}

	isShuffle := false
	var validatorsNum int

	// set isShuffle only for new catchain config, if its old, then isShuffle is always false
	switch catchainConfig.SumType {
	case "CatchainConfigNew":
		isShuffle = catchainConfig.CatchainConfigNew.ShuffleMcValidators
	}

	type validatorWithKey struct {
		descr tlb.ValidatorAddr
		key   uint16
	}
	var validatorsKeys []validatorWithKey
	switch validatorsSet.SumType {
	case "Validators":
		validatorsNum = int(validatorsSet.Validators.Main)
		validatorsKeys = make([]validatorWithKey, len(validatorsSet.Validators.List.Items()))
		for i, item := range validatorsSet.Validators.List.Items() {
			validatorsKeys[i].key = uint16(item.Key.FixedSize())
			validatorsKeys[i].descr = item.Value.ValidatorAddr
		}
	case "ValidatorsExt":
		validatorsNum = int(validatorsSet.ValidatorsExt.Main)
		totalWeight := uint64(0)
		validatorsKeys = make([]validatorWithKey, len(validatorsSet.ValidatorsExt.List.Items()))
		for i, item := range validatorsSet.ValidatorsExt.List.Items() {
			totalWeight += item.Value.ValidatorAddr.Weight
			validatorsKeys[i].key = uint16(item.Key.FixedSize())
			validatorsKeys[i].descr = item.Value.ValidatorAddr
		}

		if totalWeight != validatorsSet.ValidatorsExt.TotalWeight {
			return nil, fmt.Errorf("incorrect sum of validators weights")
		}
	default:
		return nil, fmt.Errorf("unknown validators set sumtype")
	}

	if len(validatorsKeys) == 0 {
		return nil, fmt.Errorf("zero validators found")
	}
	if validatorsNum > len(validatorsKeys) {
		validatorsNum = len(validatorsKeys)
	}

	sort.Slice(validatorsKeys, func(i, j int) bool { return validatorsKeys[i].key < validatorsKeys[j].key })

	var validators = make([]*tlb.ValidatorAddr, validatorsNum)
	if isShuffle {
		prng, err := ton.NewValidatorPRNG(nil, block.Shard, block.Workchain, catchainSeqno)
		if err != nil {
			return nil, fmt.Errorf("unable to create validator prng: %w", err)
		}

		idx := make([]uint32, validatorsNum)
		for i := 0; i < validatorsNum; i++ {
			j := prng.NextRanged(uint64(i) + 1)
			idx[i] = idx[j]
			idx[j] = uint32(i)
		}

		for i := 0; i < validatorsNum; i++ {
			validators[i] = &validatorsKeys[idx[i]].descr
		}

		return validators, nil
	}

	for i := 0; i < validatorsNum; i++ {
		validators[i] = &validatorsKeys[i].descr
	}

	return validators, nil
}

func CheckBlockSignatures(block *ton.BlockIDExt, signatures signatures, validators []*tlb.ValidatorAddr) error {
	if len(validators) == 0 {
		return fmt.Errorf("zero validators found")
	}

	if len(signatures.signatures) == 0 {
		return fmt.Errorf("zero signatures found")
	}

	validatorSetHash, err := computeValidatorSetHash(signatures.catchainSeqno, validators)
	if err != nil {
		return fmt.Errorf("unable to compute validator set hash: %w", err)
	}

	if validatorSetHash != signatures.validatorListHashShort {
		return fmt.Errorf("invalid validator set hash")
	}

	totalWeight := uint64(0)
	keyToValidator := map[[32]byte]*tlb.ValidatorAddr{}
	magicPrefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicPrefix, magicPublicKey)
	for _, v := range validators {
		pubKey := v.PublicKey.PubKey[:]
		// add some magic prefix for each validator's pub keys
		hashKey := sha256.Sum256(append(magicPrefix, pubKey...))

		totalWeight += v.Weight
		keyToValidator[hashKey] = v
	}

	sort.Slice(signatures.signatures, func(i, j int) bool {
		return bytes.Compare(signatures.signatures[i].nodeIdShort[:], signatures.signatures[j].nodeIdShort[:]) < 0
	})

	type tlBlockId struct {
		RootHash tl.Int256
		FileHash tl.Int256
	}

	magicPrefix = make([]byte, 4)
	binary.LittleEndian.PutUint32(magicPrefix, magicTonBlockID)
	blockBytes, err := tl.Marshal(tlBlockId{RootHash: tl.Int256(block.RootHash), FileHash: tl.Int256(block.FileHash)})
	if err != nil {
		return fmt.Errorf("unable to marshal block id: %w", err)
	}
	blockBytes = append(magicPrefix, blockBytes...) // add magic prefix

	signedWeight := uint64(0)
	for i, sig := range signatures.signatures {
		if i > 0 {
			prevSig := signatures.signatures[i-1]
			if sig.nodeIdShort == prevSig.nodeIdShort {
				return fmt.Errorf("duplicated node signatures found")
			}
		}

		v, ok := keyToValidator[sig.nodeIdShort]
		if !ok {
			return fmt.Errorf("unknown validator signatures: %v", hex.EncodeToString(sig.nodeIdShort[:]))
		}

		var pubKey [32]byte
		pubKey = v.PublicKey.PubKey
		if !ed25519.Verify(pubKey[:], blockBytes, sig.signature) {
			return fmt.Errorf("invalid validator signature: %v", hex.EncodeToString(sig.nodeIdShort[:]))
		}

		signedWeight += v.Weight

		if signedWeight > totalWeight {
			break
		}
	}

	if 3*signedWeight < 2*totalWeight {
		return fmt.Errorf("not enoght signed weights: %v/%v < 2/3", signedWeight, totalWeight)
	}

	return nil
}

func computeValidatorSetHash(catchainSeqno uint32, validators []*tlb.ValidatorAddr) (uint32, error) {
	type tlValidator struct {
		Key    tl.Int256
		Weight uint64
		Addr   tl.Int256
	}

	type tlValidatorSet struct {
		CatchainSeqno uint32
		Validators    []tlValidator
	}

	tlValidators := make([]tlValidator, len(validators))
	for i, currValidator := range validators {
		tlValidators[i].Key = tl.Int256(currValidator.PublicKey.PubKey)
		tlValidators[i].Weight = currValidator.Weight
		tlValidators[i].Addr = tl.Int256(currValidator.AdnlAddr)
	}

	tlValSet := tlValidatorSet{
		CatchainSeqno: catchainSeqno,
		Validators:    tlValidators,
	}

	validatorSetBytes, err := tl.Marshal(tlValSet)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal validator set: %w", err)
	}

	magicPrefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicPrefix, magicValidatorSet)
	if err != nil {
		return 0, fmt.Errorf("unable to decode magic prefix: %w", err)
	}
	validatorSetBytes = append(magicPrefix, validatorSetBytes...)

	return crc32.Checksum(validatorSetBytes, crc32.MakeTable(crc32.Castagnoli)), nil
}

func signaturesMapper(sigs liteclient.LiteServerSignatureSet) *signatures {
	sigSet := make([]signature, 0, len(sigs.Signatures))
	for _, sig := range sigs.Signatures {
		curr := signature{
			nodeIdShort: sig.NodeIdShort,
			signature:   make([]byte, 64),
		}

		copy(curr.signature, sig.Signature)
		sigSet = append(sigSet, curr)
	}

	return &signatures{
		validatorListHashShort: sigs.ValidatorSetHash,
		catchainSeqno:          sigs.CatchainSeqno,
		signatures:             sigSet,
	}
}
