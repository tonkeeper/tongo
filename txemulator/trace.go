package txemulator

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"time"
)

type Tracer struct {
	e                   *Emulator
	currentShardAccount map[tongo.AccountID]tlb.ShardAccount
	blockchain          accountGetter
	counter             int
	limit               int
}

type TxTree struct {
	TX       tlb.Transaction
	Children []*TxTree
}

type TraceOptions struct {
	config         string
	limit          int
	blockchain     accountGetter
	time           int64
	checkSignature bool
}
type accountGetter interface {
	GetAccountState(ctx context.Context, a tongo.AccountID) (tlb.ShardAccount, error)
}

type AccountGetterMixin struct {
	g     accountGetter
	mixin map[tongo.AccountID]tlb.ShardAccount
}

func NewAccountGetterMixin(g accountGetter, mixin map[tongo.AccountID]tlb.ShardAccount) AccountGetterMixin {
	return AccountGetterMixin{
		g:     g,
		mixin: mixin,
	}
}

func (g AccountGetterMixin) GetAccountState(ctx context.Context, a tongo.AccountID) (tlb.ShardAccount, error) {
	s, prs := g.mixin[a]
	if prs {
		return s, nil
	}
	return g.g.GetAccountState(ctx, a)
}

func WithConfig(c *boc.Cell) TraceOption {
	return func(o *TraceOptions) {
		o.config, _ = c.ToBocBase64()
	}
}

func WithConfigBase64(c string) TraceOption {
	return func(o *TraceOptions) {
		o.config = c
	}
}

func WithLimit(l int) TraceOption {
	return func(o *TraceOptions) {
		o.limit = l
	}
}

func WithTime(t int64) TraceOption {
	return func(o *TraceOptions) {
		o.time = t
	}
}

func WithAccountsSource(b accountGetter) TraceOption {
	return func(o *TraceOptions) {
		o.blockchain = b
	}
}

func WithSignatureCheck() TraceOption {
	return func(o *TraceOptions) {
		o.checkSignature = true
	}
}

type TraceOption func(o *TraceOptions)

func NewTraceBuilder(options ...TraceOption) (*Tracer, error) {
	var option TraceOptions
	for _, o := range options {
		o(&option)
	}
	if option.config == "" {
		option.config = DefaultConfig
	}
	if option.limit == 0 {
		option.limit = 100
	}
	if option.blockchain == nil {
		var err error
		option.blockchain, err = liteapi.NewClientWithDefaultMainnet()
		if err != nil {
			return nil, err
		}
	}
	e, err := newEmulatorBase64(option.config, LogTruncated)
	if err != nil {
		return nil, err
	}
	if option.time == 0 {
		option.time = time.Now().Unix()
	}
	err = e.SetUnixtime(uint32(option.time))
	if err != nil {
		return nil, err
	}
	err = e.SetIgnoreSignatureCheck(!option.checkSignature)
	if err != nil {
		return nil, err
	}
	return &Tracer{
		e:                   e,
		currentShardAccount: make(map[tongo.AccountID]tlb.ShardAccount),
		blockchain:          option.blockchain,
		limit:               option.limit,
	}, nil
}

func (t *Tracer) Run(ctx context.Context, message tlb.Message) (*TxTree, error) {
	if t.counter >= t.limit {
		return nil, fmt.Errorf("to many iterations: %v/%v", t.counter, t.limit)
	}
	var a tlb.MsgAddress
	switch message.Info.SumType {
	case "IntMsgInfo":
		a = message.Info.IntMsgInfo.Dest
	case "ExtInMsgInfo":
		a = message.Info.ExtInMsgInfo.Dest
	default:
		return nil, fmt.Errorf("can't emulate message with type %v", message.Info.SumType)
	}
	accountAddr, err := tongo.AccountIDFromTlb(a)
	if err != nil {
		return nil, err
	}
	if accountAddr == nil {
		return nil, fmt.Errorf("destination account is null")
	}
	state, prs := t.currentShardAccount[*accountAddr]
	if !prs {
		state, err = t.blockchain.GetAccountState(ctx, *accountAddr)
		if err != nil {
			return nil, err
		}
	}
	result, err := t.e.Emulate(state, message)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("iteration: %v, exitCode: %v, Text: %v, ", t.counter, result.Error.ExitCode, result.Error.Text)
	}
	if result.Emulation == nil {
		return nil, fmt.Errorf("empty emulation result on iteration %v", t.counter)
	}
	t.counter++
	t.currentShardAccount[*accountAddr] = result.Emulation.ShardAccount
	tree := &TxTree{
		TX: result.Emulation.Transaction,
	}
	for _, m := range result.Emulation.Transaction.Msgs.OutMsgs.Values() {
		if m.Value.Info.SumType == "ExtOutMsgInfo" {
			continue
		}
		child, err := t.Run(ctx, m.Value)
		if err != nil {
			return tree, err
		}
		tree.Children = append(tree.Children, child)
	}
	return tree, err
}

func (t *Tracer) FinalStates() map[tongo.AccountID]tlb.ShardAccount {
	return t.currentShardAccount
}
