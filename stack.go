package tongo

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"io"
	"math/big"
)

// VmStack
// vm_stack#_ depth:(## 24) stack:(VmStackList depth) = VmStack;
// vm_stk_cons#_ {n:#} rest:^(VmStackList n) tos:VmStackValue = VmStackList (n + 1);
// vm_stk_nil#_ = VmStackList 0;
type VmStack []VmStackValue

// VmCont
// _ cregs:(HashmapE 4 VmStackValue) = VmSaveList;
// vm_ctl_data$_ nargs:(Maybe uint13) stack:(Maybe VmStack) save:VmSaveList
// cp:(Maybe int16) = VmControlData;
// vmc_std$00 cdata:VmControlData code:VmCellSlice = VmCont;
// vmc_envelope$01 cdata:VmControlData next:^VmCont = VmCont;
// vmc_quit$1000 exit_code:int32 = VmCont;
// vmc_quit_exc$1001 = VmCont;
// vmc_repeat$10100 count:uint63 body:^VmCont after:^VmCont = VmCont;
// vmc_until$110000 body:^VmCont after:^VmCont = VmCont;
// vmc_again$110001 body:^VmCont = VmCont;
// vmc_while_cond$110010 cond:^VmCont body:^VmCont
// after:^VmCont = VmCont;
// vmc_while_body$110011 cond:^VmCont body:^VmCont
// after:^VmCont = VmCont;
// vmc_pushint$1111 value:int32 next:^VmCont = VmCont;
type VmCont struct {
	// TODO: implement
}

// VmStkTuple
// Custom type: len:(## 16) data:(VmTuple len). Tag excluded. Use with VmStackValue type.
// vm_tupref_nil$_ = VmTupleRef 0;
// vm_tupref_single$_ entry:^VmStackValue = VmTupleRef 1;
// vm_tupref_any$_ {n:#} ref:^(VmTuple (n + 2)) = VmTupleRef (n + 2);
// vm_tuple_nil$_ = VmTuple 0;
// vm_tuple_tcons$_ {n:#} head:(VmTupleRef n) tail:^VmStackValue = VmTuple (n + 1);
// vm_stk_tuple#07 len:(## 16) data:(VmTuple len) = VmStackValue;
type VmStkTuple []VmStackValue

// VmCellSlice
// _ cell:^Cell st_bits:(## 10) end_bits:(## 10) { st_bits <= end_bits }
// st_ref:(#<= 4) end_ref:(#<= 4) { st_ref <= end_ref } = VmCellSlice;
type VmCellSlice struct {
	cell    *boc.Cell
	stBits  int
	endBits int
	stRef   int
	endRef  int
}

// VmStackValue
// vm_stk_null#00 = VmStackValue;
// vm_stk_tinyint#01 value:int64 = VmStackValue;
// vm_stk_int#0201_ value:int257 = VmStackValue;
// vm_stk_nan#02ff = VmStackValue;
// vm_stk_cell#03 cell:^Cell = VmStackValue;
// vm_stk_slice#04 _:VmCellSlice = VmStackValue;
// vm_stk_builder#05 cell:^Cell = VmStackValue;
// vm_stk_cont#06 cont:VmCont = VmStackValue;
// vm_stk_tuple#07 len:(## 16) data:(VmTuple len) = VmStackValue;
type VmStackValue struct {
	tl.SumType
	VmStkNull    struct{} `tlbSumType:"vm_stk_null#00"`
	VmStkTinyInt int64    `tlbSumType:"vm_stk_tinyint#01"`
	VmStkInt     struct {
		Value big.Int `tlb:"256bits"`
	} `tlbSumType:"vm_stk_int$00100000000"` // vm_stk_int#0201_
	VmStkNan     struct{}          `tlbSumType:"vm_stk_nan#02ff"`
	VmStkCell    tlb.Ref[boc.Cell] `tlbSumType:"vm_stk_cell#03"`
	VmStkSlice   VmCellSlice       `tlbSumType:"vm_stk_slice#04"`
	VmStkBuilder tlb.Ref[boc.Cell] `tlbSumType:"vm_stk_builder#05"`
	VmStkCont    VmCont            `tlbSumType:"vm_stk_cont#06"`
	VmStkTuple   VmStkTuple        `tlbSumType:"vm_stk_tuple#07"`
}

func (s VmStack) MarshalTLB(c *boc.Cell, tag string) error {
	depth := uint64(len(s))
	err := c.WriteUint(depth, 24)
	if err != nil {
		return err
	}
	err = putStackListItems(c, s)
	return err
}

func (s *VmStack) UnmarshalTLB(c *boc.Cell, tag string) error {
	depth, err := c.ReadUint(24)
	if err != nil {
		return err
	}
	if depth == 0 {
		return nil
	}
	list, err := getStackListItems(c, depth)
	if err != nil {
		return err
	}
	*s = list
	return nil
}

func getStackListItems(c *boc.Cell, depth uint64) ([]VmStackValue, error) {
	var (
		res []VmStackValue
		tos VmStackValue
	)
	if depth == 0 {
		return nil, nil
	}
	restCell, err := c.NextRef()
	if err != nil {
		return nil, err
	}
	rest, err := getStackListItems(restCell, depth-1)
	if err != nil {
		return nil, err
	}
	res = append(res, rest...)
	err = tlb.Unmarshal(c, &tos)
	if err != nil {
		return nil, err
	}
	res = append(res, tos)
	return res, nil
}

func putStackListItems(c *boc.Cell, list []VmStackValue) error {
	if len(list) == 0 {
		return nil
	}
	restCell := boc.NewCell()
	err := putStackListItems(restCell, list[1:])
	if err != nil {
		return err
	}
	err = c.AddRef(restCell)
	if err != nil {
		return err
	}
	err = tlb.Marshal(c, list[0])
	return err
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
	if len(b) == 0 {
		return nil
	}
	cell, err := boc.DeserializeBoc(b)
	if err != nil {
		return err
	}
	return tlb.Unmarshal(cell[0], s)
}

func (s VmCellSlice) MarshalTLB(c *boc.Cell, tag string) error {
	if s.stBits > s.endBits {
		return fmt.Errorf("invalid StBits and EndBits for CellSlice")
	}
	if s.stRef > s.endRef {
		return fmt.Errorf("invalid StRef and EndRef for CellSlice")
	}
	if s.endBits > s.cell.BitSize() {
		return fmt.Errorf("EndBits > Cell bit len")
	}
	if s.endRef > s.cell.RefsSize() {
		return fmt.Errorf("EndRef > Cell ref qty")
	}
	err := c.AddRef(s.cell)
	if err != nil {
		return err
	}
	err = c.WriteUint(uint64(s.stBits), 10)
	if err != nil {
		return err
	}
	err = c.WriteUint(uint64(s.endBits), 10)
	if err != nil {
		return err
	}
	err = c.WriteLimUint(s.stRef, 4)
	if err != nil {
		return err
	}
	err = c.WriteLimUint(s.endRef, 4)
	return err
}

func (s *VmCellSlice) UnmarshalTLB(c *boc.Cell, tag string) error {
	cell, err := c.NextRef()
	if err != nil {
		return err
	}
	stBits, err := c.ReadUint(10)
	if err != nil {
		return err
	}
	endBits, err := c.ReadUint(10)
	if err != nil {
		return err
	}
	if stBits > endBits {
		return fmt.Errorf("invalid StBits and EndBits for CellSlice")
	}
	stRef, err := c.ReadLimUint(4)
	if err != nil {
		return err
	}
	endRef, err := c.ReadLimUint(4)
	if err != nil {
		return err
	}
	if stRef > endRef {
		return fmt.Errorf("invalid StRef and EndRef for CellSlice")
	}
	if int(endBits) > cell.BitSize() {
		return fmt.Errorf("EndBits > Cell bit len")
	}
	if int(endRef) > cell.RefsSize() {
		return fmt.Errorf("EndRef > Cell ref qty")
	}
	*s = VmCellSlice{
		cell:    cell,
		stBits:  int(stBits),
		endBits: int(endBits),
		stRef:   int(stRef),
		endRef:  int(endRef),
	}
	return nil
}

func (s VmCellSlice) Cell() *boc.Cell {
	// TODO: maybe add as a filed to VmCellSlice
	cell := boc.NewCell()
	err := s.cell.Skip(s.stBits)
	if err != nil {
		panic("not enough cell bits")
	}
	bits, err := s.cell.ReadBits(s.endBits - s.stBits)
	if err != nil {
		panic("not enough cell bits")
	}
	refs := s.cell.Refs()
	err = cell.WriteBitString(bits)
	if err != nil {
		panic("can not write bits to empty cell")
	}
	for _, ref := range refs[s.stRef:s.endRef] {
		err = cell.AddRef(ref)
		if err != nil {
			panic("can not write ref to empty cell")
		}
	}
	return cell
}

func (t VmStkTuple) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("VmStkTuple TLB marshaling not implemented")
}

func (t *VmStkTuple) UnmarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("VmStkTuple TLB unmarshaling not implemented")
}

func (ct VmCont) MarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("VmCont TLB marshaling not implemented")
}

func (ct *VmCont) UnmarshalTLB(c *boc.Cell, tag string) error {
	// TODO: implement
	return fmt.Errorf("VmCont TLB unmarshaling not implemented")
}

func TlbStructToVmCellSlice(s any) (VmCellSlice, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, s)
	if err != nil {
		return VmCellSlice{}, err
	}
	return VmCellSlice{
		cell:    cell,
		stBits:  0,
		endBits: cell.BitSize(),
		stRef:   0,
		endRef:  cell.RefsSize(),
	}, nil
}

func (s VmCellSlice) UnmarshalToTlbStruct(res any) error {
	cell := s.Cell()
	err := tlb.Unmarshal(cell, res)
	if err != nil {
		return err
	}
	return nil
}
