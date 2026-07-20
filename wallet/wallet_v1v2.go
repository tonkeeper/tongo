package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type DataV1V2 struct {
	Seqno     uint32
	PublicKey tlb.Bits256
}

type walletV1V2 struct {
	version   Version
	publicKey ed25519.PublicKey
	workchain int
}

var _ wallet = &walletV1V2{}

func newWalletV1V2(ver Version, key ed25519.PublicKey, options Options) *walletV1V2 {
	workchain := defaultOr(options.Workchain, 0)
	return &walletV1V2{
		version:   ver,
		publicKey: key,
		workchain: workchain,
	}
}

func (w *walletV1V2) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, err
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletV1V2) generateStateInit() (*tlb.StateInit, error) {
	data := DataV1V2{
		Seqno:     0,
		PublicKey: publicKeyToBits(w.publicKey),
	}
	return generateStateInit(w.version, data)
}

func (w *walletV1V2) MaxMessageNumber() int {
	return 4
}

func (w *walletV1V2) CreateMsgBodyWithoutSignature(internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	return nil, fmt.Errorf("message body creation is not implemented for v1/v2 wallets")
}

func (w *walletV1V2) AttachSignature(body *boc.Cell, signature tlb.Bits512) (*boc.Cell, error) {
	return attachSignatureAsSignedMsgBody(body, signature)
}

func (w *walletV1V2) NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error) {
	panic("implement me")
}
