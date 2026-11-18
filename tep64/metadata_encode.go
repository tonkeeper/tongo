package tep64

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// ToCell encodes non-empty metadata fields in the TEP-64 on-chain dictionary.
// A non-empty Uri makes the result semi-chain metadata.
func (m Metadata) ToCell() (*boc.Cell, error) {
	var data tlb.HashmapE[tlb.Bits256, tlb.Ref[metadataSnake]]
	for _, field := range []struct {
		key   ton.Bits256
		value []byte
	}{
		{sha256uri, []byte(m.Uri)},
		{sha256name, []byte(m.Name)},
		{sha256description, []byte(m.Description)},
		{sha256image, []byte(m.Image)},
		{sha256imageData, m.ImageData},
		{sha256symbol, []byte(m.Symbol)},
		{sha256decimals, []byte(m.Decimals)},
		{sha256renderType, []byte(m.RenderType)},
		{sha256amountStyle, []byte(m.AmountStyle)},
		{sha256customPayload, []byte(m.CustomPayloadAPIURL)},
	} {
		if len(field.value) != 0 {
			data.Put(tlb.Bits256(field.key), tlb.Ref[metadataSnake]{Value: field.value})
		}
	}

	cell := boc.NewCell()
	if err := cell.WriteUint(0, 8); err != nil { // onchain#00
		return nil, err
	}
	if err := tlb.Marshal(cell, data); err != nil {
		return nil, err
	}
	return cell, nil
}

// TEP-64 snake segments contain whole bytes, including when the first cell
// already holds the snake#00 tag.
type metadataSnake []byte

func (data metadataSnake) MarshalTLB(c *boc.Cell, _ *tlb.Encoder) error {
	if err := c.WriteUint(0, 8); err != nil {
		return err
	}
	for len(data) > 0 {
		n := min(c.BitsAvailableForWrite()/8, len(data))
		if err := c.WriteBytes(data[:n]); err != nil {
			return err
		}
		data = data[n:]
		if len(data) == 0 {
			break
		}
		next, err := c.NewRef()
		if err != nil {
			return err
		}
		c = next
	}
	return nil
}
