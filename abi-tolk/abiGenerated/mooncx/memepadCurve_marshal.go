// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *EconState) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.TokensForSale.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensForSale: %v", err)
	}
	if err = v.TokensForLp.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensForLp: %v", err)
	}
	if err = v.VirtualGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VirtualGram: %v", err)
	}
	if err = v.VirtualTokens.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .VirtualTokens: %v", err)
	}
	if err = v.TargetGram.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TargetGram: %v", err)
	}
	if err = v.TradeFeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TradeFeeBps: %v", err)
	}
	if err = v.CreatorFeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CreatorFeeBps: %v", err)
	}
	if err = v.PoolFeeBps.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PoolFeeBps: %v", err)
	}
	return nil
}
func (v EconState) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.TokensForSale.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensForSale: %v", err)
	}
	if err = v.TokensForLp.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensForLp: %v", err)
	}
	if err = v.VirtualGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VirtualGram: %v", err)
	}
	if err = v.VirtualTokens.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .VirtualTokens: %v", err)
	}
	if err = v.TargetGram.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TargetGram: %v", err)
	}
	if err = v.TradeFeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TradeFeeBps: %v", err)
	}
	if err = v.CreatorFeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CreatorFeeBps: %v", err)
	}
	if err = v.PoolFeeBps.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PoolFeeBps: %v", err)
	}
	return nil
}
func (v EconState) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveInit) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveInit); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.CpmmFactory.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CpmmFactory: %v", err)
	}
	if err = v.ProtocolRecipient.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProtocolRecipient: %v", err)
	}
	if err = v.InitialBuy.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .InitialBuy: %v", err)
	}
	if err = v.Econ.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Econ: %v", err)
	}
	if v.JwCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .JwCode: %v", err)
	}
	if v.Metadata, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .Metadata: %v", err)
	}
	return nil
}
func (v MemepadCurveInit) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveInit, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.CpmmFactory.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CpmmFactory: %v", err)
	}
	if err = v.ProtocolRecipient.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProtocolRecipient: %v", err)
	}
	if err = v.InitialBuy.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .InitialBuy: %v", err)
	}
	if err = v.Econ.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Econ: %v", err)
	}
	if err = c.AddRef(&v.JwCode); err != nil {
		return fmt.Errorf("failed to .JwCode: %v", err)
	}
	if err = c.AddRef(&v.Metadata); err != nil {
		return fmt.Errorf("failed to .Metadata: %v", err)
	}
	return nil
}
func (v MemepadCurveInit) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveDeployError) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveDeployError); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveDeployError) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveDeployError, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveDeployError) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveBuy) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveBuy); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.MinOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinOut: %v", err)
	}
	return nil
}
func (v MemepadCurveBuy) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveBuy, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.MinOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinOut: %v", err)
	}
	return nil
}
func (v MemepadCurveBuy) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveSellPayload) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveSellPayload); err != nil {
		return err
	}
	if err = v.MinOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinOut: %v", err)
	}
	return nil
}
func (v MemepadCurveSellPayload) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveSellPayload, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.MinOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinOut: %v", err)
	}
	return nil
}
func (v MemepadCurveSellPayload) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveSellPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveSellPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.GramOut.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .GramOut: %v", err)
	}
	return nil
}
func (v MemepadCurveSellPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveSellPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.GramOut.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .GramOut: %v", err)
	}
	return nil
}
func (v MemepadCurveSellPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveExcess) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveExcess); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveExcess) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveExcess, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveExcess) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveCreatorClaimFees) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveCreatorClaimFees); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveCreatorClaimFees) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveCreatorClaimFees, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveCreatorClaimFees) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveCreatorFeesPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveCreatorFeesPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveCreatorFeesPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveCreatorFeesPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveCreatorFeesPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveProtocolPayout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveProtocolPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveProtocolPayout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveProtocolPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v MemepadCurveProtocolPayout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *MemepadCurveUpgrade) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixMemepadCurveUpgrade); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if v.NewCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewCode: %v", err)
	}
	if v.NewData, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .NewData: %v", err)
	}
	return nil
}
func (v MemepadCurveUpgrade) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixMemepadCurveUpgrade, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewCode: %v", err)
	}
	if err = v.NewData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewData: %v", err)
	}
	return nil
}
func (v MemepadCurveUpgrade) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (msg MemepadCurveInit) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MemepadCurveBuy) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MemepadCurveCreatorClaimFees) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}

func (msg MemepadCurveUpgrade) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*tlb.Any]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
