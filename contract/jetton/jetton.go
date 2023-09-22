package jetton

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"time"

	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tep64"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

type blockchain interface {
	GetJettonWallet(ctx context.Context, master, owner ton.AccountID) (ton.AccountID, error)
	GetJettonData(ctx context.Context, master ton.AccountID) (tep64.Metadata, error)
	GetJettonBalance(ctx context.Context, jettonWallet ton.AccountID) (*big.Int, error)
}

type Jetton struct {
	Master     ton.AccountID
	blockchain blockchain
}

type TransferMessage struct {
	Jetton              *Jetton
	Sender              ton.AccountID
	JettonAmount        *big.Int
	Destination         ton.AccountID
	ResponseDestination *ton.AccountID
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
		payload := tlb.Any(*tm.CustomPayload)
		msgBody.CustomPayload = &payload
	}
	if tm.ForwardPayload != nil {
		msgBody.ForwardPayload.IsRight = true
		msgBody.ForwardPayload.Value = abi.JettonPayload{SumType: abi.UnknownJettonOp, Value: tm.ForwardPayload}
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

func New(master ton.AccountID, blockchain blockchain) *Jetton {
	return &Jetton{
		Master:     master,
		blockchain: blockchain,
	}
}

func (j *Jetton) GetBalance(ctx context.Context, owner ton.AccountID) (*big.Int, error) {
	if j.blockchain == nil {
		return nil, errors.New("blockchain interface is nil")
	}
	jettonWallet, err := j.blockchain.GetJettonWallet(ctx, j.Master, owner)
	if err != nil {
		return nil, err
	}
	return j.blockchain.GetJettonBalance(ctx, jettonWallet)
}

func (j *Jetton) GetJettonWallet(ctx context.Context, owner ton.AccountID) (ton.AccountID, error) {
	if j.blockchain == nil {
		return ton.AccountID{}, errors.New("blockchain interface is nil")
	}
	return j.blockchain.GetJettonWallet(ctx, j.Master, owner)
}

func (j *Jetton) GetDecimals(ctx context.Context) (int, error) {
	if j.blockchain == nil {
		return 0, errors.New("blockchain interface is nil")
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
