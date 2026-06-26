// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *MessageDescription) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.VerifierId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	if err = v.ValidUntil.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ValidUntil: %v", err)
	}
	if err = v.SourceAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SourceAddr: %v", err)
	}
	if err = v.TargetAddr.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TargetAddr: %v", err)
	}
	if v.Msg, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	return nil
}
func (v MessageDescription) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.VerifierId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierId: %v", err)
	}
	if err = v.ValidUntil.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ValidUntil: %v", err)
	}
	if err = v.SourceAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SourceAddr: %v", err)
	}
	if err = v.TargetAddr.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TargetAddr: %v", err)
	}
	if err = c.AddRef(&v.Msg); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	return nil
}
func (v MessageDescription) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VerifierSettings) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MultiSigThreshold.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MultiSigThreshold: %v", err)
	}
	if err = v.PubKeyEndpoints.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PubKeyEndpoints: %v", err)
	}
	if v.Name, err = c.ReadStringRefTail(); err != nil {
		return fmt.Errorf("failed to read .Name: %v", err)
	}
	if v.MarketingUrl, err = c.ReadStringRefTail(); err != nil {
		return fmt.Errorf("failed to read .MarketingUrl: %v", err)
	}
	return nil
}
func (v VerifierSettings) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MultiSigThreshold.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MultiSigThreshold: %v", err)
	}
	if err = v.PubKeyEndpoints.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PubKeyEndpoints: %v", err)
	}
	if err = (func() error { _, err := c.WriteStringRefTail(v.Name); return err })(); err != nil {
		return fmt.Errorf("failed to .Name: %v", err)
	}
	if err = (func() error { _, err := c.WriteStringRefTail(v.MarketingUrl); return err })(); err != nil {
		return fmt.Errorf("failed to .MarketingUrl: %v", err)
	}
	return nil
}
func (v VerifierSettings) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Verifier) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Admin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Admin: %v", err)
	}
	if err = v.Settings.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Settings: %v", err)
	}
	return nil
}
func (v Verifier) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Admin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Admin: %v", err)
	}
	if err = v.Settings.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Settings: %v", err)
	}
	return nil
}
func (v Verifier) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VerifierRegistryStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Verifiers.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Verifiers: %v", err)
	}
	return nil
}
func (v VerifierRegistryStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Verifiers.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Verifiers: %v", err)
	}
	return nil
}
func (v VerifierRegistryStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ForwardMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixForwardMessage); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Msg.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Msg: %v", err)
	}
	if v.Signatures, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Signatures: %v", err)
	}
	return nil
}
func (v ForwardMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixForwardMessage, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Msg.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Msg: %v", err)
	}
	if err = c.AddRef(&v.Signatures); err != nil {
		return fmt.Errorf("failed to .Signatures: %v", err)
	}
	return nil
}
func (v ForwardMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *UpdateVerifier) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixUpdateVerifier); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.VerifierId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	if err = v.Settings.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Settings: %v", err)
	}
	return nil
}
func (v UpdateVerifier) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixUpdateVerifier, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.VerifierId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierId: %v", err)
	}
	if err = v.Settings.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Settings: %v", err)
	}
	return nil
}
func (v UpdateVerifier) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *RemoveVerifier) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixRemoveVerifier); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Id.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Id: %v", err)
	}
	return nil
}
func (v RemoveVerifier) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixRemoveVerifier, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Id.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Id: %v", err)
	}
	return nil
}
func (v RemoveVerifier) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VerifierRegistryInternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = VerifierRegistryInternalMessageKind(prefix)
	switch v.SumType {
	case VerifierRegistryInternalMessageKind_ForwardMessage:
		v.ForwardMessage = new(ForwardMessage)
		return v.ForwardMessage.UnmarshalTLB(c, decoder)
	case VerifierRegistryInternalMessageKind_UpdateVerifier:
		v.UpdateVerifier = new(UpdateVerifier)
		return v.UpdateVerifier.UnmarshalTLB(c, decoder)
	case VerifierRegistryInternalMessageKind_RemoveVerifier:
		v.RemoveVerifier = new(RemoveVerifier)
		return v.RemoveVerifier.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}
func (v VerifierRegistryInternalMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case VerifierRegistryInternalMessageKind_ForwardMessage:
		if v.ForwardMessage == nil {
			return fmt.Errorf("VerifierRegistryInternalMessage.ForwardMessage is nil")
		}
		return v.ForwardMessage.MarshalTLB(c, encoder)
	case VerifierRegistryInternalMessageKind_UpdateVerifier:
		if v.UpdateVerifier == nil {
			return fmt.Errorf("VerifierRegistryInternalMessage.UpdateVerifier is nil")
		}
		return v.UpdateVerifier.MarshalTLB(c, encoder)
	case VerifierRegistryInternalMessageKind_RemoveVerifier:
		if v.RemoveVerifier == nil {
			return fmt.Errorf("VerifierRegistryInternalMessage.RemoveVerifier is nil")
		}
		return v.RemoveVerifier.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown VerifierRegistryInternalMessage variant: %v", v.SumType)
	}
}
func (v VerifierRegistryInternalMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VerifierInfo) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Admin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Admin: %v", err)
	}
	if v.Settings, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Settings: %v", err)
	}
	if v.Found, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Found: %v", err)
	}
	return nil
}
func (v VerifierInfo) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Admin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Admin: %v", err)
	}
	if err = c.AddRef(&v.Settings); err != nil {
		return fmt.Errorf("failed to .Settings: %v", err)
	}
	if err = c.WriteBit(v.Found); err != nil {
		return fmt.Errorf("failed to .Found: %v", err)
	}
	return nil
}
func (v VerifierInfo) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *VerifierInfo) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.Found, err = stack.ReadBool(); err != nil {
		return fmt.Errorf("failed to read .Found: %v", err)
	}
	if v.Settings, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .Settings: %v", err)
	}
	if err = v.Admin.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Admin: %v", err)
	}
	return nil
}

func (msg ForwardMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VerifierRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg UpdateVerifier) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VerifierRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg RemoveVerifier) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VerifierRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg VerifierRegistryInternalMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*VerifierRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
