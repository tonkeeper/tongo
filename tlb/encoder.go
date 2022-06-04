package tlb

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"reflect"
	"strconv"
	"strings"
)

type MarshalerTLB interface {
	MarshalTLB(c *boc.Cell, tag string) error
}

func Marshal(c *boc.Cell, o any) error {
	return encode(c, o, "")
}

func encode(c *boc.Cell, o any, tag string) error {
	if m, ok := o.(MarshalerTLB); ok {
		return m.MarshalTLB(c, tag)
	}
	val := reflect.ValueOf(o)
	switch val.Kind() {
	case reflect.Uint32, reflect.Uint64: //todo:check bit length
		var ln int
		if tag != "" {
			_, err := fmt.Sscanf(tag, "%dbits", &ln)
			if err != nil {
				return err
			}
		} else {
			ln = 32
			if val.Kind() == reflect.Uint64 {
				ln = 64
			}
		}
		err := c.WriteUint(val.Uint(), ln)
		if err != nil {
			return err
		}
		return nil
	case reflect.Int32, reflect.Int64: //todo:check bit length
		var ln int
		if tag != "" {
			_, err := fmt.Sscanf(tag, "%dbits", &ln)
			if err != nil {
				return err
			}
		} else {
			ln = 32
			if val.Kind() == reflect.Int64 {
				ln = 64
			}
		}
		err := c.WriteInt(val.Int(), ln)
		if err != nil {
			return err
		}
		return nil
	case reflect.Bool:
		err := c.WriteBit(val.Bool())
		if err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		return encodeStruct(c, o, tag)
	default:
		return fmt.Errorf("type %v not emplemented", val.Kind())
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
		err := encodeTag(c, tag)
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

func encodeTag(c *boc.Cell, tag string) error {
	t, err := parseTag(tag)
	if err != nil {
		return err
	}
	err = c.WriteUint(t.Val, t.Len)
	if err != nil {
		return err
	}
	return nil
}

func parseTag(s string) (tag, error) {
	a := strings.Split(s, "$")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 2, 32)
		if err != nil {
			return tag{}, err
		}
		return tag{a[0], len(a[1]), x}, nil
	}
	a = strings.Split(s, "#")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 16, 32)
		if err != nil {
			return tag{}, err
		}
		return tag{a[0], len(a[1]) * 4, x}, nil
	}
	return tag{}, fmt.Errorf("invalid tag")
}
