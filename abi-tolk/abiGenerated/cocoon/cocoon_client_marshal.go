// Code generated - DO NOT EDIT.

package abiCocoon

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *ExtClientChargeSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtClientChargeSigned); err != nil {
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
func (v ExtClientChargeSigned) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtClientChargeSigned, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.Rest); err != nil {
		return fmt.Errorf("failed to .Rest: %v", err)
	}
	return nil
}
func (v ExtClientChargeSigned) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtClientGrantRefundSigned) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtClientGrantRefundSigned); err != nil {
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
func (v ExtClientGrantRefundSigned) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtClientGrantRefundSigned, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = c.AddRef(&v.Rest); err != nil {
		return fmt.Errorf("failed to .Rest: %v", err)
	}
	return nil
}
func (v ExtClientGrantRefundSigned) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ExtClientTopUp) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixExtClientTopUp); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TopUpAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TopUpAmount: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtClientTopUp) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixExtClientTopUp, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TopUpAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TopUpAmount: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v ExtClientTopUp) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientChangeSecretHashAndTopUp) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientChangeSecretHashAndTopUp); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.TopUpAmount.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TopUpAmount: %v", err)
	}
	if err = v.NewSecretHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewSecretHash: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientChangeSecretHashAndTopUp) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientChangeSecretHashAndTopUp, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.TopUpAmount.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TopUpAmount: %v", err)
	}
	if err = v.NewSecretHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewSecretHash: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientChangeSecretHashAndTopUp) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientRegister) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientRegister); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.Nonce.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Nonce: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientRegister) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientRegister, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.Nonce.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Nonce: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientRegister) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientChangeSecretHash) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientChangeSecretHash); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewSecretHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewSecretHash: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientChangeSecretHash) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientChangeSecretHash, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewSecretHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewSecretHash: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientChangeSecretHash) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientIncreaseStake) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientIncreaseStake); err != nil {
		return err
	}
	if err = v.QueryId.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .QueryId: %v", err)
	}
	if err = v.NewStake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .NewStake: %v", err)
	}
	if err = v.SendExcessesTo.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientIncreaseStake) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientIncreaseStake, 32); err != nil {
		return fmt.Errorf("failed to write prefix: %v", err)
	}
	if err = v.QueryId.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .QueryId: %v", err)
	}
	if err = v.NewStake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .NewStake: %v", err)
	}
	if err = v.SendExcessesTo.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SendExcessesTo: %v", err)
	}
	return nil
}
func (v OwnerClientIncreaseStake) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientWithdraw) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientWithdraw); err != nil {
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
func (v OwnerClientWithdraw) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientWithdraw, 32); err != nil {
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
func (v OwnerClientWithdraw) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *OwnerClientRequestRefund) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err := c.ReadPrefix(32, PrefixOwnerClientRequestRefund); err != nil {
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
func (v OwnerClientRequestRefund) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteUint(PrefixOwnerClientRequestRefund, 32); err != nil {
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
func (v OwnerClientRequestRefund) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClientConstData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.ProxyAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyAddress: %v", err)
	}
	if err = v.ProxyPublicKey.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ProxyPublicKey: %v", err)
	}
	return nil
}
func (v ClientConstData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.ProxyAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyAddress: %v", err)
	}
	if err = v.ProxyPublicKey.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ProxyPublicKey: %v", err)
	}
	return nil
}
func (v ClientConstData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClientStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
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
	if err = v.UnlockTs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.SecretHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SecretHash: %v", err)
	}
	if err = v.ConstDataRef.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .ConstDataRef: %v", err)
	}
	if err = v.Params.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Params: %v", err)
	}
	return nil
}
func (v ClientStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
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
	if err = v.UnlockTs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockTs: %v", err)
	}
	if err = v.SecretHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SecretHash: %v", err)
	}
	if err = v.ConstDataRef.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .ConstDataRef: %v", err)
	}
	if err = v.Params.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Params: %v", err)
	}
	return nil
}
func (v ClientStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *ClientMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	prefix, err := c.PickUint(32)
	if err != nil {
		return err
	}
	v.SumType = ClientMessageKind(prefix)
	switch v.SumType {
	case ClientMessageKind_ExtClientChargeSigned:
		v.ExtClientChargeSigned = new(ExtClientChargeSigned)
		return v.ExtClientChargeSigned.UnmarshalTLB(c, decoder)
	case ClientMessageKind_ExtClientGrantRefundSigned:
		v.ExtClientGrantRefundSigned = new(ExtClientGrantRefundSigned)
		return v.ExtClientGrantRefundSigned.UnmarshalTLB(c, decoder)
	case ClientMessageKind_ExtClientTopUp:
		v.ExtClientTopUp = new(ExtClientTopUp)
		return v.ExtClientTopUp.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientChangeSecretHashAndTopUp:
		v.OwnerClientChangeSecretHashAndTopUp = new(OwnerClientChangeSecretHashAndTopUp)
		return v.OwnerClientChangeSecretHashAndTopUp.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientRegister:
		v.OwnerClientRegister = new(OwnerClientRegister)
		return v.OwnerClientRegister.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientChangeSecretHash:
		v.OwnerClientChangeSecretHash = new(OwnerClientChangeSecretHash)
		return v.OwnerClientChangeSecretHash.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientIncreaseStake:
		v.OwnerClientIncreaseStake = new(OwnerClientIncreaseStake)
		return v.OwnerClientIncreaseStake.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientWithdraw:
		v.OwnerClientWithdraw = new(OwnerClientWithdraw)
		return v.OwnerClientWithdraw.UnmarshalTLB(c, decoder)
	case ClientMessageKind_OwnerClientRequestRefund:
		v.OwnerClientRequestRefund = new(OwnerClientRequestRefund)
		return v.OwnerClientRequestRefund.UnmarshalTLB(c, decoder)
	default:
		return fmt.Errorf("unknown prefix: %x", prefix)
	}
}
func (v ClientMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	switch v.SumType {
	case ClientMessageKind_ExtClientChargeSigned:
		if v.ExtClientChargeSigned == nil {
			return fmt.Errorf("ClientMessage.ExtClientChargeSigned is nil")
		}
		return v.ExtClientChargeSigned.MarshalTLB(c, encoder)
	case ClientMessageKind_ExtClientGrantRefundSigned:
		if v.ExtClientGrantRefundSigned == nil {
			return fmt.Errorf("ClientMessage.ExtClientGrantRefundSigned is nil")
		}
		return v.ExtClientGrantRefundSigned.MarshalTLB(c, encoder)
	case ClientMessageKind_ExtClientTopUp:
		if v.ExtClientTopUp == nil {
			return fmt.Errorf("ClientMessage.ExtClientTopUp is nil")
		}
		return v.ExtClientTopUp.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientChangeSecretHashAndTopUp:
		if v.OwnerClientChangeSecretHashAndTopUp == nil {
			return fmt.Errorf("ClientMessage.OwnerClientChangeSecretHashAndTopUp is nil")
		}
		return v.OwnerClientChangeSecretHashAndTopUp.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientRegister:
		if v.OwnerClientRegister == nil {
			return fmt.Errorf("ClientMessage.OwnerClientRegister is nil")
		}
		return v.OwnerClientRegister.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientChangeSecretHash:
		if v.OwnerClientChangeSecretHash == nil {
			return fmt.Errorf("ClientMessage.OwnerClientChangeSecretHash is nil")
		}
		return v.OwnerClientChangeSecretHash.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientIncreaseStake:
		if v.OwnerClientIncreaseStake == nil {
			return fmt.Errorf("ClientMessage.OwnerClientIncreaseStake is nil")
		}
		return v.OwnerClientIncreaseStake.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientWithdraw:
		if v.OwnerClientWithdraw == nil {
			return fmt.Errorf("ClientMessage.OwnerClientWithdraw is nil")
		}
		return v.OwnerClientWithdraw.MarshalTLB(c, encoder)
	case ClientMessageKind_OwnerClientRequestRefund:
		if v.OwnerClientRequestRefund == nil {
			return fmt.Errorf("ClientMessage.OwnerClientRequestRefund is nil")
		}
		return v.OwnerClientRequestRefund.MarshalTLB(c, encoder)
	default:
		return fmt.Errorf("unknown ClientMessage variant: %v", v.SumType)
	}
}
func (v ClientMessage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonClientData) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
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
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Stake.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Stake: %v", err)
	}
	if err = v.TokensUsed.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .TokensUsed: %v", err)
	}
	if err = v.UnlockTs.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.SecretHash.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .SecretHash: %v", err)
	}
	return nil
}
func (v CocoonClientData) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
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
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Stake.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Stake: %v", err)
	}
	if err = v.TokensUsed.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .TokensUsed: %v", err)
	}
	if err = v.UnlockTs.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .UnlockTs: %v", err)
	}
	if err = v.SecretHash.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .SecretHash: %v", err)
	}
	return nil
}
func (v CocoonClientData) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}
func (v *CocoonClientData) ReadFromStack(stack *tlb.VmStack) (err error) {
	if err = v.SecretHash.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .SecretHash: %v", err)
	}
	if err = v.UnlockTs.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .UnlockTs: %v", err)
	}
	if err = v.TokensUsed.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .TokensUsed: %v", err)
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

func (msg ExtClientChargeSigned) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtClientGrantRefundSigned) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ExtClientTopUp) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientChangeSecretHashAndTopUp) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientRegister) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientChangeSecretHash) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientIncreaseStake) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientWithdraw) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg OwnerClientRequestRefund) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}

func (msg ClientMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*ClientStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(msg, dest, amount, bounce, init)
}
