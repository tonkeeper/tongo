package precompiled

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"math/big"
)

func nftV1getNftData(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var body struct {
		Index      uint64
		Collection tlb.MsgAddress
	}
	var body2 struct {
		Owner tlb.MsgAddress
		Data  boc.Cell `tlb:"^"`
	}
	err := tlb.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}
	var result = make([]tlb.VmStackValue, 5)
	result[0].SumType = "VmStkTinyInt"
	if int64(body.Index) < 0 {
		var b big.Int
		b.SetUint64(body.Index)
		result[1].SumType = "VmStkInt"
		result[1].VmStkInt = tlb.Int257(b)
	} else {
		result[1].SumType = "VmStkTinyInt"
		result[1].VmStkTinyInt = int64(body.Index)
	}
	collectionCells := boc.NewCell()
	err = tlb.Marshal(collectionCells, body.Collection)
	if err != nil {
		return nil, err
	}
	result[2], err = tlb.CellToVmCellSlice(collectionCells)
	if err != nil {
		return nil, err
	}
	if data.BitsAvailableForRead() > 0 {
		err = tlb.Unmarshal(data, &body2)
		if err != nil {
			return nil, err
		}
		result[0].VmStkTinyInt = -1
		ownerCell := boc.NewCell()
		err = tlb.Marshal(ownerCell, body2.Owner)
		if err != nil {
			return nil, err
		}
		result[3], err = tlb.CellToVmCellSlice(ownerCell)
		if err != nil {
			return nil, err
		}
		result[4] = tlb.VmStackValue{
			SumType:   "VmStkCell",
			VmStkCell: tlb.Ref[boc.Cell]{Value: body2.Data},
		}
	} else {
		result[0].VmStkTinyInt = 0
		result[3].SumType = "VmStkNull"
		result[4].SumType = "VmStkNull"
	}
	return result, nil
}
