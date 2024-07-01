package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type NextMsgParams struct {
	Seqno uint32
	Init  *tlb.StateInit
}

type wallet interface {
	generateAddress() (ton.AccountID, error)
	generateStateInit() (*tlb.StateInit, error)
	maxMessageNumber() int
	createSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error)
	NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error)
}

func defaultOr[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}

func publicKeyToBits(key ed25519.PublicKey) tlb.Bits256 {
	var publicKey tlb.Bits256
	copy(publicKey[:], key[:])
	return publicKey
}

func generateStateInit(ver Version, data any) (*tlb.StateInit, error) {
	dataCell := boc.NewCell()
	if err := tlb.Marshal(dataCell, data); err != nil {
		return nil, err
	}
	codeCell := GetCodeByVer(ver)
	state := tlb.StateInit{
		Code: tlb.Maybe[tlb.Ref[boc.Cell]]{Exists: true, Value: tlb.Ref[boc.Cell]{Value: *codeCell}},
		Data: tlb.Maybe[tlb.Ref[boc.Cell]]{Exists: true, Value: tlb.Ref[boc.Cell]{Value: *dataCell}},
		// Library: empty by default
	}
	return &state, nil
}

func generateAddress(workchain int, stateInit tlb.StateInit) (ton.AccountID, error) {
	stateCell := boc.NewCell()
	err := tlb.Marshal(stateCell, stateInit)
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not marshal wallet state: %v", err)
	}
	h, err := stateCell.Hash()
	if err != nil {
		return ton.AccountID{}, err
	}
	var hash tlb.Bits256
	copy(hash[:], h[:])
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not calculate state init hash: %v", err)
	}
	return ton.AccountID{Workchain: int32(workchain), Address: hash}, nil
}

func signBodyCell(bodyCell boc.Cell, privateKey ed25519.PrivateKey) (*boc.Cell, error) {
	signBytes, err := bodyCell.Sign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("can not sign wallet message body: %v", err)
	}
	bits512 := tlb.Bits512{}
	copy(bits512[:], signBytes[:])
	signedBody := SignedMsgBody{
		Sign:    bits512,
		Message: tlb.Any(bodyCell),
	}
	signedBodyCell := boc.NewCell()
	if err := tlb.Marshal(signedBodyCell, signedBody); err != nil {
		return nil, fmt.Errorf("can not marshal signed body: %v", err)
	}
	return signedBodyCell, nil
}
