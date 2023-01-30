package liteapi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"math/big"
)

// DnsResolve
// TEP-81 TON DNS Standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0081-dns-standard.md
// Returns the length (in bits) of the prefix of the internal representation of the domain that has been resolved and Cell with the TON DNS record
// Returns 0,nil,nil if DNS record is Null
func (c *Client) DnsResolve(ctx context.Context, address tongo.AccountID, domain string, category *big.Int) (int, *boc.Cell, error) {
	var params tlb.VmStack
	cell := boc.NewCell()
	err := cell.WriteBytes([]byte(domain))
	if err != nil {
		return 0, nil, err
	}
	dom, err := tlb.CellToVmCellSlice(cell)
	if err != nil {
		return 0, nil, err
	}
	params.Put(dom)
	int257 := tlb.Int257(*category)
	if err != nil {
		return 0, nil, err
	}
	cat := tlb.VmStackValue{
		SumType:  "VmStkInt",
		VmStkInt: int257,
	}
	params.Put(cat)
	errCode, stack, err := c.RunSmcMethod(ctx, address, "dnsresolve", params)
	if err != nil {
		return 0, nil, err
	}
	if errCode != 0 && errCode != 1 {
		return 0, nil, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 2 ||
		stack[0].SumType != "VmStkTinyInt" ||
		(stack[1].SumType != "VmStkCell" && stack[1].SumType != "VmStkNull") {
		return 0, nil, fmt.Errorf("invalid stack")
	}
	if stack[1].SumType == "VmStkNull" {
		return 0, nil, nil
	}
	return int(stack[0].VmStkTinyInt), &stack[1].VmStkCell.Value, err
}

func (c *Client) GetRootDNS(ctx context.Context) (tongo.AccountID, error) {
	conf, err := c.GetConfigAll(ctx, 0) // TODO: check mode
	if err != nil {
		return tongo.AccountID{}, err
	}
	for i, k := range conf.Config.Keys() {
		if k == 4 {
			addr, err := conf.Config.Values()[i].Value.ReadBytes(32)
			if err != nil {
				return tongo.AccountID{}, err
			}
			res := tongo.AccountID{
				Workchain: -1,
			}
			copy(res.Address[:], addr)
			return res, nil
		}
	}
	return tongo.AccountID{}, fmt.Errorf("config parameter not found")
}
