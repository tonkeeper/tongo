package tolk

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

type InternalAddress struct {
	Workchain int8
	Address   [32]byte
}

func (i *InternalAddress) Unmarshal(cell *boc.Cell, ty tolkParser.Address, decoder *Decoder) error {
	err := cell.Skip(3) // skip addr type ($10) and anycast (0)
	if err != nil {
		return err
	}
	workchain, err := cell.ReadInt(8)
	if err != nil {
		return err
	}
	address, err := cell.ReadBytes(32)
	if err != nil {
		return err
	}
	*i = InternalAddress{
		Workchain: int8(workchain),
		Address:   [32]byte(address),
	}
	return nil
}

func (i *InternalAddress) Marshal(cell *boc.Cell, ty tolkParser.Address, encoder *Encoder) error {
	err := cell.WriteUint(0b100, 3) // internal addr type ($10) and anycast (0)
	if err != nil {
		return err
	}
	err = cell.WriteInt(int64(i.Workchain), 8)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(i.Address[:])
	if err != nil {
		return err
	}
	return nil
}

func (i *InternalAddress) Equal(other any) bool {
	otherInternalAddress, ok := other.(InternalAddress)
	if !ok {
		return false
	}
	return *i == otherInternalAddress
}

func (i *InternalAddress) ToRaw() string {
	return fmt.Sprintf("%v:%x", i.Workchain, i.Address)
}

func (i *InternalAddress) MarshalJSON() ([]byte, error) {
	return []byte(i.ToRaw()), nil
}

type NoneAddress struct {
}

func (n *NoneAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressOpt, decoder *Decoder) error {
	_, err := cell.ReadUint(2)
	if err != nil {
		return err
	}

	return nil
}

func (n *NoneAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressOpt, encoder *Encoder) error {
	err := cell.WriteUint(0, 2) // none addr type ($00)
	if err != nil {
		return err
	}
	return nil
}

func (n *NoneAddress) MarshalJSON() ([]byte, error) {
	return []byte(""), nil
}

type OptionalAddress struct {
	SumType
	NoneAddress     *NoneAddress
	InternalAddress *InternalAddress
}

func (o *OptionalAddress) Equal(other any) bool {
	otherOptionalAddress, ok := other.(OptionalAddress)
	if !ok {
		return false
	}
	if o.SumType != otherOptionalAddress.SumType {
		return false
	}
	if o.SumType == "InternalAddress" {
		return o.InternalAddress.Equal(otherOptionalAddress.InternalAddress)
	}
	return true
}

func (o *OptionalAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressOpt, decoder *Decoder) error {
	copyCell := cell.CopyRemaining()
	tag, err := copyCell.ReadUint(2)
	if err != nil {
		return err
	}
	if tag == 0 {
		o.SumType = "NoneAddress"
		o.NoneAddress = &NoneAddress{}
		return o.NoneAddress.Unmarshal(cell, ty, decoder)
	}

	o.SumType = "InternalAddress"
	o.InternalAddress = &InternalAddress{}
	return o.InternalAddress.Unmarshal(cell, tolkParser.Address{}, decoder)
}

func (o *OptionalAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressOpt, encoder *Encoder) error {
	if o.SumType == "NoneAddress" {
		return o.NoneAddress.Marshal(cell, ty, encoder)
	} else if o.SumType == "InternalAddress" {
		return o.InternalAddress.Marshal(cell, tolkParser.Address{}, encoder)
	}
	return fmt.Errorf("unknown any address SumType: %v", o.SumType)
}

func (o *OptionalAddress) MarshalJSON() ([]byte, error) {
	if o.SumType == "NoneAddress" {
		return json.Marshal(o.NoneAddress)
	} else if o.SumType == "InternalAddress" {
		return json.Marshal(o.InternalAddress)
	}
	return nil, fmt.Errorf("unknown any address SumType: %v", o.SumType)
}

type ExternalAddress struct {
	Len     int16
	Address boc.BitString
}

func (e *ExternalAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressExt, decoder *Decoder) error {
	err := cell.Skip(2)
	if err != nil {
		return err
	}
	ln, err := cell.ReadUint(9)
	if err != nil {
		return err
	}
	bs, err := cell.ReadBits(int(ln))
	if err != nil {
		return err
	}
	*e = ExternalAddress{
		Len:     int16(ln),
		Address: bs,
	}

	return nil
}

func (e *ExternalAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressExt, encoder *Encoder) error {
	err := cell.WriteUint(1, 2) // external addr type ($01)
	if err != nil {
		return err
	}
	err = cell.WriteUint(uint64(e.Len), 9)
	if err != nil {
		return err
	}
	err = cell.WriteBitString(e.Address)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExternalAddress) Equal(other any) bool {
	otherExternalAddress, ok := other.(ExternalAddress)
	if !ok {
		return false
	}
	if e.Len != otherExternalAddress.Len {
		return false
	}
	return bytes.Equal(e.Address.Buffer(), otherExternalAddress.Address.Buffer())
}

func (e *ExternalAddress) MarshalJSON() ([]byte, error) {
	return []byte(e.Address.ToFiftHex()), nil
}

type AnyAddress struct {
	SumType
	InternalAddress *InternalAddress
	NoneAddress     *NoneAddress
	ExternalAddress *ExternalAddress
	VarAddress      *VarAddress
}

func (a *AnyAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressAny, decoder *Decoder) error {
	copyCell := cell.CopyRemaining()
	tag, err := copyCell.ReadUint(2)
	if err != nil {
		return err
	}
	switch tag {
	case 0:
		a.SumType = "NoneAddress"
		a.NoneAddress = &NoneAddress{}
		return a.NoneAddress.Unmarshal(cell, tolkParser.AddressOpt{}, decoder)
	case 1:
		a.SumType = "ExternalAddress"
		a.ExternalAddress = &ExternalAddress{}
		return a.ExternalAddress.Unmarshal(cell, tolkParser.AddressExt{}, decoder)
	case 2:
		a.SumType = "InternalAddress"
		a.InternalAddress = &InternalAddress{}
		return a.InternalAddress.Unmarshal(cell, tolkParser.Address{}, decoder)
	case 3:
		a.SumType = "VarAddress"
		a.VarAddress = &VarAddress{}
		return a.VarAddress.Unmarshal(cell, tolkParser.AddressExt{}, decoder)
	}

	return nil
}

func (a *AnyAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressAny, encoder *Encoder) error {
	switch a.SumType {
	case "NoneAddress":
		return a.NoneAddress.Marshal(cell, tolkParser.AddressOpt{}, encoder)
	case "InternalAddress":
		return a.InternalAddress.Marshal(cell, tolkParser.Address{}, encoder)
	case "ExternalAddress":
		return a.ExternalAddress.Marshal(cell, tolkParser.AddressExt{}, encoder)
	case "VarAddress":
		return a.VarAddress.Marshal(cell, tolkParser.AddressAny{}, encoder)
	}

	return nil
}

func (a *AnyAddress) Equal(other any) bool {
	otherAnyAddress, ok := other.(AnyAddress)
	if !ok {
		return false
	}
	if otherAnyAddress.SumType != a.SumType {
		return false
	}
	switch a.SumType {
	case "NoneAddress":
		return true
	case "InternalAddress":
		return a.InternalAddress.Equal(otherAnyAddress.InternalAddress)
	case "ExternalAddress":
		return a.ExternalAddress.Equal(otherAnyAddress.ExternalAddress)
	case "VarAddress":
		return a.VarAddress.Equal(otherAnyAddress.VarAddress)
	}
	return false
}

func (a *AnyAddress) MarshalJSON() ([]byte, error) {
	switch a.SumType {
	case "NoneAddress":
		return json.Marshal(a.NoneAddress)
	case "InternalAddress":
		return json.Marshal(a.InternalAddress)
	case "ExternalAddress":
		return json.Marshal(a.ExternalAddress)
	case "VarAddress":
		return json.Marshal(a.VarAddress)
	}
	return nil, fmt.Errorf("unknown any address SumType: %v", a.SumType)
}

type VarAddress struct {
	Len       int16
	Workchain int32
	Address   boc.BitString
}

func (va *VarAddress) Equal(other any) bool {
	otherVarAddress, ok := other.(VarAddress)
	if !ok {
		return false
	}
	if va.Len != otherVarAddress.Len {
		return false
	}
	if va.Workchain != otherVarAddress.Workchain {
		return false
	}
	return bytes.Equal(va.Address.Buffer(), otherVarAddress.Address.Buffer())
}

func (va *VarAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressExt, decoder *Decoder) error {
	err := cell.Skip(3) // skip var type ($11) and anycast (0)
	if err != nil {
		return err
	}
	ln, err := cell.ReadUint(9)
	if err != nil {
		return err
	}
	workchain, err := cell.ReadInt(32)
	if err != nil {
		return err
	}
	bs, err := cell.ReadBits(int(ln))
	if err != nil {
		return err
	}
	*va = VarAddress{
		Len:       int16(ln),
		Workchain: int32(workchain),
		Address:   bs,
	}

	return nil
}

func (va *VarAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressAny, encoder *Encoder) error {
	err := cell.WriteUint(0b110, 3) // var addr type ($11) and anycast (0)
	if err != nil {
		return err
	}
	err = cell.WriteUint(uint64(va.Len), 9)
	if err != nil {
		return err
	}
	err = cell.WriteInt(int64(va.Workchain), 32)
	if err != nil {
		return err
	}
	err = cell.WriteBitString(va.Address)
	if err != nil {
		return err
	}
	return nil
}

func (va *VarAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d:%s", va.Workchain, va.Address.ToFiftHex())), nil
}
