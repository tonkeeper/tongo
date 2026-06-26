// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *GramboDeploy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboDeploy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.Metadata, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Metadata: %v", err)
	}
	if err = v.TonTarget.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TonTarget: %v", err)
	}
	if err = v.PrebuyAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PrebuyAmount: %v", err)
	}
	return nil
}

func (v GramboDeploy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboDeploy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = c.AddRef(&v.Metadata); err != nil {
		return fmt.Errorf("failed to .Metadata: %v", err)
	}
	if err = v.TonTarget.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TonTarget: %v", err)
	}
	if err = v.PrebuyAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PrebuyAmount: %v", err)
	}
	return nil
}

func (v GramboDeploy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboUpdateData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboUpdateData); err != nil {
		return err
	}
	if v.NewData, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	return nil
}

func (v GramboUpdateData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboUpdateData, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.NewData); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	return nil
}

func (v GramboUpdateData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GramboUpdateCode) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixGramboUpdateCode); err != nil {
		return err
	}
	if v.NewCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	return nil
}

func (v GramboUpdateCode) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixGramboUpdateCode, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.NewCode); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	return nil
}

func (v GramboUpdateCode) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *FactoryIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = FactoryIncomingMessageKind(prefix)
	switch v.SumType {
	case FactoryIncomingMessageKind_GramboDeploy:
		v.GramboDeploy = new(GramboDeploy)
		return v.GramboDeploy.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_GramboUpdateData:
		v.GramboUpdateData = new(GramboUpdateData)
		return v.GramboUpdateData.UnmarshalTLB(c, decoder)
	case FactoryIncomingMessageKind_GramboUpdateCode:
		v.GramboUpdateCode = new(GramboUpdateCode)
		return v.GramboUpdateCode.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}

func (v FactoryIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case FactoryIncomingMessageKind_GramboDeploy:
		if v.GramboDeploy == nil {
			return fmt.Errorf("FactoryIncomingMessage.GramboDeploy is nil")
		}
		return v.GramboDeploy.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_GramboUpdateData:
		if v.GramboUpdateData == nil {
			return fmt.Errorf("FactoryIncomingMessage.GramboUpdateData is nil")
		}
		return v.GramboUpdateData.MarshalTLB(c, encoder)
	case FactoryIncomingMessageKind_GramboUpdateCode:
		if v.GramboUpdateCode == nil {
			return fmt.Errorf("FactoryIncomingMessage.GramboUpdateCode is nil")
		}
		return v.GramboUpdateCode.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown FactoryIncomingMessage variant: %v", v.SumType)
	}
}

func (v FactoryIncomingMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg GramboDeploy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboUpdateData) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg GramboUpdateCode) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg FactoryIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
