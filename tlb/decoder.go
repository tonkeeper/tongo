package tlb

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"reflect"
)

type UnmarshalerTLB interface {
	UnmarshalTLB(c *boc.Cell, tag string) error
}

func Unmarshal(c *boc.Cell, o any) error {
	//TODO: clone cell with cursor reset
	return decode(c, reflect.ValueOf(o), "")
}

func decode(c *boc.Cell, val reflect.Value, tag string) error {
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	i, ok := reflect.New(val.Type()).Interface().(UnmarshalerTLB)
	if ok {
		err := i.UnmarshalTLB(c, tag)
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
		var ln int
		if tag != "" {
			_, err := fmt.Sscanf(tag, "%dbits", &ln)
			if err != nil {
				return err
			}
			if ln > 32 {
				return fmt.Errorf("can not marshal %v bits to 32 bits integer", ln)
			}
		} else {
			ln = 32
		}
		v, err := c.ReadUint(ln)
		if err != nil {
			return err
		}
		if val.Kind() == reflect.Uint32 {
			val.SetUint(v)
		} else {
			val.SetInt(int64(v))
		}
		return nil
	case reflect.Uint64, reflect.Int64:
		var ln int
		if tag != "" {
			_, err := fmt.Sscanf(tag, "%dbits", &ln)
			if err != nil {
				return err
			}
			if ln > 64 {
				return fmt.Errorf("can not marshal %v bits to 64 bits integer", ln)
			}
		} else {
			ln = 64
		}
		v, err := c.ReadUint(ln)
		if err != nil {
			return err
		}
		if val.Kind() == reflect.Uint64 {
			val.SetUint(v)
		} else {
			val.SetInt(int64(v))
		}
		return nil
	case reflect.Bool:
		b, err := c.ReadBit()
		if err != nil {
			return err
		}
		val.SetBool(b)
		return nil
	case reflect.Struct:
		return decodeStruct(c, val, tag)
	default:
		return fmt.Errorf("type %v not emplemented", val.Kind())
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
		ok, err := compareWithTag(c, tag)
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
	return fmt.Errorf("can not decode sumtype")
}

func compareWithTag(c *boc.Cell, tag string) (bool, error) {
	t, err := parseTag(tag)
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
