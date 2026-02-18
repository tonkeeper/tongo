package precompiled

import (
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func jettonV1getWalletData(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var body struct {
		Amount tlb.VarUInteger16
		Owner  tlb.MsgAddress
		Master tlb.MsgAddress
		Code   boc.Cell `tlb:"^"`
	}
	err := tlb.Unmarshal(data, &body)
	if err != nil {
		return tlb.VmStack{}, err
	}
	var b = big.Int(body.Amount)
	ownerCell := boc.NewCell()
	masterCell := boc.NewCell()
	err = tlb.Marshal(ownerCell, body.Owner)
	if err != nil {
		return tlb.VmStack{}, err
	}
	err = tlb.Marshal(masterCell, body.Master)
	if err != nil {
		return tlb.VmStack{}, err
	}
	ownerSlice, err := tlb.CellToVmCellSlice(ownerCell)
	if err != nil {
		return tlb.VmStack{}, err
	}
	masterSlice, err := tlb.CellToVmCellSlice(masterCell)
	if err != nil {
		return tlb.VmStack{}, err
	}
	var result tlb.VmStack
	result.Put(tlb.VmStackValue{
		SumType:   "VmStkCell",
		VmStkCell: tlb.Ref[boc.Cell]{Value: body.Code},
	})
	result.Put(masterSlice)
	result.Put(ownerSlice)
	if b.IsInt64() {
		result.Put(tlb.VmStackValue{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: b.Int64(),
		})
	} else {
		result.Put(tlb.VmStackValue{
			SumType:  "VmStkInt",
			VmStkInt: tlb.Int257(b),
		})
	}
	return result, nil
}

func jettonV2getWalletData(code string) func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	if _, errP := boc.DeserializeSinglRootBase64(code); errP != nil {
		panic(errP)
	}
	return func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
		var body struct {
			Status tlb.Uint4
			Amount tlb.VarUInteger16
			Owner  tlb.MsgAddress
			Master tlb.MsgAddress
		}
		err := tlb.Unmarshal(data, &body)
		if err != nil {
			return tlb.VmStack{}, err
		}
		var b = big.Int(body.Amount)
		ownerCell := boc.NewCell()
		masterCell := boc.NewCell()
		err = tlb.Marshal(ownerCell, body.Owner)
		if err != nil {
			return tlb.VmStack{}, err
		}
		err = tlb.Marshal(masterCell, body.Master)
		if err != nil {
			return tlb.VmStack{}, err
		}
		codeCell, _ := boc.DeserializeSinglRootBase64(code)
		ownerSlice, err := tlb.CellToVmCellSlice(ownerCell)
		if err != nil {
			return tlb.VmStack{}, err
		}
		masterSlice, err := tlb.CellToVmCellSlice(masterCell)
		if err != nil {
			return tlb.VmStack{}, err
		}
		var result tlb.VmStack
		result.Put(tlb.VmStackValue{
			SumType:   "VmStkCell",
			VmStkCell: tlb.Ref[boc.Cell]{Value: *codeCell},
		})
		result.Put(masterSlice)
		result.Put(ownerSlice)
		if b.IsInt64() {
			result.Put(tlb.VmStackValue{
				SumType:      "VmStkTinyInt",
				VmStkTinyInt: b.Int64(),
			})
		} else {
			result.Put(tlb.VmStackValue{
				SumType:  "VmStkInt",
				VmStkInt: tlb.Int257(b),
			})
		}
		return result, nil
	}
}

func jettonV3getWalletData(code string) func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	if _, errP := boc.DeserializeSinglRootBase64(code); errP != nil {
		panic(errP)
	}
	return func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
		//  initialized_version:uint8
		//        blocklisted:bool
		//        balance:Coins
		//        owner_address:MsgAddressInt
		//        jetton_master_address:MsgAddressInt
		var body struct {
			InitializedVersion tlb.Uint8
			Blocklisted        bool
			Amount             tlb.VarUInteger16
			Owner              tlb.MsgAddress
			Master             tlb.MsgAddress
		}
		err := tlb.Unmarshal(data, &body)
		if err != nil {
			return tlb.VmStack{}, err
		}
		var b = big.Int(body.Amount)
		ownerCell := boc.NewCell()
		masterCell := boc.NewCell()
		err = tlb.Marshal(ownerCell, body.Owner)
		if err != nil {
			return tlb.VmStack{}, err
		}
		err = tlb.Marshal(masterCell, body.Master)
		if err != nil {
			return tlb.VmStack{}, err
		}
		codeCell, _ := boc.DeserializeSinglRootBase64(code)
		ownerSlice, err := tlb.CellToVmCellSlice(ownerCell)
		if err != nil {
			return tlb.VmStack{}, err
		}
		masterSlice, err := tlb.CellToVmCellSlice(masterCell)
		if err != nil {
			return tlb.VmStack{}, err
		}
		var result tlb.VmStack
		result.Put(tlb.VmStackValue{
			SumType:   "VmStkCell",
			VmStkCell: tlb.Ref[boc.Cell]{Value: *codeCell},
		})
		result.Put(masterSlice)
		result.Put(ownerSlice)
		if b.IsInt64() {
			result.Put(tlb.VmStackValue{
				SumType:      "VmStkTinyInt",
				VmStkTinyInt: b.Int64(),
			})
		} else {
			result.Put(tlb.VmStackValue{
				SumType:  "VmStkInt",
				VmStkInt: tlb.Int257(b),
			})
		}
		return result, nil
	}
}

func isClaimed(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var body struct {
		Status tlb.Uint4
		Amount tlb.VarUInteger16
		Owner  tlb.MsgAddress
		Master tlb.MsgAddress
	}
	err := tlb.Unmarshal(data, &body)
	if err != nil {
		return tlb.VmStack{}, err
	}
	return tlb.VmStackValue{
		SumType:      "VmStkTinyInt",
		VmStkTinyInt: int64(body.Status),
	}.ToStack(), nil
}
