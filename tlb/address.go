package tlb

import (
	"bytes"
	"fmt"

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
