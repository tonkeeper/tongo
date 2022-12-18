package tlb

//GENERATED
import (
	"github.com/startfellows/tongo/boc"
	"math/big"
)

type VarUInteger1 big.Int

func (u VarUInteger1) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 0)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger1) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(0)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger1(*val)
	return nil
}

type VarUInteger2 big.Int

func (u VarUInteger2) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 1)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger2) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(1)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger2(*val)
	return nil
}

type VarUInteger3 big.Int

func (u VarUInteger3) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 2)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger3) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(2)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger3(*val)
	return nil
}

type VarUInteger4 big.Int

func (u VarUInteger4) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 3)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger4) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(3)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger4(*val)
	return nil
}

type VarUInteger5 big.Int

func (u VarUInteger5) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 4)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger5) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(4)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger5(*val)
	return nil
}

type VarUInteger6 big.Int

func (u VarUInteger6) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 5)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger6) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(5)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger6(*val)
	return nil
}

type VarUInteger7 big.Int

func (u VarUInteger7) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 6)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger7) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(6)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger7(*val)
	return nil
}

type VarUInteger8 big.Int

func (u VarUInteger8) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 7)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger8) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(7)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger8(*val)
	return nil
}

type VarUInteger9 big.Int

func (u VarUInteger9) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 8)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger9) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(8)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger9(*val)
	return nil
}

type VarUInteger10 big.Int

func (u VarUInteger10) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 9)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger10) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(9)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger10(*val)
	return nil
}

type VarUInteger11 big.Int

func (u VarUInteger11) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 10)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger11) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(10)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger11(*val)
	return nil
}

type VarUInteger12 big.Int

func (u VarUInteger12) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 11)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger12) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(11)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger12(*val)
	return nil
}

type VarUInteger13 big.Int

func (u VarUInteger13) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 12)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger13) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(12)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger13(*val)
	return nil
}

type VarUInteger14 big.Int

func (u VarUInteger14) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 13)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger14) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(13)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger14(*val)
	return nil
}

type VarUInteger15 big.Int

func (u VarUInteger15) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 14)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger15) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(14)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger15(*val)
	return nil
}

type VarUInteger16 big.Int

func (u VarUInteger16) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 15)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger16) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(15)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger16(*val)
	return nil
}

type VarUInteger17 big.Int

func (u VarUInteger17) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 16)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger17) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(16)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger17(*val)
	return nil
}

type VarUInteger18 big.Int

func (u VarUInteger18) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 17)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger18) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(17)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger18(*val)
	return nil
}

type VarUInteger19 big.Int

func (u VarUInteger19) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 18)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger19) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(18)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger19(*val)
	return nil
}

type VarUInteger20 big.Int

func (u VarUInteger20) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 19)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger20) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(19)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger20(*val)
	return nil
}

type VarUInteger21 big.Int

func (u VarUInteger21) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 20)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger21) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(20)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger21(*val)
	return nil
}

type VarUInteger22 big.Int

func (u VarUInteger22) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 21)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger22) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(21)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger22(*val)
	return nil
}

type VarUInteger23 big.Int

func (u VarUInteger23) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 22)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger23) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(22)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger23(*val)
	return nil
}

type VarUInteger24 big.Int

func (u VarUInteger24) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 23)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger24) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(23)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger24(*val)
	return nil
}

type VarUInteger25 big.Int

func (u VarUInteger25) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 24)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger25) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(24)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger25(*val)
	return nil
}

type VarUInteger26 big.Int

func (u VarUInteger26) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 25)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger26) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(25)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger26(*val)
	return nil
}

type VarUInteger27 big.Int

func (u VarUInteger27) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 26)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger27) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(26)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger27(*val)
	return nil
}

type VarUInteger28 big.Int

func (u VarUInteger28) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 27)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger28) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(27)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger28(*val)
	return nil
}

type VarUInteger29 big.Int

func (u VarUInteger29) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 28)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger29) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(28)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger29(*val)
	return nil
}

type VarUInteger30 big.Int

func (u VarUInteger30) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 29)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger30) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(29)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger30(*val)
	return nil
}

type VarUInteger31 big.Int

func (u VarUInteger31) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 30)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger31) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(30)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger31(*val)
	return nil
}

type VarUInteger32 big.Int

func (u VarUInteger32) MarshalTLB(c *boc.Cell, tag string) error {
	i := big.Int(u)
	b := i.Bytes()
	err := c.WriteLimUint(len(b), 31)
	if err != nil {
		return err
	}
	return c.WriteBytes(b)
}

func (u *VarUInteger32) UnmarshalTLB(c *boc.Cell, tag string) error {
	ln, err := c.ReadLimUint(31)
	if err != nil {
		return err
	}
	val, err := c.ReadBigUint(int(ln) * 8)
	if err != nil {
		return err
	}
	*u = VarUInteger32(*val)
	return nil
}
