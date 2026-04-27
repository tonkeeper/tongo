// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *TextCmd) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixTextCmd); err != nil {
		return err
	}
	if err = v.Action.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Action: %v", err)
	}
	return nil
}
func (v TextCmd) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixTextCmd, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.Action.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Action: %v", err)
	}
	return nil
}
func (v TextCmd) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtProxyCloseRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtProxyCloseRequestSigned); err != nil {
		return err
	}
	if v.Rest, err = (func() (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})(); err != nil {
		return fmt.Errorf("failed to read .Rest: %v", err)
	}
	return nil
}
func (v ExtProxyCloseRequestSigned) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtProxyCloseRequestSigned, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.Rest); err != nil {
		return fmt.Errorf("failed to .Rest: %v", err)
	}
	return nil
}
func (v ExtProxyCloseRequestSigned) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtProxyCloseCompleteRequestSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtProxyCloseCompleteRequestSigned); err != nil {
		return err
	}
	if v.Rest, err = (func() (boc.Cell, error) {
		remain := c.CopyRemaining()
		if remain == nil {
			return boc.Cell{}, nil
		}
		return *remain, nil
	})(); err != nil {
		return fmt.Errorf("failed to read .Rest: %v", err)
	}
	return nil
}
func (v ExtProxyCloseCompleteRequestSigned) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtProxyCloseCompleteRequestSigned, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.Rest); err != nil {
		return fmt.Errorf("failed to .Rest: %v", err)
	}
	return nil
}
func (v ExtProxyCloseCompleteRequestSigned) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtProxyPayoutRequest) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtProxyPayoutRequest); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtProxyPayoutRequest) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtProxyPayoutRequest, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtProxyPayoutRequest) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtProxyIncreaseStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtProxyIncreaseStake); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Grams.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Grams: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtProxyIncreaseStake) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtProxyIncreaseStake, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Grams.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Grams: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtProxyIncreaseStake) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerProxyClose) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerProxyClose); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerProxyClose) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerProxyClose, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerProxyClose) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ProxyStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.ProxyPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	if err = v.RootAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RootAddress: %v", err)
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.UnlockTs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.Params.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Params: %v", err)
	}
	return nil
}
func (v ProxyStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.ProxyPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyPublicKey: %v", err)
	}
	if err = v.RootAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RootAddress: %v", err)
	}
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.UnlockTs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockTs: %v", err)
	}
	if err = v.Params.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Params: %v", err)
	}
	return nil
}
func (v ProxyStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonProxyData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.ProxyPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	if err = v.RootAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .RootAddress: %v", err)
	}
	if err = v.State.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.UnlockTs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.PricePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.MinProxyStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.ParamsVersion.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	return nil
}
func (v CocoonProxyData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.ProxyPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyPublicKey: %v", err)
	}
	if err = v.RootAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .RootAddress: %v", err)
	}
	if err = v.State.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .State: %v", err)
	}
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.UnlockTs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockTs: %v", err)
	}
	if err = v.PricePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .PricePerToken: %v", err)
	}
	if err = v.WorkerFeePerToken.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .WorkerFeePerToken: %v", err)
	}
	if err = v.MinProxyStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinProxyStake: %v", err)
	}
	if err = v.MinClientStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .MinClientStake: %v", err)
	}
	if err = v.ParamsVersion.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ParamsVersion: %v", err)
	}
	return nil
}
func (v CocoonProxyData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonProxyData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.ParamsVersion.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ParamsVersion: %v", err)
	}
	if err = v.MinClientStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinClientStake: %v", err)
	}
	if err = v.MinProxyStake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .MinProxyStake: %v", err)
	}
	if err = v.WorkerFeePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .WorkerFeePerToken: %v", err)
	}
	if err = v.PricePerToken.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .PricePerToken: %v", err)
	}
	if err = v.UnlockTs.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.Stake.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.Balance.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.State.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .State: %v", err)
	}
	if err = v.RootAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .RootAddress: %v", err)
	}
	if err = v.ProxyPublicKey.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	if err = v.OwnerAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	return nil
}

func (msg TextCmd) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtProxyCloseRequestSigned) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtProxyCloseCompleteRequestSigned) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtProxyPayoutRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtProxyIncreaseStake) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerProxyClose) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg WorkerProxyRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ClientProxyRequest) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ProxyStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
