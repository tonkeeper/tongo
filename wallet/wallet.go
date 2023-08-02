package wallet

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func DefaultWalletFromSeed(seed string, blockchain blockchain) (Wallet, error) {
	pk, err := SeedToPrivateKey(seed)
	if err != nil {
		return Wallet{}, err
	}
	return New(pk, V4R2, 0, nil, blockchain)
}

// New
// Fill new Wallet struct from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func New(key ed25519.PrivateKey, ver Version, workchain int, subWalletId *int, blockchain blockchain) (Wallet, error) {
	publicKey := key.Public().(ed25519.PublicKey)
	address, err := GenerateWalletAddress(publicKey, ver, workchain, subWalletId)
	if err != nil {
		return Wallet{}, err
	}
	if (ver == V3R1 || ver == V3R2 || ver == V4R1 || ver == V4R2 || ver == HighLoadV2R2) && subWalletId == nil {
		id := DefaultSubWallet + workchain
		subWalletId = &id
	}
	w := Wallet{
		address:     address,
		key:         key,
		ver:         ver,
		subWalletId: uint32(*subWalletId),
		blockchain:  blockchain,
	}
	return w, nil
}

// GenerateWalletAddress
// Generate wallet address from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func GenerateWalletAddress(
	key ed25519.PublicKey,
	ver Version,
	workchain int,
	subWalletId *int,
) (tongo.AccountID, error) {
	state, err := GenerateStateInit(key, ver, workchain, subWalletId)
	if err != nil {
		return tongo.AccountID{}, fmt.Errorf("can not generate wallet state: %v", err)
	}
	stateCell := boc.NewCell()
	err = tlb.Marshal(stateCell, state)
	if err != nil {
		return tongo.AccountID{}, fmt.Errorf("can not marshal wallet state: %v", err)
	}
	h, err := stateCell.Hash()
	if err != nil {
		return tongo.AccountID{}, err
	}
	var hash tlb.Bits256
	copy(hash[:], h[:])
	if err != nil {
		return tongo.AccountID{}, fmt.Errorf("can not calculate state init hash: %v", err)
	}
	return tongo.AccountID{Workchain: int32(workchain), Address: hash}, nil
}

func GenerateStateInit(
	key ed25519.PublicKey,
	ver Version,
	workchain int,
	subWalletId *int,
) (tlb.StateInit, error) {
	var (
		err       error
		publicKey tlb.Bits256
	)
	copy(publicKey[:], key[:])
	dataCell := boc.NewCell()
	switch ver {
	case V1R1, V1R2, V1R3, V2R1, V2R2:
		data := DataV1V2{0, publicKey}
		err = tlb.Marshal(dataCell, data)
	case V3R1, V3R2:
		if subWalletId == nil {
			id := DefaultSubWallet + workchain
			subWalletId = &id
		}
		data := DataV3{0, uint32(*subWalletId), publicKey}
		err = tlb.Marshal(dataCell, data)
	case V4R1, V4R2:
		if subWalletId == nil {
			id := DefaultSubWallet + workchain
			subWalletId = &id
		}
		data := DataV4{
			Seqno:       0,
			SubWalletId: uint32(*subWalletId),
			PublicKey:   publicKey,
		}
		err = tlb.Marshal(dataCell, data)
	case HighLoadV2R2:
		if subWalletId == nil {
			id := DefaultSubWallet + workchain
			subWalletId = &id
		}
		data := DataHighloadV4{
			SubWalletId: uint32(*subWalletId),
			PublicKey:   publicKey,
		}
		err = tlb.Marshal(dataCell, data)
	default:
		return tlb.StateInit{}, fmt.Errorf("address generation not implemented for this wallet ver")
	}
	if err != nil {
		return tlb.StateInit{}, fmt.Errorf("wallet data marshaling error: %v", err)
	}

	codeCell := GetCodeByVer(ver)

	state := tlb.StateInit{
		Code: tlb.Maybe[tlb.Ref[boc.Cell]]{Exists: true, Value: tlb.Ref[boc.Cell]{Value: *codeCell}},
		Data: tlb.Maybe[tlb.Ref[boc.Cell]]{Exists: true, Value: tlb.Ref[boc.Cell]{Value: *dataCell}},
		// Library: empty by default
	}
	return state, nil
}

func (w *Wallet) RawSendV2(
	ctx context.Context,
	seqno uint32,
	validUntil time.Time,
	internalMessages []RawMessage,
	init *tlb.StateInit,
	waitingConfirmation time.Duration,
) (tongo.Bits256, error) {
	if w.blockchain == nil {
		return tongo.Bits256{}, tongo.BlockchainInterfaceIsNil
	}
	err := checkMessagesLimit(len(internalMessages), w.ver)
	if err != nil {
		return tongo.Bits256{}, err
	}
	bodyCell := boc.NewCell()
	switch w.ver {
	case V3R1, V3R2:
		body := MessageV3{
			SubWalletId: w.subWalletId,
			ValidUntil:  uint32(validUntil.Unix()),
			Seqno:       seqno,
			RawMessages: PayloadV1toV4(internalMessages),
		}
		err = tlb.Marshal(bodyCell, body)
	case V4R1, V4R2:
		body := MessageV4{
			SubWalletId: w.subWalletId,
			ValidUntil:  uint32(validUntil.Unix()),
			Seqno:       seqno,
			Op:          0,
			RawMessages: PayloadV1toV4(internalMessages),
		}
		err = tlb.Marshal(bodyCell, body)
	case HighLoadV2R2:
		boundedID := uint64(time.Now().Add(DefaultMessageLifetime).UTC().Unix()<<32) + uint64(rand.Uint32())
		body := HighloadV2Message{
			SubWalletId:    w.subWalletId,
			BoundedQueryID: boundedID,
			RawMessages:    PayloadHighload(internalMessages),
		}
		err = tlb.Marshal(bodyCell, body)
	default:
		return tongo.Bits256{}, fmt.Errorf("message body generation for this wallet is not supported: %v", err)
	}
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not marshal wallet message body: %v", err)
	}

	signBytes, err := bodyCell.Sign(w.key)
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not sign wallet message body: %v", err)
	}
	bits512 := tlb.Bits512{}
	copy(bits512[:], signBytes[:])
	signedBody := SignedMsgBody{
		Sign:    bits512,
		Message: tlb.Any(*bodyCell),
	}
	signedBodyCell := boc.NewCell()
	if err = tlb.Marshal(signedBodyCell, signedBody); err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not marshal signed body: %v", err)
	}
	extMsg, err := tongo.CreateExternalMessage(w.address, signedBodyCell, init, 0)
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not create external message: %v", err)
	}
	extMsgCell := boc.NewCell()
	err = tlb.Marshal(extMsgCell, extMsg)
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not marshal wallet external message: %v", err)
	}
	msgHash, err := extMsgCell.Hash256()
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not create external message: %v", err)
	}
	payload, err := extMsgCell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return tongo.Bits256{}, fmt.Errorf("can not serialize external message cell: %v", err)
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

func (w *Wallet) getInit() (tlb.StateInit, error) {
	publicKey := w.key.Public().(ed25519.PublicKey)
	id := int(w.subWalletId)
	return GenerateStateInit(publicKey, w.ver, int(w.address.Workchain), &id)
}

func checkMessagesLimit(msgQty int, ver Version) error { // TODO: maybe return bool
	switch ver {
	case V1R1, V1R2, V1R3, V2R1, V2R2, V3R1, V3R2, V4R1, V4R2:
		if msgQty > 4 {
			return fmt.Errorf("%v wallet support up to 4 internal messages", ver.ToString())
		}
	case HighLoadV2R2:
		if msgQty > 254 {
			return fmt.Errorf("%v wallet support up to 254 internal messages", ver.ToString())
		}
	default:
		return fmt.Errorf("message qty checking is not implemented for %v wallet", ver.ToString())
	}
	return nil
}

func (w *Wallet) SendV2(
	ctx context.Context,
	waitingConfirmation time.Duration,
	messages ...Sendable,
) (tongo.Bits256, error) {
	if w.blockchain == nil {
		return tongo.Bits256{}, tongo.BlockchainInterfaceIsNil
	}

	var init *tlb.StateInit
	var seqno uint32

	state, err := w.blockchain.GetAccountState(ctx, w.GetAddress())
	if err != nil {
		return tongo.Bits256{}, err
	}
	requireInit := false
	switch w.ver {
	case HighLoadV2R2:
		// this wallet have no seqno.
		requireInit = state.Account.Status() == tlb.AccountUninit || state.Account.Status() == tlb.AccountNone
	default:
		if state.Account.Status() == tlb.AccountActive {
			seqno, err = w.blockchain.GetSeqno(ctx, w.address)
			if err != nil {
				return tongo.Bits256{}, err
			}
		}
		requireInit = seqno == 0
	}
	if requireInit {
		i, err := w.getInit()
		if err != nil {
			return tongo.Bits256{}, err
		}
		init = &i
	}
	var msgArray []RawMessage
	for _, m := range messages {
		intMsg, mode, err := m.ToInternal()
		if err != nil {
			return tongo.Bits256{}, err
		}
		cell := boc.NewCell()
		err = tlb.Marshal(cell, intMsg)
		if err != nil {
			return tongo.Bits256{}, err
		}
		msgArray = append(msgArray, RawMessage{Message: cell, Mode: mode})
	}
	validUntil := time.Now().Add(DefaultMessageLifetime)
	return w.RawSendV2(ctx, seqno, validUntil, msgArray, init, waitingConfirmation)
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
		return 0, tongo.BlockchainInterfaceIsNil
	}
	state, err := w.blockchain.GetAccountState(ctx, w.address)
	if err != nil {
		return 0, err
	}
	accInfo, err := tongo.GetAccountInfo(state.Account)
	if err != nil {
		return 0, err
	}
	return accInfo.Balance, nil
}
