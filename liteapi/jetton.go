package liteapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/tlb"
	"math/big"
)

// GetJettonWallet
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonWallet(ctx context.Context, master, owner tongo.AccountID) (tongo.AccountID, error) {
	val, err := tlb.TlbStructToVmCellSlice(owner)
	if err != nil {
		return tongo.AccountID{}, err
	}
	errCode, stack, err := c.RunSmcMethod(ctx, master, "get_wallet_address", tlb.VmStack{val})
	if err != nil {
		return tongo.AccountID{}, err
	}
	if errCode != 0 && errCode != 1 {
		return tongo.AccountID{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 1 || stack[0].SumType != "VmStkSlice" {
		return tongo.AccountID{}, fmt.Errorf("invalid stack")
	}
	var res tlb.MsgAddress
	err = stack[0].VmStkSlice.UnmarshalToTlbStruct(&res)
	if err != nil {
		return tongo.AccountID{}, err
	}
	addr, err := tongo.AccountIDFromTlb(res)
	if err != nil {
		return tongo.AccountID{}, err
	}
	if addr == nil {
		return tongo.AccountID{}, fmt.Errorf("addres none")
	}
	return *addr, nil
}

// GetJettonData
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonData(ctx context.Context, master tongo.AccountID) (tongo.JettonMetadata, error) {
	errCode, stack, err := c.RunSmcMethod(ctx, master, "get_jetton_data", tlb.VmStack{})
	if err != nil {
		return tongo.JettonMetadata{}, err
	}
	if errCode != 0 && errCode != 1 {
		return tongo.JettonMetadata{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 5 || (stack[0].SumType != "VmStkTinyInt" && stack[0].SumType != "VmStkInt") ||
		stack[1].SumType != "VmStkTinyInt" ||
		stack[2].SumType != "VmStkSlice" ||
		stack[3].SumType != "VmStkCell" ||
		stack[4].SumType != "VmStkCell" {
		return tongo.JettonMetadata{}, fmt.Errorf("invalid stack")
	}
	cell := &stack[3].VmStkCell.Value
	var content tlb.FullContent
	err = tlb.Unmarshal(cell, &content)
	if err != nil {
		return tongo.JettonMetadata{}, err
	}
	if content.SumType != "Onchain" {
		return tongo.JettonMetadata{}, fmt.Errorf("only onchain jetton data supported")
	}
	meta, err := convertOnchainData(content)
	if err != nil {
		return tongo.JettonMetadata{}, err
	}
	return meta, nil
}

// GetJettonBalance
// TEP-74 Fungible tokens (Jettons) standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
func (c *Client) GetJettonBalance(ctx context.Context, jettonWallet tongo.AccountID) (*big.Int, error) {
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

// TEP-64 Token Data Standard
// https://github.com/ton-blockchain/TEPs/blob/master/text/0064-token-data-standard.md
func convertOnchainData(content tlb.FullContent) (tongo.JettonMetadata, error) {
	if content.SumType != "Onchain" {
		return tongo.JettonMetadata{}, fmt.Errorf("not Onchain content")
	}
	var m tongo.JettonMetadata
	for i, v := range content.Onchain.Data.Values() {
		keyS := hex.EncodeToString(content.Onchain.Data.Keys()[i][:])
		switch keyS {
		case "70e5d7b6a29b392f85076fe15ca2f2053c56c2338728c4e33c9e8ddb1ee827cc": // sha256(uri)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Uri = string(b)
		case "82a3537ff0dbce7eec35d69edc3a189ee6f17d82f353a553f9aa96cb0be3ce89": // sha256(name)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Name = string(b)
		case "c9046f7a37ad0ea7cee73355984fa5428982f8b37c8f7bcec91f7ac71a7cd104": // sha256(description)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Description = string(b)
		case "6105d6cc76af400325e94d588ce511be5bfdbb73b437dc51eca43917d7a43e3d": // sha256(image)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Image = string(b)
		case "d9a88ccec79eef59c84b671136a20ece4cd00caaad5bc47e2c208829154ee9e4": // sha256(image_data)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.ImageData = b
		case "b76a7ca153c24671658335bbd08946350ffc621fa1c516e7123095d4ffd5c581": // sha256(symbol)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Symbol = string(b)
		case "ee80fd2f1e03480e2282363596ee752d7bb27f50776b95086a0279189675923e": // sha256(decimals)
			b, err := v.Value.Bytes()
			if err != nil {
				return tongo.JettonMetadata{}, err
			}
			m.Decimals = string(b)
		}
	}
	return m, nil
}
