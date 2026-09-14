package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

// Tolk serializes `array<T>` as a length followed by a chain of chunks:
//
//	array<T> = len:uint8 head:(Maybe ^Chunk)
//	Chunk = next:(Maybe ^Chunk) items:T*
//
// A chunk holds as many items as fit into a cell, the first ref of a chunk being
// occupied by the next one. The number of items in a chunk is not stored:
// chunks are filled from the tail of the array, so that every chunk but the first
// is packed to capacity, and a reader simply consumes items until a chunk runs out.
//
// For example, `array<MessageToSend>` of the Telegram wallet, whose item is
// 8 bits and one ref, gets 3 items per chunk and 4 in the last one, which has no next ref.
const (
	// MaxTolkArrayLen is the number of items an array can hold: its length is uint8.
	MaxTolkArrayLen = 255

	maxChunkBits = 1023
	maxChunkRefs = 4
)

// StoreTolkArray writes items as a Tolk array<T>.
func StoreTolkArray[T any](c *boc.Cell, encoder *Encoder, items []T) error {
	if len(items) > MaxTolkArrayLen {
		return fmt.Errorf("array of %v items is too long, max %v", len(items), MaxTolkArrayLen)
	}
	if err := c.WriteUint(uint64(len(items)), 8); err != nil {
		return err
	}
	// every item is serialized on its own first, to know how many of them fit into a chunk
	cells := make([]*boc.Cell, len(items))
	for i, item := range items {
		cell := boc.NewCell()
		if err := encode(cell, "", item, encoder); err != nil {
			return fmt.Errorf("failed to marshal item %v: %w", i, err)
		}
		cells[i] = cell
	}
	// chunks are built from the tail, so that every chunk knows its next one
	var next *boc.Cell
	for rest := cells; len(rest) > 0; {
		chunk := boc.NewCell()
		bits, refs := 1, 0 // the maybe bit of the next chunk
		if next != nil {   // reserver ref to refer next chunkl
			refs = 1
		}
		size := 0
		for ; size < len(rest); size++ {
			item := rest[len(rest)-size-1]
			itemBits := maxChunkBits - item.BitsAvailableForWrite()
			if bits+itemBits > maxChunkBits || refs+item.RefsSize() > maxChunkRefs {
				break
			}
			bits, refs = bits+itemBits, refs+item.RefsSize()
		}
		if size == 0 {
			return fmt.Errorf("item %v does not fit into a chunk", len(rest)-1)
		}
		if err := Marshal(chunk, maybeChunkRef(next)); err != nil {
			return err
		}
		for _, item := range rest[len(rest)-size:] {
			if err := chunk.WriteBitString(item.RawBitString()); err != nil {
				return err
			}
			for _, ref := range item.Refs() {
				if err := chunk.AddRef(ref); err != nil {
					return err
				}
			}
		}
		rest, next = rest[:len(rest)-size], chunk
	}
	return Marshal(c, maybeChunkRef(next))
}

// LoadTolkArray reads a Tolk array<T>.
func LoadTolkArray[T any](c *boc.Cell, decoder *Decoder) ([]T, error) {
	length, err := c.ReadUint(8)
	if err != nil {
		return nil, err
	}
	var head Maybe[Ref[boc.Cell]]
	if err := decoder.Unmarshal(c, &head); err != nil {
		return nil, err
	}
	items := make([]T, 0, length)
	for head.Exists {
		chunk := head.Value.Value
		if err := decoder.Unmarshal(&chunk, &head); err != nil {
			return nil, err
		}
		// a chunk holds nothing but items, so read until it runs out
		for !chunk.IsEmpty() {
			if len(items) == int(length) {
				return nil, fmt.Errorf("array declares %v items but holds more", length)
			}
			var item T
			if err := decoder.Unmarshal(&chunk, &item); err != nil {
				return nil, fmt.Errorf("failed to unmarshal item %v: %w", len(items), err)
			}
			items = append(items, item)
		}
	}
	if len(items) != int(length) {
		return nil, fmt.Errorf("array declares %v items but holds %v", length, len(items))
	}
	return items, nil
}

func maybeChunkRef(c *boc.Cell) Maybe[Ref[boc.Cell]] {
	if c == nil {
		return Maybe[Ref[boc.Cell]]{}
	}
	return Maybe[Ref[boc.Cell]]{Exists: true, Value: Ref[boc.Cell]{Value: *c}}
}
