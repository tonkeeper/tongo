package runtime

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
)

type InternalAddress struct {
	Workchain int8
	Address   [32]byte
}

func InternalAddressFromTLB(accID ton.AccountID) InternalAddress {
	return InternalAddress{
		Workchain: int8(accID.Workchain),
		Address:   accID.Address,
	}
}

func (i *InternalAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
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

func (i *InternalAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
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

func (i InternalAddress) MarshalJSON() ([]byte, error) {
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

func (n *NoneAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	_, err := cell.ReadUint(2)
	if err != nil {
		return fmt.Errorf("failed to read none address type: %w", err)
	}

	return nil
}

func (n *NoneAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
	err := cell.WriteUint(0, 2) // none addr type ($00)
	if err != nil {
		return fmt.Errorf("failed to write none address type: %w", err)
	}
	return nil
}

func (n NoneAddress) MarshalJSON() ([]byte, error) {
	return []byte("\"\""), nil
}

func (n *NoneAddress) UnmarshalJSON(b []byte) error {
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
	if o.SumType == SumTypeInternalAddress {
		return o.InternalAddress.Equal(*otherOptionalAddress.InternalAddress)
	}
	return true
}

func (o *OptionalAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	copyCell := cell.CopyRemaining()
	tag, err := copyCell.ReadUint(2)
	if err != nil {
		return fmt.Errorf("failed to read optional address type: %w", err)
	}
	if tag == 0 {
		o.SumType = SumTypeNoneAddress
		o.NoneAddress = &NoneAddress{}
		err = o.NoneAddress.Unmarshal(cell, decoder)
		if err != nil {
			return fmt.Errorf("failed to unmarshal optional address: %w", err)
		}
		return err
	}

	o.SumType = SumTypeInternalAddress
	o.InternalAddress = &InternalAddress{}
	err = o.InternalAddress.Unmarshal(cell, decoder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal optional address: %w", err)
	}
	return nil
}

func (o *OptionalAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
	if o.SumType == SumTypeNoneAddress {
		err := o.NoneAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address: %w", err)
		}
		return nil
	} else if o.SumType == SumTypeInternalAddress {
		err := o.InternalAddress.Marshal(cell, encoder)
		if err != nil {
			return fmt.Errorf("failed to marshal optional address: %w", err)
		}
		return nil
	}
	return fmt.Errorf("unknown optional address SumType: %v", o.SumType)
}

func (o OptionalAddress) MarshalJSON() ([]byte, error) {
	if o.SumType == SumTypeNoneAddress {
		return json.Marshal(o.NoneAddress)
	} else if o.SumType == SumTypeInternalAddress {
		return json.Marshal(o.InternalAddress)
	} else {
		return nil, fmt.Errorf("unknown optional address SumType: %v", o.SumType)
	}
}

func (o *OptionalAddress) UnmarshalJSON(b []byte) error {
	var internalAddress *InternalAddress
	if err := json.Unmarshal(b, &internalAddress); err != nil {
		o.SumType = SumTypeNoneAddress
		o.NoneAddress = &NoneAddress{}
		return nil
	}

	o.SumType = SumTypeInternalAddress
	o.InternalAddress = internalAddress
	return nil
}

type ExternalAddress struct {
	Len     int16
	Address boc.BitString
}

func (e *ExternalAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
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

func (e *ExternalAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
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

func (e ExternalAddress) MarshalJSON() ([]byte, error) {
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

func (a *AnyAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
	copyCell := cell.CopyRemaining()
	tag, err := copyCell.ReadUint(2)
	if err != nil {
		return fmt.Errorf("failed to read any address type: %w", err)
	}
	switch tag {
	case 0:
		a.SumType = SumTypeNoneAddress
		a.NoneAddress = &NoneAddress{}
		err = a.NoneAddress.Unmarshal(cell, decoder)
	case 1:
		a.SumType = SumTypeExternalAddress
		a.ExternalAddress = &ExternalAddress{}
		err = a.ExternalAddress.Unmarshal(cell, decoder)
	case 2:
		a.SumType = SumTypeInternalAddress
		a.InternalAddress = &InternalAddress{}
		err = a.InternalAddress.Unmarshal(cell, decoder)
	case 3:
		a.SumType = SumTypeVarAddress
		a.VarAddress = &VarAddress{}
		err = a.VarAddress.Unmarshal(cell, decoder)
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal any address: %w", err)
	}
	return nil
}

func (a *AnyAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
	var err error
	switch a.SumType {
	case SumTypeNoneAddress:
		err = a.NoneAddress.Marshal(cell, encoder)
	case SumTypeInternalAddress:
		err = a.InternalAddress.Marshal(cell, encoder)
	case SumTypeExternalAddress:
		err = a.ExternalAddress.Marshal(cell, encoder)
	case SumTypeVarAddress:
		err = a.VarAddress.Marshal(cell, encoder)
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
	case SumTypeNoneAddress:
		return true
	case SumTypeInternalAddress:
		return a.InternalAddress.Equal(*otherAnyAddress.InternalAddress)
	case SumTypeExternalAddress:
		return a.ExternalAddress.Equal(*otherAnyAddress.ExternalAddress)
	case SumTypeVarAddress:
		return a.VarAddress.Equal(*otherAnyAddress.VarAddress)
	}
	return false
}

func (a AnyAddress) MarshalJSON() ([]byte, error) {
	var data []byte
	var err error
	switch a.SumType {
	case SumTypeNoneAddress:
		data, err = json.Marshal(a.NoneAddress)
	case SumTypeInternalAddress:
		data, err = json.Marshal(a.InternalAddress)
	case SumTypeExternalAddress:
		data, err = json.Marshal(a.ExternalAddress)
	case SumTypeVarAddress:
		data, err = json.Marshal(a.VarAddress)
	default:
		return nil, fmt.Errorf("unknown any address SumType: %v", a.SumType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AnyAddress: %w", err)
	}
	return data, nil
}

func (a *AnyAddress) UnmarshalJSON(b []byte) error {
	var internalAddress *InternalAddress
	if err := json.Unmarshal(b, &internalAddress); err == nil {
		a.SumType = SumTypeInternalAddress
		a.InternalAddress = internalAddress
		return nil
	}

	var externalAddress *ExternalAddress
	if err := json.Unmarshal(b, &externalAddress); err == nil {
		a.SumType = SumTypeExternalAddress
		a.ExternalAddress = externalAddress
	}

	var varAddress *VarAddress
	if err := json.Unmarshal(b, &varAddress); err == nil {
		a.SumType = SumTypeVarAddress
		a.VarAddress = varAddress
	}

	a.SumType = SumTypeNoneAddress
	a.NoneAddress = &NoneAddress{}
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

func (va *VarAddress) Unmarshal(cell *boc.Cell, decoder *Decoder) error {
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

func (va *VarAddress) Marshal(cell *boc.Cell, encoder *Encoder) error {
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

func (va VarAddress) MarshalJSON() ([]byte, error) {
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
