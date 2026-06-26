package abi

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
)

type staticLibraryResolver map[ton.Bits256]*boc.Cell

func getStaticLibraryResolver(bocs []string) (staticLibraryResolver, error) {
	resolver := staticLibraryResolver{}
	for _, b64 := range bocs {
		cells, err := boc.DeserializeBocBase64(b64)
		if err != nil {
			return nil, err
		}
		hash, err := cells[0].Hash256()
		if err != nil {
			return nil, err
		}
		resolver[ton.Bits256(hash)] = cells[0]
	}
	return resolver, nil
}

func (r staticLibraryResolver) GetLibraries(ctx context.Context, hashes []ton.Bits256) (map[ton.Bits256]*boc.Cell, error) {
	out := make(map[ton.Bits256]*boc.Cell, len(hashes))
	for _, h := range hashes {
		c, ok := r[h]
		if !ok {
			return nil, fmt.Errorf("library %x is missing from inspectorLibraries", h)
		}
		out[h] = c.CopyCell()
	}
	return out, nil
}
