// Code generated - DO NOT EDIT.

package abiVerifier

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *SourcesRegistryStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MinGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinGram: %v", err)
	}
	if err = v.MaxGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxGram: %v", err)
	}
	if err = v.AdminAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .AdminAddress: %v", err)
	}
	if err = v.VerifierRegistryAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierRegistryAddress: %v", err)
	}
	if v.SourceItemCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .SourceItemCode: %v", err)
	}
	return nil
}

func (v SourcesRegistryStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MinGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinGram: %v", err)
	}
	if err = v.MaxGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxGram: %v", err)
	}
	if err = v.AdminAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .AdminAddress: %v", err)
	}
	if err = v.VerifierRegistryAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierRegistryAddress: %v", err)
	}
	if err = c.AddRef(&v.SourceItemCode); err != nil {
		return fmt.Errorf("failed to .SourceItemCode: %v", err)
	}
	return nil
}

func (v SourcesRegistryStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *DeploySourceItemPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixDeploySourceItemPayload); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.VerifierId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceContent.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SourceContent: %v", err)
	}
	return nil
}

func (v DeploySourceItemPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixDeploySourceItemPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.VerifierId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifierId: %v", err)
	}
	if err = v.VerifiedCodeCellHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VerifiedCodeCellHash: %v", err)
	}
	if err = v.SourceContent.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SourceContent: %v", err)
	}
	return nil
}

func (v DeploySourceItemPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ChangeVerifierRegistry) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeVerifierRegistry); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewVerifierRegistry.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewVerifierRegistry: %v", err)
	}
	return nil
}

func (v ChangeVerifierRegistry) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeVerifierRegistry, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewVerifierRegistry.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewVerifierRegistry: %v", err)
	}
	return nil
}

func (v ChangeVerifierRegistry) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *ChangeAdmin) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixChangeAdmin); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewAdmin.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewAdmin: %v", err)
	}
	return nil
}

func (v ChangeAdmin) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixChangeAdmin, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewAdmin.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewAdmin: %v", err)
	}
	return nil
}

func (v ChangeAdmin) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SetSourceItemCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSetSourceItemCode); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewSourceItemCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewSourceItemCode: %v", err)
	}
	return nil
}

func (v SetSourceItemCode) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSetSourceItemCode, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.NewSourceItemCode); err != nil {
		return fmt.Errorf("failed to .NewSourceItemCode: %v", err)
	}
	return nil
}

func (v SetSourceItemCode) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SetCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSetCode); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	return nil
}

func (v SetCode) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSetCode, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.NewCode); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	return nil
}

func (v SetCode) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SetDeploymentCosts) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixSetDeploymentCosts); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewMinGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewMinGram: %v", err)
	}
	if err = v.NewMaxGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewMaxGram: %v", err)
	}
	return nil
}

func (v SetDeploymentCosts) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixSetDeploymentCosts, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewMinGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewMinGram: %v", err)
	}
	if err = v.NewMaxGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewMaxGram: %v", err)
	}
	return nil
}

func (v SetDeploymentCosts) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *SourcesRegistryInternalMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = SourcesRegistryInternalMessageKind(prefix)
	switch v.SumType {
	case SourcesRegistryInternalMessageKind_DeploySourceItemPayload:
		v.DeploySourceItemPayload = new(DeploySourceItemPayload)
		return v.DeploySourceItemPayload.UnmarshalTLB(c, decoder)
	case SourcesRegistryInternalMessageKind_ChangeVerifierRegistry:
		v.ChangeVerifierRegistry = new(ChangeVerifierRegistry)
		return v.ChangeVerifierRegistry.UnmarshalTLB(c, decoder)
	case SourcesRegistryInternalMessageKind_ChangeAdmin:
		v.ChangeAdmin = new(ChangeAdmin)
		return v.ChangeAdmin.UnmarshalTLB(c, decoder)
	case SourcesRegistryInternalMessageKind_SetSourceItemCode:
		v.SetSourceItemCode = new(SetSourceItemCode)
		return v.SetSourceItemCode.UnmarshalTLB(c, decoder)
	case SourcesRegistryInternalMessageKind_SetCode:
		v.SetCode = new(SetCode)
		return v.SetCode.UnmarshalTLB(c, decoder)
	case SourcesRegistryInternalMessageKind_SetDeploymentCosts:
		v.SetDeploymentCosts = new(SetDeploymentCosts)
		return v.SetDeploymentCosts.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v SourcesRegistryInternalMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case SourcesRegistryInternalMessageKind_DeploySourceItemPayload:
		if v.DeploySourceItemPayload == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.DeploySourceItemPayload is nil")
		}
		return v.DeploySourceItemPayload.MarshalTLB(c, encoder)
	case SourcesRegistryInternalMessageKind_ChangeVerifierRegistry:
		if v.ChangeVerifierRegistry == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.ChangeVerifierRegistry is nil")
		}
		return v.ChangeVerifierRegistry.MarshalTLB(c, encoder)
	case SourcesRegistryInternalMessageKind_ChangeAdmin:
		if v.ChangeAdmin == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.ChangeAdmin is nil")
		}
		return v.ChangeAdmin.MarshalTLB(c, encoder)
	case SourcesRegistryInternalMessageKind_SetSourceItemCode:
		if v.SetSourceItemCode == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.SetSourceItemCode is nil")
		}
		return v.SetSourceItemCode.MarshalTLB(c, encoder)
	case SourcesRegistryInternalMessageKind_SetCode:
		if v.SetCode == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.SetCode is nil")
		}
		return v.SetCode.MarshalTLB(c, encoder)
	case SourcesRegistryInternalMessageKind_SetDeploymentCosts:
		if v.SetDeploymentCosts == nil {
			return fmt.Errorf("SourcesRegistryInternalMessage.SetDeploymentCosts is nil")
		}
		return v.SetDeploymentCosts.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown SourcesRegistryInternalMessage variant: %v", v.SumType)
	}
}

func (v SourcesRegistryInternalMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *DeploymentCosts) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.MinGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinGram: %v", err)
	}
	if err = v.MaxGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MaxGram: %v", err)
	}
	return nil
}

func (v DeploymentCosts) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.MinGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinGram: %v", err)
	}
	if err = v.MaxGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MaxGram: %v", err)
	}
	return nil
}

func (v DeploymentCosts) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *DeploymentCosts) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.MaxGram.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MaxGram: %v", err)
	}
	if err = v.MinGram.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinGram: %v", err)
	}
	return nil
}

func (msg DeploySourceItemPayload) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ChangeVerifierRegistry) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg ChangeAdmin) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SetSourceItemCode) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SetCode) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SetDeploymentCosts) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg SourcesRegistryInternalMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*SourcesRegistryStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
