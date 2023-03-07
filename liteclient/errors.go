package liteclient

import (
	"fmt"
)

type clientError string

func (e clientError) Error() string {
	return string(e)
}

func newClientError(msg string, args ...interface{}) clientError {
	return clientError(fmt.Sprintf(msg, args...))
}

func IsClientError(err error) bool {
	_, ok := err.(clientError)
	return ok
}

func IsNotConnectedYet(e error) bool {
	return e.Error() == "not connected yet"
}
