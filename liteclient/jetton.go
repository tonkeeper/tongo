package liteclient

import (
	"context"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/tlb"
)

func (c *Client) GetJettonWallet(ctx context.Context, master, owner tongo.AccountID) (tongo.AccountID, error) {
	slice, err := tongo.TlbStructToVmCellSlice(owner)
	if err != nil {
		return tongo.AccountID{}, err
	}
	val := tongo.VmStackValue{
		SumType:    "VmStkSlice",
		VmStkSlice: slice,
	}
	errCode, stack, err := c.RunSmcMethod(ctx, 4, master, "get_wallet_address", tongo.VmStack{val})
	if err != nil {
		return tongo.AccountID{}, err
	}
	if errCode != 0 && errCode != 1 {
		return tongo.AccountID{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 1 || stack[0].SumType != "VmStkSlice" {
		return tongo.AccountID{}, fmt.Errorf("invalid stack")
	}
	// TODO: implement universal converter
	var res tongo.MsgAddress
	cell := stack[0].VmStkSlice.Cell()
	err = tlb.Unmarshal(cell, &res)
	if err != nil {
		return tongo.AccountID{}, err
	}
	addr, err := res.AccountID()
	if err != nil {
		return tongo.AccountID{}, err
	}
	if addr == nil {
		return tongo.AccountID{}, fmt.Errorf("addres none")
	}
	return *addr, nil
}

func (c *Client) GetDecimals(ctx context.Context, master tongo.AccountID) (uint32, error) {
	// TODO: implement
	return 0, fmt.Errorf("not implemnted")
}
