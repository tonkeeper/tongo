package abiCocoon

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// UnmarshalTLB decodes a RegisteredProxy from inline cell bits.
// Binary layout (from C++ RootContractConfig::load_registered_proxies):
//
//	1 bit  - kind (0 = IPv4)
//	7 bits - address string length (max 127)
//	N*8 bits - ASCII address string
//
// Address format: "worker_addr client_addr" (space-separated) if they differ,
// or a single address if they are the same.
func (p *RegisteredProxy) UnmarshalTLB(c *boc.Cell, _ *tlb.Decoder) error {
	kind, err := c.ReadUint(1)
	if err != nil {
		return fmt.Errorf("RegisteredProxy kind: %w", err)
	}
	p.Kind = tlb.Uint1(kind)

	length, err := c.ReadUint(7)
	if err != nil {
		return fmt.Errorf("RegisteredProxy length: %w", err)
	}
	if length > 127 {
		return fmt.Errorf("RegisteredProxy address length %d exceeds maximum 127", length)
	}

	addrBytes, err := c.ReadBytes(int(length))
	if err != nil {
		return fmt.Errorf("RegisteredProxy address: %w", err)
	}
	p.Address = string(addrBytes)
	return nil
}
