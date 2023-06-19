package abi

import (
	"errors"
)

// ErrStructSizeMismatch means that a message body's cell contains more information than expected.
var ErrStructSizeMismatch = errors.New("struct size is less than available bits and refs in cell")
