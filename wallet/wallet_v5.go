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

type BlumWalletData struct {
	Version            uint16
	Seqno              uint32
	WalletID           uint32
	PublicKey          tlb.Bits256
	IsFrozen           bool
	MinSigAmount       uint8
	AdminKeys          tlb.HashmapE[tlb.Bits256, bool]
	SingleSuspendSeqno uint32
	MassSuspendSeqno   uint16
}

func (w *walletV5R1) generateStateInit() (*tlb.StateInit, error) {
	//publicKeys := make([]ed25519.PublicKey, 0, 3)
	//for _, hexStr := range []string{
	//	"9fa007ac1e01b426c413de7516e326146af91a38d6910b58e5c5b52c9a3fca40",
	//	"7083e317cf2ea0ad8cc0e9eda6bac78e41408c562d15b37473fa4c3d5b0b1184",
	//	"7acc556f12b26e086806b473dee1cc3c0595af7e9433be461025d9a973bfb0e8",
	//} {
	//	pk, err := tonconnect.DecodePublicKey(hexStr)
	//	if err != nil {
	//		panic(err)
	//	}
	//	publicKeys = append(publicKeys, pk)
	//}

	// Define the byte slices
	byteArrays := [][]byte{
		{159, 160, 7, 172, 30, 1, 180, 38, 196, 19, 222, 117, 22, 227, 38, 20, 106, 249, 26, 56, 214, 145, 11, 88, 229, 197, 181, 44, 154, 63, 202, 64},
		{112, 131, 227, 23, 207, 46, 160, 173, 140, 192, 233, 237, 166, 186, 199, 142, 65, 64, 140, 86, 45, 21, 179, 116, 115, 250, 76, 61, 91, 11, 17, 132},
		{122, 204, 85, 111, 18, 178, 110, 8, 104, 6, 180, 115, 222, 225, 204, 60, 5, 149, 175, 126, 148, 51, 190, 70, 16, 37, 217, 169, 115, 191, 176, 232},
	}
	var adminKeys tlb.HashmapE[tlb.Bits256, bool]
	for _, key := range byteArrays {
		adminKeys.Put(publicKeyToBits(key), true)
	}
	data := BlumWalletData{
		Version:            0,
		Seqno:              0,
		WalletID:           0,
		PublicKey:          publicKeyToBits(w.publicKey),
		IsFrozen:           false,
		MinSigAmount:       2,
		AdminKeys:          adminKeys,
		SingleSuspendSeqno: 0,
		MassSuspendSeqno:   0,
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

type blumSignedMessage struct {
	WalletId   uint32
	ValidUntil uint32
	Seqno      uint32
	Actions    *BlumActions `tlb:"maybe^"`
}

func (w *walletV5R1) CreateSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, extensionsActions *W5ExtendedActions, msgConfig MessageConfig) (*boc.Cell, error) {
	actions := make([]BlumSendMessageAction, 0, len(internalMessages))
	for _, msg := range internalMessages {
		actions = append(actions, BlumSendMessageAction{
			Msg:  msg.Message,
			Mode: msg.Mode,
		})
	}
	blumActions := BlumActions(actions)

	msg := blumSignedMessage{
		WalletId:   w.walletID,
		ValidUntil: uint32(msgConfig.ValidUntil.Unix()),
		Seqno:      msgConfig.Seqno,
		Actions:    &blumActions,
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
		var data BlumWalletData
		cell := boc.Cell(state.Account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value)
		if err := tlb.Unmarshal(&cell, &data); err != nil {
			return NextMsgParams{}, err
		}
		return NextMsgParams{
			Seqno: data.Seqno,
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
