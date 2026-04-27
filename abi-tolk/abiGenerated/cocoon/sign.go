package abiCocoon

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func (msg ExternalSignedMessage) Sign(signer crypto.Signer) (WalletExternalMessage, error) {
	cell := boc.NewCell()
	err := msg.MarshalTLB(cell, &tlb.Encoder{})
	if err != nil {
		return WalletExternalMessage{}, fmt.Errorf("marshaling failed: %v", err)
	}

	digest, err := cell.Hash()
	if err != nil {
		return WalletExternalMessage{}, err
	}
	signature, err := signer.Sign(rand.Reader, digest, &ed25519.Options{})
	if err != nil {
		return WalletExternalMessage{}, err
	}
	if len(signature) != 512/8 {
		return WalletExternalMessage{}, fmt.Errorf("invalid signature length, expected %d, got %d", 512/8, len(signature))
	}
	return WalletExternalMessage{
		Signature: tlb.Bits512(signature),
		Message:   msg,
	}, nil
}
