// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ClientProxyRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixClientProxyRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.StateData.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StateData: %v", err)
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}
func (v ClientProxyRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixClientProxyRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.StateData.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StateData: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}
func (v ClientProxyRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClientStateData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.TokensUsed.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensUsed: %v", err)
	}
	if err = v.SecretHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SecretHash: %v", err)
	}
	return nil
}
func (v ClientStateData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.TokensUsed.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensUsed: %v", err)
	}
	if err = v.SecretHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SecretHash: %v", err)
	}
	return nil
}
func (v ClientStateData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonParams) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.StructVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .StructVersion: %v", err)
	}
	if err = v.ParamsVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UniqueId: %v", err)
	}
	if v.IsTest, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .IsTest: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.PromptTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PromptTokensPriceMultiplier: %v", err)
	}
	if err = v.CachedTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CachedTokensPriceMultiplier: %v", err)
	}
	if err = v.CompletionTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .CompletionTokensPriceMultiplier: %v", err)
	}
	if err = v.ReasoningTokensPriceMultiplier.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ReasoningTokensPriceMultiplier: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if v.ProxyScCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .ProxyScCode: %v", err)
	}
	if v.WorkerScCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .WorkerScCode: %v", err)
	}
	if v.ClientScCode, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .ClientScCode: %v", err)
	}
	return nil
}
func (v CocoonParams) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.StructVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .StructVersion: %v", err)
	}
	if err = v.ParamsVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ParamsVersion: %v", err)
	}
	if err = v.UniqueId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UniqueId: %v", err)
	}
	if err = c.WriteBit(v.IsTest); err != nil {
		return fmt.Errorf("failed to .IsTest: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	if err = v.PromptTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PromptTokensPriceMultiplier: %v", err)
	}
	if err = v.CachedTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CachedTokensPriceMultiplier: %v", err)
	}
	if err = v.CompletionTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .CompletionTokensPriceMultiplier: %v", err)
	}
	if err = v.ReasoningTokensPriceMultiplier.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ReasoningTokensPriceMultiplier: %v", err)
	}
	if err = v.ProxyDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyDelayBeforeClose: %v", err)
	}
	if err = v.ClientDelayBeforeClose.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ClientDelayBeforeClose: %v", err)
	}
	if err = v.MinProxyStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinClientStake: %v", err)
	}
	if err = v.ProxyScCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyScCode: %v", err)
	}
	if err = v.WorkerScCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerScCode: %v", err)
	}
	if err = v.ClientScCode.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ClientScCode: %v", err)
	}
	return nil
}
func (v CocoonParams) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *Payout) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixPayout); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v Payout) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixPayout, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v Payout) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ReturnExcessesBack) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixReturnExcessesBack); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	return nil
}
func (v ReturnExcessesBack) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixReturnExcessesBack, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	return nil
}
func (v ReturnExcessesBack) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *WorkerProxyRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixWorkerProxyRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.Tokens.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Tokens: %v", err)
	}
	if v.Payload, err = tlb.UnmarshalMaybeCallback(c, func(c *boc.Cell) (boc.Cell, error) {
		return c.NextRefV()
	}); err != nil {
		return fmt.Errorf("failed to read .Payload: %v", err)
	}
	return nil
}
func (v WorkerProxyRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixWorkerProxyRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.Tokens.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Tokens: %v", err)
	}
	if err = v.Payload.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Payload: %v", err)
	}
	return nil
}
func (v WorkerProxyRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
