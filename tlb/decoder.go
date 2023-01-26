package tlb

import (
	"fmt"
	"reflect"

	"github.com/startfellows/tongo/boc"
)

type Decoder struct {
	tag string
}

type UnmarshalerTLB interface {
	UnmarshalTLB(c *boc.Cell, decoder *Decoder) error
}

func Unmarshal(c *boc.Cell, o any) error {
	//TODO: clone cell with cursor reset
	decoder := Decoder{}
	return decode(c, reflect.ValueOf(o), &decoder)
}

func decode(c *boc.Cell, val reflect.Value, decoder *Decoder) error {
	t, err := parseTag(decoder.tag)
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

	i, ok := val.Interface().(UnmarshalerTLB)
	if ok {
		return i.UnmarshalTLB(c, decoder)
	}
	if val.CanAddr() {
		i, ok := val.Addr().Interface().(UnmarshalerTLB)
		if ok {
			return i.UnmarshalTLB(c, decoder)
		}
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
		case *boc.Cell:
			return decodeCell(c, val)
		case *boc.BitString:
			return decodeBitString(c, val)
		default:
			return decodeStruct(c, val, decoder)
		}
	case reflect.Pointer:
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			return decode(c, val.Elem(), decoder)
		}
		a := reflect.New(val.Type().Elem())
		err = decode(c, a, decoder)
		if err != nil {
			return err
		}
		val.Set(a)
		return nil
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("decoding array of %v not supported", val.Type().Elem().Kind())
		}
		v, err := c.ReadBytes(val.Len())
		if err != nil {
			return err
		}
		reflect.Copy(val, reflect.ValueOf(v))
		return nil
	default:
		return fmt.Errorf("type %v not implemented", val.Kind())
	}
}

func decodeStruct(c *boc.Cell, val reflect.Value, decoder *Decoder) error {
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return decodeSumType(c, val, decoder)
	} else {
		return decodeBasicStruct(c, val, decoder)
	}
}

func decodeBasicStruct(c *boc.Cell, val reflect.Value, decoder *Decoder) error {
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		decoder.tag = val.Type().Field(i).Tag.Get("tlb")
		err := decode(c, val.Field(i), decoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeSumType(c *boc.Cell, val reflect.Value, decoder *Decoder) error {
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag := val.Type().Field(i).Tag.Get("tlbSumType")
		ok, err := compareWithSumTag(c, tag)
		if err != nil {
			return err
		}
		if ok {
			val.FieldByName("SumType").SetString(val.Type().Field(i).Name)
			decoder.tag = ""
			err := decode(c, val.Field(i), decoder)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("can not decode sumtype %v", decoder.tag)
}

func compareWithSumTag(c *boc.Cell, tag string) (bool, error) {
	t, err := ParseTag(tag)
	if err != nil {
		return false, err
	}
	if c.BitsAvailableForRead() < t.Len {
		return false, nil
	}
	y, err := c.PickUint(t.Len)
	if err != nil {
		return false, err
	}
	if t.Val == y {
		_ = c.Skip(t.Len) // already checked
		return true, nil
	}
	return false, nil
}

func decodeCell(c *boc.Cell, val reflect.Value) error {
	if !val.CanSet() {
		return fmt.Errorf("value can't be changed")
	}
	val.Set(reflect.ValueOf(c).Elem())
	return nil
}

func decodeBitString(c *boc.Cell, val reflect.Value) error {
	return fmt.Errorf("bigString decoding not supported")
}
