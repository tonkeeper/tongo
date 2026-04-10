package abiGenerated

import (
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (r *SkippedBytes) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	skipped, err := c.ReadUint(8)
	if err != nil {
		return fmt.Errorf("skip length: %v", err)
	}
	*r = SkippedBytes(skipped)
	if err = c.Skip(int(skipped) * 8); err != nil {
		return fmt.Errorf("failed to skip %v bytes: %v", skipped, err)
	}
	return nil
}

func (r *WormholeProofBytes) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	bytesToRead, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	bitsToRead := int(bytesToRead) * 8
	bits, lastCell, err := c.ReadBitsDeep(bitsToRead)
	if err != nil {
		return err
	}
	res := make([]tlb.Uint8, int(bytesToRead))
	for i, b := range bits.Buffer() {
		res[i] = tlb.Uint8(b)
	}
	*r = res
	*c = *lastCell
	return nil
}

func (r *PriceFeedIdList) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	idLen, err := c.ReadUint(8)
	if err != nil {
		return err
	}
	bits, lastCell, err := c.ReadBitsDeep(256 * int(idLen))
	if err != nil {
		return err
	}
	result := make([]tlb.Uint256, int(idLen))
	for i := 0; i < int(idLen); i++ {
		idBytes, err := bits.ReadBytes(256 / 8)
		if err != nil {
			return err
		}
		var bigInt big.Int
		bigInt.SetBytes(idBytes)
		result[i] = tlb.Uint256(bigInt)
	}
	*r = result
	*c = *lastCell
	return nil
}
