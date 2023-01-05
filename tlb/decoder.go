package tlb

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/startfellows/tongo/boc"
)

type UnmarshalerTLB interface {
	UnmarshalTLB(c *boc.Cell, tag string) error
}

func Unmarshal(c *boc.Cell, o any) error {
	//TODO: clone cell with cursor reset
	return decode(c, reflect.ValueOf(o), "")
}

func decode(c *boc.Cell, val reflect.Value, tag string) error {
	t, err := parseTag(tag)
	if err != nil {
		return err
	}
	if t.IsRef {
		c, err = c.NextRef()
		if err != nil {
			return err
		}
		if c.CellType() == boc.PrunedBranchCell {
			return nil
		}
	}
	i, ok := reflect.New(val.Type()).Interface().(UnmarshalerTLB)
	if ok {
		err := i.UnmarshalTLB(c, tag)
		if err != nil {
			return err
		}
		if !val.CanSet() {
			return fmt.Errorf("value can't be changed")
		}
		val.Set(reflect.ValueOf(i).Elem())
		return nil
	}
	if !val.CanSet() && val.Kind() != reflect.Pointer {
		return fmt.Errorf("value can't be changed")
	}
	switch val.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		l := 8
		switch val.Kind() {
		case reflect.Int16:
			l = 16
		case reflect.Int32:
			l = 32
		case reflect.Int64:
			l = 64
		}
		v, err := c.ReadInt(l)
		if err != nil {
			return err
		}
		val.SetInt(v)
		return nil
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		l := 8
		switch val.Kind() {
		case reflect.Uint16:
			l = 16
		case reflect.Uint32:
			l = 32
		case reflect.Uint64:
			l = 64
		}
		v, err := c.ReadUint(l)
		if err != nil {
			return err
		}
		val.SetUint(v)
		return nil
	case reflect.Bool:
		b, err := c.ReadBit()
		if err != nil {
			return err
		}
		val.SetBool(b)
		return nil
	case reflect.Struct:
		t := reflect.New(val.Type()).Interface()
		switch t.(type) {
		case *big.Int:
			return decodeBigInt(c, val, tag)
		default:
			return decodeStruct(c, val, tag)
		}
	case reflect.Pointer:
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			return decode(c, val.Elem(), tag)
		}
		a := reflect.New(val.Type().Elem())
		err = decode(c, a, tag)
		if err != nil {
			return err
		}
		val.Set(a)
		return nil
	default:
		return fmt.Errorf("type %v not implemented", val.Kind())
	}
}

func decodeStruct(c *boc.Cell, val reflect.Value, tag string) error {
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return decodeSumType(c, val, tag)
	} else {
		return decodeBasicStruct(c, val, tag)
	}
}

func decodeBasicStruct(c *boc.Cell, val reflect.Value, tag string) error {
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		err := decode(c, val.Field(i), val.Type().Field(i).Tag.Get("tlb"))
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeSumType(c *boc.Cell, val reflect.Value, tag string) error {
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag = val.Type().Field(i).Tag.Get("tlbSumType")
		ok, err := compareWithSumTag(c, tag)
		if err != nil {
			return err
		}
		if ok {
			val.FieldByName("SumType").SetString(val.Type().Field(i).Name)
			err := decode(c, val.Field(i), "")
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("can not decode sumtype %v", tag)
}

func compareWithSumTag(c *boc.Cell, tag string) (bool, error) {
	t, err := parseSumTag(tag)
	if err != nil {
		return false, err
	}
	y, err := c.PickUint(t.Len)
	if t.Val == y {
		_ = c.Skip(t.Len) // already checked
		return true, nil
	}
	return false, nil
}

func decodeBigInt(c *boc.Cell, val reflect.Value, tag string) error {
	// TODO: implement
	return fmt.Errorf("bigInt decoding not implemented")
}
