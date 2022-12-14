package liteclient

import "fmt"

func (t LiteServerError) Error() string {
	return fmt.Sprintf("error code: %d message: %s", t.Code, t.Message)
}
