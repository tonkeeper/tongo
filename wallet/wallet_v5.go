package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type WalletV5ID struct {
	NetworkGlobalID uint32
	Workchain       uint8
	WalletVersion   uint8
	SubWalletID     uint32
}

type DataV5 struct {
	Seqno      tlb.Uint33
	WalletID   WalletV5ID
	PublicKey  tlb.Bits256
	PluginDict tlb.HashmapE[tlb.Bits264, tlb.Uint8] // TODO: find type and check size
}

type walletV5 struct {
	version         Version
	publicKey       ed25519.PublicKey
	workchain       int
	subWalletID     uint32
	networkGlobalID uint32
}

type V5MsgType uint32

const (
	V5MsgTypeInternal V5MsgType = 0x73696e74
	V5MsgTypeExternal V5MsgType = 0x7369676e
)

const TestnetGlobalID = -3
const MainnetGlobalID = -239

type MessageConfigV5 struct {
	MsgType V5MsgType
}

var _ wallet = &walletV5{}

func newWalletV5(version Version, publicKey ed25519.PublicKey, opts Options) *walletV5 {
	workchain := defaultOr(opts.Workchain, 0)
	subWalletID := defaultOr(opts.SubWalletID, 0)

	networkGlobalID := defaultOr[int32](opts.NetworkGlobalID, MainnetGlobalID)
	return &walletV5{
		version:         version,
		publicKey:       publicKey,
		workchain:       workchain,
		subWalletID:     subWalletID,
		networkGlobalID: uint32(networkGlobalID),
	}
}

func (w *walletV5) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not generate state init: %v", err)
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletV5) generateStateInit() (*tlb.StateInit, error) {
	data := DataV5{
		Seqno: 0,
		WalletID: WalletV5ID{
			NetworkGlobalID: w.networkGlobalID,
			Workchain:       uint8(w.workchain),
			SubWalletID:     w.subWalletID,
		},
		PublicKey: publicKeyToBits(w.publicKey),
	}
	return generateStateInit(w.version, data)
}

func (w *walletV5) maxMessageNumber() int {
	return 254
}

func (w *walletV5) nextMessageParams(state tlb.ShardAccount) (nextMsgParams, error) {
	if state.Account.Status() == tlb.AccountActive {
		var data DataV5
		cell := boc.Cell(state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value)
		if err := tlb.Unmarshal(&cell, &data); err != nil {
			return nextMsgParams{}, err
		}
		return nextMsgParams{
			Seqno: uint32(data.Seqno),
		}, nil
	}
	init, err := w.generateStateInit()
	if err != nil {
		return nextMsgParams{}, err
	}
	return nextMsgParams{Init: init}, nil

}

type extSignedMessage struct {
	WalletId   WalletV5ID
	ValidUntil uint32
	Seqno      uint32
	Op         bool
	Actions    SendMessageList `tlb:"^"`
}

func (w *walletV5) createSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	actions := SendMessageList{
		Actions: make([]SendMessageAction, 0, len(internalMessages)),
	}
	for _, msg := range internalMessages {
		actions.Actions = append(actions.Actions, SendMessageAction{
			Msg:  msg.Message,
			Mode: msg.Mode,
		})
	}
	msg := extSignedMessage{
		WalletId: WalletV5ID{
			NetworkGlobalID: w.networkGlobalID,
			Workchain:       uint8(w.workchain),
			SubWalletID:     w.subWalletID,
		},
		ValidUntil: uint32(msgConfig.ValidUntil.Unix()),
		Seqno:      msgConfig.Seqno,
		Op:         false,
		Actions:    actions,
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
