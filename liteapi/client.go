package liteapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	mrand "math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/utils"
)

const (
	LiteServerEnvName = "LITE_SERVERS"
)

var (
	// ErrAccountNotFound is returned by lite server when executing a method for an account that has not been deployed to the blockchain.
	ErrAccountNotFound = errors.New("account not found")
)

type connection struct {
	workchain   int32
	shardPrefix tongo.ShardID
	client      *liteclient.Client
}

type Client struct {
	connectionPool                 []connection
	timeout                        time.Duration
	targetBlockID                  *tongo.BlockIDExt
	masterchainLastBlockCache      *liteclient.TonNodeBlockIdExtC
	masterchainLastBlockUpdateTime time.Time
}

// Options holds parameters to configure a lite api instance.
type Options struct {
	LiteServers []config.LiteServer
	Timeout     time.Duration
}

type Option func(o *Options) error

func WithLiteServers(servers []config.LiteServer) Option {
	return func(o *Options) error {
		o.LiteServers = servers
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = timeout
		return nil
	}
}

// FromEnvs configures a Client based on the following environment variables:
// LITE_SERVERS - a list of lite servers to use
func FromEnvs() Option {
	return func(o *Options) error {
		if value, ok := os.LookupEnv(LiteServerEnvName); ok {
			servers, err := config.ParseLiteServersEnvVar(value)
			if err != nil {
				return err
			}
			o.LiteServers = servers
		}
		return nil
	}
}

// Mainnet configures a client to use lite servers from the mainnet.
func Mainnet() Option {
	return func(o *Options) error {
		file, err := downloadConfig("https://ton-blockchain.github.io/global.config.json")
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
		file, err := downloadConfig("https://ton-blockchain.github.io/testnet-global.config.json")
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
func NewClient(opts ...Option) (*Client, error) {
	options := &Options{
		Timeout: 60 * time.Second,
	}
	for _, o := range opts {
		if err := o(options); err != nil {
			return nil, err
		}
	}
	if len(options.LiteServers) == 0 {
		return nil, fmt.Errorf("server list empty")
	}
	client := Client{
		timeout: options.Timeout,
	}
	for _, ls := range options.LiteServers {
		serverPubkey, err := base64.StdEncoding.DecodeString(ls.Key)
		if err != nil {
			continue
		}
		c, err := liteclient.NewConnection(context.Background(), serverPubkey, ls.Host)
		if err != nil {
			continue
		}
		client.connectionPool = append(client.connectionPool, connection{
			workchain:   0,
			shardPrefix: tongo.MustParseShardID(-0x8000000000000000),
			client:      liteclient.NewClient(c, liteclient.OptionTimeout(options.Timeout)),
		})
		go client.refreshMasterchainTask()
		return &client, nil
	}
	return nil, fmt.Errorf("all liteservers are unavailable")
}

func (c *Client) getMasterchainServer() *liteclient.Client {
	return c.connectionPool[mrand.Intn(len(c.connectionPool))].client
}

func (c *Client) getServerByAccountID(a tongo.AccountID) (*liteclient.Client, error) {
	if a.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.connectionPool {
		if server.workchain != a.Workchain {
			continue
		}
		if server.shardPrefix.MatchAccountID(a) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for account %v", a.ToRaw())
}

func (c *Client) getServerByBlockID(block tongo.BlockID) (*liteclient.Client, error) {
	if block.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.connectionPool {
		if server.shardPrefix.MatchBlockID(block) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for block %v", block.String())
}

func (c *Client) targetBlock(ctx context.Context) (tongo.BlockIDExt, error) {
	if c.targetBlockID != nil {
		return *c.targetBlockID, nil
	}
	if time.Since(c.masterchainLastBlockUpdateTime) < 20*time.Second && c.masterchainLastBlockCache != nil {
		return c.masterchainLastBlockCache.ToBlockIdExt(), nil
	}
	r, err := c.getMasterchainServer().LiteServerGetMasterchainInfo(context.TODO())
	if err != nil {
		return tongo.BlockIDExt{}, err
	}
	return r.Last.ToBlockIdExt(), nil
}

func (c *Client) getlastBlock(ctx context.Context) (liteclient.TonNodeBlockIdExtC, error) {
	info, err := c.getMasterchainServer().LiteServerGetMasterchainInfo(ctx)
	if err != nil {
		return liteclient.TonNodeBlockIdExtC{}, err
	}
	return info.Last, nil
}

func (c *Client) refreshMasterchainTask() {
	for {
		time.Sleep(time.Second) //todo: switch to wait function
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		block, err := c.getlastBlock(ctx)
		cancel()
		if err != nil {
			continue
		}
		c.masterchainLastBlockCache = &block
	}
}

func (c *Client) WithBlock(block tongo.BlockIDExt) *Client {
	return &Client{
		connectionPool:                 c.connectionPool,
		timeout:                        c.timeout,
		targetBlockID:                  &block,
		masterchainLastBlockCache:      c.masterchainLastBlockCache,
		masterchainLastBlockUpdateTime: c.masterchainLastBlockUpdateTime,
	}
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfo(ctx)
}

func (c *Client) GetMasterchainInfoExt(ctx context.Context, mode uint32) (liteclient.LiteServerMasterchainInfoExtC, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfoExt(ctx, liteclient.LiteServerGetMasterchainInfoExtRequest{Mode: mode})
}

func (c *Client) GetTime(ctx context.Context) (uint32, error) {
	res, err := c.getMasterchainServer().LiteServerGetTime(ctx)
	return res.Now, err
}

func (c *Client) GetVersion(ctx context.Context) (liteclient.LiteServerVersionC, error) {
	return c.getMasterchainServer().LiteServerGetVersion(ctx)
}

func (c *Client) GetBlock(ctx context.Context, blockID tongo.BlockIDExt) (tlb.Block, error) {
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
	var data tlb.Block
	err = tlb.NewDecoder().Unmarshal(cells[0], &data)
	if err != nil {
		return tlb.Block{}, err
	}
	return data, nil
}

func (c *Client) GetBlockRaw(ctx context.Context, blockID tongo.BlockIDExt) (liteclient.LiteServerBlockDataC, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockDataC{}, err
	}
	res, err := server.LiteServerGetBlock(ctx, liteclient.LiteServerGetBlockRequest{liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerBlockDataC{}, err
	}
	return res, err
}

func (c *Client) GetState(ctx context.Context, blockID tongo.BlockIDExt) ([]byte, tongo.Bits256, tongo.Bits256, error) {
	res, err := c.GetStateRaw(ctx, blockID)
	if err != nil {
		return nil, tongo.Bits256{}, tongo.Bits256{}, err
	}
	// TODO: implement state tlb decoding
	return res.Data, tongo.Bits256(res.RootHash), tongo.Bits256(res.FileHash), nil
}

func (c *Client) GetStateRaw(ctx context.Context, blockID tongo.BlockIDExt) (liteclient.LiteServerBlockStateC, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockStateC{}, err
	}
	res, err := server.LiteServerGetState(ctx, liteclient.LiteServerGetStateRequest{Id: liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerBlockStateC{}, err
	}
	return res, nil
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID tongo.BlockIDExt, mode uint32) (tlb.BlockInfo, error) {
	res, err := c.GetBlockHeaderRaw(ctx, blockID, mode)
	if err != nil {
		return tlb.BlockInfo{}, err
	}
	_, info, err := decodeBlockHeader(res)
	return info, err
}

func (c *Client) GetBlockHeaderRaw(ctx context.Context, blockID tongo.BlockIDExt, mode uint32) (liteclient.LiteServerBlockHeaderC, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockHeaderC{}, err
	}
	res, err := server.LiteServerGetBlockHeader(ctx, liteclient.LiteServerGetBlockHeaderRequest{
		Id:   liteclient.BlockIDExt(blockID),
		Mode: mode,
	})
	if err != nil {
		return liteclient.LiteServerBlockHeaderC{}, err
	}
	return res, nil
}

func (c *Client) LookupBlock(ctx context.Context, blockID tongo.BlockID, mode uint32, lt *uint64, utime *uint32) (tongo.BlockIDExt, tlb.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID)
	if err != nil {
		return tongo.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	res, err := server.LiteServerLookupBlock(ctx, liteclient.LiteServerLookupBlockRequest{
		Mode: mode,
		Id: liteclient.TonNodeBlockIdC{
			Workchain: uint32(blockID.Workchain),
			Shard:     blockID.Shard,
			Seqno:     blockID.Seqno,
		},
		Lt:    lt,
		Utime: utime,
	})
	if err != nil {
		return tongo.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func decodeBlockHeader(header liteclient.LiteServerBlockHeaderC) (tongo.BlockIDExt, tlb.BlockInfo, error) {
	cells, err := boc.DeserializeBoc(header.HeaderProof)
	if err != nil {
		return tongo.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	if len(cells) != 1 {
		return tongo.BlockIDExt{}, tlb.BlockInfo{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.BlockHeader]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tongo.BlockIDExt{}, tlb.BlockInfo{}, err
	}
	return header.Id.ToBlockIdExt(), proof.Proof.VirtualRoot.Info, nil // TODO: maybe decode more
}

// SendMessage verifies that the given payload contains an external message and sends it to a lite server.
func (c *Client) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	if err := VerifySendMessagePayload(payload); err != nil {
		return 0, err
	}
	res, err := c.getMasterchainServer().LiteServerSendMessage(ctx, liteclient.LiteServerSendMessageRequest{Body: payload})
	return res.Status, err
}

func (c *Client) RunSmcMethodByID(ctx context.Context, accountID tongo.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, params)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	b, err := cell.ToBoc()
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	id, err := c.targetBlock(ctx)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	req := liteclient.LiteServerRunSmcMethodRequest{
		Mode:     4,
		Id:       liteclient.BlockIDExt(id),
		Account:  liteclient.AccountID(accountID),
		MethodId: uint64(methodID),
		Params:   b,
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	res, err := server.LiteServerRunSmcMethod(ctx, req)
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
	accountID tongo.AccountID,
	method string,
	params tlb.VmStack,
) (uint32, tlb.VmStack, error) {
	return c.RunSmcMethodByID(ctx, accountID, utils.MethodIdFromName(method), params)
}

func (c *Client) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tlb.ShardAccount, error) {
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

func (c *Client) GetAccountStateRaw(ctx context.Context, accountID tongo.AccountID) (liteclient.LiteServerAccountStateC, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return liteclient.LiteServerAccountStateC{}, err
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return liteclient.LiteServerAccountStateC{}, err
	}
	res, err := server.LiteServerGetAccountState(ctx, liteclient.LiteServerGetAccountStateRequest{
		Account: liteclient.AccountID(accountID),
		Id:      liteclient.BlockIDExt(id),
	})
	if err != nil {
		return liteclient.LiteServerAccountStateC{}, err
	}
	return res, nil
}

func decodeAccountDataFromProof(bocBytes []byte, account tongo.AccountID) (uint64, tlb.Bits256, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return 0, tlb.Bits256{}, err
	}
	if len(cells) < 1 {
		return 0, tlb.Bits256{}, fmt.Errorf("must be at least one root cell")
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
	blockID tongo.BlockIDExt,
	workchain uint32,
	shard uint64,
	exact bool,
) (tongo.BlockIDExt, error) {
	res, err := c.GetShardInfoRaw(ctx, blockID, workchain, shard, exact)
	if err != nil {
		return tongo.BlockIDExt{}, err
	}
	return res.Id.ToBlockIdExt(), nil
}

func (c *Client) GetShardInfoRaw(ctx context.Context, blockID tongo.BlockIDExt, workchain uint32, shard uint64, exact bool) (liteclient.LiteServerShardInfoC, error) {
	res, err := c.getMasterchainServer().LiteServerGetShardInfo(ctx, liteclient.LiteServerGetShardInfoRequest{
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

func (c *Client) GetAllShardsInfo(ctx context.Context, blockID tongo.BlockIDExt) ([]tongo.BlockIDExt, error) {
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
	var shards []tongo.BlockIDExt
	for i, v := range inf.ShardHashes.Values() {
		wc := inf.ShardHashes.Keys()[i]
		for _, vv := range v.Value.BinTree.Values {
			shards = append(shards, tongo.ToBlockId(vv, int32(wc)))
		}
	}
	return shards, nil
}

func (c *Client) GetAllShardsInfoRaw(ctx context.Context, blockID tongo.BlockIDExt) (liteclient.LiteServerAllShardsInfoC, error) {
	res, err := c.getMasterchainServer().LiteServerGetAllShardsInfo(ctx, liteclient.LiteServerGetAllShardsInfoRequest{
		Id: liteclient.BlockIDExt(blockID)})
	if err != nil {
		return liteclient.LiteServerAllShardsInfoC{}, err
	}
	return res, nil
}

func (c *Client) GetOneTransactionFromBlock(
	ctx context.Context,
	accountID tongo.AccountID,
	blockId tongo.BlockIDExt,
	lt uint64,
) (tongo.Transaction, error) {
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return tongo.Transaction{}, err
	}
	r, err := server.LiteServerGetOneTransaction(ctx, liteclient.LiteServerGetOneTransactionRequest{
		Id:      liteclient.BlockIDExt(blockId),
		Account: liteclient.AccountID(accountID),
		Lt:      lt,
	})
	if err != nil {
		return tongo.Transaction{}, err
	}
	if len(r.Transaction) == 0 {
		return tongo.Transaction{}, fmt.Errorf("transaction not found")
	}
	cells, err := boc.DeserializeBoc(r.Transaction)
	if err != nil {
		return tongo.Transaction{}, err
	}
	if len(cells) != 1 {
		return tongo.Transaction{}, boc.ErrNotSingleRoot
	}
	var t tlb.Transaction
	err = tlb.Unmarshal(cells[0], &t)
	return tongo.Transaction{Transaction: t, BlockID: r.Id.ToBlockIdExt()}, err
}

func (c *Client) GetTransactions(
	ctx context.Context,
	count uint32,
	accountID tongo.AccountID,
	lt uint64,
	hash tongo.Bits256,
) ([]tongo.Transaction, error) {
	r, err := c.GetTransactionsRaw(ctx, count, accountID, lt, hash)
	if err != nil {
		return nil, err
	}
	if len(r.Transactions) == 0 {
		return []tongo.Transaction{}, nil
	}
	cells, err := boc.DeserializeBoc(r.Transactions)
	if err != nil {
		return nil, err
	}
	var res []tongo.Transaction
	for i, cell := range cells {
		var t tlb.Transaction
		cell.ResetCounters()
		err := tlb.Unmarshal(cell, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, tongo.Transaction{
			Transaction: t,
			BlockID:     r.Ids[i].ToBlockIdExt(),
		})
	}
	return res, nil
}

func (c *Client) GetTransactionsRaw(ctx context.Context, count uint32, accountID tongo.AccountID, lt uint64, hash tongo.Bits256) (liteclient.LiteServerTransactionListC, error) {
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return liteclient.LiteServerTransactionListC{}, err
	}
	res, err := server.LiteServerGetTransactions(ctx, liteclient.LiteServerGetTransactionsRequest{
		Count:   count,
		Account: liteclient.AccountID(accountID),
		Lt:      lt,
		Hash:    tl.Int256(hash),
	})
	if err != nil {
		return liteclient.LiteServerTransactionListC{}, err
	}
	return res, nil
}

func (c *Client) GetLastTransactions(ctx context.Context, a tongo.AccountID, limit int) ([]tongo.Transaction, error) {
	state, err := c.GetAccountState(ctx, a)
	if err != nil {
		return nil, err
	}
	lastLt, lastHash := state.LastTransLt, state.LastTransHash
	var res []tongo.Transaction
	for {
		if lastLt == 0 {
			break
		}
		txs, err := c.GetTransactions(ctx, 10, a, lastLt, tongo.Bits256(lastHash))
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
	blockID tongo.BlockIDExt,
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

func (c *Client) ListBlockTransactionsRaw(ctx context.Context, blockID tongo.BlockIDExt, mode, count uint32, after *liteclient.LiteServerTransactionId3C) (liteclient.LiteServerBlockTransactionsC, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return liteclient.LiteServerBlockTransactionsC{}, err
	}
	res, err := server.LiteServerListBlockTransactions(ctx, liteclient.LiteServerListBlockTransactionsRequest{
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
	knownBlock tongo.BlockIDExt,
	targetBlock *tongo.BlockIDExt,
) (liteclient.LiteServerPartialBlockProofC, error) {
	res, err := c.GetBlockProofRaw(ctx, knownBlock, targetBlock)
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	// TODO: decode block proof
	return res, nil
}

func (c *Client) GetBlockProofRaw(ctx context.Context, knownBlock tongo.BlockIDExt, targetBlock *tongo.BlockIDExt) (liteclient.LiteServerPartialBlockProofC, error) {
	var (
		err    error
		server *liteclient.Client
		mode   uint32 = 0
	)
	if targetBlock != nil {
		server, err = c.getServerByBlockID(targetBlock.BlockID)
		mode = 1
	} else {
		server, err = c.getServerByBlockID(knownBlock.BlockID)
	}
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	var tb *liteclient.TonNodeBlockIdExtC
	if targetBlock != nil {
		b := liteclient.BlockIDExt(*targetBlock)
		tb = &b
	}
	res, err := server.LiteServerGetBlockProof(ctx, liteclient.LiteServerGetBlockProofRequest{
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
	return decodeConfigParams(res.ConfigProof)
}

func (c *Client) GetConfigAllRaw(ctx context.Context, mode ConfigMode) (liteclient.LiteServerConfigInfoC, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return liteclient.LiteServerConfigInfoC{}, err
	}
	res, err := c.getMasterchainServer().LiteServerGetConfigAll(ctx, liteclient.LiteServerGetConfigAllRequest{
		Mode: uint32(mode),
		Id:   liteclient.BlockIDExt(id),
	})
	if err != nil {
		return liteclient.LiteServerConfigInfoC{}, err
	}
	return res, nil
}

func (c *Client) GetConfigParams(ctx context.Context, mode ConfigMode, paramList []uint32) (tlb.ConfigParams, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	r, err := c.getMasterchainServer().LiteServerGetConfigParams(ctx, liteclient.LiteServerGetConfigParamsRequest{
		Mode:      uint32(mode),
		Id:        liteclient.BlockIDExt(id),
		ParamList: paramList,
	})
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	return decodeConfigParams(r.ConfigProof)
}

func decodeConfigParams(b []byte) (tlb.ConfigParams, error) {
	cells, err := boc.DeserializeBoc(b)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	if len(cells) != 1 {
		return tlb.ConfigParams{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tlb.MerkleProof[tlb.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tlb.ConfigParams{}, err
	}
	if proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Exists {
		return proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Value.Value.Config, nil
	}
	return tlb.ConfigParams{}, fmt.Errorf("empty Custom field")
}

func (c *Client) GetValidatorStats(
	ctx context.Context,
	mode, limit uint32,
	startAfter *tongo.Bits256,
	modifiedAfter *uint32,
) (*tlb.McStateExtra, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return nil, err
	}
	var sa *tl.Int256
	if startAfter != nil {
		b := tl.Int256(*startAfter)
		sa = &b
	}
	r, err := c.getMasterchainServer().LiteServerGetValidatorStats(ctx, liteclient.LiteServerGetValidatorStatsRequest{
		Mode:          mode,
		Id:            liteclient.BlockIDExt(id),
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
		Proof tlb.MerkleProof[tlb.ShardState] // TODO: or tongo.ShardStateUnsplit
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}
	// TODO: extract validator stats params from ShardState
	// return &proof.Proof.VirtualRoot, nil //shards, nil
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetLibraries(ctx context.Context, libraryList []tongo.Bits256) ([]liteclient.LiteServerLibraryEntryC, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return nil, err
	}
	server, err := c.getServerByBlockID(id.BlockID)
	if err != nil {
		return nil, err
	}
	var ll []tl.Int256
	for _, l := range libraryList {
		ll = append(ll, tl.Int256(l))
	}
	r, err := server.LiteServerGetLibraries(ctx, liteclient.LiteServerGetLibrariesRequest{
		LibraryList: ll,
	})
	if err != nil {
		return nil, err
	}
	// TODO: replace with tongo type
	return r.Result, nil
}

func (c *Client) GetShardBlockProof(ctx context.Context) (liteclient.LiteServerShardBlockProofC, error) {
	res, err := c.GetShardBlockProofRaw(ctx)
	if err != nil {
		return liteclient.LiteServerShardBlockProofC{}, err
	}
	return res, nil
}

func (c *Client) GetShardBlockProofRaw(ctx context.Context) (liteclient.LiteServerShardBlockProofC, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return liteclient.LiteServerShardBlockProofC{}, err
	}
	server, err := c.getServerByBlockID(id.BlockID)
	if err != nil {
		return liteclient.LiteServerShardBlockProofC{}, err
	}
	return server.LiteServerGetShardBlockProof(ctx, liteclient.LiteServerGetShardBlockProofRequest{
		Id: liteclient.BlockIDExt(id),
	})
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
	return o, err
}
