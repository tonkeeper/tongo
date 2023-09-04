package liteapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/tep64"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

var (
	ErrOnchainContentOnly = fmt.Errorf("only onchain jetton data supported")
)

// GetJettonWallet
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonWallet(ctx context.Context, master, owner ton.AccountID) (ton.AccountID, error) {
	val, err := tlb.TlbStructToVmCellSlice(owner.ToMsgAddress())
	if err != nil {
		return ton.AccountID{}, err
	}
	errCode, stack, err := c.RunSmcMethod(ctx, master, "get_wallet_address", tlb.VmStack{val})
	if err != nil {
		return ton.AccountID{}, err
	}
	if errCode != 0 && errCode != 1 {
		return ton.AccountID{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 1 || stack[0].SumType != "VmStkSlice" {
		return ton.AccountID{}, fmt.Errorf("invalid stack")
	}
	var res tlb.MsgAddress
	err = stack[0].VmStkSlice.UnmarshalToTlbStruct(&res)
	if err != nil {
		return ton.AccountID{}, err
	}
	addr, err := ton.AccountIDFromTlb(res)
	if err != nil {
		return ton.AccountID{}, err
	}
	if addr == nil {
		return ton.AccountID{}, fmt.Errorf("addres none")
	}
	return *addr, nil
}

// GetJettonData
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonData(ctx context.Context, master ton.AccountID) (tep64.Metadata, error) {
	errCode, stack, err := c.RunSmcMethod(ctx, master, "get_jetton_data", tlb.VmStack{})
	if err != nil {
		return tep64.Metadata{}, err
	}
	if errCode != 0 && errCode != 1 {
		return tep64.Metadata{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 5 || (stack[0].SumType != "VmStkTinyInt" && stack[0].SumType != "VmStkInt") ||
		stack[1].SumType != "VmStkTinyInt" ||
		stack[2].SumType != "VmStkSlice" ||
		stack[3].SumType != "VmStkCell" ||
		stack[4].SumType != "VmStkCell" {
		return tep64.Metadata{}, fmt.Errorf("invalid stack")
	}
	cell := &stack[3].VmStkCell.Value
	var content tlb.FullContent
	err = tlb.Unmarshal(cell, &content)
	if err != nil {
		return tep64.Metadata{}, err
	}
	if content.SumType != "Onchain" {
		return tep64.Metadata{}, ErrOnchainContentOnly
	}
	meta, err := tep64.ConvertOnchainData(content)
	if err != nil {
		return tep64.Metadata{}, err
	}
	return meta, nil
}

// GetJettonBalance
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonBalance(ctx context.Context, jettonWallet ton.AccountID) (*big.Int, error) {
	errCode, stack, err := c.RunSmcMethod(ctx, jettonWallet, "get_wallet_data", tlb.VmStack{})
	if err != nil {
		return nil, err
	}
	if errCode == 0xFFFFFF00 { // contract not init
		return big.NewInt(0), nil
	}
	if errCode != 0 && errCode != 1 {
		return nil, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 4 || (stack[0].SumType != "VmStkTinyInt" && stack[0].SumType != "VmStkInt") ||
		stack[1].SumType != "VmStkSlice" ||
		stack[2].SumType != "VmStkSlice" ||
		stack[3].SumType != "VmStkCell" {
		return nil, fmt.Errorf("invalid stack")
	}
	if stack[0].SumType == "VmStkTinyInt" {
		return big.NewInt(stack[0].VmStkTinyInt), nil
	}
	res := big.Int(stack[0].VmStkInt)
	return &res, nil
}
