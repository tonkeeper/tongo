package tlb

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/startfellows/tongo/boc"
)

type MarshalerTLB interface {
	MarshalTLB(c *boc.Cell, tag string) error
}

func Marshal(c *boc.Cell, o any) error {
	return encode(c, o, "")
}

func encode(c *boc.Cell, o any, tag string) error {
	t, err := parseTag(tag)
	if err != nil {
		return err
	}
	if t.IsRef {
		c, err = c.NewRef()
		if err != nil {
			return err
		}
	}
	if m, ok := o.(MarshalerTLB); ok {
		return m.MarshalTLB(c, tag)
	}
	val := reflect.ValueOf(o)
	switch val.Kind() {
	case reflect.Uint8:
		return c.WriteUint(val.Uint(), 8)
	case reflect.Uint16:
		return c.WriteUint(val.Uint(), 16)
	case reflect.Uint32:
		return c.WriteUint(val.Uint(), 32)
	case reflect.Uint64:
		return c.WriteUint(val.Uint(), 64)
	case reflect.Int8:
		return c.WriteInt(val.Int(), 8)
	case reflect.Int16:
		return c.WriteInt(val.Int(), 16)
	case reflect.Int32:
		return c.WriteInt(val.Int(), 32)
	case reflect.Int64:
		return c.WriteInt(val.Int(), 64)
	case reflect.Bool:
		err := c.WriteBit(val.Bool())
		if err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		switch o.(type) {
		case big.Int:
			return encodeBigInt(c, o, tag)
		default:
			return encodeStruct(c, o, tag)
		}
	case reflect.Pointer:
		if val.IsNil() && !t.IsOptional {
			return fmt.Errorf("can't encode empty pointer %v if tlb scheme is not optional", val.Type())
		}
		return encode(c, val.Elem().Interface(), tag)
	default:
		return fmt.Errorf("type %v not implemented", val.Kind())
	}
}

func encodeStruct(c *boc.Cell, o any, tag string) error {
	val := reflect.ValueOf(o)
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return encodeSumType(c, o, tag)
	} else {
		return encodeBasicStruct(c, o, tag)
	}
}

func encodeBasicStruct(c *boc.Cell, o any, tag string) error {
	val := reflect.ValueOf(o)
	for i := 0; i < val.NumField(); i++ {
		var err error
		err = encode(c, val.Field(i).Interface(), val.Type().Field(i).Tag.Get("tlb"))
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeSumType(c *boc.Cell, o any, tag string) error {
	val := reflect.ValueOf(o)
	name := val.FieldByName("SumType").String()
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag = val.Type().Field(i).Tag.Get("tlbSumType")
		if name != val.Type().Field(i).Name {
			continue
		}
		err := encodeSumTag(c, tag)
		if err != nil {
			return err
		}
		err = encode(c, val.Field(i).Interface(), "")
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func encodeBigInt(c *boc.Cell, o any, tag string) error {
	// TODO: implement
	return fmt.Errorf("bigInt encoding not implemented")
}
