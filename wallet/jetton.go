package wallet

import (
	"github.com/startfellows/tongo"
	"math/big"
)

type Jetton struct {
	Master     tongo.AccountID
	blockchain interface {
		GetJettonWallet(master, owner tongo.AccountID)
		GetDecimals() int32
	}
}

type JettonSendMessage struct {
	Amount big.Int
	Jetton Jetton
}
