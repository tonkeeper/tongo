package jetton

import (
	"context"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"math/big"
	"strconv"
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
	JettonAmount        *big.Int
	Destination         tongo.AccountID
	ResponseDestination *tongo.AccountID
	TonAmount           int64
	ForwardTonAmount    int64
	Comment             *string
	Payload             *boc.Cell
}

func NewJetton(master tongo.AccountID, blockchain blockchain) *Jetton {
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
