package boc

import "fmt"

type prefixMismatchErr struct {
	bitLen   int
	expected uint64
	actual   uint64
}

func fmtPrefix(bits int, val uint64) string {
	if bits%8 == 0 {
		return fmt.Sprintf("0x%0*X", bits/8, val)
	}
	return fmt.Sprintf("%0*b", bits, val)
}

func (e prefixMismatchErr) Error() string {
	return fmt.Sprintf("unexpected prefix %s (exp. %s)", fmtPrefix(64, e.actual), fmtPrefix(64, e.expected))
}
