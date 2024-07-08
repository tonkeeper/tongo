package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type DataV5R1 struct {
	IsSignatureAllowed bool
	Seqno              uint32
	WalletID           uint32
	PublicKey          tlb.Bits256
	Extensions         tlb.HashmapE[tlb.Bits256, tlb.Uint1]
}

type walletV5R1 struct {
	publicKey          ed25519.PublicKey
	isSignatureAllowed bool
	walletID           uint32
	workchain          int
}

func genContextID(workchain uint32) uint32 {
	contextCell := boc.NewCell()
	if err := contextCell.WriteUint(1, 1); err != nil {
		panic(err)
	}
	if err := contextCell.WriteUint(uint64(workchain), 8); err != nil {
		panic(err)
	}
	if err := contextCell.WriteUint(0, 8); err != nil {
		panic(err)
	}
	if err := contextCell.WriteUint(0, 15); err != nil {
		panic(err)
	}
	contextCell.ResetCounters()
	cid, err := contextCell.ReadUint(32)
	if err != nil {
		panic(err)
	}
	return uint32(cid)
}

func NewWalletV5R1(publicKey ed25519.PublicKey, opts Options) *walletV5R1 {
	workchain := defaultOr(opts.Workchain, 0)
	networkGlobalID := int64(defaultOr[int32](opts.NetworkGlobalID, MainnetGlobalID))
	contextID := int64(genContextID(uint32(workchain)))
	walletID := contextID ^ networkGlobalID

	return &walletV5R1{
		publicKey:          publicKey,
		workchain:          workchain,
		walletID:           uint32(walletID), // todo: add options to configure wallet id
		isSignatureAllowed: true,             // todo: add option to disable signature
	}
}
func (w *walletV5R1) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not generate state init: %v", err)
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletV5R1) generateStateInit() (*tlb.StateInit, error) {
	data := DataV5R1{
		IsSignatureAllowed: w.isSignatureAllowed,
		Seqno:              0,
		WalletID:           w.walletID,
		PublicKey:          publicKeyToBits(w.publicKey),
	}
	return generateStateInit(V5R1, data)
}

func (w *walletV5R1) maxMessageNumber() int {
	return 255
}

type extV5R1SignedMessage struct {
	WalletId        uint32
	ValidUntil      uint32
	Seqno           uint32
	Actions         *W5Actions         `tlb:"maybe^"`
	ExtendedActions *W5ExtendedActions `tlb:"maybe"`
}

func (w *walletV5R1) CreateSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, extensionsActions *W5ExtendedActions, msgConfig MessageConfig) (*boc.Cell, error) {
	actions := make([]W5SendMessageAction, 0, len(internalMessages))
	for _, msg := range internalMessages {
		actions = append(actions, W5SendMessageAction{
			Msg:  msg.Message,
			Mode: msg.Mode,
		})
	}
	w5Actions := W5Actions(actions)
	msg := extV5R1SignedMessage{
		WalletId:        w.walletID,
		ValidUntil:      uint32(msgConfig.ValidUntil.Unix()),
		Seqno:           msgConfig.Seqno,
		Actions:         &w5Actions,
		ExtendedActions: extensionsActions,
	}
	bodyCell := boc.NewCell()
	if err := bodyCell.WriteUint(uint64(msgConfig.V5MsgType), 32); err != nil {
		return nil, err
	}
	if err := tlb.Marshal(bodyCell, msg); err != nil {
		return nil, err
	}
	signature, err := bodyCell.Sign(privateKey)
	if err != nil {
		return nil, err
	}
	if err := bodyCell.WriteBytes(signature); err != nil {
		return nil, err
	}
	return bodyCell, nil
}

func (w *walletV5R1) createSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	return w.CreateSignedMsgBodyCell(privateKey, internalMessages, nil, msgConfig)
}

func (w *walletV5R1) NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error) {
	if state.Account.Status() == tlb.AccountActive {
		var data DataV5R1
		cell := boc.Cell(state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value)
		if err := tlb.Unmarshal(&cell, &data); err != nil {
			return NextMsgParams{}, err
		}
		return NextMsgParams{
			Seqno: uint32(data.Seqno),
		}, nil
	}
	init, err := w.generateStateInit()
	if err != nil {
		return NextMsgParams{}, err
	}
	return NextMsgParams{Init: init}, nil
}

// GetW5R1ExtensionsList returns a list of wallet v5 extensions added to a specific wallet.
func GetW5R1ExtensionsList(state tlb.ShardAccount, workchain int) (map[ton.AccountID]struct{}, error) {
	if state.Account.Status() == tlb.AccountActive {
		var data DataV5R1
		cell := boc.Cell(state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value)
		if err := tlb.Unmarshal(&cell, &data); err != nil {
			return nil, err
		}
		if len(data.Extensions.Keys()) == 0 {
			return nil, nil
		}
		extensions := make(map[ton.AccountID]struct{}, len(data.Extensions.Keys()))
		for _, item := range data.Extensions.Items() {
			ext := ton.AccountID{
				Workchain: int32(workchain),
				Address:   item.Key,
			}
			extensions[ext] = struct{}{}
		}
		return extensions, nil
	}
	return nil, nil
}

var _ wallet = &walletV5R1{}
