package tolk

import (
	"bytes"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Int64 int64

func (i *Int64) Unmarshal(cell *boc.Cell, ty tolkParser.IntN, decoder *Decoder) error {
	num, err := cell.ReadInt(ty.N)
	if err != nil {
		return err
	}
	*i = Int64(num)
	return nil
}

func (i *Int64) Marshal(cell *boc.Cell, ty tolkParser.IntN, encoder *Encoder) error {
	return cell.WriteInt(int64(*i), ty.N)
}

func (i *Int64) Equal(other any) bool {
	otherInt, ok := other.(Int64)
	if !ok {
		return false
	}
	return *i == otherInt
}

type BigInt big.Int

func (b *BigInt) Unmarshal(cell *boc.Cell, ty tolkParser.IntN, decoder *Decoder) error {
	num, err := cell.ReadBigInt(ty.N)
	if err != nil {
		return err
	}
	*b = BigInt(*num)
	return nil
}

func (b *BigInt) Marshal(cell *boc.Cell, ty tolkParser.IntN, encoder *Encoder) error {
	bi := big.Int(*b)
	return cell.WriteBigInt(&bi, ty.N)
}

func (b *BigInt) Equal(other any) bool {
	otherBigInt, ok := other.(BigInt)
	if !ok {
		return false
	}
	bi := big.Int(*b)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type UInt64 uint64

func (i *UInt64) Unmarshal(cell *boc.Cell, ty tolkParser.UintN, decoder *Decoder) error {
	num, err := cell.ReadUint(ty.N)
	if err != nil {
		return err
	}
	*i = UInt64(num)
	return nil
}

func (i *UInt64) Marshal(cell *boc.Cell, ty tolkParser.UintN, encoder *Encoder) error {
	return cell.WriteUint(uint64(*i), ty.N)
}

func (i *UInt64) Equal(other any) bool {
	otherUint, ok := other.(UInt64)
	if !ok {
		return false
	}
	return *i == otherUint
}

type BigUInt big.Int

func (b *BigUInt) Unmarshal(cell *boc.Cell, ty tolkParser.UintN, decoder *Decoder) error {
	num, err := cell.ReadBigUint(ty.N)
	if err != nil {
		return err
	}
	*b = BigUInt(*num)
	return nil
}

func (b *BigUInt) Marshal(cell *boc.Cell, ty tolkParser.UintN, encoder *Encoder) error {
	bi := big.Int(*b)
	return cell.WriteBigUint(&bi, ty.N)
}

func (b *BigUInt) Equal(other any) bool {
	otherBigInt, ok := other.(BigUInt)
	if !ok {
		return false
	}
	bi := big.Int(*b)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type VarInt big.Int

func (vi *VarInt) Unmarshal(cell *boc.Cell, ty tolkParser.VarIntN, decoder *Decoder) error {
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

func (vi *VarInt) Marshal(cell *boc.Cell, ty tolkParser.VarIntN, encoder *Encoder) error {
	bi := big.Int(*vi)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), ty.N-1)
	if err != nil {
		return err
	}
	err = cell.WriteBigInt(&bi, len(num)*8)
	if err != nil {
		return err
	}

	return nil
}

func (vi *VarInt) Equal(other any) bool {
	otherBigInt, ok := other.(VarInt)
	if !ok {
		return false
	}
	bi := big.Int(*vi)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type VarUInt big.Int

func (vu *VarUInt) Unmarshal(cell *boc.Cell, ty tolkParser.VarUintN, decoder *Decoder) error {
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

func (vu *VarUInt) Marshal(cell *boc.Cell, ty tolkParser.VarUintN, encoder *Encoder) error {
	bi := big.Int(*vu)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), ty.N-1)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(num)
	if err != nil {
		return err
	}

	return nil
}

func (vu *VarUInt) Equal(other any) bool {
	otherBigInt, ok := other.(VarUInt)
	if !ok {
		return false
	}
	bi := big.Int(*vu)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type Bits boc.BitString

func (b *Bits) Unmarshal(cell *boc.Cell, ty tolkParser.BitsN, decoder *Decoder) error {
	val, err := cell.ReadBits(ty.N)
	if err != nil {
		return err
	}
	*b = Bits(val)
	return nil
}

func (b *Bits) Marshal(cell *boc.Cell, ty tolkParser.BitsN, encoder *Encoder) error {
	bi := boc.BitString(*b)
	err := cell.WriteBitString(bi)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bits) Equal(other any) bool {
	otherBits, ok := other.(Bits)
	if !ok {
		return false
	}
	bs := boc.BitString(*b)
	otherBs := boc.BitString(otherBits)
	return bytes.Equal(bs.Buffer(), otherBs.Buffer())
}

type CoinsValue big.Int

func (c *CoinsValue) Unmarshal(cell *boc.Cell, ty tolkParser.Coins, decoder *Decoder) error {
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

func (c *CoinsValue) Marshal(cell *boc.Cell, ty tolkParser.Coins, encoder *Encoder) error {
	varInt := VarInt(*c)
	return varInt.Marshal(cell, tolkParser.VarIntN{N: 16}, encoder) // coins is actually varint16
}

func (c *CoinsValue) Equal(other any) bool {
	otherBigInt, ok := other.(CoinsValue)
	if !ok {
		return false
	}
	bi := big.Int(*c)
	otherBi := big.Int(otherBigInt)
	return bi.Cmp(&otherBi) == 0
}

type BoolValue bool

func (b *BoolValue) Equal(o any) bool {
	otherBool, ok := o.(BoolValue)
	if !ok {
		return false
	}
	return *b == otherBool
}

func (b *BoolValue) Unmarshal(cell *boc.Cell, ty tolkParser.Bool, decoder *Decoder) error {
	val, err := cell.ReadBit()
	if err != nil {
		return err
	}
	*b = BoolValue(val)
	return nil
}

func (b *BoolValue) Marshal(cell *boc.Cell, ty tolkParser.Bool, encoder *Encoder) error {
	err := cell.WriteBit(bool(*b))
	if err != nil {
		return err
	}

	return nil
}
