package liteclient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/tl"
)

func pointer[T any](t T) *T {
	return &t
}

type RequestName = string

const (
	UnknownRequest RequestName = "Unknown"
)

func decodeRequest(tag uint32, name RequestName, msgType any) reqDecoderFunc {
	return func(b []byte) (*RequestName, any, error) {
		if len(b) < 4 {
			return nil, nil, errors.New("message too short")
		}
		readTag := binary.LittleEndian.Uint32(b[:4])
		if readTag != tag {
			return nil, nil, fmt.Errorf("invalid tag")
		}
		body := reflect.New(reflect.TypeOf(msgType))
		reader := bytes.NewReader(b[4:])
		err := tl.Unmarshal(reader, body.Interface())
		if err == nil {
			return pointer(name), body.Elem().Interface(), nil
		}
		return nil, nil, err
	}
}

type reqDecoderFunc func(b []byte) (*RequestName, any, error)

func LiteapiRequestDecoder(b []byte) (uint32, *RequestName, any, error) {
	if len(b) < 4 {
		return 0, nil, nil, errors.New("message too short")
	}
	tag := binary.LittleEndian.Uint32(b[:4])
	f := taggedRequestDecodeFunctions[tag]
	if f != nil {
		o, d, err := f(b)
		if err == nil {
			return tag, o, d, nil
		}
	}
	return tag, pointer(UnknownRequest), nil, nil
}
