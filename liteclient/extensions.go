package liteclient

//go:generate go run generator.go

import "fmt"

var (
	ErrBlockNotApplied = fmt.Errorf("block is not applied")
)

func (t LiteServerError) Error() string {
	return fmt.Sprintf("error code: %d message: %s", t.Code, t.Message)
}

func (t LiteServerError) IsNotApplied() bool {
	return t.Message == "block is not applied"
}
