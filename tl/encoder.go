package tl

import (
	"bytes"
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
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			b, err := encodeVector(val)
			if err != nil {
				return nil, err
			}
			return zeroPadding(b), nil
			//return nil, fmt.Errorf("encoding slice of %v not supported", val.Type().Elem().Kind())
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
		return nil, fmt.Errorf("type %v not emplemented", val.Kind())
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
	var res [4]byte
	b, err := hex.DecodeString(tag)
	if err != nil {
		return [4]byte{}, err
	}
	copy(res[:], b)
	return res, nil
}

func encodeVector(val reflect.Value) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint32(val.Len()))
	if err != nil {
		return nil, err
	}
	for i := 0; i < val.Len(); i++ {
		err = binary.Write(buf, binary.LittleEndian, val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
