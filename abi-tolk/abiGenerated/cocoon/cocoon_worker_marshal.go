// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *CocoonWorkerData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.ProxyAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyAddress: %v", err)
	}
	if err = v.ProxyPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.Tokens.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Tokens: %v", err)
	}
	return nil
}
func (v CocoonWorkerData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.ProxyAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyAddress: %v", err)
	}
	if err = v.ProxyPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyPublicKey: %v", err)
	}
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.Tokens.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Tokens: %v", err)
	}
	return nil
}
func (v CocoonWorkerData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonWorkerData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.Tokens.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Tokens: %v", err)
	}
	if err = v.State.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.ProxyPublicKey.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	if err = v.ProxyAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProxyAddress: %v", err)
	}
	if err = v.OwnerAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	return nil
}
