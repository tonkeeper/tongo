package wallet

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"time"

	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/contract/jetton"
	"github.com/startfellows/tongo/tlb"
)

func DefaultWalletFromSeed(seed string, blockchain blockchain) (Wallet, error) {
	pk, err := SeedToPrivateKey(seed)
	if err != nil {
		return Wallet{}, err
	}
	return NewWallet(pk, V4R2, 0, nil, blockchain)
}

// NewWallet
// Fill new Wallet struct from known workchain, public key and version.
// subWalletId is only used in V3 and V4 wallets. Use nil for default value.
// The version number is associated with a specific implementation of the wallet code
// (https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)
func NewWallet(key ed25519.PrivateKey, ver Version, workchain int, subWalletId *int, blockchain blockchain) (Wallet, error) {
	publicKey := key.Public().(ed25519.PublicKey)
	address, err := GenerateWalletAddress(publicKey, ver, workchain, subWalletId)
	if err != nil {
		return Wallet{}, err
	}
	if (ver == V3R1 || ver == V3R2 || ver == V4R1 || ver == V4R2) && subWalletId == nil {
		id := DefaultSubWalletIdV3V4 + workchain
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
	var hash tlb.Bits256
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
			Seqno:       0,
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
		Special: tlb.Maybe[tlb.TickTock]{Null: true},
		Code:    tlb.Maybe[tlb.Ref[boc.Cell]]{Null: false, Value: tlb.Ref[boc.Cell]{Value: *codeCell}},
		Data:    tlb.Maybe[tlb.Ref[boc.Cell]]{Null: false, Value: tlb.Ref[boc.Cell]{Value: *dataCell}},
		// Library: empty by default
	}
	state.SplitDepth.Null = true
	return state, nil
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
	if w.blockchain == nil {
		return tongo.BlockchainInterfaceIsNil
	}
	err := checkMessagesLimit(len(internalMessages), w.ver)
	if err != nil {
		return nil
	}
	bodyCell := boc.NewCell()
	switch w.ver {
	case V3R1, V3R2:
		body := MessageV3{
			SubWalletId: w.subWalletId,
			ValidUntil:  uint32(validUntil.Unix()),
			Seqno:       seqno,
			Payload:     PayloadV1toV4(internalMessages),
		}
		err = tlb.Marshal(bodyCell, body)
	case V4R1, V4R2:
		body := MessageV4{
			SubWalletId: w.subWalletId,
			ValidUntil:  uint32(validUntil.Unix()),
			Seqno:       seqno,
			Op:          0,
			Payload:     PayloadV1toV4(internalMessages),
		}
		err = tlb.Marshal(bodyCell, body)
	default:
		return fmt.Errorf("message body generation for this wallet is not supported: %v", err)
	}
	if err != nil {
		return fmt.Errorf("can not marshal wallet message body: %v", err)
	}

	signBytes, err := bodyCell.Sign(w.key)
	if err != nil {
		return fmt.Errorf("can not sign wallet message body: %v", err)
	}
	bits512 := tlb.Bits512{}
	copy(bits512[:], signBytes[:])
	signedBody := struct {
		Sign    tlb.Bits512
		Payload tlb.Any
	}{
		Sign:    bits512,
		Payload: tlb.Any(*bodyCell),
	}
	signedBodyCell := boc.NewCell()
	err = tlb.Marshal(signedBodyCell, signedBody)
	extMsg, err := tongo.CreateExternalMessage(w.address, signedBodyCell, init, 0)
	extMsgCell := boc.NewCell()
	err = tlb.Marshal(extMsgCell, extMsg)
	if err != nil {
		return fmt.Errorf("can not marshal wallet external message: %v", err)
	}
	payload, err := extMsgCell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return fmt.Errorf("can not serialize external message cell: %v", err)
	}
	_, err = w.blockchain.SendMessage(ctx, payload) // TODO: add result code check
	return err
}

func generateInternalMessage(ctx context.Context, msg Message, bc blockchain) (tlb.Message, error) {
	body := boc.NewCell()
	if msg.Comment != nil && msg.Body != nil {
		return tlb.Message{}, fmt.Errorf("only body or comment must be presented")
	} else if msg.Comment != nil {
		err := tlb.Marshal(body, TextComment(*msg.Comment))
		if err != nil {
			return tlb.Message{}, err
		}
	} else if msg.Body != nil {
		body = msg.Body
	}
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}
	info.IntMsgInfo.IhrDisabled = true
	info.IntMsgInfo.Src = tongo.MsgAddressFromAccountID(nil)
	info.IntMsgInfo.Dest = tongo.MsgAddressFromAccountID(&msg.Address)
	info.IntMsgInfo.Value.Grams = tlb.Grams(msg.Amount)

	if msg.Bounceable == nil {
		destinationState, err := bc.GetAccountState(ctx, msg.Address)
		if err != nil {
			return tlb.Message{}, err
		}
		if destinationState.Account.SumType != "AccountNone" || destinationState.Account.Account.Storage.State.SumType == "AccountActive" {
			info.IntMsgInfo.Bounce = true
		}
	} else {
		info.IntMsgInfo.Bounce = *msg.Bounceable
	}

	intMsg := tlb.Message{
		Info: info,
		Body: tlb.EitherRef[tlb.Any]{
			IsRight: true,
			Value:   tlb.Any(*body),
		},
	}

	if msg.Code == nil && msg.Data == nil {
		intMsg.Init.Null = true
	} else {
		intMsg.Init.Value.IsRight = true
		var init tlb.StateInit
		init.Special.Null = true
		init.SplitDepth.Null = true
		if msg.Code != nil {
			init.Code.Value.Value = *msg.Code
		} else {
			init.Code.Null = true
		}
		if msg.Data != nil {
			init.Data.Value.Value = *msg.Data
		} else {
			init.Data.Null = true
		}
		intMsg.Init.Value.Value = init
	}
	return intMsg, nil
}

func (w *Wallet) getInit() (tlb.StateInit, error) {
	publicKey := w.key.Public().(ed25519.PublicKey)
	id := int(w.subWalletId)
	return generateStateInit(publicKey, w.ver, int(w.address.Workchain), &id)
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

func (w *Wallet) getSeqno(ctx context.Context) (uint32, error) {
	if w.blockchain == nil {
		return 0, tongo.BlockchainInterfaceIsNil
	}
	return w.blockchain.GetSeqno(ctx, w.address)
}

// SimpleSend
// Generates a signed external message for wallet with custom internal messages and default TTL
// Gets actual seqno and attach init for wallet if it needed
// The payload is serialized into bytes and sent by the method SendRawMessage
func (w *Wallet) SimpleSend(ctx context.Context, messages []Message) error {
	if w.blockchain == nil {
		return tongo.BlockchainInterfaceIsNil
	}

	var init *tlb.StateInit

	state, err := w.blockchain.GetAccountState(ctx, w.GetAddress())
	if err != nil {
		return err
	}
	var seqno uint32
	if state.Account.Status() == tlb.AccountActive {
		seqno, err = w.getSeqno(ctx)
		if err != nil {
			return err
		}
	}

	if seqno == 0 {
		i, err := w.getInit()
		if err != nil {
			return err
		}
		init = &i
	}
	var (
		msgArray []RawMessage
		mode     byte
	)
	for _, m := range messages {
		intMsg, err := generateInternalMessage(ctx, m, w.blockchain)
		if err != nil {
			return err
		}
		if m.Mode == nil {
			mode = DefaultMessageMode
		} else {
			mode = *m.Mode
		}
		cell := boc.NewCell()
		err = tlb.Marshal(cell, intMsg)
		if err != nil {
			return err
		}
		msgArray = append(msgArray, RawMessage{Message: cell, Mode: mode})
	}
	validUntil := time.Now().Add(DefaultMessageLifetime)
	return w.RawSend(ctx, seqno, validUntil, msgArray, init)
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

// SendJetton
// Sends Jettons to recipient address
func (w *Wallet) SendJetton(ctx context.Context, messages []jetton.TransferMessage) error {
	var msgArray []Message
	for _, m := range messages {
		body, err := buildJettonTransferBody(w.GetAddress(), m)
		if err != nil {
			return err
		}
		jettonWallet, err := m.Jetton.GetJettonWallet(ctx, w.GetAddress())
		if err != nil {
			return err
		}
		msgArray = append(msgArray, Message{
			Amount:  m.TonAmount,
			Address: jettonWallet,
			Body:    body,
		})
	}
	return w.SimpleSend(ctx, msgArray)
}

func buildJettonTransferBody(owner tongo.AccountID, msg jetton.TransferMessage) (*boc.Cell, error) {
	payload := boc.NewCell()
	if msg.Comment != nil && msg.Payload != nil {
		return nil, fmt.Errorf("only payload or comment must be presented")
	} else if msg.Comment != nil {
		err := tlb.Marshal(payload, TextComment(*msg.Comment))
		if err != nil {
			return nil, err
		}
	} else if msg.Payload != nil {
		payload = msg.Payload
	}
	var responseDestination tlb.MsgAddress
	if msg.ResponseDestination == nil {
		responseDestination = tongo.MsgAddressFromAccountID(&owner) // send excess to sender wallet
	} else {
		responseDestination = tongo.MsgAddressFromAccountID(msg.ResponseDestination)
	}
	transferMsg := struct {
		Magic               tlb.Magic `tlb:"transfer#0f8a7ea5"`
		QueryId             uint64
		Amount              tlb.VarUInteger16
		Destination         tlb.MsgAddress
		ResponseDestination tlb.MsgAddress
		CustomPayload       tlb.Maybe[tlb.Ref[tlb.Any]]
		ForwardTonAmount    tlb.Grams // (VarUInteger 16)
		ForwardPayload      tlb.EitherRef[tlb.Any]
	}{
		QueryId:             rand.Uint64(),
		Amount:              tlb.VarUInteger16(*msg.JettonAmount),
		Destination:         tongo.MsgAddressFromAccountID(&msg.Destination),
		ResponseDestination: responseDestination,
		ForwardTonAmount:    tlb.Grams(msg.ForwardTonAmount),
	}
	transferMsg.CustomPayload.Null = true
	transferMsg.ForwardPayload.IsRight = true
	transferMsg.ForwardPayload.Value = tlb.Any(*payload)
	res := boc.NewCell()
	err := tlb.Marshal(res, transferMsg)
	if err != nil {
		return nil, err
	}
	return res, nil
}
