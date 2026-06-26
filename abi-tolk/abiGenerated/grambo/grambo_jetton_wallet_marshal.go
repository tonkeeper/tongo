// Code generated - DO NOT EDIT.

package abiGrambo

import (
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (v *WalletIncomingMessage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	var vx GramboActivateWallet
	if err := vx.UnmarshalTLB(c, decoder); err != nil {
		return err
	}
	*v = WalletIncomingMessage(vx)
	return nil
}

func (v WalletIncomingMessage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	return GramboActivateWallet(v).MarshalTLB(c, encoder)
}

func (v *GramboWalletStorage) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if v.Activated, err = c.ReadBit(); err != nil {
		return fmt.Errorf("failed to read .Activated: %v", err)
	}
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.Owner.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Owner: %v", err)
	}
	if err = v.Master.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Master: %v", err)
	}
	return nil
}

func (v GramboWalletStorage) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = c.WriteBit(v.Activated); err != nil {
		return fmt.Errorf("failed to .Activated: %v", err)
	}
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.Owner.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Owner: %v", err)
	}
	if err = v.Master.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Master: %v", err)
	}
	return nil
}

func (v GramboWalletStorage) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetWalletDataResult) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) (err error) {
	if err = v.Balance.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	if err = v.OwnerAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.JettonMasterAddress.UnmarshalTLB(c, decoder); err != nil {
		return fmt.Errorf("failed to read .JettonMasterAddress: %v", err)
	}
	if v.JettonWalletCode, err = c.NextRefV(); err != nil {
		return fmt.Errorf("failed to read .JettonWalletCode: %v", err)
	}
	return nil
}

func (v GetWalletDataResult) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) (err error) {
	if err = v.Balance.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .Balance: %v", err)
	}
	if err = v.OwnerAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .OwnerAddress: %v", err)
	}
	if err = v.JettonMasterAddress.MarshalTLB(c, encoder); err != nil {
		return fmt.Errorf("failed to .JettonMasterAddress: %v", err)
	}
	if err = c.AddRef(&v.JettonWalletCode); err != nil {
		return fmt.Errorf("failed to .JettonWalletCode: %v", err)
	}
	return nil
}

func (v GetWalletDataResult) ToCell() (*boc.Cell, error) {
	c := boc.NewCell()
	if err := v.MarshalTLB(c, &tlb.Encoder{}); err != nil {
		return nil, err
	}
	return c, nil
}

func (v *GetWalletDataResult) ReadFromStack(stack *tlb.VmStack) (err error) {
	if v.JettonWalletCode, err = stack.ReadCell(); err != nil {
		return fmt.Errorf("failed to read .JettonWalletCode: %v", err)
	}
	if err = v.JettonMasterAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .JettonMasterAddress: %v", err)
	}
	if err = v.OwnerAddress.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .OwnerAddress: %v", err)
	}
	if err = v.Balance.ReadFromStack(stack); err != nil {
		return fmt.Errorf("failed to read .Balance: %v", err)
	}
	return nil
}

func (msg WalletIncomingMessage) ToInternal(dest tlb.InternalAddress, amount tlb.Grams, bounce bool, init *tlb.StateInitT[*GramboWalletStorage]) (tlb.Message, error) {
	return tlb.BuildInternal(&msg, dest, amount, bounce, init)
}
