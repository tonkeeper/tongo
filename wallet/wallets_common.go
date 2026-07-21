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
	MaxMessageNumber() int
	CreateMsgBodyWithoutSignature(internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error)
	AttachSignature(body *boc.Cell, signature tlb.Bits512) (*boc.Cell, error)
	NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error)
	GetPublicKey() ed25519.PublicKey
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
		return ton.AccountID{}, fmt.Errorf("can not calculate state init hash: %v", err)
	}
	var hash tlb.Bits256
	copy(hash[:], h[:])
	return ton.AccountID{Workchain: int32(workchain), Address: hash}, nil
}

func attachSignatureToBody(body *boc.Cell, signature tlb.Bits512) (*boc.Cell, error) {
	cell := boc.NewCell()
	if err := cell.WriteBitString(body.RawBitString()); err != nil {
		return nil, err
	}
	for _, ref := range body.Refs() {
		if err := cell.AddRef(ref); err != nil {
			return nil, err
		}
	}
	if err := cell.WriteBytes(signature[:]); err != nil {
		return nil, err
	}
	return cell, nil
}

func attachSignatureAsSignedMsgBody(body *boc.Cell, signature tlb.Bits512) (*boc.Cell, error) {
	signedBody := SignedMsgBody{
		Sign:    signature,
		Message: tlb.Any(*body),
	}
	cell := boc.NewCell()
	if err := tlb.Marshal(cell, signedBody); err != nil {
		return nil, err
	}
	return cell, nil
}
