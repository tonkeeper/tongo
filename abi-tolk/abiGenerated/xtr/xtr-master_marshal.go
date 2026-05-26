// Code generated - DO NOT EDIT.

package abiXtr

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *PushXTR) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixPushXTR); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}
func (v PushXTR) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixPushXTR, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}
func (v PushXTR) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CommitXTR) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixCommitXTR); err != nil {
		return err
	}
	if err = v.Seqno.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Seqno: %v", err)
	}
	if err = v.UserAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UserAddress: %v", err)
	}
	if err = v.Amount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Amount: %v", err)
	}
	return nil
}
func (v CommitXTR) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixCommitXTR, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Seqno.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Seqno: %v", err)
	}
	if err = v.UserAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UserAddress: %v", err)
	}
	if err = v.Amount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Amount: %v", err)
	}
	return nil
}
func (v CommitXTR) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdateUser) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdateUser); err != nil {
		return err
	}
	if err = v.DestAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DestAddress: %v", err)
	}
	if v.Payload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}
func (v UpdateUser) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdateUser, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.DestAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DestAddress: %v", err)
	}
	if err = c.AddRef(&v.Payload); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}
func (v UpdateUser) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdatePayment) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdatePayment); err != nil {
		return err
	}
	if err = v.DestAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .DestAddress: %v", err)
	}
	if v.Payload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}
func (v UpdatePayment) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdatePayment, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.DestAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .DestAddress: %v", err)
	}
	if err = c.AddRef(&v.Payload); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}
func (v UpdatePayment) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdateContractAndProcessMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdateContractAndProcessMessage); err != nil {
		return err
	}
	if err = v.UpdateData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UpdateData: %v", err)
	}
	if err = v.FromAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FromAddress: %v", err)
	}
	if err = v.FromAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .FromAmount: %v", err)
	}
	if v.Payload, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}
func (v UpdateContractAndProcessMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdateContractAndProcessMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.UpdateData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UpdateData: %v", err)
	}
	if err = v.FromAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FromAddress: %v", err)
	}
	if err = v.FromAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .FromAmount: %v", err)
	}
	if err = c.AddRef(&v.Payload); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}
func (v UpdateContractAndProcessMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdateData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Code, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Code: %v", err)
	}
	if v.Data, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Data: %v", err)
	}
	if err = v.Version.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Version: %v", err)
	}
	return nil
}
func (v UpdateData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.AddRef(&v.Code); err != nil {
		return fmt.Errorf("failed to .Code: %v", err)
	}
	if err = c.AddRef(&v.Data); err != nil {
		return fmt.Errorf("failed to .Data: %v", err)
	}
	if err = v.Version.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Version: %v", err)
	}
	return nil
}
func (v UpdateData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg UpdateUser) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg UpdatePayment) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg UpdateContractAndProcessMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg PushXTR) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg CommitXTR) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
