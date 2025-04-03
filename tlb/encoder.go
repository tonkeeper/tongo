package tlb

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
)

type Encoder struct {
}

type MarshalerTLB interface {
	MarshalTLB(c *boc.Cell, encoder *Encoder) error
}

type tagEncoder interface {
	EncodeTag(c *boc.Cell, tag string) error
}

func Marshal(c *boc.Cell, o any) error {
	encoder := Encoder{}
	return encode(c, "", o, &encoder)
}

func (enc *Encoder) Marshal(c *boc.Cell, o any) error {
	return encode(c, "", o, enc)
}

func isNil(o any) bool {
	switch reflect.ValueOf(o).Kind() {
	case reflect.Interface, reflect.Slice, reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer:
		return reflect.ValueOf(o).IsNil()
	}
	return false

}

func encode(c *boc.Cell, tag string, o any, encoder *Encoder) error {
	if m, ok := o.(tagEncoder); ok {
		return m.EncodeTag(c, tag)
	}
	t, err := parseTag(tag)
	if err != nil {
		return err
	}
	tag = ""
	switch {
	case t.IsMaybeRef:
		if isNil(o) {
			err := c.WriteBit(false)
			return err
		}
		if err := c.WriteBit(true); err != nil {
			return err
		}
		c, err = c.NewRef()
		if err != nil {
			return err
		}

	case t.IsMaybe:
		if isNil(o) {
			err := c.WriteBit(false)
			return err
		}
		if err := c.WriteBit(true); err != nil {
			return err
		}
	case t.IsRef:
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
		return encode(c, tag, val.Elem().Interface(), encoder)
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
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("encoding slice of %v not supported", val.Type().Elem().Kind())
		}
		return c.WriteBytes(val.Bytes())
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
		tag := val.Type().Field(i).Tag.Get("tlb")
		if err := encode(c, tag, val.Field(i).Interface(), encoder); err != nil {
			return err
		}
	}
	return nil
}

func encodeSumType(c *boc.Cell, o any, encoder *Encoder) error {
	val := reflect.ValueOf(o)
	name := val.FieldByName("SumType").String()

	if name == "" {
		return fmt.Errorf("empty SumType value")
	}

	found := false
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().Name() == "SumType" {
			continue
		}
		if name == val.Type().Field(i).Name {
			found = true
			tag := val.Type().Field(i).Tag.Get("tlbSumType")
			if err := encodeSumTag(c, tag); err != nil {
				return err
			}
			if err := encode(c, "", val.Field(i).Interface(), encoder); err != nil {
				return err
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid SumType value: %s", name)
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
