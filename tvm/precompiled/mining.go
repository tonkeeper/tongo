package precompiled

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

var getPowParamsGram = func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error) {
	var body struct {
		Skip1  uint64
		Skip2  tlb.Bits256
		Value1 tlb.Uint128
		Value2 tlb.Uint256
		Field1 struct {
			Value3 tlb.Grams
			Value4 uint32
		} `tlb:"^"`
	}
	err := tlb.Unmarshal(data, &body)
	if err != nil {
		return tlb.VmStack{}, err
	}
	var result tlb.VmStack
	result.Put(tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: int64(body.Field1.Value4)})
	result.Put(tlb.VmStackValue{SumType: "VmStkTinyInt", VmStkTinyInt: int64(body.Field1.Value3)})
	result.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(body.Value2)})
	result.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257(body.Value1)})
	return result, nil
}
