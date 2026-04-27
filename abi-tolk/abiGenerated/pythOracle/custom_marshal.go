package abiPythOracle

import (
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

const maxBitsPerCell = 1016

func (r *WormholeProofBytes) unmarshalWithTail(c *boc.Cell) (*boc.Cell, error) {
	bytesToRead, err := c.ReadUint(16)
	if err != nil {
		return nil, fmt.Errorf("wormhole proof size: %v", err)
	}
	bitsToRead := int(bytesToRead) * 8
	bits, tail, err := c.ReadBitsDeep(bitsToRead)
	if err != nil {
		return nil, fmt.Errorf("wormhole proof bytes: %v", err)
	}
	res := make([]tlb.Uint8, int(bytesToRead))
	for i, b := range bits.Buffer() {
		res[i] = tlb.Uint8(b)
	}
	*r = res
	return tail, nil
}

func (r *WormholeProofBytes) UnmarshalTLB(c *boc.Cell, _ *tlb.Decoder) error {
	_, err := r.unmarshalWithTail(c)
	return err
}

func (r *WormholeProofBytes) marshalWithTail(c *boc.Cell) (*boc.Cell, error) {
	bytesCount := len(*r)
	if err := c.WriteUint(uint64(bytesCount), 16); err != nil {
		return nil, fmt.Errorf("failed to prepare wormhole proof header: %v", err)
	}
	bs := boc.NewBitString(8 * bytesCount)
	for _, b := range *r {
		if err := bs.WriteByte(byte(b)); err != nil {
			return nil, fmt.Errorf("failed to prepare wormhole proof bytes: %v", err)
		}
	}
	tail, err := c.WriteBitStringTail(bs, maxBitsPerCell)
	if err != nil {
		return nil, fmt.Errorf("failed to write wormhole proof bytes: %v", err)
	}
	return tail, nil
}

func (r *WormholeProofBytes) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	_, err := r.marshalWithTail(c)
	return err
}

func (r *PriceFeedProof) unmarshalWithTail(c *boc.Cell) (*boc.Cell, error) {
	count, err := c.ReadUint(8)
	if err != nil {
		return nil, fmt.Errorf("proof count: %v", err)
	}
	bits, tail, err := c.ReadBitsDeep(int(count) * 160)
	if err != nil {
		return nil, fmt.Errorf("proof bytes: %v", err)
	}
	res := make([]tlb.Bits160, count)
	for i := range res {
		b, err := bits.ReadBytes(20)
		if err != nil {
			return nil, fmt.Errorf("proof digest %d: %v", i, err)
		}
		copy(res[i][:], b)
	}
	*r = res
	return tail, nil
}

func (r *PriceFeedProof) UnmarshalTLB(c *boc.Cell, _ *tlb.Decoder) error {
	_, err := r.unmarshalWithTail(c)
	return err
}

func (r *PriceFeedProof) marshalWithTail(c *boc.Cell) (*boc.Cell, error) {
	if err := c.WriteUint(uint64(len(*r)), 8); err != nil {
		return nil, fmt.Errorf("failed to write proof count: %v", err)
	}
	bs := boc.NewBitString(len(*r) * 160)
	for _, digest := range *r {
		if err := bs.WriteBytes(digest[:]); err != nil {
			return nil, fmt.Errorf("failed to prepare proof digest: %v", err)
		}
	}
	tail, err := c.WriteBitStringTail(bs, maxBitsPerCell)
	if err != nil {
		return nil, fmt.Errorf("failed to write proof bytes: %v", err)
	}
	return tail, nil
}

func (r *PriceFeedProof) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	_, err := r.marshalWithTail(c)
	return err
}

func (r *PriceFeedUpdate) unmarshalWithTail(c *boc.Cell, decoder *tlb.Decoder) (*boc.Cell, error) {
	msgSize, err := c.ReadUint(16)
	if err != nil {
		return nil, fmt.Errorf("message size: %v", err)
	}
	msgBits, msgTail, err := c.ReadBitsDeep(int(msgSize) * 8)
	if err != nil {
		return nil, fmt.Errorf("message bytes: %v", err)
	}
	tmpCell := boc.NewCell()
	if err := tmpCell.WriteBitString(msgBits); err != nil {
		return nil, fmt.Errorf("message cell: %v", err)
	}
	if err := r.Message.UnmarshalTLB(tmpCell, decoder); err != nil {
		return nil, fmt.Errorf("message: %v", err)
	}
	proofTail, err := r.Proof.unmarshalWithTail(msgTail)
	if err != nil {
		return nil, fmt.Errorf("proof: %v", err)
	}
	return proofTail, nil
}

func (r *PriceFeedUpdate) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	_, err := r.unmarshalWithTail(c, decoder)
	return err
}

func (r *PriceFeedUpdate) marshalToCell(c *boc.Cell) (*boc.Cell, error) {
	msgCell, err := r.Message.ToCell()
	if err != nil {
		return nil, fmt.Errorf("message cell: %v", err)
	}
	msgBits := msgCell.RawBitString()
	if err := c.WriteUint(uint64(msgBits.GetWriteCursor()/8), 16); err != nil {
		return nil, fmt.Errorf("message size: %v", err)
	}
	msgTail, err := c.WriteBitStringTail(msgBits, maxBitsPerCell)
	if err != nil {
		return nil, fmt.Errorf("message bytes: %v", err)
	}
	proofTail, err := r.Proof.marshalWithTail(msgTail)
	if err != nil {
		return nil, fmt.Errorf("proof: %v", err)
	}
	return proofTail, nil
}

func (r *PriceFeedUpdate) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	_, err := r.marshalToCell(c)
	return err
}

func (r *WormholeUpdateSection) unmarshalWithTail(c *boc.Cell, decoder *tlb.Decoder) (*boc.Cell, error) {
	tail, err := r.Proof.unmarshalWithTail(c)
	if err != nil {
		return nil, fmt.Errorf("wormhole proof: %v", err)
	}
	count, err := tail.ReadUint(8)
	if err != nil {
		return nil, fmt.Errorf("update count: %v", err)
	}
	r.Updates = make([]PriceFeedUpdate, count)
	for i := range r.Updates {
		tail, err = r.Updates[i].unmarshalWithTail(tail, decoder)
		if err != nil {
			return nil, fmt.Errorf("update %d: %v", i, err)
		}
	}
	return tail, nil
}

func (r *WormholeUpdateSection) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	_, err := r.unmarshalWithTail(c, decoder)
	return err
}

func (r *WormholeUpdateSection) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	tail, err := r.Proof.marshalWithTail(c)
	if err != nil {
		return err
	}
	if err := tail.WriteUint(uint64(len(r.Updates)), 8); err != nil {
		return fmt.Errorf("failed to write update count: %v", err)
	}
	for i := range r.Updates {
		tail, err = r.Updates[i].marshalToCell(tail)
		if err != nil {
			return fmt.Errorf("update %d: %v", i, err)
		}
	}
	return nil
}

func (r *AccumulatorUpdatePayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	trailingHeaderSize, err := c.ReadUint(8)
	if err != nil {
		return fmt.Errorf("trailing header size: %v", err)
	}
	bits, tail, err := c.ReadBitsDeep(int(trailingHeaderSize) * 8)
	if err != nil {
		return fmt.Errorf("trailing header bytes: %v", err)
	}
	r.TrailingHeader = make([]tlb.Uint8, int(trailingHeaderSize))
	for i, b := range bits.Buffer() {
		r.TrailingHeader[i] = tlb.Uint8(b)
	}
	if err := r.UpdateType.UnmarshalTLB(tail, decoder); err != nil {
		return fmt.Errorf("update type: %v", err)
	}
	if err := r.WormholeSection.UnmarshalTLB(tail, decoder); err != nil {
		return fmt.Errorf("wormhole section: %v", err)
	}
	return nil
}

func (r *AccumulatorUpdatePayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(r.TrailingHeader) > 255 {
		return fmt.Errorf("trailing header is too large: %d", len(r.TrailingHeader))
	}
	if err := c.WriteUint(uint64(len(r.TrailingHeader)), 8); err != nil {
		return fmt.Errorf("trailing header size: %v", err)
	}
	tail := c
	if len(r.TrailingHeader) > 0 {
		bs := boc.NewBitString(len(r.TrailingHeader) * 8)
		for _, b := range r.TrailingHeader {
			if err := bs.WriteByte(byte(b)); err != nil {
				return fmt.Errorf("prepare trailing header: %v", err)
			}
		}
		var err error
		tail, err = c.WriteBitStringTail(bs, maxBitsPerCell)
		if err != nil {
			return fmt.Errorf("write trailing header: %v", err)
		}
	}
	if err := r.UpdateType.MarshalTLB(tail, encoder); err != nil {
		return fmt.Errorf("update type: %v", err)
	}
	if err := r.WormholeSection.MarshalTLB(tail, encoder); err != nil {
		return fmt.Errorf("wormhole section: %v", err)
	}
	return nil
}

func (r *PriceFeedIdList) unmarshalWithTail(c *boc.Cell) (*boc.Cell, error) {
	idLen, err := c.ReadUint(8)
	if err != nil {
		return nil, fmt.Errorf("price ids length: %v", err)
	}
	bits, tail, err := c.ReadBitsDeep(256 * int(idLen))
	if err != nil {
		return nil, fmt.Errorf("price ids bytes: %v", err)
	}
	result := make([]tlb.Uint256, int(idLen))
	for i := 0; i < int(idLen); i++ {
		idBytes, err := bits.ReadBytes(256 / 8)
		if err != nil {
			return nil, fmt.Errorf("price id %d: %v", i, err)
		}
		var bigInt big.Int
		bigInt.SetBytes(idBytes)
		result[i] = tlb.Uint256(bigInt)
	}
	*r = result
	return tail, nil
}

func (r *PriceFeedIdList) UnmarshalTLB(c *boc.Cell, _ *tlb.Decoder) error {
	_, err := r.unmarshalWithTail(c)
	return err
}

func (r *PriceFeedIdList) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	idLen := len(*r)
	if err := c.WriteUint(uint64(idLen), 8); err != nil {
		return fmt.Errorf("failed to write price feed id list length: %v", err)
	}
	bs := boc.NewBitString(256 * idLen)
	for _, id := range *r {
		v := big.Int(id)
		if err := bs.WriteBigUint(&v, 256); err != nil {
			return fmt.Errorf("failed to prepare price feed id buffer: %v", err)
		}
	}
	if _, err := c.WriteBitStringTail(bs, maxBitsPerCell); err != nil {
		return fmt.Errorf("failed to write price feed ids: %v", err)
	}
	return nil
}
