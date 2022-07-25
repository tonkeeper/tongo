package tongo

import (
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"io"
)

// VmStack
// vm_stack#_ depth:(## 24) stack:(VmStackList depth) = VmStack;
// vm_stk_cons#_ {n:#} rest:^(VmStackList n) tos:VmStackValue = VmStackList (n + 1);
// vm_stk_nil#_ = VmStackList 0;
type VmStack struct {
	// Draft!
	// TODO: implement
	Values []VmStackValue
}

type VmStackValue struct {
	// Draft!
	// TODO: implement
	tl.SumType
	VmStkNull    struct{} `tlbSumType:"vm_stk_null#00"`
	VmStkTinyint int64    `tlbSumType:"vm_stk_tinyint#01"`
	VmStkInt     struct{} `tlbSumType:"vm_stk_int$00100000000"` // vm_stk_int#0201_
	VmStkNan     struct{} `tlbSumType:"vm_stk_nan#02ff"`
}

func (s VmStack) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement not empty stack
	return c.WriteUint(0, 24) // depth = 0 empty stack
}

func (s *VmStack) UnmarshalTLB(c *boc.Cell, tag string) error {
	// Draft!
	// TODO: implement
	_, err := c.ReadUint(24) // depth
	if err != nil {
		return err
	}
	var val VmStackValue
	err = tlb.Unmarshal(c, &val)
	if err != nil {
		return err
	}
	s.Values = append(s.Values, val)
	return nil
}

func (s VmStack) MarshalTL() ([]byte, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, s)
	if err != nil {
		return nil, err
	}
	b, err := cell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return nil, err
	}
	return tl.Marshal(b)
}

func (s *VmStack) UnmarshalTL(r io.Reader) error {
	var b []byte
	err := tl.Unmarshal(r, &b)
	if err != nil {
		return err
	}
	cell, err := boc.DeserializeBoc(b)
	if err != nil {
		return err
	}
	return tlb.Unmarshal(cell[0], s)
}
