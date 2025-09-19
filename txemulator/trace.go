package txemulator

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/tonkeeper/tongo/liteclient"

	"github.com/tonkeeper/tongo/boc"
	codePkg "github.com/tonkeeper/tongo/code"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const basicShardDelay = 8 //difference between shards
const randomDelay = 5     // should be less than basicShardDelay

type Tracer struct {
	e                   *Emulator
	currentShardAccount map[ton.AccountID]tlb.ShardAccount
	blockchain          accountGetter
	counter             int
	limit               int
	softLimit           int
	time                uint32
	unprocessed         int
	ignoreSignDepth     int
	shards              []shard
}

type TxTree struct {
	TX       tlb.Transaction
	Children []*TxTree
}

type TraceOptions struct {
	config               string
	limit                int
	softLimit            int
	ignoreSignatureDepth int
	blockchain           accountGetter
	time                 int64
	predefinedAccounts   map[ton.AccountID]tlb.ShardAccount
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

func WithIgnoreSignatureDepth(d int) TraceOption {
	return func(o *TraceOptions) error {
		o.ignoreSignatureDepth = d
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

	block, err := option.blockchain.GetMasterchainInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %v", err)
	}

	shards, err := option.blockchain.GetAllShardsInfo(context.Background(), block.Last.ToBlockIdExt())
	if err != nil {
		return nil, fmt.Errorf("failed to get shards info: %v", err)
	}
	shardConfig := []shard{
		{workchain: -1, ShardID: ton.ShardID{}},
	}
	for _, s := range shards {
		id, err := ton.ParseShardID(int64(s.Shard))
		if err != nil {
			return nil, fmt.Errorf("failed to parse shard ID: %v", err)
		}
		shardConfig = append(shardConfig, shard{ShardID: id, workchain: s.Workchain})
	}

	// TODO: set gas limit, currently, the transaction emulator doesn't support that
	return &Tracer{
		e:                   e,
		currentShardAccount: option.predefinedAccounts,
		blockchain:          option.blockchain,
		limit:               option.limit,
		softLimit:           option.softLimit,
		ignoreSignDepth:     option.ignoreSignatureDepth,
		shards:              shardConfig,
		time:                uint32(option.time),
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
	fakeRoot := TxTree{}
	m, err := toEmulatedMessage(message, &fakeRoot)
	if err != nil {
		return nil, err
	}
	i := t.routeMessage(m)
	if i == -1 {
		return nil, fmt.Errorf("failed to route message")
	}
	t.shards[i].input = append(t.shards[i].input, m)
	t.unprocessed++

	for t.unprocessed > 0 && (t.softLimit == 0 || t.softLimit < t.counter) {
		err = t.emulationLoop(ctx)
		if err != nil {
			return nil, err
		}
	}
	if len(fakeRoot.Children) == 0 {
		return nil, fmt.Errorf("no transactions were processed")
	}
	return fakeRoot.Children[0], nil
}

func (t *Tracer) emulationLoop(ctx context.Context) error {
	for shardIndex := range t.shards {
		err := t.e.SetUnixtime(t.time + rand.Uint32N(randomDelay))
		if err != nil {
			return err
		}
		for i := 0; i < len(t.shards[shardIndex].input); i++ {
			if t.counter >= t.limit {
				return fmt.Errorf("to many iterations: %v/%v", t.counter, t.limit)
			}
			if t.softLimit > 0 && t.counter == t.softLimit {
				return nil
			}
			message := t.shards[shardIndex].input[i]
			trace, err := t.emulateMessage(ctx, message, t.ignoreSignDepth > t.counter)
			if err != nil {
				return err
			}
			message.parentTrace.Children = append(message.parentTrace.Children, trace)
			t.counter++
			t.unprocessed--
			for _, m := range trace.TX.Msgs.OutMsgs.Values() {
				if m.Value.Info.SumType == "ExtOutMsgInfo" {
					continue
				}
				msg, err := toEmulatedMessage(m.Value, trace)
				if err != nil {
					return err
				}
				t.unprocessed++
				if t.routeMessage(msg) == shardIndex {
					t.shards[shardIndex].input = append(t.shards[shardIndex].input, msg)
				} else {
					t.shards[shardIndex].output = append(t.shards[shardIndex].output, msg)
				}
			}
		}
		t.shards[shardIndex].input = t.shards[shardIndex].input[:0]
	}
	for shardIndex := range t.shards {
		for _, m := range t.shards[shardIndex].output {
			dest := t.routeMessage(m)
			t.shards[dest].input = append(t.shards[dest].input, m)
		}
		t.shards[shardIndex].output = t.shards[shardIndex].output[:0]
	}
	t.time += basicShardDelay
	return nil
}

func (t *Tracer) emulateMessage(ctx context.Context, m emulatedMessage, ignoreSignature bool) (*TxTree, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var err error
	state, prs := t.currentShardAccount[m.dest]
	if !prs {
		state, err = t.blockchain.GetAccountState(ctx, m.dest)
		if err != nil {
			return nil, err
		}
	}

	publicLibs := map[ton.Bits256]*boc.Cell{}
	for _, code := range []*boc.Cell{accountCode(state), msgStateInitCode(m.msg)} {
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

	err = t.e.SetIgnoreSignatureCheck(ignoreSignature)
	if err != nil {
		return nil, err
	}

	result, err := t.e.Emulate(state, m.msg)
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
	t.currentShardAccount[m.dest] = result.Emulation.ShardAccount

	return &TxTree{
		TX: result.Emulation.Transaction,
	}, nil
}

type shard struct {
	ton.ShardID
	workchain int32
	input     []emulatedMessage
	output    []emulatedMessage
}

type emulatedMessage struct {
	msg         tlb.Message
	dest        ton.AccountID
	parentTrace *TxTree
}

func toEmulatedMessage(m tlb.Message, parentTx *TxTree) (emulatedMessage, error) {
	var a tlb.MsgAddress
	switch m.Info.SumType {
	case "IntMsgInfo":
		a = m.Info.IntMsgInfo.Dest
	case "ExtInMsgInfo":
		a = m.Info.ExtInMsgInfo.Dest
	default:
		return emulatedMessage{}, fmt.Errorf("can't emulate message with type %v", m.Info.SumType)
	}
	account, err := ton.AccountIDFromTlb(a)
	if err != nil {
		return emulatedMessage{}, err
	}
	if account == nil {
		return emulatedMessage{}, fmt.Errorf("destination account is null")
	}
	return emulatedMessage{
		msg:         m,
		dest:        *account,
		parentTrace: parentTx,
	}, nil
}

func (t *Tracer) routeMessage(m emulatedMessage) int {
	i := -1
	for i := range t.shards {
		if t.shards[i].workchain == m.dest.Workchain && t.shards[i].MatchAccountID(m.dest) {
			return i
		}
	}
	return i
}

func (t *Tracer) FinalStates() map[ton.AccountID]tlb.ShardAccount {
	return t.currentShardAccount
}
