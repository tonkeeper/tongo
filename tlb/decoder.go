package tlb

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
)

type resolveLib func(hash Bits256) (*boc.Cell, error)

// Decoder unmarshals a cell into a golang type.
type Decoder struct {
	hasher     *boc.Hasher
	withDebug  bool
	debugPath  []string
	resolveLib resolveLib
}

func (d *Decoder) WithDebug() *Decoder {
	d.withDebug = true
	return d
}

// WithLibraryResolver provides a function which is used to fetch a library cell by its hash.
func (d *Decoder) WithLibraryResolver(resolveLib resolveLib) *Decoder {
	d.resolveLib = resolveLib
	return d
}

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		hasher: boc.NewHasher(),
	}
}

// Unmarshal decodes the give cell using TL-B schema and stores the result in the value pointed to by o.
func (dec *Decoder) Unmarshal(c *boc.Cell, o any) error {
	return decode(c, "", reflect.ValueOf(o), dec)
}

// UnmarshalerTLB contains method UnmarshalTLB that must be implemented by a struct
// if it provides specific unmarshalling code.
type UnmarshalerTLB interface {
	UnmarshalTLB(c *boc.Cell, decoder *Decoder) error
}

type tagValidator interface {
	ValidateTag(c *boc.Cell, tag string) error
}

func Unmarshal(c *boc.Cell, o any) error {
	dec := Decoder{withDebug: true}
	return decode(c, "", reflect.ValueOf(o), &dec)
}

var bocCellType = reflect.TypeOf(boc.Cell{})
var bocCellPointerType = reflect.TypeOf(&boc.Cell{})
var bocTlbANyPointerType = reflect.TypeOf(&Any{})
var bitStringType = reflect.TypeOf(boc.BitString{})

func decode(c *boc.Cell, tag string, val reflect.Value, decoder *Decoder) error {
	if decoder.withDebug {
		decoder.debugPath = append(decoder.debugPath, fmt.Sprintf("%v#%v", val.Type().Name(), tag))
		defer func() {
			decoder.debugPath = decoder.debugPath[:len(decoder.debugPath)-1]
		}()
	}
	if c.IsLibrary() {
		if val.Kind() == reflect.Ptr && val.Type() == bocCellPointerType {
			// this is a library cell, and we unmarshal it to a cell.
			// let's not resolve it and keep it as is
			val.Elem().Set(reflect.ValueOf(c).Elem())
			return nil
		}
		if val.Kind() == reflect.Ptr && val.Type() == bocTlbANyPointerType {
			//todo: remove
			a := Any(*c)
			val.Elem().Set(reflect.ValueOf(a))
			return nil
		}
		if decoder.resolveLib == nil {
			// use WithLibraryResolver to provide a library resolver
			return fmt.Errorf("library cell decoding is not configured properly")
		}
		hash, err := c.Hash256()
		if err != nil {
			return err
		}
		cell, err := decoder.resolveLib(hash)
		if err != nil {
			return err
		}
		c = cell
	}
	t, err := parseTag(tag)
	if err != nil {
		return err
	}
	switch {
	case t.IsMaybeRef:
		tag = ""
		exist, err := c.ReadBit()
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		c, err = c.NextRef()
		if err != nil {
			return err
		}
		if c.IsLibrary() {
			return fmt.Errorf("library cell as a ref is not implemented")
		}
		if c.CellType() == boc.PrunedBranchCell {
			return nil
		}
	case t.IsMaybe:
		tag = ""
		exist, err := c.ReadBit()
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
	case t.IsRef:
		tag = ""
		c, err = c.NextRef()
		if err != nil {
			return err
		}
		if c.IsLibrary() {
			if val.Kind() == reflect.Struct && val.Type() == bocCellType {
				// this is a library cell, and we unmarshal it to a cell.
				// let's not resolve it and keep it as is
				val.Set(reflect.ValueOf(c).Elem())
				return nil
			}
			return fmt.Errorf("library cell as a ref is not implemented")
		}
		if c.CellType() == boc.PrunedBranchCell {
			return nil
		}
	}
	i, ok := val.Interface().(UnmarshalerTLB)
	if ok {
		if val.IsNil() {
			p := reflect.New(val.Type().Elem())
			val.Set(p)
			i = p.Interface().(UnmarshalerTLB)
		}
		return i.UnmarshalTLB(c, decoder)
	}
	if val.CanAddr() {
		i, ok := val.Addr().Interface().(UnmarshalerTLB)
		if ok {
			return i.UnmarshalTLB(c, decoder)
		}
		v, ok := val.Addr().Interface().(tagValidator)
		if ok {
			return v.ValidateTag(c, tag)
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
		switch val.Type() {
		case bocCellType:
			return decodeCell(c, val)
		case bitStringType:
			return decodeBitString(c, val)
		default:
			return decodeStruct(c, val, decoder)
		}
	case reflect.String:
		bs := c.ReadRemainingBits()
		bytes, err := bs.GetTopUppedArray()
		if err != nil {
			return err
		}
		val.SetString(string(bytes))
		return nil
	case reflect.Pointer:
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			return decode(c, tag, val.Elem(), decoder)
		}
		a := reflect.New(val.Type().Elem())
		err = decode(c, tag, a, decoder)
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
		if decoder.withDebug {
			decoder.debugPath = append(decoder.debugPath, val.Type().Field(i).Name)
		}
		if !val.Field(i).CanSet() {
			return fmt.Errorf("can't set field %v", i)
		}
		tag := val.Type().Field(i).Tag.Get("tlb")
		err := decode(c, tag, val.Field(i), decoder)
		if err != nil {
			return err
		}
		if decoder.withDebug {
			decoder.debugPath = decoder.debugPath[:len(decoder.debugPath)-1]
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
			if decoder.withDebug {
				decoder.debugPath = append(decoder.debugPath, val.Type().Field(i).Name)
				defer func() {
					decoder.debugPath = decoder.debugPath[:len(decoder.debugPath)-1]
				}()
			}
			val.FieldByName("SumType").SetString(val.Type().Field(i).Name)
			err := decode(c, "", val.Field(i), decoder)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("can not decode sumtype %v", val.Type().Name())
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

// Hasher returns boc.Hasher that is used to calculate hashes when decoding.
func (dec *Decoder) Hasher() *boc.Hasher {
	return dec.hasher
}
