package tolk

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
)

type Int struct{}

func (i Int) SetValue(v *Value, val any) error {
	return fmt.Errorf("int is not supported")
}

func (Int) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	return fmt.Errorf("int is not supported")
}

func (Int) Equal(v Value, o Value) bool {
	return false
}

type IntN struct {
	N int `json:"n"`
}

func (i IntN) SetValue(v *Value, val any) error {
	n := i.N
	if n > 64 {
		bi, ok := val.(BigInt)
		if !ok {
			return fmt.Errorf("value is not a BigInt")
		}
		v.bigInt = &bi
	} else {
		i64, ok := val.(Int64)
		if !ok {
			bi, ok := val.(BigInt)
			if !ok {
				return fmt.Errorf("value is not a BigInt or Int64")
			}
			b := big.Int(bi)
			i64 = Int64(b.Int64())
		}
		wi64 := i64
		v.smallInt = &wi64
	}
	return nil
}

func (i IntN) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	num, err := cell.ReadBigInt(i.N)
	if err != nil {
		return err
	}
	err = v.SetValue(BigInt(*num))
	if err != nil {
		return err
	}
	return nil
}

func (i IntN) GetFixedSize() int {
	return i.N
}

func (IntN) Equal(v Value, o Value) bool {
	if v.smallInt == nil && o.smallInt == nil {
		if v.bigInt == nil || o.bigInt == nil {
			return false
		}

		vb := big.Int(*v.bigInt)
		ob := big.Int(*o.bigInt)
		return vb.Cmp(&ob) == 0
	}

	if v.smallInt == nil || o.smallInt == nil {
		return false
	}

	return *v.smallInt == *o.smallInt
}

type UintN struct {
	N int `json:"n"`
}

func (u UintN) SetValue(v *Value, val any) error {
	n := u.N
	if n > 64 {
		bi, ok := val.(BigInt)
		if !ok {
			return fmt.Errorf("value is not a BigInt")
		}
		v.bigInt = &bi
	} else {
		ui64, ok := val.(UInt64)
		if !ok {
			bi, ok := val.(BigInt)
			if !ok {
				return fmt.Errorf("value is not a BigInt or UInt64")
			}
			b := big.Int(bi)
			ui64 = UInt64(b.Uint64())
		}
		wui64 := ui64
		v.smallUint = &wui64
	}
	return nil
}

func (u UintN) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	num, err := cell.ReadBigUint(u.N)
	if err != nil {
		return err
	}
	err = v.SetValue(BigInt(*num))
	if err != nil {
		return err
	}
	return nil
}

func (u UintN) GetFixedSize() int {
	return u.N
}

func (UintN) Equal(v Value, o Value) bool {
	if v.smallUint == nil && o.smallUint == nil {
		if v.bigInt == nil || o.bigInt == nil {
			return false
		}

		vb := big.Int(*v.bigInt)
		ob := big.Int(*o.bigInt)
		return vb.Cmp(&ob) == 0
	}

	if v.smallUint == nil || o.smallUint == nil {
		return false
	}

	return *v.smallUint == *o.smallUint
}

type VarIntN struct {
	N int `json:"n"`
}

func (VarIntN) SetValue(v *Value, val any) error {
	bi, ok := val.(BigInt)
	if !ok {
		return fmt.Errorf("value is not a BigInt")
	}
	v.bigInt = &bi
	return nil
}

func (vi VarIntN) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	n := vi.N
	ln, err := cell.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return err
	}
	err = v.SetValue(BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (VarIntN) Equal(v Value, o Value) bool {
	return false
}

type VarUintN struct {
	N int `json:"n"`
}

func (VarUintN) SetValue(v *Value, val any) error {
	bi, ok := val.(BigInt)
	if !ok {
		return fmt.Errorf("value is not a BigInt")
	}
	v.bigInt = &bi
	return nil
}

func (vu VarUintN) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	n := vu.N
	ln, err := cell.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	err = v.SetValue(BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (VarUintN) Equal(v Value, o Value) bool {
	return false
}

type BitsN struct {
	N int `json:"n"`
}

func (BitsN) SetValue(v *Value, val any) error {
	bits, ok := val.(Bits)
	if !ok {
		return fmt.Errorf("value is not a Bits")
	}
	v.bits = &bits
	return nil
}

func (b BitsN) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	n := b.N
	val, err := cell.ReadBits(n)
	if err != nil {
		return err
	}
	err = v.SetValue(Bits(val))
	if err != nil {
		return err
	}
	return nil
}

func (b BitsN) GetFixedSize() int {
	return b.N
}

func (BitsN) Equal(v Value, o Value) bool {
	if v.bits == nil || o.bits == nil {
		return false
	}
	vb := boc.BitString(*v.bits)
	ob := boc.BitString(*o.bits)
	return bytes.Equal(vb.Buffer(), ob.Buffer())
}

type Coins struct {
}

func (Coins) SetValue(v *Value, val any) error {
	bi, ok := val.(BigInt)
	if !ok {
		return fmt.Errorf("value is not a BigInt")
	}
	v.bigInt = &bi
	return nil
}

func (Coins) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(15)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	err = v.SetValue(BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (Coins) Equal(v Value, o Value) bool {
	return false
}

type Bool struct{}

func (Bool) SetValue(v *Value, val any) error {
	b, ok := val.(BoolValue)
	if !ok {
		return fmt.Errorf("value is not a BoolValue")
	}
	v.bool = &b
	return nil
}

func (Bool) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	val, err := cell.ReadBit()
	if err != nil {
		return err
	}
	err = v.SetValue(BoolValue(val))
	if err != nil {
		return err
	}
	return nil
}

func (Bool) Equal(v Value, o Value) bool {
	return false
}
