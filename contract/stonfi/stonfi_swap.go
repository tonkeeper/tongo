package stonfi

import (
	"context"
	"math/big"

	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/contract/jetton"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Stonfi creates a swap message.
type Stonfi struct {
	cli *liteapi.Client

	router          ton.AccountID
	master0, token0 ton.AccountID
	master1, token1 ton.AccountID
}

var TestnetRouter = ton.MustParseAccountID("kQBsGx9ArADUrREB34W-ghgsCgBShvfUr4Jvlu-0KGc33a1n")
var MainnetRouter = ton.MustParseAccountID("EQB3ncyBUTjZUA5EnFKR5_EnOMI9V1tTEAAPaiU71gc4TiUt")
var PTON = ton.MustParseAccountID("EQCM3B12QK1e4yZSf8GtBRT0aLMNyEsBc_DhVfRRtOEffLez")

func NewStonfi(ctx context.Context, cli *liteapi.Client, router, master0, master1 ton.AccountID) (*Stonfi, error) {
	j0 := jetton.New(master0, cli)
	token0, err := j0.GetJettonWallet(ctx, router)
	if err != nil {
		return nil, err
	}
	j1 := jetton.New(master1, cli)
	token1, err := j1.GetJettonWallet(ctx, router)
	if err != nil {
		return nil, err
	}
	return &Stonfi{
		cli:     cli,
		router:  router,
		master0: master0,
		token0:  token0,
		master1: master1,
		token1:  token1,
	}, nil
}

func (s *Stonfi) EstimateMinOut(ctx context.Context, amount big.Int) (*big.Int, error) {
	_, value, err := abi.GetPoolAddress(ctx, s.cli, s.router, s.token0.ToMsgAddress(), s.token1.ToMsgAddress())
	if err != nil {
		return nil, err
	}
	result, ok := value.(abi.GetPoolAddress_StonfiResult)
	if !ok {
		return nil, err
	}
	pool, err := ton.AccountIDFromTlb(result.PoolAddress)
	if err != nil {
		return nil, err
	}
	_, output, err := abi.GetExpectedOutputs(context.Background(), s.cli, *pool, tlb.Int257(amount), s.token0.ToMsgAddress())
	if err != nil {
		return nil, err
	}
	result2, ok := output.(abi.GetExpectedOutputs_StonfiResult)
	if !ok {
		return nil, err
	}
	outputValue := big.Int(result2.Out)
	return &outputValue, nil
}

func (s *Stonfi) MakeSwapMessage(attachedTON tlb.Grams, forwardTONAmount tlb.Grams, jettonAmount big.Int, minOut big.Int, address ton.AccountID) (*jetton.TransferMessage, error) {
	payload := abi.StonfiSwapJettonPayload{
		TokenWallet: s.token1.ToMsgAddress(),
		MinOut:      tlb.VarUInteger16(minOut),
		ToAddress:   address.ToMsgAddress(),
	}
	c := boc.NewCell()
	if err := c.WriteUint(0x25938561, 32); err != nil {
		return nil, err
	}
	if err := tlb.Marshal(c, payload); err != nil {
		return nil, err
	}
	jettonTransfer := jetton.TransferMessage{
		Sender:           address,
		Jetton:           jetton.New(s.master0, s.cli),
		JettonAmount:     &jettonAmount,
		Destination:      s.router,
		AttachedTon:      attachedTON,
		ForwardTonAmount: forwardTONAmount,
		ForwardPayload:   c,
	}
	return &jettonTransfer, nil
}
