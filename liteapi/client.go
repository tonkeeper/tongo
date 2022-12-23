package liteapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	mrand "math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/config"
	"github.com/startfellows/tongo/liteclient"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
)

type connection struct {
	workchain   int32
	shardPrefix tongo.ShardID
	client      *liteclient.Client
}

type Client struct {
	connectionPool                 []connection
	targetBlockID                  *tongo.BlockIDExt
	masterchainLastBlockCache      *liteclient.TonNodeBlockIdExtC
	masterchainLastBlockUpdateTime time.Time
}

func NewClientWithDefaultMainnet() (*Client, error) {
	options, err := downloadConfig("https://ton-blockchain.github.io/global.config.json")
	if err != nil {
		return nil, err
	}
	return NewClient(*options)
}

func NewClientWithDefaultTestnet() (*Client, error) {
	options, err := downloadConfig("https://ton-blockchain.github.io/testnet-global.config.json")
	if err != nil {
		return nil, err
	}
	return NewClient(*options)
}

// NewClient
// Get options and create new lite client. If no options provided - download public config for mainnet from ton.org.
func NewClient(options config.Options) (*Client, error) {
	if len(options.LiteServers) == 0 {
		return nil, fmt.Errorf("server list empty")
	}
	client := Client{}
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
			client:      liteclient.NewClient(c),
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

func (c *Client) GetBlock(ctx context.Context, blockID tongo.BlockIDExt) (tongo.Block, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return tongo.Block{}, err
	}
	res, err := server.LiteServerGetBlock(ctx, liteclient.LiteServerGetBlockRequest{liteclient.BlockIDExt(blockID)})
	if err != nil {
		return tongo.Block{}, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return tongo.Block{}, err
	}
	if len(cells) != 1 {
		return tongo.Block{}, boc.ErrNotSingleRoot
	}
	var data tongo.Block
	err = tlb.Unmarshal(cells[0], &data)
	if err != nil {
		return tongo.Block{}, err
	}
	return data, nil
}

func (c *Client) GetState(ctx context.Context, blockID tongo.BlockIDExt) ([]byte, tongo.Hash, tongo.Hash, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return nil, tongo.Hash{}, tongo.Hash{}, err
	}
	res, err := server.LiteServerGetState(ctx, liteclient.LiteServerGetStateRequest{Id: liteclient.BlockIDExt(blockID)})
	if err != nil {
		return nil, tongo.Hash{}, tongo.Hash{}, err
	}
	// TODO: implement state tlb decoding
	return res.Data, tongo.Hash(res.RootHash), tongo.Hash(res.FileHash), nil
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID tongo.BlockIDExt, mode uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	res, err := server.LiteServerGetBlockHeader(ctx, liteclient.LiteServerGetBlockHeaderRequest{
		Id:   liteclient.BlockIDExt(blockID),
		Mode: mode,
	})
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func (c *Client) LookupBlock(ctx context.Context, blockID tongo.BlockID, mode uint32, lt *uint64, utime *uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID)
	if err != nil {
		return tongo.BlockInfo{}, err
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
		return tongo.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func decodeBlockHeader(header liteclient.LiteServerBlockHeaderC) (tongo.BlockInfo, error) {
	cells, err := boc.DeserializeBoc(header.HeaderProof)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	if len(cells) != 1 {
		return tongo.BlockInfo{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.BlockHeader]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return proof.Proof.VirtualRoot.Info, nil // TODO: maybe decode more
}

func (c *Client) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	res, err := c.getMasterchainServer().LiteServerSendMessage(ctx, liteclient.LiteServerSendMessageRequest{Body: payload})
	return res.Status, err
}

func (c *Client) RunSmcMethod(
	ctx context.Context,
	accountID tongo.AccountID,
	method string,
	params tongo.VmStack,
) (uint32, tongo.VmStack, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, params)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	b, err := cell.ToBoc()
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	id, err := c.targetBlock(ctx)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	req := liteclient.LiteServerRunSmcMethodRequest{
		Mode:     4,
		Id:       liteclient.BlockIDExt(id),
		Account:  liteclient.AccountID(accountID),
		MethodId: uint64(utils.Crc16String(method)&0xffff) | 0x10000,
		Params:   b,
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	res, err := server.LiteServerRunSmcMethod(ctx, req)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	cells, err := boc.DeserializeBoc(res.Result)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	if len(cells) != 1 {
		return 0, tongo.VmStack{}, boc.ErrNotSingleRoot
	}
	var result tongo.VmStack
	err = tlb.Unmarshal(cells[0], &result)
	return res.ExitCode, result, err
}

func (c *Client) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tongo.ShardAccount, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	res, err := server.LiteServerGetAccountState(ctx, liteclient.LiteServerGetAccountStateRequest{
		Account: liteclient.AccountID(accountID),
		Id:      liteclient.BlockIDExt(id),
	})
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	if len(res.State) == 0 {
		return tongo.ShardAccount{Account: tongo.Account{SumType: "AccountNone"}}, nil
	}
	cells, err := boc.DeserializeBoc(res.State)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	if len(cells) != 1 {
		return tongo.ShardAccount{}, boc.ErrNotSingleRoot
	}
	var acc tongo.Account
	err = tlb.Unmarshal(cells[0], &acc)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	//lt, hash, err := decodeAccountDataFromProof(res.Proof, accountID)
	// TODO: fix tlb decoding of Account
	return tongo.ShardAccount{Account: acc}, err
}

func decodeAccountDataFromProof(bocBytes []byte, account tongo.AccountID) (uint64, tongo.Hash, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return 0, tongo.Hash{}, err
	}
	if len(cells) < 1 {
		return 0, tongo.Hash{}, fmt.Errorf("must be at least one root cell")
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[1], &proof) // cells order must be strictly defined
	if err != nil {
		return 0, tongo.Hash{}, err
	}
	values := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Values()
	keys := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Keys()
	for i, k := range keys {
		keyVal, err := k.ReadBytes(32)
		if err != nil {
			return 0, tongo.Hash{}, err
		}
		if bytes.Equal(keyVal, account.Address[:]) {
			return values[i].LastTransLt, values[i].LastTransHash, nil
		}
	}
	return 0, tongo.Hash{}, fmt.Errorf("account not found in ShardAccounts")
}

func (c *Client) GetShardInfo(
	ctx context.Context,
	blockID tongo.BlockIDExt,
	workchain uint32,
	shard uint64,
	exact bool,
) (tongo.BlockIDExt, error) {
	res, err := c.getMasterchainServer().LiteServerGetShardInfo(ctx, liteclient.LiteServerGetShardInfoRequest{
		Id:        liteclient.BlockIDExt(blockID),
		Workchain: workchain,
		Shard:     shard,
		Exact:     exact,
	})
	if err != nil {
		return tongo.BlockIDExt{}, err
	}
	return res.Id.ToBlockIdExt(), nil
}

func (c *Client) GetAllShardsInfo(ctx context.Context, blockID tongo.BlockIDExt) ([]tongo.BlockIDExt, error) {
	res, err := c.getMasterchainServer().LiteServerGetAllShardsInfo(ctx, liteclient.LiteServerGetAllShardsInfoRequest{
		Id: liteclient.BlockIDExt(blockID)})
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
	var inf tongo.AllShardsInfo
	err = tlb.Unmarshal(cells[0], &inf)
	if err != nil {
		return nil, err
	}
	var shards []tongo.BlockIDExt
	for i, v := range inf.ShardHashes.Values() {
		wc, err := inf.ShardHashes.Keys()[i].ReadUint(32)
		if err != nil {
			return nil, err
		}
		for _, vv := range v.Value.BinTree.Values {
			shards = append(shards, vv.ToBlockId(int32(wc)))
		}
	}
	return shards, nil
}

func (c *Client) GetOneTransaction(
	ctx context.Context,
	accountID tongo.AccountID,
	lt uint64,
) (tongo.Transaction, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return tongo.Transaction{}, err
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return tongo.Transaction{}, err
	}
	r, err := server.LiteServerGetOneTransaction(ctx, liteclient.LiteServerGetOneTransactionRequest{
		Id:      liteclient.BlockIDExt(id),
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
	var t tongo.Transaction
	err = tlb.Unmarshal(cells[0], &t)
	return t, err
}

func (c *Client) GetTransactions(
	ctx context.Context,
	count uint32,
	accountID tongo.AccountID,
	lt uint64,
	hash tongo.Hash,
) ([]tongo.Transaction, error) {
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	r, err := server.LiteServerGetTransactions(ctx, liteclient.LiteServerGetTransactionsRequest{
		Count:   count,
		Account: liteclient.AccountID(accountID),
		Lt:      lt,
		Hash:    tl.Int256(hash),
	})
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(r.Transactions)
	if err != nil {
		return nil, err
	}
	var res []tongo.Transaction
	for _, cell := range cells {
		var t tongo.Transaction
		cell.ResetCounters()
		err := tlb.Unmarshal(cell, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
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
	server, err := c.getServerByBlockID(blockID.BlockID)
	if err != nil {
		return nil, false, err
	}
	r, err := server.LiteServerListBlockTransactions(ctx, liteclient.LiteServerListBlockTransactionsRequest{
		Id:    liteclient.BlockIDExt(blockID),
		Mode:  mode,
		Count: count,
		After: after,
	})
	if err != nil {
		return nil, false, err
	}
	return r.Ids, r.Incomplete, nil
}

func (c *Client) GetBlockProof(
	ctx context.Context,
	knownBlock tongo.BlockIDExt,
	targetBlock *tongo.BlockIDExt,
) (liteclient.LiteServerPartialBlockProofC, error) {
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
	r, err := server.LiteServerGetBlockProof(ctx, liteclient.LiteServerGetBlockProofRequest{
		Mode:        mode,
		KnownBlock:  liteclient.BlockIDExt(knownBlock),
		TargetBlock: tb,
	})
	if err != nil {
		return liteclient.LiteServerPartialBlockProofC{}, err
	}
	// TODO: decode block proof
	return r, nil
}

func (c *Client) GetConfigAll(ctx context.Context, mode uint32) (*tongo.McStateExtra, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.getMasterchainServer().LiteServerGetConfigAll(ctx, liteclient.LiteServerGetConfigAllRequest{
		Mode: mode,
		Id:   liteclient.BlockIDExt(id),
	})
	if err != nil {
		return nil, err
	}
	return decodeConfigParams(r.ConfigProof)
}

func (c *Client) GetConfigParams(ctx context.Context, mode uint32, paramList []uint32) (*tongo.McStateExtra, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.getMasterchainServer().LiteServerGetConfigParams(ctx, liteclient.LiteServerGetConfigParamsRequest{
		Mode:      mode,
		Id:        liteclient.BlockIDExt(id),
		ParamList: paramList,
	})
	if err != nil {
		return nil, err
	}
	return decodeConfigParams(r.ConfigProof)
}

func decodeConfigParams(b []byte) (*tongo.McStateExtra, error) {
	cells, err := boc.DeserializeBoc(b)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardState]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}
	// TODO: extract config params from ShardState
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetValidatorStats(
	ctx context.Context,
	mode, limit uint32,
	startAfter *tongo.Hash,
	modifiedAfter *uint32,
) (*tongo.McStateExtra, error) {
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
		Proof tongo.MerkleProof[tongo.ShardState] // TODO: or tongo.ShardStateUnsplit
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}
	// TODO: extract validator stats params from ShardState
	// return &proof.Proof.VirtualRoot, nil //shards, nil
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetLibraries(ctx context.Context, libraryList []tongo.Hash) ([]liteclient.LiteServerLibraryEntryC, error) {
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

var configCache = make(map[string]*config.Options)
var configCacheMutex sync.RWMutex

func downloadConfig(path string) (*config.Options, error) {
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
