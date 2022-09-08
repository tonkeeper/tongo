package wallet

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

func DefaultWalletFromSeed(seed string, version Version) (Wallet, error) {
	pk, err := SeedToPrivateKey(seed)
	if err != nil {
		return Wallet{}, err
	}
	return NewWallet(pk, version, 0, nil)
}

func SeedToPrivateKey(seed string) (ed25519.PrivateKey, error) {
	s := strings.Split(seed, " ")
	if len(s) < 12 {
		return nil, fmt.Errorf("seed should have at least 12 words")
	}
	mac := hmac.New(sha512.New, []byte(strings.Join(s, " ")))
	hash := mac.Sum(nil)
	p := pbkdf2.Key(hash, []byte("TON seed version"), 100000/256, 1, sha512.New)
	if p[0] != 0 {
		return nil, errors.New("invalid seed")
	}
	pk := pbkdf2.Key(hash, []byte("TON default seed"), 100000, 32, sha512.New)
	privateKey := ed25519.NewKeyFromSeed(pk)
	return privateKey, nil
}

// NewWallet
// Fill new Wallet struct from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func NewWallet(key ed25519.PrivateKey, ver Version, workchain int, subWalletId *int) (Wallet, error) {
	publicKey := key.Public().(ed25519.PublicKey)
	address, err := GenerateWalletAddress(publicKey, ver, workchain, subWalletId)
	if err != nil {
		return Wallet{}, err
	}
	if (ver == V3R1 || ver == V3R2 || ver == V4R1 || ver == V4R2) && subWalletId == nil {
		id := DefaultSubWalletIdV3V4 + workchain
		subWalletId = &id
	}
	return Wallet{
		address:     address,
		key:         key,
		ver:         ver,
		subWalletId: uint32(*subWalletId),
	}, nil
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
	state, err := generateStateInit(key, ver, workchain, subWalletId)
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
	var hash tongo.Hash
	copy(hash[:], h[:])
	if err != nil {
		return tongo.AccountID{}, fmt.Errorf("can not calculate state init hash: %v", err)
	}
	return tongo.AccountID{Workchain: int32(workchain), Address: hash}, nil
}

func generateStateInit(
	key ed25519.PublicKey,
	ver Version,
	workchain int,
	subWalletId *int,
) (tongo.StateInit, error) {
	var (
		err       error
		publicKey tongo.Hash
	)
	copy(publicKey[:], key[:])
	dataCell := boc.NewCell()
	switch ver {
	case V1R1, V1R2, V1R3, V2R1, V2R2:
		data := DataV1V2{0, publicKey}
		err = tlb.Marshal(dataCell, data)
	case V3R1, V3R2:
		if subWalletId == nil {
			id := DefaultSubWalletIdV3V4 + workchain
			subWalletId = &id
		}
		data := DataV3{0, uint32(*subWalletId), publicKey}
		err = tlb.Marshal(dataCell, data)
	case V4R1, V4R2:
		if subWalletId == nil {
			id := DefaultSubWalletIdV3V4 + workchain
			subWalletId = &id
		}
		data := DataV4{
			0,
			uint32(*subWalletId),
			publicKey,
			tlb.HashmapE[tlb.Any]{}, // TODO: clarify type
		}
		err = tlb.Marshal(dataCell, data)
	default:
		return tongo.StateInit{}, fmt.Errorf("address generation not implemented for this wallet ver")
	}
	if err != nil {
		return tongo.StateInit{}, fmt.Errorf("wallet data marshaling error: %v", err)
	}

	codeCell := GetCodeByVer(ver)

	state := tongo.StateInit{
		Special: tlb.Maybe[tongo.TickTock]{Null: true},
		Code:    tlb.Maybe[tlb.Ref[boc.Cell]]{Null: false, Value: tlb.Ref[boc.Cell]{Value: *codeCell}},
		Data:    tlb.Maybe[tlb.Ref[boc.Cell]]{Null: false, Value: tlb.Ref[boc.Cell]{Value: *dataCell}},
		// Library: empty by default
	}
	state.SplitDepth.Null = true
	return state, nil
}

// GenerateMessage
// Generate signed external message for wallet with custom internal messages and init
// Payload serialized to bytes and can be sent by liteclient.SendRawMessage method
func (w Wallet) GenerateMessage(seqno, validUntil uint32, internalMessages MessageArray, init *tongo.StateInit) ([]byte, error) {
	var err error
	err = checkMessagesLimit(len(internalMessages.Messages), w.ver)
	if err != nil {
		return nil, err
	}
	bodyCell := boc.NewCell()
	switch w.ver {
	case V3R1, V3R2:
		body := MessageV3{
			SubWalletId: w.subWalletId,
			ValidUntil:  validUntil,
			Seqno:       seqno,
			Messages:    internalMessages,
		}
		err = tlb.Marshal(bodyCell, body)
	case V4R1, V4R2:
		body := MessageV4{
			SubWalletId: w.subWalletId,
			ValidUntil:  validUntil,
			Seqno:       seqno,
			Op:          0,
			Messages:    internalMessages,
		}
		err = tlb.Marshal(bodyCell, body)
	default:
		return nil, fmt.Errorf("message body generation for this wallet is not supported: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("can not marshal wallet message body: %v", err)
	}

	sign, err := bodyCell.Sign(w.key)
	if err != nil {
		return nil, fmt.Errorf("can not sign wallet message body: %v", err)
	}
	signedBody := struct {
		Sign    boc.BitString `tlb:"512bits"`
		Payload tlb.Any
	}{
		Sign:    sign,
		Payload: tlb.Any(*bodyCell),
	}
	signedBodyCell := boc.NewCell()
	err = tlb.Marshal(signedBodyCell, signedBody)
	extMsg, err := tongo.CreateExternalMessage(w.address, signedBodyCell, init, 0)
	extMsgCell := boc.NewCell()
	err = tlb.Marshal(extMsgCell, extMsg)
	if err != nil {
		return nil, fmt.Errorf("can not marshal wallet external message: %v", err)
	}
	payload, err := extMsgCell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return nil, fmt.Errorf("can not serialize external message cell: %v", err)
	}
	return payload, err
}

// GenerateTonTransferMessage
// Generate signed external message for transfer TONs.
// Payload serialized to bytes and can be sent by liteclient.SendRawMessage method
func (w Wallet) GenerateTonTransferMessage(seqno, validUntil uint32, tonTransfers []TonTransfer) ([]byte, error) {
	if len(tonTransfers) == 0 {
		return nil, fmt.Errorf("need at least one transfer")
	}
	var msgs MessageArray
	for _, transfer := range tonTransfers {
		body := boc.NewCell()
		err := tlb.Marshal(body, TextComment(transfer.Comment))
		if err != nil {
			return nil, err
		}
		info := tongo.CommonMsgInfo{
			SumType: "IntMsgInfo",
		}
		info.IntMsgInfo.IhrDisabled = true // TODO: maybe need to add to TonTransfer struct
		info.IntMsgInfo.Bounce = transfer.Bounce
		info.IntMsgInfo.Src = tongo.MsgAddressFromAccountID(nil)
		info.IntMsgInfo.Dest = tongo.MsgAddressFromAccountID(&transfer.Recipient)
		info.IntMsgInfo.Value.Grams = transfer.Amount
		intMsg := tongo.Message[tlb.Any]{
			Info: info,
			Body: tlb.Either[tlb.Any, tlb.Ref[tlb.Any]]{
				IsRight: true,
				Right:   tlb.Ref[tlb.Any]{Value: tlb.Any(*body)},
			},
		}
		intMsg.Init.Null = true
		msgs.Messages = append(msgs.Messages, struct {
			Message tongo.Message[tlb.Any]
			Mode    byte
		}{
			Message: intMsg,
			Mode:    transfer.Mode,
		})
	}
	return w.GenerateMessage(seqno, validUntil, msgs, nil)
}

// GenerateDeployMessage
// Generate signed external message for deploying new wallet.
// Payload serialized to bytes and can be sent by liteclient.SendRawMessage method
func (w Wallet) GenerateDeployMessage() ([]byte, error) {
	publicKey := w.key.Public().(ed25519.PublicKey)
	id := int(w.subWalletId)
	init, err := generateStateInit(publicKey, w.ver, int(w.address.Workchain), &id)
	if err != nil {
		return nil, err
	}
	return w.GenerateMessage(0, 0xFFFFFFFF, MessageArray{}, &init)
}

func checkMessagesLimit(msgQty int, ver Version) error { // TODO: maybe return bool
	switch ver {
	case V1R1, V1R2, V1R3, V2R1, V2R2, V3R1, V3R2, V4R1, V4R2:
		if msgQty > 4 {
			return fmt.Errorf("%v wallet support up to 4 internal messages", ver.ToString())
		}
	default:
		return fmt.Errorf("message qty checking is not implemented for %v wallet", ver.ToString())
	}
	return nil
}
