package tolk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
		return fmt.Errorf("failed to skip internal address type and anycast: %w", err)
	}
	workchain, err := cell.ReadInt(8)
	if err != nil {
		return fmt.Errorf("failed to read internal address workchain: %w", err)
	}
	address, err := cell.ReadBytes(32)
	if err != nil {
		return fmt.Errorf("failed to read internal address hash: %w", err)
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
		return fmt.Errorf("failed to write internal address type and anycast: %w", err)
	}
	err = cell.WriteInt(int64(i.Workchain), 8)
	if err != nil {
		return fmt.Errorf("failed to write internal address workchain: %w", err)
	}
	err = cell.WriteBytes(i.Address[:])
	if err != nil {
		return fmt.Errorf("failed to write internal address hash: %w", err)
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
	return []byte(fmt.Sprintf("\"%v:%x\"", i.Workchain, i.Address)), nil
}

func (i *InternalAddress) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("invalid internal address format: %s", string(b))
	}
	addr := strings.Split(string(b[1:len(b)-1]), ":")
	if len(addr) != 2 {
		return fmt.Errorf("invalid internal address format: %s", string(b))
	}
	workchain, err := strconv.ParseInt(addr[0], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse internal address workchain: %w", err)
	}
	hexAddr, err := hex.DecodeString(addr[1])
	if err != nil {
		return fmt.Errorf("failed to parse internal address hash: %w", err)
	}
	i.Workchain = int8(workchain)
	i.Address = [32]byte(hexAddr)
	return nil
}

type NoneAddress struct {
}

func (n *NoneAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressOpt, decoder *Decoder) error {
	_, err := cell.ReadUint(2)
	if err != nil {
		return fmt.Errorf("failed to read none address type: %w", err)
	}

	return nil
}

func (n *NoneAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressOpt, encoder *Encoder) error {
	err := cell.WriteUint(0, 2) // none addr type ($00)
	if err != nil {
		return fmt.Errorf("failed to write none address type: %w", err)
	}
	return nil
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
		return o.InternalAddress.Equal(*otherOptionalAddress.InternalAddress)
	}
	return true
}

func (o *OptionalAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressOpt, decoder *Decoder) error {
	copyCell := cell.CopyRemaining()
	tag, err := copyCell.ReadUint(2)
	if err != nil {
		return fmt.Errorf("failed to read optional address type: %w", err)
	}
	if tag == 0 {
		o.SumType = "NoneAddress"
		o.NoneAddress = &NoneAddress{}
		err = o.NoneAddress.Unmarshal(cell, ty, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal optional address: %w", err)
		}
		return err
	}

	o.SumType = "InternalAddress"
	o.InternalAddress = &InternalAddress{}
	err = o.InternalAddress.Unmarshal(cell, tolkParser.Address{}, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal optional address: %w", err)
	}
	return nil
}

func (o *OptionalAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressOpt, encoder *Encoder) error {
	if o.SumType == "NoneAddress" {
		err := o.NoneAddress.Marshal(cell, ty, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address: %w", err)
		}
		return nil
	} else if o.SumType == "InternalAddress" {
		err := o.InternalAddress.Marshal(cell, tolkParser.Address{}, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address: %w", err)
		}
		return nil
	}
	return fmt.Errorf("unknown optional address SumType: %v", o.SumType)
}

func (o *OptionalAddress) MarshalJSON() ([]byte, error) {
	var optinalAddress = struct {
		SumType         SumType          `json:"sumType"`
		InternalAddress *InternalAddress `json:"internalAddress,omitempty"`
	}{}
	optinalAddress.SumType = o.SumType
	if o.SumType == "NoneAddress" {
		optinalAddress.SumType = "noneAddress"
	} else if o.SumType == "InternalAddress" {
		optinalAddress.SumType = "internalAddress"
		optinalAddress.InternalAddress = o.InternalAddress
	} else {
		return nil, fmt.Errorf("unknown optional address SumType: %v", o.SumType)
	}
	data, err := json.Marshal(optinalAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal optional address: %w", err)
	}
	return data, nil
}

func (o *OptionalAddress) UnmarshalJSON(b []byte) error {
	var optionalAddress = struct {
		SumType         SumType          `json:"sumType"`
		InternalAddress *InternalAddress `json:"internalAddress,omitempty"`
	}{}
	if err := json.Unmarshal(b, &optionalAddress); err != nil {
		return fmt.Errorf("failed to unmarshal optional address: %w", err)
	}

	if optionalAddress.SumType == "noneAddress" {
		o.SumType = "NoneAddress"
		o.NoneAddress = &NoneAddress{}
	} else if optionalAddress.SumType == "internalAddress" {
		o.SumType = "InternalAddress"
		o.InternalAddress = optionalAddress.InternalAddress
	} else {
		return fmt.Errorf("unknown optional address SumType: %v", o.SumType)
	}
	return nil
}

type ExternalAddress struct {
	Len     int16
	Address boc.BitString
}

func (e *ExternalAddress) Unmarshal(cell *boc.Cell, ty tolkParser.AddressExt, decoder *Decoder) error {
	err := cell.Skip(2)
	if err != nil {
		return fmt.Errorf("failed to skip external address type: %w", err)
	}
	ln, err := cell.ReadUint(9)
	if err != nil {
		return fmt.Errorf("failed to read external address length: %w", err)
	}
	bs, err := cell.ReadBits(int(ln))
	if err != nil {
		return fmt.Errorf("failed to read external address bytes: %w", err)
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
		return fmt.Errorf("failed to write external address type: %w", err)
	}
	err = cell.WriteUint(uint64(e.Len), 9)
	if err != nil {
		return fmt.Errorf("failed to write external address length: %w", err)
	}
	err = cell.WriteBitString(e.Address)
	if err != nil {
		return fmt.Errorf("failed to write external address bytes: %w", err)
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
	return []byte(fmt.Sprintf("\"%s\"", e.Address.ToFiftHex())), nil
}

func (e *ExternalAddress) UnmarshalJSON(b []byte) error {
	addr, err := boc.BitStringFromFiftHex(string(b[1 : len(b)-1]))
	if err != nil {
		return fmt.Errorf("invalid external address format: %v", string(b))
	}
	e.Len = int16(addr.BitsAvailableForRead())
	e.Address = *addr

	return nil
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
		return fmt.Errorf("failed to read any address type: %w", err)
	}
	switch tag {
	case 0:
		a.SumType = "NoneAddress"
		a.NoneAddress = &NoneAddress{}
		err = a.NoneAddress.Unmarshal(cell, tolkParser.AddressOpt{}, decoder)
	case 1:
		a.SumType = "ExternalAddress"
		a.ExternalAddress = &ExternalAddress{}
		err = a.ExternalAddress.Unmarshal(cell, tolkParser.AddressExt{}, decoder)
	case 2:
		a.SumType = "InternalAddress"
		a.InternalAddress = &InternalAddress{}
		err = a.InternalAddress.Unmarshal(cell, tolkParser.Address{}, decoder)
	case 3:
		a.SumType = "VarAddress"
		a.VarAddress = &VarAddress{}
		err = a.VarAddress.Unmarshal(cell, tolkParser.AddressExt{}, decoder)
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal any address: %w", err)
	}
	return nil
}

func (a *AnyAddress) Marshal(cell *boc.Cell, ty tolkParser.AddressAny, encoder *Encoder) error {
	var err error
	switch a.SumType {
	case "NoneAddress":
		err = a.NoneAddress.Marshal(cell, tolkParser.AddressOpt{}, encoder)
	case "InternalAddress":
		err = a.InternalAddress.Marshal(cell, tolkParser.Address{}, encoder)
	case "ExternalAddress":
		err = a.ExternalAddress.Marshal(cell, tolkParser.AddressExt{}, encoder)
	case "VarAddress":
		err = a.VarAddress.Marshal(cell, tolkParser.AddressAny{}, encoder)
	default:
		return fmt.Errorf("unknown any address SumType: %v", a.SumType)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal any address: %w", err)
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
		return a.InternalAddress.Equal(*otherAnyAddress.InternalAddress)
	case "ExternalAddress":
		return a.ExternalAddress.Equal(*otherAnyAddress.ExternalAddress)
	case "VarAddress":
		return a.VarAddress.Equal(*otherAnyAddress.VarAddress)
	}
	return false
}

func (a *AnyAddress) MarshalJSON() ([]byte, error) {
	var jsonAnyAddress = struct {
		SumType         string           `json:"sumType"`
		InternalAddress *InternalAddress `json:"internalAddress,omitempty"`
		ExternalAddress *ExternalAddress `json:"externalAddress,omitempty"`
		VarAddress      *VarAddress      `json:"varAddress,omitempty"`
	}{}
	switch a.SumType {
	case "NoneAddress":
		jsonAnyAddress.SumType = "noneAddress"
	case "InternalAddress":
		jsonAnyAddress.SumType = "internalAddress"
		jsonAnyAddress.InternalAddress = a.InternalAddress
	case "ExternalAddress":
		jsonAnyAddress.SumType = "externalAddress"
		jsonAnyAddress.ExternalAddress = a.ExternalAddress
	case "VarAddress":
		jsonAnyAddress.SumType = "varAddress"
		jsonAnyAddress.VarAddress = a.VarAddress
	default:
		return nil, fmt.Errorf("unknown any address SumType: %v", a.SumType)
	}
	data, err := json.Marshal(jsonAnyAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal any address: %w", err)
	}
	return data, nil
}

func (a *AnyAddress) UnmarshalJSON(b []byte) error {
	var anyAddress = struct {
		SumType         string           `json:"sumType"`
		InternalAddress *InternalAddress `json:"internalAddress,omitempty"`
		ExternalAddress *ExternalAddress `json:"externalAddress,omitempty"`
		VarAddress      *VarAddress      `json:"varAddress,omitempty"`
	}{}
	if err := json.Unmarshal(b, &anyAddress); err != nil {
		return fmt.Errorf("failed to unmarshal any address: %w", err)
	}
	switch anyAddress.SumType {
	case "noneAddress":
		a.SumType = "NoneAddress"
		a.NoneAddress = &NoneAddress{}
	case "internalAddress":
		a.SumType = "InternalAddress"
		a.InternalAddress = anyAddress.InternalAddress
	case "externalAddress":
		a.SumType = "ExternalAddress"
		a.ExternalAddress = anyAddress.ExternalAddress
	case "varAddress":
		a.SumType = "VarAddress"
		a.VarAddress = anyAddress.VarAddress
	default:
		return fmt.Errorf("unknown anyAddress SumType: %v", anyAddress.SumType)
	}
	return nil
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
		return fmt.Errorf("failed to skip var address type and anycast: %w", err)
	}
	ln, err := cell.ReadUint(9)
	if err != nil {
		return fmt.Errorf("failed to read var address length: %w", err)
	}
	workchain, err := cell.ReadInt(32)
	if err != nil {
		return fmt.Errorf("failed to read var address workchain: %w", err)
	}
	bs, err := cell.ReadBits(int(ln))
	if err != nil {
		return fmt.Errorf("failed to read var address bytes: %w", err)
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
		return fmt.Errorf("failed to write var address type and anycast: %w", err)
	}
	err = cell.WriteUint(uint64(va.Len), 9)
	if err != nil {
		return fmt.Errorf("failed to write var address length: %w", err)
	}
	err = cell.WriteInt(int64(va.Workchain), 32)
	if err != nil {
		return fmt.Errorf("failed to write var address workchain: %w", err)
	}
	err = cell.WriteBitString(va.Address)
	if err != nil {
		return fmt.Errorf("failed to write var address bytes: %w", err)
	}
	return nil
}

func (va *VarAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d:%s\"", va.Workchain, va.Address.ToFiftHex())), nil
}

func (va *VarAddress) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("invalid var address format: %v", string(b))
	}
	parts := bytes.Split(b[1:len(b)-1], []byte(":"))
	if len(parts) != 2 {
		return fmt.Errorf("invalid var address format: %v", string(b))
	}
	workchain, err := strconv.ParseInt(string(parts[0]), 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse var address workchain: %w", err)
	}
	bs, err := boc.BitStringFromFiftHex(string(parts[1]))
	if err != nil {
		return fmt.Errorf("failed to parse var address bytes: %w", err)
	}
	va.Workchain = int32(workchain)
	va.Len = int16(bs.BitsAvailableForRead())
	va.Address = *bs
	return nil
}
