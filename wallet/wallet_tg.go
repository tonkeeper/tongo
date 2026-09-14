package wallet

import (
	"crypto/ed25519"
	"fmt"

	abi "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/walletTg"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// The Telegram wallet (WalletTg)
//	https://github.com/ton-blockchain/tg-wallet-contract.

const (
	DefaultTgSubWalletIDMainnet = 0x7FFF7F11
	DefaultTgSubWalletIDTestnet = 0x7FFF7FFD
	TgMaxBatchSize              = tlb.MaxTolkArrayLen
)

// TgMsgType tells whether a message body is meant to be delivered
// as an external message (the regular path) or inside an internal one (gasless).
// Internal and external requests differ only in opcode, so that a request signed
// for a gasless relay can not be replayed as an external one.
type TgMsgType uint8

const (
	TgMsgTypeExternal TgMsgType = iota
	TgMsgTypeInternal
)

type walletTg struct {
	publicKey   ed25519.PublicKey
	subWalletID uint32
	workchain   int
}

func newWalletTg(publicKey ed25519.PublicKey, options Options) *walletTg {
	defaultSubWalletID := uint32(DefaultTgSubWalletIDMainnet)
	if options.NetworkGlobalID != nil && *options.NetworkGlobalID != MainnetGlobalID {
		defaultSubWalletID = DefaultTgSubWalletIDTestnet
	}
	return &walletTg{
		publicKey:   publicKey,
		subWalletID: defaultOr(options.SubWalletID, defaultSubWalletID),
		workchain:   defaultOr(options.Workchain, 0),
	}
}

func (w *walletTg) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("can not generate state init: %v", err)
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletTg) generateStateInit() (*tlb.StateInit, error) {
	storage := abi.Storage{
		Seqno:       0,
		SubwalletId: tlb.Uint32(w.subWalletID),
		PublicKey:   publicKeyToBits(w.publicKey).ToUint(),
	}
	return generateStateInit(WalletTg, storage)
}

func (w *walletTg) MaxMessageNumber() int {
	return TgMaxBatchSize
}

func (w *walletTg) GetPublicKey() ed25519.PublicKey {
	return w.publicKey
}

func (w *walletTg) NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error) {
	if state.Account.Status() == tlb.AccountActive {
		cell := state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value
		// todo storage layout might change on telegram wallet update
		// most safe option would be to run seqno() get method against the storage
		var storage abi.Storage
		if err := storage.UnmarshalTLB(&cell, tlb.NewDecoder()); err != nil {
			return NextMsgParams{}, err
		}
		return NextMsgParams{Seqno: uint32(storage.Seqno)}, nil
	}
	init, err := w.generateStateInit()
	if err != nil {
		return NextMsgParams{}, err
	}
	return NextMsgParams{Init: init}, nil
}

func (w *walletTg) header(msgConfig MessageConfig) abi.SeqnoHeader {
	return abi.SeqnoHeader{
		SubwalletId: tlb.Uint32(w.subWalletID),
		ValidUntil:  tlb.Uint32(msgConfig.ValidUntil.Unix()),
		Seqno:       tlb.Uint32(msgConfig.Seqno),
	}
}

// CreateMsgBodyWithoutSignature builds a request WalletTg understands.
// A single message goes with the cheaper "send one message" opcode,
// anything else as an array<MessageToSend>.
func (w *walletTg) CreateMsgBodyWithoutSignature(internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	if len(internalMessages) == 0 {
		return nil, fmt.Errorf("tg wallet requires at least one message to send")
	}
	if len(internalMessages) > TgMaxBatchSize {
		return nil, fmt.Errorf("tg wallet supports up to %v internal messages", TgMaxBatchSize)
	}
	external := msgConfig.TgMsgType == TgMsgTypeExternal
	msgArr := make(abi.RawArrayOfMessagesToSend, 0, len(internalMessages))
	for _, msg := range internalMessages {
		if external && !IsMessageModeSet(int(msg.Mode), IgnoreErrors) {
			return nil, fmt.Errorf("tg wallet requires the ignore-errors mode for external messages")
		}
		msgArr = append(msgArr, abi.MessageToSend{
			SendMode:    tlb.Uint8(msg.Mode),
			MessageCell: *msg.Message,
		})
	}
	header := w.header(msgConfig)
	if external {
		if len(msgArr) == 1 {
			return abi.AllowedExternalMsg{
				SumType:                abi.AllowedExternalMsgKind_SendOneMessageRequestE,
				SendOneMessageRequestE: &abi.SendOneMessageRequestE{Header: header, Msg: msgArr[0]},
			}.ToCell()
		}
		return abi.AllowedExternalMsg{
			SumType:                  abi.AllowedExternalMsgKind_SendBulkMessagesRequestE,
			SendBulkMessagesRequestE: &abi.SendBulkMessagesRequestE{Header: header, MsgArr: msgArr},
		}.ToCell()
	}
	if len(msgArr) == 1 {
		return abi.AllowedInternalMsg{
			SumType:                abi.AllowedInternalMsgKind_SendOneMessageRequestI,
			SendOneMessageRequestI: &abi.SendOneMessageRequestI{Header: header, Msg: msgArr[0]},
		}.ToCell()
	}
	return abi.AllowedInternalMsg{
		SumType:                  abi.AllowedInternalMsgKind_SendBulkMessagesRequestI,
		SendBulkMessagesRequestI: &abi.SendBulkMessagesRequestI{Header: header, MsgArr: msgArr},
	}.ToCell()
}

// AttachSignature prepends the signature to the request, producing a SignedRequest.
// Note that this differs from wallet v5, which keeps the signature at the end of the body.
func (w *walletTg) AttachSignature(body *boc.Cell, signature tlb.Bits512) (*boc.Cell, error) {
	cell := boc.NewCell()
	if err := tlb.Marshal(cell, signature); err != nil {
		return nil, err
	}
	if err := cell.WriteBitString(body.RawBitString()); err != nil {
		return nil, err
	}
	for _, ref := range body.Refs() {
		if err := cell.AddRef(ref); err != nil {
			return nil, err
		}
	}
	return cell, nil
}

var _ wallet = &walletTg{}
