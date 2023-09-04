package liteapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// DnsResolve is deprecated. please use github.com/tonkeeper/tongo/contract/dns
func (c *Client) DnsResolve(ctx context.Context, address ton.AccountID, domain string, category *big.Int) (int, *boc.Cell, error) {
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

func (c *Client) GetRootDNS(ctx context.Context) (ton.AccountID, error) {
	conf, err := c.GetConfigAll(ctx, 0) // TODO: check mode
	if err != nil {
		return ton.AccountID{}, err
	}
	for i, k := range conf.Config.Keys() {
		if k == 4 {
			addr, err := conf.Config.Values()[i].Value.ReadBytes(32)
			if err != nil {
				return ton.AccountID{}, err
			}
			res := ton.AccountID{
				Workchain: -1,
			}
			copy(res.Address[:], addr)
			return res, nil
		}
	}
	return ton.AccountID{}, fmt.Errorf("config parameter not found")
}
