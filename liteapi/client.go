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
	targetBlockID                  *tongo.TonNodeBlockIdExt
	masterchainLastBlockCache      *tongo.TonNodeBlockIdExt
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
	// TODO: implement multiple server support

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

func (c *Client) getServerByBlockID(block tongo.TonNodeBlockId) (*liteclient.Client, error) {
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

func (c Client) targetBlock(ctx context.Context) (tongo.TonNodeBlockIdExt, error) {
	if c.targetBlockID != nil {
		return *c.targetBlockID, nil
	}
	if time.Since(c.masterchainLastBlockUpdateTime) < 20*time.Second {
		return c.masterchainLastBlockCache, nil
	}
	return c.getlastBlock(ctx)
}

func (c *Client) refreshMasterchainTask() {
	// TODO: implement
	for {
		c.masterchainLastBlockCache = blahblah
	}
}

func (c Client) WithBlock(block tongo.TonNodeBlockIdExt) Client {
	c.targetBlockID = &block
	return c
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfo, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfo(ctx)
}

func (c *Client) GetMasterchainInfoExt(ctx context.Context, mode uint32) (liteclient.LiteServerMasterchainInfoExt, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfoExt(ctx, liteclient.LiteServerGetMasterchainInfoExtRequest(mode))
}

func (c *Client) GetTime(ctx context.Context) (uint32, error) {
	res, err := c.getMasterchainServer().LiteServerGetTime(ctx)
	return uint32(res), err
}

func (c *Client) GetVersion(ctx context.Context) (liteclient.LiteServerVersion, error) {
	return c.getMasterchainServer().LiteServerGetVersion(ctx)
}

func (c *Client) GetBlock(ctx context.Context, blockID tongo.TonNodeBlockIdExt) (tongo.Block, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.Block{}, err
	}
	res, err := server.LiteServerGetBlock(ctx, liteclient.LiteServerGetBlockRequest(blockID))
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

func (c *Client) GetState(ctx context.Context, blockID tongo.TonNodeBlockIdExt) (tongo.State, tongo.Hash, tongo.Hash, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	res, err := server.LiteServerGetState(ctx, liteclient.LiteServerGetStateRequest(blockID))
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	if len(cells) != 1 {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, boc.ErrNotSingleRoot
	}
	var state tongo.State // TODO: add State tlb type
	err = tlb.Unmarshal(cells[0], &state)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	return state, tongo.Hash(res.RootHash), tongo.Hash(res.FileHash), nil
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID tongo.TonNodeBlockIdExt, mode uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	res, err := server.LiteServerGetBlockHeader(ctx, liteclient.LiteServerGetBlockHeaderRequest{
		Id:   blockID,
		Mode: mode,
	})
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func (c *Client) LookupBlock(ctx context.Context, blockID tongo.TonNodeBlockId, mode uint32, lt *uint64, utime *uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	res, err := server.LiteServerLookupBlock(ctx, liteclient.LiteServerLookupBlockRequest{
		Mode: mode,
		Id: liteclient.TonNodeBlockId{
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

func decodeBlockHeader(header liteclient.LiteServerBlockHeader) (tongo.BlockInfo, error) {
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

func (c *Client) SendMessage(ctx context.Context, payload []byte) (int32, error) {
	res, err := c.getMasterchainServer().LiteServerSendMessage(ctx, payload)
	return int32(res), err
}

func (c *Client) RunSmcMethod(
	ctx context.Context,
	accountID tongo.AccountID,
	method string,
	params tongo.VmStack,
) (uint32, tongo.VmStack, error) {
	stack, err := tl.Marshal(params)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	id, err := c.targetBlock(ctx)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	req := liteclient.LiteServerRunSmcMethodRequest{
		Mode:     0, // TODO: default mode
		Id:       id,
		Account:  accountID,
		MethodId: uint64(utils.Crc16String(method)&0xffff) | 0x10000,
		Params:   stack,
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

func (c *Client) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tongo.Account, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return tongo.Account{}, err
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return tongo.Account{}, err
	}
	res, err := server.LiteServerGetAccountState(ctx, liteclient.LiteServerGetAccountStateRequest{
		Account: accountID,
		Id:      id,
	})
	if err != nil {
		return tongo.Account{}, err
	}
	if err.(liteclient.LiteServerError).IsNotApplied() { // TODO: add to other methods
		return tongo.Account{}, liteclient.ErrBlockNotApplied
	}
	if len(res.State) == 0 {
		return tongo.Account{SumType: "AccountNone"}, nil
	}
	cells, err := boc.DeserializeBoc(res.State)
	if err != nil {
		return tongo.Account{}, err
	}
	if len(cells) != 1 {
		return tongo.Account{}, boc.ErrNotSingleRoot
	}
	var acc tongo.Account
	err = tlb.Unmarshal(cells[0], &acc)
	if err != nil {
		return tongo.Account{}, err
	}
	return acc, nil
	// TODO: proof check and extract shard account info
}

func (c *Client) GetShardInfo(
	ctx context.Context,
	blockID tongo.TonNodeBlockIdExt,
	workchain uint32,
	shard uint64,
	exact bool,
) (liteclient.LiteServerShardInfo, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return liteclient.LiteServerShardInfo{}, err
	}
	// TODO: decode descr
	return server.LiteServerGetShardInfo(ctx, liteclient.LiteServerGetShardInfoRequest{
		Id:        blockID,
		Workchain: workchain,
		Shard:     shard,
		Exact:     exact,
	})
}

func (c *Client) GetAllShardsInfo(ctx context.Context) (liteclient.LiteServerAllShardsInfo, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return liteclient.LiteServerAllShardsInfo{}, err
	}
	// TODO: decode data
	return c.getMasterchainServer().LiteServerGetAllShardsInfo(ctx, liteclient.LiteServerGetAllShardsInfoRequest(id))
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
		Id:      id,
		Account: accountID,
		Lt:      lt,
	})
	if err != nil {
		return tongo.Transaction{}, err
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
		Account: accountID,
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
	blockID tongo.TonNodeBlockIdExt,
	mode, count uint32,
	after *liteclient.LiteServerTransactionId3,
) ([]liteclient.LiteServerTransactionId, bool, error) {
	// TODO: replace with tongo types
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return nil, false, err
	}
	r, err := server.LiteServerListBlockTransactions(ctx, liteclient.LiteServerListBlockTransactionsRequest{
		Id:    blockID,
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
	knownBlock tongo.TonNodeBlockIdExt,
	targetBlock *tongo.TonNodeBlockIdExt,
) ([]tongo.BlockProof, error) {
	var (
		err    error
		server *liteclient.Client
		mode   uint32 = 0
	)
	if targetBlock != nil {
		server, err = c.getServerByBlockID(targetBlock.TonNodeBlockId)
		mode = 1
	} else {
		server, err = c.getServerByBlockID(knownBlock.TonNodeBlockId)
	}
	if err != nil {
		return nil, err
	}
	r, err := server.LiteServerGetBlockProof(ctx, liteclient.LiteServerGetBlockProofRequest{
		Mode:        mode,
		KnownBlock:  knownBlock,
		TargetBlock: targetBlock,
	})
	if err != nil {
		return nil, err
	}
	// TODO: maybe add from, to, complete
	return decodeBlockProof(r.Steps)
}

func decodeBlockProof(steps []liteclient.LiteServerBlockLink) ([]tongo.BlockProof, error) {
	// TODO: implement
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetConfigAll(ctx context.Context, mode uint32) (*tongo.McStateExtra, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return nil, err
	}
	server, err := c.getServerByBlockID(id.TonNodeBlockId)
	if err != nil {
		return nil, err
	}
	r, err := server.LiteServerGetConfigAll(ctx, liteclient.LiteServerGetConfigAllRequest{
		Mode: mode,
		Id:   id,
	})
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(r.ConfigProof)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, boc.ErrNotSingleRoot
	} // TODO: maybe not
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

//--------------------------------------------------------------------------------------------------------------------//

//func decodeAccountDataFromProof(bocBytes []byte, account tongo.AccountID) (uint64, tongo.Hash, error) {
//	cells, err := boc.DeserializeBoc(bocBytes)
//	if err != nil {
//		return 0, tongo.Hash{}, err
//	}
//	if len(cells) < 1 {
//		return 0, tongo.Hash{}, fmt.Errorf("must be at least one root cell")
//	}
//	var proof struct {
//		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
//	}
//	err = tlb.Unmarshal(cells[1], &proof) // cells order must be strictly defined
//	if err != nil {
//		return 0, tongo.Hash{}, err
//	}
//	values := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Accounts.Values()
//	keys := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Accounts.Keys()
//	for i, k := range keys {
//		keyVal, err := k.ReadBytes(32)
//		if err != nil {
//			return 0, tongo.Hash{}, err
//		}
//		if bytes.Equal(keyVal, account.Address[:]) {
//			return values[i].LastTransLt, values[i].LastTransHash, nil
//		}
//	}
//	return 0, tongo.Hash{}, fmt.Errorf("account not found in ShardAccounts")
//}

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

func (c *Client) BlocksGetShards(ctx context.Context, last tongo.TonNodeBlockIdExt) ([]tongo.TonNodeBlockIdExt, error) {
	asReq, err := makeLiteServerAllShardsInfoRequest(last)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.getMasterchainServer().Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var response struct {
		tl.SumType
		LiteServerAllShardsInfo struct {
			Id    tongo.TonNodeBlockIdExt
			Proof []byte
			Data  []byte
		} `tlSumType:"2de78f09"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &response)
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(response.LiteServerAllShardsInfo.Data)
	if err != nil {
		return nil, err
	}
	var inf tongo.AllShardsInfo
	err = tlb.Unmarshal(cells[0], &inf)
	if err != nil {
		return nil, err
	}
	var shards []tongo.TonNodeBlockIdExt
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

func (c *Client) ValidatorStats(ctx context.Context, mode uint32, last tongo.TonNodeBlockIdExt, limit uint32, startAfter *tongo.Hash, modifiedAfter *uint32) (*tongo.ShardStateUnsplit, error) {
	asReq, err := makeLiteServerGetValidatorStatsRequest(mode, last, limit, startAfter, modifiedAfter)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.getMasterchainServer().Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var response struct {
		tl.SumType
		LiteServerValidatorStats struct {
			Mode       uint32
			Id         tongo.TonNodeBlockIdExt
			Count      uint32
			Complete   uint32
			StateProof []byte
			DataProof  []byte
		} `tlSumType:"d896f7b9"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}

	err = tl.Unmarshal(bytes.NewReader(resp), &response)
	if err != nil {
		return nil, err
	}
	if response.SumType == "Error" {
		return nil, fmt.Errorf(response.Error.Message)
	}
	cells, err := boc.DeserializeBoc(response.LiteServerValidatorStats.DataProof)
	if err != nil {
		return nil, err
	}

	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}

	// }
	return &proof.Proof.VirtualRoot, nil //shards, nil
}
