package liteapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi/pool"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/utils"
)

const (
	LiteServerEnvName           = "LITE_SERVERS"
	defaultMaxConnectionsNumber = 1

	// maxTransactionCount specifies a maximum number of transactions that can be requested from a lite server.
	// This is a limitation of lite server:
	// https://github.com/ton-blockchain/ton/blob/v2023.06/validator/impl/liteserver.hpp#L70
	maxTransactionCount = 16
)

var (
	// ErrAccountNotFound is returned by lite server when executing a method for an account that has not been deployed to the blockchain.
	ErrAccountNotFound = errors.New("account not found")
)

// ProofPolicy specifies a policy for proof checks.
// This feature is experimental and can be changed or removed in the future.
type ProofPolicy uint32

const (
	// ProofPolicyUnsafe disables proof checks.
	ProofPolicyUnsafe ProofPolicy = iota
	ProofPolicyFast
)

// Client provides a convenient way to interact with TON blockchain.
//
// By default, it uses a single connection to a lite server.
// But internally, it makes use of a failover pool,
// so it is possible to force it to use multiple connections. Take a look at WithMaxConnectionsNumber() option.
//
// When the client is configured with several connections,
// two different lite servers can be used for two consequent requests.
// Because a blockchain is inherently a distributed system,
// this could lead to some inconsistencies.
// For example,
// 1. you obtain a master head with GetMasterchainInfo,
// 2. you get an account state with GetAccountState,
// the account state can be obtained from a block that is earlier in the blockchain than the master head you obtained at step 1.
// To avoid this, you can use WithBlock() method to specify a target block for all requests.
type Client struct {
	pool *pool.ConnPool
	// proofPolicy specifies a policy for proof checks.
	proofPolicy ProofPolicy

	// archiveDetectionEnabled specifies whether
	// the underlying connections pool maintains information about which nodes are archive nodes.
	archiveDetectionEnabled bool

	// mu protects targetBlockID and networkGlobalID.
	mu              sync.RWMutex
	targetBlockID   *ton.BlockIDExt
	networkGlobalID *int32
}

// Options holds parameters to configure a lite api instance.
type Options struct {
	LiteServers []config.LiteServer
	Timeout     time.Duration
	// MaxConnections specifies a number of connections to lite servers for a connections pool.
	MaxConnections int
	// InitCtx is used when opening a new connection to lite servers during the initialization.
	InitCtx context.Context
	// ProofPolicy specifies a policy for proof checks.
	ProofPolicy ProofPolicy
	// DetectArchiveNodes specifies if a liteapi connection to a node
	// should detect if its node is an archive node.
	DetectArchiveNodes bool

	SyncConnectionsInitialization bool
	PoolStrategy                  pool.Strategy
}

type Option func(o *Options) error

func WithLiteServers(servers []config.LiteServer) Option {
	return func(o *Options) error {
		o.LiteServers = servers
		return nil
	}
}

// WithMaxConnectionsNumber specifies a number of concurrent connections to lite servers
// to be maintained by a connections pool.
// Be careful when combining WithMaxConnectionsNumber and FromEnvs() because
// MaxConnectionsNumber is set by FromEnvs() to the number of servers in LITE_SERVERS env variable.
func WithMaxConnectionsNumber(maxConns int) Option {
	return func(o *Options) error {
		o.MaxConnections = maxConns
		return nil
	}
}

func WithAsyncConnectionsInit() Option {
	return func(o *Options) error {
		o.SyncConnectionsInitialization = false
		return nil
	}
}

func WithPoolStrategy(strategy pool.Strategy) Option {
	return func(o *Options) error {
		o.PoolStrategy = strategy
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = timeout
		return nil
	}
}

func WithProofPolicy(policy ProofPolicy) Option {
	return func(o *Options) error {
		o.ProofPolicy = policy
		return nil
	}
}

func WithDetectArchiveNodes() Option {
	return func(o *Options) error {
		o.DetectArchiveNodes = true
		return nil
	}
}

// WithInitializationContext specifies a context to be used
// when opening a new connection to lite servers during the initialization.
func WithInitializationContext(ctx context.Context) Option {
	return func(o *Options) error {
		o.InitCtx = ctx
		return nil
	}
}

// FromEnvsOrMainnet configures a client to use lite servers from the LITE_SERVERS env variable.
// If LITE_SERVERS is not set, it downloads public config for mainnet from ton.org.
func FromEnvsOrMainnet() Option {
	return func(o *Options) error {
		if value, ok := os.LookupEnv(LiteServerEnvName); ok {
			servers, err := config.ParseLiteServersEnvVar(value)
			if err != nil {
				return err
			}
			o.LiteServers = servers
			o.MaxConnections = len(servers)
			return nil
		}
		file, err := downloadConfig("https://ton.org/global.config.json")
		if err != nil {
			return err
		}
		o.LiteServers = file.LiteServers
		return nil
	}
}

// FromEnvsOrTestnet configures a client to use lite servers from the LITE_SERVERS env variable.
// If LITE_SERVERS is not set, it downloads public config for testnet from ton.org.
func FromEnvsOrTestnet() Option {
	return func(o *Options) error {
		if value, ok := os.LookupEnv(LiteServerEnvName); ok {
			servers, err := config.ParseLiteServersEnvVar(value)
			if err != nil {
				return err
			}
			o.LiteServers = servers
			o.MaxConnections = len(servers)
			return nil
		}
		file, err := downloadConfig("https://ton.org/testnet-global.config.json")
		if err != nil {
			return err
		}
		o.LiteServers = file.LiteServers
		return nil
	}
}

// FromEnvs configures a Client based on the following environment variables:
// LITE_SERVERS - a list of lite servers to use.
// FromEnvs() also sets MaxConnectionsNumber to be equal to the number of servers in LITE_SERVERS.
func FromEnvs() Option {
	return func(o *Options) error {
		if value, ok := os.LookupEnv(LiteServerEnvName); ok {
			servers, err := config.ParseLiteServersEnvVar(value)
			if err != nil {
				return err
			}
			o.LiteServers = servers
			o.MaxConnections = len(servers)
		}
		return nil
	}
}

// Mainnet configures a client to use lite servers from the mainnet.
func Mainnet() Option {
	return func(o *Options) error {
		file, err := downloadConfig("https://ton.org/global.config.json")
		if err != nil {
			return err
		}
		o.LiteServers = file.LiteServers
		return nil
	}
}

// Testnet configures a client to use lite servers from the testnet.
func Testnet() Option {
	return func(o *Options) error {
		file, err := downloadConfig("https://ton.org/testnet-global.config.json")
		if err != nil {
			return err
		}
		o.LiteServers = file.LiteServers
		return nil
	}
}

func WithConfigurationFile(file config.GlobalConfigurationFile) Option {
	return func(o *Options) error {
		o.LiteServers = file.LiteServers
		return nil
	}
}

func NewClientWithDefaultMainnet() (*Client, error) {
	return NewClient(Mainnet())
}

func NewClientWithDefaultTestnet() (*Client, error) {
	return NewClient(Testnet())
}

// NewClient
// Get options and create new lite client. If no options provided - download public config for mainnet from ton.org.
func NewClient(options ...Option) (*Client, error) {
	opts := &Options{
		Timeout:                       60 * time.Second,
		MaxConnections:                defaultMaxConnectionsNumber,
		InitCtx:                       context.Background(),
		ProofPolicy:                   ProofPolicyUnsafe,
		DetectArchiveNodes:            false,
		SyncConnectionsInitialization: true,
		PoolStrategy:                  pool.BestPingStrategy,
	}
	for _, o := range options {
		if err := o(opts); err != nil {
			return nil, err
		}
	}
	if len(opts.LiteServers) == 0 {
		return nil, fmt.Errorf("server list empty")
	}
	connPool := pool.New(opts.PoolStrategy)
	initCh := connPool.InitializeConnections(opts.InitCtx, opts.Timeout, opts.MaxConnections, opts.DetectArchiveNodes, opts.LiteServers)
	if opts.SyncConnectionsInitialization {
		if err := <-initCh; err != nil {
			return nil, err
		}
	}
	client := Client{
		pool:                    connPool,
		proofPolicy:             opts.ProofPolicy,
		archiveDetectionEnabled: opts.DetectArchiveNodes,
	}
	go client.pool.Run(context.TODO())
	return &client, nil
}

func (c *Client) targetBlockOr(blockID ton.BlockIDExt) ton.BlockIDExt {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.targetBlockID == nil {
		return blockID
	}
	return *c.targetBlockID
}

func (c *Client) WithBlock(block ton.BlockIDExt) *Client {
	return &Client{
		pool:          c.pool,
		targetBlockID: &block,
	}
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error) {
	conn := c.pool.BestMasterchainInfoClient()
	if conn == nil {
		return liteclient.LiteServerMasterchainInfoC{}, pool.ErrNoConnections
	}
	return conn.LiteServerGetMasterchainInfo(ctx)
}

func (c *Client) GetMasterchainInfoExt(ctx context.Context, mode uint32) (liteclient.LiteServerMasterchainInfoExtC, error) {
	conn := c.pool.BestMasterchainInfoClient()
	if conn == nil {
		return liteclient.LiteServerMasterchainInfoExtC{}, pool.ErrNoConnections
	}
	return conn.LiteServerGetMasterchainInfoExt(ctx, liteclient.LiteServerGetMasterchainInfoExtRequest{Mode: mode})
}

func (c *Client) GetTime(ctx context.Context) (uint32, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return 0, err
	}
	res, err := client.LiteServerGetTime(ctx)
	return res.Now, err
}

func (c *Client) GetVersion(ctx context.Context) (liteclient.LiteServerVersionC, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerVersionC{}, err
	}
	return client.LiteServerGetVersion(ctx)
}

func (c *Client) GetBlock(ctx context.Context, blockID ton.BlockIDExt) (tlb.Block, error) {
	res, err := c.GetBlockRaw(ctx, blockID)
	if err != nil {
		return tlb.Block{}, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return tlb.Block{}, err
	}
	if len(cells) != 1 {
		return tlb.Block{}, boc.ErrNotSingleRoot
	}
	decoder := tlb.NewDecoder()
	var block tlb.Block
	if err := decoder.Unmarshal(cells[0], &block); err != nil {
		return tlb.Block{}, err
	}
	if c.proofPolicy == ProofPolicyUnsafe {
		return block, nil
	}
	// this should be quite fast because
	// when unmarshalling a block, we calculate hashes for transactions and messages.
	// so most of the cells' hashes should be in the cache.
	hash, err := decoder.Hasher().Hash(cells[0])
	if err != nil {
		return tlb.Block{}, fmt.Errorf("failed to calculate block hash: %w", err)
	}
	if !bytes.Equal(hash[:], blockID.RootHash[:]) {
		return tlb.Block{}, fmt.Errorf("block hash mismatch")
	}
	return block, nil
}

func (c *Client) GetBlockRaw(ctx context.Context, blockID ton.BlockIDExt) (liteclient.LiteServerBlockDataC, error) {
	client, err := c.pool.BestClientByBlockID(ctx, blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockDataC{}, err
	}
	res, err := client.LiteServerGetBlock(ctx, liteclient.LiteServerGetBlockRequest{liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerBlockDataC{}, err
	}
	return res, err
}

func (c *Client) GetState(ctx context.Context, blockID ton.BlockIDExt) ([]byte, ton.Bits256, ton.Bits256, error) {
	res, err := c.GetStateRaw(ctx, blockID)
	if err != nil {
		return nil, ton.Bits256{}, ton.Bits256{}, err
	}
	// TODO: implement state tlb decoding
	return res.Data, ton.Bits256(res.RootHash), ton.Bits256(res.FileHash), nil
}

func (c *Client) GetStateRaw(ctx context.Context, blockID ton.BlockIDExt) (liteclient.LiteServerBlockStateC, error) {
	client, err := c.pool.BestClientByBlockID(ctx, blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockStateC{}, err
	}
	res, err := client.LiteServerGetState(ctx, liteclient.LiteServerGetStateRequest{Id: liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerBlockStateC{}, err
	}
	return res, nil
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (tlb.BlockInfo, error) {
	res, err := c.GetBlockHeaderRaw(ctx, blockID, mode)
	if err != nil {
		return tlb.BlockInfo{}, err
	}
	_, info, err := decodeBlockHeader(res)
	return info, err
}

func (c *Client) GetBlockHeaderRaw(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (liteclient.LiteServerBlockHeaderC, error) {
	client, err := c.pool.BestClientByBlockID(ctx, blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockHeaderC{}, err
	}
	res, err := client.LiteServerGetBlockHeader(ctx, liteclient.LiteServerGetBlockHeaderRequest{
		Id:   liteclient.BlockIDExt(blockID),
		Mode: mode,
	})
	if err != nil {
		return liteclient.LiteServerBlockHeaderC{}, err
	}
	return res, nil
}

func (c *Client) LookupBlock(ctx context.Context, blockID ton.BlockID, mode uint32, lt *uint64, utime *uint32) (ton.BlockIDExt, tlb.BlockInfo, error) {
	client, err := c.pool.BestClientByBlockID(ctx, blockID)
	if err != nil {
		return ton.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	res, err := client.LiteServerLookupBlock(ctx, liteclient.LiteServerLookupBlockRequest{
		Mode: mode,
		Id: liteclient.TonNodeBlockIdC{
			Workchain: blockID.Workchain,
			Shard:     blockID.Shard,
			Seqno:     blockID.Seqno,
		},
		Lt:    lt,
		Utime: utime,
	})
	if err != nil {
		return ton.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func decodeBlockHeader(header liteclient.LiteServerBlockHeaderC) (ton.BlockIDExt, tlb.BlockInfo, error) {
	cells, err := boc.DeserializeBoc(header.HeaderProof)
	if err != nil {
		return ton.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	if len(cells) != 1 {
		return ton.BlockIDExt{}, tlb.BlockInfo{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.BlockHeader]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return ton.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	return header.Id.ToBlockIdExt(), proof.Proof.VirtualRoot.Info, nil // TODO: maybe decode more
}

// SendMessage verifies that the given payload contains an external message and sends it to a lite server.
func (c *Client) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	if err := VerifySendMessagePayload(payload); err != nil {
		return 0, err
	}
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return 0, err
	}
	res, err := client.LiteServerSendMessage(ctx, liteclient.LiteServerSendMessageRequest{Body: payload})
	return res.Status, err
}

func (c *Client) RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, params)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	b, err := cell.ToBoc()
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	client, masterHead, err := c.pool.BestClientByAccountID(ctx, accountID, false)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	req := liteclient.LiteServerRunSmcMethodRequest{
		Mode:     4,
		Id:       liteclient.BlockIDExt(c.targetBlockOr(masterHead)),
		Account:  liteclient.AccountID(accountID),
		MethodId: uint64(methodID),
		Params:   b,
	}
	res, err := client.LiteServerRunSmcMethod(ctx, req)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	var result tlb.VmStack
	if res.ExitCode == 4294967040 { //-256
		return res.ExitCode, nil, ErrAccountNotFound
	}
	cells, err := boc.DeserializeBoc(res.Result)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	if len(cells) != 1 {
		return 0, tlb.VmStack{}, boc.ErrNotSingleRoot
	}
	err = tlb.Unmarshal(cells[0], &result)
	return res.ExitCode, result, err
}

func (c *Client) RunSmcMethod(
	ctx context.Context,
	accountID ton.AccountID,
	method string,
	params tlb.VmStack,
) (uint32, tlb.VmStack, error) {
	return c.RunSmcMethodByID(ctx, accountID, utils.MethodIdFromName(method), params)
}

func (c *Client) GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error) {
	res, err := c.GetAccountStateRaw(ctx, accountID)
	if err != nil {
		return tlb.ShardAccount{}, err
	}
	if len(res.State) == 0 {
		return tlb.ShardAccount{Account: tlb.Account{SumType: "AccountNone"}}, nil
	}
	cells, err := boc.DeserializeBoc(res.State)
	if err != nil {
		return tlb.ShardAccount{}, err
	}
	if len(cells) != 1 {
		return tlb.ShardAccount{}, boc.ErrNotSingleRoot
	}
	var acc tlb.Account
	err = tlb.Unmarshal(cells[0], &acc)
	if err != nil {
		return tlb.ShardAccount{}, err
	}
	lt, hash, err := decodeAccountDataFromProof(res.Proof, accountID)
	return tlb.ShardAccount{Account: acc, LastTransHash: hash, LastTransLt: lt}, err
}

func (c *Client) GetAccountStateRaw(ctx context.Context, accountID ton.AccountID) (liteclient.LiteServerAccountStateC, error) {
	client, masterHead, err := c.pool.BestClientByAccountID(ctx, accountID, false)
	if err != nil {
		return liteclient.LiteServerAccountStateC{}, err
	}
	blockID := c.targetBlockOr(masterHead)
	res, err := client.LiteServerGetAccountState(ctx, liteclient.LiteServerGetAccountStateRequest{
		Account: liteclient.AccountID(accountID),
		Id:      liteclient.BlockIDExt(blockID),
	})
	if err != nil {
		return liteclient.LiteServerAccountStateC{}, err
	}
	return res, nil
}

func decodeAccountDataFromProof(bocBytes []byte, account ton.AccountID) (uint64, tlb.Bits256, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return 0, tlb.Bits256{}, err
	}
	if len(cells) < 2 {
		return 0, tlb.Bits256{}, fmt.Errorf("must be at least two root cells")
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[1], &proof) // cells order must be strictly defined
	if err != nil {
		return 0, tlb.Bits256{}, err
	}
	values := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Values()
	keys := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Keys()
	for i, k := range keys {
		if bytes.Equal(k[:], account.Address[:]) {
			return values[i].LastTransLt, values[i].LastTransHash, nil
		}
	}
	return 0, tlb.Bits256{}, fmt.Errorf("account not found in ShardAccounts")
}

func (c *Client) GetShardInfo(
	ctx context.Context,
	blockID ton.BlockIDExt,
	workchain int32,
	shard uint64,
	exact bool,
) (ton.BlockIDExt, error) {
	res, err := c.GetShardInfoRaw(ctx, blockID, workchain, shard, exact)
	if err != nil {
		return ton.BlockIDExt{}, err
	}
	return res.Id.ToBlockIdExt(), nil
}

func (c *Client) GetShardInfoRaw(ctx context.Context, blockID ton.BlockIDExt, workchain int32, shard uint64, exact bool) (liteclient.LiteServerShardInfoC, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerShardInfoC{}, err
	}
	res, err := client.LiteServerGetShardInfo(ctx, liteclient.LiteServerGetShardInfoRequest{
		Id:        liteclient.BlockIDExt(blockID),
		Workchain: workchain,
		Shard:     shard,
		Exact:     exact,
	})
	if err != nil {
		return liteclient.LiteServerShardInfoC{}, err
	}
	return res, nil
}

func (c *Client) GetAllShardsInfo(ctx context.Context, blockID ton.BlockIDExt) ([]ton.BlockIDExt, error) {
	res, err := c.GetAllShardsInfoRaw(ctx, blockID)
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, boc.ErrNotSingleRoot
	}
	var inf tlb.AllShardsInfo
	err = tlb.Unmarshal(cells[0], &inf)
	if err != nil {
		return nil, err
	}
	var shards []ton.BlockIDExt
	for i, v := range inf.ShardHashes.Values() {
		wc := inf.ShardHashes.Keys()[i]
		for _, vv := range v.Value.BinTree.Values {
			shards = append(shards, ton.ToBlockId(vv, int32(wc)))
		}
	}
	return shards, nil
}

func (c *Client) GetAllShardsInfoRaw(ctx context.Context, blockID ton.BlockIDExt) (liteclient.LiteServerAllShardsInfoC, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerAllShardsInfoC{}, err
	}
	res, err := client.LiteServerGetAllShardsInfo(ctx, liteclient.LiteServerGetAllShardsInfoRequest{
		Id: liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerAllShardsInfoC{}, err
	}
	return res, nil
}

func (c *Client) GetOneTransactionFromBlock(
	ctx context.Context,
	accountID ton.AccountID,
	blockId ton.BlockIDExt,
	lt uint64,
) (ton.Transaction, error) {
	client, _, err := c.pool.BestClientByAccountID(ctx, accountID, false)
	if err != nil {
		return ton.Transaction{}, err
	}
	r, err := client.LiteServerGetOneTransaction(ctx, liteclient.LiteServerGetOneTransactionRequest{
		Id:      liteclient.BlockIDExt(blockId),
		Account: liteclient.AccountID(accountID),
		Lt:      lt,
	})
	if err != nil {
		return ton.Transaction{}, err
	}
	if len(r.Transaction) == 0 {
		return ton.Transaction{}, fmt.Errorf("transaction not found")
	}
	cells, err := boc.DeserializeBoc(r.Transaction)
	if err != nil {
		return ton.Transaction{}, err
	}
	if len(cells) != 1 {
		return ton.Transaction{}, boc.ErrNotSingleRoot
	}
	var t tlb.Transaction
	err = tlb.Unmarshal(cells[0], &t)
	return ton.Transaction{Transaction: t, BlockID: r.Id.ToBlockIdExt()}, err
}

func (c *Client) GetTransactions(
	ctx context.Context,
	count uint32,
	accountID ton.AccountID,
	lt uint64,
	hash ton.Bits256,
) ([]ton.Transaction, error) {
	r, err := c.GetTransactionsRaw(ctx, count, accountID, lt, hash)
	if err != nil {
		return nil, err
	}
	if len(r.Transactions) == 0 {
		return []ton.Transaction{}, nil
	}
	cells, err := boc.DeserializeBoc(r.Transactions)
	if err != nil {
		return nil, err
	}
	var res []ton.Transaction
	for i, cell := range cells {
		var t tlb.Transaction
		cell.ResetCounters()
		err := tlb.Unmarshal(cell, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, ton.Transaction{
			Transaction: t,
			BlockID:     r.Ids[i].ToBlockIdExt(),
		})
	}
	return res, nil
}

func (c *Client) GetTransactionsRaw(ctx context.Context, count uint32, accountID ton.AccountID, lt uint64, hash ton.Bits256) (liteclient.LiteServerTransactionListC, error) {
	archiveRequired := false
	for {
		client, _, err := c.pool.BestClientByAccountID(ctx, accountID, archiveRequired)
		if err != nil {
			return liteclient.LiteServerTransactionListC{}, err
		}
		res, err := client.LiteServerGetTransactions(ctx, liteclient.LiteServerGetTransactionsRequest{
			Count:   count,
			Account: liteclient.AccountID(accountID),
			Lt:      lt,
			Hash:    tl.Int256(hash),
		})
		if truncatedHistory(err) {
			if !c.archiveDetectionEnabled {
				return liteclient.LiteServerTransactionListC{}, err
			}
			if archiveRequired {
				return liteclient.LiteServerTransactionListC{}, err
			}
			archiveRequired = true
			continue
		}
		if err != nil {
			return liteclient.LiteServerTransactionListC{}, err
		}
		return res, nil

	}
}

func truncatedHistory(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(liteclient.LiteServerErrorC)
	return ok && int32(e.Code) == -400
}

func (c *Client) GetLastTransactions(ctx context.Context, a ton.AccountID, limit int) ([]ton.Transaction, error) {
	state, err := c.GetAccountState(ctx, a)
	if err != nil {
		return nil, err
	}
	var res []ton.Transaction
	lastLt, lastHash := state.LastTransLt, state.LastTransHash
	for {
		if lastLt == 0 {
			break
		}
		transactionCount := maxTransactionCount
		if limit-len(res) < transactionCount {
			transactionCount = limit - len(res)
		}
		txs, err := c.GetTransactions(ctx, uint32(transactionCount), a, lastLt, ton.Bits256(lastHash))
		if err != nil {
			if e, ok := err.(liteclient.LiteServerErrorC); ok && int32(e.Code) == -400 { // liteserver can store not full history. in that case it return error -400 for old transactions
				break
			}
			return nil, err
		}
		if len(txs) == 0 {
			break
		}
		res = append(res, txs...)
		if len(res) >= limit {
			res = res[:limit]
			break
		}
		lastLt, lastHash = res[len(res)-1].PrevTransLt, res[len(res)-1].PrevTransHash
	}
	return res, nil
}

func (c *Client) ListBlockTransactions(
	ctx context.Context,
	blockID ton.BlockIDExt,
	mode, count uint32,
	after *liteclient.LiteServerTransactionId3C,
) ([]liteclient.LiteServerTransactionIdC, bool, error) {
	// TODO: replace with tongo types
	res, err := c.ListBlockTransactionsRaw(ctx, blockID, mode, count, after)
	if err != nil {
		return nil, false, err
	}
	return res.Ids, res.Incomplete, nil
}

func (c *Client) ListBlockTransactionsRaw(ctx context.Context, blockID ton.BlockIDExt, mode, count uint32, after *liteclient.LiteServerTransactionId3C) (liteclient.LiteServerBlockTransactionsC, error) {
	client, err := c.pool.BestClientByBlockID(ctx, blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockTransactionsC{}, err
	}
	res, err := client.LiteServerListBlockTransactions(ctx, liteclient.LiteServerListBlockTransactionsRequest{
		Id:    liteclient.BlockIDExt(blockID),
		Mode:  mode,
		Count: count,
		After: after,
	})
	if err != nil {
		return liteclient.LiteServerBlockTransactionsC{}, err
	}
	return res, nil
}

func (c *Client) GetBlockProof(
	ctx context.Context,
	knownBlock ton.BlockIDExt,
	targetBlock *ton.BlockIDExt,
) (liteclient.LiteServerPartialBlockProofC, error) {
	res, err := c.GetBlockProofRaw(ctx, knownBlock, targetBlock)
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	// TODO: decode block proof
	return res, nil
}

func (c *Client) GetBlockProofRaw(ctx context.Context, knownBlock ton.BlockIDExt, targetBlock *ton.BlockIDExt) (liteclient.LiteServerPartialBlockProofC, error) {
	var (
		err    error
		client *liteclient.Client
		mode   uint32 = 0
	)
	if targetBlock != nil {
		client, err = c.pool.BestClientByBlockID(ctx, targetBlock.BlockID)
		mode = 1
	} else {
		client, err = c.pool.BestClientByBlockID(ctx, knownBlock.BlockID)
	}
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	var tb *liteclient.TonNodeBlockIdExtC
	if targetBlock != nil {
		b := liteclient.BlockIDExt(*targetBlock)
		tb = &b
	}
	res, err := client.LiteServerGetBlockProof(ctx, liteclient.LiteServerGetBlockProofRequest{
		Mode:        mode,
		KnownBlock:  liteclient.BlockIDExt(knownBlock),
		TargetBlock: tb,
	})
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	return res, nil
}

// GetConfigAll returns a current configuration of the blockchain.
func (c *Client) GetConfigAll(ctx context.Context, mode ConfigMode) (tlb.ConfigParams, error) {
	res, err := c.GetConfigAllRaw(ctx, mode)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	return ton.DecodeConfigParams(res.ConfigProof)
}

func (c *Client) GetConfigAllRaw(ctx context.Context, mode ConfigMode) (liteclient.LiteServerConfigInfoC, error) {
	client, masterHead, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerConfigInfoC{}, err
	}
	res, err := client.LiteServerGetConfigAll(ctx, liteclient.LiteServerGetConfigAllRequest{
		Mode: uint32(mode),
		Id:   liteclient.BlockIDExt(c.targetBlockOr(masterHead)),
	})
	if err != nil {
		return liteclient.LiteServerConfigInfoC{}, err
	}
	return res, nil
}

func (c *Client) GetConfigParams(ctx context.Context, mode ConfigMode, paramList []uint32) (tlb.ConfigParams, error) {
	client, masterHead, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	r, err := client.LiteServerGetConfigParams(ctx, liteclient.LiteServerGetConfigParamsRequest{
		Mode:      uint32(mode),
		Id:        liteclient.BlockIDExt(c.targetBlockOr(masterHead)),
		ParamList: paramList,
	})
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	return ton.DecodeConfigParams(r.ConfigProof)
}

func (c *Client) GetValidatorStats(
	ctx context.Context,
	mode, limit uint32,
	startAfter *ton.Bits256,
	modifiedAfter *uint32,
) (*tlb.McStateExtra, error) {
	client, masterHead, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return nil, err
	}
	var sa *tl.Int256
	if startAfter != nil {
		b := tl.Int256(*startAfter)
		sa = &b
	}
	r, err := client.LiteServerGetValidatorStats(ctx, liteclient.LiteServerGetValidatorStatsRequest{
		Mode:          mode,
		Id:            liteclient.BlockIDExt(c.targetBlockOr(masterHead)),
		Limit:         limit,
		StartAfter:    sa,
		ModifiedAfter: modifiedAfter,
	})
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(r.DataProof)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.ShardState] // TODO: or ton.ShardStateUnsplit
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}
	// TODO: extract validator stats params from ShardState
	// return &proof.Proof.VirtualRoot, nil //shards, nil
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetLibraries(ctx context.Context, libraryList []ton.Bits256) (map[ton.Bits256]*boc.Cell, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return nil, err
	}
	var ll []tl.Int256
	for _, l := range libraryList {
		ll = append(ll, tl.Int256(l))
	}
	r, err := client.LiteServerGetLibraries(ctx, liteclient.LiteServerGetLibrariesRequest{
		LibraryList: ll,
	})
	if err != nil {
		return nil, err
	}
	libs := make(map[ton.Bits256]*boc.Cell, len(r.Result))
	for _, lib := range r.Result {
		data, err := boc.DeserializeBoc(lib.Data)
		if err != nil {
			return nil, err
		}
		if len(data) != 1 {
			return nil, fmt.Errorf("multiroot lib is not supported")
		}
		libs[ton.Bits256(lib.Hash)] = data[0]
	}
	return libs, nil
}

func (c *Client) GetShardBlockProof(ctx context.Context) (liteclient.LiteServerShardBlockProofC, error) {
	res, err := c.GetShardBlockProofRaw(ctx)
	if err != nil {
		return liteclient.LiteServerShardBlockProofC{}, err
	}
	return res, nil
}

func (c *Client) GetShardBlockProofRaw(ctx context.Context) (liteclient.LiteServerShardBlockProofC, error) {
	client, masterHead, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerShardBlockProofC{}, err
	}
	return client.LiteServerGetShardBlockProof(ctx, liteclient.LiteServerGetShardBlockProofRequest{
		Id: liteclient.BlockIDExt(c.targetBlockOr(masterHead)),
	})
}

// WaitMasterchainSeqno waits for a masterchain block with the given seqno.
// If any connection in the pool becomes aware of this seqno, the function returns.
// If the timeout is reached, the function returns an error.
func (c *Client) WaitMasterchainSeqno(ctx context.Context, seqno uint32, timeout time.Duration) error {
	return c.pool.WaitMasterchainSeqno(ctx, seqno, timeout)
}

func (c *Client) GetOutMsgQueueSizes(ctx context.Context) (liteclient.LiteServerOutMsgQueueSizesC, error) {
	client, _, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return liteclient.LiteServerOutMsgQueueSizesC{}, err
	}
	res, err := client.LiteServerGetOutMsgQueueSizes(ctx, liteclient.LiteServerGetOutMsgQueueSizesRequest{})
	if err != nil {
		return liteclient.LiteServerOutMsgQueueSizesC{}, err
	}
	return res, nil
}

var configCache = make(map[string]*config.GlobalConfigurationFile)
var configCacheMutex sync.RWMutex

func downloadConfig(path string) (*config.GlobalConfigurationFile, error) {
	configCacheMutex.RLock()
	o, prs := configCache[path]
	configCacheMutex.RUnlock()
	if prs {
		return o, nil
	}
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	o, err = config.ParseConfig(resp.Body)
	if err == nil {
		configCacheMutex.Lock()
		configCache[path] = o
		configCacheMutex.Unlock()
	}
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(o.LiteServers), func(i, j int) {
		o.LiteServers[i], o.LiteServers[j] = o.LiteServers[j], o.LiteServers[i]
	})
	return o, nil
}

func (c *Client) getNetworkGlobalID() *int32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.networkGlobalID
}

func (c *Client) GetNetworkGlobalID(ctx context.Context) (int32, error) {
	if networkID := c.getNetworkGlobalID(); networkID != nil {
		return *networkID, nil
	}
	_, blockID, err := c.pool.BestMasterchainClient(ctx)
	if err != nil {
		return 0, err
	}
	block, err := c.GetBlock(ctx, blockID)
	if err != nil {
		return 0, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.networkGlobalID = &block.GlobalId
	return block.GlobalId, nil
}
