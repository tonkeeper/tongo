package tvm

import "math/big"

type EntryType int

const (
	Int EntryType = iota
	Null
)

type TvmStackEntry struct {
	Type   EntryType
	IntVal big.Int
}

func NewIntStackEntry(val big.Int) TvmStackEntry {
	return TvmStackEntry{
		Type:   Int,
		IntVal: val,
	}
}

func NewNullStackEntry() TvmStackEntry {
	return TvmStackEntry{
		Type: Null,
	}
}
