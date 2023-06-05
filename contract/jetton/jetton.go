package jetton

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

type blockchain interface {
	GetJettonWallet(ctx context.Context, master, owner tongo.AccountID) (tongo.AccountID, error)
	GetJettonData(ctx context.Context, master tongo.AccountID) (tongo.JettonMetadata, error)
	GetJettonBalance(ctx context.Context, jettonWallet tongo.AccountID) (*big.Int, error)
}

type Jetton struct {
	Master     tongo.AccountID
	blockchain blockchain
}

type TransferMessage struct {
	Jetton              *Jetton
	Sender              tongo.AccountID
	JettonAmount        *big.Int
	Destination         tongo.AccountID
	ResponseDestination *tongo.AccountID
	AttachedTon         tlb.Grams
	ForwardTonAmount    tlb.Grams
	ForwardPayload      *boc.Cell
	CustomPayload       *boc.Cell
}

func (tm TransferMessage) ToInternal() (tlb.Message, uint8, error) {
	c := boc.NewCell()
	forwardTon := big.NewInt(int64(tm.ForwardTonAmount))
	msgBody := abi.JettonTransferMsgBody{
		QueryId:             uint64(time.Now().UnixNano()),
		Amount:              tlb.VarUInteger16(*tm.JettonAmount),
		Destination:         tm.Destination.ToMsgAddress(),
		ResponseDestination: tm.ResponseDestination.ToMsgAddress(),
		ForwardTonAmount:    tlb.VarUInteger16(*forwardTon),
	}
	if tm.CustomPayload != nil {
		msgBody.CustomPayload.Exists = true
		msgBody.CustomPayload.Value.Value = tlb.Any(*tm.CustomPayload)
	}
	if tm.ForwardPayload != nil {
		msgBody.ForwardPayload.IsRight = true
		msgBody.ForwardPayload.Value = tlb.Any(*tm.ForwardPayload)
	}
	if err := c.WriteUint(0xf8a7ea5, 32); err != nil {
		return tlb.Message{}, 0, err
	}
	if err := tlb.Marshal(c, msgBody); err != nil {
		return tlb.Message{}, 0, err
	}
	jettonWallet, err := tm.Jetton.GetJettonWallet(context.TODO(), tm.Sender)
	if err != nil {
		return tlb.Message{}, 0, err
	}
	m := wallet.Message{
		Amount:  tm.AttachedTon,
		Address: jettonWallet,
		Bounce:  true,
		Mode:    wallet.DefaultMessageMode,
		Body:    c,
	}
	return m.ToInternal()
}

func New(master tongo.AccountID, blockchain blockchain) *Jetton {
	return &Jetton{
		Master:     master,
		blockchain: blockchain,
	}
}

func (j *Jetton) GetBalance(ctx context.Context, owner tongo.AccountID) (*big.Int, error) {
	if j.blockchain == nil {
		return nil, tongo.BlockchainInterfaceIsNil
	}
	jettonWallet, err := j.blockchain.GetJettonWallet(ctx, j.Master, owner)
	if err != nil {
		return nil, err
	}
	return j.blockchain.GetJettonBalance(ctx, jettonWallet)
}

func (j *Jetton) GetJettonWallet(ctx context.Context, owner tongo.AccountID) (tongo.AccountID, error) {
	if j.blockchain == nil {
		return tongo.AccountID{}, tongo.BlockchainInterfaceIsNil
	}
	return j.blockchain.GetJettonWallet(ctx, j.Master, owner)
}

func (j *Jetton) GetDecimals(ctx context.Context) (int, error) {
	if j.blockchain == nil {
		return 0, tongo.BlockchainInterfaceIsNil
	}
	data, err := j.blockchain.GetJettonData(ctx, j.Master)
	if err != nil {
		return 0, err
	}
	if data.Decimals == "" {
		return 9, nil
	}
	return strconv.Atoi(data.Decimals)
}
