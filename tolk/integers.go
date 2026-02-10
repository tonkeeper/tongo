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

func (Int) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	return fmt.Errorf("int is not supported")
}

func (Int) MarshalTolk(cell *boc.Cell, v *Value) error {
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
		v.sumType = "bigInt"
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
		v.smallInt = &i64
		v.sumType = "smallInt"
	}
	return nil
}

func (i IntN) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	num, err := cell.ReadBigInt(i.N)
	if err != nil {
		return err
	}
	err = i.SetValue(v, BigInt(*num))
	if err != nil {
		return err
	}
	return nil
}

func (i *Int64) UnmarshalTolk(cell *boc.Cell, ty IntN, decoder *Decoder) error {
	num, err := cell.ReadInt(ty.N)
	if err != nil {
		return err
	}
	*i = Int64(num)
	return nil
}

func (b *BigInt) UnmarshalTolk(cell *boc.Cell, ty IntN, decoder *Decoder) error {
	num, err := cell.ReadBigInt(ty.N)
	if err != nil {
		return err
	}
	*b = BigInt(*num)
	return nil
}

func (i IntN) MarshalTolk(cell *boc.Cell, v *Value) error {
	if i.N > 64 {
		if v.bigInt == nil {
			return fmt.Errorf("big int not found")
		}
		bi := big.Int(*v.bigInt)

		err := cell.WriteBigInt(&bi, i.N)
		if err != nil {
			return err
		}
	} else {
		if v.smallInt == nil {
			return fmt.Errorf("small int not found")
		}

		i64 := int64(*v.smallInt)
		err := cell.WriteInt(i64, i.N)
		if err != nil {
			return err
		}
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
		v.sumType = "bigInt"
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
		v.smallUint = &ui64
		v.sumType = "smallUint"
	}
	return nil
}

func (u UintN) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	num, err := cell.ReadBigUint(u.N)
	if err != nil {
		return err
	}
	err = u.SetValue(v, BigInt(*num))
	if err != nil {
		return err
	}
	return nil
}

func (i *UInt64) UnmarshalTolk(cell *boc.Cell, ty UintN, decoder *Decoder) error {
	num, err := cell.ReadUint(ty.N)
	if err != nil {
		return err
	}
	*i = UInt64(num)
	return nil
}

func (b *BigUInt) UnmarshalTolk(cell *boc.Cell, ty UintN, decoder *Decoder) error {
	num, err := cell.ReadBigUint(ty.N)
	if err != nil {
		return err
	}
	*b = BigUInt(*num)
	return nil
}

func (u UintN) MarshalTolk(cell *boc.Cell, v *Value) error {
	if u.N > 64 {
		if v.bigInt == nil {
			return fmt.Errorf("big int not found")
		}
		bi := big.Int(*v.bigInt)

		err := cell.WriteBigUint(&bi, u.N)
		if err != nil {
			return err
		}
	} else {
		if v.smallUint == nil {
			return fmt.Errorf("small uint not found")
		}

		ui64 := uint64(*v.smallUint)
		err := cell.WriteUint(ui64, u.N)
		if err != nil {
			return err
		}
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
	v.sumType = "bigInt"
	return nil
}

func (vi VarIntN) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	n := vi.N
	ln, err := cell.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return err
	}
	err = vi.SetValue(v, BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (vi *VarInt) UnmarshalTolk(cell *boc.Cell, ty VarIntN, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(ty.N - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return err
	}
	*vi = VarInt(*val)
	return nil
}

func (vi VarIntN) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.bigInt == nil {
		return fmt.Errorf("BigInt is nil")
	}
	bi := big.Int(*v.bigInt)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), vi.N-1)
	if err != nil {
		return err
	}
	err = cell.WriteBigInt(&bi, len(num)*8)
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
	vu, ok := val.(BigInt)
	if !ok {
		return fmt.Errorf("value is not a BigInt")
	}
	v.bigInt = &vu
	v.sumType = "bigInt"
	return nil
}

func (vu VarUintN) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	n := vu.N
	ln, err := cell.ReadLimUint(n - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	err = vu.SetValue(v, BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (vu *VarUInt) UnmarshalTolk(cell *boc.Cell, ty VarUintN, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(ty.N - 1)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return err
	}
	*vu = VarUInt(*val)
	return nil
}

func (vu VarUintN) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.bigInt == nil {
		return fmt.Errorf("BigInt is nil")
	}
	bi := big.Int(*v.bigInt)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), vu.N-1)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(num)
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
	v.sumType = "bits"
	return nil
}

func (b BitsN) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	n := b.N
	val, err := cell.ReadBits(n)
	if err != nil {
		return err
	}
	err = b.SetValue(v, Bits(val))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bits) UnmarshalTolk(cell *boc.Cell, ty BitsN, decoder *Decoder) error {
	val, err := cell.ReadBits(ty.N)
	if err != nil {
		return err
	}
	*b = Bits(val)
	return nil
}

func (BitsN) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.bits == nil {
		return fmt.Errorf("bits is nil")
	}
	bi := boc.BitString(*v.bits)
	err := cell.WriteBitString(bi)
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
	v.sumType = "bigInt"
	return nil
}

func (c Coins) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(15)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	err = c.SetValue(v, BigInt(*val))
	if err != nil {
		return err
	}
	return nil
}

func (c *CoinsValue) UnmarshalTolk(cell *boc.Cell, ty Coins, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(15)
	if err != nil {
		return err
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return err
	}
	*c = CoinsValue(*val)
	return nil
}

func (Coins) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.bigInt == nil {
		return fmt.Errorf("BigInt is nil")
	}
	bi := big.Int(*v.bigInt)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), 15)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(num)
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
	v.sumType = "bool"
	return nil
}

func (b Bool) UnmarshalTolk(cell *boc.Cell, v *Value, decoder *Decoder) error {
	val, err := cell.ReadBit()
	if err != nil {
		return err
	}
	err = b.SetValue(v, BoolValue(val))
	if err != nil {
		return err
	}
	return nil
}

func (b *BoolValue) UnmarshalTolk(cell *boc.Cell, ty Bool, decoder *Decoder) error {
	val, err := cell.ReadBit()
	if err != nil {
		return err
	}
	*b = BoolValue(val)
	return nil
}

func (Bool) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.bool == nil {
		return fmt.Errorf("bool is nil")
	}

	err := cell.WriteBit(bool(*v.bool))
	if err != nil {
		return err
	}

	return nil
}

func (Bool) Equal(v Value, o Value) bool {
	return false
}
