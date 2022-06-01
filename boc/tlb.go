package boc

import (
	"fmt"
	"github.com/startfellows/tongo"
)

type TLBType interface {
	UnmarshalTLB(c *CellReader) error
}

func ReadAddress(c *CellReader) (*tongo.AccountID, error) {
	prefix, err := c.ReadUint(2)
	if err != nil {
		return nil, err
	}
	if prefix == 0 { // adr_none prefix
		return nil, nil
	}
	if prefix != 2 { // not adr_std prefix
		return nil, fmt.Errorf("not std address")
	}
	maybe, err := c.ReadBit()
	if err != nil {
		return nil, err
	}
	if maybe == true {
		return nil, fmt.Errorf("anycast not being processed") //TODO: add anycast processing
	}
	workchain, err := c.ReadUint(8)
	if err != nil {
		return nil, err
	}
	addr, err := c.ReadBytes(32)
	if err != nil {
		return nil, err
	}
	var address tongo.AccountID
	address.Workchain = int32(workchain)
	address.Address = addr
	return &address, nil
}
