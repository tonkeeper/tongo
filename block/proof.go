package block

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"hash/crc32"
	"sort"
)

type signature struct {
	nodeIdShort []byte
	signature   []byte
}
type signatures struct {
	validatorListHashShort uint32
	catchainSeqno          uint32
	signatures             []signature
}

// todo work on errors
// todo maybe create some struct to store client in it?
func VerifyProofChain(ctx context.Context, c *liteclient.Client, source, target *ton.BlockIDExt) error {
	isForward := source.Seqno < target.Seqno

	for source.Seqno != target.Seqno {
		sourceTL := liteclient.TonNodeBlockIdExtC{
			Workchain: uint32(source.Workchain),
			Shard:     source.Shard,
			Seqno:     source.Seqno,
			RootHash:  tl.Int256(source.RootHash),
			FileHash:  tl.Int256(source.FileHash),
		}
		partialBlockProof, err := c.LiteServerGetBlockProof(ctx, liteclient.LiteServerGetBlockProofRequest{
			Mode:       1,
			KnownBlock: sourceTL,
			TargetBlock: &liteclient.TonNodeBlockIdExtC{
				Workchain: uint32(target.Workchain),
				Shard:     target.Shard,
				Seqno:     target.Seqno,
				RootHash:  tl.Int256(target.RootHash),
				FileHash:  tl.Int256(target.FileHash),
			},
		})
		if err != nil {
			return fmt.Errorf("cannot get partial block proof source liteserver: %w", err)
		}
		// todo is ok to compare structures?
		if partialBlockProof.From != sourceTL {
			return fmt.Errorf("incorrect source partial block, got %v", partialBlockProof.From.Seqno)
		}

		for _, step := range partialBlockProof.Steps {
			switch step.SumType {
			case "LiteServerBlockLinkBack":
				if isForward {
					return fmt.Errorf("wrong direction in backward sequence")
				}
				toKeyBlock := step.LiteServerBlockLinkBack.ToKeyBlock
				sourceBlock := blockIdExtMapper(step.LiteServerBlockLinkBack.From)
				targetBlock := blockIdExtMapper(step.LiteServerBlockLinkBack.To)
				cells, err := boc.DeserializeBoc(step.LiteServerBlockLinkBack.DestProof)
				if err != nil {
					return fmt.Errorf("cannot deserialize dest proof boc: %w", err)
				}
				destProof := cells[0]
				cells, err = boc.DeserializeBoc(step.LiteServerBlockLinkBack.StateProof)
				if err != nil {
					return fmt.Errorf("cannot deserialize state proof boc: %w", err)
				}
				stateProof := cells[0]
				cells, err = boc.DeserializeBoc(step.LiteServerBlockLinkBack.Proof)
				if err != nil {
					return fmt.Errorf("cannot deserialize proof boc: %w", err)
				}
				proof := cells[0]

				if err := VerifyBackwardProofLink(toKeyBlock, sourceBlock, targetBlock, destProof, stateProof, proof); err != nil {
					return fmt.Errorf("failed to verify backward proof: %w", err)
				}

			case "LiteServerBlockLinkForward":
				toKeyBlock := step.LiteServerBlockLinkForward.ToKeyBlock
				sourceBlock := blockIdExtMapper(step.LiteServerBlockLinkForward.From)
				targetBlock := blockIdExtMapper(step.LiteServerBlockLinkForward.To)
				cells, err := boc.DeserializeBoc(step.LiteServerBlockLinkForward.DestProof)
				if err != nil {
					return fmt.Errorf("cannot deserialize dest proof boc: %w", err)
				}
				destProof := cells[0]
				cells, err = boc.DeserializeBoc(step.LiteServerBlockLinkForward.ConfigProof)
				if err != nil {
					return fmt.Errorf("cannot deserialize state proof boc: %w", err)
				}
				configProof := cells[0]
				signatures := signaturesMapper(step.LiteServerBlockLinkForward.Signatures)

				if err := VerifyForwardProofLink(toKeyBlock, *sourceBlock, *targetBlock, destProof, configProof, *signatures); err != nil {
					return fmt.Errorf("failed to verify forward proof: %w", err)
				}

				sourceBlock = blockIdExtMapper(step.LiteServerBlockLinkForward.To)
			}
		}
		source = blockIdExtMapper(partialBlockProof.To)
	}

	return nil
}

func VerifyBackwardProofLink(toKeyBlock bool, source, target *ton.BlockIDExt, destProof, stateProof, proof *boc.Cell) error {
	if source.Workchain != -1 && target.Workchain != -1 {
		return fmt.Errorf("both blocks must be source the masterchain")
	}

	if source.Seqno <= target.Seqno {
		return fmt.Errorf("source seqno must be > target seqno")
	}

	// todo have we already have serialization for multiple roots?
	prf, err := boc.SerializeBocWithMulitplRoots([]*boc.Cell{proof, stateProof}, false, false, false, 0)
	stateExtra, err := liteapi.CheckBlockShardStateProof(prf, source.RootHash, nil)
	if err != nil {
		return fmt.Errorf("cannot check proof for shard in the masterchain: %w", err)
	}
	if !stateExtra.ShardStateUnsplit.Custom.Exists {
		return fmt.Errorf("not a masterchain block")
	}
	keyBlock, ok := stateExtra.ShardStateUnsplit.Custom.Value.Value.Other.PrevBlocks.GetWithoutExtra(tlb.Uint32(target.Seqno))
	if !ok {
		return fmt.Errorf("target block not found in state proof")
	}
	if keyBlock.Key != toKeyBlock {
		// todo set more clear error value
		return fmt.Errorf("expected target block is Key block = %v, but it's %v", toKeyBlock, keyBlock.Key)
	}
	if keyBlock.BlkRef.RootHash.Equal(target.RootHash) {
		return fmt.Errorf("incorrect target block hash in proof")
	}

	// proof target block
	// https://github.com/tact-lang/lite-proof/blob/9a6f3efa938ff67d86b7c4c74e355984da962f2b/src/proof/verify.ts#L68

	_, err = liteapi.CheckBlockProof(*destProof, target.RootHash)
	if err != nil {
		return fmt.Errorf("target block proof failed: %w", err)
	}

	var targetBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(destProof, &targetBlock)
	if err != nil {
		return fmt.Errorf("failed to get target block source proof: %w", err)
	}

	if targetBlock.VirtualRoot.Info.KeyBlock != toKeyBlock {
		// todo set more clear error value
		return fmt.Errorf("expected target block is Key block = %v, but it's %v", toKeyBlock, targetBlock.VirtualRoot.Info.KeyBlock)
	}

	return nil
}

func VerifyForwardProofLink(toKeyBlock bool, source, target ton.BlockIDExt, destProof, configProof *boc.Cell, signatures signatures) error {
	if source.Workchain != -1 && target.Workchain != -1 {
		return fmt.Errorf("both blocks must be source the masterchain")
	}

	if source.Seqno >= target.Seqno {
		return fmt.Errorf("source seqno must be < target seqno")
	}

	// proof source block (source)
	_, err := liteapi.CheckBlockProof(*configProof, source.RootHash)
	if err != nil {
		return fmt.Errorf("source block proof failed: %w", err)
	}

	var sourceBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(configProof, &sourceBlock)
	if err != nil {
		return fmt.Errorf("failed to get source block source proof: %w", err)
	}

	if !sourceBlock.VirtualRoot.Extra.Custom.Exists {
		return fmt.Errorf("source block proof missing extra info")
	}

	blockValidatorsCell, ok := sourceBlock.VirtualRoot.Extra.Custom.Value.Value.Config.Config.Get(tlb.Uint32(34))
	if !ok {
		return fmt.Errorf("source block proof missing block Validators cell")
	}
	var blockValidators = &tlb.ValidatorsSet{}
	if err := tlb.Unmarshal(&blockValidatorsCell.Value, blockValidators); err != nil {
		return fmt.Errorf("cannot get Validators set: %w", err)
	}

	catchainConfigCell, ok := sourceBlock.VirtualRoot.Extra.Custom.Value.Value.Config.Config.Get(tlb.Uint32(28))
	if !ok {
		return fmt.Errorf("source block proof missing catchain config cell")
	}
	var catchainConfig = &tlb.CatchainConfig{}
	if err := tlb.Unmarshal(&catchainConfigCell.Value, catchainConfig); err != nil {
		return fmt.Errorf("cannot get catchain config: %w", err)
	}

	// proof target block
	_, err = liteapi.CheckBlockProof(*destProof, target.RootHash)
	if err != nil {
		return fmt.Errorf("target block proof failed: %w", err)
	}

	var targetBlock tlb.MerkleProof[tlb.Block]
	err = tlb.Unmarshal(destProof, &targetBlock)
	if err != nil {
		return fmt.Errorf("failed to get target block source proof: %w", err)
	}

	if targetBlock.VirtualRoot.Info.KeyBlock != toKeyBlock {
		// todo set more clear error value
		return fmt.Errorf("expected target block is Key block = %v, but it's %v", toKeyBlock, targetBlock.VirtualRoot.Info.KeyBlock)
	}

	if targetBlock.VirtualRoot.Info.GenValidatorListHashShort != signatures.validatorListHashShort {
		return fmt.Errorf("incorrect validator list hash")
	}

	if targetBlock.VirtualRoot.Info.GenCatchainSeqno != signatures.catchainSeqno {
		return fmt.Errorf("incorrect catchain seqno")
	}
	validators, err := GetMainValidators(&target, *catchainConfig, *blockValidators, targetBlock.VirtualRoot.Info.GenCatchainSeqno)
	if err != nil {
		return fmt.Errorf("failed to get main Validators: %w", err)
	}

	if err = CheckBlockSignatures(&target, signatures, validators); err != nil {
		return fmt.Errorf("block signatures verification failed: %w", err)
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
			return nil, fmt.Errorf("incorrect sum of Validators weights")
		}
	default:
		return nil, fmt.Errorf("unknown Validators set sumType")
	}

	if len(validatorsKeys) == 0 {
		return nil, fmt.Errorf("Validators set has no Validators")
	}
	if validatorsNum > len(validatorsKeys) {
		validatorsNum = len(validatorsKeys)
	}

	sort.Slice(validatorsKeys, func(i, j int) bool { return validatorsKeys[i].key < validatorsKeys[j].key })

	var validators = make([]*tlb.ValidatorAddr, validatorsNum)
	if isShuffle {
		prng, err := NewValidatorPRNG(nil, block.Shard, block.Workchain, catchainSeqno)
		if err != nil {
			return nil, fmt.Errorf("cannot create validator prng: %w", err)
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
		return fmt.Errorf("zero Validators")
	}

	if len(signatures.signatures) == 0 {
		return fmt.Errorf("zero signatures")
	}

	validatorSetHash, err := computeValidatorSetHash(signatures.catchainSeqno, validators)
	if err != nil {
		return fmt.Errorf("cannot compute validator set hash: %w", err)
	}

	if validatorSetHash != signatures.validatorListHashShort {
		return fmt.Errorf("incorrect validator set hash")
	}

	totalWeight := uint64(0)
	keyToValidator := map[string]*tlb.ValidatorAddr{}
	magicPrefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicPrefix, crc32.Checksum([]byte("pub.ed25519 key:int256 = PublicKey"), crc32.MakeTable(crc32.IEEE))) // todo why?
	for _, v := range validators {
		pubKey := v.PublicKey.PubKey[:]
		hashKey := sha256.Sum256(append(magicPrefix, pubKey...))

		totalWeight += v.Weight
		keyToValidator[string(hashKey[:])] = v
	}

	sort.Slice(signatures.signatures, func(i, j int) bool {
		return bytes.Compare(signatures.signatures[i].nodeIdShort, signatures.signatures[j].nodeIdShort) < 0
	})

	type tlBlockId struct {
		RootHash tl.Int256
		FileHash tl.Int256
	}
	magicPrefix = make([]byte, 4)
	binary.LittleEndian.PutUint32(magicPrefix, crc32.Checksum([]byte("ton.blockId root_cell_hash:int256 file_hash:int256 = ton.BlockId"), crc32.MakeTable(crc32.IEEE)))
	blockBytes, err := tl.Marshal(tlBlockId{RootHash: tl.Int256(block.RootHash), FileHash: tl.Int256(block.FileHash)})
	if err != nil {
		return fmt.Errorf("cannot marshal block id: %w", err)
	}
	blockBytes = append(magicPrefix, blockBytes...) // add magic prefix

	signedWeight := uint64(0)
	for i, sig := range signatures.signatures {
		if i > 0 {
			prevSig := signatures.signatures[i-1]
			if bytes.Equal(sig.nodeIdShort, prevSig.nodeIdShort) {
				return fmt.Errorf("duplicated node signatures")
			}
		}

		// todo maybe dont convert to string? looks weird
		v, ok := keyToValidator[string(sig.nodeIdShort)]
		if !ok {
			return fmt.Errorf("unknown validator signatures: %v", hex.EncodeToString(sig.nodeIdShort))
		}

		var pubKey [32]byte
		pubKey = v.PublicKey.PubKey
		if !ed25519.Verify(pubKey[:], blockBytes, sig.signature) {
			return fmt.Errorf("invalid signatures of validator: %v", hex.EncodeToString(sig.nodeIdShort))
		}

		signedWeight += v.Weight

		if signedWeight > totalWeight {
			break
		}
	}

	if 3*signedWeight < 2*totalWeight {
		return fmt.Errorf("not enoght signed Weight: %v/%v < 2/3", signedWeight, totalWeight)
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
		return 0, fmt.Errorf("cannot marshal validator set: %w", err)
	}

	const magicPrefix = "ed601690" // reversed 901660ed
	decodedMagicPrefix, err := hex.DecodeString(magicPrefix)
	if err != nil {
		return 0, fmt.Errorf("cannot decode magic prefix: %w", err)
	}
	validatorSetBytes = append(decodedMagicPrefix, validatorSetBytes...)

	return crc32.Checksum(validatorSetBytes, crc32.MakeTable(crc32.Castagnoli)), nil
}

func blockIdExtMapper(blockId liteclient.TonNodeBlockIdExtC) *ton.BlockIDExt {
	return &ton.BlockIDExt{
		BlockID: ton.BlockID{
			Workchain: int32(blockId.Workchain),
			Shard:     blockId.Shard,
			Seqno:     blockId.Seqno,
		},
		RootHash: ton.Bits256(blockId.RootHash),
		FileHash: ton.Bits256(blockId.FileHash),
	}
}

func signaturesMapper(sigs liteclient.LiteServerSignatureSet) *signatures {
	sigSet := make([]signature, 0, len(sigs.Signatures))
	for _, sig := range sigs.Signatures {
		curr := signature{
			nodeIdShort: make([]byte, 32),
			signature:   make([]byte, 64),
		}

		copy(curr.signature, sig.Signature)
		copy(curr.nodeIdShort, sig.NodeIdShort[:])
		sigSet = append(sigSet, curr)
	}

	return &signatures{
		validatorListHashShort: sigs.ValidatorSetHash,
		catchainSeqno:          sigs.CatchainSeqno,
		signatures:             sigSet,
	}
}
