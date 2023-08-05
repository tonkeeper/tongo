package txemulator

import (
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// FindLibraries looks for library cells inside the given cell tree and
// returns a list of hashes of found library cells.
func FindLibraries(cell *boc.Cell) ([]tongo.Bits256, error) {
	libs, err := findLibraries(cell)
	if err != nil {
		return nil, err
	}
	if len(libs) == 0 {
		return nil, nil
	}
	hashes := make([]tongo.Bits256, 0, len(libs))
	for hash := range libs {
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func findLibraries(cell *boc.Cell) (map[tongo.Bits256]struct{}, error) {
	if cell.IsExotic() {
		if cell.CellType() == boc.LibraryCell {
			bytes, err := cell.ReadBytes(33)
			if err != nil {
				return nil, err
			}
			var hash tongo.Bits256
			copy(hash[:], bytes[1:])
			return map[tongo.Bits256]struct{}{
				hash: {},
			}, nil
		}
		return nil, nil
	}
	var libs map[tongo.Bits256]struct{}
	for _, ref := range cell.Refs() {
		ref.ResetCounters()
		hashes, err := findLibraries(ref)
		if err != nil {
			return nil, err
		}
		if libs == nil {
			libs = hashes
			continue
		}
		for hash := range hashes {
			libs[hash] = struct{}{}
		}
	}
	return libs, nil
}

// LibrariesToBase64 converts a map with libraries to a base64 string.
func LibrariesToBase64(libraries map[tongo.Bits256]*boc.Cell) (string, error) {
	if len(libraries) == 0 {
		return "", nil
	}
	hashes := make([]tlb.Bits256, 0, len(libraries))
	descriptions := make([]tlb.LibDescr, 0, len(libraries))
	for hash, cell := range libraries {
		hashes = append(hashes, tlb.Bits256(hash))
		descriptions = append(descriptions, tlb.LibDescr{Lib: *cell})
	}
	hashmap := tlb.NewHashmap[tlb.Bits256, tlb.LibDescr](hashes, descriptions)
	libsCell := boc.NewCell()
	if err := tlb.Marshal(libsCell, hashmap); err != nil {
		return "", err
	}
	return libsCell.ToBocBase64()
}
