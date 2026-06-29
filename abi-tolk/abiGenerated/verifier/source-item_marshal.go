// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *SourceItemStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.VerifierId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceItemRegistry.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SourceItemRegistry: %v", err)
	}
	if err = v.Content.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	return nil
}

func (v SourceItemStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.VerifierId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceItemRegistry.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SourceItemRegistry: %v", err)
	}
	if err = v.Content.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	return nil
}

func (v SourceItemStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SourceItemData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.VerifierId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceItemRegistry.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SourceItemRegistry: %v", err)
	}
	if err = v.Content.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	return nil
}

func (v SourceItemData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.VerifierId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceItemRegistry.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SourceItemRegistry: %v", err)
	}
	if err = v.Content.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Content: %v", err)
	}
	return nil
}

func (v SourceItemData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SourceItemData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Content, err = (func() (value tlb.RefT[*SourceContent], err error) {
		var cIn boc.Cell
		cIn, err = stack.ReadCell()
		if err != nil {
			return
		}
		c := boc.NewCell()
		_ = c.AddRef(&cIn)
		err = value.UnmarshalTLB(c, tlb.NewDecoder())
		return
	})(); err != nil {
		return fmt.Errorf("failed to read .Content: %v", err)
	}
	if err = v.SourceItemRegistry.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .SourceItemRegistry: %v", err)
	}
	if err = v.VerifiedCodeCellHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .VerifiedCodeCellHash: %v", err)
	}
	if err = v.VerifierId.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	return nil
}

func (msg SourceContent) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourceItemStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
