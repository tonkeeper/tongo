package txemulator

import (
	"context"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo/liteclient"

	"math/rand"

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
	currentTime         uint32
	shardConfig         map[ton.ShardID]struct{}
	pending             []pendingTask
}

type pendingTask struct {
	parent *TxTree
	idx    int
	run    func() (*TxTree, error)
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
	predefinedAccounts map[ton.AccountID]tlb.ShardAccount
}

type accountGetter interface {
	GetAccountState(ctx context.Context, a ton.AccountID) (tlb.ShardAccount, error)
	GetLibraries(ctx context.Context, libraries []ton.Bits256) (map[ton.Bits256]*boc.Cell, error)
	GetAllShardsInfo(ctx context.Context, blockID ton.BlockIDExt) ([]ton.BlockIDExt, error)
	GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error)
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

type TraceOption func(o *TraceOptions) error

func NewTraceBuilder(options ...TraceOption) (*Tracer, error) {
	option := TraceOptions{
		config:             DefaultConfig,
		limit:              100,
		blockchain:         nil,
		time:               time.Now().Unix(),
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
	err = e.SetIgnoreSignatureCheck(true)
	if err != nil {
		return nil, err
	}

	block, err := option.blockchain.GetMasterchainInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %v", err)
	}
	shards, err := option.blockchain.GetAllShardsInfo(context.Background(), block.Last.ToBlockIdExt())
	if err != nil {
		return nil, fmt.Errorf("failed to get shards info: %v", err)
	}
	shardConfig := make(map[ton.ShardID]struct{})
	for _, shard := range shards {
		shardID, err := ton.ParseShardID(int64(shard.Shard))
		if err != nil {
			return nil, fmt.Errorf("failed to parse shard ID: %v", err)
		}
		shardConfig[shardID] = struct{}{}
	}

	// TODO: set gas limit, currently, the transaction emulator doesn't support that
	return &Tracer{
		e:                   e,
		currentShardAccount: option.predefinedAccounts,
		blockchain:          option.blockchain,
		limit:               option.limit,
		softLimit:           option.softLimit,
		currentTime:         uint32(option.time),
		shardConfig:         shardConfig,
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

func (t *Tracer) Run(ctx context.Context, message tlb.Message, signatureIgnoreDepth int) (*TxTree, error) {
	root, err := t.run(ctx, message, signatureIgnoreDepth)
	if err != nil {
		return nil, err
	}
	for len(t.pending) > 0 {
		err = t.addRandomDelay()
		if err != nil {
			return nil, err
		}
		tasks := t.pending
		t.pending = nil
		for _, task := range tasks {
			child, err := task.run()
			if err != nil {
				return nil, err
			}
			task.parent.Children[task.idx] = child
		}
	}
	return root, nil
}

func (t *Tracer) run(ctx context.Context, message tlb.Message, signatureIgnoreDepth int) (*TxTree, error) {
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
	sourceShard, err := t.getAccountShard(*accountAddr)
	if err != nil {
		return nil, err
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
	// if remaining signature ignore depth > 0 then ignore signature check, otherwise enable
	err = t.e.SetIgnoreSignatureCheck(signatureIgnoreDepth > 0)
	if err != nil {
		return nil, err
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
	emulationRes, err := t.e.Emulate(state, message)
	if err != nil {
		return nil, err
	}
	if emulationRes.Error != nil {
		return nil, ErrorWithExitCode{
			Message:   fmt.Sprintf("iteration: %v, exitCode: %v, Text: %v, ", t.counter, emulationRes.Error.ExitCode, emulationRes.Error.Text),
			ExitCode:  emulationRes.Error.ExitCode,
			Iteration: t.counter,
		}
	}
	if emulationRes.Emulation == nil {
		return nil, fmt.Errorf("empty emulation result on iteration %v", t.counter)
	}
	t.counter++
	t.currentShardAccount[*accountAddr] = emulationRes.Emulation.ShardAccount

	tree := &TxTree{
		TX: emulationRes.Emulation.Transaction,
	}

	var outs []tlb.Message
	for _, m := range emulationRes.Emulation.Transaction.Msgs.OutMsgs.Values() {
		if m.Value.Info.SumType == "ExtOutMsgInfo" {
			continue
		}
		outs = append(outs, m.Value)
	}
	if len(outs) > 0 {
		tree.Children = make([]*TxTree, len(outs))
	}

	for i, m := range outs {
		msg := m
		childDepth := signatureIgnoreDepth - 1

		destAddr, err := ton.AccountIDFromTlb(msg.Info.IntMsgInfo.Dest)
		if destAddr == nil {
			return nil, fmt.Errorf("destination account is null")
		}
		destShard, err := t.getAccountShard(*destAddr)
		if err != nil {
			return nil, err
		}
		if destShard != sourceShard {
			t.pending = append(t.pending, pendingTask{
				parent: tree,
				idx:    i,
				run: func() (*TxTree, error) {
					return t.run(ctx, msg, childDepth)
				},
			})
			continue
		}
		child, err := t.run(ctx, msg, childDepth)
		if err != nil {
			return nil, err
		}
		tree.Children[i] = child
	}
	return tree, nil
}

func (t *Tracer) addRandomDelay() error {
	delay := rand.Intn(11) + 5 // random number between 5 and 15
	t.currentTime += uint32(delay)
	return t.e.SetUnixtime(t.currentTime)
}

func (t *Tracer) getAccountShard(account ton.AccountID) (ton.ShardID, error) {
	for shardID := range t.shardConfig {
		if shardID.MatchAccountID(account) {
			return shardID, nil
		}
	}
	return ton.ShardID{}, fmt.Errorf("account %v does not belong to any known shard", account)
}

func (t *Tracer) FinalStates() map[ton.AccountID]tlb.ShardAccount {
	return t.currentShardAccount
}
