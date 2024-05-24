package precompiled

import (
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// todo: need to cover test corner cases about int size near int64 max value
// DO NOT USE UNTIL FIX TODO!!!
var jettonV1getWalletData = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var body struct {
		Amount tlb.VarUInteger16
		Owner  tlb.MsgAddress
		Master tlb.MsgAddress
		Code   boc.Cell `tlb:"^"`
	}
	err := tlb.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}
	var result = make([]tlb.VmStackValue, 4)
	var b = big.Int(body.Amount)
	if b.IsInt64() {
		result[0] = tlb.VmStackValue{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: b.Int64(),
		}
	} else {
		result[0] = tlb.VmStackValue{
			SumType:  "VmStkInt",
			VmStkInt: tlb.Int257(b),
		}
	}
	ownerCell := boc.NewCell()
	masterCell := boc.NewCell()
	err = tlb.Marshal(ownerCell, body.Owner)
	if err != nil {
		return nil, err
	}
	err = tlb.Marshal(masterCell, body.Master)
	if err != nil {
		return nil, err
	}
	result[1], err = tlb.CellToVmCellSlice(ownerCell)
	if err != nil {
		return nil, err
	}
	result[2], err = tlb.CellToVmCellSlice(masterCell)
	if err != nil {
		return nil, err
	}
	result[3] = tlb.VmStackValue{
		SumType:   "VmStkCell",
		VmStkCell: tlb.Ref[boc.Cell]{Value: body.Code},
	}
	return result, nil
}
