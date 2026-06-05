// Code generated - DO NOT EDIT.

package abiStonfi

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *BilateralLockArgs) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Resolver.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Resolver: %v", err)
	}
	if err = v.ResolverTimeoutDelta.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResolverTimeoutDelta: %v", err)
	}
	if err = v.ResolverAskAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResolverAskAmount: %v", err)
	}
	if err = v.ResolverPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ResolverPublicKey: %v", err)
	}
	if err = v.DutchSegments.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DutchSegments: %v", err)
	}
	return nil
}
func (v BilateralLockArgs) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Resolver.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Resolver: %v", err)
	}
	if err = v.ResolverTimeoutDelta.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResolverTimeoutDelta: %v", err)
	}
	if err = v.ResolverAskAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResolverAskAmount: %v", err)
	}
	if err = v.ResolverPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ResolverPublicKey: %v", err)
	}
	if err = v.DutchSegments.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DutchSegments: %v", err)
	}
	return nil
}
func (v BilateralLockArgs) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BilateralSignPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Signature.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signature: %v", err)
	}
	if err = v.Message.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Message: %v", err)
	}
	return nil
}
func (v BilateralSignPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Signature.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signature: %v", err)
	}
	if err = v.Message.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Message: %v", err)
	}
	return nil
}
func (v BilateralSignPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BilateralSignedMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.ExpirationTime.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ExpirationTime: %v", err)
	}
	if err = v.Resolver.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Resolver: %v", err)
	}
	return nil
}
func (v BilateralSignedMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.ExpirationTime.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ExpirationTime: %v", err)
	}
	if err = v.Resolver.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Resolver: %v", err)
	}
	return nil
}
func (v BilateralSignedMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *BilateralUnlockArgs) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MinOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinOut: %v", err)
	}
	if err = v.IgnoreRefundAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .IgnoreRefundAmount: %v", err)
	}
	if err = v.Signed.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Signed: %v", err)
	}
	return nil
}
func (v BilateralUnlockArgs) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MinOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinOut: %v", err)
	}
	if err = v.IgnoreRefundAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .IgnoreRefundAmount: %v", err)
	}
	if err = v.Signed.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Signed: %v", err)
	}
	return nil
}
func (v BilateralUnlockArgs) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *DutchSegments) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Num.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Num: %v", err)
	}
	if v.Segments, err = (func() (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})(); err != nil {
		return fmt.Errorf("failed to read .Segments: %v", err)
	}
	return nil
}
func (v DutchSegments) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Num.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Num: %v", err)
	}
	if err = c.AddRef(&v.Segments); err != nil {
		return fmt.Errorf("failed to .Segments: %v", err)
	}
	return nil
}
func (v DutchSegments) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ForwardParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Value.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Value: %v", err)
	}
	if v.SuccessPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .SuccessPayload: %v", err)
	}
	if v.RejectPayload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .RejectPayload: %v", err)
	}
	return nil
}
func (v ForwardParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Value.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Value: %v", err)
	}
	if err = v.SuccessPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SuccessPayload: %v", err)
	}
	if err = v.RejectPayload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RejectPayload: %v", err)
	}
	return nil
}
func (v ForwardParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *LockForwardParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Dest.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Dest: %v", err)
	}
	if err = v.ForwardParams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ForwardParams: %v", err)
	}
	return nil
}
func (v LockForwardParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Dest.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Dest: %v", err)
	}
	if err = v.ForwardParams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ForwardParams: %v", err)
	}
	return nil
}
func (v LockForwardParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VaultLockAdditionalDataMore) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OrderOwner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OrderOwner: %v", err)
	}
	if err = v.RefundTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RefundTo: %v", err)
	}
	if err = v.AskJettonMinter.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AskJettonMinter: %v", err)
	}
	return nil
}
func (v VaultLockAdditionalDataMore) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OrderOwner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OrderOwner: %v", err)
	}
	if err = v.RefundTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RefundTo: %v", err)
	}
	if err = v.AskJettonMinter.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AskJettonMinter: %v", err)
	}
	return nil
}
func (v VaultLockAdditionalDataMore) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
