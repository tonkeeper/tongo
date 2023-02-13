package tlb

import (
	"fmt"
	"reflect"
)

func (t *VmStkTuple) UnmarshalTo(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return fmt.Errorf("should be a pointer. not %v", val.Kind())
	}
	switch val.Elem().Kind() {
	case reflect.Slice:
		items, err := t.recursiveToSlice()
		if err != nil {
			return err
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

		values, err := t.Data.recursiveToSlice(int(t.Len))
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

func (t *VmStkTuple) recursiveToSlice() ([]VmStackValue, error) {
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

	values, err := t.Data.Tail.VmStkTuple.recursiveToSlice()
	if err != nil {
		return nil, err
	}
	return append(sl, values...), nil
}

func (t VmTuple) recursiveToSlice(depth int) ([]VmStackValue, error) {
	var sl []VmStackValue
	var err error
	if depth == 2 {
		if t.Head.Entry == nil {
			return nil, fmt.Errorf("stack tuple invalid depth")
		}
		sl = append(sl, *t.Head.Entry)
	} else {
		sl, err = t.Head.Ref.recursiveToSlice(depth - 1)
	}
	return append(sl, t.Tail), err
}
