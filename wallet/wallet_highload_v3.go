package wallet

/*

HighLoad wallet V3
Contract repo: https://github.com/ton-blockchain/highload-wallet-contract-v3

TLB scheme:

storage$_ public_key:bits256 subwallet_id:uint32 old_queries:(HashmapE 14 ^Cell)
          queries:(HashmapE 14 ^Cell) last_clean_time:uint64 timeout:uint22
          = Storage;

_ shift:uint13 bit_number:(## 10) { bit_number >= 0 } { bit_number < 1023 } = QueryId;

// crc32('internal_transfer n:# query_id:uint64 actions:^OutList n = InternalMsgBody n') = ae42e5a4

internal_transfer#ae42e5a4 {n:#} query_id:uint64 actions:^(OutList n) = InternalMsgBody n;

_ {n:#}  subwallet_id:uint32 message_to_send:^Cell send_mode:uint8 query_id:QueryId created_at:uint64 timeout:uint22 = MsgInner;

msg_body$_ {n:#} signature:bits512 ^(MsgInner) = ExternalInMsgBody;

*/

import (
	"crypto/ed25519"
	"fmt"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"time"
)

// DefaultSubWalletHighloadV3 https://github.com/ton-blockchain/highload-wallet-contract-v3?tab=readme-ov-file#highload-wallet-contract-v3
const DefaultSubWalletHighloadV3 = 0x10ad

const DefaultTimeoutHighloadV3 = 60 * 60 // TODO: recommended 1 hour to 24 hours

var _ wallet = &walletHighloadV3{}

type walletHighloadV3 struct {
	version     Version
	publicKey   ed25519.PublicKey
	workchain   int
	subWalletID uint32
	timeout     tlb.Uint22
}

// DataHighloadV3 represents storage data of a wallet contract.
type DataHighloadV3 struct {
	PublicKey     tlb.Bits256
	SubWalletId   uint32
	OldQueries    tlb.HashmapE[tlb.Uint14, tlb.Any]
	Queries       tlb.HashmapE[tlb.Uint14, tlb.Any]
	LastCleanTime uint64
	Timeout       tlb.Uint22
}

func newWalletHighloadV3(ver Version, key ed25519.PublicKey, options Options) *walletHighloadV3 {
	workchain := defaultOr(options.Workchain, 0)
	subWalletID := defaultOr(options.SubWalletID, uint32(DefaultSubWalletHighloadV3))
	// TODO: add custom message lifetime with size check
	return &walletHighloadV3{
		version:     ver,
		publicKey:   key,
		workchain:   workchain,
		subWalletID: subWalletID,
		timeout:     DefaultTimeoutHighloadV3,
	}
}

func (w *walletHighloadV3) generateAddress() (ton.AccountID, error) {
	stateInit, err := w.generateStateInit()
	if err != nil {
		return ton.AccountID{}, err
	}
	return generateAddress(w.workchain, *stateInit)
}

func (w *walletHighloadV3) generateStateInit() (*tlb.StateInit, error) {
	data := DataHighloadV3{
		SubWalletId: w.subWalletID,
		PublicKey:   publicKeyToBits(w.publicKey),
		Timeout:     w.timeout,
	}
	return generateStateInit(w.version, data)
}

func (w *walletHighloadV3) maxMessageNumber() int {
	return 254 * 254
}

func (w *walletHighloadV3) createSignedMsgBodyCell(privateKey ed25519.PrivateKey, internalMessages []RawMessage, msgConfig MessageConfig) (*boc.Cell, error) {
	// TODO: or use special queryID generator function as option
	now := time.Now().UnixMilli()
	queryID := tlb.Uint23((now / 100) % (1 << 23)) // allow to send every 100ms. overflows after 233 hr
	addr, err := w.generateAddress()
	if err != nil {
		return nil, err
	}
	msgInner := HighloadV3Message{
		SubwalletID: w.subWalletID,
		Messages:    internalMessages,
		SendMode:    DefaultMessageMode,
		QueryID:     queryID,
		CreatedAt:   uint64(now/1000 - 30), // TODO: fix -30 after the liteservers are updated
		Timeout:     w.timeout,
		wallet:      addr,
	}
	innerCell := boc.NewCell()
	if err := tlb.Marshal(innerCell, msgInner); err != nil {
		return nil, err
	}
	signBytes, err := innerCell.Sign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("can not sign wallet message body: %v", err)
	}
	// msg_body$_ {n:#} signature:bits512 ^(MsgInner) = ExternalInMsgBody;
	signedBodyCell := boc.NewCell()
	err = signedBodyCell.WriteBytes(signBytes)
	if err != nil {
		return nil, err
	}
	_ = signedBodyCell.AddRef(innerCell)
	return signedBodyCell, nil
}

func (w *walletHighloadV3) NextMessageParams(state tlb.ShardAccount) (NextMsgParams, error) {
	initRequired := state.Account.Status() == tlb.AccountUninit || state.Account.Status() == tlb.AccountNone
	if !initRequired {
		return NextMsgParams{}, nil
	}
	stateInit, err := w.generateStateInit()
	if err != nil {
		return NextMsgParams{}, err
	}
	return NextMsgParams{Init: stateInit}, nil
}
