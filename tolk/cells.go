package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type Cell struct{}

func (Cell) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	v.sumType = "cell"
	return nil
}

func (c Cell) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	err = c.SetValue(v, Any(*ref))
	if err != nil {
		return err
	}
	return nil
}

func (a *Any) UnmarshalTolk(cell *boc.Cell, ty Cell, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	*a = Any(*ref)

	return nil
}

func (Cell) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.cell == nil {
		return fmt.Errorf("ref not found")
	}

	c := boc.Cell(*v.cell)
	ref := c.CopyRemaining()
	err := cell.AddRef(ref)
	if err != nil {
		return err
	}

	return nil
}

func (Cell) Equal(v Value, o Value) bool {
	return false
}

type Slice struct{}

func (Slice) SetValue(v *Value, val any) error {
	return fmt.Errorf("slice is not supported")
}

func (Slice) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return fmt.Errorf("slice is not supported")
}

func (Slice) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("slice is not supported")
}

func (Slice) Equal(v Value, o Value) bool {
	return false
}

type Builder struct{}

func (Builder) SetValue(v *Value, val any) error {
	return fmt.Errorf("builder is not supported")
}

func (Builder) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return fmt.Errorf("builder is not supported")
}

func (Builder) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("builder is not supported")
}

func (Builder) Equal(v Value, o Value) bool {
	return false
}

type Callable struct{}

func (Callable) SetValue(v *Value, val any) error {
	return fmt.Errorf("callable is not supported")
}

func (Callable) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return fmt.Errorf("callable is not supported")
}

func (Callable) MarshalTolk(cell *boc.Cell, v *Value) error {
	return fmt.Errorf("callable is not supported")
}

func (Callable) Equal(v Value, o Value) bool {
	return false
}

type Remaining struct{}

func (Remaining) SetValue(v *Value, val any) error {
	a, ok := val.(Any)
	if !ok {
		return fmt.Errorf("value is not a cell")
	}
	v.cell = &a
	v.sumType = "cell"
	return nil
}

func (r Remaining) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	rem := cell.CopyRemaining()
	if rem != nil {
		err := r.SetValue(v, Any(*rem))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RemainingValue) UnmarshalTolk(cell *boc.Cell, ty Remaining, decoder *Decoder) error {
	rem := cell.CopyRemaining()
	if rem != nil {
		*r = RemainingValue(*rem)
	}
	return nil
}

func (Remaining) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.cell == nil {
		return fmt.Errorf("remaining not found")
	}
	c := boc.Cell(*v.cell)
	err := cell.WriteBitString(c.ReadRemainingBits())
	if err != nil {
		return err
	}
	for _, ref := range c.Refs() {
		err = cell.AddRef(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Remaining) Equal(v Value, o Value) bool {
	return false
}

type Nullable struct {
	Inner Ty `json:"inner"`
}

func (Nullable) SetValue(v *Value, val any) error {
	o, ok := val.(OptValue)
	if !ok {
		return fmt.Errorf("value is not an optional value")
	}
	v.optionalValue = &o
	v.sumType = "optionalValue"
	return nil
}

func (n Nullable) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	isExists, err := cell.ReadBit()
	if err != nil {
		return err
	}
	innerV := Value{}
	optV := OptValue{
		IsExists: isExists,
		Val:      innerV,
	}
	if isExists {
		err = innerV.UnmarshalTolk(cell, n.Inner, decoder)
		if err != nil {
			return err
		}
		optV.Val = innerV
	}
	err = n.SetValue(v, optV)
	if err != nil {
		return err
	}
	return nil
}

func (o *OptValue) UnmarshalTolk(cell *boc.Cell, ty Nullable, decoder *Decoder) error {
	isExists, err := cell.ReadBit()
	if err != nil {
		return err
	}
	o.IsExists = isExists
	if isExists {
		err = o.Val.UnmarshalTolk(cell, ty.Inner, decoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (Nullable) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.optionalValue == nil {
	//	return fmt.Errorf("optional value not found")
	//}
	//exists := v.optionalValue.IsExists
	//err := cell.WriteBit(exists)
	//if err != nil {
	//	return err
	//}
	//
	//if exists {
	//	val := v.optionalValue.Val
	//	err = val.valType.MarshalTolk(cell, &val)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (Nullable) Equal(v Value, o Value) bool {
	return false
}

type CellOf struct {
	Inner Ty `json:"inner"`
}

func (CellOf) SetValue(v *Value, val any) error {
	r, ok := val.(RefValue)
	if !ok {
		return fmt.Errorf("value is not a ref value")
	}
	v.refValue = &r
	v.sumType = "refValue"
	return nil
}

func (c CellOf) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	innerV := Value{}
	err = innerV.UnmarshalTolk(ref, c.Inner, decoder)
	if err != nil {
		return err
	}
	err = c.SetValue(v, RefValue(innerV))
	if err != nil {
		return err
	}

	return nil
}

func (r *RefValue) UnmarshalTolk(cell *boc.Cell, ty CellOf, decoder *Decoder) error {
	ref, err := cell.NextRef()
	if err != nil {
		return err
	}
	innerV := Value{}
	err = innerV.UnmarshalTolk(ref, ty.Inner, decoder)
	if err != nil {
		return err
	}
	*r = RefValue(innerV)

	return nil
}

func (CellOf) MarshalTolk(cell *boc.Cell, v *Value) error {
	//if v.refValue == nil {
	//	return fmt.Errorf("ref value not found")
	//}
	//ref := boc.NewCell()
	//val := Value(*v.refValue)
	//err := val.valType.MarshalTolk(ref, &val)
	//if err != nil {
	//	return err
	//}
	//
	//err = cell.AddRef(ref)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (CellOf) Equal(v Value, o Value) bool {
	return false
}
