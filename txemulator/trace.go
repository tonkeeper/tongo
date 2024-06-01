package txemulator

import (
	"context"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo/boc"
	codePkg "github.com/tonkeeper/tongo/code"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Tracer struct {
	e                   *Emulator
	currentShardAccount map[ton.AccountID]tlb.ShardAccount
	blockchain          accountGetter
	counter             int
	limit               int
	softLimit           int
}

type TxTree struct {
	TX       tlb.Transaction
	Children []*TxTree
}

type TraceOptions struct {
	config             string
	limit              int
	softLimit          int
	blockchain         accountGetter
	time               int64
	checkSignature     bool
	predefinedAccounts map[ton.AccountID]tlb.ShardAccount
}

type accountGetter interface {
	GetAccountState(ctx context.Context, a ton.AccountID) (tlb.ShardAccount, error)
	GetLibraries(ctx context.Context, libraries []ton.Bits256) (map[ton.Bits256]*boc.Cell, error)
}

func WithConfig(c *boc.Cell) TraceOption {
	return func(o *TraceOptions) error {
		var err error
		o.config, err = c.ToBocBase64()
		return err
	}
}

func WithConfigBase64(c string) TraceOption {
	return func(o *TraceOptions) error {
		o.config = c
		return nil
	}
}

func WithLimit(l int) TraceOption {
	return func(o *TraceOptions) error {
		o.limit = l
		return nil
	}
}

func WithSoftLimit(l int) TraceOption {
	return func(o *TraceOptions) error {
		o.softLimit = l
		return nil
	}
}

func WithTime(t int64) TraceOption {
	return func(o *TraceOptions) error {
		o.time = t
		return nil
	}
}

func WithAccountsMap(m map[ton.AccountID]tlb.ShardAccount) TraceOption {
	return func(o *TraceOptions) error {
		o.predefinedAccounts = m
		return nil
	}
}
func WithAccounts(accounts ...tlb.ShardAccount) TraceOption {
	return func(o *TraceOptions) error {
		for i := range accounts {
			a, err := ton.AccountIDFromTlb(accounts[i].Account.Account.Addr)
			if err != nil {
				return err
			}
			o.predefinedAccounts[*a] = accounts[i]
		}
		return nil
	}
}

func WithTestnet() TraceOption {
	return func(o *TraceOptions) error {
		var err error
		o.blockchain, err = liteapi.NewClientWithDefaultTestnet()
		return err
	}
}

func WithAccountsSource(b accountGetter) TraceOption {
	return func(o *TraceOptions) error {
		o.blockchain = b
		return nil
	}
}

func WithSignatureCheck() TraceOption {
	return func(o *TraceOptions) error {
		o.checkSignature = true
		return nil
	}
}

type TraceOption func(o *TraceOptions) error

func NewTraceBuilder(options ...TraceOption) (*Tracer, error) {
	option := TraceOptions{
		config:             DefaultConfig,
		limit:              100,
		blockchain:         nil,
		time:               time.Now().Unix(),
		checkSignature:     false,
		predefinedAccounts: make(map[ton.AccountID]tlb.ShardAccount),
	}
	for _, o := range options {
		err := o(&option)
		if err != nil {
			return nil, err
		}
	}
	if option.blockchain == nil {
		return nil, fmt.Errorf("blockchain source is not configured. please use WithAccountsSource")
	}
	e, err := newEmulatorBase64(option.config, LogTruncated)
	if err != nil {
		return nil, err
	}
	err = e.SetUnixtime(uint32(option.time))
	if err != nil {
		return nil, err
	}
	err = e.SetIgnoreSignatureCheck(!option.checkSignature)
	if err != nil {
		return nil, err
	}
	// TODO: set gas limit, currently, the transaction emulator doesn't support that
	return &Tracer{
		e:                   e,
		currentShardAccount: option.predefinedAccounts,
		blockchain:          option.blockchain,
		limit:               option.limit,
		softLimit:           option.softLimit,
	}, nil
}

func accountCode(account tlb.ShardAccount) *boc.Cell {
	if account.Account.SumType == "AccountNone" {
		return nil
	}
	if account.Account.Account.Storage.State.SumType != "AccountActive" {
		return nil
	}
	code := account.Account.Account.Storage.State.AccountActive.StateInit.Code
	if !code.Exists {
		return nil
	}
	cell := code.Value.Value
	return &cell
}

func msgStateInitCode(msg tlb.Message) *boc.Cell {
	if !msg.Init.Exists {
		return nil
	}
	code := msg.Init.Value.Value.Code
	if !code.Exists {
		return nil
	}
	cell := code.Value.Value
	return &cell
}

func (t *Tracer) Run(ctx context.Context, message tlb.Message) (*TxTree, error) {
	if t.counter >= t.limit {
		return nil, fmt.Errorf("to many iterations: %v/%v", t.counter, t.limit)
	}
	if t.softLimit > 0 && t.counter >= t.softLimit {
		return nil, nil
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
	accountAddr, err := ton.AccountIDFromTlb(a)
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
	publicLibs := map[ton.Bits256]*boc.Cell{}
	for _, code := range []*boc.Cell{accountCode(state), msgStateInitCode(message)} {
		if code == nil {
			continue
		}
		hashes, err := codePkg.FindLibraries(code)
		if err != nil {
			return nil, err
		}
		if len(hashes) > 0 {
			libs, err := t.blockchain.GetLibraries(ctx, hashes)
			if err != nil {
				return nil, err
			}
			for hash, cell := range libs {
				publicLibs[hash] = cell
			}
		}
	}
	if len(publicLibs) > 0 {
		libsBoc, err := codePkg.LibrariesToBase64(publicLibs)
		if err != nil {
			return nil, err
		}
		if err := t.e.setLibs(libsBoc); err != nil {
			return nil, err
		}
	}
	result, err := t.e.Emulate(state, message)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, ErrorWithExitCode{
			Message:   fmt.Sprintf("iteration: %v, exitCode: %v, Text: %v, ", t.counter, result.Error.ExitCode, result.Error.Text),
			ExitCode:  result.Error.ExitCode,
			Iteration: t.counter,
		}
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

func (t *Tracer) FinalStates() map[ton.AccountID]tlb.ShardAccount {
	return t.currentShardAccount
}
