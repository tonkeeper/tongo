package tolk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type Int64 int64

func (i *Int64) Unmarshal(cell *boc.Cell, ty tolkParser.IntN, decoder *Decoder) error {
	num, err := cell.ReadInt(ty.N)
	if err != nil {
		return fmt.Errorf("failed to read %v-bit integer: %w", ty.N, err)
	}
	*i = Int64(num)
	return nil
}

func (i *Int64) Marshal(cell *boc.Cell, ty tolkParser.IntN, encoder *Encoder) error {
	err := cell.WriteInt(int64(*i), ty.N)
	if err != nil {
		return fmt.Errorf("failed to write %v-bit integer: %w", ty.N, err)
	}
	return nil
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
		return fmt.Errorf("failed to read %v-bit big integer: %w", ty.N, err)
	}
	*b = BigInt(*num)
	return nil
}

func (b *BigInt) Marshal(cell *boc.Cell, ty tolkParser.IntN, encoder *Encoder) error {
	bi := big.Int(*b)
	err := cell.WriteBigInt(&bi, ty.N)
	if err != nil {
		return fmt.Errorf("failed to write %v-bit big integer: %w", ty.N, err)
	}
	return nil
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

func (b *BigInt) MarshalJSON() ([]byte, error) {
	bi := big.Int(*b)
	return []byte(`"` + bi.String() + `"`), nil
}

func (b *BigInt) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid big int format: %s", string(bytes))
	}
	bi, ok := new(big.Int).SetString(string(bytes[1:len(bytes)-1]), 10)
	if !ok {
		return fmt.Errorf("failed to parse big integer from %v", string(bytes[1:len(bytes)-1]))
	}
	*b = BigInt(*bi)
	return nil
}

type UInt64 uint64

func (i *UInt64) Unmarshal(cell *boc.Cell, ty tolkParser.UintN, decoder *Decoder) error {
	num, err := cell.ReadUint(ty.N)
	if err != nil {
		return fmt.Errorf("failed to read %v-bit unsigned integer: %w", ty.N, err)
	}
	*i = UInt64(num)
	return nil
}

func (i *UInt64) Marshal(cell *boc.Cell, ty tolkParser.UintN, encoder *Encoder) error {
	err := cell.WriteUint(uint64(*i), ty.N)
	if err != nil {
		return fmt.Errorf("failed to write %v-bit unsigned integer: %w", ty.N, err)
	}
	return nil
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
		return fmt.Errorf("failed to read %v-bit usigned big integer: %w", ty.N, err)
	}
	*b = BigUInt(*num)
	return nil
}

func (b *BigUInt) Marshal(cell *boc.Cell, ty tolkParser.UintN, encoder *Encoder) error {
	bi := big.Int(*b)
	err := cell.WriteBigUint(&bi, ty.N)
	if err != nil {
		return fmt.Errorf("failed to write %v-bit usigned big integer: %w", ty.N, err)
	}
	return nil
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

func (b *BigUInt) MarshalJSON() ([]byte, error) {
	bi := big.Int(*b)
	return []byte(`"` + bi.String() + `"`), nil
}

func (b *BigUInt) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid usigned big int format: %s", string(bytes))
	}
	bi, ok := new(big.Int).SetString(string(bytes[1:len(bytes)-1]), 10)
	if !ok {
		return fmt.Errorf("failed to parse usigned big interger from %v", string(bytes[1:len(bytes)-1]))
	}
	*b = BigUInt(*bi)
	return nil
}

type VarInt big.Int

func (vi *VarInt) Unmarshal(cell *boc.Cell, ty tolkParser.VarIntN, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(ty.N - 1)
	if err != nil {
		return fmt.Errorf("failed to read var integer length: %w", err)
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return fmt.Errorf("failed to read var integer value: %w", err)
	}
	*vi = VarInt(*val)
	return nil
}

func (vi *VarInt) Marshal(cell *boc.Cell, ty tolkParser.VarIntN, encoder *Encoder) error {
	bi := big.Int(*vi)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), ty.N-1)
	if err != nil {
		return fmt.Errorf("failed to write var integer length: %w", err)
	}
	err = cell.WriteBigInt(&bi, len(num)*8)
	if err != nil {
		return fmt.Errorf("failed to write var integer value: %w", err)
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

func (vi *VarInt) MarshalJSON() ([]byte, error) {
	bi := big.Int(*vi)
	return []byte(`"` + bi.String() + `"`), nil
}

func (vi *VarInt) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid var integer format: %s", string(bytes))
	}
	bi, ok := new(big.Int).SetString(string(bytes[1:len(bytes)-1]), 10)
	if !ok {
		return fmt.Errorf("failed to parse var integer from %v", string(bytes[1:len(bytes)-1]))
	}
	*vi = VarInt(*bi)
	return nil
}

type VarUInt big.Int

func (vu *VarUInt) Unmarshal(cell *boc.Cell, ty tolkParser.VarUintN, decoder *Decoder) error {
	ln, err := cell.ReadLimUint(ty.N - 1)
	if err != nil {
		return fmt.Errorf("failed to read usigned var integer length: %w", err)
	}
	val, err := cell.ReadBigInt(int(ln) * 8)
	if err != nil {
		return fmt.Errorf("failed to read usigned var integer value: %w", err)
	}
	*vu = VarUInt(*val)
	return nil
}

func (vu *VarUInt) Marshal(cell *boc.Cell, ty tolkParser.VarUintN, encoder *Encoder) error {
	bi := big.Int(*vu)
	num := bi.Bytes()
	err := cell.WriteLimUint(len(num), ty.N-1)
	if err != nil {
		return fmt.Errorf("failed to write usigned var integer length: %w", err)
	}
	err = cell.WriteBytes(num)
	if err != nil {
		return fmt.Errorf("failed to write usigned var integer value: %w", err)
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

func (vu *VarUInt) MarshalJSON() ([]byte, error) {
	bi := big.Int(*vu)
	return []byte(`"` + bi.String() + `"`), nil
}

func (vu *VarUInt) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid var integer format: %s", string(bytes))
	}
	bi, ok := new(big.Int).SetString(string(bytes[1:len(bytes)-1]), 10)
	if !ok {
		return fmt.Errorf("failed to parse usigned var integer from %v", string(bytes[1:len(bytes)-1]))
	}
	*vu = VarUInt(*bi)
	return nil
}

type Bits boc.BitString

func (b *Bits) Unmarshal(cell *boc.Cell, ty tolkParser.BitsN, decoder *Decoder) error {
	val, err := cell.ReadBits(ty.N)
	if err != nil {
		return fmt.Errorf("failed to read bits value: %w", err)
	}
	*b = Bits(val)
	return nil
}

func (b *Bits) Marshal(cell *boc.Cell, ty tolkParser.BitsN, encoder *Encoder) error {
	bi := boc.BitString(*b)
	err := cell.WriteBitString(bi)
	if err != nil {
		return fmt.Errorf("failed to write bits value: %w", err)
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

func (b *Bits) MarshalJSON() ([]byte, error) {
	data, err := boc.BitString(*b).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bits value: %w", err)
	}
	return data, nil
}

func (b *Bits) UnmarshalJSON(bytes []byte) error {
	bs := boc.BitString{}
	if err := json.Unmarshal(bytes, &bs); err != nil {
		return fmt.Errorf("failed to unmarshal bits value: %w", err)
	}
	*b = Bits(bs)
	return nil
}

type CoinsValue big.Int

func (c *CoinsValue) Unmarshal(cell *boc.Cell, ty tolkParser.Coins, decoder *Decoder) error {
	varUint := VarUInt{}
	if err := varUint.Unmarshal(cell, tolkParser.VarUintN{N: 16}, decoder); err != nil {
		return fmt.Errorf("failed to unmarshal coins value: %w", err)
	}
	*c = CoinsValue(varUint)
	return nil
}

func (c *CoinsValue) Marshal(cell *boc.Cell, ty tolkParser.Coins, encoder *Encoder) error {
	varUint := VarUInt(*c)
	err := varUint.Marshal(cell, tolkParser.VarUintN{N: 16}, encoder) // coins is actually varuint16
	if err != nil {
		return fmt.Errorf("failed to marshal coins value: %w", err)
	}
	return nil
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

func (c *CoinsValue) MarshalJSON() ([]byte, error) {
	bi := big.Int(*c)
	return []byte(`"` + bi.String() + `"`), nil
}

func (c *CoinsValue) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 {
		return fmt.Errorf("invalid coins format: %s", string(bytes))
	}
	bi, ok := new(big.Int).SetString(string(bytes[1:len(bytes)-1]), 10)
	if !ok {
		return fmt.Errorf("failed to parse coins from %v", string(bytes[1:len(bytes)-1]))
	}
	*c = CoinsValue(*bi)
	return nil
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
		return fmt.Errorf("failed to read bool value: %w", err)
	}
	*b = BoolValue(val)
	return nil
}

func (b *BoolValue) Marshal(cell *boc.Cell, ty tolkParser.Bool, encoder *Encoder) error {
	err := cell.WriteBit(bool(*b))
	if err != nil {
		return fmt.Errorf("failed to write bool value: %w", err)
	}
	return nil
}
