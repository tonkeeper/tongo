package tlb

import (
	"fmt"
	"reflect"

	"github.com/startfellows/tongo/boc"
)

type Encoder struct {
	tag string
}

type MarshalerTLB interface {
	MarshalTLB(c *boc.Cell, encoder *Encoder) error
}

func Marshal(c *boc.Cell, o any) error {
	encoder := Encoder{}
	return encode(c, o, &encoder)
}

func encode(c *boc.Cell, o any, encoder *Encoder) error {
	t, err := parseTag(encoder.tag)
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
		return m.MarshalTLB(c, encoder)
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
		case boc.Cell:
			return encodeCell(c, o)
		case boc.BitString:
			return encodeBitString(c, o)
		default:
			return encodeStruct(c, o, encoder)
		}
	case reflect.Pointer:
		if val.IsNil() && !t.IsOptional {
			return fmt.Errorf("can't encode empty pointer %v if tlb scheme is not optional", val.Type())
		}
		return encode(c, val.Elem().Interface(), encoder)
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("encoding array of %v not supported", val.Type().Elem().Kind())
		}
		// TODO: optimize
		b := make([]byte, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			b = append(b, uint8(val.Index(i).Uint()))
		}
		return c.WriteBytes(b)
	default:
		return fmt.Errorf("type %v not implemented", val.Kind())
	}
}

func encodeStruct(c *boc.Cell, o any, encoder *Encoder) error {
	val := reflect.ValueOf(o)
	if _, ok := val.Type().FieldByName("SumType"); ok {
		return encodeSumType(c, o, encoder)
	} else {
		return encodeBasicStruct(c, o, encoder)
	}
}

func encodeBasicStruct(c *boc.Cell, o any, encoder *Encoder) error {
	val := reflect.ValueOf(o)
	for i := 0; i < val.NumField(); i++ {
		var err error
		encoder.tag = val.Type().Field(i).Tag.Get("tlb")
		err = encode(c, val.Field(i).Interface(), encoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeSumType(c *boc.Cell, o any, encoder *Encoder) error {
	val := reflect.ValueOf(o)
	name := val.FieldByName("SumType").String()
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		tag := val.Type().Field(i).Tag.Get("tlbSumType")
		if name != val.Type().Field(i).Name {
			continue
		}
		err := encodeSumTag(c, tag)
		if err != nil {
			return err
		}
		encoder.tag = ""
		err = encode(c, val.Field(i).Interface(), encoder)
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func encodeCell(c *boc.Cell, o any) error {
	*c = o.(boc.Cell)
	return nil
}

func encodeBitString(c *boc.Cell, o any) error {
	err := c.WriteBitString(o.(boc.BitString))
	if err != nil {
		return err
	}
	return nil
}
