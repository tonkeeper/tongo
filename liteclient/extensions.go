package liteclient

//go:generate go run generator.go

import "fmt"

func (t LiteServerError) Error() string {
	return fmt.Sprintf("error code: %d message: %s", t.Code, t.Message)
}
