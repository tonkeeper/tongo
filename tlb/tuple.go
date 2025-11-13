package tlb

import (
	"fmt"
	"reflect"
)

func (t *VmStkTuple) Unmarshal(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return fmt.Errorf("should be a pointer. not %v", val.Kind())
	}
	switch val.Elem().Kind() {
	case reflect.Slice:
		items, err := t.RecursiveToSlice()
		if err != nil {
			if t.Data != nil {
				items, err = t.Data.RecursiveToSlice(int(t.Len))
				if err != nil {
					return err
				}
			} else {
				items = []VmStackValue{}
			}
		}
		sl := reflect.MakeSlice(val.Elem().Type(), 0, len(items))
		for _, i := range items {
			value := reflect.New(val.Elem().Type().Elem())
			err := i.Unmarshal(value.Interface())
			if err != nil {
				return err
			}
			sl = reflect.Append(sl, value.Elem())
		}
		val.Elem().Set(sl)
	case reflect.Struct:
		if int(t.Len) != val.Elem().Type().NumField() {
			return fmt.Errorf("mismatched fields count in tuple and struct")
		}

		values, err := t.Data.RecursiveToSlice(int(t.Len))
		if err != nil {
			return err
		}
		for i := 0; i < val.Elem().Type().NumField(); i++ {
			value := reflect.New(val.Elem().Field(i).Type())
			err = values[i].Unmarshal(value.Interface())
			val.Elem().Field(i).Set(value.Elem())
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("tuple decoding to %v not implemented", val.Elem().Kind())
	}
	return nil
}

func (t *VmStkTuple) RecursiveToSlice() ([]VmStackValue, error) {
	if t.Len != 2 {
		return nil, fmt.Errorf("recursive tuple element must have length value == 2 not %v", t.Len)
	}
	if t.Data.Head.Entry == nil {
		return nil, fmt.Errorf("unsupported tuple format")
	}
	sl := []VmStackValue{*t.Data.Head.Entry}
	if t.Data.Tail.SumType == "VmStkNull" {
		return sl, nil
	}
	if t.Data.Tail.SumType != "VmStkTuple" {
		return nil, fmt.Errorf("invalid type %v in recursive slice decoding", t.Data.Tail.SumType)
	}

	values, err := t.Data.Tail.VmStkTuple.RecursiveToSlice()
	if err != nil {
		return nil, err
	}
	return append(sl, values...), nil
}

func (t VmTuple) RecursiveToSlice(depth int) ([]VmStackValue, error) {
	if t.Head.Entry == nil && t.Head.Ref == nil {
		return nil, fmt.Errorf("can't decode tuple by unknown reason")
	}

	var sl []VmStackValue
	var err error
	if t.Head.Ref == nil {
		sl = append(sl, *t.Head.Entry)
	} else {
		sl, err = t.Head.Ref.RecursiveToSlice(depth - 1)
		if err != nil {
			return nil, err
		}
	}
	if t.Tail.SumType != "" { // tail can be null only if len(tuple) <= 1
		sl = append(sl, t.Tail)
	}
	return sl, nil
}
