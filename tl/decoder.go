package tl

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
)

type SumType string

func Unmarshal(r io.Reader, o any) error {
	return decode(r, reflect.ValueOf(o))
}

type UnmarshalerTL interface {
	UnmarshalTL(r io.Reader) error
}

func decode(buf io.Reader, val reflect.Value) error {
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	i, ok := reflect.New(val.Type()).Interface().(UnmarshalerTL)
	if ok {
		err := i.UnmarshalTL(buf)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(i).Elem())
		return nil
	}

	if !val.CanSet() {
		return fmt.Errorf("value can't be changed")
	}
	switch val.Kind() {
	case reflect.Uint32, reflect.Int32:
		b := make([]byte, 4)
		_, err := io.ReadFull(buf, b)
		if err != nil {
			return err
		}
		if val.Kind() == reflect.Uint32 {
			val.SetUint(uint64(binary.LittleEndian.Uint32(b)))
		} else {
			val.SetInt(int64(binary.LittleEndian.Uint32(b)))
		}
		return nil
	case reflect.Uint64, reflect.Int64:
		b := make([]byte, 8)
		_, err := io.ReadFull(buf, b)
		if err != nil {
			return err
		}
		if val.Kind() == reflect.Uint64 {
			val.SetUint(binary.LittleEndian.Uint64(b))
		} else {
			val.SetInt(int64(binary.LittleEndian.Uint64(b)))
		}
		return nil
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return decodeVector(buf, val)
		}
		data, err := readByteSlice(buf)
		if err != nil {
			return err
		}
		val.SetBytes(data)
		return nil
	case reflect.String:
		data, err := readByteSlice(buf)
		if err != nil {
			return err
		}
		val.SetString(string(data))
		return nil
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("decoding array of %v not supported", val.Type().Elem().Kind())
		}
		data, err := readByteSlice(buf)
		if err != nil {
			return err
		}
		if val.Type().Len() != len(data) {
			return fmt.Errorf("mismatched lenghth of decoded byte slice (%v) and array (%v)", len(data), val.Type().Len())
		}
		reflect.Copy(val, reflect.ValueOf(data))
		return nil
	case reflect.Struct:
		return decodeStruct(buf, val)

	default:
		return fmt.Errorf("type %v not emplemented", val.Kind())
	}
}

func readByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	return b[0], err
}

func readByteSlice(r io.Reader) ([]byte, error) {
	firstByte, err := readByte(r)
	if err != nil {
		return nil, err
	}
	var data []byte
	var full int
	if firstByte < 254 {
		data = make([]byte, int(firstByte))
		_, err := io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}
		full = 1 + len(data)
	} else if firstByte == 254 {
		sizeBuf := make([]byte, 4)
		_, err := io.ReadFull(r, sizeBuf[:3])
		if err != nil {
			return nil, err
		}
		data = make([]byte, binary.LittleEndian.Uint32(sizeBuf))
		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}
		full = 4 + len(data)
	} else {
		return nil, fmt.Errorf("invalid bytes prefix")
	}
	for ; full%4 != 0; full++ {
		_, err = readByte(r)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func decodeStruct(buf io.Reader, val reflect.Value) error {
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return decodeSumType(buf, val)
	} else {
		return decodeBasicStruct(buf, val)
	}
}

func decodeSumType(r io.Reader, val reflect.Value) error {
	var tagBytes [4]byte
	_, err := io.ReadFull(r, tagBytes[:])
	if err != nil {
		return err
	}
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag := val.Type().Field(i).Tag.Get("tlSumType")
		ok, err := compareWithTag(tagBytes, tag)
		if err != nil {
			return err
		}
		if ok {
			val.FieldByName("SumType").SetString(val.Type().Field(i).Name)
			err := decode(r, val.Field(i))
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("can not decode sumtype")
}

func decodeBasicStruct(r io.Reader, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		err := decode(r, val.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func compareWithTag(tagBytes [4]byte, tag string) (bool, error) {
	var a [4]byte
	t, err := hex.DecodeString(tag)
	if err != nil {
		return false, err
	}
	copy(a[:], t)
	return bytes.Equal(a[:], tagBytes[:]), nil
}

func decodeVector(r io.Reader, val reflect.Value) error {
	var b [4]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	ln := int(binary.LittleEndian.Uint32(b[:]))
	item := reflect.New(val.Type().Elem())
	slice := reflect.MakeSlice(val.Type(), 0, ln)
	for i := 0; i < ln; i++ {
		err := decode(r, item)
		slice = reflect.Append(slice, item.Elem())
		if err != nil {
			return err
		}
	}
	val.Set(slice)
	return nil
}
