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
	v.internalAddress = &a
	return nil
}

func (Address) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
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
	err = v.SetValue(InternalAddress{
		Workchain: int8(workchain),
		Address:   [32]byte(address),
	})
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
	v.optionalAddress = &a
	return nil
}

func (AddressOpt) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	if tag == 0 {
		err = v.SetValue(OptionalAddress{
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
	err = v.SetValue(OptionalAddress{
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
	return nil
}

func (AddressExt) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
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
	err = v.SetValue(ExternalAddress{
		Len:     int16(ln),
		Address: bs,
	})
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
	return nil
}

func (AddressAny) UnmarshalTolk(cell *boc.Cell, v TolkValue, decoder *Decoder) error {
	tag, err := cell.ReadUint(2)
	if err != nil {
		return err
	}
	switch tag {
	case 0:
		err = v.SetValue(AnyAddress{
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
		err = v.SetValue(AnyAddress{
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
		err = v.SetValue(AnyAddress{
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
		err = v.SetValue(AnyAddress{
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

func (AddressAny) Equal(v Value, o Value) bool {
	return false
}
