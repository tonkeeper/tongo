package tep64

import (
	"bytes"
	"crypto/sha256"
	"reflect"
	"strconv"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func TestMetadataToCell(t *testing.T) {
	for _, tt := range []struct {
		name   string
		meta   Metadata
		layout ContentLayout
	}{
		{"empty", Metadata{}, OnChain},
		{"semi-chain", Metadata{Uri: "https://example.com/token.json", Symbol: "TOK"}, SemiChain},
		{"on-chain", Metadata{
			Name: "Монета", Description: "Token description", Symbol: "TOK", Decimals: "0",
			Image: "ipfs://example", ImageData: []byte{0xff, 0, 0x81},
			RenderType: "currency", AmountStyle: "n", CustomPayloadAPIURL: "https://example.com/api",
		}, OnChain},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := tt.meta.ToCell()
			if err != nil {
				t.Fatal(err)
			}
			got, err := DecodeFullContentFromCell(cell)
			if err != nil {
				t.Fatal(err)
			}
			if got.Layout != tt.layout || !reflect.DeepEqual(got.OnchainMetadata, &tt.meta) {
				t.Errorf("got layout %v and metadata %+v, want %v and %+v", got.Layout, got.OnchainMetadata, tt.layout, tt.meta)
			}
		})
	}
}

func TestMetadataToCellSnakeBoundaries(t *testing.T) {
	// The first cell can hold 126 bytes after the tag. Subsequent cells hold 127.
	for _, size := range []int{126, 127, 253, 254} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			payload := make([]byte, size)
			for i := range payload {
				payload[i] = byte(i)
			}
			cell, err := (Metadata{ImageData: payload}).ToCell()
			if err != nil {
				t.Fatal(err)
			}
			prefix, err := cell.ReadUint(8)
			if err != nil || prefix != 0 {
				t.Fatalf("on-chain prefix: got %d, %v", prefix, err)
			}
			var data tlb.HashmapE[tlb.Bits256, tlb.Ref[tlb.Any]]
			if err := tlb.Unmarshal(cell, &data); err != nil {
				t.Fatal(err)
			}
			ref, ok := data.Get(sha256.Sum256([]byte("image_data")))
			if !ok || len(data.Keys()) != 1 {
				t.Fatalf("expected only the image_data dictionary entry, got keys %v", data.Keys())
			}
			value := boc.Cell(ref.Value)
			prefix, err = value.ReadUint(8)
			if err != nil || prefix != 0 {
				t.Fatalf("snake prefix: got %d, %v", prefix, err)
			}
			got, err := value.ReadBytesTail()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, payload) {
				t.Errorf("image_data did not preserve the %d input bytes", size)
			}
		})
	}
}
