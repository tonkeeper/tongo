package precompiled

import (
	"math/big"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
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
		return tlb.VmStack{}, err
	}

	// result[1] - index value
	var indexVal tlb.VmStackValue
	if int64(body.Index) < 0 {
		var b big.Int
		b.SetUint64(body.Index)
		indexVal.SumType = "VmStkInt"
		indexVal.VmStkInt = tlb.Int257(b)
	} else {
		indexVal.SumType = "VmStkTinyInt"
		indexVal.VmStkTinyInt = int64(body.Index)
	}

	// result[2] - collection slice
	collectionCells := boc.NewCell()
	err = tlb.Marshal(collectionCells, body.Collection)
	if err != nil {
		return tlb.VmStack{}, err
	}
	collectionSlice, err := tlb.CellToVmCellSlice(collectionCells)
	if err != nil {
		return tlb.VmStack{}, err
	}

	// result[0], result[3], result[4] - init flag, owner, data (conditional)
	var initVal, ownerVal, dataVal tlb.VmStackValue
	if data.BitsAvailableForRead() > 0 {
		err = tlb.Unmarshal(data, &body2)
		if err != nil {
			return tlb.VmStack{}, err
		}
		initVal = tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: -1}
		ownerCell := boc.NewCell()
		err = tlb.Marshal(ownerCell, body2.Owner)
		if err != nil {
			return tlb.VmStack{}, err
		}
		ownerVal, err = tlb.CellToVmCellSlice(ownerCell)
		if err != nil {
			return tlb.VmStack{}, err
		}
		dataVal = tlb.VmStackValue{
			SumType:   "VmStkCell",
			VmStkCell: tlb.Ref[boc.Cell]{Value: body2.Data},
		}
	} else {
		initVal = tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: 0}
		ownerVal = tlb.VmStackValue{SumType: "VmStkNull"}
		dataVal = tlb.VmStackValue{SumType: "VmStkNull"}
	}

	var result tlb.VmStack
	result.Put(dataVal)
	result.Put(ownerVal)
	result.Put(collectionSlice)
	result.Put(indexVal)
	result.Put(initVal)
	return result, nil
}
