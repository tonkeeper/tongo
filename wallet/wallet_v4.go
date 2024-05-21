package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type DataV4 struct {
	Seqno       uint32
	SubWalletId uint32
	PublicKey   tlb.Bits256
	PluginDict  tlb.HashmapE[tlb.Bits264, tlb.Any] // TODO: find type and check size
}

type walletV4 struct {
	version     Version
	publicKey   ed25519.PublicKey
	workchain   int
	subWalletID uint32
}

var _ wallet = &walletV4{}

func newWalletV4(version Version, publicKey ed25519.PublicKey, opts Options) *walletV4 {
	workchain := defaultOr(opts.Workchain, 0)
	subWalletID := defaultOr(opts.SubWalletID, uint32(DefaultSubWallet+workchain))
	return &walletV4{
		version:     version,
		publicKey:   publicKey,
		workchain:   workchain,
		subWalletID: subWalletID,
	}
}

func (w *walletV4) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not generate state init: %v", err)
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletV4) generateStateInit() (*tlb.StateInit, error) {
	data := DataV4{
		Seqno:       0,
		SubWalletId: w.subWalletID,
		PublicKey:   publicKeyToBits(w.publicKey),
	}
	return generateStateInit(w.version, data)
}

func (w *walletV4) maxMessageNumber() int {
	return 4
}

func (w *walletV4) createSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	body := MessageV4{
		SubWalletId: w.subWalletID,
		ValidUntil:  uint32(msgConfig.ValidUntil.Unix()),
		Seqno:       msgConfig.Seqno,
		Op:          0,
		RawMessages: PayloadV1toV4(internalMessages),
	}
	bodyCell := boc.NewCell()
	if err := tlb.Marshal(bodyCell, body); err != nil {
		return nil, err
	}
	return signBodyCell(*bodyCell, privateKey)
}

func (w *walletV4) nextMessageParams(state tlb.ShardAccount) (nextMsgParams, error) {
	if state.Account.Status() == tlb.AccountActive {
		var data DataV4
		cell := boc.Cell(state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value)
		if err := tlb.Unmarshal(&cell, &data); err != nil {
			return nextMsgParams{}, err
		}
		return nextMsgParams{
			Seqno: data.Seqno,
		}, nil
	}
	init, err := w.generateStateInit()
	if err != nil {
		return nextMsgParams{}, err
	}
	return nextMsgParams{Init: init}, nil
}
