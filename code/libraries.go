package code

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// FindLibraries looks for library cells inside the given cell tree and
// returns a list of hashes of found library cells.
func FindLibraries(cell *boc.Cell) ([]ton.Bits256, error) {
	libs := make(map[ton.Bits256]struct{})
	visited := make(map[*boc.Cell]struct{})
	if err := findLibraries(cell, libs, visited); err != nil {
		return nil, err
	}
	if len(libs) == 0 {
		return nil, nil
	}
	hashes := make([]ton.Bits256, 0, len(libs))
	for hash := range libs {
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func findLibraries(cell *boc.Cell, libs map[ton.Bits256]struct{}, visited map[*boc.Cell]struct{}) error {
	// "visited" tracks cells already walked by pointer identity.
	// Without this set, findLibraries explores every root->leaf
	// path independently, which is exponential on real code cells.
	if _, ok := visited[cell]; ok {
		return nil
	}
	visited[cell] = struct{}{}
	cell.ShallowResetCounters()
	if cell.IsExotic() {
		if cell.CellType() == boc.LibraryCell {
			bytes, err := cell.ReadBytes(33)
			if err != nil {
				return err
			}
			var hash ton.Bits256
			copy(hash[:], bytes[1:])
			libs[hash] = struct{}{}
		}
		return nil
	}
	for _, ref := range cell.Refs() {
		if err := findLibraries(ref, libs, visited); err != nil {
			return err
		}
	}
	return nil
}

// LibrariesToBase64 converts a map with libraries to a base64 string.
func LibrariesToBase64(libraries map[ton.Bits256]*boc.Cell) (string, error) {
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
