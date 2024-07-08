package wallet

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Wallet struct {
	ver        Version
	key        ed25519.PrivateKey
	address    ton.AccountID
	intWallet  wallet
	blockchain blockchain

	msgDefaultLifetime time.Duration
}

func DefaultWalletFromSeed(seed string, blockchain blockchain) (Wallet, error) {
	pk, err := SeedToPrivateKey(seed)
	if err != nil {
		return Wallet{}, err
	}
	return New(pk, V4R2, blockchain)
}

type Options struct {
	NetworkGlobalID *int32
	Workchain       *int
	SubWalletID     *uint32
	PublicKey       *ed25519.PublicKey
	PrivateKey      *ed25519.PrivateKey
	MsgLifetime     time.Duration
}

type Option func(*Options)

func WithNetworkGlobalID(id int32) Option {
	return func(o *Options) {
		o.NetworkGlobalID = &id
	}
}

func WithWorkchain(workchain int) Option {
	return func(o *Options) {
		o.Workchain = &workchain
	}
}

func WithSubWalletID(id uint32) Option {
	return func(o *Options) {
		o.SubWalletID = &id
	}
}

func WithMessageLifetime(lifetime time.Duration) Option {
	return func(o *Options) {
		o.MsgLifetime = lifetime
	}
}

func applyOptions(opts ...Option) Options {
	options := Options{
		MsgLifetime: DefaultMessageLifetime,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// New
// Fill new Wallet struct from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func New(key ed25519.PrivateKey, ver Version, blockchain blockchain, opts ...Option) (Wallet, error) {
	options := applyOptions(opts...)
	w, err := newWallet(key.Public().(ed25519.PublicKey), ver, options)
	if err != nil {
		return Wallet{}, err
	}
	address, err := w.generateAddress()
	if err != nil {
		return Wallet{}, fmt.Errorf("can not generate wallet address: %v", err)
	}
	return Wallet{
		address:            address,
		key:                key,
		ver:                ver,
		intWallet:          w,
		blockchain:         blockchain,
		msgDefaultLifetime: options.MsgLifetime,
	}, nil
}

func newWallet(key ed25519.PublicKey, version Version, options Options) (wallet, error) {
	switch version {
	case V1R1, V1R2, V1R3, V2R1, V2R2:
		return newWalletV1V2(version, key, options), nil
	case V3R1, V3R2:
		return newWalletV3(version, key, options), nil
	case V4R1, V4R2:
		return newWalletV4(version, key, options), nil
	case V5Beta:
		return NewWalletV5Beta(version, key, options), nil
	case V5R1:
		return NewWalletV5R1(key, options), nil
	case HighLoadV2R2:
		return newWalletHighloadV2(version, key, options), nil
	default:
		return nil, fmt.Errorf("unsupported wallet version: %v", version)
	}
}

// GenerateWalletAddress
// Generate wallet address from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func GenerateWalletAddress(
	key ed25519.PublicKey,
	ver Version,
	networkGlobalID *int32,
	workchain int,
	subWalletId *uint32,
) (ton.AccountID, error) {
	options := []Option{
		WithWorkchain(workchain),
	}
	if networkGlobalID != nil {
		options = append(options, WithNetworkGlobalID(*networkGlobalID))
	}
	if subWalletId != nil {
		options = append(options, WithSubWalletID(*subWalletId))
	}
	w, err := newWallet(key, ver, applyOptions(options...))
	if err != nil {
		return ton.AccountID{}, err
	}
	return w.generateAddress()
}

func GenerateStateInit(
	key ed25519.PublicKey,
	ver Version,
	networkGlobalID *int32,
	workchain int,
	subWalletId *uint32,
) (tlb.StateInit, error) {
	options := []Option{
		WithWorkchain(workchain),
	}
	if networkGlobalID != nil {
		options = append(options, WithNetworkGlobalID(*networkGlobalID))
	}
	if subWalletId != nil {
		options = append(options, WithSubWalletID(*subWalletId))
	}
	w, err := newWallet(key, ver, applyOptions(options...))
	if err != nil {
		return tlb.StateInit{}, nil
	}
	state, err := w.generateStateInit()
	if err != nil {
		return tlb.StateInit{}, err
	}
	return *state, nil
}

func (w *Wallet) RawSendV2(
	ctx context.Context,
	seqno uint32,
	validUntil time.Time,
	internalMessages []RawMessage,
	init *tlb.StateInit,
	waitingConfirmation time.Duration,
) (ton.Bits256, error) {
	if w.blockchain == nil {
		return ton.Bits256{}, errors.New("blockchain interface is nil")
	}
	if len(internalMessages) > w.intWallet.maxMessageNumber() {
		return ton.Bits256{}, fmt.Errorf("%v wallet support up to %v internal messages", w.ver, w.intWallet.maxMessageNumber())
	}
	msgConfig := MessageConfig{
		Seqno:      seqno,
		ValidUntil: validUntil,
		V5MsgType:  V5MsgTypeSignedExternal,
	}
	signedBodyCell, err := w.intWallet.createSignedMsgBodyCell(w.key, internalMessages, msgConfig)
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("can not marshal wallet message body: %v", err)
	}
	extMsg, err := ton.CreateExternalMessage(w.address, signedBodyCell, init, tlb.VarUInteger16{})
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("can not create external message: %v", err)
	}
	extMsgCell := boc.NewCell()
	err = tlb.Marshal(extMsgCell, extMsg)
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("can not marshal wallet external message: %v", err)
	}
	msgHash, err := extMsgCell.Hash256()
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("can not create external message: %v", err)
	}
	payload, err := extMsgCell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("can not serialize external message cell: %v", err)
	}
	t := time.Now()
	_, err = w.blockchain.SendMessage(ctx, payload) // TODO: add result code check
	if err != nil {
		return msgHash, err
	}
	if waitingConfirmation == 0 {
		return msgHash, nil
	}
	if w.ver == HighLoadV2R2 {
		return msgHash, fmt.Errorf("highload wallet doesn't support waiting confirmation")
	}

	for ; time.Since(t) < waitingConfirmation; time.Sleep(waitingConfirmation / 10) {
		newSeqno, err := w.blockchain.GetSeqno(ctx, w.address)
		if err == nil {
			continue
		}
		if newSeqno >= seqno {
			return msgHash, nil //todo: check if it is the same message
		}
	}
	return msgHash, fmt.Errorf("waiting confirmation timeout")
}

// RawSend
// Generates a signed external message for wallet with custom internal messages, seqno, TTL and init
// The payload is serialized into bytes and sent by the method SendRawMessage
func (w *Wallet) RawSend(
	ctx context.Context,
	seqno uint32,
	validUntil time.Time,
	internalMessages []RawMessage,
	init *tlb.StateInit,
) error {
	_, err := w.RawSendV2(ctx, seqno, validUntil, internalMessages, init, 0)
	return err
}

func (w *Wallet) SendV2(ctx context.Context, waitingConfirmation time.Duration, messages ...Sendable) (ton.Bits256, error) {
	if w.blockchain == nil {
		return ton.Bits256{}, errors.New("blockchain interface is nil")
	}
	state, err := w.blockchain.GetAccountState(ctx, w.GetAddress())
	if err != nil {
		return ton.Bits256{}, fmt.Errorf("get account state failed: %v", err)
	}
	params, err := w.intWallet.NextMessageParams(state)
	if err != nil {
		return ton.Bits256{}, err
	}
	msgArray := make([]RawMessage, 0, len(messages))
	for _, m := range messages {
		intMsg, mode, err := m.ToInternal()
		if err != nil {
			return ton.Bits256{}, err
		}
		cell := boc.NewCell()
		if err := tlb.Marshal(cell, intMsg); err != nil {
			return ton.Bits256{}, err
		}
		msgArray = append(msgArray, RawMessage{Message: cell, Mode: mode})
	}
	validUntil := time.Now().Add(w.msgDefaultLifetime)
	return w.RawSendV2(ctx, params.Seqno, validUntil, msgArray, params.Init, waitingConfirmation)
}

// Send
// Generates a signed external message for wallet with custom internal messages and default TTL
// Gets actual seqno and attach init for wallet if it needed
// The payload is serialized into bytes and sent by the method SendRawMessage
func (w *Wallet) Send(ctx context.Context, messages ...Sendable) error {
	_, err := w.SendV2(ctx, 0, messages...)
	return err
}

// GetBalance
// Gets actual TON balance for wallet
func (w *Wallet) GetBalance(ctx context.Context) (uint64, error) {
	if w.blockchain == nil {
		return 0, errors.New("blockchain interface is nil")
	}
	state, err := w.blockchain.GetAccountState(ctx, w.address)
	if err != nil {
		return 0, err
	}
	return uint64(state.Account.Account.Storage.Balance.Grams), nil
}

// GetAddress returns current wallet address but you can also call function GenerateWalletAddress
// which returns same address but doesn't require blockchain connection for calling
func (w *Wallet) GetAddress() ton.AccountID {
	return w.address
}

func (w *Wallet) StateInit() (*tlb.StateInit, error) {
	return w.intWallet.generateStateInit()
}

type MessageConfig struct {
	Seqno      uint32
	ValidUntil time.Time
	V5MsgType  V5MsgType
}

func (w *Wallet) CreateMessageBody(msgConfig MessageConfig, messages ...Sendable) (*boc.Cell, error) {
	msgArray := make([]RawMessage, 0, len(messages))
	for _, m := range messages {
		intMsg, mode, err := m.ToInternal()
		if err != nil {
			return nil, err
		}
		cell := boc.NewCell()
		if err := tlb.Marshal(cell, intMsg); err != nil {
			return nil, err
		}
		msgArray = append(msgArray, RawMessage{Message: cell, Mode: mode})
	}
	if msgConfig.ValidUntil.IsZero() {
		msgConfig.ValidUntil = time.Now().Add(w.msgDefaultLifetime)
	}
	signedBodyCell, err := w.intWallet.createSignedMsgBodyCell(w.key, msgArray, msgConfig)
	if err != nil {
		return nil, err
	}
	return signedBodyCell, nil
}
