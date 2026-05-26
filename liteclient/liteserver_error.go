package liteclient

import "fmt"

type LiteServerError struct {
	Address string
	Err     error
}

func (e *LiteServerError) Error() string {
	return fmt.Sprintf("liteserver %s: %s", e.Address, e.Err.Error())
}

func (e *LiteServerError) Unwrap() error { return e.Err }
