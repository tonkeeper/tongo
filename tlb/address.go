package tlb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

// AddressWithWorkchain is a TL-B type that represents the key in "suspended_address_list.addresses":
// suspended_address_list#00 addresses:(HashmapE AddressWithWorkchain Unit) suspended_until:uint32 = SuspendedAddressList;
type AddressWithWorkchain struct {
	Workchain int8
	Address   Bits256
}

func (addr AddressWithWorkchain) Equal(other any) bool {
	otherAddr, ok := other.(AddressWithWorkchain)
	if !ok {
		return false
	}
	return addr.Workchain == otherAddr.Workchain && addr.Address == otherAddr.Address
}

// Compare returns an integer comparing two addresses lexicographically.
// The result will be 0 if both are equal, -1 if addr < other, and +1 if addr > other.
func (addr AddressWithWorkchain) Compare(other any) (int, bool) {
	otherAddr, ok := other.(AddressWithWorkchain)
	if !ok {
		return 0, false
	}
	workchain := uint32(addr.Workchain)
	otherWorkchain := uint32(otherAddr.Workchain)
	if workchain < otherWorkchain {
		return -1, true
	}
	if workchain > otherWorkchain {
		return 1, true
	}
	return bytes.Compare(addr.Address[:], otherAddr.Address[:]), true
}

func (addr AddressWithWorkchain) FixedSize() int {
	return 288
}

func (addr *AddressWithWorkchain) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	wc, err := c.ReadInt(32)
	if err != nil {
		return err
	}
	addr.Workchain = int8(wc)
	bytes, err := c.ReadBytes(32)
	if err != nil {
		return err
	}
	copy(addr.Address[:], bytes)
	return nil
}

func (addr AddressWithWorkchain) MarshalJSON() ([]byte, error) {
	raw := fmt.Sprintf("%v:%x", addr.Workchain, addr.Address)
	return []byte(`"` + raw + `"`), nil
}

// InternalAddress is a TL-B type that represents only internal message address:
// addr_std$10 anycast:(Maybe Anycast) workchain_id:int8 address:bits256 = MsgAddressInt;
// Anycast is deprecated since TVM 10, so internal message always has `100` prefix
type InternalAddress struct {
	Workchain int8
	Address   Bits256
}

func (addr InternalAddress) Equal(other any) bool {
	otherAddr, ok := other.(InternalAddress)
	if !ok {
		return false
	}
	return addr.Workchain == otherAddr.Workchain && addr.Address == otherAddr.Address
}

// Compare returns an integer comparing two addresses lexicographically.
// The result will be 0 if both are equal, -1 if addr < other, and +1 if addr > other.
func (addr InternalAddress) Compare(other any) (int, bool) {
	otherAddr, ok := other.(InternalAddress)
	if !ok {
		return 0, false
	}
	workchain := uint32(addr.Workchain)
	otherWorkchain := uint32(otherAddr.Workchain)
	if workchain < otherWorkchain {
		return -1, true
	}
	if workchain > otherWorkchain {
		return 1, true
	}
	return bytes.Compare(addr.Address[:], otherAddr.Address[:]), true
}

func (addr InternalAddress) FixedSize() int {
	return 267
}

func (addr *InternalAddress) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	err := c.Skip(3) // $100 prefix
	if err != nil {
		return err
	}
	workchain, err := c.ReadInt(8)
	if err != nil {
		return err
	}
	address, err := c.ReadBytes(32)
	if err != nil {
		return err
	}
	addr.Workchain = int8(workchain)
	copy(addr.Address[:], address)
	return nil
}

func (addr InternalAddress) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	if err := c.WriteUint(2, 3); err != nil {
		return err
	}
	if err := c.WriteInt(int64(addr.Workchain), 8); err != nil {
		return err
	}
	return c.WriteBytes(addr.Address[:])
}

func (addr InternalAddress) MarshalJSON() ([]byte, error) {
	var anycastExtra string
	x := fmt.Sprintf("%d:%s", addr.Workchain, addr.Address.Hex())
	return []byte(fmt.Sprintf(`"%s%s"`, x, anycastExtra)), nil
}

func (addr *InternalAddress) UnmarshalJSON(b []byte) error {
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address: %s", string(b))
	}
	num, err := strconv.ParseInt(parts[0], 10, 32)
	isWorkchainInt8 := err == nil && num >= int64(math.MinInt8) && num <= int64(math.MaxInt8)
	if !isWorkchainInt8 {
		return fmt.Errorf("invalid address: %s", string(b))
	}
	var dst [32]byte
	_, err = hex.Decode(dst[:], []byte(parts[1]))
	if err != nil {
		return err
	}
	addr.Workchain = int8(num)
	addr.Address = dst
	return nil
}
