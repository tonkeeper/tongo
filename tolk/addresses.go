package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

type Address struct{}

func (Address) SetValue(v *Value, val any) error {
	a, ok := val.(InternalAddress)
	if !ok {
		return fmt.Errorf("value is not an internal address")
	}
	v.sumType = "internalAddress"
	v.internalAddress = &a
	return nil
}

func (a Address) UnmarshalTolk(cell *boc.Cell, v *Value, abiCtx *Decoder) error {
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
	err = a.SetValue(v, InternalAddress{
		Workchain: int8(workchain),
		Address:   [32]byte(address),
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *InternalAddress) UnmarshalTolk(cell *boc.Cell, ty Address, abiCtx *Decoder) error {
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

func (Address) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.internalAddress == nil {
		return fmt.Errorf("address not found")
	}

	err := cell.WriteUint(0b100, 3) // internal addr type ($10) and anycast (0)
	if err != nil {
		return err
	}
	err = cell.WriteInt(int64(v.internalAddress.Workchain), 8)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(v.internalAddress.Address[:])
	if err != nil {
		return err
	}
	return nil
}

func (Address) Equal(v Value, o Value) bool {
	if v.internalAddress == nil || o.internalAddress == nil {
		return false
	}
	vi := *v.internalAddress
	oi := *o.internalAddress
	return vi.Workchain == oi.Workchain && vi.Address == oi.Address
}

func (Address) GetFixedSize() int {
	return 267
}

type AddressOpt struct {
}

func (AddressOpt) SetValue(v *Value, val any) error {
	a, ok := val.(OptionalAddress)
	if !ok {
		return fmt.Errorf("value is not an optional address")
	}
	v.sumType = "optionalAddress"
	v.optionalAddress = &a
	return nil
}

func (a AddressOpt) UnmarshalTolk(cell *boc.Cell, v *Value, abiCtx *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	if tag == 0 {
		err = a.SetValue(v, OptionalAddress{
			SumType: "NoneAddress",
		})
		if err != nil {
			return err
		}
		return nil
	}
	err = cell.Skip(1) // skip anycast (0)
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
	err = a.SetValue(v, OptionalAddress{
		SumType: "InternalAddress",
		InternalAddress: InternalAddress{
			Workchain: int8(workchain),
			Address:   [32]byte(address),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OptionalAddress) UnmarshalTolk(cell *boc.Cell, ty AddressOpt, abiCtx *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	if tag == 0 {
		o.SumType = "NoneAddress"
		return nil
	}
	err = cell.Skip(1) // skip anycast (0)
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
	o.SumType = "InternalAddress"
	o.InternalAddress = InternalAddress{
		Workchain: int8(workchain),
		Address:   [32]byte(address),
	}

	return nil
}

func (AddressOpt) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.optionalAddress == nil {
		return fmt.Errorf("optional address not found")
	}

	if v.optionalAddress.SumType == "NoneAddress" {
		err := cell.WriteUint(0, 2)
		if err != nil {
			return err
		}

		return nil
	}

	err := cell.WriteUint(0b100, 3) // internal addr type ($10) and anycast (0)
	if err != nil {
		return err
	}
	err = cell.WriteInt(int64(v.optionalAddress.InternalAddress.Workchain), 8)
	if err != nil {
		return err
	}
	err = cell.WriteBytes(v.optionalAddress.InternalAddress.Address[:])
	if err != nil {
		return err
	}
	return nil
}

func (AddressOpt) Equal(v Value, o Value) bool {
	return false
}

type AddressExt struct{}

func (AddressExt) SetValue(v *Value, val any) error {
	a, ok := val.(ExternalAddress)
	if !ok {
		return fmt.Errorf("value is not an external address")
	}
	v.externalAddress = &a
	v.sumType = "externalAddress"
	return nil
}

func (a AddressExt) UnmarshalTolk(cell *boc.Cell, v *Value, abiCtx *Decoder) error {
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
	err = a.SetValue(v, ExternalAddress{
		Len:     int16(ln),
		Address: bs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *ExternalAddress) UnmarshalTolk(cell *boc.Cell, ty AddressExt, abiCtx *Decoder) error {
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

func (AddressExt) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.externalAddress == nil {
		return fmt.Errorf("external address not found")
	}

	err := cell.WriteUint(1, 2) // external addr type ($01)
	if err != nil {
		return err
	}
	err = cell.WriteUint(uint64(v.externalAddress.Len), 9)
	if err != nil {
		return err
	}
	err = cell.WriteBitString(v.externalAddress.Address)
	if err != nil {
		return err
	}
	return nil
}

func (AddressExt) Equal(v Value, o Value) bool {
	return false
}

type AddressAny struct{}

func (AddressAny) SetValue(v *Value, val any) error {
	a, ok := val.(AnyAddress)
	if !ok {
		return fmt.Errorf("value is not an any address")
	}
	v.anyAddress = &a
	v.sumType = "anyAddress"
	return nil
}

func (a AddressAny) UnmarshalTolk(cell *boc.Cell, v *Value, abiCtx *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	switch tag {
	case 0:
		err = a.SetValue(v, AnyAddress{
			SumType: "NoneAddress",
		})
		if err != nil {
			return err
		}
	case 1:
		ln, err := cell.ReadUint(9)
		if err != nil {
			return err
		}
		bs, err := cell.ReadBits(int(ln))
		if err != nil {
			return err
		}
		err = a.SetValue(v, AnyAddress{
			SumType: "ExternalAddress",
			ExternalAddress: &ExternalAddress{
				Len:     int16(ln),
				Address: bs,
			},
		})
		if err != nil {
			return err
		}
	case 2:
		err = cell.Skip(1) // skip anycast (0)
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
		err = a.SetValue(v, AnyAddress{
			SumType: "InternalAddress",
			InternalAddress: &InternalAddress{
				Workchain: int8(workchain),
				Address:   [32]byte(address),
			},
		})
		if err != nil {
			return err
		}
	case 3:
		err = cell.Skip(1) // skip anycast (0)
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
		err = a.SetValue(v, AnyAddress{
			SumType: "VarAddress",
			VarAddress: &VarAddress{
				Len:       int16(ln),
				Workchain: int32(workchain),
				Address:   bs,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AnyAddress) UnmarshalTolk(cell *boc.Cell, ty AddressAny, abiCtx *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	switch tag {
	case 0:
		a.SumType = "NoneAddress"
	case 1:
		ln, err := cell.ReadUint(9)
		if err != nil {
			return err
		}
		bs, err := cell.ReadBits(int(ln))
		if err != nil {
			return err
		}
		a.SumType = "ExternalAddress"
		a.ExternalAddress = &ExternalAddress{
			Len:     int16(ln),
			Address: bs,
		}
	case 2:
		err = cell.Skip(1) // skip anycast (0)
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
		a.SumType = "InternalAddress"
		a.InternalAddress = &InternalAddress{
			Workchain: int8(workchain),
			Address:   [32]byte(address),
		}
	case 3:
		err = cell.Skip(1) // skip anycast (0)
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
		a.SumType = "VarAddress"
		a.VarAddress = &VarAddress{
			Len:       int16(ln),
			Workchain: int32(workchain),
			Address:   bs,
		}
	}

	return nil
}

func (AddressAny) MarshalTolk(cell *boc.Cell, v *Value) error {
	if v.anyAddress == nil {
		return fmt.Errorf("any address not found")
	}

	switch v.anyAddress.SumType {
	case "NoneAddress":
		err := cell.WriteUint(0, 2)
		if err != nil {
			return err
		}
	case "InternalAddress":
		err := cell.WriteUint(0b100, 3) // internal addr type ($10) and anycast (0)
		if err != nil {
			return err
		}
		err = cell.WriteInt(int64(v.anyAddress.InternalAddress.Workchain), 8)
		if err != nil {
			return err
		}
		err = cell.WriteBytes(v.anyAddress.InternalAddress.Address[:])
		if err != nil {
			return err
		}
	case "ExternalAddress":
		err := cell.WriteUint(1, 2) // external addr type ($01)
		if err != nil {
			return err
		}
		err = cell.WriteUint(uint64(v.anyAddress.ExternalAddress.Len), 9)
		if err != nil {
			return err
		}
		err = cell.WriteBitString(v.anyAddress.ExternalAddress.Address)
		if err != nil {
			return err
		}
	case "VarAddress":
		err := cell.WriteUint(0b110, 3) // var addr type ($11) and anycast (0)
		if err != nil {
			return err
		}
		err = cell.WriteUint(uint64(v.anyAddress.VarAddress.Len), 9)
		if err != nil {
			return err
		}
		err = cell.WriteInt(int64(v.anyAddress.VarAddress.Workchain), 32)
		if err != nil {
			return err
		}
		err = cell.WriteBitString(v.anyAddress.VarAddress.Address)
		if err != nil {
			return err
		}
	}

	return nil
}

func (AddressAny) Equal(v Value, o Value) bool {
	return false
}
