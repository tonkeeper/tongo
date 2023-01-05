package tl

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
)

type MarshalerTL interface {
	MarshalTL() ([]byte, error)
}

func Marshal(o any) ([]byte, error) {
	if m, ok := o.(MarshalerTL); ok {
		return m.MarshalTL()
	}
	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Pointer {
		// TODO: check for nil pointer
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(val.Uint()))
		return b, nil
	case reflect.Int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(int32(val.Int())))
		return b, nil
	case reflect.Int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(val.Int()))
		return b, nil
	case reflect.Uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, val.Uint())
		return b, nil
	case reflect.Bool:
		b := make([]byte, 4)
		if val.Bool() {
			binary.BigEndian.PutUint32(b, 0xb5757299) // true
		} else {
			binary.BigEndian.PutUint32(b, 0x379779bc) // false
		}
		return b, nil
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return encodeVector(val)
		}
		data := val.Bytes()
		b := append(EncodeLength(len(data)), data...)
		return zeroPadding(b), nil
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return nil, fmt.Errorf("encoding array of %v not supported", val.Type().Elem().Kind())
		}
		b := make([]byte, 0, val.Len()+8)
		b = append(b, EncodeLength(val.Len())...)
		for i := 0; i < val.Len(); i++ {
			b = append(b, uint8(val.Index(i).Uint()))
		}
		return zeroPadding(b), nil
	case reflect.Struct:
		return encodeStruct(val)
	default:
		return nil, fmt.Errorf("type %v not implemented", val.Kind())
	}
}

func zeroPadding(b []byte) []byte {
	tail := len(b) % 4
	if tail != 0 {
		return append(b, make([]byte, 4-tail)...)
	}
	return b
}

func EncodeLength(i int) []byte {
	if i >= 254 {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(i<<8))
		b[0] = 254
		return b
	} else {
		return []byte{byte(i)}
	}
}

func encodeStruct(val reflect.Value) ([]byte, error) {
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return encodeSumType(val)
	} else {
		return encodeBasicStruct(val)
	}
}

func encodeBasicStruct(val reflect.Value) ([]byte, error) {
	var buf []byte
	for i := 0; i < val.NumField(); i++ {
		b, err := Marshal(val.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		buf = append(buf, b...)
	}
	return buf, nil
}

func encodeSumType(val reflect.Value) ([]byte, error) {
	name := val.FieldByName("SumType").String()
	var buf []byte
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag := val.Type().Field(i).Tag.Get("tlSumType")
		if name != val.Type().Field(i).Name {
			continue
		}
		t, err := encodeTag(tag)
		if err != nil {
			return nil, err
		}
		buf = append(buf, t[:]...)
		b, err := Marshal(val.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		buf = append(buf, b...)
		break
	}
	return buf, nil
}

func encodeTag(tag string) ([4]byte, error) {
	b, err := hex.DecodeString(tag)
	if err != nil {
		return [4]byte{}, err
	}
	if len(b) != 4 {
		return [4]byte{}, fmt.Errorf("invalid tag")
	}
	return [4]byte{b[3], b[2], b[1], b[0]}, nil
}

func encodeVector(val reflect.Value) ([]byte, error) {
	b := make([]byte, 4)
	ln := val.Len()
	binary.LittleEndian.PutUint32(b[:], uint32(ln))
	for i := 0; i < val.Len(); i++ {
		b1, err := Marshal(val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		b = append(b, b1...)
	}
	return b, nil
}
