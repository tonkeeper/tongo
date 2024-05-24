package precompiled

import (
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

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
