package precompiled

import (
	"fmt"
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func walletv3seqno(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV3 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
	}
	err := tlb.Unmarshal(data, &dataV3)
	if err != nil {
		return nil, err
	}
	return tlb.VmStack{
		{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: int64(dataV3.Seqno),
		},
	}, nil
}

var walletv3r2publicKey = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV3 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
	}
	err := tlb.Unmarshal(data, &dataV3)
	if err != nil {
		return nil, err
	}
	var b big.Int
	b.SetBytes(dataV3.PublicKey[:])
	return tlb.VmStack{
		{
			SumType:  "VmStkInt",
			VmStkInt: tlb.Int257(b),
		},
	}, nil
}

var walletv4r2seqno = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV4 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
	}
	err := tlb.Unmarshal(data, &dataV4)
	if err != nil {
		return nil, err
	}
	return tlb.VmStack{
		{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: int64(dataV4.Seqno),
		},
	}, nil
}

var walletv4r2SubwalletID = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV4 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
	}
	err := tlb.Unmarshal(data, &dataV4)
	if err != nil {
		return nil, err
	}
	return tlb.VmStack{
		{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: int64(dataV4.SubWalletId),
		},
	}, nil
}

var walletv4r2publicKey = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV4 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
	}
	err := tlb.Unmarshal(data, &dataV4)
	if err != nil {
		return nil, err
	}
	var b big.Int
	b.SetBytes(dataV4.PublicKey[:])
	return tlb.VmStack{
		{
			SumType:  "VmStkInt",
			VmStkInt: tlb.Int257(b),
		},
	}, nil
}

var walletv4r2getPluginList = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV4 struct {
		Seqno       uint32
		SubWalletId uint32
		PublicKey   tlb.Bits256
		PluginDict  tlb.HashmapE[tlb.Bits264, struct{}]
	}
	err := tlb.Unmarshal(data, &dataV4)
	if err != nil {
		return nil, err
	}
	if len(dataV4.PluginDict.Keys()) == 0 {
		return tlb.VmStack{
			{SumType: "VmStkNull"},
		}, nil
	}
	//todo: implement
	return nil, fmt.Errorf("not implemented not empty dict")
}

func walletv5r1seqno(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV5 struct {
		IsSignatureAllowed bool
		Seqno              uint32
		WalletID           uint32
		PublicKey          tlb.Bits256
		Extensions         tlb.HashmapE[tlb.Bits256, tlb.Uint1]
	}
	err := tlb.Unmarshal(data, &dataV5)
	if err != nil {
		return nil, err
	}
	return tlb.VmStack{
		{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: int64(dataV5.Seqno),
		},
	}, nil
}

func walletv5r1SubwalletID(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV5 struct {
		IsSignatureAllowed bool
		Seqno              uint32
		WalletID           uint32
		PublicKey          tlb.Bits256
		Extensions         tlb.HashmapE[tlb.Bits256, tlb.Uint1]
	}
	err := tlb.Unmarshal(data, &dataV5)
	if err != nil {
		return nil, err
	}
	return tlb.VmStack{
		{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: int64(dataV5.WalletID),
		},
	}, nil
}

func walletv5r1publicKey(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var dataV5 struct {
		IsSignatureAllowed bool
		Seqno              uint32
		WalletID           uint32
		PublicKey          tlb.Bits256
		Extensions         tlb.HashmapE[tlb.Bits256, tlb.Uint1]
	}
	err := tlb.Unmarshal(data, &dataV5)
	if err != nil {
		return nil, err
	}
	var b big.Int
	b.SetBytes(dataV5.PublicKey[:])
	return tlb.VmStack{
		{
			SumType:  "VmStkInt",
			VmStkInt: tlb.Int257(b),
		},
	}, nil
}
